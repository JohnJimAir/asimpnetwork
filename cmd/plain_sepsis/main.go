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
	input, _ := ReadCSVToFloat64Slice("../../data/test_data_sepsis.csv")
	output := make([][]float64, 0)
	for i:=0;i<len(input);i++ {
		// fmt.Println(KAN(input[i]))
		_, result := KAN(input[i])
		output = append(output, result)
	}
	output = Transpose(output)
	for i:=0;i<len(output);i++ {
		output[i] = SortFloat64(output[i])
	}
	output = Transpose(output)
	PrintToFile(output, "./bound/bound_output_final.txt")

}

func KAN(input []float64) (output float64, output_tmp []float64) {

	middle := make([]float64, 37) 
	middle[0] = math.Pow(-1.0* input[0] - 0.72, 2)
	middle[1] = math.Tanh(0.58* input[1] - 0.48)
	middle[2] = math.Sin(0.4* input[2] + 1.37)
	middle[3] = math.Tan(0.14* input[3] + 1.0) //
	middle[4] = math.Tanh(0.2* input[4] - 0.85)
	middle[5] = 0.0
	middle[6] = 0.0
	middle[7] = math.Sin(1.06* input[7] - 9.61)
	middle[8] = 0.0
	middle[9] = math.Tan(0.28* input[9] - 5.95) //
	middle[10] = math.Sqrt(1.0* input[10] + 0.37)
	middle[11] = 0.0
	middle[12] = 0.0
	middle[13] = math.Log(3.4* input[13] + 3.95)
	middle[14] = 0.0
	middle[15] = math.Log(-1.38* input[15] + 4.25)
	middle[16] = math.Sin(0.27* input[16] + 1.85)
	middle[17] = math.Abs(9.96* input[17] + 7.21) // -1.95 to 18.03, and only these two values
	middle[18] = math.Sin(0.31* input[18] + 5.04)
	middle[19] = math.Sin(0.89* input[19] - 0.18)
	middle[20] = math.Tanh(0.95* input[20] - 0.53)
	middle[21] = math.Sin(0.43* input[21] + 2.24)
	middle[22] = math.Sin(0.41* input[22] + 2.39)
	middle[23] = math.Sin(0.31* input[23] + 1.61)
	middle[24] = math.Sin(0.16* input[24] - 4.2)
	middle[25] = math.Sin(0.18* input[25] - 7.56)  
	middle[26] = math.Sin(0.17* input[26] + 8.5)
	middle[27] = math.Tanh(1.02* input[27] - 0.68)
	middle[28] = math.Sin(0.21* input[28] + 2.16)
	middle[29] = math.Sin(0.23* input[29] - 7.04)
	middle[30] = math.Tanh(0.35* input[30] - 1.43)
	middle[31] = math.Tanh(1.22* input[31] - 2.35)
	middle[32] = math.Tanh(1.51* input[32] - 2.12)
	middle[33] = math.Abs(6.07* input[33] + 2.42) // -2.1 to 56.08
	middle[34] = math.Sin(0.26* input[34] + 2.15)
	middle[35] = math.Sin(0.28* input[35] - 7.78)
	middle[36] = math.Sin(0.26* input[36] + 4.52)

	coefficient := []float64{0.04, -0.07, -0.05, 0.01, -0.32, 
		0.0, 0.0, 0.05, 0.0, 0.1, 
		-0.28, 0.0, 0.0, -0.24, 0.0,
		0.05, 0.12, 0.02, -0.12, 0.02,
		0.03, -0.49, 0.13, 0.24, 0.35,
		0.13, 0.5, 0.02, 0.13, -0.21,  
		-0.11, 0.04, 0.03, -0.01, 0.18,
		-0.06, 0.04,
	}
	output_final := 1.04 - 1.05* math.Sin( MultiplyAndSum(middle, coefficient) + 5.76)
	output_tmp = append(output_tmp, MultiplyAndSum(middle, coefficient) + 5.76)
	return 	output_final, output_tmp
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