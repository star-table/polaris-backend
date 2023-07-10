package consts

//Issue 关联类型
const (
	//负责人
	IssueRelationTypeOwner = 1
	//参与人
	IssueRelationTypeParticipant = 2
	//关注人
	IssueRelationTypeFollower = 3
	//关联任务
	IssueRelationTypeIssue = 4
	//资源
	IssueRelationTypeResource = 5
	//收藏项目
	IssueRelationTypeStar = 6
	//前后置任务
	IssueRelationTypeBeforeAfter = 7

	//状态
	IssueRelationTypeStatus = 20
	//飞书日历日程
	IssueRelationTypeCalendar = 21
	//飞书群聊外部id
	IssueRelationTypeChat = 22
	//飞书群聊外部id（主群聊，只用于记录主群聊，不会涉及删除和新增）
	IssueRelationTypeMainChat = 23
	//自定义字段
	IssueRelationTypeCustomField = 24
	//参与部门（目前针对于项目）
	IssueRelationTypeDepartmentParticipant = 25
	//确认人（审批）
	IssueRelationTypeAuditor = 26
	// IssueRelationTypeBelongManyPro 任务属于多项目的关联
	IssueRelationTypeBelongManyPro = 27

	// project chat 中的 chat_type 设置
	ChatTypeMain = 1 // 对应于主群聊，或者创建项目时，自动创建的群聊
	ChatTypeOut  = 2
)

var (
	// 建议仅用于查询**项目的成员**查询，不可用于任务的成员查询。
	MemberRelationTypeList = []int{IssueRelationTypeOwner, IssueRelationTypeParticipant}

	// 任务的参与人，包含：负责人、关注人、确认人
	IssueParticipantTypeList = []int{IssueRelationTypeOwner, IssueRelationTypeFollower, IssueRelationTypeAuditor}
)

//项目对象类型LangCode
const (
	// 敏捷项目的对象类型
	ProjectObjectTypeLangCodeIteration = "Project.ObjectType.Iteration"
	ProjectObjectTypeLangCodeBug       = "Project.ObjectType.Bug"
	ProjectObjectTypeLangCodeTestTask  = "Project.ObjectType.TestTask"
	ProjectObjectTypeLangCodeFeature   = "Project.ObjectType.Feature"
	ProjectObjectTypeLangCodeDemand    = "Project.ObjectType.Demand"
	ProjectObjectTypeLangCodeTask      = "Project.ObjectType.Task"

	// 通用项目的对象类型
	ProjectObjectTypeLangCodeCommonTask = "Project.ObjectType.Common.Task"
)

//项目类型LangCode
const (
	//普通
	ProjectTypeLangCodeNormalTask = "ProjectType.NormalTask"
	//敏捷
	ProjectTypeLangCodeAgile = "ProjectType.Agile"

	// 过滤内置字段、自定义字段信息的场景，不同场景会有一些特殊操作
	FilterColumnUsingTypeImport = "import"
	FilterColumnUsingTypeExport = "export"
)

//项目基础类型
const (
	//通用项目id
	ProjectTypeNormalId = 1
	//敏捷项目id
	ProjectTypeAgileId = 2
	// 新项目类型2022
	ProjectTypeCommon2022V47 = 47
	// 空项目
	ProjectTypeEmpty = 48
)

//流程langCode
const (
	ProcessLangCodeDefaultTask      = "Process.Issue.DefaultTask"      //默认任务流程
	ProcessLangCodeDefaultAgileTask = "Process.Issue.DefaultAgileTask" //默认敏捷项目任务流程
	ProcessLangCodeDefaultBug       = "Process.Issue.DefaultBug"       //默认缺陷流程
	ProcessLangCodeDefaultTestTask  = "Process.Issue.DefaultTestTask"  //默认测试任务流程
	ProcessLangCodeDefaultProject   = "Process.DefaultProject"         //默认项目流程
	ProcessLangCodeDefaultIteration = "Process.DefaultIteration"       //默认迭代流程
	ProcessLangCodeDefaultFeature   = "Process.Issue.DefaultFeature"   //默认特性流程
	ProcessLangCodeDefaultDemand    = "Process.Issue.DefaultDemand"    //默认需求流程
)

//项目对象类型ObjectType
const (
	//ProjectObjectTypeIteration = 1
	ProjectObjectTypeTask = 2
)

const (
	//项目流程
	ProcessPrject = 1
	//迭代流程
	ProcessIteration = 2
	//问题流程
	ProcessIssue = 3
)

const (
	//优先级类型-项目
	PriorityTypeProject = 1
	//优先级类型-需求/任务等优先级
	PriorityTypeIssue = 2
)

const (
	PriorityIssueCommon = "Priority.Issue.P2"

	// 项目隐私模式状态值 启用1；不启用2
	ProSetPrivacyEnable  = 1
	ProSetPrivacyDisable = 2
)

const (
	//项目状态
	ProcessStatusCategoryProject = 1
	//迭代状态
	ProcessStatusCategoryIteration = 2
	//任务状态
	ProcessStatusCategoryIssue = 3
)

//user team relation Type
const (
	UserTeamRelationTypeLeader = 1
	UserTeamRelationTypeMember = 2
)

const (
	DailyReport   = 1
	WeeklyReport  = 2
	MonthlyReport = 3
)

const TemplateDirPrefix = "resources/template/"

//资源存储方式
const (
	//本地
	LocalResource = 1
	//oss
	OssResource = 2
	//钉盘
	DingDiskResource = 3
	//集群
	PrivateResource = 4
	//飞书云文档
	FsResource = 5
)

var ImgTypeMap = map[string]bool{
	"GIF":  true,
	"JPG":  true,
	"JPEG": true,
	"PNG":  true,
	"TIFF": true,
	"WEBP": true,
}

//url类型
const (
	//http路径
	UrlTypeHttpPath = 1
	//dist路径
	UrlTypeDistPath = 2
)

const (
	//公用项目
	PublicProject = 1
	//私有项目
	PrivateProject = 2
)

const (
	ProjectMemberEffective = 1
	ProjectMemberDisabled  = 2
)

const (
	ProjectIsFiling     = 1
	ProjectIsNotFiling  = 2
	ProjectIsFilingDesc = "已归档"
)

//消息类型定义
const (
	MessageTypeIssueTrends   = 20
	MessageTypeProjectTrends = 21
)

//消息状态
const (
	MessageStatusWait      = 1
	MessageStatusDoing     = 2
	MessageStatusCompleted = 3
	MessageStatusFail      = 4
)

//是否导出飞书日历
const (
	IsSyncOutCalendar    = 1
	IsNotSyncOutCalendar = 2
	// 该值通过二进制的位运算进行计算 128 64 32 16 8 4 2 1 以此类推
	// 4：负责人；8关注人；其他暂时预留。
	IsSyncOutCalendarForOwner    = 4
	IsSyncOutCalendarForFollower = 8
	// 16 同步到订阅日历
	IsSyncOutCalendarForSubCalendar      = 16
	IsSyncOutCalendarForOwnerAndFollower = 12
	// 该值用于位运算的参数之一。业务上不做使用。
	IsSyncOutCalendarForAllPosition = 255
)

//性别 1:男性 2:女性
const (
	Male   = 1
	Female = 2
)

//联系地址类型
const (
	ContactAddressTypeMobile = 1
	ContactAddressTypeEmail  = 2
)

//验证码验证类型1: 登录验证码，2：注册验证码，3：修改密码验证码，4：找回密码验证码，5：绑定，6，解绑
const (
	AuthCodeTypeLogin              = 1
	AuthCodeTypeRegister           = 2
	AuthCodeTypeResetPwd           = 3
	AuthCodeTypeRetrievePwd        = 4
	AuthCodeTypeBind               = 5
	AuthCodeTypeUnBind             = 6
	AuthCodeTypeInviteUserByPhones = 7 // 通过手机号批量邀请成员加入团队
)

//登录类型
const (
	//短信验证码登录
	LoginTypeSMSCode = 1
	//账号密码登录
	LoginTypePwd = 2
	//邮箱登录
	LoginTypeMail = 3
)

//注册类型
const (
	//短信注册
	RegisterTypeSMSCode = 1
	//账号密码注册
	RegisterTypePwd = 2
	//邮箱注册
	RegisterTypeMail = 3
)

const (
	WeiXinScanCodeType     = 1
	WeiXinPlatformCodeType = 2
)

const (
	DingAuthError = "invalidAuthInfo"
)

//mqtt channelType
const (
	//项目
	MQTTChannelTypeProject = 1
	//应用
	MQTTChannelTypeApp = 4
	//组织
	MQTTChannelTypeOrg = 2
	//个人
	MQTTChannelTypeUser = 3
)

//mqtt channel
const (
	MQTTChannelPrefix  = "MQTT/"
	MQTTChannelRoot    = MQTTChannelPrefix + "#/"
	MQTTChannelOrgRoot = MQTTChannelPrefix + "org/{{.OrgId}}/#/"
	MQTTChannelOrg     = MQTTChannelPrefix + "org/{{.OrgId}}/channel/"
	MQTTChannelProject = MQTTChannelPrefix + "org/{{.OrgId}}/project/{{.ProjectId}}/channel/"
	MQTTChannelApp     = MQTTChannelPrefix + "org/{{.OrgId}}/app/{{.AppId}}/channel/"
	MQTTChannelTable   = MQTTChannelPrefix + "org/{{.OrgId}}/app/{{.AppId}}/table/{{.TableId}}/channel/"
	MQTTChannelUser    = MQTTChannelPrefix + "org/{{.OrgId}}/user/{{.UserId}}/channel/"

	MQTTChannelKeyOrg     = "OrgId"
	MQTTChannelKeyProject = "ProjectId"
	MQTTChannelKeyTable   = "TableId"
	MQTTChannelKeyApp     = "AppId"
	MQTTChannelKeyUser    = "UserId"
)

//mqtt默认ttl(单位：秒)
const MQTTDefaultTTL = 0

//mqtt默认key ttl(单位：秒)
const MQTTDefaultKeyTTLOneDay = 86400

//mqtt默认权限
const MQTTDefaultPermissions = "r"

//mqtt默认全局权限
const MQTTDefaultRootPermissions = "rwlsp"

//mqtt 通知类型
const (
	// 提醒
	MQTTNoticeTypeRemind = 1
	// 数据刷新
	MQTTNoticeTypeDataRefresh = 2
)

//mqtt 数据刷新通知类型
const (
	//新增
	MQTTDataRefreshActionAdd = "ADD"
	//删除
	MQTTDataRefreshActionDel = "DEL"
	//更新
	MQTTDataRefreshActionModify = "MODIFY"
	//更新任务sort
	MQTTDataRefreshActionModifySort = "MODIFYSORT"
	//移动
	MQTTDataRefreshActionMove = "MOVE"
	//归档
	MQTTDataRefreshActionArchive = "Archive"
	//取消归档
	MQTTDataRefreshActionCancelArchive = "CancelArchive"
)

//mqtt 数据刷新对象类型
const (
	//项目
	MQTTDataRefreshTypePro = "PRO"
	//任务
	MQTTDataRefreshTypeIssue = "ISSUE"
	//标签
	MQTTDataRefreshTypeTag = "TAG"
	//人员
	MQTTDataRefreshTypeMember = "MEMBER"
)

//支付等级
const (
	PayLevelStandard      = 1 // 标准版
	PayLevelEnterprise    = 2 // 企业版
	PayLevelFlagship      = 3 // 旗舰版
	PayLevelPrivateDeploy = 4 // 私有化部署

	PayLevelDouble11Activity = 33 // 双11大促 活动

	PayLevelDingFreeTrial = 33 // 免费试用 （在试用时间期限内 相当于旗舰版）
)

var (
	PayLevelDescMap = map[int]string{
		PayLevelStandard:      "标准版",
		PayLevelEnterprise:    "企业版",
		PayLevelFlagship:      "旗舰版",
		PayLevelPrivateDeploy: "私有化版",

		PayLevelDouble11Activity: "双11旗舰版6.66折",
	}
)

const (
	ActivityFlag = "activity20221111"

	ActivityFinished    = 2
	ActivityNotFinished = 1
)

const (
	IsPlatformAdmin  = 1
	NotPlatformAdmin = 2
)

//上传文件大小限制
const (
	// 单文件大小限制，默认的 1G
	MaxFileSizeStandard = 1073741824

	// 任务评论上传附件大小限制
	MaxIssueAttachmentSize = 10 * 1024 * 1024
)

//回收站对象类型
const (
	RecycleTypeIssue = 1 //任务
	//RecycleTypeTag        = 2 //标签
	RecycleTypeFolder     = 3 //文件夹
	RecycleTypeResource   = 4 //文件
	RecycleTypeAttachment = 5 //附件
)

//回收站版本id
const RecycleVersion = "recyle_version"

const (
	//任务层级限制
	IssueLevel = 9

	// MaxRawSupportForImportIssue 导入任务的最大任务数
	MaxRawSupportForImportIssue = 5000
)

const (
	//任务附件引用和上传区分（用于动态）
	IssueAttachmentReferRelationCode = "refer"
	// 任务属于多个项目的关联关系 code
	IssueRelationBelongManyProCode = "IssueRelationBelongManyPro"
)

var (
	// 任务、项目优先级名称多语言版本
	LANG_PRIORITIES_MAP = map[string]map[string]string{
		"en": {
			"Priority.Project.P0": "important",
			"Priority.Project.P1": "normal",
			"Priority.Project.P2": "lowest",
			"Priority.Issue.P0":   "highest",
			"Priority.Issue.P1":   "higher",
			"Priority.Issue.P2":   "normal",
			"Priority.Issue.P3":   "lower",
			"Priority.Issue.P4":   "lowest",
		},
	}
	LANG_PRIORITIES_NAME_MAP = map[string]map[string]string{
		"en": {
			"最高": "highest",
			"较高": "higher",
			"普通": "normal",
			"较低": "lower",
			"最低": "lowest",
		},
	}
	LANG_STATUS_MAP = map[string]map[string]string{
		"en": {
			"未开始": "Pending",
			"进行中": "In Progress",

			"待处理":  "Todo",
			"设计中":  "Designing",
			"研发中":  "Developing",
			"测试中":  "Testing",
			"已发布":  "Published",
			"已关闭":  "Closed",
			"已确认":  "Confirmed",
			"修复中":  "Repairing",
			"复测通过": "Passed",
			"已解决":  "Resolved",
			"未完成":  "Incomplete",
			"已完成":  "Completed",
			"待确认":  "Unconfirmed",
		},
	}
	LANG_ISSUE_SRC_MAP = map[string]map[string]string{
		"en": {
			"市场需求": "Market",
			"活动需求": "Activity",
			"客户":   "Customer",
			"老板":   "Boss",
			"产品经理": "PM",
			"其他":   "Other",
		},
	}
	LANG_ISSUE_TYPE_MAP = map[string]map[string]string{
		"en": {
			// 需求
			"新功能":  "New Func",
			"需求优化": "Demand Improve",
			"交互优化": "Interaction Improve",
			"视觉优化": "Visual Improve",
			"其他":   "Other",
			// 缺陷
			"功能问题":  "Function problem",
			"性能问题":  "Performance problem",
			"接口问题":  "The interface problem",
			"界面问题":  "The vision problem",
			"兼容问题":  "Compatibility problem",
			"易用性问题": "Usability issues",
			"逻辑问题":  "Logical problem",
			"需求问题":  "Demand",
			"不确定":   "Other",
		},
	}
	LANG_ISSUE_PROPERTY_MAP = map[string]map[string]string{
		"en": {
			"建议": "Advice",
			"一般": "General",
			"严重": "Severity",
			"致命": "Deadly",
		},
	}
	LANG_PROJECT_STATUS_MAP = map[string]map[string]string{
		"en": {
			"进行中": "In Progress",
			"已完成": "Completed",
		},
	}
	LANG_ISSUE_STAT_DESC_MAP = map[string]map[string]string{
		"en": {
			"待处理": "Todo",
			"未完成": "Incomplete",
			"进行中": "In Progress",
			"已完成": "Completed",
			"已逾期": "Overdue",
			// 另一部分的中英文转换
			"设计中": "Designing",
			"研发中": "Developing",
			"测试中": "Testing",
			"已发布": "Published",
			"已关闭": "Closed",
			// 对一些内置状态的中英文的补充
			"未逾期":  "Not Overdue",
			"未开始":  "Pending",
			"已确认":  "Confirmed",
			"修复中":  "Repairing",
			"复测通过": "Passed",
			"已解决":  "Resolved",
			"待确认":  "Unconfirmed",
		},
	}
	LANG_ROLE_NAME_MAP = map[string]map[string]string{
		"en": {
			"组织超级管理员": "Super Admin",
			"组织管理员":   "Admin",
			"组织成员":    "Member",
		},
	}
	LANG_PRO_ROLE_NAME_MAP = map[string]map[string]string{
		"en": {
			"负责人":  "Admin",
			"项目成员": "Project Member",
		},
	}
	LANG_PRO_ITER_NAME_MAP = map[string]map[string]string{
		"en": {
			"未规划": "Not planning",
		},
	}
	LANG_ISSUE_IMPORT_TITLE_MAP = map[string]map[string]string{
		"en": {
			"任务栏":    "Section",
			"任务栏名称":  "Section",
			"任务标题":   "Task Title",
			"标题":     "Task Title",
			"*标题":    "*Task Title",
			"*任务栏名称": "*Section",
			"*任务标题":  "*Task Title",
			"负责人":    "Assignee",
			"父子类型":   "Superior or sub",
			"优先级":    "Priority",
			"任务状态":   "Task Status",
			"任务描述":   "Task Description",
			"开始时间":   "Start Date",
			"截止时间":   "Due Date",
			"标签":     "Tags",
			"关注人":    "Collaborators",
			"确认人":    "Auditor",
			"任务创建时间": "Creation Date",
			"任务创建者":  "Creator",
			"编号ID":   "ID Number",
			"编号":     "ID Number",
			"实际完成时间": "Completed Date",
			"迭代":     "Iteration",

			"任务类型":  "Task Type",
			"*任务类型": "*Task Type",
			"需求类型":  "Demand Type",
			"需求来源":  "Demand Source",
			"缺陷类型":  "Bug Type",
			"严重程度":  "Severity Degree",

			"需求/缺陷类型":   "Demand/Bug Type",
			"需求来源/严重程度": "Demand Source/Severity Degree",

			"重名部门/成员": "Duplicate departments/members",
			"类型":      "Type",
			"所处部门":    "Department of affiliation",
			"部门/成员ID": "Departments/members ID",
			"建议导入样式":  "Recommended import styles",
			"用户":      "User",
			"部门":      "Department",
		},
	}
	LANG_ISSUE_OTHER_MAP = map[string]map[string]string{
		"en": {
			"子任务": "Sub Task",
			"父任务": "Superior Task",
		},
	}
	// 敏捷项目，任务类型（需求、任务、缺陷）的映射
	LANG_ISSUE_OBJECT_TYPE_MAP = map[string]map[string]string{
		"en": {
			"需求": "Demand",
			"任务": "Task",
			"缺陷": "Bug",
		},
	}
	// 用户自定义字段，预定义的系统字段多语言映射
	LANG_CUSTOMER_SYS_FIELD_MAP = map[string]map[string]string{
		"en": {
			"进度": "Progress",
			"评分": "Score",
		},
	}
	LANG_IMPORT_ERR_MAP = map[string]map[string]string{
		"en": {
			"任务统计":  "task_statistics",
			"统计":    "statistics",
			"全部工作项": "all_work_items",
			"父任务":   "Superior Task",
			"子任务":   "Sub Task",

			// 导入，异常提示的文案英文映射
			" 第%d行,第%d列错误：%s":                " Row %d, Col %d error: %s",
			ImportUserErrOfEndTimeLessErr:    "Plan end time must greater than plan start time! ",
			ImportUserErrOfStartTimeStyleErr: "Plan end/start time style error. ",
			ImportIssueErrOfEmail:            "Email style error",
			ImportIssueErrOfLink:             "Link style error",
			ImportIssueErrOfPhone:            "Phone number style error",
			ImportIssueErrOfIssueType:        "Task type error",
			ImportIssueErrOfDateTime:         "Date or time style error",
		},
	}
)

const (
	TriggerByNormal        = "normal"
	TriggerByCopy          = "copy"
	TriggerByImport        = "import"        // 导入
	TriggerByApplyTemplate = "applyTemplate" // 应用模板
	TriggerByAutomation    = "automation"    // 自动化(工作流)
	TriggerByAutoSchedule  = "autoSchedule"  // 自动排期
)

const (
	TabularViewName    = "表格视图"
	OwnerViewName      = "我负责的"
	AllOverDueSoonView = "所有即将逾期"
	AllOverDueView     = "所有已逾期"
	AllUnFinishedView  = "所有未完成"
	AllIssueView       = "所有任务"
	MyOverDueView      = "我的已逾期"
	MyDailyDueOfToday  = "我的今日截止"
	MyOverDueSoonView  = "我的即将逾期"
	MyPending          = "我的待完成"

	OverDueSoonViewName = "即将逾期"
)

// 表格视图type
const (
	TabularViewType = 1 // 表格视图
	KanbanViewType  = 2 // 看板视图
	GanttViewType   = 6 // 甘特视图
)

//无分配项目的任务编号前缀
const NoProjectPreCode = "NPC"

//未分配用户默认头像
const AvatarForUnallocated = "http://attachments.startable.cn/front_resources/Unallocated.png"

const (
	//不需要确认的初始状态
	AuditStatusNoNeed = -1
	//未查看（任务里面需要审批的默认为此状态）
	AuditStatusNotView = 1
	//已查看未审核
	AuditStatusView = 2
	//审批通过
	AuditStatusPass = 3
	//审批驳回
	AuditStatusReject = 4
)

const (
	WaitConfirmStatusName      = "待确认"
	WaitConfirmStatusEnName    = "To confirm" // 英文名
	CompleteStatusEnName       = "Completed"
	WaitConfirmStatusFontStyle = "#FF8040"
	WaitConfirmStatusBgStyle   = "#FCE7DC"

	// excel 导入成员错误信息
	ImportUserErrOfPhoneFormatErr       = "手机号码格式错误"
	ImportUserErrOfPhoneRegionFormatErr = "国家代码格式错误"
	ImportUserErrOfPhoneSendLimitErr    = "该手机号发送邀请过于频繁，请稍后再试。"
	ImportUserErrOfPhoneMultipleErr     = "手机号码重复（在提交的表格中重复）"
	ImportUserErrOfPhoneExistErr        = "当前团队内已存在相同的手机号码，请更换后重试"
	ImportUserErrOfEndTimeLessErr       = "截止时间必须大于开始时间！"
	ImportUserErrOfStartTimeStyleErr    = "开始/截止时间格式不正确！"

	// 导入任务时的一些提示信息
	ImportIssueErrOfEmail     = "邮箱格式不正确"
	ImportIssueErrOfLink      = "链接格式不正确"
	ImportIssueErrOfPhone     = "手机号格式不正确"
	ImportIssueErrOfIssueType = "任务类型不正确"
	ImportIssueErrOfDateTime  = "日期/时间格式不正确"
)

const (
	RemindPopUp = "remindPopUp"

	NeedRemindPopUp    = 1 // 个人中心需要弹窗弹出
	NotNeedRemindPopUp = 2 // 个人中心不需要弹窗弹出
)

// 飞书权限标识
const (
	// 飞书云文档的关键词就是 drive
	FsScopeDriveRead = "drive:drive:readonly"

	// 机器人个人指令名
	BotPrivateInsNameOfSetting = "设置"

	// 异步任务，已处理的任务 redis key
	AsyncTaskHashPartKeyOfProcessed = "processed" // 处理成功的数量
	AsyncTaskHashPartKeyOfFailed    = "failed"    // 处理失败的数量
	AsyncTaskHashPartKeyOfTotal     = "total"     // 记录总数量
	AsyncTaskHashPartKeyOfStartTime = "startTime" // 异步任务开始时间戳
	AsyncTaskHashPartKeyOfCover     = "cover"     // 处理应用模板时，返回给前端的模板封面图
	AsyncTaskHashPartKeyOfErrCode   = "errCode"   // 异步任务异常时，保存异常 code
	AsyncTaskHashPartKeyOfCardSend  = "cardSend"  // 异步任务完成后，是否已发送卡片通知
	AsyncTaskHashPartKeyOfTableIds  = "tableIds"  // 异步任务相关的 tableId。仅仅用于配合前端（lianBin）：需要一个 tableIds
	//AsyncTaskHashPartKeyOfParentTotalCount = "parentTotalCount" // 父任务总计数
	//AsyncTaskHashPartKeyOfParentSucCount   = "parentSucCount"   // 父任务成功的计数
	//AsyncTaskHashPartKeyOfParentErrCount   = "parentErrCount"   // 父任务失败的计数

	// 任务导入/应用模板等异步任务处理，设定一个进度步进时间阈值，超过该时间，视为失败
	AsyncTaskProcessTimeout = 60 * 5
)

// 和无码交互的一些常量/变量值
var (
	// 应用类型，1：表单，2：仪表盘，3：文件夹，4：项目, 5: 汇总表
	// 表单类型为 1
	LcAppTypeForForm   = 1
	LcAppTypeForFolder = 3
	// 创建极星的项目，对应到无码中是一个应用，该应用的类型是 4。
	LcAppTypeForPolaris = 4
	// 汇总表 5
	LcAppTypeForSummaryTable = 5
	// 视图镜像
	LcAppTypeForViewMirror = 6

	// 创建组织时，创建默认示例项目，任务栏单选选项的 id 值设定
	DefaultProjectTaskBarIds = map[int]int64{
		0: 1,
		1: 2,
		2: 3,
		3: 4,
		4: 5,
	}

	// 特殊处理：新项目类型，任务栏不会再出现
	// Common2022 指代"新项目类型"
	HiddenColumnNamesForCommon2022Type = []string{
		BasicFieldProjectObjectTypeId,
		BasicFieldIterationId,
		BasicFieldAuditorIds,
		BasicFieldProjectId,
		BasicFieldParentId,
	}

	HiddenColumnNameForCommon = []string{
		BasicFieldIterationId,
		BasicFieldProjectId,
		BasicFieldParentId,
	}

	HiddenColumnNameForAgileId = []string{
		BasicFieldProjectObjectTypeId,
		BasicFieldProjectId,
		BasicFieldAuditorIds,
		BasicFieldParentId,
	}

	HiddenColumnNameForEmpty = []string{
		BasicFieldCode,
		BasicFieldProjectId,           //所属项目
		BasicFieldIterationId,         //迭代
		BasicFieldParentId,            //父任务ID
		BasicFieldRemark,              //任务描述
		BasicFieldOwnerId,             //负责人
		BasicFieldProjectObjectTypeId, //任务栏 // 任务栏不再是基础字段，先注释
		BasicFieldIssueStatus,         //任务状态
		BasicFieldPlanStartTime,       //开始时间
		BasicFieldPlanEndTime,         //截止时间
		BasicFieldFollowerIds,         //关注人
		BasicFieldAuditorIds,          //确认人
	}

	IssueChatTitlePrefix      = "[任务讨论] "
	IssueChatTitleNoNameTitle = "未命名任务"
)

func IssueChatTitle(title string) string {
	if title == "" {
		title = IssueChatTitleNoNameTitle
	}
	return IssueChatTitlePrefix + title
}

var BasicFieldsMap = map[string]string{
	BasicFieldIterationId:   TcIterationId,
	BasicFieldPlanStartTime: TcPlanStartTime,
	BasicFieldPlanEndTime:   TcPlanEndTime,
	BasicFieldRemark:        TcRemark,
}

var DefaultWorkHour = map[string]interface{}{
	"planHour":        "0",
	"actualHour":      "0",
	"collaboratorIds": []string{},
}

const (
	IncludeAdmin = 1
)

var AllTranslate = map[string]map[string]string{}

func init() {
	list := []map[string]map[string]string{
		LANG_PRIORITIES_NAME_MAP, LANG_STATUS_MAP, LANG_ISSUE_SRC_MAP, LANG_ISSUE_TYPE_MAP, LANG_ROLE_NAME_MAP,
		LANG_PRO_ROLE_NAME_MAP, LANG_PRO_ITER_NAME_MAP, LANG_ISSUE_IMPORT_TITLE_MAP, LANG_ISSUE_OTHER_MAP,
		LANG_ISSUE_OBJECT_TYPE_MAP, LANG_CUSTOMER_SYS_FIELD_MAP, LANG_IMPORT_ERR_MAP,
	}
	for _, m := range list {
		for s, m2 := range m {
			if AllTranslate[s] == nil {
				AllTranslate[s] = make(map[string]string)
			}
			for s2, s3 := range m2 {
				AllTranslate[s][s2] = s3
			}
		}
	}
}

const (
	WeiXinStandardId   = "sp28f962a931196966"
	WeiXinEnterpriseId = "sp27595b32ebd809b2"
	WeiXinFlagshipId   = "spb2f970fbc6267e95"

	TestWeiXinStandardId   = "spfc96b11066ab0fe3"
	TestWeiXinEnterpriseId = "sp0daf1130d48ea70a"
	TestWeiXinFlagshipId   = "spdfd120ffb5c71e07"
)

const (
	FsStandardId   = "price_a15167c3baae1013"
	FsEnterpriseId = "price_a15167c3baaf5013"
	FsFlagshipId   = "price_a387472ac97b900b"
)

const (
	DingFreeCode       = "DT_GOODS_881666335729668_1808020"
	DingTrailCode      = "DT_GOODS_881666335729668_1808016"
	DingEnterpriseCode = "DT_GOODS_881666335729668_1808018"
	DingFlagshipCode   = "DT_GOODS_881666335729668_1808015"
)

const (
	SetFsPayRange    = 1 // 飞书设置了付费范围
	NotSetFsPayRange = 2 // 飞书没有设置付费范围
)

// 企微订单状态
const (
	// 企微订单状态
	WeiXinOrderStatusPaySuccess = 1
)

func GetDingOrderLevel(itemCode string) int64 {
	switch itemCode {
	case DingFreeCode:
		return PayLevelStandard
	case DingEnterpriseCode:
		return PayLevelEnterprise
	case DingFlagshipCode, DingTrailCode:
		return PayLevelFlagship
	default:
		return PayLevelStandard
	}
}

func GetFsOrderLevel(pricePlanId string) int64 {
	switch pricePlanId {
	case FsStandardId:
		return PayLevelStandard
	case FsEnterpriseId:
		return PayLevelEnterprise
	case FsFlagshipId:
		return PayLevelFlagship
	default:
		return PayLevelStandard
	}
}

func GetWeiXinOrderLevel(editionId string) int64 {
	switch editionId {
	case WeiXinStandardId:
		return PayLevelStandard
	case WeiXinEnterpriseId:
		return PayLevelEnterprise
	case WeiXinFlagshipId:
		return PayLevelFlagship
	default:
		return PayLevelStandard
	}
}
