package bo

import (
	"time"
)

type PpmPriIterationStatusRelationBo struct {
	Id            int64     `db:"id,omitempty" json:"id"`
	OrgId         int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId     int64     `db:"project_id,omitempty" json:"projectId"`
	IterationId   int64     `db:"iteration_id,omitempty" json:"iterationId"`
	StatusId      int64     `db:"status_id,omitempty" json:"statusId"`
	PlanStartTime time.Time `db:"plan_start_time,omitempty" json:"planStartTime"`
	PlanEndTime   time.Time `db:"plan_end_time,omitempty" json:"planEndTime"`
	StartTime     time.Time `db:"start_time,omitempty" json:"startTime"`
	EndTime       time.Time `db:"end_time,omitempty" json:"endTime"`
	Creator       int64     `db:"creator,omitempty" json:"creator"`
	CreateTime    time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator       int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime    time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version       int       `db:"version,omitempty" json:"version"`
	IsDelete      int       `db:"is_delete,omitempty" json:"isDelete"`
}
