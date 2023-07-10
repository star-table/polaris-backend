package format

import "regexp"

//资源名
func VerifyResourceNameFormat(input string) bool {
	reg := regexp.MustCompile(ResourceNamePattern)
	return reg.MatchString(input)
}

//文件夹名
func VerifyFolderNameFormat(input string) bool {
	reg := regexp.MustCompile(FolderNamePattern)
	return reg.MatchString(input)
}
