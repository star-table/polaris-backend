package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type CreateProjectResourceReqVo struct {
	Input  vo.CreateProjectResourceReq `json:"input"`
	UserId int64                       `json:"userId"`
	OrgId  int64                       `json:"orgId"`
}

type UpdateProjectResourceNameReqVo struct {
	Input  vo.UpdateProjectResourceNameReq `json:"input"`
	UserId int64                           `json:"userId"`
	OrgId  int64                           `json:"orgId"`
}

type UpdateProjectFileResourceReqVo struct {
	Input  vo.UpdateProjectFileResourceReq `json:"input"`
	UserId int64                           `json:"userId"`
	OrgId  int64                           `json:"orgId"`
}

type UpdateProjectResourceFolderReqVo struct {
	Input  vo.UpdateProjectResourceFolderReq `json:"input"`
	UserId int64                             `json:"userId"`
	OrgId  int64                             `json:"orgId"`
}
type UpdateProjectResourceFolderRespVo struct {
	Output *vo.UpdateProjectResourceFolderResp `json:"onput"`
	vo.Err
}

type DeleteProjectResourceReqVo struct {
	Input  vo.DeleteProjectResourceReq `json:"input"`
	UserId int64                       `json:"userId"`
	OrgId  int64                       `json:"orgId"`
}

type DeleteProjectResourceRespVo struct {
	Output *vo.DeleteProjectResourceResp `json:"onput"`
	vo.Err
}

type GetProjectResourceReqVo struct {
	Input  vo.ProjectResourceReq `json:"input"`
	UserId int64                 `json:"userId"`
	OrgId  int64                 `json:"orgId"`
	Page   int                   `json:"page"`
	Size   int                   `json:"size"`
}

type GetProjectResourceResVo struct {
	vo.Err
	Output *vo.ResourceList `json:"data"`
}

type GetProjectResourceInfoReqVo struct {
	Input  vo.ProjectResourceInfoReq `json:"input"`
	UserId int64                     `json:"userId"`
	OrgId  int64                     `json:"orgId"`
}

type GetProjectResourceInfoRespVo struct {
	vo.Err
	Output *vo.Resource `json:"data"`
}
