package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) Projects(reqVo projectvo.ProjectsRepVo) projectvo.ProjectsRespVo {
	res, err := service.Projects(reqVo)
	return projectvo.ProjectsRespVo{Err: vo.NewErr(err), ProjectList: res}
}

func (PostGreeter) CreateProject(reqVo projectvo.CreateProjectReqVo) projectvo.ProjectRespVo {
	res, err := service.CreateProject(reqVo)
	return projectvo.ProjectRespVo{Err: vo.NewErr(err), Project: res}
}

func (PostGreeter) UpdateProject(reqVo projectvo.UpdateProjectReqVo) projectvo.ProjectRespVo {
	res, err := service.UpdateProject(reqVo)
	return projectvo.ProjectRespVo{Err: vo.NewErr(err), Project: res}
}

func (PostGreeter) QuitProject(reqVo projectvo.ProjectIdReqVo) projectvo.QuitProjectRespVo {
	res, err := service.QuitProject(reqVo)
	return projectvo.QuitProjectRespVo{Err: vo.NewErr(err), QuitProject: res}
}

func (PostGreeter) StarProject(reqVo projectvo.ProjectIdReqVo) projectvo.OperateProjectRespVo {
	res, err := service.StarProject(reqVo)
	return projectvo.OperateProjectRespVo{Err: vo.NewErr(err), OperateProject: res}
}

func (PostGreeter) UnstarProject(reqVo projectvo.ProjectIdReqVo) projectvo.OperateProjectRespVo {
	res, err := service.UnstarProject(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	return projectvo.OperateProjectRespVo{Err: vo.NewErr(err), OperateProject: res}
}

//func (GetGreeter) ProjectStatistics(reqVo projectvo.ProjectIdReqVo) projectvo.ProjectStatisticsRespVo {
//	res, err := service.ProjectStatistics(reqVo.OrgId, reqVo.ProjectId)
//	return projectvo.ProjectStatisticsRespVo{Err: vo.NewErr(err), ProjectStatistics: res}
//}

func (PostGreeter) UpdateProjectStatus(reqVo projectvo.UpdateProjectStatusReqVo) vo.CommonRespVo {
	res, err := service.UpdateProjectStatus(reqVo)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) ProjectInfo(reqVo projectvo.ProjectInfoReqVo) projectvo.ProjectInfoRespVo {
	res, err := service.ProjectInfo(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.ProjectInfoRespVo{Err: vo.NewErr(err), ProjectInfo: res}
}

//通过项目类型langCode获取项目列表
func (GetGreeter) GetProjectBoListByProjectTypeLangCode(req projectvo.GetProjectBoListByProjectTypeLangCodeReqVo) projectvo.GetProjectBoListByProjectTypeLangCodeRespVo {
	res, err := service.GetProjectBoListByProjectTypeLangCode(req.OrgId, req.ProjectTypeLangCode)
	return projectvo.GetProjectBoListByProjectTypeLangCodeRespVo{ProjectBoList: res, Err: vo.NewErr(err)}
}

func (PostGreeter) GetSimpleProjectInfo(req projectvo.GetSimpleProjectInfoReqVo) projectvo.GetSimpleProjectInfoRespVo {
	res, err := service.GetSimpleProjectInfo(req.OrgId, req.Ids)
	return projectvo.GetSimpleProjectInfoRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) GetProjectDetails(req projectvo.GetSimpleProjectInfoReqVo) projectvo.GetProjectDetailsRespVo {
	res, err := service.GetProjectDetails(req.OrgId, req.Ids)
	return projectvo.GetProjectDetailsRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) GetProjectRelation(req projectvo.GetProjectRelationReqVo) projectvo.GetProjectRelationRespVo {
	res, err := service.GetProjectRelation(req.ProjectId, req.RelationType)
	return projectvo.GetProjectRelationRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) GetProjectRelationBatch(req projectvo.GetProjectRelationBatchReqVo) projectvo.GetProjectRelationBatchRespVo {
	res, err := service.GetProjectRelationBatch(req.OrgId, req.Data)
	return projectvo.GetProjectRelationBatchRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) ArchiveProject(reqVo projectvo.ProjectIdReqVo) vo.CommonRespVo {
	res, err := service.ArchiveProject(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) CancelArchivedProject(reqVo projectvo.ProjectIdReqVo) vo.CommonRespVo {
	res, err := service.CancelArchivedProject(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (GetGreeter) GetCacheProjectInfo(reqVo projectvo.GetCacheProjectInfoReqVo) projectvo.GetCacheProjectInfoRespVo {
	res, err := service.GetCacheProjectInfo(reqVo)
	return projectvo.GetCacheProjectInfoRespVo{Err: vo.NewErr(err), ProjectCacheBo: res}
}

//通过组织id集合获取未删除 未归档的项目
func (PostGreeter) GetProjectInfoByOrgIds(req projectvo.GetProjectInfoListByOrgIdsReqVo) projectvo.GetProjectInfoListByOrgIdsListRespVo {
	res, err := service.GetProjectInfoByOrgIds(req.OrgIds)
	return projectvo.GetProjectInfoListByOrgIdsListRespVo{ProjectInfoListByOrgIdsRespVo: res, Err: vo.NewErr(err)}
}

func (GetGreeter) OrgProjectMember(reqVo projectvo.OrgProjectMemberReqVo) projectvo.OrgProjectMemberListRespVo {
	res, err := service.OrgProjectMembers(reqVo)
	return projectvo.OrgProjectMemberListRespVo{Err: vo.NewErr(err), OrgProjectMemberRespVo: res}
}

func (PostGreeter) DeleteProject(reqVo projectvo.ProjectIdReqVo) vo.CommonRespVo {
	res, err := service.DeleteProject(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

//
//func (PostGreeter) GetProjectInfoForFeishuAction(req projectvo.GetProjectInfoForFeishuActionReqVo) projectvo.GetProjectInfoForFeishuActionRespVo {
//	res, err := service.GetProjectInfoForFeishuAction(req.OrgId, req.ProjectId)
//	return projectvo.GetProjectInfoForFeishuActionRespVo{
//		Err:  vo.NewErr(err),
//		NewData: res,
//	}
//}

//func (PostGreeter) GetLatestOperatedProjectInfo(req vo.CommonReqVo) projectvo.GetLatestOperatedProjectInfoResp {
//	res, err := service.GetLatestOperatedProjectInfo(req.OrgId, req.UserId)
//	return projectvo.GetLatestOperatedProjectInfoResp{
//		Err:  vo.NewErr(err),
//		Data: res,
//	}
//}

func (PostGreeter) GetSimpleProjectsByOrgId(req projectvo.GetSimpleProjectsByOrgIdReq) projectvo.GetSimpleProjectsByOrgIdResp {
	res, err := service.GetSimpleProjectsByOrgId(req.OrgId)
	return projectvo.GetSimpleProjectsByOrgIdResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (GetGreeter) GetOrgIssueAndProjectCount(req projectvo.GetOrgIssueAndProjectCountReq) projectvo.GetOrgIssueAndProjectCountResp {
	res, err := service.GetOrgIssueAndProjectCount(req.OrgId)
	return projectvo.GetOrgIssueAndProjectCountResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

// QueryProcessForAsyncTask 异步任务查询进度条
func (PostGreeter) QueryProcessForAsyncTask(req projectvo.QueryProcessForAsyncTaskReqVo) projectvo.QueryProcessForAsyncTaskRespVo {
	res, err := service.QueryProcessForAsyncTask(req.OrgId, &req.Input)
	return projectvo.QueryProcessForAsyncTaskRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
