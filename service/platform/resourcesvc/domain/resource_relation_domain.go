package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/dao"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func DeleteResourceRelation(orgId, userId, projectId int64, issueIds []int64, resourceIds []int64, sourceTypes []int,
	recycleVersionId int64, columnId string, isDeleteResource bool, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	cond := db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
	}
	if sourceTypes != nil && len(sourceTypes) > 0 {
		cond[consts.TcSourceType] = db.In(sourceTypes)
	}
	if issueIds != nil && len(issueIds) > 0 {
		cond[consts.TcIssueId] = db.In(issueIds)
	}
	if resourceIds != nil {
		cond[consts.TcResourceId] = db.In(resourceIds)
	}
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  userId,
	}
	if recycleVersionId != 0 {
		upd[consts.TcVersion] = recycleVersionId
	}
	if columnId != "" {
		upd[consts.TcColumnId] = columnId
	}

	err := dao.UpdateResourceRelationByCond(cond, upd, tx...)
	if err != nil {
		log.Errorf("[DeleteResourceRelation] UpdateResourceRelationByCond err:%v", err)
		return err
	}
	if isDeleteResource {
		// 飞书云文档等网络资源不删除
		isDelete := consts.AppIsNoDelete
		resourceBoList, _, errSys := GetResourceBoList(0, 0, resourcevo.GetResourceBoListCond{
			OrgId:       orgId,
			ResourceIds: &resourceIds,
			IsDelete:    &isDelete,
		})
		if errSys != nil {
			log.Errorf("[DeleteResourceRelation] GetResourceBoList err:%v, orgId:%d, resourceIds:%v",
				errSys, orgId, resourceIds)
			return errSys
		}
		needDeleteIds := []int64{}
		for _, r := range *resourceBoList {
			if r.Type != consts.FsResource {
				needDeleteIds = append(needDeleteIds, r.Id)
			}
		}
		if len(needDeleteIds) > 0 {
			_, err3 := dao.UpdateResourceByCond(db.Cond{
				consts.TcIsDelete: consts.AppIsNoDelete,
				consts.TcOrgId:    orgId,
				consts.TcId:       db.In(needDeleteIds),
			}, mysql.Upd{
				consts.TcIsDelete: consts.AppIsDeleted,
				consts.TcUpdator:  userId,
			}, tx...)
			if err3 != nil {
				log.Errorf("[DeleteResourceRelation] UpdateResourceByCond err:%v", err3)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err3)
			}
		}
	}

	return nil
}

func UpdateResourceRelationProjectId(orgId, userId, projectId int64, issueIds []int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcIssueId:  db.In(issueIds),
	}
	err := dao.UpdateResourceRelationByCond(cond, mysql.Upd{
		consts.TcProjectId: projectId,
		consts.TcUpdator:   userId,
	}, tx...)
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}
	return nil
}

func RecoverResourceRelation(orgId, userId, projectId int64, resourceIds []int64, sourceType int, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableResourceRelation, db.Cond{
		consts.TcIsDelete:   consts.AppIsDeleted,
		consts.TcOrgId:      orgId,
		consts.TcProjectId:  projectId,
		consts.TcSourceType: sourceType,
		consts.TcResourceId: db.In(resourceIds),
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUpdator:  userId,
	})
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}

	return nil
}

func GetResourceRelationsByProjectId(orgId, userId, projectId int64, sourceTypes []int32) ([]resourcevo.ResourceRelationVo, errs.SystemErrorInfo) {
	var resourceRelations []po.PpmResResourceRelation
	err := mysql.SelectAllByCond((&po.PpmResResourceRelation{}).TableName(), db.Cond{
		consts.TcIsDelete:   db.Eq(consts.AppIsNoDelete),
		consts.TcOrgId:      orgId,
		consts.TcProjectId:  projectId,
		consts.TcSourceType: db.In(sourceTypes),
	}, &resourceRelations)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	var vos []resourcevo.ResourceRelationVo
	_ = copyer.Copy(&resourceRelations, &vos)
	return vos, nil
}

func GetResourceRelations(orgId, projectId int64, versionId, isDelete int, sourceTypes []int32, resourceIds []int64, isNeedResourceType bool) ([]resourcevo.ResourceRelationVo, errs.SystemErrorInfo) {
	resourceRelations := []po.PpmResResourceRelation{}
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: isDelete,
	}
	if projectId != 0 {
		cond[consts.TcProjectId] = projectId
	}
	if versionId != 0 {
		cond[consts.TcVersion] = versionId
	}
	if sourceTypes != nil && len(sourceTypes) > 0 {
		cond[consts.TcSourceType] = db.In(sourceTypes)
	}
	if resourceIds != nil && len(resourceIds) > 0 {
		cond[consts.TcResourceId] = db.In(resourceIds)
	}
	err := mysql.SelectAllByCond(consts.TableResourceRelation, cond, &resourceRelations)
	if err != nil {
		log.Errorf("[GetResourceRelations]SelectAllByCond err:%v", err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	resources := []po.PpmResResource{}
	if isNeedResourceType {
		err := mysql.SelectAllByCond(consts.TableResource, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
			consts.TcId:       db.In(resourceIds),
		}, &resources)
		if err != nil {
			log.Errorf("[GetResourceRelations]SelectAllByCond err:%v", err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}
	vos := []resourcevo.ResourceRelationVo{}
	err = copyer.Copy(&resourceRelations, &vos)
	if err != nil {
		log.Errorf("[GetResourceRelations] copy err:%v", err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	if len(resources) > 0 {
		resourceMap := map[int64]po.PpmResResource{}
		for _, rs := range resources {
			resourceMap[rs.Id] = rs
		}
		for i, r := range vos {
			vos[i].ResourceType = resourceMap[r.ResourceId].Type
		}
	}

	return vos, nil
}

func InsertResourceRelationWithColumnId(orgId, userId, projectId, issueId int64, resourceIds []int64, sourceType int, columnId string) errs.SystemErrorInfo {
	ids, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableResourceRelation, len(resourceIds))
	if idErr != nil {
		log.Error(idErr)
		return idErr
	}
	insertArr := []po.PpmResResourceRelation{}
	for i, id := range resourceIds {
		insertArr = append(insertArr, po.PpmResResourceRelation{
			Id:         ids.Ids[i].Id,
			OrgId:      orgId,
			ProjectId:  projectId,
			IssueId:    issueId,
			ResourceId: id,
			Creator:    userId,
			Updator:    userId,
			SourceType: sourceType,
			ColumnId:   columnId,
		})
	}

	insertErr := mysql.BatchInsert(&po.PpmResResourceRelation{}, slice.ToSlice(insertArr))
	if insertErr != nil {
		log.Error(insertErr)
		return errs.MysqlOperateError
	}
	return nil
}

func UpdateAttachmentRelation(orgId, userId, projectId int64, issueIds, resourceIds []int64, recycleVersion int64,
	columnId string, isDeleteResource bool) errs.SystemErrorInfo {
	pos := po.PpmResResourceRelation{}
	err := mysql.SelectOneByCond(consts.TableResourceRelation, db.Cond{
		consts.TcOrgId:      orgId,
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcProjectId:  projectId,
		consts.TcColumnId:   columnId,
		consts.TcResourceId: db.In(resourceIds),
	}, &pos)
	if err != nil && err != db.ErrNoMoreRows {
		log.Errorf("[UpdateAttachmentRelation] err:%v", err)
		return errs.MysqlOperateError
	}
	cond := db.Cond{
		consts.TcOrgId:      orgId,
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcProjectId:  projectId,
		consts.TcIssueId:    db.In(issueIds),
		consts.TcResourceId: db.In(resourceIds),
	}
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  userId,
	}
	if pos.ColumnId == "" {
		upd[consts.TcColumnId] = columnId
		upd[consts.TcVersion] = recycleVersion
	} else {
		cond[consts.TcColumnId] = columnId
		upd[consts.TcVersion] = recycleVersion
	}

	err = dao.UpdateResourceRelationByCond(cond, upd)
	if err != nil {
		log.Errorf("[UpdateAttachmentRelation] UpdateResourceRelationByCond err:%v", err)
		return errs.MysqlOperateError
	}
	if isDeleteResource {
		// 飞书云文档等网络资源不删除
		isDelete := consts.AppIsNoDelete
		resourceBoList, _, errSys := GetResourceBoList(0, 0, resourcevo.GetResourceBoListCond{
			OrgId:       orgId,
			ResourceIds: &resourceIds,
			IsDelete:    &isDelete,
		})
		if errSys != nil {
			log.Errorf("[DeleteResourceRelation] GetResourceBoList err:%v, orgId:%d, resourceIds:%v",
				errSys, orgId, resourceIds)
			return errSys
		}
		needDeleteIds := []int64{}
		for _, r := range *resourceBoList {
			if r.SourceType != consts.OssPolicyTypeProjectResource && r.Type != consts.FsResource {
				needDeleteIds = append(needDeleteIds, r.Id)
			}
		}
		if len(needDeleteIds) > 0 {
			_, err3 := dao.UpdateResourceByCond(db.Cond{
				consts.TcIsDelete: consts.AppIsNoDelete,
				consts.TcOrgId:    orgId,
				consts.TcId:       db.In(needDeleteIds),
			}, mysql.Upd{
				consts.TcIsDelete: consts.AppIsDeleted,
				consts.TcUpdator:  userId,
			})
			if err3 != nil {
				log.Errorf("[DeleteResourceRelation] UpdateResourceByCond err:%v", err3)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err3)
			}
		}
	}
	return nil
}
