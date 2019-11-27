package soter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	balanceCheck = &BalanceCheck{
		Address:       "address",
		Balance:       100,
		FreezeBalance: 100,
		Timestamp:     1218154088,
	}
	check = "f04da9dbef60fde70a45611a8e806ac9"
)

func TestBalanceCheck_GetBalanceCheck(t *testing.T) {
	currentBalanceLock, err := balanceCheck.GetBalanceCheck()
	assert.NoError(t, err)
	assert.EqualValues(t, currentBalanceLock, check)
}

func TestBalanceCheck_VerifyBalanceCheck(t *testing.T) {
	assert.True(t, true, balanceCheck.VerifyBalanceCheck(check))
}
