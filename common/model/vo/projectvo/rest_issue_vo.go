package projectvo

import (
	"github.com/star-table/polaris-backend/common/model/bo/status"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type GetIssueIdsByOrgIdReqVo struct {
	OrgId  int64                  `json:"orgId"`
	UserId int64                  `json:"userId"`
	Input  *GetIssueIdsByOrgIdReq `json:"input"`
}

type GetIssueIdsByOrgIdReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type GetIssueIdsByOrgIdRespVo struct {
	vo.Err
	Data *GetIssueIdsByOrgIdResp `json:"data"`
}

type GetIssueIdsByOrgIdResp struct {
	List  []int64 `json:"list"`
	Total int64   `json:"total"`
}

type InsertIssueProRelationReqVo struct {
	OrgId  int64                      `json:"orgId"`
	UserId int64                      `json:"userId"`
	Input  *InsertIssueProRelationReq `json:"input"`
}

type InsertIssueProRelationReq struct {
	IssueIds []int64 `json:"issueIds"`
}

type GetIssueAllStatusNewReq struct {
	OrgId int64                    `json:"orgId"`
	Input GetIssueAllStatusNewData `json:"input"`
}

type GetIssueAllStatusNewData struct {
	ProjectIds []int64 `json:"projectIds"`
	TableIds   []int64 `json:"tableIds"`
}

type GetIssueAllStatusNewResp struct {
	vo.Err
	Data map[int64][]status.StatusInfoBo `json:"data"`
}
