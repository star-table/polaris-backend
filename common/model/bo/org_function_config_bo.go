package bo

import "time"

type PpmOrcFunctionConfig struct {
	Id           int64     `db:"id,omitempty" json:"id"`
	OrgId        int64     `db:"org_id,omitempty" json:"orgId"`
	FunctionCode string    `db:"function_code,omitempty" json:"functionCode"`
	Remark       string    `db:"remark,omitempty" json:"remark"`
	Creator      int64     `db:"creator,omitempty" json:"creator"`
	CreateTime   time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime   time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int       `db:"version,omitempty" json:"version"`
	IsDelete     int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrcFunctionConfig) TableName() string {
	return "ppm_orc_function_config"
}

