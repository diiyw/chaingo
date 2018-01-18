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
	"crypto/ecdsa"
	"fmt"
	"crypto/rand"
)

const (
	subsidy = 25
	utxo    = "data/utxo"
)

type Transaction struct {
	Id      []byte    // 交易的hash值
	Inputs  []TXInput // 交易的所有输入(指明引用了哪个输出)
	Outputs TXOutputs // 交易的所有输出
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.Id = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

//  检查是否是coinbase交易
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxId) == 0 && tx.Inputs[0].TxOut == -1
}

// 创建交易
func NewTransaction(w *wallet.Wallet, to string, amount int) *Transaction {
	var (
		inputs  []TXInput
		outputs = TXOutputs{[]TXOutput{}}
		utxoSet = NewUTXOSet()
	)
	balance, validOutputs := utxoSet.FindSpendableUTXO(w.PublicKey, amount)
	if balance < amount {
		log.Fatal("ERROR: Not enough funds")
	}
	// 使用输出作为输入
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Fatal(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}
	// 接受收者的输入
	outputs.Outputs = append(outputs.Outputs, *NewTXOutput(amount, to))
	// 余额大于发送的币，返还剩余币给发送者
	if balance > amount {
		outputs.Outputs = append(outputs.Outputs, *NewTXOutput(balance-amount, w.GetAddress()))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.Id = tx.Hash()
	// 交易签名
	tx.Sign(w.PrivateKey)
	return &tx
}

// 对交易进行私钥签名(证明该交易是合法发起的)
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey) {
	chain := OpenChain()
	if tx.IsCoinbase() {
		return
	}
	for idx, in := range tx.Inputs {
		// 链中是否存在此交易
		_, err := chain.GetTransaction(in.TxId)
		if err != nil {
			log.Fatal(err)
		}
		// 私钥签名，确认是拥有者发送的
		dataToSign := fmt.Sprintf("%x\n", *tx)

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
		if err != nil {
			log.Fatal(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[idx].Signature = signature
	}
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
	V          int    // 输出多少币
	PubKeyHash []byte // 输出者的公钥
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

// 所有输出
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

// 未花掉的输出（相对接受者就是输入）
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
		// 统计未消费掉的输出
		for outIdx, out := range outputs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && balance < amount {
				balance += out.V
				// outIdx是未花费掉的输出索引
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
			}
		}
	}
	iter.Release()
	return balance, unspentOutputs
}

// 添加花费的输出（没有就创建）
func (u UTXOSet) AddTXOutputs(tx *Transaction) {
	u.Put(tx.Id, tx.Outputs.Serialize(), nil)
}
