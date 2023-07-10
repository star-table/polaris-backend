package handler

import (
	"strconv"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/gin-gonic/gin"
)

type issueViewHandlers struct{}

var IssueViewHandles issueViewHandlers

// @Security PM-TOEKN
// @Summary 新增任务视图
// @Description 新增任务视图
// @Tags 任务视图
// @accept application/json
// @Produce application/json
// @param input body vo.CreateIssueViewReq true "入参"
// @Success 200 {object} int64
// @Failure 400
// @Router /api/rest/issueView/create [post]
func (issueViewHandlers) CreateTaskView(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	inputReqVo := vo.CreateIssueViewReq{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	resp := projectfacade.CreateTaskView(&projectvo.CreateTaskViewReqVo{
		Input:  inputReqVo,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if resp.Failure() {
		Fail(c, resp.Error())
	} else {
		Success(c, resp.Data)
	}
}

// 更新任务视图
func (issueViewHandlers) UpdateTaskView(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	inputReqVo := vo.UpdateIssueViewReq{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	resp := projectfacade.UpdateTaskView(&projectvo.UpdateTaskViewReqVo{
		Input:  inputReqVo,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if resp.Failure() {
		Fail(c, resp.Error())
	} else {
		respVo := &vo.BoolRespVoData{
			IsTrue: resp.Void.ID > 0,
		}
		Success(c, respVo)
	}
}

// 获取任务视图列表
func (issueViewHandlers) GetTaskViewList(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	inputReqVo := vo.GetIssueViewListReq{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	resp := projectfacade.GetTaskViewList(&projectvo.GetTaskViewListReqVo{
		Input:  inputReqVo,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if resp.Failure() {
		Fail(c, resp.Error())
	} else {
		Success(c, resp.Data)
	}
}

// 删除任务视图
func (issueViewHandlers) DeleteTaskView(c *gin.Context) {
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	viewIdStr := c.Param("viewId")
	viewId, _ := strconv.ParseInt(viewIdStr, 10, 64)
	resp := projectfacade.DeleteTaskView(&projectvo.DeleteTaskViewReqVo{
		Input: vo.DeleteIssueViewReq{
			ID: viewId,
		},
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if resp.Failure() {
		Fail(c, resp.Error())
	} else {
		respVo := &vo.BoolRespVoData{
			IsTrue: resp.Void.ID > 0,
		}
		Success(c, respVo)
	}
}

// 筛选时，获取可选值列表
func (issueViewHandlers) GetOptionList(c *gin.Context) {
	//cacheUserInfo, err := GetCacheUserInfo(c)
	//if err != nil {
	//	Fail(c, err)
	//	return
	//}
	//inputReq := projectvo.GetOptionListReqInput{}
	//err1 := c.ShouldBindJSON(&inputReq)
	//if err1 != nil {
	//	Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
	//	return
	//}
	//resp := projectfacade.GetOptionList(&projectvo.GetOptionListReqVo{
	//	Input:  inputReq,
	//	OrgId:  cacheUserInfo.OrgId,
	//	UserId: cacheUserInfo.UserId,
	//})
	//if resp.Failure() {
	//	Fail(c, resp.Error())
	//} else {
	//	Success(c, resp.Data)
	//}
}
