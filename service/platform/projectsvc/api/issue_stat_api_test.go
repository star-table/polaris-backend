package api

import (
	"context"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetGreeter_IssueAssignRank(t *testing.T) {

	convey.Convey("Test 任务负责数量rank统计", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
		if err != nil {
			log.Error(err)
			return
		}

		reqVo := projectvo.IssueAssignRankReqVo{
			Input: vo.IssueAssignRankReq{
				ProjectID: 1010,
			},
			OrgId: cacheUserInfo.OrgId,
		}

		respVo := postGreeter.IssueAssignRank(reqVo)
		t.Log(json.ToJsonIgnoreError(respVo))
		convey.So(respVo.Failure(), convey.ShouldEqual, false)
	}))
}
