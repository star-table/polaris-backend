package api

import (
	"context"
	"testing"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

var postGreeter = PostGreeter{}

var getGreeter = GetGreeter{}

func TestPostGreeter_Departments(t *testing.T) {
	convey.Convey("Test 获取部门列表", t, test.StartUp(func(ctx context.Context) {

		user := getGreeter.GetCurrentUser(ctx)

		convey.Convey("Test 获取部门列表", func() {
			isTop := 1

			respVo := postGreeter.Departments(orgvo.DepartmentsReqVo{
				Params: &vo.DepartmentListReq{
					IsTop: &isTop,
				},
				OrgId:         user.CacheInfo.OrgId,
				CurrentUserId: user.CacheInfo.UserId,
			})
			t.Log(json.ToJsonIgnoreError(respVo))
			convey.So(respVo.Failure(), convey.ShouldEqual, false)
		})
	}))
}

func TestPostGreeter_DepartmentMembers(t *testing.T) {
	convey.Convey("Test 获取部门成员列表", t, test.StartUp(func(ctx context.Context) {

		departmentId := int64(6)

		user := getGreeter.GetCurrentUser(ctx)

		convey.Convey("Test 获取部门成员列表", func() {
			respVo := postGreeter.DepartmentMembers(orgvo.DepartmentMembersReqVo{
				Params: vo.DepartmentMemberListReq{
					DepartmentID: &departmentId,
				},
				OrgId:         user.CacheInfo.OrgId,
				CurrentUserId: user.CacheInfo.UserId,
			})
			t.Log(json.ToJsonIgnoreError(respVo))
			convey.So(respVo.Failure(), convey.ShouldEqual, false)
		})
	}))
}
