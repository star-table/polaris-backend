package pinyin

import (
	"github.com/mozillazg/go-pinyin"
	"regexp"
)

var args = pinyin.NewArgs()
var hzRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]$")

func ConvertToPinyin(str string) string {
	if str == "" {
		return str
	}
	result := ""
	for _, c := range str {
		target := string(c)
		strs := pinyin.Pinyin(target, args)
		if len(strs) != 0 {
			target = Capitalize(strs[0][0])
		}
		result += target
	}
	return result
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	if str == "" {
		return str
	}
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
