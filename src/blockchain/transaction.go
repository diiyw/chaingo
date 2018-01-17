package blockchain

import (
	"core"
	"wallet"
	"bytes"
	"encoding/gob"
	"log"
	"github.com/syndtr/goleveldb/leveldb"
	"encoding/hex"
	"crypto/sha256"
)

const (
	subsidy = 25
	utxo    = "data/utxo"
)

type Transaction struct {
	Id      []byte    // 交易的hash值
	Inputs  []TXInput // 交易的所有收入(指明引用了哪个支出)
	Outputs TXOutputs // 交易的所有支出
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.Id = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// 创建coinbase交易
func NewCoinbaseTx(to, data string) *Transaction {
	utxoSet := NewUTXOSet()
	defer utxoSet.Close()
	tx := &Transaction{
		Id:     []byte{},
		Inputs: []TXInput{{[]byte{}, -1, nil, []byte(data)}},
		Outputs: TXOutputs{
			[]TXOutput{
				*NewTXOutput(subsidy, to),
			},
		},
	}
	tx.Id = tx.Hash()
	utxoSet.AddTXOutputs(tx)
	return tx
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
	pubKeyHash := core.Base58Decode(address)
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

// 所有支出
type TXOutputs struct {
	Outputs []TXOutput
}

// 编码输出
func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Fatal(err)
	}

	return buff.Bytes()
}

// 解码输出
func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Fatal(err)
	}

	return outputs
}

// 未花掉的支出（相对接受者就是收入）
type UTXOSet struct {
	*leveldb.DB
}

// 创建UTXO集合
func NewUTXOSet() UTXOSet {
	db, err := leveldb.OpenFile(utxo, nil)
	if err != nil {
		log.Fatal(err)
	}
	return UTXOSet{
		DB: db,
	}
}

// 找到所有可用的UTXO（获取余额）
func (u UTXOSet) FindSpendableUTXO(pubKeyHash []byte, amount int) (int, map[string][]int) {
	var (
		unspentOutputs = make(map[string][]int)
		balance        = 0
	)
	iter := u.NewIterator(nil, nil)
	for iter.Next() {
		k, v := iter.Key(), iter.Value()
		txID := hex.EncodeToString(k)
		outputs := DeserializeOutputs(v)
		// 统计未消费掉的支出
		for outIdx, out := range outputs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && balance < amount {
				balance += out.V
				// outIdx是未花费掉的支出索引
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
			}
		}
	}
	iter.Release()
	return balance, unspentOutputs
}

// 添加花费的支出（没有就创建）
func (u UTXOSet) AddTXOutputs(tx *Transaction) {
	u.Put(tx.Id, tx.Outputs.Serialize(), nil)
}
