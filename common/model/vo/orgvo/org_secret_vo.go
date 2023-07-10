package orgvo

import "github.com/star-table/polaris-backend/common/model/vo"

type OpenAPIAuthReq struct {
	OrgID int64 `json:"orgId"`
	AccessToken string `json:"accessToken"`
}

type OpenAPIAuthResp struct {
	vo.Err
	Data *OpenAPIAuthData `json:"data"`
}

type OpenAPIAuthData struct {
	// 组织ID
	OrgID int64 `json:"orgId"`
	// 所拥有的权限(暂时不用)
	Permissions []string `json:"permissions"`
}