package cli

import (
	"wallet"
	"fmt"
	"blockchain"
	"log"
	"network"
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

func Fund(address string) {
	if wallet.ValidateAddress(address) {
		fmt.Println("balance:", blockchain.NewUTXOSet().FindUTXO(address))
		return
	}
	fmt.Println("ERROR:", "address error")
}

func Send(amount int, from, to string) {
	ws := wallet.NewWallets()
	for addr, w := range ws.Sets {
		if addr == from {
			tx := blockchain.NewTransaction(w, to, amount)
			network.NewNode("127.0.0.1", 2048).Broadcasting(&network.Tx{
				From:        network.P2PNode.String(),
				Transaction: tx.Serialize(),
			}, network.P2PNode.String())
			return
		}
	}
	log.Println("ERROR:", "no permit to send funds")
}
