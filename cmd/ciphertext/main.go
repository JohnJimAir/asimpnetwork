// Package main implements an example of smooth function approximation using Chebyshev polynomial interpolation.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"

	"github.com/JohnJimAir/asimpnetwork/src"
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"github.com/tuneinsight/lattigo/v5/he/hefloat/bootstrapping"
	"github.com/tuneinsight/lattigo/v5/ring"
	"github.com/tuneinsight/lattigo/v5/utils"
	"github.com/tuneinsight/lattigo/v5/utils/bignum"
)

var flagShort = flag.Bool("short", false, "run the example with a smaller and insecure ring degree.")

func main() {

	flag.Parse()

	// Default LogN, which with the following defined parameters
	// provides a security of 128-bit.
	LogN := 16

	if *flagShort {
		LogN -= 3
	}

	
	var err error
	input, _ := ReadCSVToFloat64Slice("../../data/test_data_breast-cancer.csv")
	input = Transpose(input)
	

	params, err := hefloat.NewParametersFromLiteral(hefloat.ParametersLiteral{
		LogN:            LogN,                                              // Log2 of the ring degree
		LogQ:            []int{55, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40}, // Log2 of the ciphertext prime moduli
		LogP:            []int{61, 61, 61},                                 // Log2 of the key-switch auxiliary prime moduli
		LogDefaultScale: 40,                                                // Log2 of the scale
		Xs:              ring.Ternary{H: 192},
	})

	if err != nil {
		panic(err)
	}

	btpParametersLit := bootstrapping.ParametersLiteral{
		LogN: utils.Pointy(LogN),
		LogP: []int{61, 61, 61, 61},
		Xs: params.Xs(),
	}

	btpParams, err := bootstrapping.NewParametersFromLiteral(params, btpParametersLit)
	if err != nil {
		panic(err)
	}

	if *flagShort {
		btpParams.Mod1ParametersLiteral.LogMessageRatio += 16 - params.LogN()
	}

	// Key Generator
	kgen := rlwe.NewKeyGenerator(params)

	sk, pk := kgen.GenKeyPairNew()

	encoder := hefloat.NewEncoder(params)
	decryptor := rlwe.NewDecryptor(params, sk)
	encryptor := rlwe.NewEncryptor(params, pk)

	fmt.Println()
	fmt.Println("Generating bootstrapping evaluation keys...")
	evk_boot, _, err := btpParams.GenEvaluationKeys(sk)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done")

	var eval_boot *bootstrapping.Evaluator
	if eval_boot, err = bootstrapping.NewEvaluator(btpParams, evk_boot); err != nil {
		panic(err)
	}

	// Relinearization Key
	rlk := kgen.GenRelinearizationKeyNew(sk)

	// Evaluation Key Set with the Relinearization Key
	evk := rlwe.NewMemEvaluationKeySet(rlk)

	// Evaluator
	eval := hefloat.NewEvaluator(params, evk)

	// Samples values in [-K, K]
	K := 16.0

	// Allocates a plaintext at the max level.
	pt := hefloat.NewPlaintext(params, params.MaxLevel())

	values := make([]float64, pt.Slots())

	for i := range values {
		values[i] = 4.3
	}
	if err = encoder.Encode(values, pt); err != nil {
		panic(err)
	}
	var ct_in *rlwe.Ciphertext
	if ct_in, err = encryptor.EncryptNew(pt); err != nil {
		panic(err)
	}



	tanh := func(x float64) (y float64) {
		return math.Tanh(x)
	}
	sin := func(x float64) (y float64) {
		return math.Sin(x)
	}
	abs := func(x float64) (y float64) {
		return math.Abs(x)
	}
	exp := func(x float64) (y float64) {
		return math.Exp(x)
	}
	pow2 := func(x float64) (y float64) {
		return math.Pow(x, 2)
	}
	pow3 := func(x float64) (y float64) {
		return math.Pow(x, 3)
	}
	identity := func(x float64) (y float64) {
		return x
	}
	contract := func(x float64) (y float64) {
		return 0.000001* x
	}
	
	for i := range values {
		values[i] = -2.236
	}
	if err = encoder.Encode(values, pt); err != nil {
		panic(err)
	}
	var ct_0_in *rlwe.Ciphertext
	if ct_0_in, err = encryptor.EncryptNew(pt); err != nil {
		panic(err)
	}

	// for i := range values {
	// 	values[i] = 1.1
	// }
	// if err = encoder.Encode(values, pt); err != nil {
	// 	panic(err)
	// }
	// var ct_1_in *rlwe.Ciphertext
	// if ct_1_in, err = encryptor.EncryptNew(pt); err != nil {
	// 	panic(err)
	// }

	nn := src.Node{
		Coefficients_mult: []float64{1.0001},
		Coefficient_add: 1.0001,
		Activation : identity,
		Input: []*rlwe.Ciphertext{ct_0_in},
	}

	out := nn.Forward(16, 31, eval, params)
	fmt.Println("nnnnnnnnnn")
	PrintValues(params, out, encoder, decryptor)



	// Chebyhsev approximation of the sigmoid in the domain [-K, K] of degree 63.
	poly := hefloat.NewPolynomial(GetChebyshevPoly(16, 15, tanh))

	// Instantiates the polynomial evaluator
	polyEval := hefloat.NewPolynomialEvaluator(params, eval)

	// Retrieves the change of basis y = scalar * x + constant
	scalar, constant := poly.ChangeOfBasis()

	// Performes the change of basis Standard -> Chebyshev
	var ct_out *rlwe.Ciphertext
	if ct_out, err = encryptor.EncryptNew(pt); err != nil {
		panic(err)
	}
	if err := eval.Mul(ct_in, scalar, ct_out); err != nil {
		panic(err)
	}

	if err := eval.Add(ct_out, constant, ct_out); err != nil {
		panic(err)
	}

	if err := eval.Rescale(ct_out, ct_out); err != nil {
		panic(err)
	}

	// Evaluates the polynomial
	if ct_out, err = polyEval.Evaluate(ct_out, poly, params.DefaultScale()); err != nil {
		panic(err)
	}

	// Allocates a vector for the reference values and
	// evaluates the same circuit on the plaintext values
	want := make([]float64, ct_out.Slots())
	for i := range want {
		want[i], _ = poly.Evaluate(values[i])[0].Float64()
		want[i] = tanh(values[i])
	}

	// Decrypts and print the stats about the precision.
	PrintPrecisionStats(params, ct_out, want, encoder, decryptor)


	input_ct := EncryptMany(params, encoder, encryptor, [][]float64{input[0], input[1], input[2], input[3], input[4], input[5], input[6], input[7], input[8]})
	input_ct_2d := [][]*rlwe.Ciphertext{{input_ct[0]}, {input_ct[1]}, {input_ct[2]}, {input_ct[3]}, {input_ct[4]}, {input_ct[5]}, {input_ct[6]}, {input_ct[7]}, {input_ct[8]}}

	var blo_top_0 *src.Block = new(src.Block)
	blo_top_0.Initialize(9, [][]float64{{3.77}, {7.07}, {9.52}, {9.96}, {3.64}, {2.24}, {10.000001}, {7.85}, {7.94}}, 
		[]float64{-1.01, -6.21, -8.15, -3.26, -0.62, 8.2, -8.2, 7.58, -0.2}, 
		[]func (float64) (float64){tanh, sin, sin, abs, sin, sin, tanh, sin, abs},
		input_ct_2d,
	)
	out_blo_top_0 := blo_top_0.Forward([]float64{K,K,K, 8,K,K, K,K,8}, []int{31,31,31, 31,31,31, 31,31,31}, eval, params)
	out_blo_top_0_BTS := BTSmany(eval_boot, out_blo_top_0)
	// fmt.Println("0000000")
	// PrintValuesMany(params, out_blo_top_0_BTS, encoder, decryptor)


	var blo_top_1 *src.Block = new(src.Block)
	blo_top_1.Initialize(9, [][]float64{{0.000001}, {7.4}, {-1.000001}, {9.6}, {6.44}, {6.11}, {5.2}, {4.95}, {5.89}}, 
		[]float64{0.000001, 1.19, 0.33, -2.47, -2.23, -0.73, 1.18, 9.62, -2.45}, 
		[]func (float64) (float64){contract, sin, pow2, tanh, sin, sin, sin, sin, tanh},
		input_ct_2d,
	)
	out_blo_top_1 := blo_top_1.Forward([]float64{K,K,K, K,K,K, K,K,K}, []int{31,31,31, 31,31,31, 31,31,31}, eval, params)
	out_blo_top_1_BTS := BTSmany(eval_boot, out_blo_top_1)
	// fmt.Println("1111111")
	// PrintValuesMany(params, out_blo_top_1_BTS, encoder, decryptor)


	var blo_top_2 *src.Block = new(src.Block)
	blo_top_2.Initialize(9, [][]float64{{-1.000001}, {5.08}, {6.62}, {7.21}, {2.2}, {3.24}, {-1.000001}, {-1.000001}, {0.28}}, 
		[]float64{0.43, -2.22, 2.99, -5.79, -9.64, -2.6, 0.24, 0.37, 1.0}, 
		[]func (float64) (float64){pow3, sin, sin, sin, contract, tanh, pow2, pow3, contract},
		input_ct_2d,
	)
	out_blo_top_2 := blo_top_2.Forward([]float64{K,K,K, K,K,K, K,K,K}, []int{31,31,31, 31,31,31, 31,31,31}, eval, params)
	out_blo_top_2_BTS := BTSmany(eval_boot, out_blo_top_2)
	// fmt.Println("2222222")
	// PrintValuesMany(params, out_blo_top_2_BTS, encoder, decryptor)


	var blo_top_3 *src.Block = new(src.Block)
	blo_top_3.Initialize(9, [][]float64{{1.13}, {3.94}, {3.89}, {3.86}, {1.49}, {10.000001}, {3.65}, {9.79}, {7.8}}, 
		[]float64{-9.75, -0.58, -7.86, -8.02, 2.53, -2.6, -1.43, 4.21, -0.84}, 
		[]func (float64) (float64){contract, tanh, sin, sin, contract, tanh, sin, sin, tanh},
		input_ct_2d,
	)
	out_blo_top_3 := blo_top_3.Forward([]float64{K,K,K, K,K,K, K,K,K}, []int{31,31,31, 31,31,31, 31,31,31}, eval, params)
	out_blo_top_3_BTS := BTSmany(eval_boot, out_blo_top_3)
	// fmt.Println("333333")
	// PrintValuesMany(params, out_blo_top_3_BTS, encoder, decryptor)


	var blo_middle *src.Block = new(src.Block)
	coefficient_middle := [][]float64{
		{0.08, 0.39, 0.09, 0.000001/*-0.e-2*/, 0.21, -0.13, 0.19, 0.01, 0.01},
		{0.000001, 1.38, 2.28, 0.27, 1.64, 0.72, -0.37, -0.87, 0.29},
		{-0.61, -0.01, -0.04, 0.05, 0.000001/*tan -0.e-2*/, 0.07, 0.05, -0.3, 0.000001/*tan*/},
		{0.000001/*tan*/, 25.59, 21.1, 19.17, 0.000001/*tan*/, 12.94, 4.29, 12.29, 9.55},
	}
	blo_middle.Initialize(4, coefficient_middle, 
		[]float64{-1.0, 5.55, 4.04, 85.59}, 
		[]func (float64) (float64){pow2, sin, sin, identity},
		[][]*rlwe.Ciphertext{out_blo_top_0_BTS, out_blo_top_1_BTS, out_blo_top_2_BTS, out_blo_top_3_BTS},
	)
	out_blo_middle := blo_middle.Forward([]float64{K,K,K, K}, []int{31,31,31, 31}, eval, params)
	out_blo_middle_BTS := BTSmany(eval_boot, out_blo_middle)
	// fmt.Println("middle")
	// PrintValuesMany(params, out_blo_middle_BTS, encoder, decryptor)


	var blo_bottom *src.Block = new(src.Block)
	coefficient_bottom := [][]float64 {
		{0.13, -0.09, 2.93, -0.01},
		{0.31, -0.21, 7.04, -0.03},
	}
	blo_bottom.Initialize(2, coefficient_bottom, 
		[]float64{0.0, 10.72}, 
		[]func (float64) (float64){exp, tanh},
		[][]*rlwe.Ciphertext{out_blo_middle_BTS, out_blo_middle_BTS},
	)
	out_blo_bottom := blo_bottom.Forward([]float64{8,8}, []int{31,31}, eval, params)
	// fmt.Println("bottom")
	// PrintValuesMany(params, out_blo_bottom, encoder, decryptor)

	var blo_trick *src.Block = new(src.Block)
	blo_trick.Initialize(2, [][]float64{{1988.48}, {-7.34}}, 
		[]float64{-31.97, 1.99}, 
		[]func (float64) (float64){identity, identity},
		[][]*rlwe.Ciphertext{ {out_blo_bottom[0]}, {out_blo_bottom[1]} },
	)
	out_blo_trick := blo_trick.Forward([]float64{8,8}, []int{1,1}, eval, params)
	// fmt.Println("trick")
	// PrintValuesMany(params, out_blo_trick, encoder, decryptor)

	re := PrintValuesMany(params, out_blo_trick, encoder, decryptor)
	re = Transpose(re)
	for i:=0;i<140;i++ {
		fmt.Printf("%.8f,%.8f", re[i][0], re[i][1])
		fmt.Println()
	}

}

// GetChebyshevPoly returns the Chebyshev polynomial approximation of f the
// in the interval [-K, K] for the given degree.
func GetChebyshevPoly(K float64, degree int, f64 func(x float64) (y float64)) bignum.Polynomial {

	FBig := func(x *big.Float) (y *big.Float) {
		xF64, _ := x.Float64()
		return new(big.Float).SetPrec(x.Prec()).SetFloat64(f64(xF64))
	}

	var prec uint = 128

	interval := bignum.Interval{
		A:     *bignum.NewFloat(-K, prec),
		B:     *bignum.NewFloat(K, prec),
		Nodes: degree,
	}

	// Returns the polynomial.
	return bignum.ChebyshevApproximation(FBig, interval)
}

// PrintPrecisionStats decrypts, decodes and prints the precision stats of a ciphertext.
func PrintPrecisionStats(params hefloat.Parameters, ct *rlwe.Ciphertext, want []float64, ecd *hefloat.Encoder, dec *rlwe.Decryptor) {

	var err error

	// Decrypts the vector of plaintext values
	pt := dec.DecryptNew(ct)

	// Decodes the plaintext
	have := make([]float64, ct.Slots())
	if err = ecd.Decode(pt, have); err != nil {
		panic(err)
	}

	// Pretty prints some values
	fmt.Printf("Have: ")
	for i := 0; i < 4; i++ {
		fmt.Printf("%20.15f ", have[i])
	}
	fmt.Printf("...\n")

	fmt.Printf("Want: ")
	for i := 0; i < 4; i++ {
		fmt.Printf("%20.15f ", want[i])
	}
	fmt.Printf("...\n")

	// Pretty prints the precision stats
	// fmt.Println(hefloat.GetPrecisionStats(params, ecd, dec, have, want, 0, false).String())
}

func PrintValues(params hefloat.Parameters, ct *rlwe.Ciphertext, ecd *hefloat.Encoder, dec *rlwe.Decryptor) {
	
	var err error
	// Decrypts the vector of plaintext values
	pt := dec.DecryptNew(ct)

	// Decodes the plaintext
	have := make([]float64, ct.Slots())
	if err = ecd.Decode(pt, have); err != nil {
		panic(err)
	}

	for i := 0; i < 4; i++ {
		fmt.Printf("%.8f ", have[i])
	}
	fmt.Printf("...\n")
}

func PrintValuesMany(params hefloat.Parameters, ct []*rlwe.Ciphertext, ecd *hefloat.Encoder, dec *rlwe.Decryptor) (output [][]float64) {
	
	var err error
	num := len(ct)
	output = make([][]float64, num)

	for i:=0;i<num;i++ {
		pt := dec.DecryptNew(ct[i])
		values := make([]float64, ct[i].Slots())
		if err = ecd.Decode(pt, values); err != nil {
			panic(err)
		}
		output[i] = values
	
		for i := 0; i < 4; i++ {
			fmt.Printf("%.8f ", values[i])
		}
		fmt.Printf("...\n")
	}	

	return output
}



func BTSmany(eval_boot *bootstrapping.Evaluator, input []*rlwe.Ciphertext) (output []*rlwe.Ciphertext) {
	
	var err error
	output = make([]*rlwe.Ciphertext, len(input))
	for i:=0;i<len(input);i++ {
		output[i], err = eval_boot.Bootstrap(input[i])
		if err != nil {
			panic(err)
		}
	}
	return  output
}

func EncryptMany(params hefloat.Parameters, encoder *hefloat.Encoder, encryptor *rlwe.Encryptor, value_input [][]float64) (output []*rlwe.Ciphertext) {
 
	var err error
	num := len(value_input)
	output = make([]*rlwe.Ciphertext, num)

	pt := hefloat.NewPlaintext(params, params.MaxLevel())
	for i:=0;i<num;i++ {
		if err = encoder.Encode(value_input[i], pt); err != nil {
			panic(err)
		}
		if output[i], err = encryptor.EncryptNew(pt); err != nil {
			panic(err)
		}
	}
	return output
}

func Transpose(data [][]float64) [][]float64 {
	if len(data) == 0 {
		return nil
	}

	rows := len(data)
	cols := len(data[0])

	transposed := make([][]float64, cols)
	for i := range transposed {
		transposed[i] = make([]float64, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			transposed[j][i] = data[i][j]
		}
	}

	return transposed
}

func ReadCSVToFloat64Slice(filename string) ([][]float64, error) {
    // 打开 CSV 文件
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // 创建一个 CSV 读取器
    reader := csv.NewReader(file)

    // 读取所有行
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    // 创建一个二维切片来存储 float64 数据
    var data [][]float64

    // 遍历 CSV 内容，跳过第一行，并忽略每一行的最后一列
    for i, record := range records {
        if i == 0 {
            // 跳过第一行（通常是标题行）
            continue
        }

        // 创建一个切片来存储当前行的 float64 数据
        var row []float64
        for j := 0; j < len(record)-1; j++ { // 忽略最后一列
            // 将字符串转换为 float64
            value, err := strconv.ParseFloat(record[j], 64)
            if err != nil {
                return nil, err
            }
            row = append(row, value)
        }
        data = append(data, row)
    }

    return data, nil
}