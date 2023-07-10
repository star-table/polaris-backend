package permissionfacade

import (
	"fmt"
	"strconv"

	"github.com/star-table/common/core/util/slice"

	"github.com/star-table/polaris-backend/facade"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"
)

// 获取应用操作权限列表（使用传入的协作人角色）
func GetAppAuthByCollaboratorRoles(orgId, appId, userId int64, collaboratorRoleIds []int64) *permissionvo.GetAppAuthResp {
	respVo := &permissionvo.GetAppAuthResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/appAuthByCollaboratorRoles", config.GetPreUrl(consts.ServiceNamePermission))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	queryParams["userId"] = userId
	queryParams["appId"] = appId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, collaboratorRoleIds, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 初始化应用权限
func InitAppPermission(req *permissionvo.InitAppPermissionReq) *permissionvo.InitAppPermissionResp {
	respVo := &permissionvo.InitAppPermissionResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/init", config.GetPreUrl(consts.ServiceNamePermission))
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 更新或创建管理组
func CreateOrUpdateManageGroup(req *permissionvo.CreateOrUpdateManageGroupReq) *permissionvo.CreateOrUpdateManageGroupResp {
	respVo := &permissionvo.CreateOrUpdateManageGroupResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/manage-group/save", config.GetPreUrl(consts.ServiceNamePermission))
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 更新权限组的权限项
func UpdateLcAppPermissionGroupOpConfig(req *permissionvo.UpdateLcAppPermissionGroupOpConfigReq) *permissionvo.UpdateLcAppPermissionGroupOpConfigResp {
	respVo := &permissionvo.UpdateLcAppPermissionGroupOpConfigResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/updateByLangCode", config.GetPreUrl(consts.ServiceNamePermission))
	err := facade.Request(consts.HttpMethodPut, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// UpdateLcPermissionGroup todo 更新管理组 http://192.168.88.164:10999/swagger-ui.html?urls.primaryName=permission#/%E5%BA%94%E7%94%A8%E6%9D%83%E9%99%90%E7%BB%84%EF%BC%88%E5%86%85%E9%83%A8%E8%B0%83%E7%94%A8%EF%BC%89/createOrUpdateManageGroupUsingPOST
// http://192.168.88.164:10999/swagger-ui.html?urls.primaryName=usercenter#/%E7%AE%A1%E7%90%86%E5%91%98/put_usercenter_api_v1_adminGroup_updateContents
func UpdateLcPermissionGroup(req *permissionvo.UpdateLcPermissionGroupReq) *permissionvo.UpdateLcPermissionGroupResp {
	respVo := &permissionvo.UpdateLcPermissionGroupResp{}
	// 先获取原有的管理树，获取系统管理组id，获取组ID对应的详情，将新管理员加入，保存
	adminGroup := GetLcPermissionGroupTree(req.OrgId, req.NewUserId)
	sysAdminGroupId := adminGroup.Data.SysGroup.Id
	detailResp := GetLcPermissionAdminGroupDetail(req.OrgId, req.NewUserId, sysAdminGroupId)
	if detailResp.Failure() {
		respVo.Err = vo.NewErr(detailResp.Error())
		return respVo
	}
	req.Values = append(detailResp.Data.AdminGroup.UserIds, strconv.FormatInt(req.NewUserId, 10))
	req.Values = slice.SliceUniqueString(req.Values)
	req.Key = "user_ids"
	req.Id = sysAdminGroupId

	reqUrl := fmt.Sprintf("%s/usercenter/api/v1/adminGroup/updateContents", config.GetPreUrl(consts.ServiceNamePermission))
	err := facade.Request(consts.HttpMethodPut, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetLcPermissionGroupTree 获取管理组信息（树状结构）
func GetLcPermissionGroupTree(orgId, userId int64) *permissionvo.GetLcPermissionGroupTreeResp {
	respVo := &permissionvo.GetLcPermissionGroupTreeResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/adminGroup/tree", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = userId
	queryParams["orgId"] = orgId
	req := &permissionvo.GetLcPermissionGroupTreeReq{
		UserId: userId,
		OrgId:  orgId,
	}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetDataAuthBatch(req permissionvo.GetDataAuthBatchReq) *permissionvo.GetDataAuthBatchResp {
	respVo := &permissionvo.GetDataAuthBatchResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/collaborators/dataAuths", config.GetPreUrl(consts.ServiceNamePermission))
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 获取应用权限组
func GetAppAuths(req permissionvo.GetAppAuthBatchReq) *permissionvo.GetAppAuthBatchResp {
	respVo := &permissionvo.GetAppAuthBatchResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/appAuths", config.GetPreUrl(consts.ServiceNamePermission))
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 获取权限最后更新时间
func GetPermissionUpdateTime(req permissionvo.GetPermissionUpdateTimeReq) *permissionvo.GetPermissionUpdateTimeResp {
	respVo := &permissionvo.GetPermissionUpdateTimeResp{}
	reqUrl := fmt.Sprintf("%s/permission/inner/api/v1/app-permission/get-permission-update-time", config.GetPreUrl(consts.ServiceNamePermission))
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}
