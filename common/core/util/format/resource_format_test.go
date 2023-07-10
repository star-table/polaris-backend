package format

import "testing"
import "github.com/magiconair/properties/assert"

func TestVerifyResourceNameFormat(t *testing.T) {
	testStr := ""
	for i := 0; i < 300; i++ {
		testStr = testStr + "1"
	}
	testStr = testStr + ".jpg"
	suc := VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 300; i++ {
		testStr = testStr + "1"
	}
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 300; i++ {
		testStr = testStr + "a"
	}
	testStr = testStr + ".jpg"
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 300; i++ {
		testStr = testStr + "好"
	}
	testStr = testStr + ".jpg"
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 300; i++ {
		testStr = testStr + "好"
	}
	testStr = testStr + "1111"
	testStr = testStr + ".jpg"
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "/"
	testStr = testStr + ".jpg"
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "\\"
	testStr = testStr + ".jpg"
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = ":"
	testStr = testStr + ".jpg"
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "?"
	testStr = testStr + ".jpg"
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "\""
	testStr = testStr + ".jpg"
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "<"
	testStr = testStr + ".jpg"
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = ">"
	testStr = testStr + ".jpg"
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "|"
	testStr = testStr + ".jpg"
	suc = VerifyResourceNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)
}

func TestVerifyFolderNameFormat(t *testing.T) {
	testStr := ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "1"
	}
	suc := VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "a"
	}
	suc = VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "好"
	}
	suc = VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "好"
	}
	testStr = testStr + "1111"
	suc = VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "/"
	suc = VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "\\"
	suc = VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = ":"
	suc = VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "?"
	suc = VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "\""
	suc = VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "<"
	suc = VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = ">"
	suc = VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "|"
	suc = VerifyFolderNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)
}
