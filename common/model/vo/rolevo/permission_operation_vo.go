package rolevo

import "github.com/star-table/polaris-backend/common/model/vo"

type PermissionOperationListReqVo struct {
	OrgId     int64  `json:"orgId"`
	UserId    int64  `json:"userId"`
	RoleId    int64  `json:"roleId"`
	ProjectId *int64 `json:"projectId"`
}

type PermissionOperationListRespVo struct {
	vo.Err
	Data []*vo.PermissionOperationListResp `json:"data"`
}

type UpdateRolePermissionOperationReqVo struct {
	OrgId  int64                               `json:"orgId"`
	UserId int64                               `json:"userId"`
	Input  vo.UpdateRolePermissionOperationReq `json:"input"`
}

type GetPersonalPermissionInfoReqVo struct {
	OrgId     int64  `json:"orgId"`
	UserId    int64  `json:"userId"`
	ProjectId *int64 `json:"projectId"`
	IssueId *int64 `json:"issueId"`
	SourceChannel string `json:"sourceChannel"`
}

type GetPersonalPermissionInfoRespVo struct {
	vo.Err
	Data map[string]interface{} `json:"data"`
}