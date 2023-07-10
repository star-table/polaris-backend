package appvo

import "github.com/star-table/polaris-backend/common/model/vo"

type AppInfoRespVo struct {
	vo.Err
	AppInfo *vo.AppInfo `json:"data"`
}

type AppInfoReqVo struct {
	AppCode string `json:"appCode"`
}

type CreateAppInfoReqVo struct {
	CreateAppInfo vo.CreateAppInfoReq `json:"input"`
	UserId        int64               `json:"userId"`
	OrgId         int64               `json:"orgId"`
}

type UpdateAppInfoReqVo struct {
	Input  vo.UpdateAppInfoReq `json:"input"`
	UserId int64               `json:"userId"`
	OrgId  int64               `json:"orgId"`
}

type DeleteAppInfoReqVo struct {
	Input  vo.DeleteAppInfoReq `json:"input"`
	UserId int64               `json:"userId"`
	OrgId  int64               `json:"orgId"`
}
