package projectvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type ImportIssuesReqVo struct {
	UserId int64              `json:"userId"`
	OrgId  int64              `json:"orgId"`
	Input  vo.ImportIssuesReq `json:"data"`
}

type ExportIssueReqVo struct {
	UserId int64                `json:"userId"`
	OrgId  int64                `json:"orgId"`
	Input  ExportIssueReqVoData `json:"input"`
}

type ExportIssueReqVoData struct {
	ProjectId      int64             `json:"projectId"` // 项目 id
	TableId        string            `json:"tableId"`   // 表 id
	ViewId         string            `json:"viewId"`    // 视图 id
	IterationId    int64             `json:"iterationId"`
	FilterColumns  []string          `json:"filterColumns"` // 导出的列
	Condition      *vo.LessCondsData `json:"condition"`     // 筛选条件
	Orders         []*vo.LessOrder   `json:"orders"`        // 排序
	Page           *int              `json:"page"`          // 页码
	Size           *int              `json:"size"`
	IsNeedDocument bool              `json:"isNeedDocument"` // 是否需要导出附件
}

type ExportIssueTemplateReqVo struct {
	UserId              int64 `json:"userId"`
	OrgId               int64 `json:"orgId"`
	ProjectId           int64 `json:"projectId"`
	ProjectObjectTypeId int64 `json:"projectObjectTypeId"`
	IterationId         int64 `json:"iterationId"`
	TableId             int64 `json:"tableId"`
}

type ExportIssueTemplateRespVo struct {
	vo.Err
	Data *vo.ExportIssueTemplateResp `json:"data"`
}

type ImportIssuesRespVo struct {
	vo.Err
	Data ImportIssuesRespVoData `json:"data"`
}

type ImportIssuesRespVoData struct {
	Count                int64    `json:"count"`
	TaskId               string   `json:"taskId"`
	UnknownImportColumns []string `json:"unknownImportColumns"`
}

type ImportIssuesErrorRespVo struct {
	vo.Err
	Data interface{} `json:"data"`
}

type DeleteAgileProjectExcelReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
	ProcessId int64 `json:"processId"`
}

type ExportUserOrDeptSameNameListReqVo struct {
	OrgId     int64 `json:"orgId"`
	UserId    int64 `json:"userId"`
	ProjectId int64 `json:"projectId"`
}

// IssueInfoMapForExportIssue 导出任务数据时，需要查询的信息
type IssueInfoMapForExportIssue struct {
	ProjectInfo   *bo.ProjectBo `json:"projectInfo"`
	AllColumnKeys []string      `json:"allColumnKeys"`
}

type GetUserOrDeptSameNameListReq struct {
	OrgId    int64  `json:"orgId"`
	UserId   int64  `json:"userId"`
	DataType string `json:"dataType"` // 获取的数据类型. user 同名用户;dept 同名部门
}

type GetUserOrDeptSameNameListRespVo struct {
	vo.Err
	Data GetUserOrDeptSameNameListRespData `json:"list"`
}

type GetUserOrDeptSameNameListRespData struct {
	List []GetUserOrDeptSameNameListItem `json:"list"`
}

type GetUserOrDeptSameNameListItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ImportIssues2022ReqVo struct {
	UserId int64                     `json:"userId"`
	OrgId  int64                     `json:"orgId"`
	Input  ImportIssues2022ReqVoData `json:"input"`
}

type ImportIssues2022ReqVoData struct {
	List       []map[string]interface{} `json:"list"`
	ColumnIds  []string                 `json:"columnIds"`
	AppIdStr   string                   `json:"appIdStr"`
	TableIdStr string                   `json:"tableIdStr"`
	Action     string                   `json:"action"` // 新增/更新：create/update
}
