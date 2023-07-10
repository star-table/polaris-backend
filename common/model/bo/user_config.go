package bo

type UserConfigBo struct {
	ID                              int64  `json:"id"`
	UserId                          int64  `json:"userId"`
	DailyReportMessageStatus        int    `json:"dailyReportMessageStatus"`        // 个人日报
	OwnerRangeStatus                int    `json:"ownerRangeStatus"`                // 我负责的
	ParticipantRangeStatus          int    `json:"participantRangeStatus"`          // 我关注的 // 该配置，前端请求接口的 key 有误导致，语义化有问题。
	AttentionRangeStatus            int    `json:"attentionRangeStatus"`            // [页面无此配置]
	CreateRangeStatus               int    `json:"createRangeStatus"`               // [页面无此配置] 猜测是：我创建的
	RemindMessageStatus             int    `json:"remindMessageStatus"`             // 任务提醒，被添加到记录
	CommentAtMessageStatus          int    `json:"commentAtMessageStatus"`          // “评论和 at 我的”
	ModifyMessageStatus             int    `json:"modifyMessageStatus"`             // “任务更新动态”
	RelationMessageStatus           int    `json:"relationMessageStatus"`           // 关联内容动态
	CollaborateMessageStatus        int    `json:"collaborateMessageStatus"`        // 我协作的  1表示开启我协作的任务相关推送 // New for 20220531
	DailyProjectReportMessageStatus int    `json:"dailyProjectReportMessageStatus"` // 项目日报
	DefaultProjectId                int64  `json:"defaultProjectId"`
	DefaultProjectObjectTypeId      int64  `db:"default_project_object_type_id,omitempty" json:"defaultProjectObjectTypeId"`
	PcNoticeOpenStatus              int    `db:"pc_notice_open_status,omitempty" json:"pcNoticeOpenStatus"`
	PcIssueRemindMessageStatus      int    `db:"pc_issue_remind_message_status,omitempty" json:"pcIssueRemindMessageStatus"`
	PcOrgMessageStatus              int    `db:"pc_org_message_status,omitempty" json:"pcOrgMessageStatus"`
	PcProjectMessageStatus          int    `db:"pc_project_message_status,omitempty" json:"pcProjectMessageStatus"`
	PcCommentAtMessageStatus        int    `db:"pc_comment_at_message_status,omitempty" json:"pcCommentAtMessageStatus"`
	UserViewLocationConfig          string `db:"user_view_location_config,omitempty" json:"userViewLocationConfig"`
	Ext                             string `db:"ext,omitempty" json:"ext"`
	RemindExpiring                  int    `json:"remindExpiring"` // 即将逾期提醒
}

type UserConfigExt struct {
	// 用户是否浏览过“新手指引”
	VisitedNewUserGuide bool `json:"visitedNewUserGuide"`
	// 版本信息是否展示
	Version *VersionConfig `json:"version"`

	// 2022-11-11活动信息标识
	Activity20221111 int `json:"activity20221111"`

	// 2023-02-07 个人中心 不再提醒弹窗
	RemindPopUp int `json:"remindPopUp"`
}

type VersionConfig struct {
	VersionInfoVisible bool `json:"versionInfoVisible"`
	//Version string `json:"version"`
}

type UserViewLocation struct {
	AppId       string `json:"appId"`
	ProjectId   int64  `json:"projectId"`
	IterationId int64  `json:"iterationId"`
	TableId     string `json:"tableId"`
	ViewId      string `json:"viewId"`
	MenuId      string `json:"menuId"`
	DashboardId string `json:"dashboardId"`
}

type UserViewLocationConfig struct {
	Location UserViewLocation `json:"location"`
}
