package Controller

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"net/http/httptest"
	"os"
	"testing"
)

func TestReceiveFile(t *testing.T) {
	file, _ := os.Open("./files/rabbit.png")
	if file != nil {
		fmt.Println("OK")
	}
	router := ReceivePublicKey()
	req := CreateSendFileReq(file, "rabbit.png", "http://localhost:8080/receive")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
