package commonfacade

import (
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/facade"
)


func AreaInfo(req commonvo.AreaInfoReqVo) commonvo.AreaInfoRespVo {
	respVo := &commonvo.AreaInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/commonsvc/areaInfo", config.GetPreUrl("commonsvc"))
	queryParams := map[string]interface{}{}
	queryParams["industryID"] = req.IndustryID
	queryParams["countryID"] = req.CountryID
	queryParams["provinceID"] = req.ProvinceID
	queryParams["cityID"] = req.CityID
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AreaLinkageList(req commonvo.AreaLinkageListReqVo) commonvo.AreaLinkageListRespVo {
	respVo := &commonvo.AreaLinkageListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/commonsvc/areaLinkageList", config.GetPreUrl("commonsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DingTalkInfo(req commonvo.DingTalkInfoReqVo) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/commonsvc/dingTalkInfo", config.GetPreUrl("commonsvc"))
	requestBody := &req.LogData
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UploadOssByFsImageKey(req commonvo.UploadOssByFsImageKeyReq) commonvo.UploadOssByFsImageKeyResp {
	respVo := &commonvo.UploadOssByFsImageKeyResp{}
	reqUrl := fmt.Sprintf("%s/api/commonsvc/uploadOssByFsImageKey", config.GetPreUrl("commonsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["imageKey"] = req.ImageKey
	queryParams["isApp"] = req.IsApp
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

