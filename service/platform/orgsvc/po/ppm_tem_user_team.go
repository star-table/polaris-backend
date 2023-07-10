package po

import "time"

type PpmTemUserTeam struct {
	Id           int64     `db:"id,omitempty" json:"id"`
	OrgId        int64     `db:"org_id,omitempty" json:"orgId"`
	TeamId       int64     `db:"team_id,omitempty" json:"teamId"`
	UserId       int64     `db:"user_id,omitempty" json:"userId"`
	RelationType int       `db:"relation_type,omitempty" json:"relationType"`
	Status       int       `db:"status,omitempty" json:"status"`
	Creator      int64     `db:"creator,omitempty" json:"creator"`
	CreateTime   time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime   time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int       `db:"version,omitempty" json:"version"`
	IsDelete     int       `db:"is_delete" json:"isDelete"`
}

func (*PpmTemUserTeam) TableName() string {
	return "ppm_tem_user_team"
}
