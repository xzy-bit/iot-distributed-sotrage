package Secret_Share

import (
	"IOT_Storage/src/SM4"
	"encoding/json"
	"fmt"
	"github.com/tjfoc/gmsm/sm4"
	"math/big"
	"os"
	"testing"
)

type Student struct {
	Name  string
	Age   int
	StuId int
}

func TestSliceAndEncrypt(t *testing.T) {
	stu := Student{
		Name:  "XiaoMing",
		Age:   18,
		StuId: 1748526,
	}

	stuInfo, _ := json.Marshal(stu)

	matrix := MatrixInit()
	ciphertext, p := SliceAndEncrypt(matrix, stuInfo)
	//fmt.Println(ciphertext)
	cipher := []*big.Int{
		ciphertext[0],
		ciphertext[1],
		ciphertext[2],
		ciphertext[3],
	}
	MsgBytes := ResotreMsg(cipher, p, []int{0, 1, 2, 3})

	var result Student
	json.Unmarshal(MsgBytes, &result)
	fmt.Println(result)
}

func TestSlice(t *testing.T) {
	data := make([]byte, 60000000)

	//data := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
	key := SM4.GenerateSM4Key("123456")
	input, _ := os.Open("1.txt")
	output, _ := os.OpenFile("indexes.idx", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer input.Close()
	defer output.Close()

	input.Read(data)

	iv := []byte("0000000000000000")
	err := sm4.SetIV(iv)
	fmt.Printf("err = %v\n", err)

	cbcDec, err := sm4.Sm4Cbc(key, data, true)

	matrix := MatrixInit()
	ciphertext, _ := SliceAndEncrypt(matrix, cbcDec)
	for i := 0; i < 7; i++ {
		temp := ciphertext[i].Bytes()
		output.Write(temp)
	}
}
