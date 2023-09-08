package Secret_Share

import (
	"fmt"
	"math/big"
	"testing"
)

func TestSliceAndEncrypt(t *testing.T) {
	text := "Hello World!"
	message := []byte(text)
	matrix := MatrixInit()
	ciphertext, p := SliceAndEncrypt(matrix, message)
	//fmt.Println(ciphertext)
	cipher := []*big.Int{
		ciphertext[0],
		ciphertext[2],
		ciphertext[3],
		ciphertext[5],
	}
	MsgBytes := restoreMsg(cipher, p, []int{0, 2, 3, 5})
	fmt.Println(string(MsgBytes))
	if string(MsgBytes) != text {
		t.Errorf("Slice and then restore failed\n")
	}
	fmt.Println("Slice and then restore succeed")
}
