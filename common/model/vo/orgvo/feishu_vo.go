package orgvo

import "github.com/star-table/polaris-backend/common/model/vo"

type FeiShuAuthRespVo struct {
	vo.Err
	Auth *vo.FeiShuAuthResp `json:"data"`
}

type FeiShuAuthCodeRespVo struct {
	vo.Err
	Auth *vo.FeiShuAuthCodeResp `json:"data"`
}

type GetFsAccessTokenReqVo struct {
	UserId int64 `json:"userId"`
	OrgId  int64 `json:"orgId"`
}

type GetFsAccessTokenRespVo struct {
	AccessToken string `json:"accessToken"`
	vo.Err
}

type FeishuDeptChangeReqVo struct {
	OrgId         int64  `json:"orgId"`
	DeptOpenId    string `json:"deptOpenId"`
	EventType     string `json:"eventType"`
	SourceChannel string `json:"sourceChannel"`
}

type BoundFsReqVo struct {
	UserId int64             `json:"userId"`
	OrgId  int64             `json:"orgId"`
	Token  string            `json:"token"`
	Input  vo.BoundFeiShuReq `json:"input"`
}

type BoundFsAccountReqVo struct {
	UserId int64                    `json:"userId"`
	OrgId  int64                    `json:"orgId"`
	Token  string                   `json:"token"`
	Input  vo.BoundFeiShuAccountReq `json:"input"`
}

type InitFsAccountReqVo struct {
	Input vo.InitFeiShuAccountReq `json:"input"`
}

type InitFsAccountRespVo struct {
	Auth *vo.FeiShuAuthCodeResp `json:"data"`
	vo.Err
}

type GetJsAPITicketResp struct {
	vo.Err
	Data *vo.GetJsAPITicketResp `json:"data"`
}
