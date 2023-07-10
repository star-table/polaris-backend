package po

import "time"

type LcAppView struct {
	ID         int64     `db:"id,omitempty" json:"id,omitempty"`
	OrgID      int64     `db:"org_id,omitempty" json:"orgId,omitempty"`
	AppId      int64     `db:"app_id,omitempty" json:"appId,omitempty"`
	Type       int8      `db:"type,omitempty" json:"type,omitempty"`
	ViewName   string    `db:"view_name,omitempty" json:"viewName,omitempty"`
	Config     string    `db:"config,omitempty" json:"config,omitempty"`
	Remark     string    `db:"remark,omitempty" json:"remark,omitempty"`
	Owner      int64     `db:"owner,omitempty" json:"owner,omitempty"`
	Sort       int64     `db:"sort,omitempty" json:"sort,omitempty"`
	Version    int32     `db:"version,omitempty" json:"version,omitempty"`
	DelFlag    int8      `db:"del_flag,omitempty" json:"delFlag,omitempty"`
	Creator    int64     `db:"creator,omitempty" json:"creator,omitempty"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime,omitempty"`
	Updator    int64     `db:"updator,omitempty" json:"updator,omitempty"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime,omitempty"`
}

func (*LcAppView) TableName() string {
	return "lc_app_view"
}