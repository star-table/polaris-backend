package projectvo

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
)

// 批量创建
type BatchCreateIssueReqVo struct {
	OrgId         int64                    `json:"orgId"`
	UserId        int64                    `json:"userId"`
	AppId         int64                    `json:"appId"`
	ProjectId     int64                    `json:"projectId"`
	TableId       int64                    `json:"tableId"`
	BeforeDataId  int64                    `json:"beforeDataId"`
	AfterDataId   int64                    `json:"afterDataId"`
	Data          []map[string]interface{} `json:"data"`
	IsIdGenerated bool                     `json:"isIdGenerated"`
}

type BatchCreateIssueRespVo struct {
	vo.Err
	Data       []map[string]interface{}          `json:"data"`
	UserDepts  map[string]*uservo.MemberDept     `json:"userDepts,omitempty"`
	RelateData map[string]map[string]interface{} `json:"relateData,omitempty"`
}

// 批量更新
type BatchUpdateIssueReqInnerVo struct {
	OrgId          int64                       `json:"orgId"`
	UserId         int64                       `json:"userId"`
	AppId          int64                       `json:"appId"`
	ProjectId      int64                       `json:"projectId"`
	TableId        int64                       `json:"tableId"`
	Data           []map[string]interface{}    `json:"data"`
	BeforeDataId   int64                       `json:"beforeDataId"`
	AfterDataId    int64                       `json:"afterDataId"`
	TodoId         int64                       `json:"todoId"`
	TodoOp         int                         `json:"todoOp"`
	TodoMsg        string                      `json:"todoMsg"`
	TrendPushType  *consts.IssueNoticePushType `json:"push_type"`       // 额外动态信息：如果携带，则会覆盖原始的批量更新的动态信息
	TrendExtension *bo.TrendExtensionBo        `json:"trend_extension"` // 额外动态信息：如果携带，则会覆盖原始的批量更新的动态信息
}

type BatchUpdateIssueInput struct {
	AppId        int64                    `json:"appId"`
	ProjectId    int64                    `json:"projectId"`
	TableId      int64                    `json:"tableId"`
	Data         []map[string]interface{} `json:"data"`
	BeforeDataId int64                    `json:"beforeDataId"`
	AfterDataId  int64                    `json:"afterDataId"`
	TodoId       int64                    `json:"todoId"`
	TodoOp       int                      `json:"todoOp"`
	TodoMsg      string                   `json:"todoMsg"`
}

type BatchUpdateIssueReqVo struct {
	OrgId  int64                  `json:"orgId"`
	UserId int64                  `json:"userId"`
	Input  *BatchUpdateIssueInput `json:"input"`
}

// 批量审批
type BatchAuditIssueReqVo struct {
	OrgId  int64                      `json:"orgId"`
	UserId int64                      `json:"userId"`
	Input  *vo.LessBatchAuditIssueReq `json:"input"`
}

type BatchUrgeIssueReq struct {
	IssueIds []int64 `json:"issueIds"`
	Message  string  `json:"message"`
}

// 批量催办
type BatchUrgeIssueReqVo struct {
	OrgId  int64                     `json:"orgId"`
	UserId int64                     `json:"userId"`
	Input  *vo.LessBatchUrgeIssueReq `json:"input"`
}
