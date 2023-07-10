package domain

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
)

func CreateProjectFolder(reqBo bo.CreateFolderBo) (int64, errs.SystemErrorInfo) {

	respVo := resourcefacade.CreateFolder(resourcevo.CreateFolderReqVo{Input: reqBo})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return 0, respVo.Error()
	}

	asyn.Execute(func() {
		//新增动态
		ext := bo.TrendExtensionBo{ObjName: reqBo.Name, FolderId: respVo.Void.ID}
		trendBo := bo.ProjectTrendsBo{
			PushType:   consts.PushTypeCreateProjectFolder,
			OrgId:      reqBo.OrgId,
			ProjectId:  reqBo.ProjectId,
			OperatorId: reqBo.UserId,
			Ext:        ext,
		}
		PushProjectTrends(trendBo)
	})

	folderId := respVo.Void.ID
	//time.Sleep(15 * time.Second)
	return folderId, nil
}

func UpdateProjectFolder(updateBo bo.UpdateFolderBo) (*int64, errs.SystemErrorInfo) {
	respVo := resourcefacade.UpdateFolder(resourcevo.UpdateFolderReqVo{
		Input: updateBo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	changes := bo.TrendChangeListBo{
		OldValue: *respVo.OldValue,
		NewValue: *respVo.NewValue,
	}
	if respVo.UpdateFields[0] == "name" {
		changes.Field = "folderName"
		changes.FieldName = consts.ProjectFolderName
	} else if respVo.UpdateFields[0] == "parentId" {
		changes.Field = "folderParentId"
		changes.FieldName = consts.ProjectFolderParentId
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{ChangeList: []bo.TrendChangeListBo{
			changes,
		},
			ObjName:  json.ToJsonIgnoreError(*respVo.FolderName),
			FolderId: updateBo.FolderID,
		}
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeUpdateProjectFolder,
			OrgId:      updateBo.OrgId,
			ProjectId:  updateBo.ProjectID,
			OperatorId: updateBo.UserId,
			Ext:        ext,
		})
	})
	return &updateBo.FolderID, nil
}

func DeleteProjectFolder(deleteBo bo.DeleteFolderBo) ([]int64, errs.SystemErrorInfo) {
	respVo := resourcefacade.DeleteFolder(resourcevo.DeleteFolderReqVo{
		Input: deleteBo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}

	asyn.Execute(func() {
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeDeleteProjectFolder,
			OrgId:      deleteBo.OrgId,
			ProjectId:  deleteBo.ProjectId,
			OperatorId: deleteBo.UserId,
			NewValue:   json.ToJsonIgnoreError(respVo.DeleteFolderData.FolderNames),
			Ext:        bo.TrendExtensionBo{CommonChange: respVo.DeleteFolderData.FolderNames},
		})
	})

	folderIds := respVo.FolderIds
	return folderIds, nil
}

func GetProjectFolder(bo bo.GetFolderBo) (*vo.FolderList, errs.SystemErrorInfo) {
	respVo := resourcefacade.GetFolder(resourcevo.GetFolderReqVo{
		Input: bo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.FolderList, nil
}
