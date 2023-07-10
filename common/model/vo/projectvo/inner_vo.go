package projectvo

import (
	"github.com/star-table/polaris-backend/common/model/vo"
)

type TriggerBy struct {
	TriggerBy        string `json:"triggerBy"`
	WorkflowId       int64  `json:"workflowId,omitempty"`
	WorkflowName     string `json:"workflowName,omitempty"`
	ExecutionId      int64  `json:"executionId,omitempty"`
	TriggerUserId    int64  `json:"triggerUserId,omitempty"`
	IsCreateTemplate bool   `json:"isCreateTemplate,omitempty"`
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Inner Issue Filter
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type InnerIssueFilterInput struct {
	AppId     int64             `json:"appId"`
	TableId   int64             `json:"tableId"`
	Columns   []string          `json:"columns"`
	Condition *vo.LessCondsData `json:"condition"`
	Orders    []*vo.LessOrder   `json:"orders"`
}

type InnerIssueFilterReq struct {
	OrgId  int64                  `json:"orgId"`
	UserId int64                  `json:"userId"`
	Page   int                    `json:"page"`
	Size   int                    `json:"size"`
	Input  *InnerIssueFilterInput `json:"input"`
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Inner Issue Create
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type InnerIssueCreateInput struct {
	TriggerBy *TriggerBy               `json:"triggerBy"`
	AppId     int64                    `json:"appId"`
	TableId   int64                    `json:"tableId"`
	Data      []map[string]interface{} `json:"data"`
}

type InnerIssueCreateReq struct {
	OrgId  int64                  `json:"orgId"`
	UserId int64                  `json:"userId"`
	Input  *InnerIssueCreateInput `json:"input"`
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Inner Issue Create By Copy
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type InnerIssueCreateByCopyInput struct {
	TriggerBy    *TriggerBy               `json:"triggerBy"`
	AppId        int64                    `json:"appId"`
	TableId      int64                    `json:"tableId"`
	IssueIds     []int64                  `json:"issueIds"`
	Data         []map[string]interface{} `json:"data"`
	IsStaticCopy bool                     `json:"isStaticCopy"`
}

type InnerIssueCreateByCopyReq struct {
	OrgId  int64                        `json:"orgId"`
	UserId int64                        `json:"userId"`
	Input  *InnerIssueCreateByCopyInput `json:"input"`
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Inner Issue Update
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type InnerIssueUpdateInput struct {
	TriggerBy *TriggerBy               `json:"triggerBy"`
	AppId     int64                    `json:"appId"`
	TableId   int64                    `json:"tableId"`
	Data      []map[string]interface{} `json:"data"`
}

type InnerIssueUpdateReq struct {
	OrgId  int64                  `json:"orgId"`
	UserId int64                  `json:"userId"`
	Input  *InnerIssueUpdateInput `json:"input"`
}

//type ColumnSetting struct {
//	CanRead    bool `json:"canRead"`
//	CanWrite   bool `json:"canWrite"`
//	IsRequired bool `json:"isRequired"`
//}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Inner Create Todo FillIn
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//type TodoOperator struct {
//	Type string      `json:"type"`
//	Ids  interface{} `json:"ids"`
//}

//type TodoFillInParameters struct {
//	FormSettings map[string]*ColumnSetting `json:"formSettings"`
//}

//type InnerCreateTodoFillInInput struct {
//	TriggerBy              *TriggerBy           `json:"triggerBy"`
//	AppId                  int64                `json:"appId"`
//	TableId                int64                `json:"tableId"`
//	DataId                 int64                `json:"dataId"`
//	TriggerUserId          int64                `json:"triggerUserId"`
//	AllowWithdrawByTrigger int                  `json:"allowWithdrawByTrigger"`
//	AllowUrgeByTrigger     int                  `json:"allowUrgeByTrigger"`
//	Operators              []*TodoOperator      `json:"operators"`
//	Parameters             TodoFillInParameters `json:"parameters"`
//}
//
//type InnerCreateTodoFillInReq struct {
//	OrgId  int64                       `json:"orgId"`
//	UserId int64                       `json:"userId"`
//	Input  *InnerCreateTodoFillInInput `json:"input"`
//}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Inner Create Todo Audit
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//type TodoAuditParameters struct {
//	SignMode     string                    `json:"signMode"`
//	FormSettings map[string]*ColumnSetting `json:"formSettings"`
//}

//type InnerCreateTodoAuditInput struct {
//	TriggerBy              *TriggerBy          `json:"triggerBy"`
//	AppId                  int64               `json:"appId"`
//	TableId                int64               `json:"tableId"`
//	DataId                 int64               `json:"dataId"`
//	TriggerUserId          int64               `json:"triggerUserId"`
//	AllowWithdrawByTrigger int                 `json:"allowWithdrawByTrigger"`
//	AllowUrgeByTrigger     int                 `json:"allowUrgeByTrigger"`
//	Operators              []*TodoOperator     `json:"operators"`
//	Parameters             TodoAuditParameters `json:"parameters"`
//}
//
//type InnerCreateTodoAuditReq struct {
//	OrgId  int64                      `json:"orgId"`
//	UserId int64                      `json:"userId"`
//	Input  *InnerCreateTodoAuditInput `json:"input"`
//}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Inner Create Todo Hook
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type InnerCreateTodoHookInput struct {
	TodoId int64 `json:"todoId"`
}

type InnerCreateTodoHookReq struct {
	OrgId  int64                     `json:"orgId"`
	UserId int64                     `json:"userId"`
	Input  *InnerCreateTodoHookInput `json:"input"`
}
