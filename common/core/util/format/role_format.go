package format

import "regexp"

//角色名
func VerifyRoleNameFormat(input string) bool {
	reg := regexp.MustCompile(ChinesePattern)
	formInput := reg.ReplaceAllString(input, "a")
	reg = regexp.MustCompile(RoleNamePattern)
	return reg.MatchString(formInput)
}
