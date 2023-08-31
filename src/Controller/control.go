package Controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func UploadFile() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	//router.Static("/", "./static")
	router.POST("/upload", func(context *gin.Context) {
		file, _ := context.FormFile("file")
		log.Println(file.Filename)
		dst := "./" + file.Filename
		context.SaveUploadedFile(file, dst)
		context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
		context.String(200, "uploaded")
	})
	return router
}
