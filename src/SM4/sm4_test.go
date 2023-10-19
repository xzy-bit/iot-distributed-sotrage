package SM4

import (
	"fmt"
	"testing"
)

func TestPaddingWithBytes(t *testing.T) {
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, byte(1))
	}
	fmt.Println(data)
	fmt.Println(len(data))
	data = PaddingWithBytes(data)
	fmt.Println(data)
	fmt.Println(len(data))
}

func TestEncryptWithPadding(t *testing.T) {
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, byte(1))
	}
	data = PaddingWithBytes(data)
	password := "123456"
	output := EncryptWithPadding(data, password)
	fmt.Println(len(output))
	for i := 0; i < len(output); i++ {
		fmt.Println(output[i])
	}
}

func TestDecryptWithPadding(t *testing.T) {
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, byte(1))
	}
	data = PaddingWithBytes(data)
	fmt.Println(data)
	password := "123456"
	output := EncryptWithPadding(data, password)
	plainPadding := DecryptWithPadding(output, password)
	fmt.Println(plainPadding)
	fmt.Println(len(plainPadding))
	plain := WithdrawPadding(plainPadding)
	fmt.Println(len(plain))
	fmt.Println(plain)
}
