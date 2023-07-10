package api

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) ProjectDayStats(params *projectvo.ProjectDayStatsReqVo) projectvo.ProjectDayStatsRespVo {
	list, err := service.ProjectDayStats(params.OrgId, params.Page, params.Size, params.Params)
	return projectvo.ProjectDayStatsRespVo{Err: vo.NewErr(err), ProjectDayStatList: list}
}

//date: yyyy-MM-dd
func (PostGreeter) AppendProjectDayStat(req projectvo.AppendProjectDayStatReqVo) vo.VoidErr {
	err := service.AppendProjectDayStat(req.ProjectBo, req.Date)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) PayLimitNum(req projectvo.PayLimitNumReq) projectvo.PayLimitNumResp {
	res, err := service.PayLimitNum(req.OrgId)
	respData := vo.PayLimitNumResp{}
	copyer.Copy(res, &respData)

	return projectvo.PayLimitNumResp{
		Err:  vo.NewErr(err),
		Data: &respData,
	}
}

func (PostGreeter) PayLimitNumForRest(req projectvo.PayLimitNumReq) projectvo.PayLimitNumForRestResp {
	res, err := service.PayLimitNum(req.OrgId)

	return projectvo.PayLimitNumForRestResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetProjectStatistics(req projectvo.GetProjectStatisticsReqVo) projectvo.GetProjectStatisticsResp {
	res, err := service.GetProjectStatistics(req.OrgId, req.UserId, req.Input)
	return projectvo.GetProjectStatisticsResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
