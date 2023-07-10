package rolevo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type GetOrgRoleListReqVo struct {
	OrgId int64 `json:"orgId"`
}

type GetOrgRoleListRespVo struct {
	vo.Err
	Data []*vo.Role `json:"data"`
}

type CreateOrgReqVo struct {
	OrgId  int64            `json:"orgId"`
	UserId int64            `json:"userId"`
	Input  vo.CreateRoleReq `json:"input"`
}

type DelRoleReqVo struct {
	OrgId  int64         `json:"orgId"`
	UserId int64         `json:"userId"`
	Input  vo.DelRoleReq `json:"input"`
}

type UpdateRoleReqVo struct {
	OrgId  int64            `json:"orgId"`
	UserId int64            `json:"userId"`
	Input  vo.UpdateRoleReq `json:"input"`
}

type GetUserAdminFlagReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetUserAdminFlagRespVo struct {
	Data *bo.UserAdminFlagBo `json:"data"`
	vo.Err
}

type ClearUserRoleReqVo struct {
	OrgId     int64   `json:"orgId"`
	UserIds   []int64 `json:"userIds"`
	ProjectId int64   `json:"projectId"`
}

type GetProjectRoleListReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type UpdateOrgOwnerReqVo struct {
	OrgId int64 `json:"orgId"`
	UserId int64 `json:"userId"`
	OldOwnerId int64 `json:"oldOwnerId"`
	NewOwnerId int64 `json:"newOwnerId"`
} 
