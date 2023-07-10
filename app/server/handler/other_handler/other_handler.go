package other_handler

import (
	"github.com/star-table/polaris-backend/app/server/handler"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/othervo"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/gin-gonic/gin"
)

type otherHandlers struct {}
var OtherHandler otherHandlers

// AddScriptTask 调用 schedulesvc 的一个接口，新增脚本任务并执行
func (otherHandlers) AddScriptTask(c *gin.Context) {
	inputReqVo := &othervo.AddScriptTaskReq{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil{
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	resp := userfacade.AddScriptTask(inputReqVo.OrgId, inputReqVo.UserId, uservo.AddScriptTaskReq{
		AppId:    inputReqVo.AppId,
		TaskId:   inputReqVo.TaskId,
		TaskName: inputReqVo.TaskName,
		CronSpec: inputReqVo.CronSpec,
		ParamMap: inputReqVo.ParamMap,
	})
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}
