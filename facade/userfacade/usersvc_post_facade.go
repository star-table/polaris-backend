package userfacade

import (
	"fmt"

	"github.com/star-table/polaris-backend/facade"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
)

// GetUserAuthority 获取成员权限信息，调用的是 usercenter 的 inner `getUserAuthoritySimple` 接口。
// Notice: 该方法调用 usercenter 的实现，由于要优化耗时，去除不必要的信息查询，在 usercenter 中做了裁剪，无需用到的字段则无需查询。
// 调整后，返回的 `*uservo.GetUserAuthorityResp` 中可用信息有：`NewData.OptAuth`, `NewData.IsSysAdmin`, `NewData.IsSubAdmin`。
// 如果需要用到该接口返回的其他字段，需检查 usercenter 实现中，是否进行查询和处理。
func GetUserAuthority(orgId, userId int64) *uservo.GetUserAuthorityResp {
	respVo := &uservo.GetUserAuthorityResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/user/getUserAuthoritySimple", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = userId
	queryParams["orgId"] = orgId
	req := &uservo.GetUserAuthorityReq{
		UserId: userId,
		OrgId:  orgId,
	}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// InitDefaultManageGroup 初始化权限组
func InitDefaultManageGroup(orgId, userId int64) *uservo.InitDefaultManageGroupResp {
	respVo := &uservo.InitDefaultManageGroupResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/manage-group/init", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = userId
	queryParams["orgId"] = orgId
	queryParams["sourceFrom"] = "polaris"
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 给系统管理组里面添加人员（没有的时候会新建）
func AddUserToSysManageGroup(orgId, userId int64, req uservo.AddUserToSysManageGroupReq) *uservo.AddUserToSysManageGroupResp {
	respVo := &uservo.AddUserToSysManageGroupResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/manage-group/add-user-to-sys-manage-group", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = userId
	queryParams["orgId"] = orgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// AddScriptTask 向 schedule 服务发送新增脚本任务的请求
func AddScriptTask(orgId, userId int64, req uservo.AddScriptTaskReq) *uservo.AddScriptTaskResp {
	respVo := &uservo.AddScriptTaskResp{}
	reqUrl := fmt.Sprintf("%s/api/schedule/addTask", config.GetPreUrl("schedule"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = userId
	queryParams["orgId"] = orgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// 获取该组织具有管理员权限（超管、普通管理员）的人员
func GetUsersCouldManage(orgId int64, appId int64) *uservo.GetUsersCouldManageResp {
	respVo := &uservo.GetUsersCouldManageResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/user/get-users-could-manage", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	queryParams["appId"] = appId
	req := map[string]interface{}{
		"orgId": orgId,
		"appId": appId,
	}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// CheckAndSetSuperAdmin 检查组织拥有者是否是超管，如果不是，则设置为超管。该方法一般被系统调用，如事件回调。
func CheckAndSetSuperAdmin(orgId int64) *uservo.CheckAndSetSuperAdminResp {
	respVo := &uservo.CheckAndSetSuperAdminResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/org/check-and-set-super-admin", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	req := map[string]interface{}{
		"orgId":  orgId,
		"userId": -1,
	}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// ClearAdminGroupForOneUser 将一个用户从该组织的所有管理组中移除
func ClearAdminGroupForOneUser(orgId, operateUid, targetUid int64) *uservo.CheckAndSetSuperAdminResp {
	respVo := &uservo.CheckAndSetSuperAdminResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/adminGroup/deleteOneUserFromOrg", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	queryParams["operateUid"] = operateUid
	queryParams["targetUserId"] = targetUid
	req := map[string]interface{}{
		"orgId": orgId,
	}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetAllUserByIds 根据id，批量获取用户信息
func GetAllUserByIds(orgId int64, userIds []int64) *uservo.GetAllUserByIdsResp {
	respVo := &uservo.GetAllUserByIdsResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/user/getAllListByIds", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	req := map[string]interface{}{
		"orgId": orgId,
		"ids":   userIds,
	}
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetAllDeptByIds 根据id，批量获取部门信息
func GetAllDeptByIds(orgId int64, deptIds []int64) *uservo.GetAllDeptByIdsResp {
	respVo := &uservo.GetAllDeptByIdsResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/dept/getAllListByIds", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	req := map[string]interface{}{
		"orgId": orgId,
		"ids":   deptIds,
	}
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetFullDeptNameArrByIds 根据id，获取全路径部门名称信息，如：["职能部", "技术部", "运维组"]
func GetFullDeptNameArrByIds(orgId int64, userIds []int64) *uservo.GetFullDeptNameArrByIdsResp {
	respVo := &uservo.GetFullDeptNameArrByIdsResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/dept/getFullDeptByIds", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	req := map[string]interface{}{
		"orgId":   orgId,
		"deptIds": userIds,
	}
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetDeptList 获取组织的部门列表
func GetDeptList(orgId int64) *uservo.GetDeptListResp {
	respVo := &uservo.GetDeptListResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/dept/getDeptList", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	req := map[string]interface{}{
		"orgId": orgId,
	}
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// AddNewMenuToRole 向组织角色中追加新的菜单权限项
func AddNewMenuToRole(req *uservo.AddNewMenuToRoleReqInput) *uservo.BoolResp {
	respVo := &uservo.BoolResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/adminGroup/addNewMenuToRole", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetNormalAdminManageApps(orgId int64) *uservo.GetCommAdminMangeAppsResp {
	respVo := &uservo.GetCommAdminMangeAppsResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/manage-group/getCommonAdminManageApps", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = orgId
	req := map[string]interface{}{
		"orgId": orgId,
	}
	respVo.Err = facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	return respVo
}

func GetUserIdsByDeptIds(req *uservo.GetUserIdsByDeptIdsReq) *uservo.GetUserIdsByDeptIdsResp {
	respVo := &uservo.GetUserIdsByDeptIdsResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/dept/getUserIdsByDeptIds", config.GetPreUrl(consts.ServiceNameLcUserCenter))
	respVo.Err = facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	return respVo
}
