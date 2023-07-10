package permissionfacade

import (
	"fmt"

	"github.com/star-table/polaris-backend/facade"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"
)

// 获取用户管理组权限信息
func GetUserManageAuthInfo(orgId, userId int64) *permissionvo.ManageAuthInfoResp {
	respVo := &permissionvo.ManageAuthInfoResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/manage-group/user-manage-auth-info", config.GetPreUrl(consts.ServiceNamePermission))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	queryParams["userId"] = userId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 获取应用操作权限列表
func GetOptAuthList(orgId, appPackageId, appId, userId int64) *permissionvo.GetOptAuthListResp {
	respVo := &permissionvo.GetOptAuthListResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/optAuthMap", config.GetPreUrl(consts.ServiceNamePermission))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	queryParams["userId"] = userId
	queryParams["appPackageId"] = appPackageId
	queryParams["appId"] = appId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 获取应用操作权限列表
func GetAppAuth(orgId, appId, tableId, userId int64) *permissionvo.GetAppAuthResp {
	respVo := &permissionvo.GetAppAuthResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/appAuth", config.GetPreUrl(consts.ServiceNamePermission))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	queryParams["userId"] = userId
	queryParams["appId"] = appId
	queryParams["tableId"] = tableId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 获取应用操作权限列表-不合并协作人角色（用于判断单挑数据的基本权限）
func GetAppAuthWithoutCollaborator(orgId, appId, userId int64) *permissionvo.GetAppAuthResp {
	respVo := &permissionvo.GetAppAuthResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/appAuthWithoutCollaborator", config.GetPreUrl(consts.ServiceNamePermission))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	queryParams["userId"] = userId
	queryParams["appId"] = appId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetAppRoleList 获取应用的角色列表
func GetAppRoleList(orgId, appId int64) *permissionvo.GetAppRoleListResp {
	respVo := &permissionvo.GetAppRoleListResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/getAppPermissionGroupList", config.GetPreUrl(consts.ServiceNamePermission))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	queryParams["appId"] = appId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 获取应用操作权限列表
func GetDataAuth(orgId, appId, userId int64, dataId int64) *permissionvo.GetAppAuthResp {
	respVo := &permissionvo.GetAppAuthResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/collaborators/dataAuth", config.GetPreUrl(consts.ServiceNamePermission))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	queryParams["userId"] = userId
	queryParams["appId"] = appId
	queryParams["dataId"] = dataId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 获取成员与权限组的映射关系
func GetUserGroupMappings(orgId, userId, appId int64) *permissionvo.GetUserGroupMappingsResp {
	respVo := &permissionvo.GetUserGroupMappingsResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/member-group-mappings", config.GetPreUrl(consts.ServiceNamePermission))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	queryParams["userId"] = userId
	queryParams["appId"] = appId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetLcPermissionAdminGroupDetail 获取管理组详情
// userId 表示操作人，此处待优化。
func GetLcPermissionAdminGroupDetail(orgId, userId int64, adminGroupId string) *permissionvo.GetLcPermissionAdminGroupDetailResp {
	respVo := &permissionvo.GetLcPermissionAdminGroupDetailResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/adminGroup/detail", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	queryParams := map[string]interface{}{}
	queryParams["id"] = adminGroupId
	queryParams["orgId"] = orgId
	queryParams["userId"] = userId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}
