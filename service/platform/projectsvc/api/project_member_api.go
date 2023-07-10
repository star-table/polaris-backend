package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) RemoveProjectMember(req projectvo.RemoveProjectMemberReqVo) vo.CommonRespVo {
	res, err := service.RemoveProjectMember(req.OrgId, req.UserId, req.SourceChannel, req.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) ProjectUserList(req projectvo.ProjectUserListReq) projectvo.ProjectUserListRespVo {
	res, err := service.ProjectUserList(req.OrgId, req.Page, req.Size, req.Input)
	return projectvo.ProjectUserListRespVo{Err: vo.NewErr(err), Data: res}
}

// 融合版：获取项目的成员列表
func (PostGreeter) ProjectUserListForFuse(req projectvo.ProjectUserListReq) projectvo.ProjectUserListForFuseRespVo {
	res, err := service.ProjectUserListForFuse(req.OrgId, req.Page, req.Size, req.Input)
	return projectvo.ProjectUserListForFuseRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) AddProjectMember(req projectvo.RemoveProjectMemberReqVo) vo.CommonRespVo {
	res, err := service.AddProjectMember(req.OrgId, req.UserId, req.SourceChannel, req.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) OrgProjectMemberList(req projectvo.OrgProjectMemberListReq) projectvo.OrgProjectMemberListResp {
	res, err := service.OrgProjectMemberList(req.OrgId, req.SourceChannel, req.Page, req.Size, req.Params)
	return projectvo.OrgProjectMemberListResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetProjectRelationUserIds(req projectvo.GetProjectRelationUserIdsReq) projectvo.GetProjectRelationUserIdsResp {
	res, err := service.GetProjectRelationUserIds(req.ProjectId, req.RelationType)
	return projectvo.GetProjectRelationUserIdsResp{
		Err:     vo.NewErr(err),
		UserIds: res,
	}
}

func (PostGreeter) ProjectMemberIdList(req projectvo.ProjectMemberIdListReq) projectvo.ProjectMemberIdListResp {
	res, err := service.ProjectMemberIdList(req.OrgId, req.ProjectId, req.Data)
	return projectvo.ProjectMemberIdListResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetProjectMemberIds(req projectvo.GetProjectMemberIdsReqVo) projectvo.GetProjectMemberIdsResp {
	res, err := service.GetProjectMemberIds(req.OrgId, req.Input.ProjectId, req.Input.IncludeAdmin)
	return projectvo.GetProjectMemberIdsResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetTrendsMembers(req projectvo.GetTrendListMembersReqVo) projectvo.GetTrendListMembersResp {
	res, err := service.GetTrendsMembers(req.OrgId, req.UserId, req.Input.ProjectId)
	return projectvo.GetTrendListMembersResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
