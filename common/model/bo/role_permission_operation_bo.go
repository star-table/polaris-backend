package bo

import "time"

type RolePermissionOperationBo struct {
	Id					int64		`db:"id,omitempty" json:"id"`
	OrgId					int64		`db:"org_id,omitempty" json:"orgId"`
	RoleId					int64		`db:"role_id,omitempty" json:"roleId"`
	ProjectId					int64		`db:"project_id,omitempty" json:"projectId"`
	PermissionId					int64		`db:"permission_id,omitempty" json:"permissionId"`
	PermissionPath					string		`db:"permission_path,omitempty" json:"permissionPath"`
	OperationCodes					string		`db:"operation_codes,omitempty" json:"operationCodes"`
	Creator					int64		`db:"creator,omitempty" json:"creator"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"createTime"`
	Updator					int64		`db:"updator,omitempty" json:"updator"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"updateTime"`
	Version					int		`db:"version,omitempty" json:"version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"isDelete"`
}

func (*RolePermissionOperationBo) TableName() string {
	return "ppm_rol_role_permission_operation"
}