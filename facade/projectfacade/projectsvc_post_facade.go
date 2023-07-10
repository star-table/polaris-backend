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

func AddIssueAttachment(req projectvo.AddIssueAttachmentReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/addIssueAttachment", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AddIssueAttachmentFs(req projectvo.AddIssueAttachmentFsReq) projectvo.AddIssueAttachmentFsRespVo {
	respVo := &projectvo.AddIssueAttachmentFsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/addIssueAttachmentFs", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AddProjectChat(req projectvo.UpdateRelateChatReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/addProjectChat", config.GetPreUrl("projectsvc"))
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

func AddProjectMember(req projectvo.RemoveProjectMemberReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/addProjectMember", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AddRecycleRecord(req projectvo.AddRecycleRecordReqVo) projectvo.AddRecycleRecordRespVo {
	respVo := &projectvo.AddRecycleRecordRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/addRecycleRecord", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AppendIterationStat(req projectvo.AppendIterationStatReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/appendIterationStat", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["date"] = req.Date
	requestBody := &req.IterationBo
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AppendProjectDayStat(req projectvo.AppendProjectDayStatReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/appendProjectDayStat", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["date"] = req.Date
	requestBody := &req.ProjectBo
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ApplyProjectTemplateInner(req projectvo.ApplyProjectTemplateReq) projectvo.ApplyProjectTemplateResp {
	respVo := &projectvo.ApplyProjectTemplateResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/applyProjectTemplateInner", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["appId"] = req.AppId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ArchiveProject(req projectvo.ProjectIdReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/archiveProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["projectId"] = req.ProjectId
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AuditIssue(req projectvo.AuditIssueReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/auditIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Params
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AuthProject(req projectvo.AuthProjectReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/authProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["projectId"] = req.ProjectId
	queryParams["path"] = req.Path
	queryParams["operation"] = req.Operation
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AuthProjectPermission(req projectvo.AuthProjectPermissionReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/authProjectPermission", config.GetPreUrl("projectsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func BatchAuditIssue(req *projectvo.BatchAuditIssueReqVo) *vo.DataRespVo {
	respVo := &vo.DataRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/batchAuditIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func BatchCreateIssue(req *projectvo.BatchCreateIssueReqVo) *projectvo.BatchCreateIssueRespVo {
	respVo := &projectvo.BatchCreateIssueRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/batchCreateIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["appId"] = req.AppId
	queryParams["projectId"] = req.ProjectId
	queryParams["tableId"] = req.TableId
	queryParams["beforeDataId"] = req.BeforeDataId
	queryParams["afterDataId"] = req.AfterDataId
	queryParams["isIdGenerated"] = req.IsIdGenerated
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func BatchUpdateIssue(req *projectvo.BatchUpdateIssueReqVo) *vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/batchUpdateIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func BatchUrgeIssue(req *projectvo.BatchUrgeIssueReqVo) *vo.DataRespVo {
	respVo := &vo.DataRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/batchUrgeIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func CancelArchivedProject(req projectvo.ProjectIdReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/cancelArchivedProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["projectId"] = req.ProjectId
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ChangeParentIssue(req projectvo.ChangeParentIssueReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/changeParentIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ChangeProjectChatMember(req projectvo.ChangeProjectMemberReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/changeProjectChatMember", config.GetPreUrl("projectsvc"))
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

func CheckIsEnableWorkHour(req projectvo.CheckIsEnableWorkHourReqVo) projectvo.CheckIsEnableWorkHourRespVo {
	respVo := &projectvo.CheckIsEnableWorkHourRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/checkIsEnableWorkHour", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CheckIsIssueMember(req projectvo.CheckIsIssueMemberReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/checkIsIssueMember", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CheckIsIssueRelatedPeople(req projectvo.CheckIsIssueMemberReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/checkIsIssueRelatedPeople", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CheckIsShowProChatIcon(req projectvo.CheckIsShowProChatIconReq) projectvo.CheckIsShowProChatIconResp {
	respVo := &projectvo.CheckIsShowProChatIconResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/checkIsShowProChatIcon", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CheckShareViewPassword(req *projectvo.CheckShareViewPasswordReq) *projectvo.CheckShareViewPasswordResp {
	respVo := &projectvo.CheckShareViewPasswordResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/checkShareViewPassword", config.GetPreUrl("projectsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func CompleteDelete(req projectvo.CompleteDeleteReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/completeDelete", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ConvertCode(req projectvo.ConvertCodeReqVo) projectvo.ConvertCodeRespVo {
	respVo := &projectvo.ConvertCodeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/convertCode", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ConvertIssueToParent(req projectvo.ConvertIssueToParentReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/convertIssueToParent", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CopyColumn(req projectvo.CopyColumnReqVo) projectvo.CopyColumnRespVo {
	respVo := &projectvo.CopyColumnRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/copyColumn", config.GetPreUrl("projectsvc"))
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

func CopyIssue(req projectvo.CopyIssueReqVo) *projectvo.LcDataListRespVo {
	respVo := &projectvo.LcDataListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/copyIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func CopyIssueBatch(req projectvo.CopyIssueBatchReqVo) projectvo.CopyIssueBatchRespVo {
	respVo := &projectvo.CopyIssueBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/copyIssueBatch", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateColumn(req projectvo.CreateColumnReqVo) projectvo.CreateColumnRespVo {
	respVo := &projectvo.CreateColumnRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createColumn", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateIssue(req *projectvo.CreateIssueReqVo) *projectvo.IssueRespVo {
	respVo := &projectvo.IssueRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["inputAppId"] = req.InputAppId
	queryParams["tableId"] = req.TableId
	requestBody := &req.CreateIssue
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func CreateIssueComment(req projectvo.CreateIssueCommentReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createIssueComment", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateIssueResource(req projectvo.CreateIssueResourceReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createIssueResource", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateIssueSource(req projectvo.CreateIssueSourceReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createIssueSource", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateIssueWorkHours(req projectvo.CreateIssueWorkHoursReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createIssueWorkHours", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateIteration(req projectvo.CreateIterationReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createIteration", config.GetPreUrl("projectsvc"))
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

func CreateMultiIssueWorkHours(req projectvo.CreateMultiIssueWorkHoursReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createMultiIssueWorkHours", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateOrgDirectoryAppsAndViews(req projectvo.CreateOrgDirectoryAppsReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createOrgDirectoryAppsAndViews", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateProject(req projectvo.CreateProjectReqVo) projectvo.ProjectRespVo {
	respVo := &projectvo.ProjectRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["version"] = req.Version
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateProjectDetail(req projectvo.CreateProjectDetailReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createProjectDetail", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateProjectFolder(req projectvo.CreateProjectFolderReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createProjectFolder", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateProjectResource(req projectvo.CreateProjectResourceReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createProjectResource", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateShareView(req *projectvo.CreateShareViewReq) *projectvo.GetShareViewInfoResp {
	respVo := &projectvo.GetShareViewInfoResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createShareView", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func CreateTable(req projectvo.CreateTableReq) projectvo.CreateTableRespVo {
	respVo := &projectvo.CreateTableRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createTable", config.GetPreUrl("projectsvc"))
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

func CreateTaskView(req *projectvo.CreateTaskViewReqVo) projectvo.CreateTaskViewRespVo {
	respVo := &projectvo.CreateTaskViewRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/createTaskView", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteChatCallback(req projectvo.DeleteChatReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteChatCallback", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["outOrgId"] = req.OutOrgId
	queryParams["projectId"] = req.ProjectId
	queryParams["chatId"] = req.ChatId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteColumn(req projectvo.DeleteColumnReqVo) projectvo.DeleteColumnRespVo {
	respVo := &projectvo.DeleteColumnRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteColumn", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteIssue(req projectvo.DeleteIssueReqVo) projectvo.IssueRespVo {
	respVo := &projectvo.IssueRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteIssueBatch(req projectvo.DeleteIssueBatchReqVo) projectvo.DeleteIssueBatchRespVo {
	respVo := &projectvo.DeleteIssueBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteIssueBatch", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["inputAppId"] = req.InputAppId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteIssueResource(req projectvo.DeleteIssueResourceReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteIssueResource", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteIssueSource(req projectvo.DeleteIssueSourceReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteIssueSource", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteIssueWorkHours(req projectvo.DeleteIssueWorkHoursReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteIssueWorkHours", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteIteration(req projectvo.DeleteIterationReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteIteration", config.GetPreUrl("projectsvc"))
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

func DeleteProject(req projectvo.ProjectIdReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["projectId"] = req.ProjectId
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteProjectAttachment(req projectvo.DeleteProjectAttachmentReqVo) projectvo.DeleteProjectAttachmentRespVo {
	respVo := &projectvo.DeleteProjectAttachmentRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteProjectAttachment", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteProjectBatchInner(req projectvo.DeleteProjectInnerReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteProjectBatchInner", config.GetPreUrl("projectsvc"))
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

func DeleteProjectDetail(req projectvo.DeleteProjectDetailReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteProjectDetail", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteProjectFolder(req projectvo.DeleteProjectFolerReqVo) projectvo.DeleteProjectFolerRespVo {
	respVo := &projectvo.DeleteProjectFolerRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteProjectFolder", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteProjectResource(req projectvo.DeleteProjectResourceReqVo) projectvo.DeleteProjectResourceRespVo {
	respVo := &projectvo.DeleteProjectResourceRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteProjectResource", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteShareView(req *projectvo.DeleteShareViewReq) *vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteShareView", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func DeleteTable(req projectvo.DeleteTableReq) projectvo.DeleteTableResp {
	respVo := &projectvo.DeleteTableResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteTable", config.GetPreUrl("projectsvc"))
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

func DeleteTaskView(req *projectvo.DeleteTaskViewReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/deleteTaskView", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DisOrEnableIssueWorkHours(req projectvo.DisOrEnableIssueWorkHoursReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/disOrEnableIssueWorkHours", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DisbandProjectChat(req projectvo.UpdateRelateChatReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/disbandProjectChat", config.GetPreUrl("projectsvc"))
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

func ExportData(req projectvo.ExportIssueReqVo) projectvo.ExportIssueTemplateRespVo {
	respVo := &projectvo.ExportIssueTemplateRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/exportData", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ExportWorkHourStatistic(req projectvo.GetWorkHourStatisticReqVo) projectvo.ExportWorkHourStatisticRespVo {
	respVo := &projectvo.ExportWorkHourStatisticRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/exportWorkHourStatistic", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func FsChatDisbandCallback(req projectvo.FsChatDisbandCallbackReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/fsChatDisbandCallback", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["chatId"] = req.ChatId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetFsProjectChatPushSettings(req projectvo.GetFsProjectChatPushSettingsReq) projectvo.GetFsProjectChatPushSettingsResp {
	respVo := &projectvo.GetFsProjectChatPushSettingsResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getFsProjectChatPushSettings", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["chatId"] = req.ChatId
	queryParams["projectId"] = req.ProjectId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetIssueInfo(req projectvo.GetIssueInfoReqVo) projectvo.GetIssueInfoRespVo {
	respVo := &projectvo.GetIssueInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getIssueInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["issueId"] = req.IssueId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetIssueInfoByDataIdsList(req projectvo.IssueInfoListByDataIdsReqVo) projectvo.IssueInfoListByDataIdsRespVo {
	respVo := &projectvo.IssueInfoListByDataIdsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getIssueInfoByDataIdsList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.DataIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetIssueInfoList(req projectvo.IssueInfoListReqVo) projectvo.IssueInfoListRespVo {
	respVo := &projectvo.IssueInfoListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getIssueInfoList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.IssueIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetIssueLinks(req projectvo.GetIssueLinksReqVo) projectvo.GetIssueLinksRespVo {
	respVo := &projectvo.GetIssueLinksRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getIssueLinks", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["orgId"] = req.OrgId
	queryParams["issueId"] = req.IssueId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetIssueRelationResource(req projectvo.GetIssueRelationResourceReqVo) projectvo.GetIssueRelationResourceRespVo {
	respVo := &projectvo.GetIssueRelationResourceRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getIssueRelationResource", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetIssueRowList(req projectvo.IssueRowListReq) projectvo.IssueRowListResp {
	respVo := &projectvo.IssueRowListResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getIssueRowList", config.GetPreUrl("projectsvc"))
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

func GetIssueWorkHoursInfo(req projectvo.GetIssueWorkHoursInfoReqVo) projectvo.GetIssueWorkHoursInfoRespVo {
	respVo := &projectvo.GetIssueWorkHoursInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getIssueWorkHoursInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetIssueWorkHoursList(req projectvo.GetIssueWorkHoursListReqVo) projectvo.GetIssueWorkHoursListRespVo {
	respVo := &projectvo.GetIssueWorkHoursListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getIssueWorkHoursList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetLcIssueInfoBatch(req projectvo.GetLcIssueInfoBatchReqVo) projectvo.GetLcIssueInfoBatchRespVo {
	respVo := &projectvo.GetLcIssueInfoBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getLcIssueInfoBatch", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.IssueIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetMenu(req projectvo.GetMenuReqVo) projectvo.GetMenuRespVo {
	respVo := &projectvo.GetMenuRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getMenu", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["appId"] = req.AppId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOneTableColumns(req projectvo.GetTableColumnReq) projectvo.TableColumnsResp {
	respVo := &projectvo.TableColumnsResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getOneTableColumns", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["projectId"] = req.ProjectId
	queryParams["tableId"] = req.TableId
	queryParams["notAllIssue"] = req.NotAllIssue
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectAttachment(req projectvo.GetProjectAttachmentReqVo) projectvo.GetProjectAttachmentRespVo {
	respVo := &projectvo.GetProjectAttachmentRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectAttachment", config.GetPreUrl("projectsvc"))
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

func GetProjectAttachmentInfo(req projectvo.GetProjectAttachmentInfoReqVo) projectvo.GetProjectAttachmentInfoRespVo {
	respVo := &projectvo.GetProjectAttachmentInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectAttachmentInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectDetails(req projectvo.GetSimpleProjectInfoReqVo) projectvo.GetProjectDetailsRespVo {
	respVo := &projectvo.GetProjectDetailsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectDetails", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Ids
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectFolder(req projectvo.GetProjectFolderReqVo) projectvo.GetProjectFolderRespVo {
	respVo := &projectvo.GetProjectFolderRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectFolder", config.GetPreUrl("projectsvc"))
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

func GetProjectIdByChatId(req projectvo.GetProjectIdByChatIdReqVo) projectvo.GetProjectIdByChatIdRespVo {
	respVo := &projectvo.GetProjectIdByChatIdRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectIdByChatId", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["openChatId"] = req.OpenChatId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectIdsByChatId(req projectvo.GetProjectIdsByChatIdReqVo) projectvo.GetProjectIdsByChatIdRespVo {
	respVo := &projectvo.GetProjectIdsByChatIdRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectIdsByChatId", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["openChatId"] = req.OpenChatId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectInfoByOrgIds(req projectvo.GetProjectInfoListByOrgIdsReqVo) projectvo.GetProjectInfoListByOrgIdsListRespVo {
	respVo := &projectvo.GetProjectInfoListByOrgIdsListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectInfoByOrgIds", config.GetPreUrl("projectsvc"))
	requestBody := &req.OrgIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectMainChatId(req projectvo.GetProjectMainChatIdReq) projectvo.GetProjectMainChatIdResp {
	respVo := &projectvo.GetProjectMainChatIdResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectMainChatId", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectMemberIds(req projectvo.GetProjectMemberIdsReqVo) projectvo.GetProjectMemberIdsResp {
	respVo := &projectvo.GetProjectMemberIdsResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectMemberIds", config.GetPreUrl("projectsvc"))
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

func GetProjectRelation(req projectvo.GetProjectRelationReqVo) projectvo.GetProjectRelationRespVo {
	respVo := &projectvo.GetProjectRelationRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectRelation", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["projectId"] = req.ProjectId
	requestBody := &req.RelationType
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectRelationBatch(req projectvo.GetProjectRelationBatchReqVo) projectvo.GetProjectRelationBatchRespVo {
	respVo := &projectvo.GetProjectRelationBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectRelationBatch", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectRelationUserIds(req projectvo.GetProjectRelationUserIdsReq) projectvo.GetProjectRelationUserIdsResp {
	respVo := &projectvo.GetProjectRelationUserIdsResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectRelationUserIds", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["projectId"] = req.ProjectId
	queryParams["relationType"] = req.RelationType
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectResource(req projectvo.GetProjectResourceReqVo) projectvo.GetProjectResourceResVo {
	respVo := &projectvo.GetProjectResourceResVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectResource", config.GetPreUrl("projectsvc"))
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

func GetProjectResourceInfo(req projectvo.GetProjectResourceInfoReqVo) projectvo.GetProjectResourceInfoRespVo {
	respVo := &projectvo.GetProjectResourceInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectResourceInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectStatistics(req projectvo.GetProjectStatisticsReqVo) projectvo.GetProjectStatisticsResp {
	respVo := &projectvo.GetProjectStatisticsResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getProjectStatistics", config.GetPreUrl("projectsvc"))
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

func GetRecycleList(req projectvo.GetRecycleListReqVo) projectvo.GetRecycleListRespVo {
	respVo := &projectvo.GetRecycleListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getRecycleList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetShareViewInfo(req *projectvo.GetShareViewInfoReq) *projectvo.GetShareViewInfoResp {
	respVo := &projectvo.GetShareViewInfoResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getShareViewInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetShareViewInfoByKey(req *projectvo.GetShareViewInfoByKeyReq) *projectvo.GetShareViewInfoResp {
	respVo := &projectvo.GetShareViewInfoResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getShareViewInfoByKey", config.GetPreUrl("projectsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetSimpleIterationInfo(req projectvo.GetSimpleIterationInfoReqVo) projectvo.GetSimpleIterationInfoRespVo {
	respVo := &projectvo.GetSimpleIterationInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getSimpleIterationInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.IterationIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetSimpleProjectInfo(req projectvo.GetSimpleProjectInfoReqVo) projectvo.GetSimpleProjectInfoRespVo {
	respVo := &projectvo.GetSimpleProjectInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getSimpleProjectInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Ids
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetSimpleProjectsByOrgId(req projectvo.GetSimpleProjectsByOrgIdReq) projectvo.GetSimpleProjectsByOrgIdResp {
	respVo := &projectvo.GetSimpleProjectsByOrgIdResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getSimpleProjectsByOrgId", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetTable(req projectvo.GetTableInfoReq) projectvo.GetTableInfoResp {
	respVo := &projectvo.GetTableInfoResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getTable", config.GetPreUrl("projectsvc"))
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

func GetTableStatusByTableId(req projectvo.GetTableStatusReq) projectvo.GetTableStatusResp {
	respVo := &projectvo.GetTableStatusResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getTableStatusByTableId", config.GetPreUrl("projectsvc"))
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

func GetTables(req projectvo.GetTablesReqVo) projectvo.GetTablesDataResp {
	respVo := &projectvo.GetTablesDataResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getTables", config.GetPreUrl("projectsvc"))
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

func GetTablesByApps(req projectvo.ReadTablesByAppsReqVo) projectvo.ReadTablesByAppsRespVo {
	respVo := &projectvo.ReadTablesByAppsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getTablesByApps", config.GetPreUrl("projectsvc"))
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

func GetTablesByOrg(req projectvo.GetTablesByOrgReq) projectvo.GetTablesByOrgRespVo {
	respVo := &projectvo.GetTablesByOrgRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getTablesByOrg", config.GetPreUrl("projectsvc"))
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

func GetTablesColumns(req projectvo.GetTablesColumnsReq) projectvo.TablesColumnsResp {
	respVo := &projectvo.TablesColumnsResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getTablesColumns", config.GetPreUrl("projectsvc"))
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

func GetTaskViewList(req *projectvo.GetTaskViewListReqVo) projectvo.GetTaskViewListRespVo {
	respVo := &projectvo.GetTaskViewListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getTaskViewList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetTrendsMembers(req projectvo.GetTrendListMembersReqVo) projectvo.GetTrendListMembersResp {
	respVo := &projectvo.GetTrendListMembersResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getTrendsMembers", config.GetPreUrl("projectsvc"))
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

func GetWorkHourStatistic(req projectvo.GetWorkHourStatisticReqVo) projectvo.GetWorkHourStatisticRespVo {
	respVo := &projectvo.GetWorkHourStatisticRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/getWorkHourStatistic", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func HandleGroupChatUserInsAtUserName(req projectvo.HandleGroupChatUserInsAtUserNameReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/handleGroupChatUserInsAtUserName", config.GetPreUrl("projectsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func HandleGroupChatUserInsAtUserNameWithIssueTitle(req projectvo.HandleGroupChatUserInsAtUserNameWithIssueTitleReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/handleGroupChatUserInsAtUserNameWithIssueTitle", config.GetPreUrl("projectsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func HandleGroupChatUserInsProIssue(req projectvo.HandleGroupChatUserInsProIssueReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/handleGroupChatUserInsProIssue", config.GetPreUrl("projectsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func HandleGroupChatUserInsProProgress(req projectvo.HandleGroupChatUserInsProProgressReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/handleGroupChatUserInsProProgress", config.GetPreUrl("projectsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func HandleGroupChatUserInsProjectSettings(req projectvo.HandleGroupChatUserInsProjectSettingsReq) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/handleGroupChatUserInsProjectSettings", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["openChatId"] = req.OpenChatId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ImportIssues(req projectvo.ImportIssuesReqVo) projectvo.ImportIssuesRespVo {
	respVo := &projectvo.ImportIssuesRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/importIssues", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func InnerCreateTodoHook(req *projectvo.InnerCreateTodoHookReq) *vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/innerCreateTodoHook", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func InnerIssueCreate(req *projectvo.InnerIssueCreateReq) *projectvo.LcDataListRespVo {
	respVo := &projectvo.LcDataListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/innerIssueCreate", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func InnerIssueCreateByCopy(req *projectvo.InnerIssueCreateByCopyReq) *projectvo.LcDataListRespVo {
	respVo := &projectvo.LcDataListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/innerIssueCreateByCopy", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func InnerIssueFilter(req *projectvo.InnerIssueFilterReq) string {
	respVo := ""
	reqUrl := fmt.Sprintf("%s/api/projectsvc/innerIssueFilter", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, &respVo)
	if err.Failure() {
		return facade.ToJsonString(err)
	}
	return respVo
}

func InnerIssueUpdate(req *projectvo.InnerIssueUpdateReq) *vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/innerIssueUpdate", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func IssueAssignRank(req projectvo.IssueAssignRankReqVo) projectvo.IssueAssignRankRespVo {
	respVo := &projectvo.IssueAssignRankRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueAssignRank", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IssueDailyPersonalWorkCompletionStat(req projectvo.IssueDailyPersonalWorkCompletionStatReqVo) projectvo.IssueDailyPersonalWorkCompletionStatRespVo {
	respVo := &projectvo.IssueDailyPersonalWorkCompletionStatRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueDailyPersonalWorkCompletionStat", config.GetPreUrl("projectsvc"))
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

func IssueListSimpleByDataIds(req projectvo.GetIssueListSimpleByDataIdsReqVo) projectvo.SimpleIssueListRespVo {
	respVo := &projectvo.SimpleIssueListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueListSimpleByDataIds", config.GetPreUrl("projectsvc"))
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

func IssueListSimpleByTableIds(req projectvo.GetIssueListWithConditionsReqVo) projectvo.IssueListWithConditionsResp {
	respVo := &projectvo.IssueListWithConditionsResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueListSimpleByTableIds", config.GetPreUrl("projectsvc"))
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

func IssueListStat(req projectvo.IssueListStatReq) projectvo.IssueListStatResp {
	respVo := &projectvo.IssueListStatResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueListStat", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IssueResources(req projectvo.IssueResourcesReqVo) projectvo.IssueResourcesRespVo {
	respVo := &projectvo.IssueResourcesRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueResources", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IssueShareCard(req projectvo.IssueCardShareReqVo) projectvo.IssueCardShareResp {
	respVo := &projectvo.IssueCardShareResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueShareCard", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IssueSourceList(req projectvo.IssueSourceListReqVo) projectvo.IssueSourceListRespVo {
	respVo := &projectvo.IssueSourceListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueSourceList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IssueStartChat(req projectvo.IssueStartChatReqVo) projectvo.IssueStartChatRespVo {
	respVo := &projectvo.IssueStartChatRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueStartChat", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IssueStatusTypeStat(req projectvo.IssueStatusTypeStatReqVo) projectvo.IssueStatusTypeStatRespVo {
	respVo := &projectvo.IssueStatusTypeStatRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueStatusTypeStat", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IssueStatusTypeStatDetail(req projectvo.IssueStatusTypeStatReqVo) projectvo.IssueStatusTypeStatDetailRespVo {
	respVo := &projectvo.IssueStatusTypeStatDetailRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/issueStatusTypeStatDetail", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IterationInfo(req projectvo.IterationInfoReqVo) projectvo.IterationInfoRespVo {
	respVo := &projectvo.IterationInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/iterationInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IterationList(req projectvo.IterationListReqVo) projectvo.IterationListRespVo {
	respVo := &projectvo.IterationListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/iterationList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IterationStats(req projectvo.IterationStatsReqVo) projectvo.IterationStatsRespVo {
	respVo := &projectvo.IterationStatsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/iterationStats", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func IterationStatusTypeStat(req projectvo.IterationStatusTypeStatReqVo) projectvo.IterationStatusTypeStatRespVo {
	respVo := &projectvo.IterationStatusTypeStatRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/iterationStatusTypeStat", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func LcHomeIssues(req *projectvo.HomeIssuesReqVo) string {
	respVo := ""
	reqUrl := fmt.Sprintf("%s/api/projectsvc/lcHomeIssues", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, &respVo)
	if err.Failure() {
		return facade.ToJsonString(err)
	}
	return respVo
}

func LcHomeIssuesForAll(req *projectvo.HomeIssuesReqVo) string {
	respVo := ""
	reqUrl := fmt.Sprintf("%s/api/projectsvc/lcHomeIssuesForAll", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, &respVo)
	if err.Failure() {
		return facade.ToJsonString(err)
	}
	return respVo
}

func LcHomeIssuesForIssue(req *projectvo.IssueDetailReqVo) *projectvo.IssueDetailRespVo {
	respVo := &projectvo.IssueDetailRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/lcHomeIssuesForIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["appId"] = req.AppId
	queryParams["tableId"] = req.TableId
	queryParams["issueId"] = req.IssueId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func LcViewStatForAll(req *projectvo.LcViewStatReqVo) *projectvo.LcViewStatRespVo {
	respVo := &projectvo.LcViewStatRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/lcViewStatForAll", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func MirrorsStat(req *projectvo.MirrorsStatReq) projectvo.MirrorsStatResp {
	respVo := &projectvo.MirrorsStatResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/mirrorsStat", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func MoveIssue(req projectvo.MoveIssueReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/moveIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["inputAppId"] = req.InputAppId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func MoveIssueBatch(req projectvo.MoveIssueBatchReqVo) projectvo.MoveIssueBatchRespVo {
	respVo := &projectvo.MoveIssueBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/moveIssueBatch", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenCreateIssue(req projectvo.CreateIssueReqVo) projectvo.IssueRespVo {
	respVo := &projectvo.IssueRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openCreateIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["inputAppId"] = req.InputAppId
	queryParams["tableId"] = req.TableId
	requestBody := &req.CreateIssue
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenCreateIssueComment(req projectvo.CreateIssueCommentReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openCreateIssueComment", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenCreateProject(req projectvo.CreateProjectReqVo) projectvo.ProjectRespVo {
	respVo := &projectvo.ProjectRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openCreateProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["version"] = req.Version
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenDeleteIssue(req projectvo.OpenDeleteIssueReqVo) projectvo.DeleteIssueRespVo {
	respVo := &projectvo.DeleteIssueRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openDeleteIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenDeleteProject(req projectvo.ProjectIdReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openDeleteProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["projectId"] = req.ProjectId
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenIssueList(req projectvo.OpenIssueListReqVo) projectvo.LcHomeIssuesRespVo {
	respVo := &projectvo.LcHomeIssuesRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openIssueList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenIssueStatusTypeStat(req projectvo.IssueStatusTypeStatReqVo) projectvo.IssueStatusTypeStatRespVo {
	respVo := &projectvo.IssueStatusTypeStatRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openIssueStatusTypeStat", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenLcHomeIssues(req *projectvo.HomeIssuesReqVo) string {
	respVo := ""
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openLcHomeIssues", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, &respVo)
	if err.Failure() {
		return facade.ToJsonString(err)
	}
	return respVo
}

func OpenLcHomeIssuesForAll(req *projectvo.HomeIssuesReqVo) string {
	respVo := ""
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openLcHomeIssuesForAll", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, &respVo)
	if err.Failure() {
		return facade.ToJsonString(err)
	}
	return respVo
}

func OpenProjectInfo(req projectvo.ProjectInfoReqVo) projectvo.ProjectInfoRespVo {
	respVo := &projectvo.ProjectInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openProjectInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenProjects(req projectvo.ProjectsRepVo) projectvo.ProjectsRespVo {
	respVo := &projectvo.ProjectsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openProjects", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.ProjectExtraBody
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenUpdateIssue(req projectvo.OpenUpdateIssueReqVo) projectvo.UpdateIssueRespVo {
	respVo := &projectvo.UpdateIssueRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openUpdateIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenUpdateProject(req projectvo.UpdateProjectReqVo) projectvo.ProjectRespVo {
	respVo := &projectvo.ProjectRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/openUpdateProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OrgProjectMemberList(req projectvo.OrgProjectMemberListReq) projectvo.OrgProjectMemberListResp {
	respVo := &projectvo.OrgProjectMemberListResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/orgProjectMemberList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	requestBody := &req.Params
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PayLimitNum(req projectvo.PayLimitNumReq) projectvo.PayLimitNumResp {
	respVo := &projectvo.PayLimitNumResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/payLimitNum", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PayLimitNumForRest(req projectvo.PayLimitNumReq) projectvo.PayLimitNumForRestResp {
	respVo := &projectvo.PayLimitNumForRestResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/payLimitNumForRest", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ProjectChatList(req projectvo.ProjectChatListReqVo) projectvo.ProjectChatListRespVo {
	respVo := &projectvo.ProjectChatListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/projectChatList", config.GetPreUrl("projectsvc"))
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

func ProjectDayStats(req *projectvo.ProjectDayStatsReqVo) projectvo.ProjectDayStatsRespVo {
	respVo := &projectvo.ProjectDayStatsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/projectDayStats", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Params
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ProjectInfo(req projectvo.ProjectInfoReqVo) projectvo.ProjectInfoRespVo {
	respVo := &projectvo.ProjectInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/projectInfo", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ProjectInit(req projectvo.ProjectInitReqVo) projectvo.ProjectInitRespVo {
	respVo := &projectvo.ProjectInitRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/projectInit", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ProjectIssueRelatedStatus(req projectvo.ProjectIssueRelatedStatusReqVo) projectvo.ProjectIssueRelatedStatusRespVo {
	respVo := &projectvo.ProjectIssueRelatedStatusRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/projectIssueRelatedStatus", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ProjectMemberIdList(req projectvo.ProjectMemberIdListReq) projectvo.ProjectMemberIdListResp {
	respVo := &projectvo.ProjectMemberIdListResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/projectMemberIdList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ProjectUserList(req projectvo.ProjectUserListReq) projectvo.ProjectUserListRespVo {
	respVo := &projectvo.ProjectUserListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/projectUserList", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ProjectUserListForFuse(req projectvo.ProjectUserListReq) projectvo.ProjectUserListForFuseRespVo {
	respVo := &projectvo.ProjectUserListForFuseRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/projectUserListForFuse", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func Projects(req projectvo.ProjectsRepVo) projectvo.ProjectsRespVo {
	respVo := &projectvo.ProjectsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/projects", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.ProjectExtraBody
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func QueryProcessForAsyncTask(req projectvo.QueryProcessForAsyncTaskReqVo) projectvo.QueryProcessForAsyncTaskRespVo {
	respVo := &projectvo.QueryProcessForAsyncTaskRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/queryProcessForAsyncTask", config.GetPreUrl("projectsvc"))
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

func QuitProject(req projectvo.ProjectIdReqVo) projectvo.QuitProjectRespVo {
	respVo := &projectvo.QuitProjectRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/quitProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["projectId"] = req.ProjectId
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func RecoverRecycleBin(req projectvo.CompleteDeleteReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/recoverRecycleBin", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func RemoveProjectMember(req projectvo.RemoveProjectMemberReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/removeProjectMember", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func RenameTable(req projectvo.RenameTableReq) projectvo.RenameTableResp {
	respVo := &projectvo.RenameTableResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/renameTable", config.GetPreUrl("projectsvc"))
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

func ReportAppEvent(req *commonvo.ReportAppEventReq) *vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/reportAppEvent", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["eventType"] = req.EventType
	queryParams["traceId"] = req.TraceId
	requestBody := &req.AppEvent
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func ReportTableEvent(req *commonvo.ReportTableEventReq) *vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/reportTableEvent", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["eventType"] = req.EventType
	queryParams["traceId"] = req.TraceId
	requestBody := &req.TableEvent
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func ResetShareKey(req *projectvo.ResetShareKeyReq) *projectvo.GetShareViewInfoResp {
	respVo := &projectvo.GetShareViewInfoResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/resetShareKey", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func SaveMenu(req projectvo.SaveMenuReqVo) projectvo.SaveMenuRespVo {
	respVo := &projectvo.SaveMenuRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/saveMenu", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SetAutoSchedule(req projectvo.SetAutoScheduleReq) projectvo.SetAutoScheduleResp {
	respVo := &projectvo.SetAutoScheduleResp{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/setAutoSchedule", config.GetPreUrl("projectsvc"))
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

func SetUserJoinIssue(req projectvo.SetUserJoinIssueReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/setUserJoinIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func StarProject(req projectvo.ProjectIdReqVo) projectvo.OperateProjectRespVo {
	respVo := &projectvo.OperateProjectRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/starProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["projectId"] = req.ProjectId
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func TodoUrge(req *projectvo.TodoUrgeReq) *vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/todoUrge", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func UnrelatedChatList(req projectvo.UnrelatedChatListReqVo) projectvo.ProjectChatListRespVo {
	respVo := &projectvo.ProjectChatListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/unrelatedChatList", config.GetPreUrl("projectsvc"))
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

func UnstarProject(req projectvo.ProjectIdReqVo) projectvo.OperateProjectRespVo {
	respVo := &projectvo.OperateProjectRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/unstarProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["projectId"] = req.ProjectId
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateColumn(req projectvo.UpdateColumnReqVo) projectvo.UpdateColumnRespVo {
	respVo := &projectvo.UpdateColumnRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateColumn", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateColumnDescription(req projectvo.UpdateColumnDescriptionReqVo) projectvo.UpdateColumnDescriptionRespVo {
	respVo := &projectvo.UpdateColumnDescriptionRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateColumnDescription", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateFsProjectChatPushSettings(req projectvo.UpdateFsProjectChatPushSettingsReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateFsProjectChatPushSettings", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateIssueSource(req projectvo.UpdateIssueSourceReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateIssueSource", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateIssueWorkHours(req projectvo.UpdateIssueWorkHoursReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateIssueWorkHours", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateIteration(req projectvo.UpdateIterationReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateIteration", config.GetPreUrl("projectsvc"))
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

func UpdateIterationSort(req projectvo.UpdateIterationSortReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateIterationSort", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Params
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateIterationStatus(req projectvo.UpdateIterationStatusReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateIterationStatus", config.GetPreUrl("projectsvc"))
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

func UpdateIterationStatusTime(req projectvo.UpdateIterationStatusTimeReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateIterationStatusTime", config.GetPreUrl("projectsvc"))
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

func UpdateMultiIssueWorkHours(req projectvo.UpdateMultiIssueWorkHoursReqVo) projectvo.BoolRespVo {
	respVo := &projectvo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateMultiIssueWorkHours", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateProject(req projectvo.UpdateProjectReqVo) projectvo.ProjectRespVo {
	respVo := &projectvo.ProjectRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateProject", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateProjectDetail(req projectvo.UpdateProjectDetailReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateProjectDetail", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateProjectFileResource(req projectvo.UpdateProjectFileResourceReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateProjectFileResource", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateProjectFolder(req projectvo.UpdateProjectFolderReqVo) projectvo.UpdateProjectFolderRespVo {
	respVo := &projectvo.UpdateProjectFolderRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateProjectFolder", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateProjectResourceFolder(req projectvo.UpdateProjectResourceFolderReqVo) projectvo.UpdateProjectResourceFolderRespVo {
	respVo := &projectvo.UpdateProjectResourceFolderRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateProjectResourceFolder", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateProjectResourceName(req projectvo.UpdateProjectResourceNameReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateProjectResourceName", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateProjectStatus(req projectvo.UpdateProjectStatusReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateProjectStatus", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateShareConfig(req *projectvo.UpdateShareConfigReq) *vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateShareConfig", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func UpdateSharePassword(req *projectvo.UpdateSharePasswordReq) *vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateSharePassword", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func UpdateTaskView(req *projectvo.UpdateTaskViewReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/updateTaskView", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UrgeAuditIssue(req projectvo.UrgeAuditIssueReq) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/urgeAuditIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UrgeIssue(req projectvo.UrgeIssueReqVo) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/urgeIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ViewAuditIssue(req projectvo.ViewAuditIssueReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/viewAuditIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["issueId"] = req.IssueId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func WithdrawIssue(req projectvo.WithdrawIssueReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/projectsvc/withdrawIssue", config.GetPreUrl("projectsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["issueId"] = req.IssueId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

