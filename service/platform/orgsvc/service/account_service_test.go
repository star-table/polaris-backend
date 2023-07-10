package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"gotest.tools/assert"
)

func TestSetPassword(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		err := SetPassword(orgvo.SetPasswordReqVo{
			UserId: userId,
			OrgId:  orgId,
			Input: vo.SetPasswordReq{
				"adf",
			},
		})
		t.Log(err)
		err = SetPassword(orgvo.SetPasswordReqVo{
			UserId: userId,
			OrgId:  orgId,
			Input: vo.SetPasswordReq{
				"123134646",
			},
		})
		t.Log(err)
		err = SetPassword(orgvo.SetPasswordReqVo{
			UserId: 1001,
			OrgId:  orgId,
			Input: vo.SetPasswordReq{
				"helloworld",
			},
		})
		t.Log(err)
	}))
}

func TestUserLogin(t *testing.T) {
	convey.Convey("TestUserLogin", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		password := "helloworld"
		user, err := UserLogin(vo.UserLoginReq{
			LoginType: 2,
			LoginName: "18221304331",
			Password:  &password,
		})
		t.Log(err)
		t.Log(json.ToJsonIgnoreError(user))

		user, err = UserLogin(vo.UserLoginReq{
			LoginType: 2,
			LoginName: "ainililia@163.com",
			Password:  &password,
		})
		t.Log(err)
		t.Log(json.ToJsonIgnoreError(user))

		password = "helloworld1"
		user, err = UserLogin(vo.UserLoginReq{
			LoginType: 2,
			LoginName: "18221304331",
			Password:  &password,
		})
		t.Log(err)
		t.Log(json.ToJsonIgnoreError(user))

		user, err = UserLogin(vo.UserLoginReq{
			LoginType: 2,
			LoginName: "ainililia@163.com",
			Password:  &password,
		})
		t.Log(err)
		t.Log(json.ToJsonIgnoreError(user))
	}))
}

func TestResetPassword(t *testing.T) {
	convey.Convey("TestResetPassword", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		password := "helloworld1"
		err := ResetPassword(orgvo.ResetPasswordReqVo{
			UserId: userId,
			OrgId:  orgId,
			Input: vo.ResetPasswordReq{
				CurrentPassword: "helloworld",
				NewPassword:     password,
			},
		})
		t.Log(err)
	}))
}

func TestBindLoginName(t *testing.T) {
	convey.Convey("TestBindLoginName", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		err := SendAuthCode(orgvo.SendAuthCodeReqVo{
			Input: vo.SendAuthCodeReq{
				AuthType:    5,
				AddressType: 2,
				Address:     "ainililia@163.com",
			},
		})
		assert.Equal(t, err, nil)
		time.Sleep(3 * time.Second)

		err = BindLoginName(orgvo.BindLoginNameReqVo{
			UserId: userId,
			OrgId:  orgId,
			Input: vo.BindLoginNameReq{
				Address:     "ainililia@163.com",
				AddressType: consts.ContactAddressTypeEmail,
				AuthCode:    "000000",
			},
		})
		t.Log(err)

		//err = SendAuthCode(orgvo.SendAuthCodeReqVo{
		//	Input: vo.SendAuthCodeReq{
		//		AuthType: 6,
		//		AddressType: 2,
		//		Address:"ainililia@163.com",
		//	},
		//})
		//assert.Equal(t, err, nil)
		//time.Sleep(3 * time.Second)
		//
		//err = UnbindLoginName(orgvo.UnbindLoginNameReqVo{
		//	UserId: userId,
		//	OrgId:  orgId,
		//	Input: vo.UnbindLoginNameReq{
		//		AddressType: consts.ContactAddressTypeEmail,
		//		AuthCode: "000000",
		//	},
		//})
		//t.Log(err)
	}))
}

func TestRetrievePassword(t *testing.T) {
	convey.Convey("TestBindLoginName", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		//err := SendAuthCode(orgvo.SendAuthCodeReqVo{
		//	Input: vo.SendAuthCodeReq{
		//		AuthType:    4,
		//		AddressType: 2,
		//		Address:     "ainililia@163.com",
		//	},
		//})
		//assert.Equal(t, err, nil)
		//time.Sleep(3 * time.Second)
		//
		//err = RetrievePassword(orgvo.RetrievePasswordReqVo{
		//	Input: vo.RetrievePasswordReq{
		//		Username:    "ainililia@163.com",
		//		AuthCode:    "000000",
		//		NewPassword: "helloworld",
		//	},
		//})
		//assert.Equal(t, err, nil)
		//time.Sleep(3 * time.Second)

	}))
}

func TestVerifyOldName(t *testing.T) {
	convey.Convey("TestVerifyOldName", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(VerifyOldName(orgvo.UnbindLoginNameReqVo{Input: vo.UnbindLoginNameReq{
			AddressType: 1,
			AuthCode:    "000000",
		},
			UserId: 1284,
			OrgId:  1111}))
	}))
}

func TestUserLogin2(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		password := "zZ123456"
		res, err := UserLogin(vo.UserLoginReq{
			LoginType: 2,
			LoginName: "suhanyu4@bjx.cloud",
			Password:  &password,
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}

func TestThirdAccountList(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		list, _ := ThirdAccountBindList(0, 34845)
		fmt.Println(list.List[0])
	}))
}

func TestUnbindThirdAccount(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		err := UnbindThirdAccount(orgvo.UnbindAccountReq{
			OrgId:  0,
			UserId: 29598,
			Input:  &orgvo.UnbindAccountData{},
		})
		fmt.Println(err)
	}))
}

func TestCheckMobileHasBind(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		list, _ := CheckMobileHasBind(&orgvo.CheckMobileHasBindData{
			Mobile:         "+86-17674120696",
			SourcePlatform: "",
		})
		fmt.Println(list)
	}))
}
