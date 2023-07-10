package handler

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/gin-gonic/gin"
)

type todoHandler struct{}

var TodoHandler todoHandler

//func (todoHandler) Filter(c *gin.Context) {
//	cacheUserInfo, sysErr := GetCacheUserInfo(c)
//	if sysErr != nil {
//		Fail(c, sysErr)
//		return
//	}
//
//	req := vo.TodoFilterReq{}
//	if err := c.BindJSON(&req); err != nil {
//		Fail(c, errs.ReqParamsValidateError)
//		return
//	}
//
//	respVo := projectfacade.TodoFilter(&projectvo.TodoFilterReq{
//		OrgId:      cacheUserInfo.OrgId,
//		UserId:     cacheUserInfo.UserId,
//		Page:       req.Page,
//		Size:       req.Size,
//		FilterType: req.FilterType,
//	})
//	if respVo.Failure() {
//		Fail(c, respVo.Error())
//	} else {
//		Success(c, respVo.Data)
//	}
//}

func (todoHandler) Urge(c *gin.Context) {
	cacheUserInfo, sysErr := GetCacheUserInfo(c)
	if sysErr != nil {
		Fail(c, sysErr)
		return
	}

	req := vo.TodoUrgeReq{}
	if err := c.BindJSON(&req); err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	respVo := projectfacade.TodoUrge(&projectvo.TodoUrgeReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input: &projectvo.TodoUrgeInput{
			TodoId: req.TodoId,
			Msg:    req.Msg,
		},
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Success(c, respVo.Void)
	}
}
