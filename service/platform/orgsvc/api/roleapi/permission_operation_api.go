package roleapi

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	service "github.com/star-table/polaris-backend/service/platform/orgsvc/service/roleservice"
)

// TODO 替换为lesscode-permission
func (GetGreeter) PermissionOperationList(req rolevo.PermissionOperationListReqVo) rolevo.PermissionOperationListRespVo {
	res, err := service.PermissionOperationList(req.OrgId, req.RoleId, req.UserId, req.ProjectId)
	return rolevo.PermissionOperationListRespVo{Err: vo.NewErr(err), Data: res}
}

// TODO 替换为lesscode-permission
func (PostGreeter) UpdateRolePermissionOperation(req rolevo.UpdateRolePermissionOperationReqVo) vo.CommonRespVo {
	res, err := service.UpdateRolePermissionOperation(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

// TODO 替换为lesscode-permission
func (GetGreeter) GetPersonalPermissionInfo(req rolevo.GetPersonalPermissionInfoReqVo) rolevo.GetPersonalPermissionInfoRespVo {
	res, err := service.GetPersonalPermissionInfoForFuse(req.OrgId, req.UserId, req.ProjectId, req.IssueId, req.SourceChannel)
	return rolevo.GetPersonalPermissionInfoRespVo{Err: vo.NewErr(err), Data: res}
}
