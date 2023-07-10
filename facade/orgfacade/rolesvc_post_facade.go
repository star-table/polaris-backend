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

func Authenticate(req rolevo.AuthenticateReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/authenticate", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["path"] = req.Path
	queryParams["operation"] = req.Operation
	requestBody := &req.AuthInfoReqVo
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ClearUserRoleList(req rolevo.ClearUserRoleReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/clearUserRoleList", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	requestBody := &req.UserIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateRole(req rolevo.CreateOrgReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/createRole", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DelRole(req rolevo.DelRoleReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/delRole", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgAdminUserBatch(req rolevo.GetOrgAdminUserBatchReqVo) rolevo.GetOrgAdminUserBatchRespVo {
	respVo := &rolevo.GetOrgAdminUserBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getOrgAdminUserBatch", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func RemoveRoleDepartmentRelation(req rolevo.RemoveRoleDepartmentRelationReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/removeRoleDepartmentRelation", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["operatorId"] = req.OperatorId
	requestBody := &req.DeptIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func RemoveRoleUserRelation(req rolevo.RemoveRoleUserRelationReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/removeRoleUserRelation", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["operatorId"] = req.OperatorId
	requestBody := &req.UserIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func RoleInit(req rolevo.RoleInitReqVo) rolevo.RoleInitRespVo {
	respVo := &rolevo.RoleInitRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/roleInit", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func RoleUserRelation(req rolevo.RoleUserRelationReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/roleUserRelation", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["roleId"] = req.RoleId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateOrgOwner(req rolevo.UpdateOrgOwnerReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/updateOrgOwner", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["oldOwnerId"] = req.OldOwnerId
	queryParams["newOwnerId"] = req.NewOwnerId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateRole(req rolevo.UpdateRoleReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/updateRole", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateRolePermissionOperation(req rolevo.UpdateRolePermissionOperationReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/updateRolePermissionOperation", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateUserOrgRoleBatch(req rolevo.UpdateUserOrgRoleBatchReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/rolesvc/updateUserOrgRoleBatch", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["currentUserId"] = req.CurrentUserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

