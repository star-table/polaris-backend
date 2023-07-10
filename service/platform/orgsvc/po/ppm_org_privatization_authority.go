package po

import "time"

type PpmOrgPrivatizationAuthority struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	Secret     string    `db:"secret,omitempty" json:"secret"`
	OrgName    string    `db:"org_name,omitempty" json:"orgName"`
	StartTime  time.Time `db:"start_time,omitempty" json:"startTime"`
	EndTime    time.Time `db:"end_time,omitempty" json:"endTime"`
	Status     int       `db:"status" json:"status"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete" json:"isDelete"`
}

func (*PpmOrgPrivatizationAuthority) TableName() string {
	return "ppm_org_privatization_authority"
}
