package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/trendsvo"
	"github.com/star-table/polaris-backend/service/platform/trendssvc/service"
)

func (PostGreeter) TrendList(reqVo trendsvo.TrendListReqVo) trendsvo.TrendListRespVo {
	trendList, err := service.TrendList(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return trendsvo.TrendListRespVo{Err: vo.NewErr(err), TrendList: trendList}
}

func (PostGreeter) AddIssueTrends(req trendsvo.AddIssueTrendsReqVo) vo.VoidErr {
	service.AddIssueTrends(req.IssueTrendsBo)
	return vo.VoidErr{Err: vo.NewErr(nil)}
}

//添加项目趋势
func (PostGreeter) AddProjectTrends(req trendsvo.AddProjectTrendsReqVo) vo.VoidErr {
	service.AddProjectTrends(req.ProjectTrendsBo)
	return vo.VoidErr{Err: vo.NewErr(nil)}
}

func (PostGreeter) AddOrgTrends(req trendsvo.AddOrgTrendsReqVo) vo.VoidErr {
	service.AddOrgTrends(req.OrgTrendsBo)
	return vo.VoidErr{Err: vo.NewErr(nil)}
}

func (PostGreeter) CreateTrends(req trendsvo.CreateTrendsReqVo) trendsvo.CreateTrendsRespVo {
	trendsId, err := service.CreateTrends(req.TrendsBo)
	return trendsvo.CreateTrendsRespVo{TrendsId: trendsId, Err: vo.NewErr(err)}
}

func (PostGreeter) GetProjectLatestUpdateTime(req trendsvo.GetProjectLatestUpdateTimeReqVo) trendsvo.GetProjectLatestUpdateTimeRespVo {
	res, err := service.GetProjectLatestUpdateTime(req.OrgId, req.Input.ProjectIds)
	return trendsvo.GetProjectLatestUpdateTimeRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetIssueCommentCount(req trendsvo.GetIssueCommentCountReqVo) trendsvo.GetIssueCommentCountRespVo {
	res, err := service.GetIssueCommentCount(req.OrgId, req.Input.IssueIds)
	return trendsvo.GetIssueCommentCountRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetRecentChangeTime(req trendsvo.GetRecentChangeTimeReqVo) trendsvo.GetRecentChangeTimeRespVo {
	res, err := service.GetRecentChangeTime(req.Data.OrgIds)
	return trendsvo.GetRecentChangeTimeRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetLatestOperatedProjectId(req trendsvo.GetLatestOperatedProjectIdReq) trendsvo.GetLatestOperatedProjectIdResp {
	res, err := service.GetLatestOperatedProjectId(req.OrgId, req.UserId, req.ProjectIds)
	return trendsvo.GetLatestOperatedProjectIdResp{
		Err:       vo.NewErr(err),
		ProjectId: res,
	}
}

func (PostGreeter) DeleteTrends(req trendsvo.DeleteTrendsReq) vo.CommonRespVo {
	err := service.DeleteTrends(req)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: 0},
	}
}
