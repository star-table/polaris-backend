package org_handler

import (
	"github.com/star-table/polaris-backend/app/server/handler"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/gin-gonic/gin"
)


func InviteWithPhones(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	inputReqVo := orgvo.InviteUserByPhonesReqVoData{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil{
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	respVo := orgfacade.InviteUserByPhones(orgvo.InviteUserByPhonesReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Input:  inputReqVo,
	})
	if respVo.Failure() {
		handler.Fail(c, respVo.Error())
	} else {
		handler.Success(c, respVo.Data)
	}
}

// 下载邀请新成员的模板
func ExportInviteTemplate(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	respVo := orgfacade.ExportInviteTemplate(orgvo.ExportInviteTemplateReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if respVo.Failure() {
		handler.Fail(c, respVo.Error())
	} else {
		handler.Success(c, respVo.Data)
	}
}

// 邀请-excel 导入组织成员
func ImportMembers(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	inputReqVo := orgvo.ImportMembersReqVoData{}
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil{
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	respVo := orgfacade.ImportMembers(orgvo.ImportMembersReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input: &inputReqVo,
	})
	if respVo.Failure() {
		handler.Fail(c, respVo.Error())
	} else {
		handler.Success(c, respVo.Data)
	}
}

