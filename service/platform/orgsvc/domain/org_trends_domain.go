package domain

import (
	"time"

	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/trendsvo"
	"github.com/star-table/polaris-backend/facade/trendsfacade"
)

func PushOrgTrends(orgTrendsBo bo.OrgTrendsBo) {
	orgTrendsBo.OperateTime = time.Now()
	//动态改成同步的
	resp := trendsfacade.AddOrgTrends(trendsvo.AddOrgTrendsReqVo{OrgTrendsBo: orgTrendsBo})
	if resp.Failure() {
		log.Error(resp.Message)
	}
}
