package projectvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type IterationListReqVo struct {
	Page  uint                 `json:"page"`
	Size  uint                 `json:"size"`
	Input *vo.IterationListReq `json:"input"`
	OrgId int64
}

type IterationListRespVo struct {
	vo.Err
	IterationList *vo.IterationList `json:"data"`
}

type IterationStatusTypeStatRespVo struct {
	vo.Err
	IterationStatusTypeStat *vo.IterationStatusTypeStatResp `json:"data"`
}

type IterationInfoRespVo struct {
	vo.Err
	IterationInfo *vo.IterationInfoResp `json:"data"`
}

type IterationIssueRelateReqVo struct {
	Input  vo.IterationIssueRealtionReq `json:"input"`
	OrgId  int64                        `json:"orgId"`
	UserId int64                        `json:"userId"`
}

type UpdateIterationStatusReqVo struct {
	Input  vo.UpdateIterationStatusReq `json:"input"`
	OrgId  int64                       `json:"orgId"`
	UserId int64                       `json:"userId"`
}

type CreateIterationReqVo struct {
	Input  vo.CreateIterationReq `json:"input"`
	OrgId  int64                 `json:"orgId"`
	UserId int64                 `json:"userId"`
}

type UpdateIterationReqVo struct {
	Input  vo.UpdateIterationReq `json:"input"`
	OrgId  int64                 `json:"orgId"`
	UserId int64                 `json:"userId"`
}

type DeleteIterationReqVo struct {
	Input  vo.DeleteIterationReq `json:"input"`
	OrgId  int64                 `json:"orgId"`
	UserId int64                 `json:"userId"`
}

type IterationStatusTypeStatReqVo struct {
	Input *vo.IterationStatusTypeStatReq `json:"input"`
	OrgId int64                          `json:"orgId"`
}

type IterationInfoReqVo struct {
	Input vo.IterationInfoReq `json:"input"`
	OrgId int64               `json:"orgId"`
}

type GetNotCompletedIterationBoListReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type GetNotCompletedIterationBoListRespVo struct {
	IterationBoList []bo.IterationBo `json:"data"`
	vo.Err
}

type UpdateIterationStatusTimeReqVo struct {
	OrgId  int64                 `json:"orgId"`
	UserId int64                 `json:"userId"`
	Input vo.UpdateIterationStatusTimeReq `json:"data"`
}

type GetSimpleIterationInfoReqVo struct {
	OrgId int64 `json:"orgId"`
	IterationIds []int64 `json:"iterationIds"`
}

type GetSimpleIterationInfoRespVo struct {
	vo.Err
	Data []bo.IterationBo `json:"data"`
}

type UpdateIterationSortReqVo struct {
	OrgId int64 `json:"orgId"`
	UserId int64 `json:"userId"`
	Params vo.UpdateIterationSortReq `json:"params"`
}
