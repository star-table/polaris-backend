package format

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

//组织名
func VerifyOrgNameFormat(input string) bool {
	length := utf8.RuneCountInString(strings.TrimSpace(input))
	if length == 0 || length > 30 {
		return false
	}
	return true
}

//网址后缀
func VerifyOrgCodeFormat(input string) bool {
	reg := regexp.MustCompile(OrgCodePattern)
	return reg.MatchString(input)
}

//详细地址
func VerifyOrgAdressFormat(input string) bool {
	reg := regexp.MustCompile(OrgAdressPattern)
	return reg.MatchString(input)
}
