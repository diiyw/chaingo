package blockchain

import (
	"testing"
	"fmt"
	"wallet"
)

func TestCreateChain(t *testing.T) {
	w := wallet.NewWallet()
	c := CreateChain(w.GetAddress())
	c.Close()
}

func TestOpenChain(t *testing.T) {
	c := OpenChain()
	fmt.Printf("Hash: %x", c.prevHash)
}

func TestChain_GetBlockByHash(t *testing.T) {
	c := OpenChain()
	block := c.GetBlockByHash(c.prevHash)
	fmt.Printf("Nonce:  %d\n", block.Nonce)
	fmt.Printf("Heigth: %d\n", block.Height)
	fmt.Printf("Hash:   %x\n", block.Hash)
	fmt.Printf("Timestamp:  %s\n", block.Timestamp)
}
