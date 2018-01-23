package network

import (
	"net"
	"blockchain"
	"proof"
	"encoding/json"
)

type Message interface {
	Resolve() []byte
}

// 广播节点
type RelNode struct {
	Ip   string
	Port int
}

func (r RelNode) Resolve() []byte {
	node := &net.TCPAddr{
		IP:   net.ParseIP(r.Ip),
		Port: r.Port,
	}
	for _, n := range P2PNode.nodes {
		if n == node.String() {
			return nil
		}
	}
	// 广播到其他节点
	go P2PNode.Broadcasting(r, node.String())
	P2PNode.AddNode(node.String())
	return nil
}

// 广播交易
type Tx struct {
	From        string
	Transaction []byte
}

func (tx Tx) Resolve() []byte {
	P2PNode.Broadcasting(tx, tx.From)
	return nil
}

// 广播区块
type Block struct {
	Block *blockchain.Block
}

func (b Block) Resolve() []byte {
	chain := blockchain.OpenChain()
	if proof.NewPow().Validate(b.Block.Ore(), b.Block.Nonce) {
		chain.AppendBlock(b.Block)
	}
	return nil
}

func Unmarshal(b []byte) Message {
	var (
		messages = []Message{RelNode{}, Tx{}, Block{}}
	)
	for _, m := range messages {
		switch mType := m.(type) {
		case RelNode:
			if json.Unmarshal(b, &mType) == nil {
				return mType
			}
		case Tx:
			if json.Unmarshal(b, &mType) == nil {
				return mType
			}
		case Block:
			if json.Unmarshal(b, &mType) == nil {
				return mType
			}
		}
	}

	return nil
}
