package po

import "time"

type PpmOrgUserInvite struct {
	Id             int64     `db:"id,omitempty" json:"id"`
	OrgId          int64     `db:"org_id,omitempty" json:"orgId"`
	Name           string    `db:"name,omitempty" json:"name"`
	Email          string    `db:"email,omitempty" json:"email"`
	InviteUserId   int64     `db:"invite_user_id,omitempty" json:"inviteUserId"`
	IsRegister     int       `db:"is_register,omitempty" json:"isRegister"`
	LastInviteTime time.Time `db:"last_invite_time,omitempty" json:"lastInviteTime"`
	Creator        int64     `db:"creator,omitempty" json:"creator"`
	CreateTime     time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator        int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime     time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version        int       `db:"version,omitempty" json:"version"`
	IsDelete       int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrgUserInvite) TableName() string {
	return "ppm_org_user_invite"
}
