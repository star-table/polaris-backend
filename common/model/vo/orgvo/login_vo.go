package orgvo

import (
	"github.com/star-table/polaris-backend/common/model/vo"
)

type SendSMSLoginCodeReqVo struct {
	Input vo.SendSmsLoginCodeReq `json:"input"`
}

type SendAuthCodeReqVo struct {
	Input vo.SendAuthCodeReq `json:"input"`
}

type UserLoginReqVo struct {
	UserLoginReq vo.UserLoginReq `json:"userLoginReq"`
}

type UserSMSLoginRespVo struct {
	vo.Err
	Data *vo.UserLoginResp `json:"data"`
}

type GetInviteCodeReqVo struct {
	CurrentUserId  int64  `json:"userId"`
	OrgId          int64  `json:"orgId"`
	SourcePlatform string `json:"sourcePlatform"`
}

type InviteUserByPhonesReqVo struct {
	UserId int64                       `json:"userId"`
	OrgId  int64                       `json:"orgId"`
	Input  InviteUserByPhonesReqVoData `json:"input"`
}

type InviteUserByPhonesReqVoData struct {
	InviteCode string                             `json:"inviteCode"`
	Phones     []InviteUserByPhonesReqVoDataPhone `json:"phones"`
}

type InviteUserByPhonesReqVoDataPhone struct {
	Origin string `json:"origin"`
	Number string `json:"number"`
}

type GetInviteInfoReqVo struct {
	InviteCode string `json:"inviteCode"`
}

type GetInviteInfoRespVo struct {
	Data *vo.GetInviteInfoResp `json:"data"`
	vo.Err
}

type GetInviteCodeRespVo struct {
	Data *GetInviteCodeRespVoData `json:"data"`
	vo.Err
}

type GetInviteCodeRespVoData struct {
	InviteCode string `json:"inviteCode"`
	Expire     int    `json:"expire"`
}

type InviteUserByPhonesRespVo struct {
	vo.Err
	Data *InviteUserByPhonesRespVoData `json:"data"`
}

type InviteUserByPhonesRespVoData struct {
	ErrorList   []InviteUserByPhonesRespVoDataErrItem     `json:"errorList"`
	SuccessList []InviteUserByPhonesRespVoDataSuccessItem `json:"successList"`
}

type InviteUserByPhonesRespVoDataErrItem struct {
	Number string `json:"number"`
	Reason string `json:"reason"`
}

type InviteUserByPhonesRespVoDataSuccessItem struct {
	Number string `json:"number"`
}

type ExportInviteTemplateReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type ExportInviteTemplateRespVo struct {
	vo.Err
	Data *ExportInviteTemplateRespVoData `json:"data"`
}

type ExportInviteTemplateRespVoData struct {
	// 下载模板的 url
	Url string `json:"url"`
}

type ImportMembersReqVo struct {
	OrgId  int64                   `json:"orgId"`
	UserId int64                   `json:"userId"`
	Input  *ImportMembersReqVoData `json:"input"`
}

type ImportMembersReqVoData struct {
	// excel 文件地址
	URL string `json:"url"`
	// url类型, 1 网址，2 本地dist路径
	URLType int `json:"urlType"`
	// 方案2：前端直接传入数据进行导入
	ImportUserList []ImportMembersReqVoDataUserItem `json:"importUserList"`
}

type ImportMembersReqVoDataUserItem struct {
	// 姓名
	Name string `json:"name"`
	// 国家代码，比如中国为 `+86`
	Origin string `json:"origin"`
	// 手机号码
	Phone string `json:"phone"`
}

type ImportMembersRespVo struct {
	vo.Err
	Data *ImportMembersRespVoData `json:"data"`
}

type ImportMembersRespVoData struct {
	ErrCount     int    `json:"errCount"`
	SuccessCount int    `json:"successCount"`
	SkipCount    int    `json:"skipCount"`
	ErrExcelUrl  string `json:"errorUrl"`
}

type ImportMembersResultErrItem struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	Phone  string `json:"phone"`
	Reason string `json:"reason"`
}

type ImportUserInfo struct {
	OperateUid    int64  `json:"operateUid"` // 操作人、导入人
	OrgId         int64  `json:"orgId"`
	SourceChannel string `json:"sourceChannel"`
}

type UserQuitReqVo struct {
	Token string `json:"token"`
}

type GetPwdLoginCodeReqVo struct {
	CaptchaId string `json:"captchaId"`
}

type GetPwdLoginCodeRespVo struct {
	CaptchaPassword string `json:"data"`
	vo.Err
}

type SetPwdLoginCodeReqVo struct {
	CaptchaId       string `json:"captchaId"`
	CaptchaPassword string `json:"captchaPassword"`
}

type JoinOrgByInviteCodeReq struct {
	OrgId      int64  `json:"org_id"`
	UserId     int64  `json:"user_id"`
	InviteCode string `json:"invite_code"`
}

type AddUserReqVo struct {
	OrgId  int64                           `json:"orgId"`
	UserId int64                           `json:"userId"`
	Input  *ImportMembersReqVoDataUserItem `json:"input"`
}
