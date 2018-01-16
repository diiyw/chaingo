package blockchain

import (
	"log"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	tipKey              = "tip"
	dbDir               = "data"
	genesisCoinbaseData = "The Times 16/Jan/2018 Chancellor on brink of second bailout for world"
)

type Chain struct {
	prevHash []byte // 最靠前的区块哈希值
	db       *leveldb.DB
}

// 创建区块链
func CreateChain(address string) *Chain {
	// 创建链数据库
	db, err := leveldb.OpenFile(dbDir, nil)
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
		db:       db,
	}
}

// 打开区块链
func OpenChain() *Chain {
	db, err := leveldb.OpenFile(dbDir, nil)
	if err != nil {
		log.Fatal(err)
	}
	prevHash, err := db.Get([]byte(tipKey), nil)
	if err != nil {
		log.Fatal("Database was destroyed.")
	}
	return &Chain{
		prevHash: prevHash,
		db:       db,
	}
}

// 添加区块到链中
func (c *Chain) AppendBlock(b *Block) error {
	// 存在区块不添加
	if exits, _ := c.db.Has(b.Hash, nil); exits {
		return nil
	}
	// 将区块追加到链中
	blockData := b.Serialize()
	err := c.db.Put(b.Hash, blockData, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 最新哈希值
	tipHash, err := c.db.Get([]byte(tipKey), nil)
	if err != nil {
		log.Fatal(err)
	}
	// 取出最新区块的数据
	tipBlock := c.GetBlockByHash(tipHash)
	// 需要添加区块高度比当前区块的高度大（新区块）
	if b.Height > tipBlock.Height {
		// 更新区块链的最新哈希值
		err = c.db.Put([]byte(tipKey), b.Hash, nil)
		if err != nil {
			log.Panic(err)
		}
		c.prevHash = b.Hash
	}

	return nil
}

// 通过区块的哈希值获取区块
func (c *Chain) GetBlockByHash(hash []byte) *Block {
	blockData, err := c.db.Get(hash, nil)
	if err != nil {
		log.Fatal(err)
	}
	return DeserializeBlock(blockData)
}

// 关闭链
func (c *Chain) Close() {
	c.db.Close()
}
