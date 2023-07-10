package api

import (
	"context"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
)

// 获取任务的工时信息
func (r *queryResolver) GetIssueWorkHoursInfo(ctx context.Context, params vo.GetIssueWorkHoursInfoReq) (*vo.GetIssueWorkHoursInfoResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.GetIssueWorkHoursInfo(projectvo.GetIssueWorkHoursInfoReqVo{
		Input:  params,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

func (r *mutationResolver) CreateIssueWorkHours(ctx context.Context, params *vo.CreateIssueWorkHoursReq) (*vo.BoolResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.CreateIssueWorkHours(projectvo.CreateIssueWorkHoursReqVo{
		Input:  *params,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

// 创建详细的预估工时记录
func (r *mutationResolver) CreateMultiIssueWorkHours(ctx context.Context, params *vo.CreateMultiIssueWorkHoursReq) (*vo.BoolResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.CreateMultiIssueWorkHours(projectvo.CreateMultiIssueWorkHoursReqVo{
		Input:  *params,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

// 更新一条工时记录
func (r *mutationResolver) UpdateIssueWorkHours(ctx context.Context, params *vo.UpdateIssueWorkHoursReq) (*vo.BoolResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.UpdateIssueWorkHours(projectvo.UpdateIssueWorkHoursReqVo{
		Input:  *params,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

// 编辑详细预估工时信息
func (r *mutationResolver) UpdateMultiIssueWorkHours(ctx context.Context, params *vo.UpdateMultiIssueWorkHoursReq) (*vo.BoolResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.UpdateMultiIssueWorkHours(projectvo.UpdateMultiIssueWorkHoursReqVo{
		Input:  *params,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

// 删除一条工时记录
func (r *mutationResolver) DeleteIssueWorkHours(ctx context.Context, params *vo.DeleteIssueWorkHoursReq) (*vo.BoolResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.DeleteIssueWorkHours(projectvo.DeleteIssueWorkHoursReqVo{
		Input:  *params,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

// 获取工时统计信息
func (r *queryResolver) GetWorkHourStatistic(ctx context.Context, params vo.GetWorkHourStatisticReq) (*vo.GetWorkHourStatisticResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.GetWorkHourStatistic(projectvo.GetWorkHourStatisticReqVo{
		Input:  params,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

// 检查员工是否是任务成员
func (r *queryResolver) CheckIsIssueMember(ctx context.Context, params vo.CheckIsIssueMemberReq) (*vo.BoolResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.CheckIsIssueMember(projectvo.CheckIsIssueMemberReqVo{
		Input:  params,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

// 将一个员工加入到任务成员中
func (r *mutationResolver) SetUserJoinIssue(ctx context.Context, params vo.SetUserJoinIssueReq) (*vo.BoolResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.SetUserJoinIssue(projectvo.SetUserJoinIssueReqVo{
		Input:  params,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

func (r *queryResolver) CheckIsEnableWorkHour(ctx context.Context, params vo.CheckIsEnableWorkHourReq) (*vo.CheckIsEnableWorkHourResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.CheckIsEnableWorkHour(projectvo.CheckIsEnableWorkHourReqVo{
		Input:  params,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

func (r *queryResolver) ExportWorkHourStatistic(ctx context.Context, params vo.GetWorkHourStatisticReq) (*vo.ExportWorkHourStatisticResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.ExportWorkHourStatistic(projectvo.GetWorkHourStatisticReqVo{
		Input:  params,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}
