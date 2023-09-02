package User

import (
	"IOT_Storage/src/Controller"
	"IOT_Storage/src/Identity_Verify"
	"IOT_Storage/src/Node"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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

func SignForRandom(url string) bool {
	reqForChallenge, _ := http.NewRequest("GET", url+"/challenge", nil)
	resp := Controller.SendRequest(reqForChallenge)
	if resp.StatusCode != 200 {
		log.Fatal("cannot get random")
		return false
	}

	body, _ := io.ReadAll(resp.Body)
	//str := string(body)
	//random := new(big.Int)
	//random, _ = random.SetString(str, 10)
	//fmt.Println(random)
	rText, sText := Identity_Verify.Sign(body, "private.pem")
	sign := Node.Sign{
		RText: rText,
		SText: sText,
	}
	signBytes, _ := json.Marshal(sign)
	reader := bytes.NewReader(signBytes)
	reqForSign, _ := http.NewRequest("POST", url+"/sign", reader)
	reqForSign.Header.Set("Content-Type", "application/json")
	resp = Controller.SendRequest(reqForSign)
	if resp.StatusCode != 200 {
		log.Fatal("cannot get random")
		return false
	}
	return true
}
