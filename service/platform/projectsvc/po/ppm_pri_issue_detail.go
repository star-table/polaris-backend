package po

import "time"

type PpmPriIssueDetail struct {
	Id         int64 `db:"id,omitempty" json:"id"`
	OrgId      int64 `db:"org_id,omitempty" json:"orgId"`
	IssueId    int64 `db:"issue_id,omitempty" json:"issueId"`
	ProjectId  int64 `db:"project_id,omitempty" json:"projectId"`
	StoryPoint int   `db:"story_point" json:"storyPoint"`
	//Tags         string    `db:"tags,omitempty" json:"tags"`
	Remark       *string   `db:"remark,omitempty" json:"remark"`
	RemarkDetail *string   `db:"remark_detail,omitempty" json:"remarkDetail"`
	Status       int64     `db:"status,omitempty" json:"status"`
	Creator      int64     `db:"creator,omitempty" json:"creator"`
	CreateTime   time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime   time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int       `db:"version,omitempty" json:"version"`
	IsDelete     int       `db:"is_delete,omitempty" json:"isDelete"`
}

//func (*PpmPriIssueDetail) TableName() string {
//	return "ppm_pri_issue_detail"
//}
