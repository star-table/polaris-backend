package consts

const (
	TrendsIssueExt = "{\"issueType\":\"T\"}"
)

const (
	TrendsModuleOrg     = "Org"
	TrendsModuleProject = "Project"
	TrendsModuleIssue   = "Issue"
	TrendsModuleRole    = "Role"
)

const (
	TrendsOperObjectTypeIssue    = "Issue"
	TrendsOperObjectTypeComment  = "Comment"
	TrendsOperObjectTypeWorkHour = "WorkHour"
	TrendsOperObjectTypeProject  = "Project"
	TrendsOperObjectTypeRole     = "Role"
	TrendsOperObjectTypeUser     = "User"
	TrendsOperObjectTypeOrg      = "Org"
	TrendsOperObjectTypeField    = "Field"
)

// 动态事件，放于relationType字段，长度限制32（任务）
const (
	//新增任务
	TrendsRelationTypeCreateIssue = "AddIssue"
	//添加评论
	TrendsRelationTypeCreateIssueComment = "AddIssueComment"
	//新增任务（子）
	TrendsRelationTypeCreateChildIssue = "AddChildIssue"
	//更新任务
	TrendsRelationTypeUpdateIssue = "UpdIssue"
	//更新任务状态
	TrendsRelationTypeUpdateIssueStatus = "UpdIssueStatus"
	//删除任务
	TrendsRelationTypeDeleteIssue = "DelIssue"
	//删除任务（子）
	TrendsRelationTypeDeleteChildIssue = "DelChildIssue"
	//添加关注人
	TrendsRelationTypeAddedIssueFollower = "AddIssueFollower"
	//去除关注人
	TrendsRelationTypeDeleteIssueFollower = "DelIssueFollower"
	//添加成员
	TrendsRelationTypeAddedIssueParticipant = "AddIssueParticipant"
	//修改任务负责人
	TrendsRelationTypeUpdateIssueOwner = "UpdateIssueOwner"
	// 添加负责人
	TrendsRelationTypeAddedIssueOwner = "AddIssueOwner"
	// 去除负责人
	TrendsRelationTypeDeleteIssueOwner = "DelIssueOwner"
	//删除成员
	TrendsRelationTypeDeletedIssueParticipant = "DelIssueParticipant"
	//添加关联任务
	TrendsRelationTypeAddRelationIssue = "AddRelationIssue"
	//添加单向关联任务
	TrendsRelationTypeAddSingleRelationIssue = "AddSingleRelationIssue"
	//被添加关联任务
	TrendsRelationTypeAddRelationIssueByOther = "AddRelationIssueByOther"
	//删除关联任务
	TrendsRelationTypeDeleteRelationIssue = "DelRelationIssue"
	//删除单向关联任务
	TrendsRelationTypeDeleteSingleRelationIssue = "DelSingleRelationIssue"
	//被删除关联任务
	TrendsRelationTypeDeleteRelationIssueByOther = "DelRelationIssueByOther"
	//上传附件
	TrendsRelationTypeUploadResource = "UploadResource"
	//删除附件
	TrendsRelationTypeDeleteResource = "DeleteResource"
	// 上传图片
	TrendsRelationTypeUploadPicture = "UploadPicture"
	// 删除图片
	TrendsRelationTypeDeletePicture = "DeletePicture"
	//删除项目附件
	TrendsRelationTypeDeleteProjectResource = "DeleteProjectResource"
	//变更任务栏
	//TrendsRelationTypeUpdateIssueProjectObjectType = "UpdateIssueProjectObjectType"
	//增加任务标签
	TrendsRelationTypeAddIssueTag = "AddIssueTag"
	//删除任务标签
	TrendsRelationTypeDeleteIssueTag = "DelIssueTag"
	//恢复任务
	TrendsRelationTypeRecoverIssue = "RecoverIssue"
	//引用附件
	TrendsRelationTypeReferResource = "ReferResource"
	//删除引用附件
	TrendsRelationTypeDelReferResource = "DelReferResource"
	// 新增工时
	TrendsRelationTypeCreateWorkHour = "AddIssueWorkHour"
	// 编辑工时
	TrendsRelationTypeUpdateIssueWorkHour = "UpdateIssueWorkHour"
	// 删除工时
	TrendsRelationTypeDeleteIssueWorkHour = "DeleteIssueWorkHour"
	//添加前置任务
	TrendsRelationTypeAddBeforeIssue = "AddBeforeIssue"
	//被添加后置任务
	TrendsRelationTypeAddAfterIssueByOther = "AddAfterIssueByOther"
	//添加后置任务
	TrendsRelationTypeAddAfterIssue = "AddAfterIssue"
	//被添加前置任务
	TrendsRelationTypeAddBeforeIssueByOther = "AddBeforeIssueByOther"
	//删除前置任务
	TrendsRelationTypeDeleteBeforeIssue = "DeleteBeforeIssue"
	//被删除后置任务
	TrendsRelationTypeDeleteAfterIssueByOther = "DeleteAfterIssueByOther"
	//删除后置任务
	TrendsRelationTypeDeleteAfterIssue = "DeleteAfterIssue"
	//被删除前置任务
	TrendsRelationTypeDeleteBeforeIssueByOther = "DeleteBeforeIssueByOther"
	//审核任务
	TrendsRelationTypeAuditIssue = "AuditIssue"
	//撤回任务
	TrendsRelationTypeWithdrawIssue = "WithdrawIssue"
	//添加确认人
	TrendsRelationTypeAddedIssueAuditor = "AddIssueAuditor"
	//删除确认人
	TrendsRelationTypeDeleteIssueAuditor = "DelIssueAuditor"
	TrendsRelationTypeStartIssueChat     = "StartIssueChat" // 发起任务群聊
	//变更项目表
	TrendsRelationTypeUpdateIssueProjectTable = "UpdateIssueProjectTable"
	//变更父任务
	TrendsRelationTypeUpdateIssueParent = "UpdateIssueParent"
)

// 动态事件，放于relationType字段，长度限制32（项目）
const (
	//创建项目
	TrendsRelationTypeCreateProject = "AddProject"
	//更新项目
	TrendsRelationTypeUpdateProject = "UpdProject"
	//关注项目
	TrendsRelationTypeStarProject = "StarProject"
	//取关项目
	TrendsRelationTypeUnstarProject = "UnstarProject"
	//退出项目
	TrendsRelationTypeUnbindProject = "UnbindProject"
	//更新项目状态
	TrendsRelationTypeUpdateProjectStatus = "UpdProjectStatus"
	//修改项目负责人
	TrendsRelationTypeUpdateProjectOwner = "UpdateProjectOwner"
	//添加关注人
	TrendsRelationTypeAddedProjectOwner = "AddProjectOwner"
	//去除关注人
	TrendsRelationTypeDeleteProjectOwner = "DelProjectOwner"
	//添加关注人
	TrendsRelationTypeAddedProjectFollower = "AddProjectFollower"
	//去除关注人
	TrendsRelationTypeDeleteProjectFollower = "DelProjectFollower"
	//添加成员
	TrendsRelationTypeAddedProjectParticipant = "AddProjectParticipant"
	//删除成员
	TrendsRelationTypeDeletedProjectParticipant = "DelProjectParticipant"
	//批量插入任务
	TrendsRelationTypeCreateIssueBatch = "CreateIssueBatch"
	//上传文件
	TrendsRelationTypeUploadProjectFile = "UploadProjectFile"
	//更新文件
	TrendsRelationTypeUpDateProjectFile = "UpdateProjectFile"
	//更新文件所属文件夹
	TrendsRelationTypeUpDateProjectFileFolder = "UpdateProjectFileFolder"
	//删除文件
	TrendsRelationTypeDeleteProjectFile = "DeleteProjectFile"
	//创建文件夹
	TrendsRelationTypeCreateProjectFolder = "CreateProjectFolder"
	//更新文件夹
	TrendsRelationTypeUpdateProjectFolder = "UpdateProjectFolder"
	//删除文件夹
	TrendsRelationTypeDeleteProjectFolder = "DeleteProjectFolder"
	//删除项目
	TrendsRelationTypeDeleteProject = "DeleteProject"
	//恢复标签
	TrendsRelationTypeRecoverTag = "RecoverTag"
	//删除项目标签
	TrendsRelationTypeDeleteProjectTag = "DeleteProjectTag"
	//恢复文件夹
	TrendsRelationTypeRecoverFolder = "RecoverFolder"
	//恢复文件
	TrendsRelationTypeRecoverProjectFile = "RecoverProjectFile"
	//恢复附件
	TrendsRelationTypeRecoverProjectAttachment = "RecoverProjectAttachment"
	//创建自定义字段
	TrendsRelationTypeAddCustomField = "AddCustomField"
	//删除自定义字段
	TrendsRelationTypeDeleteCustomField = "DeleteCustomField"
	//启用自定义字段
	TrendsRelationTypeUseCustomField = "UseCustomField"
	//禁用自定义字段
	TrendsRelationTypeForbidCustomField = "ForbidCustomField"
	//编辑自定义字段
	TrendsRelationTypeUpdateCustomField = "UpdateCustomField"
	//使用自定义字段
	TrendsRelationTypeUseOrgCustomField = "UseOrgCustomField"
	// 启用工时
	TrendsRelationTypeEnableWorkHour = "EnableWorkHour"
	// 禁用工时
	TrendsRelationTypeDisableWorkHour = "DisableWorkHour"
	// 操作字段变动
	TrendsRelationTypeUpdateFormField = "UpdateFormField"
)

// @相关事件（放于notice表中，但独立于其余类型的notice）
const (
	NoticeTypeIssueCommentAtSomebody = "IssueCommentAtSomebody"
	NoticeTypeIssueRemarkAtSomebody  = "IssueRemarkAtSomebody"
)

// 组织动态事件
const (
	TrendsRelationTypeApplyJoinOrg = "ApplyJoinOrg"
)

// 项目动态查询有效的动态类型
var ValidRelationTypesOfProject = []string{
	TrendsRelationTypeCreateProject,
	TrendsRelationTypeUpdateProject,
	TrendsRelationTypeStarProject,
	TrendsRelationTypeUnstarProject,
	TrendsRelationTypeUnbindProject,
	TrendsRelationTypeUpdateProjectStatus,
	TrendsRelationTypeAddedProjectFollower,
	TrendsRelationTypeDeleteProjectFollower,
	TrendsRelationTypeAddedProjectParticipant,
	TrendsRelationTypeDeletedProjectParticipant,
	TrendsRelationTypeCreateIssue,
	TrendsRelationTypeDeleteIssue,
	TrendsRelationTypeCreateIssueBatch,
	TrendsRelationTypeUpdateProjectOwner,
	TrendsRelationTypeAddedProjectOwner,
	TrendsRelationTypeDeleteProjectOwner,
	TrendsRelationTypeAddCustomField,
	TrendsRelationTypeDeleteCustomField,
	TrendsRelationTypeUseCustomField,
	TrendsRelationTypeForbidCustomField,
	TrendsRelationTypeUpdateCustomField,
	TrendsRelationTypeUseOrgCustomField,
	TrendsRelationTypeEnableWorkHour,
	TrendsRelationTypeDisableWorkHour,
}

var TrendsOperObjPropertyMap = map[string]string{}

// 动态修改字段名称定义
const (
	TrendsOperObjPropertyNameStatus      = "status"
	TrendsOperObjPropertyNameFollower    = "follower"
	TrendsOperObjPropertyNameParticipant = "participant"
	TrendsOperObjPropertyNameOwner       = "owner"
	TrendsOperObjPropertyNameAuditor     = "auditor"
)

const (
	TrendsRelationTypeUser = "User"
)

const (
	ShowTrendsEmptyTime = "未设置"
)

// 别名
func init() {
	TrendsOperObjPropertyMap[TrendsOperObjPropertyNameStatus] = "状态"
}
