package api

import (
	"context"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

func (r *queryResolver) Departments(ctx context.Context, page *int, size *int, params *vo.DepartmentListReq) (*vo.DepartmentList, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := orgfacade.Departments(orgvo.DepartmentsReqVo{
		Page:          page,
		Size:          size,
		Params:        params,
		CurrentUserId: cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
	})

	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.DepartmentList, nil
}

func (r *queryResolver) DepartmentMembers(ctx context.Context, params vo.DepartmentMemberListReq) ([]*vo.DepartmentMemberInfo, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := orgfacade.DepartmentMembers(orgvo.DepartmentMembersReqVo{
		Params:        params,
		CurrentUserId: cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.DepartmentMemberInfos, nil
}

func (r *queryResolver) DepartmentMembersList(ctx context.Context, page *int, size *int, params *vo.DepartmentMembersListReq) (*vo.DepartmentMembersListResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := orgfacade.DepartmentMembersList(orgvo.DepartmentMembersListReq{
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Page:          page,
		Size:          size,
		Params:        params,
		IgnoreDelete:  false,
	})

	return respVo.Data, respVo.Error()
}
