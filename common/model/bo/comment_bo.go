package bo

import "time"

type CommentBo struct {
	Id					int64		`db:"id,omitempty" json:"id"`
	OrgId					int64		`db:"org_id,omitempty" json:"orgId"`
	ProjectId					int64		`db:"project_id,omitempty" json:"projectId"`
	TrendsId					int64		`db:"trends_id,omitempty" json:"trendsId"`
	ObjectId					int64		`db:"object_id,omitempty" json:"objectId"`
	ObjectType					string		`db:"object_type,omitempty" json:"objectType"`
	Content					string		`db:"content,omitempty" json:"content"`
	ParentId					int64		`db:"parent_id,omitempty" json:"parentId"`
	Creator					int64		`db:"creator,omitempty" json:"creator"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"createTime"`
	Updator					int64		`db:"updator,omitempty" json:"updator"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"updateTime"`
	Version					int		`db:"version,omitempty" json:"version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"isDelete"`
}

func (*CommentBo) TableName() string {
	return "ppm_tre_comment"
}