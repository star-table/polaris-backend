package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) AddRecycleRecord(req projectvo.AddRecycleRecordReqVo) projectvo.AddRecycleRecordRespVo {
	res, err := service.AddRecycleRecord(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.RelationIds, req.Input.RelationType)
	return projectvo.AddRecycleRecordRespVo{Data: res, Err: vo.NewErr(err)}
}

func (PostGreeter) GetRecycleList(req projectvo.GetRecycleListReqVo) projectvo.GetRecycleListRespVo {
	res, err := service.GetRecycleList(req.OrgId, req.UserId, req.Input.ProjectID, req.Input.RelationType, req.Page, req.Size)
	return projectvo.GetRecycleListRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) CompleteDelete(req projectvo.CompleteDeleteReqVo) vo.CommonRespVo {
	res, err := service.CompleteDelete(req.OrgId, req.UserId, req.Input.ProjectID, req.Input.RecycleID, req.Input.RelationID, req.Input.RelationType)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
}

func (PostGreeter) RecoverRecycleBin(req projectvo.CompleteDeleteReqVo) vo.CommonRespVo {
	res, err := service.RecoverRecycleBin(req.OrgId, req.UserId, req.Input.ProjectID, req.Input.RecycleID, req.Input.RelationID, req.Input.RelationType, req.SourceChannel)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
}
