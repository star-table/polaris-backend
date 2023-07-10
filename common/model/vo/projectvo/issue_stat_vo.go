package projectvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type StatisticDailyTaskCompletionProgressReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type IssueAssignRankReqVo struct {
	Input vo.IssueAssignRankReq
	OrgId int64
}

type IssueAndProjectCountStatReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type HandleGroupChatUserInsProIssueReqVo struct {
	Input *HandleGroupChatUserInsProIssueReq `json:"input"`
}

type HandleGroupChatUserInsProIssueReq struct {
	OpenChatId    string `json:"openChatId"`
	OpenId        string `json:"openid"`
	SourceChannel string `json:"sourceChannel"`
}

type HandleGroupChatUserInsProProgressReqVo struct {
	Input *HandleGroupChatUserInsProProgressReq `json:"input"`
}

type UrgeIssueReqVo struct {
	Input         *vo.UrgeIssueReq `json:"input"`
	OrgId         int64            `json:"orgId"`
	UserId        int64            `json:"userId"`
	SourceChannel string           `json:"sourceChannel"`
}

type HandleGroupChatUserInsProProgressReq struct {
	OpenChatId    string `json:"openChatId"`
	OpenId        string `json:"openid"`
	SourceChannel string `json:"sourceChannel"`
}

type HandleGroupChatUserInsProjectSettingsReq struct {
	OpenChatId    string `json:"openChatId"`
	SourceChannel string `json:"sourceChannel"`
}

type HandleGroupChatUserInsAtUserNameReqVo struct {
	Input *HandleGroupChatUserInsAtUserNameReq `json:"input"`
}

type HandleGroupChatUserInsAtUserNameReq struct {
	OpenChatId   string `json:"openChatId"`
	AtUserOpenId string `json:"atUserOpenId"`
	OpUserOpenId string `json:"opUserOpenId"`
}

type HandleGroupChatUserInsAtUserNameWithIssueTitleReqVo struct {
	Input *HandleGroupChatUserInsAtUserNameWithIssueTitleReq `json:"input"`
}

type HandleGroupChatUserInsAtUserNameWithIssueTitleReq struct {
	OpenChatId    string `json:"openChatId"`
	OpUserOpenId  string `json:"opUserOpenId"`
	AtUserOpenId  string `json:"atUserOpenId"`
	IssueTitleStr string `json:"issueTitleStr"`
}

type IssueDailyPersonalWorkCompletionStatReqVo struct {
	Input *vo.IssueDailyPersonalWorkCompletionStatReq `json:"input"`

	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type StatisticDailyTaskCompletionProgressRespVo struct {
	IssueDailyNoticeBo *bo.IssueDailyNoticeBo `json:"data"`
	vo.Err
}

type IssueAssignRankRespVo struct {
	IssueAssignRankResp []*vo.IssueAssignRankInfo `json:"data"`
	vo.Err
}

type IssueAndProjectCountStatRespVo struct {
	Data *vo.IssueAndProjectCountStatResp `json:"data"`

	vo.Err
}

type IssueDailyPersonalWorkCompletionStatRespVo struct {
	Data *vo.IssueDailyPersonalWorkCompletionStatResp `json:"data"`

	vo.Err
}

type GetIssueCountByStatusReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
	StatusId  int64 `json:"statusId"`
}

type GetIssueCountByStatusRespVo struct {
	vo.Err
	Data GetIssueCountByStatusData `json:"data"`
}

type GetIssueCountByStatusData struct {
	Count uint64 `json:"count"`
}

type AuthProjectReqVo struct {
	OrgId     int64  `json:"orgId"`
	UserId    int64  `json:"userId"`
	ProjectId int64  `json:"projectId"`
	Path      string `json:"path"`
	Operation string `json:"operation"`
}

type AuditIssueReq struct {
	OrgId  int64            `json:"orgId"`
	UserId int64            `json:"userId"`
	Params vo.AuditIssueReq `json:"params"`
}

type ViewAuditIssueReq struct {
	OrgId   int64 `json:"orgId"`
	UserId  int64 `json:"userId"`
	IssueId int64 `json:"issueId"`
}

type WithdrawIssueReq struct {
	OrgId   int64 `json:"orgId"`
	UserId  int64 `json:"userId"`
	IssueId int64 `json:"issueId"`
}

type UrgeAuditIssueReq struct {
	OrgId         int64                `json:"orgId"`
	UserId        int64                `json:"userId"`
	SourceChannel string               `json:"sourceChannel"`
	Input         vo.UrgeAuditIssueReq `json:"input"`
}
