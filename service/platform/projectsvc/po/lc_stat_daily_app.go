package po

import "time"

type LcStatDailyApp struct {
	Id                 int64     `db:"id,omitempty" json:"id"`
	OrgId              int64     `db:"org_id,omitempty" json:"orgId"`
	AppId              int64     `db:"app_id,omitempty" json:"appId"`
	ProjectId          int64     `db:"project_id,omitempty" json:"projectId"`
	TableId            int64     `db:"table_id,omitempty" json:"tableId"`
	IterationId        int64     `db:"iteration_id,omitempty" json:"iterationId"`
	Date               time.Time `db:"date,omitempty" json:"date"`
	IssueCount         int32     `db:"issue_count,omitempty" json:"issueCount"`
	IssueNotStartCount int32     `db:"issue_not_start_count,omitempty" json:"issueNotStartCount"`
	IssueRunningCount  int32     `db:"issue_running_count,omitempty" json:"issueRunningCount"`
	IssueOverdueCount  int32     `db:"issue_overdue_count,omitempty" json:"issueOverdueCount"`
	IssueCompleteCount int32     `db:"issue_complete_count,omitempty" json:"issueCompleteCount"`
	Creator            int64     `db:"creator,omitempty" json:"creator"`
	CreateTime         time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator            int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime         time.Time `db:"update_time,omitempty" json:"updateTime"`
	IsDelete           int32     `db:"is_delete,omitempty" json:"isDelete"`
}

func (*LcStatDailyApp) TableName() string {
	return "lc_stat_daily_app"
}
