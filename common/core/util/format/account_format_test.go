package format

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestVerifyPwdFormat(t *testing.T) {
	suc := VerifyPwdFormat("a123123")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyPwdFormat("a123123")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyPwdFormat("a123123a")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyPwdFormat("a123a123a")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyPwdFormat("a123a123")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyPwdFormat("a123a123")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyPwdFormat("1a")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyPwdFormat("aAAAAA1")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyPwdFormat("a")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyPwdFormat("1")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyPwdFormat("A1a&")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyPwdFormat("a1a.")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyPwdFormat("a%!#12")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyPwdFormat("*")
	t.Log(suc)
	assert.Equal(t, suc, false)
}

func TestVerifyUserNameFormat(t *testing.T) {
	suc := VerifyUserNameFormat("好好好好好好好好好好")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyUserNameFormat("11111111111111111111")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyUserNameFormat("llllllllllllllllllll")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyUserNameFormat("好好好好好好好好好1K")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyUserNameFormat("好好好好好好好好好1Kl")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyUserNameFormat("*")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyUserNameFormat("hao好*")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyUserNameFormat("hao1265")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyUserNameFormat("hao1265             ")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyUserNameFormat("hao1265              ")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyUserNameFormat("  ")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyUserNameFormat(" a哈1和🏃1sssssssssssssss ")
	t.Log(suc)
	assert.Equal(t, suc, false)
}
