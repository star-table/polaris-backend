package domain

import (
	"strconv"
	"time"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
)

func CreateProjectResource(createResourceReqBo bo.ProjectCreateResourceReqBo) (int64, errs.SystemErrorInfo) {
	var resourceId int64 = 0

	bucketName := ""
	md5 := ""
	if createResourceReqBo.BucketName != nil {
		bucketName = *createResourceReqBo.BucketName
	}
	if createResourceReqBo.Md5 != nil {
		md5 = *createResourceReqBo.Md5
	}
	resourceType := createResourceReqBo.ResourceType
	sourcetype := createResourceReqBo.PolicyType
	createResourceBo := bo.CreateResourceBo{
		ProjectId:  createResourceReqBo.ProjectId,
		OrgId:      createResourceReqBo.OrgId,
		Path:       createResourceReqBo.ResourcePath,
		Name:       createResourceReqBo.FileName,
		Suffix:     createResourceReqBo.FileSuffix,
		Bucket:     bucketName,
		Type:       resourceType,
		Size:       createResourceReqBo.ResourceSize,
		Md5:        md5,
		OperatorId: createResourceReqBo.OperatorId,
		FolderId:   &createResourceReqBo.FolderId,
		SourceType: &sourcetype,
	}

	respVo := resourcefacade.CreateResource(resourcevo.CreateResourceReqVo{CreateResourceBo: createResourceBo})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return 0, respVo.Error()
	}

	asyn.Execute(func() {
		//新增动态
		resourceTrend := []bo.ResourceInfoBo{}
		suffix := createResourceReqBo.FileSuffix
		if suffix == "" {
			suffix = util.ParseFileSuffix(createResourceBo.Name)
		}
		resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
			Name:       createResourceReqBo.FileName,
			Url:        createResourceReqBo.ResourcePath,
			Size:       createResourceReqBo.ResourceSize,
			UploadTime: time.Now(),
			Suffix:     suffix,
		})
		ext := bo.TrendExtensionBo{
			ObjName:      createResourceBo.Name,
			FolderId:     createResourceReqBo.FolderId,
			ResourceInfo: resourceTrend,
		}
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeCreateProjectFile,
			OrgId:      createResourceBo.OrgId,
			ProjectId:  createResourceBo.ProjectId,
			OperatorId: createResourceBo.OperatorId,
			Ext:        ext,
		})
	})

	resourceId = respVo.ResourceId
	return resourceId, nil
}

func UpdateProjectResourceName(updateResourceBo bo.UpdateResourceInfoBo) (int64, errs.SystemErrorInfo) {
	respVo := resourcefacade.UpdateResourceInfo(resourcevo.UpdateResourceInfoReqVo{
		Input: updateResourceBo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return 0, respVo.Error()
	}
	//新增动态
	changes := bo.TrendChangeListBo{
		Field:     "resourceName",
		FieldName: consts.ProjectResourceName,
		OldValue:  respVo.OldBo[0].Name,
		NewValue:  respVo.NewBo[0].Name,
	}

	resourceTrend := []bo.ResourceInfoBo{}
	resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
		Name:       respVo.NewBo[0].Name,
		Url:        respVo.NewBo[0].Host + respVo.NewBo[0].Path,
		Size:       respVo.NewBo[0].Size,
		UploadTime: respVo.NewBo[0].UpdateTime,
		Suffix:     respVo.NewBo[0].Suffix,
	})
	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{
			ChangeList:   []bo.TrendChangeListBo{changes},
			ObjName:      json.ToJsonIgnoreError(respVo.NewBo[0].Name),
			ResourceInfo: resourceTrend,
		}
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeUpdateProjectFile,
			OrgId:      updateResourceBo.OrgId,
			ProjectId:  updateResourceBo.ProjectId,
			OperatorId: updateResourceBo.UserId,
			Ext:        ext,
		})
	})

	resourceId := updateResourceBo.ResourceId
	//time.Sleep(10 * time.Second)
	return resourceId, nil
}

func UpdateProjectFileResource(updateResourceBo bo.UpdateResourceInfoBo) (int64, errs.SystemErrorInfo) {
	log.Infof("[UpdateProjectFileResource] project_resource_domain updateResourceBo: %v", updateResourceBo)

	respVo := resourcefacade.UpdateResourceInfo(resourcevo.UpdateResourceInfoReqVo{
		Input: updateResourceBo,
	})
	if respVo.Failure() {
		//log.Error(respVo.Message)
		log.Errorf("[UpdateProjectFileResource] project_resource_domain UpdateResourceInfo Failed: %s, input: %v", strs.ObjectToString(respVo.Error()), updateResourceBo)
		return 0, respVo.Error()
	}
	//新增动态
	changes := bo.TrendChangeListBo{
		Field:     "resourceName",
		FieldName: consts.ProjectResourceName,
		OldValue:  respVo.OldBo[0].Name,
		NewValue:  respVo.NewBo[0].Name,
	}

	resourceTrend := []bo.ResourceInfoBo{}
	resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
		Name:       respVo.NewBo[0].Name,
		Url:        respVo.NewBo[0].Host + respVo.NewBo[0].Path,
		Size:       respVo.NewBo[0].Size,
		UploadTime: respVo.NewBo[0].UpdateTime,
		Suffix:     respVo.NewBo[0].Suffix,
	})
	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{
			ChangeList:   []bo.TrendChangeListBo{changes},
			ObjName:      json.ToJsonIgnoreError(respVo.NewBo[0].Name),
			ResourceInfo: resourceTrend,
		}
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeUpdateProjectFile,
			OrgId:      updateResourceBo.OrgId,
			ProjectId:  updateResourceBo.AppId,
			OperatorId: updateResourceBo.UserId,
			Ext:        ext,
		})
	})

	resourceId := updateResourceBo.ResourceId
	//time.Sleep(10 * time.Second)
	return resourceId, nil
}

func UpdateProjectResourceFolder(updateResourceBo bo.UpdateResourceFolderBo) ([]int64, errs.SystemErrorInfo) {
	respVo := resourcefacade.UpdateResourceFolder(resourcevo.UpdateResourceFolderReqVo{
		Input: updateResourceBo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	//新增动态
	changes := bo.TrendChangeListBo{
		Field:     "resourceFolder",
		FieldName: consts.ProjectResourceFolder,
		OldValue:  *respVo.CurrentFolderName,
		NewValue:  *respVo.TargetFolderName,
	}
	resourceNames := []string{}
	resourceIds := []int64{}
	for _, value := range respVo.OldBo {
		resourceNames = append(resourceNames, value.Name)
		resourceIds = append(resourceIds, value.Id)
	}

	asyn.Execute(func() {
		resourceTrend := []bo.ResourceInfoBo{}
		for _, resourceBo := range respVo.NewBo {
			resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
				Name:       resourceBo.Name,
				Url:        resourceBo.Host + resourceBo.Path,
				Size:       resourceBo.Size,
				UploadTime: resourceBo.UpdateTime,
				Suffix:     resourceBo.Suffix,
			})
		}

		ext := bo.TrendExtensionBo{
			ChangeList: []bo.TrendChangeListBo{
				changes,
			},
			ObjName:      json.ToJsonIgnoreError(resourceNames),
			ResourceInfo: resourceTrend,
		}
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeUpdateProjectFileFolder,
			OrgId:      updateResourceBo.OrgId,
			ProjectId:  updateResourceBo.ProjectId,
			OperatorId: updateResourceBo.UserId,
			OldValue:   strconv.FormatInt(updateResourceBo.CurrentFolderId, 10),
			NewValue:   strconv.FormatInt(updateResourceBo.TargetFolderID, 10),
			Ext:        ext,
		})
	})
	return resourceIds, nil
}

func DeleteProjectResource(deleteBo bo.DeleteResourceBo) ([]int64, errs.SystemErrorInfo) {
	respVo := resourcefacade.DeleteResource(resourcevo.DeleteResourceReqVo{
		Input: deleteBo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}

	resourceNames := []string{}
	resourceIds := []int64{}
	for _, value := range respVo.OldBo {
		resourceNames = append(resourceNames, value.Name)
		resourceIds = append(resourceIds, value.Id)
	}
	asyn.Execute(func() {
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeDeleteProjectFile,
			OrgId:      deleteBo.OrgId,
			ProjectId:  deleteBo.ProjectId,
			OperatorId: deleteBo.UserId,
			NewValue:   json.ToJsonIgnoreError(resourceNames),
			Ext:        bo.TrendExtensionBo{CommonChange: resourceNames},
		})
	})
	return resourceIds, nil
}

func GetProjectResource(bo bo.GetResourceBo) (*vo.ResourceList, errs.SystemErrorInfo) {
	respVo := resourcefacade.GetResource(resourcevo.GetResourceReqVo{
		Input: bo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.ResourceList, nil
}

func GetProjectResourceInfo(bo bo.GetResourceInfoBo) (*vo.Resource, errs.SystemErrorInfo) {
	respVo := resourcefacade.GetResourceInfo(resourcevo.GetResourceInfoReqVo{
		Input: bo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.Resource, nil
}
