package po

import "time"

const TableNameShareView = "ppm_pri_share_view"

type PpmShareView struct {
	Id            int64     `gorm:"column:id" db:"id,omitempty" json:"id" form:"id"`
	ShareKey      string    `gorm:"column:share_key" db:"share_key" json:"share_key" form:"share_key"`
	SharePassword string    `gorm:"column:share_password" db:"share_password" json:"share_password" form:"share_password"`
	OrgId         int64     `gorm:"column:org_id" db:"org_id" json:"org_id" form:"org_id"`
	UserId        int64     `gorm:"column:user_id" db:"user_id" json:"user_id" form:"user_id"`
	ProjectId     int64     `gorm:"column:project_id" db:"project_id" json:"project_id" form:"project_id"`
	AppId         int64     `gorm:"column:app_id" db:"app_id" json:"app_id" form:"app_id"`
	ViewId        int64     `gorm:"column:view_id" db:"view_id" json:"view_id" form:"view_id"`
	TableId       int64     `gorm:"column:table_id" db:"table_id" json:"table_id" form:"table_id"`
	Config        string    `gorm:"column:config" db:"config" json:"config" form:"config"`
	CreateTime    time.Time `gorm:"column:create_time" db:"create_time,omitempty" json:"create_time" form:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time" db:"update_time,omitempty" json:"update_time" form:"update_time"`
}

func (*PpmShareView) TableName() string {
	return TableNameShareView
}
