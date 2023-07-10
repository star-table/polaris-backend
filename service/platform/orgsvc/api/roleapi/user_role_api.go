package roleapi

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	service "github.com/star-table/polaris-backend/service/platform/orgsvc/service/roleservice"
)

func (GetGreeter) GetOrgRoleUser(req rolevo.GetOrgRoleUserReqVo) rolevo.GetOrgRoleUserRespVo {
	res, err := service.GetOrgRoleUser(req.OrgId, req.ProjectId)
	return rolevo.GetOrgRoleUserRespVo{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) GetOrgAdminUser(req rolevo.GetOrgAdminUserReqVo) rolevo.GetOrgAdminUserRespVo {
	res, err := service.GetOrgAdminUser(req.OrgId)
	return rolevo.GetOrgAdminUserRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) GetOrgAdminUserBatch(req rolevo.GetOrgAdminUserBatchReqVo) rolevo.GetOrgAdminUserBatchRespVo {
	res, err := service.GetOrgAdminUserBatch(req.Input.OrgIds)
	return rolevo.GetOrgAdminUserBatchRespVo{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) UpdateUserOrgRole(reqVo rolevo.UpdateUserOrgRoleReqVo) vo.CommonRespVo {
	res, err := service.UpdateUserOrgRole(reqVo)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (GetGreeter) GetOrgRoleList(reqVo rolevo.GetOrgRoleListReqVo) rolevo.GetOrgRoleListRespVo {
	res, err := service.GetOrgRoleList(reqVo.OrgId)
	return rolevo.GetOrgRoleListRespVo{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) GetUserAdminFlag(reqVo rolevo.GetUserAdminFlagReqVo) rolevo.GetUserAdminFlagRespVo {
	res, err := service.GetUserAdminFlag(reqVo.OrgId, reqVo.UserId)
	return rolevo.GetUserAdminFlagRespVo{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) GetProjectRoleList(reqVo rolevo.GetProjectRoleListReqVo) rolevo.GetOrgRoleListRespVo {
	res, err := service.GetProjectRoleList(reqVo.OrgId, reqVo.ProjectId)
	return rolevo.GetOrgRoleListRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) UpdateOrgOwner(req rolevo.UpdateOrgOwnerReqVo) vo.CommonRespVo {
	res, err := service.UpdateOrgAdmin(req.UserId, req.OrgId, req.OldOwnerId, req.NewOwnerId)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateUserOrgRoleBatch(reqVo rolevo.UpdateUserOrgRoleBatchReqVo) vo.CommonRespVo {
	res, err := service.UpdateUserOrgRoleBatch(reqVo)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}
