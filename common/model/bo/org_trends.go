package bo

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"time"
)
//组织动态结构体
type OrgTrendsBo struct {
	OrgId                 int64
	PushType              consts.IssueNoticePushType //推送类型
	TargetMembers []int64	//被操作的成员，根据pushType去区分业务场景
	SourceChannel string	//来源通道
	OperatorId            int64
	OperateTime        time.Time `json:"operateTime"`
}

//组织动态结构体ext
type OrgTrendsBoExt struct {
	MemberInfo    []SimpleUserInfoBo  `json:"memberInfo"`
	OrgName string `json:"orgName"`
}