package orgfacade

import (
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	"github.com/star-table/polaris-backend/facade"
)

var _ = commonvo.DataEvent{}

func ChangeDefaultRole() vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/changeDefaultRole", config.GetPreUrl("rolesvc"))
	err := facade.Request(consts.HttpMethodGet, reqUrl, nil, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgAdminUser(req rolevo.GetOrgAdminUserReqVo) rolevo.GetOrgAdminUserRespVo {
	respVo := &rolevo.GetOrgAdminUserRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getOrgAdminUser", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgRoleList(req rolevo.GetOrgRoleListReqVo) rolevo.GetOrgRoleListRespVo {
	respVo := &rolevo.GetOrgRoleListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getOrgRoleList", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgRoleUser(req rolevo.GetOrgRoleUserReqVo) rolevo.GetOrgRoleUserRespVo {
	respVo := &rolevo.GetOrgRoleUserRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getOrgRoleUser", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetPersonalPermissionInfo(req rolevo.GetPersonalPermissionInfoReqVo) rolevo.GetPersonalPermissionInfoRespVo {
	respVo := &rolevo.GetPersonalPermissionInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getPersonalPermissionInfo", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["projectId"] = req.ProjectId
	queryParams["issueId"] = req.IssueId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetProjectRoleList(req rolevo.GetProjectRoleListReqVo) rolevo.GetOrgRoleListRespVo {
	respVo := &rolevo.GetOrgRoleListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getProjectRoleList", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserAdminFlag(req rolevo.GetUserAdminFlagReqVo) rolevo.GetUserAdminFlagRespVo {
	respVo := &rolevo.GetUserAdminFlagRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getUserAdminFlag", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PermissionOperationList(req rolevo.PermissionOperationListReqVo) rolevo.PermissionOperationListRespVo {
	respVo := &rolevo.PermissionOperationListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/permissionOperationList", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["roleId"] = req.RoleId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateUserOrgRole(req rolevo.UpdateUserOrgRoleReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/updateUserOrgRole", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["userId"] = req.UserId
	queryParams["roleId"] = req.RoleId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

