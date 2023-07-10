package vo

import (
	"github.com/star-table/common/core/types"
)

type CreateIssueReq struct {
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 标题
	Title string `json:"title"`
	// 优先级
	PriorityID int64 `json:"priorityId"`
	// 类型id，问题，需求....
	TypeID *int64 `json:"typeId"`
	// 负责人
	OwnerID []int64 `json:"ownerId"`
	// 参与人
	ParticipantIds []int64 `json:"participantIds"`
	// 关注人
	FollowerIds []int64 `json:"followerIds"`
	// 关注人部门（后端实际转化为人）
	FollowerDeptIds []int64 `json:"followerDeptIds"`
	// 计划开始时间
	PlanStartTime *types.Time `json:"planStartTime"`
	// 计划结束时间
	PlanEndTime *types.Time `json:"planEndTime"`
	// 计划工作时长
	PlanWorkHour *int `json:"planWorkHour"`
	// 所属版本id
	VersionID *int64 `json:"versionId"`
	// 所属模块id
	ModuleID *int64 `json:"moduleId"`
	// 父任务id
	ParentID *int64 `json:"parentId"`
	// 备注
	Remark *string `json:"remark"`
	// 备注详情
	RemarkDetail *string `json:"remarkDetail"`
	// 备注提及人
	MentionedUserIds []int64 `json:"mentionedUserIds"`
	// 所属迭代id
	IterationID *int64 `json:"iterationId"`
	// 问题对象类型id
	IssueObjectID *int64 `json:"issueObjectId"`
	// 来源id
	IssueSourceID *int64 `json:"issueSourceId"`
	// 性质id
	IssuePropertyID *int64 `json:"issuePropertyId"`
	// 状态id
	StatusID *int64 `json:"statusId"`
	// 子任务列表
	Children []*IssueChildren `json:"children"`
	// 关联的标签列表
	Tags []*IssueTagReqInfo `json:"tags"`
	// 关联的附件id列表
	ResourceIds []int64 `json:"resourceIds"`
	// 自定义字段
	CustomField []*UpdateIssueCustionFieldData `json:"customField"`
	// 审批人
	AuditorIds []int64 `json:"auditorIds"`
	// 无码入参
	LessCreateIssueReq map[string]interface{} `json:"lessCreateIssueReq"`
	// 前面的任务id
	BeforeID *int64 `json:"beforeId"`
	// 后面的任务id
	AfterID *int64 `json:"afterId"`
	// 前面的无码数据id
	BeforeDataID *string `json:"beforeDataId"`
	// 后面的无码数据id
	AfterDataID *string `json:"afterDataId"`
	// 排序
	Asc *bool `json:"asc"`
	// 是否可忽略字段准确性（目前主要用于模板生成任务）
	IsImport *bool `json:"isImport"`
	// 表id
	TableID string `json:"tableId"`
	// 处理业务时的辅助参数，比如导入任务时，向其中存储 import 表示导入来源，便于后续业务处理
	//ExtraInfo map[string]interface{} `json:"extraInfo"`
}

// 更新任务请求结构体
type UpdateIssueReq struct {
	// 要更新的任务id
	ID int64 `json:"id"`
	// 标题
	Title *string `json:"title"`
	// 负责人
	OwnerID []int64 `json:"ownerId"`
	// 优先级id
	PriorityID *int64 `json:"priorityId"`
	// 计划开始时间
	PlanStartTime *types.Time `json:"planStartTime"`
	// 计划结束时间
	PlanEndTime *types.Time `json:"planEndTime"`
	// 计划工作时长
	PlanWorkHour *int `json:"planWorkHour"`
	// 备注
	Remark *string `json:"remark"`
	// 备注详情
	RemarkDetail *string `json:"remarkDetail"`
	// 备注提及人
	MentionedUserIds []int64 `json:"mentionedUserIds"`
	// 迭代
	IterationID *int64 `json:"iterationId"`
	// 来源
	SourceID *int64 `json:"sourceId"`
	// 问题对象类型id
	IssueObjectTypeID *int64 `json:"issueObjectTypeId"`
	// 问题性质id
	IssuePropertyID *int64 `json:"issuePropertyId"`
	// 参与人
	ParticipantIds []int64 `json:"participantIds"`
	// 关注人
	FollowerIds []int64 `json:"followerIds"`
	// 关注人部门（后端实际转化为人）
	FollowerDeptIds []int64 `json:"followerDeptIds"`
	// 审核人
	AuditorIds []int64 `json:"auditorIds"`
	// 变动的字段列表
	UpdateFields []string `json:"updateFields"`
	// 无码更新入参
	LessUpdateIssueReq map[string]interface{} `json:"lessUpdateIssueReq"`
}

type IssueReportResp struct {
}

type IssueListStatReq struct {
	// 项目id
	ProjectID int64 `json:"projectId"`
}

type IssueListStatResp struct {
	List []*IssueListStatData `json:"list"`
}

// 每日个人完成图数据统计请求结构体
type IssueDailyPersonalWorkCompletionStatReq struct {
	// 开始时间, 开始时间和结束时间可以不传，默认七天
	StartDate *types.Time `json:"startDate"`
	// 结束时间
	EndDate *types.Time `json:"endDate"`
}

// 项目和任务信息统计
type IssueAndProjectCountStatResp struct {
	// 项目未完成的数量
	ProjectNotCompletedCount int64 `json:"projectNotCompletedCount"`
	// 任务未完成的数量
	IssueNotCompletedCount int64 `json:"issueNotCompletedCount"`
	// 参与项目数
	ParticipantsProjectCount int64 `json:"participantsProjectCount"`
	// 参与归档项目数
	FilingParticipantsProjectCount int64 `json:"filingParticipantsProjectCount"`
}

// 每日个人完成图数据统计响应结构体
type IssueDailyPersonalWorkCompletionStatResp struct {
	// 数据列表
	List []*IssueDailyPersonalWorkCompletionStatData `json:"list"`
}

// 创建问题对象类型请求结构体
type CreateIssueObjectTypeReq struct {
	// 组织id
	OrgID int64 `json:"orgId"`
	// 语言编号
	LangCode string `json:"langCode"`
	// 类型名称
	Name string `json:"name"`
	// 排序
	Sort int `json:"sort"`
	// 对象类型id
	ProjectObjectTypeID int64 `json:"projectObjectTypeId"`
	// 描述
	Remark string `json:"remark"`
	// 状态,  1可用,2禁用
	Status int `json:"status"`
	// 创建人
	Creator int64 `json:"creator"`
	// 创建时间
	CreateTime types.Time `json:"createTime"`
	// 更新人
	Updator int64 `json:"updator"`
	// 更新时间
	UpdateTime types.Time `json:"updateTime"`
	// 乐观锁
	Version int `json:"version"`
	// 是否删除,1是,2否
	IsDelete int `json:"isDelete"`
}

// 更新问题对象类型请求结构体
type UpdateIssueObjectTypeReq struct {
	// 主键
	ID int64 `json:"id"`
	// 组织id
	OrgID int64 `json:"orgId"`
	// 语言编号
	LangCode string `json:"langCode"`
	// 类型名称
	Name string `json:"name"`
	// 排序
	Sort int `json:"sort"`
	// 对象类型id
	ProjectObjectTypeID int64 `json:"projectObjectTypeId"`
	// 描述
	Remark string `json:"remark"`
	// 状态,  1可用,2禁用
	Status int `json:"status"`
	// 创建人
	Creator int64 `json:"creator"`
	// 创建时间
	CreateTime types.Time `json:"createTime"`
	// 更新人
	Updator int64 `json:"updator"`
	// 更新时间
	UpdateTime types.Time `json:"updateTime"`
	// 乐观锁
	Version int `json:"version"`
	// 是否删除,1是,2否
	IsDelete int `json:"isDelete"`
}

// 删除问题对象类型请求结构体
type DeleteIssueObjectTypeReq struct {
	// 主键
	ID int64 `json:"id"`
	// 组织Id 暂时不用传 之后用户校验的时候比较是否包含这个orgId 操作的时候是否有当前orgId的权限
	OrgID *int64 `json:"orgId"`
}

type HandleOldIssueToNewReq struct {
	OrgID     int64   `json:"orgId"`
	ProjectID int64   `json:"projectId"`
	IssueIds  []int64 `json:"issueIds"`
}

type CreateProjectDetailReq struct {
}

// 项目任务关联的状态请求结构体
type ProjectIssueRelatedStatusReq struct {
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 表 id
	TableID string `json:"tableId"`
}

type ProjectObjectTypeList struct {
}

type ProjectSupportObjectTypeListResp struct {
}

type CreateProjectObjectTypeReq struct {
}

type UpdateProjectObjectTypeReq struct {
}

type DeleteProjectObjectTypeReq struct {
}

type ProjectSupportObjectTypeListReq struct {
}

type ProjectObjectTypeWithProjectList struct {
}
type ProjectResourceInfoReq struct {
	// 资源id
	ResourceID int64 `json:"resourceId"`
	// 应用id
	AppID int64 `json:"appId"`
}

type ProjectTypeCategoryList struct {
}

type ProjectTypeListResp struct {
}

type OperateProjectResp struct {
	// 是否成功
	IsSuccess interface{} `json:"isSuccess"`
}

type ProjectStatisticsResp struct {
}

type ProjectObjectTypesReq struct {
}

// 移出项目成员
type RemoveProjectMemberReq struct {
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 要移除的项目成员列表
	MemberIds []int64 `json:"memberIds"`
	// 要移除的项目部门列表
	MemberForDepartmentID []int64 `json:"memberForDepartmentId"`
}

type ProjectUserListReq struct {
	// 项目id
	ProjectID int64 `json:"projectId"`
}

type ProjectUserListResp struct {
	Total int64          `json:"total"`
	List  []*ProjectUser `json:"list"`
}

type OrgProjectMemberListReq struct {
	// 关联类型(1负责人2关注人3全部，默认全部)
	RelationType *int64 `json:"relationType"`
	// 名字
	Name *string `json:"name"`
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 是否忽略离职人员(默认否)
	IgnoreDelete *bool `json:"ignoreDelete"`
}

type OrgProjectMemberListResp struct {
	// 总数
	Total int64 `json:"total"`
	// 数据
	List []*OrgProjectMemberInfoResp `json:"list"`
}

// 更新角色操作权限
type UpdateRolePermissionOperationReq struct {
	// 角色id
	RoleID int64 `json:"roleId"`
	// 项目id
	ProjectID *int64 `json:"projectId"`
	// 所有涉及到的更改的权限组
	UpdatePermissions []*EveryPermission `json:"updatePermissions"`
}

type CancelArchiveIssueReq struct {
}

type CancelArchiveIssueBatchResp struct {
}

// 导入任务
type ImportIssuesReq struct {
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 项目表id
	TableID string `json:"tableId"`
	// 迭代id
	IterationID *int64 `json:"iterationId"`
	// excel地址
	URL string `json:"url"`
	// url类型, 1 网址，2 本地dist路径
	URLType int `json:"urlType"`
}

// 获取任务性质列表请求结构体
type IssuePropertysReq struct {
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 项目对象类型id
	ProjectObjectTypeID *int64 `json:"projectObjectTypeId"`
}

type IssueRestInfoResp struct {
}

// 更新任务标签关联请求结构体
type UpdateIssueTagsReq struct {
	// 任务id
	ID int64 `json:"id"`
	// 新关联的标签列表，addTags和delTags可以同时存在
	AddTags []*IssueTagReqInfo `json:"addTags"`
	// 要取消关联的标签列表
	DelTags []*IssueTagReqInfo `json:"delTags"`
}

// 获取任务资源请求结构体
type GetIssueResourcesReq struct {
	// 任务id
	IssueID int64 `json:"issueId"`
}

// 任务添加文件资源
type CreateIssueResourceReq struct {
	// 项目id
	ProjectId int64 `json:"projectId"`
	// 任务id
	IssueId int64 `json:"issueId"`
	// 资源路径
	ResourcePath string `json:"resourcePath"`
	// 资源大小，单位B
	ResourceSize int64 `json:"resourceSize"`
	// 文件名
	FileName string `json:"fileName"`
	// 文件后缀
	FileSuffix string `json:"fileSuffix"`
	// md5
	Md5 *string `json:"md5"`
	// bucketName
	BucketName *string `json:"bucketName"`
	// policyType
	PolicyType int `json:"policyType"`
	// 资源存储方式
	ResourceType int `json:"resourceType"`
}

// 删除子任务请求结构体
type DeleteIssueResourceReq struct {
	// 任务id'
	IssueID int64 `json:"issueId"`
	// 关联资源id列表
	DeletedResourceIds []int64 `json:"deletedResourceIds"`
}
type RelatedIssueListReq struct {
}

type AddIssueAttachmentReq struct {
	IssueID     int64   `json:"issueId"`
	ResourceIds []int64 `json:"resourceIds"`
}

type UpdateIssueBeforeAfterIssuesReq struct {
}

type BeforeAfterIssueListResp struct {
}

// 获取任务来源列表请求结构体
type IssueSourcesReq struct {
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 项目对象类型id
	ProjectObjectTypeID *int64 `json:"projectObjectTypeId"`
}

// 任务来源列表响应结构体
type IssueSourceList struct {
	Total int64          `json:"total"`
	List  []*IssueSource `json:"list"`
}

// 创建任务来源请求结构体
type CreateIssueSourceReq struct {
	// 组织id
	OrgID int64 `json:"orgId"`
	// 语言编号
	LangCode string `json:"langCode"`
	// 名称
	Name string `json:"name"`
	// 排序
	Sort int `json:"sort"`
	// 对象类型id
	ProjectObjectTypeID int64 `json:"projectObjectTypeId"`
	// 描述
	Remark string `json:"remark"`
	// 状态,  1可用,2禁用
	Status int `json:"status"`
	// 创建人
	Creator int64 `json:"creator"`
	// 创建时间
	CreateTime types.Time `json:"createTime"`
	// 更新人
	Updator int64 `json:"updator"`
	// 更新时间
	UpdateTime types.Time `json:"updateTime"`
	// 乐观锁
	Version int `json:"version"`
	// 是否删除,1是,2否
	IsDelete int `json:"isDelete"`
}

// 更新任务来源请求结构体
type UpdateIssueSourceReq struct {
	// 主键
	ID int64 `json:"id"`
	// 组织id
	OrgID int64 `json:"orgId"`
	// 语言编号
	LangCode string `json:"langCode"`
	// 名称
	Name string `json:"name"`
	// 排序
	Sort int `json:"sort"`
	// 对象类型id
	ProjectObjectTypeID int64 `json:"projectObjectTypeId"`
	// 描述
	Remark string `json:"remark"`
	// 状态,  1可用,2禁用
	Status int `json:"status"`
	// 创建人
	Creator int64 `json:"creator"`
	// 创建时间
	CreateTime types.Time `json:"createTime"`
	// 更新人
	Updator int64 `json:"updator"`
	// 更新时间
	UpdateTime types.Time `json:"updateTime"`
	// 乐观锁
	Version int `json:"version"`
	// 是否删除,1是,2否
	IsDelete int `json:"isDelete"`
}

// 删除任务来源请求结构体
type DeleteIssueSourceReq struct {
	// 主键
	ID int64 `json:"id"`
	// 组织Id 暂时不用传 之后用户校验的时候比较是否包含这个orgId 操作的时候是否有当前orgId的权限
	OrgID *int64 `json:"orgId"`
}

type CreateIssueViewReq struct {
	// 项目 id
	ProjectID *int64 `json:"projectId"`
	// 视图配置
	Config string `json:"config"`
	// 视图备注
	Remark *string `json:"remark"`
	// 是否私有
	IsPrivate *bool `json:"isPrivate"`
	// 所属任务类型 id：需求、任务、缺陷的 id 值
	ProjectObjectTypeID *int64 `json:"projectObjectTypeId"`
	// 视图名称
	ViewName string `json:"viewName"`
	// 类型，1：表格视图，2：看板视图，3：照片视图
	Type *int `json:"type"`
	// 视图排序
	Sort *int64 `json:"sort"`
}

type DeleteIssueViewReq struct {
	// 主键id，根据主键删除
	ID int64 `json:"id"`
}

// 获取工时记录的列表请求参数
type GetIssueWorkHoursListReq struct {
	// 页码
	Page *int64 `json:"page"`
	// 每页条数
	Size *int64 `json:"size"`
	// 关联的任务id
	IssueID int64 `json:"issueId"`
	// 记录类型：1预估工时记录，2实际工时记录，3子预估工时
	Type int64 `json:"type"`
}

// 获取的工时记录列表结果
type GetIssueWorkHoursListResp struct {
	// 总数量
	Total int64 `json:"total"`
	// 任务列表
	List []*IssueWorkHours `json:"list"`
}

// 启用/关闭工时记录功能接口请求体
type DisOrEnableIssueWorkHoursReq struct {
	// 项目id，针对这个项目关闭/启用工时功能
	ProjectID int64 `json:"projectId"`
	// 是否关闭工时功能：1启用,2关闭
	Enable int64 `json:"enable"`
}

// 迭代状态类型统计响应结构体
type IterationStatusTypeStatResp struct {
	// 状态为未开始的数量
	NotStartTotal int64 `json:"notStartTotal"`
	// 状态为进行中的数量
	ProcessingTotal int64 `json:"processingTotal"`
	// 状态为已完成的数量
	CompletedTotal int64 `json:"completedTotal"`
	// 总数量
	Total int64 `json:"total"`
}

// 迭代和任务关联请求结构体
type IterationIssueRealtionReq struct {
	// 迭代id
	IterationID int64 `json:"iterationId"`
	// 要添加的任务id列表（除特性任务）
	AddIssueIds []int64 `json:"addIssueIds"`
	// 要移除的任务id列表
	DelIssueIds []int64 `json:"delIssueIds"`
}

// 创建优先级请求结构体
type CreatePriorityReq struct {
	// 组织id,全局的填0
	OrgID int64 `json:"orgId"`
	// 语言编号
	LangCode string `json:"langCode"`
	// 名称
	Name string `json:"name"`
	// 类型,1项目优先级,2:需求/任务等优先级
	Type int `json:"type"`
	// 排序
	Sort int `json:"sort"`
	// 背景颜色
	BgStyle string `json:"bgStyle"`
	// 字体颜色
	FontStyle string `json:"fontStyle"`
	// 是否默认,1是,2否
	IsDefault int `json:"isDefault"`
	// 描述
	Remark string `json:"remark"`
	// 状态,  1可用,2禁用
	Status int `json:"status"`
	// 创建人
	Creator int64 `json:"creator"`
	// 创建时间
	CreateTime types.Time `json:"createTime"`
	// 更新人
	Updator int64 `json:"updator"`
	// 更新时间
	UpdateTime types.Time `json:"updateTime"`
	// 乐观锁
	Version int `json:"version"`
	// 是否删除,1是,2否
	IsDelete int `json:"isDelete"`
}

// 更新优先级请求结构体
type UpdatePriorityReq struct {
	// 主键
	ID int64 `json:"id"`
	// 组织id,全局的填0
	OrgID int64 `json:"orgId"`
	// 语言编号
	LangCode string `json:"langCode"`
	// 名称
	Name string `json:"name"`
	// 类型,1项目优先级,2:需求/任务等优先级
	Type int `json:"type"`
	// 排序
	Sort int `json:"sort"`
	// 背景颜色
	BgStyle string `json:"bgStyle"`
	// 字体颜色
	FontStyle string `json:"fontStyle"`
	// 是否默认,1是,2否
	IsDefault int `json:"isDefault"`
	// 描述
	Remark string `json:"remark"`
	// 状态,  1可用,2禁用
	Status int `json:"status"`
	// 创建人
	Creator int64 `json:"creator"`
	// 创建时间
	CreateTime types.Time `json:"createTime"`
	// 更新人
	Updator int64 `json:"updator"`
	// 更新时间
	UpdateTime types.Time `json:"updateTime"`
	// 乐观锁
	Version int `json:"version"`
	// 是否删除,1是,2否
	IsDelete int `json:"isDelete"`
}

// 删除优先级请求结构体
type DeletePriorityReq struct {
	// 主键
	ID int64 `json:"id"`
	// 组织Id 暂时不用传 之后用户校验的时候比较是否包含这个orgId 操作的时候是否有当前orgId的权限
	OrgID *int64 `json:"orgId"`
}

type ConvertCodeResp struct {
	// 项目code
	Code string `json:"code"`
}

type ProjectChatListReq struct {
	// 每页数量(不传默认10条)
	PageSize *int64 `json:"pageSize"`
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 上次分页最后一条关联id（不传默认第一页）
	LastRelationID *int64 `json:"lastRelationId"`
	// 搜索内容
	Name *string `json:"name"`
}

type ChatListResp struct {
	// 数量
	Total int64 `json:"total"`
	// 列表
	List []*ChatData `json:"list"`
}

type UnrelatedChatListReq struct {
	// 每页数量(不传默认10条)
	PageSize *int64 `json:"pageSize"`
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 上次分页最后一条外部群聊id（不传默认第一页）
	LastOutChatID *string `json:"lastOutChatId"`
	// 搜索内容
	Name *string `json:"name"`
}

type UpdateRelateChat struct {
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 群聊外部id
	OutChatIds []string `json:"outChatIds"`
}
type SimpleProjectInfo struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	ProjectTypeID int64  `json:"projectTypeId"`
}

// 更新任务响应结构体
type UpdateIssueResp struct {
	// 任务id
	ID int64 `json:"id"`
}

type IssueRestInfoReq struct {
}

// 新增项目资源
type CreateProjectResourceReq struct {
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 文件夹id
	FolderID int64 `json:"folderId"`
	// 资源路径
	ResourcePath string `json:"resourcePath"`
	// 资源大小，单位B
	ResourceSize int64 `json:"resourceSize"`
	// 文件名
	FileName string `json:"fileName"`
	// 文件后缀
	FileSuffix string `json:"fileSuffix"`
	// md5
	Md5 *string `json:"md5"`
	// bucketName
	BucketName *string `json:"bucketName"`
	// policyType
	PolicyType int `json:"policyType"`
	// 资源存储方式
	ResourceType int `json:"resourceType"`
}

type QuitResult struct {
	// 是否退出
	IsQuitted interface{} `json:"isQuitted"`
}

// Oss申请signUrl响应结构体
type OssApplySignURLResp struct {
	// signUrl
	SignURL string `json:"signUrl"`
}

// Oss申请signUrl请求结构体
type OssApplySignURLReq struct {
	// 文件url
	URL string `json:"url"`
}

type NoticeListReq struct {
	// 通知类型, 1项目通知,2组织通知,
	Type *int `json:"type"`
	// 是否是@我的相关notice(不传或传1表示普通通知，2表示@我的)
	IsCallMe *int `json:"isCallMe"`
	// 上次分页的最后一条id
	LastID *int64 `json:"lastId"`
	// 项目id
	ProjectID *int64 `json:"projectId"`
	// 任务id
	IssueID *int64 `json:"issueId"`
}

// 列表响应结构体
type NoticeList struct {
	Total int64     `json:"total"`
	List  []*Notice `json:"list"`
}

// 创建角色请求结构体
type CreateRoleReq struct {
	// 角色组1组织角色2项目角色
	RoleGroupType int `json:"roleGroupType"`
	// 名称
	Name string `json:"name"`
	// 描述
	Remark *string `json:"remark"`
	// 是否只读 1只读 2可编辑
	IsReadonly *int `json:"isReadonly"`
	// 是否可以变更权限,1可以,2不可以
	IsModifyPermission *int `json:"isModifyPermission"`
	// 是否默认角色,1是,2否
	IsDefault *int `json:"isDefault"`
	// 状态,  1可用,2禁用
	Status *int `json:"status"`
	// 项目id
	ProjectID *int64 `json:"projectId"`
}

// 删除角色请求结构体
type DelRoleReq struct {
	// 角色ID
	RoleIds []int64 `json:"roleIds"`
	// 如果删除项目角色，projectId必填，如果删除组织角色，projectId可以不传或者传0
	ProjectID *int64 `json:"projectId"`
}

// 更新角色请求结构体
type UpdateRoleReq struct {
	// 角色ID
	RoleID int64 `json:"roleId"`
	// 角色名称, 非必填，为空则不更新
	Name *string `json:"name"`
}

// 移除组织
type RemoveOrgMemberReq struct {
	// 要移除的组织成员列表
	MemberIds []int64 `json:"memberIds"`
}

type FunctionConfigResp struct {
	// 功能
	FunctionCodes []string `json:"functionCodes"`
}

type UpdateOrgBasicShowSettingReq struct {
}

// 从飞书方同步成员、部门等信息
type SyncUserInfoFromFeiShuReq struct {
	// 是否需要同步用户姓名
	NeedSyncName bool `json:"needSyncName"`
	// 是否需要同步用户头像
	NeedSyncAvatar bool `json:"needSyncAvatar"`
	// 是否需要同步部门架构信息
	NeedSyncDepartment bool `json:"needSyncDepartment"`
	// 用户来源、渠道。如果是飞书，则传 `fs`
	SourceChannel string `json:"sourceChannel"`
}

type UpdateDepartmentForInviteReq struct {
}

type RegisterWebSiteContactReq struct {
}

type ProcessStatusList struct {
}

type CreateProcessStatusReq struct {
}

type UpdateProcessStatusReq struct {
}

type DeleteProcessStatusReq struct {
}

type CreateCustomFieldReq struct {
	// 名称
	Name string `json:"name"`
	// 类型(1文本类型2单选框3多选框4日期选框5人员选择6是非选择7数字框)
	FieldType int `json:"fieldType"`
	// 选项值
	FieldValue interface{} `json:"fieldValue"`
	// 是否加入组织字段库(1是2否，不选默认为否)
	IsOrgField *int `json:"isOrgField"`
	// 字段描述
	Remark *string `json:"remark"`
	// 项目id
	ProjectID *int64 `json:"projectId"`
	// 任务类型
	ProjectObjectTypeID *int64 `json:"projectObjectTypeId"`
}

type DeleteCustomFieldReq struct {
	// 自定义字段id
	FieldID int64 `json:"fieldId"`
	// 项目id(传了表示删除项目关联)
	ProjectID *int64 `json:"projectId"`
	// 任务类型
	ProjectObjectTypeID *int64 `json:"projectObjectTypeId"`
}

type UseOrgCustomFieldReq struct {
	// 自定义字段id集合
	FieldIds []int64 `json:"fieldIds"`
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 任务类型
	ProjectObjectTypeID *int64 `json:"projectObjectTypeId"`
}
type UpdateCustomFieldReq struct {
	// 自定义字段id
	FieldID int64 `json:"fieldId"`
	// 名称
	Name *string `json:"name"`
	// 值
	FieldValue interface{} `json:"fieldValue"`
	// 是否加入组织字段库(1是2否)
	IsOrgField *int `json:"isOrgField"`
	// 字段描述
	Remark *string `json:"remark"`
	// 更新的字段
	UpdateField []string `json:"updateField"`
}

type ChangeProjectCustomFieldStatusReq struct {
	// 自定义字段id
	FieldID int64 `json:"fieldId"`
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 任务类型
	ProjectObjectTypeID int64 `json:"projectObjectTypeId"`
	// 状态(1启用2禁用)
	Status int `json:"status"`
}

type CustomFieldListReq struct {
	// 页码
	Page *int `json:"page"`
	// 每页条数
	Size *int `json:"size"`
	// 项目id
	ProjectID *int64 `json:"projectId"`
	// 任务类型
	ProjectObjectTypeID *int64 `json:"projectObjectTypeId"`
	// 是否是没有被项目(敏捷包括任务类型)的字段（和项目id联合使用，只有传了项目id才生效，1是2否，默认否）
	IsUsedCurrentProject *int `json:"isUsedCurrentProject"`
	// 是否属于组织字段库(1组织2项目3系统)
	IsOrgField []int `json:"isOrgField"`
	// 名称
	Name *string `json:"name"`
	// 排序类型（1创建时间正序2创建时间倒序3添加到项目时间正序4添加到项目时间倒序。默认创建时间正序）
	OrderType *int `json:"orderType"`
}

type CustomFieldListResp struct {
	Total int64          `json:"total"`
	List  []*CustomField `json:"list"`
}

// vo.UpdateIssueCustomFieldReq
type UpdateIssueCustomFieldReq struct {
}

type ProjectFieldViewReq struct {
	// 视图类型（1看板2列表3表格）
	ViewType int `json:"viewType"`
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 任务类型
	ProjectObjectTypeID int64 `json:"projectObjectTypeId"`
}

type ProjectFieldViewResp struct {
	// 关闭的默认的字段code
	ClosedDefaultFields []string `json:"closedDefaultFields"`
	// 自定义字段
	CustomFields []*CustomField `json:"customFields"`
}

type UpdateProjectFieldViewReq struct {
	// 视图类型（1看板2列表3表格）
	ViewType int `json:"viewType"`
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 任务类型
	ProjectObjectTypeID int64 `json:"projectObjectTypeId"`
	// 关闭的默认字段
	ClosedDefaultFields []string `json:"closedDefaultFields"`
	// 关闭的自定义字段
	ClosedCustomFields []int64 `json:"closedCustomFields"`
}

type CustomField struct {
	// id
	ID int64 `json:"id"`
	// 名称
	Name string `json:"name"`
	// 类型
	FieldType int `json:"fieldType"`
	// 值
	FieldValue []map[string]interface{} `json:"fieldValue"`
	// 是否是组织字段(1组织2项目3系统)
	IsOrgField int `json:"isOrgField"`
	// 启用状态（1启用2禁用）对于项目而言
	Status int `json:"status"`
	// 字段描述
	Remark string `json:"remark"`
	// 使用项目
	ProjectList []*SimpleProjectInfo `json:"projectList"`
	// 编辑人信息
	UpdatorInfo *UserIDInfo `json:"updatorInfo"`
	// 编辑时间
	UpdateTime types.Time `json:"updateTime"`
	// 是否展示(1展示2否，在项目视图里使用该字段)
	IsDisplay int `json:"isDisplay"`
}

type InitExistOrgReq struct {
	NeedOcrConfig    int  `json:"needOcrConfig"`
	NeedSummaryTable int  `json:"needSummaryTable"`
	NeedPriority     int  `json:"needPriority"`
	NeedSetToPaid    *int `json:"needSetToPaid"`
	// 登录状态。-1 表示非登录状态下调用接口；1 登录状态。默认 1。
	LoginMode *int   `json:"loginMode"`
	OrgID     *int64 `json:"OrgId"`
	UserID    *int64 `json:"UserId"`
	// 查询 token 对应的一些信息。0不查询，1表示查询。默认 0。
	QueryAuthTokenInfo *int `json:"QueryAuthTokenInfo"`
}

// 向钉钉群发送告警日志接口请求参数
type DingTalkInfoReq struct {
	// 日志信息。没有则传空字符串。
	Content string `json:"content"`
	// 其他信息。没有则传空字符串。
	Other string `json:"other"`
}

// DingTalk第三方扫码登录
type AuthDingCodeReq struct {
	Code string `json:"code"`
}

// 选择组织返回登录信息
type AuthForChosenOrgReq struct {
	OutUserID string `json:"outUserId"`
	OrgID     int64  `json:"orgId"`
}

// DingTalk免登陆 Code 登录验证响应结构体
type AuthResp struct {
	// 持久化登录信息的Token
	Token string `json:"token"`
	// 组织id
	OrgID int64 `json:"orgId"`
	// 组织名称
	OrgName string `json:"orgName"`
	// 组织code
	OrgCode string `json:"orgCode"`
	// 用户id
	UserID int64 `json:"userId"`
	// 用户姓名
	Name string `json:"name"`
}
