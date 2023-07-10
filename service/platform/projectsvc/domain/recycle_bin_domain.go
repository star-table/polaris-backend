package domain

import (
	"fmt"

	tablePb "github.com/star-table/interface/golang/table/v1"

	"github.com/spf13/cast"

	"github.com/star-table/common/core/util/slice"
	slice2 "github.com/star-table/polaris-backend/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/facade/formfacade"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetRecycleInfo(orgId, recycleId, relationId int64, relationType int) (*bo.PpmPrsRecycleBin, errs.SystemErrorInfo) {
	recycleInfo := po.PpmPrsRecycleBin{}
	err := mysql.SelectOneByCond(consts.TableRecycleBin, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: relationType,
		consts.TcRelationId:   relationId,
		consts.TcOrgId:        orgId,
		consts.TcId:           recycleId,
	}, &recycleInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.RecycleObjectNotExist
		} else {
			log.Errorf("[GetRecycleInfo] err:%v, orgId:%d, recycleId:%d", err, orgId, recycleId)
			return nil, errs.MysqlOperateError
		}
	}
	res := &bo.PpmPrsRecycleBin{}
	errCopy := copyer.Copy(&recycleInfo, res)
	if errCopy != nil {
		log.Errorf("[GetRecycleInfo] copy err:%v", errCopy)
		return nil, errs.ObjectCopyError
	}
	return res, nil
}

func GetRecoverRecord(orgId, relationId int64, relationType int) (*bo.PpmPrsRecycleBin, errs.SystemErrorInfo) {
	info := &po.PpmPrsRecycleBin{}
	err := mysql.SelectOneByCond(consts.TableRecycleBin, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: relationType,
		consts.TcRelationId:   relationId,
		consts.TcOrgId:        orgId,
	}, info)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.RecycleObjectNotExist
		} else {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}
	}

	res := &bo.PpmPrsRecycleBin{}
	_ = copyer.Copy(info, res)
	return res, nil
}

// 和AddRecycleRecord方法效果一样
func AddRecycleForProjectResource(orgId, userId int64, projectId int64, relationIds []int64, resourceTypeMap map[int64]int, tx ...sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	if len(relationIds) == 0 {
		return 0, nil
	}
	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableRecycleBin, len(relationIds))
	if err != nil {
		log.Error(err)
		return 0, err
	}
	versionId, err := idfacade.ApplyPrimaryIdRelaxed(consts.RecycleVersion)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	pos := []po.PpmPrsRecycleBin{}
	for i, id := range relationIds {
		pos = append(pos, po.PpmPrsRecycleBin{
			Id:           ids.Ids[i].Id,
			OrgId:        orgId,
			ProjectId:    projectId,
			RelationId:   id,
			RelationType: resourceTypeMap[id],
			Creator:      userId,
			Version:      int(versionId),
		})
	}

	insertErr := dao.InsertRecycleBinBatch(pos, tx...)
	if insertErr != nil {
		log.Error(insertErr)
		return 0, errs.MysqlOperateError
	}
	return versionId, nil
}

func AddRecycleRecord(orgId, userId int64, projectId int64, relationIds []int64, relationType int, tx ...sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	if len(relationIds) == 0 {
		return 0, nil
	}
	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableRecycleBin, len(relationIds))
	if err != nil {
		log.Error(err)
		return 0, err
	}

	versionId, err := idfacade.ApplyPrimaryIdRelaxed(consts.RecycleVersion)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	pos := []po.PpmPrsRecycleBin{}
	for i, id := range relationIds {
		pos = append(pos, po.PpmPrsRecycleBin{
			Id:           ids.Ids[i].Id,
			OrgId:        orgId,
			ProjectId:    projectId,
			RelationId:   id,
			RelationType: relationType,
			Creator:      userId,
			Version:      int(versionId),
		})
	}

	insertErr := dao.InsertRecycleBinBatch(pos, tx...)
	if insertErr != nil {
		log.Error(insertErr)
		return 0, errs.MysqlOperateError
	}

	return versionId, nil
}

func DeleteRecycleBin(id, userId int64) errs.SystemErrorInfo {
	_, err := mysql.UpdateSmartWithCond(consts.TableRecycleBin, db.Cond{
		consts.TcId: id,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  userId,
	})

	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}

	return nil
}

func RecoverIssue(orgId, userId, relationId int64, recycleId int64, sourceChannel string) ([]int64, int64, errs.SystemErrorInfo) {
	//查询任务是否存在且已被删除
	//info := &po.PpmPriIssue{}
	//err := mysql.SelectOneByCond(consts.TableIssue, db.Cond{
	//	consts.TcId:       relationId,
	//	consts.TcOrgId:    orgId,
	//	consts.TcIsDelete: consts.AppIsDeleted,
	//}, info)
	//if err != nil {
	//	if err == db.ErrNoMoreRows {
	//		return nil, 0, errs.RecycleObjectNotExist
	//	} else {
	//		log.Error(err)
	//		return nil, 0, errs.MysqlOperateError
	//	}
	//}
	condition := &tablePb.Condition{Type: tablePb.ConditionType_and}
	condition.Conditions = []*tablePb.Condition{
		GetRowsCondition(consts.BasicFieldIssueId, tablePb.ConditionType_equal, relationId, nil),
		GetRowsCondition(consts.BasicFieldRecycleFlag, tablePb.ConditionType_equal, consts.AppIsDeleted, nil),
	}
	lcIssueInfos, errSys := GetIssueInfosMapLc(orgId, userId, condition, nil, -1, -1)
	if errSys != nil {
		log.Errorf("[RecoverIssue] GetIssueInfosMapLc err:%v, orgId:%v, issueId:%v", errSys, orgId, relationId)
		return nil, 0, errSys
	}
	if len(lcIssueInfos) < 1 {
		return nil, 0, errs.RecycleObjectNotExist
	}
	info, errSys := ConvertIssueDataToIssueBo(lcIssueInfos[0])
	if errSys != nil {
		log.Errorf("[RecoverIssue] GetIssueInfosMapLc err:%v, orgId:%v, issueId:%v", errSys, orgId, relationId)
		return nil, 0, errSys
	}

	//查询表是否存在，以及表下面的状态是否存在（表不存在，就放到空项目中，状态不存在，就放到默认第一个状态中）
	needChangeProject := false
	needChangeStatus := false
	appId := int64(0)
	if info.ProjectId != 0 {
		projectInfo, projectInfoErr := GetProjectSimple(orgId, info.ProjectId)
		if projectInfoErr != nil {
			log.Errorf("[RecoverIssue] GetProject failed:%d", projectInfoErr)
			return nil, 0, projectInfoErr
		}
		appId = projectInfo.AppId
		_, tableInfoErr := GetTableInfo(orgId, projectInfo.AppId, info.TableId)
		if tableInfoErr != nil {
			if tableInfoErr == errs.TableNotExist {
				needChangeProject = true
			} else {
				log.Errorf("[RecoverIssue] GetTableInfo failed:%d", tableInfoErr)
				return nil, 0, tableInfoErr
			}
		}
	}
	tableId := info.TableId
	if needChangeProject {
		tableId = 0
		summaryId, summaryIdErr := GetAppIdFromProjectId(orgId, 0)
		if summaryIdErr != nil {
			log.Errorf("[RecoverIssue] GetAppIdFromProjectId failed:%d", summaryIdErr)
			return nil, 0, summaryIdErr
		}
		appId = summaryId
	}
	statusIds := []int64{}
	if info.Status != 0 {
		statusList, statusListErr := GetTableStatus(orgId, tableId)
		if statusListErr != nil {
			log.Errorf("[RecoverIssue] GetTableStatus failed:%d", statusListErr)
			return nil, 0, statusListErr
		}
		for _, infoBo := range statusList {
			statusIds = append(statusIds, infoBo.ID)
		}
		if ok, _ := slice.Contain(statusIds, info.Status); !ok {
			needChangeStatus = true
		}
	}

	//children := &[]po.PpmPriIssue{}
	//err1 := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
	//	consts.TcPath:     db.Like(fmt.Sprintf("%s%d,%s", info.Path, relationId, "%")),
	//	consts.TcOrgId:    orgId,
	//	consts.TcVersion:  recycleId,
	//	consts.TcIsDelete: consts.AppIsDeleted,
	//}, children)
	//if err1 != nil {
	//	log.Error(err1)
	//	return nil, 0, errs.MysqlOperateError
	//}

	conds := []*tablePb.Condition{
		GetRowsCondition(consts.BasicFieldPath, tablePb.ConditionType_like, fmt.Sprintf(",%d,", relationId), nil),
		GetRowsCondition(consts.BasicFieldRecycleFlag, tablePb.ConditionType_equal, consts.AppIsDeleted, nil),
	}
	idList, errSys := GetIssueIdList(orgId, conds, -1, -1)
	if errSys != nil {
		log.Errorf("[RecoverIssue] GetIssueIdList err:%v, orgId:%v, issueId:%v", errSys, orgId, relationId)
		return nil, 0, errSys
	}

	allIssueIds := []int64{relationId}
	if len(idList) > 0 {
		for _, child := range idList {
			allIssueIds = append(allIssueIds, child.Id)
		}
	}
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//恢复任务
		upd := mysql.Upd{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcUpdator:  userId,
		}
		if needChangeProject {
			upd[consts.TcProjectId] = 0
			upd[consts.TcTableId] = 0
		}
		if needChangeStatus {
			upd[consts.TcStatus] = statusIds[0]
		}
		//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcId: db.In(allIssueIds),
		//}, upd)
		//if err != nil {
		//	log.Error(err)
		//	return err
		//}
		//无码
		updForm := []map[string]interface{}{}
		baseUpd := slice2.CaseCamelCopy(upd)
		if status, ok := baseUpd["status"]; ok {
			baseUpd["issueStatus"] = status
			delete(baseUpd, "status")
		}
		if needChangeProject {
			baseUpd["appId"] = "0"
		}
		for _, id := range allIssueIds {
			tmp := baseUpd
			tmp["issueId"] = id
			updForm = append(updForm, tmp)
		}
		//更新到无码
		resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
			AppId:   appId,
			OrgId:   orgId,
			UserId:  userId,
			TableId: tableId,
			Form:    updForm,
		})
		if resp.Failure() {
			log.Error(resp.Error())
			return resp.Error()
		}

		//恢复任务详情
		//_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueDetail, db.Cond{
		//	consts.TcOrgId:   orgId,
		//	consts.TcIssueId: db.In(allIssueIds),
		//}, mysql.Upd{
		//	consts.TcIsDelete: consts.AppIsNoDelete,
		//	consts.TcUpdator:  userId,
		//})
		//if err1 != nil {
		//	log.Error(err1)
		//	return err1
		//}

		//恢复任务关联
		// 只恢复一部分的任务关联实体，因为有些不能恢复，比如日程的关联。
		// 恢复日程关联，而实际上对应的日程已在飞书方删除了。
		noIssueRelationTypes := []int{
			consts.IssueRelationTypeCalendar,
		}
		_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcIssueId:      db.In(allIssueIds),
			consts.TcVersion:      recycleId,
			consts.TcRelationType: db.NotIn(noIssueRelationTypes),
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcUpdator:  userId,
		})
		if err2 != nil {
			log.Error(err2)
			return err2
		}

		// 恢复任务对应的工时记录
		if err := RecoverWorkHours(orgId, allIssueIds, recycleId); err != nil {
			log.Errorf("[RecoverIssue] orgId: %d, err: %v", orgId, err)
			return err
		}

		return nil
	})

	if transErr != nil {
		log.Error(transErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}

	//恢复文件关联
	resp := resourcefacade.RecoverResource(resourcevo.RecoverResourceReqVo{
		OrgId:  orgId,
		UserId: userId,
		Input: resourcevo.RecoverResourceData{
			ProjectId:        info.ProjectId,
			IssueIds:         allIssueIds,
			RecycleVersionId: recycleId,
		},
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, 0, resp.Error()
	}

	//动态
	asyn.Execute(func() {
		blank := []int64{}
		issueTrendsBo := &bo.IssueTrendsBo{
			PushType:      consts.PushTypeRecoverIssue,
			OrgId:         orgId,
			OperatorId:    userId,
			DataId:        info.DataId,
			IssueId:       relationId,
			ParentIssueId: info.ParentId,
			ProjectId:     info.ProjectId,
			PriorityId:    info.PriorityId,
			ParentId:      info.ParentId,
			IssueTitle:    info.Title,
			IssueStatusId: info.Status,
			BeforeOwner:   info.OwnerIdI64,
			AfterOwner:    blank,
			SourceChannel: sourceChannel,
			TableId:       info.TableId,
		}
		asyn.Execute(func() {
			PushIssueTrends(issueTrendsBo)
		})
		asyn.Execute(func() {
			PushIssueThirdPlatformNotice(issueTrendsBo)
		})
	})
	return allIssueIds, tableId, nil
}

func RecoverFolder(orgId, userId, projectId, relationId int64, recycleId int64) errs.SystemErrorInfo {
	resp := resourcefacade.RecoverFolder(resourcevo.RecoverFolderReqVo{
		OrgId:  orgId,
		UserId: userId,
		Input: resourcevo.RecoverFolderData{
			ProjectId: projectId,
			FolderId:  relationId,
			RecycleId: recycleId,
		},
	})

	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}

	//动态
	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{ObjName: resp.Data.Name, FolderId: relationId}
		projectTrendsBo := bo.ProjectTrendsBo{
			PushType:   consts.PushTypeRecoverFolder,
			OrgId:      orgId,
			ProjectId:  resp.Data.ProjectId,
			OperatorId: userId,
			Ext:        ext,
		}
		PushProjectTrends(projectTrendsBo)
	})

	return nil
}

func RecoverResource(orgId, userId, projectId, relationId int64, recycleId int64) errs.SystemErrorInfo {
	resp := resourcefacade.RecoverResource(resourcevo.RecoverResourceReqVo{
		OrgId:  orgId,
		UserId: userId,
		Input: resourcevo.RecoverResourceData{
			ProjectId:        projectId,
			ResourceId:       relationId,
			RecycleVersionId: recycleId,
		},
	})

	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}

	//动态
	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		resourceInfo := bo.ResourceInfoBo{
			Url:        resp.Data.Host + resp.Data.Path,
			Name:       resp.Data.Name,
			Size:       resp.Data.Size,
			UploadTime: resp.Data.CreateTime,
			Suffix:     resp.Data.Suffix,
		}
		ext.ResourceInfo = append(ext.ResourceInfo, resourceInfo)
		projectTrendsBo := bo.ProjectTrendsBo{
			PushType:   consts.PushTypeRecoverProjectFile,
			OrgId:      orgId,
			ProjectId:  resp.Data.ProjectId,
			OperatorId: userId,
			Ext:        ext,
		}
		PushProjectTrends(projectTrendsBo)
	})

	return nil
}

func RecoverAttachment(orgId, userId, projectId, relationId int64, recycleId int64) errs.SystemErrorInfo {
	projectInfo, err := GetProjectSimple(orgId, projectId)
	if err != nil {
		return err
	}
	resp := resourcefacade.RecoverResource(resourcevo.RecoverResourceReqVo{
		OrgId:  orgId,
		UserId: userId,
		Input: resourcevo.RecoverResourceData{
			ProjectId:        projectId,
			ResourceId:       relationId,
			RecycleVersionId: recycleId,
			AppId:            projectInfo.AppId,
		},
	})

	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}

	//恢复任务附件关联
	_, err1 := mysql.UpdateSmartWithCond(consts.TableIssueRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsDeleted,
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationId:   relationId,
		consts.TcRelationType: consts.IssueRelationTypeResource,
		consts.TcVersion:      recycleId,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUpdator:  userId,
	})

	if err1 != nil {
		log.Error(err1)
		return errs.MysqlOperateError
	}

	//动态
	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		resourceInfo := bo.ResourceInfoBo{
			Url:        resp.Data.Host + resp.Data.Path,
			Name:       resp.Data.Name,
			Size:       resp.Data.Size,
			UploadTime: resp.Data.CreateTime,
			Suffix:     resp.Data.Suffix,
		}
		ext.ResourceInfo = append(ext.ResourceInfo, resourceInfo)
		projectTrendsBo := bo.ProjectTrendsBo{
			PushType:   consts.PushTypeRecoverProjectAttachment,
			OrgId:      orgId,
			ProjectId:  resp.Data.ProjectId,
			OperatorId: userId,
			Ext:        ext,
		}
		PushProjectTrends(projectTrendsBo)
	})

	return nil
}

// 从回收站恢复无码附件
func RecoverLessAttachments(orgId, userId, projectId, resourceId int64, recycleId int64, sourceChannel string) errs.SystemErrorInfo {
	projectInfo, err := GetProjectSimple(orgId, projectId)
	if err != nil {
		return err
	}

	resp := resourcefacade.RecoverResource(resourcevo.RecoverResourceReqVo{
		OrgId:         orgId,
		UserId:        userId,
		SourceChannel: sourceChannel,
		Input: resourcevo.RecoverResourceData{
			ProjectId:        projectId,
			ResourceId:       resourceId,
			RecycleVersionId: recycleId,
			AppId:            projectInfo.AppId,
		},
	})

	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}

	//恢复任务附件关联
	//_, err1 := mysql.UpdateSmartWithCond(consts.TableIssueRelation, db.Cond{
	//	consts.TcIsDelete:     consts.AppIsDeleted,
	//	consts.TcOrgId:        orgId,
	//	consts.TcProjectId:    projectId,
	//	consts.TcRelationId:   resourceId,
	//	consts.TcRelationType: consts.IssueRelationTypeResource,
	//}, mysql.Upd{
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//	consts.TcUpdator:  userId,
	//})
	//
	//if err1 != nil {
	//	log.Error(err1)
	//	return errs.MysqlOperateError
	//}

	return nil
}

func GetRecycleList(orgId, projectId int64, relationType int, page, size int) (uint64, []bo.PpmPrsRecycleBin, errs.SystemErrorInfo) {
	info := &[]po.PpmPrsRecycleBin{}
	cond := db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
	}
	if relationType != 0 {
		cond[consts.TcRelationType] = relationType
	}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableRecycleBin, cond, nil, page, size, "create_time desc", info)
	if err != nil {
		log.Error(err)
		return 0, nil, errs.MysqlOperateError
	}

	bos := &[]bo.PpmPrsRecycleBin{}
	_ = copyer.Copy(info, bos)

	return total, *bos, nil
}

// GetRecycleListWithAuth 只获取当前登陆者有权访问的资源。不相关的数据不展示
func GetRecycleListWithAuth(orgId, currentUserId, projectId int64, relationType int, page, size int) (uint64, []bo.PpmPrsRecycleBin, errs.SystemErrorInfo) {
	info := &[]po.PpmPrsRecycleBin{}
	cond := db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
	}
	if relationType != 0 {
		cond[consts.TcRelationType] = relationType
	}

	total, oriErr := mysql.SelectAllByCondWithPageAndOrder(consts.TableRecycleBin, cond, nil, page, size, "create_time desc", info)
	if oriErr != nil {
		log.Error(oriErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}

	bos := &[]bo.PpmPrsRecycleBin{}
	_ = copyer.Copy(info, bos)

	return total, *bos, nil
}

//// 从回收站恢复附件资源
//func RecoverAttachmentsResource(orgId, userId, projectId, resourceId, recycleId int64) errs.SystemErrorInfo {
//	// 恢复无码的recycleFlag为1
//	// 恢复任务资源关联表
//
//	// 恢复资源关联表
//	_, err := mysql.UpdateSmartWithCond(consts.TableIssueRelation, db.Cond{
//		consts.TcIsDelete:     consts.AppIsDeleted,
//		consts.TcOrgId:        orgId,
//		consts.TcProjectId:    projectId,
//		consts.TcRelationId:   resourceId,
//		consts.TcRelationType: consts.IssueRelationTypeResource,
//		consts.TcVersion:      recycleId,
//	}, mysql.Upd{
//		consts.TcIsDelete: consts.AppIsNoDelete,
//		consts.TcUpdator:  userId,
//	})
//	if err != nil {
//		log.Errorf("[RecoverAttachmentsResource] update TableIssueRelation err:%v", err)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	return nil
//}

// DeleteAttachmentsForOneTable 删除一个table，把所有附件都删除，不进入回收站
func DeleteAttachmentsForOneTable(orgId, userId, projectId, tableId int64) errs.SystemErrorInfo {
	// 查询是否有附件字段
	tableColumn, errSys := GetTableColumnConfig(orgId, tableId, nil, false)
	if errSys != nil {
		log.Errorf("[DeleteAttachmentsForOneTable] GetTableColumnConfig error:%v, orgId:%d, tableId:%d",
			errSys, orgId, tableId)
		return errSys
	}
	columnIds := []string{}
	for _, column := range tableColumn.Columns {
		if column.Field.Type == consts.BasicFieldDocument {
			// 处理附件字段
			columnIds = append(columnIds, column.Name)
		}
	}

	if len(columnIds) < 1 {
		return nil
	}

	condition := GetRowsCondition(consts.BasicFieldTableId, tablePb.ConditionType_equal, cast.ToString(tableId), nil)
	filterColumns := []string{}
	for _, col := range columnIds {
		filterColumns = append(filterColumns, lc_helper.ConvertToFilterColumn(col))
	}
	infosMapLc, err := GetIssueInfosMapLc(orgId, userId, condition, filterColumns, 1, 2000)
	if err != nil {
		log.Errorf("[DeleteAttachmentsForOneTable] GetIssueInfosMapLc err:%v, orgId:%d, userId:%d, projectId:%d, tableId:%d",
			err, orgId, userId, projectId, tableId)
		return err
	}

	//attachmentsMap := map[string]map[string]bo.Attachments{}
	resourceIds := []int64{}
	issueIds := []int64{}
	for _, issueMap := range infosMapLc {
		for _, columnId := range columnIds {
			attach := map[string]bo.Attachments{}
			if documentM, ok := issueMap[columnId]; ok {
				copyer.Copy(documentM, &attach)
			}
			for _, r := range attach {
				if r.Id != 0 {
					issueId := cast.ToInt64(issueMap[consts.BasicFieldIssueId])
					issueIds = append(issueIds, issueId)
				}
				resourceIds = append(resourceIds, r.Id)
			}
		}
	}
	resourceIds = slice.SliceUniqueInt64(resourceIds)
	issueIds = slice.SliceUniqueInt64(issueIds)

	if len(resourceIds) < 1 {
		// 附件资源没有就不处理了
		return nil
	}

	// 查询该附件列中在回收站的文件过滤掉
	relationList := resourcefacade.GetResourceRelationList(resourcevo.GetResourceRelationListReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &resourcevo.GetResourceRelationList{
			ProjectId:   projectId,
			ResourceIds: resourceIds,
			IsDelete:    consts.AppIsDeleted,
		},
	})
	if relationList.Failure() {
		log.Errorf("[DeleteAttachmentsForColumn] err:%v, orgId:%d, userId:%d, projectId:%d, resourceIds:%v",
			relationList.Error(), orgId, userId, projectId, resourceIds)
		return relationList.Error()
	}
	needDeleteResourceIds := []int64{}
	for _, rr := range relationList.Data {
		// 在附件字段上传的附件可以在回收站删除
		if rr.ColumnId != "" && rr.SourceType == consts.OssPolicyTypeLesscodeResource {
			needDeleteResourceIds = append(needDeleteResourceIds, rr.ResourceId)
		}
	}

	if len(needDeleteResourceIds) < 1 {
		return nil
	}

	err2 := deleteResourceRelation(orgId, userId, projectId, issueIds, resourceIds)
	if err2 != nil {
		log.Errorf("[DeleteAttachmentsForOneTable] err:%v", err2)
		return err2
	}

	return nil
}

// DeleteAttachmentsForColumn 删除附件字段列，把涉及到的附件资源全部删除 不进入回收站
func DeleteAttachmentsForColumn(orgId, userId, projectId, tableId int64, columnId string) errs.SystemErrorInfo {
	columnsMap, err := GetTableColumnsMap(orgId, tableId, []string{columnId})
	if err != nil {
		log.Errorf("[DeleteAttachmentsForColumn] GetTableColumnsMap err:%v, orgId:%d, projectId:%d, tableId:%d",
			err, orgId, projectId, tableId)
		return err
	}

	if columnsMap[columnId].Field.Type != consts.LcColumnFieldTypeDocument {
		return nil
	}

	condition := &tablePb.Condition{Type: tablePb.ConditionType_and}
	condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldTableId, tablePb.ConditionType_equal, cast.ToString(tableId), nil))
	columns := []string{
		lc_helper.ConvertToFilterColumn(columnId),
	}
	infosMapLc, err := GetIssueInfosMapLc(orgId, userId, condition, columns, 1, 2000)
	if err != nil {
		log.Errorf("[DeleteAttachmentsForColumn] GetIssueInfosMapLc err:%v, orgId:%d, userId:%d, projectId:%d, tableId:%d",
			err, orgId, userId, projectId, tableId)
		return err
	}
	if len(infosMapLc) < 1 {
		log.Errorf("[DeleteAttachmentsForColumn] deleteIssue not found, orgId:%d, userId:%d, appId:%d, projectId:%d",
			orgId, userId, projectId, tableId)
		return nil
	}
	issueIds := []int64{}
	resourceIds := []int64{}
	for _, issue := range infosMapLc {
		attachments := map[string]bo.Attachments{}
		if attachM, ok := issue[columnId]; ok {
			copyer.Copy(attachM, &attachments)
		}
		for _, r := range attachments {
			if r.Id != 0 {
				issueId := cast.ToInt64(issue[consts.BasicFieldIssueId])
				issueIds = append(issueIds, issueId)
			}
			resourceIds = append(resourceIds, r.Id)
		}
	}

	issueIds = slice.SliceUniqueInt64(issueIds)
	resourceIds = slice.SliceUniqueInt64(resourceIds)

	if len(resourceIds) < 1 {
		// 附件没有数据就跳过了
		return nil
	}

	// 查询该附件列中在回收站的文件过滤掉
	relationList := resourcefacade.GetResourceRelationList(resourcevo.GetResourceRelationListReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &resourcevo.GetResourceRelationList{
			ProjectId:   projectId,
			ResourceIds: resourceIds,
			IsDelete:    consts.AppIsDeleted,
		},
	})
	if relationList.Failure() {
		log.Errorf("[DeleteAttachmentsForColumn] err:%v, orgId:%d, userId:%d, projectId:%d, resourceIds:%v",
			relationList.Error(), orgId, userId, projectId, resourceIds)
		return relationList.Error()
	}
	needDeleteResourceIds := []int64{}
	for _, rr := range relationList.Data {
		// 在附件字段上传的附件可以在回收站删除
		if rr.ColumnId != "" && rr.SourceType == consts.OssPolicyTypeLesscodeResource {
			needDeleteResourceIds = append(needDeleteResourceIds, rr.ResourceId)
		}
	}

	if len(needDeleteResourceIds) < 1 {
		return nil
	}

	err2 := deleteResourceRelation(orgId, userId, projectId, issueIds, needDeleteResourceIds)
	if err2 != nil {
		log.Error(err2)
		return err2
	}

	return nil
}

func deleteIssueResourceRelation(orgId, userId, projectId int64, issueIds, resourceIds []int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:      orgId,
		consts.TcProjectId:  projectId,
		consts.TcIssueId:    db.In(issueIds),
		consts.TcRelationId: db.In(resourceIds),
		consts.TcIsDelete:   consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  userId,
	})
	if err != nil {
		log.Errorf("[deleteIssueResourceRelation] err:%v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func deleteRecycleResource(orgId, userId, projectId int64, resourceIds []int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableRecycleBin, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: consts.RecycleTypeAttachment,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationId:   db.In(resourceIds),
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  userId,
	})
	if err != nil {
		log.Errorf("[deleteRecycleResource]TransUpdateSmartWithCond err:%v, orgId:%d, projectId:%d, resourceIds:%v",
			err, orgId, projectId, resourceIds)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func deleteResourceRelation(orgId, userId, projectId int64, issueIds, resourceIds []int64) errs.SystemErrorInfo {
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(issueIds) > 0 && len(resourceIds) > 0 {
			// 删除资源和任务的关联关系
			//err2 := deleteIssueResourceRelation(orgId, userId, projectId, issueIds, resourceIds, tx)
			//if err2 != nil {
			//	return err2
			//}
			// 删除资源表的关联关系
			resp := resourcefacade.DeleteAttachmentRelation(resourcevo.DeleteAttachmentRelationReq{
				OrgId:  orgId,
				UserId: userId,
				Input: &resourcevo.DeleteAttachmentRelationData{
					ProjectId:        projectId,
					IssueIds:         issueIds,
					ResourceIds:      resourceIds,
					IsDeleteResource: true,
				},
			})
			if resp.Failure() {
				log.Errorf("[deleteResourceRelation] DeleteAttachmentRelation err:%v, orgId:%d, userId:%d, projectId:%d, resourceIds:%v",
					resp.Error(), orgId, userId, projectId, resourceIds)
				return resp.Error()
			}
			// 删除回收站的附件
			err2 := deleteRecycleResource(orgId, userId, projectId, resourceIds, tx)
			if err2 != nil {
				return err2
			}
		}
		return nil
	})
	if err != nil {
		log.Errorf("[deleteResourceRelation] err:%v", err)
		return errs.MysqlOperateError
	}

	return nil
}
