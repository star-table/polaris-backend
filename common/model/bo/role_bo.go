package bo

import "time"

type RoleBo struct {
	Id					int64		`db:"id,omitempty" json:"id"`
	OrgId					int64		`db:"org_id,omitempty" json:"orgId"`
	ProjectId 				int64 		`db:"project_id,omitempty" json:"projectId"`
	LangCode					string		`db:"lang_code,omitempty" json:"langCode"`
	Name					string		`db:"name,omitempty" json:"name"`
	Remark					string		`db:"remark,omitempty" json:"remark"`
	IsReadonly					int		`db:"is_readonly,omitempty" json:"isReadonly"`
	IsModifyPermission					int		`db:"is_modify_permission,omitempty" json:"isModifyPermission"`
	IsDefault					int		`db:"is_default,omitempty" json:"isDefault"`
	RoleGroupId					int64		`db:"role_group_id,omitempty" json:"roleGroupId"`
	Status					int		`db:"status,omitempty" json:"status"`
	Creator					int64		`db:"creator,omitempty" json:"creator"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"createTime"`
	Updator					int64		`db:"updator,omitempty" json:"updator"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"updateTime"`
	Version					int		`db:"version,omitempty" json:"version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"isDelete"`
}

func (*RoleBo) TableName() string {
	return "ppm_rol_role"
}


type RoleInitResp struct {
	OrgSuperAdminRoleId  int64
	OrgNormalAdminRoleId int64
}

type UpdateRoleBo struct {
	Name *string `json:"name"`
}