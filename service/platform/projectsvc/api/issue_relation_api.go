package api

import (
	"strings"

	"github.com/star-table/common/core/util/copyer"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) CreateIssueComment(reqVo projectvo.CreateIssueCommentReqVo) vo.CommonRespVo {
	reqVo.Input.Comment = strings.TrimSpace(reqVo.Input.Comment)
	attachmentIds := reqVo.Input.AttachmentIds
	isCommentRight := format.VerifyIssueCommenFormat(reqVo.Input.Comment)
	if !isCommentRight && (attachmentIds == nil || len(attachmentIds) == 0) {
		log.Error(errs.IssueCommentLenError)
		return vo.CommonRespVo{Err: vo.NewErr(errs.IssueCommentLenError), Void: nil}
	}

	res, err := service.CreateIssueComment(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) CreateIssueResource(reqVo projectvo.CreateIssueResourceReqVo) vo.CommonRespVo {
	res, err := service.CreateIssueResource(reqVo.OrgId, reqVo.UserId, &reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

//func (PostGreeter) CreateIssueRelationIssue(reqVo projectvo.CreateIssueRelationIssueReqVo) vo.CommonRespVo {
//	res, err := service.CreateIssueRelationIssue(reqVo.OrgId, reqVo.UserId, reqVo.Input)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}

func (PostGreeter) DeleteIssueResource(reqVo projectvo.DeleteIssueResourceReqVo) vo.CommonRespVo {
	res, err := service.DeleteIssueResource(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

//func (PostGreeter) RelatedIssueList(reqVo projectvo.RelatedIssueListReqVo) projectvo.RelatedIssueListRespVo {
//	res, err := service.RelatedIssueList(reqVo.OrgId, reqVo.Input, reqVo.SourceChannel)
//	return projectvo.RelatedIssueListRespVo{Err: vo.NewErr(err), RelatedIssueList: res}
//}

func (PostGreeter) IssueResources(reqVo projectvo.IssueResourcesReqVo) projectvo.IssueResourcesRespVo {
	res, err := service.IssueResources(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.IssueResourcesRespVo{Err: vo.NewErr(err), IssueResources: res}
}

//func (PostGreeter) CreateIssueRelationTags(reqVo projectvo.CreateIssueRelationTagsReqVo) vo.CommonRespVo {
//	err := service.CreateIssueRelationTags(reqVo)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: &vo.Void{
//		ID: reqVo.Input.ID,
//	}}
//}

func (GetGreeter) GetIssueMembers(reqVo projectvo.GetIssueMembersReqVo) projectvo.GetIssueMembersRespVo {
	res, err := service.GetIssueMembers(reqVo.OrgId, reqVo.IssueId)
	return projectvo.GetIssueMembersRespVo{Err: vo.NewErr(err), Data: *res}
}

func (PostGreeter) GetIssueRelationResource(req projectvo.GetIssueRelationResourceReqVo) projectvo.GetIssueRelationResourceRespVo {
	res, err := service.GetIssueRelationResource(req.Page, req.Size)
	return projectvo.GetIssueRelationResourceRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) AddIssueAttachmentFs(req projectvo.AddIssueAttachmentFsReq) projectvo.AddIssueAttachmentFsRespVo {
	res, err := service.AddIssueAttachmentFs(req.OrgId, req.UserId, req.Input)
	var resources []*vo.Resource
	copyer.Copy(&res, &resources)
	return projectvo.AddIssueAttachmentFsRespVo{
		Err: vo.NewErr(err),
		Data: &vo.AddIssueAttachmentFsResp{
			Resources: resources,
		},
	}
}

func (PostGreeter) AddIssueAttachment(req projectvo.AddIssueAttachmentReq) vo.CommonRespVo {
	res, err := service.AddIssueAttachment(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
}

//func (PostGreeter) UpdateIssueBeforeAfterIssues(req projectvo.UpdateIssueBeforeAfterIssuesReq) vo.CommonRespVo {
//	res, err := service.UpdateIssueBeforeAfterIssues(req.OrgId, req.UserId, req.Input)
//	return vo.CommonRespVo{
//		Err:  vo.NewErr(err),
//		Void: res,
//	}
//}

//func (PostGreeter) BeforeAfterIssueList(req projectvo.BeforeAfterIssueListReq) projectvo.BeforeAfterIssueListResp {
//	res, err := service.BeforeAfterIssueList(req.OrgId, req.IssueId, req.SourceChannel)
//	return projectvo.BeforeAfterIssueListResp{
//		Err:  vo.NewErr(err),
//		NewData: res,
//	}
//}

// 分享任务卡片消息
func (PostGreeter) IssueShareCard(req projectvo.IssueCardShareReqVo) projectvo.IssueCardShareResp {
	resp, err := service.ShareIssueCard(req.OrgId, req.UserId, req.Input)
	return projectvo.IssueCardShareResp{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}
