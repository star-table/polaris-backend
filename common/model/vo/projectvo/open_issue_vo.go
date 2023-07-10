package projectvo

import (
	"github.com/star-table/common/core/types"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type OpenCreateIssueReq struct {
	vo.LessCreateIssueReq
	// 操作人
	OperatorId int64 `json:"operatorId"`
}

type OpenUpdateIssueReqVo struct {
	Data   *OpenUpdateIssueReq `json:"data"`
	UserId int64               `json:"userId"`
	OrgId  int64               `json:"orgId"`
}

type OpenDeleteIssueReqVo struct {
	Data   *OpenDeleteIssueReq `json:"data"`
	UserId int64               `json:"userId"`
	OrgId  int64               `json:"orgId"`
}

type OpenIssueListReqVo struct {
	Data   *HomeIssueInfoReq `json:"data"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	UserId int64             `json:"userId"`
	OrgId  int64             `json:"orgId"`
}

type OpenDeleteIssueReq struct {
	IssueId int64 `json:"issueId"`
}

type OpenUpdateIssueReq struct {
	vo.LessUpdateIssueReq

	// 要更新的状态id
	NextStatusID *int64 `json:"nextStatusId"`
	// 要更新的状态类型,1: 未开始，2：进行中，3：已完成
	NextStatusType *int `json:"nextStatusType"`
	// 操作人
	OperatorId int64 `json:"operatorId"`
}

type OpenIssueListReq struct {
	// 操作人
	OperatorId int64 `json:"operatorId"`
	// 项目编号
	ProjectId int64 `json:"projectId"`

	vo.HomeIssuesRestReq
}

type OpenCreateIssueCommentReq struct {
	OperatorId int64  `json:"operatorId"`
	Comment    string `json:"comment"`
	// 提及的用户id
	MentionedUserIds []int64 `json:"mentionedUserIds"`
}

type OpenMirrorsStatReq struct {
	OperatorId int64    `json:"operatorId"`
	AppIds     []string `json:"appIds"`
}

type OpenIssueCommentsListResp struct {
	Page int64 `json:"page"`
	Size int64 `json:"size"`
	//LastId int64 `json:"lastId"`
	Total int64                  `json:"total"`
	List  []OpenIssueCommentInfo `json:"list"`
}

type OpenIssueCommentInfo struct {
	Id          int64          `json:"id"`
	CreatorInfo *vo.UserIDInfo `json:"creatorInfo"`
	Comment     *string        `json:"comment"`
	CreateTime  types.Time     `json:"createTime"`
}

type OpenIssueCommentsListReq struct {
	vo.TrendReq
	OperatorId int64 `json:"operatorId"`
}
