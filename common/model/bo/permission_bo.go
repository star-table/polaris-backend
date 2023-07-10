package bo

import "time"

type PermissionBo struct {
	Id					int64		`db:"id,omitempty" json:"id"`
	OrgId					int64		`db:"org_id,omitempty" json:"orgId"`
	LangCode					string		`db:"lang_code,omitempty" json:"langCode"`
	Code					string		`db:"code,omitempty" json:"code"`
	Name					string		`db:"name,omitempty" json:"name"`
	ParentId					int64		`db:"parent_id,omitempty" json:"parentId"`
	Type					int		`db:"type,omitempty" json:"type"`
	Path					string		`db:"path,omitempty" json:"path"`
	IsShow					int		`db:"is_show,omitempty" json:"isShow"`
	Remark					string		`db:"remark,omitempty" json:"remark"`
	Status					int		`db:"status,omitempty" json:"status"`
	Creator					int64		`db:"creator,omitempty" json:"creator"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"createTime"`
	Updator					int64		`db:"updator,omitempty" json:"updator"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"updateTime"`
	Version					int		`db:"version,omitempty" json:"version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"isDelete"`
}

func (*PermissionBo) TableName() string {
	return "ppm_rol_permission"
}