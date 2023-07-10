package roleapi

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	service "github.com/star-table/polaris-backend/service/platform/orgsvc/service/roleservice"
)

func (PostGreeter) Authenticate(req rolevo.AuthenticateReqVo) vo.VoidErr {
	err := service.Authenticate(req.OrgId, req.UserId, req.AuthInfoReqVo.ProjectAuthInfo, req.AuthInfoReqVo.IssueAuthInfo, req.Path, req.Operation, req.AuthInfoReqVo.Fields)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) RoleUserRelation(req rolevo.RoleUserRelationReqVo) vo.VoidErr {
	err := service.RoleUserRelation(req.OrgId, req.UserId, req.RoleId)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) RemoveRoleUserRelation(req rolevo.RemoveRoleUserRelationReqVo) vo.VoidErr {
	err := service.RemoveRoleUserRelation(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) RoleInit(req rolevo.RoleInitReqVo) rolevo.RoleInitRespVo {
	respVo := rolevo.RoleInitRespVo{}
	roleInitResp, err := service.RoleInit(req.OrgId)
	respVo.RoleInitResp = roleInitResp
	respVo.Err = vo.NewErr(err)
	return respVo
}

func (GetGreeter) ChangeDefaultRole() vo.CommonRespVo {
	err := service.ChangeDefaultRole()
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}

func (PostGreeter) RemoveRoleDepartmentRelation(req rolevo.RemoveRoleDepartmentRelationReqVo) vo.VoidErr{
	err := service.RemoveRoleDepartmentRelation(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}
