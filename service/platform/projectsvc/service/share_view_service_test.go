package service

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/star-table/polaris-backend/common/model/vo/projectvo"

	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestCreateShareView(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(CreateShareView(&projectvo.CreateShareViewReq{
			OrgId:  3082,
			UserId: 34845,
			Input: &projectvo.CreateShareViewData{
				AppId:     1643095250803195906,
				ProjectId: 62889,
				TableId:   1643095251545493504,
				ViewId:    1643095250803195908,
			},
		}))
	}))
}

func TestGetShareViewInfo(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {
		regp := regexp.MustCompile(`/\d*/`)
		fmt.Println(regp.ReplaceAllString("/api/rest/project/1/issue/filter?withoutInfo=1&isFiling=3", "/*/"))
		//fmt.Println(GetShareViewInfo(&projectvo.GetShareViewInfoReq{
		//	OrgId:  3082,
		//	UserId: 34845,
		//	ViewId: 1643095250803195908,
		//}))
	}))
}

func TestGetShareViewInfoByKey(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(GetShareViewInfoByKey(&projectvo.GetShareViewInfoByKeyReq{
			Input: &projectvo.ShareKeyData{ShareKey: "59oXwpRu"},
		}))
	}))
}

func TestUpdatePassword(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(UpdateSharePassword(&projectvo.UpdateSharePasswordReq{
			OrgId:  3082,
			UserId: 34845,
			Input: &projectvo.UpdateData{
				ViewId:   1643095250803195908,
				Config:   "{}",
				Password: "%IFDjfdsk",
			},
		}))
	}))
}

func TestUpdateConfig(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(UpdateShareConfig(&projectvo.UpdateShareConfigReq{
			OrgId:  3082,
			UserId: 34845,
			Input: &projectvo.UpdateData{
				ViewId: 1643095250803195908,
				Config: `{"key":"1232"}`,
			},
		}))
	}))
}

func TestResetShareKey(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(ResetShareKey(&projectvo.ResetShareKeyReq{
			OrgId:  3082,
			UserId: 34845,
			Input:  &projectvo.ShareViewIdData{ViewId: 1643095250803195908},
		}))
	}))
}

func TestCheckShareViewPassword(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(CheckShareViewPassword(&projectvo.CheckPasswordData{
			ShareKey: "TOXSeAaM",
			Password: "%IFDjfdsk",
		}))
	}))
}
