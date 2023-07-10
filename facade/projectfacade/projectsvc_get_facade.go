package projectfacade

import (
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade"
)

var _ = commonvo.DataEvent{}

func AuthCreateProject(req projectvo.AuthCreateProjectReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/authCreateProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ExportIssueTemplate(req projectvo.ExportIssueTemplateReqVo) projectvo.ExportIssueTemplateRespVo {
	respVo := &projectvo.ExportIssueTemplateRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/exportIssueTemplate", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	queryParams["projectObjectTypeId"] = req.ProjectObjectTypeId
	queryParams["iterationId"] = req.IterationId
	queryParams["tableId"] = req.TableId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ExportUserOrDeptSameNameList(req projectvo.ExportUserOrDeptSameNameListReqVo) projectvo.ExportIssueTemplateRespVo {
	respVo := &projectvo.ExportIssueTemplateRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/exportUserOrDeptSameNameList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetCacheProjectInfo(req projectvo.GetCacheProjectInfoReqVo) projectvo.GetCacheProjectInfoRespVo {
	respVo := &projectvo.GetCacheProjectInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getCacheProjectInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["projectId"] = req.ProjectId
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetIssueMembers(req projectvo.GetIssueMembersReqVo) projectvo.GetIssueMembersRespVo {
	respVo := &projectvo.GetIssueMembersRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getIssueMembers", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["issueId"] = req.IssueId
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetNotCompletedIterationBoList(req projectvo.GetNotCompletedIterationBoListReqVo) projectvo.GetNotCompletedIterationBoListRespVo {
	respVo := &projectvo.GetNotCompletedIterationBoListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getNotCompletedIterationBoList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgIssueAndProjectCount(req projectvo.GetOrgIssueAndProjectCountReq) projectvo.GetOrgIssueAndProjectCountResp {
	respVo := &projectvo.GetOrgIssueAndProjectCountResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getOrgIssueAndProjectCount", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectBoListByProjectTypeLangCode(req projectvo.GetProjectBoListByProjectTypeLangCodeReqVo) projectvo.GetProjectBoListByProjectTypeLangCodeRespVo {
	respVo := &projectvo.GetProjectBoListByProjectTypeLangCodeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectBoListByProjectTypeLangCode", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectTypeLangCode"] = req.ProjectTypeLangCode
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectTemplateInner(req projectvo.GetProjectTemplateReq) projectvo.GetProjectTemplateResp {
	respVo := &projectvo.GetProjectTemplateResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectTemplateInner", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserOrDeptSameNameList(req projectvo.GetUserOrDeptSameNameListReq) projectvo.GetUserOrDeptSameNameListRespVo {
	respVo := &projectvo.GetUserOrDeptSameNameListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getUserOrDeptSameNameList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["dataType"] = req.DataType
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IssueAndProjectCountStat(req projectvo.IssueAndProjectCountStatReqVo) projectvo.IssueAndProjectCountStatRespVo {
	respVo := &projectvo.IssueAndProjectCountStatRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueAndProjectCountStat", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IssueReportDetail(req projectvo.IssueReportDetailReqVo) projectvo.IssueReportDetailRespVo {
	respVo := &projectvo.IssueReportDetailRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueReportDetail", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["shareID"] = req.ShareID
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenGetIssueSourceList(req projectvo.OpenGetDemandSourceListReqVo) projectvo.OpenSomeAttrListRespVo {
	respVo := &projectvo.OpenSomeAttrListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openGetIssueSourceList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenGetIterationList(req projectvo.OpenGetIterationListReqVo) projectvo.OpenGetIterationListRespVo {
	respVo := &projectvo.OpenGetIterationListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openGetIterationList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenGetPriorityList(req projectvo.OpenPriorityListReqVo) projectvo.OpenSomeAttrListRespVo {
	respVo := &projectvo.OpenSomeAttrListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openGetPriorityList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenGetProjectObjectTypeList(req projectvo.OpenPriorityListReqVo) projectvo.OpenSomeAttrListRespVo {
	respVo := &projectvo.OpenSomeAttrListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openGetProjectObjectTypeList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenGetPropertyList(req projectvo.OpenGetPropertyListReqVo) projectvo.OpenSomeAttrListRespVo {
	respVo := &projectvo.OpenSomeAttrListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openGetPropertyList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenIssueInfo(req projectvo.IssueInfoReqVo) *projectvo.IssueDetailRespVo {
	respVo := &projectvo.IssueDetailRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openIssueInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["issueID"] = req.IssueID
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["includeDeletedStatus"] = req.IncludeDeletedStatus
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func OrgProjectMember(req projectvo.OrgProjectMemberReqVo) projectvo.OrgProjectMemberListRespVo {
	respVo := &projectvo.OrgProjectMemberListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/orgProjectMember", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ProjectDetail(req projectvo.ProjectDetailReqVo) projectvo.ProjectDetailRespVo {
	respVo := &projectvo.ProjectDetailRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/projectDetail", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["projectId"] = req.ProjectId
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ProjectTypes(req projectvo.ProjectTypesReqVo) projectvo.ProjectTypesRespVo {
	respVo := &projectvo.ProjectTypesRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/projectTypes", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

