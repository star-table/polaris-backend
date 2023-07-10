package domain

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
)

func TestGetUserInfoListByOrg(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		res, err := GetUserInfoListByOrg(2588)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))

}

func TestDetectUserInfoInUser(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(DetectUserInfoInUser(14396, []string{
			"+86-18917630563",
			"+86-18917630564",
			"+86-18917630565",
			"+86-18917630566",
			"+86-18917630567",
			"+86-18917630568",
			"+86-18917630569",
			"+86-18917630570",
		}))
	}))
}

func TestBatchGetUserDetailInfoWithMobile(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(BatchGetUserDetailInfoWithMobile([]int64{
			33965,
			1245489,
			33963,
			1112525,
			1111849,
			1511974,
			1511975,
		}))
	}))
}

func TestGetOrgIdListBySourceChannel(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		resp, err := GetOrgIdListBySourceChannel([]string{"ding"}, 1, 20, 0)
		if err != nil {
			t.Error(err)
		}
		t.Log(resp)
	}))
}

func TestGetBaseOrgInfo2(t *testing.T)  {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		baseOrgInfo, err := GetBaseOrgInfo(2588)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(baseOrgInfo)
	}))
}
