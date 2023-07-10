package commonvo

import "github.com/star-table/polaris-backend/common/model/vo"

type AreaLinkageListReqVo struct {
	Input vo.AreaLinkageListReq
}

type AreaLinkageListRespVo struct {
	vo.Err
	AreaLinkageListResp *vo.AreaLinkageListResp `json:"data"`
}

type AreaInfoReqVo struct {
	// 所属行业
	IndustryID int64 `json:"industryId"`
	// 所在国家
	CountryID int64 `json:"countryId"`
	// 所在省份
	ProvinceID int64 `json:"provinceId"`
	// 所在城市
	CityID int64 `json:"cityId"`
}

type AreaInfoRespVo struct {
	vo.Err
	AreaInfoResp *AreaInfoResp `json:"data"`
}

type AreaInfoResp struct {
	// 所属行业中文名
	IndustryName string `json:"industryName"`
	// 所在国家中文名
	CountryCname string `json:"countryCname"`
	// 所在省份中文名
	ProvinceCname string `json:"provinceCname"`
	// 所在城市中文名
	CityCname string `json:"cityCname"`
}

type DingTalkInfoReqVo struct {
	LogData map[string]interface{} `json:"logData"`
}
