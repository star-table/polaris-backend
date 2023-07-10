package domain

import (
	"time"

	tablePb "github.com/star-table/interface/golang/table/v1"

	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/date"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/bo/status"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func UpdateIterationStatus(iterationBo bo.IterationBo, nextStatusId, operationId int64, beforeEndTime, nextStartTime types.Time) errs.SystemErrorInfo {
	orgId := iterationBo.OrgId
	iterationId := iterationBo.Id

	if iterationBo.Status == nextStatusId {
		log.Error("更新迭代状态-要更新的状态和当前状态一样")
		return errs.BuildSystemErrorInfo(errs.IterationStatusUpdateError)
	}

	//验证状态有效性
	if isOk, _ := slice.Contain([]int64{
		consts.StatusNotStart.ID,
		consts.StatusRunning.ID,
		consts.StatusComplete.ID,
	}, nextStatusId); !isOk {
		err := errs.NotSupportIterStatus
		log.Errorf("[UpdateIterationStatus] err: %v", err)
		return err
	}
	nextStatusBo := GetIterationStatusInfo(nextStatusId)
	nextStatus := &bo.CacheProcessStatusBo{
		StatusId:    nextStatusBo.ID,
		StatusType:  nextStatusBo.Type,
		Category:    consts.ProcessStatusCategoryIteration,
		Name:        nextStatusBo.Name,
		DisplayName: nextStatusBo.Name,
	}
	// 任务状态改造后，无需查询 process
	// deleted code

	projectId := iterationBo.ProjectId
	needInsertStat := false
	//判断迭代下是否有未完成的任务
	if nextStatus.StatusType == consts.StatusTypeComplete {
		// 直接从无码中查询该迭代下， 未完成类型的任务是否存在
		count, err := getRowsCountByCondition(orgId, operationId, &tablePb.Condition{
			Type: tablePb.ConditionType_and,
			Conditions: GetNoRecycleCondition(
				GetRowsCondition(consts.BasicFieldIterationId, tablePb.ConditionType_equal, iterationId, nil),
				GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_equal, projectId, nil),
				GetRowsCondition(consts.BasicFieldIssueStatusType, tablePb.ConditionType_in, nil, []interface{}{consts.StatusTypeNotStart, consts.StatusTypeRunning}),
			),
		})

		if err != nil {
			log.Errorf("[UpdateIterationStatus] err: %c\n", err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}

		if count > 0 {
			return errs.BuildSystemErrorInfo(errs.IterationExistingNotFinishedTask)
		}

		//是否创建统计基准数据
		needInsertStat = true
	}

	//处理db
	err3 := dealTx(needInsertStat, iterationBo, iterationId, orgId, nextStatusId, beforeEndTime, nextStartTime, operationId)

	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err3)
	}
	return nil
}

func GetIterationStatusList() []*vo.HomeIssueStatusInfo {
	allIterStatusArr := consts.IterationStatusList
	list := make([]*vo.HomeIssueStatusInfo, 0, 3)
	for i, item := range allIterStatusArr {
		list = append(list, &vo.HomeIssueStatusInfo{
			ID:          item.ID,
			Name:        item.Name,
			DisplayName: &allIterStatusArr[i].Name,
			Type:        item.Type,
			FontStyle:   item.FontStyle,
			BgStyle:     item.BgStyle,
			Sort:        item.Sort,
		})
	}

	return list
}

func GetIterationStatusInfo(statusId int64) status.StatusInfoBo {
	map1 := GetIterationStatusMapById()
	return map1[statusId]
}

func GetIterationStatusMapById() map[int64]status.StatusInfoBo {
	allIterStatusArr := consts.IterationStatusList
	map1 := make(map[int64]status.StatusInfoBo, 3)
	for _, item := range allIterStatusArr {
		map1[item.ID] = item
	}
	return map1
}

func GetIterationStatusMap2Type() map[int64]int {
	iterStatusBoMap := GetIterationStatusMapById()
	map2Type := make(map[int64]int, len(iterStatusBoMap))
	for statusId, statusBo := range iterStatusBoMap {
		map2Type[statusId] = statusBo.Type
	}

	return map2Type
}

func dealTx(needInsertStat bool, iterationBo bo.IterationBo, iterationId, orgId, nextStatusId int64, beforeEndTime, nextStartTime types.Time, operationId int64) error {
	allStatusId := []int64{
		consts.StatusNotStart.ID, consts.StatusRunning.ID, consts.StatusComplete.ID,
	}

	// 任务状态改造后，无需查询 process 的状态列表
	// deleted code

	allIterationProjectStatusMap := GetIterationStatusMap2Type()
	err3 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if needInsertStat {
			//创建统计基准数据
			err1 := AppendIterationStat(iterationBo, consts.BlankDate, tx)
			if err1 != nil {
				log.Error(err1)
				return errs.BuildSystemErrorInfo(errs.IterationDomainError, err1)
			}
		}
		_, err2 := dao.UpdateIterationByOrg(iterationId, orgId, mysql.Upd{
			consts.TcStatus: nextStatusId,
		}, tx)
		if err2 != nil {
			log.Error(err2)
			return errs.BuildSystemErrorInfo(errs.IterationStatusUpdateError)
		}

		beforeIsNotStarted := false
		nextIsComleted := false
		if typeId, ok := allIterationProjectStatusMap[iterationBo.Status]; ok {
			if typeId == consts.StatusTypeNotStart {
				beforeIsNotStarted = true
			}
		}
		if typeId, ok := allIterationProjectStatusMap[nextStatusId]; ok {
			if typeId == consts.StatusTypeComplete {
				nextIsComleted = true
			}
		}
		if ok, _ := slice.Contain(allStatusId, iterationBo.Status); ok {
			upd := mysql.Upd{
				consts.TcUpdator: operationId,
			}
			if !beforeIsNotStarted {
				upd[consts.TcEndTime] = date.FormatTime(beforeEndTime)
			} else {
				upd[consts.TcStartTime] = date.FormatTime(beforeEndTime)
			}
			_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableIterationStatusRelation, db.Cond{
				consts.TcOrgId:       orgId,
				consts.TcIterationId: iterationId,
				consts.TcStatusId:    iterationBo.Status,
				consts.TcIsDelete:    consts.AppIsNoDelete,
			}, upd)
			if err1 != nil {
				log.Error(err1)
				return err1
			}
		} else {
			id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIterationStatusRelation)
			if err != nil {
				log.Error(err)
				return err
			}
			end, _ := time.ParseInLocation(consts.AppTimeFormat, beforeEndTime.String(), time.Local)
			insertPo := &po.PpmPriIterationStatusRelation{
				Id:          id,
				OrgId:       orgId,
				ProjectId:   iterationBo.ProjectId,
				IterationId: iterationId,
				StatusId:    iterationBo.Status,
				Creator:     operationId,
			}
			if beforeIsNotStarted {
				insertPo.StartTime = end
			} else {
				insertPo.EndTime = end
			}
			err1 := mysql.TransInsert(tx, insertPo)
			if err1 != nil {
				log.Error(err1)
				return err1
			}
		}

		if ok, _ := slice.Contain(allStatusId, nextStatusId); ok {
			upd := mysql.Upd{
				consts.TcUpdator: operationId,
			}
			if nextIsComleted {
				upd[consts.TcEndTime] = date.FormatTime(nextStartTime)
			} else {
				upd[consts.TcStartTime] = date.FormatTime(nextStartTime)
			}
			_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableIterationStatusRelation, db.Cond{
				consts.TcOrgId:       orgId,
				consts.TcIterationId: iterationId,
				consts.TcStatusId:    nextStatusId,
				consts.TcIsDelete:    consts.AppIsNoDelete,
			}, upd)
			if err1 != nil {
				log.Error(err1)
				return err1
			}
		} else {
			id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIterationStatusRelation)
			if err != nil {
				log.Error(err)
				return err
			}
			start, _ := time.ParseInLocation(consts.AppTimeFormat, nextStartTime.String(), time.Local)
			insertPo := &po.PpmPriIterationStatusRelation{
				Id:          id,
				OrgId:       orgId,
				ProjectId:   iterationBo.ProjectId,
				IterationId: iterationId,
				StatusId:    nextStatusId,
				Creator:     operationId,
			}
			if nextIsComleted {
				insertPo.EndTime = start
			} else {
				insertPo.StartTime = start
			}
			err1 := mysql.TransInsert(tx, insertPo)
			if err1 != nil {
				log.Error(err1)
				return err1
			}
		}

		return nil
	})

	return err3
}
