package Secret_Share

import (
	"fmt"
	"math/big"
)

// calculate a*b^{-1} mod p
func mulMod(prime *big.Int, number *big.Int, inverse *big.Int) {
	inverse.ModInverse(inverse, prime)
	number.Mul(number, inverse)
	number.Mod(number, prime)
}

// restore message from ciphertext and choice (which slices to choose)
func restoreMsg(ciphertext []*big.Int, p big.Int, choice []int) []byte {
	matrix := MatrixInit()
	advance := make([][]*big.Int, 4)
	for i := 0; i < 4; i++ {
		advance[i] = make([]*big.Int, 5)
		for j := 0; j < 4; j++ {
			advance[i][j] = matrix[i][j]
		}
		advance[i][4] = ciphertext[choice[i]]
	}
	for i := 1; i < 4; i++ {
		divide := big.NewInt(0)
		for j := i; j < 4; j++ {
			divide.Div(advance[j][i-1], advance[i-1][i-1])
			for k := 0; k < 5; k++ {
				tempMul := big.NewInt(0)
				tempMul.Mul(advance[i-1][k], divide)
				advance[j][k].Sub(advance[j][k], tempMul)
			}
		}
	}
	fmt.Println(advance)

	for i := 3; i > 0; i-- {
		prime := &p
		mulMod(prime, advance[i][4], advance[i][i])
		//advance[i][4].Div(advance[i][4], advance[i][i])
		advance[i][i] = big.NewInt(1)
		for j := i - 1; j >= 0; j-- {
			divide := big.NewInt(0)
			divide.Div(advance[j][i], advance[i][i])
			for k := 0; k < 5; k++ {
				tempMul := big.NewInt(0)
				tempMul.Mul(advance[i][k], divide)
				advance[j][k].Sub(advance[j][k], tempMul)
			}
		}
	}

	fmt.Println(advance)
	var MsgBin string
	for i := 0; i < 4; i++ {
		binaryStr := fmt.Sprintf("%b", advance[i][4])
		MsgBin += binaryStr[1:]
	}
	fmt.Println(MsgBin)

	MsgInt := big.NewInt(0)
	MsgInt.SetString(MsgBin, 2)

	return MsgInt.Bytes()
}
