package po

import "time"

type PpmPriTag struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId  int64     `db:"project_id,omitempty" json:"projectId"`
	Name       string    `db:"name,omitempty" json:"name"`
	NamePinyin string    `db:"name_pinyin,omitempty" json:"namePinyin"`
	BgStyle    string    `db:"bg_style,omitempty" json:"bgStyle"`
	FontStyle  string    `db:"font_style,omitempty" json:"fontStyle"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmPriTag) TableName() string {
	return "ppm_pri_tag"
}

type TagInfoWithIssue struct {
	IssueId   int64  `db:"issue_id,omitempty" json:"issueId"`
	Name      string `db:"name,omitempty" json:"name"`
	BgStyle   string `db:"bg_style,omitempty" json:"bgStyle"`
	FontStyle string `db:"font_style,omitempty" json:"fontStyle"`
	TagId     int64  `db:"tag_id,omitempty" json:"tagId"`
}

type IssueTagStat struct {
	Total int64 `db:"total" json:"total"`
	TagId int64 `db:"tag_id" json:"tagId"`
}
