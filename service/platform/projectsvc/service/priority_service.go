package service

//
//import (
//	"github.com/star-table/common/core/util/copyer"
//	"github.com/star-table/polaris-backend/common/core/consts"
//	"github.com/star-table/polaris-backend/common/core/errs"
//	"github.com/star-table/polaris-backend/common/core/util/lang"
//	"github.com/star-table/polaris-backend/common/model/bo"
//	"github.com/star-table/polaris-backend/common/model/vo"
//	"github.com/star-table/polaris-backend/facade/idfacade"
//	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
//	"upper.io/db.v3"
//)
//
//func PriorityList(orgId int64, page uint, size uint, cond db.Cond) (*vo.PriorityList, errs.SystemErrorInfo) {
//	priorityList, err := domain.GetPriorityListByType(orgId, consts.PriorityTypeIssue)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	bo.SortPriorityBo(*priorityList)
//
//	resultList := &[]*vo.Priority{}
//	copyErr := copyer.Copy(priorityList, resultList)
//	if copyErr != nil {
//		log.Errorf("对象copy异常: %v", copyErr)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
//	}
//	lang := lang.GetLang()
//	if len(lang) > 0 {
//		if tmpMap, ok := consts.LANG_PRIORITIES_NAME_MAP[lang]; ok {
//			for index, item := range *resultList {
//				if tmpVal1, ok2 := tmpMap[item.Name]; ok2 {
//					(*resultList)[index].Name = tmpVal1
//				}
//			}
//		}
//	}
//
//	return &vo.PriorityList{
//		Total: int64(len(*resultList)),
//		List:  *resultList,
//	}, nil
//}
//
//func CreatePriority(currentUserId int64, input vo.CreatePriorityReq) (*vo.Void, errs.SystemErrorInfo) {
//	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
//	//if err != nil {
//	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
//	//}
//	//currentUserId := cacheUserInfo.UserId
//
//	//TODO 权限
//	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
//	//if err != nil {
//	//	log.Error(err)
//	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
//	//}
//
//	entity := &bo.PriorityBo{}
//	copyErr := copyer.Copy(input, entity)
//	if copyErr != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
//	}
//
//	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TablePriority)
//	if err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
//	}
//	entity.Id = id
//	entity.Creator = currentUserId
//	entity.Updator = currentUserId
//
//	//清楚缓存
//	err1 := domain.DeletePriorityListCache(input.OrgID)
//
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err1)
//	}
//
//	err1 = domain.CreatePriority(entity)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
//	}
//
//	return &vo.Void{
//		ID: id,
//	}, nil
//}
//
//func UpdatePriority(currentUserId int64, input vo.UpdatePriorityReq) (*vo.Void, errs.SystemErrorInfo) {
//	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
//	//if err != nil {
//	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
//	//}
//	//currentUserId := cacheUserInfo.UserId
//
//	//TODO 权限
//	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
//	//if err != nil {
//	//	log.Error(err)
//	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
//	//}
//
//	entity := &bo.PriorityBo{}
//	copyErr := copyer.Copy(input, entity)
//	if copyErr != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
//	}
//	entity.Updator = currentUserId
//
//	//是否存在
//	_, err2 := domain.GetPriorityBo(entity.Id)
//	if err2 != nil {
//		log.Error(err2)
//		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
//	}
//
//	//清楚缓存
//	err1 := domain.DeletePriorityListCache(input.OrgID)
//
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err1)
//	}
//
//	err1 = domain.UpdatePriority(entity)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
//	}
//
//	return &vo.Void{
//		ID: input.ID,
//	}, nil
//}
//
//func DeletePriority(orgId, currentUserId int64, input vo.DeletePriorityReq) (*vo.Void, errs.SystemErrorInfo) {
//	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
//	//if err != nil {
//	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
//	//}
//	//currentUserId := cacheUserInfo.UserId
//	targetId := input.ID
//
//	//TODO 权限
//	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
//	//if err != nil {
//	//	log.Error(err)
//	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
//	//}
//
//	bo, err1 := domain.GetPriorityBo(targetId)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
//	}
//
//	//清楚缓存 暂时用ctx 传进来的当前用户orgId 等多组织了用input中的 在获取用户中做org的包含校验
//	err1 = domain.DeletePriorityListCache(orgId)
//
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err1)
//	}
//
//	err2 := domain.DeletePriority(bo, currentUserId)
//	if err2 != nil {
//		log.Error(err2)
//		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
//	}
//
//	return &vo.Void{
//		ID: targetId,
//	}, nil
//}
//
//func VerifyPriority(orgId int64, typ int, beVerifyId int64) (bool, errs.SystemErrorInfo) {
//	return domain.VerifyPriority(orgId, typ, beVerifyId)
//}
//
////func GetPriorityById(orgId int64, id int64) (*bo.PriorityBo, errs.SystemErrorInfo) {
////	return domain.GetPriorityById(orgId, id)
////}
