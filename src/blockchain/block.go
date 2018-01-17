package blockchain

import (
	"encoding/gob"
	"log"
	"bytes"
	"time"
	"core"
	"proof"
)

// 区块
type Block struct {
	PrevHash     []byte         // 上一个区块的哈希值
	Timestamp    time.Time      // 区块创建时间
	Hash         []byte         // 区块哈希
	Nonce        int            // 当前计数器（挖矿使用）
	Height       int            // 当前区块高度
	Transactions []*Transaction // 区块中的交易记录
}

// 创建区块
func NewBlock(prevHash []byte, height int) *Block {
	return &Block{
		PrevHash:  prevHash,
		Timestamp: time.Now(),
		Height:    height,
	}
}

// 创建创世块
func NewGenesisBlock(address string) *Block {
	block := NewBlock(nil, 0)
	block.mining([]*Transaction{NewCoinbaseTx(address, genesisCoinbaseData)})
	return block
}

func (b *Block) mining(txs []*Transaction) {
	b.Transactions = txs
	pow := proof.NewPow()
	nonce, hash := pow.Mining(b.Ore())
	b.Hash = hash
	b.Nonce = nonce
}

// 提供给挖矿的数据
func (b *Block) Ore() []byte {
	var (
		payload bytes.Buffer
		txBytes [][]byte
	)
	payload.Write(b.PrevHash)
	payload.Write(core.I64Hex(b.Timestamp.Unix()))
	for _, tx := range b.Transactions {
		txBytes = append(txBytes, tx.Serialize())
	}
	mTree := NewMerkle(txBytes)
	payload.Write(mTree.RootNode.Hash)
	return payload.Bytes()
}

// 编码区块
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 解码区块
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Fatal(err)
	}

	return &block
}
