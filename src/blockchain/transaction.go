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
	"crypto/elliptic"
	"math/big"
)

const (
	subsidy = 25
	utxo    = "data/utxo"
)

type Transaction struct {
	Id     []byte    // 交易的hash值
	Inputs []TXInput // 交易的所有输入(指明引用了哪个输出)
	TXOutputs        // 交易的所有输出
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
	chain := OpenChain()
	defer chain.Close()
	var (
		inputs  []TXInput
		outputs = TXOutputs{[]TXOutput{}}
		utxoSet = NewUTXOSet()
	)
	balance, validOutputs := utxoSet.FindSpendableUTXO(wallet.HashPubKey(w.PublicKey), amount)
	if balance < amount {
		log.Fatal("ERROR: Not enough funds")
	}
	// 使用输出作为输入
	for txId, outs := range validOutputs {
		txID, err := hex.DecodeString(txId)
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
	chain.SignTransaction(&tx, w.PrivateKey)
	return &tx
}

// 对交易进行私钥签名(要取输出的公钥一起加密)
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}
	txCopy := tx.TrimmedCopy()
	for idx, in := range tx.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.TxId)]
		txCopy.Inputs[idx].Signature = nil
		txCopy.Inputs[idx].PubKey = prevTx.Outputs[in.TxOut].PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
		if err != nil {
			log.Fatal(err)
		}
		// 验证的时候，对半取出
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[idx].Signature = signature
		txCopy.Inputs[idx].PubKey = nil
	}
}

// 复制交易（避免修改到原交易信息）
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TXInput{in.TxId, in.TxOut, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TXOutput{out.V, out.PubKeyHash})
	}

	return Transaction{
		Id:        tx.Id,
		Inputs:    inputs,
		TXOutputs: TXOutputs{outputs},
	}
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	curve := elliptic.P256()
	for idx, in := range tx.Inputs {
		if prevTx, ok := prevTXs[hex.EncodeToString(in.TxId)]; ok {
			if prevTx.Id == nil {
				return false
			}
			txCopy := tx.TrimmedCopy()
			txCopy.Inputs[idx].Signature = nil
			txCopy.Inputs[idx].PubKey = prevTx.Outputs[in.TxOut].PubKeyHash

			r := big.Int{}
			s := big.Int{}
			sigLen := len(in.Signature)
			r.SetBytes(in.Signature[:(sigLen / 2)])
			s.SetBytes(in.Signature[(sigLen / 2):])

			x := big.Int{}
			y := big.Int{}
			keyLen := len(in.PubKey)
			x.SetBytes(in.PubKey[:(keyLen / 2)])
			y.SetBytes(in.PubKey[(keyLen / 2):])

			dataToVerify := fmt.Sprintf("%x\n", txCopy)

			rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
			if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
				return false
			}
			txCopy.Inputs[idx].PubKey = nil
		}
	}

	return true
}

// 创建coinbase交易
func NewCoinbaseTx(to, data string) *Transaction {
	utxoSet := NewUTXOSet()
	defer utxoSet.Close()
	tx := &Transaction{
		Id:     []byte{},
		Inputs: []TXInput{{[]byte{}, -1, nil, []byte(data)}},
		TXOutputs: TXOutputs{
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
	TxOut     int    // 输出的索引
	Signature []byte // 签名
	PubKey    []byte // 公钥
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

// 找到可用的UTXO
func (u UTXOSet) FindSpendableUTXO(pubKeyHash []byte, amount int) (int, map[string][]int) {
	var (
		unspentOutputs = make(map[string][]int)
		sum            = 0
	)
	iter := u.NewIterator(nil, nil)
	for iter.Next() {
		k, v := iter.Key(), iter.Value()
		txID := hex.EncodeToString(k)
		outputs := DeserializeOutputs(v)
		// 统计未消费掉的输出
		for outIdx, out := range outputs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && sum < amount {
				sum += out.V
				// outIdx是未花费掉的输出索引
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
			}
		}
	}
	iter.Release()
	return sum, unspentOutputs
}

// 获取余额
func (u UTXOSet) FindUTXO(address string) int {
	var (
		balance    = 0
		pubKeyHash = wallet.GetPublicKey([]byte(address))
	)
	iter := u.NewIterator(nil, nil)
	for iter.Next() {
		_, v := iter.Key(), iter.Value()
		outputs := DeserializeOutputs(v)
		// 统计未消费掉的输出
		for _, out := range outputs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				balance += out.V
			}
		}
	}
	iter.Release()
	return balance
}

// 添加花费的输出（没有就创建）
func (u UTXOSet) AddTXOutputs(tx *Transaction) {
	u.Put(tx.Id, tx.TXOutputs.Serialize(), nil)
}
