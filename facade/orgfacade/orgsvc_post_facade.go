package orgfacade

import (
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade"
)

var _ = commonvo.DataEvent{}

func AddUser(req orgvo.AddUserReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/addUser", config.GetPreUrl("orgsvc"))
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

func AllocateDepartment(req orgvo.AllocateDepartmentReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/allocateDepartment", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ApplyScopes(req orgvo.ApplyScopesReqVo) orgvo.ApplyScopesRespVo {
	respVo := &orgvo.ApplyScopesRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/applyScopes", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AuthDingCode(req orgvo.AuthDingCodeReqVo) orgvo.AuthDingCodeRespVo {
	respVo := &orgvo.AuthDingCodeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/authDingCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["code"] = req.Code
	queryParams["codeType"] = req.CodeType
	queryParams["corpId"] = req.CorpId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AuthWeiXinCode(req orgvo.AuthWeiXinCodeReqVo) orgvo.AuthDingCodeRespVo {
	respVo := &orgvo.AuthDingCodeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/authWeiXinCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["code"] = req.Code
	queryParams["codeType"] = req.CodeType
	queryParams["corpId"] = req.CorpId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func BatchGetUserDetailInfo(req orgvo.BatchGetUserInfoReq) orgvo.BatchGetUserInfoResp {
	respVo := &orgvo.BatchGetUserInfoResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/batchGetUserDetailInfo", config.GetPreUrl("orgsvc"))
	requestBody := &req.UserIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func BindCoolApp(req orgvo.BindCoolAppReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/bindCoolApp", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func BindLoginName(req orgvo.BindLoginNameReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/bindLoginName", config.GetPreUrl("orgsvc"))
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

func BoundFeiShu(req orgvo.BoundFsReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/boundFeiShu", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["token"] = req.Token
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func BoundFeiShuAccount(req orgvo.BoundFsAccountReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/boundFeiShuAccount", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["token"] = req.Token
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ChangeDeptScope(req vo.CommonReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/changeDeptScope", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ChangeLoginName(req orgvo.BindLoginNameReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/changeLoginName", config.GetPreUrl("orgsvc"))
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

func CheckLoginName(req orgvo.CheckLoginNameReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/checkLoginName", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CheckMobileHasBind(req orgvo.CheckMobileHasBindReq) orgvo.CheckMobileHasBindResp {
	respVo := &orgvo.CheckMobileHasBindResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/checkMobileHasBind", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CheckSpecificScope(req orgvo.CheckSpecificScopeReqVo) orgvo.CheckSpecificScopeRespVo {
	respVo := &orgvo.CheckSpecificScopeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/checkSpecificScope", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["powerFlag"] = req.PowerFlag
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ClearOrgUsersPayCache(req orgvo.GetBaseOrgInfoReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/clearOrgUsersPayCache", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateCoolApp(req orgvo.CreateCoolAppReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/createCoolApp", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateDepartment(req orgvo.CreateDepartmentReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/createDepartment", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateOrg(req orgvo.CreateOrgReqVo) orgvo.CreateOrgRespVo {
	respVo := &orgvo.CreateOrgRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/createOrg", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateOrgColumn(req orgvo.CreateOrgColumnReq) orgvo.CreateOrgColumnRespVo {
	respVo := &orgvo.CreateOrgColumnRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/createOrgColumn", config.GetPreUrl("orgsvc"))
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

func CreateOrgUser(req orgvo.CreateOrgMemberReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/createOrgUser", config.GetPreUrl("orgsvc"))
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

func DeleteAppUserViewLocation(req orgvo.UpdateUserViewLocationReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/deleteAppUserViewLocation", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["appId"] = req.AppId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteCoolApp(req orgvo.DeleteCoolAppReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/deleteCoolApp", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteCoolAppByProject(req orgvo.DeleteCoolAppByProjectReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/deleteCoolAppByProject", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteDepartment(req orgvo.DeleteDepartmentReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/deleteDepartment", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteOrgColumn(req orgvo.DeleteOrgColumnReq) orgvo.DeleteOrgColumnRespVo {
	respVo := &orgvo.DeleteOrgColumnRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/deleteOrgColumn", config.GetPreUrl("orgsvc"))
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

func DepartmentMembers(req orgvo.DepartmentMembersReqVo) orgvo.DepartmentMembersRespVo {
	respVo := &orgvo.DepartmentMembersRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/departmentMembers", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Params
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DepartmentMembersList(req orgvo.DepartmentMembersListReq) orgvo.DepartmentMembersListResp {
	respVo := &orgvo.DepartmentMembersListResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/departmentMembersList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["ignoreDelete"] = req.IgnoreDelete
	requestBody := &req.Params
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func Departments(req orgvo.DepartmentsReqVo) orgvo.DepartmentsRespVo {
	respVo := &orgvo.DepartmentsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/departments", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Params
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeptChange(req orgvo.FeishuDeptChangeReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/deptChange", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["deptOpenId"] = req.DeptOpenId
	queryParams["eventType"] = req.EventType
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DisbandThirdAccount(req vo.CommonReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/disbandThirdAccount", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func EmptyUser(req orgvo.EmptyUserReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/emptyUser", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ExchangeShareToken(req orgvo.ExchangeShareTokenReq) orgvo.UserSMSLoginRespVo {
	respVo := &orgvo.UserSMSLoginRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/exchangeShareToken", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ExportAddressList(req orgvo.ExportAddressListReqVo) orgvo.ExportAddressListRespVo {
	respVo := &orgvo.ExportAddressListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/exportAddressList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ExportInviteTemplate(req orgvo.ExportInviteTemplateReqVo) orgvo.ExportInviteTemplateRespVo {
	respVo := &orgvo.ExportInviteTemplateRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/exportInviteTemplate", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func FeiShuAuth(req vo.FeiShuAuthReq) orgvo.FeiShuAuthRespVo {
	respVo := &orgvo.FeiShuAuthRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/feiShuAuth", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["code"] = req.Code
	queryParams["codeType"] = req.CodeType
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func FeiShuAuthCode(req vo.FeiShuAuthReq) orgvo.FeiShuAuthCodeRespVo {
	respVo := &orgvo.FeiShuAuthCodeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/feiShuAuthCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["code"] = req.Code
	queryParams["codeType"] = req.CodeType
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func FilterResignedUserIds(req orgvo.FilterResignedUserIdsReqVo) orgvo.FilterResignedUserIdsResp {
	respVo := &orgvo.FilterResignedUserIdsResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/filterResignedUserIds", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.UserIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetAppTicket(req orgvo.GetAppTicketReq) orgvo.GetAppTicketResp {
	respVo := &orgvo.GetAppTicketResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getAppTicket", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetBaseUserInfoBatch(req orgvo.GetBaseUserInfoBatchReqVo) orgvo.GetBaseUserInfoBatchRespVo {
	respVo := &orgvo.GetBaseUserInfoBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getBaseUserInfoBatch", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.UserIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetBaseUserInfoByEmpIdBatch(req orgvo.GetBaseUserInfoByEmpIdBatchReqVo) orgvo.GetBaseUserInfoByEmpIdBatchRespVo {
	respVo := &orgvo.GetBaseUserInfoByEmpIdBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getBaseUserInfoByEmpIdBatch", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetCoolAppInfo(req orgvo.GetCoolAppInfoReq) orgvo.GetCoolAppInfoResp {
	respVo := &orgvo.GetCoolAppInfoResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getCoolAppInfo", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetDeptByIds(req orgvo.GetDeptByIdsReq) orgvo.GetDeptByIdsResp {
	respVo := &orgvo.GetDeptByIdsResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getDeptByIds", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.DeptIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetDingJsAPISign(req orgvo.GetDingApiSignReq) orgvo.GetJsAPISignRespVo {
	respVo := &orgvo.GetJsAPISignRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getDingJsAPISign", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetFsAccessToken(req orgvo.GetFsAccessTokenReqVo) orgvo.GetFsAccessTokenRespVo {
	respVo := &orgvo.GetFsAccessTokenRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getFsAccessToken", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetFunctionConfig(req orgvo.GetOrgFunctionConfigReq) orgvo.GetOrgFunctionConfigResp {
	respVo := &orgvo.GetOrgFunctionConfigResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getFunctionConfig", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetFunctionObjArrByOrg(req orgvo.GetOrgFunctionConfigReq) orgvo.GetFunctionArrByOrgResp {
	respVo := &orgvo.GetFunctionArrByOrgResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getFunctionObjArrByOrg", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetJsAPITicket(req vo.CommonReqVo) orgvo.GetJsAPITicketResp {
	respVo := &orgvo.GetJsAPITicketResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getJsAPITicket", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgBoListByPage(req orgvo.GetOrgIdListByPageReqVo) orgvo.GetOrgBoListByPageRespVo {
	respVo := &orgvo.GetOrgBoListByPageRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgBoListByPage", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgColumns(req orgvo.GetOrgColumnsReq) orgvo.GetOrgColumnsRespVo {
	respVo := &orgvo.GetOrgColumnsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgColumns", config.GetPreUrl("orgsvc"))
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

func GetOrgIdListByPage(req orgvo.GetOrgIdListByPageReqVo) orgvo.GetOrgIdListBySourceChannelRespVo {
	respVo := &orgvo.GetOrgIdListBySourceChannelRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgIdListByPage", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgIdListBySourceChannel(req orgvo.GetOrgIdListBySourceChannelReqVo) orgvo.GetOrgIdListBySourceChannelRespVo {
	respVo := &orgvo.GetOrgIdListBySourceChannelRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgIdListBySourceChannel", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgInfo(req orgvo.GetOrgInfoReq) orgvo.GetOrgInfoResp {
	respVo := &orgvo.GetOrgInfoResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgOutInfoByOutOrgId(req orgvo.GetOutOrgInfoByOutOrgIdReqVo) orgvo.GetOutOrgInfoByOutOrgIdRespVo {
	respVo := &orgvo.GetOutOrgInfoByOutOrgIdRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgOutInfoByOutOrgId", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["outOrgId"] = req.OutOrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgOutInfoByOutOrgIdBatch(req orgvo.GetOrgOutInfoByOutOrgIdBatchReqVo) orgvo.GetOrgOutInfoByOutOrgIdBatchRespVo {
	respVo := &orgvo.GetOrgOutInfoByOutOrgIdBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgOutInfoByOutOrgIdBatch", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgSuperAdminInfo(req orgvo.GetOrgSuperAdminInfoReq) orgvo.GetOrgSuperAdminInfoResp {
	respVo := &orgvo.GetOrgSuperAdminInfoResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgSuperAdminInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgUserIds(req orgvo.GetOrgUserIdsReq) orgvo.GetOrgUserIdsResp {
	respVo := &orgvo.GetOrgUserIdsResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgUserIds", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgUserIdsByEmIds(req orgvo.GetOrgUserIdsByEmIdsReq) orgvo.GetOrgUserIdsByEmIdsResp {
	respVo := &orgvo.GetOrgUserIdsByEmIdsResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgUserIdsByEmIds", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.EmpIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgUserInfoListBySourceChannel(req orgvo.GetOrgUserInfoListBySourceChannelReq) orgvo.GetOrgUserInfoListBySourceChannelResp {
	respVo := &orgvo.GetOrgUserInfoListBySourceChannelResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgUserInfoListBySourceChannel", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgUsersInfoByEmIds(req orgvo.GetOrgUsersInfoByEmIdsReq) orgvo.GetBaseUserInfoBatchRespVo {
	respVo := &orgvo.GetBaseUserInfoBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgUsersInfoByEmIds", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.EmpIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOutOrgInfoByOrgIdBatch(req orgvo.GetOutOrgInfoByOrgIdBatchReqVo) orgvo.GetOutOrgInfoByOrgIdBatchRespVo {
	respVo := &orgvo.GetOutOrgInfoByOrgIdBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOutOrgInfoByOrgIdBatch", config.GetPreUrl("orgsvc"))
	requestBody := &req.OrgIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOutUserInfoListByUserIds(req orgvo.GetOutUserInfoListByUserIdsReqVo) orgvo.GetOutUserInfoListByUserIdsRespVo {
	respVo := &orgvo.GetOutUserInfoListByUserIdsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOutUserInfoListByUserIds", config.GetPreUrl("orgsvc"))
	requestBody := &req.UserIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetPayRemind(req orgvo.GetPayRemindReq) orgvo.GetPayRemindResp {
	respVo := &orgvo.GetPayRemindResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getPayRemind", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetPwdLoginCode(req orgvo.GetPwdLoginCodeReqVo) orgvo.GetPwdLoginCodeRespVo {
	respVo := &orgvo.GetPwdLoginCodeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getPwdLoginCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["captchaId"] = req.CaptchaId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetShareUrl(req orgvo.GetShareUrlReq) orgvo.GetShareUrlResp {
	respVo := &orgvo.GetShareUrlResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getShareUrl", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["key"] = req.Key
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetSpaceList(req orgvo.GetSpaceListReq) orgvo.GetSpaceListResp {
	respVo := &orgvo.GetSpaceListResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getSpaceList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUpdateCoolAppTopCardData(req orgvo.GetCoolAppTopCardDataReq) orgvo.GetCoolAppTopCardDataResp {
	respVo := &orgvo.GetCoolAppTopCardDataResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUpdateCoolAppTopCardData", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserConfigInfoBatch(req orgvo.GetUserConfigInfoBatchReqVo) orgvo.GetUserConfigInfoBatchRespVo {
	respVo := &orgvo.GetUserConfigInfoBatchRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserConfigInfoBatch", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserIds(req orgvo.GetUserIdsReqVo) orgvo.GetUserIdsRespVo {
	respVo := &orgvo.GetUserIdsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserIds", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["corpId"] = req.CorpId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.EmpIdsBody
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserInfoByFeishuTenantKey(req orgvo.GetUserInfoByFeishuTenantKeyReq) orgvo.GetUserInfoByFeishuTenantKeyResp {
	respVo := &orgvo.GetUserInfoByFeishuTenantKeyResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserInfoByFeishuTenantKey", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["tenantKey"] = req.TenantKey
	queryParams["openId"] = req.OpenId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserInfoByUserIds(req orgvo.GetUserInfoByUserIdsReqVo) orgvo.GetUserInfoByUserIdsListRespVo {
	respVo := &orgvo.GetUserInfoByUserIdsListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserInfoByUserIds", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.UserIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserListWithCreateTimeRange(req orgvo.GetUserListWithCreateTimeRangeReqVo) orgvo.GetUserListWithCreateTimeRangeRespVo {
	respVo := &orgvo.GetUserListWithCreateTimeRangeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserListWithCreateTimeRange", config.GetPreUrl("orgsvc"))
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ImportMembers(req orgvo.ImportMembersReqVo) orgvo.ImportMembersRespVo {
	respVo := &orgvo.ImportMembersRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/importMembers", config.GetPreUrl("orgsvc"))
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

func InitFeiShuAccount(req orgvo.InitFsAccountReqVo) orgvo.InitFsAccountRespVo {
	respVo := &orgvo.InitFsAccountRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/initFeiShuAccount", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func InitOrg(req orgvo.InitOrgReqVo) orgvo.OrgInitRespVo {
	respVo := &orgvo.OrgInitRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/initOrg", config.GetPreUrl("orgsvc"))
	requestBody := &req.InitOrg
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func InnerUserInfo(req *orgvo.InnerUserInfosReq) *vo.DataRespVo {
	respVo := &vo.DataRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/innerUserInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func InviteUser(req orgvo.InviteUserReqVo) orgvo.InviteUserRespVo {
	respVo := &orgvo.InviteUserRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/inviteUser", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func InviteUserByPhones(req orgvo.InviteUserByPhonesReqVo) orgvo.InviteUserByPhonesRespVo {
	respVo := &orgvo.InviteUserByPhonesRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/inviteUserByPhones", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func InviteUserList(req orgvo.InviteUserListReqVo) orgvo.InviteUserListRespVo {
	respVo := &orgvo.InviteUserListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/inviteUserList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func JoinOrgByInviteCode(req orgvo.JoinOrgByInviteCodeReq) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/joinOrgByInviteCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["inviteCode"] = req.InviteCode
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func JudgeUserIsAdmin(req orgvo.JudgeUserIsAdminReq) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/judgeUserIsAdmin", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["outUserId"] = req.OutUserId
	queryParams["outOrgId"] = req.OutOrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenOrgUserList(req orgvo.OpenOrgUserListReqVo) orgvo.OpenOrgUserListRespVo {
	respVo := &orgvo.OpenOrgUserListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/openOrgUserList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OrgUserList(req orgvo.OrgUserListReq) orgvo.OrgUserListResp {
	respVo := &orgvo.OrgUserListResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/orgUserList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OrganizationInfo(req orgvo.OrganizationInfoReqVo) orgvo.OrganizationInfoRespVo {
	respVo := &orgvo.OrganizationInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/organizationInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PAReport(req orgvo.PAReportMsgReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/pAReport", config.GetPreUrl("orgsvc"))
	requestBody := &req.Body
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PersonWeiXinBind(req orgvo.PersonWeiXinBindReq) orgvo.PersonWeiXinBindResp {
	respVo := &orgvo.PersonWeiXinBindResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/personWeiXinBind", config.GetPreUrl("orgsvc"))
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PersonWeiXinBindExistAccount(req orgvo.PersonWeiXinBindExistAccountReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/personWeiXinBindExistAccount", config.GetPreUrl("orgsvc"))
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

func PersonWeiXinLogin(req orgvo.PersonWeiXinLoginReqVo) orgvo.PersonWeiXinLoginRespVo {
	respVo := &orgvo.PersonWeiXinLoginRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/personWeiXinLogin", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["code"] = req.Code
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PrepareTransferOrg(req orgvo.PrepareTransferOrgReq) orgvo.PrepareTransferOrgRespVo {
	respVo := &orgvo.PrepareTransferOrgRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/prepareTransferOrg", config.GetPreUrl("orgsvc"))
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

func RemoveInviteUser(req orgvo.RemoveInviteUserReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/removeInviteUser", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ResetPassword(req orgvo.ResetPasswordReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/resetPassword", config.GetPreUrl("orgsvc"))
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

func RetrievePassword(req orgvo.RetrievePasswordReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/retrievePassword", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SaveOrgSummaryTableAppId(req orgvo.SaveOrgSummaryTableAppIdReqVo) orgvo.VoidRespVo {
	respVo := &orgvo.VoidRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/saveOrgSummaryTableAppId", config.GetPreUrl("orgsvc"))
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

func ScheduleOrganizationPageList(req orgvo.ScheduleOrganizationPageListReqVo) orgvo.ScheduleOrganizationPageListRespVo {
	respVo := &orgvo.ScheduleOrganizationPageListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/scheduleOrganizationPageList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SearchUser(req orgvo.SearchUserReqVo) orgvo.SearchUserRespVo {
	respVo := &orgvo.SearchUserRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/searchUser", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SendAuthCode(req orgvo.SendAuthCodeReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/sendAuthCode", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SendCardToAdminForUpgrade(req orgvo.SendCardToAdminForUpgradeReqVo) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/sendCardToAdminForUpgrade", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SendFeishuMemberHelpMsg(req orgvo.SendFeishuMemberHelpMsgReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/sendFeishuMemberHelpMsg", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["tenantKey"] = req.TenantKey
	queryParams["ownerOpenId"] = req.OwnerOpenId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SendSMSLoginCode(req orgvo.SendSMSLoginCodeReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/sendSMSLoginCode", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SetLabConfig(req orgvo.SetLabReqVo) orgvo.SetLabRespVo {
	respVo := &orgvo.SetLabRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/setLabConfig", config.GetPreUrl("orgsvc"))
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

func SetPassword(req orgvo.SetPasswordReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/setPassword", config.GetPreUrl("orgsvc"))
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

func SetPwdLoginCode(req orgvo.SetPwdLoginCodeReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/setPwdLoginCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["captchaId"] = req.CaptchaId
	queryParams["captchaPassword"] = req.CaptchaPassword
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SetShareUrl(req orgvo.SetShareUrlReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/setShareUrl", config.GetPreUrl("orgsvc"))
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SetUserActivity(req orgvo.SetUserActivityReq) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/setUserActivity", config.GetPreUrl("orgsvc"))
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

func SetUserDepartmentLevel(req orgvo.SetUserDepartmentLevelReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/setUserDepartmentLevel", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SetUserViewLocation(req orgvo.SaveViewLocationReqVo) orgvo.SaveViewLocationRespVo {
	respVo := &orgvo.SaveViewLocationRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/setUserViewLocation", config.GetPreUrl("orgsvc"))
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

func SetVersionVisible(req orgvo.SetVersionReq) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/setVersionVisible", config.GetPreUrl("orgsvc"))
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

func StopThirdIntegration(req vo.CommonReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/stopThirdIntegration", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SwitchUserOrganization(req orgvo.SwitchUserOrganizationReqVo) orgvo.SwitchUserOrganizationRespVo {
	respVo := &orgvo.SwitchUserOrganizationRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/switchUserOrganization", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["token"] = req.Token
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SyncUserInfoFromFeiShu(req orgvo.SyncUserInfoFromFeiShuReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/syncUserInfoFromFeiShu", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.SyncUserInfoFromFeiShuReq
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ThirdAccountBindList(req vo.CommonReqVo) orgvo.ThirdAccountListResp {
	respVo := &orgvo.ThirdAccountListResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/thirdAccountBindList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func TransferOrg(req orgvo.PrepareTransferOrgReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/transferOrg", config.GetPreUrl("orgsvc"))
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

func UnbindLoginName(req orgvo.UnbindLoginNameReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/unbindLoginName", config.GetPreUrl("orgsvc"))
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

func UnbindThirdAccount(req orgvo.UnbindAccountReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/unbindThirdAccount", config.GetPreUrl("orgsvc"))
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

func UpdateCoolAppTopCard(req orgvo.UpdateCoolAppTopCardReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateCoolAppTopCard", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateDepartment(req orgvo.UpdateDepartmentReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateDepartment", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateOrgFunctionConfig(req orgvo.UpdateOrgFunctionConfigReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateOrgFunctionConfig", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateOrgMemberCheckStatus(req orgvo.UpdateOrgMemberCheckStatusReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateOrgMemberCheckStatus", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateOrgMemberStatus(req orgvo.UpdateOrgMemberStatusReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateOrgMemberStatus", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateOrgRemarkSetting(req orgvo.UpdateOrgRemarkSettingReqVo) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateOrgRemarkSetting", config.GetPreUrl("orgsvc"))
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

func UpdateOrgUser(req orgvo.UpdateOrgMemberReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateOrgUser", config.GetPreUrl("orgsvc"))
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

func UpdateOrganizationSetting(req orgvo.UpdateOrganizationSettingReqVo) orgvo.UpdateOrganizationSettingRespVo {
	respVo := &orgvo.UpdateOrganizationSettingRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateOrganizationSetting", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateUserConfig(req orgvo.UpdateUserConfigReqVo) orgvo.UpdateUserConfigRespVo {
	respVo := &orgvo.UpdateUserConfigRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateUserConfig", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.UpdateUserConfigReq
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateUserDefaultProjectIdConfig(req orgvo.UpdateUserDefaultProjectIdConfigReqVo) orgvo.UpdateUserConfigRespVo {
	respVo := &orgvo.UpdateUserConfigRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateUserDefaultProjectIdConfig", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.UpdateUserDefaultProjectIdConfigReq
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateUserInfo(req orgvo.UpdateUserInfoReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateUserInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.UpdateUserInfoReq
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateUserPcConfig(req orgvo.UpdateUserPcConfigReqVo) orgvo.UpdateUserConfigRespVo {
	respVo := &orgvo.UpdateUserConfigRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateUserPcConfig", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.UpdateUserPcConfigReq
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateUserToSysManageGroup(req orgvo.UpdateUserToSysManageGroupReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateUserToSysManageGroup", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UserList(req orgvo.UserListReqVo) orgvo.UserListRespVo {
	respVo := &orgvo.UserListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Data
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UserLogin(req orgvo.UserLoginReqVo) orgvo.UserSMSLoginRespVo {
	respVo := &orgvo.UserSMSLoginRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userLogin", config.GetPreUrl("orgsvc"))
	requestBody := &req.UserLoginReq
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UserOrganizationList(req orgvo.UserOrganizationListReqVo) orgvo.UserOrganizationListRespVo {
	respVo := &orgvo.UserOrganizationListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userOrganizationList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UserQuit(req orgvo.UserQuitReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userQuit", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["token"] = req.Token
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UserRegister(req orgvo.UserRegisterReqVo) orgvo.UserRegisterRespVo {
	respVo := &orgvo.UserRegisterRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userRegister", config.GetPreUrl("orgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UserStat(req orgvo.UserStatReqVo) orgvo.UserStatRespVo {
	respVo := &orgvo.UserStatRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userStat", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func VerifyDepartments(req orgvo.VerifyDepartmentsReq) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/verifyDepartments", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.DepartmentIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func VerifyOldName(req orgvo.UnbindLoginNameReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/verifyOldName", config.GetPreUrl("orgsvc"))
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

func VerifyOrgUsers(req orgvo.VerifyOrgUsersReqVo) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/verifyOrgUsers", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func VerifyOrgUsersReturnValid(req orgvo.VerifyOrgUsersReqVo) vo.DataRespVo {
	respVo := &vo.DataRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/verifyOrgUsersReturnValid", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

