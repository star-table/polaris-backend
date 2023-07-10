package org_handler

import (
	"github.com/star-table/polaris-backend/app/server/handler"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/gin-gonic/gin"
)

type dingTalkHandlers struct{}

var DingTalkHandlers dingTalkHandlers

func (d dingTalkHandlers) GetJsApiSign(c *gin.Context) {
	inputReqVo := orgvo.GetDingApiSignReq{Input: &orgvo.JsAPISignReq{}}
	err1 := c.ShouldBindJSON(&inputReqVo.Input)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	resp := orgfacade.GetDingJsAPISign(inputReqVo)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.GetJsAPISign)
	}
}

func (d dingTalkHandlers) AuthCode(c *gin.Context) {
	inputReqVo := orgvo.AuthDingCodeReqVo{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	resp := orgfacade.AuthDingCode(inputReqVo)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

//func (d dingTalkHandlers) GetCoolAppInfo(c *gin.Context) {
//	inputReqVo := orgvo.DeleteCoolAppReq{}
//	err1 := c.ShouldBindJSON(&inputReqVo)
//	if err1 != nil {
//		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
//		return
//	}
//
//	resp := orgfacade.GetCoolAppInfo(inputReqVo)
//	if resp.Failure() {
//		handler.Fail(c, resp.Error())
//	} else {
//		handler.Success(c, resp.Data)
//	}
//}
