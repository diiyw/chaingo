package cli

import (
	"blockchain"
	"fmt"
)

func PrintChain() {
	chain := blockchain.OpenChain()
	defer chain.Close()
	iter := chain.NewIterator(nil, nil)
	for iter.Next() {
		k, v := iter.Key(), iter.Value()
		if string(k) != "tip" {
			block := blockchain.DeserializeBlock(v)
			fmt.Printf("============ Block %x ============\n", block.Hash)
			fmt.Printf("Height: %d\n", block.Height)
			fmt.Printf("Prev. block: %x\n", block.PrevHash)
			for _, tx := range block.Transactions {
				fmt.Println(tx)
			}
			fmt.Printf("\n\n")
		}
	}
	iter.Release()
}
