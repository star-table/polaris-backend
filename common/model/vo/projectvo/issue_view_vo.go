package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type CreateTaskViewReqVoOld struct {
	ProjectId int64                     `json:"projectId"`
	Config    CreateTaskViewReqVoConfig `json:"config"`
	IsPrivate bool                      `json:"isPrivate"`
	Remark    string                    `json:"remark"`
	ViewName  string                    `json:"viewName"`
	Type      int                       `json:"type"`
	Sort      int64                     `json:"sort"`
	OrgId     int64                     `json:"orgId"`
	CurUserId int64                     `json:"curUserId"`
}

type CreateTaskViewReqVo struct {
	Input  vo.CreateIssueViewReq `json:"input"`
	UserId int64                 `json:"userId"`
	OrgId  int64                 `json:"orgId"`
}

type CreateTaskViewReqVoConfig struct {
	Condition []CreateTaskViewReqVoConfigCondItem  `json:"condition"`
	Orders    []CreateTaskViewReqVoConfigOrderItem `json:"orders"`
}

type CreateTaskViewReqVoConfigCondItem struct {
	Column string      `json:"column"`
	Type   string      `json:"type"`
	Values interface{} `json:"values"`
}

type CreateTaskViewReqVoConfigOrderItem struct {
	Asc    bool   `json:"asc"`
	Column string `json:"column"`
}

type CreateTaskViewRespVo struct {
	vo.Err
	Data *vo.GetIssueViewListItem `json:"data"`
}

type GetOneTaskViewReqVo struct {
	Id int64 `json:"id"`
}

type GetOneTaskViewRespVo struct {
	vo.Err
	Data *vo.GetIssueViewListItem `json:"data"`
}

type GetOneTaskViewRespVoData struct {
	ProjectId int64  `json:"projectId"`
	Config    string `json:"config"`
	IsPrivate bool   `json:"isPrivate"`
	Remark    string `json:"remark"`
	ViewName  string `json:"viewName"`
	Type      int    `json:"type"`
	Sort      int64  `json:"sort"`
	OrgId     int64  `json:"orgId"`
}

type GetTaskViewListReqVo struct {
	Input  vo.GetIssueViewListReq `json:"input"`
	UserId int64                  `json:"userId"`
	OrgId  int64                  `json:"orgId"`
}

type GetOptionListReqVo struct {
	Input  GetOptionListReqInput `json:"input"`
	UserId int64                 `json:"userId"`
	OrgId  int64                 `json:"orgId"`
}

type GetOptionListReqInput struct {
	ProjectId           int64  `json:"projectId"`
	ProjectObjectTypeId int64  `json:"projectObjectTypeId"`
	FieldName           int64  `json:"fieldName"`
	TableId             string `json:"tableId"`
}

type GetOptionListRespVo struct {
	vo.Err
	Data []GetOptionListRespVoDataItem `json:"data"`
}

type GetOptionListRespVoDataItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type GetTaskViewListRespVo struct {
	vo.Err
	Data *vo.GetIssueViewListResp `json:"data"`
}

type UpdateTaskViewReqVo struct {
	Input  vo.UpdateIssueViewReq `json:"input"`
	UserId int64                 `json:"userId"`
	OrgId  int64                 `json:"orgId"`
}

type DeleteTaskViewReqVo struct {
	Input  vo.DeleteIssueViewReq `json:"input"`
	UserId int64                 `json:"userId"`
	OrgId  int64                 `json:"orgId"`
}

type MirrorsStatReq struct {
	UserId int64             `json:"userId"`
	OrgId  int64             `json:"orgId"`
	Input  vo.MirrorCountReq `json:"input"`
}

type MirrorsStatResp struct {
	vo.Err
	Data vo.MirrorsStatResp `json:"data"`
}
