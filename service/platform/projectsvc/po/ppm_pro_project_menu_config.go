package po

import "time"

type PpmProProjectMenuConfig struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	AppId      int64     `db:"app_id,omitempty" json:"appId"`
	Config     string    `db:"config,omitempty" json:"config"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmProProjectMenuConfig) TableName() string {
	return "ppm_pro_project_menu_config"
}
