package service

import (
	"context"
	"testing"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/bo/operation"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

//从缓存获取操作的集合
func TestGetRoleOperationListFromCache(t *testing.T) {

	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("加载配置", func() {

			convey.Convey("test", func() {

				roleOperationList, _ := GetRoleOperationList()

				convey.So(roleOperationList, convey.ShouldNotBeNil)

			})
		})
	}))
}

func TestAuthenticateLocal(t *testing.T) {
	convey.Convey("Test AppendIterationStat", t, test.StartUp(func(ctx context.Context) {
		t.Log(json.ToJsonIgnoreError(config.GetMysqlConfig()))

		orgId := int64(1065)
		projectId := int64(1144)
		issueId := int64(1511)
		projectAuthInfo := &bo.ProjectAuthBo{
			Id:    projectId,
			Owner: 1233,
		}
		issueAuthInfo := &bo.IssueAuthBo{
			Id: issueId,
			//Owner: 1233,
		}

		err := Authenticate(orgId, 1235, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/4", operation.Delete, nil)
		t.Log(err)
		err = Authenticate(orgId, 1233, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/4", operation.Delete, nil)
		t.Log(err)
		err = Authenticate(orgId, 1235, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/4", operation.Delete, nil)
		t.Log(err)
		err = Authenticate(orgId, 1233, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/4", operation.Modify, nil)
		t.Log(err)
		err = Authenticate(orgId, 1233, projectAuthInfo, nil, "/Org/{org_id}/Pro/{pro_id}/ProConfig", operation.ModifyStatus, nil)
		t.Log(err)
		err = Authenticate(orgId, 1233, projectAuthInfo, nil, "/Org/{org_id}/Pro/{pro_id}", operation.ModifyStatus, nil)
		t.Log(err)
		//err = Authenticate(orgId, 1029, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/T", operation.Delete)
		//t.Log(err)
		//err = Authenticate(orgId, 1029, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/T", operation.Delete)
		//t.Log(err)
		//err = Authenticate(orgId, 1029, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/T", operation.Delete)
	}))
}

//未知的缓存
func TestGetRoleOperationListUnknowCache(t *testing.T) {

	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("加载配置", func() {

			convey.Convey("test", func() {

				roleOperationList, _ := GetRoleOperationList()

				convey.So(roleOperationList, convey.ShouldBeNil)

			})
		})
	}))
}

func TestClearRolePermissionOperationList(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		t.Log(ClearRolePermissionOperationList(10101, 10231, 0))
	}))
}
