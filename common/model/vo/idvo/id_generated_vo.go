package idvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

//申请ID
type ApplyPrimaryIdRespVo struct {
	vo.Err
	Id int64 `json:"data"`
}

type ApplyPrimaryIdReqVo struct {
	Code string `json:"code"`
}

//申请Code
type ApplyCodeRespVo struct {
	vo.Err
	Code string `json:"data"`
}

type ApplyCodeReqVo struct {
	OrgId   int64  `json:"orgId"`
	Code    string `json:"code"`
	PreCode string `json:"preCode"`
}

//批量申请
type ApplyMultipleIdRespVo struct {
	vo.Err
	IdCodes *bo.IdCodes `json:"data"`
}

type ApplyMultipleIdRespByCodesVo struct {
	vo.Err
	CodesIds *bo.CodesIds `json:"data"`
}

type ApplyMultipleIdReqVo struct {
	OrgId   int64  `json:"orgId"`
	Code    string `json:"code"`
	PreCode string `json:"preCode"`
	Count   int    `json:"count"`
}

type ApplyMultiplePrimaryIdReqVo struct {
	Code  string `json:"code"`
	Count int    `json:"count"`
}

type ApplyMultiplePrimaryIdByCodesReqVo struct {
	CodeInfos []ApplyMultiplePrimaryIdReqVo `json:"codeInfos"`
}
