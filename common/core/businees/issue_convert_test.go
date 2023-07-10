package businees

import (
	"fmt"
	"testing"
)

//func TestConvertIssueStatusFilterReq(t *testing.T) {
//fmt.Println(json.ToJsonIgnoreError(ConvertIssueStatusFilterReq(vo.LessCondsData{
//	Type:      "and",
//	Conds:     []*vo.LessCondsData{
//		&vo.LessCondsData{
//			Type:      "in",
//			Values:    []interface{}{7},
//			Column:    "issueStatus",
//			Conds:     nil,
//		},
//		&vo.LessCondsData{
//			Type:      "equal",
//			Value:    12,
//			Column:    "projectId",
//			Conds:     nil,
//		},
//	},
//})))
//}

func TestLcMemberToUserIdsWithError(t *testing.T) {
	fmt.Println(LcMemberToUserIdsWithError([]string{"U_123", "0"}, true))
	fmt.Println(LcMemberToUserIdsWithError([]string{"U_123", "0", "sdfsf"}, true))
	fmt.Println(LcMemberToUserIdsWithError([]string{"U_123", "U_0"}))

}
