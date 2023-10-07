package SearchableEncrypt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Document struct {
	I11 mat.VecDense
	I12 mat.VecDense
	I21 mat.VecDense
	I22 mat.VecDense
}

type QueryRequest struct {
	T11 mat.VecDense
	T12 mat.VecDense
	T21 mat.VecDense
	T22 mat.VecDense
}

type SecretKey struct {
	M1 mat.Dense
	M2 mat.Dense
	S  mat.VecDense
}

type SecretKeyBinary struct {
	M1 []byte
	M2 []byte
	S  []byte
}

type DocumentRank struct {
	UserID    string
	TimeStamp time.Time
	Score     float64
}

func MatrixPrint(m mat.Matrix) {
	formattedMatrix := mat.Formatted(m, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", formattedMatrix)
}

func GenerateRandomMatrix(n int) *mat.Dense {
	data := make([]float64, n*n)
	for i := range data {
		data[i] = rand.NormFloat64()
	}
	matrix := mat.NewDense(n, n, data)
	return matrix
}

func GenerateInvertibleMatrix(n int) *mat.Dense {
	matrix := GenerateRandomMatrix(n)
	var Inv mat.Dense
	for {
		err := Inv.Inverse(matrix)
		if err != nil {
			matrix.Reset()
			matrix = GenerateRandomMatrix(n)
			continue
		}
		break
	}
	return matrix
}

func GenerateVector(n int) *mat.VecDense {
	data := make([]float64, n)
	for i := range data {
		flag := rand.Intn(10)
		//println("flag", flag)
		if flag >= 5 {
			data[i] = 1
		}
	}
	vector := mat.NewVecDense(n, data)
	return vector
}

func SetUp(n int) *SecretKey {
	var sk SecretKey
	m1 := GenerateInvertibleMatrix(n + 2)
	m2 := GenerateInvertibleMatrix(n + 2)
	s := GenerateVector(n + 2)
	sk.M1 = *m1
	sk.M2 = *m2
	sk.S = *s
	return &sk
}

func AddKeyWords(words []string) {
	file, err := os.OpenFile("keywords.txt", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	data := bytes.NewBuffer([]byte{})
	for i := range words {
		data.WriteString(words[i] + "\n")
	}
	file.Write(data.Bytes())
	file.Close()
}

func ReadKeyWords() []string {
	var keyWords []string
	file, err := os.Open("keywords.txt")
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	for {
		currentLine, err := reader.ReadBytes('\n')
		word := strings.ReplaceAll(string(currentLine), "\n", "")
		keyWords = append(keyWords, word)
		if err == io.EOF {
			break
		}
	}
	return keyWords
}

func GenerateDocumentVector(subKey []string) *mat.VecDense {
	keyWords := ReadKeyWords()
	n := len(keyWords)
	data := make([]float64, n+2)
	for i := range subKey {
		for j := range keyWords {
			if subKey[i] == keyWords[j] {
				data[j] = 1
				break
			} else {
				continue
			}
		}
	}
	data[n] = 0.5
	data[n+1] = 1
	docVec := mat.NewVecDense(n+2, data)
	return docVec
}

func GenerateQueryVector(subKey []string) *mat.VecDense {
	var r float64
	keyWords := ReadKeyWords()
	n := len(keyWords)
	data := make([]float64, n+2)
	for {
		r = rand.NormFloat64()
		if r > 0 {
			break
		}
	}
	//r = 1
	for i := range subKey {
		for j := range keyWords {
			if subKey[i] == keyWords[j] {
				data[j] = r
				break
			} else {
				continue
			}
		}
	}
	data[n] = r
	data[n+1] = rand.NormFloat64()
	docVec := mat.NewVecDense(n+2, data)
	return docVec
}

func splitVec(vec *mat.VecDense, S *mat.VecDense, flag int) (*mat.VecDense, *mat.VecDense) {
	length := S.Len()
	data1 := make([]float64, length)
	data2 := make([]float64, length)
	for i := 0; i < length; i++ {
		if S.AtVec(i) == 0 {
			if flag == 0 {
				data1[i] = vec.AtVec(i)
				data2[i] = vec.AtVec(i)
			} else {
				data1[i] = rand.NormFloat64()
				data2[i] = vec.AtVec(i) - data1[i]
			}
		} else {
			if flag == 0 {
				data1[i] = rand.NormFloat64()
				data2[i] = vec.AtVec(i) - data1[i]
			} else {
				data1[i] = vec.AtVec(i)
				data2[i] = vec.AtVec(i)
			}
		}
	}
	vec1 := mat.NewVecDense(length, data1)
	vec2 := mat.NewVecDense(length, data2)
	return vec1, vec2
}

func BuildIndex(documentKeyWords []string, sk *SecretKey) *Document {
	var i11 mat.VecDense
	var i12 mat.VecDense
	var i21 mat.VecDense
	var i22 mat.VecDense
	var document Document

	docVec := GenerateDocumentVector(documentKeyWords)

	//println("docVec")
	//MatrixPrint(docVec)

	p1, p2 := splitVec(docVec, &sk.S, 0)
	p11, p12 := splitVec(p1, &sk.S, 0)
	p21, p22 := splitVec(p2, &sk.S, 0)

	i11.MulVec(sk.M1.T(), p11)
	i12.MulVec(sk.M2.T(), p12)
	i21.MulVec(sk.M1.T(), p21)
	i22.MulVec(sk.M2.T(), p22)

	document.I11 = i11
	document.I12 = i12
	document.I21 = i21
	document.I22 = i22

	return &document
}

func Trapdoor(queryKeywords []string, sk *SecretKey) *QueryRequest {
	var t11 mat.VecDense
	var t12 mat.VecDense
	var t21 mat.VecDense
	var t22 mat.VecDense
	var query QueryRequest
	var M1_Inverse mat.Dense
	var M2_Inverse mat.Dense

	qryVec := GenerateQueryVector(queryKeywords)
	//println("qryVec")
	//MatrixPrint(qryVec)

	q1, q2 := splitVec(qryVec, &sk.S, 1)
	q11, q12 := splitVec(q1, &sk.S, 1)
	q21, q22 := splitVec(q2, &sk.S, 1)

	M1_Inverse.Inverse(&sk.M1)
	M2_Inverse.Inverse(&sk.M2)

	t11.MulVec(&M1_Inverse, q11)
	t12.MulVec(&M2_Inverse, q12)
	t21.MulVec(&M1_Inverse, q21)
	t22.MulVec(&M2_Inverse, q22)

	query.T11 = t11
	query.T12 = t12
	query.T21 = t21
	query.T22 = t22
	return &query
}

func Query(doc *Document, query *QueryRequest) float64 {
	var result1 mat.VecDense
	var result2 mat.VecDense
	var result3 mat.VecDense
	var result4 mat.VecDense
	//var result mat.VecDense

	result1.MulVec(query.T11.T(), &doc.I11)
	result2.MulVec(query.T12.T(), &doc.I12)
	result3.MulVec(query.T21.T(), &doc.I21)
	result4.MulVec(query.T22.T(), &doc.I22)

	return result1.At(0, 0) + result2.At(0, 0) + result3.At(0, 0) + result4.At(0, 0)
}

func QueryForUser(subKey []string, documents []Document, sk *SecretKey) []float64 {
	query := Trapdoor(subKey, sk)
	results := []float64{}
	for i := range documents {
		result := Query(&documents[i], query)
		results = append(results, result)
	}

	return results
}

func GenerateSk() {
	var sk *SecretKey
	var sk_binary SecretKeyBinary
	keywords := ReadKeyWords()
	sk = SetUp(len(keywords))

	m1info, _ := sk.M1.MarshalBinary()
	m2info, _ := sk.M1.MarshalBinary()
	sinfo, _ := sk.S.MarshalBinary()

	sk_binary.M1 = m1info
	sk_binary.M2 = m2info
	sk_binary.S = sinfo

	skInfo, _ := json.Marshal(sk_binary)
	file, err := os.OpenFile("user.sk", os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	data := bytes.NewBuffer(skInfo)
	file.Write(data.Bytes())
	file.Close()
}

func ReadSk() *SecretKey {
	data, err := os.ReadFile("user.sk")
	if err != nil {
		panic(err)
		log.Fatal("Read sk error!")
	}
	sk_binary := new(SecretKeyBinary)
	json.Unmarshal(data, sk_binary)

	sk := new(SecretKey)

	sk.M1.UnmarshalBinary(sk_binary.M1)
	sk.M2.UnmarshalBinary(sk_binary.M2)
	sk.S.UnmarshalBinary(sk_binary.S)

	return sk
}

func SendIndex(nodes []string, dVector []string, iotId string, timeStamp time.Time) {
	var indexes [4][]byte
	sk := ReadSk()

	docInx := BuildIndex(dVector, sk)

	I11, _ := docInx.I11.MarshalBinary()
	I12, _ := docInx.I12.MarshalBinary()
	I21, _ := docInx.I21.MarshalBinary()
	I22, _ := docInx.I22.MarshalBinary()

	indexes[0] = I11
	indexes[1] = I12
	indexes[2] = I21
	indexes[3] = I22

	//fmt.Println(timeStamp.Format("2006-01-02 15:04:05"))
	for index, node := range nodes {
		body := url.Values{
			"vector":    {string(indexes[index])},
			"iotId":     {iotId},
			"address":   {node},
			"timeStamp": {timeStamp.Format("2006-01-02 15:04:05")},
		}
		resp, _ := http.PostForm(node+"/getIndex", body)
		if resp.StatusCode != 200 {
			log.Fatal("can not send data to nodes")
		}
	}
}

func SendQuery(queryKeyWords []string) {
	sk := ReadSk()
	queryInx := Trapdoor(queryKeyWords, sk)
	node := "http://192.168.42.129"

	T11, _ := queryInx.T11.MarshalBinary()
	T12, _ := queryInx.T12.MarshalBinary()
	T21, _ := queryInx.T21.MarshalBinary()
	T22, _ := queryInx.T22.MarshalBinary()

	body := url.Values{
		"t11": {string(T11)},
		"t12": {string(T12)},
		"t21": {string(T21)},
		"t22": {string(T22)},
	}
	resp, _ := http.PostForm(node+"/queryByKeyWords", body)
	if resp.StatusCode != 200 {
		log.Fatal("can not send data to nodes")
	}
}
