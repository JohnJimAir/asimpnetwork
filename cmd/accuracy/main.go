package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	label_true, _ := ReadCSVToFloat64Slice_lastcolumn("../../data/test_data_breast-cancer.csv")
	label_true = Transpose(label_true)
	fmt.Println(label_true)

	result_plain, _ := ReadCSVToFloat64Slice("../../result/KAN_plaintext.csv")
	result_plain = Transpose(result_plain)

	result_cipher, _ := ReadCSVToFloat64Slice("../../result/KAN_ciphertext.csv")
	result_cipher = Transpose(result_cipher)

	
	label_plain := Compare(result_plain[0], result_plain[1])
	label_cipher := Compare(result_cipher[0], result_cipher[1])

	accuracy_plain := CountAccuracy(label_plain, label_true[0])
	accuracy_cipher := CountAccuracy(label_cipher, label_true[0])

	fmt.Println(accuracy_plain, accuracy_cipher)

}

func ReadCSVToFloat64Slice_lastcolumn(filename string) ([][]float64, error) {
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

    // 遍历 CSV 内容，跳过第一行，并只关注最后一列
    for i, record := range records {
        if i == 0 {
            // 跳过第一行（通常是标题行）
            continue
        }

        // 创建一个切片来存储当前行的 float64 数据
        var row []float64
        for j := len(record)-1; j < len(record); j++ { // 只要最后一列
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

    // 遍历 CSV 内容
    for _, record := range records {

        // 创建一个切片来存储当前行的 float64 数据
        var row []float64
        for j := 0; j < len(record); j++ { 
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

func Compare(input_0 []float64, input_1 []float64) (output []float64) {
	
	num := len(input_0)
	output = make([]float64, num)
	for i:=0;i<num;i++ {
		if input_0[i] < input_1[i] {
			output[i] = 1.0
		}
	}
	return output
}

func CountAccuracy(input_0, input_1 []float64) (accuracy float64) {

	num := len(input_0)
	match := 0.0

	for i:=0;i<num;i++ {
		if input_0[i] == input_1[i] {
			match++
		}
	}
	accuracy = match / float64(num)
	return accuracy
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