package po

import "time"

type LcApp struct {
	ID         int64     `db:"id,omitempty" json:"id,omitempty"`
	OrgID      int64     `db:"org_id,omitempty" json:"orgId,omitempty"`
	PkgId      int64     `db:"pkg_id,omitempty" json:"pkgId,omitempty"`
	Type       int8      `db:"type,omitempty" json:"type,omitempty"`
	Name   string    `db:"name,omitempty" json:"name,omitempty"`
	Icon   string    `db:"icon,omitempty" json:"icon,omitempty"`
	Status   int    `db:"status,omitempty" json:"status,omitempty"`
	Config     string    `db:"config,omitempty" json:"config,omitempty"`
	Remark     *string    `db:"remark,omitempty" json:"remark,omitempty"`
	Owner      int64     `db:"owner,omitempty" json:"owner,omitempty"`
	Sort       int64     `db:"sort,omitempty" json:"sort,omitempty"`
	Version    int32     `db:"version,omitempty" json:"version,omitempty"`
	DelFlag    int8      `db:"del_flag,omitempty" json:"delFlag,omitempty"`
	Creator    int64     `db:"creator,omitempty" json:"creator,omitempty"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime,omitempty"`
	Updator    int64     `db:"updator,omitempty" json:"updator,omitempty"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime,omitempty"`
	WorkflowFlag int `db:"workflow_flag,omitempty" json:"workflowFlag,omitempty"`
	TemplateFlag int `db:"template_flag,omitempty" json:"templateFlag,omitempty"`
	ExternalApp int `db:"external_app,omitempty" json:"externalApp,omitempty"`
	LinkUrl int `db:"external_app,omitempty" json:"externalApp,omitempty"`
}

func (*LcApp) TableName() string {
	return "lc_app"
}