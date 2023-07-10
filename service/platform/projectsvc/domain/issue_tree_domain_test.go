package domain

import (
	"testing"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/bo"
)

func TestMakeTree(t *testing.T) {
	res := []*bo.IssueNode{
		&bo.IssueNode{
			IssueId: 1,
			Path:    "0,",
			Child:   nil,
		},
		&bo.IssueNode{
			IssueId: 2,
			Path:    "0,1,",
			Child:   nil,
		},
		&bo.IssueNode{
			IssueId: 3,
			Path:    "0,1,2,",
			Child:   nil,
		},
		&bo.IssueNode{
			IssueId: 4,
			Path:    "0,1,2,3,",
			Child:   nil,
		},
		&bo.IssueNode{
			IssueId: 5,
			Path:    "0,1,",
			Child:   nil,
		},
	}
	node := res[0]
	MakeTree(res, node)
	t.Log(json.ToJsonIgnoreError(node))
}

//func TestCreateProjectDefaultView(t *testing.T) {
//	convey.Convey("TestGetIssueWorkHoursList", t, test.StartUp(func(ctx context.Context) {
//		trueFlag := true
//		tableId := int64(1)
//		t.Log(CreateProjectDefaultView(1749, 11257, 1470309924822155266, &trueFlag, &tableId))
//	}))
//}
