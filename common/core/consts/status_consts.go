package consts

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo/status"
)

const (
	// 项目和迭代 固定的状态值
	StatusIdNotStart = 1 // 未开始
	StatusIdRunning  = 2 // 进行中
	StatusIdComplete = 3 // 已完成

	// 项目和迭代 固定的状态分类
	StatusTypeNotStart = 1 // 未开始
	StatusTypeRunning  = 2 // 进行中
	StatusTypeComplete = 3 // 已完成

	// 项目和迭代 固定的状态名
	StatusNameNotStart = "未开始"
	StatusNameRunning  = "进行中"
	StatusNameComplete = "已完成"

	// 项目
	StatusLangCodeProjectNotStart = "ProcessStatus.Project.NotStart"
	StatusLangCodeProjectRunning  = "ProcessStatus.Project.Running"
	StatusLangCodeProjectComplete = "ProcessStatus.Project.Complete"

	// 迭代
	StatusLangCodeIterationNotStart = "ProcessStatus.Iteration.NotStart"
	StatusLangCodeIterationRunning  = "ProcessStatus.Iteration.Running"
	StatusLangCodeIterationComplete = "ProcessStatus.Iteration.Complete"

	// 任务
	StatusLangCodeIssueNotStart       = "ProcessStatus.Issue.NotStart"
	StatusLangCodeIssueWaitEvaluate   = "ProcessStatus.Issue.WaitEvaluate"
	StatusLangCodeIssueWaitConfirmBug = "ProcessStatus.Issue.WaitConfirmBug"
	StatusLangCodeIssueConfirmedBug   = "ProcessStatus.Issue.ConfirmedBug"
	StatusLangCodeIssueReOpen         = "ProcessStatus.Issue.ReOpen"
	StatusLangCodeIssueEvaluated      = "ProcessStatus.Issue.Evaluated"
	StatusLangCodeIssuePlanning       = "ProcessStatus.Issue.Planning"
	StatusLangCodeIssueDesign         = "ProcessStatus.Issue.Design"
	StatusLangCodeIssueProcessing     = "ProcessStatus.Issue.Processing"
	StatusLangCodeIssueDevelopment    = "ProcessStatus.Issue.Development"
	StatusLangCodeIssueWaitTest       = "ProcessStatus.Issue.WaitTest"
	StatusLangCodeIssueTesting        = "ProcessStatus.Issue.Testing"
	StatusLangCodeIssueWaitRelease    = "ProcessStatus.Issue.WaitRelease"
	StatusLangCodeIssueWait           = "ProcessStatus.Issue.Wait"
	StatusLangCodeIssueRepair         = "ProcessStatus.Issue.Repair"
	StatusLangCodeIssueSuccess        = "ProcessStatus.Issue.Success"
	StatusLangCodeIssueFail           = "ProcessStatus.Issue.Fail"
	StatusLangCodeIssueReleased       = "ProcessStatus.Issue.Released"
	StatusLangCodeIssueClosed         = "ProcessStatus.Issue.Closed"
	StatusLangCodeIssueComplete       = "ProcessStatus.Issue.Complete"
	StatusLangCodeIssueConfirmed      = "ProcessStatus.Issue.Confirmed"
)

const (
	FontStyleNotStart = "#5f5f5f"
	FontStyleRunning  = "#377aff"
	FontStyleComplete = "#54a944"
	BgStyleNotStart   = "#f5f6f5"
	BgStyleRunning    = "#ecf5fc"
	BgStyleComplete   = "#edf8ed"
)

const (

	//-1: 待确认 1:未完成，2：已完成，3：未开始，4：进行中 5: 已逾期
	// 待确认
	IssueStatusNone = -1
	// 未完成
	IssueStatusUnfinished = 1
	// 已完成
	IssueStatusFinished = 2
	// 未开始
	IssueStatusNotStarted = 3
	// 进行中
	IssueStatusProcessing = 4
	// 已逾期
	IssueStatusOverdue = 5
)

var StatusNotStart = status.StatusInfoBo{
	ID:          StatusIdNotStart,
	Type:        StatusTypeNotStart,
	Name:        StatusNameNotStart,
	DisplayName: StatusNameNotStart,
	FontStyle:   FontStyleNotStart,
	BgStyle:     BgStyleNotStart,
	Sort:        1,
}

var StatusRunning = status.StatusInfoBo{
	ID:          StatusIdRunning,
	Type:        StatusTypeRunning,
	Name:        StatusNameRunning,
	DisplayName: StatusNameRunning,
	FontStyle:   FontStyleRunning,
	BgStyle:     BgStyleRunning,
	Sort:        2,
}

var StatusComplete = status.StatusInfoBo{
	ID:          StatusIdComplete,
	Type:        StatusTypeComplete,
	Name:        StatusNameComplete,
	DisplayName: StatusNameComplete,
	FontStyle:   FontStyleComplete,
	BgStyle:     BgStyleComplete,
	Sort:        3,
}

// 默认的迭代状态对象
var StatusNotStartOfIter = status.StatusInfoBo{
	ID:          StatusIdNotStart,
	Type:        StatusTypeNotStart,
	Name:        StatusNameNotStart,
	DisplayName: StatusNameNotStart,
	FontStyle:   "#FFFFFF",
	BgStyle:     "#377AFF",
	Sort:        1,
}

var StatusRunningOfIter = status.StatusInfoBo{
	ID:          StatusIdRunning,
	Type:        StatusTypeRunning,
	Name:        StatusNameRunning,
	DisplayName: StatusNameRunning,
	FontStyle:   "#FFFFFF",
	BgStyle:     "#F0A100",
	Sort:        2,
}

var StatusCompleteOfIter = status.StatusInfoBo{
	ID:          StatusIdComplete,
	Type:        StatusTypeComplete,
	Name:        StatusNameComplete,
	DisplayName: StatusNameComplete,
	FontStyle:   "#FFFFFF",
	BgStyle:     "#25B47E",
	Sort:        3,
}

// ProjectStatusList 项目状态
var ProjectStatusList = []status.StatusInfoBo{StatusRunning, StatusComplete}

// IterationStatusList 迭代状态
var IterationStatusList = []status.StatusInfoBo{StatusNotStartOfIter, StatusRunningOfIter, StatusCompleteOfIter}

// DefaultProTypeMap2TableList 创建项目时，不同项目类型，默认创建的 table list，在这里进行配置：
// map: {projectTypeId: []ProTable}
var DefaultProTypeMap2TableList = map[int64][]ProTable{
	1: []ProTable{{Name: "任务"}},                             // 通用
	2: []ProTable{{Name: "需求"}, {Name: "任务"}, {Name: "缺陷"}}, // 敏捷
}

type ProTable struct {
	Name string `json:"name"`
}

// GetStatusIdsByCategory 根据 category 和 typ 获取对应的状态 ids
// 任务状态改造后，这个方法不再适用于查询任务的状态，请使用 GetIssueStatusWithTypeGroupByTableId 函数（全局搜）
func GetStatusIdsByCategory(orgId int64, category int, statysType int) ([]int64, errs.SystemErrorInfo) {
	statusIds := make([]int64, 0)
	switch category {
	case ProcessStatusCategoryProject: // 项目
		allStatus := ProjectStatusList
		for _, item := range allStatus {
			if item.Type == statysType {
				statusIds = append(statusIds, item.ID)
			}
		}
	case ProcessStatusCategoryIteration: // 迭代
		allStatus := IterationStatusList
		for _, item := range allStatus {
			if item.Type == statysType {
				statusIds = append(statusIds, item.ID)
			}
		}
	case ProcessStatusCategoryIssue: // todo
		//allStatusGroup :=
	}
	//respVo := GetProcessStatusIds(processvo.GetProcessStatusIdsReqVo{OrgId: orgId, Typ: typ, Category: category})
	//if respVo.Failure() {
	//	return nil, respVo.Error()
	//}

	return statusIds, nil
}

func GetProjectStatusById(statusId int64) *status.StatusInfoBo {
	for _, item := range IterationStatusList {
		if item.ID == statusId {
			return &item
		}
	}
	return nil
}

func GetIterationStatusById(statusId int64) *status.StatusInfoBo {
	for _, item := range IterationStatusList {
		if item.ID == statusId {
			return &item
		}
	}
	return nil
}
