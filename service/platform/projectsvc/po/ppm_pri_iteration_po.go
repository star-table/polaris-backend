package po

import "time"

type PpmPriIteration struct {
	Id            int64     `db:"id,omitempty" json:"id"`
	OrgId         int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId     int64     `db:"project_id,omitempty" json:"projectId"`
	Name          string    `db:"name,omitempty" json:"name"`
	Sort		  int64     `db:"sort,omitempty" json:"sort"`
	Owner         int64     `db:"owner,omitempty" json:"owner"`
	VersionId     int64     `db:"version_id,omitempty" json:"versionId"`
	PlanStartTime time.Time `db:"plan_start_time,omitempty" json:"planStartTime"`
	PlanEndTime   time.Time `db:"plan_end_time,omitempty" json:"planEndTime"`
	PlanWorkHour  int       `db:"plan_work_hour,omitempty" json:"planWorkHour"`
	StoryPoint    int       `db:"story_point,omitempty" json:"storyPoint"`
	Remark        *string   `db:"remark,omitempty" json:"remark"`
	Status        int64     `db:"status,omitempty" json:"status"`
	Creator       int64     `db:"creator,omitempty" json:"creator"`
	CreateTime    time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator       int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime    time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version       int       `db:"version,omitempty" json:"version"`
	IsDelete      int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmPriIteration) TableName() string {
	return "ppm_pri_iteration"
}
