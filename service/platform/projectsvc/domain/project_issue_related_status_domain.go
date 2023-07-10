package domain

import (
	"sort"

	"github.com/star-table/polaris-backend/common/model/bo/status"

	"github.com/star-table/polaris-backend/facade/tablefacade"

	tableV1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
)

// GetProjectRelatedStatus 获取项目对象类型对应的任务状态。，
func GetProjectRelatedStatus(orgId int64, projectId int64, tableId int64) ([]status.StatusInfoBo, errs.SystemErrorInfo) {
	allStatus, err := GetIssueAllStatus(orgId, []int64{projectId}, []int64{tableId})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if info, ok := allStatus[tableId]; ok {
		return info, nil
	}
	return []status.StatusInfoBo{}, nil
}

func SortStatus(statusInfos []status.StatusInfoBo) []status.StatusInfoBo {
	sort.Slice(statusInfos, func(i, j int) bool {
		if statusInfos[i].Type == statusInfos[j].Type {
			return statusInfos[i].Sort < statusInfos[j].Sort
		} else {
			return statusInfos[i].Type < statusInfos[j].Type
		}
	})

	return statusInfos
}

//func GetIssueAllStatus(orgId int64, projectIds, projectObjectTypeIds []int64) (map[int64][]bo.StatusInfoBo, errs.SystemErrorInfo) {
//	//获取所有项目对象类型的流程id
//	allProcessInfo, err := GetProjectObjectTypeProcessCache(orgId, projectObjectTypeIds)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	//获取所有关联项目的所有状态
//	allProjectStatusResp := processfacade.GetProjectAllStatus(processvo.GetProjectAllStatusReqVo{
//		OrgId:      orgId,
//		ProjectIds: projectIds,
//	})
//	if allProjectStatusResp.Failure() {
//		log.Error(allProjectStatusResp.Error())
//		return nil, allProjectStatusResp.Error()
//	}
//	lang := lang2.GetLang()
//	isOtherLang := lang2.IsEnglish()
//	otherLanguageMap := make(map[string]string, 0)
//	if tmpMap, ok1 := consts.LANG_STATUS_MAP[lang]; ok1 {
//		otherLanguageMap = tmpMap
//	}
//	projectMap := map[int64]map[int64][]bo.StatusInfoBo{}
//	for i, bos := range allProjectStatusResp.ProcessStatusBoList {
//		projectMap[i] = map[int64][]bo.StatusInfoBo{}
//		for _, statusBo := range bos {
//			if isOtherLang {
//				if tmpVal, ok2 := otherLanguageMap[statusBo.Name]; ok2 {
//					statusBo.Name = tmpVal
//				}
//			}
//			if statusBo.StatusId != 15 {
//				temp := bo.StatusInfoBo{
//					ID:        statusBo.StatusId,
//					Name:      statusBo.Name,
//					BgStyle:   statusBo.BgStyle,
//					FontStyle: statusBo.FontStyle,
//					Type:      statusBo.StatusType,
//					Sort:      statusBo.Sort,
//				}
//				projectMap[i][statusBo.ProcessId] = append(projectMap[i][statusBo.ProcessId], temp)
//			}
//		}
//	}
//
//	res := map[int64][]bo.StatusInfoBo{}
//	for _, process := range allProcessInfo {
//		if temp, ok := projectMap[process.ProjectId]; ok {
//			if temp1, ok1 := temp[process.ProcessId]; ok1 {
//				tmpList := SortStatus(temp1)
//				if isOtherLang {
//					for index, item := range tmpList {
//						if tmpVal, ok2 := otherLanguageMap[item.Name]; ok2 {
//							tmpList[index].Name = tmpVal
//						}
//					}
//				}
//				res[process.TableId] = tmpList
//			}
//		}
//	}
//	if ok, _ := slice.Contain(projectIds, int64(0)); ok {
//		defaultProcess := processfacade.GetProcessByLangCode(processvo.GetProcessByLangCodeReqVo{
//			OrgId:    orgId,
//			LangCode: consts.ProcessLangCodeDefaultTask,
//		})
//		if defaultProcess.Failure() {
//			log.Error(defaultProcess.Error())
//			return nil, defaultProcess.Error()
//		}
//		if noProjectStatus, ok1 := projectMap[0]; ok1 {
//			if info, ok2 := noProjectStatus[defaultProcess.ProcessBo.Id]; ok2 {
//				tmpList := SortStatus(info)
//				if isOtherLang {
//					for index, item := range tmpList {
//						if tmpVal, ok2 := otherLanguageMap[item.Name]; ok2 {
//							tmpList[index].Name = tmpVal
//						}
//					}
//				}
//				res[0] = tmpList
//			}
//		}
//	}
//
//	return res, nil
//}

// GetIssueAllStatus 先通过 projectObjectTypeIds 查询 tableId 再查询对应的表头中的状态配置
// 改造后，projectObjectTypeIds 的意义变成 tableIds
// @param projectIds 这个参数暂且用不上，可以传 0
// 返回值是按 tableId 分组
func GetIssueAllStatus(orgId int64, projectIds, tableIds []int64) (map[int64][]status.StatusInfoBo, errs.SystemErrorInfo) {
	statusGroupMap := make(map[int64][]status.StatusInfoBo, 0)
	if len(tableIds) == 0 {
		return statusGroupMap, nil
	}
	// tableIds := GetTableIdsByProjectObjectIds(orgId, projectObjectTypeIds)
	tableColumnsResp := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
		OrgId:  orgId,
		UserId: 0,
		Input: &tableV1.ReadTableSchemasRequest{
			TableIds:  tableIds,
			ColumnIds: []string{consts.BasicFieldIssueStatus},
		},
	})
	if tableColumnsResp.Failure() {
		log.Errorf("[GetIssueAllStatus] projectIds: %s, projectObjectTypeIds: %s, err: %v",
			json.ToJsonIgnoreError(projectIds), json.ToJsonIgnoreError(tableIds), tableColumnsResp.Error())
		return statusGroupMap, tableColumnsResp.Error()
	}
	for _, tableSchema := range tableColumnsResp.Data.Tables {
		curTableId := tableSchema.TableId
		columns := tableSchema.Columns
		columnMap := GetColumnMapByColumnsForTableMode(columns)
		if statusColumn, ok := columnMap[consts.BasicFieldIssueStatus]; ok {
			statusGroupMap[curTableId] = GetStatusListFromStatusColumn(statusColumn)
		}
	}

	return statusGroupMap, nil
}

// GetIssueAllStatusByType 获取表中的某个类型的状态，按表分组返回
func GetIssueAllStatusByType(orgId int64, tableIds []int64, statusType int) (map[int64][]status.StatusInfoBo, errs.SystemErrorInfo) {
	resMap := map[int64][]status.StatusInfoBo{}
	statusListGroup, err := GetIssueAllStatus(orgId, []int64{}, tableIds)
	if err != nil {
		log.Errorf("[GetIssueAllStatusByType] tableIds: %s, err: %v", json.ToJsonIgnoreError(tableIds), err)
		return nil, err
	}
	for tableId, statusList := range statusListGroup {
		tmpStatusList := make([]status.StatusInfoBo, 0, len(statusList))
		for _, status := range statusList {
			if status.Type == statusType {
				tmpStatusList = append(tmpStatusList, status)
			}
		}
		resMap[tableId] = tmpStatusList
	}

	return resMap, nil
}

func GetIssueAllStatusIdsByType(status []status.StatusInfoBo, statusType int) []int64 {
	var statusIds []int64
	for _, s := range status {
		if s.Type == statusType {
			statusIds = append(statusIds, s.ID)
		}
	}
	return statusIds
}

func GetIssueAllStatusIdsByTypes(status []status.StatusInfoBo, statusTypes []int) []int64 {
	var statusIds []int64
	for _, s := range status {
		for _, t := range statusTypes {
			if s.Type == t {
				statusIds = append(statusIds, s.ID)
				break
			}
		}
	}
	return statusIds
}

// GetStatusListFromStatusColumn 将表头信息中状态列，获取状态值列表
func GetStatusListFromStatusColumn(column *projectvo.TableColumnData) []status.StatusInfoBo {
	statusBoList := make([]status.StatusInfoBo, 0)
	field := column.Field
	// 任务状态是分组单选类型
	if field.Type != tableV1.ColumnType_groupSelect.String() {
		return statusBoList
	}
	fieldProps := projectvo.FormConfigColumnFieldGroupSelectProps{}
	propsJson := json.ToJsonIgnoreError(field.Props)
	json.FromJson(propsJson, &fieldProps)
	for _, option := range fieldProps.GroupSelect.Options {
		statusBoList = append(statusBoList, status.StatusInfoBo{
			ID:          option.ID,
			Name:        option.Value,
			DisplayName: option.Value,
			Type:        option.ParentID, // 改造后，状态字段类型是分组单选，type 暂且用 option 的 parent 替代
		})
	}

	return statusBoList
}

func GetColumnMapByColumnsForTableMode(columns []*projectvo.TableColumnData) map[string]*projectvo.TableColumnData {
	columnMap := make(map[string]*projectvo.TableColumnData, len(columns))
	for _, item := range columns {
		columnMap[item.Name] = item
	}

	return columnMap
}
