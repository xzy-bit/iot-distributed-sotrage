package Node

import (
	"IOT_Storage/src/Block_Chain"
	"IOT_Storage/src/Controller"
	"IOT_Storage/src/File_Index"
	"IOT_Storage/src/IOT_Device"
	"IOT_Storage/src/Identity_Verify"
	"IOT_Storage/src/Secret_Share"
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Head *Block_Chain.DataNode
var Tail *Block_Chain.DataNode
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

func NodeIsALive(url string) bool {
	pipe := make(chan bool)
	go func() {
		req := Controller.CreatePingReq(url)

		resp := Controller.SendRequest(req)
		if resp == nil {
			log.Printf("Can not get connection with %s\n", url)
			time.Sleep(time.Second)
			pipe <- false
		} else {
			log.Printf("%s is alive\n", url)
		}
		pipe <- true
	}()
	select {
	case result := <-pipe:
		return result
	}
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

func NodeGetSlice() *gin.Engine {
	router := gin.Default()
	router.POST("slice", func(context *gin.Context) {
		cipherStr := context.PostForm("cipher")
		iotId := context.PostForm("iotId")
		serialStr := context.PostForm("serial")
		address := context.PostForm("address")
		modNumStr := context.PostForm("modNum")
		timeStamp := context.PostForm("timeStamp")
		hash := context.PostForm("hash")
		index := context.PostForm("indexOfGroup")

		num, _ := strconv.Atoi(index)

		dataIndex := GenerateDATA(iotId, serialStr, address, modNumStr, timeStamp, hash, num)
		AddDataToCache(dataIndex)
		log.Println("Add data index to cache...")
		log.Println(Head.Data)

		fileName := "./slices/" + hash + "/_" + index + ".slc"
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
		//if Head == nil {
		//	log.Println("Head is nil!")
		//} else {
		//	log.Println("Head is not nil")
		//}

		// This place exists some problem need to be solved
		data := GetAllDataInCache()
		if len(data) != 0 {
			log.Println("Handling data")
			HandleData(data, nodeConfig.NodeId)
		} else {
			log.Println("Data is nil put token to the next node!")
			time.Sleep(time.Second * 5)
		}

		index := (nodeConfig.NodeId + 1) % 7
		trueUrl := nodeConfig.AddressBook[index] + ":" + strconv.Itoa(nodeConfig.PortForToken+index)
		pingUrl := nodeConfig.AddressBook[index] + ":" + strconv.Itoa(nodeConfig.PortForPIng+index)

		count := 0
		for NodeIsALive(pingUrl) == false && count < 7 {
			index = (index + 1) & 7
			count++
			pingUrl = nodeConfig.AddressBook[index] + ":" + strconv.Itoa(nodeConfig.PortForPIng+index)
			trueUrl = nodeConfig.AddressBook[index] + ":" + strconv.Itoa(nodeConfig.PortForToken+index)
		}

		req, _ := http.NewRequest("GET", trueUrl+"/token", nil)

		log.Printf("Send token to %s\n", trueUrl)

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

		File_Index.InsertBlock(&block, tree)

		Block_Chain.StoreBlock(block)

		context.String(200, "Get block")
	})
	return router
}

func NodeSendSlice() *gin.Engine {
	router := gin.Default()
	router.POST("userGetSlice", func(context *gin.Context) {
		filename := "./slices/" + context.PostForm("filename") + ".slc"
		file, err := os.Open(filename)
		defer file.Close()
		stat, err := file.Stat()

		if err != nil {
			log.Println(err)
			context.String(502, "Can not open file")
		}
		body := make([]byte, stat.Size())
		_, err = bufio.NewReader(file).Read(body)

		log.Println(body)
		if err != nil {
			log.Println(err)
			context.String(502, "Can not read file")
		} else {
			context.Data(200, "text/plain", body)
		}
	})
	return router
}

func SendSlice(rg *gin.RouterGroup) {
	router := rg.Group("/userGetSlice")
	router.POST("/", func(context *gin.Context) {
		filename := "./slices/" + context.PostForm("filename") + ".slc"
		file, err := os.Open(filename)
		defer file.Close()
		stat, err := file.Stat()

		if err != nil {
			log.Println(err)
			context.String(502, "Can not open file")
		}
		body := make([]byte, stat.Size())
		_, err = bufio.NewReader(file).Read(body)

		log.Println(body)
		if err != nil {
			log.Println(err)
			context.String(502, "Can not read file")
		} else {
			context.Data(200, "text/plain", body)
		}
	})
}

func SendSliceSm4(rg *gin.RouterGroup) {
	router := rg.Group("/userGetSliceSM4")
	router.POST("/", func(context *gin.Context) {
		index := context.PostForm("index")
		filename := "./slices/" + context.PostForm("filename") + "/" + index + ".slc"
		file, err := os.Open(filename)
		defer file.Close()
		stat, err := file.Stat()

		if err != nil {
			log.Println(err)
			context.String(502, "Can not open file")
		}
		body := make([]byte, stat.Size())
		_, err = bufio.NewReader(file).Read(body)

		log.Println(body)
		if err != nil {
			log.Println(err)
			context.String(502, "Can not read file")
		} else {
			context.Data(200, "text/plain", body)
		}
	})
}

func NodeGetQuery(rg *gin.RouterGroup) {
	router := rg.Group("/query")
	router.POST("/", func(context *gin.Context) {
		log.Println("Receive query request")
		iotId := context.PostForm("iotId")
		startTime := context.PostForm("startTime")
		endTIme := context.PostForm("endTime")

		start, _ := time.Parse("2006-01-02 15:04:05", startTime)
		end, _ := time.Parse("2006-01-02 15:04:05", endTIme)

		indexes := File_Index.QueryData(tree, iotId, start, end)

		body, err := json.Marshal(indexes)
		if err != nil {
			context.String(502, "Can not get slice")
		} else {
			context.Data(200, "application/json", body)
		}
	})
}

func NodeQuerySlice(address string, hash []byte) []byte {

	body := url.Values{
		"filename": {hex.EncodeToString(hash)},
	}
	resp, _ := http.PostForm(address, body)
	if resp.StatusCode == 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return data
	} else {
		log.Println("Can not get the file")
		return []byte{}
	}
}

func NodeQueryDataForUSer(rg *gin.RouterGroup) {
	var patient IOT_Device.Patient
	router := rg.Group("/queryForUser")
	router.POST("/", func(context *gin.Context) {
		log.Println("Receive query request")
		iotId := context.PostForm("iotId")
		startTime := context.PostForm("startTime")
		endTIme := context.PostForm("endTime")

		start, _ := time.Parse("2006-01-02 15:04:05", startTime)
		end, _ := time.Parse("2006-01-02 15:04:05", endTIme)

		indexes := File_Index.QueryData(tree, iotId, start, end)

		port := 9000
		for i := 0; i < len(indexes); i += 7 {
			count := 0
			var cipher []*big.Int
			var p big.Int
			var choice []int
			p = *indexes[i].ModNum
			for j := 0; j < 7; j++ {
				temp := strings.Split(indexes[i+j].StoreOn, ":")
				trueUrl := temp[0] + ":" + temp[1] + ":" + strconv.Itoa(port+j) + "/userGetSlice"
				slice := NodeQuerySlice(trueUrl, indexes[i+j].Hash)
				//fmt.Println(slice)
				if len(slice) == 0 {
					fmt.Printf("Can not get slice from %s\n", trueUrl)
				} else {
					choice = append(choice, indexes[i+j].Serial)
					num := big.NewInt(1)
					num.SetString(string(slice), 10)
					//fmt.Println(num)
					cipher = append(cipher, num)
					//cipher = append(cipher)
					count++
				}
				if count == 4 {
					msgBytes := Secret_Share.ResotreMsg(cipher, p, choice)
					json.Unmarshal(msgBytes, &patient)
					break
				}
			}
		}
		body, err := json.Marshal(patient)
		if err != nil {
			context.String(502, "Can not get slice")
		} else {
			context.Data(200, "application/json", body)
		}
	})
}
