package api

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) LcViewStatForAll(reqVo *projectvo.LcViewStatReqVo) *projectvo.LcViewStatRespVo {
	viewStats, err := service.LcViewStatForAll(reqVo.OrgId, reqVo.UserId)
	return &projectvo.LcViewStatRespVo{Err: vo.NewErr(err), Data: viewStats}
}

func (PostGreeter) LcHomeIssues(reqVo *projectvo.HomeIssuesReqVo) string {
	if reqVo.Size > 30000 {
		reqVo.Size = 30000
	}
	res, err := service.LcHomeIssuesForProject(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input, false)
	if err != nil {
		return json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"})
	}

	return res.Data
}

func (PostGreeter) LcHomeIssuesForAll(reqVo *projectvo.HomeIssuesReqVo) string {
	if reqVo.Size > 30000 {
		reqVo.Size = 30000
	}
	res, err := service.LcHomeIssuesForAll(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input, false)
	if err != nil {
		return json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"})
	}

	return res.Data
}

func (PostGreeter) LcHomeIssuesForIssue(reqVo *projectvo.IssueDetailReqVo) *projectvo.IssueDetailRespVo {
	res, err := service.LcHomeIssuesForIssue(reqVo.OrgId, reqVo.UserId, reqVo.AppId, reqVo.TableId, reqVo.IssueId, false)
	if err != nil {
		return &projectvo.IssueDetailRespVo{Err: vo.NewErr(err)}
	}
	return &projectvo.IssueDetailRespVo{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) IssueReportDetail(reqVo projectvo.IssueReportDetailReqVo) projectvo.IssueReportDetailRespVo {
	res, err := service.IssueReportDetail(reqVo.ShareID)
	return projectvo.IssueReportDetailRespVo{Err: vo.NewErr(err), IssueReportDetail: res}
}

func (PostGreeter) IssueStatusTypeStat(reqVo projectvo.IssueStatusTypeStatReqVo) projectvo.IssueStatusTypeStatRespVo {
	res, err := service.IssueStatusTypeStat(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.IssueStatusTypeStatRespVo{Err: vo.NewErr(err), IssueStatusTypeStat: res}
}

func (PostGreeter) IssueStatusTypeStatDetail(reqVo projectvo.IssueStatusTypeStatReqVo) projectvo.IssueStatusTypeStatDetailRespVo {
	res, err := service.IssueStatusTypeStatDetail(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.IssueStatusTypeStatDetailRespVo{Err: vo.NewErr(err), IssueStatusTypeStatDetail: res}
}

//func (PostGreeter) GetSimpleIssueInfoBatch(reqVo projectvo.GetSimpleIssueInfoBatchReqVo) projectvo.GetSimpleIssueInfoBatchRespVo {
//	res, err := service.GetSimpleIssueInfoBatch(reqVo.OrgId, reqVo.Ids)
//	return projectvo.GetSimpleIssueInfoBatchRespVo{Err: vo.NewErr(err), Data: res}
//}

func (PostGreeter) GetLcIssueInfoBatch(reqVo projectvo.GetLcIssueInfoBatchReqVo) projectvo.GetLcIssueInfoBatchRespVo {
	res, err := service.GetLcIssueInfoBatch(reqVo.OrgId, reqVo.IssueIds)
	return projectvo.GetLcIssueInfoBatchRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

//func (PostGreeter) GetIssueRemindInfoList(reqVo projectvo.GetIssueRemindInfoListReqVo) projectvo.GetIssueRemindInfoListRespVo {
//	res, err := service.GetIssueRemindInfoList(reqVo)
//	return projectvo.GetIssueRemindInfoListRespVo{Err: vo.NewErr(err), Data: res}
//}

func (PostGreeter) IssueListStat(reqVo projectvo.IssueListStatReq) projectvo.IssueListStatResp {
	res, err := service.IssueListStat(reqVo.OrgId, reqVo.UserId, reqVo.Input.ProjectID)
	return projectvo.IssueListStatResp{Err: vo.NewErr(err), Data: res}
}

//func (PostGreeter) HomeIssuesGroup(reqVo projectvo.HomeIssuesGroupReqVo) projectvo.HomeIssuesGroupRespVo {
//	res, err := service.HomeIssuesGroup(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
//	return projectvo.HomeIssuesGroupRespVo{Err: vo.NewErr(err), HomeIssueInfo: res}
//}

func (PostGreeter) IssueListSimpleByDataIds(reqVo projectvo.GetIssueListSimpleByDataIdsReqVo) projectvo.SimpleIssueListRespVo {
	res, err := service.IssueListSimpleByDataIds(reqVo.OrgId, reqVo.UserId, reqVo.Input.AppId, reqVo.Input.DataIds)
	return projectvo.SimpleIssueListRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) IssueListSimpleByTableIds(req projectvo.GetIssueListWithConditionsReqVo) projectvo.IssueListWithConditionsResp {
	res, err := service.IssueListSimpleByTableIds(req.OrgId, req.UserId, req.Input)
	return projectvo.IssueListWithConditionsResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

//func (PostGreeter) IssueFilterIncremental(req projectvo.IssueFilterIncrementalReq) vo.DataRespVo {
//	return vo.DataRespVo{Err: vo.NewErr(nil), Data: nil}
//}
