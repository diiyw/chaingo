package cli

import (
	"fmt"
	"github.com/diiyw/chaingo/blockchain"
	"github.com/diiyw/chaingo/network"
	"github.com/diiyw/chaingo/wallet"
	"log"
)

func ListAccount() {
	ws := wallet.NewWallets()
	for _, w := range ws.List {
		fmt.Println("->", w.GetAddress())
	}
}

func NewAccount() error {
	wallets, err := wallet.LoadWallets()
	if err != nil {
		return err
	}
	w := wallets.NewWallet()
	fmt.Println("address:", w.GetAddress())
	wallets.Storage()
	return nil
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
	for addr, w := range ws.List {
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
