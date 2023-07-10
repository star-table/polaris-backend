package bo

import "time"

type CustomFieldBo struct {
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

func (*CustomFieldBo) TableName() string {
	return "ppm_pro_custom_field"
}

type CustomFieldItemFrontend struct {
	// 展示的自定义字段值
	Title   string      `json:"title"`
	Value   interface{} `json:"value"`
	FieldId int64       `json:"fieldId"`
}

type CustomFieldItemFrontendForLc struct {
	// 展示的自定义字段值
	Title   string      `json:"title"`
	Value   interface{} `json:"value"`
	FieldId string      `json:"fieldId"`
}

type CustomFieldPersonValue struct {
	Id string `json:"id"`
	UserId int64 `json:"userId"`
	Name string `json:"name"`
}
