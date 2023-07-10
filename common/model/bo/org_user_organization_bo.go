package bo

import "time"

type PpmOrgUserOrganizationBo struct {
	Id               int64     `db:"id,omitempty" json:"id"`
	OrgId            int64     `db:"org_id,omitempty" json:"orgId"`
	UserId           int64     `db:"user_id,omitempty" json:"userId"`
	CheckStatus      int       `db:"check_status,omitempty" json:"checkStatus"`
	UseStatus        int       `db:"use_status,omitempty" json:"useStatus"`
	Status           int       `db:"status,omitempty" json:"status"`
	StatusChangeTime time.Time `db:"status_change_time,omitempty" json:"statusChangeTime"`
	AuditorId        int64     `db:"auditor_id,omitempty" json:"auditorId"`
	AuditTime        time.Time `db:"audit_time,omitempty" json:"auditTime"`
	Creator          int64     `db:"creator,omitempty" json:"creator"`
	CreateTime       time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator          int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime       time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version          int       `db:"version,omitempty" json:"version"`
	IsDelete         int       `db:"is_delete,omitempty" json:"isDelete"`
}

type OrgMemberChangeBo struct {
	//变动类型：1 禁用，2 启用，3 新加入组织(正常)，4 从组织移除,5:新加入组织（禁用）
	ChangeType int `json:"changeType"`
	//组织id
	OrgId int64 `json:"orgId"`
	//变动人员id
	UserId int64 `json:"userId"`
	// 所在部门
	DeptIds []string `json:"deptIds"`
	//变动OpenIds
	OpenId string `json:"openId"`
	// 新的outUserId,微信可以改一次
	NewOutUserId string `json:"newOutUserId"`
	//来源
	SourceChannel string `json:"sourceChannel"`
}

type FeishuHelpObjectBo struct {
	OrgId     int64    `json:"orgId"`
	TenantKey string   `json:"tenantKey"`
	OpenIds   []string `json:"openIds"`
}

type OrgMemberBaseInfoBo struct {
	OrgId            int64     `db:"org_id,omitempty" json:"orgId"`
	OrgName          string    `db:"org_name,omitempty" json:"orgName"` // OrgName
	UserId           int64     `db:"user_id,omitempty" json:"userId"`
	LoginName        string    `db:"login_name,omitempty" json:"loginName"`
	Name             string    `db:"name,omitempty" json:"name"`
	NamePinyin       string    `db:"name_pinyin,omitempty" json:"namePinyin"`
	Email            string    `db:"email,omitempty" json:"email"`
	MobileRegion     string    `db:"mobile_region,omitempty" json:"mobileRegion"`
	Mobile           string    `db:"mobile,omitempty" json:"mobile"`
	Sex              string    `db:"sex,omitempty" json:"sex"`           // Sex
	Birthday         time.Time `db:"birthday,omitempty" json:"birthday"` // Birthday
	Language         string    `db:"language,omitempty" json:"language"` // Language
	Avatar           string    `db:"avatar,omitempty" json:"avatar"`
	EmpNo            string    `db:"emp_no,omitempty" json:"empNo"`       // 工号
	WeiboIds         string    `db:"weibo_ids,omitempty" json:"weiboIds"` // 微博ID array string ,号分割
	CheckStatus      int       `db:"check_status,omitempty" json:"checkStatus"`
	Status           int       `db:"status,omitempty" json:"status"`
	Type             int       `db:"type,omitempty" json:"type"`
	StatusChangeTime time.Time `db:"status_change_time,omitempty" json:"statusChangeTime"` // status修改时间
	UseStatus        int       `db:"use_status,omitempty" json:"useStatus"`                // 使用过该组织
	AuditorId        int64     `db:"auditor_id,omitempty" json:"auditorId"`                // 审核人
	AuditTime        time.Time `db:"audit_time,omitempty" json:"auditTime"`                // 审核时间
	OrgCreator       int64     `db:"org_creator,omitempty" json:"orgCreator"`              // 组织创建者
	OrgOwner         int64     `db:"org_owner,omitempty" json:"orgOwner"`                  // 组织拥有者
	Creator          int64     `db:"creator,omitempty" json:"creator"`
	CreateTime       time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator          int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime       time.Time `db:"update_time,omitempty" json:"updateTime"`
	IsDelete         int       `db:"is_delete,omitempty" json:"isDelete"` // IsDelete
}
