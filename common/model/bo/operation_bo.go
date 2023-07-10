package bo

import "time"

type OperationBo struct {
	Id					int64		`db:"id,omitempty" json:"id"`
	Code					string		`db:"code,omitempty" json:"code"`
	Name					string		`db:"name,omitempty" json:"name"`
	Remark					string		`db:"remark,omitempty" json:"remark"`
	Status					int		`db:"status,omitempty" json:"status"`
	Creator					int64		`db:"creator,omitempty" json:"creator"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"createTime"`
	Updator					int64		`db:"updator,omitempty" json:"updator"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"updateTime"`
	Version					int		`db:"version,omitempty" json:"version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"isDelete"`
}

func (*OperationBo) TableName() string {
	return "ppm_rol_operation"
}