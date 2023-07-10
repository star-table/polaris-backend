package bo

import "time"

type AppInfoBo struct {
	Id          int64     `db:"id,omitempty" json:"id"`
	Name        string    `db:"name,omitempty" json:"name"`
	Code        string    `db:"code,omitempty" json:"code"`
	Secret1     string    `db:"secret1,omitempty" json:"secret1"`
	Secret2     string    `db:"secret2,omitempty" json:"secret2"`
	Owner       string    `db:"owner,omitempty" json:"owner"`
	CheckStatus int       `db:"check_status,omitempty" json:"checkStatus"`
	Remark      string    `db:"remark,omitempty" json:"remark"`
	Status      int       `db:"status,omitempty" json:"status"`
	Creator     int64     `db:"creator,omitempty" json:"creator"`
	CreateTime  time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator     int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime  time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version     int       `db:"version,omitempty" json:"version"`
	IsDelete    int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*AppInfoBo) TableName() string {
	return "ppm_bas_app_info"
}
