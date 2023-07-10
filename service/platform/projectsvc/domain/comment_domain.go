package domain

import (
	"regexp"
	"strings"

	"github.com/star-table/polaris-backend/common/core/consts"
)

func CutComment(input string, limit int) string {
	re := regexp.MustCompile(`(<at id=ou_[A-Za-z0-9]+></at>)|( \[附件\])`)
	match := re.FindAllIndex([]byte(input), -1)

	var idx []int
	for _, m := range match {
		idx = append(idx, m...)
	}

	total := 0
	last := 0
	var cut strings.Builder
	for i := 0; i < len(idx); i += 2 {
		idxFrom := idx[i]
		idxTo := idx[i+1]
		if idxFrom > last {
			s := input[last:idxFrom]
			length := len(s)
			if total+length > limit {
				cut.WriteString(s[0 : limit-total])
				cut.WriteString(consts.CardIssueChangeDescTextOverflow)
				total = limit
				break
			} else {
				cut.WriteString(input[last:idxFrom])
				cut.WriteString(input[idxFrom:idxTo])
				last = idxTo
				total += length
			}
		} else {
			last = idxTo
			cut.WriteString(input[idxFrom:idxTo])
		}
	}
	if len(input) > last && total < limit {
		s := input[last:]
		length := len(s)
		if total+length > limit {
			cut.WriteString(s[0 : limit-total])
			cut.WriteString(consts.CardIssueChangeDescTextOverflow)
			total = limit
		} else {
			cut.WriteString(input[last:])
			total += length
		}
	}
	return cut.String()
}
