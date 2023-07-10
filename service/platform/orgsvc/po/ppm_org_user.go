package po

import "time"

type PpmOrgUser struct {
	Id                 int64     `db:"id,omitempty" json:"id"`
	OrgId              int64     `db:"org_id,omitempty" json:"orgId"`
	Name               string    `db:"name,omitempty" json:"name"`
	NamePinyin         string    `db:"name_pinyin,omitempty" json:"namePinyin"`
	LoginName          string    `db:"login_name,omitempty" json:"loginName"`
	LoginNameEditCount int       `db:"login_name_edit_count,omitempty" json:"loginNameEditCount"`
	Email              string    `db:"email,omitempty" json:"email"`
	MobileRegion       string    `db:"mobile_region,omitempty" json:"mobileRegion"`
	Mobile             string    `db:"mobile,omitempty" json:"mobile"`
	Avatar             string    `db:"avatar,omitempty" json:"avatar"`
	Birthday           time.Time `db:"birthday,omitempty" json:"birthday"`
	Sex                int       `db:"sex,omitempty" json:"sex"`
	Password           string    `db:"password,omitempty" json:"password"`
	PasswordSalt       string    `db:"password_salt,omitempty" json:"passwordSalt"`
	SourcePlatform     string    `db:"source_platform,omitempty" json:"sourcePlatform"`
	SourceChannel      string    `db:"source_channel,omitempty" json:"sourceChannel"`
	SourceObjId        string    `db:"source_obj_id,omitempty" json:"sourceObjId"`
	Language           string    `db:"language,omitempty" json:"language"`
	Motto              string    `db:"motto,omitempty" json:"motto"`
	LastLoginIp        string    `db:"last_login_ip,omitempty" json:"lastLoginIp"`
	LastLoginTime      time.Time `db:"last_login_time,omitempty" json:"lastLoginTime"`
	LoginFailCount     int       `db:"login_fail_count,omitempty" json:"loginFailCount"`
	Status             int       `db:"status,omitempty" json:"status"`
	RemindBindPhone    int       `db:"remind_bind_phone,omitempty" json:"remindBindPhone"`
	Creator            int64     `db:"creator,omitempty" json:"creator"`
	CreateTime         time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator            int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime         time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version            int       `db:"version,omitempty" json:"version"`
	IsDelete           int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrgUser) TableName() string {
	return "ppm_org_user"
}
