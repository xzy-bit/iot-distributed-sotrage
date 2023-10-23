package Patient_Data

import (
	"IOT_Storage/src/SM4"
	"IOT_Storage/src/SearchableEncrypt"
	"IOT_Storage/src/Secret_Share"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tjfoc/gmsm/sm3"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Patient_Test struct {
	Name      string
	Age       int
	PatientId int
	KeyWords  []string
}

type Patient struct {
	Identity    string
	Name        string
	Age         int
	Country     string
	Nation      string
	Sex         string
	Match       string
	Career      string
	TimeStamp   time.Time
	Description string
	Image       []byte
	Features    []string
}

func SendSliceAndIndexWithSplitMat(nodes []string, portForSlice int, portForIndex int) {
	var sliceNode [7]string
	var indexNode [7]string

	for i := 0; i < 7; i++ {
		sliceNode[i] = nodes[i] + ":" + strconv.Itoa(i+portForSlice)
		indexNode[i] = nodes[i] + ":" + strconv.Itoa(i+portForIndex)
	}

	document := []string{
		"精神科",
		"食欲不振",
		"记忆力衰退",
	}
	//documentCompare := []string{
	//	"精神科",
	//	"心律不齐",
	//	"记忆力衰退",
	//}

	//document1 := []string{
	//	"胃科",
	//	"食欲不振",
	//	"四肢乏力",
	//}
	//document1Compare := []string{
	//	"胃科",
	//	"心律不齐",
	//	"食欲不振",
	//}

	patient := Patient_Test{
		Name:      "Li",
		Age:       60,
		PatientId: 1002,
		KeyWords:  document,
	}

	patientInfo, _ := json.Marshal(patient)
	matrix := Secret_Share.MatrixInit()
	ciphertext, p := Secret_Share.SliceAndEncrypt(matrix, patientInfo)

	file, _ := os.Open("public.pem")
	iotId := GenerateIotId(file)
	file.Close()

	timeStamp := time.Now()
	fmt.Println(timeStamp.Format("2006-01-02 15:04:05"))

	SearchableEncrypt.SendIndexWithSplitMat(indexNode[:4], patient.KeyWords, iotId, timeStamp)

	for index, node := range sliceNode {
		body := url.Values{
			"cipher":    {ciphertext[index].String()},
			"modNum":    {p.String()},
			"iotId":     {iotId},
			"serial":    {strconv.Itoa(index)},
			"address":   {node},
			"timeStamp": {timeStamp.Format("2006-01-02 15:04:05")},
		}
		resp, _ := http.PostForm(node+"/slice", body)
		if resp.StatusCode != 200 {
			log.Fatal("can not send data to nodes")
		}
	}

	fmt.Println("Indexes and slices were successfully sent to nodes")
}

func Sm3(str string) string {
	sm3.New()
	hw := sm3.New()
	hw.Write([]byte(str))
	hash := hw.Sum(nil)
	return hex.EncodeToString(hash)
}

func UploadSliceAndIndexWithSplitMat(patient Patient) {
	nodes := []string{
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
	}
	portForSlice := 10080
	portForIndex := 8040

	var sliceNode [7]string
	var indexNode [7]string

	for i := 0; i < 7; i++ {
		sliceNode[i] = nodes[i] + ":" + strconv.Itoa(i+portForSlice)
		indexNode[i] = nodes[i] + ":" + strconv.Itoa(i+portForIndex)
	}

	patientInfo, _ := json.Marshal(patient)
	matrix := Secret_Share.MatrixInit()
	patientInfo = SM4.PaddingWithBytes(patientInfo)

	iotId := Sm3(patient.Identity)
	timeStamp := patient.TimeStamp

	SearchableEncrypt.SendIndexWithSplitMat(indexNode[:4], patient.Features, iotId, timeStamp)

	sm4Msg := SM4.EncryptWithPadding(patientInfo, "123456")
	var hashTable [][]byte
	for _, node := range sliceNode {
		hash := SliceHash(iotId, timeStamp, node)
		hashTable = append(hashTable, hash)
	}
	modNum := Secret_Share.FixedPara()

	numOfGroup := len(sm4Msg)
	for i := 0; i < numOfGroup; i++ {
		ciphertext, _ := Secret_Share.SliceAndEncryptWithFixedPara(matrix, sm4Msg[i], modNum)
		for index, node := range sliceNode {
			body := url.Values{
				"cipher":       {ciphertext[index].String()},
				"numOfGroup":   {strconv.Itoa(numOfGroup)},
				"iotId":        {iotId},
				"serial":       {strconv.Itoa(index)},
				"address":      {node},
				"timeStamp":    {timeStamp.Format("2006-01-02 15:04:05")},
				"indexOfGroup": {strconv.Itoa(i)},
				"hash":         {hex.EncodeToString(hashTable[index])},
			}
			resp, _ := http.PostForm(node+"/slice", body)
			if resp.StatusCode != 200 {
				log.Println("group:", i, "slice:", index, "can not send to ", node)
			}
		}
	}
	fmt.Println("Sending complete")
	fmt.Println("NumOfGroup:", numOfGroup)
	fmt.Println("Indexes and slices were successfully sent to nodes")
}

func SendSliceAndIndexToNode(nodes []string, portForSlice int, portForIndex int) {
	var sliceNode [7]string
	var indexNode [7]string

	for i := 0; i < 7; i++ {
		sliceNode[i] = nodes[i] + ":" + strconv.Itoa(i+portForSlice)
		indexNode[i] = nodes[i] + ":" + strconv.Itoa(i+portForIndex)
	}

	//document := []string{
	//	"精神科",
	//	"食欲不振",
	//	"记忆力衰退",
	//}
	//documentCompare := []string{
	//	"精神科",
	//	"心律不齐",
	//	"记忆力衰退",
	//}

	//document1 := []string{
	//	"胃科",
	//	"食欲不振",
	//	"四肢乏力",
	//}
	document1Compare := []string{
		"胃科",
		"心律不齐",
		"食欲不振",
	}

	patient := Patient_Test{
		Name:      "QIAN",
		Age:       20,
		PatientId: 1002,
		KeyWords:  document1Compare,
	}

	patientInfo, _ := json.Marshal(patient)
	matrix := Secret_Share.MatrixInit()
	ciphertext, p := Secret_Share.SliceAndEncrypt(matrix, patientInfo)

	file, _ := os.Open("public.pem")
	iotId := GenerateIotId(file)
	file.Close()

	timeStamp := time.Now()
	fmt.Println(timeStamp.Format("2006-01-02 15:04:05"))

	SearchableEncrypt.SendIndex(indexNode[:4], patient.KeyWords, iotId, timeStamp)

	for index, node := range sliceNode {
		body := url.Values{
			"cipher":    {ciphertext[index].String()},
			"modNum":    {p.String()},
			"iotId":     {iotId},
			"serial":    {strconv.Itoa(index)},
			"address":   {node},
			"timeStamp": {timeStamp.Format("2006-01-02 15:04:05")},
		}
		resp, _ := http.PostForm(node+"/slice", body)
		if resp.StatusCode != 200 {
			log.Fatal("can not send data to nodes")
		}
	}

	fmt.Println("Indexes and slices were successfully sent to nodes")
}

func SliceHash(DeviceID string, TimeStamp time.Time, address string) []byte {
	info := bytes.Join([][]byte{
		[]byte(DeviceID),
		[]byte(TimeStamp.String()),
		[]byte(address),
	}, []byte{})
	hash := sha256.Sum256(info)
	Hash := hash[:]
	return Hash
}

func SendSliceWithSM4(data []byte, nodes []string, password string, portForSlice int) int {
	var sliceNode [7]string

	for i := 0; i < 7; i++ {
		sliceNode[i] = nodes[i] + ":" + strconv.Itoa(i+portForSlice)
	}

	padding := SM4.PaddingWithBytes(data)
	sm4Msg := SM4.EncryptWithPadding(padding, password)
	//fmt.Println(SM4.WithdrawPadding(SM4.DecryptWithPadding(sm4Msg, password)))
	matrix := Secret_Share.MatrixInit()

	file, _ := os.Open("public.pem")

	iotId := GenerateIotId(file)
	file.Close()
	timeStamp := time.Now()
	fmt.Println(timeStamp.Format("2006-01-02 15:04:05"))

	var hashTable [][]byte
	for _, node := range sliceNode {
		hash := SliceHash(iotId, timeStamp, node)
		hashTable = append(hashTable, hash)
	}
	modNum := Secret_Share.FixedPara()

	//var test [][]byte
	for i := 0; i < len(sm4Msg); i++ {
		ciphertext, p := Secret_Share.SliceAndEncryptWithFixedPara(matrix, sm4Msg[i], modNum)
		//choice := []int{0, 1, 2, 3}
		//msg := Secret_Share.ResotreMsg(ciphertext, *modNum, choice)
		//test = append(test, msg)

		for index, node := range sliceNode {
			body := url.Values{
				"cipher":       {ciphertext[index].String()},
				"modNum":       {p.String()},
				"iotId":        {iotId},
				"serial":       {strconv.Itoa(index)},
				"address":      {node},
				"timeStamp":    {timeStamp.Format("2006-01-02 15:04:05")},
				"indexOfGroup": {strconv.Itoa(i)},
				"hash":         {hex.EncodeToString(hashTable[index])},
			}
			resp, _ := http.PostForm(node+"/slice", body)
			if resp.StatusCode != 200 {
				log.Println("group:", i, "slice:", index, "can not send to ", node)
			}
		}
	}
	fmt.Println("Sending complete")

	return len(sm4Msg)
}

//func SendSM4Slice(data []byte, nodes []string, password string, portForSlice int) {
//	var sliceNode [7]string
//
//	for i := 0; i < 7; i++ {
//		sliceNode[i] = nodes[i] + ":" + strconv.Itoa(i+portForSlice)
//	}
//
//	sm4Msg := SM4.Encrypt(data, password)
//	matrix := Secret_Share.MatrixInit()
//
//	file, _ := os.Open("public.pem")
//
//	iotId := GenerateIotId(file)
//	file.Close()
//	timeStamp := time.Now()
//	fmt.Println(timeStamp.Format("2006-01-02 15:04:05"))
//
//	var hashTable [][]byte
//	for _, node := range sliceNode {
//		hash := SliceHash(iotId, timeStamp, node)
//		hashTable = append(hashTable, hash)
//	}
//
//	modNum := Secret_Share.FixedPara()
//	ciphertext, p := Secret_Share.SliceAndEncryptWithFixedPara(matrix, sm4Msg, modNum)
//	for index, node := range sliceNode {
//		body := url.Values{
//			"cipher":    {ciphertext[index].String()},
//			"modNum":    {p.String()},
//			"iotId":     {iotId},
//			"serial":    {strconv.Itoa(index)},
//			"address":   {node},
//			"timeStamp": {timeStamp.Format("2006-01-02 15:04:05")},
//			"hash":      {hex.EncodeToString(hashTable[index])},
//		}
//		resp, _ := http.PostForm(node+"/slice", body)
//		if resp.StatusCode != 200 {
//			log.Println("slice:", index, "can not send to ", node)
//		}
//	}
//	fmt.Println("Sending complete")
//}
