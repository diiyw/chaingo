package blockchain

import (
	"fmt"
	"github.com/diiyw/chaingo/wallet"
	"testing"
)

func TestUTXOSet_FindSpendableUTXO(t *testing.T) {
	wallets := wallet.NewWallets()
	utxo := NewUTXOSet()
	defer utxo.Close()
	for _, w := range wallets.List {
		fmt.Println("address:" + w.GetAddress())
		b := utxo.FindUTXO(w.GetAddress())
		fmt.Println("balance:", b)
	}
}

func TestNewTransaction(t *testing.T) {
	wallets := wallet.NewWallets()

	for _, w := range wallets.List {
		tx := NewTransaction(w, "SrQ5LKvsFSW4Yb7h6mEpQ341JspcRqHER", 10)
		chain := OpenChain()
		fmt.Println("tx:", chain.VerifyTransaction(tx))
		chain.Close()
	}
}
