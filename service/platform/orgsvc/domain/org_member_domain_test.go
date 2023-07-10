package domain

import (
	"context"
	"testing"

	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestUpdateUserInfoWithFeiShu(t *testing.T) {
	convey.Convey("TestGetDepartmentIdsByOutDepartmentIds", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		//needUpdateUserOpenId := "ou_83776f77689ab1a5f399d5aafa27a0a4"
		//orgOutInfo, err := GetBaseOrgInfo(sdk_const.SourceChannelFeishu, orgId)
		//if err != nil {
		//	log.Error(err)
		//	return
		//}
		//tenant, err := feishu.GetTenant(orgOutInfo.OutOrgId)
		//if err != nil {
		//	log.Error(err)
		//	return
		//}
		//userBatchResp, err2 := tenant.GetUserBatchGetV2(nil, []string{needUpdateUserOpenId})
		//if err2 != nil {
		//	log.Error(err2)
		//	return
		//}
		//usersKeyByOpenId := map[string]sdkVo.UserDetailInfoV2{}
		//for _, tmpUser := range userBatchResp.NewData.Users {
		//	usersKeyByOpenId[tmpUser.OpenId] = tmpUser
		//}
		//var baseUserInfo sdkVo.UserDetailInfoV2
		//if oneUser, ok := usersKeyByOpenId[needUpdateUserOpenId]; !ok {
		//	log.Errorf("用户信息不存在")
		//	return
		//} else {
		//	baseUserInfo = oneUser
		//	convey.So(baseUserInfo, convey.ShouldNotEqual, nil)
		//}
		//sysErr := UpdateUserInfoWithPlatformUserInfo(orgId, 1305473375526715392, baseUserInfo)
		//t.Logf("%v\n", sysErr)
		//convey.So(sysErr, convey.ShouldEqual, nil)
	}))
}

func TestPayLevelIsExist(t *testing.T)  {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		exist, err := PayLevelIsExist(3)
		if err != nil {
			t.Error(err)
		}
		t.Log(exist)
	}))
}

func TestGetLocalOrgUserIdMap(t *testing.T)  {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		resp, err := GetLocalOrgUserIdMap([]int64{34221, 34220})
		if err != nil {
			t.Error(err)
		}
		t.Log(resp)
	}))
}
