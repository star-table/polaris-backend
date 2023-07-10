package orgvo

import (
	"github.com/star-table/polaris-backend/common/model/vo"
)

type AuthWeiXinCodeReqVo struct {
	Code     string `json:"code"`
	CodeType int    `json:"codeType"`
	CorpId   string `json:"corpId"`
}

type AuthWeiXinCodeRespVo struct {
	vo.Err
	Data *PlatformAuthCodeData `json:"data"`
}

type GetDingApiSignReq struct {
	Input *JsAPISignReq `json:"input"`
}

// 获取JSApi签名请求结构体
type JsAPISignReq struct {
	// 类型:目前只支持:jsapi
	Type string `json:"type"`
	// 路由url
	URL string `json:"url"`
	// 企业id
	CorpID string `json:"corpId"`
}

type JsAPISignResp struct {
	// 应用代理id
	AgentID int64 `json:"agentId"`
	// 时间戳
	TimeStamp string `json:"timeStamp"`
	// 随机字符串
	NoceStr string `json:"noceStr"`
	// 签名
	Signature string `json:"signature"`
}

type PersonWeiXinLoginReqVo struct {
	Code string `json:"code"`
}

type PersonWeiXinLoginRespVo struct {
	vo.Err
	Data *PersonWeiXinLoginData `json:"data"`
}

type PersonWeiXinLoginData struct {
	// 是否被绑定
	Binding bool `json:"binding"`
	// 如果已经绑定了用户，则会返回token
	Token string `json:"token"`
	// 组织id，如果为0，则要创建组织
	OrgID int64 `json:"orgId"`
	// 组织名称
	OrgName string `json:"orgName"`
	// 用户id
	UserID int64 `json:"userId"`
	// 用户姓名
	Name string `json:"name"`
}

type AccessTokenResp struct {
	ErrorInfo
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`
}

type ErrorInfo struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (e *ErrorInfo) GetCode() int {
	return e.ErrCode
}

type PersonWeiXinQrCodeReqVo struct {
	Code string
}

type PersonWeiXinQrCodeRespVo struct {
	vo.Err
	Data *PersonWeiXinRqCodeData `json:"data"`
}

type PersonWeiXinRqCodeData struct {
	Ticket   string `json:"ticket"`
	SceneKey string `json:"sceneKey"`
}

type WeiXinQrCodeReq struct {
	ExpireSeconds int         `json:"expire_seconds"`
	ActionName    string      `json:"action_name"`
	ActionInfo    *ActionInfo `json:"action_info"`
}

type ActionInfo struct {
	Scene *Scene `json:"scene"`
}

type Scene struct {
	SceneId  int    `json:"scene_id"`
	SceneStr string `json:"scene_str"`
}
type WeiXinQrCodeReply struct {
	ErrorInfo
	Ticket        string `json:"ticket"`
	ExpireSeconds int    `json:"expire_seconds"`
	Url           string `json:"url"`
}

type QrCodeScanReq struct {
	OpenId   string `json:"openId"`
	SceneKey string `json:"token"`
}

type CheckQrCodeScanReq struct {
	SceneKey string `json:"sceneKey"`
}

type CheckQrCodeScanResp struct {
	vo.Err
	Data *CheckQrCodeScanData `json:"data"`
}

type CheckQrCodeScanData struct {
	// 是否已经扫描
	IsScan bool `json:"isScan"`
	// 是否被绑定
	Binding bool `json:"binding"`
	// 如果已经绑定了用户，则会返回token
	Token string `json:"token"`
	// 组织id，如果为0，则要创建组织
	OrgID int64 `json:"orgId"`
	// 组织名称
	OrgName string `json:"orgName"`
	// 用户id
	UserID int64 `json:"userId"`
	// 用户姓名
	Name string `json:"name"`
}

type PersonWeiXinBindReq struct {
	Data *PersonWeiXinBindData `json:"data"`
}

type PersonWeiXinBindData struct {
	SceneKey string `json:"sceneKey"`
	Mobile   string `json:"mobile"`
	AuthCode string `json:"authCode"`
}

type PersonWeiXinBindResp struct {
	vo.Err
	Data *PersonWeiXinLoginData `json:"data"`
}

type CreateCoolAppData struct {
	OpenConversationId string `json:"openConversationId"`
	CorpId             string `json:"corpId"`
	RobotCode          string `json:"robotCode"`
}

type CreateCoolAppReq struct {
	Input *CreateCoolAppData `json:"input"`
}

type DeleteCoolAppReq struct {
	Input *CoolAppBaseData `json:"input"`
}

type DeleteCoolAppByProjectReq struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type BindCoolAppData struct {
	OpenConversationId string `json:"openConversationId"`
	ProjectId          int64  `json:"projectId"`
	AppId              int64  `json:"appId"`
}

type BindCoolAppReq struct {
	OrgId int64            `json:"orgId"`
	Input *BindCoolAppData `json:"input"`
}

type UpdateCoolAppTopCardReq struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type GetCoolAppTopCardDataReq struct {
	Input *CoolAppBaseData `json:"input"`
}

type GetCoolAppTopCardDataResp struct {
	vo.Err
	Data string `json:"data"`
}

type CoolAppBaseData struct {
	OpenConversationId string `json:"openConversationId"`
}

type GetCoolAppInfoReq struct {
	Input *CoolAppBaseData `json:"input"`
}

type GetCoolAppInfoResp struct {
	vo.Err
	Data *CoolAppInfo `json:"data"`
}

type CoolAppInfo struct {
	ProjectId int64 `json:"projectId"`
	AppId     int64 `json:"appId,string"`
}

type CheckMobileHasBindData struct {
	Mobile         string `json:"mobile"`
	SourcePlatform string `json:"sourcePlatform"`
}

type CheckMobileHasBindReq struct {
	Input *CheckMobileHasBindData `json:"input"`
}

type CheckMobileHasBindResult struct {
	IsBind bool `json:"isBind"`
}

type CheckMobileHasBindResp struct {
	vo.Err
	Data *CheckMobileHasBindResult `json:"data"`
}

type UnbindAccountData struct {
	SourcePlatform string `json:"sourcePlatform"`
}

type UnbindAccountReq struct {
	OrgId  int64              `json:"orgId"`
	UserId int64              `json:"userId"`
	Input  *UnbindAccountData `json:"input"`
}

type ThirdAccountData struct {
	BindPlatform string `json:"bindPlatform"`
	BindStatus   int32  `json:"bindStatus"`
	BindTime     int64  `json:"bindTime"`
	BindName     string `json:"bindName"`
}

type ThirdAccountListData struct {
	IsBindMobile bool                `json:"isBindMobile"`
	List         []*ThirdAccountData `json:"list"`
}

type ThirdAccountListResp struct {
	vo.Err
	Data *ThirdAccountListData `json:"data"`
}

type PersonWeiXinBindExistAccountData struct {
	SceneKey string `json:"sceneKey"`
}

type PersonWeiXinBindExistAccountReq struct {
	OrgId  int64                             `json:"orgId"`
	UserId int64                             `json:"userId"`
	Input  *PersonWeiXinBindExistAccountData `json:"input"`
}

type GetSpaceListReq struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type SpaceInfo struct {
	SpaceId   string `json:"spaceId"`
	SpaceName string `json:"spaceName"`
}

type SpaceList struct {
	List []*SpaceInfo `json:"list"`
}

type GetSpaceListResp struct {
	vo.Err
	Data *SpaceList `json:"data"`
}

type GetRegisterUrlReq struct {
}

type GetRegisterUrlData struct {
	Url string `json:"url"`
}

type GetRegisterUrlResp struct {
	vo.Err
	Data *GetRegisterUrlData `json:"data"`
}
