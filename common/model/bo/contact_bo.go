package bo

import "time"

type ContactBo struct {
	Id           int64     `db:"id,omitempty" json:"id"`
	Code         string    `db:"code,omitempty" json:"code"`
	Name         string    `db:"name,omitempty" json:"name"`
	Sex          int       `db:"sex,omitempty" json:"sex"`
	Email        string    `db:"email,omitempty" json:"email"`
	MobileRegion string    `db:"mobile_region,omitempty" json:"mobileRegion"`
	Mobile       string    `db:"mobile,omitempty" json:"mobile"`
	Remark       string    `db:"remark,omitempty" json:"remark"`
	Intention    int       `db:"intention,omitempty" json:"intention"`
	Source       string    `db:"source,omitempty" json:"source"`
	ResourceInfo string    `db:"resource_info,omitempty" json:"resourceInfo"`
	Status       int       `db:"status,omitempty" json:"status"`
	Creator      int64     `db:"creator,omitempty" json:"creator"`
	CreateTime   time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime   time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int       `db:"version,omitempty" json:"version"`
	IsDelete     int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*ContactBo) TableName() string {
	return "ppm_wst_contact"
}
