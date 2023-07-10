package bo

import "time"

type RoleDepartmentBo struct {
	Id           int64     `db:"id,omitempty" json:"id"`
	OrgId        int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId    int64     `db:"project_id,omitempty" json:"projectId"`
	RoleId       int64     `db:"role_id,omitempty" json:"roleId"`
	DepartmentId int64     `db:"department_id,omitempty" json:"departmentId"`
	Creator      int64     `db:"creator,omitempty" json:"creator"`
	CreateTime   time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime   time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int       `db:"version,omitempty" json:"version"`
	IsDelete     int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*RoleDepartmentBo) TableName() string {
	return "ppm_rol_role_department"
}

type DepartmentAdminFlagBo struct {
	IsAdmin bool `json:"isAdmin"`  //超级管理员
	IsManager bool `json:"isManager"`	//是管理员
}

type RoleDepartment struct {
	RoleId       int64  `db:"role_id,omitempty" json:"roleId"`
	RoleName     string `db:"role_name,omitempty" json:"roleName"`
	DepartmentId       int64  `db:"department_id,omitempty" json:"departmentId"`
	RoleLangCode string `db:"lang_code,omitempty" json:"roleLangCode"`
}
