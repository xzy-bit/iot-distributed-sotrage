package Secret_Share

import (
	"fmt"
	"math/big"
)

func MatrixInit() [][]*big.Int {
	tempMatrix := [7][4]int64{
		{1, 1, 1, 1},
		{1, 2, 3, 4},
		{1, 4, 9, 16},
		{1, 8, 27, 64},
		{1, 16, 81, 256},
		{1, 32, 243, 1024},
		{1, 64, 729, 4096},
	}
	matrix := make([][]*big.Int, 7)
	for i := 0; i < 7; i++ {
		matrix[i] = make([]*big.Int, 4)
	}
	for i := 0; i < 7; i++ {
		for j := 0; j < 4; j++ {
			matrix[i][j] = new(big.Int).SetInt64(tempMatrix[i][j])
		}
	}
	return matrix
}

func Data2String(data []byte) string {
	intFormat := new(big.Int)
	intFormat.SetBytes(data)
	binStr := fmt.Sprintf("%b", intFormat)
	return binStr
}
