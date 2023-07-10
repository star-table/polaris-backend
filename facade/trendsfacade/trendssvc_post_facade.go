package trendsfacade

import (
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/trendsvo"
	"github.com/star-table/polaris-backend/facade"
)

func AddIssueTrends(req trendsvo.AddIssueTrendsReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/addIssueTrends", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["key"] = req.Key
	requestBody := &req.IssueTrendsBo
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AddOrgTrends(req trendsvo.AddOrgTrendsReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/addOrgTrends", config.GetPreUrl("trendssvc"))
	requestBody := &req.OrgTrendsBo
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AddProjectTrends(req trendsvo.AddProjectTrendsReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/addProjectTrends", config.GetPreUrl("trendssvc"))
	requestBody := &req.ProjectTrendsBo
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CallMeCount(req trendsvo.CallMeCountReqVo) trendsvo.UnreadNoticeCountRespVo {
	respVo := &trendsvo.UnreadNoticeCountRespVo{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/callMeCount", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["projectId"] = req.ProjectId
	queryParams["issueId"] = req.IssueId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateComment(req trendsvo.CreateCommentReqVo) trendsvo.CreateCommentRespVo {
	respVo := &trendsvo.CreateCommentRespVo{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/createComment", config.GetPreUrl("trendssvc"))
	requestBody := &req.CommentBo
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateTrends(req trendsvo.CreateTrendsReqVo) trendsvo.CreateTrendsRespVo {
	respVo := &trendsvo.CreateTrendsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/createTrends", config.GetPreUrl("trendssvc"))
	requestBody := &req.TrendsBo
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetIssueCommentCount(req trendsvo.GetIssueCommentCountReqVo) trendsvo.GetIssueCommentCountRespVo {
	respVo := &trendsvo.GetIssueCommentCountRespVo{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/getIssueCommentCount", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetLatestOperatedProjectId(req trendsvo.GetLatestOperatedProjectIdReq) trendsvo.GetLatestOperatedProjectIdResp {
	respVo := &trendsvo.GetLatestOperatedProjectIdResp{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/getLatestOperatedProjectId", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.ProjectIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectLatestUpdateTime(req trendsvo.GetProjectLatestUpdateTimeReqVo) trendsvo.GetProjectLatestUpdateTimeRespVo {
	respVo := &trendsvo.GetProjectLatestUpdateTimeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/getProjectLatestUpdateTime", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetRecentChangeTime(req trendsvo.GetRecentChangeTimeReqVo) trendsvo.GetRecentChangeTimeRespVo {
	respVo := &trendsvo.GetRecentChangeTimeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/getRecentChangeTime", config.GetPreUrl("trendssvc"))
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func NoticeList(req trendsvo.NoticeListReqVo) trendsvo.NoticeListRespVo {
	respVo := &trendsvo.NoticeListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/noticeList", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func TrendList(req trendsvo.TrendListReqVo) trendsvo.TrendListRespVo {
	respVo := &trendsvo.TrendListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/trendList", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UnreadNoticeCount(req trendsvo.UnreadNoticeCountReqVo) trendsvo.UnreadNoticeCountRespVo {
	respVo := &trendsvo.UnreadNoticeCountRespVo{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/unreadNoticeCount", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteTrends(req trendsvo.DeleteTrendsReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/trendssvc/deleteTrends", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}
