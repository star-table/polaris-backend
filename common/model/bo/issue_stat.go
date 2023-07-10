package bo

import "time"

type IssueStatusStatCondBo struct {
	OrgId        int64
	UserId       int64
	ProjectId    *int64
	IterationId  *int64
	RelationType *int
}

type IssueStatusStatBo struct {
	ProjectTypeId             int64
	ProjectTypeName           string
	ProjectTypeLangCode       string
	IssueCount                int
	IssueWaitCount            int
	IssueRunningCount         int
	IssueEndCount             int
	IssueEndTodayCount        int //今日完成
	IssueOverdueCount         int //逾期数量
	IssueOverdueTodayCount    int //今日逾期
	IssueOverdueEndCount      int //逾期完成
	IssueOverdueTomorrowCount int //明日逾期/即将逾期
	TodayCount                int //今日指派给我
	TodayCreateCount          int //今日创建的任务
	StoryPointCount           int
	StoryPointWaitCount       int
	StoryPointRunningCount    int
	StoryPointEndCount        int
	WaitConfirmedCount        int
}

type StatExtBo struct {
	Issue StatIssueExtBo `json:"issue"`
}

type StatIssueExtBo struct {
	Data map[interface{}]interface{} `json:"data"`
}

type IssueStatistic struct {
	ProjectId      int64  `db:"project_id" json:"projectId"`
	IterationId    int64  `db:"iteration_id" json:"iterationId"`
	IterationName  string `json:"iterationName" json:"iterationName"`
	Overdue        int64  `db:"overdue" json:"overdue"`
	All            int64  `db:"all" json:"all"`
	Finish         int64  `db:"finish" json:"finish"`
	RelateUnfinish int64  `json:"relateUnfinish" json:"relateUnfinish"`
}

type RelateIssueTotal struct {
	ProjectId int64 `db:"project_id" json:"projectId"`
	Total     int64 `db:"total" json:"total"`
}

type IssueStatByObjectType struct {
	ProjectObjectTypeId int64  `db:"project_object_type_id,omitempty" json:"projectObjectTypeId"`
	Count               int64  `db:"count" json:"count"`
	LangCode            string `db:"lang_code" json:"langCode"`
}

type IssueStatByStatus struct {
	Status int64 `db:"status" json:status`
	Count  int64 `db:"count" json:"count"`
}

type IssueStatByProjectIdAndObjectId struct {
	Name                  string `db:"name" json:"name"`
	Count                 int64  `db:"count" json:"count"`
	StatusType            int64  `db:"status_type" json:statusType`
	ProjectObjectTypeId   int64  `db:"project_object_type_id,omitempty" json:"projectObjectTypeId"`
	ProjectObjectTypeName string
}

type IssueAssignCountBo struct {
	Owner int64 `db:"owner" json:"owner"`
	Count int64 `db:"count" json:"count"`
}

type IssueDailyPersonalWorkCompletionStatBo struct {
	Count int64     `db:"count" json:"count"`
	Date  time.Time `db:"date" json:"date"`
}

type CommonCountBo struct {
	Count int64 `json:"count"`
}
