package po

import "time"

type PpmProProjectChat struct {
	Id           int64  `db:"id,omitempty" json:"id"`
	OrgId        int64  `db:"org_id,omitempty" json:"orgId"`
	ProjectId    int64  `db:"project_id,omitempty" json:"projectId"`
	TableId      int64  `db:"table_id,omitempty" json:"tableId"`
	ChatId       string `db:"chat_id,omitempty" json:"chatId"`
	ChatType     int    `db:"chat_type,omitempty" json:"chatType"` // 该字段可忽略，主动创建的群聊关系存储在 project relation 表中
	ChatSettings string `db:"chat_settings,omitempty" json:"chatSettings"`
	IsEnable     int    `db:"is_enable,omitempty" json:"isEnable"`

	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmProProjectChat) TableName() string {
	return "ppm_pro_project_chat"
}

type ProjectChatObj struct {
	Id        int64  `db:"id,omitempty" json:"id"`
	OrgId     int64  `db:"org_id,omitempty" json:"orgId"`
	ProjectId int64  `db:"project_id,omitempty" json:"projectId"`
	ChatId    string `db:"chat_id,omitempty" json:"chatId"`
	ChatType  int    `db:"chat_type,omitempty" json:"chatType"`
}
