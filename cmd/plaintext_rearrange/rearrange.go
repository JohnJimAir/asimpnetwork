package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

func main() {
	input, _ := ReadCSVToFloat64Slice("../../data/test_data_breast-cancer.csv")
	output := make([][]float64, 0)
	for i:=0;i<len(input);i++ {
		// fmt.Println(KAN(input[i]))
		_, _, result := KAN(input[i])
		output = append(output, result)
	}
	output = Transpose(output)
	for i:=0;i<len(output);i++ {
		output[i] = SortFloat64(output[i])
	}
	output = Transpose(output)
	PrintToFile(output, "./bound/bound_output_final.txt")

}

func KAN(input []float64) (output_0, output_1 float64, output_tmp []float64) {

	middle_0 := make([]float64, 9)
	middle_0[0] = math.Tanh(3.77* input[0] - 1.01)
	middle_0[1] = math.Sin(7.07* input[1] - 6.21)
	middle_0[2] = math.Sin(9.52* input[2] - 8.15)
	middle_0[3] = math.Abs(9.96* input[3] - 3.26) // -3.26 to 6.7
	middle_0[4] = math.Sin(3.64* input[4] - 0.62)
	middle_0[5] = math.Sin(2.24* input[5] + 8.2)
	middle_0[6] = math.Tanh(10.0* input[6] - 8.2)
	middle_0[7] = math.Sin(7.85* input[7] + 7.58)
	middle_0[8] = math.Abs(7.94* input[8] - 0.2) // -0.2 to 7.74

	middle_1 := make([]float64, 9)
	middle_1[0] = 0
	middle_1[1] = math.Sin(7.4* input[1] + 1.19)
	middle_1[2] = math.Pow(0.33 - input[2], 2)
	middle_1[3] = math.Tanh(9.6* input[3] - 2.47)
	middle_1[4] = math.Sin(6.44* input[4] - 2.23)
	middle_1[5] = math.Sin(6.11* input[5] - 0.73)
	middle_1[6] = math.Sin(5.2* input[6] + 1.18)
	middle_1[7] = math.Sin(4.95* input[7] + 9.62)
	middle_1[8] = math.Tanh(5.89* input[8] - 2.45)

	middle_2 := make([]float64, 9)
	middle_2[0] = math.Pow(0.43 - input[0], 3)
	middle_2[1] = math.Sin(5.08* input[1] - 2.22)
	middle_2[2] = math.Sin(6.62* input[2] + 2.99)
	middle_2[3] = math.Sin(7.21* input[3] - 5.79)
	middle_2[4] = math.Tan(2.2* input[4] - 9.64) //
	middle_2[5] = math.Tanh(3.24* input[5] - 2.6)
	middle_2[6] = math.Pow(0.24 - input[6], 2)
	middle_2[7] = math.Pow(0.37 - input[7], 3)
	middle_2[8] = math.Tan(0.28* input[8] + 1.0) //

	middle_3 := make([]float64, 9)
	middle_3[0] = math.Tan(1.13* input[0] - 9.75) //
	middle_3[1] = math.Tanh(3.94* input[1] - 0.58)
	middle_3[2] = math.Sin(3.89* input[2] - 7.86)
	middle_3[3] = math.Sin(3.86* input[3] - 8.02)
	middle_3[4] = math.Tan(1.49* input[4] + 2.53) // dump "tan" by setting a very small coefficient before it. But to maintain the structural consistency of Node, use "contract" to replace "Tan", such as to best eliminate its effect.
	middle_3[5] = math.Tanh(10.0* input[5] - 2.6)
	middle_3[6] = math.Sin(3.65* input[6] - 1.43)
	middle_3[7] = math.Sin(9.79* input[7] + 4.21)
	middle_3[8] = math.Tanh(7.8* input[8] - 0.84)

	output := make([]float64, 4)
	coefficients_middle_0 := []float64{0.08, 0.39, 0.09, -0.e-2, 0.21, -0.13, 0.19, 0.01, 0.01}
	coefficients_middle_1 := []float64{0.0, 1.38, 2.28, 0.27, 1.64, 0.72, -0.37, -0.87, 0.29}
	coefficients_middle_2 := []float64{-0.61, -0.01, -0.04, 0.05, -0.e-2, 0.07, 0.05, -0.3, 0.08}
	coefficients_middle_3 := []float64{2.98, 25.59, 21.1, 19.17, 2.38, 12.94, 4.29, 12.29, 9.55}
	output[0] = math.Pow( MultiplyAndSum(middle_0, coefficients_middle_0) - 1, 2)
	output[1] = math.Sin( MultiplyAndSum(middle_1, coefficients_middle_1) + 5.55)
	output[2] = math.Sin( MultiplyAndSum(middle_2, coefficients_middle_2) + 4.04)
	output[3] = math.Abs( MultiplyAndSum(middle_3, coefficients_middle_3) + 85.59) // -4.3 to 182.42

	output_final := make([]float64, 2)
	coefficient_0 := []float64{0.13, -0.09, 2.93, -0.01}
	coefficient_1 := []float64{0.31, -0.21, 7.04, -0.03}
	output_final[0] = 1988.48* math.Exp( MultiplyAndSum(output, coefficient_0) + 0.0) - 31.97
	output_final[1] = -7.34* math.Tanh( MultiplyAndSum(output, coefficient_1) + 10.72) + 1.99

	return 	output_final[0], output_final[1], output_tmp
}

func SumFloats(numbers []float64) float64 {
    var sum float64
    for _, number := range numbers {
        sum += number
    }
    return sum
}

func MultiplyAndSum(slice1, slice2 []float64) float64 {
    if len(slice1) != len(slice2) {
        fmt.Println("different length")
        return 0
    }

    sum := 0.0
    for i := range slice1 {
        sum += slice1[i] * slice2[i]
    }
    return sum
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

func SortFloat64(slice []float64) []float64 {

	sortedSlice := make([]float64, len(slice))
	copy(sortedSlice, slice)

	sort.Slice(sortedSlice, func(i, j int) bool {
		return sortedSlice[i] < sortedSlice[j]
	})

	return sortedSlice
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

func PrintToFile(data [][]float64, filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, row := range data {
		for _, value := range row {
			_, err := fmt.Fprintf(file, "%.8f\t", value)
			if err != nil {
				return err
			}
		}
		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}