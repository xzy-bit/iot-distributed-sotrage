package User

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ReceiveKeys() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	//router.Static("/", "./static")
	router.POST("/receive", func(context *gin.Context) {
		file, _ := context.FormFile("file")
		log.Println(file.Filename)
		dst := "./" + file.Filename
		if file.Size == 0 {
			context.String(http.StatusNotFound, fmt.Sprintf("no file get", file.Filename))
			context.String(404, "user: can not receive the file")
		} else {
			context.SaveUploadedFile(file, dst)
			context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
			context.String(200, "user: receive Key from iot device successfully!")
		}
	})
	return router
}

func Ping() *gin.Engine {
	router := gin.Default()
	router.GET("ping", func(context *gin.Context) {
		context.String(200, "pong")

	})
	return router
}
