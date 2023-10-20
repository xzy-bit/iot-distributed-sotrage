package SM4

import (
	"bytes"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
)

func GenerateSM4Key(password string) []byte {
	h := sm3.New()
	h.Write([]byte(password))
	key := h.Sum(nil)
	return key[16:]
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func NumOfPadding(num int) []byte {
	result := make([]byte, num)
	for i := 0; i < num; i++ {
		result[i] = uint8(num)
	}
	return result
}

func PaddingWithBytes(data []byte) []byte {
	var result []byte
	length := len(data)
	numberOfGroup := length/63 + 1

	for i := 0; i < numberOfGroup; i++ {
		padding := uint8(i)
		if i == numberOfGroup-1 {
			result = BytesCombine(result, data[i*63:])
			result = BytesCombine(result, []byte{padding})
		} else {
			result = BytesCombine(result, data[i*63:(i+1)*63])
			result = BytesCombine(result, []byte{padding})
		}
	}
	return result
}

func Encrypt(data []byte, password string) []byte {
	key := GenerateSM4Key(password)
	result, _ := sm4.Sm4Cbc(key, data, true)
	return result
}

func EncryptWithPadding(data []byte, password string) [][]byte {
	var output [][]byte
	key := GenerateSM4Key(password)
	length := len(data)
	numOfGroup := length / 64
	result, _ := sm4.Sm4Cbc(key, data, true)
	for i := 0; i < numOfGroup; i++ {
		if i == numOfGroup-1 {
			output = append(output, result[i*64:])
		} else {
			output = append(output, result[i*64:(i+1)*64])
		}
	}
	return output
}

func DecryptWithPadding(cipher [][]byte, password string) []byte {
	var input []byte
	for i := 0; i < len(cipher); i++ {
		input = BytesCombine(input, cipher[i])
	}
	key := GenerateSM4Key(password)
	result, _ := sm4.Sm4Cbc(key, input, false)
	return result
}

func WithdrawPadding(plain []byte) []byte {
	var result []byte
	numOfGroup := len(plain)/64 + 1
	for i := 0; i < numOfGroup; i++ {
		if i == numOfGroup-1 {
			result = BytesCombine(result, plain[i*64:len(plain)-1])
		} else {
			result = BytesCombine(result, plain[i*64:(i+1)*64-1])
		}
	}
	return result
}
