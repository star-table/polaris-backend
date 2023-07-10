package domain

import (
	"context"
	"fmt"
	"testing"

	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetBaseUserOutInfoByUserIds(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		//t.Log(GetBaseUserOutInfoByUserIds("fs", 1082, []int64{1167,1168,1166,1862,1871}))
		t.Log(GetBaseUserOutInfoByUserIds(1113, []int64{1289}))
	}))
}

func TestGetBaseUserInfoBatch(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		//t.Log(GetBaseUserInfoBatch("", 1003, []int64{1006,1007,1008,1012,1013}))
		t.Log(GetBaseUserInfoBatch(1113, []int64{1289}))
		//t.Log(GetBaseUserOutInfoBatch("ding," 1))
	}))
}

func TestGetChangeLoginNameSign(t *testing.T) {
	convey.Convey("TestGetChangeLoginNameSign", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		//t.Log(GetChangeLoginNameSign("11"))
	}))

}

func TestGetWhiteListVipOrg(t *testing.T) {
	convey.Convey("TestGetChangeLoginNameSign", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(GetWhiteListVipOrg(1013))
	}))

}

func TestCheckUserPay(t *testing.T) {
	convey.Convey("TestGetChangeLoginNameSign", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(cache.HSet("abc", "a", "a"))
		t.Log(cache.Expire("abc", 50))
		t.Log(cache.Get("abc"))
		t.Log(cache.Exist("abc"))
		//t.Log(cache.HGet("abc","a"))
	}))
}

func TestGetBaseOrgInfo(t *testing.T) {
	convey.Convey("TestGetChangeLoginNameSign", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(GetBaseOrgInfo(14396))
	}))
}

func TestClearAllOrgUserPayCache(t *testing.T) {
	convey.Convey("ClearAllOrgUserPayCache", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(ClearAllOrgUserPayCache(2373))
	}))
}

func TestRemoveOrgMember(t *testing.T) {
	convey.Convey("ClearAllOrgUserPayCache", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(RemoveOrgMember(3151, []int64{34906}, 0))
	}))
}
