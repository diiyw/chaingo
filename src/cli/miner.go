package cli

import "blockchain"

const TXNumbers = 1

var Txs = make(chan blockchain.Transaction, 1)

func StartMine(address string) {
	var txs []*blockchain.Transaction
	chain := blockchain.OpenChain()
	for {
		select {
		case tx := <-Txs:
			txs = append(txs, &tx)
			if len(txs) == TXNumbers {
				chain.MineBlock(address, txs)
				txs = nil
			}
		}
	}
}
