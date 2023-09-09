package File_Index

import (
	"IOT_Storage/src/Block_Chain"
	"bufio"
	"encoding/json"
	"github.com/emirpasic/gods/trees/avltree"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type TreeKey struct {
	DeviceId  string
	TimeStamp time.Time
	Serial    int
}

// implement a comparator
func ComparatorForTreeKey(a, b interface{}) int {
	v1 := a.(TreeKey)
	v2 := b.(TreeKey)
	s1 := strings.Compare(v1.DeviceId, v2.DeviceId)
	if s1 == 0 {
		// compare timestamp
		if v1.TimeStamp.Before(v2.TimeStamp) {
			return -1
		}
		if v1.TimeStamp.After(v2.TimeStamp) {
			return 1
		}

		// compare serial
		if v1.Serial < v2.Serial {
			return -1
		}
		if v1.Serial > v2.Serial {
			return 1
		}
		return 0
	}
	return s1
}

func GetNextBlock(reader *bufio.Reader) *Block_Chain.Block {

	currentLine, err := reader.ReadBytes('\n')
	if err == io.EOF {
		return nil
	}

	block := Block_Chain.Block{}
	err = json.Unmarshal(currentLine, &block)
	if err != nil {
		log.Fatal("Json to block error！\n")
	}
	//fmt.Println(block)
	return &block
}
func BuildTraverser(backupFilePath string) *avltree.Tree {
	tree := avltree.NewWith(ComparatorForTreeKey)
	var err error
	file, err := os.OpenFile(backupFilePath, os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		log.Print(err)
		return nil
	}
	reader := bufio.NewReader(file)
	for {
		block := GetNextBlock(reader)
		if block == nil {
			break
		}
		InsertBlock(block, tree)
	}
	return tree
}

func InsertBlock(block *Block_Chain.Block, tree *avltree.Tree) {
	//fmt.Println(block.Data)
	//tree := avltree.NewWithIntComparator()
	for _, data := range block.Data {
		if data.Hash == nil {
			break
		}
		tree.Put(TreeKey{DeviceId: data.DeviceID, TimeStamp: data.TimeStamp, Serial: data.Serial}, data)
	}
}
