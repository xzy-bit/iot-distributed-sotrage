package Node

import (
	"IOT_Storage/src/Block_Chain"
	"bytes"
	"io"
	"math/big"
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
