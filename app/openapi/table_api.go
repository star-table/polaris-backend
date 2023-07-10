package openapi

import (
	"strconv"

	tableV1 "github.com/star-table/interface/golang/table/v1"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/app/server/handler"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/facade/tablefacade"
	"github.com/gin-gonic/gin"
)

func TablesByAppId(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}

	appId, err1 := strconv.ParseInt(c.Param("appId"), 10, 64)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}

	req := projectvo.OpenOperatorReq{}
	err1 = c.BindJSON(&req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}

	appIds := make([]int64, 0)
	appIds = append(appIds, appId)

	respVo := tablefacade.ReadTablesByApps(projectvo.ReadTablesByAppsReqVo{
		OrgId:  authData.OrgID,
		UserId: 0,
		Input: &tableV1.ReadTablesByAppsRequest{
			AppIds: appIds,
		},
	})

	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.Data.AppsTables)
	}
}

func TablesBaseColumns(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}

	req := projectvo.OpenTablesColumnsReq{}
	err1 := c.BindJSON(&req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}

	columnsResp := projectfacade.GetTablesColumns(projectvo.GetTablesColumnsReq{
		OrgId:  authData.OrgID,
		UserId: 0,
		Input: &projectvo.TablesColumnsInput{
			TableIds:  req.TableIds,
			ColumnIds: req.ColumnIds,
		},
	})

	baseColumnsResp := projectvo.TablesBaseColumnsResp{
		Err:  columnsResp.Err,
		Data: &projectvo.TablesBaseColumnsRespData{},
	}

	baseColumnsResp.Code = columnsResp.Code

	if columnsResp.Data != nil && len(columnsResp.Data.Tables) > 0 {
		baseColumnsResp.Data.Tables = make([]*projectvo.TableBaseColumnsTable, len(columnsResp.Data.Tables))
		copyer.Copy(columnsResp.Data.Tables, &baseColumnsResp.Data.Tables)
	}

	if columnsResp.Failure() {
		handler.Fail(c, columnsResp.Error())
		return
	} else {
		handler.Success(c, baseColumnsResp.Data)
	}

}
