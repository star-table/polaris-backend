package po

import "time"

type PpmResResourceRelation struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId  int64     `db:"project_id,omitempty" json:"projectId"`
	IssueId    int64     `db:"issue_id,omitempty" json:"issueId"`
	ResourceId int64     `db:"resource_id,omitempty" json:"resourceId"`
	SourceType int       `db:"source_type,omitempty" json:"sourceType"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	ColumnId   string    `db:"column_id,omitempty" json:"columnId"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmResResourceRelation) TableName() string {
	return "ppm_res_resource_relation"
}
