package Block_Chain

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func GetPrevBlock() Block {
	file, err := os.Open("backup.json")
	if err != nil {
		log.Fatal("Open backup error!\n")
	}

	reader := bufio.NewReader(file)
	var lastLine []byte
	for {
		currentLine, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		lastLine = currentLine
	}
	block := Block{}
	err = json.Unmarshal(lastLine, &block)
	if err != nil {
		log.Fatal("Json to block error！\n")
	}
	//fmt.Println(block)
	return block
}

func StoreBlock(newBlock Block) {
	fd, _ := os.OpenFile("backup.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	blockInfo, _ := json.Marshal(newBlock)

	//block := Block{}
	//json.Unmarshal(blockInfo, &block)
	//fmt.Println(block)

	fd.Write(blockInfo)
	fd.Write([]byte("\n"))
	fd.Close()
	return
}

func GenesisBlock() *Block {
	block := Block{0, time.Now().UTC(), nil, nil, -1, 1, nil}
	return &block
}

func SetHash(block *Block) {
	info := bytes.Join([][]byte{
		[]byte(block.TimeStamp.String()),
		block.PrevHash,
	}, []byte{})
	hash := sha256.Sum256(info)
	block.Hash = hash[:]
	return
}

func DataHash(data *DATA) {
	info := bytes.Join([][]byte{
		[]byte(data.DeviceID),
		[]byte(data.TimeStamp.String()),
		[]byte(data.StoreOn),
		data.ModNum.Bytes(),
	}, []byte{})
	hash := sha256.Sum256(info)
	data.Hash = hash[:]
	return
}

func GenerateBlock(oldBlock Block, Data []DATA) Block {
	var newBlock Block
	t := time.Now().UTC()

	newBlock.Index = oldBlock.Index + 1
	newBlock.TimeStamp = t
	newBlock.PrevHash = oldBlock.Hash
	if oldBlock.Hash == nil {
		fmt.Println("Genius block!")
	}
	newBlock.Data = Data
	SetHash(&newBlock)
	return newBlock
}

func CreateBlockChain() *BlockChain {
	blockchain := BlockChain{}
	blockchain.Blocks = append(blockchain.Blocks, GenesisBlock())
	blockchain.size = 0
	return &blockchain
}

func (bc *BlockChain) AddBlock(Data []DATA) {
	newBlock := GenerateBlock(*bc.Blocks[bc.size], Data)
	bc.Blocks = append(bc.Blocks, &newBlock)
	bc.size++
}
