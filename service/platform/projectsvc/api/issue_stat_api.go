package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) IssueAssignRank(reqVo projectvo.IssueAssignRankReqVo) projectvo.IssueAssignRankRespVo {
	input := reqVo.Input
	projectId := input.ProjectID
	rankTop := 5
	if input.RankTop != nil {
		rt := *input.RankTop
		if rt >= 1 && rt <= 100 {
			rankTop = rt
		}
	}
	res, err := service.IssueAssignRank(reqVo.OrgId, projectId, rankTop)
	return projectvo.IssueAssignRankRespVo{Err: vo.NewErr(err), IssueAssignRankResp: res}
}

func (GetGreeter) IssueAndProjectCountStat(reqVo projectvo.IssueAndProjectCountStatReqVo) projectvo.IssueAndProjectCountStatRespVo {
	res, err := service.IssueAndProjectCountStat(reqVo)
	return projectvo.IssueAndProjectCountStatRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) IssueDailyPersonalWorkCompletionStat(reqVo projectvo.IssueDailyPersonalWorkCompletionStatReqVo) projectvo.IssueDailyPersonalWorkCompletionStatRespVo {
	res, err := service.IssueDailyPersonalWorkCompletionStat(reqVo)
	return projectvo.IssueDailyPersonalWorkCompletionStatRespVo{Err: vo.NewErr(err), Data: res}
}

//func (PostGreeter) GetIssueCountByStatus(req projectvo.GetIssueCountByStatusReqVo) projectvo.GetIssueCountByStatusRespVo {
//	res, err := service.GetIssueCountByStatus(req.OrgId, req.ProjectId, req.StatusId)
//	return projectvo.GetIssueCountByStatusRespVo{
//		Err: vo.NewErr(err),
//		Data: projectvo.GetIssueCountByStatusData{
//			Count: res,
//		},
//	}
//}

func (PostGreeter) AuthProject(req projectvo.AuthProjectReqVo) vo.CommonRespVo {
	err := service.AuthProject(req.OrgId, req.UserId, req.ProjectId, req.Path, req.Operation)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}
