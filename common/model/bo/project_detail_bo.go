package bo

import (
	"time"

	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
)

type ProjectDetailBo struct {
	Id                       int64     `db:"id,omitempty" json:"id"`
	OrgId                    int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId                int64     `db:"project_id,omitempty" json:"projectId"`
	Notice                   string    `db:"notice,omitempty" json:"notice"`
	IsEnableWorkHours        int       `db:"is_enable_work_hours,omitempty" json:"isEnableWorkHours"`
	IsSyncOutCalendar        int       `db:"is_sync_out_calendar,omitempty" json:"isSyncOutCalendar"`
	SyncOutCalendarStatusArr []*int    `db:"-" json:"syncOutCalendarStatusArr"` // 日历同步的配置
	Creator                  int64     `db:"creator,omitempty" json:"creator"`
	CreateTime               time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator                  int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime               time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version                  int       `db:"version,omitempty" json:"version"`
	IsDelete                 int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*ProjectDetailBo) TableName() string {
	return "ppm_pro_project_detail"
}

// 针对群聊的推送配置数据。一个群聊可以绑定多个项目
type GetFsProjectChatPushSettingsForChatBo struct {
	Tables []GetFsProjectChatPushSettingsForChatBoTable `json:"tables"`
	// 添加任务时是否推送(1开，2关)
	CreateIssue int `json:"createIssue"`
	// 任务有新的评论时是否推送
	CreateIssueComment int `json:"createIssueComment"`
	// 任务被更新时是否推送
	UpdateIssueCase int `json:"updateIssueCase"`
	// 任务被更新时，那些字段更新时才推送
	ModifyColumnsOfSend []string `json:"modifyColumnsOfSend"`
}

type GetFsProjectChatPushSettingsForChatBoTable struct {
	ProjectId int64  `json:"projectId"`
	TableId   string `json:"tableId"`
}

// 针对项目的群聊配置数据
type GetFsProjectChatPushSettingsBo struct {
	ProjectId           int64                              `json:"projectId"`
	IsProjectChatNative int                                `json:"isProjectChatNative"`
	Tables              []GetFsTableSettingOfProjectChatBo `json:"tables"`
}

// 飞书群聊的配置项
type FsGroupChatSettingItems struct {
	// 添加任务时是否推送(1开2关)   默认值：1
	CreateIssue int `json:"createIssue"`
	// 任务有新的评论时是否推送
	CreateIssueComment int `json:"createIssueComment"`
	// 任务被更新时是否推送
	UpdateIssueCase int `json:"updateIssueCase"`
	// 任务被更新时，那些字段更新时才推送
	ModifyColumnsOfSend []string `json:"modifyColumnsOfSend"`
}

// 项目下表的群聊配置
type GetFsTableSettingOfProjectChatBo struct {
	// project chat 的 id 值
	SettingId  int64  `json:"settingId"`  // 该值实例化后，需要手动赋值才会有值
	ChatId     string `json:"chatId"`     // 该值实例化后，需要手动赋值才会有值
	TableIdStr string `json:"tableIdStr"` // 该值实例化后，需要手动赋值才会有值
	FsGroupChatSettingItems
}

/// 为 GetFsTableSettingOfProjectChatBo 挂载一些方法，便于调用

// CheckUpdate 检查任务的列更新后，是否要推送卡片消息到群聊中
func (setting *GetFsTableSettingOfProjectChatBo) CheckUpdate(columnName string) bool {
	if setting.UpdateIssueCase != 1 {
		return false
	}
	if isOk, _ := slice.Contain(setting.ModifyColumnsOfSend, consts.GroupChatAllIssueColumnFlag); isOk {
		return isOk
	}
	isOk, _ := slice.Contain(setting.ModifyColumnsOfSend, columnName)

	return isOk
}

type GetFsProjectChatPushSettingsOldBo struct {
	// 添加任务(1开2关)   默认值：1
	CreateIssue int `json:"createIssue"`
	// 任务负责人变更
	UpdateIssueOwner int `json:"updateIssueOwner"`
	// 任务状态变更
	UpdateIssueStatus int `json:"updateIssueStatus"`
	// 任务栏变更
	UpdateIssueProjectObjectType int `json:"updateIssueProjectObjectType"`
	// 任务标题被修改
	UpdateIssueTitle int `json:"updateIssueTitle"`
	// 任务时间被修改
	UpdateIssueTime int `json:"updateIssueTime"`
	// 任务有新的评论
	CreateIssueComment int `json:"createIssueComment"`
	// 任务有新的附件
	UploadNewAttachment int `json:"uploadNewAttachment"`
	// 创建群聊的开关标识
	EnableGroupChat int `json:"enableGroupChat"`
	// 项目隐私模式状态值
	ProPrivacyStatus int `json:"proPrivacyStatus"`
}

type OutGroupChatSettingBo struct {
	ChatId string                          `json:"chatId"`
	Tables []OutGroupChatSettingOneTableBo `json:"tables"`
}

type OutGroupChatSettingOneTableBo struct {
	// 表的 id
	TableIDStr string `json:"tableIdStr"`
	// 创建任务时，是否推送。1：推送
	CreateIssueCase int `json:"createIssueCase"`
	// 任务被评论时是否推送。1：推送
	CommentIssueCase int `json:"commentIssueCase"`
	// 任务的某些字段被修改时，是否推送。1：推送（字段请看 modifyColumnsOfSend）
	UpdateIssueCase int `json:"updateIssueCase"`
	// 如果开启“任务修改推送”时，这些字段被修改时，才推送消息
	ModifyColumnsOfSend []*string `json:"modifyColumnsOfSend"`
}

type ProjectOutChatSettingBo struct {
	// 项目 id
	ProjectID int64 `json:"projectId"`
	// 是否是项目创建的群聊。1是；2机器人被拉入普通群的群聊
	IsProjectChatNative int `json:"isProjectChatNative"`
	// 表的配置，一个项目下有多个表
	Tables []ProjectOutChatSettingOneTableBo `json:"tables"`
}

type ProjectOutChatSettingOneTableBo struct {
	// 表 id 整形
	TableID int64 `json:"tableId"`
	// 表的 id
	TableIDStr string `json:"tableIdStr"`
	// 创建任务时，是否推送。1：推送
	CreateIssueCase int `json:"createIssueCase"`
	// 任务被评论时是否推送。1：推送
	CommentIssueCase int `json:"commentIssueCase"`
	// 任务的某些字段被修改时，是否推送。1：推送（字段请看 modifyColumnsOfSend）
	UpdateIssueCase int `json:"updateIssueCase"`
	// 如果开启“任务修改推送”时，这些字段被修改时，才推送消息
	ModifyColumnsOfSend []*string `json:"modifyColumnsOfSend"`
}

// ToGroupChatConfig 将以项目为主体的配置转为以群聊为主体的配置对象 todo
func (chatConf *GetFsProjectChatPushSettingsBo) ToGroupChatConfig() {

}
