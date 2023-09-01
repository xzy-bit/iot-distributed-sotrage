package Controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// when iot device send file to node the node save the file
func ReceivePublicKey() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	//router.Static("/", "./static")
	router.POST("/receive", func(context *gin.Context) {
		file, _ := context.FormFile("file")
		log.Println(file.Filename)
		dst := "./" + file.Filename
		context.SaveUploadedFile(file, dst)
		context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
		context.String(200, "server: receive publicKey successfully!")
	})
	return router
}

// when iot device need to create a request to send file to node
func CreateSendFileReq(file *os.File, fileName string, url string) *http.Request {
	bodyBuf := &bytes.Buffer{}
	bodyWrite := multipart.NewWriter(bodyBuf)
	// file 为key
	fileWrite, err := bodyWrite.CreateFormFile("file", fileName)
	_, err = io.Copy(fileWrite, file)
	if err != nil {
		log.Println("err")
	}
	bodyWrite.Close() //要关闭，会将w.w.boundary刷写到w.writer中
	if err != nil {
		log.Println("err")
	}
	req, _ := http.NewRequest("POST", url, bodyBuf)
	// 设置头
	req.Header.Set("Content-Type", bodyWrite.FormDataContentType())
	return req
}

// proxy iot device to send the request
func SendRequest(r *http.Request) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(r)
	if err == nil {
		return resp
	} else {
		log.Fatal(err)
		return nil
	}
}
