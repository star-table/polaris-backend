package orgvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type GetBaseOrgInfoReqVo struct {
	OrgId int64 `json:"orgId"`
}

type GetBaseOrgInfoByOutOrgIdReqVo struct {
	OutOrgId string `json:"outOrgId"`
}

type GetBaseOrgInfoRespVo struct {
	vo.Err
	BaseOrgInfo *bo.BaseOrgInfoBo `json:"data"`
}

type GetBaseOrgInfoByOutOrgIdRespVo struct {
	vo.Err
	BaseOrgInfo *bo.BaseOrgInfoBo `json:"data"`
}

type GetOrgOutInfoByOutOrgIdBatchReqVo struct {
	Input GetOrgOutInfoByOutOrgIdBatchReqInput `json:"input"`
}

type GetOrgOutInfoByOutOrgIdBatchReqInput struct {
	OutOrgIds []string `json:"outOrgIds"`
}

type GetOrgOutInfoByOutOrgIdBatchRespVo struct {
	vo.Err
	Data []bo.BaseOrgOutInfoBo `json:"data"`
}

type GetDingTalkBaseUserInfoByEmpIdReqVo struct {
	OrgId int64  `json:"orgId"`
	EmpId string `json:"empId"`
}

type GetBaseUserInfoByEmpIdReqVo struct {
	OrgId int64  `json:"orgId"`
	EmpId string `json:"empId"`
}

type GetBaseUserInfoByEmpIdBatchReqVo struct {
	OrgId int64                                 `json:"orgId"`
	Input GetBaseUserInfoByEmpIdBatchReqVoInput `json:"input"`
}
type GetBaseUserInfoByEmpIdBatchReqVoInput struct {
	OpenIds []string `json:"openIds"`
}

type GetBaseUserInfoByEmpIdRespVo struct {
	vo.Err
	BaseUserInfo *bo.BaseUserInfoBo `json:"data"`
}

type GetBaseUserInfoByEmpIdBatchRespVo struct {
	vo.Err
	Data []bo.BaseUserInfoBo `json:"data"`
}

type GetDingTalkBaseUserInfoByEmpIdRespVo struct {
	vo.Err
	DingTalkBaseUserInfo *bo.BaseUserInfoBo `json:"data"`
}

type GetUserConfigInfoReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetUserConfigInfoBatchReqVo struct {
	OrgId int64                            `json:"orgId"`
	Input GetUserConfigInfoBatchReqVoInput `json:"input"`
}

type GetUserConfigInfoBatchReqVoInput struct {
	UserIds []int64 `json:"userIds"`
}

type GetUserConfigInfoRespVo struct {
	vo.Err
	UserConfigInfo *bo.UserConfigBo `json:"data"`
}

type GetUserConfigInfoBatchRespVo struct {
	vo.Err
	Data []bo.UserConfigBo `json:"data"`
}

type GetBaseUserInfoReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetBaseUserInfoRespVo struct {
	vo.Err
	BaseUserInfo *bo.BaseUserInfoBo `json:"data"`
}

type GetDingTalkBaseUserInfoReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetBaseUserInfoBatchReqVo struct {
	OrgId   int64   `json:"orgId"`
	UserIds []int64 `json:"userIds"`
}

type GetBaseUserInfoBatchRespVo struct {
	vo.Err
	BaseUserInfos []bo.BaseUserInfoBo `json:"data"`
}

type GetShareUrlReq struct {
	Key string `json:"key"`
}

type GetShareUrlResp struct {
	vo.Err
	Url string `json:"url"`
}

type SetShareUrlReq struct {
	Data ShareUrlData `json:"data"`
}

type ShareUrlData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
