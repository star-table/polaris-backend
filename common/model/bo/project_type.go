package bo

import "time"

type ProjectTypeBo struct {
	Id               int64     `db:"id,omitempty" json:"id"`
	OrgId            int64     `db:"org_id,omitempty" json:"orgId"`
	LangCode         string    `db:"lang_code,omitempty" json:"langCode"`
	Name             string    `db:"name,omitempty" json:"name"`
	Sort             int       `db:"sort,omitempty" json:"sort"`
	Cover            string    `db:"cover,omitempty" json:"cover"`
	DefaultProcessId int64     `db:"default_process_id,omitempty" json:"defaultProcessId"`
	Category         int64     `db:"category,omitempty" json:"category"`
	Mode             int       `db:"mode,omitempty" json:"mode"`
	IsReadonly       int       `db:"is_readonly,omitempty" json:"isReadonly"`
	Remark           string    `db:"remark,omitempty" json:"remark"`
	Status           int       `db:"status,omitempty" json:"status"`
	Creator          int64     `db:"creator,omitempty" json:"creator"`
	CreateTime       time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator          int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime       time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version          int       `db:"version,omitempty" json:"version"`
	IsDelete         int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*ProjectTypeBo) TableName() string {
	return "ppm_prs_project_type"
}
