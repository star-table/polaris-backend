package bo

import "time"

type TagBo struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	Name       string    `db:"name,omitempty" json:"name"`
	ProjectId  int64	 `db:"project_id,omitempty" json:"projectId"`
	NamePinyin string    `db:"name_pinyin,omitempty" json:"namePinyin"`
	BgStyle    string    `db:"bg_style,omitempty" json:"bgStyle"`
	FontStyle  string    `db:"font_style,omitempty" json:"fontStyle"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

type IssueTagBo struct {
	IssueId   int64
	TagId     int64
	TagName   string
	BgStyle   string
	FontStyle string
}

func (*TagBo) TableName() string {
	return "ppm_pri_tag"
}

type IssueTagStatBo struct {
	TagId int64 `json:"tagId"`
	Total int64 `json:"total"`
}