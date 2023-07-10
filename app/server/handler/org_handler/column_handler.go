package org_handler

import (
	"github.com/star-table/go-common/pkg/encoding"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/gin-gonic/gin"

	tablev1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/polaris-backend/app/server/handler"
)

type columnHandlers struct{}

var ColumnHandler columnHandlers

func (co columnHandlers) GetOrgColumns(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	//var inputReqVo tablev1.ReadOrgColumnsRequest
	//err1 := co.unmarshal(c, &inputReqVo)
	//if err1 != nil {
	//	handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
	//	return
	//}
	req := orgvo.GetOrgColumnsReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  nil,
	}
	resp := orgfacade.GetOrgColumns(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (co columnHandlers) CreateOrgColumn(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo orgvo.CreateOrgColumnRequest
	err1 := co.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	req := orgvo.CreateOrgColumnReq{
		OrgId:  cacheUserInfo.UserId,
		UserId: cacheUserInfo.UserId,
		Input:  &inputReqVo,
	}
	resp := orgfacade.CreateOrgColumn(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, tablev1.CreateOrgColumnReply{})
	}
}

func (co columnHandlers) DeleteOrgColumn(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo orgvo.DeleteOrgColumnRequest
	err1 := co.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	req := orgvo.DeleteOrgColumnReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  &inputReqVo,
	}
	resp := orgfacade.DeleteOrgColumn(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, tablev1.DeleteOrgColumnReply{})
	}
}

func (co columnHandlers) unmarshal(c *gin.Context, v interface{}) error {
	bts, err := c.GetRawData()
	if err != nil {
		return err
	}

	return encoding.GetJsonCodec().Unmarshal(bts, v)
}
