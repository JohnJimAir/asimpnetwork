package src

import (
	"math/big"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"github.com/tuneinsight/lattigo/v5/utils/bignum"
)

type Node struct {
	Coefficients_mult []float64  // must not be integer
	Coefficient_add float64
	Activation func (float64) (float64)	
	Input []*rlwe.Ciphertext
}

func (n Node) Forward(interval []float64, degree int, eval *hefloat.Evaluator, params hefloat.Parameters) (output *rlwe.Ciphertext) {
	
	var err error
	output = Innerproduct(n.Coefficients_mult, n.Coefficient_add, n.Input, eval)
	
	poly := hefloat.NewPolynomial(GetChebyshevPoly(interval[0], interval[1], degree, n.Activation))
	polyEval := hefloat.NewPolynomialEvaluator(params, eval)

	scalar, constant := poly.ChangeOfBasis()

	if err := eval.Mul(output, scalar, output); err != nil {
		panic(err)
	}
	if err := eval.Add(output, constant, output); err != nil {
		panic(err)
	}
	if err := eval.Rescale(output, output); err != nil {
		panic(err)
	}
	if output, err = polyEval.Evaluate(output, poly, params.DefaultScale()); err != nil {
		panic(err)
	}
	return output
}

func Innerproduct(coefficients_mult []float64, coefficient_add float64, input []*rlwe.Ciphertext, eval *hefloat.Evaluator) (output *rlwe.Ciphertext) {

	var err error
	num := len(input)
	tmp := make([]*rlwe.Ciphertext, len(input))
	for i:=0;i<num;i++ {
		if tmp[i], err = eval.MulNew(input[i], coefficients_mult[i]); err != nil {
			panic(err)
		}
		// if err := eval.Rescale(tmp[i], tmp[i]); err != nil {
		// 	panic(err)
		// }
	}

	output = tmp[0]
	for i:=1;i<num;i++ {
		if err := eval.Add(output, tmp[i], output); err != nil {
			panic(err)
		}
	}

	if err := eval.Add(output, coefficient_add, output); err != nil {
		panic(err)
	}

	if err := eval.Rescale(output, output); err != nil {
		panic(err)
	}

	return output
}

func GetChebyshevPoly(K_left, K_right float64, degree int, f64 func(x float64) (y float64)) bignum.Polynomial {

	FBig := func(x *big.Float) (y *big.Float) {
		xF64, _ := x.Float64()
		return new(big.Float).SetPrec(x.Prec()).SetFloat64(f64(xF64))
	}

	var prec uint = 128

	interval := bignum.Interval{
		A:     *bignum.NewFloat(K_left, prec),
		B:     *bignum.NewFloat(K_right, prec),
		Nodes: degree,
	}

	// Returns the polynomial.
	return bignum.ChebyshevApproximation(FBig, interval)
}