package rolevo

import "github.com/star-table/polaris-backend/common/model/vo"

type RoleUser struct {
	RoleId       int64  `json:"roleId"`
	RoleName     string `json:"roleName"`
	UserId       int64  `json:"userId"`
	RoleLangCode string `json:"roleLangCode"`
}

type GetOrgRoleUserReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type GetOrgAdminUserReqVo struct {
	OrgId int64 `json:"orgId"`
}

type GetOrgAdminUserBatchReqVo struct {
	OrgId int64                        `json:"orgId"`
	Input GetOrgAdminUserBatchReqInput `json:"input"`
}

type GetOrgAdminUserBatchReqInput struct {
	OrgIds []int64 `json:"orgIds"`
}

type GetOrgRoleUserRespVo struct {
	vo.Err
	Data []RoleUser `json:"data"`
}

type GetOrgAdminUserRespVo struct {
	vo.Err
	//用户id
	Data []int64 `json:"data"`
}

type GetOrgAdminUserBatchRespVo struct {
	vo.Err
	// orgId => [用户id]
	Data map[int64][]int64 `json:"data"`
}

type UpdateUserOrgRoleReqVo struct {
	OrgId         int64 `json:"orgId"`
	CurrentUserId int64 `json:"currentUserId"`
	UserId        int64 `json:"userId"`
	RoleId        int64 `json:"roleId"`
	ProjectId     *int64 `json:"projectId"`
}


type UpdateUserOrgRoleBatchReqVo struct {
	OrgId         int64 `json:"orgId"`
	CurrentUserId int64 `json:"currentUserId"`
	Input UpdateUserOrgRoleBatchData `json:"input"`
}

type UpdateUserOrgRoleBatchData struct {
	UserIds []int64 `json:"userIds"`
	RoleId int64 `json:"roleId"`
}
