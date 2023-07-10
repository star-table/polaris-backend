package service

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/domain"
)

func CreateResourceRelation(orgId, userId int64, projectId int64, issueId int64, resourceIds []int64) ([]int64, errs.SystemErrorInfo) {
	//查询资源是否存在
	resourceInfo, err := domain.GetResourceByIds(resourceIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var needIds []int64
	resourceIdSourceTypes := make(map[int64]int, 0)
	for _, resourceBo := range resourceInfo {
		needIds = append(needIds, resourceBo.Id)
		resourceIdSourceTypes[resourceBo.Id] = resourceBo.Type
	}

	if len(needIds) == 0 {
		return needIds, nil
	}

	err = domain.InsertResourceRelation(orgId, projectId, issueId, userId, needIds, resourceIdSourceTypes)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return needIds, nil
}

func DeleteResourceRelation(orgId, userId, projectId int64, issueIds []int64, sourceTypes []int, recycleVersionId int64) errs.SystemErrorInfo {
	return domain.DeleteResourceRelation(orgId, userId, projectId, issueIds, nil, sourceTypes, recycleVersionId, "", false)
}

func UpdateResourceRelationProjectId(orgId, userId, projectId int64, issueIds []int64) errs.SystemErrorInfo {
	return domain.UpdateResourceRelationProjectId(orgId, userId, projectId, issueIds)
}

func GetResourceRelationsByProjectId(orgId, userId, projectId int64, sourceTypes []int32) ([]resourcevo.ResourceRelationVo, errs.SystemErrorInfo) {
	return domain.GetResourceRelationsByProjectId(orgId, userId, projectId, sourceTypes)
	//return domain.GetResourceRelations(orgId, projectId, 0, 2, sourceTypes, nil)
}

// 资源关联关系处理，存入回收站的versionId
func UpdateAttachmentRelation(orgId, userId, projectId int64, issueIds, resourceIds []int64, recycleVersionId int64, columnId string, isDeleteResource bool) errs.SystemErrorInfo {
	return domain.UpdateAttachmentRelation(orgId, userId, projectId, issueIds, resourceIds, recycleVersionId, columnId, isDeleteResource)
}

func GetResourceRelationList(orgId, userId, projectId int64, versionId, isDelete int, sourceType []int32, resourceIds []int64, isNeedResourceType bool) ([]resourcevo.ResourceRelationVo, errs.SystemErrorInfo) {
	return domain.GetResourceRelations(orgId, projectId, versionId, isDelete, sourceType, resourceIds, isNeedResourceType)
}

func AddResourceRelations(orgId, userId, projectId, issueId int64, resourceIds []int64, sourceType int, columnId string) errs.SystemErrorInfo {
	return domain.InsertResourceRelationWithColumnId(orgId, userId, projectId, issueId, resourceIds, sourceType, columnId)
}
