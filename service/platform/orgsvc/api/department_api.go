package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) Departments(req orgvo.DepartmentsReqVo) orgvo.DepartmentsRespVo {
	page := req.Page
	size := req.Size
	params := req.Params
	orgId := req.OrgId

	pageA := uint(0)
	sizeA := uint(0)
	if page != nil && size != nil && *page > 0 && *size > 0 {
		pageA = uint(*page)
		sizeA = uint(*size)
	}
	res, err := service.Departments(pageA, sizeA, params, orgId)
	return orgvo.DepartmentsRespVo{Err: vo.NewErr(err), DepartmentList: res}
}

func (PostGreeter) DepartmentMembers(req orgvo.DepartmentMembersReqVo) orgvo.DepartmentMembersRespVo {
	res, err := service.DepartmentMembers(req.Params, req.OrgId)
	return orgvo.DepartmentMembersRespVo{Err: vo.NewErr(err), DepartmentMemberInfos: res}
}

func (PostGreeter) CreateDepartment(req orgvo.CreateDepartmentReqVo) vo.CommonRespVo {
	res, err := service.CreateDepartment(req.Data, req.OrgId, req.CurrentUserId)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateDepartment(req orgvo.UpdateDepartmentReqVo) vo.CommonRespVo {
	res, err := service.UpdateDepartment(req.Data, req.OrgId, req.CurrentUserId)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) DeleteDepartment(req orgvo.DeleteDepartmentReqVo) vo.CommonRespVo {
	res, err := service.DeleteDepartment(req.Data, req.OrgId, req.CurrentUserId)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) AllocateDepartment(req orgvo.AllocateDepartmentReqVo) vo.CommonRespVo {
	err := service.AllocateDepartment(req.Data, req.OrgId, req.CurrentUserId)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}

func (PostGreeter) SetUserDepartmentLevel(req orgvo.SetUserDepartmentLevelReqVo) vo.CommonRespVo {
	err := service.SetUserDepartmentLevel(req.OrgId, req.CurrentUserId, req.Data)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}

func (PostGreeter) DepartmentMembersList(req orgvo.DepartmentMembersListReq) orgvo.DepartmentMembersListResp {
	page := req.Page
	size := req.Size

	pageA := 0
	sizeA := 0
	if page != nil && size != nil && *page > 0 && *size > 0 {
		pageA = *page
		sizeA = *size
	}
	var name *string
	var userIds []int64
	var excludeProjectId int64 = 0
	var relationType int64 = 0
	if req.Params != nil {
		if req.Params.Name != nil {
			name = req.Params.Name
		}
		userIds = req.Params.UserIds
		if req.Params.ExcludeProjectID != nil {
			excludeProjectId = *req.Params.ExcludeProjectID
		}
		if req.Params.RelationType != nil {
			relationType = *req.Params.RelationType
		}
	}
	res, err := service.DepartmentMembersList(req.OrgId, req.SourceChannel, name, pageA, sizeA, userIds, req.IgnoreDelete, excludeProjectId, relationType)

	return orgvo.DepartmentMembersListResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) VerifyDepartments(req orgvo.VerifyDepartmentsReq) vo.BoolRespVo {
	return vo.BoolRespVo{
		Err:    vo.NewErr(nil),
		IsTrue: service.VerifyDepartments(req.OrgId, req.DepartmentIds),
	}
}

func (PostGreeter) GetDeptByIds(req orgvo.GetDeptByIdsReq) orgvo.GetDeptByIdsResp {
	res, err := service.GetDeptByIds(req.OrgId, req.DeptIds)
	return orgvo.GetDeptByIdsResp{
		Err:  vo.NewErr(err),
		List: res,
	}
}
