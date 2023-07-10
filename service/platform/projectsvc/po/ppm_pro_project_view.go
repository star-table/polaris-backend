package po

import "time"

type PpmProProjectView struct {
	Id                  int64     `db:"id,omitempty" json:"id"`
	OrgId               int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId           int64     `db:"project_id,omitempty" json:"projectId"`
	ProjectObjectTypeId int64     `db:"project_object_type_id,omitempty" json:"projectObjectTypeId"`
	ViewType            int       `db:"view_type,omitempty" json:"viewType"`
	ClosedDefaultField  string    `db:"closed_default_field" json:"closedDefaultField"`
	ClosedCustomField   string    `db:"closed_custom_field" json:"closedCustomField"`
	Creator             int64     `db:"creator,omitempty" json:"creator"`
	CreateTime          time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator             int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime          time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version             int       `db:"version,omitempty" json:"version"`
	IsDelete            int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmProProjectView) TableName() string {
	return "ppm_pro_project_view"
}
