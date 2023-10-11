package Node

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NodeIndexPageForUser(rg *gin.RouterGroup) {
	router := rg.Group("/index")
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{})
	})
}
