package IOT_Device

import (
	"IOT_Storage/src/Secret_Share"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Student struct {
	Name  string
	Age   int
	StuId int
}

func SendSliceToNode(nodes []string) {
	stu := Student{
		Name:  "dwjklas",
		Age:   34,
		StuId: 3242,
	}

	stuInfo, _ := json.Marshal(stu)
	matrix := Secret_Share.MatrixInit()
	ciphertext, p := Secret_Share.SliceAndEncrypt(matrix, stuInfo)

	file, _ := os.Open("public.pem")
	iotId := GenerateIotId(file)
	file.Close()

	timeStamp := time.Now()
	fmt.Println(timeStamp.Format("2006-01-02 15:04:05"))
	for index, node := range nodes {
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

}
