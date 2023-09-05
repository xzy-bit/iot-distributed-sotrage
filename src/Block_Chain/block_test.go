package Block_Chain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

// Create blockChain and store in backup.json , then read and reconstruct block
func TestBlockGenerateAndStore(t *testing.T) {
	blockchain := CreateBlockChain()
	time.Sleep(time.Second)
	file, _ := os.Open("data.txt")
	reader := bufio.NewReader(file)
	for {
		currentLine, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		data := DATA{}
		json.Unmarshal(currentLine, &data)
		blockchain.AddBlock([]DATA{data})
	}

	for _, block := range blockchain.Blocks {
		fmt.Printf("Time:%s\n", block.TimeStamp.String())
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("PrevHash:%x\n", block.PrevHash)
		fmt.Printf("\n")
		//fmt.Println(block)
		StoreBlock(*block)
	}

	size := blockchain.size
	block := GetPrevBlock()
	if string(block.Hash) != string(blockchain.Blocks[size].Hash) {
		t.Errorf("Hash:%x\n", blockchain.Blocks[size].Hash)
		t.Errorf("Hash:%x\n", block.Hash)
		t.Errorf("Read and reconstruct error!\n")
	}
}

func TestDataHash(t *testing.T) {
	fd, _ := os.OpenFile("data.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)

	data1 := DATA{
		"MacBook",
		time.Now().UTC(),
		0,
		nil,
		"",
		nil,
	}
	DataHash(&data1)
	info1, _ := json.Marshal(data1)
	fd.Write(info1)
	fd.Write([]byte("\n"))
	data2 := DATA{
		"1",
		time.Now().UTC(),
		1,
		nil,
		"",
		nil,
	}
	DataHash(&data2)
	info2, _ := json.Marshal(data2)
	fd.Write(info2)
	fd.Write([]byte("\n"))
	data3 := DATA{
		"2",
		time.Now().UTC(),
		2,
		nil,
		"",
		nil,
	}
	DataHash(&data3)
	info3, _ := json.Marshal(data3)
	fd.Write(info3)
	fd.Write([]byte("\n"))
	data4 := DATA{
		"3",
		time.Now().UTC(),
		3,
		nil,
		"",
		nil,
	}
	DataHash(&data4)
	info4, _ := json.Marshal(data4)
	fd.Write(info4)
	fd.Write([]byte("\n"))
	data5 := DATA{
		"4",
		time.Now().UTC(),
		4,
		nil,
		"",
		nil,
	}
	DataHash(&data5)
	info5, _ := json.Marshal(data5)
	fd.Write(info5)
	fd.Write([]byte("\n"))
	data6 := DATA{
		"5",
		time.Now().UTC(),
		5,
		nil,
		"",
		nil,
	}
	DataHash(&data6)
	info6, _ := json.Marshal(data6)
	fd.Write(info6)
	fd.Write([]byte("\n"))
	fd.Close()
}
