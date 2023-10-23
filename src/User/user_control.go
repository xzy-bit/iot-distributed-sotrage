package User

import (
	"IOT_Storage/src/Block_Chain"
	"IOT_Storage/src/Patient_Data"
	"IOT_Storage/src/SM4"
	"IOT_Storage/src/SearchableEncrypt"
	"IOT_Storage/src/Secret_Share"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

type PatientRank struct {
	Score   float64
	Patient Patient_Data.Patient
}

//func ReceiveKeys() *gin.Engine {
//	router := gin.Default()
//	router.MaxMultipartMemory = 8 << 20
//	//router.Static("/", "./static")
//	router.POST("/receive", func(context *gin.Context) {
//		file, _ := context.FormFile("file")
//		log.Println(file.Filename)
//		dst := "./" + file.Filename
//		if file.Size == 0 {
//			context.String(http.StatusNotFound, fmt.Sprintf("no file get", file.Filename))
//			context.String(404, "user: can not receive the file")
//		} else {
//			context.SaveUploadedFile(file, dst)
//			context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
//			context.String(200, "user: receive Key from iot device successfully!")
//		}
//	})
//	return router
//}

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

func QueryData(node string, startTime string, endTime string, port int) []Patient_Data.Patient_Test {
	var patient Patient_Data.Patient_Test
	var patients []Patient_Data.Patient_Test
	file, _ := os.Open("public.pem")
	iotId := Patient_Data.GenerateIotId(file)
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

	p := Secret_Share.FixedPara()
	for i := 0; i < len(indexes); i += 7 {
		count := 0
		var cipher []*big.Int
		var choice []int
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
				msgBytes := Secret_Share.ResotreMsg(cipher, *p, choice)
				json.Unmarshal(msgBytes, &patient)
				patients = append(patients, patient)
				log.Println(patient)
				break
			}
		}
	}
	fmt.Println("End of data querying!")
	return patients
}

func QueryDataWithSM4(identity string, node string, startTime string, endTime string, port int, password string) [][]byte {
	var msg [][]byte

	body := url.Values{
		"iotId":     {identity},
		"startTime": {startTime},
		"endTime":   {endTime},
	}
	resp, _ := http.PostForm(node+"/query", body)
	if resp.StatusCode != 200 {
		log.Fatal("can not send query to nodes")
	}
	data, err := ioutil.ReadAll(resp.Body)
	//resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("can not get the data")
	}

	var indexes []Block_Chain.DATA
	json.Unmarshal(data, &indexes)

	//for i := 0; i < len(indexes); i++ {
	//	fmt.Println(indexes[i])
	//}
	p := Secret_Share.FixedPara()
	numOfGroup := indexes[1].NumOfGroup

	for i := 0; i < len(indexes); i += 7 {
		for k := 0; k < numOfGroup; k++ {
			count := 0
			var cipher []*big.Int
			var choice []int
			for j := 0; j < 7; j++ {
				temp := strings.Split(indexes[i+j].StoreOn, ":")
				trueUrl := temp[0] + ":" + temp[1] + ":" + strconv.Itoa(port+j) + "/userGetSlice"
				indexOfGroup := strconv.Itoa(k)
				slice := UserGetSliceWithSM4(trueUrl, indexes[i+j].Hash, indexOfGroup)
				//fmt.Println(slice)
				if len(slice) == 0 {
					fmt.Printf("Can not get slice from %s\n", trueUrl)
				} else {
					choice = append(choice, indexes[i+j].Serial)
					num := big.NewInt(1)
					num.SetString(string(slice), 10)
					cipher = append(cipher, num)
					count++
				}
				if count == 4 {
					msgBytes := Secret_Share.ResotreMsg(cipher, *p, choice)
					if len(msgBytes) < 64 && k != numOfGroup-1 {
						padding := make([]byte, 64-len(msgBytes))
						msgBytes = SM4.BytesCombine(padding, msgBytes)
					}
					msg = append(msg, msgBytes)
					//fmt.Println(msgBytes)
					break
				}
			}
		}
	}

	fmt.Println("End of data querying!")
	//final := SM4.DecryptWithPadding(msg, "123456")
	//finalFile, _ := os.OpenFile("final.jpg", os.O_RDWR|os.O_CREATE, 0755)

	//defer finalFile.Close()
	//fmt.Println(final)
	//plain := SM4.WithdrawPadding(final)
	//fmt.Println(plain)
	//finalFile.Write(plain)
	return msg
}

func RestoreStructFromMsg(Msg [][]byte) Patient_Data.Patient {
	var patient Patient_Data.Patient
	final := SM4.DecryptWithPadding(Msg, "123456")
	plain := SM4.WithdrawPadding(final)
	json.Unmarshal(plain, &patient)
	return patient
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

func UserGetSliceWithSM4(address string, hash []byte, index string) []byte {

	body := url.Values{
		"filename": {hex.EncodeToString(hash) + "/" + index},
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

func QueryDocumentRank(scores []SearchableEncrypt.DocumentRank) []PatientRank {
	var patients []PatientRank
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
		//fmt.Println("identity:", document.UserID)
		msg := QueryDataWithSM4(document.UserID, nodeToQuery, startTime, endTime, portForSendSlice, "123456")
		patient := RestoreStructFromMsg(msg)
		temp := PatientRank{
			Score:   document.Score,
			Patient: patient,
		}
		patients = append(patients, temp)
		fmt.Println(patient)
	}
	return patients
}

func QueryByKeyWords(query []string) {
	fmt.Println("query key words:", query)
	documentScores := SearchableEncrypt.QueryByKeyWords(query)
	QueryDocumentRank(documentScores)
}

func QueryByKeyWordsWithSplitMat(query []string) {
	documentScores := SearchableEncrypt.QueryByKeyWordsWithSplitMat(query)
	QueryDocumentRank(documentScores)
}

func QueryByKeyWordsWithSm4(query []string) []PatientRank {
	documentScores := SearchableEncrypt.QueryByKeyWordsWithSplitMat(query)
	return QueryDocumentRank(documentScores)
}
