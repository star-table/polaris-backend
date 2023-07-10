package bo

import "time"

type IssueWorkHoursBo struct {
	// Id 自增id
	Id uint64 `db:"id,omitempty" json:"id"`
	// OrgId 组织id
	OrgId     int64 `db:"org_id,omitempty" json:"orgId"`
	ProjectId int64 `db:"project_id,omitempty" json:"projectId"`
	// IssueId 关联的任务id
	IssueId int64 `db:"issue_id,omitempty" json:"issueId"`
	// Type 记录类型：1总预估工时，2实际工时记录，3子预估工时
	Type int8 `db:"type,omitempty" json:"type"`
	// 工作者id
	WorkerId int64 `db:"worker_id,omitempty" json:"workerId"`
	// NeedTime 工时的时间，单位分钟
	NeedTime uint32 `db:"need_time,omitempty" json:"needTime"`
	// StartTime 开始时间，时间戳
	StartTime uint32 `db:"start_time,omitempty" json:"startTime"`
	// EndTime 工时记录的结束时间，时间戳
	EndTime uint32 `db:"end_time,omitempty" json:"endTime"`
	// Desc 工时记录的内容，工作内容
	Desc string `db:"desc,omitempty" json:"desc"`
	// Creator 工时记录创建者id
	Creator int64 `db:"creator,omitempty" json:"creator"`
	// Updator 工时记录更新者的id
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int32     `db:"version,omitempty" json:"version"`
	IsDelete   int8      `db:"is_delete,omitempty" json:"isDelete"`
}

type IssueWorkHoursBoListCondBo struct {
	Ids        []uint64 `json:"ids"`
	OrgId      int64    `json:"orgId"`
	ProjectId  int64    `json:"projectId"`
	ProjectIds []int64  `json:"projectIds"`
	IssueId    int64    `json:"issueId"`
	// 查询多个任务下的工时
	IssueIds  []int64 `json:"issueIds"`
	Type      uint8   `json:"type"` // 记录类型：1预估工时记录，2实际工时记录，3子预估工时。目前主要可分为2、3；类型 1 的已经弱化掉了。
	Types     []uint8 `json:"types"`
	WorkerIds []int64 `json:"workerIds"`
	// 查询 start_time 的区间数据：start_time>StartTimeT1 && start_time<=StartTimeT2
	StartTimeT1 int64 `json:"startTimeT1"`
	StartTimeT2 int64 `json:"startTimeT2"`
	// 查询 end_time 的区间数据
	EndTimeT1 int64 `json:"endTimeT1"`
	EndTimeT2 int64 `json:"endTimeT2"`
	IsDelete  uint8 `json:"isDelete"`
	Page      int   `json:"page"`
	Size      int   `json:"size"`
}

type UpdateOneIssueWorkHoursBo struct {
	Id                uint64 `json:"id"`
	OrgId             int64  `json:"orgId"`
	ProjectId         int64  `json:"projectId"`
	IssueId           int64  `json:"issueId"`
	WorkerId          int64  `json:"workerId"`
	Type              int8   `json:"type"`
	Desc              string `json:"desc"`
	StartTime         int64  `json:"startTime"`
	EndTime           int64  `json:"endTime"`
	NeedTime          uint32 `json:"needTime"`
	RemainTimeCalType int8   `json:"remainTimeCalType"`
	RemainTime        uint32 `json:"remainTime"`
	IsDelete          int8   `json:"isDelete"`
}

type DeleteIssueWorkHoursBo struct {
	Id       uint64 `json:"id"`
	IsDelete int8   `json:"isDelete"`
}

type DisOrEnableIssueWorkHoursBo struct {
	Id        uint64 `json:"id"`
	ProjectId int64  `json:"projectId"`
}

type MoveIssueWorkHoursByProjectIdBo struct {
	Operator     int64 `json:"operator"`
	SrcProjectId int64 `json:"srcProjectId"`
	DstProjectId int64 `json:"dstProjectId"`
}

type CheckWorkHoursUpdatePowerBo struct {
	CurrentUserId    int64  `json:"currentUserId"`
	IssueId          int64  `json:"issueId"`
	IssueWorkHoursId uint64 `json:"issueWorkHoursId"`
	WorkerId         int64  `json:"workerId"`
}

type SimpleWorkHourBo struct {
	// Id 自增id
	Id uint64 `db:"id,omitempty" json:"id"`
	// OrgId 组织id
	OrgId     int64 `db:"org_id,omitempty" json:"orgId"`
	ProjectId int64 `db:"project_id,omitempty" json:"projectId"`
	// IssueId 关联的任务id
	IssueId int64 `db:"issue_id,omitempty" json:"issueId"`
	// Type 记录类型：1总预估工时，2实际工时记录，3子预估工时
	Type int8 `db:"type,omitempty" json:"type"`
	// 工作者id
	WorkerId int64 `db:"worker_id,omitempty" json:"workerId"`
	// NeedTime 工时的时间，单位分钟
	NeedTime uint32 `db:"need_time,omitempty" json:"needTime"`
	// StartTime 开始时间，时间戳
	StartTime uint32 `db:"start_time,omitempty" json:"startTime"`
	// EndTime 工时记录的结束时间，时间戳
	EndTime uint32 `db:"end_time,omitempty" json:"endTime"`
	// Desc 工时记录的内容，工作内容
	Desc string `db:"desc,omitempty" json:"desc"`
}

// 存储某个人的工时记录。
type WorkHourForSomeoneBo struct {
	TotalPredict []*SimpleWorkHourBo `json:"totalPredict"`
	SubPredict   []*SimpleWorkHourBo `json:"subPredict"`
	Actual       []*SimpleWorkHourBo `json:"actual"`
}
