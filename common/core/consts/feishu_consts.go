package consts

//卡片推送自定义字段
const (
	FsCardValueAppId                  = "appId"
	FsCardValueIssueId                = "issueId"
	FsCardValueTableId                = "tableId" // 任务所在的表id
	FsCardValueProjectId              = "projectId"
	FsCardValueOrgId                  = "orgId"
	FsCardValueUserId                 = "userId"
	FsCardValueCardType               = "cardType"
	FsCardValueStatusType             = "statusType"
	FsCardValueAction                 = "action"
	FsCardValueUrgeIssueBeUrgedUserId = "fsCardValueUrgeIssueBeUrgedUserId" // 任务催办：被催者id
	FsCardValueUrgeIssueUrgeText      = "fsCardValueUrgeIssueUrgeText"
	FsCardValueIsGroupChat            = "fsCardValueIsGroupChat"         // 是否是群聊
	FsCardValueGroupChatId            = "fsCardValueGroupChatId"         // 群聊的 out id
	FsCardValueAuditIssueResult       = "fsCardValueAuditIssueResult"    // 通用项目的任务审批，3审批通过;4驳回
	FsCardValueAuditIssueStarterId    = "fsCardValueAuditIssueStarterId" // 任务审批的发起人 id
	FsCardValueCardTitle              = "cardTitle"
	FsCardDatePickerInit              = "initialDatetime" // datePicker组件的初始值key
)

//卡片类型
const (
	FsCardTypeIssueRemind     = "IssueRemind"          // 即将到达截止时间
	FsCardTypeIssueUrge       = "FsCardTypeIssueUrge"  // 任务催办卡片
	FsCardTypeIssueAudit      = "FsCardTypeIssueAudit" // 任务审核卡片
	FsCardTypeNoticeSubscribe = "NoticeSubscribe"      //通知订阅
	FsCardTypeIssueCreate     = "FsCardIssueCreate"    // 通过机器人指令创建记录
	FsCardTypeDealTime        = "FsCardTypeDealTime"   // 飞书卡片交互的时候选择时间
)

//卡片动作
const (
	FsCardActionUpdatePlanEndTime   = "FsCardActionUpdatePlanEndTime"
	FsCardActionUpdatePlanStartTime = "FsCardActionUpdatePlanStartTime"
	FsCardActionUpdateIssueStatus   = "FsCardActionUpdateIssueStatus"
	FsCardUnsubscribePersonalReport = "FsCardActionUnsubscribePersonalReport" //退订个人日报
	FsCardUnsubscribeProjectReport  = "FsCardActionUnsubscribeProjectReport"  //退订项目日报
	FsCardUnsubscribeIssueRemind    = "FsCardActionUnsubscribeIssueRemind"    //退订任务提醒
	FsCardUnsubscribeMyCollaborate  = "FsCardActionUnsubscribeMyCollaborate"  //退订我负责的
	FsCardUnsubscribeMyOwn          = "FsCardActionUnsubscribeMyOwn"          //退订我协作的
)

//飞书订单类型
const (
	FsOrderBuy     = "buy"
	FsOrderUpgrade = "upgrade"
	FsOrderRenew   = "renew"
)

//飞书价格方案类型
const (
	FsPerSeatPerYear  = "per_seat_per_year"
	FsPerSeatPerMonth = "per_seat_per_month"
	FsTrial           = "trial"
	FsActiveDay       = "active_day"
	FsActiveEndDate   = "active_end_date"
	FsPermanent       = "permanent"

	// 飞书订单试用天数
	FsOrderTrailDays = 15
)

const FsUserStatusValid = "valid"
const FsUserStatusNotInScope = "not_in_scope"
const FsUserStatusNotActiveLicense = "no_active_license"

//飞书文档类型
const FsDocumentTypeDoc = "doc"
const FsDocumentTypeDocx = "docx"
const FsDocumentTypeSheet = "sheet"
const FsDocumentTypeSlide = "slide"
const FsDocumentTypeBiTable = "bitable"
const FsDocumentTypeMindNote = "mindnote"
const FsDocumentTypeFile = "file"

//飞书文档地址前缀
const FsDocumentDomain = "https://feishu.cn"

//飞书捷径trigger订阅
const FsShortcutSubscribe = 1

//飞书捷径trigger退订
const FsShortcutUnsubscribe = 2

//飞书捷径推送地址
//const FsShortcutPushUrl = "https://www.feishu.cn/flow/api/subscriber/ptllbvEXgP7QX9Lp"
const FsShortcutPushUrl = "https://www.feishu.cn/flow/api/subscriber/ptleyJMdX8DJXlD1"

//飞书捷径triggerId
const FsTriggerCreateIssue = "createIssue"
const FsTriggerUpdateIssue = "updateIssue"
const FsTriggerFinishIssue = "finishIssue"
const FsTriggerDoProjectObjectType = "doProjectObjectType"
const FsTriggerUpdateProject = "updateProject"

//飞书捷径事件类型
const FsEventUpdateProjectName = "updateProjectName"     //重命名项目
const FsEventUpdateProjectRemark = "updateProjectRemark" //修改项目描述
const FsEventUpdateProjectOwner = "updateProjectOwner"   //修改项目负责人
const FsEventAddProjectMember = "addProjectMember"       //添加项目成员
const FsEventDeleteProjectMember = "deleteProjectMember" //移除项目成员

const FsEventCreateProjectObjectType = "createProjectObjectType"         //新建任务栏
const FsEventUpdateProjectObjectTypeName = "updateProjectObjectTypeName" //重命名任务栏
const FsEventDeleteProjectObjectType = "deleteProjectObjectType"         //删除任务栏

const FsEventMoveIssue = "moveIssue"                               //移动任务
const FsEventUpdateIssueName = "updateIssueName"                   //重命名任务
const FsEventUpdateIssueStatus = "updateIssueStatus"               //更新任务状态
const FsEventUpdateIssueRemark = "updateIssueRemark"               //更新任务描述
const FsEventUpdateIssueOwner = "updateIssueOwner"                 //修改负责人
const FsEventUpdateIssueFollower = "updateIssueFollower"           //修改关注人
const FsEventUpdateIssuePlanStartTime = "updateIssuePlanStartTime" //更新开始时间
const FsEventUpdateIssuePlanEndTime = "updateIssuePlanEndTime"     //更新截止时间
const FsEventCreateIssueComment = "createIssueComment"             //任务新评论
const FsEventAddIssueAttachment = "addIssueAttachment"             //任务新附件
const FsEventDeleteIssue = "deleteIssue"                           //删除任务

var FsEventDefine = map[string]string{
	FsEventUpdateProjectName:           "重命名项目",
	FsEventUpdateProjectRemark:         "修改项目描述",
	FsEventUpdateProjectOwner:          "修改项目负责人",
	FsEventAddProjectMember:            "添加项目成员",
	FsEventDeleteProjectMember:         "移除项目成员",
	FsEventCreateProjectObjectType:     "新建任务栏",
	FsEventUpdateProjectObjectTypeName: "重命名任务栏",
	FsEventDeleteProjectObjectType:     "删除任务栏",
	FsEventMoveIssue:                   "移动任务",
	FsEventUpdateIssueName:             "重命名任务",
	FsEventUpdateIssueStatus:           "更新任务状态",
	FsEventUpdateIssueRemark:           "更新任务描述",
	FsEventUpdateIssueOwner:            "修改负责人",
	FsEventUpdateIssueFollower:         "修改关注人",
	FsEventUpdateIssuePlanStartTime:    "更新开始时间",
	FsEventUpdateIssuePlanEndTime:      "更新截止时间",
	FsEventCreateIssueComment:          "任务新评论",
	FsEventAddIssueAttachment:          "任务新附件",
	FsEventDeleteIssue:                 "删除任务",
}

const (
	ScopeNameContactAccessAsApp   = "contact:contact:access_as_app"   //以应用身份访问通讯录
	ScopeNameContactReadOnly      = "contact:contact:readonly"        //以应用身份访问通讯录（旧版）
	ScopeNameContactReadOnlyAsApp = "contact:contact:readonly_as_app" //以应用身份访问通讯录，2021-11-22
	ScopeNameCalendarCalendar     = "calendar:calendar"               //更新日历及日程信息v4
	ScopeNameCalendarReadOnly     = "calendar:calendar:readonly"      //获取日历、日程及忙闲信息v4
	ScopeNameCalendarAccess       = "calendar:calendar:access"        //日历相关权限v3
	ScopeNameDocOld               = "docs:docs:operate_as_user"       //老版云文档(操作)
	ScopeNameDocNew               = "drive:drive:readonly"            //新版云文档（查看）
	ScopeNameDocRead              = "drive:drive.search:readonly"     //云文档查看
)
