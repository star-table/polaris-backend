package trendsvo

import (
	"github.com/star-table/polaris-backend/common/model/vo"
)

type UnreadNoticeCountReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type UnreadNoticeCountRespVo struct {
	vo.Err
	Count int64 `json:"count"`
}

type CallMeCountReqVo struct {
	OrgId     int64  `json:"orgId"`
	UserId    int64  `json:"userId"`
	ProjectId *int64 `json:"projectId"`
	IssueId   *int64 `json:"issueId"`
}

type NoticeListReqVo struct {
	UserId int64             `json:"userId"`
	OrgId  int64             `json:"orgId"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Input  *vo.NoticeListReq `json:"input"`
}

type NoticeListRespVo struct {
	vo.Err
	Data *vo.NoticeList `json:"data"`
}
