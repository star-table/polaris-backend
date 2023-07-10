package po

import "time"

type PpmProCustomField struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	Name       string    `db:"name,omitempty" json:"name"`
	FieldType  int       `db:"field_type,omitempty" json:"fieldType"`
	FieldValue string    `db:"field_value,omitempty" json:"fieldValue"`
	IsOrgField int       `db:"is_org_field,omitempty" json:"isOrgField"`
	Remark     string    `db:"remark,omitempty" json:"remark"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmProCustomField) TableName() string {
	return "ppm_pro_custom_field"
}

type CustomFieldStat struct {
	FieldId int64 `db:"field_id,omitempty" json:"fieldId"`
	Total   int64 `db:"total,omitempty" json:"total"`
}
