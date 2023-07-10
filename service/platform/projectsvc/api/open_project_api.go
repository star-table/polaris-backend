package api

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) OpenCreateProject(reqVo projectvo.CreateProjectReqVo) projectvo.ProjectRespVo {
	// 校验操作人是否存在
	_, err := orgfacade.GetBaseUserInfoRelaxed(reqVo.OrgId, reqVo.UserId)
	if err != nil {
		log.Error(err)
		return projectvo.ProjectRespVo{Err: vo.NewErr(errs.OperatorInvalid), Project: nil}
	}

	// // common/core/consts/business_consts.go:83
	// consts.ProjectTypeCommon2022V47
	// 空白普通项目
	projectTypeID := int64(consts.ProjectTypeCommon2022V47)
	// 私有项目,公开项目暂时不支持
	projectStatus := int(consts.PrivateProject)

	//
	reqVo.Input.ProjectTypeID = &projectTypeID
	reqVo.Input.PublicStatus = projectStatus

	res, err := service.CreateProjectWithoutAuth(reqVo)
	return projectvo.ProjectRespVo{Err: vo.NewErr(err), Project: res}
}

func (PostGreeter) OpenProjects(reqVo projectvo.ProjectsRepVo) projectvo.ProjectsRespVo {
	res, err := service.Projects(reqVo)
	return projectvo.ProjectsRespVo{
		Err:         vo.NewErr(err),
		ProjectList: res,
	}
}

func (PostGreeter) OpenProjectInfo(reqVo projectvo.ProjectInfoReqVo) projectvo.ProjectInfoRespVo {
	res, err := service.ProjectInfo(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.ProjectInfoRespVo{
		Err:         vo.NewErr(err),
		ProjectInfo: res,
	}
}

func (PostGreeter) OpenUpdateProject(reqVo projectvo.UpdateProjectReqVo) projectvo.ProjectRespVo {
	// 校验操作人是否存在
	_, err := orgfacade.GetBaseUserInfoRelaxed(reqVo.OrgId, reqVo.UserId)
	if err != nil {
		log.Error(err)
		return projectvo.ProjectRespVo{Err: vo.NewErr(errs.OperatorInvalid), Project: nil}
	}

	//兼容updateFields
	input := reqVo.Input
	updateFields := make([]string, 0)
	if input.Name != nil {
		updateFields = append(updateFields, "name")
	}
	if input.Owner != nil {
		updateFields = append(updateFields, "owner")
	}
	if input.PriorityID != nil {
		updateFields = append(updateFields, "priorityId")
	}
	if input.PlanStartTime != nil {
		updateFields = append(updateFields, "planStartTime")
	}
	if input.PlanEndTime != nil {
		updateFields = append(updateFields, "planEndTime")
	}
	if input.PublicStatus != nil {
		updateFields = append(updateFields, "publicStatus")
	}
	if input.ResourceID != nil {
		updateFields = append(updateFields, "resourceId")
	}
	if input.ResourcePath != nil {
		updateFields = append(updateFields, "resourcePath")
	}
	if input.Remark != nil {
		updateFields = append(updateFields, "remark")
	}
	if input.Status != nil {
		updateFields = append(updateFields, "status")
	}
	if input.MemberIds != nil {
		updateFields = append(updateFields, "memberIds")
	}
	if input.MemberForDepartmentID != nil {
		updateFields = append(updateFields, "memberForDepartmentId")
	}
	if input.IsAllMember != nil {
		//这里也是用memberIds
		updateFields = append(updateFields, "memberIds")
	}
	if input.FollowerIds != nil {
		updateFields = append(updateFields, "followerIds")
	}
	if input.IsSyncOutCalendar != nil {
		updateFields = append(updateFields, "isSyncOutCalendar")
	}
	if input.SyncCalendarStatusList != nil {
		updateFields = append(updateFields, "syncCalendarStatusList")
	}
	if input.IsCreateFsChat != nil {
		updateFields = append(updateFields, "isCreateFsChat")
	}
	reqVo.Input.UpdateFields = updateFields

	judgeErr := service.JudgeProjectFiling(reqVo.OrgId, reqVo.Input.ID)
	if judgeErr != nil {
		log.Error(judgeErr)
		return projectvo.ProjectRespVo{Err: vo.NewErr(judgeErr), Project: nil}
	}

	res, err := service.UpdateProjectWithoutAuth(reqVo)
	return projectvo.ProjectRespVo{Err: vo.NewErr(err), Project: res}
}

func (PostGreeter) OpenDeleteProject(reqVo projectvo.ProjectIdReqVo) vo.CommonRespVo {
	// 校验操作人是否存在
	_, err := orgfacade.GetBaseUserInfoRelaxed(reqVo.OrgId, reqVo.UserId)
	if err != nil {
		log.Error(err)
		return vo.CommonRespVo{Err: vo.NewErr(errs.OperatorInvalid), Void: nil}
	}

	judgeErr := service.JudgeProjectFiling(reqVo.OrgId, reqVo.ProjectId)
	if judgeErr != nil {
		log.Error(judgeErr)
		return vo.CommonRespVo{
			Err:  vo.NewErr(judgeErr),
			Void: nil,
		}
	}

	res, err := service.DeleteProjectWithoutAuth(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
}

func (GetGreeter) OpenGetPriorityList(reqVo projectvo.OpenPriorityListReqVo) projectvo.OpenSomeAttrListRespVo {
	resp, err := service.OpenPriorityList(reqVo.OrgId)
	return projectvo.OpenSomeAttrListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (GetGreeter) OpenGetProjectObjectTypeList(reqVo projectvo.OpenPriorityListReqVo) projectvo.OpenSomeAttrListRespVo {
	resp, err := service.OpenGetProjectObjectTypeList(reqVo)
	return projectvo.OpenSomeAttrListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (GetGreeter) OpenGetIterationList(reqVo projectvo.OpenGetIterationListReqVo) projectvo.OpenGetIterationListRespVo {
	resp, err := service.OpenGetIterationList(reqVo)
	return projectvo.OpenGetIterationListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (GetGreeter) OpenGetIssueSourceList(reqVo projectvo.OpenGetDemandSourceListReqVo) projectvo.OpenSomeAttrListRespVo {
	resp, err := service.OpenGetDemandSourceList(reqVo)
	return projectvo.OpenSomeAttrListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (GetGreeter) OpenGetPropertyList(reqVo projectvo.OpenGetPropertyListReqVo) projectvo.OpenSomeAttrListRespVo {
	resp, err := service.OpenGetPropertyList(reqVo)
	return projectvo.OpenSomeAttrListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}
