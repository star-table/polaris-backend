package bo

import (
	"github.com/star-table/polaris-backend/common/core/consts"
)

//type DailyProjectReportProjectInfoMsgBo struct {
//	ScheduleTraceId string                     `json:"traceId"`
//	PushType        consts.IssueNoticePushType `json:"pushType"` //推送类型
//	OrgId           int64                      `json:"orgId"`
//	ProjectId       int64                      `json:"projectId"`
//	Owner           int64                      `json:"owner"`
//}

type DailyProjectReportMsgBo struct {
	ScheduleTraceId  string                     `json:"traceId"`
	PushType         consts.IssueNoticePushType `json:"pushType"` //推送类型
	OrgId            int64                      `json:"orgId"`
	SourceChannel    string                     `json:"sourceChannel"`
	OutOrgId         string                     `json:"outOrgId"`
	OpenIds          []string                   `json:"openIds"`
	AppId            int64                      `json:"appId"`
	ProjectId        int64                      `json:"projectId"` //项目id
	ProjectTypeId    int64                      `json:"projectTypeId"`
	ProjectName      string                     `json:"projectName"`         //项目名称
	DailyFinishCount int64                      `json:"dailyFinishCount"`    //今日完成
	RemainingCount   int64                      `json:"dailyRemainingCount"` //剩余未完成
	OverdueCount     int64                      `json:"dailyOverdueCount"`   //逾期任务数量
}

//type RetryPushMsg struct {
//	Num       int       `json:"num"`
//	OrgId     int64     `json:"orgId"`
//	RetryTime time.Time `json:"retryTime"`
//}

// 个人日报
type DailyIssueReportMsgBo struct {
	ScheduleTraceId           string                     `json:"traceId"`
	PushType                  consts.IssueNoticePushType `json:"pushType"`                  //推送类型
	OrgId                     int64                      `json:"orgId"`                     //org id
	SourceChannel             string                     `json:"sourceChannel"`             //source channel
	OutOrgId                  string                     `json:"outOrgId"`                  //out org id
	OpenIds                   []string                   `json:"openIds"`                   //open ids
	IssueOverdueCount         int64                      `json:"issueOverdueCount"`         // 已逾期数量
	IssueOverdueTodayCount    int64                      `json:"issueOverdueTodayCount"`    //今日逾期
	IssueOverdueTomorrowCount int64                      `json:"issueOverdueTomorrowCount"` // 即将逾期
	IssueToBeCompleted        int64                      `json:"issueToBeCompleted"`        // 待完成数量
}

type IssueRemindMsg struct {
	ScheduleTraceId string                     `json:"traceId"`
	PushType        consts.IssueNoticePushType `json:"pushType"`
	OrgId           int64                      `json:"orgId"`
	SourceChannel   string                     `json:"sourceChannel"`
	Operator        int64                      `json:"operator"`
	OwnerId         []int64                    `json:"ownerId"`
	OutOrgId        string                     `json:"outOrgId"` //out org id
	OpenIds         []string                   `json:"openIds"`  //open ids
	AppId           int64                      `json:"appId"`
	ProjectId       int64                      `json:"projectId"`
	TableId         int64                      `json:"tableId"`
	IssueId         int64                      `json:"issueId"`
	ParentId        int64                      `json:"parentId"`
	Title           string                     `json:"title"`
	ProjectName     string                     `json:"projectName"`
	TableName       string                     `json:"tableName"`
	PlanEndTime     string                     `json:"planEndTime"`
}
