package po

import "time"

type PpmProChat struct {
	Id             int64     `db:"id,omitempty" json:"id"`
	SourceChannel  string    `db:"source_channel,omitempty" json:"sourceChannel"`
	SourcePlatform string    `db:"source_platform,omitempty" json:"sourcePlatform"`
	OrgId          int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId      int64     `db:"project_id,omitempty" json:"projectId"`
	Name           string    `db:"name,omitempty" json:"name"`
	Avatar         string    `db:"avatar,omitempty" json:"avatar"`
	Description    string    `db:"description,omitempty" json:"description"`
	OutChatId      string    `db:"out_chat_id,omitempty" json:"outChatId"`
	Creator        int64     `db:"creator,omitempty" json:"creator"`
	CreateTime     time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator        int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime     time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version        int       `db:"version,omitempty" json:"version"`
	IsDelete       int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmProChat) TableName() string {
	return "ppm_pro_chat"
}
