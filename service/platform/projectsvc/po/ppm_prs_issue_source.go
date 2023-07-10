package po

import "time"

type PpmPrsIssueSource struct {
	Id                  int64     `db:"id,omitempty" json:"id"`
	OrgId               int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId           int64     `db:"project_id,omitempty" json:"projectId"`
	LangCode            string    `db:"lang_code,omitempty" json:"langCode"`
	Name                string    `db:"name,omitempty" json:"name"`
	Sort                int       `db:"sort,omitempty" json:"sort"`
	ProjectObjectTypeId int64     `db:"project_object_type_id,omitempty" json:"projectObjectTypeId"`
	Remark              string    `db:"remark,omitempty" json:"remark"`
	Status              int       `db:"status,omitempty" json:"status"`
	Creator             int64     `db:"creator,omitempty" json:"creator"`
	CreateTime          time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator             int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime          time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version             int       `db:"version,omitempty" json:"version"`
	IsDelete            int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmPrsIssueSource) TableName() string {
	return "ppm_prs_issue_source"
}
