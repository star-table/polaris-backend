package rolevo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type AuthenticateReqVo struct {
	OrgId int64	`json:"orgId"`
	UserId int64 `json:"userId"`
	Path string `json:"path"`
	Operation string `json:"operation"`
	AuthInfoReqVo AuthenticateAuthInfoReqVo `json:"authInfoVo"`
}

type AuthenticateAuthInfoReqVo struct {
	ProjectAuthInfo *bo.ProjectAuthBo `json:"projectAuthInfo"`
	IssueAuthInfo *bo.IssueAuthBo `json:"issueAuthInfo"`
	Fields []string `json:"fields"`
}

type RoleUserRelationReqVo struct {
	OrgId int64 `json:"orgId"`
	UserId int64 `json:"userId"`
	RoleId int64 `json:"roleId"`
}

type RemoveRoleUserRelationReqVo struct {
	OrgId int64 `json:"orgId"`
	UserIds []int64 `json:"userIds"`
	OperatorId int64 `json:"operatorId"`
}

type RoleInitReqVo struct {
	OrgId int64 `json:"orgId"`
}

type RoleInitRespVo struct {
	RoleInitResp *bo.RoleInitResp `json:"data"`

	vo.Err
}

type RemoveRoleDepartmentRelationReqVo struct {
	OrgId int64 `json:"orgId"`
	DeptIds []int64 `json:"userIds"`
	OperatorId int64 `json:"operatorId"`
}
