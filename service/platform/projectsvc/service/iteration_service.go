package service

import (
	"strings"
	"time"

	"github.com/spf13/cast"

	msgPb "github.com/star-table/interface/golang/msg/v1"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/facade/common/report"

	"github.com/star-table/polaris-backend/common/core/util/asyn"

	"github.com/star-table/polaris-backend/common/model/bo/status"

	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/date"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func IterationList(orgId int64, page uint, size uint, params *vo.IterationListReq) (*vo.IterationList, errs.SystemErrorInfo) {
	cond := db.Cond{}
	if params != nil {

		isError, err := dealParam(params, &cond, orgId)
		if isError {
			return nil, err
		}

	}
	cond[consts.TcOrgId] = orgId
	cond[consts.TcIsDelete] = consts.AppIsNoDelete

	bos, total, err := domain.GetIterationBoList(page, size, cond, params.OrderBy)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}

	resultList := &[]*vo.Iteration{}
	copyErr := copyer.Copy(bos, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	//状态列表
	statusList, err := domain.GetStatusListByCategoryRelaxed(orgId, consts.ProcessStatusCategoryIteration)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	ownerIds := make([]int64, 0)
	iterationIds := []int64{}
	for _, result := range *resultList {
		ownerId := result.Owner
		ownerIdExist, _ := slice.Contain(ownerIds, ownerId)
		if !ownerIdExist {
			ownerIds = append(ownerIds, ownerId)
		}
		iterationIds = append(iterationIds, result.ID)
	}
	ownerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, ownerIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//转换map
	ownerMap := maps.NewMap("UserId", ownerInfos)
	statusMap := maps.NewMap("StatusId", statusList)

	//获取任务信息
	projectId := int64(0)
	if params.ProjectID != nil {
		projectId = *params.ProjectID
	}
	issueStat, _ := domain.GetIssueInfoByIterationForLc(iterationIds, orgId, projectId)

	//获取负责人
	for _, result := range *resultList {
		ownerInfo := &vo.HomeIssueOwnerInfo{}
		if userCacheInfo, ok := ownerMap[result.Owner]; ok {
			baseUserInfo := userCacheInfo.(bo.BaseUserInfoBo)
			ownerInfo.ID = baseUserInfo.UserId
			ownerInfo.Name = baseUserInfo.Name
			ownerInfo.Avatar = &baseUserInfo.Avatar
		}
		result.OwnerInfo = ownerInfo

		statusInfo := &vo.HomeIssueStatusInfo{}
		if statusCacheInfo, ok := statusMap[result.Status]; ok {
			statusCacheBo := statusCacheInfo.(bo.CacheProcessStatusBo)
			homeStatusInfoBo := domain.ConvertStatusInfoToHomeIssueStatusInfo(statusCacheBo)
			_ = copyer.Copy(homeStatusInfoBo, statusInfo)
		}
		result.StatusInfo = statusInfo

		if stat, ok := issueStat[result.ID]; ok {
			result.AllIssueCount = stat.All
			result.FinishedIssueCount = stat.Finish
		}
	}

	return &vo.IterationList{
		Total: total,
		List:  *resultList,
	}, nil
}

func dealParam(params *vo.IterationListReq, cond *db.Cond, orgId int64) (isError bool, eoor errs.SystemErrorInfo) {

	if params.Name != nil {
		(*cond)[consts.TcName] = db.Like("%" + *params.Name + "%")
	}
	if params.ProjectID != nil {
		projectId := *params.ProjectID
		_, err := domain.LoadProjectAuthBo(orgId, projectId)
		if err != nil {
			log.Error(err)
			return true, errs.BuildSystemErrorInfo(errs.IllegalityProject, err)
		}
		(*cond)[consts.TcProjectId] = projectId
	}
	if params.StatusType != nil {
		err := domain.IterationCondStatusAssembly(cond, orgId, *params.StatusType)
		if err != nil {
			log.Error(err)
			return true, errs.BuildSystemErrorInfo(errs.IterationDomainError, err)
		}
	}

	return false, nil
}

func CreateIteration(orgId, currentUserId int64, input vo.CreateIterationReq) (*vo.Void, errs.SystemErrorInfo) {
	projectId := input.ProjectID
	ownerId := input.Owner

	err := domain.AuthProject(orgId, currentUserId, input.ProjectID, consts.RoleOperationPathOrgProIteration, consts.OperationProIterationCreate)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	authFunctionErr := domain.AuthPayProjectIteration(orgId, consts.FunctionAgile, projectId)
	if authFunctionErr != nil {
		log.Error(authFunctionErr)
		return nil, authFunctionErr
	}

	//校验负责人
	if ownerId != currentUserId {
		suc := orgfacade.VerifyOrgRelaxed(orgId, ownerId)
		if !suc {
			return nil, errs.BuildSystemErrorInfo(errs.IllegalityOwner)
		}
	}

	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIteration)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}

	// 任务状态改造后，无需再查询流程状态
	// deleted code

	//创建迭代状态时间关联（只创建开始和结束的）
	var beginStatusId, endStatusId int64
	statusInsert := []po.PpmPriIterationStatusRelation{}
	relationIdResp, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIterationStatusRelation, 2)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	beginStatusId = consts.StatusNotStart.ID
	endStatusId = consts.StatusComplete.ID
	if beginStatusId != 0 {
		statusInsert = append(statusInsert, po.PpmPriIterationStatusRelation{
			Id:            relationIdResp.Ids[0].Id,
			OrgId:         orgId,
			ProjectId:     projectId,
			IterationId:   id,
			StatusId:      beginStatusId,
			PlanStartTime: time.Time(input.PlanStartTime),
			Creator:       currentUserId,
		})
	}
	if endStatusId != 0 {
		statusInsert = append(statusInsert, po.PpmPriIterationStatusRelation{
			Id:          relationIdResp.Ids[1].Id,
			OrgId:       orgId,
			ProjectId:   projectId,
			IterationId: id,
			StatusId:    endStatusId,
			PlanEndTime: time.Time(input.PlanEndTime),
			Creator:     currentUserId,
		})
	}

	//获取初始化状态
	initStatusId := consts.StatusNotStart.ID

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		entity := &bo.IterationBo{
			Id:            id,
			OrgId:         orgId,
			ProjectId:     projectId,
			Name:          input.Name,
			Sort:          id * 65536,
			Owner:         ownerId,
			PlanStartTime: input.PlanStartTime,
			PlanEndTime:   input.PlanEndTime,
			Status:        initStatusId,
			Creator:       currentUserId,
			Updator:       currentUserId,
			IsDelete:      consts.AppIsNoDelete,
		}
		err1 := domain.CreateIteration(entity, tx)
		if err1 != nil {
			log.Error(err1)
			return err1
		}

		if len(statusInsert) > 0 {
			err := mysql.TransBatchInsert(tx, &po.PpmPriIterationStatusRelation{}, slice.ToSlice(statusInsert))
			if err != nil {
				log.Error(err)
				return err
			}
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return nil, errs.CreateIterationErr
	}

	// mqtt
	asyn.Execute(func() {
		projectBo, err := domain.GetProjectSimple(orgId, projectId)
		if err != nil {
			log.Errorf("[GetProjectSimple] failed, err:%v", err)
			return
		}
		columnsResp, err := domain.GetTableColumnByAppId(orgId, currentUserId, projectBo.AppId, []string{consts.BasicFieldIterationId})
		if err != nil {
			log.Errorf("[GetTableColumnByAppId] failed, err:%v", err)
			return
		}
		for _, table := range columnsResp.Tables {
			if len(table.Columns) > 0 && table.Columns[0].Name == consts.BasicFieldIterationId {
				column := table.Columns[0]
				err = addIterationOptions(orgId, projectId, column)
				if err != nil {
					log.Errorf("[addIterationOptions] failed, err:%v", err)
					return
				}

				e := &commonvo.TableEvent{}
				e.OrgId = orgId
				e.AppId = projectBo.AppId
				e.ProjectId = projectId
				e.TableId = table.TableId
				e.New = column

				openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
				openTraceIdStr := cast.ToString(openTraceId)

				report.ReportTableEvent(msgPb.EventType_TableColumnRefresh, openTraceIdStr, e)
			}
		}
	})

	return &vo.Void{
		ID: id,
	}, nil
}

func UpdateIteration(orgId, currentUserId int64, input vo.UpdateIterationReq) (*vo.Void, errs.SystemErrorInfo) {
	targetId := input.ID

	iterationBo, err1 := domain.GetIterationBo(targetId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}
	err := domain.AuthProject(orgId, currentUserId, iterationBo.ProjectId, consts.RoleOperationPathOrgProIteration, consts.OperationProIterationModify)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	iterationNewBo := &bo.IterationBo{}
	_ = copyer.Copy(iterationBo, iterationNewBo)
	iterationUpdateBo := &bo.IterationUpdateBo{}
	upd := mysql.Upd{}

	isErr1, err1 := dealNeedUpdate(input, &upd, iterationNewBo, currentUserId, orgId)

	if isErr1 {
		return nil, err1
	}

	var planStartTime = iterationBo.PlanStartTime
	var planEndTime = iterationBo.PlanEndTime
	// && checkNeedUpdate(input.UpdateColumns, "planEndTime")
	if checkNeedUpdate(input.UpdateFields, "planStartTime") {
		if input.PlanStartTime == nil {
			return nil, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "计划开始时间不能为空")
		}
		planStartTime = *input.PlanStartTime
		upd[consts.TcPlanStartTime] = date.FormatTime(planStartTime)
		iterationNewBo.PlanStartTime = planStartTime
	}
	if checkNeedUpdate(input.UpdateFields, "planEndTime") {
		if input.PlanEndTime == nil {
			return nil, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "计划结束时间不能为空")
		}
		planEndTime = *input.PlanEndTime
		upd[consts.TcPlanEndTime] = date.FormatTime(planEndTime)
		iterationNewBo.PlanEndTime = planEndTime
	}
	if time.Time(planEndTime).Before(time.Time(planStartTime)) {
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "开始时间应该在结束时间之前")
	}
	upd[consts.TcUpdator] = currentUserId

	iterationUpdateBo.Id = targetId
	iterationUpdateBo.Upd = upd
	iterationUpdateBo.IterationNewBo = *iterationNewBo

	err1 = domain.UpdateIteration(iterationUpdateBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	// 任务状态改造后，迭代的状态不再从对象流程中查询。而是固定的状态值。
	// deleted code

	//创建迭代状态时间关联（只创建开始和结束的）
	var beginStatusId, endStatusId int64
	beginStatusId = consts.StatusNotStart.ID
	endStatusId = consts.StatusComplete.ID
	if checkNeedUpdate(input.UpdateFields, "planStartTime") {
		err := domain.UpdateOrInsertIterationRelation(orgId, iterationBo.ProjectId, input.ID, beginStatusId, input.PlanStartTime, nil, nil, nil, currentUserId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	if checkNeedUpdate(input.UpdateFields, "planEndTime") {
		err := domain.UpdateOrInsertIterationRelation(orgId, iterationBo.ProjectId, input.ID, endStatusId, nil, input.PlanEndTime, nil, nil, currentUserId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	return &vo.Void{
		ID: input.ID,
	}, nil
}

//判断处理是否需要更新
func dealNeedUpdate(input vo.UpdateIterationReq, upd *mysql.Upd, iterationNewBo *bo.IterationBo, currentUserId int64, orgId int64) (bool, errs.SystemErrorInfo) {

	if checkNeedUpdate(input.UpdateFields, "name") {
		if input.Name == nil {
			return true, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "迭代名称不能为空")
		}
		//input.Name指针指向的值 去空格给name
		name := strings.Trim(*input.Name, " ")
		if name == "" || strs.Len(name) > 200 {
			return true, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "迭代名称不能为空且限制在200字以内")
		}
		(*upd)[consts.TcName] = name
		(*iterationNewBo).Name = name
	}

	if checkNeedUpdate(input.UpdateFields, "owner") {
		if input.Owner == nil {
			return true, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "负责人不能为空")
		}
		ownerId := *input.Owner
		if ownerId != iterationNewBo.Owner {
			suc := orgfacade.VerifyOrgRelaxed(orgId, ownerId)
			if !suc {
				return true, errs.BuildSystemErrorInfo(errs.IllegalityOwner)
			}

			(*upd)[consts.TcOwner] = ownerId
			(*iterationNewBo).Owner = ownerId
		}
	}

	return false, nil

}

func DeleteIteration(orgId, currentUserId int64, input vo.DeleteIterationReq) (*vo.Void, errs.SystemErrorInfo) {

	targetId := input.ID

	bo, err1 := domain.GetIterationBo(targetId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	err := domain.AuthProject(orgId, currentUserId, bo.ProjectId, consts.RoleOperationPathOrgProIteration, consts.OperationProIterationDelete)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	err2 := domain.DeleteIteration(bo, currentUserId)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}

	return &vo.Void{
		ID: targetId,
	}, nil
}

func IterationStatusTypeStat(orgId int64, input *vo.IterationStatusTypeStatReq) (*vo.IterationStatusTypeStatResp, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacorgIdade.GetCurrentUserRelaxed(ctx)
	//	//if err != nil {
	//	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//	//}
	//	//orgId := cacheUserInfo.OrgId

	projectId := int64(0)
	if input != nil {
		if input.ProjectID != nil {
			projectId := *input.ProjectID
			_, err := domain.LoadProjectAuthBo(orgId, projectId)
			if err != nil {
				log.Error(err)
				return nil, errs.BuildSystemErrorInfo(errs.IllegalityProject, err)
			}
		}
	}

	countBo, err1 := domain.StatisticIterationCountGroupByStatus(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IterationDomainError, err1)
	}

	return &vo.IterationStatusTypeStatResp{
		NotStartTotal:   countBo.NotStartTotal,
		ProcessingTotal: countBo.ProcessingTotal,
		CompletedTotal:  countBo.FinishedTotal,
		Total:           countBo.NotStartTotal + countBo.ProcessingTotal + countBo.FinishedTotal,
	}, nil
}

////规划迭代
//func IterationIssueRelate(orgId, currentUserId int64, input vo.IterationIssueRealtionReq) (*vo.Void, errs.SystemErrorInfo) {
//	iterationBo, err1 := domain.GetIterationBoByOrgId(input.IterationID, orgId)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIteration)
//	}
//
//	err := domain.AuthProject(orgId, currentUserId, iterationBo.ProjectId, consts.RoleOperationPathOrgProIteration, consts.OperationProIterationBind)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
//	}
//
//	err2 := domain.RelateIteration(orgId, iterationBo.ProjectId, iterationBo.Id, currentUserId, input.AddIssueIds, input.DelIssueIds)
//	if err2 != nil {
//		log.Error(err2)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
//	}
//
//	return &vo.Void{
//		ID: iterationBo.Id,
//	}, nil
//}

func UpdateIterationStatus(orgId, currentUserId int64, input vo.UpdateIterationStatusReq) (*vo.Void, errs.SystemErrorInfo) {

	iterationBo, err1 := domain.GetIterationBoByOrgId(input.ID, orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIteration)
	}

	err := domain.AuthProject(orgId, currentUserId, iterationBo.ProjectId, consts.RoleOperationPathOrgProIteration, consts.OperationProIterationModifyStatus)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	err2 := domain.UpdateIterationStatus(*iterationBo, input.NextStatusID, currentUserId, input.BeforeStatusEndTime, input.NextStatusStartTime)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IterationDomainError, err2)
	}

	return &vo.Void{
		ID: iterationBo.Id,
	}, nil
}

func IterationInfo(orgId int64, input vo.IterationInfoReq) (*vo.IterationInfoResp, errs.SystemErrorInfo) {

	iterationBo, err1 := domain.GetIterationBoByOrgId(input.ID, orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIteration)
	}

	iterationInfoBo, err2 := domain.GetIterationInfoBo(*iterationBo)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IterationDomainError, err2)
	}

	iterationInfoResp := &vo.IterationInfoResp{}
	err3 := copyer.Copy(iterationInfoBo, iterationInfoResp)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	//获取所有流程状态
	allStatus := consts.IterationStatusList

	//获取关联数据
	relationInfo, err := domain.GetIterationRelationStatusInfo(orgId, input.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	relationMap := map[int64]vo.StatusTimeInfo{}
	for _, relationBo := range relationInfo {
		relationMap[relationBo.StatusId] = vo.StatusTimeInfo{
			StatusID:      relationBo.StatusId,
			PlanStartTime: types.Time(relationBo.PlanStartTime),
			PlanEndTime:   types.Time(relationBo.PlanEndTime),
			StartTime:     types.Time(relationBo.StartTime),
			EndTime:       types.Time(relationBo.EndTime),
		}
	}
	for _, status := range allStatus {
		temp := vo.StatusTimeInfo{}
		if relation, ok := relationMap[status.ID]; ok {
			temp = relation
			temp.StatusType = status.Type
			temp.StatusName = status.Name
		} else {
			defaultTime := types.Time(consts.BlankTimeObject)
			temp = vo.StatusTimeInfo{
				StatusID:      status.ID,
				StatusName:    status.Name,
				StatusType:    status.Type,
				PlanStartTime: defaultTime,
				PlanEndTime:   defaultTime,
				StartTime:     defaultTime,
				EndTime:       defaultTime,
			}
		}

		iterationInfoResp.StatusTimeInfo = append(iterationInfoResp.StatusTimeInfo, &temp)
	}
	ownerInfo := &vo.HomeIssueOwnerInfo{}
	if iterationInfoResp.Owner != nil {
		ownerInfo = &vo.HomeIssueOwnerInfo{
			ID:         iterationInfoResp.Owner.UserID,
			Name:       iterationInfoResp.Owner.Name,
			Avatar:     &iterationInfoResp.Owner.Avatar,
			IsDeleted:  iterationInfoResp.Owner.IsDeleted,
			IsDisabled: iterationInfoResp.Owner.IsDisabled,
		}
	}
	iterationInfoResp.Iteration.OwnerInfo = ownerInfo

	iterationInfoResp.Iteration.StatusInfo = iterationInfoResp.Status

	//获取任务信息
	issueStat, err := domain.GetIssueInfoByIterationForLc([]int64{input.ID}, orgId, iterationBo.ProjectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if stat, ok := issueStat[input.ID]; ok {
		iterationInfoResp.Iteration.AllIssueCount = stat.All
		iterationInfoResp.Iteration.FinishedIssueCount = stat.Finish
	}
	iterationInfoResp.IterStatusList = domain.GetIterationStatusList()

	return iterationInfoResp, nil
}

//获取未完成的的迭代列表
func GetNotCompletedIterationBoList(orgId int64, projectId int64) ([]bo.IterationBo, errs.SystemErrorInfo) {
	return domain.GetNotCompletedIterationBoList(orgId, projectId)
}

func UpdateIterationStatusTime(orgId, userId int64, input vo.UpdateIterationStatusTimeReq) (*vo.Void, errs.SystemErrorInfo) {
	iterationBo, err1 := domain.GetIterationBo(input.IterationID)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	//获取所有流程状态
	allStatus := consts.IterationStatusList
	allStatusMap := map[int64]status.StatusInfoBo{}
	for _, status := range allStatus {
		allStatusMap[status.ID] = status
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		upd := mysql.Upd{}
		var planStartTime = iterationBo.PlanStartTime
		var planEndTime = iterationBo.PlanEndTime
		for _, req := range input.StatusUpdate {
			if req == nil {
				continue
			}
			if statusInfo, ok := allStatusMap[req.StatusID]; ok {
				err := domain.UpdateOrInsertIterationRelation(orgId, iterationBo.ProjectId, input.IterationID, req.StatusID, req.PlanStartTime, req.PlanEndTime, req.StartTime, req.EndTime, userId)
				if err != nil {
					log.Error(err)
					return err
				}
				if statusInfo.Type == consts.StatusTypeNotStart && req.PlanStartTime != nil {
					planStartTime = *req.PlanStartTime
					upd[consts.TcPlanStartTime] = date.FormatTime(planStartTime)
				}

				if statusInfo.Type == consts.StatusTypeComplete && req.PlanEndTime != nil {
					planEndTime = *req.PlanEndTime
					upd[consts.TcPlanEndTime] = date.FormatTime(planEndTime)
				}
			}
		}
		if planStartTime != iterationBo.PlanStartTime || planEndTime != iterationBo.PlanEndTime {
			if time.Time(planEndTime).Before(time.Time(planStartTime)) {
				return errs.PlanEndTimeInvalidError
			}
			_, updateErr := mysql.UpdateSmartWithCond(consts.TableIteration, db.Cond{
				consts.TcId: input.IterationID,
			}, upd)
			if updateErr != nil {
				log.Error(updateErr)
				return errs.MysqlOperateError
			}

		}
		return nil
	})

	if transErr != nil {
		log.Error(transErr)
		return nil, errs.BuildSystemErrorInfo(errs.IterationDomainError, transErr)
	}
	return &vo.Void{ID: input.IterationID}, nil
}

func GetSimpleIterationInfo(orgId int64, ids []int64) ([]bo.IterationBo, errs.SystemErrorInfo) {
	return domain.GetIterationsInfo(orgId, ids)
}

func UpdateIterationSort(orgId int64, userId int64, params vo.UpdateIterationSortReq) (*vo.Void, errs.SystemErrorInfo) {
	iterationBo, err1 := domain.GetIterationBoByOrgId(params.IterationID, orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIteration)
	}

	err := domain.AuthProject(orgId, userId, iterationBo.ProjectId, consts.RoleOperationPathOrgProIteration, consts.OperationProIterationModify)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	uid := uuid.NewUuid()
	lockKey := consts.IterationSortLock
	suc, cacheErr := cache.TryGetDistributedLock(lockKey, uid)
	if cacheErr != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, cacheErr)
		return nil, errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	}

	//获取前一个迭代
	beforeSort := int64(0)
	afterSort := int64(0)
	if params.AfterID != nil {
		afterInfo, err1 := domain.GetIterationBoByOrgId(*params.AfterID, orgId)
		if err1 != nil {
			log.Error(err1)
			return nil, errs.BuildSystemErrorInfo(errs.IllegalityIteration)
		}

		afterSort = afterInfo.Sort

		before, err := domain.GetBeforeIterationSort(orgId, iterationBo.ProjectId, afterSort)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		beforeSort = before
	} else if params.BeforeID != int64(0) {
		beforeInfo, err1 := domain.GetIterationBoByOrgId(params.BeforeID, orgId)
		if err1 != nil {
			log.Error(err1)
			return nil, errs.BuildSystemErrorInfo(errs.IllegalityIteration)
		}
		beforeSort = beforeInfo.Sort

		after, err := domain.GetAfterIterationSort(orgId, iterationBo.ProjectId, beforeSort)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		afterSort = after
	} else {
		//如果移到最前面
		//查找项目下最小的sort
		min, err := domain.GetMinIterationSort(orgId, iterationBo.ProjectId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		afterSort = min
		beforeSort = min - 65536
	}

	if afterSort-beforeSort < 2 {
		//空隙不够，则前移sort 位置
		err := mysql.TransX(func(tx sqlbuilder.Tx) error {
			//前移所有之前的sort
			_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIteration, db.Cond{
				consts.TcOrgId:     orgId,
				consts.TcProjectId: iterationBo.ProjectId,
				consts.TcSort:      db.Lte(beforeSort),
			}, mysql.Upd{
				consts.TcSort: db.Raw("sort - 65535"),
			})
			if err != nil {
				log.Error(err)
				return err
			}

			//把当前任务插入
			_, err1 := mysql.UpdateSmartWithCond(consts.TableIteration, db.Cond{
				consts.TcOrgId: orgId,
				consts.TcId:    params.IterationID,
			}, mysql.Upd{
				consts.TcUpdator: userId,
				consts.TcSort:    (afterSort + beforeSort - 65535) / 2,
			})
			if err1 != nil {
				log.Error(err1)
				return err1
			}
			return nil
		})
		if err != nil {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}
	} else {
		_, err := mysql.UpdateSmartWithCond(consts.TableIteration, db.Cond{
			consts.TcOrgId: orgId,
			consts.TcId:    params.IterationID,
		}, mysql.Upd{
			consts.TcUpdator: userId,
			consts.TcSort:    (afterSort + beforeSort) / 2,
		})
		if err != nil {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}
	}

	return &vo.Void{ID: params.IterationID}, nil
}
