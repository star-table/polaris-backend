package orderfacade

import (
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/ordervo"
	"github.com/star-table/polaris-backend/facade"
)

func AddDingOrder(req ordervo.AddDingOrderReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/ordersvc/addDingOrder", config.GetPreUrl("ordersvc"))
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AddWeiXinOrder(req ordervo.AddWeiXinOrderReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/ordersvc/addWeiXinOrder", config.GetPreUrl("ordersvc"))
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AddFsOrder(req ordervo.AddFsOrderReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/ordersvc/addFsOrder", config.GetPreUrl("ordersvc"))
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CheckFsOrgIsInTrial(req ordervo.CheckFsOrgIsInTrialReq) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/ordersvc/checkFsOrgIsInTrial", config.GetPreUrl("ordersvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateWeiXinLicenceOrder(req ordervo.CreateWeiXinLicenceOrderReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/ordersvc/createWeiXinLicenceOrder", config.GetPreUrl("ordersvc"))
	queryParams := map[string]interface{}{}
	queryParams["suitId"] = req.SuitId
	queryParams["corpId"] = req.CorpId
	queryParams["orderId"] = req.OrderId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteDingOrder(req ordervo.DeleteDingOrderReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/ordersvc/deleteDingOrder", config.GetPreUrl("ordersvc"))
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetFunctionByLevel(req ordervo.FunctionReq) ordervo.FunctionResp {
	respVo := &ordervo.FunctionResp{}
	reqUrl := fmt.Sprintf("%s/api/ordersvc/getFunctionByLevel", config.GetPreUrl("ordersvc"))
	queryParams := map[string]interface{}{}
	queryParams["level"] = req.Level
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateDingOrder(req ordervo.UpdateDingOrderReq) ordervo.UpdateDingOrderDataResp {
	respVo := &ordervo.UpdateDingOrderDataResp{}
	reqUrl := fmt.Sprintf("%s/api/ordersvc/updateDingOrder", config.GetPreUrl("ordersvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrderPayInfo(req ordervo.GetOrderPayInfoReq) ordervo.GetOrderPayInfoResp {
	respVo := &ordervo.GetOrderPayInfoResp{}
	reqUrl := fmt.Sprintf("%s/api/ordersvc/getOrderPayInfo", config.GetPreUrl("ordersvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}
