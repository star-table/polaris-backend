package api

import (
	"context"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
)

func (r *queryResolver) ProjectDayStats(ctx context.Context, page *int, size *int, params *vo.ProjectDayStatReq) (*vo.ProjectDayStatList, error) {
	pageA := uint(0)
	sizeA := uint(0)
	if page != nil && size != nil && *page > 0 && *size > 0 {
		pageA = uint(*page)
		sizeA = uint(*size)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.ProjectDayStats(&projectvo.ProjectDayStatsReqVo{
		Page:   pageA,
		Size:   sizeA,
		Params: params,
		OrgId:  cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.ProjectDayStatList, nil
}

func (r *queryResolver) PayLimitNum(ctx context.Context) (*vo.PayLimitNumResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.PayLimitNum(projectvo.PayLimitNumReq{
		OrgId: cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.Data, nil
}
