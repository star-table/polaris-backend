package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type ConvertCodeRespVo struct {
	vo.Err
	ConvertCode *vo.ConvertCodeResp `json:"data"`
}

type ConvertCodeReqVo struct {
	Input vo.ConvertCodeReq `json:"input"`
	OrgId int64             `json:"orgId"`
}
