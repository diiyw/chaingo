package network

import (
	"net"
	"blockchain"
	"proof"
	"encoding/json"
	"fmt"
)

type Message interface {
	Resolve() []byte
}

// 伙伴连接
type RelNode struct {
	Ip   string
	Port int
}

func (p RelNode) Resolve() []byte {
	partner := &net.TCPAddr{
		IP:   net.ParseIP(p.Ip),
		Port: p.Port,
	}
	for _, node := range P2PNode.nodes {
		if node == partner.String() {
			return nil
		}
	}
	// 广播到其他节点
	go P2PNode.Broadcasting(p)
	P2PNode.AddNode(partner.String())
	return nil
}

// 交易
type Tx struct {
	From        string
	Transaction []byte
}

func (tx Tx) Resolve() []byte {
	P2PNode.Broadcasting(tx)
	return nil
}

// 区块链
type Blockchain struct {
	Block *blockchain.Block
}

func (b Blockchain) Resolve() []byte {
	chain := blockchain.OpenChain()
	if proof.NewPow().Validate(b.Block.Ore(), b.Block.Nonce) {
		chain.AppendBlock(b.Block)
	}
	return nil
}

func Unmarshal(b []byte) Message {
	var (
		messages = []interface{}{RelNode{}, Tx{}, Blockchain{}}
	)
	for _, m := range messages {
		fmt.Println(json.Unmarshal(b, &m))
		if json.Unmarshal(b, &m) == nil {
			return m
		}
	}

	return nil
}
