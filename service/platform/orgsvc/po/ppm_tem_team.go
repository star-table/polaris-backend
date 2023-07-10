package po

import "time"

type PpmTemTeam struct {
	Id           int64     `db:"id,omitempty" json:"id"`
	OrgId        int64     `db:"org_id,omitempty" json:"orgId"`
	Name         string    `db:"name,omitempty" json:"name"`
	NickName     string    `db:"nick_name,omitempty" json:"nickName"`
	Owner        int64     `db:"owner,omitempty" json:"owner"`
	DepartmentId int64     `db:"department_id,omitempty" json:"departmentId"`
	IsDefault    int       `db:"is_default" json:"isDefault"`
	Remark       string    `db:"remark,omitempty" json:"remark"`
	Status       int       `db:"status" json:"status"`
	Creator      int64     `db:"creator,omitempty" json:"creator"`
	CreateTime   time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime   time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int       `db:"version,omitempty" json:"version"`
	IsDelete     int       `db:"is_delete" json:"isDelete"`
}

func (*PpmTemTeam) TableName() string {
	return "ppm_tem_team"
}
