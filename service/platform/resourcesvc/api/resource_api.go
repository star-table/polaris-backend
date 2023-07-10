package api

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/domain"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/service"
	"upper.io/db.v3/lib/sqlbuilder"
)

func (PostGreeter) CreateResource(req resourcevo.CreateResourceReqVo) resourcevo.CreateResourceRespVo {
	respVo := resourcevo.CreateResourceRespVo{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		resourceId, err := service.CreateResource(req.CreateResourceBo, tx)
		respVo.ResourceId = resourceId
		respVo.Err = vo.NewErr(err)
		return err
	})
	return respVo
}

func (PostGreeter) UpdateResourceInfo(reqVo resourcevo.UpdateResourceInfoReqVo) resourcevo.UpdateResourceInfoResVo {
	res, err := service.UpdateResourceInfo(reqVo.Input)
	return resourcevo.UpdateResourceInfoResVo{UpdateResourceData: res, Err: vo.NewErr(err)}
}

func (PostGreeter) UpdateResourceFolder(reqVo resourcevo.UpdateResourceFolderReqVo) resourcevo.UpdateResourceInfoResVo {
	res, err := service.UpdateResourceFolder(reqVo.Input)
	return resourcevo.UpdateResourceInfoResVo{UpdateResourceData: res, Err: vo.NewErr(err)}
}

//func (PostGreeter) UpdateResourceParentId(reqVo resourcevo.UpdateResourceParentIdReqVo) vo.CommonRespVo {
//	res, err := service.UpdateResourceParentId(reqVo.Input)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}

func (PostGreeter) DeleteResource(reqVo resourcevo.DeleteResourceReqVo) resourcevo.UpdateResourceInfoResVo {
	res, err := service.DeleteResource(reqVo.Input)
	return resourcevo.UpdateResourceInfoResVo{Err: vo.NewErr(err), UpdateResourceData: res}
}

func (PostGreeter) CompleteDeleteResource(reqVo resourcevo.CompleteDeleteResourceReq) vo.CommonRespVo {
	err := service.CompleteDeleteResource(reqVo.OrgId, reqVo.UserId, []int64{reqVo.ResourceId})
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}

func (PostGreeter) GetResource(reqVo resourcevo.GetResourceReqVo) resourcevo.GetResourceVoListRespVo {
	res, err := service.GetResource(reqVo.Input)
	return resourcevo.GetResourceVoListRespVo{ResourceList: res, Err: vo.NewErr(err)}
}

func (PostGreeter) GetResourceInfo(reqVo resourcevo.GetResourceInfoReqVo) resourcevo.GetResourceVoInfoRespVo {
	res, err := service.GetResourceInfo(reqVo.Input)
	return resourcevo.GetResourceVoInfoRespVo{Resource: res, Err: vo.NewErr(err)}
}

//func (PostGreeter) InsertResource(req resourcevo.InsertResourceReqVo) resourcevo.InsertResourceRespVo {
//	input := req.Input
//	respVo := resourcevo.InsertResourceRespVo{}
//	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//		resourceId, err := service.InsertResource(tx, input.ResourcePath, input.OrgId, input.CurrentUserId, input.ResourceType, input.FileName)
//		respVo.ResourceId = resourceId
//		respVo.Err = vo.NewErr(err)
//		return err
//	})
//	return respVo
//}

func (PostGreeter) GetResourceById(req resourcevo.GetResourceByIdReqVo) resourcevo.GetResourceByIdRespVo {
	resourceBos, err := domain.GetResourceByIds(req.GetResourceByIdReqBody.ResourceIds)
	return resourcevo.GetResourceByIdRespVo{ResourceBos: resourceBos, Err: vo.NewErr(err)}
}

func (GetGreeter) GetIdByPath(req resourcevo.GetIdByPathReqVo) resourcevo.GetIdByPathRespVo {
	resourceId, err := service.GetIdByPath(req.OrgId, req.ResourcePath, req.ResourceType)
	return resourcevo.GetIdByPathRespVo{ResourceId: resourceId, Err: vo.NewErr(err)}
}

func (PostGreeter) GetResourceBoList(req resourcevo.GetResourceBoListReqVo) resourcevo.GetResourceBoListRespVo {
	list, total, err := domain.GetResourceBoList(req.Page, req.Size, req.Input)
	return resourcevo.GetResourceBoListRespVo{GetResourceBoListRespData: resourcevo.GetResourceBoListRespData{ResourceBos: list, Total: total}, Err: vo.NewErr(err)}
}

func (PostGreeter) GetIssueIdsByResourceIds(req resourcevo.GetIssueIdsByResourceIdsReqVo) resourcevo.GetIssueIdsByResourceIdsRespVo {
	resp, err := service.GetIssueIdsByResource(req.OrgId, req.Input.ProjectId, req.Input.ResourceIds)
	return resourcevo.GetIssueIdsByResourceIdsRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (PostGreeter) RecoverResource(req resourcevo.RecoverResourceReqVo) resourcevo.RecoverResourceRespVo {
	res, err := service.RecoverResource(req.OrgId, req.Input.AppId, req.UserId, req.Input.ProjectId, req.Input.ResourceId, req.Input.RecycleVersionId, req.Input.IssueIds, req.SourceChannel)
	return resourcevo.RecoverResourceRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) CacheResourceSize(req resourcevo.CacheResourceSizeReq) resourcevo.CacheResourceSizeResp {
	size, err := service.CacheResourceSize(req.OrgId)
	return resourcevo.CacheResourceSizeResp{
		Err:  vo.NewErr(err),
		Size: size,
	}
}

func (PostGreeter) FsDocumentList(req resourcevo.FsDocumentListReq) resourcevo.FsDocumentListResp {
	res, err := service.FsDocumentList(req.OrgId, req.UserId, req.Page, req.Size, req.Params.SearchKey)
	return resourcevo.FsDocumentListResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) AddFsResourceBatch(req resourcevo.AddFsResourceBatchReq) resourcevo.AddFsResourceBatchResp {
	res, err := service.AddFsResourceBatch(req.OrgId, req.UserId, req.ProjectId, req.IssueId, req.FolderId, req.Input)
	return resourcevo.AddFsResourceBatchResp{
		Err: vo.NewErr(err),
		Ids: res,
	}
}

// 获取钉盘团队文件下 所有层级的文件夹、文件列表信息
func (PostGreeter) DingFileList(req resourcevo.DingFileListReq) resourcevo.DingFileListRespVo {
	res, err := service.DingFileList(req)
	return resourcevo.DingFileListRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) DingDocumentList(req resourcevo.DingDocReq) resourcevo.DingDocumentRespVo {
	res, err := service.DingDocumentList(req.OrgId, req.UserId, req.Page, req.Size)
	return resourcevo.DingDocumentRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

// 获取钉钉空间列表
func (PostGreeter) DingSpaceList(req resourcevo.DingSpaceReqVo) resourcevo.DingSpaceListResp {
	res, err := service.DingSpaceList(req)
	return resourcevo.DingSpaceListResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

// 通过spaceId、dirId获取到文件(夹)列表信息
func (PostGreeter) DingFileListById(req resourcevo.DingSpaceFileReqVo) resourcevo.DingSpaceFileResp {
	res, err := service.DingFileListById(req)
	return resourcevo.DingSpaceFileResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
