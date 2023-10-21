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
	//data, _ = os.ReadFile("test.jpg")
	data = PaddingWithBytes(data)
	fmt.Println(data)
	password := "123456"
	output := EncryptWithPadding(data, password)
	plainPadding := DecryptWithPadding(output, password)
	plain := WithdrawPadding(plainPadding)
	fmt.Println(plain)
	//outFile, _ := os.OpenFile("output.jpg", os.O_RDWR|os.O_CREATE, 0755)
	//outFile.Write(plain)
	//defer outFile.Close()
}
