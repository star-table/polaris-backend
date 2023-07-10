package handler

import (
	"strconv"

	"github.com/star-table/go-common/pkg/encoding"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type issueListHandler struct{}

var IssueListHandler issueListHandler

//Title = -1, // 标题
//Number_ID = -2, // 编号ID
//Owner = -3, // 负责人
//Status = -4, //状态，
//Plan_EndTime = -6, //截止时间
//Priority = -7, // 优先级
//Follower = -8, // 关注人
//Plan_StartTime = -9, // 开始时间
//Iteration = -10, // 迭代
//Task_Column = -11, // 任务栏
//Demand_Type = -12, // 需求类型
//Demand_Source = -13, // 需求来源
//Bug_Type = -14, // 缺陷类型
//Severity = -15, // 严重程度
//Tag = -16, // 标签
//Auditor = -26, //确认人

// 性质、缺陷、类型 直接作为自定义字段
var commonFields = []string{
	"title",
	"code",
	"ownerId",
	"issueStatus",
	"planStartTime",
	"planEndTime",
	"remark",
	"remarkDetail",
	"mentionUserIds",
	"iterationId",
	"followerIds",
	"followerDeptIds",
	"auditorIds",
	"sort",
}

// @Security ApiKeyAuth
// @Summary 任务列表
// @Description 任务列表
// @Tags 任务
// @accept application/json
// @Produce application/json
// @Param projectId path int64 true "项目id"
// @Param input body vo.HomeIssuesRestReq true "入参"
// @Success 200 {object} vo.HomeIssueInfoResp
// @Failure 400
// @Router /api/rest/project/{projectId}/issue/filter [post]
func (issueListHandler) Filter(c *gin.Context) {
	cacheUserInfo, sysErr := GetCacheUserInfo(c)
	if sysErr != nil {
		Fail(c, sysErr)
		return
	}

	req := &projectvo.HomeIssueInfoReq{}
	if err := c.BindJSON(req); err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	zero := 0
	if req.Page == nil {
		req.Page = &zero
	}
	if req.Size == nil {
		req.Size = &zero
	}

	projectId, err := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}
	if projectId < 0 {
		projectId = 0
	}
	req.ProjectID = &projectId

	if req.TableID == nil {
		temp := "0"
		req.TableID = &temp
	}

	params := &projectvo.HomeIssuesReqVo{
		Page:   *req.Page,
		Size:   *req.Size,
		Input:  req,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	}

	withoutInfo := c.Query("withoutInfo")
	if withoutInfo == "4" {
		// 主页菜单任务列表
		respVo := projectfacade.LcHomeIssuesForAll(params)
		SuccessJson(c, respVo)
	} else {
		// 项目任务列表
		respVo := projectfacade.LcHomeIssues(params)
		SuccessJson(c, respVo)
	}
}

func (issueListHandler) FilterForIssue(c *gin.Context) {
	cacheUserInfo, sysErr := GetCacheUserInfo(c)
	if sysErr != nil {
		Fail(c, sysErr)
		return
	}

	req := vo.IssueDetailRestReq{}
	if err := c.BindJSON(&req); err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}
	if req.IssueId <= 0 {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	var appId, tableId int64
	if req.AppId != nil {
		appId = cast.ToInt64(*req.AppId)
	}
	if req.TableId != nil {
		tableId = cast.ToInt64(*req.TableId)
	}

	respVo := projectfacade.LcHomeIssuesForIssue(&projectvo.IssueDetailReqVo{
		OrgId:   cacheUserInfo.OrgId,
		UserId:  cacheUserInfo.UserId,
		AppId:   appId,
		TableId: tableId,
		IssueId: int64(req.IssueId),
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Success(c, respVo.Data)
	}
}

func (h issueListHandler) BatchAuditIssue(c *gin.Context) {
	cacheUserInfo, errSys := GetCacheUserInfo(c)
	if errSys != nil {
		Fail(c, errSys)
		return
	}
	req := vo.LessBatchAuditIssueReq{}
	err := c.BindJSON(&req)
	if err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	resp := projectfacade.BatchAuditIssue(&projectvo.BatchAuditIssueReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input:  &req,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		Fail(c, resp.Error())
		return
	}
	Success(c, resp.Data)
}

func (h issueListHandler) BatchUrgeIssue(c *gin.Context) {
	cacheUserInfo, errSys := GetCacheUserInfo(c)
	if errSys != nil {
		Fail(c, errSys)
		return
	}
	req := vo.LessBatchUrgeIssueReq{}
	err := c.BindJSON(&req)
	if err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	resp := projectfacade.BatchUrgeIssue(&projectvo.BatchUrgeIssueReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input:  &req,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		Fail(c, resp.Error())
		return
	}
	Success(c, resp.Data)
}

func (h issueListHandler) BatchUpdateIssue(c *gin.Context) {
	cacheUserInfo, errSys := GetCacheUserInfo(c)
	if errSys != nil {
		Fail(c, errSys)
		return
	}
	req := vo.LessBatchUpdateIssueReq{}
	err := c.BindJSON(&req)
	if err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	appId, err1 := cast.ToInt64E(req.AppId)
	tableId, err2 := cast.ToInt64E(req.TableId)
	projectId, err3 := cast.ToInt64E(c.Param("projectId"))
	if err1 != nil || err2 != nil || err3 != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	resp := projectfacade.BatchUpdateIssue(&projectvo.BatchUpdateIssueReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input: &projectvo.BatchUpdateIssueInput{
			AppId:     appId,
			ProjectId: projectId,
			TableId:   tableId,
			Data:      req.Data,
		},
	})
	if resp.Failure() {
		log.Error(resp.Error())
		Fail(c, resp.Error())
		return
	}
	Success(c, true)
}

func (issueListHandler) UpdateIssueNew(c *gin.Context) {
	cacheUserInfo, errSys := GetCacheUserInfo(c)
	if errSys != nil {
		Fail(c, errSys)
		return
	}
	req := vo.LessUpdateIssueReq{}
	err := c.BindJSON(&req)
	if err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	resp := projectfacade.BatchUpdateIssue(&projectvo.BatchUpdateIssueReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input: &projectvo.BatchUpdateIssueInput{
			AppId:        req.MenuAppId,
			ProjectId:    -1,
			TableId:      -1,
			Data:         req.Form,
			BeforeDataId: req.BeforeDataId,
			AfterDataId:  req.AfterDataId,
			TodoId:       req.TodoId,
			TodoOp:       req.TodoOp,
			TodoMsg:      req.TodoMsg,
		},
	})
	if resp.Failure() {
		log.Error(resp.Error())
		Fail(c, resp.Error())
		return
	}
	Success(c, true)
}

func (h issueListHandler) CreateIssueNew(c *gin.Context) {
	cacheUserInfo, errSys := GetCacheUserInfo(c)
	if errSys != nil {
		Fail(c, errSys)
		return
	}
	req := vo.LessCreateIssueReq{}
	err := h.unmarshal(c, &req)
	if err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	appId := cast.ToInt64(req.MenuAppId)

	tableId, err := cast.ToInt64E(req.TableId)
	if err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	var beforeDataId, afterDataId int64
	if req.BeforeDataId != nil {
		beforeDataId = cast.ToInt64(*req.BeforeDataId)
	}
	if req.AfterDataId != nil {
		afterDataId = cast.ToInt64(*req.AfterDataId)
	}

	res := make(map[string]interface{}, 0)
	resp := projectfacade.BatchCreateIssue(&projectvo.BatchCreateIssueReqVo{
		OrgId:        cacheUserInfo.OrgId,
		UserId:       cacheUserInfo.UserId,
		AppId:        appId,
		ProjectId:    req.ProjectId,
		TableId:      tableId,
		BeforeDataId: beforeDataId,
		AfterDataId:  afterDataId,
		Data:         req.Form,
	})
	if resp.Failure() {
		Fail(c, resp.Error())
		return
	}

	res["total"] = len(resp.Data)
	res["actualTotal"] = len(resp.Data)
	res["list"] = resp.Data
	res["userDepts"] = resp.UserDepts
	res["relateData"] = resp.RelateData
	Success(c, res)
}

// @Security ApiKeyAuth
// @Summary 任务删除
// @Description 任务删除
// @Tags 任务
// @accept application/json
// @Produce application/json
// @Param projectId path int64 false "项目id"
// @Param req body vo.LessDeleteIssueBatchReq true "入参"
// @Success 200 {object} vo.DeleteIssueBatchResp
// @Failure 400
// @Router /api/rest/project/{projectId}/values/delete [post]
func (issueListHandler) DeleteIssue(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	req := vo.LessDeleteIssueBatchReq{}
	err1 := c.BindJSON(&req)
	if err1 != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}
	inputAppIdStr := req.MenuAppId
	if inputAppIdStr == "" {
		inputAppIdStr = "0"
	}
	inputAppId, inputAppIdErr := strconv.ParseInt(inputAppIdStr, 10, 64)
	if inputAppIdErr != nil {
		Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, inputAppIdErr.Error()))
		return
	}

	input := vo.DeleteIssueBatchReq{
		ProjectID: 0,
		Ids:       req.AppValueIds,
	}

	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}
	input.ProjectID = projectId
	input.TableID = req.TableId

	respVo := projectfacade.DeleteIssueBatch(projectvo.DeleteIssueBatchReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
		InputAppId:    inputAppId,
	})

	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Success(c, respVo.Data)
	}
}

func (issueListHandler) CopyIssue(c *gin.Context) {
	cacheUserInfo, sysErr := GetCacheUserInfo(c)
	if sysErr != nil {
		Fail(c, sysErr)
		return
	}
	req := vo.LessCopyIssueReq{}
	err := c.BindJSON(&req)
	if err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	respVo := projectfacade.CopyIssue(projectvo.CopyIssueReqVo{
		Input:         &req,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})

	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Success(c, map[string]interface{}{
			"list":       respVo.Data,
			"relateData": respVo.RelateData,
			"userDept":   respVo.UserDept,
			"total":      len(respVo.Data),
		})
	}
}

// @Security ApiKeyAuth
// @Summary 任务统计
// @Description 任务统计
// @Tags 任务
// @accept application/json
// @Produce application/json
// @Param projectId query int64 false "项目id"
// @Param iterationId query int64 false "迭代id"
// @Param relationType query int false "关联类型：1我负责的2我参与的3我关注的4我发起的5我确认的"
// @Success 200 {object} vo.IssueStatusTypeStatResp
// @Failure 400
// @Router /api/rest/issue/stat [get]
func (issueListHandler) IssueStat(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	req := &vo.IssueStatusTypeStatReq{}
	if c.Query("projectId") != "" {
		projectId, parseErr := strconv.ParseInt(c.Query("projectId"), 10, 64)
		if parseErr != nil {
			Fail(c, errs.ParamError)
			return
		}
		req.ProjectID = &projectId
	}
	if c.Query("iterationId") != "" {
		iterationId, parseErr := strconv.ParseInt(c.Query("iterationId"), 10, 64)
		if parseErr != nil {
			Fail(c, errs.ParamError)
			return
		}
		req.IterationID = &iterationId
	}
	if c.Query("relationType") != "" {
		relationType, parseErr := strconv.Atoi(c.Query("relationType"))
		if parseErr != nil {
			Fail(c, errs.ParamError)
			return
		}
		req.RelationType = &relationType
	}

	respVo := projectfacade.IssueStatusTypeStat(projectvo.IssueStatusTypeStatReqVo{
		Input:  req,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})

	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Success(c, respVo.IssueStatusTypeStat)
	}
}

func (issueListHandler) ViewStat(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}

	respVo := projectfacade.LcViewStatForAll(&projectvo.LcViewStatReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Success(c, respVo.Data)
	}
}

// // @Security ApiKeyAuth
// // @Summary 视图镜像统计
// // @Description 视图镜像统计
// // @Tags 镜像
// // @accept application/json
// // @Produce application/json
// // @Param req body vo.MirrorCountReq true "入参"
// // @Success 200 {object} vo.MirrorsStatResp
// // @Failure 400
// // @Router /api/rest/mirrors/stat [post]
func (issueListHandler) MirrorsStat(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}

	req := vo.MirrorCountReq{}
	err1 := c.BindJSON(&req)
	if err1 != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	respVo := projectfacade.MirrorsStat(&projectvo.MirrorsStatReq{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input:  req,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Success(c, respVo.Data)
	}
}

// 根据TableIds查询任务列表基本信息
func (i issueListHandler) GetIssueListSimpleByTableIds(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	inputReq := projectvo.IssueListWithConditionsReq{}
	err2 := i.unmarshal(c, &inputReq)
	if err2 != nil {
		Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err2.Error()))
		return
	}

	req := projectvo.GetIssueListWithConditionsReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  inputReq,
	}

	resp := projectfacade.IssueListSimpleByTableIds(req)
	if resp.Failure() {
		Fail(c, resp.Error())
	} else {
		Success(c, resp.Data)
	}
}

// 根据DataIds查询任务列表基本信息
func (i issueListHandler) GetIssueListSimpleByDataIds(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	inputReq := projectvo.GetIssueListSimpleByDataIdsReq{}
	err2 := i.unmarshal(c, &inputReq)
	if err2 != nil {
		Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err2.Error()))
		return
	}

	req := projectvo.GetIssueListSimpleByDataIdsReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  inputReq,
	}
	resp := projectfacade.IssueListSimpleByDataIds(req)
	if resp.Failure() {
		Fail(c, resp.Error())
	} else {
		Success(c, resp.Data)
	}
}

func (i issueListHandler) unmarshal(c *gin.Context, v interface{}) error {
	bts, err := c.GetRawData()
	if err != nil {
		return err
	}

	return encoding.GetJsonCodec().Unmarshal(bts, v)
}

// StartChat 任务-发起群聊
func (i issueListHandler) StartChat(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	inputReq := projectvo.IssueStartChatReq{}
	err2 := i.unmarshal(c, &inputReq)
	if err2 != nil {
		Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err2.Error()))
		return
	}
	req := projectvo.IssueStartChatReqVo{
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Input:         inputReq,
	}
	resp := projectfacade.IssueStartChat(req)
	if resp.Failure() {
		Fail(c, resp.Error())
	} else {
		Success(c, resp.Data)
	}
}

func (i issueListHandler) IssueCardShare(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}

	inputReq := projectvo.IssueCardShareReq{}
	err2 := i.unmarshal(c, &inputReq)
	if err2 != nil {
		Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err2.Error()))
		return
	}

	req := projectvo.IssueCardShareReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input:  inputReq,
	}

	resp := projectfacade.IssueShareCard(req)
	if resp.Failure() {
		Fail(c, resp.Error())
	} else {
		Success(c, resp.Data)
	}
}
