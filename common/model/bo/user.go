package bo

import (
	"time"

	"github.com/star-table/common/core/types"
)

type UserNoticeInfoBo struct {
	UserId       int64
	OutUserId    string
	Name         string
	OutOrgUserId string
}

type UserInfoBo struct {
	ID                 int64      `json:"id"`
	EmplID             *string    `json:"emplId"`
	OrgID              int64      `json:"orgId"`
	OrgName            string     `json:"orgName"`
	Name               string     `json:"name"`
	NamePinyin         string     `db:"name_pinyin,omitempty" json:"namePinyin"`
	ThirdName          string     `json:"thirdName"`
	LoginName          string     `json:"loginName"`
	LoginNameEditCount int        `json:"loginNameEditCount"`
	Email              string     `json:"email"`
	Mobile             string     `json:"mobile"`
	Birthday           types.Time `json:"birthday"`
	Sex                int        `json:"sex"`
	Rimanente          int        `json:"rimanente"`
	Level              int        `json:"level"`
	LevelName          string     `json:"levelName"`
	Password           string     `db:"password,omitempty" json:"password"`
	PasswordSalt       string     `db:"password_salt,omitempty" json:"passwordSalt"`
	Avatar             string     `json:"avatar"`
	SourceChannel      string     `json:"sourceChannel"`
	Language           string     `json:"language"`
	Motto              string     `json:"motto"`
	LastLoginIP        string     `json:"lastLoginIp"`
	LastLoginTime      types.Time `json:"lastLoginTime"`
	LoginFailCount     int        `json:"loginFailCount"`
	RemindBindPhone    int        `json:"remindBindPhone"`
	CreateTime         types.Time `json:"createTime"`
	UpdateTime         types.Time `json:"updateTime"`
	GlobalUserId       int64
	UserStatus         int // 是否启用
}

type UserOutInfoBo struct {
	Id            int64     `db:"id,omitempty" json:"id"`
	OrgId         int64     `db:"org_id,omitempty" json:"orgId"`
	UserId        int64     `db:"user_id,omitempty" json:"userId"`
	SourceChannel string    `db:"source_channel,omitempty" json:"sourceChannel"`
	OutOrgUserId  string    `db:"out_org_user_id,omitempty" json:"outOrgUserId"`
	OutUserId     string    `db:"out_user_id,omitempty" json:"outUserId"`
	Name          string    `db:"name,omitempty" json:"name"`
	Avatar        string    `db:"avatar,omitempty" json:"avatar"`
	IsActive      int       `db:"is_active,omitempty" json:"isActive"`
	JobNumber     string    `db:"job_number,omitempty" json:"jobNumber"`
	Status        int       `db:"status" json:"status"`
	Creator       int64     `db:"creator,omitempty" json:"creator"`
	CreateTime    time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator       int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime    time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version       int       `db:"version,omitempty" json:"version"`
	IsDelete      int       `db:"is_delete" json:"isDelete"`
}

type SimpleUserInfoBo struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Status int    `json:"status"`
}

type SimpleUserWithUserIdBo struct {
	Id     int64  `json:"id"` // 为适配前端格式，先用 string。
	UserId int64  `json:"userId"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Status int    `json:"status"`
	Type   string `json:"type"`
}

type UserEmpIdInfo struct {
	EmpId  string `json:"empId"`
	UserId int64  `json:"userId"`
}

type UserSMSRegisterInfo struct {
	OrgId          int64  `json:"orgId"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	SourceChannel  string `json:"sourceChannel"`
	SourcePlatform string `json:"sourcePlatform"`
	Name           string `json:"name"`
	InviteCode     string `json:"inviteCode"`
	// 手机归属代码。如中国为 `+86`。如果有，则使用
	MobileRegion string `json:"mobileRegion"`
	Password     string `json:"password"`
	Status       int    `json:"status"`
	AccountName  string `json:"accountName"`
}

type OrgUserInfo struct {
	UserId    int64  `db:"user_id" json:"userId"`
	OutUserId string `db:"out_user_id" json:"outUserId"` //有可能为空
	OrgId     int64  `db:"org_id" json:"orgId"`

	OrgUserStatus      int `db:"org_user_status" json:"orgUserStatus"`            //用户组织状态
	OrgUserCheckStatus int `db:"org_user_check_status" json:"orgUserCheckStatus"` //用户组织审核状态
}
type OrgProjectMemberBo struct {
	ID                 int64      `json:"id"`
	OrgID              int64      `json:"orgId"`
	OrgName            string     `json:"orgName"`
	Name               string     `json:"name"`
	LoginName          string     `json:"loginName"`
	LoginNameEditCount int        `json:"loginNameEditCount"`
	Email              string     `json:"email"`
	Mobile             string     `json:"mobile"`
	Birthday           types.Time `json:"birthday"`
	Sex                int        `json:"sex"`
	Level              int        `json:"level"`
	LevelName          string     `json:"levelName"`
	SourceChannel      string     `json:"sourceChannel"`
	OutOrgUserId       string     `json:"outOrgUserId"`
	OutUserId          string     `json:"outUserId"`
	OutUserName        string     `json:"outUserName"`
}

type CacheUserCheckPayBo struct {
	UserId          int64     `json:"userId"`
	ServiceStopTime time.Time `json:"serviceStopTime"`
	Status          string    `json:"status"`
}

type UserOrDeptInfoItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type PlatformAdminBo struct {
	OrgId         int64  `json:"orgId"`
	UserId        int64  `json:"userId"`
	OutOrgId      string `json:"outOrgId"`
	OpenId        string `json:"openId"`
	SourceChannel string `json:"sourceChannel"`
	IsAdmin       int    `json:"isAdmin"`
}
