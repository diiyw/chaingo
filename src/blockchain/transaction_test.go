package blockchain

import (
	"testing"
	"wallet"
	"fmt"
)

func TestUTXOSet_FindSpendableUTXO(t *testing.T) {
	wallets := wallet.NewWallets()
	utxo := NewUTXOSet()
	defer utxo.Close()
	for _, w := range wallets.Sets {
		fmt.Println("Tx:" + w.GetAddress())
		b, m := utxo.FindSpendableUTXO(w.GetPubKey(), 10)
		fmt.Println("balance:", b)
		fmt.Println("utxoes:", m)
	}
}
