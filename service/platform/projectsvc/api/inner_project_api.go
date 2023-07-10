package api

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

// GetProjectTemplateInner 获取项目模板内部接口
func (GetGreeter) GetProjectTemplateInner(req projectvo.GetProjectTemplateReq) projectvo.GetProjectTemplateResp {
	res, err := service.GetProjectTemplateInner(req)
	return projectvo.GetProjectTemplateResp{Err: vo.NewErr(err), Data: res}
}

// ApplyProjectTemplateInner 应用项目模板
func (PostGreeter) ApplyProjectTemplateInner(req projectvo.ApplyProjectTemplateReq) projectvo.ApplyProjectTemplateResp {
	data, err := service.ApplyProjectTemplateInner(&req)
	return projectvo.ApplyProjectTemplateResp{
		Err:  vo.NewErr(err),
		Data: data,
	}
}

//更新项目群聊成员
func (PostGreeter) ChangeProjectChatMember(req projectvo.ChangeProjectMemberReq) vo.CommonRespVo {
	err := service.ChangeProjectChatMember(req)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}

//批量删除模板项目
func (PostGreeter) DeleteProjectBatchInner(req projectvo.DeleteProjectInnerReq) vo.CommonRespVo {
	err := service.DeleteProjectBatchInner(req.OrgId, req.UserId, req.ProjectIds)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}

//是否允许创建项目
func (GetGreeter) AuthCreateProject(req projectvo.AuthCreateProjectReq) vo.CommonRespVo {
	err := domain.AuthPayProjectNum(req.OrgId, consts.FunctionProjectCreate)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}
