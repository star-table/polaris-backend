package api

import (
	"context"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
)

func (r *queryResolver) RecycleBinList(ctx context.Context, page *int, size *int, params vo.RecycleBinListReq) (*vo.RecycleBinList, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	defaultPage := 0
	defaultSize := 5
	if page == nil || *page < 0 {
		page = &defaultPage
	}

	if size == nil || *size <= 0 {
		size = &defaultSize
	}

	resp := projectfacade.GetRecycleList(projectvo.GetRecycleListReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Page:   *page,
		Size:   *size,
		Input: vo.RecycleBinListReq{
			ProjectID:    params.ProjectID,
			RelationType: params.RelationType,
		},
	})

	return resp.Data, resp.Error()
}

func (r *mutationResolver) RecoverRecycleBinRecord(ctx context.Context, input vo.RecoverRecycleBinRecordReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.RecoverRecycleBin(projectvo.CompleteDeleteReqVo{
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Input: vo.RecoverRecycleBinRecordReq{
			ProjectID:    input.ProjectID,
			RelationType: input.RelationType,
			RelationID:   input.RelationID,
			RecycleID:    input.RecycleID,
		},
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) CompleteDelete(ctx context.Context, input vo.RecoverRecycleBinRecordReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.CompleteDelete(projectvo.CompleteDeleteReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input: vo.RecoverRecycleBinRecordReq{
			ProjectID:    input.ProjectID,
			RelationType: input.RelationType,
			RelationID:   input.RelationID,
			RecycleID:    input.RecycleID,
		},
	})
	return resp.Void, resp.Error()
}
