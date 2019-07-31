package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"golang.org/x/crypto/sha3"
)

const (
	TRX_MESSAGE_HEADER = "\x19TRON Signed Message:\n32"
	AddressPrefixMain  = "41" //41 + address
)

// VerifySign 验证签名
//	transaction: 交易对象
func VerifySignature(sign []byte, rawData string, addr string) bool {
	if len(sign) != 65 { // sign check
		return false
	}
	if sign[64] >= 27 {
		sign[64] = sign[64] - 27
	}
	rawData = TRX_MESSAGE_HEADER + rawData

	pubKey, err := GetSignedPubKey(rawData, sign)
	if err != nil {
		return false
	}
	signAddr, err := GetTronBase58Address(HexEncode(pubKey))
	if err != nil {
		return false
	}

	if signAddr != addr {
		return false
	}
	return true
}

// GetSignedPublicKey 获取交易签名账户的公钥
func GetSignedPubKey(rawData string, sign []byte) ([]byte, error) {
	if len(sign) != 65 { // sign check
		return nil, errors.New("Invalid Transaction Signature, should be 65 length bytes")
	}
	rawByte := []byte(rawData)
	hash := ethcrypto.Keccak256(rawByte)

	return ethcrypto.Ecrecover(hash[:], sign)
}

// GetTronBase58Address 根据 hex encoding public key生成 hex encoding address
//	in: hex encoding public key (uncompressed public key)
//	out: base58 encoding address
func GetTronBase58Address(in string) (out string, err error) {
	hexAddr, err := GetTronHexAddress(in)
	if nil != err {
		return "", err
	}

	out = Base58EncodeAddr(HexDecode(hexAddr))

	return
}

// GetTronHexAddress 根据 hex encoding public key生成 hex encoding address
//	in: hex encoding public key (uncompressed public key)
//	out: hex encoding address
func GetTronHexAddress(in string) (out string, err error) {
	pubBytes, err := hex.DecodeString(in)
	if nil != err {
		return "", err
	}
	if 1 > len(pubBytes) {
		return "", fmt.Errorf("Invalid address")
	}
	rawPubKey := pubBytes[1:] // remove prefix byte

	sha3Hash := sha3.NewLegacyKeccak256() // use sha3 keccad256
	sha3Hash.Write(rawPubKey)
	hashRet := sha3Hash.Sum(nil)

	hashRetStr := HexEncode(hashRet) // covert to hex string
	addrPrefix := AddressPrefixMain

	out = fmt.Sprintf("%s%s", addrPrefix, hashRetStr[24:]) // address prefix + hash remove first 24 length

	return
}

// Base58EncodeAddr 将地址字节码编码为base58字符串
func Base58EncodeAddr(in []byte) string {
	if len(in) < 2 {
		return ""
	}
	return base58.CheckEncode(in[1:], in[0]) // first byte is version, reset is data
}

// HexDecode ...
func HexDecode(in string) []byte {
	ret, _ := hex.DecodeString(in)
	return ret
}

// HexEncode ...
func HexEncode(in []byte) string {
	return hex.EncodeToString(in)
}
