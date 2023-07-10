package bo

import (
	"github.com/star-table/common/core/types"
	"github.com/star-table/polaris-backend/common/core/consts"
)

//
//type IssueNoticeBo struct {
//	UpdateUserName string
//	UpdateTime     date.Time
//	IssueTitle     string
//	IssueId        int64
//	StatusName     string
//	NoticeTitle    string
//	NoticeContent  string
//}

type IssueNoticeBo struct {
	PushType   consts.IssueNoticePushType `json:"pushType"` //推送类型
	OrgId      int64                      `json:"orgId"`
	OperatorId int64                      `json:"operatorId"`
	DataId     int64                      `json:"dataId"`
	IssueId    int64                      `json:"issueId"`
	ProjectId  int64                      `json:"projectId"`
	PriorityId int64                      `json:"priorityId"`
	ParentId   int64                      `json:"parentId"` //父任务id
	TableId    int64                      `json:"tableId"`

	IssueTitle         string      `json:"issueTitle"`         //更新后的任务标题
	IssueRemark        string      `json:"issueRemark"`        //任务描述
	IssueStatusId      int64       `json:"issueStatusId"`      //更新后的任务状态
	IssuePlanStartTime *types.Time `json:"issuePlanStartTime"` //更新后的任务开始时间
	IssuePlanEndTime   *types.Time `json:"issuePlanEndTime"`   //更新后的任务结束时间

	SourceChannel string `json:"sourceChannel"` //来源通道

	BeforeOwner              []int64 `json:"beforeOwner"`
	AfterOwner               []int64 `json:"afterOwner"`
	BeforeChangeFollowers    []int64 `json:"beforeChangeFollowers"`
	AfterChangeFollowers     []int64 `json:"afterChangeFollowers"`
	BeforeChangeParticipants []int64 `json:"beforeChangeParticipants"`
	AfterChangeParticipants  []int64 `json:"afterChangeParticipants"`
	IssueChildren            []int64 `json:"issueChildren"` //相关联的子任务
	UpdateOwner              bool    `json:"updateOwner"`
	UpdateFollower           bool    `json:"updateFollower"`
}
