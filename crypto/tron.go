package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	"golang.org/x/crypto/sha3"
)

const (
	TrxMessageHeader  = "\x19TRON Signed Message:\n32"
	AddressPrefixMain = "41"
)

// TronlinkSignature
//  signType: if add TrxMessageHeader to data
//  data: data to be signed
//  hexPrivateKey: sign private key
func Signature(signType bool, data string, hexPrivateKey string) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid raw data")
	}

	if signType {
		data = TrxMessageHeader + data
	}

	hash0 := ethCrypto.Keccak256([]byte(data))

	privateKey, err := ethCrypto.HexToECDSA(hexPrivateKey)
	if nil != err || nil == privateKey {
		return nil, err
	}

	signData, err := ethCrypto.Sign(hash0[:], privateKey)
	if err == nil && len(signData) == 65 && signData[64] < 27 {
		signData[64] = signData[64] + 27
	}
	return signData, err
}

// VerifyTronlinkSignature
//  signType: if add TrxMessageHeader to data
//  sign: signature obtained by signing rawData
//  rawData: Raw data to be signed
//  addr: Tron address, prefixed with "T"
func VerifySignature(signType bool, sign []byte, rawData string, addr string) bool {
	if len(sign) != 65 { // sign check
		return false
	}

	if sign[64] >= 27 {
		sign[64] = sign[64] - 27
	}

	if signType {
		rawData = TrxMessageHeader + rawData
	}

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

// GetSignedPublicKey get the public key of the transaction signing account
//  rawData: Raw data to be signed
//  sign: signature obtained by signing rawData
func GetSignedPubKey(rawData string, sign []byte) ([]byte, error) {
	if len(sign) != 65 { // sign check
		return nil, errors.New("invalid transaction signature, should be 65 length bytes")
	}
	rawByte := []byte(rawData)
	hash := ethCrypto.Keccak256(rawByte)

	return ethCrypto.Ecrecover(hash[:], sign)
}

// GetTronBase58Address Generate hex encoding address according to hex encoding public key
//	in: hex encoding public key (uncompressed public key)
//	out: base58 encoding address
func GetTronBase58Address(in string) (out string, err error) {
	hexAddr, err := GetTronHexAddress(in)
	if nil != err {
		return "", err
	}

	bytes, err := HexDecode(hexAddr)
	if err != nil {
		return "", err
	}
	out = Base58EncodeAddr(bytes)
	return
}

// GetTronHexAddress Generate hex encoding address according to hex encoding public key
//	in: hex encoding public key (uncompressed public key)
//	out: hex encoding address
func GetTronHexAddress(in string) (out string, err error) {
	pubBytes, err := hex.DecodeString(in)
	if nil != err {
		return "", err
	}
	if 1 > len(pubBytes) {
		return "", fmt.Errorf("invalid address")
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

// Base58EncodeAddr Encode byte address to base58 string
//  in: byte array address
//  out: Base58 string address
func Base58EncodeAddr(in []byte) string {
	if len(in) < 2 {
		return ""
	}
	return base58.CheckEncode(in[1:], in[0]) // first byte is version, reset is data
}

// Base58DecodeAddr Decode base58 string to byte address
//  in: Base58 string address
//  out: byte array address
func Base58DecodeAddr(in string) ([]byte, error) {
	decodeCheck := base58.Decode(in)
	if len(decodeCheck) <= 4 {
		return nil, errors.New("base58 decode length error")
	}
	decodeData := decodeCheck[:len(decodeCheck)-4]
	hash0, err := Hash(decodeData)
	if err != nil {
		return nil, err
	}
	hash1, err := Hash(hash0)
	if hash1 == nil {
		return nil, err
	}
	if hash1[0] == decodeCheck[len(decodeData)] && hash1[1] == decodeCheck[len(decodeData)+1] &&
		hash1[2] == decodeCheck[len(decodeData)+2] && hash1[3] == decodeCheck[len(decodeData)+3] {
		return decodeData, nil
	}
	return nil, errors.New("base58 check failed")
}

// HexDecode ...
func HexDecode(in string) ([]byte, error) {
	return hex.DecodeString(in)
}

// HexEncode ...
func HexEncode(in []byte) string {
	return hex.EncodeToString(in)
}

// Package goLang sha256 hash algorithm.
func Hash(s []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(s)
	if err != nil {
		return nil, err
	}
	bs := h.Sum(nil)
	return bs, nil
}
