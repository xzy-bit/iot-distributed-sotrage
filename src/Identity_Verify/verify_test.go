package Identity_Verify

import (
	"fmt"
	"testing"
)

func TestVerification(t *testing.T) {
	GenerateKey(false)
	message := []byte("Hello World!")
	rText, sText := Sign(message, "private.pem")
	result := Verify(message, rText, sText, "public.pem")
	fmt.Println(result)
}
