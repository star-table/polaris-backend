package bo

import (
"time"
)

type OrderFsBo struct {
	Id              int64     `db:"id,omitempty" json:"id"`
	OrgId           int64     `db:"org_id,omitempty" json:"orgId"`
	OrderId         string    `db:"order_id,omitempty" json:"orderId"`
	PricePlanId     string    `db:"price_plan_id,omitempty" json:"pricePlanId"`
	PricePlanType   string    `db:"price_plan_type,omitempty" json:"pricePlanType"`
	Seats           int       `db:"seats,omitempty" json:"seats"`
	BuyCount        int       `db:"buy_count,omitempty" json:"buyCount"`
	PaidTime        time.Time `db:"paid_time,omitempty" json:"paidTime"`
	OrderCreateTime time.Time `db:"order_create_time,omitempty" json:"orderCreateTime"`
	Status          string    `db:"status,omitempty" json:"status"`
	BuyType         string    `db:"buy_type,omitempty" json:"buyType"`
	SrcOrderId      string    `db:"src_order_id,omitempty" json:"srcOrderId"`
	DstOrderId      string    `db:"dst_order_id,omitempty" json:"dstOrderId"`
	OrderPayPrice   int64     `db:"order_pay_price,omitempty" json:"orderPayPrice"`
	TenantKey       string    `db:"tenant_key,omitempty" json:"tenantKey"`
	Creator         int64     `db:"creator,omitempty" json:"creator"`
	CreateTime      time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator         int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime      time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version         int       `db:"version,omitempty" json:"version"`
	IsDelete        int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*OrderFsBo) TableName() string {
	return "ppm_ord_order_fs"
}

