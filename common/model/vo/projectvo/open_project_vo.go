package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type OpenCreateProjectReq struct {
	vo.CreateProjectReq

	// 操作人
	OperatorId int64 `json:"operatorId"`
}

type OpenProjectsReq struct {
	vo.ProjectsReq
	OperatorId int64 `json:"operatorId"`
}

type OpenUpdateProjectReq struct {
	vo.UpdateProjectReq
	// 操作人
	OperatorId int64 `json:"operatorId"`
}

type OpenOperatorReq struct {
	// 操作人
	OperatorId int64 `json:"operatorId"`
}

type TablesColumnsReq struct {
	TableIds  []int64  `json:"tableIds"`
	ColumnIds []string `json:"columnIds"`
}

type OpenTablesColumnsReq struct {
	TablesColumnsReq
	// 操作人
	OperatorId int64 `json:"operatorId"`
}

type TableBaseColumnField struct {
	Type string `json:"type,omitempty"`
	//CustomType string `json:"customType"`
	DataType string                 `json:"dataType,omitempty"`
	Props    map[string]interface{} `json:"props"`
}

type TableBaseColumnData struct {
	Name  string               `json:"name"`
	Field TableBaseColumnField `json:"field"`
}

type TableBaseColumnsTable struct {
	AppId   int64                  `json:"appId,string"`
	TableId int64                  `json:"tableId,string"`
	Name    string                 `json:"name"`
	Columns []*TableBaseColumnData `json:"columns"`
}

type TablesBaseColumnsRespData struct {
	Tables []*TableBaseColumnsTable `json:"tables"`
}

type TablesBaseColumnsResp struct {
	vo.Err
	Data *TablesBaseColumnsRespData `json:"data"`
}

type OpenGetIssueStatusListReq struct {
	TablesColumnsReq
	// 操作人
	OperatorId int64 `json:"operatorId"`
}

type IssueTableColumnField struct {
	Type     string                 `json:"type,omitempty"`
	DataType string                 `json:"dataType,omitempty"`
	Props    map[string]interface{} `json:"props"`
}

type IssueStatusColumnData struct {
	Name  string                `json:"name"`
	Field IssueTableColumnField `json:"field"`
}

type IssueStatusListColumnsTable struct {
	TableId int64       `json:"tableId,string"`
	Name    string      `json:"name"`
	Props   interface{} `json:"props"`
}

type IssueStatusListRespData struct {
	Tables []*IssueStatusListColumnsTable `json:"tables"`
}

type OpenGetIssueStatusListResp struct {
	vo.Err
	Data *IssueStatusListRespData `json:"data"`
}

type IssueStatusTables struct {
}
