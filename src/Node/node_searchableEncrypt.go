package Node

import (
	"IOT_Storage/src/Controller"
	"IOT_Storage/src/SearchableEncrypt"
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type SearchIndex struct {
	PatientId string
	TimeStamp time.Time
	Address   string
}

type SplitIndex struct {
	DocSplit1 string
	DocSplit2 string
}

type SplitMat struct {
	MatSplit1 string
	MatSplit2 string
}

func StoreIndex(index SearchIndex) {
	fd, _ := os.OpenFile("indexTable.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	indexInfo, _ := json.Marshal(index)

	fd.Write(indexInfo)
	fd.Write([]byte("\n"))
	fd.Close()
	return
}

func ReadIndexes() {
	file, err := os.Open("indexTable.json")
	if err != nil {
		return
	}
	reader := bufio.NewReader(file)
	for {
		currentLine, fileErr := reader.ReadBytes('\n')
		if fileErr == io.EOF {
			break
		}
		index := SearchIndex{}
		json.Unmarshal(currentLine, &index)
		table = append(table, index)
	}
	log.Println("indexes number:", len(table))
}

func BroadcastIndex(Index SearchIndex) {
	indexInfo, _ := json.Marshal(&Index)

	for index, node := range nodeConfig.AddressBook {
		if index == nodeConfig.NodeId {
			continue
		}
		trueUrl := node + ":" + strconv.Itoa(nodeConfig.PortForIndexBroad+index)

		reader := bytes.NewReader(indexInfo)

		req, _ := http.NewRequest("GET", trueUrl+"/indexes", reader)
		req.Header.Set("Content-Type", "application/json")
		Controller.SendRequest(req)
	}
}

func NodeForIndexBroad() *gin.Engine {
	router := gin.Default()
	router.GET("indexes", func(context *gin.Context) {
		//var blockInfo []byte
		var index SearchIndex

		body, err := context.GetRawData() // 读取 request body 的内容
		if err != nil {
			log.Println("failed to get body")
		}
		context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // 创建 io.ReadCloser 对象传给 request body
		err = json.Unmarshal(body, &index)
		if err != nil {
			log.Println("failed to create index")
		}

		table = append(table, index)
		StoreIndex(index)

		context.String(200, "Get index")
	})
	return router
}

func GenerateFileName(index SearchIndex) string {
	info := bytes.Join([][]byte{
		[]byte(index.PatientId),
		[]byte(index.TimeStamp.String()),
		[]byte(index.Address),
	}, []byte{})
	h := sha256.Sum256(info)
	temp := h[:]
	hash := hex.EncodeToString(temp)
	return hash
}

func NodeGetIndex(rg *gin.RouterGroup) {
	router := rg.Group("/getIndex")
	router.POST("/", func(context *gin.Context) {
		vector := context.PostForm("vector")
		iotId := context.PostForm("iotId")
		address := context.PostForm("address")
		timeStamp := context.PostForm("timeStamp")

		stamp, _ := time.Parse("2006-01-02 15:04:05", timeStamp)

		index := SearchIndex{
			PatientId: iotId,
			TimeStamp: stamp,
			Address:   address,
		}

		table = append(table, index)
		StoreIndex(index)
		log.Println("Add index to table...")

		BroadcastIndex(index)
		log.Println("Broadcasting index to nodes...")

		hash := GenerateFileName(index)
		fileName := "./indexes/" + hash + ".idx"
		//context.SaveUploadedFile()
		SaveSlice(vector, fileName)
		context.String(200, "Get index")
	})
}

func NodeGetSpltMat(rg *gin.RouterGroup) {
	router := rg.Group("/getSplitMat")
	router.POST("/", func(context *gin.Context) {

		mat1 := context.PostForm("mat_split1")
		mat2 := context.PostForm("mat_split2")

		//log.Println(mat1)
		//log.Println(mat2)

		splitMat := SplitMat{
			MatSplit1: mat1,
			MatSplit2: mat2,
		}

		//log.Println(splitMat.matSplit1)
		//log.Println(splitMat.matSplit2)

		indexInfo, _ := json.Marshal(splitMat)
		//log.Println(err)
		fileName := "splitMat.sk"
		//context.SaveUploadedFile()
		SaveJson(indexInfo, fileName)
		context.String(200, "Get the SplitMat")
	})
}

func NodeGetIndexWithSplitMat(rg *gin.RouterGroup) {
	router := rg.Group("/getIndexWithSplitMat")
	router.POST("/", func(context *gin.Context) {

		doc_split1 := context.PostForm("doc_split1")
		doc_split2 := context.PostForm("doc_split2")
		iotId := context.PostForm("iotId")
		address := context.PostForm("address")
		timeStamp := context.PostForm("timeStamp")

		stamp, _ := time.Parse("2006-01-02 15:04:05", timeStamp)

		index := SearchIndex{
			PatientId: iotId,
			TimeStamp: stamp,
			Address:   address,
		}

		split_index := SplitIndex{
			DocSplit1: doc_split1,
			DocSplit2: doc_split2,
		}

		table = append(table, index)
		StoreIndex(index)
		log.Println("Add index to table...")

		BroadcastIndex(index)
		log.Println("Broadcasting index to nodes...")

		indexInfo, _ := json.Marshal(&split_index)
		hash := GenerateFileName(index)
		fileName := "./indexes/" + hash + ".idx"
		log.Println(fileName)
		//context.SaveUploadedFile()
		SaveJson(indexInfo, fileName)
		context.String(200, "Get the IndexWithSplitMat")
	})
}

func NodeSendIndex() *gin.Engine {
	router := gin.Default()
	router.POST("nodeQueryIndex", func(context *gin.Context) {
		filename := "./indexes/" + context.PostForm("filename") + ".idx"
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

func NodeForKeyWordsQueryWithSplitMat(rg *gin.RouterGroup) {
	router := rg.Group("/queryByKeyWordsWithSplitMat")
	router.POST("/", func(context *gin.Context) {
		var query SearchableEncrypt.QueryRequest
		var documnetScores []SearchableEncrypt.DocumentRank

		t11 := context.PostForm("t11")
		t12 := context.PostForm("t12")
		t21 := context.PostForm("t21")
		t22 := context.PostForm("t22")

		query.T11.UnmarshalBinary([]byte(t11))
		query.T12.UnmarshalBinary([]byte(t12))
		query.T21.UnmarshalBinary([]byte(t21))
		query.T22.UnmarshalBinary([]byte(t22))
		log.Println("Receive query vector")

		nodes := nodeConfig.AddressBook[:4]
		for i := 0; i < len(table); i += 4 {
			var document SearchableEncrypt.Document
			for j := 0; j < 4; j++ {
				var data []byte
				temp := table[i+j]
				fileName := GenerateFileName(temp)
				if j == 0 {
					data, _ = os.ReadFile("./indexes/" + fileName + ".idx")
					document.I11.UnmarshalBinary(data)
					continue
				}

				trueUrl := nodes[j] + ":" + strconv.Itoa(nodeConfig.PortForSendIndex+j)

				body := url.Values{
					"filename": {fileName},
				}
				resp, _ := http.PostForm(trueUrl+"/nodeQueryIndex", body)

				if resp.StatusCode != 200 {
					log.Println("can not send data to nodes")
				}
				data, err := ioutil.ReadAll(resp.Body)
				//resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
				if err != nil {
					log.Println("can not get the data")
				}
				if j == 1 {
					document.I12.UnmarshalBinary(data)
				} else if j == 2 {
					document.I21.UnmarshalBinary(data)
				} else {
					document.I22.UnmarshalBinary(data)
				}
			}
			result := SearchableEncrypt.Query(&document, &query)
			rank := SearchableEncrypt.DocumentRank{
				UserID:    table[i].PatientId,
				TimeStamp: table[i].TimeStamp,
				Score:     result,
			}
			documnetScores = append(documnetScores, rank)
			log.Println("table hash:", table[i].PatientId)
			log.Println("time stamp:", table[i].TimeStamp)
			log.Println("query score:", result)
		}

		body, err := json.Marshal(documnetScores)
		if err != nil {
			context.String(502, "Can not get indexes")
		} else {
			context.Data(200, "application/json", body)
		}
	})
}
func NodeForKeyWordsQuery(rg *gin.RouterGroup) {
	router := rg.Group("/queryByKeyWords")
	router.POST("/", func(context *gin.Context) {
		var query SearchableEncrypt.QueryRequest
		var documnetScores []SearchableEncrypt.DocumentRank

		t11 := context.PostForm("t11")
		t12 := context.PostForm("t12")
		t21 := context.PostForm("t21")
		t22 := context.PostForm("t22")

		query.T11.UnmarshalBinary([]byte(t11))
		query.T12.UnmarshalBinary([]byte(t12))
		query.T21.UnmarshalBinary([]byte(t21))
		query.T22.UnmarshalBinary([]byte(t22))
		log.Println("Receive query vector")

		nodes := nodeConfig.AddressBook[:4]
		for i := 0; i < len(table); i += 4 {
			var document SearchableEncrypt.Document
			for j := 0; j < 4; j++ {
				var data []byte
				temp := table[i+j]
				fileName := GenerateFileName(temp)
				if j == 0 {
					data, _ = os.ReadFile("./indexes/" + fileName + ".idx")
					document.I11.UnmarshalBinary(data)
					continue
				}

				trueUrl := nodes[j] + ":" + strconv.Itoa(nodeConfig.PortForSendIndex+j)

				body := url.Values{
					"filename": {fileName},
				}
				resp, _ := http.PostForm(trueUrl+"/nodeQueryIndex", body)

				if resp.StatusCode != 200 {
					log.Println("can not send data to nodes")
				}
				data, err := ioutil.ReadAll(resp.Body)
				//resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
				if err != nil {
					log.Println("can not get the data")
				}
				if j == 1 {
					document.I12.UnmarshalBinary(data)
				} else if j == 2 {
					document.I21.UnmarshalBinary(data)
				} else {
					document.I22.UnmarshalBinary(data)
				}
			}
			result := SearchableEncrypt.Query(&document, &query)
			rank := SearchableEncrypt.DocumentRank{
				UserID:    table[i].PatientId,
				TimeStamp: table[i].TimeStamp,
				Score:     result,
			}
			documnetScores = append(documnetScores, rank)
			log.Println("table hash:", table[i].PatientId)
			log.Println("time stamp:", table[i].TimeStamp)
			log.Println("query score:", result)
		}

		body, err := json.Marshal(documnetScores)
		if err != nil {
			context.String(502, "Can not get indexes")
		} else {
			context.Data(200, "application/json", body)
		}
	})
}
