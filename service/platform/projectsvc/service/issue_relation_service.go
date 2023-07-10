package service

import (
	"time"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func CreateIssueComment(orgId, currentUserId int64, input vo.CreateIssueCommentReq) (*vo.Void, errs.SystemErrorInfo) {
	issueId := input.IssueID

	// 如果任务有附件传过来，限制一下数量
	if input.AttachmentIds != nil && len(input.AttachmentIds) > 9 {
		return nil, errs.IssueCommentImageLimitsError
	}

	issueBo, err1 := domain.GetIssueInfoLc(orgId, 0, issueId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	err := domain.AuthIssue(orgId, currentUserId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Comment)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	return CreateIssueCommentWithoutAuth(orgId, currentUserId, input, issueBo)
}

//func OpenCreateIssueComment(orgId, currentUserId int64, input vo.CreateIssueCommentReq) (*vo.Void, errs.SystemErrorInfo) {
//	issueId := input.IssueID
//
//	issueBo, err1 := domain.GetIssueBo(orgId, issueId)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
//	}
//
//	return CreateIssueCommentWithoutAuth(orgId, currentUserId, input, issueBo)
//}

func CreateIssueCommentWithoutAuth(orgId, currentUserId int64, input vo.CreateIssueCommentReq, issueBo *bo.IssueBo) (*vo.Void, errs.SystemErrorInfo) {
	commentId, err2 := domain.CreateIssueComment(*issueBo, input.Comment, input.MentionedUserIds, currentUserId, input.AttachmentIds)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
	}

	return &vo.Void{
		ID: commentId,
	}, nil
}

func CreateIssueResource(orgId, userId int64, input *vo.CreateIssueResourceReq) (*vo.Void, errs.SystemErrorInfo) {
	resourceId, err := domain.CreateIssueResource(orgId, userId, input)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}
	return &vo.Void{ID: resourceId}, nil
}

func AddIssueAttachmentFs(orgId, operatorId int64, input vo.AddIssueAttachmentFsReq) ([]bo.ResourceBo, errs.SystemErrorInfo) {
	issueId := input.IssueID
	projectId := input.ProjectID
	folderId := input.FolderID

	issueBo := &bo.IssueBo{}
	if issueId > 0 {
		issueBos, err := domain.GetIssueInfosLc(orgId, operatorId, []int64{issueId})
		if err != nil {
			log.Errorf("[AddIssueAttachmentFs] GetIssueInfosLc err:%v, orgId:%d, operatorId:%d, issueId:%d",
				err, orgId, operatorId, issueId)
			return nil, err
		}
		if len(issueBos) < 1 {
			return nil, nil
		}
		issueBo = issueBos[0]
	}

	//issueBo, err := domain.GetIssueBo(orgId, issueId)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.IllegalityIssue)
	//}

	if len(input.Data) == 0 {
		return []bo.ResourceBo{}, nil
	}
	if projectId == 0 {
		projectId = issueBo.ProjectId
	}
	resourceResp := resourcefacade.AddFsResourceBatch(resourcevo.AddFsResourceBatchReq{
		OrgId:     orgId,
		UserId:    operatorId,
		ProjectId: projectId,
		IssueId:   issueId,
		FolderId:  folderId,
		Input:     input.Data,
	})
	if resourceResp.Failure() {
		log.Error(resourceResp.Error())
		return nil, resourceResp.Error()
	}
	//if issueId > 0 {
	//	_, err2 := domain.UpdateIssueRelation(operatorId, *issueBo, consts.IssueRelationTypeResource, resourceResp.Ids, consts.IssueAttachmentReferRelationCode)
	//	if err2 != nil {
	//		log.Error(err2)
	//		return nil, err2
	//	}
	//}

	asyn.Execute(func() {
		resourceTrend := []bo.ResourceInfoBo{}
		for _, createResourceReqBo := range input.Data {
			resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
				Name:       createResourceReqBo.Title,
				Url:        createResourceReqBo.URL,
				Size:       0,
				UploadTime: time.Now(),
				Suffix:     util.ParseFileSuffix(createResourceReqBo.URL),
			})
		}
		if issueId > 0 {
			issueTrendsBo := &bo.IssueTrendsBo{
				PushType:      consts.PushTypeReferResource,
				OrgId:         issueBo.OrgId,
				OperatorId:    operatorId,
				DataId:        issueBo.DataId,
				IssueId:       issueBo.Id,
				ParentIssueId: issueBo.ParentId,
				ProjectId:     issueBo.ProjectId,
				TableId:       issueBo.TableId,
				PriorityId:    issueBo.PriorityId,
				ParentId:      issueBo.ParentId,
				IssueTitle:    issueBo.Title,
				IssueStatusId: issueBo.Status,
				SourceChannel: "", //这里加了只是为了飞书群聊用到
				Ext: bo.TrendExtensionBo{
					ObjName:      issueBo.Title,
					ResourceInfo: resourceTrend,
				},
			}
			//asyn.Execute(func() {
			//	domain.PushIssueTrends(issueTrendsBo)
			//})

			domain.PushIssueThirdPlatformNotice(issueTrendsBo)
		}

	})

	// asyn.Execute(func() {
	// 	attachmentList := []string{}
	// 	for _, createResourceReqBo := range input.NewData {
	// 		attachmentList = append(attachmentList, createResourceReqBo.URL)
	// 	}
	// 	domain.PushMessageToFeishuShortcut(bo.ShortcutPushBo{
	// 		TriggerType:         consts.FsTriggerUpdateIssue,
	// 		EventType:           []string{consts.FsEventAddIssueAttachment},
	// 		OrgId:               orgId,
	// 		ProjectId:           issueBo.ProjectId,
	// 		ProjectObjectTypeId: issueBo.ProjectObjectTypeId,
	// 		IssueId:             issueBo.Id,
	// 		Operator:            operatorId,
	// 		Attachment:          attachmentList,
	// 	})
	// })

	isDelete := consts.AppIsNoDelete
	resourceBoListResp := resourcefacade.GetResourceBoList(resourcevo.GetResourceBoListReqVo{
		Page: 0,
		Size: 0,
		Input: resourcevo.GetResourceBoListCond{
			OrgId:       orgId,
			ResourceIds: &resourceResp.Ids,
			IsDelete:    &isDelete,
		},
	})
	if resourceBoListResp.Failure() {
		log.Error(resourceBoListResp.Error())
		return nil, resourceBoListResp.Error()
	}

	return *resourceBoListResp.ResourceBos, nil
}

func AddIssueAttachment(orgId, operatorId int64, input vo.AddIssueAttachmentReq) (*vo.Void, errs.SystemErrorInfo) {
	issueId := input.IssueID

	issueBo, err := domain.GetIssueInfoLc(orgId, 0, issueId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIssue)
	}

	if len(input.ResourceIds) == 0 {
		return &vo.Void{ID: 0}, nil
	}

	//先插入资源关联
	resourceRelationResp := resourcefacade.CreateResourceRelation(resourcevo.CreateResourceRelationReqVo{
		Input: resourcevo.CreateResourceRelationData{
			ProjectId:   issueBo.ProjectId,
			IssueId:     issueId,
			ResourceIds: input.ResourceIds,
		},
		UserId: issueBo.Creator,
		OrgId:  issueBo.OrgId,
	})
	if resourceRelationResp.Failure() {
		log.Error(resourceRelationResp.Error())
		return nil, resourceRelationResp.Error()
	}
	//再插入任务关联
	if len(resourceRelationResp.ResourceIds) > 0 {
		_, err2 := domain.UpdateIssueRelation(operatorId, issueBo, consts.IssueRelationTypeResource, resourceRelationResp.ResourceIds, consts.IssueAttachmentReferRelationCode)
		if err2 != nil {
			log.Error(err2)
			return nil, err2
		}
	}

	asyn.Execute(func() {
		resourceResp := resourcefacade.GetResourceById(resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: resourcevo.GetResourceByIdReqBody{ResourceIds: resourceRelationResp.ResourceIds}})
		if resourceResp.Failure() {
			log.Error(resourceResp.Error())
			return
		}

		attachmentList := []string{}
		for _, createResourceReqBo := range resourceResp.ResourceBos {
			resourceTrend := []bo.ResourceInfoBo{}
			attachmentList = append(attachmentList, createResourceReqBo.Host+createResourceReqBo.Path)
			resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
				Name:       createResourceReqBo.Name,
				Url:        createResourceReqBo.Host + createResourceReqBo.Path,
				Size:       createResourceReqBo.Size,
				UploadTime: time.Now(),
				Suffix:     util.ParseFileSuffix(createResourceReqBo.Suffix),
			})
			issueTrendsBo := &bo.IssueTrendsBo{
				PushType:      consts.PushTypeReferResource,
				OrgId:         issueBo.OrgId,
				OperatorId:    operatorId,
				DataId:        issueBo.DataId,
				IssueId:       issueBo.Id,
				ParentIssueId: issueBo.ParentId,
				ProjectId:     issueBo.ProjectId,
				TableId:       issueBo.TableId,
				PriorityId:    issueBo.PriorityId,
				ParentId:      issueBo.ParentId,
				IssueTitle:    issueBo.Title,
				IssueStatusId: issueBo.Status,
				Ext: bo.TrendExtensionBo{
					ObjName:      issueBo.Title,
					ResourceInfo: resourceTrend,
				},
			}
			//asyn.Execute(func() {
			//	domain.PushIssueTrends(issueTrendsBo)
			//})
			asyn.Execute(func() {
				domain.PushIssueThirdPlatformNotice(issueTrendsBo)
			})

			// asyn.Execute(func() {
			// 	domain.PushMessageToFeishuShortcut(bo.ShortcutPushBo{
			// 		TriggerType:         consts.FsTriggerUpdateIssue,
			// 		EventType:           []string{consts.FsEventAddIssueAttachment},
			// 		OrgId:               orgId,
			// 		ProjectId:           issueBo.ProjectId,
			// 		ProjectObjectTypeId: issueBo.ProjectObjectTypeId,
			// 		IssueId:             issueBo.Id,
			// 		Operator:            operatorId,
			// 		Attachment:          attachmentList,
			// 	})
			// })
		}
	})

	return &vo.Void{
		ID: int64(len(resourceRelationResp.ResourceIds)),
	}, nil
}

//func CreateIssueRelationIssue(orgId, operatorId int64, input vo.UpdateIssueAndIssueRelateReq) (*vo.Void, errs.SystemErrorInfo) {
//	issueId := input.IssueID
//
//	issueBo, err := domain.GetIssueBo(orgId, issueId)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIssue)
//	}
//
//	err = domain.AuthIssue(orgId, operatorId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
//	}
//
//	//真正增加的关联问题id
//	realAddRelationIds := []int64{}
//	if len(input.AddRelateIssueIds) > 0 {
//		relationInfo, err1 := domain.UpdateIssueRelation(operatorId, *issueBo, consts.IssueRelationTypeIssue, input.AddRelateIssueIds, "")
//		if err1 != nil {
//			log.Error(err1)
//			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
//		}
//
//		for _, v := range relationInfo {
//			realAddRelationIds = append(realAddRelationIds, v.RelationId)
//		}
//	}
//
//	if len(input.DelRelateIssueIds) > 0 {
//		err1 := domain.DeleteIssueRelationByIds(operatorId, *issueBo, consts.IssueRelationTypeIssue, input.DelRelateIssueIds)
//		if err1 != nil {
//			log.Error(err1)
//			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
//		}
//	}
//
//	asyn.Execute(func() {
//		ext := bo.TrendExtensionBo{}
//		ext.IssueType = "T"
//		ext.ObjName = issueBo.Title
//
//		issueTrendsBo := bo.IssueTrendsBo{
//			PushType:      consts.PushTypeUpdateRelationIssue,
//			OrgId:         orgId,
//			UserId:    operatorId,
//			IssueId:       issueBo.Id,
//			ParentIssueId: issueBo.ParentId,
//			ProjectId:     issueBo.ProjectId,
//			IssueTitle:    issueBo.Title,
//			IssueStatusId: issueBo.Status,
//			//BeforeOwner:   issueBo.Owner,
//			ParentId: issueBo.ParentId,
//
//			BindIssues:   realAddRelationIds,
//			UnbindIssues: input.DelRelateIssueIds,
//			Ext:          ext,
//			TableId:      issueBo.TableId,
//		}
//
//		asyn.Execute(func() {
//			domain.PushIssueTrends(issueTrendsBo)
//		})
//		asyn.Execute(func() {
//			domain.PushIssueThirdPlatformNotice(issueTrendsBo)
//		})
//	})
//
//	return &vo.Void{
//		ID: input.IssueID,
//	}, nil
//}

func DeleteIssueResource(orgId, operatorId int64, input vo.DeleteIssueResourceReq) (*vo.Void, errs.SystemErrorInfo) {
	issueId := input.IssueID

	issueBo, err := domain.GetIssueInfoLc(orgId, 0, issueId)
	if err != nil {
		log.Errorf("[DeleteIssueResource] GetIssueInfosLc err:%v, orgId:%v, issueId:%v", err, orgId, issueId)
		return nil, err
	}

	//projectInfo, err := domain.LoadProjectAuthBo(orgId, issueBo.ProjectId)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}

	//去重
	input.DeletedResourceIds = slice.SliceUniqueInt64(input.DeletedResourceIds)

	//这里权限写死:已上传文件，可删除其文件，拥有删除权限的人员为，任务负责人、附件上传者，项目负责人，超级管理员
	hasPermission := true

	//err = AuthIssue(orgId, operatorId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationUnbind)
	if !hasPermission {
		log.Error(err)
		//小逻辑：允许创建者删除文件资源
		relationIds, err := domain.GetIssueResourceIdsByCreator(orgId, issueId, input.DeletedResourceIds, operatorId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		//判断要删除的资源是否是这个人创建的
		if relationIds != nil && len(*relationIds) == len(input.DeletedResourceIds) {
			hasPermission = true
		}
	}

	if !hasPermission {
		return nil, errs.BuildSystemErrorInfo(errs.NoOperationPermissions)
	}

	if issueBo.ProjectId != int64(0) {
		judgeErr := JudgeProjectFiling(orgId, issueBo.ProjectId)
		if judgeErr != nil {
			log.Error(judgeErr)
			return nil, judgeErr
		}
	}

	err1 := domain.DeleteIssueResource(*issueBo, input.DeletedResourceIds, operatorId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	asyn.Execute(func() {
		bos, total, err1 := resourcefacade.GetResourceBoListRelaxed(0, 0, resourcevo.GetResourceBoListCond{
			ResourceIds: &input.DeletedResourceIds,
			OrgId:       orgId,
		})
		if err1 != nil {
			log.Error(err1)
			return
		}
		if total == 0 {
			log.Error("无更新")
			return
		}

		relationPos := &[]po.PpmPriIssueRelation{}
		err := mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcRelationId:   db.In(input.DeletedResourceIds),
			consts.TcRelationType: consts.IssueRelationTypeResource,
			consts.TcIssueId:      issueBo.Id,
		}, relationPos)
		if err != nil {
			log.Error(err)
			return
		}
		//引用的附件id
		referIds := []int64{}
		normalIds := []int64{}
		for _, relation := range *relationPos {
			if relation.RelationCode == consts.IssueAttachmentReferRelationCode {
				referIds = append(referIds, relation.RelationId)
			} else {
				normalIds = append(normalIds, relation.RelationId)
			}
		}
		referIds = slice.SliceUniqueInt64(referIds)
		normalIds = slice.SliceUniqueInt64(normalIds)

		if len(referIds) > 0 {
			handleDeleteResourceTrend(*issueBo, operatorId, *bos, consts.PushTypeDeleteReferResource, referIds)
		}
		if len(normalIds) > 0 {
			handleDeleteResourceTrend(*issueBo, operatorId, *bos, consts.PushTypeDeleteResource, normalIds)
		}
	})

	return &vo.Void{
		ID: issueId,
	}, nil
}

func handleDeleteResourceTrend(issueBo bo.IssueBo, operatorId int64, resources []bo.ResourceBo, pushType consts.IssueNoticePushType, resourceIds []int64) {
	resourceTrend := []bo.ResourceInfoBo{}
	for _, resourceBo := range resources {
		if ok, _ := slice.Contain(resourceIds, resourceBo.Id); !ok {
			continue
		}
		resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
			Name:       resourceBo.Name,
			Url:        resourceBo.Host + resourceBo.Path,
			Size:       resourceBo.Size,
			UploadTime: time.Now(),
			Suffix:     resourceBo.Suffix,
		})
	}

	issueTrendsBo := &bo.IssueTrendsBo{
		PushType:      pushType,
		OrgId:         issueBo.OrgId,
		OperatorId:    operatorId,
		DataId:        issueBo.DataId,
		IssueId:       issueBo.Id,
		ParentIssueId: issueBo.ParentId,
		ProjectId:     issueBo.ProjectId,
		TableId:       issueBo.TableId,
		PriorityId:    issueBo.PriorityId,
		ParentId:      issueBo.ParentId,
		IssueTitle:    issueBo.Title,
		IssueStatusId: issueBo.Status,
		Ext: bo.TrendExtensionBo{
			ObjName:      issueBo.Title,
			ResourceInfo: resourceTrend,
		},
	}
	asyn.Execute(func() {
		domain.PushIssueTrends(issueTrendsBo)
	})
	asyn.Execute(func() {
		domain.PushIssueThirdPlatformNotice(issueTrendsBo)
	})
}

func IssueResources(orgId int64, page uint, size uint, input *vo.GetIssueResourcesReq) (*vo.ResourceList, errs.SystemErrorInfo) {
	issueId := input.IssueID

	//_, err := domain.GetIssueBo(orgId, issueId)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.IllegalityIssue)
	//}

	resourceIds, err1 := domain.GetIssueRelationIdsByRelateType(orgId, issueId, consts.IssueRelationTypeResource)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError)
	}

	isDelete := consts.AppIsNoDelete
	bos, total, err1 := resourcefacade.GetResourceBoListRelaxed(page, size, resourcevo.GetResourceBoListCond{
		ResourceIds: resourceIds,
		OrgId:       orgId,
		IsDelete:    &isDelete,
	})
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ResourceDomainError, err1)
	}
	if bos != nil && len(*bos) > 0 {
		for _, resBo := range *bos {
			resBo.Path = resBo.Host + resBo.Path
		}
	}

	resultList := &[]*vo.Resource{}
	copyErr := copyer.Copy(bos, resultList)

	creatorIds := make([]int64, 0)
	for _, res := range *resultList {
		creatorIds = append(creatorIds, res.Creator)
	}

	creatorInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, creatorIds)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, err
	}

	creatorMap := maps.NewMap("UserId", creatorInfos)

	for _, res := range *resultList {
		if userInfo, ok := creatorMap[res.Creator]; ok {
			res.CreatorName = userInfo.(bo.BaseUserInfoBo).Name
		}
		//拼接host
		res.Path = util.JointUrl(res.Host, res.Path)
		res.PathCompressed = util.GetCompressedPath(res.Path, res.Type)
	}

	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &vo.ResourceList{
		Total: total,
		List:  *resultList,
	}, nil

}

func GetIssueMembers(orgId int64, issueId int64) (*projectvo.GetIssueMembersRespData, errs.SystemErrorInfo) {
	issueMembersBo, err := domain.GetIssueMembers(orgId, issueId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	respVo := &projectvo.GetIssueMembersRespData{}
	_ = copyer.Copy(issueMembersBo, respVo)
	return respVo, nil
}

func GetIssueRelationResource(page, size int) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	return domain.GetIssueRelationResource(page, size)
}

//func UpdateIssueBeforeAfterIssues(orgId, userId int64, input vo.UpdateIssueBeforeAfterIssuesReq) (*vo.Void, errs.SystemErrorInfo) {
//	issueId := input.IssueID
//
//	issueBo, err := domain.GetIssueBo(orgId, issueId)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIssue)
//	}
//
//	err = domain.AuthIssue(orgId, userId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
//	}
//
//	//真正增加的关联问题id
//	realAddRelationIds := []int64{}
//	if len(input.AddRelateIssueIds) > 0 {
//		addIds, err1 := domain.AddBeforeAfterIssue(orgId, userId, issueId, issueBo.ProjectId, input.Type, input.AddRelateIssueIds)
//		if err1 != nil {
//			log.Error(err1)
//			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
//		}
//
//		realAddRelationIds = addIds
//	}
//
//	if len(input.DelRelateIssueIds) > 0 {
//		err1 := domain.DeleteBeforeAfterIssue(orgId, userId, issueId, input.Type, input.DelRelateIssueIds)
//		if err1 != nil {
//			log.Error(err1)
//			return nil, err1
//		}
//	}
//
//	asyn.Execute(func() {
//		ext := bo.TrendExtensionBo{}
//		ext.IssueType = "T"
//		ext.ObjName = issueBo.Title
//
//		pushType := consts.PushTypeUpdateAfterIssue
//		if input.Type == 1 {
//			pushType = consts.PushTypeUpdateBeforeIssue
//		}
//		issueTrendsBo := bo.IssueTrendsBo{
//			PushType:      pushType,
//			OrgId:         orgId,
//			UserId:    userId,
//			IssueId:       issueBo.Id,
//			ParentIssueId: issueBo.ParentId,
//			ProjectId:     issueBo.ProjectId,
//			IssueTitle:    issueBo.Title,
//			IssueStatusId: issueBo.Status,
//			//BeforeOwner:   issueBo.Owner,
//			ParentId:     issueBo.ParentId,
//			BindIssues:   realAddRelationIds,
//			UnbindIssues: input.DelRelateIssueIds,
//			Ext:          ext,
//			TableId:      issueBo.TableId,
//		}
//
//		asyn.Execute(func() {
//			domain.PushIssueTrends(issueTrendsBo)
//		})
//		//asyn.Execute(func() {
//		//	domain.PushIssueThirdPlatformNotice(issueTrendsBo)
//		//})
//	})
//
//	return &vo.Void{
//		ID: input.IssueID,
//	}, nil
//}

func IssueListSimpleByDataIds(orgId, userId, appId int64, dataIds []int64) ([]*projectvo.IssueSimpleInfo, errs.SystemErrorInfo) {
	issueList := []*projectvo.IssueSimpleInfo{}
	// 构造无码查询条件
	//values := []interface{}{}
	//lessReq := vo.LessCondsData{Type: consts.ConditionAnd}
	//lessReq.Conds = append(lessReq.Conds, &vo.LessCondsData{
	//	Type:   "equal",
	//	Value:  2,
	//	Column: "recycleFlag",
	//})
	//lessReq.Conds = append(lessReq.Conds, &vo.LessCondsData{
	//	Type:   "equal",
	//	Value:  2,
	//	Column: "delFlag",
	//})
	//for _, dataId := range dataIds {
	//	values = append(values, dataId)
	//}
	//
	//lessReq.Conds = append(lessReq.Conds, &vo.LessCondsData{
	//	Type:   consts.ConditionIn,
	//	Column: consts.BasicFieldId,
	//	Values: values,
	//})
	//filterColumns := []string{consts.BasicFieldId}
	//for _, column := range getSimpleIssueFilterColumns() {
	//	if column == consts.BasicFieldTitle {
	//		filterColumns = append(filterColumns, fmt.Sprintf("\"data\" :: jsonb ->> '%s' \"%s\"", column, column))
	//	} else {
	//		filterColumns = append(filterColumns, fmt.Sprintf("\"data\" :: jsonb -> '%s' \"%s\"", column, column))
	//	}
	//}
	//page, size := 1, 10000
	//issueListResp := formfacade.LessIssueList(formvo.LessIssueListReq{
	//	Condition:     lessReq,
	//	AppId:         appId,
	//	OrgId:         orgId,
	//	UserId:        userId,
	//	FilterColumns: filterColumns,
	//	Page:          int64(page),
	//	Size:          int64(size),
	//})
	//
	//if issueListResp.Failure() {
	//	log.Errorf("[IssueListSimpleByDataIds]formfacade.LessIssueList failed, orgId:%v, userId:%v, appId:%v,dataIds:%v ,err:%v",
	//		orgId, userId, appId, dataIds, issueListResp.Message,
	//	)
	//}
	//
	//err := copyer.Copy(issueListResp.Data.List, &issueList)
	//if err != nil {
	//	log.Errorf("[IssueListSimpleByDataIds] json转换错误")
	//	return nil, errs.JSONConvertError
	//}

	return issueList, nil
}

func IssueListSimpleByTableIds(orgId, userId int64, req projectvo.IssueListWithConditionsReq) ([]*projectvo.IssueSimpleInfo, errs.SystemErrorInfo) {
	//condition := vo.LessCondsData{Type: consts.ConditionAnd}
	//condition.Conds = append(condition.Conds, &vo.LessCondsData{
	//	Type:   "equal",
	//	Value:  2,
	//	Column: "recycleFlag",
	//})
	//condition.Conds = append(condition.Conds, &vo.LessCondsData{
	//	Type:   "equal",
	//	Value:  2,
	//	Column: "delFlag",
	//})
	//values := []interface{}{}
	//for _, tableId := range req.TableIds {
	//	values = append(values, tableId)
	//}
	//if len(req.TableIds) > 0 {
	//	condition.Conds = append(condition.Conds, &vo.LessCondsData{
	//		Type:   consts.ConditionIn,
	//		Column: consts.BasicFieldTableId,
	//		Values: values,
	//	})
	//}
	//issueList := []*projectvo.IssueSimpleInfo{}
	//filterColumns := []string{consts.BasicFieldId}
	//for _, column := range getSimpleIssueFilterColumns() {
	//	if column == consts.BasicFieldTitle {
	//		filterColumns = append(filterColumns, fmt.Sprintf("\"data\" :: jsonb ->> '%s' \"%s\"", column, column))
	//	} else {
	//		filterColumns = append(filterColumns, fmt.Sprintf("\"data\" :: jsonb -> '%s' \"%s\"", column, column))
	//	}
	//}
	//issueListResp := formfacade.LessIssueList(formvo.LessIssueListReq{
	//	Condition:     condition,
	//	AppId:         req.AppId,
	//	OrgId:         orgId,
	//	UserId:        userId,
	//	FilterColumns: filterColumns,
	//	Page:          int64(req.Page),
	//	Size:          int64(req.Size),
	//})
	//if issueListResp.Failure() {
	//	log.Errorf("[IssueListSimpleByTableIds]formfacade.LessIssueList failed, orgId:%v, userId:%v, appId:%v ,err:%v",
	//		orgId, userId, req.AppId, issueListResp.Message,
	//	)
	//}
	//err := copyer.Copy(issueListResp.Data.List, &issueList)
	//if err != nil {
	//	log.Errorf("[IssueListSimpleByTableIds] json转换错误")
	//	return nil, errs.JSONConvertError
	//}

	return nil, nil
}

func getSimpleIssueFilterColumns() []string {
	filterColumns := []string{consts.BasicFieldIssueId,
		consts.BasicFieldParentId, consts.BasicFieldProjectId,
		consts.BasicFieldTitle, consts.BasicFieldTableId, consts.BasicFieldIssueStatusType,
		consts.BasicFieldOwnerId, consts.BasicFieldFollowerIds,
		consts.BasicFieldAuditorIds, consts.BasicFieldPlanStartTime, consts.BasicFieldPlanEndTime}
	return filterColumns
}

func ShareIssueCard(orgId, currentUserId int64, input projectvo.IssueCardShareReq) (bool, errs.SystemErrorInfo) {
	return true, nil
}
