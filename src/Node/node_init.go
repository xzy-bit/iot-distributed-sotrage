package Node

import (
	"github.com/gin-gonic/gin"
)

func Ping() *gin.Engine {
	router := gin.Default()
	router.GET("ping", func(context *gin.Context) {
		context.String(200, "pong")

	})
	return router
}
