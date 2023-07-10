package errs

import "github.com/star-table/common/core/errors"

var (
	// 成功
	OK = errors.OK

	// token错误
	RequestError = errors.RequestError

	// 认证错误
	Unauthorized = errors.Unauthorized

	// 禁止访问
	ForbiddenAccess = errors.ForbiddenAccess

	// 请求地址不存在
	PathNotFound = errors.PathNotFound

	// 不支持该方法
	MethodNotAllowed = errors.MethodNotAllowed

	// Token过期
	TokenExpires = errors.TokenExpires

	// 请求参数错误
	ServerError = errors.ServerError

	// 过载保护,服务暂不可用
	ServiceUnavailable = errors.ServiceUnavailable

	// 服务调用超时
	Deadline = errors.Deadline

	// 超出限制
	LimitExceed = errors.LimitExceed

	// 参数错误
	ParamError = errors.ParamError

	// 文件过大
	FileTooLarge = errors.FileTooLarge

	// 文件类型错误
	FileTypeError = errors.FileTypeError

	// 文件或目录不存在
	FileNotExist = errors.FileNotExist

	// 文件路径为空
	FilePathIsNull = errors.FilePathIsNull

	// 读取文件失败
	FileReadFail = errors.FileReadFail
	// 文件保存失败
	FileSaveFail = errors.AddResultCodeInfoWithSentry(615, "文件保存失败！", "ResultCode.FileSaveFail")

	// 错误未定义
	ErrorUndefined = errors.ErrorUndefined

	// 业务失败
	BusinessFail = errors.BusinessFail

	// 系统异常
	SystemError = errors.SystemError

	// 未知错误
	UnknownError = errors.UnknownError

	// Panic错误
	PanicError = errors.PanicError

	// 数据库错误
	DatabaseError = errors.DatabaseError

	// 业务中详细异常定义
	// 工具异常
	MysqlOperateError          = errors.MysqlOperateError
	RedisOperateError          = errors.RedisOperateError
	GetDistributedLockError    = errors.GetDistributedLockError
	OssError                   = errors.AddResultCodeInfo(300103, "Oss异常", "ResultCode.OssError")
	RocketMQProduceInitError   = errors.RocketMQProduceInitError
	RocketMQSendMsgError       = errors.RocketMQSendMsgError
	RocketMQConsumerInitError  = errors.RocketMQConsumerInitError
	RocketMQConsumerStartError = errors.RocketMQConsumerStartError
	RocketMQConsumerStopError  = errors.RocketMQConsumerStopError

	DbMQSendMsgError         = errors.DbMQSendMsgError
	DbMQCreateConsumerError  = errors.DbMQCreateConsumerError
	DbMQConsumerStartedError = errors.DbMQConsumerStartedError

	KafkaMqSendMsgError           = errors.KafkaMqSendMsgError
	KafkaMqSendMsgCantBeNullError = errors.KafkaMqSendMsgCantBeNullError
	KafkaMqConsumeMsgError        = errors.KafkaMqConsumeMsgError
	KafkaMqConsumeStartError      = errors.KafkaMqConsumeStartError

	RunModeUnsupportUpload = errors.AddResultCodeInfo(300104, "该部署模式不支持本地上传", "ResultCode.RunModeUnsupportUpload")

	SystemBusy = errors.AddResultCodeInfo(300105, "系统繁忙，请稍后重试", "ResultCode.SystemBusy")

	JSONConvertError = errors.AddResultCodeInfo(300201, "Json转换出现异常", "ResultCode.JSONConvertError")
	ObjectCopyError  = errors.AddResultCodeInfo(300202, "对象copy出现异常", "ResultCode.ObjectCopyError")
	CacheProxyError  = errors.AddResultCodeInfo(300203, "缓存代理出现异常", "ResultCode.CacheProxyError")
	ObjectTypeError  = errors.AddResultCodeInfo(300204, "对象类型错误", "ResultCode.ObjectTypeError")

	ApplyIdError         = errors.AddResultCodeInfoWithSentry(300205, "ID申请异常", "ResultCode.ApplyIdError")
	ApplyIdCountTooMany  = errors.AddResultCodeInfoWithSentry(300206, "申请id数量过多", "ResultCode.ApplyIdCountTooMany")
	TypeConvertError     = errors.AddResultCodeInfo(300207, "类型转换出现异常", "ResultCode.TypeConvertError")
	UpdateFiledIsEmpty   = errors.AddResultCodeInfo(300208, "未更新任何信息", "ResultCode.UpdateFiledIsEmpty")
	ProtobufMarshalError = errors.AddResultCodeInfo(300209, "Proto转换出现异常", "ResultCode.ProtobufMarshalError")

	TokenAuthError       = errors.AddResultCodeInfo(300301, "身份认证异常，请重新登录", "ResultCode.TokenAuthError")
	TokenNotExist        = errors.AddResultCodeInfo(300302, "身份认证失败，请重新登录", "ResultCode.TokenNotExist")
	SuiteTicketError     = errors.AddResultCodeInfoWithSentry(300303, "获取SuiteTicket异常", "ResultCode.SuiteTicketError")
	GetContextError      = errors.AddResultCodeInfo(300304, "获取请求上下文异常", "ResultCode.GetContextError")
	TemplateRenderError  = errors.AddResultCodeInfo(300305, "模板解析失败", "ResultCode.TemplateRenderError")
	DecryptError         = errors.AddResultCodeInfo(300401, "参数解密异常", "ResultCode.DecryptError")
	CaptchaError         = errors.AddResultCodeInfo(300402, "图形验证码错误", "ResultCode.CaptchaError")
	DingCodeCacheInvalid = errors.AddResultCodeInfo(300403, "扫码认证已失效，请重新扫码", "ResultCode.DingCodeCacheInvalid")

	MQTTKeyGenError        = errors.AddResultCodeInfoWithSentry(300501, "生成key发生异常", "ResultCode.MQTTKeyGenError")
	MQTTPublishError       = errors.AddResultCodeInfoWithSentry(300502, "MQTT推送消息发生异常", "ResultCode.MQTTPublishError")
	MQTTConnectError       = errors.AddResultCodeInfoWithSentry(300503, "MQTT连接时发生异常", "ResultCode.MQTTConnectError")
	MQTTMissingConfigError = errors.AddResultCodeInfoWithSentry(300504, "MQTT缺少配置", "ResultCode.MQTTMissingConfigError")

	TryDistributedLockError = errors.TryDistributedLockError

	// 业务异常
	InitDbFail                     = errors.AddResultCodeInfoWithSentry(400000, "初始化db失败", "ResultCode.InitDbFail")
	ObjectRecordNotFoundError      = errors.AddResultCodeInfo(400001, "对象记录不存在", "ResultCode.ObjectRecordNotFoundError")
	DingTalkUserInfoNotInitedError = errors.AddResultCodeInfo(400002, "钉钉用户没有初始化", "ResultCode.DingTalkUserInfoNotInitedError")
	UserNotFoundError              = errors.AddResultCodeInfo(400003, "用户信息不存在或已经删除", "ResultCode.UserNotFoundError")
	CacheUserInfoNotExistError     = errors.AddResultCodeInfo(400004, "令牌对应的用户信息不存在", "ResultCode.CacheUserInfoNotExistError")
	PageSizeOverflowMaxSizeError   = errors.AddResultCodeInfo(400005, "请求页长超出最大页长限制", "ResultCode.PageSizeOverflowMaxSizeError")
	OutOfConditionError            = errors.AddResultCodeInfo(400006, "请求条件超出限制", "ResultCode.OutOfConditionError")
	ConditionHandleError           = errors.AddResultCodeInfo(400007, "条件处理异常", "ResultCode.ConditionHandleError")
	ReqParamsValidateError         = errors.AddResultCodeInfo(400008, "请求参数校验异常", "ResultCode.ReqParamsValidateError")
	OrgNotInitError                = errors.AddResultCodeInfoWithSentry(400009, "组织未初始化", "ResultCode.OrgNotInitError")
	UserConfigNotExist             = errors.AddResultCodeInfo(400010, "用户配置不存在", "Result.UserConfigNotExist")
	OrgNotExist                    = errors.AddResultCodeInfo(400011, "组织不存在", "ResultCode.OrgNotExist")
	OrgInitError                   = errors.AddResultCodeInfoWithSentry(400012, "组织初始化异常", "ResultCode.OrgInitError")
	OrgOwnTransferError            = errors.AddResultCodeInfo(400013, "非组织创建者不能更改信息", "ResultCode.OrgOwnTransferError")
	OrgOutInfoNotExist             = errors.AddResultCodeInfo(400014, "组织外部信息不存在", "ResultCode.OrgOutInfoNotExist")
	UserOutInfoNotExist            = errors.AddResultCodeInfo(400015, "用户外部信息不存在", "ResultCode.UserOutInfoNotExist")
	UserOutInfoNotError            = errors.AddResultCodeInfo(400016, "用户外部信息错误", "ResultCode.UserOutInfoNotError")
	OrgCodeAlreadySetError         = errors.AddResultCodeInfo(400017, "组织网址不能二次修改", "ResultCode.OrgCodeAlreadySetError")
	OrgCodeLenError                = errors.AddResultCodeInfo(400018, "组织网址后缀只能输入20个字符,包含数字和英文", "ResultCode.OrgWebSiteSettingLenError")
	OrgCodeExistError              = errors.AddResultCodeInfo(400019, "组织网址后缀已被占用，请重新输入", "ResultCode.OrgCodeExistError")
	OrgAddressLenError             = errors.AddResultCodeInfo(400020, "详情地址不得超过100字", "ResultCode.OrgAddressLenError")
	OrgLogoLenError                = errors.AddResultCodeInfo(400021, "组织logo路径长度不能超过512字", "ResultCode.OrgLogoLenError")
	OrgUserRoleModifyError         = errors.AddResultCodeInfo(400022, "无权修改当前角色", "ResultCode.OrgUserRoleModifyError")
	OrgRoleGroupNotExist           = errors.AddResultCodeInfo(400023, "角色分组不存在", "ResultCode.OrgRoleGroupNotExist")
	OrgUserUnabled                 = errors.AddResultCodeInfo(400024, "您已被当前组织禁止访问，请联系管理员解除限制", "ResultCode.OrgUserUnabled")
	OrgRoleNoExist                 = errors.AddResultCodeInfo(400025, "角色不存在", "ResultCode.OrgRoleNoExist")
	OrgUserDeleted                 = errors.AddResultCodeInfo(400026, "您已被该组织移除", "ResultCode.OrgUserDeleted")
	OrgUserCheckStatusUnabled      = errors.AddResultCodeInfo(400027, "您未通过该组织的审核", "ResultCode.OrgUserCheckStatusUnabled")
	OrgFunctionInvalid             = errors.AddResultCodeInfo(400028, "存在无效功能，请确认功能项", "ResultCode.OrgFunctionInvalid")
	PayLevelNotExist               = errors.AddResultCodeInfo(400029, "该付费等级不存在", "ResultCode.PayLevelNotExist")
	FunctionIsLimitPayLevel        = errors.AddResultCodeInfo(400030, "该功能仅对付费用户开放，请升级后使用", "ResultCode.FunctionIsLimitPayLevel")
	CommonUserCreateProjectLimit   = errors.AddResultCodeInfo(400031, "免费版用户可创建项目上限为6个，请升级至标准版继续创建项目", "ResultCode.CommonUserCreateProjectLimit")
	CommonUserCreateIterationLimit = errors.AddResultCodeInfo(400032, "免费版用户单个项目可创建迭代上限为2个，请升级至标准版继续创建迭代", "ResultCode.CommonUserCreateIterationLimit")
	CommonUserCreateTaskLimit      = errors.AddResultCodeInfo(400033, "标准版用户可创建记录上限为1000条（包括已删除记录），请升级版本以便继续创建任务", "ResultCode.CommonUserCreateTaskLimit")
	OrgUserInvalid                 = errors.AddResultCodeInfo(400034, "抱歉，您不在应用付费授权范围内，请联系管理员开通", "ResultCode.OrgUserInvalid")
	ExportFieldIsNull              = errors.AddResultCodeInfo(400035, "请选择导出字段", "ResultCode.ExportFieldIsNull")
	OrgAdminRoleCannotModify       = errors.AddResultCodeInfo(400036, "无权变更组织超管", "ResultCode.OrgAdminRoleCannotModify")
	OrgManagerRoleCannotModify     = errors.AddResultCodeInfo(400037, "无权变更管理员，请联系组织超管进行变更", "ResultCode.OrgManagerRoleCannotModify")
	OrgInitDoing                   = errors.AddResultCodeInfo(400038, "组织正在初始化，请稍后重试", "ResultCode.OrgInitDoing")

	StatusMustExistInType                  = errors.AddResultCodeInfo(400100, "每项状态分类中都不能为空", "Result.StatusMustExistInType")
	RepeatProjectName                      = errors.AddResultCodeInfo(400101, "项目名重复", "ResultCode.RepeatProjectName")
	NotAllInitOrgUser                      = errors.AddResultCodeInfo(400102, "当前成员或负责人不属于组织成员", "Result.NotAllInitOrgUser")
	ExistingNotFinishedSubTask             = errors.AddResultCodeInfo(400103, "当前任务下还有未完成的子任务", "Result.ExistingSubTask")
	VerifyOrgError                         = errors.AddResultCodeInfo(400104, "存在无效用户，请刷新重试", "Result.VerifyOrgError")
	ProcessNotExist                        = errors.AddResultCodeInfo(400105, "流程不存在", "Result.ProcessNotExist")
	IssueNotExist                          = errors.AddResultCodeInfo(400106, "记录不存在", "Result.IssueNotExist")
	ProcessStatusNotExist                  = errors.AddResultCodeInfo(400107, "流程状态不存在", "Result.ProcessStatusNotExist")
	NotAllowQuitProject                    = errors.AddResultCodeInfo(400108, "管理员不允许退出项目", "Result.NotAllowQuitProject")
	NotProjectParticipant                  = errors.AddResultCodeInfo(400109, "抱歉，您不是当前项目成员", "Result.NotProjectParticipant")
	PriorityNotExist                       = errors.AddResultCodeInfo(400110, "优先级不存在", "Result.PriorityNotExist")
	ProjectPreCodeExist                    = errors.AddResultCodeInfo(400111, "项目前缀编号已存在，请手动输入", "Result.ProjectPreCodeExist")
	RepeatProjectPrecode                   = errors.AddResultCodeInfo(400112, "项目前缀编号重复", "ResultCode.RepeatProjectPrecode")
	CreateProjectTimeError                 = errors.AddResultCodeInfo(400113, "项目截至时间必须大于开始时间", "ResultCode.CreateProjectTimeError")
	ParentIssueNotExist                    = errors.AddResultCodeInfo(400114, "父任务不存在", "ResultCode.ParentIssueNotExist")
	ExistingSubTask                        = errors.AddResultCodeInfo(400115, "删除失败，当前任务下还有未删除的子任务", "Result.ExistingSubTask")
	IssueAlreadyBeDeleted                  = errors.AddResultCodeInfo(400116, "任务不存在或已被删除，无法操作", "Result.IssueAlreadyBeDeleted")
	ProcessProcessStatusRelationError      = errors.AddResultCodeInfo(400117, "流程状态关联异常", "Result.ProcessProcessStatusRelationError")
	ProcessProcessStatusInitStatueNotExist = errors.AddResultCodeInfo(400118, "任务初始状态不存在", "Result.ProcessProcessStatusInitStatueNotExist")
	ProjectNotExist                        = errors.AddResultCodeInfo(400119, "项目不存在或已被删除", "Result.ProjectNotExist")
	RoleNotExist                           = errors.AddResultCodeInfo(400120, "角色不存在", "Result.RoleNotExist")
	RoleOperationNotExist                  = errors.AddResultCodeInfo(400121, "角色操作不存在", "Result.RoleOperationNotExist")
	GetUserRoleError                       = errors.AddResultCodeInfo(400122, "获取用户角色时发生异常", "Result.GetUserRoleError")
	ProjectNotInit                         = errors.AddResultCodeInfo(400123, "项目尚未初始化", "Result.ProjectNotInit")
	GetUserInfoError                       = errors.AddResultCodeInfo(400124, "获取用户信息异常", "Result.GetUserInfoError")
	IssueCondAssemblyError                 = errors.AddResultCodeInfo(400125, "任务查询条件封装异常", "Result.IssueCondAssemblyError")
	IssueDetailNotExist                    = errors.AddResultCodeInfo(400126, "任务详情不存在", "Result.IssueDetailNotExist")
	AlreadyStarProject                     = errors.AddResultCodeInfo(400127, "项目已关注", "Result.AlreadyStarProject")
	NotYetStarProject                      = errors.AddResultCodeInfo(400128, "项目尚未关注", "Result.NotYetStarProject")
	TargetNotExist                         = errors.AddResultCodeInfo(400129, "操作对象不存在", "Result.TargetNotExist")
	InvalidResourceType                    = errors.AddResultCodeInfo(400130, "资源类型有误", "Result.InvalidResourceType")
	ProjectObjectTypeProcessNotExist       = errors.AddResultCodeInfo(400131, "项目对象类型对应的流程不存在", "Result.ProjectObjectTypeProcessNotExist")
	IterationExistingNotFinishedTask       = errors.AddResultCodeInfo(400132, "当前迭代存在未完成的任务", "Result.IterationExistingNotFinishedTask")
	ProjectTypeNotExist                    = errors.AddResultCodeInfo(400133, "项目类型不存在", "Result.ProjectTypeNotExist")
	OssPolicyTypeError                     = errors.AddResultCodeInfo(400134, "错误的策略类型", "Result.OssPolicyTypeError")
	RelationIssueError                     = errors.AddResultCodeInfo(400135, "关联的任务有误", "Result.RelationIssueError")
	ParentIssueRelationChildIssueError     = errors.AddResultCodeInfo(400136, "父子任务不能关联", "Result.ParentIssueRelationChildIssueError")
	IterationNotExist                      = errors.AddResultCodeInfo(400137, "迭代不存在", "Result.IterationNotExist")
	ProjectNotRelatedError                 = errors.AddResultCodeInfo(400138, "项目未关联对应的资源", "Result.ProjectNotRelatedError")
	SourceNotExist                         = errors.AddResultCodeInfo(400139, "来源不存在", "Result.SourceNotExist")
	IssueObjectTypeNotExist                = errors.AddResultCodeInfo(400140, "任务类型不存在", "Result.IssueObjectTypeNotExist")
	ResourceNotExist                       = errors.AddResultCodeInfo(400141, "资源不存在", "Result.ResourceNotExist")
	ProjectTypeNormalError                 = errors.AddResultCodeInfo(400142, "项目不是普通任务", "Result.ProjectTypeNormalError")
	InviteCodeInvalid                      = errors.AddResultCodeInfo(400143, "邀请链接失效", "Result.InviteCodeInvalid")
	UnSupportLoginType                     = errors.AddResultCodeInfo(400144, "不支持的登录方式", "Result.UnSupportLoginType")
	ProjectIsFilingYet                     = errors.AddResultCodeInfo(400145, "项目已归档", "Result.ProjectIsFilingYet")
	LastProjectObjectType                  = errors.AddResultCodeInfo(400146, "最后一个任务清单无法删除", "Result.LastProjectObjectType")
	PasswordEmptyError                     = errors.AddResultCodeInfo(400147, "请输入密码", "Result.PasswordEmptyError")
	PasswordNotSetError                    = errors.AddResultCodeInfo(400148, "密码未设置", "Result.PasswordNotSetError")
	PasswordNotMatchError                  = errors.AddResultCodeInfo(400149, "密码验证错误", "Result.PasswordNotMatchError")
	ParentIssueHasParent                   = errors.AddResultCodeInfo(400150, "子任务不允许创建子任务", "Result.ParentIssueHasParent")
	CreateIssueFail                        = errors.AddResultCodeInfo(400151, "创建任务失败", "Result.CreateIssueFail")
	CommonStatusCannotDelete               = errors.AddResultCodeInfo(400156, "通用状态不能删除", "Result.CommonStatusCannotDelete")
	ProcessStatusHasIssue                  = errors.AddResultCodeInfo(400157, "任务栏中存在任务不可删除，请先将该任务栏中的任务移出", "Result.ProcessStatusHasIssue")
	ProcessStatusNotBind                   = errors.AddResultCodeInfo(400158, "任务栏未绑定该项目", "Result.ProcessStatusNotBind")
	CannotDeleteUniqueStatus               = errors.AddResultCodeInfo(400159, "该任务栏是当前栏目的唯一状态，不能删除", "Result.CannotDeleteUniqueStatus")
	CannotMoveUniqueStatus                 = errors.AddResultCodeInfo(400160, "该任务栏是当前栏目的唯一状态，不能拖动", "Result.CannotMoveUniqueStatus")
	ProcessStatusIsInvaild                 = errors.AddResultCodeInfo(400161, "流程状态不存在", "Result.ProcessStatusIsInvaild")
	PropertyIdNotExist                     = errors.AddResultCodeInfo(400162, "严重程度类型不存在", "Result.PropertyIdNotExist")
	BeforeAfterIssueConflict               = errors.AddResultCodeInfo(400163, "前后置任务不能存在相同任务", "Result.BeforeAfterIssueConflict")
	CannotBeforeAfterSelf                  = errors.AddResultCodeInfo(400164, "不能将自己前置和后置", "Result.CannotBeforeAfterSelf")
	NotNeedAuditIssueNow                   = errors.AddResultCodeInfo(400165, "当前状态不需要确认", "Result.NotNeedAuditIssueNow")
	IssueIsAuditPass                       = errors.AddResultCodeInfo(400166, "当前任务已经被确认", "Result.IssueIsAuditPass")
	NotIssueAuditor                        = errors.AddResultCodeInfo(400167, "您不是当前任务的确认人", "Result.NotIssueAuditor")
	CannotAuditTwice                       = errors.AddResultCodeInfo(400168, "当前审批已处理，请勿重复操作！", "Result.CannotAuditTwice")
	OnlyOwnerCanWithdrawIssue              = errors.AddResultCodeInfo(400169, "只有管理员可以撤回", "Result.OnlyOwnerCanWithdrawIssue")
	NotFinishIssue                         = errors.AddResultCodeInfo(400170, "很抱歉，当前任务尚未完成", "Result.NotFinishIssue")
	UrgeIssueOnlyOwner                     = errors.AddResultCodeInfo(400171, "只有任务负责人可以催办", "Result.UrgeIssueOnlyOwner")
	NeedProjectOwner                       = errors.AddResultCodeInfo(400172, "请选择项目管理员", "Result.NeedProjectOwner")
	NotIssueParticipantWillDeny            = errors.AddResultCodeInfo(400173, "您还不是任务参与者，请联系任务负责人", "Result.NotIssueParticipantWillDeny")
	ForbidInviteSourcePlatform             = errors.AddResultCodeInfo(400174, "该渠道组织不能通过此方式邀请成员。", "Result.ForbidInviteSourcePlatform")
	ProjectPreCodeCannotModify             = errors.AddResultCodeInfo(400175, "项目前缀编号不允许修改", "Result.ProjectPreCodeCannotModify")
	AppNotExist                            = errors.AddResultCodeInfo(400176, "应用不存在", "Result.AppNotExist")
	InviteCodeEmpty                        = errors.AddResultCodeInfo(400177, "邀请码不能为空", "Result.InviteCodeEmpty")
	InviteImportTplGenExcelErr             = errors.AddResultCodeInfo(400178, "生成并保存 excel 文件失败。", "Result.InviteImportTplGenExcelErr")
	NoSupportProjectType                   = errors.AddResultCodeInfo(400179, "不支持的项目类型", "Result.NoSupportProjectType")
	ParamTableIdIsMust                     = errors.AddResultCodeInfo(400180, "请传入正确的表 id", "Result.ParamTableIdIsMust")
	ProjectHasNoTableList                  = errors.AddResultCodeInfo(400181, "项目没有对应的表", "Result.ProjectHasNoTableList")
	TableNotExist                          = errors.AddResultCodeInfo(400182, "表头不存在", "Result.TableNotExist")
	InvalidTableId                         = errors.AddResultCodeInfo(400183, "非法的表ID", "Result.InvalidTableId")
	TablesNotExist                         = errors.AddResultCodeInfo(400184, "表不存在或已被删除", "Result.TablesNotExist")
	TableColumnNotExist                    = errors.AddResultCodeInfo(400185, "表头不存在", "Result.TableColumnNotExist")
	IssueStatusNotExist                    = errors.AddResultCodeInfo(400186, "任务状态不存在", "Result.IssueStatusNotExist")
	BatchUpdateForbidenColumn              = errors.AddResultCodeInfo(400187, "存在批量更新不允许更新的字段", "Result.BatchUpdateForbidenColumn")
	BatchOperateTooManyRows                = errors.AddResultCodeInfo(400188, "批量操作的任务数已超过上限", "Result.BatchOperateTooManyRows")
	DenyDeleteTableWhenAsyncTask           = errors.AddResultCodeInfo(400189, "正在批量导入数据，暂时无法删除数据表", "Result.DenyDeleteTableWhenAsyncTask")
	DenyDeleteProWhenAsyncTask             = errors.AddResultCodeInfo(400190, "正在批量导入数据，暂时无法删除项目", "Result.DenyDeleteProWhenAsyncTask")
	ImportIssueCellInvalid                 = errors.AddResultCodeInfo(400191, "单元格值不合法，请检查。", "Result.ImportIssueCellInvalid")
	ReadExcelFailed                        = errors.AddResultCodeInfo(400192, "读取文件失败，请重新下载导入模板。", "Result.ReadExcelFailed")

	IssueStatusUpdateError                     = errors.AddResultCodeInfo(400202, "任务状态更新失败", "Result.IssueStatusUpdateError")
	UserConfigUpdateError                      = errors.AddResultCodeInfo(400203, "用户设置更新失败", "Result.UserConfigUpdateError")
	UserConfigInsertError                      = errors.AddResultCodeInfo(400204, "用户设置失败", "Result.UserConfigUpdateError")
	IterationIssueRelateError                  = errors.AddResultCodeInfo(400205, "迭代和任务关联失败", "Result.IterationIssueRelateError")
	IterationStatusUpdateError                 = errors.AddResultCodeInfo(400206, "迭代状态更新失败", "Result.IterationStatusUpdateError")
	IssueRelationUpdateError                   = errors.AddResultCodeInfo(400207, "任务关联更新失败", "Result.RelationUpdateError")
	ProjectStatusUpdateError                   = errors.AddResultCodeInfo(400208, "项目状态更新失败", "Result.ProjectStatusUpdateError")
	IssueProjectObjectTypeNotParttenError      = errors.AddResultCodeInfo(400209, "任务项目对象类型不匹配", "Result.IssueProjectObjectTypeNotParttenError")
	IssueOwnerCantBeNull                       = errors.AddResultCodeInfo(400301, "任务负责人不能为空", "Result.IssueOwnerCantBeNull")
	DepartmentNotExist                         = errors.AddResultCodeInfo(400302, "部门不存在", "Result.DepartmentNotExist")
	ParentDepartmentNotExist                   = errors.AddResultCodeInfo(400303, "父部门不存在", "Result.ParentDepartmentNotExist")
	TopDepartmentNotExist                      = errors.AddResultCodeInfo(400304, "顶级部门不存在", "Result.TopDepartmentNotExist")
	ProjectObjectTypeCantBeNullError           = errors.AddResultCodeInfo(400305, "项目对象类型不能为空", "Result.ProjectObjectTypeCantBeNullError")
	PlanEndTimeInvalidError                    = errors.AddResultCodeInfo(400306, "计划结束时间需要大于开始时间", "Result.PlanEndTimeInvalidError")
	OrgNameLenError                            = errors.AddResultCodeInfo(400307, "组织名称为空或超出30个字符", "Result.OrgNameLenError")
	UserNameLenError                           = errors.AddResultCodeInfo(400308, "姓名为空或超出30个字符", "Result.UserNameLenError")
	RepeatTag                                  = errors.AddResultCodeInfo(400309, "标签已存在", "Result.RepeatTag")
	IssueSortReferenceError                    = errors.AddResultCodeInfo(400310, "任务排序参照物不能为空", "Result.IssueSortReferenceError")
	IssueSortReferenceInvalidError             = errors.AddResultCodeInfo(400311, "任务排序参照物无效", "Result.IssueSortReferenceInvalidError")
	DateRangeError                             = errors.AddResultCodeInfo(400312, "时间范围错误", "Result.DateRangeError")
	ImportDataEmpty                            = errors.AddResultCodeInfo(400313, "导入数据为空", "Result.ImportDataEmpty")
	ImportFileNotExist                         = errors.AddResultCodeInfo(400314, "未上传数据文件", "Result.ImportFileNotExist")
	NotDefaultStyle                            = errors.AddResultCodeInfo(400315, "样式无效", "Result.NotDefaultStyle")
	LengthOutOfLimit                           = errors.AddResultCodeInfo(400316, "标签为空或长度超出限制", "Result.LengthOutOfLimit")
	DateParseError                             = errors.AddResultCodeInfo(400317, "时间格式不正确", "Result.DateParseError")
	DailyProjectReportError                    = errors.AddResultCodeInfo(400318, "当日项目已发送", "Result.DailyProjectReportError")
	PageInvalidError                           = errors.AddResultCodeInfo(400319, "页码无效", "Result.PageInvalidError")
	PageSizeInvalidError                       = errors.AddResultCodeInfo(400320, "页长无效", "Result.PageSizeInvalidError")
	UserOrgNotRelation                         = errors.AddResultCodeInfo(400321, "用户不是该组织成员", "Result.UserOrgNotRelation")
	UserDisabledError                          = errors.AddResultCodeInfo(400322, "已经被组织禁用", "Result.UserDisabledError")
	InvalidImportFile                          = errors.AddResultCodeInfo(400323, "文件格式有误，请上传xls、xlsx格式的文件", "Result.InvalidImportFile")
	FileParseFail                              = errors.AddResultCodeInfo(400324, "文件解析失败,请下载最新文件模板或检查文件内容", "Result.FileParseFail")
	TooLargeImportData                         = errors.AddResultCodeInfo(400325, "导入任务数据过大", "Result.TooLargeImportData")
	TooLongProjectRemark                       = errors.AddResultCodeInfo(400326, "项目简介应少于500字", "Result.TooLongProjectRemark")
	ProjectCodeLenError                        = errors.AddResultCodeInfo(400327, "项目编号长度不得超过64个字", "Result.ProjectCodeLenError")
	ProjectNameLenError                        = errors.AddResultCodeInfo(400328, "项目名称长度不得超过256个字", "Result.ProjectNameLenError")
	ProjectPreCodeLenError                     = errors.AddResultCodeInfo(400329, "项目前缀编号长度不得超过16个字", "Result.ProjectPreCodeLenError")
	ProjectRemarkLenError                      = errors.AddResultCodeInfo(400330, "项目描述长度不得超过512个字", "Result.ProjectRemarkLenError")
	ProjectIsArchivedWhenModifyIssue           = errors.AddResultCodeInfo(400331, "不允许操作归档项目下的任务", "Result.ProjectIsArchivedWhenModifyIssue")
	NoPrivateProjectPermissions                = errors.AddResultCodeInfo(400332, "没有私有项目操作权限", "Result.NoPrivateProjectPermissions")
	ChildIssueForFirst                         = errors.AddResultCodeInfo(400333, "第一条任务不能是子任务", "Result.ChildIssueForFirst")
	ProjectObjectTypeSameName                  = errors.AddResultCodeInfo(400334, "任务清单名字重复", "Result.ProjectObjectTypeSameName")
	ProjectNameEmpty                           = errors.AddResultCodeInfo(400335, "项目名称不能为空", "Result.ProjectNameEmpty")
	UpdateMemberIdsIsEmptyError                = errors.AddResultCodeInfo(400336, "变动的成员列表为空", "Result.UpdateMemberIdsIsEmptyError")
	UpdateMemberStatusFail                     = errors.AddResultCodeInfo(400337, "修改成员状态失败", "Result.UpdateMemberStatusFail")
	CantUpdateStatusWhenParentIssueIsCompleted = errors.AddResultCodeInfo(400338, "父任务已完成，无法修改子任务状态", "Result.CantUpdateStatusWhenParentIssueIsCompleted")
	RoleNameLenErr                             = errors.AddResultCodeInfo(400339, "角色名包含非法字符或长度超出10个字符", "Result.RoleNameLenErr")
	DefaultRoleCantModify                      = errors.AddResultCodeInfo(400340, "默认角色不允许编辑", "Result.DefaultRoleCantModify")
	RoleModifyBusy                             = errors.AddResultCodeInfo(400341, "角色更新繁忙", "Result.RoleEditBusy")
	RoleNameRepeatErr                          = errors.AddResultCodeInfo(400342, "角色名称重复", "Result.RoleNameRepeatErr")
	CannotRemoveProjectOwner                   = errors.AddResultCodeInfo(400343, "项目管理员不能被移除", "Result.CannotRemoveProjectOwner")
	DefaultRoleNameErr                         = errors.AddResultCodeInfo(400344, "与系统角色名称冲突", "Result.DefaultRoleNameErr")
	SourceChannelNotDefinedError               = errors.AddResultCodeInfo(400355, "来源通道未定义", "Result.SourceChannelNotDefinedError")
	OrgNotNeedInitError                        = errors.AddResultCodeInfo(400356, "组织已存在，不需要初始化", "Result.OrgNotNeedInitError")
	IssueCommentLenError                       = errors.AddResultCodeInfo(400357, "评论不得为空且不能超过2000字", "Result.IssueCommentLenError")
	IssueRemarkLenError                        = errors.AddResultCodeInfo(400358, "描述不能超过10000字", "Result.IssueRemarkLenError")
	AuthCodeIsNull                             = errors.AddResultCodeInfo(400359, "验证码不得为空", "Result.AuthCodeIsNull")
	ContactRemarkLenErr                        = errors.AddResultCodeInfo(400360, "问题反馈描述不得超过512字", "Result.ContactRemarkLenErr")
	ContactResourceInfoLenErr                  = errors.AddResultCodeInfo(400361, "问题反馈资源信息不得超过2048字", "Result.ContactResourceInfoLenErr")
	ContactResourceSizeErr                     = errors.AddResultCodeInfo(400362, "问题反馈图片数量不能超过5个", "Result.ContactResourceSizeErr")
	PwdAlreadySettingsErr                      = errors.AddResultCodeInfo(400363, "密码已设置过", "Result.PwdAlreadySettingsErr")
	PwdFormatError                             = errors.AddResultCodeInfo(400364, "密码需要以字母开头，长度在6~18之间，只能包含字母、数字和下划线", "Result.PwdLengthError")
	TagNotExist                                = errors.AddResultCodeInfo(400365, "标签不存在", "Result.TagNotExist")
	InvalidProjectNameError                    = errors.AddResultCodeInfo(400366, "项目名不能超出20个字符", "Result.InvalidProjectNameError")
	InvalidProjectPreCodeError                 = errors.AddResultCodeInfo(400367, "项目前缀编号不能为空且最多输入10个字符,包含数字和英文", "Result.InvalidProjectPreCodeError")
	InvalidProjectRemarkError                  = errors.AddResultCodeInfo(400368, "项目简介不能超出500个字符", "Result.InvalidProjectRemarkError")
	IssueRelateTagFail                         = errors.AddResultCodeInfo(400369, "任务关联标签失败", "Result.IssueRelateTagFail")
	CreateTagFail                              = errors.AddResultCodeInfo(400370, "创建标签失败", "Result.CreateTagFail")
	ProjectNoticeLenError                      = errors.AddResultCodeInfo(400371, "项目公告不能超出2000字", "Result.ProjectNoticeLenError")
	AlreadyBindChat                            = errors.AddResultCodeInfo(400372, "项目已绑定该群聊", "Result.AlreadyBindChat")
	CannotBindChat                             = errors.AddResultCodeInfo(400373, "该功能暂只支持飞书用户", "Result.CannotBindChat")
	NotBindChatYet                             = errors.AddResultCodeInfo(400374, "项目尚未绑定该群聊", "Result.NotBindChatYet")
	CannotDisbandMainChat                      = errors.AddResultCodeInfo(400375, "您无权解绑项目主群聊，请联系项目管理员", "Result.NotBindChatYet")
	DeleteProjectErr                           = errors.AddResultCodeInfo(400376, "删除失败", "Result.DeleteProjectErr")
	CreateIterationErr                         = errors.AddResultCodeInfo(400377, "迭代创建失败", "Result.CreateIterationErr")
	CreateIterationRelationErr                 = errors.AddResultCodeInfo(400378, "关联迭代状态失败", "Result.CreateIterationRelationErr")
	OrgConfigNotExist                          = errors.AddResultCodeInfo(400379, "组织配置不存在", "Result.OrgConfigNotExist")
	DepartmentNameInvalid                      = errors.AddResultCodeInfo(400381, "部门名称应为1~20个字符", "Result.DepartmentNameInvalid")
	DepartmentExistAlready                     = errors.AddResultCodeInfo(400382, "部门已存在", "Result.DepartmentExistAlready")
	DenyUpdateIssueWorkHours                   = errors.AddResultCodeInfo(400383, "抱歉，您无权更改工时记录。", "Result.DenyUpdateIssueWorkHours")
	DenyEnableFuncIssueWorkHours               = errors.AddResultCodeInfo(400384, "抱歉，您无权启用/关闭当前项目的工时功能。", "Result.DenyEnableFuncIssueWorkHours")
	DenyCreateIssueWorkHours                   = errors.AddResultCodeInfo(400385, "抱歉，您无权增加工时记录。", "Result.DenyCreateIssueWorkHours")
	SimplePredictIssueWorkHourExist            = errors.AddResultCodeInfo(400386, "任务的预估工时记录已存在。", "Result.SimplePredictIssueWorkHourExist")
	DetailPredictIssueWorkHourNeedWorker       = errors.AddResultCodeInfo(400387, "详细预估工时需要执行人。", "Result.DetailPredictIssueWorkHourNeedWorker")
	IssueWorkHourNotExist                      = errors.AddResultCodeInfo(400388, "工时记录不存在。", "Result.IssueWorkHourNotExist")
	IssueRelationExist                         = errors.AddResultCodeInfo(400389, "任务关联实体已存在。", "Result.IssueRelationExist")
	ProjectDisabledWorkHour                    = errors.AddResultCodeInfo(400390, "项目没有开启工时功能。", "Result.ProjectDisabledWorkHour")
	WorkHourHasSubRecordDisableDel             = errors.AddResultCodeInfo(400391, "该任务下还有工时日志或子预估工时，不能删除总预估工时。", "Result.WorkHourHasSubRecordDisableDel")
	WorkHourTimeRangeForNeedTimeInvalid        = errors.AddResultCodeInfo(400392, "起止时间内的工时不合法。", "Result.WorkHourTimeRangeForNeedTimeInvalid")
	WorkHourMaxNeedTime                        = errors.AddResultCodeInfo(400393, "您输入的工时太多了哦。不允许超过 100 万小时。", "Result.WorkHourMaxNeedTime")
	WorkHourNeedTimeInvalid                    = errors.AddResultCodeInfo(400394, "请输入合法的工时。", "Result.WorkHourNeedTimeInvalid")
	ConditionInvalid                           = errors.AddResultCodeInfo(400395, "条件限制不合法。", "Result.ConditionInvalid")
	TimeRangeLimitHalfYearInvalid              = errors.AddResultCodeInfo(400396, "查询的时间范围请不要超过半年。", "Result.TimeRangeLimitHalfYearInvalid")
	IssueColumnIsEmpty                         = errors.AddResultCodeInfo(400397, "任务的字段为空。", "Result.IssueColumnIsEmpty")
	AddUserToProButNoPower                     = errors.AddResultCodeInfo(400398, "将执行人加入此项目中，但您无权为项目新增成员，请联系管理员。", "Result.AddUserToProButNoPower")
	DenyUpdateWorkHourForIssueHashNoPro        = errors.AddResultCodeInfo(400399, "抱歉，不能给没有项目归属的任务更新工时。", "Result.DenyUpdateWorkHourForIssueHashNoPro")
	IssueHasNoBelongPro                        = errors.AddResultCodeInfo(400400, "抱歉，这个任务没有项目归属。", "Result.IssueHasNoBelongPro")
	InvalidSex                                 = errors.AddResultCodeInfo(400401, "性别不在正常范围内", "Result.InvalidSex")
	IssueTitleError                            = errors.AddResultCodeInfo(400402, "任务标题包含非法字符或超出500个字符", "Result.IssueTitleError")
	FolderIdNotExistError                      = errors.AddResultCodeInfo(400403, "文件夹不存在", "Result.FolderIdNotExistError")
	InvalidResourceNameError                   = errors.AddResultCodeInfo(400404, "文件名包含非法字符或超出300个字符", "Result.InvalidResourceNameError")
	InvalidFolderNameError                     = errors.AddResultCodeInfo(400405, "文件夹名包含非法字符或超出30个字符", "Result.InvalidResourceNameError")
	InvalidFolderIdsError                      = errors.AddResultCodeInfo(400406, "无效的文件夹ids", "Result.InvalidFolderIdsError")
	InvalidResourceIdsError                    = errors.AddResultCodeInfo(400407, "无效的文件ids", "Result.InvalidResourceIdsError")
	ParentIdIsItselfError                      = errors.AddResultCodeInfo(400409, "目标文件夹是自己本身,无需移动", "Result.ParentIdIsItselfError")
	ResouceNotInFolderError                    = errors.AddResultCodeInfo(400410, "文件不在该文件夹下", "Result.ResouceNotInFolderError")
	ReourceTypeMismatchType                    = errors.AddResultCodeInfo(400411, "文件类型不匹配", "Result.ReourceTypeMismatchType")
	EncodeNotSupport                           = errors.AddResultCodeInfo(400412, "不支持的编码类型", "Result.EncodeNotSupport")
	SetUserPasswordError                       = errors.AddResultCodeInfo(400413, "设置密码失败", "Result.SetUserPasswordError")
	UnBindLoginNameFail                        = errors.AddResultCodeInfo(400414, "解绑登录方式失败", "Result.UnBindLoginNameFail")
	BindLoginNameFail                          = errors.AddResultCodeInfo(400415, "绑定登录方式失败", "Result.BindLoginNameFail")
	NotBindAccountError                        = errors.AddResultCodeInfo(400416, "手机号未注册，请先注册或填写其他账号", "Result.NotBindAccountError")
	AccountAlreadyBindError                    = errors.AddResultCodeInfo(400417, "该登录方式已绑定其它账号", "Result.AccountAlreadyBindError")
	EmailNotBindAccountError                   = errors.AddResultCodeInfo(400418, "该邮箱未绑定任何账户，请重新输入或使用手机验证码登录", "Result.EmailNotBindAccountError")
	MobileNotBindAccountError                  = errors.AddResultCodeInfo(400419, "该手机号未绑定任何账号", "Result.MobileNotBindAccountError")
	DisbandThirdAccountError                   = errors.AddResultCodeInfo(400420, "当前组织已经绑定第三方平台，不允许解绑", "Result.DisbandThirdAccountError")
	NotDisbandCurrentSourceChannel             = errors.AddResultCodeInfo(400421, "账号未绑定当前平台", "Result.NotDisbandCurrentSourceChannel")
	AccountNotBelongToCurrentFs                = errors.AddResultCodeInfo(400422, "该手机号不属于当前飞书组织用户", "Result.AccountNotBelongToCurrentFs")
	DeptAndUserNoSameName                      = errors.AddResultCodeInfo(400423, "不存在同名部门和同名用户", "Result.DeptAndUserNoSameName")
	MemberDuplicateWhenImport                  = errors.AddResultCodeInfo(400424, "成员名字重名了，请下载重名 excel，并使用 excel 中推荐的值填写", "Result.MemberDuplicateWhenImport")
	DeptDuplicateWhenImport                    = errors.AddResultCodeInfo(400425, "部门名字重名了，请下载重名 excel，并使用 excel 中推荐的值填写", "Result.DeptDuplicateWhenImport")
	TableSameName                              = errors.AddResultCodeInfo(400426, "表格名字重复", "Result.TableSameName")
	ImportExcelNotMatchedWithTable             = errors.AddResultCodeInfo(400427, "导入文件列与视图列数量不一致，请检查后重新导入。", "Result.ImportExcelNotMatchedWithTable")
	DeleteAttachmentError                      = errors.AddResultCodeInfo(400428, "没有权限删除附件", "Result.DeleteAttachmentError")
	UrgeOwnersNoPermission                     = errors.AddResultCodeInfo(400429, "项目管理员才能催办负责人", "Result.UrgeOwnersNoPermission")
	UrgeAuditorsNoPermission                   = errors.AddResultCodeInfo(400430, "项目管理员、负责人才能催办确认人", "Result.UrgeAuditorsNoPermission")
	ImportIssueFailed                          = errors.AddResultCodeInfo(400431, "导入任务失败", "Result.ImportIssueFailed")
	AsyncTaskNotExist                          = errors.AddResultCodeInfo(400432, "后台任务不存在", "Result.AsyncTaskNotExist")
	RecoverAttachmentError                     = errors.AddResultCodeInfo(400433, "该附件对应的记录已被删除，无法恢复该附件，请先恢复对应的记录再恢复该附件", "Result.RecoverAttachmentError")
	RecoverResourceFailed                      = errors.AddResultCodeInfo(400434, "恢复资源失败", "Result.RecoverResourceFailed")
	RecoverDocumentFailedWithNoFolder          = errors.AddResultCodeInfo(400435, "原路径不可用，无法恢复", "Result.RecoverDocumentFailedWithNoFolder")
	ImportAsyncTaskIsExecuting                 = errors.AddResultCodeInfo(400436, "该表正在导入，导入完成才能执行下一次导入", "Result.ImportAsyncTaskIsExecuting")
	ProChatNotInChat                           = errors.AddResultCodeInfo(400437, "你还不在项目群内，如需进群，请联系管理员", "Result.ProChatNotInChat")
	HasNotProChat                              = errors.AddResultCodeInfo(400438, "该项目没有群聊，请去“项目设置”开启", "Result.HasNotProChat")
	OrgNotSupportProChat                       = errors.AddResultCodeInfo(400439, "当前组织暂不支持群聊功能", "Result.OrgNotSupportProChat")
	HasNoMatchedColumn                         = errors.AddResultCodeInfo(400440, "没有匹配到列，请重新下载导入模板进行导入。", "Result.HasNoMatchedColumn")
	AccountNameLenError                        = errors.AddResultCodeInfo(400441, "账号名30字符内", "Result.AccountNameLenError")
	AccountAllReadyExist                       = errors.AddResultCodeInfo(400442, "账号名已存在", "Result.AccountAllReadyExist")
	AccountHadBindWeiXin                       = errors.AddResultCodeInfo(400443, "该手机号已经绑定过微信，请先解绑", "Result.AccountHadBindWeiXin")
	ShareViewNotExist                          = errors.AddResultCodeInfo(400444, "分享的视图不存在", "Result.ShareViewNotExist")
	ShareViewPasswordError                     = errors.AddResultCodeInfo(400445, "密码错误", "Result.ShareViewPasswordError")
	AccountNotBindWeiXin                       = errors.AddResultCodeInfo(400446, "没有绑定过手机号", "Result.AccountNotBindWeiXin")
	WeiXinAlreadyBindError                     = errors.AddResultCodeInfo(400447, "该微信已绑定其它账号", "Result.WeiXinAlreadyBindError")

	// User
	UserInitError                       = errors.AddResultCodeInfoWithSentry(400501, "用户初始化失败", "Result.UserInitError")
	UserNotInitError                    = errors.AddResultCodeInfoWithSentry(400502, "用户未初始化", "Result.UserNotInitError")
	UserNotExist                        = errors.AddResultCodeInfo(400503, "用户不存在", "Result.UserNotExist")
	UserInfoGetFail                     = errors.AddResultCodeInfo(400504, "用户信息获取失败", "Result.UserInfoGetFail")
	UserRegisterError                   = errors.AddResultCodeInfo(400505, "用户注册失败", "Result.UserRegisterError")
	LarkInitError                       = errors.AddResultCodeInfo(400506, "示例数据已初始化", "Result.LarkInitError")
	UserSexFail                         = errors.AddResultCodeInfo(400507, "用户性别错误", "Result.UserSexFail")
	UserNameEmpty                       = errors.AddResultCodeInfo(400508, "用户姓名不能为空串", "Result.UserNameEmpty")
	EmailNotRegisterError               = errors.AddResultCodeInfo(400509, "当前邮箱未注册", "Result.EmailNotRegisterError")
	EmailNotBindError                   = errors.AddResultCodeInfo(400510, "邮箱未绑定", "Result.EmailNotBindError")
	MobileNotBindError                  = errors.AddResultCodeInfo(400511, "手机号未绑定", "Result.MobileNotBindError")
	EmailAlreadyBindError               = errors.AddResultCodeInfo(400512, "邮箱已绑定, 请先解绑", "Result.EmailAlreadyBindError")
	MobileAlreadyBindError              = errors.AddResultCodeInfo(400513, "手机号已绑定， 请先解绑", "Result.MobileAlreadyBindError")
	EmailAlreadyBindByOtherAccountError = errors.AddResultCodeInfo(400514, "该邮箱已被其他账户绑定", "Result.EmailAlreadyBindByOtherAccountError")
	MobileAlreadyBindOtherAccountError  = errors.AddResultCodeInfo(400515, "该手机号已被其他账户绑定", "Result.MobileAlreadyBindOtherAccountError")
	AccountNotRegister                  = errors.AddResultCodeInfo(400516, "手机号未注册，请先注册或填写其他账号", "Result.AccountNotRegister")
	InputParamEmpty                     = errors.AddResultCodeInfo(400517, "入参为空", "Result.InputParamEmpty")
	SyncUserHasNoUserUnderPermission    = errors.AddResultCodeInfoWithSentry(400518, "同步用户数据时，授权范围内无用户", "Result.SyncUserHasNoUserUnderPermission")
	SyncUserUpdateUserFail              = errors.AddResultCodeInfoWithSentry(400519, "同步用户数据时，更新用户失败", "Result.SyncUserUpdateUserFail")
	SyncDepartmentError                 = errors.AddResultCodeInfoWithSentry(400520, "同步部门信息时，异常", "Result.SyncDepartmentError")
	OperatorInvalid                     = errors.AddResultCodeInfoWithSentry(400521, "操作人无效", "Result.OperatorInvalid")
	ImportUserNumTooMany                = errors.AddResultCodeInfoWithSentry(400522, "抱歉，导入成员数量必须小于 200。", "Result.ImportUserNumTooMany")
	NotSupportIterStatus                = errors.AddResultCodeInfoWithSentry(400523, "该迭代状态不支持。", "Result.ErrIterStatus")
	IterationInitStatusIdsError         = errors.AddResultCodeInfo(400524, "迭代初始化状态ids错误", "Result.IterationInitStatusIdsError")
	IterationInitStatusIdsNotExist      = errors.AddResultCodeInfo(400525, "迭代初始化状态ids不存在", "Result.IterationInitStatusIdNotExist")
	UserOrgConflict                     = errors.AddResultCodeInfo(400526, "用户组织冲突，同一个用户不能加入同一个组织", "Result.UserOrgConflict")
	MobileSameError                     = errors.AddResultCodeInfo(400527, "与已有手机号重复，请重新填写", "Result.MobileSamError")
	MobileInvalidError                  = errors.AddResultCodeInfo(400528, "手机号码格式不正确", "Result.MobileSamError")
	ManageGroupNotExist                 = AddResultCodeInfo(400529, "管理组不存在或已删除", "Result.ManageGroupNotExist")
	DenyChangeSysAdminGroupOfUser       = AddResultCodeInfo(400530, "不能更改组织超级管理员的管理组", "Result.DenyChangeSysAdminGroupOfUser")
	CannotChangeSelfStatus              = AddResultCodeInfo(400531, "不允许变更自己的状态", "Result.CannotChangeSelfStatus")
	CannotEditSuperAdminInfo            = AddResultCodeInfo(400532, "非超管无法修改超管的个人信息", "Result.CannotEditSuperAdminInfo")

	// 动态
	TrendsCreateError            = errors.AddResultCodeInfoWithSentry(401001, "动态创建失败", "Result.TrendsCreateError")
	TrendsObjTypeNilError        = errors.AddResultCodeInfo(401002, "对象id有值的情况下对象类型不能为空", "Result.TrendsObjTypeNilError")
	TrendsObjIdNilError          = errors.AddResultCodeInfo(401003, "对象类型有值的情况下对象id不能为空", "Result.TrendsObjIdNilError")
	IssueCommentImageLimitsError = errors.AddResultCodeInfo(401004, "图片数量限制为9张", "Result.IssueCommentAttachmentsLimitsError")

	// 项目对象类型不存在
	ProjectObjectTypeNotExist             = errors.AddResultCodeInfo(402001, "任务栏不存在", "Result.ProjectObjectTypeNotExist")
	ProjectTypeProjectObjectTypeNotExist  = errors.AddResultCodeInfo(402002, "项目类型与项目对象类型关联不存在", "Result.ProjectTypeProjectObjectTypeNotExist")
	ProjectObjectTypeDeleteFailExistIssue = errors.AddResultCodeInfo(402003, "任务栏中存在任务不可删除，请先将该任务栏中的任务移出", "Result.ProjectTypeDeleteFailExistIssue")
	InvalidProjectObjectTypeName          = errors.AddResultCodeInfo(402004, "任务栏名称不能为空且不能超过30字", "Result.InvalidProjectObjectTypeName")
	CannotMoveChildIssue                  = errors.AddResultCodeInfo(402005, "子任务不可单独移动任务栏", "Result.CannotMoveChildIssue")
	MoveIssueFail                         = errors.AddResultCodeInfo(402006, "移动记录失败", "Result.MoveIssueFail")
	IssueLevelOutLimit                    = errors.AddResultCodeInfo(402007, "转化后目标任务超过九级", "Result.IssueLevelOutLimit")
	IssuePlanWorkHourNegativeError        = errors.AddResultCodeInfo(402009, "任务工时不允许为负数", "Result.IssuePlanWorkHourNegativeError")
	ParentIsChildIssue                    = errors.AddResultCodeInfo(402010, "父任务属于当前任务的子任务", "Result.ParentIsChild")
	RemainIssuesInStatus                  = errors.AddResultCodeInfo(402011, "删除的状态中仍存在任务，请移除相应的任务再进行操作", "Result.RemainIssuesInStatus")
	ProjectMenuConfigNotExist             = errors.AddResultCodeInfo(402012, "项目的菜单配置不存在", "Result.ProjectMenuConfigNotExist")
	InvalidProjectTableName               = errors.AddResultCodeInfo(402013, "数据表名称长度上限200字", "Result.InvalidProjectTableName")
	InvalidUserIdsError                   = errors.AddResultCodeInfo(402014, "用户id数据参数异常", "Result.InvalidUserIdsError")
	UpdateResourceFolderError             = errors.AddResultCodeInfo(402015, "文件已变化，请刷新后重试！", "Result.UpdateResourceFolderError")

	// domain
	ProjectDomainError    = errors.AddResultCodeInfo(405001, "项目领域出错", "Result.ProjectDomainError")
	IssueDomainError      = errors.AddResultCodeInfo(405002, "任务领域出错", "Result.IssueDomainError")
	UserDomainError       = errors.AddResultCodeInfo(405003, "用户领域出错", "Result.UserDomainError")
	BaseDomainError       = errors.AddResultCodeInfo(405004, "领域出错", "Result.BaseDomainError")
	TrendDomainError      = errors.AddResultCodeInfo(405005, "动态领域出错", "Result.TrendDomainError")
	IterationDomainError  = errors.AddResultCodeInfo(405006, "迭代领域出错", "Result.IterationDomainError")
	ObjectTypeDomainError = errors.AddResultCodeInfo(405007, "对象类型领域出错", "Result.ObjectTypeDomainError")
	ResourceDomainError   = errors.AddResultCodeInfo(405008, "资源领域出错", "Result.ResourceDomainError")
	ProcessDomainError    = errors.AddResultCodeInfo(405009, "流程领域出错", "Result.ProcessDomainError")
	DepartmentDomainError = errors.AddResultCodeInfo(405010, "部门领域出错", "Result.DepartmentDomainError")
	FormDomainError       = errors.AddResultCodeInfo(405011, "表单服务出错", "Result.FormDomainError")
	DatacenterDomainError = errors.AddResultCodeInfo(405012, "数据中心服务出错", "Result.DatacenterDomainError")
	TableDomainError      = errors.AddResultCodeInfo(405013, "table服务出错", "Result.TableDomainError")
	DefaultFieldError     = errors.AddResultCodeInfo(405014, "默认字段不允许删除", "Result.DefaultFieldError")

	// 权限验证领域
	IllegalityRoleOperation                       = errors.AddResultCodeInfo(407001, "非法的操作code", "Result.IllegalityRoleOperation")
	UserRoleNotDefinition                         = errors.AddResultCodeInfo(407002, "用户角色未定义", "Result.UserRoleNotDefinition")
	NoOperationPermissions                        = errors.AddResultCodeInfo(407003, "你没有创建项目的权限，如需要请联系管理员", "Result.NoOperationPermissions")
	PermissionNotExist                            = errors.AddResultCodeInfo(407004, "权限项不存在", "Result.PermissionNotExist")
	NoOperationPermissionForProject               = errors.AddResultCodeInfo(407005, "暂无权限操作，请联系项目管理员", "Result.NoOperationPermissionForProject")
	NoOperationPermissionForIssue                 = errors.AddResultCodeInfo(407006, "暂无权限操作，请联系任务负责人", "Result.NoOperationPermissionForIssue")
	NoOperationPermissionForIssueNotProjectMember = errors.AddResultCodeInfo(407007, "你不是项目成员，请联系项目管理员", "Result.NoOperationPermissionForIssueNotProjectMember")
	NoOperationPermissionForIssueUpdate           = errors.AddResultCodeInfo(407008, "你没有相关字段编辑权限，如需要请联系管理员。", "Result.NoOperationPermissionForIssueUpdate")
	NoOperationPermissionForOrgColumn             = errors.AddResultCodeInfo(407009, "你没有组织字段编辑权限，如需要请联系管理员。", "Result.NoOperationPermissionForOrgColumn")
	TodoFillInMissRequiredColumn                  = errors.AddResultCodeInfo(407010, "填写内容缺少必填字段", "Result.TodoFillInMissRequiredColumn")
	TodoFillInChangeReadOnlyColumn                = errors.AddResultCodeInfo(407011, "填写内容不允许修改只读字段", "Result.TodoFillInChangeReadOnlyColumn")

	// 待办相关
	TodoInvalidParameter = errors.AddResultCodeInfo(408000, "待办配置参数异常", "Result.TodoInvalidParameter")
	TodoInvalidOperator  = errors.AddResultCodeInfo(408001, "待办操作人非法", "Result.TodoInvalidOperator")
	TodoInvalidOp        = errors.AddResultCodeInfo(408003, "待办处理操作异常", "Result.TodoInvalidOp")
	TodoIsDone           = errors.AddResultCodeInfo(408004, "待办已完成或撤销，无法重复处理", "Result.TodoIsDone")

	// LessCode fuse 无码融合
	LcOrgInitError                    = errors.AddResultCodeInfo(5901204, "组织没有合理初始化，缺少汇总表 id。", "Result.LcOrgInitError")
	LcUpdateAppPermissionGroupOptAuth = errors.AddResultCodeInfo(5901205, "更新项目对应的应用权限组 optAuth 配置失败。", "Result.LcUpdateAppPermissionGroupOptAuth")
	LcAppIdInvalid                    = errors.AddResultCodeInfo(5901206, "应用的 appId 不合法。", "Result.LcAppIdInvalid")
	LcCanNotDeleteByUse               = errors.AddResultCodeInfo(451026, "有表单使用该团队字段", "Result.LcCanNotDeleteByUse")

	// dingtalk open api error
	SuiteTicketNotExistError      = errors.AddResultCodeInfoWithSentry(600001, "suiteTicket失效或不存在", "ResultCode.SuiteTicketNotExistError")
	DingTalkOpenApiCallError      = errors.AddResultCodeInfoWithSentry(600002, "钉钉OpenApi调用异常", "ResultCode.DingTalkOpenApiCallError")
	DingTalkAvoidCodeInvalidError = errors.AddResultCodeInfoWithSentry(600003, "钉钉免登code失效", "ResultCode.DingTalkAvoidCodeInvalidError")
	DingTalkClientError           = errors.AddResultCodeInfoWithSentry(600004, "钉钉Client获取失败", "ResultCode.DingTalkClientError")
	DingTalkGetUserInfoError      = errors.AddResultCodeInfoWithSentry(600005, "钉钉获取用户信息失败", "ResultCode.DingTalkGetUserInfoError")
	DingTalkOrgInitError          = errors.AddResultCodeInfoWithSentry(600006, "钉钉企业初始化失败", "ResultCode.DingTalkOrgInitError")
	DingTalkConfigError           = errors.AddResultCodeInfoWithSentry(600007, "钉钉配置错误", "ResultCode.DingTalkConfigError")
	DingTalkFinishOrderError      = errors.AddResultCodeInfoWithSentry(600008, "通知钉钉处理钉钉完成失败", "ResultCode.DingTalkFinishOrderError")
	DingTalkLogBotConfigError     = errors.AddResultCodeInfoWithSentry(600009, "钉钉日志告警配置异常。", "ResultCode.DingTalkLogBotConfigError")

	PlatFormOpenApiCallError     = errors.AddResultCodeInfoWithSentry(600010, "第三方平台OpenApi调用异常", "ResultCode.ThirdPlatformOpenApiCallError")
	PlatFormAppUnauthorizedError = errors.AddResultCodeInfo(600011, "该第三方组织未注册极星协作应用，请在第三方平台应用内完成注册并授权", "Result.PlatFormAppUnauthorizedError")

	// Login Error
	SMSLoginCodeSendError                = errors.AddResultCodeInfo(601001, "登录验证码发送失败", "ResultCode.SMSLoginCodeSendError")
	SMSPhoneNumberFormatError            = errors.AddResultCodeInfo(601002, "手机号格式错误，请重新输入", "ResultCode.SMSPhoneNumberFormatError")
	SMSSendLimitError                    = errors.AddResultCodeInfo(601003, "发送过于频繁（服务商）", "ResultCode.SMSSendLimitError")
	SMSSendTimeLimitError                = errors.AddResultCodeInfo(601004, "发送过于频繁", "ResultCode.SMSSendTimeLimitError")
	SMSLoginCodeInvalid                  = errors.AddResultCodeInfo(601005, "验证码已失效，请重新获取", "ResultCode.SMSLoginCodeInvalid")
	SMSLoginCodeNotMatch                 = errors.AddResultCodeInfo(601006, "验证码错误，请重新获取", "ResultCode.SMSLoginCodeNotMatch")
	SMSLoginCodeVerifyFailTimesOverLimit = errors.AddResultCodeInfo(601007, "验证码错误，失败次数过多，请重新发送", "ResultCode.SMSLoginCodeVerifyFailTimesOverLimit")
	PwdLoginCodeNotMatch                 = errors.AddResultCodeInfo(601008, "图形验证码错误", "ResultCode.PwdLoginCodeNotMatch")
	PwdLoginUsrOrPwdNotMatch             = errors.AddResultCodeInfo(601009, "用户名或密码错误", "ResultCode.PwdLoginUsrOrPwdNotMatch")
	ChangeLoginNameInvalid               = errors.AddResultCodeInfo(601010, "换绑操作已过期，请重新进行换绑操作", "ResultCode.ChangeLoginNameInvalid")
	PwdLoginLimitError                   = errors.AddResultCodeInfo(601011, "登录过于频繁，请一小时后再试", "ResultCode.PwdLoginLimitError")
	PhoneIsRegisterError                 = errors.AddResultCodeInfo(601012, "手机号已经注册，请登录", "ResultCode.PhoneIsRegisterError")

	EmailFormatErr       = errors.AddResultCodeInfo(602001, "邮箱格式错误", "ResultCode.EmailFormatErr")
	EmailSubjectEmptyErr = errors.AddResultCodeInfo(602002, "邮箱标题不能为空", "ResultCode.EmailSubjectEmptyErr")
	EmailSendErr         = errors.AddResultCodeInfo(602003, "邮件发送失败", "ResultCode.EmailSendErr")

	NotSupportedContactAddressType = errors.AddResultCodeInfo(603001, "不支持的联系方式类型", "ResultCode.NotSupportedContactAddressType")
	NotSupportedAuthCodeType       = errors.AddResultCodeInfo(603002, "不支持的验证码类型", "ResultCode.NotSupportedAuthCodeType")
	NotSupportedRegisterType       = errors.AddResultCodeInfo(603003, "暂时不支持该注册方式", "ResultCode.NotSupportedRegisterType")
	HaveNoContract                 = errors.AddResultCodeInfo(603004, "请保证至少一种联系方式", "ResultCode.HaveNoContract")
	CanNotBindSameMobile           = errors.AddResultCodeInfo(603005, "不能绑定相同的手机号", "ResultCode.CanNotBindSameMobile")
	SameUserInOrg                  = errors.AddResultCodeInfo(603006, "相同用户在同一个组织内，不允许绑定", "ResultCode.SameUserInOrg")

	// 飞书 open api err
	FeiShuOpenApiCallError                = errors.AddResultCodeInfoWithSentry(606001, "飞书OpenApi调用异常", "ResultCode.FeiShuOpenApiCallError")
	FeiShuAppTicketNotExistError          = errors.AddResultCodeInfoWithSentry(606002, "AppTicket不存在", "ResultCode.FeiShuAppTicketNotExistError")
	FeiShuConfigNotExistError             = errors.AddResultCodeInfoWithSentry(606003, "飞书配置不存在", "ResultCode.FeiShuConfigNotExistError")
	FeiShuClientTenantError               = errors.AddResultCodeInfoWithSentry(606004, "飞书客户端获取失败", "ResultCode.FeiShuClientTenantError")
	FeiShuGetAppAccessTokenError          = errors.AddResultCodeInfoWithSentry(606005, "飞书获取AppAccessToken失败", "ResultCode.FeiShuGetAppAccessTokenError")
	FeiShuGetTenantAccessTokenError       = errors.AddResultCodeInfoWithSentry(606006, "飞书获取TenantAccessToken失败", "ResultCode.FeiShuGetTenantAccessTokenError")
	FeiShuAuthCodeInvalid                 = errors.AddResultCodeInfoWithSentry(606007, "飞书用户授权失败", "ResultCode.FeiShuAuthCodeInvalid")
	FeiShuCardCallSignVerifyError         = errors.AddResultCodeInfoWithSentry(606008, "飞书卡片回调签名校验失败", "ResultCode.FeiShuCardCallSigVerifyError")
	FeiShuCardCallMsgRepetError           = errors.AddResultCodeInfoWithSentry(606009, "飞书卡片消息重复推送", "ResultCode.FeiShuCardCallMsgRepetError")
	FeiShuUserNotInAppUseScopeOfAuthority = errors.AddResultCodeInfoWithSentry(606010, "抱歉，您当前无权使用，请联系管理员编辑开通范围。操作路径：飞书管理后台 - 工作台 - 应用列表 - 极星协作，选择极星协作后配置使用范围即可", "ResultCode.FeiShuUserNotInAppUseScopeOfAuthority")
	FeishuImageIsEmpty                    = errors.AddResultCodeInfoWithSentry(606011, "文件格式不支持", "ResultCode.FeishuImageIsEmpty")
	NotSupportTypeForPushFeishuChat       = errors.AddResultCodeInfo(606012, "当前内容不需要推送至群聊", "ResultCode.NotSupportTypeForPushFeishuChat")
	OnlyFeishuCanUseFunction              = errors.AddResultCodeInfo(606013, "当前功能仅支持飞书用户", "ResultCode.OnlyFeishuCanUseFunction")
	FeiShuNoPowerToApply                  = errors.AddResultCodeInfo(606016, "暂无要申请的权限。", "ResultCode.FeiShuNoPowerToApply")
	FeiShuScopeNeedApply                  = errors.AddResultCodeInfo(606017, "暂无此权限，需要向管理员申请。", "ResultCode.FeiShuScopeNeedApply")
	FeiShuScopeOtherHasApply              = errors.AddResultCodeInfo(606018, "已有其他人申请过该权限，请等待管理员审核。", "ResultCode.FeiShuScopeOtherHasApply")
	FeiShuScopeUserHasApply               = errors.AddResultCodeInfo(606019, "已经申请过该权限，请等待管理员审核。", "ResultCode.FeiShuScopeUserHasApply")
	FeiShuNotScopeInCalendar              = errors.AddResultCodeInfo(606020, "没有日历相关权限", "ResultCode.FeiShuNotScopeInCalendar")
	FeiShuEventNotSupport                 = errors.AddResultCodeInfo(606021, "事件不在正常处理范围内。", "ResultCode.FeiShuEventNotSupport")
	FeiShuEventCheckTokenFailed           = errors.AddResultCodeInfo(606022, "无效的 token。", "ResultCode.FeiShuEventCheckTokenFailed")
	FeiShuGetInfoNoPermission             = errors.AddResultCodeInfo(606023, "无权限获取企业信息", "ResultCode.FeiShuGetInfoNoPermission")
	FeiShuAuthCodeExpired                 = errors.AddResultCodeInfo(606024, "授权码已失效，请重新授权登录", "ResultCode.FeiShuAuthCodeExpired")
	CardTitleEmpty                        = errors.AddResultCodeInfo(606025, "任务标题为空，不需要推送", "ResultCode.CardTitleEmpty")
	IssueTitleCannotEmpty                 = errors.AddResultCodeInfo(606026, "标题不能为空且不能超过500字符。", "ResultCode.IssueTitleCannotEmpty")
	CardColumnEmpty                       = errors.AddResultCodeInfo(606027, "没有需要推送的字段", "Result.CardColumnEmpty")
	CardNoNeedPush                        = errors.AddResultCodeInfo(606028, "不需要推送消息卡片", "Result.CardNoNeedPush")

	// 订单相关
	FsPricePlanNotExist   = errors.AddResultCodeInfo(700001, "飞书方案尚未与系统关联", "Result.FsPricePlanNotExist")
	DingPricePlanNotExist = errors.AddResultCodeInfo(700002, "钉钉方案尚未与系统关联", "Result.DingPricePlanNotExist")

	// 非法对象
	IllegalityPriority        = errors.AddResultCodeInfo(800001, "非法优先级", "Result.IllegalityPriority")
	IllegalityOwner           = errors.AddResultCodeInfo(800002, "非法负责人", "Result.IllegalityOwner")
	IllegalityFollower        = errors.AddResultCodeInfo(800003, "非法关注人", "Result.IllegalityFollower")
	IllegalityParticipant     = errors.AddResultCodeInfo(800004, "非法参与人", "Result.IllegalityParticipant")
	IllegalityProject         = errors.AddResultCodeInfo(800005, "非法项目", "Result.IllegalityProject")
	IllegalityOrg             = errors.AddResultCodeInfo(800006, "非法组织", "Result.IllegalityOrg")
	IllegalityIteration       = errors.AddResultCodeInfo(800007, "非法的迭代", "Result.IllegalityIteration")
	IllegalityIssue           = errors.AddResultCodeInfo(800008, "任务不存在或已被删除，无法操作", "Result.IllegalityIssue")
	IllegalityMQTTChannelType = errors.AddResultCodeInfo(800009, "非法的通道类型", "Result.IllegalityMQTTChannelType")

	// 项目关联对象
	ProjectRelationNotExist     = errors.AddResultCodeInfo(900001, "项目关联对象不存在", "Result.ProjectRelationNotExist")
	NeedChooseProjectObjectType = errors.AddResultCodeInfo(900002, "请选择任务类型", "Result.NeedChooseProjectObjectType")
	// 回收站任务
	RecycleObjectNotExist = errors.AddResultCodeInfo(900003, "资源已被恢复或不存在", "Result.RecycleObjectNotExist")
	IssueIdsNotBeenChosen = errors.AddResultCodeInfo(900004, "请选择需要复制的任务", "Result.IssueIdsNotBeenChosen")

	// 自定义字段
	CustomFieldNameLengthError  = errors.AddResultCodeInfo(900005, "自定义字段名称不能为空", "Result.CustomFieldNameLengthError")
	CustomFieldTypeError        = errors.AddResultCodeInfo(900006, "字段类型有误", "Result.CustomFieldTypeError")
	CustomFieldNotExist         = errors.AddResultCodeInfo(900007, "字段不存在", "Result.CustomFieldNotExist")
	SysCustomFieldCannotOperate = errors.AddResultCodeInfo(900008, "无法操作系统字段", "Result.SysCustomFieldCannotOperate")
	CustomFieldUseMore          = errors.AddResultCodeInfo(900009, "已有多个项目在使用此字段，无法删除", "Result.CustomFieldUseMore")
	DefaultFieldMustExist       = errors.AddResultCodeInfo(900010, "默认字段必须存在", "Result.DefaultFieldMustExist")
	CanNotRecoverDocuments      = errors.AddResultCodeInfo(900011, "对应的记录已被删除，无法恢复该附件", "Result.CanNotRecoverDocuments")

	// OpenAPI
	OpenAccessTokenIsEmpty = errors.AddResultCodeInfo(901001, "AccessToken为空", "Result.OpenAccessTokenIsEmpty")
	OpenAccessTokenInvalid = errors.AddResultCodeInfo(901002, "AccessToken无效", "Result.OpenAccessTokenInvalid")
	OpenAccessTokenExpired = errors.AddResultCodeInfo(901003, "AccessToken失效", "Result.OpenAccessTokenExpired")
	AppTicketNotAllocated  = errors.AddResultCodeInfo(901004, "尚未配置开发平台信息", "Result.AppTicketNotAllocated")

	// 视图
	IssueViewNotExist       = errors.AddResultCodeInfo(901100, "任务视图不存在", "Result.IssueViewNotExist")
	IssueViewNameLenInValid = errors.AddResultCodeInfo(901101, "抱歉，任务视图名称长度不能超过 20 个汉字长度。", "Result.IssueViewNameLenInValid")
	IssueViewOpDeny         = errors.AddResultCodeInfo(901102, "抱歉，您没有权限，请联系项目管理员！", "Result.IssueViewOpDeny")
	TableViewNotExist       = errors.AddResultCodeInfo(901103, "抱歉，当前项目缺少“表格视图”！", "Result.TableViewNotExist")

	ChangeProIntoAgileProWhenBelongManyPro = errors.AddResultCodeInfo(901201, "抱歉，不能将归属于多个项目的任务移动到敏捷项目中！", "Result.ChangeProIntoAgileProWhenBelongManyPro")
	IssueHasBelongThePro                   = errors.AddResultCodeInfo(901202, "抱歉，这个任务已经属于该项目！", "Result.IssueHasBelongThePro")
	IssueNotBelongThePro                   = errors.AddResultCodeInfo(901203, "抱歉，这个任务不属于该项目！", "Result.IssueNotBelongThePro")
	BindFeiShuNoPermissionErr              = errors.AddResultCodeInfo(901204, "需要系统管理员权限", "Result.BindFeiShuNoPermissionErr")
	OrgAlreadyBindPlatform                 = errors.AddResultCodeInfo(901205, "组织已绑定过其它外部平台", "Result.OrgAlreadyBindPlatform")
	TenantKeyAlreadyBindOrg                = errors.AddResultCodeInfo(901206, "外部团队已被其它组织绑定", "Result.TenantAlreadyBindOrg")
	OrgNotBindPlatform                     = errors.AddResultCodeInfo(901207, "组织未绑定过外部平台", "Result.OrgNotBindPlatform")
	UserAlreadyBindPlatform                = errors.AddResultCodeInfo(901208, "用户已绑定该平台", "Result.UserAlreadyBindPlatform")
	PANotExist                             = errors.AddResultCodeInfo(901209, "私有化授权码无效", "Result.PANotExist")
	CodeTokenInvalid                       = errors.AddResultCodeInfo(901210, "第三方平台登录信息失效", "Result.CodeTokenInvalid")
	OutUserAlreadyBind                     = errors.AddResultCodeInfo(901211, "第三方账号已被绑定", "Result.OutUserAlreadyBind")
	OrgNotInit                             = errors.AddResultCodeInfo(901212, "抱歉，您所在的组织尚未初始化。", "Result.OrgNotInit")
	GetCollaboratorRoleIdsFailed           = errors.AddResultCodeInfo(901213, "", "Result.GetCollaboratorRoleIdsFailed")
	DenyStartIssueChat                     = errors.AddResultCodeInfo(901214, "抱歉，您无权发起任务讨论，请联系管理员。", "Result.DenyStartIssueChat")
	CreateIssueChatDuplicate               = errors.AddResultCodeInfo(901215, "发起群聊操作重复，请重试。", "Result.CreateIssueChatDuplicate")
	RequestFrequentError                   = errors.AddResultCodeInfo(901216, "请求过于频繁，请稍后再试", "Result.RequestFrequentError")
	SetLabNoPermissionErr                  = errors.AddResultCodeInfo(901217, "没有权限操作实验室", "Result.SetLabNoPermissionErr")
)
