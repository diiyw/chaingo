package proof

import (
	"bytes"
	"crypto/sha256"
	"github.com/diiyw/chaingo/utils"
	"log"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

// 难度系数
const targetBits = 16

type Pow struct {
	target *big.Int
}

func NewPow() Pow {
	// 初始化
	target := big.NewInt(1)
	// 左移256-24=232位
	// 0x100000000000000000000000000000000000000000000000000000000000000
	// 0x000010000000000000000000000000000000000000000000000000000000000
	target.Lsh(target, uint(256-targetBits))
	return Pow{target}
}

// 取上一个块的哈希值和收集到的交易加上时间戳、难度系数、计数器合并成一个数据
func (pow Pow) prepareData(data []byte, nonce int) []byte {
	var buf bytes.Buffer
	buf.Write(data)
	buf.Write(utils.I64Hex(int64(targetBits)))
	buf.Write(utils.I64Hex(int64(nonce)))
	return buf.Bytes()
}

// 开始挖矿
func (pow Pow) Mining(data []byte) (int, []byte) {
	var (
		hashInt big.Int
		hash    [32]byte
		nonce   = 0
	)
	log.Println("INFO:", "Mining a new block")
	for nonce < maxNonce {
		data = pow.prepareData(data, nonce)

		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

// 验证工作
func (pow Pow) Validate(data []byte, nonce int) bool {
	var hashInt big.Int

	data = pow.prepareData(data, nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
