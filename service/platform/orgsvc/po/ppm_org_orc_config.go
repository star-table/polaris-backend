package po

type ScheduleOrganizationListPo struct {
	OrgId                      int64  `db:"id,omitempty" json:"orgId"`
	ProjectDailyReportSendTime string `db:"project_daily_report_send_time,omitempty" json:"projectDailyReportSendTime"`
}
