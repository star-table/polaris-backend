package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type ProjectDetailReqVo struct {
	ProjectId int64 `json:"projectId"`
	OrgId     int64 `json:"orgId"`
}

type ProjectDetailRespVo struct {
	vo.Err
	ProjectDetail *vo.ProjectDetail `json:"data"`
}

type CreateProjectDetailReqVo struct {
	Input  vo.CreateProjectDetailReq `json:"input"`
	UserId int64                     `json:"userId"`
}

type UpdateProjectDetailReqVo struct {
	Input  vo.UpdateProjectDetailReq `json:"input"`
	UserId int64                     `json:"userId"`
	OrgId  int64                     `json:"orgId"`
}

type DeleteProjectDetailReqVo struct {
	Input  vo.DeleteProjectDetailReq `json:"input"`
	UserId int64                     `json:"userId"`
}
