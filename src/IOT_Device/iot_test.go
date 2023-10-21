package IOT_Device

import (
	"IOT_Storage/src/SM4"
	"IOT_Storage/src/Secret_Share"
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestSendSliceWithSM4(t *testing.T) {
	file, _ := os.OpenFile("test.jpg", os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Read the file into a byte slice
	buffer := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	nodes := []string{
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
	}
	portForSlice := 10080
	password := "123456"
	SendSliceWithSM4(buffer, nodes, password, portForSlice)
}

func TestSendSM4Slice(t *testing.T) {
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, byte(1))
	}
	password := "123456"
	nodes := []string{
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
	}
	portForSlice := 10080
	SendSliceWithSM4(data, nodes, password, portForSlice)
}

func TestSliceHash(t *testing.T) {
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, byte(1))
	}
	password := "123456"
	padding := SM4.PaddingWithBytes(data)
	sm4Msg := SM4.EncryptWithPadding(padding, password)
	//fmt.Println(SM4.WithdrawPadding(SM4.DecryptWithPadding(sm4Msg, password)))
	matrix := Secret_Share.MatrixInit()

	modNum := Secret_Share.FixedPara()

	var test [][]byte
	for i := 0; i < len(sm4Msg); i++ {
		ciphertext, _ := Secret_Share.SliceAndEncryptWithFixedPara(matrix, sm4Msg[i], modNum)
		choice := []int{0, 1, 2, 3}
		msg := Secret_Share.ResotreMsg(ciphertext, *modNum, choice)
		test = append(test, msg)
	}

	//for j := 0; j < len(sm4Msg); j++ {
	//	for k := 0; k < len(sm4Msg[j]); k++ {
	//		fmt.Println(sm4Msg[j][k] == test[j][k])
	//	}
	//}
	fmt.Println(modNum.String())
	padding = SM4.DecryptWithPadding(test, password)
	sm4padding := SM4.DecryptWithPadding(sm4Msg, password)
	sm4plain := SM4.WithdrawPadding(sm4padding)
	fmt.Println(sm4plain)
	plain := SM4.WithdrawPadding(padding)
	fmt.Println(plain)
}
