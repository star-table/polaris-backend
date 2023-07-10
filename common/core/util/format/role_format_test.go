package format

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestVerifyRoleNameFormat(t *testing.T) {
	suc := VerifyRoleNameFormat("好好好好好好好好好好")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyRoleNameFormat("11111111111111111111")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyRoleNameFormat("llllllllllllllllllll")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyRoleNameFormat("好好好好好好好好好1K")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyRoleNameFormat("好好好好好好好好好1Kl")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyRoleNameFormat("*")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyRoleNameFormat("hao好*")
	t.Log(suc)
	assert.Equal(t, suc, false)
}
