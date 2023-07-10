package api

import (
	"context"
	"strconv"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
)

func (r *mutationResolver) DeleteIssue(ctx context.Context, input vo.DeleteIssueReq) (*vo.Issue, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.DeleteIssue(projectvo.DeleteIssueReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	return respVo.Issue, respVo.Error()
}

func (r *queryResolver) IssueInfoNotDelete(ctx context.Context, param vo.IssueInfoNotDeleteReq) (*vo.IssueInfo, error) {
	return nil, nil
}

func (r *mutationResolver) MoveIssue(ctx context.Context, input vo.MoveIssueReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.MoveIssue(projectvo.MoveIssueReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})

	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) MoveIssueBatch(ctx context.Context, input vo.MoveIssueBatchReq) (*vo.MoveIssueBatchResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.MoveIssueBatch(projectvo.MoveIssueBatchReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})

	return resp.Data, resp.Error()
}

func (r *queryResolver) IssueStatusTypeStat(ctx context.Context, input *vo.IssueStatusTypeStatReq) (*vo.IssueStatusTypeStatResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IssueStatusTypeStat(projectvo.IssueStatusTypeStatReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.IssueStatusTypeStat, respVo.Error()
}

func (r *queryResolver) IssueAssignRank(ctx context.Context, input vo.IssueAssignRankReq) ([]*vo.IssueAssignRankInfo, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	reqVo := projectvo.IssueAssignRankReqVo{
		Input: input,
		OrgId: cacheUserInfo.OrgId,
	}

	respVo := projectfacade.IssueAssignRank(reqVo)
	return respVo.IssueAssignRankResp, respVo.Error()
}

func (r *queryResolver) IssueStatusTypeStatDetail(ctx context.Context, input *vo.IssueStatusTypeStatReq) (*vo.IssueStatusTypeStatDetailResp, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IssueStatusTypeStatDetail(projectvo.IssueStatusTypeStatReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.IssueStatusTypeStatDetail, respVo.Error()
}

func (r *mutationResolver) CreateIssueComment(ctx context.Context, input vo.CreateIssueCommentReq) (*vo.Void, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.CreateIssueComment(projectvo.CreateIssueCommentReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.Void, respVo.Error()
}

func (r *queryResolver) ExportIssueTemplate(ctx context.Context, projectID int64, tableIdStr string) (*vo.ExportIssueTemplateResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Errorf("[ExportIssueTemplate] err: %v", err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	tableId, oriErr := strconv.ParseInt(tableIdStr, 10, 64)
	if oriErr != nil {
		log.Errorf("[ExportIssueTemplate] parse tableIdStr err: %v", oriErr)
		return nil, errs.BuildSystemErrorInfo(errs.TypeConvertError, oriErr)
	}
	objectTypeId := tableId
	respVo := projectfacade.ExportIssueTemplate(projectvo.ExportIssueTemplateReqVo{
		OrgId:               cacheUserInfo.OrgId,
		ProjectId:           projectID,
		ProjectObjectTypeId: objectTypeId,
		TableId:             tableId,
	})

	return respVo.Data, respVo.Error()
}

func (r *queryResolver) ExportData(ctx context.Context, projectID int64, iterationID *int64, tableIdStr string, isNeedDocument *bool) (*vo.ExportIssueTemplateResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	iterationId := int64(0)
	if iterationID != nil {
		iterationId = *iterationID
	}
	req := projectvo.ExportIssueReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input: projectvo.ExportIssueReqVoData{
			ProjectId:   projectID,
			TableId:     tableIdStr,
			ViewId:      "",
			IterationId: iterationId,
		},
	}
	if isNeedDocument != nil {
		req.Input.IsNeedDocument = *isNeedDocument
	}
	respVo := projectfacade.ExportData(req)

	return respVo.Data, respVo.Error()
}

func (r *mutationResolver) UpdateIssueSort(ctx context.Context, input vo.UpdateIssueSortReq) (*vo.Void, error) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//
	//respVo := projectfacade.UpdateIssueSort(projectvo.UpdateIssueSortReqVo{
	//	Input:         input,
	//	OrgId:         cacheUserInfo.OrgId,
	//	UserId:        cacheUserInfo.UserId,
	//	SourceChannel: cacheUserInfo.SourceChannel,
	//})
	//return respVo.Void, respVo.Error()
	return nil, nil
}

func (r *mutationResolver) CopyIssue(ctx context.Context, input vo.CopyIssueReq) (*vo.Issue, error) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//
	//resp := projectfacade.CopyIssue(projectvo.CopyIssueReqVo{
	//	Input:         input,
	//	UserId:        cacheUserInfo.UserId,
	//	OrgId:         cacheUserInfo.OrgId,
	//	SourceChannel: cacheUserInfo.SourceChannel,
	//})
	//
	//return &vo.Issue{}, resp.Error()
	return nil, nil
}

func (r *mutationResolver) CopyIssueBatch(ctx context.Context, input vo.CopyIssueBatchReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return &vo.Void{ID: 0}, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.CopyIssueBatch(projectvo.CopyIssueBatchReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})

	return &vo.Void{ID: resp.Data}, resp.Error()
}

func (r *mutationResolver) DeleteIssueBatch(ctx context.Context, input vo.DeleteIssueBatchReq) (*vo.DeleteIssueBatchResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	req := projectvo.DeleteIssueBatchReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	}
	if input.MenuAppID != nil && *input.MenuAppID != "" {
		menuAppId, err := strconv.ParseInt(*input.MenuAppID, 10, 64)
		if err != nil {
			log.Error(err)
			return nil, errs.JSONConvertError
		}
		req.InputAppId = menuAppId
	}
	resp := projectfacade.DeleteIssueBatch(req)

	return resp.Data, resp.Error()
}

func (r *mutationResolver) ConvertIssueToParent(ctx context.Context, input vo.ConvertIssueToParentReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.ConvertIssueToParent(projectvo.ConvertIssueToParentReq{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input:  input,
	})

	return resp.Void, resp.Error()
}

func (r *mutationResolver) ChangeParentIssue(ctx context.Context, input vo.ChangeParentIssueReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.ChangeParentIssue(projectvo.ChangeParentIssueReq{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input:  input,
	})

	return resp.Void, resp.Error()
}

// 任务催促
func (r *mutationResolver) UrgeIssue(ctx context.Context, params vo.UrgeIssueReq) (*vo.BoolResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.UrgeIssue(projectvo.UrgeIssueReqVo{
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Input:         &params,
	})
	return &vo.BoolResp{
		IsTrue: resp.IsTrue,
	}, resp.Error()
}

func (r *mutationResolver) AuditIssue(ctx context.Context, params vo.AuditIssueReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.AuditIssue(projectvo.AuditIssueReq{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Params: params,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) ViewAuditIssue(ctx context.Context, params vo.ViewAuditIssueReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.ViewAuditIssue(projectvo.ViewAuditIssueReq{
		UserId:  cacheUserInfo.UserId,
		OrgId:   cacheUserInfo.OrgId,
		IssueId: params.IssueID,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) WithdrawIssue(ctx context.Context, params vo.WithdrawIssueReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.WithdrawIssue(projectvo.WithdrawIssueReq{
		UserId:  cacheUserInfo.UserId,
		OrgId:   cacheUserInfo.OrgId,
		IssueId: params.IssueID,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) UrgeAuditIssue(ctx context.Context, params vo.UrgeAuditIssueReq) (*vo.BoolResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.UrgeAuditIssue(projectvo.UrgeAuditIssueReq{
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Input:         params,
	})
	return &vo.BoolResp{
		IsTrue: resp.IsTrue,
	}, resp.Error()
}
