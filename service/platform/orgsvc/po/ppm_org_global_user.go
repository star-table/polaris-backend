package po

import "time"

const TableNamePpmOrgGlobalUser = "ppm_org_global_user"

type PpmOrgGlobalUser struct {
	Id              int64     `db:"id,omitempty" json:"id"`
	Mobile          string    `db:"mobile,omitempty" json:"mobile"`
	LastLoginUserId int64     `db:"last_login_user_id,omitempty" json:"last_login_user_id"`
	LastLoginOrgId  int64     `db:"last_login_org_id,omitempty" json:"last_login_org_id"`
	CreateTime      time.Time `db:"create_time,omitempty" json:"createTime"`
	UpdateTime      time.Time `db:"update_time,omitempty" json:"updateTime"`
	IsDelete        int       `db:"is_delete,omitempty" json:"isDelete"`
	Password        string    `db:"password,omitempty" json:"password"`
	PasswordSalt    string    `db:"password_salt,omitempty" json:"passwordSalt"`
	WeiXinOpenId    string    `db:"wei_xin_open_id,omitempty" json:"wei_xin_open_id"`
}

func (*PpmOrgGlobalUser) TableName() string {
	return TableNamePpmOrgGlobalUser
}
