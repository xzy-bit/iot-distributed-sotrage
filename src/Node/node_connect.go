package Node

import (
	"IOT_Storage/src/Block_Chain"
	"IOT_Storage/src/Identity_Verify"
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"log"
	"math/big"
)

var head *Block_Chain.DataNode
var tail *Block_Chain.DataNode

type Sign struct {
	RText []byte
	SText []byte
}

func Challenge() *gin.Engine {
	var sign Sign
	var random *big.Int
	router := gin.Default()
	router.GET("challenge", func(context *gin.Context) {
		random, _ = rand.Int(rand.Reader, big.NewInt(1073741824))
		context.String(200, random.String())
	})
	router.POST("sign", func(context *gin.Context) {
		if context.ShouldBindJSON(&sign) == nil {
			log.Println(sign.RText)
			log.Println(sign.SText)
		}
		randomBytes, _ := random.MarshalJSON()
		result := Identity_Verify.Verify(randomBytes, sign.RText, sign.SText, "public.pem")
		if result == false {
			context.String(502, "Your identification's verification does not pass!")
		} else {
			context.String(200, "OK")
		}
	})
	return router
}

func GetSlice() *gin.Engine {
	router := gin.Default()
	router.POST("slice", func(context *gin.Context) {
		cipherStr, _ := context.Get("cipher")
		pStr, _ := context.Get("modNum")
		iotId, _ := context.Get("iotId")

		log.Println(cipherStr)
		log.Println(pStr)
		log.Println(iotId)

		context.String(200, "Get slice")
	})
	return router
}
