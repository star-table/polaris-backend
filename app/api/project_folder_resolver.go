package api

import (
	"context"
	"github.com/star-table/common/core/util/validator"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
)

func (r *mutationResolver) CreateProjectFolder(ctx context.Context, input vo.CreateProjectFolderReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.CreateProjectFolder(projectvo.CreateProjectFolderReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Void, resp.Error()
}

func (r *mutationResolver) DeleteProjectFolder(ctx context.Context, input vo.DeleteProjectFolderReq) (*vo.DeleteProjectFolderResp, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.DeleteProjectFolder(projectvo.DeleteProjectFolerReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Output, resp.Error()
}

func (r *mutationResolver) UpdateProjectFolder(ctx context.Context, input vo.UpdateProjectFolderReq) (*vo.UpdateProjectFolderResp, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.UpdateProjectFolder(projectvo.UpdateProjectFolderReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Output, resp.Error()
}

func (r *queryResolver) ProjectFolder(ctx context.Context, page *int, size *int, params vo.ProjectFolderReq) (*vo.FolderList, error) {
	validate, err := validator.Validate(params)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.GetProjectFolder(projectvo.GetProjectFolderReqVo{
		Input:  params,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Page:   *page,
		Size:   *size,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Output, resp.Error()
}
