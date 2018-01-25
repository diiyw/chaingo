package cli

import (
	"blockchain"
	"network"
	"wallet"
	"log"
)

const TXNumbers = 1

func StartMine(ip string, port int, address string) {
	if !wallet.ValidateAddress(address) {
		log.Fatal("ERROR: error address")
	}
	var txs []*blockchain.Transaction
	chain := blockchain.OpenChain()
	go network.NewNode(ip, port).Serving()
	for {
		select {
		case tx := <-network.Txs:
			txs = append(txs, &tx)
			if len(txs) == TXNumbers {
				chain.MineBlock(address, txs)
				txs = nil
			}
		}
	}
}
