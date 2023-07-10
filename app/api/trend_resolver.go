package api

import (
	"context"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/trendsvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/trendsfacade"
)

func (r *queryResolver) TrendList(ctx context.Context, input *vo.TrendReq) (*vo.TrendsList, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	req := trendsvo.TrendListReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	}

	resp := trendsfacade.TrendList(req)

	return resp.TrendList, resp.Error()

}
