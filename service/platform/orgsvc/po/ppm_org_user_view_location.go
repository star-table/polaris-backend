package po

import "time"

type PpmOrgUserViewLocation struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	UserId     int64     `db:"user_id,omitempty" json:"userId"`
	AppId      int64     `db:"app_id,omitempty" json:"appId"`
	Config     string    `db:"config,omitempty" json:"config"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrgUserViewLocation) TableName() string {
	return "ppm_org_user_view_location"
}
