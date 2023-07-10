package api

import (
	"context"

	"github.com/star-table/common/core/util/validator"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
)

func (r *queryResolver) ProjectAttachment(ctx context.Context, page *int, size *int, params vo.ProjectAttachmentReq) (*vo.AttachmentList, error) {
	validate, err := validator.Validate(params)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	truePage := 0
	trueSize := 0
	if page != nil {
		truePage = *page
	}
	if size != nil {
		trueSize = *size
	}
	resp := projectfacade.GetProjectAttachment(projectvo.GetProjectAttachmentReqVo{
		Input:  params,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Page:   truePage,
		Size:   trueSize,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Output, resp.Error()
}

func (r *queryResolver) ProjectAttachmentInfo(ctx context.Context, input vo.ProjectAttachmentInfoReq) (*vo.Attachment, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.GetProjectAttachmentInfo(projectvo.GetProjectAttachmentInfoReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	if resp.Output == nil {
		return nil, nil
	}

	return resp.Output, resp.Error()
}

func (r *mutationResolver) DeleteProjectAttachment(ctx context.Context, input vo.DeleteProjectAttachmentReq) (*vo.DeleteProjectAttachmentResp, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.DeleteProjectAttachment(projectvo.DeleteProjectAttachmentReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Output, resp.Error()
}

func (r *queryResolver) FsDocumentList(ctx context.Context, page *int, size *int, input vo.FsDocumentListReq) (*vo.FsDocumentListResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	req := resourcevo.FsDocumentListReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Page:   0,
		Size:   0,
		Params: resourcevo.FsDocumentLisReqData{SearchKey: ""},
	}
	if page != nil {
		req.Page = *page
	}
	if size != nil {
		req.Size = *size
	}
	if input.SearchKey != nil {
		req.Params.SearchKey = *input.SearchKey
	}
	resp := resourcefacade.FsDocumentList(req)
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Data, resp.Error()
}
