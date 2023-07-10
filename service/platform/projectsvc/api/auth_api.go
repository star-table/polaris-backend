package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) AuthProjectPermission(reqVo projectvo.AuthProjectPermissionReqVo) vo.CommonRespVo {
	input := reqVo.Input
	err := service.AuthProjectPermission(input.OrgId, input.UserId, input.ProjectId, input.Path, input.Operation, input.AuthFiling)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: nil}
}
