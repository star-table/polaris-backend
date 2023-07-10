package domain

import (
	int642 "github.com/star-table/polaris-backend/common/core/util/slice/int64"
	"strconv"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/appvo"
	"github.com/star-table/polaris-backend/facade/appfacade"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetProjectRelationBoIgnoreStatus(orgId, projectId int64, relationType int, relationId int64, projectObjectTypeId int64) (*bo.ProjectRelationBo, errs.SystemErrorInfo) {
	projectRelation, err := dao.SelectOneProjectRelation(db.Cond{
		consts.TcOrgId:               orgId,
		consts.TcProjectId:           projectId,
		consts.TcRelationId:          relationId,
		consts.TcRelationType:        relationType,
		consts.TcProjectObjectTypeId: projectObjectTypeId,
		consts.TcIsDelete:            consts.AppIsNoDelete,
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotRelatedError)
	}

	bo := &bo.ProjectRelationBo{}
	err1 := copyer.Copy(projectRelation, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func GetProjectRelationByType(projectId int64, relationTypes []int64) (*[]bo.ProjectRelationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcProjectId:    projectId,
		consts.TcRelationType: db.In(relationTypes),
		consts.TcStatus:       consts.AppStatusEnable,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotRelatedError)
	}
	bos := &[]bo.ProjectRelationBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bos, nil
}

func GetProjectRelationByTypeAndUserId(orgId, userId int64, relationTypes []int64) (*[]bo.ProjectRelationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcRelationId:   userId,
		consts.TcRelationType: db.In(relationTypes),
		consts.TcStatus:       consts.AppStatusEnable,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotRelatedError)
	}
	bos := &[]bo.ProjectRelationBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bos, nil
}

//更新项目关联，多关联类型check
//relationTypes: 使用该类型范围内的数据与relationIds做过滤，得到真正需要关联的id
//targetRelationType: 目标关联类型，新增的关联以此类型为准
func UpdateProjectRelationWithRelationTypes(operatorId, orgId, projectId int64, relationTypes []int, targetRelationType int, relationIds []int64) ([]int64, errs.SystemErrorInfo) {
	//防止项目成员重复插入
	uid := uuid.NewUuid()
	projectIdStr := strconv.FormatInt(projectId, 10)
	lockKey := consts.AddProjectRelationLock + projectIdStr
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, err)
		return nil, errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	} else {
		return nil, errs.BuildSystemErrorInfo(errs.GetDistributedLockError)
	}

	//预先查询已有的关联
	projectRelations := &[]po.PpmProProjectRelation{}
	err5 := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcProjectId:    projectId,
		consts.TcRelationId:   db.In(relationIds),
		consts.TcRelationType: db.In(relationTypes),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, projectRelations)
	if err5 != nil {
		log.Error(err5)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	//check，去掉已有的关联
	if len(*projectRelations) > 0 {
		notRelationUserIds := make([]int64, 0)
		allExistIds := []int64{}
		for _, issueRelation := range *projectRelations {
			allExistIds = append(allExistIds, issueRelation.RelationId)
		}
		for _, id := range relationIds {
			exist, err := slice.Contain(allExistIds, id)
			if err != nil {
				log.Error(err)
				continue
			}
			if !exist {
				notRelationUserIds = append(notRelationUserIds, id)
			}
		}

		relationIds = notRelationUserIds
	}
	relationIds = slice.SliceUniqueInt64(relationIds)

	relationIdsSize := len(relationIds)
	if relationIdsSize == 0 {
		return relationIds, nil
	}

	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, relationIdsSize)
	if err != nil {
		log.Errorf("id generate: %q\n", err)
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}

	projectRelationPos := make([]po.PpmProProjectRelation, relationIdsSize)
	for i, relationId := range relationIds {
		id := ids.Ids[i].Id
		issueRelation := &po.PpmProProjectRelation{}
		issueRelation.Id = id
		issueRelation.OrgId = orgId
		issueRelation.ProjectId = projectId
		issueRelation.RelationId = relationId
		issueRelation.RelationType = targetRelationType
		issueRelation.Creator = operatorId
		issueRelation.Updator = operatorId
		issueRelation.IsDelete = consts.AppIsNoDelete
		projectRelationPos[i] = *issueRelation
	}

	err2 := mysql.BatchInsert(&po.PpmProProjectRelation{}, slice.ToSlice(projectRelationPos))
	if err2 != nil {
		log.Errorf("mysql.BatchInsert(): %q\n", err2)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	return relationIds, nil
}

//更新项目关联，带分布式锁
func UpdateProjectRelation(operatorId, orgId, projectId int64, relationType int, relationIds []int64) errs.SystemErrorInfo {
	_, err := UpdateProjectRelationWithRelationTypes(operatorId, orgId, projectId, []int{relationType}, relationType, relationIds)
	return err
}

func GetProjectRelationByCond(cond db.Cond) (*[]bo.ProjectRelationBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	pos := &[]po.PpmProProjectRelation{}
	err = conn.Collection(consts.TableProjectRelation).Find(cond).OrderBy("create_time asc").All(pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotRelatedError)
	}
	bos := &[]bo.ProjectRelationBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bos, nil
}

func GetProjectRelationByCondSort(cond db.Cond, order interface{}) (*[]bo.ProjectRelationBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	pos := &[]po.PpmProProjectRelation{}
	err = conn.Collection(consts.TableProjectRelation).Find(cond).OrderBy(order).All(pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotRelatedError)
	}
	bos := &[]bo.ProjectRelationBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bos, nil
}

func IsProjectOwner(orgId, userId, projectId int64) (bool, errs.SystemErrorInfo) {
	info := &po.PpmProProject{}
	err := mysql.SelectOneByCond(consts.TableProject, db.Cond{
		consts.TcOrgId: orgId,
		consts.TcId:    projectId,
	}, info)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return false, nil
		}
		log.Error(err)
		return false, errs.MysqlOperateError
	}

	if info.Owner == userId {
		return true, nil
	}
	return false, nil
}

// GetProjectIdByOpenChatId 通过飞书的 chatId 获取项目id。如果项目不存在或者删除了，则返回 0 值。
func GetProjectIdByOpenChatId(openChatId string) (int64, int64, errs.SystemErrorInfo) {
	settingPo := po.PpmProProjectChat{}
	err := mysql.SelectOneByCond(consts.TableProjectChat, db.Cond{
		consts.TcChatId:   openChatId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &settingPo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return 0, 0, nil
		}
		log.Errorf("[GetProjectIdByOpenChatId] query setting po err: %v, openChatId: %v", err, openChatId)
		return 0, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	_, busiErr := GetProjectSimple(settingPo.OrgId, settingPo.ProjectId)
	if busiErr != nil {
		if busiErr == errs.ProjectNotExist {
			return 0, 0, nil
		} else {
			log.Errorf("[GetProjectIdByOpenChatId] GetProject err: %v, openChatId: %s", busiErr, openChatId)
			return 0, 0, busiErr
		}
	}
	return settingPo.OrgId, settingPo.ProjectId, nil
}

// IsProjectParticipant 检查一个用户是否是项目/应用参与人
func IsProjectParticipant(orgId, userId, projectId int64) (bool, errs.SystemErrorInfo) {
	project, err := GetProjectSimple(orgId, projectId)
	if err != nil {
		log.Errorf("[IsProjectParticipant] GetProjectSimple err: %v, orgId: %d, userId: %d, projectId: %d", err, orgId, userId, projectId)
		return false, err
	}
	checkResp := appfacade.IsAppMember(appvo.IsAppMemberReq{
		AppId:  project.AppId,
		OrgId:  orgId,
		UserId: userId,
	})
	if checkResp.Failure() {
		log.Errorf("[IsProjectParticipant] err: %v, orgId: %d, userId: %d, projectId: %d", checkResp.Err, orgId, userId, projectId)
		return false, checkResp.Error()
	}

	manageAuthInfoResp := userfacade.GetUserAuthority(orgId, userId)
	if manageAuthInfoResp.Failure() {
		log.Errorf("[IsProjectParticipant] GetUserAuthority err:%v, orgId:%v, projectId:%v",
			manageAuthInfoResp.Error(), orgId, projectId)
		return false, manageAuthInfoResp.Error()
	}
	isSysAdmin := manageAuthInfoResp.Data.IsSysAdmin
	isSubAdmin := manageAuthInfoResp.Data.IsSubAdmin
	isOrgOwner := manageAuthInfoResp.Data.IsOrgOwner
	manageApps := manageAuthInfoResp.Data.ManageApps

	isAdmin := isSysAdmin || isOrgOwner || (isSubAdmin && len(manageApps) > 0 && manageApps[0] == -1)
	canManagePartialApp := false
	if isSubAdmin && len(manageApps) > 0 && manageApps[0] != -1 {
		canManagePartialApp = int642.InArray(project.AppId, manageApps)
	}
	return checkResp.Data || isAdmin || canManagePartialApp, nil
}

func GetProjectIdByAppId(appId string) (int64, int64, errs.SystemErrorInfo) {
	project := &po.PpmProProject{}
	err := mysql.SelectOneByCond(consts.TableProject, db.Cond{
		consts.TcAppId:    appId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, project)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return 0, 0, nil
		}
		return 0, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return project.OrgId, project.Id, nil
}

func GetProjectInfoByAppId(appId int64) (*po.PpmProProject, errs.SystemErrorInfo) {
	project := &po.PpmProProject{}
	err := mysql.SelectOneByCond(consts.TableProject, db.Cond{
		consts.TcAppId:    appId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, project)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return project, nil
}
