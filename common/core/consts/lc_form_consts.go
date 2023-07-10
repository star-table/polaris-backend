package consts

const (
	LcJsonColumn = "data"

	TempFieldChildren      = "children"
	TempFieldChildrenCount = "childrenCount"
	TempFieldIssueId       = "_issueId"
	TempFieldCode          = "_code"
	TempFieldTodoResult    = "todoResult"

	BasicFieldId                  = "id"
	BasicFieldDataId              = "dataId"
	BasicFieldIssueId             = "issueId"
	BasicFieldTitle               = "title"
	BasicFieldOrgId               = "orgId"
	BasicFieldAppId               = "appId"
	BasicFieldProjectId           = "projectId"
	BasicFieldTableId             = "tableId"
	BasicFieldIterationId         = "iterationId"
	BasicFieldParentId            = "parentId"
	BasicFieldPath                = "path"
	BasicFieldCreator             = "creator"
	BasicFieldCreateTime          = "createTime"
	BasicFieldUpdator             = "updator"
	BasicFieldUpdateTime          = "updateTime"
	BasicFieldMoveTime            = "moveTime"
	BasicFieldDelFlag             = "delFlag"
	BasicFieldRecycleFlag         = "recycleFlag"
	BasicFieldCustomField         = "customField"
	BasicFieldCode                = "code"
	BasicFieldOrder               = "order"
	BasicFieldProjectObjectTypeId = "projectObjectTypeId"
	BasicFieldIssueStatusType     = "issueStatusType" // 任务状态分类：1未开始；2进行中；3已完成
	BasicFieldIssueStatus         = "issueStatus"
	BasicFieldStatus              = "status"
	BasicFieldAuditStatus         = "auditStatus"
	BasicFieldAuditStatusDetail   = "auditStatusDetail"
	BasicFieldPlanStartTime       = "planStartTime"
	BasicFieldPlanEndTime         = "planEndTime"
	BasicFieldEndTime             = "endTime"
	BasicFieldOwnerId             = "ownerId"
	BasicFieldFollowerIds         = "followerIds"
	BasicFieldAuditorIds          = "auditorIds"
	BasicFieldRemark              = "remark"
	BasicFieldRemarkDetail        = "remarkDetail"
	BasicFieldOwnerChangeTime     = "ownerChangeTime"
	BasicFieldTag                 = "_field_tag"
	BasicFieldPriority            = "_field_priority"
	BasicFieldDemandType          = "_field_demand_type"
	BasicFieldDemandSource        = "_field_demand_source"
	BasicFieldBugProperty         = "_field_bug_property"
	BasicFieldBugType             = "_field_bug_type"
	BasicFieldRelating            = "relating"
	BasicFieldBaRelating          = "baRelating"
	BasicFieldRelatingLinkTo      = "linkTo"
	BasicFieldRelatingLinkFrom    = "linkFrom"
	BasicFieldDocument            = "document"
	BasicFieldWorkHour            = "workHour"
	BasicFieldCollaborators       = "collaborators"
	RelateTableId                 = "relateTableId"
	BasicFieldColla
	//BasicFieldIssuePropertyId     = "issuePropertyId"
	//BasicFieldIssueObjectTypeId   = "issueObjectTypeId"
	//BasicFieldSourceId         = "sourceId"
	//BasicFieldPlanWorkHour        = "planWorkHour" // 已弃用；改为使用 ProBasicFieldWorkHour

	BasicFieldAuditStarter = "auditStarter" // 审批发起人
	BasicFieldOutChatId    = "outChatId"    // 任务讨论群id

	//一些其他辅助的任务字段 key
	HelperFieldCreateFrom    = "createFrom"    // 任务创建的来源，枚举值有：普通创建、导入、复制等
	HelperFieldIsApplyTpl    = "isApplyTpl"    // 是否是应用模板
	HelperFieldOperateUserId = "operateUserId" // 操作人 key
	HelperFieldIsParentDesc  = "isParentDesc"  // excel 列，是否是父任务的描述：父任务，子任务

	// 项目数据同步到无码，字段汇总
	ProBasicFieldName                = "name"
	ProBasicFieldProId               = "projectId"
	ProBasicFieldAppId               = "appId"
	ProBasicFieldCode                = "code"
	ProBasicFieldPreCode             = "preCode"
	ProBasicFieldOwnerIds            = "ownerIds" // 支持多个。值为 `U_23310`
	ProBasicFieldProTypeId           = "projectTypeId"
	ProBasicFieldPriorityId          = "priorityId" // 优先级 单选
	ProBasicFieldPlanStartTime       = "planStartTime"
	ProBasicFieldPlanEndTime         = "planEndTime"
	ProBasicFieldStatus              = "proStatus"    // 项目状态。因为无码的数据本身自带 status 字段，与北极星项目的 status 有冲突，因此，极星项目的状态字段用 proStatus 表示。
	ProBasicFieldPublicStatus        = "publicStatus" // 项目的公有、私有
	ProBasicFieldTemplateFlag        = "templateFlag" // 模板标识：1：是，2：否
	ProBasicFieldResource            = "resource"     // 项目**封面图**，在无码中，可以存储成 string
	ProBasicFieldIsFiling            = "isFiling"     // 1归档；2正常（不归档）
	ProBasicFieldRemark              = "remark"
	ProBasicFieldNotice              = "notice"
	ProBasicFieldIsEnableWorkHours   = "isEnableWorkHours"     // 是否启用工时。1启用；2未启用
	ProBasicFieldIsStared            = "isStared"              // 是否星标。1启用；2未启用
	ProBasicFieldOutCalendar         = "outCalendar"           // 外部日历日程标识
	ProBasicFieldOutCalendarSettings = "syncOutCalendarStatus" // 日历日程设置   类型是 int
	ProBasicFieldOutChat             = "outChat"               // 外部群聊标识。可以有多个
	ProBasicFieldChatSettings        = "outChatSettings"       // 群聊设置  类型是 map[string]interface{}
	ProBasicFieldParticipantIds      = "participantIds"        // 参与人/部门。如果是人，值为 `U_23310`；如果是部门，则为 `D_2312312312312`
	ProBasicFieldSort                = "sort"
	ProBasicFieldCreator             = "creator"
	ProBasicFieldUpdator             = "updator"
	ProBasicFieldCreateTime          = "createTime"
	ProBasicFieldUpdateTime          = "updateTime"
	ProBasicFieldVersion             = "version"  // 版本、乐观锁
	ProBasicFieldIsDelete            = "isDelete" // 1删除；2正常
	ProBasicFieldWorkHour            = "workHour" // 工时 存储任务的工时信息，属于**任务**字段
)

var (
	ProFormFieldsMap = map[string]string{
		ProBasicFieldName:                "项目名",
		ProBasicFieldAppId:               "应用 id",
		ProBasicFieldProId:               "项目 id",
		ProBasicFieldCode:                "项目 code",
		ProBasicFieldPreCode:             "项目 code 前缀",
		ProBasicFieldOwnerIds:            "项目负责人",
		ProBasicFieldParticipantIds:      "参与人、参与部门",
		ProBasicFieldProTypeId:           "项目类型",
		ProBasicFieldPriorityId:          "优先级",
		ProBasicFieldPlanStartTime:       "开始时间",
		ProBasicFieldPlanEndTime:         "结束时间",
		ProBasicFieldPublicStatus:        "项目的公有私有",
		ProBasicFieldStatus:              "项目状态",
		ProBasicFieldTemplateFlag:        "模板标识",
		ProBasicFieldResource:            "项目资源",
		ProBasicFieldIsFiling:            "是否归档",
		ProBasicFieldRemark:              "项目描述",
		ProBasicFieldNotice:              "项目公告",
		ProBasicFieldIsEnableWorkHours:   "是否启用工时",
		ProBasicFieldOutCalendar:         "外部日历标识",
		ProBasicFieldOutCalendarSettings: "日历设置状态",
		ProBasicFieldOutChat:             "外部群聊",
		ProBasicFieldChatSettings:        "外部群聊设置",
		ProBasicFieldSort:                "项目排序值",
		ProBasicFieldCreator:             "创建人",
		ProBasicFieldUpdator:             "更新人",
		ProBasicFieldCreateTime:          "创建时间",
		ProBasicFieldUpdateTime:          "更新时间",
		ProBasicFieldVersion:             "版本、乐观锁",
		ProBasicFieldIsDelete:            "是否删除",
	}
)

const (
	CopyIssueWorkHour = "copyIssueWorkHour"
)

const (
	DataTypeCascader = "cascader"
)

const KeyCurrentUser = "U_${current_user}"

const DefaultAvatar = "https://polaris-hd2.oss-cn-shanghai.aliyuncs.com/front_resources/picture.svg"

// 基础字段
var DefaultField = []interface{}{
	BasicFieldTitle,       //标题
	BasicFieldProjectId,   //所属项目
	BasicFieldIterationId, //迭代
	BasicFieldParentId,    //父任务ID
	BasicFieldCode,        //编号
	BasicFieldRemark,      //任务描述
	BasicFieldOwnerId,     //负责人
	// BasicFieldProjectObjectTypeId, //任务栏 // 任务栏不再是基础字段，先注释
	BasicFieldIssueStatus,   //任务状态
	BasicFieldPlanStartTime, //开始时间
	BasicFieldPlanEndTime,   //截止时间
	BasicFieldFollowerIds,   //关注人
	BasicFieldAuditorIds,    //确认人
}
