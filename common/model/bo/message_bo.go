package bo

import "time"

type MessageBo struct {
	Id					int64		`db:"id,omitempty" json:"id"`
	OrgId					int64		`db:"org_id,omitempty" json:"orgId"`
	Topic					string		`db:"topic,omitempty" json:"topic"`
	Type					int		`db:"type,omitempty" json:"type"`
	ProjectId					int64		`db:"project_id,omitempty" json:"projectId"`
	IssueId					int64		`db:"issue_id,omitempty" json:"issueId"`
	TrendsId					int64		`db:"trends_id,omitempty" json:"trendsId"`
	Info					string		`db:"info,omitempty" json:"info"`
	Content					*string		`db:"content,omitempty" json:"content"`
	FailCount					int		`db:"fail_count,omitempty" json:"failCount"`
	FailTime					time.Time		`db:"fail_time,omitempty" json:"failTime"`
	FailMsg					string		`db:"fail_msg,omitempty" json:"failMsg"`
	FinishStatus					string		`db:"finish_status,omitempty" json:"finishStatus"`
	FinishMsg					string		`db:"finish_msg,omitempty" json:"finishMsg"`
	StartTime					time.Time		`db:"start_time,omitempty" json:"startTime"`
	Status					int		`db:"status,omitempty" json:"status"`
	Creator					int64		`db:"creator,omitempty" json:"creator"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"createTime"`
	Updator					int64		`db:"updator,omitempty" json:"updator"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"updateTime"`
	Version					int		`db:"version,omitempty" json:"version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"isDelete"`
}

func (*MessageBo) TableName() string {
	return "ppm_tak_message"
}