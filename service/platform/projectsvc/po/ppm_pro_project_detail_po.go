package po

import "time"

type PpmProProjectDetail struct {
	Id                int64     `db:"id,omitempty" json:"id"`
	OrgId             int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId         int64     `db:"project_id,omitempty" json:"projectId"`
	Notice            string    `db:"notice,omitempty" json:"notice"`
	IsEnableWorkHours int       `db:"is_enable_work_hours,omitempty" json:"isEnableWorkHours"`
	IsSyncOutCalendar int       `db:"is_sync_out_calendar,omitempty" json:"isSyncOutCalendar"`
	Creator           int64     `db:"creator,omitempty" json:"creator"`
	CreateTime        time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator           int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime        time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version           int       `db:"version,omitempty" json:"version"`
	IsDelete          int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmProProjectDetail) TableName() string {
	return "ppm_pro_project_detail"
}
