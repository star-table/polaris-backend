package po

import "time"

type PpmBasPayLevel struct {
	Id          int64     `db:"id,omitempty" json:"id"`
	LangCode    string    `db:"lang_code,omitempty" json:"langCode"`
	Name        string    `db:"name,omitempty" json:"name"`
	Storage     int64     `db:"storage,omitempty" json:"storage"`
	MemberCount int       `db:"member_count,omitempty" json:"memberCount"`
	Price       int64     `db:"price,omitempty" json:"price"`
	MemberPrice int64     `db:"member_price,omitempty" json:"memberPrice"`
	Duration    int64     `db:"duration,omitempty" json:"duration"`
	IsShow      int       `db:"is_show,omitempty" json:"isShow"`
	Sort        int       `db:"sort,omitempty" json:"sort"`
	Status      int       `db:"status,omitempty" json:"status"`
	Creator     int64     `db:"creator,omitempty" json:"creator"`
	CreateTime  time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator     int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime  time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version     int       `db:"version,omitempty" json:"version"`
	IsDelete    int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmBasPayLevel) TableName() string {
	return "ppm_bas_pay_level"
}
