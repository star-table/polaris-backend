package service

import (
	"context"
	"testing"

	"github.com/star-table/polaris-backend/facade/permissionfacade"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/util/jsonx"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"gotest.tools/assert"
)

func TestPermissionOperationList(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		//var projectId int64 = 10101
		data, err := PermissionOperationList(1002, 1022, 1081, nil)
		t.Log(json.ToJsonIgnoreError(data))
		assert.Equal(t, err, nil)
	}))
}

func TestUpdateRolePermissionOperation(t *testing.T) {
	convey.Convey("Test UpdateRolePermissionOperation", t, test.StartUp(func(ctx context.Context) {
		updArr := []*vo.EveryPermission{}
		upd := vo.EveryPermission{PermissionID: 15, OperationIds: []int64{51}}
		data, err := UpdateRolePermissionOperation(1004, 1023, vo.UpdateRolePermissionOperationReq{
			RoleID:            1056,
			UpdatePermissions: append(updArr, &upd),
		})
		t.Log(json.ToJsonIgnoreError(data))
		assert.Equal(t, err, nil)
	}))
}

//func TestGetPersonalPermissionInfo(t *testing.T) {
//	convey.Convey("Test UpdateRolePermissionOperation", t, test.StartUp(func(ctx context.Context) {
//		var projectId int64 = 1045
//		var issueId int64 = 1281
//		//t.Log(GetPersonalPermissionInfo(1001, 1002, &projectId, nil, ""))
//		t.Log(GetPersonalPermissionInfo(1004, 1023, &projectId, &issueId, ""))
//	}))
//}

func TestUpdateRolePermissionOperation2(t *testing.T) {
	convey.Convey("Test UpdateRolePermissionOperation", t, test.StartUp(func(ctx context.Context) {
		updArr := []*vo.EveryPermission{}
		updArr = append(updArr, &vo.EveryPermission{PermissionID: 15, OperationIds: []int64{}})
		updArr = append(updArr, &vo.EveryPermission{PermissionID: 21, OperationIds: []int64{}})
		updArr = append(updArr, &vo.EveryPermission{PermissionID: 27, OperationIds: []int64{}})
		updArr = append(updArr, &vo.EveryPermission{PermissionID: 33, OperationIds: []int64{}})
		updArr = append(updArr, &vo.EveryPermission{PermissionID: 38, OperationIds: []int64{}})
		updArr = append(updArr, &vo.EveryPermission{PermissionID: 39, OperationIds: []int64{}})
		updArr = append(updArr, &vo.EveryPermission{PermissionID: 40, OperationIds: []int64{}})
		updArr = append(updArr, &vo.EveryPermission{PermissionID: 41, OperationIds: []int64{}})
		projectId := int64(1504)
		data, err := UpdateRolePermissionOperation(1113, 1289, vo.UpdateRolePermissionOperationReq{
			RoleID:            1682,
			ProjectID:         &projectId,
			UpdatePermissions: updArr,
		})
		t.Log(json.ToJsonIgnoreError(data))
		assert.Equal(t, err, nil)
	}))
}

func TestGetOrgRoleDepartment(t *testing.T) {
	//convey.Convey("Test UpdateRolePermissionOperation", t, test.StartUp(func(ctx context.Context) {
	//	t.Log(GetOrgRoleDepartment(1004, 1736))
	//}))
}

func TestOperationCodeToSlice(t *testing.T) {
	a := "{\"code\":0,\"message\":\"\",\"data\":{\"deptIds\":[\"1557\"]}}"
	respVo := &orgvo.GetUserDeptIdsResp{}

	t.Log(jsonx.FromJson(a, respVo))
	t.Log(respVo)
}

func TestGetPersonalPermissionInfoForFuse1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		resp, err := GetPersonalPermissionInfoForFuse(1574, 24370, Int64Ptr(-1), Int64Ptr(0), "fs")
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func Int64Ptr(num int64) *int64 {
	return &num
}

func TestGetAppAuth(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		optAuthResp := permissionfacade.GetAppAuth(999, 1452840747366711298, 1452840747366711298, 2)
		//optAuthResp := permissionfacade.GetAppAuth(999, 1452873242346881026, 2)
		if optAuthResp.Failure() {
			t.Error(optAuthResp.Error())
			return
		}
		t.Log(json.ToJsonIgnoreError(optAuthResp))
	}))
}
