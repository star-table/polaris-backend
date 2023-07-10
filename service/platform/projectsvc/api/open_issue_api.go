package api

import (
	"strings"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) OpenLcHomeIssues(reqVo *projectvo.HomeIssuesReqVo) string {
	if reqVo.Size > 30000 {
		reqVo.Size = 30000
	}
	res, err := service.LcHomeIssuesForProject(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input, true)
	if err != nil {
		return json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"})
	}

	return res.Data
}

func (PostGreeter) OpenLcHomeIssuesForAll(reqVo *projectvo.HomeIssuesReqVo) string {
	if reqVo.Size > 30000 {
		reqVo.Size = 30000
	}
	res, err := service.LcHomeIssuesForAll(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input, true)
	if err != nil {
		return json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"})
	}

	return res.Data
}

func (GetGreeter) OpenIssueInfo(reqVo projectvo.IssueInfoReqVo) *projectvo.IssueDetailRespVo {
	data, err := service.LcHomeIssuesForIssue(reqVo.OrgId, reqVo.UserId, -1, -1, reqVo.IssueID, true)
	return &projectvo.IssueDetailRespVo{Err: vo.NewErr(err), Data: data}
}

func (PostGreeter) OpenIssueStatusTypeStat(reqVo projectvo.IssueStatusTypeStatReqVo) projectvo.IssueStatusTypeStatRespVo {
	res, err := service.IssueStatusTypeStat(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.IssueStatusTypeStatRespVo{Err: vo.NewErr(err), IssueStatusTypeStat: res}
}

func (PostGreeter) OpenCreateIssue(reqVo projectvo.CreateIssueReqVo) projectvo.IssueRespVo {
	return projectvo.IssueRespVo{}
}

func (PostGreeter) OpenUpdateIssue(reqVo projectvo.OpenUpdateIssueReqVo) projectvo.UpdateIssueRespVo {
	return projectvo.UpdateIssueRespVo{Err: vo.NewErr(nil), UpdateIssue: &vo.UpdateIssueResp{}}
}

func (PostGreeter) OpenDeleteIssue(reqVo projectvo.OpenDeleteIssueReqVo) projectvo.DeleteIssueRespVo {
	// 校验操作人是否存在
	_, userErr := orgfacade.GetBaseUserInfoRelaxed(reqVo.OrgId, reqVo.UserId)
	if userErr != nil {
		log.Error(userErr)
		return projectvo.DeleteIssueRespVo{Err: vo.NewErr(userErr)}
	}

	judgeErr := service.JudgeProjectFilingByIssueId(reqVo.OrgId, reqVo.Data.IssueId)
	if judgeErr != nil {
		log.Error(judgeErr)
		return projectvo.DeleteIssueRespVo{Err: vo.NewErr(judgeErr)}
	}
	_, err := service.DeleteIssueWithoutAuth(projectvo.DeleteIssueReqVo{
		Input: vo.DeleteIssueReq{
			ID: reqVo.Data.IssueId,
		},
		OrgId:  reqVo.OrgId,
		UserId: reqVo.UserId,
	})
	if err != nil {
		return projectvo.DeleteIssueRespVo{Err: vo.NewErr(err)}
	}
	return projectvo.DeleteIssueRespVo{}
}

func (PostGreeter) OpenIssueList(reqVo projectvo.OpenIssueListReqVo) projectvo.LcHomeIssuesRespVo {
	res, err := service.LcHomeIssuesForProject(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Data, false)
	if err != nil {
		return projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"}
	}
	return projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: json.ToJsonIgnoreError(map[string]interface{}{
		"code":    0,
		"success": true,
		"data":    res,
	})}
}

func (PostGreeter) OpenCreateIssueComment(reqVo projectvo.CreateIssueCommentReqVo) vo.CommonRespVo {
	reqVo.Input.Comment = strings.TrimSpace(reqVo.Input.Comment)
	isCommentRight := format.VerifyIssueCommenFormat(reqVo.Input.Comment)
	if !isCommentRight {
		log.Error(errs.IssueCommentLenError)
		return vo.CommonRespVo{Err: vo.NewErr(errs.IssueCommentLenError), Void: nil}
	}

	res, err := service.CreateIssueComment(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}
