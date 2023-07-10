package bo

import (
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/model/bo/status"
)

type IterationBo struct {
	Id            int64      `db:"id,omitempty" json:"id"`
	OrgId         int64      `db:"org_id,omitempty" json:"orgId"`
	ProjectId     int64      `db:"project_id,omitempty" json:"projectId"`
	Name          string     `db:"name,omitempty" json:"name"`
	Sort          int64      `db:"sort,omitempty" json:"sort"`
	Owner         int64      `db:"owner,omitempty" json:"owner"`
	VersionId     int64      `db:"version_id,omitempty" json:"versionId"`
	PlanStartTime types.Time `db:"plan_start_time,omitempty" json:"planStartTime"`
	PlanEndTime   types.Time `db:"plan_end_time,omitempty" json:"planEndTime"`
	PlanWorkHour  int        `db:"plan_work_hour,omitempty" json:"planWorkHour"`
	StoryPoint    int        `db:"story_point,omitempty" json:"storyPoint"`
	Remark        string     `db:"remark,omitempty" json:"remark"`
	Status        int64      `db:"status,omitempty" json:"status"`
	Creator       int64      `db:"creator,omitempty" json:"creator"`
	CreateTime    types.Time `db:"create_time,omitempty" json:"createTime"`
	Updator       int64      `db:"updator,omitempty" json:"updator"`
	UpdateTime    types.Time `db:"update_time,omitempty" json:"updateTime"`
	Version       int        `db:"version,omitempty" json:"version"`
	IsDelete      int        `db:"is_delete,omitempty" json:"isDelete"`
}

type IterationStatusTypeCountBo struct {
	NotStartTotal   int64 `json:"notStartTotal"`
	ProcessingTotal int64 `json:"processingTotal"`
	FinishedTotal   int64 `json:"finishedTotal"`
}

type IterationStatusTypeCountSelectBo struct {
	Count int64 `db:"count"`
	Type  int   `db:"type"`
}

type IterationUpdateBo struct {
	Id  int64
	Upd mysql.Upd

	IterationNewBo IterationBo
}

type IterationInfoBo struct {
	// 迭代信息
	Iteration IterationBo `json:"iteration"`
	// 项目信息
	Project *HomeIssueProjectInfoBo `json:"project"`
	// 状态信息
	Status *status.StatusInfoBo `json:"status"`
	// 负责人信息
	Owner *UserIDInfoBo `json:"owner"`
	// 下一步骤状态列表
	NextStatus *[]status.StatusInfoBo `json:"nextStatus"`
}

func (*IterationBo) TableName() string {
	return "ppm_pri_iteration"
}
