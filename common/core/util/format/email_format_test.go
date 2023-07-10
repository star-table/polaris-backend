package format

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestVerifyEmailFormat(t *testing.T) {
	suc := VerifyEmailFormat("xxxx@163.com")
	assert.Equal(t, suc, true)
	suc = VerifyEmailFormat("x@1.c")
	assert.Equal(t, suc, true)
	suc = VerifyEmailFormat("@163.com")
	assert.Equal(t, suc, false)
	suc = VerifyEmailFormat("xxxx@163")
	assert.Equal(t, suc, false)
	suc = VerifyEmailFormat("xxxx163.com")
	assert.Equal(t, suc, false)
}