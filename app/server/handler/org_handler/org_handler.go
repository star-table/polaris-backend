package org_handler

import (
	"github.com/star-table/polaris-backend/app/server/handler"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/gin-gonic/gin"
)

type orgHandlers struct{}

var OrgHandler orgHandlers

func (co orgHandlers) TransferOrg(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	inputReqVo := orgvo.PrepareTransferOrgInput{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	respVo := orgfacade.TransferOrg(orgvo.PrepareTransferOrgReq{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input:  &inputReqVo,
	})
	if respVo.Failure() {
		handler.Fail(c, respVo.Error())
	} else {
		handler.Success(c, respVo.Void)
	}
}

func (co orgHandlers) PrepareTransferOrg(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	inputReqVo := orgvo.PrepareTransferOrgInput{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	respVo := orgfacade.PrepareTransferOrg(orgvo.PrepareTransferOrgReq{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input:  &inputReqVo,
	})
	if respVo.Failure() {
		handler.Fail(c, respVo.Error())
	} else {
		handler.Success(c, respVo.Data)
	}
}

// GetOrgConfig 获取组织配置信息
func (orgHandlers) GetOrgConfig(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}

	req := orgvo.GetOrgConfigReq{
		OrgId: cacheUserInfo.OrgId,
	}
	resp := orgfacade.GetOrgConfigRest(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

// SendCardToAdminForUpgrade 给组织管理员发送卡片，有成员申请升级极星
func (orgHandlers) SendCardToAdminForUpgrade(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}

	req := orgvo.SendCardToAdminForUpgradeReqVo{
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		SourceChannel: cacheUserInfo.SourceChannel,
	}
	resp := orgfacade.SendCardToAdminForUpgrade(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.IsTrue)
	}
}

func (orgHandlers) ExchangeShareToken(c *gin.Context) {
	inputReqVo := orgvo.ExchangeShareTokenData{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	resp := orgfacade.ExchangeShareToken(orgvo.ExchangeShareTokenReq{Input: inputReqVo})
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}
