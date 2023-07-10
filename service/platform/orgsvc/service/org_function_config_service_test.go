package service

import (
	"testing"

	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestUpdateOrgFunctionConfig(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		//payTime, err := strconv.ParseInt("1598862827", 10, 64)
		//if err != nil {
		//	log.Error(err)
		//}
		//paidTime := time.Unix(payTime, 0)
		//t.Log(UpdateOrgFunctionConfig(1303, sdk_const.SourceChannelFeishu, 2, "buy", "per_seat_per_month", paidTime))
	}))

}
