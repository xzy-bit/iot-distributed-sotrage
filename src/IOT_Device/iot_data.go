package IOT_Device

import (
	"IOT_Storage/src/Secret_Share"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Student struct {
	name  string
	age   int
	stuId int
}

func SendSliceToNode(nodes []string) {
	stu := Student{
		name:  "XiaoMing",
		age:   18,
		stuId: 1748526,
	}

	stuInfo, _ := json.Marshal(stu)
	matrix := Secret_Share.MatrixInit()
	ciphertext, p := Secret_Share.SliceAndEncrypt(matrix, stuInfo)

	file, _ := os.Open("public.pem")
	iotId := GenerateIotId(file)
	file.Close()

	for index, node := range nodes {
		body := url.Values{
			"cipher":  {ciphertext[index].String()},
			"modNum":  {p.String()},
			"iotId":   {iotId},
			"serial":  {strconv.Itoa(index)},
			"address": {node},
		}
		resp, _ := http.PostForm(node+"/slice", body)
		if resp.StatusCode != 200 {
			log.Fatal("can not send data to nodes")
		}
	}

}
