package blockchain

import (
	"core"
	"wallet"
	"bytes"
	"encoding/gob"
	"log"
)

const subsidy = 25

type Transaction struct {
	Id  []byte     // 交易的hash值
	In  []TXInput  // 交易的所有收入
	Out []TXOutput // 交易的所有支出
}

// 创建coinbase交易
func NewCoinbaseTx(to, data string) *Transaction {
	return &Transaction{
		Id:  []byte{},
		In:  []TXInput{{[]byte{}, -1, nil, []byte(data)}},
		Out: []TXOutput{*NewTXOutput(subsidy, to)},
	}
}

// 编码交易
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

type TXInput struct {
	TxId      []byte // 交易的hash值
	TxOut     int
	Signature []byte
	PubKey    []byte
}

type TXOutput struct {
	V          int    // 支出多少币
	PubKeyHash []byte // 支出者的公钥
}

// 输出签名（锁定）
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := core.Base58Encode(address)
	// 地址的0位是版本号去掉，最后4位是校验位也去掉
	pubKeyHash = pubKeyHash[1: len(pubKeyHash)-wallet.AddressChecksumLen]
	out.PubKeyHash = pubKeyHash
}

// 检测输出是否是某签名签名过的
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// 创建一个新的输出
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}
