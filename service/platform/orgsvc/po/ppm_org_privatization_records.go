package po

import "time"

type PpmOrgPrivatizationRecords struct {
	Id              int64     `db:"id,omitempty" json:"id"`
	AuthorityId     int64     `db:"authority_id,omitempty" json:"authorityId"`
	OrgUserCount    int64     `db:"org_user_count,omitempty" json:"orgUserCount"`
	OrgDeptCount    int64     `db:"org_dept_count,omitempty" json:"orgDeptCount"`
	OrgProjectCount int64     `db:"org_project_count,omitempty" json:"orgProjectCount"`
	OrgIssueCount   int64     `db:"org_issue_count,omitempty" json:"orgIssueCount"`
	OutSideIp       string    `db:"out_side_ip,omitempty" json:"outSideIp"`
	Creator         int64     `db:"creator,omitempty" json:"creator"`
	CreateTime      time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator         int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime      time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version         int       `db:"version,omitempty" json:"version"`
	IsDelete        int       `db:"is_delete" json:"isDelete"`
}

func (*PpmOrgPrivatizationRecords) TableName() string {
	return "ppm_org_privatization_records"
}
