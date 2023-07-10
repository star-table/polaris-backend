package domain

import (
	"context"
	"fmt"
	"testing"

	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestUpdateIssueRelationSingle(t *testing.T) {

	//convey.Convey("更新任务关联测试", t, tests.StartUp(func() {
	//	convey.Convey("更新任务关联测试", func() {
	//		issueBo, err1 := GetIssueBo(17, 1011)
	//		if err1 != nil {
	//			return
	//		}
	//
	//		_, err3 := UpdateIssueRelationSingle(1, *issueBo, consts.IssueRelationTypeParticipant, []int64{1083})
	//		t.Log(err3)
	//
	//	})
	//}))

}

func TestUpdateIssueSort(t *testing.T) {
	//convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
	//	issue, err := GetIssueBo(1111, 3587)
	//	t.Log(err)
	//	beforeId := int64(3589)
	//	t.Log(UpdateIssueSort(*issue, 1007, &beforeId, nil, issue.Status))
	//}))

}

func TestAllIssueForProject(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		//t.Log(AllIssueForProject(1004, 1844, true, 0,0, []int64{}, nil))
	}))
}

func TestGetIssueListWithRefColumn(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(getTablesDataCount(2578, []int64{1558017733801545728, 1579751411611078656, 1525670440150765568}))
	}))

}
