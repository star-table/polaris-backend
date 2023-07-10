package v2

import (
	"fmt"
	"strconv"

	"github.com/star-table/polaris-backend/app/openapi"

	"github.com/star-table/common/core/logger"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/businees"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/trendsvo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/facade/trendsfacade"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

var log = logger.GetDefaultLogger()

func CreateIssue(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}

	openCreateIssueReq := projectvo.OpenCreateIssueReq{}
	err1 := c.BindJSON(&openCreateIssueReq)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}

	if openCreateIssueReq.OperatorId == int64(0) {
		openapi.Fail(c, errs.OperatorInvalid)
		return
	}

	var beforeDataId, afterDataId int64
	if openCreateIssueReq.BeforeDataId != nil {
		beforeDataId = cast.ToInt64(*openCreateIssueReq.BeforeDataId)
	}
	if openCreateIssueReq.AfterDataId != nil {
		afterDataId = cast.ToInt64(*openCreateIssueReq.AfterDataId)
	}

	res := make(map[string]interface{}, 0)
	resp := projectfacade.BatchCreateIssue(&projectvo.BatchCreateIssueReqVo{
		OrgId:        authData.OrgID,
		UserId:       openCreateIssueReq.OperatorId,
		AppId:        cast.ToInt64(openCreateIssueReq.MenuAppId),
		ProjectId:    openCreateIssueReq.ProjectId,
		TableId:      cast.ToInt64(openCreateIssueReq.TableId),
		BeforeDataId: beforeDataId,
		AfterDataId:  afterDataId,
		Data:         openCreateIssueReq.Form,
	})
	if resp.Failure() {
		openapi.Fail(c, resp.Error())
		return
	}

	if len(resp.Data) > 0 {
		res = resp.Data[0]
	}

	openapi.Suc(c, res)
}

func UpdateIssue(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}

	openUpdateIssueReq := &projectvo.OpenUpdateIssueReq{}
	err1 := c.BindJSON(openUpdateIssueReq)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}

	if openUpdateIssueReq.OperatorId == int64(0) {
		openapi.Fail(c, errs.OperatorInvalid)
		return
	}

	if len(openUpdateIssueReq.Form) == 0 {
		openapi.Suc(c, true)
		return
	}

	//issueId, err1 := cast.ToInt64E(openUpdateIssueReq.Form[0][consts.BasicFieldId])
	//if err1 != nil {
	//	openapi.Fail(c, errs.ReqParamsValidateError)
	//	return
	//}

	// app 应当尽量不处理业务逻辑，尽快把请求传递给后端服务
	resp := projectfacade.BatchUpdateIssue(&projectvo.BatchUpdateIssueReqVo{
		UserId: openUpdateIssueReq.OperatorId,
		OrgId:  authData.OrgID,
		Input: &projectvo.BatchUpdateIssueInput{
			AppId:        openUpdateIssueReq.MenuAppId,
			ProjectId:    -1,
			TableId:      -1,
			Data:         openUpdateIssueReq.Form,
			BeforeDataId: openUpdateIssueReq.BeforeDataId,
			AfterDataId:  openUpdateIssueReq.AfterDataId,
		},
	})
	if resp.Failure() {
		log.Error(resp.Error())
		openapi.Fail(c, resp.Error())
		return
	}

	//// 处理排序
	//if openUpdateIssueReq.BeforeDataId != nil || openUpdateIssueReq.AfterDataId != nil {
	//	moveResp := projectfacade.UpdateIssueSort(projectvo.UpdateIssueSortReqVo{
	//		Input: vo.UpdateIssueSortReq{
	//			ID:           issueId,
	//			BeforeDataID: openUpdateIssueReq.BeforeDataId,
	//			AfterDataID:  openUpdateIssueReq.AfterDataId,
	//			Asc:          openUpdateIssueReq.Asc,
	//		},
	//		UserId:        openUpdateIssueReq.OperatorId,
	//		OrgId:         authData.OrgID,
	//		SourceChannel: "openapi",
	//	})
	//
	//	if moveResp.Failure() {
	//		//这里报错只是排序有问题
	//		log.Error(resp.Error())
	//		openapi.Fail(c, moveResp.Error())
	//		return
	//	}
	//}

	openapi.Suc(c, true)
}

func IssueInfo(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}

	issueId, err1 := strconv.ParseInt(c.Param("issueId"), 10, 64)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}

	respVo := projectfacade.OpenIssueInfo(projectvo.IssueInfoReqVo{
		IssueID: issueId,
		OrgId:   authData.OrgID,
	})
	if respVo.Failure() {
		openapi.Fail(c, respVo.Error())
	} else {
		openapi.Suc(c, respVo.Data)
	}
}

func IssueStatusTypeStat(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}
	req := &vo.IssueStatusTypeStatReq{}
	err1 := c.BindQuery(req)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}

	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
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
		openapi.Fail(c, respVo.Error())
	} else {
		openapi.Suc(c, respVo.IssueStatusTypeStat)
	}
}

func IssueList(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}

	openIssueListReq := &projectvo.OpenIssueListReq{}
	err1 := c.BindJSON(openIssueListReq)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}

	zero := int(0)
	if openIssueListReq.Page == nil {
		openIssueListReq.Page = &zero
	}
	if openIssueListReq.Size == nil {
		openIssueListReq.Size = &zero
	}

	projectId := openIssueListReq.ProjectId
	if projectId < 0 {
		projectId = 0
	}

	if openIssueListReq.TableId == nil {
		temp := "0"
		openIssueListReq.TableId = &temp
	}

	params := &projectvo.HomeIssuesReqVo{
		Page: *openIssueListReq.Page,
		Size: *openIssueListReq.Size,
		Input: &projectvo.HomeIssueInfoReq{
			MenuAppID:     openIssueListReq.MenuAppId,
			ProjectID:     &projectId,
			LessOrder:     openIssueListReq.Orders,
			TableID:       openIssueListReq.TableId,
			FilterColumns: openIssueListReq.FilterColumns,
		},
		UserId: openIssueListReq.OperatorId,
		OrgId:  authData.OrgID,
	}

	//类型(between,equal,gt,gte,in,like,lt,lte,not_in,not_like,not_null,is_null,all_in,values_in)
	if openIssueListReq.Condition != nil && len(openIssueListReq.Condition.Conds) > 0 {
		var lessConds []*vo.LessCondsData
		for _, data := range openIssueListReq.Condition.Conds {
			if data.Column == "fuzzyCode" {
			} else {
				lessConds = append(lessConds, data)
			}
		}
		if lessConds != nil {
			params.Input.LessConds = &vo.LessCondsData{
				Type:  "and",
				Conds: lessConds,
			}
		}
	}

	withoutInfo := c.Query("withoutInfo")
	if withoutInfo == "4" {
		// 主页菜单任务列表
		respVo := projectfacade.OpenLcHomeIssuesForAll(params)
		openapi.SuccessJson(c, respVo)
	} else {
		// 项目任务列表
		respVo := projectfacade.OpenLcHomeIssues(params)
		openapi.SuccessJson(c, respVo)
	}
}

func DeleteIssue(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}
	req := projectvo.OpenOperatorReq{}
	err1 := c.BindJSON(&req)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}
	if req.OperatorId == int64(0) {
		openapi.Fail(c, errs.OperatorInvalid)
		return
	}
	issueId := openapi.ParseInt64(c.Param("issueId"))
	respVo := projectfacade.OpenDeleteIssue(projectvo.OpenDeleteIssueReqVo{
		Data: &projectvo.OpenDeleteIssueReq{
			IssueId: issueId,
		},
		UserId: req.OperatorId,
		OrgId:  authData.OrgID,
	})
	if respVo.Failure() {
		openapi.Fail(c, respVo.Error())
	} else {
		openapi.Suc(c, vo.Void{ID: issueId})
	}
}

func GetPriorityList(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}
	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}
	respVo := projectfacade.OpenGetPriorityList(projectvo.OpenPriorityListReqVo{
		OrgId:     authData.OrgID,
		ProjectId: projectId,
	})
	if respVo.Failure() {
		openapi.Fail(c, respVo.Error())
	} else {
		openapi.Suc(c, respVo.Data)
	}
}

func GetProjectObjectTypeList(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}
	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}
	respVo := projectfacade.OpenGetProjectObjectTypeList(projectvo.OpenPriorityListReqVo{
		OrgId:     authData.OrgID,
		ProjectId: projectId,
	})
	if respVo.Failure() {
		openapi.Fail(c, respVo.Error())
	} else {
		openapi.Suc(c, respVo.Data)
	}
}

func GetIterationList(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}
	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}
	respVo := projectfacade.OpenGetIterationList(projectvo.OpenGetIterationListReqVo{
		OrgId:     authData.OrgID,
		ProjectId: projectId,
	})
	if respVo.Failure() {
		openapi.Fail(c, respVo.Error())
	} else {
		openapi.Suc(c, respVo.Data)
	}
}

func GetIssueSourceList(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}
	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}
	respVo := projectfacade.OpenGetIssueSourceList(projectvo.OpenGetDemandSourceListReqVo{
		OrgId:     authData.OrgID,
		ProjectId: projectId,
	})
	if respVo.Failure() {
		openapi.Fail(c, respVo.Error())
	} else {
		openapi.Suc(c, respVo.Data)
	}
}

func GetPropertyList(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}
	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}
	respVo := projectfacade.OpenGetPropertyList(projectvo.OpenGetPropertyListReqVo{
		OrgId:     authData.OrgID,
		ProjectId: projectId,
	})
	if respVo.Failure() {
		openapi.Fail(c, respVo.Error())
	} else {
		openapi.Suc(c, respVo.Data)
	}
}

func CreateIssueComment(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}
	req := &projectvo.OpenCreateIssueCommentReq{}
	err1 := c.BindJSON(req)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}
	if req.OperatorId == int64(0) {
		openapi.Fail(c, errs.OperatorInvalid)
		return
	}
	issueId := openapi.ParseInt64(c.Param("issueId"))
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
		openapi.Fail(c, respVo.Error())
	} else {
		openapi.Suc(c, vo.Void{ID: issueId})
	}
}

//func MirrorsStat(c *gin.Context) {
//	authData, err := openapi.ParseOpenAuthInfo(c)
//	if err != nil {
//		openapi.Fail(c, err)
//		return
//	}
//
//	req := projectvo.OpenMirrorsStatReq{}
//	err1 := c.BindJSON(&req)
//	if err1 != nil {
//		openapi.Fail(c, errs.ReqParamsValidateError)
//		return
//	}
//
//	if req.UserId == int64(0) {
//		openapi.Fail(c, errs.OperatorInvalid)
//		return
//	}
//
//	respVo := projectfacade.MirrorsStat(&projectvo.MirrorsStatReq{
//		UserId: req.UserId,
//		OrgId:  authData.OrgID,
//		Input:  vo.MirrorCountReq{AppIds: req.AppIds},
//	})
//	if respVo.Failure() {
//		openapi.Fail(c, respVo.Error())
//	} else {
//		openapi.Suc(c, respVo.NewData)
//	}
//}

func IssueCommentsList(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}

	req := &projectvo.OpenIssueCommentsListReq{}
	err1 := c.BindJSON(req)
	if err1 != nil {
		openapi.Fail(c, errs.ReqParamsValidateError)
		return
	}

	if req.OperatorId == int64(0) {
		openapi.Fail(c, errs.OperatorInvalid)
		return
	}

	issueId := openapi.ParseInt64(c.Param("issueId"))
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
		UserId: req.OperatorId,
	})

	if respVo.Failure() {
		openapi.Fail(c, respVo.Error())
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
		openapi.Suc(c, res)
	}
}

func handleChildIssue(req interface{}) ([]*vo.IssueChildren, errs.SystemErrorInfo) {
	children := make([]map[string]interface{}, 0)
	copyErr := copyer.Copy(req, &children)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ObjectCopyError
	}
	res := make([]*vo.IssueChildren, 0)
	for _, m := range children {
		//人员单独处理（负责人，关注人，确认人）
		if idArr, ok := m[consts.BasicFieldOwnerId]; ok {
			ownerIds, errConvert := businees.LcMemberToUserIdsWithError(idArr)
			if errConvert != nil {
				log.Errorf("[handleChildIssue]error, ownerIds:%v, err:%v", idArr, errConvert)
				return nil, errConvert
			}
			m[consts.BasicFieldOwnerId] = ownerIds
		}
		if idArr, ok := m[consts.BasicFieldFollowerIds]; ok {
			followers, errConvert := businees.LcMemberToUserIdsWithError(idArr)
			if errConvert != nil {
				log.Errorf("[handleChildIssue]error, ownerIds:%v, err:%v", idArr, errConvert)
				return nil, errConvert
			}
			m[consts.BasicFieldFollowerIds] = followers
		}
		if idArr, ok := m[consts.BasicFieldAuditorIds]; ok {
			auditors, errConvert := businees.LcMemberToUserIdsWithError(idArr)
			if errConvert != nil {
				log.Errorf("[handleChildIssue]error, ownerIds:%v, err:%v", idArr, errConvert)
				return nil, errConvert
			}
			m[consts.BasicFieldAuditorIds] = auditors
		}
		temp := vo.IssueChildren{}
		err := json.FromJson(json.ToJsonIgnoreError(m), &temp)
		if err != nil {
			log.Error(err)
			return nil, errs.ReqParamsValidateError
		}
		if projectObjectTypeId, ok := m[consts.BasicFieldProjectObjectTypeId]; ok {
			typeId, err := strconv.ParseInt(fmt.Sprintf("%v", projectObjectTypeId), 10, 64)
			if err == nil {
				temp.TypeID = &typeId
			} else {
				log.Error(err)
			}
		}
		if status, ok := m[consts.BasicFieldIssueStatus]; ok {
			statusId, err := cast.ToInt64E(status)
			if err == nil {
				temp.StatusID = &statusId
			} else {
				log.Error(err)
			}
		}
		temp.LessCreateIssueReq = m
		res = append(res, &temp)
	}

	return res, nil
}
