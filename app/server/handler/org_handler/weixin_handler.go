package org_handler

import (
	"net/http"

	"github.com/star-table/polaris-backend/app/server/handler"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/gin-gonic/gin"
)

type weiXinHandlers struct{}

var WeiXinHandlers weiXinHandlers

func (d weiXinHandlers) GetJsApiSign(c *gin.Context) {
	inputReqVo := orgvo.JsAPISignReq{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	resp := orgfacade.GetWeiXinJsAPISign(inputReqVo)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.GetJsAPISign)
	}
}

func (d weiXinHandlers) AuthCode(c *gin.Context) {
	inputReqVo := orgvo.AuthWeiXinCodeReqVo{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	resp := orgfacade.AuthWeiXinCode(inputReqVo)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (d weiXinHandlers) GetAgentId(c *gin.Context) {
	inputReqVo := orgvo.GetCorpAgentIdReq{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	resp := orgfacade.GetOrgOutInfoByOutOrgId(orgvo.GetOutOrgInfoByOutOrgIdReqVo{
		OutOrgId: inputReqVo.CorpId,
	})
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, &orgvo.GetCorpAgentIdResp{AgentId: resp.Data.TenantCode})
	}
}

func (d weiXinHandlers) GetQrCode(c *gin.Context) {
	inputReqVo := orgvo.PersonWeiXinQrCodeReqVo{}
	resp := orgfacade.PersonWeiXinQrCode(inputReqVo)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (d weiXinHandlers) CheckQrCodeScan(c *gin.Context) {
	inputReqVo := orgvo.CheckQrCodeScanReq{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	resp := orgfacade.CheckPersonWeiXinQrCode(inputReqVo)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (d weiXinHandlers) PersonWeiXinBind(c *gin.Context) {
	inputReqVo := orgvo.PersonWeiXinBindData{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	resp := orgfacade.PersonWeiXinBind(orgvo.PersonWeiXinBindReq{Data: &inputReqVo})
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (d weiXinHandlers) CheckMobileHasBind(c *gin.Context) {
	inputReqVo := orgvo.CheckMobileHasBindData{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	resp := orgfacade.CheckMobileHasBind(orgvo.CheckMobileHasBindReq{Input: &inputReqVo})
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (d weiXinHandlers) RedirectInstall(c *gin.Context) {
	resp := orgfacade.GetWeiXinInstallUrl(orgvo.GetRegisterUrlReq{})
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		c.Redirect(http.StatusTemporaryRedirect, resp.Data.Url)
	}
}

func (d weiXinHandlers) RedirectRegister(c *gin.Context) {
	resp := orgfacade.GetWeiXinRegisterUrl(orgvo.GetRegisterUrlReq{})
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		c.Redirect(http.StatusTemporaryRedirect, resp.Data.Url)
	}
}
