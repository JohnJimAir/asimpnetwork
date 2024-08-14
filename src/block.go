package src

import (
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

type Block struct {
	Num_node int
	Nodes []Node
}

func (bl *Block) Initialize(num_node int, coefficients_mult [][]float64, coefficient_add []float64, activation []func (float64) (float64), input [][]*rlwe.Ciphertext) {

	bl.Num_node = num_node
	bl.Nodes = make([]Node, num_node)
	for i:=0;i<num_node;i++ {
		bl.Nodes[i] = Node{
			Coefficients_mult: coefficients_mult[i],
			Coefficient_add: coefficient_add[i],
			Activation: activation[i],
			Input: input[i],
		}
	}
}

func (bl Block) Forward(intervals []float64, degrees []int, eval *hefloat.Evaluator, params hefloat.Parameters) (output []*rlwe.Ciphertext) {

	output = make([]*rlwe.Ciphertext, bl.Num_node)
	for i:=0;i<bl.Num_node;i++ {
		output[i] = bl.Nodes[i].Forward(intervals[i], degrees[i], eval, params)
	}
	return output
}