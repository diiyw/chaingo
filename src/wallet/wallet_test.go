package wallet

import (
	"testing"
	"fmt"
	"core"
)

func TestWallet(t *testing.T) {
	wallets := NewWallets()
	for addr, wallet := range wallets.Sets {
		pubKeyHash := core.Base58Decode(wallet.Address)
		fmt.Println(addr)
		fmt.Println("addr:", pubKeyHash)
		fmt.Println("version:", pubKeyHash[0])
		fmt.Println("checksum:", pubKeyHash[len(pubKeyHash)-AddressChecksumLen:])
		fmt.Println("public key:", pubKeyHash[1: len(pubKeyHash)-AddressChecksumLen])
		fmt.Println("ok:", ValidateAddress(addr))
	}
}
