package bo

import "time"

type RoleUserBo struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId  int64     `db:"project_id,omitempty" json:"projectId"`
	RoleId     int64     `db:"role_id,omitempty" json:"roleId"`
	UserId     int64     `db:"user_id,omitempty" json:"userId"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*RoleUserBo) TableName() string {
	return "ppm_rol_role_user"
}

type RoleUser struct {
	RoleId       int64  `db:"role_id,omitempty" json:"roleId"`
	RoleName     string `db:"role_name,omitempty" json:"roleName"`
	UserId       int64  `db:"user_id,omitempty" json:"userId"`
	RoleLangCode string `db:"lang_code,omitempty" json:"roleLangCode"`
}

type UserAdminFlagBo struct {
	IsAdmin         bool `json:"isAdmin"`         //超级管理员
	IsManager       bool `json:"isManager"`       //是管理员
	IsPlatformAdmin bool `json:"isPlatformAdmin"` // 是否是三方平台的管理验，比如：飞书的应用管理员、钉钉应用管理员
}

type RoleUserWithOrgId struct {
	OrgId        int64  `db:"org_id,omitempty" json:"orgId"`
	RoleId       int64  `db:"role_id,omitempty" json:"roleId"`
	RoleName     string `db:"role_name,omitempty" json:"roleName"`
	UserId       int64  `db:"user_id,omitempty" json:"userId"`
	RoleLangCode string `db:"lang_code,omitempty" json:"roleLangCode"`
}
