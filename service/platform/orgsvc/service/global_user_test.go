package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/star-table/polaris-backend/common/core/consts"

	"github.com/star-table/polaris-backend/common/model/vo"

	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestLoginMobile(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(loginMobile(vo.UserLoginReq{
			LoginType:      consts.LoginTypeSMSCode,
			LoginName:      "+86-18917630569",
			SourceChannel:  consts.AppSourceChannelWeb,
			SourcePlatform: consts.AppSourcePlatformLarkXYJH2019,
		}, "我的天啊", ""))
	}))
}

func TestBindMobile(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(domain.BindUserName(1004, 1014, consts.ContactAddressTypeMobile, "+86-18917630568"))
	}))
}
