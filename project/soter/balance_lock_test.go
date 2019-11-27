package soter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	lock = &Lock{
		Address:       "address",
		Balance:       100,
		FreezeBalance: 100,
		UpdateTime:    1218154088,
	}
	balanceLock = "fdddf1603f98d20921a0de46e51a4774"
)

func TestLock_GetBalanceLock(t *testing.T) {
	currentBalanceLock, err := lock.GetBalanceLock()
	assert.NoError(t, err)
	assert.EqualValues(t, currentBalanceLock, balanceLock)
}

func TestLock_CheckBalanceLock(t *testing.T) {
	assert.True(t, true, lock.CheckBalanceLock(balanceLock))
}
