package domain

import (
	"github.com/spf13/cast"
	"strconv"
	"time"

	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"

	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
)

func DeleteAllIssueRelation(tx sqlbuilder.Tx, operatorId int64, orgId int64, issueIds []int64, recycleVersionId int64) errs.SystemErrorInfo {
	//删除之前的关联
	_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIssueId:  db.In(issueIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcVersion:  recycleVersionId,
		consts.TcUpdator:  operatorId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func DeleteAllIssueRelationByIds(orgId, operatorId int64, relationIds []int64, recycleVersionId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	//删除之前的关联
	_, err := mysql.TransUpdateSmartWithCond(tx[0], consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       db.In(relationIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete:   consts.AppIsDeleted,
		consts.TcUpdator:    operatorId,
		consts.TcUpdateTime: time.Now(),
		consts.TcVersion:    recycleVersionId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func DeleteIssueRelation(operatorId int64, issueBo bo.IssueBo, relationType int) errs.SystemErrorInfo {
	orgId := issueBo.OrgId
	issueId := issueBo.Id
	//删除之前的关联
	_, err := mysql.UpdateSmartWithCond(consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: relationType,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operatorId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func DeleteIssueRelationByIds(operatorId int64, issueBo bo.IssueBo, relationType int, relationIds []int64) errs.SystemErrorInfo {
	if relationIds == nil || len(relationIds) == 0 {
		return nil
	}

	orgId := issueBo.OrgId
	issueId := issueBo.Id

	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除之前的关联
		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcIssueId:      issueId,
			consts.TcRelationId:   db.In(relationIds),
			consts.TcRelationType: relationType,
			consts.TcIsDelete:     consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  operatorId,
		})
		if err != nil {
			log.Error(err)
			return err
		}

		if relationType == consts.IssueRelationTypeResource {
			//删除文件
			resp := resourcefacade.DeleteResource(resourcevo.DeleteResourceReqVo{
				Input: bo.DeleteResourceBo{
					ResourceIds: relationIds,
					UserId:      operatorId,
					OrgId:       orgId,
					ProjectId:   issueBo.ProjectId,
					IssueId:     issueId,
				},
			})
			if resp.Failure() {
				log.Error(resp.Message)
				return resp.Error()
			}
		}

		return nil
	})

	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

func UpdateIssueRelationSingle(operatorId int64, issueBo *bo.IssueBo, relationType int, newUserIds []int64) (*bo.IssueRelationBo, errs.SystemErrorInfo) {
	bos, err := UpdateIssueRelation(operatorId, issueBo, relationType, newUserIds, "")
	if err != nil {
		return nil, err
	}
	if len(bos) == 0 {
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError)
	}
	return &bos[0], nil
}

func UpdateIssueRelation(operatorId int64, issueBo *bo.IssueBo, relationType int, newUserIds []int64, relationCode string) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	orgId := issueBo.OrgId
	issueId := issueBo.Id

	//防止项目成员重复插入
	uid := uuid.NewUuid()
	issueIdStr := strconv.FormatInt(issueId, 10)
	relationTypeStr := strconv.Itoa(relationType)
	lockKey := consts.AddIssueRelationLock + issueIdStr + "#" + relationTypeStr
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
	issueRelations := &[]po.PpmPriIssueRelation{}
	err5 := mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
		consts.TcIssueId:      issueBo.Id,
		consts.TcRelationId:   db.In(newUserIds),
		consts.TcRelationType: relationType,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, issueRelations)
	if err5 != nil {
		log.Error(err5)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	//check，去掉已有的关联
	if len(*issueRelations) > 0 {
		alreadyExistRelationIdMap := map[int64]bool{}
		for _, issueRelation := range *issueRelations {
			alreadyExistRelationIdMap[issueRelation.RelationId] = true
		}
		notRelationUserIds := make([]int64, 0)
		for _, newUserId := range newUserIds {
			if _, ok := alreadyExistRelationIdMap[newUserId]; !ok {
				notRelationUserIds = append(notRelationUserIds, newUserId)
			}
		}
		newUserIds = notRelationUserIds
	}
	newUserIds = slice.SliceUniqueInt64(newUserIds)

	issueRelationBos := make([]bo.IssueRelationBo, len(newUserIds))
	if len(newUserIds) == 0 {
		return issueRelationBos, nil
	}

	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, len(newUserIds))
	if err != nil {
		log.Errorf("id generate: %q\n", err)
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}

	issueRelationPos := make([]po.PpmPriIssueRelation, len(newUserIds))
	for i, newUserId := range newUserIds {
		id := ids.Ids[i].Id
		issueRelation := &po.PpmPriIssueRelation{}
		issueRelation.Id = id
		issueRelation.OrgId = orgId
		issueRelation.ProjectId = issueBo.ProjectId
		issueRelation.IssueId = issueBo.Id
		issueRelation.RelationId = newUserId
		issueRelation.RelationType = relationType
		issueRelation.Creator = operatorId
		issueRelation.Updator = operatorId
		issueRelation.IsDelete = consts.AppIsNoDelete
		issueRelation.RelationCode = relationCode
		issueRelationPos[i] = *issueRelation

		issueRelationBos[i] = bo.IssueRelationBo{
			Id:           id,
			OrgId:        issueBo.OrgId,
			IssueId:      issueBo.Id,
			RelationId:   newUserId,
			RelationType: consts.IssueRelationTypeOwner,
			Creator:      operatorId,
			CreateTime:   types.NowTime(),
			Updator:      operatorId,
			UpdateTime:   types.NowTime(),
			Version:      1,
		}
	}

	err2 := PaginationInsert(slice.ToSlice(issueRelationPos), &po.PpmPriIssueRelation{})
	if err2 != nil {
		log.Errorf("mysql.BatchInsert(): %q\n", err2)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}
	return issueRelationBos, nil
}

func GetIssueRelationIdsByRelateType(orgId int64, issueId int64, relationType int) (*[]int64, errs.SystemErrorInfo) {
	issueParticipantRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: relationType,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, bo.PageBo{
		Order: consts.TcCreateTime + " asc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	relationIds := make([]int64, len(*issueParticipantRelations))
	for i, participantRelation := range *issueParticipantRelations {
		relationIds[i] = participantRelation.RelationId
	}
	return &relationIds, nil
}

func GetIssueRelationByRelateTypeList(orgId int64, issueId int64, relationTypes []int) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	issueParticipantRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: db.In(relationTypes),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, bo.PageBo{
		Order: consts.TcCreateTime + " desc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	relationBos := &[]bo.IssueRelationBo{}
	_ = copyer.Copy(issueParticipantRelations, relationBos)
	return *relationBos, nil
}

//创建者也可以删除资源
func GetIssueResourceIdsByCreator(orgId int64, issueId int64, ids []int64, creatorId int64) (*[]int64, errs.SystemErrorInfo) {
	issueRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcId:           db.In(ids),
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: consts.IssueRelationTypeResource,
		consts.TcCreator:      creatorId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, bo.PageBo{
		Order: consts.TcCreateTime + " desc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	relationIds := make([]int64, len(*issueRelations))
	for i, participantRelation := range *issueRelations {
		relationIds[i] = participantRelation.Id
	}
	return &relationIds, nil
}

//func VerifyRelationIssue(issueIds []int64, projectObjectTypeId int64, orgId int64) errs.SystemErrorInfo {
//	issueList := &[]po.PpmPriIssue{}
//	err := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
//		consts.TcOrgId: orgId,
//		consts.TcId:    db.In(issueIds),
//		//consts.TcProjectObjectTypeId: projectObjectTypeId,
//	}, issueList)
//	if err != nil {
//		log.Error(err)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//
//	issueIds = slice.SliceUniqueInt64(issueIds)
//	trueIssueIds := []int64{}
//	for _, v := range *issueList {
//		trueIssueIds = append(trueIssueIds, v.Id)
//	}
//	//去重后的数组长度相同则表示传递的任务id都有效
//	if len(issueIds) == len(trueIssueIds) {
//		return nil
//	}
//
//	return errs.BuildSystemErrorInfo(errs.RelationIssueError)
//}

//func RelationIssueList(orgId, issueId int64, pushType int) ([]bo.IssueBo, errs.SystemErrorInfo) {
//	issueRelationList := &[]po.PpmPriIssueRelation{}
//	issueList := &[]po.PpmPriIssue{}
//	issueBo := &[]bo.IssueBo{}
//	cond := db.Cond{
//		consts.TcOrgId:        orgId,
//		consts.TcIsDelete:     consts.AppIsNoDelete,
//		consts.TcRelationType: pushType,
//	}
//	union := db.Or(db.Cond{
//		consts.TcIssueId: issueId,
//	}).Or(db.Cond{
//		consts.TcRelationId: issueId,
//	})
//	_, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIssueRelation, cond, union, 0, 0, "create_time asc", issueRelationList)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//	if len(*issueRelationList) == 0 {
//		return *issueBo, nil
//	}
//
//	issueIds := []int64{}
//	relateIssueIds := []int64{}
//	beRelatedIssueIds := []int64{}
//	for _, v := range *issueRelationList {
//		if v.IssueId == issueId {
//			//关联
//			relateIssueIds = append(relateIssueIds, v.RelationId)
//			issueIds = append(issueIds, v.RelationId)
//		} else {
//			//被关联
//			beRelatedIssueIds = append(beRelatedIssueIds, v.IssueId)
//			issueIds = append(issueIds, v.IssueId)
//		}
//	}
//
//	//按照传入id排序
//	orderBySort := ""
//	for _, id := range issueIds {
//		orderBySort += fmt.Sprintf(",%d", id)
//	}
//	_, err = mysql.SelectAllByCondWithPageAndOrder(consts.TableIssue, db.Cond{
//		consts.TcOrgId:    orgId,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//		consts.TcId:       db.In(issueIds),
//	}, nil, 0, 0, db.Raw("FIELD(id"+orderBySort+")"), issueList)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//
//	copyer.Copy(issueList, issueBo)
//	result := []bo.IssueBo{}
//	for _, b := range *issueBo {
//		if ok, _ := slice.Contain(relateIssueIds, b.Id); ok {
//			temp := b
//			temp.TypeForRelate = 1
//			result = append(result, temp)
//		}
//		if ok, _ := slice.Contain(beRelatedIssueIds, b.Id); ok {
//			temp := b
//			temp.TypeForRelate = 2
//			result = append(result, temp)
//		}
//	}
//
//	return result, nil
//}

func GetRelateIssueCount(orgId int64, issueIds []int64, relationType int) (map[int64]int64, errs.SystemErrorInfo) {
	//"select issue_id, count(issue_id) from ppm_pri_issue_relation where org_id=? and is_delete=2 and relation_type=4 and issue_id in(issueIds) group by issue_id"
	pos := []po.PpmPriIssueRelationCount{}
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Errorf("[GetRelateIssueCount] mysql链接异常: %v", err)
		return nil, errs.MysqlOperateError
	}
	conds := db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: relationType,
		consts.TcIssueId:      db.In(issueIds),
	}
	err = conn.Select(db.Raw("issue_id, count(issue_id) as total")).
		From(consts.TableIssueRelation).
		Where(conds).
		GroupBy(consts.TcIssueId).
		All(&pos)

	if err != nil {
		log.Errorf("[GetRelateIssueCount] mysql查询错误: %v", err)
		return nil, errs.MysqlOperateError
	}

	result := make(map[int64]int64, len(pos))
	for _, relation := range pos {
		result[relation.IssueId] = relation.Total
	}

	return result, nil
}

func GetIssueRelationByResource(orgId int64, projectId int64, resourceIds []int64) (*[]po.PpmPriIssueRelation, errs.SystemErrorInfo) {
	issueParticipantRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: consts.IssueRelationTypeResource,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationId:   db.In(resourceIds),
	}, bo.PageBo{
		Order: consts.TcId + " desc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return issueParticipantRelations, nil
}

func GetTotalResourceByRelationCond(cond db.Cond) (*[]po.PpmPriIssueRelation, errs.SystemErrorInfo) {
	pos := &[]po.PpmPriIssueRelation{}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	resourceIds := []int64{}
	for _, value := range *pos {
		if isContain, _ := slice.Contain(resourceIds, value.RelationId); !isContain {
			resourceIds = append(resourceIds, value.RelationId)
		}
	}
	return pos, nil
}

func DeleteProjectAttachment(orgId, operatorId, projectId int64, resourceIds []int64) errs.SystemErrorInfo {
	issueRelationPos, err := GetIssueRelationByResource(orgId, projectId, resourceIds)
	if err != nil {
		log.Error(err)
		return err
	}
	realResourceIds := make([]int64, 0)
	relationIds := make([]int64, len(*issueRelationPos))
	realResourceMap := make(map[int64]bool)
	for index, value := range *issueRelationPos {
		realResourceMap[value.RelationId] = true
		relationIds[index] = value.Id
	}
	for key, _ := range realResourceMap {
		realResourceIds = append(realResourceIds, key)
	}
	if len(realResourceIds) != len(resourceIds) {
		return errs.InvalidResourceIdsError
	}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		recycleVersionId, versionErr := AddRecycleRecord(orgId, operatorId, projectId, resourceIds, consts.RecycleTypeAttachment, tx)
		if versionErr != nil {
			log.Error(versionErr)
			return versionErr
		}
		err := DeleteAllIssueRelationByIds(orgId, operatorId, relationIds, recycleVersionId, tx)
		if err != nil {
			log.Error(err)
			return nil
		}

		projectinfo, err := GetProjectSimple(orgId, projectId)
		if err != nil {
			return err
		}
		deleteInput := bo.DeleteResourceBo{
			ResourceIds:      resourceIds,
			UserId:           operatorId,
			OrgId:            orgId,
			ProjectId:        projectId,
			RecycleVersionId: recycleVersionId,
			AppId:            projectinfo.AppId,
		}
		resp := resourcefacade.DeleteResource(resourcevo.DeleteResourceReqVo{Input: deleteInput})
		if resp.Failure() {
			log.Error(resp.Error())
			return resp.Error()
		}
		return nil
	})

	asyn.Execute(func() {
		reqVo := resourcevo.GetResourceByIdReqBody{
			ResourceIds: resourceIds,
		}
		req := resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: reqVo}
		resp := resourcefacade.GetResourceById(req)
		resourceBos := resp.ResourceBos
		resourceNames := make([]string, len(resourceBos))
		resourceTrend := []bo.ResourceInfoBo{}
		for index, value := range resourceBos {
			resourceNames[index] = value.Name
			resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
				Name:   value.Name,
				Url:    value.Host + value.Path,
				Size:   value.Size,
				Suffix: value.Suffix,
			})
		}

		trendBo := bo.ProjectTrendsBo{
			PushType:   consts.PushTypeDeleteResource,
			OrgId:      orgId,
			ProjectId:  projectId,
			OperatorId: operatorId,
			NewValue:   json.ToJsonIgnoreError(resourceNames),
			Ext: bo.TrendExtensionBo{
				ResourceInfo: resourceTrend,
			},
		}

		asyn.Execute(func() {
			PushProjectTrends(trendBo)
		})
		//asyn.Execute(func() {
		//	PushProjectThirdPlatformNotice(trendBo)
		//})
	})

	return nil
}

func GetRelationInfoByIssueIds(issueIds []int64, relationTypes []int) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	relationInfos := &[]po.PpmPriIssueRelation{}
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcIssueId:  db.In(issueIds),
	}
	if len(relationTypes) != 0 {
		cond[consts.TcRelationType] = db.In(relationTypes)
	}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, cond, relationInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.IssueRelationBo{}
	copyErr := copyer.Copy(relationInfos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ObjectCopyError
	}

	return *bos, nil
}

// GetIssueRelationsByCond 根据条件查询 issue relation
func GetIssueRelationsByCond(issueIds []int64, relationTypes []int, baseCond db.Cond) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	relationInfos := &[]po.PpmPriIssueRelation{}
	cond := baseCond
	cond[consts.TcIssueId] = db.In(issueIds)
	if len(relationTypes) != 0 {
		cond[consts.TcRelationType] = db.In(relationTypes)
	}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, cond, relationInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.IssueRelationBo{}
	copyErr := copyer.Copy(relationInfos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ObjectCopyError
	}

	return *bos, nil
}

// GetIssueRelationsByCondAndProId 通过条件查询 issue relations
func GetIssueRelationsByCondAndProId(projectId int64, relationTypes []int, baseCond db.Cond) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	relationInfos := &[]po.PpmPriIssueRelation{}
	cond := baseCond
	cond[consts.TcProjectId] = projectId
	if len(relationTypes) != 0 {
		cond[consts.TcRelationType] = db.In(relationTypes)
	}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, cond, relationInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.IssueRelationBo{}
	copyErr := copyer.Copy(relationInfos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ObjectCopyError
	}

	return *bos, nil
}

// GetIssueRelationsByCondAndOrgId 通过条件和组织 id 查询 issue relations
func GetIssueRelationsByCondAndOrgId(orgId int64, relationTypes []int, baseCond db.Cond) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	relationInfos := &[]po.PpmPriIssueRelation{}
	cond := baseCond
	cond[consts.TcOrgId] = orgId
	if len(relationTypes) != 0 {
		cond[consts.TcRelationType] = db.In(relationTypes)
	}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, cond, relationInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.IssueRelationBo{}
	copyErr := copyer.Copy(relationInfos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ObjectCopyError
	}

	return *bos, nil
}

func GetIssueMembers(orgId int64, issueId int64) (*bo.IssueMembersBo, errs.SystemErrorInfo) {
	relationBos, err := GetIssueRelationByRelateTypeList(orgId, issueId, consts.IssueParticipantTypeList)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	ownerMap := map[int64]bool{}
	participantMap := map[int64]bool{}
	followerMap := map[int64]bool{}
	memberMap := map[int64]bool{}
	for _, v := range relationBos {
		if v.RelationType == consts.IssueRelationTypeFollower {
			followerMap[v.RelationId] = true
		} else if v.RelationType == consts.IssueRelationTypeParticipant {
			participantMap[v.RelationId] = true
		} else if v.RelationType == consts.IssueRelationTypeOwner {
			ownerMap[v.RelationId] = true
		}
		memberMap[v.RelationId] = true
	}

	var ownerIds, followerIds, participantIds, memberIds []int64
	for o, _ := range ownerMap {
		ownerIds = append(ownerIds, o)
	}
	for k, _ := range followerMap {
		followerIds = append(followerIds, k)
	}

	for k, _ := range participantMap {
		participantIds = append(participantIds, k)
	}

	for k, _ := range memberMap {
		memberIds = append(memberIds, k)
	}

	return &bo.IssueMembersBo{
		MemberIds:      memberIds,
		OwnerId:        ownerIds,
		ParticipantIds: participantIds,
		FollowerIds:    followerIds,
	}, nil
}

func GetIssueRelationResource(page, size int) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmPriIssueRelation{}
	_, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIssueRelation, db.Cond{
		consts.TcRelationType: consts.IssueRelationTypeResource,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, nil, page, size, nil, pos)

	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	res := &[]bo.IssueRelationBo{}
	_ = copyer.Copy(pos, res)
	return *res, nil
}

//任务的前置任务就是issueId=当前任务，relationId=前置任务的id。后置相反
//func AddBeforeAfterIssue(orgId, userId, issueId int64, projectId int64, addType int, relatedIssueIds []int64) ([]int64, errs.SystemErrorInfo) {
//	relatedIssueIds = slice.SliceUniqueInt64(relatedIssueIds)
//	if ok, _ := slice.Contain(relatedIssueIds, issueId); ok {
//		//不能把自己作为前后置任务
//		return nil, errs.CannotBeforeAfterSelf
//	}
//
//	issuePos := &[]po.PpmPriIssue{}
//	issuePosErr := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
//		consts.TcOrgId:    orgId,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//		consts.TcId:       db.In(relatedIssueIds),
//	}, issuePos)
//	if issuePosErr != nil {
//		log.Error(issuePosErr)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//	allTrulyIssueIds := []int64{}
//	for _, issue := range *issuePos {
//		allTrulyIssueIds = append(allTrulyIssueIds, issue.Id)
//	}
//	if len(relatedIssueIds) != len(allTrulyIssueIds) {
//		return nil, errs.IllegalityIssue
//	}
//
//	if len(allTrulyIssueIds) == 0 {
//		return allTrulyIssueIds, nil
//	}
//
//	pos := &[]po.PpmPriIssueRelation{}
//	cond := db.Cond{
//		consts.TcOrgId:        orgId,
//		consts.TcIsDelete:     consts.AppIsNoDelete,
//		consts.TcRelationType: consts.IssueRelationTypeBeforeAfter,
//	}
//	union := db.Or(db.Cond{
//		consts.TcIssueId: issueId,
//	}).Or(db.Cond{
//		consts.TcRelationId: issueId,
//	})
//	_, posErr := mysql.SelectAllByCondWithPageAndOrder(consts.TableIssueRelation, cond, union, 0, 0, "create_time asc", pos)
//	if posErr != nil {
//		log.Error(posErr)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//	beforeIds := []int64{}
//	afterIds := []int64{}
//	for _, relation := range *pos {
//		if relation.IssueId == issueId {
//			beforeIds = append(beforeIds, relation.RelationId)
//		} else {
//			afterIds = append(afterIds, relation.IssueId)
//		}
//	}
//
//	addIds := []int64{}
//	insertPos := []po.PpmPriIssueRelation{}
//	if addType == 1 {
//		//前置
//		_, addIds = util.GetDifMemberIds(beforeIds, allTrulyIssueIds)
//		if len(addIds) > 0 {
//			for _, id := range addIds {
//				if ok, _ := slice.Contain(afterIds, id); ok {
//					//前置任务不能是后置任务
//					return nil, errs.BeforeAfterIssueConflict
//				}
//				insertPos = append(insertPos, po.PpmPriIssueRelation{
//					Id:           0,
//					OrgId:        orgId,
//					ProjectId:    projectId, //这里正向反向可能导致projectId没啥用，不用来进行筛选条件
//					IssueId:      issueId,
//					RelationId:   id,
//					RelationType: consts.IssueRelationTypeBeforeAfter,
//					Creator:      userId,
//				})
//			}
//		}
//	} else {
//		//后置
//		_, addIds = util.GetDifMemberIds(afterIds, allTrulyIssueIds)
//		if len(addIds) > 0 {
//			for _, id := range addIds {
//				if ok, _ := slice.Contain(beforeIds, id); ok {
//					//后置任务不能是前置任务
//					return nil, errs.BeforeAfterIssueConflict
//				}
//				insertPos = append(insertPos, po.PpmPriIssueRelation{
//					Id:           0,
//					OrgId:        orgId,
//					ProjectId:    projectId, //这里正向反向可能导致projectId没啥用，不用来进行筛选条件
//					IssueId:      id,
//					RelationId:   issueId,
//					RelationType: consts.IssueRelationTypeBeforeAfter,
//					Creator:      userId,
//				})
//			}
//		}
//	}
//	if len(insertPos) == 0 {
//		return []int64{}, nil
//	}
//	idResp, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, len(insertPos))
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	for i, _ := range insertPos {
//		insertPos[i].Id = idResp.Ids[i].Id
//	}
//
//	insertErr := mysql.BatchInsert(&po.PpmPriIssueRelation{}, slice.ToSlice(insertPos))
//	if insertErr != nil {
//		log.Error(insertErr)
//		return nil, errs.MysqlOperateError
//	}
//
//	return addIds, nil
//}

//任务的前置任务就是issueId=当前任务，relationId=前置任务的id。后置相反
//func DeleteBeforeAfterIssue(orgId, userId, issueId int64, deleteType int, relatedIssueIds []int64) errs.SystemErrorInfo {
//	if len(relatedIssueIds) == 0 {
//		return nil
//	}
//	cond := db.Cond{
//		consts.TcIsDelete:     consts.AppIsNoDelete,
//		consts.TcOrgId:        orgId,
//		consts.TcRelationType: consts.IssueRelationTypeBeforeAfter,
//	}
//
//	if deleteType == 1 {
//		cond[consts.TcIssueId] = issueId
//		cond[consts.TcRelationId] = db.In(relatedIssueIds)
//	} else {
//		cond[consts.TcRelationId] = issueId
//		cond[consts.TcIssueId] = db.In(relatedIssueIds)
//	}
//
//	_, err := mysql.UpdateSmartWithCond(consts.TableIssueRelation, cond, mysql.Upd{
//		consts.TcIsDelete: consts.AppIsDeleted,
//		consts.TcUpdator:  userId,
//	})
//
//	if err != nil {
//		log.Error(err)
//		return errs.MysqlOperateError
//	}
//
//	return nil
//}

//func GetIssueAuditorRelation(orgId, issueId, userId int64) (*bo.IssueRelationBo, errs.SystemErrorInfo) {
//	info := &po.PpmPriIssueRelation{}
//	err := mysql.SelectOneByCond(consts.TableIssueRelation, db.Cond{
//		consts.TcIsDelete:     consts.AppIsNoDelete,
//		consts.TcOrgId:        orgId,
//		consts.TcIssueId:      issueId,
//		consts.TcRelationId:   userId,
//		consts.TcRelationType: consts.IssueRelationTypeAuditor,
//	}, info)
//	if err != nil {
//		if err == db.ErrNoMoreRows {
//			return nil, errs.NotIssueAuditor
//		} else {
//			log.Error(err)
//			return nil, errs.MysqlOperateError
//		}
//	}
//
//	infoBo := &bo.IssueRelationBo{}
//	_ = copyer.Copy(info, infoBo)
//
//	return infoBo, nil
//}

// IsIssueParticipant 检查当前用户是否是项目参与人，参与人的关系自定义（根据入参 relationTypes 而定）
func IsIssueParticipant(orgId, curUserId, projectId, issueId int64, relationTypes []int) (bool, errs.SystemErrorInfo) {
	relationBos, err := GetIssueRelationByRelateTypeList(orgId, issueId, relationTypes)
	if err != nil {
		log.Error(err)
		return false, err
	}
	userIds := make([]int64, 0)
	for _, item := range relationBos {
		userIds = append(userIds, item.RelationId)
	}
	isExist, _ := slice.Contain(userIds, curUserId)
	return isExist, nil
}

// GetIssueParticipantBatch 获取多个任务的参与人 id。参与人的定义根据传入的 relationTypes 而定。
func GetIssueParticipantBatch(orgId, projectId int64, issueIds []int64, relationTypes []int) (map[int64][]int64, errs.SystemErrorInfo) {
	mapArr := make(map[int64][]int64, 0)
	relations, err := GetIssueRelationsByCond(issueIds, relationTypes, db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	})
	if err != nil {
		return mapArr, err
	}
	for _, item := range relations {
		if _, ok := mapArr[item.IssueId]; !ok {
			mapArr[item.IssueId] = []int64{item.RelationId}
		} else {
			mapArr[item.IssueId] = append(mapArr[item.IssueId], item.RelationId)
		}
	}

	return mapArr, nil
}

// GetBelongProjectRelationsByIssueId 获取一个任务所属的项目关联，被存在 issue_relation 表中。
func GetBelongProjectRelationsByIssueId(issueIds []int64) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	relateTypes := []int{consts.IssueRelationTypeBelongManyPro}

	relationInfos := &[]po.PpmPriIssueRelation{}
	cond := db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: db.In(relateTypes),
		consts.TcIssueId:      db.In(issueIds),
	}
	oriErr := mysql.SelectAllByCond(consts.TableIssueRelation, cond, relationInfos)
	if oriErr != nil {
		log.Error(oriErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	bos := &[]bo.IssueRelationBo{}
	copyErr := copyer.Copy(relationInfos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ObjectCopyError
	}

	return *bos, nil
}

// GetBelongProjectIdsByIssueId 获取一个任务所属的项目ids
// [warning] 截至 20211112，这种关联关系只开放一个关联关系，即，一个任务只会属于一个项目
func GetBelongProjectIdsByIssueId(issueId int64) ([]int64, errs.SystemErrorInfo) {
	resultProIds := make([]int64, 0)
	relations, err := GetBelongProjectRelationsByIssueId([]int64{issueId})
	if err != nil {
		return resultProIds, err
	}
	for _, item := range relations {
		resultProIds = append(resultProIds, item.ProjectId)
	}

	return resultProIds, nil
}

// GetIssueIdsByProject 获取一个项目下的所有的任务 id，被存在 issue_relation 表中。
func GetIssueIdsByProject(projectId int64) ([]int64, errs.SystemErrorInfo) {
	resultIds := make([]int64, 0)
	list, err := GetIssueProRelationsByProject(projectId)
	if err != nil {
		log.Error(err)
		return resultIds, err
	}
	for _, item := range list {
		resultIds = append(resultIds, item.IssueId)
	}
	return resultIds, nil
}

// GetIssueProRelationsByProject 获取任务和项目的关联关系集合
func GetIssueProRelationsByProject(projectId int64) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	resultList := make([]bo.IssueRelationBo, 0)
	relateTypes := []int{consts.IssueRelationTypeBelongManyPro}
	resultList, err = GetIssueRelationsByCondAndProId(projectId, relateTypes, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		log.Error(err)
		return resultList, err
	}

	return resultList, nil
}

// GetIssueProRelationsByProIds 获取任务和项目的关联关系集合
// 该逻辑已经不再支持。不再维护 IssueRelationTypeBelongManyPro 的关联关系
//func GetIssueProRelationsByProIds(projectIds []int64) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
//	relateTypes := []int{consts.IssueRelationTypeBelongManyPro}
//
//	relationInfos := &[]po.PpmPriIssueRelation{}
//	cond := db.Cond{
//		consts.TcIsDelete:     consts.AppIsNoDelete,
//		consts.TcRelationType: db.In(relateTypes),
//		consts.TcProjectId:    db.In(projectIds),
//	}
//	oriErr := mysql.SelectAllByCond(consts.TableIssueRelation, cond, relationInfos)
//	if oriErr != nil {
//		log.Error(oriErr)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
//	}
//	bos := &[]bo.IssueRelationBo{}
//	copyErr := copyer.Copy(relationInfos, bos)
//	if copyErr != nil {
//		log.Error(copyErr)
//		return nil, errs.ObjectCopyError
//	}
//
//	return *bos, nil
//}

// GetIssueIdsByProIds 通过项目id，获取项目下所有的任务 id
func GetIssueIdsByProIds(orgId int64, projectIds []int64) ([]int64, errs.SystemErrorInfo) {
	issueIds := make([]int64, 0)
	conds := []*tablePb.Condition{
		GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_in, nil, projectIds),
	}
	issueInfoDatas, err := GetIssueInfosMapLc(orgId, 0, &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: conds,
	}, []string{lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId)}, 0, 0)
	if err != nil {
		log.Errorf("[GetIssueIdsByProIds]GetIssueInfosMapLc err:%v, orgId:%v, projectIds:%v", err, orgId, projectIds)
		return nil, err
	}

	for _, data := range issueInfoDatas {
		issueId := cast.ToInt64(data[consts.BasicFieldIssueId])
		issueIds = append(issueIds, issueId)
	}

	return issueIds, nil
}

// GetIssueProRelationsByIssueIdAndProId 获取任务和项目的关联关系集合
func GetIssueProRelationsByIssueIdAndProId(projectId, issueId int64) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	resultList := make([]bo.IssueRelationBo, 0)
	relateTypes := []int{consts.IssueRelationTypeBelongManyPro}
	resultList, err = GetIssueRelationsByCondAndProId(projectId, relateTypes, db.Cond{
		consts.TcIssueId:  issueId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		log.Error(err)
		return resultList, err
	}

	return resultList, nil
}

func AddResourceRelation(orgId, userId, projectId, issueId int64, resourceIds []int64, sourceType int, columnId string) errs.SystemErrorInfo {
	resourceRelationResp := resourcefacade.AddResourceRelationWithType(resourcevo.AddResourceRelationReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &resourcevo.AddResourceRelationData{
			ProjectId:   projectId,
			IssueId:     issueId,
			ResourceIds: resourceIds,
			SourceType:  sourceType,
			ColumnId:    columnId,
		},
	})
	if resourceRelationResp.Failure() {
		log.Errorf("[AddResourceRelation] CreateResourceRelation err:%v, orgId:%v, userId:%v, projectId:%v, issueId:%v, resourceIds:%v",
			resourceRelationResp.Error(), orgId, userId, projectId, issueId, resourceIds)
		return resourceRelationResp.Error()
	}
	return nil
}

// GetResourceTypeByResourceIds 获取资源id和sourceType的对应关系
func GetResourceTypeByResourceIds(orgId, projectId int64, resourceIds []int64) (map[int64]int, errs.SystemErrorInfo) {
	relationList := resourcefacade.GetResourceRelationsByProjectId(resourcevo.GetResourceRelationsByProjectIdReqVo{
		UserId: 0,
		OrgId:  orgId,
		Input: resourcevo.GetResourceRelationsByProjectIdData{
			ProjectId:   projectId,
			SourceTypes: []int32{consts.OssPolicyTypeProjectResource, consts.OssPolicyTypeLesscodeResource},
		},
	})
	if relationList.Failure() {
		log.Errorf("[GetResourceTypeByResourceIds] err:%v, orgId:%d, projectId:%d, resourceIds:%v",
			relationList.Error(), orgId, projectId, resourceIds)
		return nil, relationList.Error()
	}
	relationMap := make(map[int64]resourcevo.ResourceRelationVo, len(relationList.ResourceRelations))
	for _, r := range relationList.ResourceRelations {
		relationMap[r.ResourceId] = r
	}

	// map[resourceId][sourceType]
	resourceIdType := map[int64]int{}
	for _, resourceId := range resourceIds {
		if _, ok := relationMap[resourceId]; ok {
			resourceIdType[resourceId] = relationMap[resourceId].SourceType
		}
	}

	return resourceIdType, nil
}

// 添加附件需要处理的资源关联关系
// 引用项目附件和飞书云文档需要手工处理资源关联关系
func AddAttachmentsRelations(orgId, userId, projectId, issueId int64, resourceIds []int64, columnId string) errs.SystemErrorInfo {
	if len(resourceIds) == 0 {
		return nil
	}
	err := AddResourceRelation(orgId, userId, projectId, issueId, resourceIds, consts.OssPolicyTypeLesscodeResource, columnId)
	if err != nil {
		log.Errorf("[AddAttachmentsRelations] AddResourceRelation err:%v", err)
		return err
	}
	return nil
}
