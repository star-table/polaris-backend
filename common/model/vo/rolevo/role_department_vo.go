package rolevo

import "github.com/star-table/polaris-backend/common/model/vo"

type UpdateDepartmentOrgRoleReqVo struct {
	OrgId         int64 `json:"orgId"`
	CurrentUserId int64 `json:"currentUserId"`
	DepartmentId        int64 `json:"departmentId"`
	RoleId        int64 `json:"roleId"`
	ProjectId     *int64 `json:"projectId"`
}

type RoleDepartment struct {
	RoleId       int64  `json:"roleId"`
	RoleName     string `json:"roleName"`
	DepartmentId       int64  `json:"departmentId"`
	RoleLangCode string `json:"roleLangCode"`
}

type GetOrgRoleDepartmentReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type GetOrgRoleDepartmentRespVo struct {
	vo.Err
	Data []RoleDepartment `json:"data"`
}
