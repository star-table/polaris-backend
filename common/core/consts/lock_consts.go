package consts

const (
	LockCacheKeyPrefix = CacheKeyPrefix + "lock:"

	//添加用户配置的锁 + 用户id
	AddUserConfigLock = LockCacheKeyPrefix + "AddUserConfigLock:"
	//初始化用户锁 + 用户id
	InitUserLock = LockCacheKeyPrefix + "InitUserLock:"
	//添加任务时同步项目成员的锁 + projectId
	AddIssueScheduleProjectMemberLock = LockCacheKeyPrefix + "AddIssueScheduleProjectMemberLock:"

	//任务关联锁 + issueId + relationType
	AddIssueRelationLock = LockCacheKeyPrefix + "AddIssueRelationLock:"

	//文件资源锁 + orgId + relationType
	AddResourceLock = LockCacheKeyPrefix + "AddResourceLock:"

	//项目关联锁 + projectId
	AddProjectRelationLock = LockCacheKeyPrefix + "AddProjectRelationLock:"

	//项目自定义字段关联锁 + projectId
	AddProjectCustomFieldLock = LockCacheKeyPrefix + "AddProjectCustomFieldLock:"

	//任务和标签关联 + issueId
	AddIssueTagsLock = LockCacheKeyPrefix + "AddIssueTagsLock:"

	//更新用户角色权限锁
	UpdateUserOrgRoleLock = LockCacheKeyPrefix + "LockCacheKeyPrefix:"

	//ID申请时的分布式锁前缀
	IdServiceIdKey = LockCacheKeyPrefix + "IdServiceIdKey:"

	//钉钉企业初始化
	DingTalkCorpInitKey = LockCacheKeyPrefix + "DingTalkCorpInitKey:"

	// 企业微信
	WeiXinTalkCorpInitKey = LockCacheKeyPrefix + "WeiXinTalkCorpInitKey:"

	//飞书企业初始化
	FeiShuCorpInitKey = LockCacheKeyPrefix + "FeiShuCorpInitKey:"
	//泳道锁
	ProjectObjectTypeLockKey = LockCacheKeyPrefix + "DeleteProjectObjectType:"

	//飞书获取AppAccessToken时的锁
	FeiShuGetAppAccessTokenLockKey = LockCacheKeyPrefix + "FeiShuGetAppAccessTokenLockKey"

	//飞书获取TenantAccessToken时的锁
	FeiShuGetTenantAccessTokenLockKey = LockCacheKeyPrefix + "FeiShuGetTenantAccessTokenLockKey:"

	//组织权限补偿 + orgId
	RolePermissionPathCompensatoryLockKey = LockCacheKeyPrefix + "RolePermissionPathCompensatoryLockKey:"

	//用户和组织关联锁 + orgId + userId
	UserAndOrgRelationLockKey = LockCacheKeyPrefix + "UserAndOrgRelationLockEky:"

	//用户和部门关联锁 + orgId + departmentId
	UserAndDepartmentRelationLockKey = LockCacheKeyPrefix + "UserAndDepartmentRelationLockKey:"

	//创建标签锁 + projectId
	CreateTagLock = LockCacheKeyPrefix + "CreateTagLock:"

	//编辑角色锁 + orgId + projectId(如果projectId没有则为0)
	ModifyRoleLock = LockCacheKeyPrefix + "ModifyRoleLock:"

	//编辑角色权限锁 + orgId + roleId
	ModifyRolePermissionLock = LockCacheKeyPrefix + "ModifyRolePermissionLock:"

	//用户登录账户（手机号，邮箱）绑定锁 + loginName
	UserBindLoginNameLock = LockCacheKeyPrefix + "UserBindLoginNameLock:"
	UserRegisterNameLock  = LockCacheKeyPrefix + "UserRegisterNameLock:"

	//飞书日历创建锁 + projectId
	CreateCalendarLock = LockCacheKeyPrefix + "CreateCalendarLock:"

	//飞书日程创建锁 + issueId
	CreateCalendarEventLock = LockCacheKeyPrefix + "CreateCalendarEventLock:"
	// 任务创建群聊加锁 + issueId
	CreateIssueChatLock = LockCacheKeyPrefix + "CreateIssueChatLock:"

	//任务相关操作锁 + issueId（包括移动任务，创建子任务，创建标签关联）
	IssueRelateOperationLock = LockCacheKeyPrefix + "IssueRelateOperationLock:"

	//项目标签锁 + projectId
	CreateProjectTagLock = LockCacheKeyPrefix + "CreateProjectTagLock:"

	//项目文件和目录关联的锁 + targetFolderId
	UpdateResourceFolderLock = LockCacheKeyPrefix + "UpdateResourceFolderLock:"

	//添加项目群聊锁 + projectId
	AddProjectChatLock = LockCacheKeyPrefix + "AddProjectChatLock:"

	//添加迭代状态关联锁 + iterationId
	AddIterationRelationLock = LockCacheKeyPrefix + "AddIterationRelationLock:"

	//添加资源关联锁 + projectId + issueId
	AddResourceRelationLock = LockCacheKeyPrefix + "AddResourceRelationLock:"

	//检验飞书用户是否在授权范围内
	CheckFsUserLock = LockCacheKeyPrefix + "CheckFsUserLock:"

	//飞书捷径 + trigger_key
	FeishuShortcutLock = LockCacheKeyPrefix + "FeishuShortcutLock:"

	//更新部门角色权限锁
	UpdateDepartmentOrgRoleLock = LockCacheKeyPrefix + "UpdateDepartmentOrgRoleLock:"

	//获取组织开放平台信息锁
	GetOrgAppTicketLock = LockCacheKeyPrefix + "GetOrgAppTicketLock:"

	//飞书同步部门加锁
	AddFeishuDepartmentLock = LockCacheKeyPrefix + "AddFeishuDepartmentLock:"

	//迭代加锁
	IterationSortLock = LockCacheKeyPrefix + "IterationSortLock:"

	//新版本飞书组织初始化锁
	NewFeiShuCorpInitKey = LockCacheKeyPrefix + "NewFeiShuCorpInitKey:"

	//飞书获取jsapi_ticket时的锁
	FeiShuGetJSApiTicketLockKey = LockCacheKeyPrefix + "FeiShuGetJSApiTicketLockKey:"

	TransferOrgLockKey = LockCacheKeyPrefix + "TransferOrg:"

	// 批量创建任务异步任务处理
	CreateIssueAsyncTaskLockKey = LockCacheKeyPrefix + "CreateIssueAsyncTaskLockKey:"
)
