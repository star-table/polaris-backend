package orgvo

import "github.com/star-table/polaris-backend/common/model/vo"

type SetPasswordReqVo struct {
	OrgId  int64             `json:"orgId"`
	UserId int64             `json:"userId"`
	Input  vo.SetPasswordReq `json:"input"`
}

type ResetPasswordReqVo struct {
	OrgId  int64               `json:"orgId"`
	UserId int64               `json:"userId"`
	Input  vo.ResetPasswordReq `json:"input"`
}

type RetrievePasswordReqVo struct {
	OrgId int64                  `json:"OrgId"`
	Input vo.RetrievePasswordReq `json:"input"`
}

type UnbindLoginNameReqVo struct {
	OrgId  int64                 `json:"orgId"`
	UserId int64                 `json:"userId"`
	Input  vo.UnbindLoginNameReq `json:"input"`
}

type BindLoginNameReqVo struct {
	OrgId  int64               `json:"orgId"`
	UserId int64               `json:"userId"`
	Input  vo.BindLoginNameReq `json:"input"`
}

type CheckLoginNameReqVo struct {
	Input vo.CheckLoginNameReq `json:"input"`
}

type ThirdAccountBindListResp struct {
	vo.Err
	Data []*vo.ThirdAccountBindListResp `json:"data"`
}
