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
	c = OpenChain()
	fmt.Printf("prevHash: %s", c.prevHash)
}
