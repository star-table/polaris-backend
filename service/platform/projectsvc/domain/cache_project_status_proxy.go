package domain

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

// GetProjectStatus 所有的项目状态都一样，因此 projectId 可以不传
func GetProjectStatus(orgId, projectId int64) ([]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	statusBoMap := make(map[int64]bo.CacheProcessStatusBo, 0)
	statusBoList := make([]bo.CacheProcessStatusBo, 0)
	allProStatusList := consts.ProjectStatusList
	for _, item := range allProStatusList {
		tmpBo := bo.CacheProcessStatusBo{
			StatusId:    item.ID,
			StatusType:  item.Type,
			Category:    consts.ProcessStatusCategoryProject,
			Name:        item.Name,
			DisplayName: item.DisplayName,
			BgStyle:     item.BgStyle,
			FontStyle:   item.FontStyle,
		}
		statusBoMap[item.ID] = tmpBo
		statusBoList = append(statusBoList, tmpBo)
	}

	return statusBoList, nil
}

func GetProjectStatusBoMap() map[int64]bo.CacheProcessStatusBo {
	allStatusList, _ := GetProjectStatus(0, 0)
	statusMapById := make(map[int64]bo.CacheProcessStatusBo, len(allStatusList))
	for _, item := range allStatusList {
		statusMapById[item.StatusId] = item
	}

	return statusMapById
}

// 批量查询项目中的状态列表
func GetProjectStatusBatch(orgId int64, projectIds []int64) (map[int64][]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	result := make(map[int64][]bo.CacheProcessStatusBo, 0)
	projectRelationList := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcProjectId:    db.In(projectIds),
		consts.TcStatus:       consts.AppStatusEnable,
		consts.TcRelationType: consts.IssueRelationTypeStatus,
	}, projectRelationList)
	if err != nil {
		log.Error(err)
		return result, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if len(*projectRelationList) == 0 {
		return result, errs.BuildSystemErrorInfo(errs.ProcessStatusNotExist)
	}

	// 对 projectRelationList 按 projectId 进行分组
	projectRelationGroups := make(map[int64][]po.PpmProProjectRelation, 0)
	for _, item := range *projectRelationList {
		if _, ok := projectRelationGroups[item.ProjectId]; ok {
			projectRelationGroups[item.ProjectId] = append(projectRelationGroups[item.ProjectId], item)
		} else {
			initList := make([]po.PpmProProjectRelation, 0)
			initList = append(initList, item)
			projectRelationGroups[item.ProjectId] = initList
		}
	}

	statusBoGroupMap, err := GetProcessStatusWithGroup(projectRelationGroups, orgId)
	if err != nil {
		log.Error(err)
		return result, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}

	return statusBoGroupMap, nil
}

// 组装项目的状态列表
func assemblyProcessStatusList(projectRelationList *[]po.PpmProProjectRelation, processStatusList *[]bo.CacheProcessStatusBo) {
	statusBoList := make([]bo.CacheProcessStatusBo, 0)
	statusBoMap := make(map[int64]bo.CacheProcessStatusBo, 0)
	allProStatusList := consts.ProjectStatusList
	for _, item := range allProStatusList {
		statusBoMap[item.ID] = bo.CacheProcessStatusBo{
			StatusId:    item.ID,
			StatusType:  item.Type,
			Category:    consts.ProcessStatusCategoryProject,
			Name:        item.Name,
			DisplayName: item.Name,
		}
	}
	for _, v := range *projectRelationList {
		if item, ok := statusBoMap[v.RelationId]; ok {
			statusBoList = append(statusBoList, item)
		}
	}

	*processStatusList = statusBoList
}

// 查询分组了的项目的任务”状态列表“
func GetProcessStatusWithGroup(projectRelationListGroup map[int64][]po.PpmProProjectRelation, orgId int64) (map[int64][]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	result := make(map[int64][]bo.CacheProcessStatusBo, 0)
	//respVo := processfacade.GetProcessStatusBatch(processvo.GetProcessStatusBatchReqVo{OrgId: orgId})
	//if respVo.Failure() {
	//	return result, respVo.Error()
	//}
	//
	//statusBoMap := make(map[int64]bo.CacheProcessStatusBo, 0)
	//if respVo.Successful() {
	//	s := respVo.NewData
	//	for _, item := range *s {
	//		statusBoMap[item.StatusId] = item
	//	}
	//}
	//for projectId, proRelations := range projectRelationListGroup {
	//	statusBoList := make([]bo.CacheProcessStatusBo, 0)
	//	for _, v := range proRelations {
	//		if item, ok := statusBoMap[v.RelationId]; ok {
	//			statusBoList = append(statusBoList, item)
	//		}
	//	}
	//	result[projectId] = statusBoList
	//}

	return result, nil
}
