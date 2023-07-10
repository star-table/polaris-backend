package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type GetIssueWorkHoursListReqVo struct {
	Input  vo.GetIssueWorkHoursListReq `json:"input"`
	UserId int64                       `json:"userId"`
	OrgId  int64                       `json:"orgId"`
}

type GetIssueWorkHoursListRespVo struct {
	vo.Err
	Data *vo.GetIssueWorkHoursListResp `json:"data"`
}

type BoolRespVo struct {
	vo.Err
	Data *vo.BoolResp `json:"data"`
}

type CreateIssueWorkHoursReqVo struct {
	Input  vo.CreateIssueWorkHoursReq `json:"input"`
	UserId int64                      `json:"userId"`
	OrgId  int64                      `json:"orgId"`
}

type CreateMultiIssueWorkHoursReqVo struct {
	Input  vo.CreateMultiIssueWorkHoursReq `json:"input"`
	UserId int64                           `json:"userId"`
	OrgId  int64                           `json:"orgId"`
}

type UpdateIssueWorkHoursReqVo struct {
	Input  vo.UpdateIssueWorkHoursReq `json:"input"`
	UserId int64                      `json:"userId"`
	OrgId  int64                      `json:"orgId"`
}

type UpdateMultiIssueWorkHoursReqVo struct {
	Input  vo.UpdateMultiIssueWorkHoursReq `json:"input"`
	UserId int64                           `json:"userId"`
	OrgId  int64                           `json:"orgId"`
}

type DeleteIssueWorkHoursReqVo struct {
	Input  vo.DeleteIssueWorkHoursReq `json:"input"`
	UserId int64                      `json:"userId"`
	OrgId  int64                      `json:"orgId"`
}

type DisOrEnableIssueWorkHoursReqVo struct {
	Input  vo.DisOrEnableIssueWorkHoursReq `json:"input"`
	UserId int64                           `json:"userId"`
	OrgId  int64                           `json:"orgId"`
}

type GetIssueWorkHoursInfoReqVo struct {
	Input  vo.GetIssueWorkHoursInfoReq `json:"input"`
	UserId int64                       `json:"userId"`
	OrgId  int64                       `json:"orgId"`
}

type GetIssueWorkHoursInfoRespVo struct {
	vo.Err
	Data *vo.GetIssueWorkHoursInfoResp `json:"data"`
}

type GetWorkHourStatisticReqVo struct {
	Input  vo.GetWorkHourStatisticReq `json:"input"`
	UserId int64                      `json:"userId"`
	OrgId  int64                      `json:"orgId"`
}

type GetWorkHourStatisticRespVo struct {
	vo.Err
	Data *vo.GetWorkHourStatisticResp `json:"data"`
}

type CheckIsIssueMemberReqVo struct {
	Input  vo.CheckIsIssueMemberReq `json:"input"`
	UserId int64                    `json:"userId"`
	OrgId  int64                    `json:"orgId"`
}

type SetUserJoinIssueReqVo struct {
	Input  vo.SetUserJoinIssueReq `json:"input"`
	UserId int64                  `json:"userId"`
	OrgId  int64                  `json:"orgId"`
}

type CheckIsEnableWorkHourReqVo struct {
	Input  vo.CheckIsEnableWorkHourReq `json:"input"`
	UserId int64                       `json:"userId"`
	OrgId  int64                       `json:"orgId"`
}

type CheckIsEnableWorkHourRespVo struct {
	vo.Err
	Data *vo.CheckIsEnableWorkHourResp `json:"data"`
}

type ExportWorkHourStatisticRespVo struct {
	vo.Err
	Data *vo.ExportWorkHourStatisticResp `json:"data"`
}
