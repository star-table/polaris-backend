package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service/work_hour"
)

// 获取列表
func (PostGreeter) GetIssueWorkHoursList(reqVo projectvo.GetIssueWorkHoursListReqVo) projectvo.GetIssueWorkHoursListRespVo {
	resp, err := work_hour.GetIssueWorkHoursList(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.GetIssueWorkHoursListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 创建工时记录
func (PostGreeter) CreateIssueWorkHours(reqVo projectvo.CreateIssueWorkHoursReqVo) projectvo.BoolRespVo {
	resp, err := work_hour.CreateIssueWorkHours(reqVo.OrgId, reqVo.UserId, reqVo.Input, true)
	return projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 创建详细工时记录，含多条预估记录
func (PostGreeter) CreateMultiIssueWorkHours(reqVo projectvo.CreateMultiIssueWorkHoursReqVo) projectvo.BoolRespVo {
	resp, err := work_hour.CreateMultiIssueWorkHours(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 编辑工时
func (PostGreeter) UpdateIssueWorkHours(reqVo projectvo.UpdateIssueWorkHoursReqVo) projectvo.BoolRespVo {
	resp, err := work_hour.UpdateIssueWorkHours(reqVo.OrgId, reqVo.UserId, reqVo.Input, true, true, true, true)
	return projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 编辑详细预估工时记录
func (PostGreeter) UpdateMultiIssueWorkHours(reqVo projectvo.UpdateMultiIssueWorkHoursReqVo) projectvo.BoolRespVo {
	resp, err := work_hour.UpdateMultiIssueWorkHourWithDelete(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 删除工时列表
func (PostGreeter) DeleteIssueWorkHours(reqVo projectvo.DeleteIssueWorkHoursReqVo) projectvo.BoolRespVo {
	resp, err := work_hour.DeleteIssueWorkHours(reqVo.OrgId, reqVo.UserId, reqVo.Input, true, true)
	return projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 启用、关闭一个项目的工时功能
func (PostGreeter) DisOrEnableIssueWorkHours(reqVo projectvo.DisOrEnableIssueWorkHoursReqVo) projectvo.BoolRespVo {
	resp, err := work_hour.DisOrEnableIssueWorkHours(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 获取一个任务的工时信息
func (PostGreeter) GetIssueWorkHoursInfo(reqVo projectvo.GetIssueWorkHoursInfoReqVo) projectvo.GetIssueWorkHoursInfoRespVo {
	resp, err := work_hour.GetIssueWorkHoursInfo(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.GetIssueWorkHoursInfoRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 获取工时统计数据
func (PostGreeter) GetWorkHourStatistic(reqVo projectvo.GetWorkHourStatisticReqVo) projectvo.GetWorkHourStatisticRespVo {
	resp, err := work_hour.GetWorkHourStatistic(reqVo.OrgId, reqVo.Input)
	return projectvo.GetWorkHourStatisticRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 检查用户是否是任务的关注人
func (PostGreeter) CheckIsIssueMember(reqVo projectvo.CheckIsIssueMemberReqVo) projectvo.BoolRespVo {
	resp, err := work_hour.CheckIsIssueMember(reqVo.OrgId, reqVo.Input)
	return projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 将用户加入到任务成员中
func (PostGreeter) SetUserJoinIssue(reqVo projectvo.SetUserJoinIssueReqVo) projectvo.BoolRespVo {
	resp, err := work_hour.SetUserJoinIssue(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 检查项目的工时是否开启
func (PostGreeter) CheckIsEnableWorkHour(reqVo projectvo.CheckIsEnableWorkHourReqVo) projectvo.CheckIsEnableWorkHourRespVo {
	resp, err := work_hour.CheckIsEnableWorkHour(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.CheckIsEnableWorkHourRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 工时统计数据的导出
func (PostGreeter) ExportWorkHourStatistic(reqVo projectvo.GetWorkHourStatisticReqVo) projectvo.ExportWorkHourStatisticRespVo {
	resp, err := work_hour.ExportWorkHourStatistic(reqVo.OrgId, reqVo.Input)
	return projectvo.ExportWorkHourStatisticRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// 检查是否是任务相关人员（创建人，关注人，负责人）
func (PostGreeter) CheckIsIssueRelatedPeople(reqVo projectvo.CheckIsIssueMemberReqVo) projectvo.BoolRespVo {
	resp, err := work_hour.CheckIsIssueRelatedPeople(reqVo.OrgId, reqVo.Input)
	return projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}
