package po

import "time"

type PpmPriIssue struct {
	Id                  int64     `db:"id,omitempty" json:"id"`
	OrgId               int64     `db:"org_id,omitempty" json:"orgId"`
	Code                string    `db:"code,omitempty" json:"code"`
	ProjectId           int64     `db:"project_id,omitempty" json:"projectId"`
	ProjectObjectTypeId int64     `db:"project_object_type_id,omitempty" json:"projectObjectTypeId"`
	Title               string    `db:"title,omitempty" json:"title"`
	Path                string    `db:"path,omitempty" json:"path"`
	IsFiling            int       `db:"is_filing,omitempty" json:"isFiling"`
	Owner               int64     `db:"owner,omitempty" json:"owner"`
	OwnerChangeTime     time.Time `db:"owner_change_time,omitempty" json:"ownerChangeTime"`
	PriorityId          int64     `db:"priority_id,omitempty" json:"priorityId"`
	SourceId            int64     `db:"source_id,omitempty" json:"sourceId"`
	PropertyId          int64     `db:"property_id,omitempty" json:"propertyId"`
	IssueObjectTypeId   int64     `db:"issue_object_type_id,omitempty" json:"issueObjectTypeId"`
	PlanStartTime       time.Time `db:"plan_start_time,omitempty" json:"planStartTime"`
	PlanEndTime         time.Time `db:"plan_end_time,omitempty" json:"planEndTime"`
	StartTime           time.Time `db:"start_time,omitempty" json:"startTime"`
	EndTime             time.Time `db:"end_time,omitempty" json:"endTime"`
	PlanWorkHour        int       `db:"plan_work_hour,omitempty" json:"planWorkHour"`
	IterationId         int64     `db:"iteration_id,omitempty" json:"iterationId"`
	VersionId           int64     `db:"version_id,omitempty" json:"versionId"`
	ModuleId            int64     `db:"module_id,omitempty" json:"moduleId"`
	ParentId            int64     `db:"parent_id,omitempty" json:"parentId"`
	Status              int64     `db:"status,omitempty" json:"status"`
	AuditStatus         int64     `db:"audit_status,omitempty" json:"auditStatus"`
	Sort                int64     `db:"sort,omitempty" json:"sort"`
	Creator             int64     `db:"creator,omitempty" json:"creator"`
	CreateTime          time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator             int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime          time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version             int       `db:"version,omitempty" json:"version"`
	IsDelete            int       `db:"is_delete,omitempty" json:"isDelete"`
	TableId             int64     `db:"table_id,omitempty" json:"tableId"`
}

//func (*PpmPriIssue) TableName() string {
//	return "ppm_pri_issue"
//}
