package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testData = "test"
	testMd5  = "098f6bcd4621d373cade4e832627b4f6"
)

func TestMd5(t *testing.T) {
	assert.EqualValues(t, Md5(testData), testMd5)
}
