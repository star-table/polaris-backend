package api

import (
	"context"

	"github.com/star-table/common/core/util/validator"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/extra/gin/util"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
)

func (r *queryResolver) Projects(ctx context.Context, page int, size int, params map[string]interface{}, order []*string, input *vo.ProjectsReq) (*vo.ProjectList, error) {
	maxPageSize := 100
	if size > maxPageSize {
		return nil, errs.PageSizeOverflowMaxSizeError
	}

	err := validator.ValidateConds("Projects", &params)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.OutOfConditionError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.Projects(projectvo.ProjectsRepVo{
		Page: page,
		Size: size,
		ProjectExtraBody: projectvo.ProjectExtraBody{
			Params: params,
			Order:  order,
			Input:  input,
		},
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.ProjectList, nil
}

func (r *mutationResolver) CreateProject(ctx context.Context, input vo.CreateProjectReq) (*vo.Project, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	req := projectvo.CreateProjectReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	}
	version, err1 := util.GetCtxVersion(ctx)
	if err != nil {
		//打日志
		log.Error(err1)
	}

	if version != "" {
		req.Version = version
	}

	resp := projectfacade.CreateProject(req)
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.Project, nil
}

func (r *mutationResolver) UpdateProject(ctx context.Context, input vo.UpdateProjectReq) (*vo.Project, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.UpdateProject(projectvo.UpdateProjectReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.Project, nil
}

func (r *mutationResolver) UpdateProjectStatus(ctx context.Context, input vo.UpdateProjectStatusReq) (*vo.Void, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.UpdateProjectStatus(projectvo.UpdateProjectStatusReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})

	return resp.Void, resp.Error()
}

func (r *queryResolver) ProjectInfo(ctx context.Context, input vo.ProjectInfoReq) (*vo.ProjectInfo, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.ProjectInfo(projectvo.ProjectInfoReqVo{
		Input:         input,
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})

	return resp.ProjectInfo, resp.Error()
}

func (r *mutationResolver) ArchiveProject(ctx context.Context, projectID int64) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.ArchiveProject(projectvo.ProjectIdReqVo{
		OrgId:     cacheUserInfo.OrgId,
		UserId:    cacheUserInfo.UserId,
		ProjectId: projectID,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) CancelArchivedProject(ctx context.Context, projectID int64) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.CancelArchivedProject(projectvo.ProjectIdReqVo{
		OrgId:     cacheUserInfo.OrgId,
		UserId:    cacheUserInfo.UserId,
		ProjectId: projectID,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) DeleteProject(ctx context.Context, projectID int64) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.DeleteProject(projectvo.ProjectIdReqVo{
		OrgId:     cacheUserInfo.OrgId,
		UserId:    cacheUserInfo.UserId,
		ProjectId: projectID,
	})
	return resp.Void, resp.Error()
}
