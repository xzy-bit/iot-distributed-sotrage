package Secret_Share

import "math/big"

func SliceAndEncrypt(matrix [][]*big.Int, message []byte) ([]*big.Int, big.Int) {
	messageStr := Data2String(message)
	L := len(messageStr)
	modNum := big.NewInt(1)
	baseNum := big.NewInt(2)
	for i := 0; i < L; i++ {
		modNum.Mul(modNum, baseNum)
	}

	ONE := big.NewInt(1)
	TWO := big.NewInt(2)
	p := modNum
	p.Add(p, ONE)

	for {
		if p.ProbablyPrime(10) {
			break
		}
		p.Add(p, TWO)
	}

	var splitMessage [4]*big.Int
	splitLen := (L / 4) + 1
	for i := 0; i < 4; i++ {
		tempStr := ""
		tempInt := big.NewInt(0)
		if (i+1)*splitLen >= L {
			tempStr = "1" + messageStr[i*splitLen:]
		} else {
			tempStr = "1" + messageStr[i*splitLen:(i+1)*splitLen]
		}
		tempInt, _ = tempInt.SetString(tempStr, 2)
		splitMessage[i] = tempInt
	}

	ciphertext := make([]*big.Int, 7)
	for i := 0; i < 7; i++ {
		var temp *big.Int
		temp = new(big.Int)
		ciphertext[i] = big.NewInt(0)
		for j := 0; j < 4; j++ {
			temp = big.NewInt(0)
			temp.Mul(matrix[i][j], splitMessage[j])
			ciphertext[i].Add(ciphertext[i], temp)
		}
		ciphertext[i].Mod(ciphertext[i], p)
	}
	return ciphertext, *p
}
