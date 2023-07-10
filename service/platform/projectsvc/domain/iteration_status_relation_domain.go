package domain

import (
	"fmt"
	"time"

	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/date"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetIterationRelationStatusInfo(orgId int64, iterationId int64) ([]bo.PpmPriIterationStatusRelationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmPriIterationStatusRelation{}
	err := mysql.SelectAllByCond(consts.TableIterationStatusRelation, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcOrgId:       orgId,
		consts.TcIterationId: iterationId,
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.PpmPriIterationStatusRelationBo{}
	copyerErr := copyer.Copy(pos, bos)
	if copyerErr != nil {
		log.Error(copyerErr)
		return nil, errs.ObjectCopyError
	}

	return *bos, nil
}

func GetIterationRelationStatusInfoBatch(orgId int64, iterationIds []int64) ([]bo.PpmPriIterationStatusRelationBo, errs.SystemErrorInfo) {
	pos := []po.PpmPriIterationStatusRelation{}
	err := mysql.SelectAllByCond(consts.TableIterationStatusRelation, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcOrgId:       orgId,
		consts.TcIterationId: db.In(iterationIds),
	}, &pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := []bo.PpmPriIterationStatusRelationBo{}
	copyerErr := copyer.Copy(pos, &bos)
	if copyerErr != nil {
		log.Error(copyerErr)
		return nil, errs.ObjectCopyError
	}

	return bos, nil
}

//有则更新，无则添加
func UpdateOrInsertIterationRelation(orgId int64, projectId int64, iterationId int64, statusId int64, planStartTime, planEndTime, startTime, endTime *types.Time, operatorId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	lockKey := fmt.Sprintf("%s:%d", consts.AddIterationRelationLock, iterationId)
	uid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uid)
	if lockErr != nil {
		log.Error(lockErr)
		return errs.CreateIterationRelationErr
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	} else {
		return errs.CreateIterationRelationErr
	}

	info := &po.PpmPriIterationStatusRelation{}
	var err error
	if tx != nil && len(tx) > 0 {
		err = mysql.TransSelectOneByCond(tx[0], consts.TableIterationStatusRelation, db.Cond{
			consts.TcIsDelete:    consts.AppIsNoDelete,
			consts.TcOrgId:       orgId,
			consts.TcIterationId: iterationId,
			consts.TcProjectId:   projectId,
			consts.TcStatusId:    statusId,
		}, info)
	} else {
		err = mysql.SelectOneByCond(consts.TableIterationStatusRelation, db.Cond{
			consts.TcIsDelete:    consts.AppIsNoDelete,
			consts.TcOrgId:       orgId,
			consts.TcIterationId: iterationId,
			consts.TcProjectId:   projectId,
			consts.TcStatusId:    statusId,
		}, info)
	}

	if err != nil {
		if err != db.ErrNoMoreRows {
			log.Error(err)
		} else {
			//没有查询到，新增
			id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIterationStatusRelation)
			if err != nil {
				return err
			}
			insertPo := &po.PpmPriIterationStatusRelation{
				Id:          id,
				OrgId:       orgId,
				ProjectId:   projectId,
				IterationId: iterationId,
				StatusId:    statusId,
				Creator:     operatorId,
			}
			if planStartTime != nil {
				formatTime, _ := time.ParseInLocation(consts.AppTimeFormat, planStartTime.String(), time.Local)
				insertPo.PlanStartTime = formatTime
			}
			if planEndTime != nil {
				formatTime, _ := time.ParseInLocation(consts.AppTimeFormat, planEndTime.String(), time.Local)
				insertPo.PlanEndTime = formatTime
			}
			if startTime != nil {
				formatTime, _ := time.ParseInLocation(consts.AppTimeFormat, startTime.String(), time.Local)
				insertPo.StartTime = formatTime
			}
			if endTime != nil {
				formatTime, _ := time.ParseInLocation(consts.AppTimeFormat, endTime.String(), time.Local)
				insertPo.EndTime = formatTime
			}

			var insertErr error
			if tx != nil && len(tx) > 0 {
				insertErr = mysql.TransInsert(tx[0], insertPo)
			} else {
				insertErr = mysql.Insert(insertPo)
			}
			if insertErr != nil {
				log.Error(insertErr)
				return errs.MysqlOperateError
			}
		}
	}

	upd := mysql.Upd{}
	if planStartTime != nil {
		upd[consts.TcPlanStartTime] = date.FormatTime(*planStartTime)
	}
	if planEndTime != nil {
		upd[consts.TcPlanEndTime] = date.FormatTime(*planEndTime)
	}
	if startTime != nil {
		upd[consts.TcStartTime] = date.FormatTime(*startTime)
	}
	if endTime != nil {
		upd[consts.TcEndTime] = date.FormatTime(*endTime)
	}
	var updateErr error
	if tx != nil && len(tx) > 0 {
		_, updateErr = mysql.TransUpdateSmartWithCond(tx[0], consts.TableIterationStatusRelation, db.Cond{
			consts.TcId: info.Id,
		}, upd)
	} else {
		_, updateErr = mysql.UpdateSmartWithCond(consts.TableIterationStatusRelation, db.Cond{
			consts.TcId: info.Id,
		}, upd)
	}

	if updateErr != nil {
		log.Error(updateErr)
		return errs.MysqlOperateError
	}

	return nil
}

//获取项目所有的迭代状态关联
func GetProjectIterationRelationStatusInfo(orgId int64, projectId int64) ([]bo.PpmPriIterationStatusRelationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmPriIterationStatusRelation{}
	err := mysql.SelectAllByCond(consts.TableIterationStatusRelation, db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.PpmPriIterationStatusRelationBo{}
	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ObjectCopyError
	}

	return *bos, nil
}
