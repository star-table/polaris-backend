package po

import (
	"time"
)

type PpmOrgSecret struct{
	Id int64 `json:"id" db:"id,omitempty"`
	OrgId int64 `json:"orgId" db:"org_id,omitempty"`
	Key string `json:"key" db:"key,omitempty"`
	Secret string `json:"secret" db:"secret,omitempty"`
	Creator int64 `json:"creator" db:"creator,omitempty"`
	CreateTime time.Time `json:"createTime" db:"create_time,omitempty"`
	Updator int64 `json:"updator" db:"updator,omitempty"`
	UpdateTime time.Time `json:"updateTime" db:"update_time,omitempty"`
	Version int `json:"version" db:"version,omitempty"`
	IsDelete int `json:"isDelete" db:"is_delete,omitempty"`
}

func (*PpmOrgSecret) TableName() string {
	return "ppm_org_secret"
}