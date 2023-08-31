package Controller

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUploadFile(t *testing.T) {
	bodyBuf := &bytes.Buffer{}
	bodyWrite := multipart.NewWriter(bodyBuf)
	file, err := os.Open("./images/rabbit.png")
	defer file.Close()
	if err != nil {
		log.Println("err")
	}
	// file 为key
	fileWrite, err := bodyWrite.CreateFormFile("file", "rabbit.png")
	_, err = io.Copy(fileWrite, file)
	if err != nil {
		log.Println("err")
	}
	bodyWrite.Close() //要关闭，会将w.w.boundary刷写到w.writer中
	if err != nil {
		log.Println("err")
	}
	req, _ := http.NewRequest("POST", "/upload", bodyBuf)
	// 设置头
	req.Header.Set("Content-Type", bodyWrite.FormDataContentType())

	router := UploadFile()
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
