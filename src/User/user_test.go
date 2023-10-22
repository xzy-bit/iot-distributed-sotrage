package User

import (
	"IOT_Storage/src/IOT_Device"
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
	msg := QueryDataWithSM4("6866974dc2b54eb8c5363f174f6f0c0e8d6a69ba3d9035957dceff69992402f2", node, startTime, endTime, 9000, "123456")
	for j := 0; j < len(sm4Msg); j++ {
		for k := 0; k < len(sm4Msg[j]); k++ {
			fmt.Println(sm4Msg[j][k] == msg[j][k])
		}
	}
}

func TestQueryData(t *testing.T) {
	var patient IOT_Device.Patient
	node := "http://192.168.42.129:8000"

	startTime := "2023-10-22 14:14:00"
	endTime := "2023-10-22 14:14:00"
	msg := QueryDataWithSM4("6866974dc2b54eb8c5363f174f6f0c0e8d6a69ba3d9035957dceff69992402f2", node, startTime, endTime, 9000, "123456")
	patient = RestoreStructFromMsg(msg)
	fmt.Println(patient)
}

func TestQueryByKeyWordsWithSplitMat(t *testing.T) {
	//nodes := []string{
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//}
	//SearchableEncrypt.SendSplitMat(nodes)
	query := []string{
		"内科",
		"心率失常",
		"胸闷",
	}
	QueryByKeyWordsWithSm4(query)
}
