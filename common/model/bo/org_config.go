package bo

import "time"

type OrgConfigBo struct {
	Id                         int64     `db:"id,omitempty" json:"id"`
	OrgId                      int64     `db:"org_id,omitempty" json:"orgId"`
	TimeZone                   string    `db:"time_zone,omitempty" json:"timeZone"`
	TimeDifference             string    `db:"time_difference,omitempty" json:"timeDifference"`
	PayLevel                   int       `db:"pay_level,omitempty" json:"payLevel"`
	PayStartTime               time.Time `db:"pay_start_time,omitempty" json:"payStartTime"`
	PayEndTime                 time.Time `db:"pay_end_time,omitempty" json:"payEndTime"`
	Seats                      int       `db:"seats,omitempty" json:"seats"`
	BasicShowSetting           string    `db:"basic_show_setting,omitempty" json:"basicShowSetting"`
	WebSite                    string    `db:"web_site,omitempty" json:"webSite"`
	Language                   string    `db:"language,omitempty" json:"language"`
	RemindSendTime             string    `db:"remind_send_time,omitempty" json:"remindSendTime"`
	ProjectDailyReportSendTime string    `db:"project_daily_report_send_time,omitempty" json:"projectDailyReportSendTime"`
	DatetimeFormat             string    `db:"datetime_format,omitempty" json:"datetimeFormat"`
	PasswordLength             int       `db:"password_length,omitempty" json:"passwordLength"`
	PasswordRule               int       `db:"password_rule,omitempty" json:"passwordRule"`
	MaxLoginFailCount          int       `db:"max_login_fail_count,omitempty" json:"maxLoginFailCount"`
	Status                     int       `db:"status,omitempty" json:"status"`
	Creator                    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime                 time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator                    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime                 time.Time `db:"update_time,omitempty" json:"updateTime"`
}
