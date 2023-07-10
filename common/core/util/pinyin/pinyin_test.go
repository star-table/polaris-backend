package pinyin

import (
	"github.com/mozillazg/go-pinyin"
	"testing"
)

func TestConvertToPinyin(t *testing.T) {

	strs := pinyin.Pinyin("刘s2千源", args)
	for k, str := range strs {
		t.Log(k, str)
	}
	t.Log(ConvertToPinyin("s23千源"))
	t.Log(ConvertToPinyin("abc"))
	t.Log(ConvertToPinyin("s23千源"))
	t.Log(ConvertToPinyin("123刘3123"))
	t.Log(ConvertToPinyin("1233123刘"))


}
