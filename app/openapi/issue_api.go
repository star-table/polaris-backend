package openapi

import (
	"strconv"

	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/trendsvo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/facade/trendsfacade"
	"github.com/gin-gonic/gin"
)

func CreateIssue(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	openCreateIssueReq := projectvo.OpenCreateIssueReq{}
	err1 := c.BindJSON(&openCreateIssueReq)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	if openCreateIssueReq.OperatorId == int64(0) {
		Fail(c, errs.OperatorInvalid)
		return
	}
	respVo := projectfacade.OpenCreateIssue(projectvo.CreateIssueReqVo{
		// CreateIssue: openCreateIssueReq.CreateIssueReq,
		OrgId:  authData.OrgID,
		UserId: openCreateIssueReq.OperatorId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.Issue)
	}
}

func UpdateIssue(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	openCreateIssueReq := &projectvo.OpenUpdateIssueReq{}
	err1 := c.BindJSON(openCreateIssueReq)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}

	if openCreateIssueReq.OperatorId == int64(0) {
		Fail(c, errs.OperatorInvalid)
		return
	}

	// issueId, err1 := strconv.ParseInt(c.Param("issueId"), 10, 64)
	// if err1 != nil {
	// 	Fail(c, errs.ReqParamsValidateError)
	// 	return
	// }
	//openCreateIssueReq.ID = issueId
	respVo := projectfacade.OpenUpdateIssue(projectvo.OpenUpdateIssueReqVo{
		Data:   openCreateIssueReq,
		OrgId:  authData.OrgID,
		UserId: openCreateIssueReq.OperatorId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.UpdateIssue)
	}
}

func IssueStatusTypeStat(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	req := &vo.IssueStatusTypeStatReq{}
	err1 := c.BindQuery(req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}

	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	if projectId != 0 {
		req.ProjectID = &projectId
	}
	operatorId, _ := strconv.ParseInt(c.Query("operatorId"), 10, 64)

	relateType, _ := strconv.Atoi(c.Query("relateType"))
	//这里做个映射关系，作为openApi尽量让参数统一一些：统计接口里面relateType 和 homeIssues里面不一致。这里外层传入以homeIssues为准，里面做一层转换
	realType := 0
	switch relateType {
	case 1:
		//我发起的
		realType = 4
	case 2:
		//我负责的
		realType = 1
	case 3:
		//我参与的
		realType = 2
	case 4:
		//我关注的
		realType = 3
	case 5:
		//我确认的
		realType = 5
	case 6:
		//我负责的+我关注的
		realType = 6
	}
	req.RelationType = &realType

	respVo := projectfacade.OpenIssueStatusTypeStat(projectvo.IssueStatusTypeStatReqVo{
		Input:  req,
		OrgId:  authData.OrgID,
		UserId: operatorId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.IssueStatusTypeStat)
	}
}

func IssueList(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	req := &projectvo.OpenIssueListReq{}
	err1 := c.BindJSON(req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	page := 1
	size := 10
	if c.Query("page") != "" {
		page = ParseInt(c.Query("page"))
	}
	if c.Query("size") != "" {
		size = ParseInt(c.Query("size"))
	}
	respVo := projectfacade.OpenIssueList(projectvo.OpenIssueListReqVo{
		//NewData:   &req.HomeIssueInfoReq,
		Page:   page,
		Size:   size,
		OrgId:  authData.OrgID,
		UserId: req.OperatorId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.Data)
	}
}

func DeleteIssue(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	req := projectvo.OpenOperatorReq{}
	err1 := c.BindJSON(&req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	if req.OperatorId == int64(0) {
		Fail(c, errs.OperatorInvalid)
		return
	}
	issueId := ParseInt64(c.Param("issueId"))
	respVo := projectfacade.OpenDeleteIssue(projectvo.OpenDeleteIssueReqVo{
		Data: &projectvo.OpenDeleteIssueReq{
			IssueId: issueId,
		},
		UserId: req.OperatorId,
		OrgId:  authData.OrgID,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, vo.Void{ID: issueId})
	}
}

func GetPriorityList(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	respVo := projectfacade.OpenGetPriorityList(projectvo.OpenPriorityListReqVo{
		OrgId:     authData.OrgID,
		ProjectId: projectId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.Data)
	}
}

func GetProjectObjectTypeList(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	respVo := projectfacade.OpenGetProjectObjectTypeList(projectvo.OpenPriorityListReqVo{
		OrgId:     authData.OrgID,
		ProjectId: projectId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.Data)
	}
}

func GetIterationList(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	respVo := projectfacade.OpenGetIterationList(projectvo.OpenGetIterationListReqVo{
		OrgId:     authData.OrgID,
		ProjectId: projectId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.Data)
	}
}

func GetIssueSourceList(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	respVo := projectfacade.OpenGetIssueSourceList(projectvo.OpenGetDemandSourceListReqVo{
		OrgId:     authData.OrgID,
		ProjectId: projectId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.Data)
	}
}

func GetPropertyList(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	respVo := projectfacade.OpenGetPropertyList(projectvo.OpenGetPropertyListReqVo{
		OrgId:     authData.OrgID,
		ProjectId: projectId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.Data)
	}
}

func CreateIssueComment(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	req := &projectvo.OpenCreateIssueCommentReq{}
	err1 := c.BindJSON(req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	if req.OperatorId == int64(0) {
		Fail(c, errs.OperatorInvalid)
		return
	}
	issueId := ParseInt64(c.Param("issueId"))
	respVo := projectfacade.OpenCreateIssueComment(projectvo.CreateIssueCommentReqVo{
		Input: vo.CreateIssueCommentReq{
			Comment:          req.Comment,
			IssueID:          issueId,
			MentionedUserIds: req.MentionedUserIds,
		},
		UserId: req.OperatorId,
		OrgId:  authData.OrgID,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, vo.Void{ID: issueId})
	}
}

//func MirrorsStat(c *gin.Context) {
//	authData, err := ParseOpenAuthInfo(c)
//	if err != nil {
//		Fail(c, err)
//		return
//	}
//
//	req := projectvo.OpenMirrorsStatReq{}
//	err1 := c.BindJSON(&req)
//	if err1 != nil {
//		Fail(c, errs.ReqParamsValidateError)
//		return
//	}
//
//	if req.UserId == int64(0) {
//		Fail(c, errs.OperatorInvalid)
//		return
//	}
//
//	respVo := projectfacade.MirrorsStat(&projectvo.MirrorsStatReq{
//		UserId: req.UserId,
//		OrgId:  authData.OrgID,
//		Input:  vo.MirrorCountReq{AppIds: req.AppIds},
//	})
//	if respVo.Failure() {
//		Fail(c, respVo.Error())
//	} else {
//		Suc(c, respVo.NewData)
//	}
//}

func IssueCommentsList(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	req := &vo.TrendReq{}
	err1 := c.BindJSON(req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}

	issueId := ParseInt64(c.Param("issueId"))
	trendType := 2 //目前固定只能查询评论
	objType := consts.TrendsOperObjectTypeIssue
	respVo := trendsfacade.TrendList(trendsvo.TrendListReqVo{
		Input: &vo.TrendReq{
			LastTrendID: req.LastTrendID,
			ObjType:     &objType,
			ObjID:       &issueId,
			OperID:      nil,
			StartTime:   req.StartTime,
			EndTime:     req.EndTime,
			Type:        &trendType,
			Page:        req.Page,
			Size:        req.Size,
			OrderType:   req.OrderType,
		},
		OrgId:  authData.OrgID,
		UserId: 0,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		info := respVo.TrendList
		res := projectvo.OpenIssueCommentsListResp{
			Page: info.Page,
			Size: info.Size,
			//LastId: info.LastTrendID,
			Total: info.Total,
			List:  nil,
		}
		for _, trend := range info.List {
			res.List = append(res.List, projectvo.OpenIssueCommentInfo{
				Id:          trend.ID,
				CreatorInfo: trend.CreatorInfo,
				Comment:     trend.Comment,
				CreateTime:  trend.CreateTime,
			})
		}
		Suc(c, res)
	}
}

/**
 * 任务状态列表
 *
 **/
func GetIssueStatusList(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}

	req := projectvo.OpenGetIssueStatusListReq{}
	err1 := c.BindJSON(&req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}

	columnIds := make([]string, 0)
	// 添加查询字段
	columnIds = append(columnIds, "issueStatus")

	columnsResp := projectfacade.GetTablesColumns(projectvo.GetTablesColumnsReq{
		OrgId:  authData.OrgID,
		UserId: req.OperatorId,
		Input: &projectvo.TablesColumnsInput{
			TableIds:  req.TableIds,
			ColumnIds: columnIds,
		},
	})

	baseColumnsResp := projectvo.OpenGetIssueStatusListResp{
		Err:  columnsResp.Err,
		Data: &projectvo.IssueStatusListRespData{},
	}

	baseColumnsResp.Code = columnsResp.Code

	if columnsResp.Data != nil && len(columnsResp.Data.Tables) > 0 {
		baseColumnsResp.Data.Tables = make([]*projectvo.IssueStatusListColumnsTable, 0)

		for _, table := range columnsResp.Data.Tables {
			newTable := &projectvo.IssueStatusListColumnsTable{}
			for _, column := range table.Columns {
				props := column.Field.Props
				newTable.Props = props["groupSelect"]
				newTable.Name = column.Name
				newTable.TableId = table.TableId
			}
			baseColumnsResp.Data.Tables = append(baseColumnsResp.Data.Tables, newTable)
		}
	}

	if baseColumnsResp.Failure() {
		Fail(c, baseColumnsResp.Error())
	} else {
		Suc(c, baseColumnsResp.Data)
	}
}
