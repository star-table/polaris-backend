package po

import "time"

const TableNamePpmOrgDingCoolApp = "ppm_org_ding_cool_app"

type PpmOrgDingCoolApp struct {
	Id             int64     `db:"id,omitempty" json:"id"`
	CorpId         string    `db:"corp_id,omitempty" json:"corp_id"`
	OrgId          int64     `db:"org_id,omitempty" json:"org_id"`
	AppId          int64     `db:"app_id,omitempty" json:"app_id"`
	ProjectId      int64     `db:"project_id,omitempty" json:"project_id"`
	ConversationId string    `db:"conversation_id,omitempty" json:"conversation_id"`
	RobotCode      string    `db:"robot_code,omitempty" json:"robot_code"`
	TopTraceId     string    `db:"top_trace_id,omitempty" json:"top_trace_id"`
	FirstTraceId   string    `db:"first_trace_id,omitempty" json:"first_trace_id"`
	CreateTime     time.Time `db:"create_time,omitempty" json:"create_time"`
	UpdateTime     time.Time `db:"update_time,omitempty" json:"update_time"`
}

func (*PpmOrgDingCoolApp) TableName() string {
	return TableNamePpmOrgDingCoolApp
}
