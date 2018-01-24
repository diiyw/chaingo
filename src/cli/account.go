package cli

import (
	"wallet"
	"fmt"
)

func ListAccount() {
	ws := wallet.NewWallets()
	for idx, w := range ws.Sets {
		fmt.Println(idx, "->", w.GetAddress())
	}
}

func NewAccount() {
	wallets, _ := wallet.LoadWallets()
	w := wallets.NewWallet()
	fmt.Println("address:", w.GetAddress())
	wallets.Storage()
}
