package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/star-table/common/core/types"

	"github.com/star-table/polaris-backend/common/model/vo/orgvo"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetDingTalkBaseUserInfo(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		res, err := GetDingTalkBaseUserInfo(1113, 1289)
		t.Log(json.ToJsonIgnoreError(res))
		t.Log(err)
	}))
}

func TestVerifyDepartments(t *testing.T) {
	convey.Convey("TestVerifyDepartments", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(VerifyDepartments(1004, []int64{1004}))
	}))
}

func TestGetBaseUserInfo1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		info, err := GetBaseUserInfo(2373, 29611)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(info))
	}))
}

func TestGetUserConfigInfoBatch1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		resData, err := GetUserConfigInfoBatch(2373, &orgvo.GetUserConfigInfoBatchReqVoInput{
			UserIds: []int64{29611},
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resData))
	}))
}

func TestGetBaseUserInfoByEmpIdBatch1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		list, err := GetBaseUserInfoByEmpIdBatch(2373, orgvo.GetBaseUserInfoByEmpIdBatchReqVoInput{
			OpenIds: []string{"ou_3ab7fe596cf91692218f744558ae157f"},
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(list))
	}))
}

func TestGetBaseOrg1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(2373)
		orgOutInfo, err := GetBaseOrgInfo(orgId)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(orgOutInfo))
	}))
}

func TestTime(t *testing.T) {
	fmt.Println(time.Time(types.NowTime()).Unix())
	fmt.Println(time.Now().Unix())
}
