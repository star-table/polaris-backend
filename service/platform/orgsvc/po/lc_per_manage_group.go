package po

import "time"

type LcPerManageGroup struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	LangCode   string    `db:"lang_code,omitempty" json:"langCode"`
	Name       string    `db:"name,omitempty" json:"name"`
	UserIds    *string   `db:"user_ids,omitempty" json:"userIds"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*LcPerManageGroup) TableName() string {
	return "lc_per_manage_group"
}
