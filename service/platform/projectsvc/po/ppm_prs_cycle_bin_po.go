package po

import "time"

type PpmPrsRecycleBin struct {
	Id           int64     `db:"id,omitempty" json:"id"`
	OrgId        int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId    int64     `db:"project_id,omitempty" json:"projectId"`
	RelationId   int64     `db:"relation_id,omitempty" json:"relationId"`
	RelationType int       `db:"relation_type,omitempty" json:"relationType"`
	Creator      int64     `db:"creator,omitempty" json:"creator"`
	CreateTime   time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime   time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int       `db:"version,omitempty" json:"version"`
	IsDelete     int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmPrsRecycleBin) TableName() string {
	return "ppm_prs_recycle_bin"
}
