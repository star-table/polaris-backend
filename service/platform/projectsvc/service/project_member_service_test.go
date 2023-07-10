package service

import (
	"context"
	"testing"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"gotest.tools/assert"
)

func TestProjectUserList(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		data, err := ProjectUserList(1242, 1, 0, vo.ProjectUserListReq{ProjectID: 8374})
		t.Log(json.ToJsonIgnoreError(data))
		assert.Equal(t, err, nil)

		//t.Log(AddProjectMember(1001, 1001, "", vo.RemoveProjectMemberReq{ProjectID:1001, MemberIds:[]int64{1002}}))
		//t.Log(RemoveProjectMember(1001, 1001, "", vo.RemoveProjectMemberReq{ProjectID:1001, MemberIds:[]int64{1002}}))
	}))
}

func TestAddProjectMember(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		//t.Log(AddProjectMember(10101, 10201, vo.RemoveProjectMemberReq{ProjectID:10101, MemberIds:[]int64{10201, 10203}}))
		t.Log(AddProjectMember(1323, 1608, "", vo.RemoveProjectMemberReq{ProjectID: 1521, MemberIds: []int64{1464, 1609}}))
	}))
}

//func TestHomeIssuesGroup(t *testing.T) {
//	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
//		relationType := 2
//		projectId := int64(1704)
//		res, err := HomeIssuesGroup(1004, 1021, -1, -1, &vo.HomeIssueInfoReq{
//			ProjectID: &projectId,
//			GroupType: &relationType,
//		})
//		t.Log(json.ToJsonIgnoreError(res))
//		t.Log(err)
//	}))
//}
