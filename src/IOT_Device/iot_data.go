package IOT_Device

import (
	"IOT_Storage/src/SM4"
	"IOT_Storage/src/SearchableEncrypt"
	"IOT_Storage/src/Secret_Share"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Patient struct {
	Name      string
	Age       int
	PatientId int
	KeyWords  []string
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

	patient := Patient{
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

	patient := Patient{
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

func SendIndexWithSplitMat(nodes []string, portForIndex int) {
	var indexNode [7]string

	for i := 0; i < 7; i++ {
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

	patient := Patient{
		Name:      "Li",
		Age:       60,
		PatientId: 1002,
		KeyWords:  document,
	}

	file, _ := os.Open("public.pem")
	iotId := GenerateIotId(file)
	file.Close()

	timeStamp := time.Now()
	fmt.Println(timeStamp.Format("2006-01-02 15:04:05"))

	SearchableEncrypt.SendIndexWithSplitMat(indexNode[:4], patient.KeyWords, iotId, timeStamp)

	fmt.Println("Indexes were successfully sent to nodes")
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

func SendSliceWithSM4(data []byte, nodes []string, password string, portForSlice int) {
	var sliceNode [7]string

	for i := 0; i < 7; i++ {
		sliceNode[i] = nodes[i] + ":" + strconv.Itoa(i+portForSlice)
	}

	padding := SM4.PaddingWithBytes(data)
	sm4Msg := SM4.EncryptWithPadding(padding, password)
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

	for i := 0; i < len(sm4Msg); i++ {
		ciphertext, p := Secret_Share.SliceAndEncrypt(matrix, sm4Msg[i])
		for index, node := range sliceNode {
			body := url.Values{
				"cipher":    {ciphertext[index].String()},
				"modNum":    {p.String()},
				"iotId":     {iotId},
				"serial":    {strconv.Itoa(index)},
				"address":   {node},
				"timeStamp": {timeStamp.Format("2006-01-02 15:04:05")},
				"index":     {strconv.Itoa(i)},
				"hash":      {hex.EncodeToString(hashTable[index])},
			}
			resp, _ := http.PostForm(node+"/slice", body)
			if resp.StatusCode != 200 {
				log.Fatal("can not send data to nodes")
			}
		}
	}
	fmt.Println("Slices were successfully sent to nodes")

}
