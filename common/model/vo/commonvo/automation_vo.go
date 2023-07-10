package commonvo

import (
	automationPb "github.com/star-table/interface/golang/automation/v1"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type GlobalSwitchOffResp struct {
	vo.Err
	Data *automationPb.GlobalSwitchOffReply `json:"data"`
}

// ApplyTemplate
type ApplyTemplateReq struct {
	OrgId  int64                          `json:"orgId"`
	UserId int64                          `json:"userId"`
	Input  *automationPb.ApplyTemplateReq `json:"input"`
}

type ApplyTemplateResp struct {
	vo.Err
	Data *automationPb.ApplyTemplateReply `json:"data"`
}

// GetTodo
type GetTodoResp struct {
	vo.Err
	Data *automationPb.GetTodoReply `json:"data"`
}

// UpdateTodo
type UpdateTodoReq struct {
	OrgId  int64                       `json:"orgId"`
	UserId int64                       `json:"userId"`
	Input  *automationPb.UpdateTodoReq `json:"input"`
}

type UpdateTodoResp struct {
	vo.Err
	Data *automationPb.UpdateTodoReply `json:"data"`
}
