package bo

import "time"

type ShareBo struct {
	Id					int64		`db:"id,omitempty" json:"id"`
	OrgId					int64		`db:"org_id,omitempty" json:"orgId"`
	ProjectId					int64		`db:"project_id,omitempty" json:"projectId"`
	Name					string		`db:"name,omitempty" json:"name"`
	Remark					string		`db:"remark,omitempty" json:"remark"`
	Logo					string		`db:"logo,omitempty" json:"logo"`
	Type					int		`db:"type,omitempty" json:"type"`
	Content					*string		`db:"content,omitempty" json:"content"`
	ContentMd5					string		`db:"content_md5,omitempty" json:"contentMd5"`
	FinishTime					time.Time		`db:"finish_time,omitempty" json:"finishTime"`
	Status					int		`db:"status,omitempty" json:"status"`
	Creator					int64		`db:"creator,omitempty" json:"creator"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"createTime"`
	Updator					int64		`db:"updator,omitempty" json:"updator"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"updateTime"`
	Version					int		`db:"version,omitempty" json:"version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"isDelete"`
}

func (*ShareBo) TableName() string {
	return "ppm_sha_share"
}