package cli

import (
	"wallet"
	"fmt"
	"blockchain"
)

func ListAccount() {
	ws := wallet.NewWallets()
	for _, w := range ws.Sets {
		fmt.Println("->", w.GetAddress())
	}
}

func NewAccount() {
	wallets, _ := wallet.LoadWallets()
	w := wallets.NewWallet()
	fmt.Println("address:", w.GetAddress())
	wallets.Storage()
}

func Balance(address string) {
	fmt.Println("balance:", blockchain.NewUTXOSet().FindUTXO(address))
}
