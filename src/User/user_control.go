package User

import (
	"IOT_Storage/src/Block_Chain"
	"IOT_Storage/src/IOT_Device"
	"IOT_Storage/src/SearchableEncrypt"
	"IOT_Storage/src/Secret_Share"
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
	"sort"
	"strconv"
	"strings"
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

//func SignForRandom(url string) bool {
//	reqForChallenge, _ := http.NewRequest("GET", url+"/challenge", nil)
//	resp := Controller.SendRequest(reqForChallenge)
//	if resp.StatusCode != 200 {
//		log.Fatal("cannot get random")
//		return false
//	}
//
//	body, _ := io.ReadAll(resp.Body)
//	//str := string(body)
//	//random := new(big.Int)
//	//random, _ = random.SetString(str, 10)
//	//fmt.Println(random)
//	rText, sText := Identity_Verify.Sign(body, "private.pem")
//	sign := Node.Sign{
//		RText: rText,
//		SText: sText,
//	}
//	signBytes, _ := json.Marshal(sign)
//	reader := bytes.NewReader(signBytes)
//	reqForSign, _ := http.NewRequest("POST", url+"/sign", reader)
//	reqForSign.Header.Set("Content-Type", "application/json")
//	resp = Controller.SendRequest(reqForSign)
//	if resp.StatusCode != 200 {
//		log.Fatal("cannot get random")
//		return false
//	}
//	fmt.Println(resp.StatusCode)
//	return true
//}

func QueryData(node string, startTime string, endTime string, port int) {
	file, _ := os.Open("public.pem")
	iotId := IOT_Device.GenerateIotId(file)
	println(iotId)
	defer file.Close()
	body := url.Values{
		"iotId":     {iotId},
		"startTime": {startTime},
		"endTime":   {endTime},
	}
	resp, _ := http.PostForm(node+"/query", body)
	if resp.StatusCode != 200 {
		log.Fatal("can not send data to nodes")
	}
	data, err := ioutil.ReadAll(resp.Body)
	//resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("can not get the data")
	}

	var indexes []Block_Chain.DATA
	json.Unmarshal(data, &indexes)

	for i := 0; i < len(indexes); i += 7 {
		count := 0
		var cipher []*big.Int
		var p big.Int
		var choice []int
		p = *indexes[i].ModNum
		for j := 0; j < 7; j++ {
			temp := strings.Split(indexes[i+j].StoreOn, ":")
			trueUrl := temp[0] + ":" + temp[1] + ":" + strconv.Itoa(port+j) + "/userGetSlice"
			slice := UserGetSlice(trueUrl, indexes[i+j].Hash)
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
				var patient IOT_Device.Patient
				json.Unmarshal(msgBytes, &patient)
				fmt.Println(patient)
				break
			}
		}
	}
	fmt.Println("End of data querying!")
}

func UserGetSlice(address string, hash []byte) []byte {

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

func QueryDocumentRank(scores []SearchableEncrypt.DocumentRank) {
	sort.Sort(SearchableEncrypt.DocumentScores(scores))
	for index, document := range scores {
		if index == 3 {
			break
		}
		portForSendSlice := 9000
		nodeToQuery := "http://192.168.42.129:8000"
		startTime := document.TimeStamp.Format("2006-01-02 15:04:05")
		endTime := document.TimeStamp.Format("2006-01-02 15:04:05")
		fmt.Println("document score:", document.Score)
		QueryData(nodeToQuery, startTime, endTime, portForSendSlice)
	}
}

func QueryByKeyWords(query []string) {
	fmt.Println("query key words:", query)
	documentScores := SearchableEncrypt.QueryByKeyWords(query)
	QueryDocumentRank(documentScores)
}

func QueryByKeyWorsWithSplitMat(query []string) {
	documentScores := SearchableEncrypt.QueryByKeyWords(query)
	QueryDocumentRank(documentScores)
}
