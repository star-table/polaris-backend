package appfacade

import (
	"fmt"

	"github.com/star-table/polaris-backend/facade"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/appvo"
)

// FormCreateIssue 新增任务数据
func FormCreateIssue(req *appvo.FormCreateIssueReq) *appvo.FormCreateOneResp {
	respVo := &appvo.FormCreateOneResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func FormCreatePriority(req *appvo.FormCreatePriorityReq) *appvo.FormCreateOneResp {
	respVo := &appvo.FormCreateOneResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}
