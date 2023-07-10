package bo

import "time"

type DepartmentBo struct {
	Id					int64		`db:"id,omitempty" json:"id"`
	OrgId					int64		`db:"org_id,omitempty" json:"orgId"`
	Name					string		`db:"name,omitempty" json:"name"`
	Code					string		`db:"code,omitempty" json:"code"`
	ParentId					int64		`db:"parent_id,omitempty" json:"parentId"`
	Path          			string    `db:"path,omitempty" json:"path"`
	Sort					int		`db:"sort,omitempty" json:"sort"`
	IsHide					int		`db:"is_hide,omitempty" json:"isHide"`
	SourceChannel					string		`db:"source_channel,omitempty" json:"sourceChannel"`
	Status					int		`db:"status,omitempty" json:"status"`
	Creator					int64		`db:"creator,omitempty" json:"creator"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"createTime"`
	Updator					int64		`db:"updator,omitempty" json:"updator"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"updateTime"`
	Version					int		`db:"version,omitempty" json:"version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"isDelete"`
}

func (*DepartmentBo) TableName() string {
	return "ppm_org_department"
}

type UserDepartmentInfo struct {
	UserId         int64  `db:"user_id,omitempty" json:"userId"`
	DepartmentId   int64  `db:"department_id,omitempty" json:"departmentId"`
	IsLeader       int    `db:"is_leader,omitempty" json:"isLeader"`
	DepartmentName string `db:"name,omitempty" json:"name"`
}