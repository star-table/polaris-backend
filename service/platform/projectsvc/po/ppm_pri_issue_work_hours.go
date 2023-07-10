package po

import "time"

type PpmPriIssueWorkHours struct {
	// Id 自增id
	Id uint64 `db:"id,omitempty" json:"id"`
	// OrgId 组织id
	OrgId int64 `db:"org_id,omitempty" json:"orgId"`
	// 关联的项目id
	ProjectId int64 `db:"project_id,omitempty" json:"projectId"`
	// IssueId 关联的任务id
	IssueId int64 `db:"issue_id,omitempty" json:"issueId"`
	// Type 记录类型：1预估工时记录，2实际工时记录，3详细预估工时
	Type int8 `db:"type,omitempty" json:"type"`
	// 工作者id
	WorkerId int64 `db:"worker_id,omitempty" json:"workerId"`
	// NeedTime 工时的时间，单位分钟
	NeedTime uint32 `db:"need_time,omitempty" json:"needTime"`
	// 剩余工时计算方式：1动态计算；2手动填写
	RemainTimeCalType uint32 `db:"remain_time_cal_type,omitempty" json:"remainTimeCalType"`
	// "手动填写"剩余工时时的剩余工时的值
	RemainTime uint32 `db:"remain_time,omitempty" json:"remainTime"`
	// StartTime 开始时间，时间戳
	StartTime uint32 `db:"start_time,omitempty" json:"startTime"`
	// EndTime 工时记录的结束时间，时间戳
	EndTime uint32 `db:"end_time,omitempty" json:"endTime"`
	// Desc 工时记录的内容，工作内容
	Desc string `db:"desc,omitempty" json:"desc"`
	// Creator 工时记录创建者id
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	// Updator 工时记录更新者的id
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int32     `db:"version,omitempty" json:"version"`
	IsDelete   int8      `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmPriIssueWorkHours) TableName() string {
	return "ppm_pri_issue_work_hours"
}

type SimpleWorkHourPo struct {
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

type OnePersonWorkHourData struct {
	TotalPredict []*PpmPriIssueWorkHours `json:"totalPredict"`
	SubPredict   []*PpmPriIssueWorkHours `json:"subPredict"`
	Actual       []*PpmPriIssueWorkHours `json:"actual"`
}
