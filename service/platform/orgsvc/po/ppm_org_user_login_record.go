package po

import "time"

type PpmOrgUserLoginRecord struct {
	Id             int64     `db:"id,omitempty" json:"id"`
	OrgId          int64     `db:"org_id,omitempty" json:"orgId"`
	UserId         int64     `db:"user_id,omitempty" json:"userId"`
	LoginIp        string    `db:"login_ip,omitempty" json:"loginIp"`
	SourcePlatform string    `db:"source_platform,omitempty" json:"sourcePlatform"`
	SourceChannel  string    `db:"source_channel,omitempty" json:"sourceChannel"`
	LoginTime      time.Time `db:"login_time,omitempty" json:"loginTime"`
	Creator        int64     `db:"creator,omitempty" json:"creator"`
	CreateTime     time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator        int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime     time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version        int       `db:"version,omitempty" json:"version"`
	IsDelete       int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrgUserLoginRecord) TableName() string {
	return "ppm_org_user_login_record"
}
