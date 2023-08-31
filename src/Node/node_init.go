package Node

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping() *gin.Engine {
	router := gin.Default()
	router.GET("ping", func(context *gin.Context) {
		context.String(200, "pong")

	})
	return router
}

func Login() *gin.Engine {
	router := gin.Default()
	router.GET("login", func(context *gin.Context) {
		context.String(200, "login")
	})
	return router
}

func Send() {
	url := "http://localhost:8080/ping"
	resp, _ := http.Get(url)
	fmt.Println(resp.StatusCode)
}
