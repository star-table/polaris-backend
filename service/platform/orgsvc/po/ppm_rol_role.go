package po

import "time"

type PpmRolRole struct {
	Id                 int64     `db:"id,omitempty" json:"id"`
	OrgId              int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId          int64     `db:"project_id,omitempty" json:"projectId"`
	LangCode           string    `db:"lang_code,omitempty" json:"langCode"`
	Name               string    `db:"name,omitempty" json:"name"`
	Remark             string    `db:"remark,omitempty" json:"remark"`
	IsReadonly         int       `db:"is_readonly,omitempty" json:"isReadonly"`
	IsModifyPermission int       `db:"is_modify_permission,omitempty" json:"isModifyPermission"`
	IsDefault          int       `db:"is_default,omitempty" json:"isDefault"`
	RoleGroupId        int64     `db:"role_group_id,omitempty" json:"roleGroupId"`
	Status             int       `db:"status,omitempty" json:"status"`
	Creator            int64     `db:"creator,omitempty" json:"creator"`
	CreateTime         time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator            int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime         time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version            int       `db:"version,omitempty" json:"version"`
	IsDelete           int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmRolRole) TableName() string {
	return "ppm_rol_role"
}

type RoleUser struct {
	RoleId       int64  `db:"role_id,omitempty" json:"roleId"`
	RoleName     string `db:"role_name,omitempty" json:"roleName"`
	UserId       int64  `db:"user_id,omitempty" json:"userId"`
	RoleLangCode string `db:"lang_code,omitempty" json:"roleLangCode"`
}

type RoleUserWithOrgId struct {
	OrgId        int64  `db:"org_id,omitempty" json:"orgId"`
	RoleId       int64  `db:"role_id,omitempty" json:"roleId"`
	RoleName     string `db:"role_name,omitempty" json:"roleName"`
	UserId       int64  `db:"user_id,omitempty" json:"userId"`
	RoleLangCode string `db:"lang_code,omitempty" json:"roleLangCode"`
}
