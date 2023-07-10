package idfacade

import (
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/idvo"
	"github.com/star-table/polaris-backend/facade"
)

func ApplyCode(req idvo.ApplyCodeReqVo) idvo.ApplyCodeRespVo {
	respVo := &idvo.ApplyCodeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/idsvc/applyCode", config.GetPreUrl("idsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["code"] = req.Code
	queryParams["preCode"] = req.PreCode
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ApplyMultipleId(req idvo.ApplyMultipleIdReqVo) idvo.ApplyMultipleIdRespVo {
	respVo := &idvo.ApplyMultipleIdRespVo{}
	reqUrl := fmt.Sprintf("%s/api/idsvc/applyMultipleId", config.GetPreUrl("idsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["code"] = req.Code
	queryParams["preCode"] = req.PreCode
	queryParams["count"] = req.Count
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ApplyMultiplePrimaryId(req idvo.ApplyMultiplePrimaryIdReqVo) idvo.ApplyMultipleIdRespVo {
	respVo := &idvo.ApplyMultipleIdRespVo{}
	reqUrl := fmt.Sprintf("%s/api/idsvc/applyMultiplePrimaryId", config.GetPreUrl("idsvc"))
	queryParams := map[string]interface{}{}
	queryParams["code"] = req.Code
	queryParams["count"] = req.Count
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ApplyMultiplePrimaryIdByCodes(req idvo.ApplyMultiplePrimaryIdByCodesReqVo) idvo.ApplyMultipleIdRespByCodesVo {
	respVo := &idvo.ApplyMultipleIdRespByCodesVo{}
	reqUrl := fmt.Sprintf("%s/api/idsvc/applyMultiplePrimaryIdByCodes", config.GetPreUrl("idsvc"))
	requestBody := &req.CodeInfos
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ApplyPrimaryId(req idvo.ApplyPrimaryIdReqVo) idvo.ApplyPrimaryIdRespVo {
	respVo := &idvo.ApplyPrimaryIdRespVo{}
	reqUrl := fmt.Sprintf("%s/api/idsvc/applyPrimaryId", config.GetPreUrl("idsvc"))
	queryParams := map[string]interface{}{}
	queryParams["code"] = req.Code
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}
