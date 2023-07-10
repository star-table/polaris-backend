package roleapi

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	service "github.com/star-table/polaris-backend/service/platform/orgsvc/service/roleservice"
)

func (PostGreeter) CreateRole(req rolevo.CreateOrgReqVo) vo.CommonRespVo {
	res, err := service.CreateRole(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) DelRole(req rolevo.DelRoleReqVo) vo.CommonRespVo {
	err := service.DelRole(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{Void: &vo.Void{ID: 0}, Err: vo.NewErr(err)}
}

func (PostGreeter) UpdateRole(req rolevo.UpdateRoleReqVo) vo.CommonRespVo {
	id, err := service.UpdateRole(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{Void: &vo.Void{ID: id}, Err: vo.NewErr(err)}
}

func (PostGreeter) ClearUserRoleList(req rolevo.ClearUserRoleReqVo) vo.CommonRespVo {
	err := service.ClearUserRoleList(req.OrgId, req.UserIds, req.ProjectId)
	return vo.CommonRespVo{Void: &vo.Void{ID: 0}, Err: vo.NewErr(err)}
}
