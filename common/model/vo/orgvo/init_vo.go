package orgvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type OrgInitReqVo struct {
	CorpId        string `json:"corpId"`
	PermanentCode string `json:"permanentCode"`
}

type InitOrgReqVo struct {
	InitOrg bo.InitOrgBo `json:"initOrg"`
}

type OrgInitRespVo struct {
	vo.Err
	OrgId int64 `json:"data"`
}

type OrgOwnerInitReqVo struct {
	OrgId int64 `json:"orgId"`
	Owner int64 `json:"owner"`
}

type OrgVo struct {
	OrgId int64 `json:"orgId"`
}

type TeamInitRespVo struct {
	vo.Err
	TeamId int64 `json:"data"`
}

type TeamOwnerInitReqVo struct {
	TeamId int64 `json:"teamId"`
	Owner  int64 `json:"owner"`
}

type TeamUserInitReqVo struct {
	OrgId  int64 `json:"orgId"`
	TeamId int64 `json:"teamId"`
	UserId int64 `json:"userId"`
	IsRoot bool  `json:"isRoot"`
}

type UserInitByOrgReqVo struct {
	UserId string `json:"userId"`
	CorpId string `json:"corpId"`
	OrgId  int64  `json:"orgId"`
}

type UserInitByOrgRespVo struct {
	vo.Err
	UserId int64 `json:"data"`
}

type DepartmentInitReqVo struct {
	OrgId         int64  `json:"orgId"`
	CorpId        string `json:"corpId"`
	SourceChannel string `json:"sourceChannel"`
}

type SendFeishuMemberHelpMsgReqVo struct {
	OrgId       int64  `json:"orgId"`
	TenantKey   string `json:"tenantKey"`
	OwnerOpenId string `json:"ownerOpenId"`
}

type ScheduleOrgUseMobileAndEmailReqVo struct {
	TenantKey string `json:"tenantKey"`
}
