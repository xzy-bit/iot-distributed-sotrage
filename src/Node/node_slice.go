package Node

import (
	"IOT_Storage/src/Block_Chain"
	"IOT_Storage/src/Controller"
	"IOT_Storage/src/File_Index"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// save the sclie and generate DATA struct
func GenerateDATA(iotId string, serial string, address string, modNum string) *Block_Chain.DATA {
	num, _ := strconv.Atoi(serial)
	mod := new(big.Int)
	mod.SetString(modNum, 10)
	//log.Println("modNum:" + mod.String())
	dataIndex := Block_Chain.DATA{
		DeviceID:  iotId,
		TimeStamp: time.Now().UTC(),
		Serial:    num,
		Hash:      nil,
		StoreOn:   address,
		ModNum:    mod,
	}
	Block_Chain.DataHash(&dataIndex)
	return &dataIndex
}

func AddDataToCache(head *Block_Chain.DataNode, tail *Block_Chain.DataNode, newData *Block_Chain.DATA) {
	dataNode := Block_Chain.DataNode{
		Data: *newData,
		Next: nil,
	}
	if head == nil {
		head = &dataNode
		tail = head
		return
	}
	tail.Next = &dataNode
	tail = tail.Next
	return
}

func SaveSlice(cipher string, fileName string) {
	os.MkdirAll(filepath.Dir(fileName), 0750)
	out, _ := os.Create(fileName)
	defer out.Close()
	reader := bytes.NewReader([]byte(cipher))
	io.Copy(out, reader)
}

func GetAllDataInCache(head *Block_Chain.DataNode, tail *Block_Chain.DataNode) []Block_Chain.DATA {
	var data []Block_Chain.DATA

	if head == nil {
		return nil
	}

	for head != tail {
		data = append(data, head.Data)
		head = head.Next
	}

	return data
}

func HandleData(data []Block_Chain.DATA) {
	mutex.Lock()

	oldBlock := Block_Chain.GetPrevBlock()
	log.Println(oldBlock)

	newBlock := Block_Chain.GenerateBlock(oldBlock, data)
	log.Println(newBlock)

	File_Index.InsertBlock(&newBlock, tree)
	Block_Chain.StoreBlock(newBlock)

	BroadCastBlock(newBlock)

	mutex.Unlock()
}

func BroadCastBlock(block Block_Chain.Block) {
	blockInfo, _ := json.Marshal(&block)

	for index, node := range nodeConfig.AddressBook {
		if index == nodeConfig.NodeId {
			continue
		}
		trueUrl := node + ":" + strconv.Itoa(nodeConfig.PortForBlock+index)

		reader := bytes.NewReader(blockInfo)

		req, _ := http.NewRequest("GET", trueUrl+"/block", reader)
		req.Header.Set("Content-Type", "application/json")
		Controller.SendRequest(req)
	}
}
