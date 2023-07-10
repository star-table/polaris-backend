package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/service"
)

func (PostGreeter) CreateFolder(reqVo resourcevo.CreateFolderReqVo) vo.CommonRespVo {
	res, err := service.CreateFolder(reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateFolder(reqVo resourcevo.UpdateFolderReqVo) resourcevo.UpdateFolderRespVo {
	res, err := service.UpdateFolder(reqVo.Input)
	return resourcevo.UpdateFolderRespVo{Err: vo.NewErr(err), UpdateFolderData: res}
}

func (PostGreeter) DeleteFolder(reqVo resourcevo.DeleteFolderReqVo) resourcevo.DeleteFolderRespVo {
	res, err := service.DeleteFolder(reqVo.Input)
	return resourcevo.DeleteFolderRespVo{Err: vo.NewErr(err), DeleteFolderData: res}
}

func (PostGreeter) CompleteDeleteFolder(req resourcevo.CompleteDeleteFolderReq) vo.CommonRespVo {
	err := service.CompleteDeleteFolder(req.OrgId, req.UserId, req.ProjectId, req.FolderId, req.RecycleId)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}

func (PostGreeter) GetFolder(reqVo resourcevo.GetFolderReqVo) resourcevo.GetFolderVoListRespVo {
	res, err := service.GetFolder(reqVo.Input)
	return resourcevo.GetFolderVoListRespVo{FolderList: res, Err: vo.NewErr(err)}
}

func (PostGreeter) RecoverFolder(req resourcevo.RecoverFolderReqVo) resourcevo.RecoverFolderRespVo {
	res, err := service.RecoverFolder(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.FolderId, req.Input.RecycleId)
	return resourcevo.RecoverFolderRespVo{
		Err:  vo.NewErr(err),
		Data:res,
	}
}

func (PostGreeter) GetFolderInfoBasic(req resourcevo.GetFolderInfoBasicReqVo)resourcevo.GetFolderInfoBasicRespVo {
	res, err := service.GetFolderInfoBasic(req.OrgId, req.FolderIds)
	return resourcevo.GetFolderInfoBasicRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
