package api

import (
	"context"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

func (r *mutationResolver) UpdateOrgMemberStatus(ctx context.Context, input vo.UpdateOrgMemberStatusReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := orgfacade.UpdateOrgMemberStatus(orgvo.UpdateOrgMemberStatusReq{
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Input:         input,
	})
	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) UpdateOrgMemberCheckStatus(ctx context.Context, input vo.UpdateOrgMemberCheckStatusReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := orgfacade.UpdateOrgMemberCheckStatus(orgvo.UpdateOrgMemberCheckStatusReq{
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Input:         input,
	})
	return respVo.Void, respVo.Error()
}

func (r *queryResolver) OrgUserList(ctx context.Context, page *int, size *int, input vo.OrgUserListReq) (*vo.UserOrganizationList, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	defaultPage := 1
	defaultSize := 20
	if page == nil || *page < 1 {
		page = &defaultPage
	}

	if size == nil || *size <= 0 {
		size = &defaultSize
	}

	resp := orgfacade.OrgUserList(orgvo.OrgUserListReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Page:   *page,
		Size:   *size,
		Input:  &input,
	})
	return resp.Data, resp.Error()
}
