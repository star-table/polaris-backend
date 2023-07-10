package resourcevo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type CreateFolderReqVo struct {
	Input bo.CreateFolderBo `json:"createfolder"`
}

type UpdateFolderReqVo struct {
	Input bo.UpdateFolderBo `json:"updateFolderBo"`
}

type UpdateFolderRespVo struct {
	*UpdateFolderData
	vo.Err
}

type GetFolderReqVo struct {
	Input bo.GetFolderBo `json:"getFolderBo"`
}

type GetFolderVoListRespVo struct {
	*vo.FolderList `json:"data"`

	vo.Err
}

type UpdateFolderData struct {
	FolderId     int64
	FolderName   *string
	UpdateFields []string
	OldValue     *string
	NewValue     *string
}

type DeleteFolderReqVo struct {
	Input bo.DeleteFolderBo `json:"deleteFolder"`
}

type DeleteFolderRespVo struct {
	*DeleteFolderData
	vo.Err
}

type DeleteFolderData struct {
	FolderIds   []int64
	FolderNames []string
}

type RecoverFolderReqVo struct {
	OrgId int64 `json:"orgId"`
	UserId int64 `json:"userId"`
	Input RecoverFolderData `json:"input"`
}

type RecoverFolderData struct {
	ProjectId int64 `json:"projectId"`
	FolderId int64 `json:"folderId"`
	RecycleId int64 `json:"recycleId"`
}

type GetFolderInfoBasicReqVo struct {
	OrgId int64 `json:"orgId"`
	FolderIds []int64 `json:"folderIds"`
}

type GetFolderInfoBasicRespVo struct {
	vo.Err
	Data []vo.Folder `json:"data"`
}

type RecoverFolderRespVo struct {
	vo.Err
	Data *bo.FolderBo `json:"data"`
}

type CompleteDeleteFolderReq struct {
	OrgId int64 `json:"orgId"`
	UserId int64 `json:"userId"`
	ProjectId int64 `json:"projectId"`
	FolderId int64 `json:"folderId"`
	RecycleId int64 `json:"recycleId"`
}
