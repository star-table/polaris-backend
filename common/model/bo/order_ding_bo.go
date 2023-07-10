package bo

import (
	"time"
)

type OrderDingBo struct {
	Id                 int64     `db:"id,omitempty" json:"id"`
	OrgId              int64     `db:"org_id,omitempty" json:"orgId"`
	OutOrgId           string    `db:"out_org_id,omitempty" json:"outOrgId"`
	OrderId            string    `db:"order_id,omitempty" json:"orderId"`
	GoodsName          string    `db:"goods_name,omitempty" json:"goodsName"`
	GoodsCode          string    `db:"goods_code,omitempty" json:"goodsCode"`
	ItemName           string    `db:"item_name,omitempty" json:"itemName"`
	ItemCode           string    `db:"item_code,omitempty" json:"itemCode"`
	Quantity           int       `db:"quantity,omitempty" json:"quantity"`
	OrderPayPrice      int64     `db:"order_pay_price,omitempty" json:"orderPayPrice"`
	OrderType          string    `db:"order_type,omitempty" json:"orderType"`
	PaidTime           time.Time `db:"paid_time,omitempty" json:"paidTime"`
	OrderCreateTime    time.Time `db:"order_create_time,omitempty" json:"orderCreateTime"`
	StartTime          time.Time `db:"start_time,omitempty" json:"startTime"`
	EndTime            time.Time `db:"end_time,omitempty" json:"endTime"`
	Status             int       `db:"status,omitempty" json:"status"`
	OrderSource        string    `db:"order_source,omitempty" json:"orderSource"`
	SaleModelType      string    `db:"sale_model_type,omitempty" json:"saleModelType"`
	AutoChangeFreeItem int       `db:"auto_change_free_item" json:"autoChangeFreeItem"`
	OrderChargeType    string    `db:"order_charge_type" json:"orderChargeType"`
	Creator            int64     `db:"creator,omitempty" json:"creator"`
	CreateTime         time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator            int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime         time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version            int       `db:"version,omitempty" json:"version"`
	IsDelete           int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*OrderDingBo) TableName() string {
	return "ppm_ord_order_ding"
}

type OrderWeiXinBo struct {
	Id            int64     `db:"id,omitempty" json:"id"`
	OrgId         int64     `db:"org_id,omitempty" json:"orgId"`
	OutOrgId      string    `db:"out_org_id,omitempty" json:"outOrgId"`
	OrderId       string    `db:"order_id,omitempty" json:"orderId"`
	EditionId     string    `db:"edition_id,omitempty" json:"editionId"`
	EditionName   string    `db:"edition_name,omitempty" json:"editionName"`
	OrderType     int       `db:"order_type,omitempty" json:"orderType"`
	OrderStatus   int       `db:"order_status,omitempty" json:"orderStatus"`
	UserCount     int       `db:"user_count,omitempty" json:"userCount"`
	OrderPeriod   int       `db:"order_period,omitempty" json:"orderPeriod"`
	PaidTime      time.Time `db:"paid_time,omitempty" json:"paidTime"`
	BeginTime     time.Time `db:"begin_time,omitempty" json:"beginTime"`
	EndTime       time.Time `db:"end_time,omitempty" json:"endTime"`
	OrderPayPrice int64     `db:"order_pay_price,omitempty" json:"orderPayPrice"`
	IsDelete      int       `db:"is_delete,omitempty" json:"isDelete"`
}
