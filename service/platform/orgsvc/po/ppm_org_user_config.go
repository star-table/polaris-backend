package po

import "time"

type PpmOrgUserConfig struct {
	Id                              int64     `db:"id,omitempty" json:"id"`
	OrgId                           int64     `db:"org_id,omitempty" json:"orgId"`
	UserId                          int64     `db:"user_id,omitempty" json:"userId"`
	DailyReportMessageStatus        int       `db:"daily_report_message_status,omitempty" json:"dailyReportMessageStatus"`
	OwnerRangeStatus                int       `db:"owner_range_status,omitempty" json:"ownerRangeStatus"`
	ParticipantRangeStatus          int       `db:"participant_range_status,omitempty" json:"participantRangeStatus"`
	AttentionRangeStatus            int       `db:"attention_range_status,omitempty" json:"attentionRangeStatus"`
	CreateRangeStatus               int       `db:"create_range_status,omitempty" json:"createRangeStatus"`
	RemindMessageStatus             int       `db:"remind_message_status,omitempty" json:"remindMessageStatus"`
	CommentAtMessageStatus          int       `db:"comment_at_message_status,omitempty" json:"commentAtMessageStatus"`
	ModifyMessageStatus             int       `db:"modify_message_status,omitempty" json:"modifyMessageStatus"`
	RelationMessageStatus           int       `db:"relation_message_status,omitempty" json:"relationMessageStatus"`
	DailyProjectReportMessageStatus int       `db:"daily_project_report_message_status,omitempty" json:"dailyProjectReportMessageStatus"`
	DefaultProjectId                int64     `db:"default_project_id,omitempty" json:"defaultProjectId"`
	DefaultProjectObjectTypeId      int64     `db:"default_project_object_type_id,omitempty" json:"defaultProjectObjectTypeId"`
	PcNoticeOpenStatus              int       `db:"pc_notice_open_status,omitempty" json:"pcNoticeOpenStatus"`
	PcIssueRemindMessageStatus      int       `db:"pc_issue_remind_message_status,omitempty" json:"pcIssueRemindMessageStatus"`
	PcOrgMessageStatus              int       `db:"pc_org_message_status,omitempty" json:"pcOrgMessageStatus"`
	PcProjectMessageStatus          int       `db:"pc_project_message_status,omitempty" json:"pcProjectMessageStatus"`
	PcCommentAtMessageStatus        int       `db:"pc_comment_at_message_status,omitempty" json:"pcCommentAtMessageStatus"`
	CollaborateMessageStatus        int       `db:"collaborate_message_status,omitempty" json:"collaborateMessageStatus"`
	UserViewLocationConfig          *string   `db:"user_view_location_config,omitempty" json:"userViewLocationConfig"`
	Ext                             string    `db:"ext,omitempty" json:"ext"`
	Creator                         int64     `db:"creator,omitempty" json:"creator"`
	CreateTime                      time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator                         int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime                      time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version                         int       `db:"version,omitempty" json:"version"`
	IsDelete                        int       `db:"is_delete,omitempty" json:"isDelete"`
	RemindExpiring                  int       `db:"remind_expiring,omitempty" json:"remindExpiring"`
}

func (*PpmOrgUserConfig) TableName() string {
	return "ppm_org_user_config"
}
