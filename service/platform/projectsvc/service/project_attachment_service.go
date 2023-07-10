package service

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"upper.io/db.v3"
)

func DeleteProjectAttachment(orgId, operatorId int64, input vo.DeleteProjectAttachmentReq) (*vo.DeleteProjectAttachmentResp, errs.SystemErrorInfo) {
	userAuthorityResp := userfacade.GetUserAuthority(orgId, operatorId)
	if userAuthorityResp.Failure() {
		log.Errorf("[DeleteProjectAttachment.GetUserAuthority] failed, orgId:%v, userId:%v, err:%v", orgId, operatorId, userAuthorityResp.Error())
		return nil, errs.BuildSystemErrorInfo(errs.UserDomainError, userAuthorityResp.Error())
	}
	isOrgOwner := userAuthorityResp.Data.IsOrgOwner
	isSysAdmin := userAuthorityResp.Data.IsSysAdmin
	isSubAdmin := userAuthorityResp.Data.IsSubAdmin
	if !isOrgOwner && !isSysAdmin && !isSubAdmin {
		return nil, errs.DeleteAttachmentError
	}

	err := domain.DeleteProjectAttachment(orgId, operatorId, input.ProjectID, input.ResourceIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.DeleteProjectAttachmentResp{
		ResourceIds: input.ResourceIds,
	}, nil
}

func GetProjectAttachment(orgId, operatorId int64, page, size int, input vo.ProjectAttachmentReq) (*vo.AttachmentList, errs.SystemErrorInfo) {
	// [权限迭代20211026] 权限项被删除，对应的鉴权由字段权限控制。
	//err := domain.AuthProjectWithOutPermission(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProAttachment, consts.OperationProAttachmentView)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}

	// 获取resources
	resourceInput := bo.GetResourceBo{
		UserId:     operatorId,
		OrgId:      orgId,
		ProjectId:  input.ProjectID,
		Page:       page,
		Size:       size,
		SourceType: []int{consts.OssPolicyTypeIssueResource, consts.OssPolicyTypeLesscodeResource},
	}
	if input.FileType != nil {
		resourceInput.FileType = input.FileType
	}
	if input.KeyWord != nil {
		resourceInput.KeyWord = input.KeyWord
	}
	resp := resourcefacade.GetResource(resourcevo.GetResourceReqVo{
		Input: resourceInput,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}
	resourceList := resp.ResourceList.List
	total := resp.ResourceList.Total
	resourceIds := make([]int64, len(resourceList))
	for i, value := range resourceList {
		resourceIds[i] = value.ID
	}

	// 获取relations
	//cond := db.Cond{
	//	consts.TcIsDelete:     consts.AppIsNoDelete,
	//	consts.TcProjectId:    input.ProjectID,
	//	consts.TcOrgId:        orgId,
	//	consts.TcRelationType: consts.IssueRelationTypeResource,
	//	consts.TcRelationId:   db.In(resourceIds),
	//}
	//issueRelationPos, err := domain.GetTotalResourceByRelationCond(cond)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	////获取resourceId和issueId的映射关系
	//relationMap := map[int64][]int64{}
	//issueIdMap := map[int64]bool{}
	//for _, value := range *issueRelationPos {
	//	issueIdMap[value.IssueId] = true
	//	issueId := value.IssueId
	//	resourceId := value.RelationId
	//	relationMap[resourceId] = append(relationMap[resourceId], issueId)
	//	//fmt.Printf("relationId:%d,issueId:%d,resourceId:%d\n", value.Id, value.IssueId, value.RelationId)
	//}
	//issueIds := make([]int64, len(issueIdMap))
	//for issueId, _ := range issueIdMap {
	//	issueIds = append(issueIds, issueId)
	//}

	resourceResp := resourcefacade.GetIssueIdsByResourceIds(resourcevo.GetIssueIdsByResourceIdsReqVo{
		OrgId:  orgId,
		UserId: 0,
		Input: resourcevo.GetIssueIdsByResourceIdsReq{
			ProjectId:   input.ProjectID,
			ResourceIds: resourceIds,
		},
	})
	if resourceResp.Failure() {
		log.Errorf("[GetProjectAttachment] GetIssueIdsByResourceIds err:%, projectId:%v", resourceResp.Error(), input.ProjectID)
		return nil, resourceResp.Error()
	}

	issueIds := []int64{}
	relationMap := map[int64][]int64{}
	for _, v := range resourceResp.Data {
		ids := slice.SliceUniqueInt64(v.IssueIds)
		issueIds = append(issueIds, ids...)
		relationMap[v.ResourceId] = ids
	}

	issueIds = slice.SliceUniqueInt64(issueIds)

	issueBoList, err := domain.GetIssueInfosLc(orgId, operatorId, issueIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	attachmentBos, err := domain.AssemblyAttachmentList(resourceList, issueBoList, relationMap)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	list := make([]*vo.Attachment, 0)
	err1 := copyer.Copy(attachmentBos, &list)
	if err1 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	for _, attachmentBo := range attachmentBos {
		log.Infof("[GetProjectAttachment] attachmentBo: %v", attachmentBo)
	}
	for _, attachment := range list {
		log.Infof("[GetProjectAttachment] attachment: %v", *attachment)
	}

	return &vo.AttachmentList{Total: total, List: list}, nil
}

func GetProjectAttachmentInfo(orgId, operatorId int64, input vo.ProjectAttachmentInfoReq) (*vo.Attachment, errs.SystemErrorInfo) {
	log.Infof("[GetProjectAttachmentInfo] GetProjectAttachmentInfo orgId: %d, operatorId: %d, input: %v", orgId, operatorId, input)
	resourceInput := bo.GetResourceInfoBo{
		UserId:      operatorId,
		OrgId:       orgId,
		AppId:       input.AppID,
		ResourceId:  input.ResourceID,
		IssueId:     input.IssueID,
		SourceTypes: []int{consts.OssPolicyTypeIssueResource, consts.OssPolicyTypeLesscodeResource},
	}

	if input.FileType != nil {
		resourceInput.FileType = input.FileType
	}

	resp := resourcefacade.GetResourceInfo(resourcevo.GetResourceInfoReqVo{
		Input: resourceInput,
	})

	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}

	resourceList := make([]*vo.Resource, 1)
	resourceList[0] = &vo.Resource{
		ID:             resp.ID,
		OrgID:          resp.OrgID,
		Host:           resp.Host,
		Path:           resp.Path,
		PathCompressed: resp.PathCompressed,
		Name:           resp.Name,
		Type:           resp.Type,
		Size:           resp.Size,
		CreatorName:    resp.CreatorName,
		Suffix:         resp.Suffix,
		Md5:            resp.Md5,
		FileType:       resp.FileType,
		Creator:        resp.Creator,
		CreateTime:     resp.CreateTime,
		Updator:        resp.Updator,
		UpdateTime:     resp.UpdateTime,
		Version:        resp.Version,
		IsDelete:       resp.IsDelete,
	}

	resourceIds := make([]int64, 1)
	resourceIds[0] = resp.ID

	cond := db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcProjectId:    input.AppID,
		consts.TcOrgId:        orgId,
		consts.TcRelationType: consts.IssueRelationTypeResource,
		consts.TcRelationId:   db.In(resourceIds),
	}
	issueRelationPos, err := domain.GetTotalResourceByRelationCond(cond)
	if err != nil {
		//log.Error(err)
		log.Errorf("[GetProjectAttachment] GetTotalResourceByRelationCond Failed: %s, Cond: %v", strs.ObjectToString(err), cond)
		return nil, err
	}
	//获取resourceId和issueId的映射关系
	relationMap := map[int64][]int64{}
	issueIdMap := map[int64]bool{}
	for _, value := range *issueRelationPos {
		issueIdMap[value.IssueId] = true
		issueId := value.IssueId
		resourceId := value.RelationId
		relationMap[resourceId] = append(relationMap[resourceId], issueId)
		//fmt.Printf("relationId:%d,issueId:%d,resourceId:%d\n", value.Id, value.IssueId, value.RelationId)
	}
	issueIds := make([]int64, len(issueIdMap))
	for issueId, _ := range issueIdMap {
		issueIds = append(issueIds, issueId)
	}
	issueBoList, err := domain.GetIssueInfosLc(orgId, operatorId, issueIds)
	if err != nil {
		//log.Error(err)
		log.Errorf("[GetProjectAttachment] GetIssueInfoList Failed: %s, issueIds: %v", strs.ObjectToString(err), issueIds)
		return nil, err
	}
	attachmentBos, err := domain.AssemblyAttachmentList(resourceList, issueBoList, relationMap)
	if err != nil {
		//log.Error(err)
		log.Errorf("[GetProjectAttachment] AssemblyAttachmentList Failed: %s, resourceList: %v, issueBoList: %v, relationMap: %v", strs.ObjectToString(err), resourceList, issueBoList, relationMap)
		return nil, err
	}
	list := make([]*vo.Attachment, 0)
	err1 := copyer.Copy(attachmentBos, &list)
	if err1 != nil {
		//log.Error(err)
		log.Errorf("[GetProjectAttachment] Copy Failed: %s, attachmentBos: %v, list: %v", strs.ObjectToString(err), attachmentBos, list)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}

	if len(attachmentBos) == 0 {
		log.Info("[GetProjectAttachment] attachmentBos == 0")
		return nil, nil
	}

	issueList := make([]*vo.Issue, 0)
	err1 = copyer.Copy(issueBoList, &issueList)
	if err1 != nil {
		//log.Error(err)
		log.Errorf("[GetProjectAttachment] Copy Failed: %s, issueList: %v, list: %v", strs.ObjectToString(err), issueList, list)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}

	return &vo.Attachment{
		ID:             attachmentBos[0].ID,
		OrgID:          attachmentBos[0].OrgID,
		Host:           attachmentBos[0].Host,
		Path:           attachmentBos[0].Path,
		OfficeURL:      "",
		PathCompressed: attachmentBos[0].PathCompressed,
		Name:           attachmentBos[0].Name,
		Type:           attachmentBos[0].Type,
		Size:           attachmentBos[0].Size,
		CreatorName:    attachmentBos[0].CreatorName,
		Suffix:         attachmentBos[0].Suffix,
		Md5:            attachmentBos[0].Md5,
		FileType:       attachmentBos[0].FileType,
		Creator:        attachmentBos[0].Creator,
		CreateTime:     attachmentBos[0].CreateTime,
		Updator:        attachmentBos[0].Updator,
		UpdateTime:     attachmentBos[0].UpdateTime,
		Version:        attachmentBos[0].Version,
		// IsDelete:       attachmentBos[0].IsDelete,
		IssueList: issueList,
	}, nil
}
