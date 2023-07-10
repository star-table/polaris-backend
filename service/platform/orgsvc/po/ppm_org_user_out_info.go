package po

import "time"

type PpmOrgUserOutInfo struct {
	Id             int64     `db:"id,omitempty" json:"id"`
	OrgId          int64     `db:"org_id,omitempty" json:"orgId"`
	UserId         int64     `db:"user_id,omitempty" json:"userId"`
	SourcePlatform string    `db:"source_platform,omitempty" json:"sourcePlatform"`
	SourceChannel  string    `db:"source_channel,omitempty" json:"sourceChannel"`
	OutOrgUserId   string    `db:"out_org_user_id,omitempty" json:"outOrgUserId"`
	OutUserId      string    `db:"out_user_id,omitempty" json:"outUserId"`
	Name           string    `db:"name,omitempty" json:"name"`
	Avatar         string    `db:"avatar,omitempty" json:"avatar"`
	IsActive       int       `db:"is_active,omitempty" json:"isActive"`
	JobNumber      string    `db:"job_number,omitempty" json:"jobNumber"`
	Status         int       `db:"status,omitempty" json:"status"`
	Creator        int64     `db:"creator,omitempty" json:"creator"`
	CreateTime     time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator        int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime     time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version        int       `db:"version,omitempty" json:"version"`
	IsDelete       int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrgUserOutInfo) TableName() string {
	return "ppm_org_user_out_info"
}
