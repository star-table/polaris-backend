package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) IterationList(reqVo projectvo.IterationListReqVo) projectvo.IterationListRespVo {
	res, err := service.IterationList(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.IterationListRespVo{Err: vo.NewErr(err), IterationList: res}
}

func (PostGreeter) CreateIteration(reqVo projectvo.CreateIterationReqVo) vo.CommonRespVo {
	res, err := service.CreateIteration(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateIteration(reqVo projectvo.UpdateIterationReqVo) vo.CommonRespVo {
	res, err := service.UpdateIteration(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) DeleteIteration(reqVo projectvo.DeleteIterationReqVo) vo.CommonRespVo {
	res, err := service.DeleteIteration(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) IterationStatusTypeStat(reqVo projectvo.IterationStatusTypeStatReqVo) projectvo.IterationStatusTypeStatRespVo {
	res, err := service.IterationStatusTypeStat(reqVo.OrgId, reqVo.Input)
	return projectvo.IterationStatusTypeStatRespVo{Err: vo.NewErr(err), IterationStatusTypeStat: res}
}

//func (PostGreeter) IterationIssueRelate(reqVo projectvo.IterationIssueRelateReqVo) vo.CommonRespVo {
//	res, err := service.IterationIssueRelate(reqVo.OrgId, reqVo.UserId, reqVo.Input)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}

func (PostGreeter) UpdateIterationStatus(reqVo projectvo.UpdateIterationStatusReqVo) vo.CommonRespVo {
	res, err := service.UpdateIterationStatus(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) IterationInfo(reqVo projectvo.IterationInfoReqVo) projectvo.IterationInfoRespVo {
	res, err := service.IterationInfo(reqVo.OrgId, reqVo.Input)
	return projectvo.IterationInfoRespVo{Err: vo.NewErr(err), IterationInfo: res}
}

//获取未完成的的迭代列表
func (GetGreeter) GetNotCompletedIterationBoList(req projectvo.GetNotCompletedIterationBoListReqVo) projectvo.GetNotCompletedIterationBoListRespVo {
	res, err := service.GetNotCompletedIterationBoList(req.OrgId, req.ProjectId)
	return projectvo.GetNotCompletedIterationBoListRespVo{Err: vo.NewErr(err), IterationBoList: res}
}

func (PostGreeter) UpdateIterationStatusTime(req projectvo.UpdateIterationStatusTimeReqVo) vo.CommonRespVo {
	res, err := service.UpdateIterationStatusTime(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) GetSimpleIterationInfo(req projectvo.GetSimpleIterationInfoReqVo) projectvo.GetSimpleIterationInfoRespVo {
	res, err := service.GetSimpleIterationInfo(req.OrgId, req.IterationIds)
	return projectvo.GetSimpleIterationInfoRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) UpdateIterationSort(req projectvo.UpdateIterationSortReqVo) vo.CommonRespVo {
	res, err := service.UpdateIterationSort(req.OrgId, req.UserId, req.Params)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
}
