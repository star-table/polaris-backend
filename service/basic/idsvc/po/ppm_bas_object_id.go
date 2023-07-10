package po

import "time"

type PpmBasObjectId struct {
	Id         int64     `db:"id,omitempty" json:"id"`                  // 主键
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`           // 组织id
	Code       string    `db:"code,omitempty" json:"code"`              // 对象编号
	MaxId      int64     `db:"max_id,omitempty" json:"maxId"`           // 当前最大编号
	Step       int       `db:"step,omitempty" json:"step"`              // 步长
	Creator    int64     `db:"creator,omitempty" json:"creator"`        // 创建人
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"` // 创建时间
	Updator    int64     `db:"updator,omitempty" json:"updator"`        // 更新人
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"` // 更新时间
	Version    int       `db:"version,omitempty" json:"version"`        // 乐观锁
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`     // 是否删除,1是,0否
}

func (*PpmBasObjectId) TableName() string {
	return "ppm_bas_object_id"
}
