package Block_Chain

import (
	"time"
)

type DataNode struct {
	Data DATA
	Next *DataNode
}

type Block struct {
	Index          int
	TimeStamp      time.Time
	Hash           []byte
	PrevHash       []byte
	BlockGenerator int
	Data           []DATA
}

type BlockChain struct {
	Blocks []*Block
	size   int
}

type DATA struct {
	DeviceID     string
	TimeStamp    time.Time
	Serial       int
	IndexOfGroup int
	Hash         []byte
	StoreOn      string
	NumOfGroup   int
}
