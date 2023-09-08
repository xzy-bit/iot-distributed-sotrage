package Node

import (
	"IOT_Storage/src/Block_Chain"
	"IOT_Storage/src/Controller"
	"IOT_Storage/src/File_Index"
	"IOT_Storage/src/Identity_Verify"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var head *Block_Chain.DataNode
var tail *Block_Chain.DataNode
var mutex sync.Mutex

type Sign struct {
	RText []byte
	SText []byte
}

func Ping() *gin.Engine {
	router := gin.Default()
	router.GET("ping", func(context *gin.Context) {
		context.String(200, "pong")

	})
	return router
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
		log.Println(result)
	})
	return router
}

func ServerGetSlice() *gin.Engine {
	router := gin.Default()
	router.POST("slice", func(context *gin.Context) {
		cipherStr := context.PostForm("cipher")
		iotId := context.PostForm("iotId")
		serialStr := context.PostForm("serial")
		address := context.PostForm("address")
		modNumStr := context.PostForm("modNum")

		log.Println(cipherStr)
		log.Println(iotId)
		log.Println(serialStr)
		log.Println(address)
		log.Println(modNumStr)

		dataIndex := GenerateDATA(iotId, serialStr, address, modNumStr)
		AddDataToCache(head, tail, dataIndex)

		hash := hex.EncodeToString(dataIndex.Hash)
		fileName := "./slices/" + hash + ".slc"
		//context.SaveUploadedFile()
		SaveSlice(cipherStr, fileName)
		context.String(200, "Get slice")
	})
	return router
}

func NodeGetToken() *gin.Engine {
	router := gin.Default()
	router.GET("token", func(context *gin.Context) {
		log.Println("Receive token...")

		// Generate block from cache
		data := GetAllDataInCache(head, tail)
		if data != nil {
			log.Println("Handling data")
			HandleData(data)
		} else {
			log.Println("Data is nil put token to the next node!")
			time.Sleep(time.Second)
		}

		index := (nodeConfig.NodeId + 1) % 7
		trueUrl := nodeConfig.AddressBook[index] + ":" + strconv.Itoa(nodeConfig.PortForToken+index)

		req, _ := http.NewRequest("GET", trueUrl+"/token", nil)
		Controller.SendRequest(req)

		context.String(200, "Send token to next node")
	})
	return router
}

func NodeGetBlock() *gin.Engine {
	router := gin.Default()
	router.GET("block", func(context *gin.Context) {
		//var blockInfo []byte
		var block Block_Chain.Block

		body, err := context.GetRawData() // 读取 request body 的内容
		if err != nil {
			log.Println("failed to get body")
		}
		context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // 创建 io.ReadCloser 对象传给 request body
		err = json.Unmarshal(body, &block)
		if err != nil {
			log.Println("failed to create block")
		}
		log.Println(block)

		File_Index.InsertBlock(&block, tree)
		Block_Chain.StoreBlock(block)

		context.String(200, "Get block")
	})
	return router
}
