package po

import (
	"time"
)

type PpmProProject struct {
	Id            int64     `db:"id,omitempty" json:"id"`
	AppId         int64     `db:"app_id,omitempty" json:"appId,string"`
	OrgId         int64     `db:"org_id,omitempty" json:"orgId"`
	Code          string    `db:"code,omitempty" json:"code"`
	Name          string    `db:"name,omitempty" json:"name"`
	PreCode       string    `db:"pre_code,omitempty" json:"preCode"`
	Owner         int64     `db:"owner,omitempty" json:"owner"`
	ProjectTypeId int64     `db:"project_type_id,omitempty" json:"projectTypeId"`
	PriorityId    int64     `db:"priority_id,omitempty" json:"priorityId"`
	PlanStartTime time.Time `db:"plan_start_time,omitempty" json:"planStartTime"`
	PlanEndTime   time.Time `db:"plan_end_time,omitempty" json:"planEndTime"`
	PublicStatus  int       `db:"public_status,omitempty" json:"publicStatus"`
	TemplateFlag  int       `db:"template_flag,omitempty" json:"templateFlag"`
	ResourceId    int64     `db:"resource_id,omitempty" json:"resourceId"`
	IsFiling      int       `db:"is_filing,omitempty" json:"isFiling"`
	Remark        string    `db:"remark,omitempty" json:"remark"`
	Status        int64     `db:"status,omitempty" json:"status"`
	Creator       int64     `db:"creator,omitempty" json:"creator"`
	CreateTime    time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator       int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime    time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version       int       `db:"version,omitempty" json:"version"`
	IsDelete      int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmProProject) TableName() string {
	return "ppm_pro_project"
}
