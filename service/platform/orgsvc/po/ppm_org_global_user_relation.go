package po

import "time"

const TableNamePpmOrgGlobalUserRelation = "ppm_org_global_user_relation"

type PpmOrgGlobalUserRelation struct {
	Id           int64     `db:"id,omitempty" json:"id"`
	GlobalUserId int64     `db:"global_user_id,omitempty" json:"global_user_id"`
	UserId       int64     `db:"user_id,omitempty" json:"user_id"`
	CreateTime   time.Time `db:"create_time,omitempty" json:"createTime"`
	UpdateTime   time.Time `db:"update_time,omitempty" json:"updateTime"`
	IsDelete     int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrgGlobalUserRelation) TableName() string {
	return TableNamePpmOrgGlobalUserRelation
}
