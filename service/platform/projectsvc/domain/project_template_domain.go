package domain

import (
	"time"

	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

/**
处理任务栏和状态
projectObjectTypeMap 旧任务栏->新任务栏
statusMap 旧任务状态->新任务状态
*/
//func CreateProjectObjectTypeProcessFromTemplate(tx sqlbuilder.Tx, orgId, userId, projectId int64, projectTypeId int64, projectObjectTypeTemplates []projectvo.ProjectObjectTypeTemplate) (map[int64]int64, map[int64]int64, errs.SystemErrorInfo) {
//	projectObjectTypeMap := make(map[int64]int64, 0)
//	statusMap := make(map[int64]int64, 0)
//
//	if len(projectObjectTypeTemplates) == 0 {
//		return projectObjectTypeMap, statusMap, nil
//	}
//	pos := make([]interface{}, 0)
//	processPos := make([]interface{}, 0)
//
//	insertProcessStatusList := make([]bo.ProcessStatusBo, 0)
//	insertProcessProcessStatusList := make([]bo.ProcessProcessStatusBo, 0)
//
//	length := len(projectObjectTypeTemplates)
//	posIds, posIdsErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectObjectType, length)
//	if posIdsErr != nil {
//		log.Error(posIdsErr)
//		return nil, nil, posIdsErr
//	}
//
//	processPosIds, processPosIdsErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectObjectTypeProcess, length)
//	if processPosIdsErr != nil {
//		log.Error(processPosIdsErr)
//		return nil, nil, processPosIdsErr
//	}
//
//	//状态（通用项目不需要，这里会是0）
//	statusCount := 0
//	for _, template := range projectObjectTypeTemplates {
//		statusCount += len(template.Status)
//	}
//
//	statusIds := &bo.IdCodes{}
//	processStatusIds := &bo.IdCodes{}
//	if projectTypeId == consts.ProjectTypeAgileId {
//		statusIdsMid, statusIdsErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProcessStatus, statusCount)
//		if statusIdsErr != nil {
//			log.Error(statusIdsErr)
//			return nil, nil, statusIdsErr
//		}
//		statusIds = statusIdsMid
//
//		processStatusIdsMid, processStatusIdsErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProcessProcessStatus, statusCount)
//		if processStatusIdsErr != nil {
//			log.Error(processStatusIdsErr)
//			return nil, nil, processStatusIdsErr
//		}
//		processStatusIds = processStatusIdsMid
//	}
//
//	statusIdIndex := 0
//	for i, template := range projectObjectTypeTemplates {
//		pos = append(pos, po.PpmPrsProjectObjectType{
//			Id:         posIds.Ids[i].Id,
//			OrgId:      orgId,
//			Name:       template.Name,
//			ObjectType: template.ObjectType,
//			LangCode:   template.LangCode,
//			Sort:       template.Sort,
//			Creator:    userId,
//			Updator:    userId,
//		})
//		projectObjectTypeMap[template.ID] = posIds.Ids[i].Id
//
//		processPos = append(processPos, po.PpmPrsProjectObjectTypeProcess{
//			Id:                  processPosIds.Ids[i].Id,
//			OrgId:               orgId,
//			ProjectId:           projectId,
//			TableId: posIds.Ids[i].Id,
//			ProcessId:           template.ProcessId,
//			Sort:                template.Sort,
//			Creator:             userId,
//			Updator:             userId,
//		})
//
//		if projectTypeId == consts.ProjectTypeAgileId {
//			for _, status := range template.Status {
//				insertProcessStatusList = append(insertProcessStatusList, bo.ProcessStatusBo{
//					Id:        statusIds.Ids[statusIdIndex].Id,
//					OrgId:     orgId,
//					ProjectId: projectId,
//					LangCode:  "",
//					Name:      status.Name,
//					Sort:      status.Sort,
//					BgStyle:   status.BgStyle,
//					FontStyle: status.FontStyle,
//					Type:      status.Type,
//					Category:  status.Category,
//					Creator:   userId,
//					Updator:   userId,
//				})
//				statusMap[status.ID] = statusIds.Ids[statusIdIndex].Id
//
//				insertProcessProcessStatusList = append(insertProcessProcessStatusList, bo.ProcessProcessStatusBo{
//					Id:              processStatusIds.Ids[statusIdIndex].Id,
//					OrgId:           orgId,
//					ProcessId:       template.ProcessId,
//					ProjectId:       projectId,
//					ProcessStatusId: statusIds.Ids[statusIdIndex].Id,
//					IsInitStatus:    status.IsInitStatus,
//					Sort:            status.Sort,
//					Creator:         userId,
//					Updator:         userId,
//				})
//
//				statusIdIndex++
//			}
//		}
//	}
//	//if projectTypeId == consts.ProjectTypeAgileId {
//	//	//敏捷带上迭代
//	//	projectTypeProjectObjectTypeList, err := GetProjectTypeProjectObjectTypeList(orgId)
//	//	if err != nil {
//	//		return nil, nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//	//	}
//	//	//查通用的项目对象类型
//	//	projectObjectTypeList, err := GetProjectObjectTypeList(0, 0)
//	//	if err != nil {
//	//		log.Error(err)
//	//		return nil, nil, err
//	//	}
//	//	projectObjectTypeMap := map[int64]bo.ProjectObjectTypeBo{}
//	//	for _, typeBo := range *projectObjectTypeList {
//	//		//只取迭代
//	//		if typeBo.LangCode == consts.ProjectObjectTypeLangCodeIteration {
//	//			projectObjectTypeMap[typeBo.Id] = typeBo
//	//		}
//	//	}
//	//	for _, v := range *projectTypeProjectObjectTypeList {
//	//		projectObjectType, ok := projectObjectTypeMap[v.TableId]
//	//		if !ok {
//	//			continue
//	//		}
//	//		projectObjectTypePo := &po.PpmPrsProjectObjectType{}
//	//		_ = copyer.Copy(projectObjectType, projectObjectTypePo)
//	//
//	//		id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectObjectType)
//	//		if err != nil {
//	//			log.Error(err)
//	//			return nil, nil, errs.ApplyIdError
//	//		}
//	//		projectObjectTypePo.Id = id
//	//		projectObjectTypePo.OrgId = orgId
//	//		pos = append(pos, *projectObjectTypePo)
//	//
//	//		processId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectObjectTypeProcess)
//	//		if err != nil {
//	//			return nil, nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
//	//		}
//	//		processPos = append(processPos, po.PpmPrsProjectObjectTypeProcess{
//	//			Id:                  processId,
//	//			OrgId:               orgId,
//	//			ProjectId:           projectId,
//	//			TableId: projectObjectTypePo.Id,
//	//			ProcessId:           v.DefaultProcessId,
//	//			Creator:             userId,
//	//			Updator:             userId,
//	//		})
//	//	}
//	//}
//
//	// 此处删除通过模版建立流程状态
//
//	insertErr1 := mysql.TransBatchInsert(tx, &po.PpmPrsProjectObjectType{}, slice.ToSlice(pos))
//	if insertErr1 != nil {
//		log.Error(insertErr1)
//		return nil, nil, errs.MysqlOperateError
//	}
//	insertErr2 := mysql.TransBatchInsert(tx, &po.PpmPrsProjectObjectTypeProcess{}, slice.ToSlice(processPos))
//	if insertErr2 != nil {
//		log.Error(insertErr2)
//		return nil, nil, errs.MysqlOperateError
//	}
//
//	return projectObjectTypeMap, statusMap, nil
//}

func GetProjectIterationInfo(orgId, projectId int64) ([]projectvo.ProjectIterationTemplate, errs.SystemErrorInfo) {
	projectIterationTemplates := make([]projectvo.ProjectIterationTemplate, 0)

	//获取迭代所有流程状态
	allStatus := consts.IterationStatusList
	statusTypeMap := make(map[int64]int, 0)
	for _, status := range allStatus {
		statusTypeMap[status.ID] = status.Type
	}

	//获取所有迭代
	iterations, _, err := GetIterationBoList(0, 0, db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
		consts.TcOrgId:     orgId,
	}, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//获取所有迭代的状态详情
	iterationRelation, relationErr := GetProjectIterationRelationStatusInfo(orgId, projectId)
	if relationErr != nil {
		log.Error(relationErr)
		return nil, relationErr
	}

	relationMap := make(map[int64][]projectvo.IterationStatus, 0)
	for _, relationBo := range iterationRelation {
		st := relationBo.PlanStartTime
		et := relationBo.PlanEndTime

		pst := relationBo.PlanStartTime
		pet := relationBo.PlanEndTime
		startTime, endTime, planStartTime, planEndTime := "", "", "", ""
		if !st.IsZero() {
			startTime = st.Format(consts.AppTimeFormat)
		}
		if !et.IsZero() {
			endTime = et.Format(consts.AppTimeFormat)
		}
		if !pst.IsZero() {
			planStartTime = pst.Format(consts.AppTimeFormat)
		}
		if !pet.IsZero() {
			planEndTime = pet.Format(consts.AppTimeFormat)
		}
		temp := projectvo.IterationStatus{
			StartTime:     startTime,
			EndTime:       endTime,
			PlanStartTime: planStartTime,
			PlanEndTime:   planEndTime,
			StatusType:    statusTypeMap[relationBo.StatusId],
		}
		relationMap[relationBo.IterationId] = append(relationMap[relationBo.IterationId], temp)
	}

	for _, iteration := range *iterations {
		st := iteration.PlanStartTime
		et := iteration.PlanEndTime
		startTime, endTime := "", ""
		if !time.Time(st).IsZero() {
			startTime = time.Time(st).Format(consts.AppTimeFormat)
		}
		if !time.Time(et).IsZero() {
			endTime = time.Time(et).Format(consts.AppTimeFormat)
		}
		current := projectvo.ProjectIterationTemplate{
			ID:         iteration.Id,
			Name:       iteration.Name,
			StartTime:  startTime,
			EndTime:    endTime,
			Sort:       iteration.Sort,
			Owner:      iteration.Owner,
			StatusType: statusTypeMap[iteration.Status],
		}
		if relations, ok := relationMap[iteration.Id]; ok {
			current.IterationStatus = relations
		}
		projectIterationTemplates = append(projectIterationTemplates, current)
	}

	return projectIterationTemplates, nil
}

//旧迭代->新迭代
func CreateIterationFromTemplate(orgId, userId, projectId int64, projectIterationTemplates []projectvo.ProjectIterationTemplate) (map[int64]int64, errs.SystemErrorInfo) {
	if len(projectIterationTemplates) == 0 {
		return nil, nil
	}
	//获取迭代所有流程状态
	allStatus := consts.IterationStatusList

	typeToStatusMap := make(map[int]int64, 0)
	for _, status := range allStatus {
		typeToStatusMap[status.Type] = status.ID
	}

	iterationCount := len(projectIterationTemplates)
	relationCount := 0
	for _, template := range projectIterationTemplates {
		relationCount += len(template.IterationStatus)
	}

	iterationIds, iterationIdsErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIteration, iterationCount)
	if iterationIdsErr != nil {
		log.Error(iterationIdsErr)
		return nil, iterationIdsErr
	}

	relationIds := &bo.IdCodes{}
	if relationCount > 0 {
		ids, relationIdsErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIterationStatusRelation, relationCount)
		if relationIdsErr != nil {
			log.Error(relationIdsErr)
			return nil, relationIdsErr
		}
		relationIds = ids
	}

	// 取出第一个作为初始化状态id
	initStatusId := int64(consts.StatusIdNotStart)

	pos := make([]interface{}, 0)
	relationPos := make([]interface{}, 0)
	relationIndex := 0
	iterationMap := make(map[int64]int64, 0) //旧迭代-》新迭代
	for i, template := range projectIterationTemplates {
		temp := po.PpmPriIteration{
			Id:            iterationIds.Ids[i].Id,
			OrgId:         orgId,
			ProjectId:     projectId,
			Name:          template.Name,
			Sort:          template.Sort,
			Owner:         template.Owner,
			PlanStartTime: time.Time{},
			PlanEndTime:   time.Time{},
			Status:        0,
			Creator:       userId,
			Updator:       userId,
		}
		if statusId, ok := typeToStatusMap[template.StatusType]; ok {
			temp.Status = statusId
		} else {
			temp.Status = initStatusId
		}
		if template.StartTime != "" {
			start, _ := time.ParseInLocation(consts.AppTimeFormat, template.StartTime, time.Local)
			temp.PlanStartTime = start
		}
		if template.EndTime != "" {
			end, _ := time.ParseInLocation(consts.AppTimeFormat, template.EndTime, time.Local)
			temp.PlanEndTime = end
		}
		pos = append(pos, temp)

		iterationMap[template.ID] = iterationIds.Ids[i].Id

		for _, status := range template.IterationStatus {
			relation := po.PpmPriIterationStatusRelation{
				Id:            relationIds.Ids[relationIndex].Id,
				OrgId:         orgId,
				ProjectId:     projectId,
				IterationId:   iterationIds.Ids[i].Id,
				PlanStartTime: time.Time{},
				PlanEndTime:   time.Time{},
				StartTime:     time.Time{},
				EndTime:       time.Time{},
				Creator:       userId,
				Updator:       userId,
			}
			if statusId, ok := typeToStatusMap[status.StatusType]; ok {
				relation.StatusId = statusId
			} else {
				relation.StatusId = initStatusId
			}
			if status.StartTime != "" {
				start, _ := time.ParseInLocation(consts.AppTimeFormat, status.StartTime, time.Local)
				relation.StartTime = start
			}
			if status.EndTime != "" {
				end, _ := time.ParseInLocation(consts.AppTimeFormat, status.EndTime, time.Local)
				relation.EndTime = end
			}
			if status.PlanStartTime != "" {
				start, _ := time.ParseInLocation(consts.AppTimeFormat, status.PlanStartTime, time.Local)
				relation.PlanStartTime = start
			}
			if status.PlanEndTime != "" {
				end, _ := time.ParseInLocation(consts.AppTimeFormat, status.PlanEndTime, time.Local)
				relation.PlanEndTime = end
			}
			relationPos = append(relationPos, relation)
			relationIndex++
		}
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(pos) > 0 {
			err := mysql.TransBatchInsert(tx, &po.PpmPriIteration{}, slice.ToSlice(pos))
			if err != nil {
				log.Error(err)
				return err
			}
		}

		if len(relationPos) > 0 {
			err := mysql.TransBatchInsert(tx, &po.PpmPriIterationStatusRelation{}, slice.ToSlice(relationPos))
			if err != nil {
				log.Error(err)
				return err
			}
		}

		return nil
	})

	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}

	return iterationMap, nil
}
