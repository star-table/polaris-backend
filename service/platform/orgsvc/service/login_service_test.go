package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/star-table/polaris-backend/common/model/vo/orgvo"

	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestUserLogin1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		pwd := "Ab123456"
		res, err := UserLogin(vo.UserLoginReq{
			LoginType:      2,
			LoginName:      "admin@admin.com",
			Password:       &pwd,
			AuthCode:       nil,
			Name:           nil,
			InviteCode:     nil,
			SourceChannel:  "",
			SourcePlatform: "",
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	}))
}

func TestExchangeShareToken(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(ExchangeShareToken(orgvo.ExchangeShareTokenReq{Input: orgvo.ExchangeShareTokenData{
			ShareKey: "TOXSeAa",
			Password: "%IFDjfdsk",
		}}))
	}))
}
