package orgfacade

import (
	"context"
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/extra/gin/util"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade"
)

var _ = commonvo.DataEvent{}

func CheckAndSetSuperAdmin(req orgvo.CheckAndSetSuperAdminReq) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/checkAndSetSuperAdmin", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CheckPersonWeiXinQrCode(req orgvo.CheckQrCodeScanReq) orgvo.CheckQrCodeScanResp {
	respVo := &orgvo.CheckQrCodeScanResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/checkPersonWeiXinQrCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sceneKey"] = req.SceneKey
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetBaseOrgInfo(req orgvo.GetBaseOrgInfoReqVo) orgvo.GetBaseOrgInfoRespVo {
	respVo := &orgvo.GetBaseOrgInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getBaseOrgInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetBaseOrgInfoByOutOrgId(req orgvo.GetBaseOrgInfoByOutOrgIdReqVo) orgvo.GetBaseOrgInfoByOutOrgIdRespVo {
	respVo := &orgvo.GetBaseOrgInfoByOutOrgIdRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getBaseOrgInfoByOutOrgId", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["outOrgId"] = req.OutOrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetBaseUserInfo(req orgvo.GetBaseUserInfoReqVo) orgvo.GetBaseUserInfoRespVo {
	respVo := &orgvo.GetBaseUserInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getBaseUserInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetBaseUserInfoByEmpId(req orgvo.GetBaseUserInfoByEmpIdReqVo) orgvo.GetBaseUserInfoByEmpIdRespVo {
	respVo := &orgvo.GetBaseUserInfoByEmpIdRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getBaseUserInfoByEmpId", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["empId"] = req.EmpId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetCurrentUser(ctx context.Context) orgvo.CacheUserInfoVo {
	respVo := &orgvo.CacheUserInfoVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getCurrentUser", config.GetPreUrl("orgsvc"))
	headerOptions, _ := util.BuildHeaderOptions(ctx)
	err := facade.Request(consts.HttpMethodGet, reqUrl, nil, headerOptions, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetCurrentUserWithoutOrgVerify(ctx context.Context) orgvo.CacheUserInfoVo {
	respVo := &orgvo.CacheUserInfoVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getCurrentUserWithoutOrgVerify", config.GetPreUrl("orgsvc"))
	headerOptions, _ := util.BuildHeaderOptions(ctx)
	err := facade.Request(consts.HttpMethodGet, reqUrl, nil, headerOptions, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetCurrentUserWithoutPayVerify(ctx context.Context) orgvo.CacheUserInfoVo {
	respVo := &orgvo.CacheUserInfoVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getCurrentUserWithoutPayVerify", config.GetPreUrl("orgsvc"))
	headerOptions, _ := util.BuildHeaderOptions(ctx)
	err := facade.Request(consts.HttpMethodGet, reqUrl, nil, headerOptions, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetDingTalkBaseUserInfo(req orgvo.GetDingTalkBaseUserInfoReqVo) orgvo.GetBaseUserInfoRespVo {
	respVo := &orgvo.GetBaseUserInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getDingTalkBaseUserInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetInviteCode(req orgvo.GetInviteCodeReqVo) orgvo.GetInviteCodeRespVo {
	respVo := &orgvo.GetInviteCodeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getInviteCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourcePlatform"] = req.SourcePlatform
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetInviteInfo(req orgvo.GetInviteInfoReqVo) orgvo.GetInviteInfoRespVo {
	respVo := &orgvo.GetInviteInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getInviteInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["inviteCode"] = req.InviteCode
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetLabConfig(req orgvo.GetLabReqVo) orgvo.GetLabRespVo {
	respVo := &orgvo.GetLabRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getLabConfig", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgBoList() orgvo.GetOrgBoListRespVo {
	respVo := &orgvo.GetOrgBoListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgBoList", config.GetPreUrl("orgsvc"))
	err := facade.Request(consts.HttpMethodGet, reqUrl, nil, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgConfig(req orgvo.GetOrgConfigReq) orgvo.GetOrgConfigResp {
	respVo := &orgvo.GetOrgConfigResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgConfig", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgConfigRest(req orgvo.GetOrgConfigReq) orgvo.GetOrgConfigRestResp {
	respVo := &orgvo.GetOrgConfigRestResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgConfigRest", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOrgUserAndDeptCount(req orgvo.GetOrgUserAndDeptCountReq) orgvo.GetOrgUserAndDeptCountResp {
	respVo := &orgvo.GetOrgUserAndDeptCountResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgUserAndDeptCount", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOutUserInfoListBySourceChannel(req orgvo.GetOutUserInfoListBySourceChannelReqVo) orgvo.GetOutUserInfoListBySourceChannelRespVo {
	respVo := &orgvo.GetOutUserInfoListBySourceChannelRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOutUserInfoListBySourceChannel", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserConfigInfo(req orgvo.GetUserConfigInfoReqVo) orgvo.GetUserConfigInfoRespVo {
	respVo := &orgvo.GetUserConfigInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserConfigInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserId(req orgvo.GetUserIdReqVo) orgvo.GetUserIdRespVo {
	respVo := &orgvo.GetUserIdRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserId", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["empId"] = req.EmpId
	queryParams["corpId"] = req.CorpId
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserInfo(req orgvo.GetUserInfoReqVo) orgvo.GetUserInfoRespVo {
	respVo := &orgvo.GetUserInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserInfoListByOrg(req orgvo.GetUserInfoListReqVo) orgvo.GetUserInfoListRespVo {
	respVo := &orgvo.GetUserInfoListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserInfoListByOrg", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetUserViewLastLocation(req orgvo.GetViewLocationReq) orgvo.GetViewLocationRespVo {
	respVo := &orgvo.GetViewLocationRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserViewLastLocation", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetVersion(req orgvo.GetVersionReq) orgvo.GetVersionRespVo {
	respVo := &orgvo.GetVersionRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getVersion", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetWeiXinInstallUrl(req orgvo.GetRegisterUrlReq) orgvo.GetRegisterUrlResp {
	respVo := &orgvo.GetRegisterUrlResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getWeiXinInstallUrl", config.GetPreUrl("orgsvc"))
	err := facade.Request(consts.HttpMethodGet, reqUrl, nil, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetWeiXinJsAPISign(req orgvo.JsAPISignReq) orgvo.GetJsAPISignRespVo {
	respVo := &orgvo.GetJsAPISignRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getWeiXinJsAPISign", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["type"] = req.Type
	queryParams["uRL"] = req.URL
	queryParams["corpID"] = req.CorpID
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetWeiXinRegisterUrl(req orgvo.GetRegisterUrlReq) orgvo.GetRegisterUrlResp {
	respVo := &orgvo.GetRegisterUrlResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getWeiXinRegisterUrl", config.GetPreUrl("orgsvc"))
	err := facade.Request(consts.HttpMethodGet, reqUrl, nil, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func OpenAPIAuth(req orgvo.OpenAPIAuthReq) orgvo.OpenAPIAuthResp {
	respVo := &orgvo.OpenAPIAuthResp{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/openAPIAuth", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgID"] = req.OrgID
	queryParams["accessToken"] = req.AccessToken
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PersonWeiXinQrCode(req orgvo.PersonWeiXinQrCodeReqVo) orgvo.PersonWeiXinQrCodeRespVo {
	respVo := &orgvo.PersonWeiXinQrCodeRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/personWeiXinQrCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["code"] = req.Code
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PersonWeiXinQrCodeScan(req orgvo.QrCodeScanReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/personWeiXinQrCodeScan", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["openId"] = req.OpenId
	queryParams["sceneKey"] = req.SceneKey
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PersonalInfo(req orgvo.PersonalInfoReqVo) orgvo.PersonalInfoRespVo {
	respVo := &orgvo.PersonalInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/personalInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PersonalInfoRest(req orgvo.PersonalInfoReqVo) orgvo.PersonalInfoRestRespVo {
	respVo := &orgvo.PersonalInfoRestRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/personalInfoRest", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func ScheduleOrgUseMobileAndEmail(req orgvo.ScheduleOrgUseMobileAndEmailReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/scheduleOrgUseMobileAndEmail", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["tenantKey"] = req.TenantKey
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SetVisitUserGuideStatus(req orgvo.SetVisitUserGuideStatusReq) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/setVisitUserGuideStatus", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["flag"] = req.Flag
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UserConfigInfo(req orgvo.UserConfigInfoReqVo) orgvo.UserConfigInfoRespVo {
	respVo := &orgvo.UserConfigInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userConfigInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func VerifyOrg(req orgvo.VerifyOrgReqVo) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	reqUrl := fmt.Sprintf("%s/api/orgsvc/verifyOrg", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

