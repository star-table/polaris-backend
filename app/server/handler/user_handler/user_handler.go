package user_handler

import (
	"github.com/star-table/go-common/pkg/encoding"
	"github.com/star-table/polaris-backend/app/server/handler"
	"github.com/star-table/polaris-backend/common/core/errs"
	_ "github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/gin-gonic/gin"
)

type userHandlers struct{}

var UserHandler userHandlers

func (u userHandlers) unmarshal(c *gin.Context, v interface{}) error {
	bts, err := c.GetRawData()
	if err != nil {
		return err
	}

	return encoding.GetJsonCodec().Unmarshal(bts, v)
}

// @Security PM-TOEKN
// @Summary 设定用户浏览过新用户指引
// @Description 设定用户浏览过新用户指引
// @Tags 用户
// @accept application/json
// @Produce application/json
// @Param flag query string true "状态类型、标识"
// @Success 200 {object} vo.BoolRespVo
// @Failure 400
// @Router /api/rest/user/setUserGuideStatus [post]
func (userHandlers) SetVisitUserGuideStatus(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	flag := c.Query("flag")
	resp := orgfacade.SetVisitUserGuideStatus(orgvo.SetVisitUserGuideStatusReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Flag:   flag,
	})
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.IsTrue)
	}
}

// 设置卡片弹窗显示flag
func (userHandlers) SetVersionVisible(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}

	var inputReqVo orgvo.SetVersionData
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	req := orgvo.SetVersionReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  inputReqVo,
	}
	resp := orgfacade.SetVersionVisible(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.IsTrue)
	}

}

func (userHandlers) SetUserActivity(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo orgvo.SetActivityData
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := orgvo.SetUserActivityReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  inputReqVo,
	}
	resp := orgfacade.SetUserActivity(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.IsTrue)
	}

}

func (userHandlers) GetVersion(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}

	req := orgvo.GetVersionReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	}

	resp := orgfacade.GetVersion(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (u userHandlers) SetUserViewLocation(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}

	var inputReqVo orgvo.SaveViewLocationReq
	err1 := u.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	req := orgvo.SaveViewLocationReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  &inputReqVo,
	}
	resp := orgfacade.SetUserViewLocation(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (userHandlers) GetUserViewLocation(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}

	req := orgvo.GetViewLocationReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	}
	resp := orgfacade.GetUserViewLastLocation(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (userHandlers) GetSameNameUserOrDept(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}

	var inputReqVo orgvo.GetUserOrDeptSameNameListReqVoData
	err1 := c.ShouldBindJSON(&inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	resp := projectfacade.GetUserOrDeptSameNameList(projectvo.GetUserOrDeptSameNameListReq{
		OrgId:    cacheUserInfo.OrgId,
		UserId:   cacheUserInfo.UserId,
		DataType: inputReqVo.DataType,
	})
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

// GetUserInfo 查询用户信息
func (userHandlers) GetUserInfo(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}

	req := orgvo.PersonalInfoReqVo{
		SourceChannel: cacheUserInfo.SourceChannel,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
	}
	resp := orgfacade.PersonalInfoRest(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}
