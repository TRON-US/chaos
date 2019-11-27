package crypto

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	privateKey         = "ec4c14e1e6c3c163d4a7edb253b2cac3f8e13369f36eec107dcac6452bf74086"
	publicKeyBash58    = "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh"
	publicKeyHexString = "41CF5D9F0BCBD34852AE1AA1F7AD1A73532E28427A"
	testRawData        = "test"
	testSign           = "e5e8df08423f501f5a7db498258a341c2faff6f86133ffe091d76129135d274a1a7cb4e63157d361d9e0012186a78acf9352d0d1af745589bdd260aeab35bbaf1b"
)

func TestSignature(t *testing.T) {
	sign, err := Signature(true, testRawData, privateKey)
	assert.NoError(t, err)

	assert.EqualValues(t, hex.EncodeToString(sign), testSign)
}

func TestVerifySignature(t *testing.T) {
	sign, err := hex.DecodeString(testSign)
	assert.NoError(t, err)

	assert.True(t, true, VerifySignature(true, sign, testRawData, publicKeyBash58))
}

func TestBase58EncodeAddr(t *testing.T) {
	publicKeyByte, err := hex.DecodeString(publicKeyHexString)
	assert.NoError(t, err)

	assert.EqualValues(t, Base58EncodeAddr(publicKeyByte), publicKeyBash58)
}

func TestBase58DecodeAddr(t *testing.T) {
	publicKeyByte, err := Base58DecodeAddr(publicKeyBash58)
	assert.NoError(t, err)

	assert.EqualValues(t, strings.ToTitle(hex.EncodeToString(publicKeyByte)), publicKeyHexString)
}
