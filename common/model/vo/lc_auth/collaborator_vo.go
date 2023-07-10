package lc_auth

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

/*
{
    "code": 500,
    "reason": "PERMISSION_ORG_ID_PARSE_FAILED",
    "message": "TableId 不能小于0",
    "metadata": {}
}

{
    "id": "39"
}
*/

type KratosVo struct {
}

type CollaboratorUserRequest struct {
}
type CollaboratorUserRespVo struct {
	vo.Err
	Reason   string `json:"reason"`
	Metadata string `json:"message"`
	Id       int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

type CollaboratorUsersRespVo struct {
	vo.Err
	Reason   string `json:"reason"`
	Metadata string `json:"message"`
	Ok       bool   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

type CollaboratorUsersRequest struct {
	Users []*bo.Collaborator `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

type MemberFieldRequest struct {
	OrgId     int64  `json:"org_id,omitempty"`
	TableId   int64  `json:"table_id,omitempty"`
	FieldName string `json:"field_name,omitempty"`
}

type GetUserCollaboratorRoleIdsRequest struct {
	AppId   int64 `protobuf:"varint,1,opt,name=appId,proto3" json:"appId,omitempty"`
	OrgId   int64 `protobuf:"varint,2,opt,name=orgId,proto3" json:"orgId,omitempty"`
	UserId  int64 `protobuf:"varint,3,opt,name=userId,proto3" json:"userId,omitempty"`
	TableId int64 `protobuf:"varint,4,opt,name=tableId,proto3" json:"tableId,omitempty"`
}

type GetUserCollaboratorRoleIdsRespVo struct {
	vo.Err
	Reason   string  `json:"reason"`
	Metadata string  `json:"message"`
	RoleIds  []int64 `protobuf:"varint,1,rep,packed,name=role_ids,json=roleIds,proto3" json:"roleIds,omitempty"`
}
