package wallet

import (
	"testing"
	"fmt"
	"core"
)

func TestWallet(t *testing.T) {
	wallet := NewWallet()
	addr := wallet.GetAddress()
	pubKeyHash := core.Base58Decode(wallet.Address)
	fmt.Println(addr)
	fmt.Println("addr:", pubKeyHash)
	fmt.Println("version:", pubKeyHash[0])
	fmt.Println("checksum:", pubKeyHash[len(pubKeyHash)-AddressChecksumLen:])
	fmt.Println("public key:", pubKeyHash[1: len(pubKeyHash)-AddressChecksumLen])
	fmt.Println("ok:", ValidateAddress(addr))
}
