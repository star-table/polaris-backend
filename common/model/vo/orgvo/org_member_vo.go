package orgvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type UpdateOrgMemberStatusReq struct {
	UserId        int64  `json:"userId"`
	OrgId         int64  `json:"orgId"`
	SourceChannel string `json:"sourceChannel"`

	Input vo.UpdateOrgMemberStatusReq `json:"input"`
}

type UpdateOrgMemberCheckStatusReq struct {
	UserId        int64  `json:"userId"`
	OrgId         int64  `json:"orgId"`
	SourceChannel string `json:"sourceChannel"`

	Input vo.UpdateOrgMemberCheckStatusReq `json:"input"`
}

type RemoveOrgMemberReq struct {
	UserId        int64  `json:"userId"`
	OrgId         int64  `json:"orgId"`
	SourceChannel string `json:"sourceChannel"`

	Input vo.RemoveOrgMemberReq `json:"input"`
}

type OrgUserListReq struct {
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	OrgId  int64              `json:"orgId"`
	UserId int64              `json:"userId"`
	Input  *vo.OrgUserListReq `json:"input"`
}

type OrgUserListResp struct {
	vo.Err
	Data *vo.UserOrganizationList `json:"data"`
}

type GetOrgUserInfoListBySourceChannelReq struct {
	Page          int    `json:"page"`
	Size          int    `json:"size"`
	OrgId         int64  `json:"orgId"`
	SourceChannel string `json:"sourceChannel"`
}

type GetOrgUserInfoListBySourceChannelResp struct {
	Data *GetOrgUserInfoListBySourceChannelRespData `json:"data"`
	vo.Err
}

type GetOrgUserInfoListBySourceChannelRespData struct {
	Total int64            `json:"total"`
	List  []bo.OrgUserInfo `json:"list"`
}

type BatchGetUserInfoReq struct {
	UserIds []int64 `json:"userIds"`
}

type BatchGetUserInfoResp struct {
	vo.Err
	Data []vo.PersonalInfo `json:"data"`
}

type JudgeUserIsAdminReq struct {
	SourceChannel string `json:"sourceChannel"`
	OutUserId     string `json:"outUserId"`
	OutOrgId      string `json:"outOrgId"`
}
