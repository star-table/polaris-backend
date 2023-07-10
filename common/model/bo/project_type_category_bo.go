package bo

import "time"

type ProjectTypeCategoryBo struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	Name       string    `db:"name,omitempty" json:"name"`
	Icon       string    `db:"icon,omitempty" json:"icon"`
	ObjectType int       `db:"object_type,omitempty" json:"objectType"`
	Sort       int       `db:"sort,omitempty" json:"sort"`
	Remark     string    `db:"remark,omitempty" json:"remark"`
	IsReadonly int       `db:"is_readonly,omitempty" json:"isReadonly"`
	Status     int       `db:"status,omitempty" json:"status"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
	ProjectTypeTotal int64 `json:"projectTypeTotal"`
	ProjectTypeList []ProjectTypeBo `json:"projectTypeList"`
}

func (*ProjectTypeCategoryBo) TableName() string {
	return "ppm_prs_project_type_category"
}
