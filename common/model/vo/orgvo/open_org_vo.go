package orgvo

import "github.com/star-table/polaris-backend/common/model/vo"

type OpenOrgUserListReqVo struct {
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	OrgId  int64              `json:"orgId"`
	UserId int64              `json:"userId"`
	Input  *vo.OrgUserListReq `json:"input"`
}

type OpenOrgUserListRespVo struct {
	vo.Err
	Data *vo.UserOrganizationList `json:"data"`
}
