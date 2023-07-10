package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type IssueSourceListReqVo struct {
	Page  uint
	Size  uint
	Input vo.IssueSourcesReq `json:"input"`
	//UserId  int64           `json:"userId"`
	OrgId int64 `json:"orgId"`
}

type IssueSourceListRespVo struct {
	vo.Err
	IssueSourceList *vo.IssueSourceList `json:"data"`
}

type CreateIssueSourceReqVo struct {
	Input  vo.CreateIssueSourceReq `json:"input"`
	UserId int64                   `json:"userId"`
}

type UpdateIssueSourceReqVo struct {
	Input  vo.UpdateIssueSourceReq `json:"input"`
	UserId int64                   `json:"userId"`
}

type DeleteIssueSourceReqVo struct {
	Input  vo.DeleteIssueSourceReq `json:"input"`
	UserId int64                   `json:"userId"`
	OrgId  int64                   `json:"orgId"`
}
