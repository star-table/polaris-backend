package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/service"
)

func (PostGreeter) CreateResourceRelation(req resourcevo.CreateResourceRelationReqVo) resourcevo.CreateResourceRelationRespVo {
	resourceIds, err := service.CreateResourceRelation(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.IssueId, req.Input.ResourceIds)
	return resourcevo.CreateResourceRelationRespVo{
		ResourceIds: resourceIds,
		Err:         vo.NewErr(err),
	}
}

func (PostGreeter) DeleteResourceRelation(req resourcevo.DeleteResourceRelationReqVo) vo.CommonRespVo {
	err := service.DeleteResourceRelation(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.IssueIds, req.Input.SourceTypes, req.Input.RecycleVersionId)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}

func (PostGreeter) UpdateResourceRelationProjectId(req resourcevo.UpdateResourceRelationProjectIdReqVo) vo.CommonRespVo {
	err := service.UpdateResourceRelationProjectId(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.IssueIds)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}

func (PostGreeter) GetResourceRelationsByProjectId(req resourcevo.GetResourceRelationsByProjectIdReqVo) resourcevo.GetResourceRelationsByProjectIdRespVo {
	resourceRelations, err := service.GetResourceRelationsByProjectId(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.SourceTypes)
	return resourcevo.GetResourceRelationsByProjectIdRespVo{
		ResourceRelations: resourceRelations,
		Err:               vo.NewErr(err),
	}
}

func (PostGreeter) DeleteAttachmentRelation(req resourcevo.DeleteAttachmentRelationReq) vo.CommonRespVo {
	err := service.UpdateAttachmentRelation(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.IssueIds,
		req.Input.ResourceIds, req.Input.RecycleVersionId, req.Input.ColumnId, req.Input.IsDeleteResource)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}

func (PostGreeter) GetResourceRelationList(req resourcevo.GetResourceRelationListReq) resourcevo.GetResourceRelationListResp {
	res, err := service.GetResourceRelationList(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.VersionId,
		req.Input.IsDelete, req.Input.SourceTypes, req.Input.ResourceIds, req.Input.IsNeedResourceType)
	return resourcevo.GetResourceRelationListResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) AddResourceRelationWithType(req resourcevo.AddResourceRelationReq) vo.CommonRespVo {
	err := service.AddResourceRelations(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.IssueId,
		req.Input.ResourceIds, req.Input.SourceType, req.Input.ColumnId)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}
