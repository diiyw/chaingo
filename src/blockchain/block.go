package blockchain

// 区块
type Block struct {
	PrevHash     []byte         // 上一个区块的哈希值
	Timestamp    int64          // 区块创建时间
	Hash         []byte         // 区块哈希
	Nonce        int            // 当前计数器（挖矿使用）
	Height       int            // 当前区块高度
	Transactions []*Transaction // 区块中的交易记录
}
