package format

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestVerifyOrgNameFormat(t *testing.T) {
	suc := VerifyOrgNameFormat("11111111111111111111")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyOrgNameFormat("")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyOrgNameFormat("aaaaaaaaaaaaaaaaaaaa")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyOrgNameFormat("你你你你你你你你你你你你你你你你你你你你")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyOrgNameFormat("你你你111你你你你你你你aaa你你你你")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyOrgNameFormat("你你你111你你你你你你你aaa你你你你a")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyOrgNameFormat("你你你111你你你你你你你aaa你你你你1")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyOrgNameFormat("你你你111你你你你你你你aaa你你你你你")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyOrgNameFormat("你你你111你你你你你你你aaa你你你")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyOrgNameFormat("\\.*")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyOrgNameFormat("^-")
	t.Log(suc)
	assert.Equal(t, suc, false)
}

func TestVerifyOrgCodeFormat(t *testing.T) {
	suc := VerifyOrgCodeFormat("11111111111111111111")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyOrgCodeFormat("aaaaaaaaaaaaaaaaaaaa")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyOrgCodeFormat("aaaaaaaaaaaaaaaaa111")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyOrgCodeFormat("aaaaaaaaaaaaaaaaa111a")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyOrgCodeFormat("aaaaaaaaaaaaaaa.*")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyOrgCodeFormat("/")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyOrgCodeFormat("啊啊")
	t.Log(suc)
	assert.Equal(t, suc, false)
}

func TestVerifyOrgAdressFormat(t *testing.T) {
	testStr := ""
	for i := 0; i < 100; i++ {
		testStr = testStr + "1"
	}
	suc := VerifyOrgAdressFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 100; i++ {
		testStr = testStr + "a"
	}
	suc = VerifyOrgAdressFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 100; i++ {
		testStr = testStr + "好"
	}
	suc = VerifyOrgAdressFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 100; i++ {
		testStr = testStr + "/"
	}
	suc = VerifyOrgAdressFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 100; i++ {
		testStr = testStr + "好"
	}
	testStr = testStr + "1"
	suc = VerifyOrgAdressFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)
}
