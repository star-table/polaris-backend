package service

import (
	"fmt"

	"github.com/star-table/polaris-backend/facade/tablefacade"

	tablePb "github.com/star-table/interface/golang/table/v1"

	"github.com/star-table/polaris-backend/common/model/vo/commonvo"

	"github.com/star-table/common/core/threadlocal"

	msgPb "github.com/star-table/interface/golang/msg/v1"
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	"github.com/star-table/polaris-backend/facade/common/report"
	"github.com/star-table/polaris-backend/facade/formfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/spf13/cast"
)

// 恢复资源
func RecoverRecycleBin(orgId, userId, projectId int64, recycleId, relationId int64, relationType int, sourceChannel string) (*vo.Void, errs.SystemErrorInfo) {
	//info, err := domain.GetRecoverRecord(orgId, relationId, relationType)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	info, err := domain.GetRecycleInfo(orgId, recycleId, relationId, relationType)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//权限判断(自己删除，超管，项目管理员)
	authErr := AuthRecover(orgId, userId, projectId, info.Creator)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}

	recycleVersionId := int64(info.Version)
	allIssueIds := []int64{}
	var recoverErr errs.SystemErrorInfo
	var tableId int64
	switch relationType {
	case consts.RecycleTypeIssue:
		allIssueIds, tableId, recoverErr = domain.RecoverIssue(orgId, userId, relationId, recycleVersionId, sourceChannel)
	case consts.RecycleTypeFolder:
		recoverErr = domain.RecoverFolder(orgId, userId, projectId, relationId, recycleVersionId)
	case consts.RecycleTypeResource:
		recoverErr = domain.RecoverResource(orgId, userId, projectId, relationId, recycleVersionId)
	case consts.RecycleTypeAttachment:
		//recoverErr = domain.RecoverAttachment(orgId, userId, projectId, relationId, recycleId)
		recoverErr = domain.RecoverLessAttachments(orgId, userId, projectId, relationId, recycleVersionId, sourceChannel)
	}
	if recoverErr != nil {
		log.Error(recoverErr)
		return nil, recoverErr
	}

	//删除回收站
	deleteErr := domain.DeleteRecycleBin(info.Id, userId)
	if deleteErr != nil {
		log.Error(deleteErr)
		return nil, deleteErr
	}
	if len(allIssueIds) == 0 {
		return &vo.Void{ID: relationId}, nil
	}

	//无码恢复
	appId, appIdErr := domain.GetAppIdFromProjectId(orgId, projectId)
	if appIdErr != nil {
		log.Error(appIdErr)
		return nil, appIdErr
	}
	resp := formfacade.LessRecoverIssue(formvo.LessRecoverIssueReq{
		AppId:    appId,
		OrgId:    orgId,
		UserId:   userId,
		IssueIds: allIssueIds,
		TableId:  tableId,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}

	// 上报事件
	asyn.Execute(func() {
		list, err := domain.GetIssueInfosMapLcByIssueIds(orgId, userId, allIssueIds)
		if err != nil {
			return
		}

		openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
		openTraceIdStr := cast.ToString(openTraceId)

		for i := len(list) - 1; i >= 0; i-- {
			data := list[i]

			// 转成前端识别的模式。。
			dataId := cast.ToInt64(data[consts.BasicFieldId])
			issueId := cast.ToInt64(data[consts.BasicFieldIssueId])
			domain.AssembleDataIds(data)
			e := &commonvo.DataEvent{
				OrgId:     orgId,
				AppId:     appId,
				ProjectId: projectId,
				TableId:   tableId,
				DataId:    dataId,
				IssueId:   issueId,
				UserId:    userId,
				New:       data,
			}
			// 上报事件
			report.ReportDataEvent(msgPb.EventType_DataRecoverd, openTraceIdStr, e)
		}
	})

	return &vo.Void{ID: relationId}, nil
}

func AuthRecover(orgId, userId, projectId, deleteUserId int64) errs.SystemErrorInfo {
	if userId == deleteUserId {
		return nil
	}

	isProjectOwner, err := domain.IsProjectOwner(orgId, userId, projectId)
	if err != nil {
		log.Error(err)
		return err
	}
	if isProjectOwner {
		return nil
	}

	roleResp := orgfacade.GetUserAdminFlag(rolevo.GetUserAdminFlagReqVo{
		OrgId:  orgId,
		UserId: userId,
	})
	if roleResp.Failure() {
		log.Error(roleResp.Error())
		return roleResp.Error()
	}

	if roleResp.Data.IsAdmin {
		return nil
	}
	return errs.NoOperationPermissions
}

// 彻底删除
func CompleteDelete(orgId, userId, projectId int64, recycleId, relationId int64, relationType int) (*vo.Void, errs.SystemErrorInfo) {
	//info, err := domain.GetRecoverRecord(orgId, relationId, relationType)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	info, err := domain.GetRecycleInfo(orgId, recycleId, relationId, relationType)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//权限判断(自己删除，超管，项目管理员)
	authErr := AuthRecover(orgId, userId, projectId, info.Creator)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}

	//删除回收站
	deleteErr := domain.DeleteRecycleBin(info.Id, userId)
	if deleteErr != nil {
		log.Error(deleteErr)
		return nil, deleteErr
	}

	//由於目前普通组织只有附件，所以不用关心文件的限制，暂时先不处理文件相关的逻辑(废弃)
	//把文件也删除 2020-12-18
	if ok, _ := slice.Contain([]int{consts.RecycleTypeAttachment, consts.RecycleTypeResource}, relationType); ok {
		deleteResp := resourcefacade.CompleteDeleteResource(resourcevo.CompleteDeleteResourceReq{
			OrgId:      orgId,
			UserId:     userId,
			ResourceId: relationId,
		})
		if deleteResp.Failure() {
			log.Errorf("[CompleteDelete] deleteResource error:%v, orgId:%d, resourceId:%d", deleteResp.Error(), orgId, relationId)
			return nil, deleteResp.Error()
		}
	}

	//如果是文件夹 要删除文件夹里面的文件
	if relationType == consts.RecycleTypeFolder {
		deleteResp := resourcefacade.CompleteDeleteFolder(resourcevo.CompleteDeleteFolderReq{
			OrgId:     orgId,
			UserId:    userId,
			ProjectId: projectId,
			FolderId:  relationId,
			RecycleId: int64(info.Version),
		})
		if deleteResp.Failure() {
			log.Error(deleteResp.Error())
			return nil, deleteResp.Error()
		}
	}

	//如果是任务的话要把回收站的子任务也删除掉
	if info.RelationType == consts.RecycleTypeIssue {
		reply := tablefacade.DeleteRows(orgId, userId, &tablePb.DeleteRequest{Condition: &tablePb.Condition{Type: tablePb.ConditionType_and, Conditions: []*tablePb.Condition{
			{
				Type: tablePb.ConditionType_or, Conditions: []*tablePb.Condition{
					domain.GetRowsCondition(consts.BasicFieldPath, tablePb.ConditionType_like, fmt.Sprintf(",%d,", relationId), nil),
					domain.GetRowsCondition(consts.BasicFieldIssueId, tablePb.ConditionType_equal, relationId, nil),
				},
			},
			domain.GetRowsCondition(consts.BasicFieldOrgId, tablePb.ConditionType_equal, orgId, nil),
		}}})
		if reply.Failure() {
			return nil, reply.Error()
		}
		projectSimple, errSys := domain.GetProjectSimple(orgId, projectId)
		if errSys != nil {
			log.Errorf("[CompleteDelete] GetProjectSimple err:%v, orgId:%v, projectId:%v", errSys, orgId, projectId)
			return nil, errSys
		}
		if projectSimple.TemplateFlag == consts.TemplateFalse {
			errCache := domain.SetDeletedIssuesNum(orgId, reply.Data.Count)
			if errCache != nil {
				log.Errorf("[CompleteDelete] SetDeletedIssuesNum err:%v, orgId:%v", errCache, orgId)
			}
		}
	}

	asyn.Execute(func() {
		// 补充逻辑，把无码附件删除
		if info.RelationType == consts.RecycleTypeAttachment {
			errSys := DeleteLessAttachments(orgId, userId, projectId, relationId, info.Version)
			if errSys != nil {
				log.Error(errSys)
				return
			}
		}
	})

	return &vo.Void{ID: relationId}, nil
}

// 把无码那 附件字段的数据彻底删除
func DeleteLessAttachments(orgId, userId int64, projectId, resourceId int64, versionId int) errs.SystemErrorInfo {
	// 查询附件字段的columnId
	resourceRelationList := resourcefacade.GetResourceRelationList(resourcevo.GetResourceRelationListReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &resourcevo.GetResourceRelationList{
			ProjectId:          projectId,
			ResourceIds:        []int64{resourceId},
			IsDelete:           consts.AppIsDeleted,
			VersionId:          versionId,
			IsNeedResourceType: true,
		},
	})
	if resourceRelationList.Failure() {
		log.Errorf("[DeleteLessAttachments]GetResourceRelationList err:%v, orgId:%d, projectId:%d, resourceId:%d",
			resourceRelationList.Error(), orgId, projectId, resourceId)
		return resourceRelationList.Error()
	}

	if len(resourceRelationList.Data) < 1 {
		log.Error("[DeleteLessAttachments] resourceRelation not found")
		return nil
	}

	//if resourceRelationList.NewData[0].ResourceType == consts.FsResource {
	//	// 网络资源就不在无码删除了
	//	return nil
	//}

	columnId := resourceRelationList.Data[0].ColumnId
	issueId := resourceRelationList.Data[0].IssueId

	err := domain.DeleteLessIssueAttachmentsData(orgId, userId, projectId, issueId, []int64{resourceId}, columnId)
	if err != nil {
		log.Errorf("[DeleteLessAttachments] deleteAttachments error:%v", err)
		return err
	}
	return nil
}

func AddRecycleRecord(orgId, userId int64, projectId int64, relationIds []int64, relationType int) (int64, errs.SystemErrorInfo) {
	return domain.AddRecycleRecord(orgId, userId, projectId, relationIds, relationType)
}

func GetRecycleList(orgId, userId, projectId int64, relationType int, page, size int) (*vo.RecycleBinList, errs.SystemErrorInfo) {
	total, info, err := domain.GetRecycleListWithAuth(orgId, userId, projectId, relationType, page, size)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if len(info) == 0 {
		return &vo.RecycleBinList{Total: int64(total), List: []*vo.RecycleBin{}}, nil
	}

	//获取对应资源的信息
	issueIds := []int64{}
	folderIds := []int64{}
	resourceIds := []int64{}
	creatorIds := []int64{}
	for _, bin := range info {
		creatorIds = append(creatorIds, bin.Creator)
		switch bin.RelationType {
		case consts.RecycleTypeIssue:
			issueIds = append(issueIds, bin.RelationId)
		case consts.RecycleTypeFolder:
			folderIds = append(folderIds, bin.RelationId)
		case consts.RecycleTypeResource:
			resourceIds = append(resourceIds, bin.RelationId)
		case consts.RecycleTypeAttachment:
			resourceIds = append(resourceIds, bin.RelationId)
		}
	}

	creatorIds = slice.SliceUniqueInt64(creatorIds)
	baseUserInfos := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: creatorIds,
	})
	if baseUserInfos.Failure() {
		log.Error(baseUserInfos.Error())
		return nil, baseUserInfos.Error()
	}
	userMap := make(map[int64]*bo.BaseUserInfoBo)
	for i, user := range baseUserInfos.BaseUserInfos {
		userMap[user.UserId] = &baseUserInfos.BaseUserInfos[i]
	}

	issueMap := map[int64]*bo.IssueBo{}
	if len(issueIds) > 0 {
		issues, err := domain.GetIssueInfosLc(orgId, userId, issueIds)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for _, item := range issues {
			issueMap[item.Id] = item
		}
		//issueMap = maps.NewMap("Id", issues)
	}

	folderMap := maps.LocalMap{}
	if len(folderIds) > 0 {
		folderResp := resourcefacade.GetFolderInfoBasic(resourcevo.GetFolderInfoBasicReqVo{FolderIds: folderIds, OrgId: orgId})
		if folderResp.Failure() {
			log.Error(folderResp.Error())
			return nil, folderResp.Error()
		}

		folderMap = maps.NewMap("ID", folderResp.Data)
	}

	resourceMap := maps.LocalMap{}
	if len(resourceIds) > 0 {
		resourceResp := resourcefacade.GetResourceBoList(resourcevo.GetResourceBoListReqVo{
			Input: resourcevo.GetResourceBoListCond{
				ResourceIds: &resourceIds,
				OrgId:       orgId,
			},
		})
		if resourceResp.Failure() {
			log.Error(resourceResp.Error())
			return nil, resourceResp.Error()
		}

		resourceMap = maps.NewMap("Id", *resourceResp.ResourceBos)
	}

	res := make([]*vo.RecycleBin, 0, len(info))
	_ = copyer.Copy(info, &res)

	hasPermission := false
	isProjectOwner, err := domain.IsProjectOwner(orgId, userId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if isProjectOwner {
		hasPermission = true
	}

	roleResp := orgfacade.GetUserAdminFlag(rolevo.GetUserAdminFlagReqVo{
		OrgId:  orgId,
		UserId: userId,
	})
	if roleResp.Failure() {
		log.Error(roleResp.Error())
		return nil, roleResp.Error()
	}
	if roleResp.Data.IsAdmin {
		hasPermission = true
	}

	deleteIssueIdsMap := make(map[int64]bool)
	for i, bin := range res {
		if hasPermission || userId == bin.Creator {
			res[i].IsCanDo = true
		}
		if baseUserInfo, ok := userMap[bin.Creator]; ok {
			ownerInfo := &vo.UserIDInfo{}
			ownerInfo.UserID = baseUserInfo.UserId
			ownerInfo.Name = baseUserInfo.Name
			ownerInfo.Avatar = baseUserInfo.Avatar
			res[i].CreatorInfo = ownerInfo
		}
		res[i].TagInfo = &vo.Tag{}
		res[i].ResourceInfo = &vo.ResourceInfo{}

		switch bin.RelationType {
		case consts.RecycleTypeIssue:
			if issue, ok := issueMap[bin.RelationID]; ok {
				//issueInfo := issue.(bo.IssueBo)
				res[i].Name = issue.Title
			} else {
				deleteIssueIdsMap[bin.RelationID] = true
			}
		case consts.RecycleTypeFolder:
			if folder, ok := folderMap[bin.RelationID]; ok {
				folderInfo := folder.(vo.Folder)
				res[i].Name = folderInfo.Name
			}
		case consts.RecycleTypeResource:
			if resource, ok := resourceMap[bin.RelationID]; ok {
				resourceInfo := resource.(bo.ResourceBo)
				res[i].Name = resourceInfo.Name
				url := resourceInfo.Host + resourceInfo.Path
				createTime := types.Time(resourceInfo.CreateTime)
				temp := vo.ResourceInfo{
					URL:         &url,
					Name:        &resourceInfo.Name,
					Size:        &resourceInfo.Size,
					Suffix:      &resourceInfo.Suffix,
					UploadTime:  &createTime,
					CreatorName: &resourceInfo.CreatorName,
					Creator:     &resourceInfo.Creator,
				}
				res[i].ResourceInfo = &temp
			}
		case consts.RecycleTypeAttachment:
			if resource, ok := resourceMap[bin.RelationID]; ok {
				resourceInfo := resource.(bo.ResourceBo)
				res[i].Name = resourceInfo.Name
				url := resourceInfo.Host + resourceInfo.Path
				temp := vo.ResourceInfo{
					URL:    &url,
					Name:   &resourceInfo.Name,
					Size:   &resourceInfo.Size,
					Suffix: &resourceInfo.Suffix,
				}
				res[i].ResourceInfo = &temp
			}
		}
	}

	newRes := res
	if len(deleteIssueIdsMap) > 0 {
		total = total - uint64(len(deleteIssueIdsMap))
		newRes = make([]*vo.RecycleBin, 0, len(res))
		for _, re := range res {
			if re.RelationType != consts.RecycleTypeIssue || !deleteIssueIdsMap[re.RelationID] {
				newRes = append(newRes, re)
			}
		}
	}

	return &vo.RecycleBinList{
		Total: int64(total),
		List:  newRes,
	}, nil
}
