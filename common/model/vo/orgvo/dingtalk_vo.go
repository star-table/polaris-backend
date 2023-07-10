package orgvo

import "github.com/star-table/polaris-backend/common/model/vo"

type GetJsAPISignRespVo struct {
	vo.Err
	GetJsAPISign *JsAPISignResp `json:"data"`
}

//type AuthRespVo struct {
//	vo.Err
//	Auth *vo.AuthResp `json:"data"`
//}

type AuthDingCodeReqVo struct {
	Code     string `json:"code"`
	CodeType int    `json:"codeType"`
	CorpId   string `json:"corpId"`
}

type AuthDingCodeRespVo struct {
	vo.Err
	Data *PlatformAuthCodeData `json:"data"`
}

type PlatformAuthCodeData struct {
	// 企业ID
	CorpId string `json:"corpId"`
	// OutUserId
	OutUserId string `json:"outUserId"`
	// 是否为企业管理
	IsAdmin bool `json:"isAdmin"`
	// 是否被绑定
	Binding bool `json:"binding"`
	// refreshToken
	RefreshToken string `json:"refreshToken"`
	// accessToken
	AccessToken string `json:"accessToken"`
	// token
	Token string `json:"token"`
	// codeToken
	CodeToken string `json:"codeToken"`
	// 组织id
	OrgID int64 `json:"orgId"`
	// 组织名称
	OrgName string `json:"orgName"`
	// 外部组织名称
	OutOrgName string `json:"outOrgName"`
	// 组织code
	OrgCode string `json:"orgCode"`
	// 用户id
	UserID int64 `json:"userId"`
	// 用户姓名
	Name string `json:"name"`
}
