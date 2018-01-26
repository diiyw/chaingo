package network

import (
	"net"
	"blockchain"
	"proof"
	"encoding/json"
	"log"
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

var Txs = make(chan blockchain.Transaction, 1)

func (tx Tx) Resolve() []byte {
	log.Println("INFO: Received transaction")
	Txs <- blockchain.DeserializeTransaction(tx.Transaction)
	P2PNode.Broadcasting(tx, tx.From)
	return nil
}

// 广播区块
type Block struct {
	Block *blockchain.Block
}

func (b Block) Resolve() []byte {
	log.Println("INFO:", "new block", b.Block.Height)
	chain := blockchain.OpenChain()
	if proof.NewPow().Validate(b.Block.Ore(), b.Block.Nonce) {
		if err := chain.AppendBlock(b.Block); err != nil {
			log.Println("INFO: Append block error ->", err)
		}
	}
	return nil
}

func Unmarshal(b []byte) Message {
	switch  b[9] {
	case 'T':
		tx := Tx{}
		if json.Unmarshal(b[11:], &tx) == nil {
			return tx
		}
	case 'B':
		block := Block{}
		if json.Unmarshal(b[14:], &block) == nil {
			return block
		}
	case 'R':
		node := RelNode{}
		if json.Unmarshal(b[16:], &node) == nil {
			return node
		}
	}
	return nil
}
