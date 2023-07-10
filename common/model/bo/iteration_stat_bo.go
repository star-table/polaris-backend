package bo

import "time"

type IterationStatBo struct {
	Id                     int64     `db:"id,omitempty" json:"id"`
	OrgId                  int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId              int64     `db:"project_id,omitempty" json:"projectId"`
	IterationId            int64     `db:"iteration_id,omitempty" json:"iterationId"`
	IssueCount             int       `db:"issue_count,omitempty" json:"issueCount"`
	IssueWaitCount         int       `db:"issue_wait_count,omitempty" json:"issueWaitCount"`
	IssueRunningCount      int       `db:"issue_running_count,omitempty" json:"issueRunningCount"`
	IssueOverdueCount      int       `db:"issue_overdue_count,omitempty" json:"issueOverdueCount"`
	IssueEndCount          int       `db:"issue_end_count,omitempty" json:"issueEndCount"`
	DemandCount            int       `db:"demand_count,omitempty" json:"demandCount"`
	DemandWaitCount        int       `db:"demand_wait_count,omitempty" json:"demandWaitCount"`
	DemandRunningCount     int       `db:"demand_running_count,omitempty" json:"demandRunningCount"`
	DemandOverdueCount     int       `db:"demand_overdue_count,omitempty" json:"demandOverdueCount"`
	DemandEndCount         int       `db:"demand_end_count,omitempty" json:"demandEndCount"`
	StoryPointCount        int       `db:"story_point_count,omitempty" json:"storyPointCount"`
	StoryPointWaitCount    int       `db:"story_point_wait_count,omitempty" json:"storyPointWaitCount"`
	StoryPointRunningCount int       `db:"story_point_running_count,omitempty" json:"storyPointRunningCount"`
	StoryPointOverdueCount int       `db:"story_point_overdue_count,omitempty" json:"storyPointOverdueCount"`
	StoryPointEndCount     int       `db:"story_point_end_count,omitempty" json:"storyPointEndCount"`
	TaskCount              int       `db:"task_count,omitempty" json:"taskCount"`
	TaskWaitCount          int       `db:"task_wait_count,omitempty" json:"taskWaitCount"`
	TaskRunningCount       int       `db:"task_running_count,omitempty" json:"taskRunningCount"`
	TaskOverdueCount       int       `db:"task_overdue_count,omitempty" json:"taskOverdueCount"`
	TaskEndCount           int       `db:"task_end_count,omitempty" json:"taskEndCount"`
	BugCount               int       `db:"bug_count,omitempty" json:"bugCount"`
	BugWaitCount           int       `db:"bug_wait_count,omitempty" json:"bugWaitCount"`
	BugRunningCount        int       `db:"bug_running_count,omitempty" json:"bugRunningCount"`
	BugOverdueCount        int       `db:"bug_overdue_count,omitempty" json:"bugOverdueCount"`
	BugEndCount            int       `db:"bug_end_count,omitempty" json:"bugEndCount"`
	TesttaskCount          int       `db:"testtask_count,omitempty" json:"testtaskCount"`
	TesttaskWaitCount      int       `db:"testtask_wait_count,omitempty" json:"testtaskWaitCount"`
	TesttaskRunningCount   int       `db:"testtask_running_count,omitempty" json:"testtaskRunningCount"`
	TesttaskOverdueCount   int       `db:"testtask_overdue_count,omitempty" json:"testtaskOverdueCount"`
	TesttaskEndCount       int       `db:"testtask_end_count,omitempty" json:"testtaskEndCount"`
	Ext                    string    `db:"ext,omitempty" json:"ext"`
	StatDate               time.Time `db:"stat_date,omitempty" json:"statDate"`
	Status                 int64     `db:"status,omitempty" json:"status"`
	Creator                int64     `db:"creator,omitempty" json:"creator"`
	CreateTime             time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator                int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime             time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version                int       `db:"version,omitempty" json:"version"`
	IsDelete               int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*IterationStatBo) TableName() string {
	return "ppm_sta_iteration_stat"
}
