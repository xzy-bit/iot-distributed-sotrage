package Secret_Share

import (
	"encoding/json"
	"fmt"
	"math/big"
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
