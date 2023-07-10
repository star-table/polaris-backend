package lab_handler

import (
	"github.com/star-table/polaris-backend/app/server/handler"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/gin-gonic/gin"
)

type labHandlers struct{}

var LabHandler labHandlers

func (labHandlers) SetLab(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo orgvo.SetLabReq
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := orgvo.SetLabReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  inputReqVo,
	}

	resp := orgfacade.SetLabConfig(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (labHandlers) GetLab(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	req := orgvo.GetLabReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	}

	resp := orgfacade.GetLabConfig(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}
