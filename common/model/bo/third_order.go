package bo

import (
	"github.com/star-table/polaris-backend/common/core/consts"
)

type ThirdOrderPush struct {
	//traceId
	TraceId string `json:"traceId"`
	//推送类型
	PushType      consts.IssueNoticePushType `json:"pushType"` //推送类型
	Msg           OrderDingBo                `json:"msg"`
	SourceChannel string                     `json:"sourceChannel"`
}

type PayRangeData struct {
	OrgId           int64   `json:"orgId"`
	OutOrgId        string  `json:"outOrgId"`
	SourceChannel   string  `json:"sourceChannel"`
	ScopeNum        int     `json:"scopeNum"`        // 授权可见范围人数
	PayNum          int     `json:"payNum"`          // 付费人数
	EndTime         string  `json:"endTime"`         // 付费版本结束时间
	UserIds         []int64 `json:"userIds"`         // 已经登录过极星的userId
	IsFsSetPayRange int     `json:"isFsSetPayRange"` // 飞书是否设置过了可见范围
}

type OrgAndUserInfo struct {
	OrgName       string `json:"orgName"`
	UserName      string `json:"userName"`
	SourceChannel string `json:"sourceChannel"`
}
