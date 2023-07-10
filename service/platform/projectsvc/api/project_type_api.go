package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (GetGreeter) ProjectTypes(req projectvo.ProjectTypesReqVo) projectvo.ProjectTypesRespVo {
	resp, err := service.ProjectTypes(req.OrgId)
	return projectvo.ProjectTypesRespVo{Err: vo.NewErr(err), ProjectTypes: resp}
}
