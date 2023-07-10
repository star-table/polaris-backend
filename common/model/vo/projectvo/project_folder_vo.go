package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type CreateProjectFolderReqVo struct {
	Input  vo.CreateProjectFolderReq `json:"input"`
	UserId int64                     `json:"userId"`
	OrgId  int64                     `json:"orgId"`
}

type UpdateProjectFolderReqVo struct {
	Input  vo.UpdateProjectFolderReq `json:"input"`
	UserId int64                     `json:"userId"`
	OrgId  int64                     `json:"orgId"`
}

type UpdateProjectFolderRespVo struct {
	Output *vo.UpdateProjectFolderResp `json:"input"`
	vo.Err
}

type DeleteProjectFolerReqVo struct {
	Input  vo.DeleteProjectFolderReq `json:"input"`
	UserId int64                     `json:"userId"`
	OrgId  int64                     `json:"orgId"`
}

type DeleteProjectFolerRespVo struct {
	vo.Err
	Output *vo.DeleteProjectFolderResp `json:"data"`
}

type GetProjectFolderReqVo struct {
	Input  vo.ProjectFolderReq `json:"input"`
	UserId int64               `json:"userId"`
	OrgId  int64               `json:"orgId"`
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
}

type GetProjectFolderRespVo struct {
	vo.Err
	Output *vo.FolderList `json:"data"`
}
