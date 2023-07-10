package api

import (
	"context"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

func (r *queryResolver) PermissionOperationList(ctx context.Context, roleID int64, projectID *int64) ([]*vo.PermissionOperationListResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := orgfacade.PermissionOperationList(rolevo.PermissionOperationListReqVo{
		OrgId:     cacheUserInfo.OrgId,
		RoleId:    roleID,
		UserId:    cacheUserInfo.UserId,
		ProjectId: projectID,
	})

	return resp.Data, resp.Error()
}

func (r *queryResolver) GetPersonalPermissionInfo(ctx context.Context, projectID *int64, issueID *int64) (*vo.GetPersonalPermissionInfoResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := orgfacade.GetPersonalPermissionInfo(rolevo.GetPersonalPermissionInfoReqVo{
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		ProjectId:     projectID,
		IssueId:       issueID,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	return &vo.GetPersonalPermissionInfoResp{Data: resp.Data}, resp.Error()
}

func (r *mutationResolver) ApplyScopes(ctx context.Context) (*vo.ApplyScopesResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	output := &vo.ApplyScopesResp{}
	resp := orgfacade.ApplyScopes(orgvo.ApplyScopesReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if resp.Successful() {
		output.ThirdCode = resp.Data.ThirdCode
		output.ThirdMsg = resp.Data.ThirdMsg
	}
	return output, resp.Error()
}

func (r *queryResolver) CheckSpecificScope(ctx context.Context, params vo.CheckSpecificScopeReq) (*vo.CheckSpecificScopeResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	result := &vo.CheckSpecificScopeResp{
		HasPower: false,
	}
	resp := orgfacade.CheckSpecificScope(orgvo.CheckSpecificScopeReqVo{
		OrgId:     cacheUserInfo.OrgId,
		UserId:    cacheUserInfo.UserId,
		PowerFlag: params.PowerFlag,
	})
	if resp.Failure() {
		log.Errorf("[CheckSpecificScope] CheckSpecificScope err: %v", resp.Error())
	} else {
		result.HasPower = resp.Data.HasPower
	}

	return result, resp.Error()
}
