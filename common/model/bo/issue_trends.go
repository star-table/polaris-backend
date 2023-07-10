package bo

import (
	"time"

	"github.com/star-table/common/core/types"
	"github.com/star-table/polaris-backend/common/core/consts"
)

type IssueTrendsBo struct {
	PushType      consts.IssueNoticePushType `json:"pushType"`   //推送类型
	OrgId         int64                      `json:"orgId"`      //组织id
	OperatorId    int64                      `json:"operatorId"` //操作人id
	DataId        int64                      `json:"dataId"`
	IssueId       int64                      `json:"issueId"`       //任务id
	ParentIssueId int64                      `json:"parentIssueId"` //父任务id
	ProjectId     int64                      `json:"projectId"`     //项目id
	TableId       int64                      `json:"tableId"`
	PriorityId    int64                      `json:"priorityId"` //优先级id
	ParentId      int64                      `json:"parentId"`   //父任务id

	IssueTitle         string      `json:"issueTitle"`         //更新后的任务标题
	IssueRemark        string      `json:"issueRemark"`        //任务描述
	IssueStatusId      int64       `json:"issueStatusId"`      //更新后的任务状态
	IssuePlanStartTime *types.Time `json:"issuePlanStartTime"` //更新后的任务开始时间
	IssuePlanEndTime   *types.Time `json:"issuePlanEndTime"`   //更新后的任务结束时间

	SourceChannel string `json:"sourceChannel"` //来源通道

	BeforeOwner              []int64                      `json:"beforeOwner"`              //之前的负责人
	AfterOwner               []int64                      `json:"afterOwner"`               //新的负责人
	BeforeChangeFollowers    []int64                      `json:"beforeChangeFollowers"`    //之前的关注人
	AfterChangeFollowers     []int64                      `json:"afterChangeFollowers"`     //之后的负责人
	BeforeChangeParticipants []int64                      `json:"beforeChangeParticipants"` //之前的参与人
	AfterChangeParticipants  []int64                      `json:"afterChangeParticipants"`  //之后的参与人
	BeforeChangeAuditors     []int64                      `json:"beforeChangeAuditors"`     //之前的确认人
	AfterChangeAuditors      []int64                      `json:"afterChangeAuditors"`      //之后的确认人
	RelatingChange           map[string]*RelatingChangeBo `json:"relatingChange"`           // 关联字段的变化
	BaRelatingChange         map[string]*RelatingChangeBo `json:"baRelatingChange"`         // 前后置字段的变化
	SingleRelatingChange     map[string]*RelatingChangeBo `json:"singleRelatingChange"`     // 单向关联字段的变化
	BeforeChangeDocuments    map[string]interface{}       `json:"beforeChangeDocuments"`    //之前的附件
	AfterChangeDocuments     map[string]interface{}       `json:"afterChangeDocuments"`     //之后的附件
	BeforeChangeImages       map[string]interface{}       `json:"beforeChangeImages"`       //之前的图片
	AfterChangeImages        map[string]interface{}       `json:"afterChangeImages"`        //之后的图片
	BeforeWorkHourIds        []int64                      `json:"beforeWorkHourIds"`
	AfterWorkHourIds         []int64                      `json:"afterWorkHourIds"`
	UpdateOwner              bool                         `json:"updateOwner"`
	UpdateFollower           bool                         `json:"updateFollower"`
	UpdateAuditor            bool                         `json:"updateAuditor"`
	UpdateWorkHour           bool                         `json:"updateWorkHour"` // 标识是否变更了任务工时

	IssueChildren []int64 `json:"issueChildren"` //相关联的子任务

	OnlyNotice bool `json:"onlyNotice"` //如果为true，表示只处理通知

	OperateObjProperty string `json:"operateObjProperty"` //操作属性

	// 任务
	NewValue string `json:"newValue"` //新值
	OldValue string `json:"oldValue"` //老值

	Ext         TrendExtensionBo `json:"ext"`
	OperateTime time.Time        `json:"operateTime"` //操作时间
}

type TrendExtensionBo struct {
	IssueType           string               `json:"issueType"`
	ObjName             string               `json:"objName"`
	ChangeList          []TrendChangeListBo  `json:"changeList,omitempty"`
	MemberInfo          []SimpleUserInfoBo   `json:"memberInfo,omitempty"`
	RelationIssue       *RelationIssue       `json:"relationIssue,omitempty"`
	CommonChange        []string             `json:"commonChange,omitempty"`
	FolderId            int64                `json:"folderId"`
	MentionedUserIds    []int64              `json:"mentionedUserIds,omitempty"` //提及的用户列表
	CommentBo           *CommentBo           `json:"commentBo,omitempty"`
	ResourceInfo        []ResourceInfoBo     `json:"resourceInfo,omitempty"`
	Remark              string               `json:"remark"`
	FieldIds            []int64              `json:"fieldIds,omitempty"`
	ProjectObjectTypeId int64                `json:"projectObjectTypeId"`
	AuditInfo           *TrendAuditInfo      `json:"auditInfo,omitempty"`
	AutomationInfo      *TrendAutomationInfo `json:"automationInfo,omitempty"`

	// 新增自定义字段
	AddedFormFields []string `json:"addedFormFields,omitempty"`
	// 删除自定义字段
	DeletedFormFields []string `json:"deletedFormFields,omitempty"`
	// 修改的自定义字段
	UpdatedFormFields []string `json:"updatedFormFields,omitempty"`
}

type TrendAutomationInfo struct {
	WorkflowId   int64             `json:"workflowId,omitempty"`
	WorkflowName string            `json:"workflowName,omitempty"`
	ExecutionId  int64             `json:"executionId,omitempty"`
	TriggerUser  *SimpleUserInfoBo `json:"triggerUser,omitempty"`
}

type TrendAuditInfo struct {
	Status      int              `json:"status"`
	Remark      string           `json:"remark"`
	Attachments []ResourceInfoBo `json:"attachments"`
}

type SimpleTagInfoBo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ResourceInfoBo struct {
	Id         int64     `json:"id"`
	Url        string    `json:"url"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	UploadTime time.Time `json:"uploadTime"`
	Suffix     string    `json:"suffix"`
}

type TrendChangeListBo struct {
	Field                    string   `json:"field"`     // 字段名，字段id
	FieldType                string   `json:"fieldType"` // 字段类型
	FieldName                string   `json:"fieldName"` // 字段名称
	FieldNameValue           string   `json:"fieldNameValue"`
	AliasTitle               string   `json:"aliasTitle"` // 字段别名
	OldValue                 string   `json:"oldValue"`
	NewValue                 string   `json:"newValue"`
	OldUserIdsOrDeptIdsValue []string `json:"oldUserIdsOrDeptIdsValue"`
	NewUserIdsOrDeptIdsValue []string `json:"newUserIdsOrDeptIdsValue"`
	// 描述变更类，如：新增、删除、修改等
	// 主要针对自定义字段，比如多选、成员等。尚未用上
	ChangeTypeDesc string `json:"changeTypeDesc"`
	// 由于历史原因，工时变更时的数据使用了特殊的方式存储在 TrendChangeListBo，因此它和其他所有的 changeList 使用方式有所不同。这里做个标记
	IsForWorkHour bool `json:"isForWorkHour"`
}

type RelationIssue struct {
	// 关联信息id
	ID int64 `json:"id"`
	// 关联信息名称
	Title string `json:"title"`
}
