package blockchain

import (
	"testing"
	"fmt"
	"wallet"
)

func TestCreateChain(t *testing.T) {
	ws := wallet.NewWallets()
	for addr, w := range ws.Sets {
		fmt.Println("Mining to: " + w.GetAddress())
		c := CreateChain(addr)
		c.Close()
	}
}

func TestOpenChain(t *testing.T) {
	c := OpenChain()
	defer c.Close()
	iter := c.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Printf("Hash: %x %x \n", iter.Key(),iter.Value())
	}
	iter.Release()

}

func TestChain_GetBlockByHash(t *testing.T) {
	c := OpenChain()
	defer c.Close()
	block := c.GetBlock(c.prevHash)
	fmt.Printf("Nonce:  %d\n", block.Nonce)
	fmt.Printf("Heigth: %d\n", block.Height)
	fmt.Printf("BlockHash:  %x\n", block.Hash)
	fmt.Printf("Timestamp:  %s\n", block.Timestamp)
}
