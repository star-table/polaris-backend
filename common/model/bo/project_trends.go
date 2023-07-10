package bo

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"time"
)

type ProjectTrendsBo struct {
	PushType              consts.IssueNoticePushType //推送类型
	OrgId                 int64
	ProjectId             int64
	OperatorId            int64
	BeforeChangeMembers   []int64
	AfterChangeMembers    []int64
	BeforeOwner           int64
	AfterOwner            int64
	BeforeOwnerIds        []int64
	AfterOwnerIds         []int64
	BeforeChangeFollowers []int64
	AfterChangeFollowers  []int64
	FieldId               int64

	SourceChannel string //来源通道

	OperateObjProperty string
	NewValue           string
	OldValue           string
	Ext                TrendExtensionBo
	OperateTime        time.Time `json:"operateTime"`
}
