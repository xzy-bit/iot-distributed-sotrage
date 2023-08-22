package Block_Chain

import (
	"math/big"
	"time"
)

type DataNode struct {
	data DATA
	next *DataNode
}

type Block struct {
	Index              int
	TimeStamp          time.Time
	Hash               []byte
	PrevHash           []byte
	BlockGenerator     int
	NextBlockGenerator int
	Data               []DATA
}

type BlockChain struct {
	Blocks []*Block
	size   int
}

type DATA struct {
	DeviceID  string
	UserId    string
	PubKey    string
	TimeStamp time.Time
	Serial    int
	Hash      []byte
	StoreOn   string
	ModNum    *big.Int
}
