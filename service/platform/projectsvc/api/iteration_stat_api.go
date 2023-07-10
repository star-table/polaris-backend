package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) IterationStats(reqVo projectvo.IterationStatsReqVo) projectvo.IterationStatsRespVo {
	res, err := service.IterationStats(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.IterationStatsRespVo{Err: vo.NewErr(err), IterationStats: res}
}

func (PostGreeter) AppendIterationStat(req projectvo.AppendIterationStatReqVo) vo.VoidErr {
	err := service.AppendIterationStat(req.IterationBo, req.Date)
	return vo.VoidErr{Err: vo.NewErr(err)}
}
