package bo

import (
	"time"
)

type OrdPricePlanFs struct {
	Id            int64     `db:"id,omitempty" json:"id"`
	OutPlanId     string    `db:"out_plan_id,omitempty" json:"outPlanId"`
	Name          string    `db:"name,omitempty" json:"name"`
	PricePlanType string    `db:"price_plan_type,omitempty" json:"pricePlanType"`
	Seats         int       `db:"seats,omitempty" json:"seats"`
	TrialDays     int       `db:"trial_days,omitempty" json:"trialDays"`
	ExpireDays    int       `db:"expire_days,omitempty" json:"expireDays"`
	EndDate       time.Time `db:"end_date,omitempty" json:"endDate"`
	MonthPrice    int64     `db:"month_price,omitempty" json:"monthPrice"`
	YearPrice     int64     `db:"year_price,omitempty" json:"yearPrice"`
	Level         int64     `db:"level,omitempty" json:"level"`
	Creator       int64     `db:"creator,omitempty" json:"creator"`
	CreateTime    time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator       int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime    time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version       int       `db:"version,omitempty" json:"version"`
	IsDelete      int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*OrdPricePlanFs) TableName() string {
	return "ppm_ord_price_plan_fs"
}
