package bo

import "time"

type RoleGroupBo struct {
	Id					int64		`db:"id,omitempty" json:"id"`
	OrgId					int64		`db:"org_id,omitempty" json:"orgId"`
	LangCode					string		`db:"lang_code,omitempty" json:"langCode"`
	Name					string		`db:"name,omitempty" json:"name"`
	Remark					string		`db:"remark,omitempty" json:"remark"`
	Type					int		`db:"type,omitempty" json:"type"`
	IsReadonly					int		`db:"is_readonly,omitempty" json:"isReadonly"`
	IsShow					int		`db:"is_show,omitempty" json:"isShow"`
	IsDefault					int		`db:"is_default,omitempty" json:"isDefault"`
	Status					int		`db:"status,omitempty" json:"status"`
	Creator					int64		`db:"creator,omitempty" json:"creator"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"createTime"`
	Updator					int64		`db:"updator,omitempty" json:"updator"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"updateTime"`
	Version					int		`db:"version,omitempty" json:"version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"isDelete"`
}

func (*RoleGroupBo) TableName() string {
	return "ppm_rol_role_group"
}