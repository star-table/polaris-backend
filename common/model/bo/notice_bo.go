package bo

import "time"

type NoticeBo struct {
	Id           int64     `db:"id,omitempty" json:"id"`
	OrgId        int64     `db:"org_id,omitempty" json:"orgId"`
	Type         int       `db:"type,omitempty" json:"type"`
	ProjectId    int64     `db:"project_id,omitempty" json:"projectId"`
	IssueId      int64     `db:"issue_id,omitempty" json:"issueId"`
	TrendsId     int64     `db:"trends_id,omitempty" json:"trendsId"`
	Content      string    `db:"content,omitempty" json:"content"`
	Noticer      int64     `db:"noticer,omitempty" json:"noticer"`
	Status       int       `db:"status,omitempty" json:"status"`
	RelationType string    `db:"relation_type,omitempty" json:"relationType"`
	Ext          string    `db:"ext,omitempty" json:"ext"`
	Creator      int64     `db:"creator,omitempty" json:"creator"`
	CreateTime   time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime   time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int       `db:"version,omitempty" json:"version"`
	IsDelete     int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*NoticeBo) TableName() string {
	return "ppm_tre_notice"
}
