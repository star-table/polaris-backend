package trendsvo

import (
	"github.com/star-table/common/core/types"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"time"
)

type TrendListRespVo struct {
	vo.Err
	TrendList *vo.TrendsList `json:"trendList"`
}

type TrendListReqVo struct {
	Input  *vo.TrendReq
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type AddIssueTrendsReqVo struct {
	IssueTrendsBo bo.IssueTrendsBo `json:"issueTrendsBo"`
	Key           string           `json:"key"`
}

type AddProjectTrendsReqVo struct {
	ProjectTrendsBo bo.ProjectTrendsBo `json:"projectTrendsBo"`
}

type AddOrgTrendsReqVo struct {
	OrgTrendsBo bo.OrgTrendsBo `json:"orgTrendsBo"`
}

type CreateTrendsReqVo struct {
	TrendsBo *bo.TrendsBo `json:"trendsBo"`
}

type CreateTrendsRespVo struct {
	TrendsId *int64 `json:"data"`

	vo.Err
}

type GetProjectLatestUpdateTimeReqVo struct {
	OrgId int64                          `json:"orgId"`
	Input GetProjectLatestUpdateTimeData `json:"input"`
}

type GetProjectLatestUpdateTimeData struct {
	ProjectIds []int64 `json:"projectIds"`
}

type GetProjectLatestUpdateTimeRespVo struct {
	vo.Err
	Data map[int64]types.Time `json:"data"`
}

type GetIssueCommentCountReqVo struct {
	OrgId int64                    `json:"orgId"`
	Input GetIssueCommentCountData `json:"input"`
}

type GetIssueCommentCountData struct {
	IssueIds []int64 `json:"issueIds"`
}

type GetIssueCommentCountRespVo struct {
	vo.Err
	Data map[int64]int64 `json:"data"`
}

type GetRecentChangeTimeReqVo struct {
	Data OrgSlice `json:"data"`
}

type OrgSlice struct {
	OrgIds []int64
}

type GetRecentChangeTimeRespVo struct {
	vo.Err
	Data map[int64]time.Time `json:"data"`
}

type GetLatestOperatedProjectIdReq struct {
	OrgId      int64   `json:"orgId"`
	UserId     int64   `json:"userId"`
	ProjectIds []int64 `json:"projectIds"`
}

type GetLatestOperatedProjectIdResp struct {
	vo.Err
	ProjectId int64 `json:"projectId"`
}

type DeleteTrendsReq struct {
	OrgId int64        `json:"orgId"`
	Input DeleteTrends `json:"input"`
}

type DeleteTrends struct {
	ProjectId int64   `json:"projectId"`
	IssueIds  []int64 `json:"issueIds"`
}
