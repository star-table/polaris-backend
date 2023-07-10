package service

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/lang"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

func IssueSourceList(orgId int64, page uint, size uint, params vo.IssueSourcesReq) (*vo.IssueSourceList, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//orgId := cacheUserInfo.OrgId

	var sourceList interface{} = nil
	var err1 error = nil
	if params.ProjectObjectTypeID != nil {
		sourceList, err1 = domain.GetIssueSourceListByProjectObjectTypeId(orgId, *params.ProjectObjectTypeID, params.ProjectID)
	} else {
		sourceList, err1 = domain.GetIssueSourceList(orgId, params.ProjectID)
	}
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err1)
	}

	resultList := &[]*vo.IssueSource{}
	copyErr := copyer.Copy(sourceList, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	lang := lang.GetLang()
	if len(lang) > 0 {
		if tmpMap, ok := consts.LANG_ISSUE_SRC_MAP[lang]; ok {
			for index, item := range *resultList {
				if tmpVal1, ok2 := tmpMap[item.LangCode]; ok2 {
					(*resultList)[index].Name = tmpVal1
				}
			}
		}
	}

	return &vo.IssueSourceList{
		Total: int64(len(*resultList)),
		List:  *resultList,
	}, nil
}

func CreateIssueSource(currentUserId int64, input vo.CreateIssueSourceReq) (*vo.Void, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//currentUserId := cacheUserInfo.UserId

	//TODO 权限
	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	//}

	entity := &bo.IssueSourceBo{}
	copyErr := copyer.Copy(input, entity)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueSource)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	entity.Id = id
	entity.Creator = currentUserId
	entity.Updator = currentUserId

	//添加清除缓存逻辑
	err = domain.DeleteIssueSourceListCache(input.OrgID)

	if err != nil {
		return nil, err
	}

	err1 := domain.CreateIssueSource(entity)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	return &vo.Void{
		ID: id,
	}, nil
}

func UpdateIssueSource(currentUserId int64, input vo.UpdateIssueSourceReq) (*vo.Void, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//currentUserId := cacheUserInfo.UserId

	//TODO 权限
	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	//}

	entity := &bo.IssueSourceBo{}
	copyErr := copyer.Copy(input, entity)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	entity.Updator = currentUserId

	//是否存在
	_, err2 := domain.GetIssueSourceBo(entity.Id)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}

	//添加清除缓存逻辑
	err := domain.DeleteIssueSourceListCache(input.OrgID)

	if err != nil {
		return nil, err
	}

	err1 := domain.UpdateIssueSource(entity)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	return &vo.Void{
		ID: input.ID,
	}, nil
}

func DeleteIssueSource(orgId, currentUserId int64, input vo.DeleteIssueSourceReq) (*vo.Void, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//currentUserId := cacheUserInfo.UserId
	targetId := input.ID

	//TODO 权限
	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	//}

	bo, err1 := domain.GetIssueSourceBo(targetId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	//添加清除缓存逻辑 这里暂时用传进来的orgId 后面有校验的时候用input里面的
	err := domain.DeleteIssueSourceListCache(orgId)

	if err != nil {
		return nil, err
	}

	err2 := domain.DeleteIssueSource(bo, currentUserId)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}

	return &vo.Void{
		ID: targetId,
	}, nil
}
