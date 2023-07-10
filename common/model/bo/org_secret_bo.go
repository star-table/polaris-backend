package bo

// 开放平台认证信息
type OpenAuthInfo struct {
	Exp int64 `json:"exp"`
	Iat int64 `json:"iat"`
	Key string `json:"_key"`
}

type AppTicketBo struct {
	AppId string `json:"appId"`
	AppSecret string `json:"appSecret"`
}
