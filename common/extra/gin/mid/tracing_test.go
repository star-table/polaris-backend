package mid

import (
	"fmt"
	"github.com/star-table/common/core/util/strs"
	"strings"
	"testing"
)

func TestStartTrace(t *testing.T) {
	fmt.Println(1)

	//s := "/api/task/health"
	s := "health"
	urlCount := strs.Len(s)
	if urlCount <= 6 || s[urlCount-6:] == "health" {
		fmt.Println("asasdfasdf")
	} else {
		fmt.Println("123123123")
	}

	fmt.Printf("%d", strings.LastIndex(s, "health"))

	sa := "/a\npi/task/health"
	sb := strings.ReplaceAll(sa, "\n", "\\n")

	fmt.Println(sa)
	fmt.Println(sb)
}
