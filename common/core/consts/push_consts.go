package consts

//推送常量定义
const (
	PushIssueRemind            = "PushIssueRemind"
	PushIssueCommentAndAt      = "PushIssueCommentAndAt"
	PushIssueUpdate            = "PushIssueUpdate"
	PushRelatedContentDynamics = "PushRelatedContentDynamics"
)

type IssueNoticePushType int32

// 任务动态推送类型
const (
	// 创建任务
	PushTypeCreateIssue IssueNoticePushType = 0
	// 更新任务
	PushTypeUpdateIssue IssueNoticePushType = 1
	// 更新任务状态
	PushTypeUpdateIssueStatus IssueNoticePushType = 2
	// 删除任务
	PushTypeDeleteIssue IssueNoticePushType = 3
	// 更新任务成员
	//PushTypeUpdateIssueMembers IssueNoticePushType = 4
	// 更新任务关联
	//PushTypeUpdateRelationIssue IssueNoticePushType = 5
	// 每日项目日报推送- 通知项目推送
	PushTypeDailyProjectReportNoticProject IssueNoticePushType = 6
	// 每日项目日报推送- 通知消息服务
	PushTypeDailyProjectReportNoticMsg IssueNoticePushType = 7
	// 任务提醒
	PushTypeIssueRemind IssueNoticePushType = 8
	// 任务评论
	PushTypeIssueComment IssueNoticePushType = 9
	// 上传附件
	PushTypeUploadResource IssueNoticePushType = 10
	// 移动任务栏
	//PushTypeUpdateIssueProjectObjectType IssueNoticePushType = 11
	// 删除附件
	PushTypeDeleteResource IssueNoticePushType = 12
	// 任务描述通知
	PushTypeIssueRemarkRemind = 13
	// 任务日报推送
	PushTypeDailyIssueReport = 14
	// 迭代燃尽图
	PushTypeIterationBurnDownChart = 15
	// 项目燃尽图
	PushTypeProjectBurnDownChart = 16
	// 飞书新版提示进入方式
	PushTypeFeishuEntrance = 17
	// 新增任务标签
	PushTypeAddIssueTag = 18
	// 删除任务标签
	PushTypeDeleteIssueTag = 19
	// 恢复任务
	PushTypeRecoverIssue = 20
	// 第三方订单推送
	PushTypeThirdOrder = 21
	// 归档任务
	PushTypeArchiveIssue IssueNoticePushType = 21
	// 任务取消归档
	PushTypeCancelArchiveIssue IssueNoticePushType = 22
	// 引用附件
	PushTypeReferResource IssueNoticePushType = 22
	// 删除引用附件
	PushTypeDeleteReferResource IssueNoticePushType = 23
	// 新增工时
	PushTypeCreateWorkHour IssueNoticePushType = 24
	// 编辑工时
	PushTypeUpdateWorkHour IssueNoticePushType = 25
	// 删除工时
	PushTypeDeleteWorkHour IssueNoticePushType = 26
	// 更新前置任务
	PushTypeUpdateBeforeIssue IssueNoticePushType = 27
	// 更新后置任务
	PushTypeUpdateAfterIssue IssueNoticePushType = 28
	// 撤回任务
	PushTypeWithdrawIssue IssueNoticePushType = 29
	// 审核任务
	PushTypeAuditIssue IssueNoticePushType = 30
	// 任务发起群聊讨论
	PushTypeIssueStartChat IssueNoticePushType = 31
	// 移动项目表格
	PushTypeUpdateIssueProjectTable IssueNoticePushType = 32
	// 变更父任务
	PushTypeUpdateIssueParent IssueNoticePushType = 33

	// 接收飞书回调，处理失败的信息记录
	PushTypeFsCallbackHandleFailed IssueNoticePushType = 100
)

//项目动态推送类型
const (
	PushTypeCreateProject            IssueNoticePushType = 1000
	PushTypeUpdateProject            IssueNoticePushType = 1001
	PushTypeUpdateProjectMembers     IssueNoticePushType = 1002
	PushTypeStarProject              IssueNoticePushType = 1003
	PushTypeUnstarProject            IssueNoticePushType = 1004
	PushTypeUnbindProject            IssueNoticePushType = 1005
	PushTypeUpdateProjectStatus      IssueNoticePushType = 1006
	PushTypeCreateIssueBatch         IssueNoticePushType = 1007
	PushTypeCreateProjectFile        IssueNoticePushType = 1008
	PushTypeUpdateProjectFile        IssueNoticePushType = 1009
	PushTypeDeleteProjectFile        IssueNoticePushType = 1010
	PushTypeCreateProjectFolder      IssueNoticePushType = 1011
	PushTypeUpdateProjectFolder      IssueNoticePushType = 1012
	PushTypeDeleteProjectFolder      IssueNoticePushType = 1013
	PushTypeDeleteProject            IssueNoticePushType = 1014
	PushTypeUpdateProjectFileFolder  IssueNoticePushType = 1015
	PushTypeRecoverTag               IssueNoticePushType = 1016
	PushTypeRecoverFolder            IssueNoticePushType = 1017
	PushTypeRecoverProjectFile       IssueNoticePushType = 1018
	PushTypeRecoverProjectAttachment IssueNoticePushType = 1019
	PushTypeDeleteProjectTag         IssueNoticePushType = 1020
	PushTypeAddCustomField           IssueNoticePushType = 1021
	PushTypeDeleteCustomField        IssueNoticePushType = 1022
	PushTypeUseCustomField           IssueNoticePushType = 1023
	PushTypeForbidCustomField        IssueNoticePushType = 1024
	PushTypeUpdateCustomField        IssueNoticePushType = 1025
	PushTypeUseOrgCustomField        IssueNoticePushType = 1026
	PushTypeEnableWorkHour           IssueNoticePushType = 1027
	PushTypeDisableWorkHour          IssueNoticePushType = 1028

	PushTypeUpdateFormField IssueNoticePushType = 1029
)

//组织动态推送类型
const (
	//申请加入组织
	PushTypeApplyJoinOrg IssueNoticePushType = 2001
	//通过申请
	PushTypeApplicationApproved IssueNoticePushType = 2002
	//提升为管理员
	PushTypePromotionToOrgManager IssueNoticePushType = 2003
)

const (
	Code                  = "编号"
	Title                 = "标题"
	Owner                 = "负责人"
	Follower              = "关注人"
	Auditor               = "确认人"
	Status                = "状态"
	Relating              = "关联"
	BaRelating            = "前后置"
	PlanStartTime         = "计划开始时间"
	PlanEndTime           = "截止时间"
	PlanWorkHour          = "计划工作时间"
	Priority              = "优先级"
	Source                = "来源"
	IssueObjectTypeId     = "类型"
	Remark                = "备注"
	PublicStatus          = "项目公开性"
	ProjectNotice         = "项目公告"
	ProjectResourcePath   = "项目封面"
	ProjectResourceName   = "项目文件名"
	ProjectResourceFolder = "项目文件父级文件夹名"
	ProjectFolderName     = "项目文件夹名"
	ProjectFolderParentId = "项目文件夹父文件夹"
	ProjectObjectType     = "任务栏"
	Project               = "项目"
	Table                 = "表格"
	Parent                = "父记录"
	ProjectPreCode        = "前缀编号"
	IssuePropertyId       = "严重程度"
	Iteration             = "迭代"
	Name                  = "名称"
	OrgField              = "组织字段"
	Value                 = "值"
	Description           = "描述"
	WorkId                = "执行者"
	StartTime             = "开始时间"
	WorkTime              = "工作时间"
	WorkContent           = "工作内容"
)
