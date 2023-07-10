package po

import "time"

const TableNamePpmOrgThirdLoginUser = "ppm_org_third_login_user"

type PpmOrgThirdLoginUser struct {
	Id             int64     `db:"id,omitempty" json:"id"`
	UserId         int64     `db:"user_id,omitempty" json:"userId"`
	ThirdUserId    string    `db:"third_user_id,omitempty" json:"thirdUserId"`
	SourcePlatform string    `db:"source_platform,omitempty" json:"source_platform"`
	CreateTime     time.Time `db:"create_time,omitempty" json:"createTime"`
	UpdateTime     time.Time `db:"update_time,omitempty" json:"updateTime"`
}

func (*PpmOrgThirdLoginUser) TableName() string {
	return TableNamePpmOrgThirdLoginUser
}
