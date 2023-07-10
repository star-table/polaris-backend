package resourcevo

import "github.com/star-table/polaris-backend/common/model/vo"

type GetOssSignURLRespVo struct {
	vo.Err
	GetOssSignURL *vo.OssApplySignURLResp `json:"data"`
}

type GetOssPostPolicyRespVo struct {
	vo.Err
	GetOssPostPolicy *vo.OssPostPolicyResp `json:"data"`
}

type OssApplySignURLReqVo struct {
	OrgId int64 `json:"orgId"`
	UserId int64 `json:"userId"`
	Input vo.OssApplySignURLReq `json:"input"`
}

type GetOssPostPolicyReqVo struct {
	Input vo.OssPostPolicyReq `json:"input"`
	OrgId int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}
