package bo

import "time"

type OrgOutInfoBo struct {
	Id              int64     `db:"id,omitempty" json:"id"`
	OrgId           int64     `db:"org_id,omitempty" json:"orgId"`
	OutOrgId        string    `db:"out_org_id,omitempty" json:"outOrgId"`
	SourceChannel   string    `db:"source_channel,omitempty" json:"sourceChannel"`
	SourcePlatform  string    `db:"source_platform,omitempty" json:"sourcePlatform"`
	TenantCode      string    `db:"tenant_code,omitempty" json:"tenantCode"`
	Name            string    `db:"name,omitempty" json:"name"`
	Industry        string    `db:"industry,omitempty" json:"industry"`
	IsAuthenticated int       `db:"is_authenticated,omitempty" json:"isAuthenticated"`
	AuthTicket      string    `db:"auth_ticket,omitempty" json:"authTicket"`
	AuthLevel       string    `db:"auth_level,omitempty" json:"authLevel"`
	Status          int       `db:"status" json:"status"`
	Creator         int64     `db:"creator,omitempty" json:"creator"`
	CreateTime      time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator         int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime      time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version         int       `db:"version,omitempty" json:"version"`
	IsDelete        int       `db:"is_delete" json:"isDelete"`
}

func (*OrgOutInfoBo) TableName() string {
	return "ppm_org_organization_out_info"
}
