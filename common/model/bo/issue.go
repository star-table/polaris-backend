package bo

import (
	"time"

	tableV1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/model/bo/status"
)

var log = logger.GetDefaultLogger()

type IssueRelationBo struct {
	Id           int64      `db:"id,omitempty" json:"id"`
	OrgId        int64      `db:"org_id,omitempty" json:"orgId"`
	ProjectId    int64      `db:"project_id,omitempty" json:"projectId"`
	IssueId      int64      `db:"issue_id,omitempty" json:"issueId"`
	RelationId   int64      `db:"relation_id,omitempty" json:"relationId"`
	RelationCode string     `db:"relation_code,omitempty" json:"relationCode"`
	RelationType int        `db:"relation_type,omitempty" json:"relationType"`
	Status       int        `db:"status,omitempty" json:"status"`
	Creator      int64      `db:"creator,omitempty" json:"creator"`
	CreateTime   types.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64      `db:"updator,omitempty" json:"updator"`
	UpdateTime   types.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int        `db:"version,omitempty" json:"version"`
	IsDelete     int        `db:"is_delete,omitempty" json:"isDelete"`
}

type RelatingChangeBo struct {
	LinkToAdd   []int64
	LinkToDel   []int64
	LinkFromAdd []int64
	LinkFromDel []int64
}

type RelatingDesc struct {
	AddToDesc   string
	DelToDesc   string
	AddFromDesc string
	DelFromDesc string
}

type IssueBo struct {
	Id                  int64                  `json:"id"`
	OrgId               int64                  `json:"orgId"`
	TableId             int64                  `json:"tableId"`
	ProjectId           int64                  `json:"projectId"`
	ProjectObjectTypeId int64                  `json:"projectObjectTypeId"`
	Title               string                 `json:"title"`
	Code                string                 `json:"code"`
	Path                string                 `json:"path"`
	PriorityId          int64                  `json:"priorityId"`
	PlanStartTime       types.Time             `json:"planStartTime"`
	PlanEndTime         types.Time             `json:"planEndTime"`
	StartTime           types.Time             `json:"startTime"`
	EndTime             types.Time             `json:"endTime"`
	OwnerChangeTime     types.Time             `json:"ownerChangeTime"`
	PlanWorkHour        int                    `json:"planWorkHour"`
	IterationId         int64                  `json:"iterationId"`
	VersionId           int64                  `json:"versionId"`
	ParentId            int64                  `json:"parentId"`
	AuditStatus         int                    `json:"auditStatus"`
	AuditStatusDetail   map[string]int         `json:"auditStatusDetail"`
	Status              int64                  `json:"status"` // 对应无码字段issueStatus
	IssueStatusType     int32                  `json:"issueStatusType"`
	Creator             int64                  `json:"creator"`
	CreateTime          types.Time             `json:"createTime"`
	Updator             int64                  `json:"updator"`
	UpdateTime          types.Time             `json:"updateTime"`
	Sort                int64                  `json:"sort"`
	Order               int64                  `json:"order"`
	Version             int                    `json:"version"`
	IsDelete            int                    `json:"isDelete"`
	Remark              string                 `json:"remark"`
	OwnerId             []string               `json:"ownerId"`
	FollowerIds         []string               `json:"followerIds"`
	AuditorIds          []string               `json:"auditorIds"`
	RelatingIssue       RelatingIssue          `json:"relating"`
	BaRelatingIssue     RelatingIssue          `json:"baRelating"`
	OutChatId           string                 `json:"outChatId"` // 第三方平台群聊id
	LessData            map[string]interface{} `json:"lessData"`

	AppId          int64                  `json:"appId"`
	DataId         int64                  `json:"dataId"` // 对应无码的主键id，区别于极星的issueId
	OwnerIdI64     []int64                `json:"-"`
	FollowerIdsI64 []int64                `json:"-"`
	AuditorIdsI64  []int64                `json:"-"`
	Documents      map[string]interface{} `json:"-"`
	Images         map[string]interface{} `json:"-"`

	//IssueDetailBo  `json:"-"`
	//CustomField         string                 `json:"customField"`
	//OwnerInfos     []IssueUserBo          `json:"ownerInfos"`
	//FollowerInfos  []IssueUserBo          `json:"followerInfos"`
	//AuditorInfos   []IssueUserBo          `json:"auditorInfos"`
	//ExtraInfo map[string]interface{} `json:"extraInfo"` // 处理创建任务时用到的辅助信息
	//ModuleId            int64                  `json:"moduleId"`
	//ParticipantInfos *[]IssueUserBo `json:"-"`
	//IsFiling            int                    `json:"isFiling" db:"is_filing,omitempty"`
	//SourceId            int64      `json:"sourceId"`
	//IssueObjectTypeId   int64                  `json:"issueObjectTypeId"`
	//PropertyId        int64                  `json:"propertyId"`
	//Tags             []IssueTagReqBo `json:"tags"`
	//TypeForRelate int    `json:"typeForRelate"`
	//ParentTitle string `json:"parentTitle"`
	//ParentInfo       []*ParentInfo          `json:"parentInfo"`
	//ProjectTypeId int64 `json:"projectTypeId"`
	//ResourceIds      []int64                `json:"resourceIds"`
	//IssueProRelation *[]IssueRelationBo     `json:"issueProRelation"` // 任务所属项目的关联
	//FieldPriorityId  *int64                 `json:"_field_priority"`
}

type RelatingIssue struct {
	LinkTo   []string `json:"linkTo,omitempty"`
	LinkFrom []string `json:"linkFrom,omitempty"`
}

type ParentInfo struct {
	// id
	ID int64 `json:"id"`
	// 标题
	Title string `json:"title"`
	// code
	Code string `json:"code"`
}

// 任务和明细联合的Bo
type IssueAndDetailUnionBo struct {
	IssueId                  int64     `db:"issueId,omitempty"`
	IssueAuditStatus         int       `db:"issueAuditStatus,omitempty"`
	IssueStatusId            int64     `db:"issueStatusId,omitempty"`
	IssueProjectObjectTypeId int64     `db:"issueProjectObjectTypeId,omitempty"`
	StoryPoint               int       `db:"storyPoint,omitempty"`
	PlanEndTime              time.Time `db:"planEndTime,omitempty"`
	EndTime                  time.Time `db:"endTime,omitempty"`
	OwnerChangeTime          time.Time `db:"ownerChangeTime,omitempty"`
	CreateTime               time.Time `db:"createTime,omitempty"`

	// 新增的一些 field
	IssueStatusType int   `json:"issueStatusType"` // 任务状态的分类
	TableId         int64 `json:"tableId"`         // 任务所属的表
}

type IssueUpdateBo struct {
	IssueBo    IssueBo
	NewIssueBo IssueBo

	IssueUpdateCond         mysql.Upd
	IssueDetailRemark       *string
	IssueDetailRemarkDetail *string
	OwnerIds                []int64
	UpdateOwner             bool
	UpdateParticipant       bool
	UpdateFollower          bool
	UpdateAuditor           bool
	UpdateRelating          bool
	UpdateBeforeRelating    bool
	UpdateAfterRelating     bool
	DeleteResource          bool
	Participants            []int64
	Followers               []int64
	BeforeFollowers         []int64
	Auditors                []int64
	BindIssues              []int64
	UnbindIssues            []int64

	OperatorId int64
}

type IssueDailyNoticeBo struct {
	//所有未完成的本人负责的任务数量
	PendingSum uint64
	//截止时间为今天的本人负责的任务数量
	DueOfTodaySum uint64
	//已逾期的本人负责的任务数量
	OverdueSum uint64
	//即将逾期的本人负责的任务数量
	BeAboutToOverdueSum uint64
}

//func BuildIssueRelationBosFromUserBos(bos []IssueUserBo) []IssueRelationBo {
//	if bos == nil {
//		return []IssueRelationBo{}
//	}
//	relationBos := make([]IssueRelationBo, 0, len(bos))
//	for _, v := range bos {
//		relationBos = append(relationBos, v.IssueRelationBo)
//	}
//	return relationBos
//}

// 任务验证bo
type IssueAuthBo struct {
	Id           int64
	Owner        []int64
	ProjectId    int64
	Creator      int64
	Status       int64
	Participants []int64
	Followers    []int64
	TableId      int64
}

type HomeIssueOwnerInfoBo struct {
	ID     int64   `json:"id"`
	UserId int64   `json:"userId"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
	// 是否已被删除，为true则代表被组织移除
	IsDeleted bool `json:"isDeleted"`
	// 是否已被禁用, 为true则代表被组织禁用
	IsDisabled bool `json:"isDisabled"`
}

type HomeIssuePriorityInfoBo struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	BgStyle   string `json:"bgStyle"`
	FontStyle string `json:"fontStyle"`
}

type HomeIssueProjectInfoBo struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	IsFilling     int    `json:"isFilling"`
	ProjectTypeId int64  `json:"projectTypeId"`
}

type HomeIssueInfoPlusRespBo struct {
	Total int64 `json:"total"`
	// 实际总数量
	ActualTotal int64 `json:"actualTotal"`
	// 首页任务列表
	List []map[string]interface{} `json:"list"`
}

type HomeIssueInfoRespBo struct {
	Total int64 `json:"total"`
	// 实际总数量
	ActualTotal int64 `json:"actualTotal"`
	// 首页任务列表
	List []HomeIssueInfoBo `json:"list"`
}

type HomeIssueInfoBo struct {
	IssueId int64 `json:"issueId"`
	// 父任务id
	ParentId int64 `json:"parentId"`
	// 父任务信息
	ParentInfo []*ParentInfo `json:"parentInfo"`
	// 任务标题
	Title                 string                          `json:"title"`
	Issue                 IssueBo                         `json:"issue"`
	Project               *HomeIssueProjectInfoBo         `json:"project"`
	Owner                 []*HomeIssueOwnerInfoBo         `json:"owner"`
	Status                *status.StatusInfoBo            `json:"status"`
	Priority              *HomeIssuePriorityInfoBo        `json:"priority"`
	Tags                  []HomeIssueTagInfoBo            `json:"tags"`
	ChildsNum             int64                           `json:"childsNum"`
	ChildsFinishedNum     int64                           `json:"childsFinishedNum"`
	ProjectOBjectTypeName string                          `json:"projectOBjectTypeName"`
	AllStatus             []status.StatusInfoBo           `json:"allStatus"`
	SourceInfo            SimpleBasicNameBo               `json:"sourceInfo"`
	PropertyInfo          SimpleBasicNameBo               `json:"propertyInfo"`
	TypeInfo              SimpleBasicNameBo               `json:"typeInfo"`
	IterationName         string                          `json:"iterationName"`
	FollowerInfos         []*UserIDInfoBo                 `json:"followerInfos"`
	AuditorsInfo          []*UserIDInfoExtraForIssueAudit `json:"auditorsInfo"`
	RelateIssueCount      int64                           `json:"relateIssueCount"`
	RelateResourceCount   int64                           `json:"relateResourceCount"`
	RelateCommentCount    int64                           `json:"relateCommentCount"`
	// 自定义字段结果
	CustomField []*CustomValue `json:"customField"`
	// 任务的预估工时和实际工时
	WorkHourInfo  HomeIssueWorkHourInfo  `json:"workHourInfo"`
	AfterIssueIds []int64                `json:"afterIssueIds"`
	LessData      map[string]interface{} `json:"lessData"`
}

type CustomValue struct {
	// 字段id
	ID int64 `json:"id"`
	// 字段名称
	Name string `json:"name"`
	// 字段值
	Value interface{} `json:"value"`
	// 类型(1文本类型2单选框3多选框4日期选框5人员选择6是非选择7数字框)
	FieldType int `json:"fieldType"`
	// 值
	FieldValue []map[string]interface{} `json:"fieldValue"`
	// 是否属于组织字段库(1是2否，不选默认为否)
	IsOrgField int `json:"isOrgField"`
	// 字段描述
	Remark string `json:"remark"`
	// 字段
	Title string `json:"title"`
	// 启用状态（1启用2禁用）对于项目而言
	Status int `json:"status"`
}

type HomeIssueWorkHourInfo IssueWorkHour

type SimpleBasicNameBo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type HomeIssueTagInfoBo struct {
	// 标签id
	ID int64 `json:"id"`
	// 标签名
	Name string `json:"name"`
	// 背景颜色
	BgStyle string `json:"bgStyle"`
	// 字体颜色
	FontStyle string `json:"fontStyle"`
}

type UserIDInfoBo struct {
	//用户id
	Id int64 `json:"id"`
	//用户id
	UserID int64 `json:"userId"`
	//姓名
	Name string `json:"name"`
	//姓名拼音
	NamePy string `json:"namePy"`
	//头像
	Avatar string `json:"avatar"`
	//外部用户id
	EmplID string `json:"emplId"`
	//唯一id(弃用，统一emplId，也就是外部用户id)
	UnionID string `json:"unionId"`
	// 是否已被删除，为true则代表被组织移除
	IsDeleted bool `json:"isDeleted"`
	// 是否已被禁用, 为true则代表被组织禁用
	IsDisabled bool `json:"isDisabled"`
}

type UserIDInfoExtraForIssueAudit struct {
	//用户id
	Id int64 `json:"id"`
	//用户id
	UserID int64 `json:"userId"`
	//姓名
	Name string `json:"name"`
	//姓名拼音
	NamePy string `json:"namePy"`
	//头像
	Avatar string `json:"avatar"`
	//外部用户id
	EmplID string `json:"emplId"`
	//唯一id(弃用，统一emplId，也就是外部用户id)
	UnionID string `json:"unionId"`
	// 是否已被删除，为true则代表被组织移除
	IsDeleted bool `json:"isDeleted"`
	// 是否已被禁用, 为true则代表被组织禁用
	IsDisabled bool `json:"isDisabled"`
	// 状态(1未查看2已查看未审核3审核通过4驳回)
	AuditStatus int `json:"auditStatus"`
}

type AuditorAuditStatusBo struct {
	IssueId     int64 `json:"issueId"`
	AuditStatus int   `json:"auditStatus"`
}

type DepartmentMemberInfoBo struct {
	// 用户id
	UserID int64 `json:"userId"`
	// 用户名称
	Name string `json:"name"`
	//姓名拼音
	NamePy string `json:"namePy"`
	// 用户头像
	Avatar string `json:"avatar"`
	// 工号：企业下唯一
	EmplID string `json:"emplId"`
	// unionId： 开发者账号下唯一
	UnionID string `json:"unionId"`
	// 用户部门id
	DepartmentID int64 `json:"departmentId"`
	//用户组织状态
	OrgUserStatus int `json:"orgUserStatus"`
}

type IssueInfoBo struct {
	Issue    IssueBo                  `json:"issue"`
	Project  *HomeIssueProjectInfoBo  `json:"project"`
	Owner    *HomeIssueOwnerInfoBo    `json:"owner"`
	Status   *status.StatusInfoBo     `json:"status"`
	Priority *HomeIssuePriorityInfoBo `json:"priority"`

	ParticipantInfos []*UserIDInfoBo        `json:"participantInfos"`
	FollowerInfos    []*UserIDInfoBo        `json:"followerInfos"`
	NextStatus       []*status.StatusInfoBo `json:"nextStatus"`
}

type IssueReportDetailBo struct {
	Total          int64
	ReportUserName string
	StartTime      string
	EndTime        string
	List           []HomeIssueInfoBo
	ShareID        string
}

type TaskStatBo struct {
	NotStart   int64 `json:"notStart"`
	Processing int64 `json:"processing"`
	Completed  int64 `json:"completed"`
}

type IssueBoListCond struct {
	OrgId        int64
	ProjectId    *int64
	ProjectIds   []int64 // 项目的多选
	ProAppId     *int64
	IterationId  *int64
	RelationType *int
	UserId       *int64
	Ids          []int64
}

type IssueChildContStatus struct {
	ParentIssueId int64 `db:"parentIssueId" json:"parentIssueId"`
	Status        int64 `db:"status" json:"status"`
}

type IssueChildCountBo struct {
	ParentIssueId int64 `db:"parentIssueId"`
	Count         int64 `db:"count"`
}

type SelectIssueIdsCondBo struct {
	//计划开始时间开始范围
	BeforePlanEndTime *string `json:"beforePlanEndTime"`
	//计划开始时间结束范围
	AfterPlanEndTime *string `json:"afterPlanEndTime"`
	//计划开始时间开始范围
	BeforePlanStartTime *string `json:"beforePlanStartTime"`
	//计划开始时间结束范围
	AfterPlanStartTime *string `json:"afterPlanStartTime"`
}

type IssueIdBo struct {
	Id int64 `json:"id" db:"id"`
}

type IssueIdData struct {
	Id int64 `json:"issueId"`
}

type ParentIdData struct {
	Id int64 `json:"parentId"`
}

type IssueAndDetailInfoBo struct {
	Id                  int64      `db:"id,omitempty" json:"id"`
	Code                string     `db:"code,omitempty" json:"code"`
	ProjectObjectTypeId int64      `db:"project_object_type_id,omitempty" json:"projectObjectTypeId"`
	TableId             int64      `db:"table_id,omitempty" json:"tableId"`
	Title               string     `db:"title,omitempty" json:"title"`
	PriorityId          int64      `db:"priority_id,omitempty" json:"priorityId"`
	SourceId            int64      `db:"source_id,omitempty" json:"sourceId"`
	PropertyId          int64      `db:"property_id,omitempty" json:"propertyId"`
	IssueObjectTypeId   int64      `db:"issue_object_type_id,omitempty" json:"issueObjectTypeId"`
	IterationId         int64      `db:"iteration_id,omitempty" json:"oterationId"`
	Owner               int64      `db:"owner,omitempty" json:"owner"`
	PlanStartTime       types.Time `db:"plan_start_time,omitempty" json:"planStartTime"`
	PlanEndTime         types.Time `db:"plan_end_time,omitempty" json:"planEndTime"`
	EndTime             types.Time `db:"end_time,omitempty" json:"endTime"`
	ParentId            int64      `db:"parent_id,omitempty" json:"parentId"`
	Status              int64      `db:"status,omitempty" json:"status"`
	AuditStatus         int        `db:"audit_status,omitempty" json:"auditStatus"`
	CustomField         *string    `db:"custom_field,omitempty" json:"customField"`
	Creator             int64      `db:"creator,omitempty" json:"creator"`
	CreateTime          types.Time `db:"create_time,omitempty" json:"createTime"`
	Remark              *string    `db:"remark,omitempty" json:"remark"`
}

type IssueCombineBo struct {
	IssueAndDetailInfoBo
	Children []IssueAndDetailInfoBo `json:"children"`
}

type IssueMembersBo struct {
	MemberIds      []int64 `json:"memberIds"`
	OwnerId        []int64 `json:"ownerId"`
	ParticipantIds []int64 `json:"participantIds"`
	FollowerIds    []int64 `json:"followerIds"`
}

type SimpleIssueInfoForMqtt struct {
	Id                  int64 `json:"id"`
	ProjectId           int64 `json:"projectId"`
	Status              int64 `json:"status"`
	ProjectObjectTypeId int64 `json:"projectObjectTypeId"`
	Sort                int64 `json:"sort"`
}

type DeleteIssueBatchBo struct {
	SuccessIssues        []*IssueBo `json:"successIssues"`
	NoAuthIssues         []*IssueBo `json:"noAuthIssues"`
	RemainChildrenIssues []*IssueBo `json:"remainChildrenIssues"`
}

type MoveIssueBatchBo struct {
	SuccessIssues        []*IssueBo `json:"successIssues"`
	NoAuthIssues         []*IssueBo `json:"noAuthIssues"`
	RemainChildrenIssues []*IssueBo `json:"remainChildrenIssues"`
	ChildrenIssues       []*IssueBo `json:"childrenIssues"`
}

type IssueGroupSimpleInfoBo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

//type ArchiveIssueBo struct {
//	SuccessIssues []IssueBo `json:"successIssues"`
//	NoAuthIssues  []IssueBo `json:"noAuthIssues"`
//}

//type CancelArchiveIssueBo struct {
//	SuccessIssues []IssueBo `json:"successIssues"`
//	NoAuthIssues  []IssueBo `json:"noAuthIssues"`
//	// 父任务已被删除导致无法删除的子任务
//	ParentDeletedIssues []IssueBo `json:"parentDeletedIssues"`
//	// 父任务处于归档中导致无法删除的子任务
//	ParentFilingIssues []IssueBo `json:"parentFilingIssues"`
//}

type IssueNode struct {
	IssueId int64        `json:"issueId"`
	Path    string       `json:"path"`
	Child   []*IssueNode `json:"child"`
}

type NeedCopyColumnsValue struct {
	OwnerIds      []int64
	FollowerIds   []int64
	AuditorIds    []int64
	Remark        string
	PlanStartTime types.Time
	PlanEndTime   types.Time
	IterationId   int64
	PriorityId    interface{}
	WorkHour      interface{}
	Relating      interface{}
	BaRelating    interface{}
}

type Issue struct {
	Id       int64
	Children []*Issue
}

type IssueWorkHour struct {
	PredictWorkHour string `json:"predictWorkHour"`
	ActualWorkHour  string `json:"actualWorkHour"`
	// 预估工时详情列表
	PredictList []*PredictListItem `json:"predictList"`
	// 实际工时详情列表
	ActualList []*ActualListItem `json:"actualList"`
}

// 预估工时详情列表单个对象
type PredictListItem struct {
	// 工时执行人名字
	Name string `json:"name"`
	// 工时，单位：小时。
	WorkHour string `json:"workHour"`
}

// 实际工时详情列表单个对象
type ActualListItem struct {
	// 工时执行人名字
	Name                   string                    `json:"name"`
	ActualWorkHourDateList []*ActualWorkHourDateItem `json:"actualWorkHourDateList"`
}

type ActualWorkHourDateItem struct {
	// 实际工时的日期，开始日期。
	Date string `json:"date"`
	// 工时，单位：小时。
	WorkHour string `json:"workHour"`
}

type ActualWorkHourDateItemForInt struct {
	// 实际工时的日期，开始日期。
	Date string `json:"date"`
	// 工时，单位：小时。
	WorkHourInt uint32 `json:"workHourInt"`
}

type TmpLcData struct {
	IssueId int64                  `json:"issueId"`
	Data    map[string]interface{} `json:"data"`
}

type TableColumn struct {
	Name              string             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`               // 其实就内部列的唯一id
	Label             string             `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`             // 标题
	AliasTitle        string             `protobuf:"bytes,3,opt,name=aliasTitle,proto3" json:"aliasTitle,omitempty"`   // 别名
	Description       string             `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"` // 描述
	IsSys             bool               `protobuf:"varint,5,opt,name=isSys,proto3" json:"isSys,omitempty"`
	IsOrg             bool               `protobuf:"varint,6,opt,name=isOrg,proto3" json:"isOrg,omitempty"`
	Writable          bool               `protobuf:"varint,7,opt,name=writable,proto3" json:"writable,omitempty"`
	Editable          bool               `protobuf:"varint,8,opt,name=editable,proto3" json:"editable,omitempty"`
	Unique            bool               `protobuf:"varint,9,opt,name=unique,proto3" json:"unique,omitempty"`                       // 该类型是否是唯一的列
	UniquePreHandler  string             `protobuf:"bytes,10,opt,name=uniquePreHandler,proto3" json:"uniquePreHandler,omitempty"`   // 唯一预处理函数，"urlDomainParse": url域名解析，保留一位
	SensitiveStrategy string             `protobuf:"bytes,11,opt,name=sensitiveStrategy,proto3" json:"sensitiveStrategy,omitempty"` // 脱敏策略
	SensitiveFlag     int32              `protobuf:"varint,12,opt,name=sensitiveFlag,proto3" json:"sensitiveFlag,omitempty"`        // 是否脱敏
	Field             *TableColumnOption `protobuf:"bytes,13,opt,name=field,proto3" json:"field,omitempty"`
}

type TableColumnOption struct {
	Type       *tableV1.ColumnType        `protobuf:"varint,1,opt,name=type,proto3,enum=table.v1.ColumnType" json:"type,omitempty"`                // 列类型
	CustomType string                     `protobuf:"bytes,2,opt,name=customType,proto3" json:"customType,omitempty"`                              // 自定义类型
	DataType   *tableV1.StorageColumnType `protobuf:"varint,3,opt,name=dataType,proto3,enum=table.v1.StorageColumnType" json:"dataType,omitempty"` // 数据类型
	Props      interface{}                `protobuf:"bytes,4,opt,name=props,proto3" json:"props,omitempty"`                                        // 列类型下特有的数据
}

type WorkHour struct {
	ActualHour string `json:"actualHour"`
	PlanHour   string `json:"planHour"`
}
