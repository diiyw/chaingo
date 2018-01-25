package blockchain

import (
	"log"
	"github.com/syndtr/goleveldb/leveldb"
	"bytes"
	"errors"
	"crypto/ecdsa"
	"encoding/hex"
)

const (
	tipKey              = "tip"
	blocks              = "data/blocks"
	genesisCoinbaseData = "The Times 16/Jan/2018 Chancellor on brink of second bailout for world"
)

var chainSingleton *Chain

type Chain struct {
	*leveldb.DB
	prevHash []byte // 最靠前的区块哈希值
	height   int
}

// 创建区块链
func CreateChain(address string) *Chain {
	// 创建链数据库
	db, err := leveldb.OpenFile(blocks, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Put([]byte(tipKey), nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 挖出创世块
	block := NewGenesisBlock(address)
	err = db.Put([]byte(tipKey), block.Hash, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 加入创世块
	db.Put(block.Hash, block.Serialize(), nil)
	return &Chain{
		prevHash: block.Hash,
		DB:       db,
		height:   block.Height,
	}
}

// 打开区块链
func OpenChain() *Chain {
	if chainSingleton == nil {
		db, err := leveldb.OpenFile(blocks, nil)
		if err != nil {
			log.Fatal(err)
		}
		prevHash, err := db.Get([]byte(tipKey), nil)
		if err != nil {
			log.Fatal("ERROR: Blockchain need update.")
		}
		prevBlock, err := db.Get(prevHash, nil)
		if err != nil {
			log.Fatal("ERROR: Last block not found.")
		}
		block := DeserializeBlock(prevBlock)
		chainSingleton = &Chain{
			prevHash: prevHash,
			DB:       db,
			height:   block.Height,
		}
	}
	return chainSingleton
}

// 添加区块到链中
func (c *Chain) AppendBlock(b *Block) error {
	// 存在区块不添加
	if exits, _ := c.Has(b.Hash, nil); exits {
		return nil
	}
	// 将区块追加到链中
	blockData := b.Serialize()
	err := c.Put(b.Hash, blockData, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 最新哈希值
	tipHash, err := c.Get([]byte(tipKey), nil)
	if err != nil {
		log.Fatal(err)
	}
	// 取出最新区块的数据
	tipBlock := c.GetBlock(tipHash)
	// 需要添加区块高度比当前区块的高度大（新区块）
	if b.Height > tipBlock.Height {
		// 更新区块链的最新哈希值
		err = c.Put([]byte(tipKey), b.Hash, nil)
		if err != nil {
			log.Fatal(err)
		}
		c.prevHash = b.Hash
	}
	NewUTXOSet().Update(b)
	return nil
}

// 通过区块的哈希值获取区块
func (c *Chain) GetBlock(hash []byte) *Block {
	blockData, err := c.Get(hash, nil)
	if err != nil {
		log.Fatal(err)
	}
	return DeserializeBlock(blockData)
}

// 通过交易TxID获取交易
func (c *Chain) GetTransaction(txId []byte) (Transaction, error) {
	iter := c.NewIterator(nil, nil)
	if iter.Next() {
		k, v := iter.Key(), iter.Value()
		if bytes.Compare(k, []byte(tipKey)) != 0 {
			block := DeserializeBlock(v)
			for _, tx := range block.Transactions {
				if bytes.Compare(tx.Id, txId) == 0 {
					return *tx, nil
				}
			}
		}
	}
	iter.Release()
	return Transaction{}, errors.New("Transaction is not found ")
}

// 交易签名
func (c *Chain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		// 链中是否存在此交易
		prevTX, err := c.GetTransaction(in.TxId)
		if err != nil {
			log.Fatal(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Id)] = prevTX
	}

	tx.Sign(privateKey, prevTXs)
}

// 验证交易
func (c *Chain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := c.GetTransaction(in.TxId)
		if err != nil {
			log.Fatal(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Id)] = prevTX
	}

	return tx.Verify(prevTXs)
}

// 挖矿
func (c *Chain) MineBlock(address string, txs []*Transaction) {
	block := NewBlock(c.prevHash, c.height+1)
	block.Mining(txs)
	log.Println("INFO:", "new block", block.Height, block.Hash)
	c.AppendBlock(block)
}
