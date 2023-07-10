package service

import (
	"context"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/trendssvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUnreadNoticeCount(t *testing.T) {
	convey.Convey("UnreadNoticeCount", t, test.StartUp(func(ctx context.Context) {
		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		log.Info("缓存用户信息" + cacheUserInfoJson)
		count, err := UnreadNoticeCount(cacheUserInfo.OrgId, cacheUserInfo.UserId)
		t.Log(count, err)
		//assert.Assert(t, count, uint64(0))

		t.Log(NoticeList(cacheUserInfo.OrgId, cacheUserInfo.UserId, 0, 0, nil))
	}))
}
