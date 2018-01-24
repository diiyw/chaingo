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
		fmt.Println("address:" + w.GetAddress())
		b, m := utxo.FindSpendableUTXO(wallet.GetPublicKey(w.Address), 10)
		fmt.Println("balance:", b)
		fmt.Println("utxoes:", m)
	}
}

func TestNewTransaction(t *testing.T) {
	wallets := wallet.NewWallets()

	for _, w := range wallets.Sets {
		tx := NewTransaction(w, "SrQ5LKvsFSW4Yb7h6mEpQ341JspcRqHER", 10)
		chain := OpenChain()
		fmt.Println("tx:", chain.VerifyTransaction(tx))
		chain.Close()
	}
}
