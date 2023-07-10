package service

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

func CreateProjectResource(orgId, operatorId int64, input vo.CreateProjectResourceReq) (*vo.Void, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.OperationProFileUpload)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resourceId, err := domain.CreateProjectResource(bo.ProjectCreateResourceReqBo{
		ProjectId:    input.ProjectID,
		OrgId:        orgId,
		ResourcePath: input.ResourcePath,
		ResourceSize: input.ResourceSize,
		FileName:     input.FileName,
		FileSuffix:   input.FileSuffix,
		Md5:          input.Md5,
		BucketName:   input.BucketName,
		OperatorId:   operatorId,
		FolderId:     input.FolderID,
		PolicyType:   input.PolicyType,
		ResourceType: input.ResourceType,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.Void{
		ID: resourceId,
	}, nil
}

func UpdateProjectResourceName(orgId, operatorId int64, input vo.UpdateProjectResourceNameReq) (*vo.Void, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.OperationProFileModify)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	updateBo := &bo.UpdateResourceInfoBo{}
	err1 := copyer.Copy(&input, updateBo)
	if err1 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	updateBo.UserId = operatorId
	updateBo.OrgId = orgId
	resourceId, err := domain.UpdateProjectResourceName(*updateBo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.Void{
		ID: resourceId,
	}, nil
}

func UpdateProjectFileResource(orgId, operatorId int64, input vo.UpdateProjectFileResourceReq) (*vo.Void, errs.SystemErrorInfo) {
	log.Infof("[UpdateProjectFileResource] project_resource_service orgId: %d, operatorId: %d, input: %v", orgId, operatorId, input)

	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.AppID, consts.RoleOperationPathOrgProFile, consts.OperationProFileModify)
	if err != nil {
		log.Errorf("[UpdateProjectFileResource] project_resource_service AuthProject Failed: %s, orgId: %d, operatorId: %d, appId: %d ,path: %s , operation: %s", strs.ObjectToString(err), orgId, operatorId, input.AppID, consts.RoleOperationPathOrgProFile, consts.OperationProFileModify)
		return nil, err
	}
	updateBo := &bo.UpdateResourceInfoBo{}
	err1 := copyer.Copy(&input, updateBo)
	if err1 != nil {
		log.Errorf("[UpdateProjectFileResource] project_resource_service copyer.Copy Failed: %s, input: %v, updateBo: %v", strs.ObjectToString(err), input, updateBo)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	updateBo.UserId = operatorId
	updateBo.OrgId = orgId
	resourceId, err := domain.UpdateProjectFileResource(*updateBo)
	if err != nil {
		log.Errorf("[UpdateProjectFileResource] project_resource_service Failed: %s, updateBo: %v", strs.ObjectToString(err), updateBo)
		return nil, err
	}
	return &vo.Void{
		ID: resourceId,
	}, nil
}

func UpdateProjectResourceFolder(orgId, operatorId int64, input vo.UpdateProjectResourceFolderReq) (*vo.UpdateProjectResourceFolderResp, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.OperationProFileModify)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	bo := &bo.UpdateResourceFolderBo{}
	err1 := copyer.Copy(&input, bo)
	if err1 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	bo.UserId = operatorId
	bo.OrgId = orgId
	resourceIds, err := domain.UpdateProjectResourceFolder(*bo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.UpdateProjectResourceFolderResp{
		ResourceIds: resourceIds,
	}, nil
}

// 目前在项目文件管理中使用这个方法删除 文件、网络资源（飞书云文档）
func DeleteProjectResource(orgId, operatorId int64, input vo.DeleteProjectResourceReq) (*vo.DeleteProjectResourceResp, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.OperationProFileDelete)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// 区别文件和云文档的sourceType
	resourceIdsTypeMap, errSys := domain.GetResourceTypeByResourceIds(orgId, input.ProjectID, input.ResourceIds)
	if errSys != nil {
		log.Errorf("[DeleteProjectResource]GetResourceTypeByResourceIds err:%v", errSys)
		return nil, errSys
	}

	// 资源的sourceType 转化为 回收站资源id和类型的映射
	resourceIdRecycleType := map[int64]int{}
	for reId, sourceType := range resourceIdsTypeMap {
		recycleType := 0
		if sourceType == consts.OssPolicyTypeLesscodeResource {
			recycleType = consts.RecycleTypeAttachment
		} else if sourceType == consts.OssPolicyTypeProjectResource {
			recycleType = consts.RecycleTypeResource
		}
		resourceIdRecycleType[reId] = recycleType
	}

	recycleVersionId, err := domain.AddRecycleForProjectResource(orgId, operatorId, input.ProjectID, input.ResourceIds, resourceIdRecycleType)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resourceIds, err := domain.DeleteProjectResource(bo.DeleteResourceBo{
		ResourceIds:      input.ResourceIds,
		FolderId:         &input.FolderID,
		OrgId:            orgId,
		UserId:           operatorId,
		ProjectId:        input.ProjectID,
		RecycleVersionId: recycleVersionId,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.DeleteProjectResourceResp{
		ResourceIds: resourceIds,
	}, nil
}

func GetProjectResource(orgId, operatorId int64, page, size int, input vo.ProjectResourceReq) (*vo.ResourceList, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProjectWithOutPermission(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.OperationProFileView)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resourceList, err := domain.GetProjectResource(bo.GetResourceBo{
		FolderId:  &input.FolderID,
		OrgId:     orgId,
		UserId:    operatorId,
		ProjectId: input.ProjectID,
		Page:      page,
		Size:      size,
		//新增文件来源类型 2019/12/27
		SourceType: []int{consts.OssPolicyTypeProjectResource},
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return resourceList, nil
}

func GetProjectResourceInfo(orgId, operatorId int64, input vo.ProjectResourceInfoReq) (*vo.Resource, errs.SystemErrorInfo) {
	log.Infof("[GetProjectResourceInfo] project_resource_service GetProjectResourceInfo, orgId: %d, operatorId: %d,input: %v", orgId, operatorId, input)

	//此处的path需要更改
	err := domain.AuthProjectWithOutPermission(orgId, operatorId, input.AppID, consts.RoleOperationPathOrgProFile, consts.OperationProFileView)
	if err != nil {
		log.Errorf("[GetProjectResourceInfo] project_resource_service GetProjectResourceInfo Failed: %s, orgId: %d, operatorId: %d, appId:%d, path: %s, operation: %s", strs.ObjectToString(err), orgId, operatorId, input.AppID, consts.RoleOperationPathOrgProFile, consts.OperationProFileView)
		return nil, err
	}
	resource, err := domain.GetProjectResourceInfo(bo.GetResourceInfoBo{
		OrgId:      orgId,
		UserId:     operatorId,
		AppId:      input.AppID,
		ResourceId: input.ResourceID,
		//新增文件来源类型 2019/12/27
		SourceTypes: []int{consts.OssPolicyTypeProjectResource},
	})
	if err != nil {
		log.Errorf("[GetProjectResourceInfo] domain.GetProjectResourceInfo Failed: %s, orgId: %d, userId: %d, appId:%d, resourceId:%d, SourceType:%d", strs.ObjectToString(err), orgId, operatorId, input.AppID, input.ResourceID, consts.OssPolicyTypeProjectResource)
		return nil, err
	}

	return resource, nil
}
