package format

import "unicode/utf8"

func VerifyDepartmentName(name string) bool {
	if name == "" || utf8.RuneCountInString(name) > 20 {
		return false
	}

	return true
}
