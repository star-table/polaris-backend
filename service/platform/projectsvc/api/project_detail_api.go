package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (GetGreeter) ProjectDetail(reqVo projectvo.ProjectDetailReqVo) projectvo.ProjectDetailRespVo {
	res, err := service.ProjectDetail(reqVo.OrgId, reqVo.ProjectId)
	return projectvo.ProjectDetailRespVo{Err: vo.NewErr(err), ProjectDetail: res}
}

func (PostGreeter) CreateProjectDetail(reqVo projectvo.CreateProjectDetailReqVo) vo.CommonRespVo {
	res, err := service.CreateProjectDetail(reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateProjectDetail(reqVo projectvo.UpdateProjectDetailReqVo) vo.CommonRespVo {
	res, err := service.UpdateProjectDetail(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) DeleteProjectDetail(reqVo projectvo.DeleteProjectDetailReqVo) vo.CommonRespVo {
	res, err := service.DeleteProjectDetail(reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}
