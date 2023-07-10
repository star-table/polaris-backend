package util

import (
	"github.com/star-table/common/core/util/file"
	"strings"
)

func GetCurrentPath() string {
	return file.GetCurrentPath()
}

func GetMobile(mobile string) string{
	strs := strings.Split(mobile, "-")
	if len(strs) == 1 {
		return strs[0]
	}
	return strs[1]
}