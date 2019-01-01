package blockchain

import "crypto/sha256"

// 默克林树形结构
type Merkle struct {
	RootNode *MerkleNode
}

// 默克林树的节点
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte //节点的哈希值
}

// 从字节数组中创建默克林树
func NewMerkle(data [][]byte) *Merkle {
	var nodes []*MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	// 创建每个默克林节点
	for _, datum := range data {
		nodes = append(nodes, NewMerkleNode(nil, nil, datum))
	}

	// 上面的节点创建后还不是树形，这里把他成构造成一棵树(两个分支)
	for i := 0; i < len(data)/2; i++ {
		var newLevel []*MerkleNode
		// 每两个一组向上构造一棵树，最终只会剩余一个root节点
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(nodes[j], nodes[j+1], nil)
			newLevel = append(newLevel, node)
		}
		// 覆盖原先nodes继续两个一组向上生成一棵树
		nodes = newLevel
	}

	mTree := &Merkle{nodes[0]}

	return mTree
}

// 创建默克林节点
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{}

	if left == nil && right == nil {
		// 左右叶节点为空，即创建新节点，算出sha256哈希值
		hash := sha256.Sum256(data)
		node.Hash = hash[:]
	} else {
		// 不是创建新节点，左右子叶合并后计算哈希值
		prevHashes := append(left.Hash, right.Hash...)
		hash := sha256.Sum256(prevHashes)
		node.Hash = hash[:]
	}

	node.Left = left
	node.Right = right

	return node
}
