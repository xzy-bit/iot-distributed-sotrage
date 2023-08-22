package Secret_Share

import (
	"fmt"
	"testing"
)

func TestSliceAndEncrypt(t *testing.T) {
	text := "Hello World!"
	message := []byte(text)
	matrix := MatrixInit()
	ciphertext, p := SliceAndEncrypt(matrix, message)
	//fmt.Println(ciphertext)
	MsgBytes := restoreMsg(ciphertext, p, []int{0, 1, 2, 3})
	fmt.Println(string(MsgBytes))
	if string(MsgBytes) != text {
		t.Errorf("Slice and then restore failed\n")
	}
	fmt.Println("Slice and then restore succeed")
}
