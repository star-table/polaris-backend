package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type AddRecycleRecordReqVo struct {
	UserId int64           `json:"userId"`
	OrgId  int64           `json:"orgId"`
	Input  AddRecycleRecordDataVo `json:"data"`
}

type AddRecycleRecordDataVo struct {
	ProjectId int64 `json:"projectId"`
	RelationIds []int64 `json:"relationIds"`
	RelationType int `json:"relationType"`
}

type AddRecycleRecordRespVo struct {
	Data int64 `json:"data"`
	vo.Err
}

type GetRecycleListReqVo struct {
	OrgId  int64           `json:"orgId"`
	UserId int64 `json:"userId"`
	Page             int              `json:"page"`
	Size             int              `json:"size"`
	Input vo.RecycleBinListReq `json:"input"`
}

type GetRecycleListRespVo struct {
	vo.Err
	Data *vo.RecycleBinList `json:"data"`
}

type CompleteDeleteReqVo struct {
	OrgId int64 `json:"orgId"`
	UserId int64 `json:"userId"`
	SourceChannel string `json:"sourceChannel"`
	Input vo.RecoverRecycleBinRecordReq `json:"input"`
}
