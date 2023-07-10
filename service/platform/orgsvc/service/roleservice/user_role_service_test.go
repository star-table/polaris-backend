package service

import (
	"context"
	"testing"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetOrgRoleUser(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		t.Log(GetOrgRoleUser(1001, 0))
	}))
}

func TestMap(t *testing.T) {
	a := []map[string]string{}
	a = append(a, map[string]string{
		"a":    "a",
		"name": "jj",
	})
	a = append(a, map[string]string{
		"a":    "b",
		"name": "dd",
	})
	t.Log(maps.NewMap("a", a))
}

func TestUpdateUserOrgRole(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {

		req := rolevo.UpdateUserOrgRoleReqVo{
			OrgId:         1002,
			CurrentUserId: 1081,
			UserId:        926,
			RoleId:        8,
		}

		t.Log(UpdateUserOrgRole(req))
	}))
}

func TestGetProjectRoleList(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {

		data, err := GetProjectRoleList(10101, 10101)
		t.Log(json.ToJsonIgnoreError(data))
		t.Log(len(data))
		t.Log(err)
	}))
}

func TestUpdateUserOrgRoleBatch(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		t.Log(UpdateUserOrgRoleBatch(rolevo.UpdateUserOrgRoleBatchReqVo{
			OrgId:         1004,
			CurrentUserId: 1023,
			Input: rolevo.UpdateUserOrgRoleBatchData{
				UserIds: []int64{1011},
				RoleId:  1050,
			},
		}))
	}))
}

func TestGetOrgAdminUserBatch(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		resp, err := GetOrgAdminUserBatch([]int64{1004, 1417, 1420, 1428})
		if err != nil {
			t.Error(resp)
		}
		log.Info(json.ToJsonIgnoreError(resp))
	}))
}
