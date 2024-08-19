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
	input, _ := ReadCSVToFloat64Slice("../../data/test_data_sepsis.csv")
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
	// abs := func(x float64) (y float64) {
	// 	return math.Abs(x)
	// }
	// exp := func(x float64) (y float64) {
	// 	return math.Exp(x)
	// }
	pow2 := func(x float64) (y float64) {
		return math.Pow(x, 2)
	}
	// pow3 := func(x float64) (y float64) {
	// 	return math.Pow(x, 3)
	// }
	identity := func(x float64) (y float64) {
		return x
	}
	contract := func(x float64) (y float64) {
		return 0.000001* x
	}
	log := func(x float64) (y float64) {
		return math.Log(x)
	}
	sqrt := func(x float64) (y float64) {
		return math.Sqrt(x)
	}
	
	for i := range values {
		values[i] = -2.23
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

	out := nn.Forward([]float64{-16.0, 16.0}, 31, eval, params)
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


	input_ct := EncryptMany(params, encoder, encryptor, input)
	input_ct_2d := make([][]*rlwe.Ciphertext, 0)
	for i:=0;i<len(input_ct);i++ {
		input_ct_2d = append(input_ct_2d, []*rlwe.Ciphertext{input_ct[i]})
	}
	fmt.Println("lenlllllllllll")
	fmt.Println(len(input_ct_2d))

	var blo_top *src.Block = new(src.Block)
	coefficient_top_mult := [][]float64{{-1.0000001},{0.58},{0.4},{0.14},{0.2},
		{0.0000001},{0.0000001},{1.06},{0.0000001},{0.28},
		{1.0000001},{0.0000001},{0.0000001},{3.4},{0.0000001},
		{-1.38},{0.27},{9.96},{0.31},{0.89},
		{0.95},{0.43},{0.41},{0.31},{0.16},
		{0.18},{0.17},{1.02},{0.21},{0.23},
		{0.35},{1.22},{1.51},{6.07},{0.26},
		{0.28},{0.26}}
	coefficient_top_add := []float64{-0.72, -0.48, 1.37, 1.0, -0.85,
		0.0, 0.0, -9.61, 0.0, -5.95, 
		0.37, 0.0, 0.0, 3.95, 0.0,
		4.25, 1.85, 7.21, 5.04, -0.18, 
		-0.53, 2.24, 2.39, 1.61, -4.2,
		-7.56, 8.5, -0.68, 2.16, -7.04,
		-1.43, -2.35, -2.12, 2.42, 2.15,
		-7.78, 4.52}

	blo_top.Initialize(37, coefficient_top_mult, coefficient_top_add, 
		[]func (float64) (float64){pow2, tanh, sin, contract, tanh,
			contract, contract, sin, contract, contract,
			sqrt, contract, contract, log, contract,
			log, sin, identity, sin, sin,
			tanh, sin, sin, sin, sin,
			sin, sin, tanh, sin,sin,
			tanh, tanh, tanh, identity, sin,
			sin, sin},
		input_ct_2d )
	out_blo_top := blo_top.Forward([][]float64{{-K,K},{-K,K},{-K,K},{-K,K},{-K,K}, {-K,K},{-K,K},{-K,K},{-K,K},{-K,K}, {0.0,K},{-K,K},{-K,K},{0.0,K},{-K,K},
		{0.0,K},{-K,K},{-K,K},{-K,K},{-K,K}, {-K,K},{-K,K},{-K,K},{-K,K},{-K,K}, {-K,K},{-K,K},{-K,K},{-K,K},{-K,K}, {-K,K},{-K,K},{-K,K},{-K,K},{-K,K}, {-K,K},{-K,K}}, 
		[]int{31,31,31,31,31, 31,31,31,31,31, 31,31,31,31,31, 31,31,31,31,31, 31,31,31,31,31, 31,31,31,31,31, 31,31,31,31,31, 31,31}, 
		eval, params)
	out_blo_top_BTS := BTSmany(eval_boot, out_blo_top)
	// fmt.Println("topppppppppp")
	// PrintValuesMany(params, out_blo_top_BTS, encoder, decryptor)


	var blo_bottom *src.Block = new(src.Block)
	coefficient_bottom_mult := []float64{0.04, -0.07, -0.05, 0.01, -0.32, 
		0.0, 0.0, 0.05, 0.0, 0.1, 
		-0.28, 0.0, 0.0, -0.24, 0.0,
		0.05, 0.12, 0.02, -0.12, 0.02,
		0.03, -0.49, 0.13, 0.24, 0.35,
		0.13, 0.5, 0.02, 0.13, -0.21,  
		-0.11, 0.04, 0.03, -0.01, 0.18,
		-0.06, 0.04}
	blo_bottom.Initialize(1, [][]float64{coefficient_bottom_mult}, 
		[]float64{5.76}, 
		[]func (float64) (float64){sin},
		[][]*rlwe.Ciphertext{out_blo_top_BTS})
	out_blo_bottom := blo_bottom.Forward([][]float64{{K,-K}}, []int{31}, eval, params)
	// fmt.Println("bottom")
	// PrintValuesMany(params, out_blo_bottom, encoder, decryptor)

	var blo_trick *src.Block = new(src.Block)
	blo_trick.Initialize(1, [][]float64{{-1.05}}, 
		[]float64{1.04}, 
		[]func (float64) (float64){identity},
		[][]*rlwe.Ciphertext{ {out_blo_bottom[0]} },
	)
	out_blo_trick := blo_trick.Forward([][]float64{{-K,K}}, []int{1,1}, eval, params)
	// fmt.Println("trick")
	// PrintValuesMany(params, out_blo_trick, encoder, decryptor)

	re := PrintValuesMany(params, out_blo_trick, encoder, decryptor)
	re = Transpose(re)
	for i:=0;i<138;i++ {
		fmt.Printf("%.8f", re[i][0])
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