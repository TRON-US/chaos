package soter

import (
	"encoding/json"
	"strings"
)

type Lock struct {
	Address       string `json:"address"`
	Balance       int    `json:"balance"`
	FreezeBalance int    `json:"freeze_balance"`
	UpdateTime    int    `json:"update_time"`
}

// Get balance lock by Lock struct.
func (lock *Lock) GetBalanceLock() (string, error) {
	lockJsonString, err := json.Marshal(lock)
	if err != nil {
		return "", err
	}

	return Md5(string(lockJsonString)), nil
}

// Verify user balance illegal.
func (lock *Lock) CheckBalanceLock(balanceLock string) bool {
	currentBalanceLock, err := lock.GetBalanceLock()
	if err != nil {
		return false
	}

	if !strings.EqualFold(balanceLock, currentBalanceLock) {
		return false
	}
	return true
}
