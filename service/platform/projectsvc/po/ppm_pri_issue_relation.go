package po

import "time"

type PpmPriIssueRelation struct {
	Id           int64     `db:"id,omitempty" json:"id"`
	OrgId        int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId    int64     `db:"project_id,omitempty" json:"projectId"`
	IssueId      int64     `db:"issue_id,omitempty" json:"issueId"`
	RelationId   int64     `db:"relation_id,omitempty" json:"relationId"`
	RelationCode string    `db:"relation_code,omitempty" json:"relationCode"`
	RelationType int       `db:"relation_type,omitempty" json:"relationType"`
	Status       int       `db:"status,omitempty" json:"status"`
	Creator      int64     `db:"creator,omitempty" json:"creator"`
	CreateTime   time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime   time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int       `db:"version,omitempty" json:"version"`
	IsDelete     int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmPriIssueRelation) TableName() string {
	return "ppm_pri_issue_relation"
}

type PpmPriIssueRelationCount struct {
	IssueId int64 `db:"issue_id,omitempty" json:"issueId"`
	Total   int64 `db:"total,omitempty" json:"total"`
}
