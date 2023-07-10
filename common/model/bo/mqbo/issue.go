package mqbo

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
)

type PushBatchCreateIssue struct {
	TraceId     string                           `json:"traceId"`
	WithoutAuth bool                             `json:"withoutAuth"`
	TriggerBy   *projectvo.TriggerBy             `json:"triggerBy"`   // normal,import,applyTemplate
	AsyncTaskId string                           `json:"asyncTaskId"` // 异步任务id（进度更新/查询用）
	Count       int64                            `json:"count"`       // 创建总数（含子任务）
	Req         *projectvo.BatchCreateIssueReqVo `json:"req"`
}

type PushCreateIssueBo struct {
	CreateIssueReqVo projectvo.CreateIssueReqVo `json:"createIssueReqVo"`
	IssueId          int64                      `json:"issueId"`
	// traceId 导入任务时，保证使用同一个 traceId
	TraceId string `json:"traceId"`
}

type DailyTaskMqBo struct {
	//组织id
	OrgId int64 `json:"orgId"`
	//人员id
	UserId int64 `json:"userId"`
	//平台
	PlatForm string `json:"platForm"`
	//飞书组织tenaut
	Tenant interface{} `json:"tenant"`
	//钉钉组织
	//DingTalkClient sdk2.DingTalkClient `json:"dingTalkClient"`
	//traceId
	TraceId string `json:"traceId"`
	//推送类型
	PushType consts.IssueNoticePushType `json:"pushType"`
	//外部组织信息
	OutOrgId string `json:"outOrgId"`
	// 额外信息
	ExtInfo map[string]interface{} `json:"extInfo"`
}

type IterationStatMqBo struct {
	IterBo bo.IterationBo `json:"iterBo"`
	//统计时间
	StatDate string `json:"statDate"`
	//traceId
	TraceId string `json:"traceId"`
	//推送类型
	PushType consts.IssueNoticePushType `json:"pushType"`
}

type ProjectStatMqBo struct {
	ProjectBo bo.ProjectBo `json:"projectBo"`
	//统计时间
	StatDate string `json:"statDate"`
	//traceId
	TraceId string `json:"traceId"`
	//推送类型
	PushType consts.IssueNoticePushType `json:"pushType"`
}
