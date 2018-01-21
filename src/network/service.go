package network

import (
	"net"
	"blockchain"
)

type Message interface {
	Resolve() []byte
}

// 伙伴连接
type Partner struct {
	Ip   string
	Port int
}

func (p Partner) Resolve() []byte {
	partner := &net.TCPAddr{
		IP:   net.ParseIP(p.Ip),
		Port: p.Port,
	}
	for _, node := range P2PNode.nodes {
		if node.addr.String() == partner.String() {
			return nil
		}
	}
	P2PNode.AddPartner(partner)
	// 广播到其他节点
	P2PNode.Broadcasting(p)
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
	chain.AppendBlock(b.Block)
	return nil
}
