package po

import "time"

type PpmPriIssueTag struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId  int64     `db:"project_id,omitempty" json:"projectId"`
	IssueId    int64     `db:"issue_id,omitempty" json:"issueId"`
	TagId      int64     `db:"tag_id,omitempty" json:"tagId"`
	TagName    string    `db:"tag_name,omitempty" json:"tagName"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmPriIssueTag) TableName() string {
	return "ppm_pri_issue_tag"
}
