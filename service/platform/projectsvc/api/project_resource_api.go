package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) CreateProjectResource(reqVo projectvo.CreateProjectResourceReqVo) vo.CommonRespVo {
	res, err := service.CreateProjectResource(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateProjectResourceName(reqVo projectvo.UpdateProjectResourceNameReqVo) vo.CommonRespVo {
	res, err := service.UpdateProjectResourceName(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateProjectFileResource(reqVo projectvo.UpdateProjectFileResourceReqVo) vo.CommonRespVo {
	res, err := service.UpdateProjectFileResource(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateProjectResourceFolder(reqVo projectvo.UpdateProjectResourceFolderReqVo) projectvo.UpdateProjectResourceFolderRespVo {
	res, err := service.UpdateProjectResourceFolder(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.UpdateProjectResourceFolderRespVo{Err: vo.NewErr(err), Output: res}
}

func (PostGreeter) DeleteProjectResource(reqVo projectvo.DeleteProjectResourceReqVo) projectvo.DeleteProjectResourceRespVo {
	res, err := service.DeleteProjectResource(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.DeleteProjectResourceRespVo{Err: vo.NewErr(err), Output: res}
}

func (PostGreeter) GetProjectResource(reqVo projectvo.GetProjectResourceReqVo) projectvo.GetProjectResourceResVo {
	res, err := service.GetProjectResource(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.GetProjectResourceResVo{Err: vo.NewErr(err), Output: res}
}

func (PostGreeter) GetProjectResourceInfo(reqVo projectvo.GetProjectResourceInfoReqVo) projectvo.GetProjectResourceInfoRespVo {
	res, err := service.GetProjectResourceInfo(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.GetProjectResourceInfoRespVo{Err: vo.NewErr(err), Output: res}
}
