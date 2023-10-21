package User

import (
	"IOT_Storage/src/SM4"
	"IOT_Storage/src/Secret_Share"
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestQueryDataWithSM4(t *testing.T) {
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

	padding := SM4.PaddingWithBytes(buffer)
	sm4Msg := SM4.EncryptWithPadding(padding, "123456")
	matrix := Secret_Share.MatrixInit()
	modNum := Secret_Share.FixedPara()
	var test [][]byte
	for i := 0; i < len(sm4Msg); i++ {
		ciphertext, _ := Secret_Share.SliceAndEncryptWithFixedPara(matrix, sm4Msg[i], modNum)
		choice := []int{0, 1, 2, 3}
		msg_ := Secret_Share.ResotreMsg(ciphertext, *modNum, choice)

		if len(msg_) < 64 && i != len(sm4Msg)-1 {
			padding := make([]byte, 64-len(msg_))
			msg_ = SM4.BytesCombine(padding, msg_)
		}

		test = append(test, msg_)
	}

	for j := 0; j < len(sm4Msg); j++ {
		count := 0
		fmt.Println("line", j)
		for k := 0; k < len(sm4Msg[j]); k++ {
			if sm4Msg[j][k] != test[j][k] {
				count++
			}
		}
		if count == 0 {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
	}

	node := "http://192.168.42.129:8000"
	startTime := "2023-10-21 18:56:10"
	endTime := "2023-10-21 18:56:10"
	msg := QueryDataWithSM4(node, startTime, endTime, 2218, 9000, "123456")
	for j := 0; j < len(sm4Msg); j++ {
		for k := 0; k < len(sm4Msg[j]); k++ {
			fmt.Println(sm4Msg[j][k] == msg[j][k])
		}
	}
}

func TestQueryData(t *testing.T) {
	node := "http://192.168.42.129:8000"

	startTime := "2023-10-21 18:52:54"
	endTime := "2023-10-21 18:52:54"
	QueryDataWithSM4(node, startTime, endTime, 17, 9000, "123456")
}
