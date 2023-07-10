package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) ImportIssues(reqVo projectvo.ImportIssuesReqVo) projectvo.ImportIssuesRespVo {
	data, err := service.ImportIssues(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.ImportIssuesRespVo{Err: vo.NewErr(err), Data: data}
}

func (GetGreeter) ExportIssueTemplate(reqVo projectvo.ExportIssueTemplateReqVo) projectvo.ExportIssueTemplateRespVo {
	url, err := service.ExportIssueTemplate(reqVo.OrgId, reqVo.ProjectId, reqVo.ProjectObjectTypeId)
	return projectvo.ExportIssueTemplateRespVo{Data: &vo.ExportIssueTemplateResp{URL: url}, Err: vo.NewErr(err)}
}

func (PostGreeter) ExportData(reqVo projectvo.ExportIssueReqVo) projectvo.ExportIssueTemplateRespVo {
	url, err := service.ExportData(reqVo.OrgId, reqVo)
	return projectvo.ExportIssueTemplateRespVo{Data: &vo.ExportIssueTemplateResp{URL: url}, Err: vo.NewErr(err)}
}

// ExportUserOrDeptSameNameList 导出同名部门和用户列表
func (GetGreeter) ExportUserOrDeptSameNameList(reqVo projectvo.ExportUserOrDeptSameNameListReqVo) projectvo.ExportIssueTemplateRespVo {
	url, err := service.ExportUserOrDeptSameNameList(reqVo.OrgId, reqVo.ProjectId)
	return projectvo.ExportIssueTemplateRespVo{Data: &vo.ExportIssueTemplateResp{URL: url}, Err: vo.NewErr(err)}
}

// GetUserOrDeptSameNameList 获取同名部门或用户列表
func (GetGreeter) GetUserOrDeptSameNameList(reqVo projectvo.GetUserOrDeptSameNameListReq) projectvo.GetUserOrDeptSameNameListRespVo {
	res, err := service.GetUserOrDeptSameNameList(reqVo.OrgId, reqVo.DataType)
	return projectvo.GetUserOrDeptSameNameListRespVo{Data: res, Err: vo.NewErr(err)}
}

// ImportIssues2022 2022 实现的任务导入
// 无需解析 excel，直接根据入参进行任务新增/更新
//func (PostGreeter) ImportIssues2022(reqVo projectvo.ImportIssues2022ReqVo) projectvo.ImportIssuesRespVo {
//	data, err := service.ImportIssues2022(reqVo.OrgId, reqVo.UserId, &reqVo.Input)
//	return projectvo.ImportIssuesRespVo{Err: vo.NewErr(err), Data: data}
//}
