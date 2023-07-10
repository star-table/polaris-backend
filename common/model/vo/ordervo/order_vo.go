package ordervo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type AddFsOrderReq struct {
	Data bo.OrderFsBo `json:"data"`
}

type AddDingOrderReq struct {
	Data bo.OrderDingBo `json:"data"`
}

type AddWeiXinOrderReq struct {
	Data bo.OrderWeiXinBo `json:"data"`
}

type FixOrderOfZeroOrgIdReq struct {
	Input FixOrderOfZeroOrgIdReqInput `json:"input"`
}

type FixOrderOfZeroOrgIdReqInput struct {
}

type CheckFsOrgIsInTrialReq struct {
	OrgId int64 `json:"orgId"`
}

type UpdateDingOrderReq struct {
	Input UpdateDingOrderData `json:"input"`
}

type UpdateDingOrderData struct {
	OrgId    int64  `json:"orgId"`
	OutOrgId string `json:"outOrgId"`
}

type UpdateDingOrderDataResp struct {
	vo.Err
	Data *UpdateDingOrderRespData `json:"data"`
}

type UpdateDingOrderRespData struct {
	Level    int64          `json:"level"`
	DingData bo.OrderDingBo `json:"dingData"`
}

type DeleteDingOrderReqVo struct {
	Data DeleteDingOrderReq `json:"data"`
}

type DeleteDingOrderReq struct {
	OutOrgId string `json:"outOrgId"`
	OrderId  string `json:"orderId"`
}

type CreateWeiXinLicenceOrderReq struct {
	SuitId  string `json:"suit_id"`
	CorpId  string `json:"corp_id"`
	OrderId string `json:"order_id"`
}

type GetOrderPayInfoReq struct {
	OrgId int64 `json:"orgId"`
}

type GetOrderPayInfoResp struct {
	vo.Err
	Data *GetOrderPayInfo `json:"data"`
}

type GetOrderPayInfo struct {
	PayNum int `json:"payNum"`
}
