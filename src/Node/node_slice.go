package Node

import (
	"IOT_Storage/src/Block_Chain"
	"IOT_Storage/src/Controller"
	"IOT_Storage/src/File_Index"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// save the sclie and generate DATA struct
func GenerateDATA(iotId string, serial string, address string, modNum string, timeStamp string, hash string, indexOfGroup int) *Block_Chain.DATA {
	num, _ := strconv.Atoi(serial)
	mod := new(big.Int)
	mod.SetString(modNum, 10)
	//log.Println("modNum:" + mod.String())
	stamp, _ := time.Parse("2006-01-02 15:04:05", timeStamp)
	h, _ := hex.DecodeString(hash)
	dataIndex := Block_Chain.DATA{
		DeviceID:     iotId,
		TimeStamp:    stamp,
		Serial:       num,
		Hash:         h,
		StoreOn:      address,
		ModNum:       mod,
		IndexOfGroup: indexOfGroup,
	}
	return &dataIndex
}

func AddDataToCache(newData *Block_Chain.DATA) {
	dataNode := Block_Chain.DataNode{
		Data: *newData,
		Next: nil,
	}
	if Head == nil {
		Head = &dataNode
		Tail = Head.Next
		return
	}
	Tail = &dataNode
	Tail = Tail.Next
	return
}

func SaveSlice(cipher string, fileName string) {
	os.MkdirAll(filepath.Dir(fileName), 0750)
	out, _ := os.Create(fileName)
	defer out.Close()
	reader := bytes.NewReader([]byte(cipher))
	io.Copy(out, reader)
}

func SaveJson(info []byte, fileName string) {
	os.MkdirAll(filepath.Dir(fileName), 0750)
	out, _ := os.Create(fileName)
	defer out.Close()
	reader := bytes.NewReader(info)
	io.Copy(out, reader)
}

func GetAllDataInCache() []Block_Chain.DATA {
	data := []Block_Chain.DATA{}

	if Head == nil {
		return nil
	}

	for Head != Tail {
		data = append(data, Head.Data)
		Head = Head.Next
	}

	return data
}

func HandleData(data []Block_Chain.DATA, generator int) {
	mutex.Lock()

	oldBlock := Block_Chain.GetPrevBlock()
	//log.Println(oldBlock)

	newBlock := Block_Chain.GenerateBlock(*oldBlock, data)
	newBlock.BlockGenerator = generator
	//log.Println("block:")
	//log.Println(newBlock)

	File_Index.InsertBlock(&newBlock, tree)

	//log.Println("tree:")
	//log.Println(tree)
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
