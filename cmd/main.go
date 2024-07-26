package main

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	dModel = 64
	nHeads = 8
	dFF    = 256
)

type Transformer struct {
	EncoderLayers []EncoderLayer
	DecoderLayers []DecoderLayer
}

type EncoderLayer struct {
	Attention MultiHeadAttention
	FFN       FeedForward
}

type DecoderLayer struct {
	Attention1 MultiHeadAttention
	Attention2 MultiHeadAttention
	FFN        FeedForward
}

type MultiHeadAttention struct {
	WQ [][][]float64
	WK [][][]float64
	WV [][][]float64
	WO [][]float64
}

type FeedForward struct {
	W1 [][]float64
	W2 [][]float64
}

func NewMultiHeadAttention() MultiHeadAttention {
	return MultiHeadAttention{
		WQ: random3DMatrix(nHeads, dModel, dModel/nHeads),
		WK: random3DMatrix(nHeads, dModel, dModel/nHeads),
		WV: random3DMatrix(nHeads, dModel, dModel/nHeads),
		WO: randomMatrix(dModel, dModel),
	}
}

func NewFeedForward() FeedForward {
	return FeedForward{
		W1: randomMatrix(dModel, dFF),
		W2: randomMatrix(dFF, dModel),
	}
}

func randomMatrix(rows, cols int) [][]float64 {
	matrix := make([][]float64, rows)
	for i := range matrix {
		matrix[i] = make([]float64, cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Float64()*2 - 1
		}
	}
	return matrix
}

func random3DMatrix(depth, rows, cols int) [][][]float64 {
	matrix := make([][][]float64, depth)
	for i := range matrix {
		matrix[i] = randomMatrix(rows, cols)
	}
	return matrix
}

func (mha MultiHeadAttention) forward(Q, K, V [][]float64) [][]float64 {
	heads := make([][][]float64, nHeads)
	for i := 0; i < nHeads; i++ {
		heads[i] = scaledDotProductAttention(matMul(Q, mha.WQ[i]), matMul(K, mha.WK[i]), matMul(V, mha.WV[i]))
	}
	concatenated := concatenate(heads)
	return matMul(concatenated, mha.WO)
}

func (ff FeedForward) forward(x [][]float64) [][]float64 {
	return relu(matMul(relu(matMul(x, ff.W1)), ff.W2))
}

func (el EncoderLayer) forward(x [][]float64) [][]float64 {
	x = el.Attention.forward(x, x, x)
	return el.FFN.forward(x)
}

func (dl DecoderLayer) forward(x, encOut [][]float64) [][]float64 {
	x = dl.Attention1.forward(x, x, x)
	x = dl.Attention2.forward(x, encOut, encOut)
	return dl.FFN.forward(x)
}

func scaledDotProductAttention(Q, K, V [][]float64) [][]float64 {
	dk := float64(len(K[0]))
	scores := matMul(Q, transpose(K))
	for i := range scores {
		for j := range scores[i] {
			scores[i][j] /= math.Sqrt(dk)
		}
	}
	return matMul(softmax(scores), V)
}

func transpose(matrix [][]float64) [][]float64 {
	t := make([][]float64, len(matrix[0]))
	for i := range t {
		t[i] = make([]float64, len(matrix))
		for j := range t[i] {
			t[i][j] = matrix[j][i]
		}
	}
	return t
}

func matMul(a, b [][]float64) [][]float64 {
	result := make([][]float64, len(a))
	for i := range result {
		result[i] = make([]float64, len(b[0]))
		for j := range result[i] {
			for k := range a[0] {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

func relu(matrix [][]float64) [][]float64 {
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] < 0 {
				matrix[i][j] = 0
			}
		}
	}
	return matrix
}

func softmax(matrix [][]float64) [][]float64 {
	for i := range matrix {
		max := matrix[i][0]
		for j := range matrix[i] {
			if matrix[i][j] > max {
				max = matrix[i][j]
			}
		}
		sum := 0.0
		for j := range matrix[i] {
			matrix[i][j] = math.Exp(matrix[i][j] - max) // 减去最大值以增强数值稳定性
			sum += matrix[i][j]
		}
		for j := range matrix[i] {
			matrix[i][j] /= sum
		}
	}
	return matrix
}

func concatenate(matrices [][][]float64) [][]float64 {
	rows := len(matrices[0])
	cols := 0
	for _, matrix := range matrices {
		cols += len(matrix[0])
	}
	result := make([][]float64, rows)
	for i := range result {
		result[i] = make([]float64, cols)
		offset := 0
		for _, matrix := range matrices {
			copy(result[i][offset:], matrix[i])
			offset += len(matrix[i])
		}
	}
	return result
}

func main() {
	encoderLayer := EncoderLayer{NewMultiHeadAttention(), NewFeedForward()}
	decoderLayer := DecoderLayer{NewMultiHeadAttention(), NewMultiHeadAttention(), NewFeedForward()}
	transformer := Transformer{
		EncoderLayers: []EncoderLayer{encoderLayer},
		DecoderLayers: []DecoderLayer{decoderLayer},
	}

	input := randomMatrix(1, dModel)
	fmt.Println("Input:", input)
	encOutput := transformer.EncoderLayers[0].forward(input)
	fmt.Println("Encoder Output:", encOutput)
	decOutput := transformer.DecoderLayers[0].forward(input, encOutput)
	fmt.Println("Decoder Output:", decOutput)
}
