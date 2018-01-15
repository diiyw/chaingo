package wallet

import (
	"crypto/sha256"
	"log"
	"golang.org/x/crypto/ripemd160"
	"crypto/ecdsa"
	"bytes"
	"crypto/elliptic"
	"core"
	"crypto/rand"
)

const (
	version            = 0x01
	AddressChecksumLen = 4
)

// 钱包，存储公钥和私钥
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
	Address    []byte
}

// 新建钱包
func NewWallet() *Wallet {
	// 生成私钥和公钥
	private, public := newKeyPair()
	w := Wallet{
		PrivateKey: private,
		PublicKey:  public,
	}
	// 对公钥哈希化
	pubKeyHash := HashPubKey(w.PublicKey)
	// 再加上版本号，完整的钱包地址
	versionedPayload := append([]byte{version}, pubKeyHash...)
	// 生成一个校验码
	checksum := checksum(versionedPayload)
	// 校验码追加到公钥的尾部
	fullPayload := append(versionedPayload, checksum...)
	// 以可读的形式base58编码
	w.Address = core.Base58Encode(fullPayload)
	return &w
}

// 获取钱包地址 base58(版本号+公钥+校验码)
func (w Wallet) GetAddress() string {
	return string(w.Address)
}

// 获取公钥哈希值
func HashPubKey(pubKey []byte) []byte {
	// sha256算法加密公钥
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	// 通过RIPEMD160算法得出的公钥
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// 验证钱包地址是否可用
func ValidateAddress(address string) bool {
	pubKeyHash := core.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-AddressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1: len(pubKeyHash)-AddressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// 生成公钥的校验
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:AddressChecksumLen]
}

// 生成公钥和私钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
