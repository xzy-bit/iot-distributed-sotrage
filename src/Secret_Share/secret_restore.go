package Secret_Share

import (
	"fmt"
	"math/big"
)

func mulMod(p *big.Int, a *big.Int, b *big.Int) *big.Int {
	divide := big.NewInt(0)
	divide.ModInverse(b, p)
	divide.Mul(a, divide)
	divide.Mod(divide, p)
	return divide
}

// restore message from ciphertext and choice (which slices to choose)
func ResotreMsg(ciphertext []*big.Int, p big.Int, choice []int) []byte {
	matrix := MatrixInit()
	advance := make([][]*big.Int, 4)
	for i := 0; i < 4; i++ {
		advance[i] = make([]*big.Int, 5)
		for j := 0; j < 4; j++ {
			advance[i][j] = matrix[choice[i]][j]
		}
		advance[i][4] = ciphertext[i]
	}
	//fmt.Println("p:" + p.String())
	for i := 1; i < 4; i++ {
		for j := i; j < 4; j++ {
			//divide.ModInverse(advance[i-1][i-1], &p)
			//fmt.Println("inverse:" + divide.String())
			//divide.Mul(advance[j][i-1], divide)
			//divide.Mod(divide, &p)
			//fmt.Println("a*b^{-1}mod p:" + divide.String())

			divide := mulMod(&p, advance[j][i-1], advance[i-1][i-1])
			for k := 0; k < 5; k++ {
				tempMul := big.NewInt(0)
				tempMul.Mul(advance[i-1][k], divide)
				advance[j][k].Sub(advance[j][k], tempMul)
				advance[j][k].Mod(advance[j][k], &p)
			}
		}
	}
	//fmt.Println(advance)

	for i := 3; i > 0; i-- {
		prime := &p
		advance[i][4] = mulMod(prime, advance[i][4], advance[i][i])
		//advance[i][4].Div(advance[i][4], advance[i][i])
		advance[i][i] = big.NewInt(1)
		for j := i - 1; j >= 0; j-- {

			//divide.ModInverse(advance[i][i], &p)
			//divide.Mul(advance[j][i], divide)
			divide := mulMod(&p, advance[j][i], advance[i][i])
			for k := 0; k < 5; k++ {
				tempMul := big.NewInt(0)
				tempMul.Mul(advance[i][k], divide)
				advance[j][k].Sub(advance[j][k], tempMul)
				advance[j][k].Mod(advance[j][k], &p)
			}
		}
	}

	//fmt.Println(advance)
	var MsgBin string
	for i := 0; i < 4; i++ {
		binaryStr := fmt.Sprintf("%b", advance[i][4])
		MsgBin += binaryStr[1:]
	}
	//fmt.Println(MsgBin)

	MsgInt := big.NewInt(0)
	MsgInt.SetString(MsgBin, 2)

	return MsgInt.Bytes()
}
