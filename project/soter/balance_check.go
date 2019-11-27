package soter

import (
	"encoding/json"
	"strings"

	"github.com/TRON-US/chaos/utils"
)

type BalanceCheck struct {
	Address       string `json:"address"`
	Balance       int64  `json:"balance"`
	FreezeBalance int64  `json:"freeze_balance"`
	Timestamp     int    `json:"timestamp"`
}

// Get balance check by Lock struct.
func (balanceCheck *BalanceCheck) GetBalanceCheck() (string, error) {
	lockJsonString, err := json.Marshal(balanceCheck)
	if err != nil {
		return "", err
	}

	return utils.Md5(string(lockJsonString)), nil
}

// Verify user balance illegal.
func (balanceCheck *BalanceCheck) VerifyBalanceCheck(balanceLock string) bool {
	currentBalanceLock, err := balanceCheck.GetBalanceCheck()
	if err != nil {
		return false
	}

	if !strings.EqualFold(balanceLock, currentBalanceLock) {
		return false
	}
	return true
}
