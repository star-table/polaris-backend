package websitevo

import "github.com/star-table/polaris-backend/common/model/vo"

type RegisterWebSiteContactReqVo struct {
	//可以为空
	OrgId int64 `json:"orgId"`
	//可以为空
	UserId int64                        `json:"userId"`
	Input  vo.RegisterWebSiteContactReq `json:"input"`
}
