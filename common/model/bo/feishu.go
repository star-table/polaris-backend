package bo

import "time"

type FeiShuAuthCodeInfo struct {
	TenantKey string `json:"tenantKey"`
	OpenId    string `json:"openId"`
	Binding   bool   `json:"binding"`
}

type PlatformAppTicketCacheBo struct {
	AppId          string `json:"appId"`
	AppTicket      string `json:"appTicket"`
	LastUpdateTime string `json:"lastUpdateTime"`
}

type PlatformAppAccessTokenCacheBo struct {
	AppAccessToken string `json:"appAccessToken"`
	LastUpdateTime string `json:"lastUpdateTime"`
}

type PlatformCorpAccessTokenCacheBo struct {
	TenantAccessToken string `json:"tenantAccessToken"`
	TenantKey         string `json:"tenantKey"`
	LastUpdateTime    string `json:"lastUpdateTime"`
}

type FeiShuCallBackMqBo struct {
	//类型
	Type string `json:"type"`
	//数据
	Data string `json:"data"`
	//消息时间
	MsgTime time.Time `json:"msgTime"`
}
