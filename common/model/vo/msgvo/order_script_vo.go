package msgvo

import "time"

type FixAddOrderForFeiShuReqVo struct {
	Input FixAddOrderForFeiShuReqData `json:"input"`
}

type FixAddOrderForFeiShuReqData struct {
	StartTime *time.Time `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
}

type OrderContent struct {
	Type string `json:"type"`
	Data string `json:"data"`
	MsgTime string `json:"msgTime"`
}

type OrderContentData struct {
	Uuid string `json:"uuid"`
	Event OrderContentDataEvent `json:"event"`
	Token string `json:"token"`
	Ts string `json:"ts"`
	Type string `json:"type"`
}

type OrderContentDataEvent struct {
	AppId         string `json:"app_id"`
	BuyCount      int  `json:"buy_count"`
	BuyType       string `json:"buy_type"`
	CreateTime    string `json:"create_time"`
	OrderId       string `json:"order_id"`
	OrderPayPrice int64  `json:"order_pay_price"`
	PayTime       string `json:"pay_time"`
	PricePlanId   string `json:"price_plan_id"`
	PricePlanType string `json:"price_plan_type"`
	Seats         int  `json:"seats"`
	SrcOrderId    string `json:"src_order_id"`
	TenantKey     string `json:"tenant_key"`
	Type          string `json:"type"`
}
