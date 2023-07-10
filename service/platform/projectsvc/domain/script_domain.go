package domain

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetTrashMenuItem() projectvo.ScheduleOrgMenuOption {
	return projectvo.ScheduleOrgMenuOption{
		Name:    "回收站",
		Icon:    "dustbin,dustbin_light",
		LinkUrl: "/trash/:appId",
	}
}

// AddOneNewMenu 给组织新增一个左侧菜单。这次增加回收站
func AddOneNewMenu(orgId int64) (int64, errs.SystemErrorInfo) {
	options := []projectvo.ScheduleOrgMenuOption{
		GetTrashMenuItem(),
	}
	newFolderId := int64(0)
	for _, option := range options {
		resp := CreateExternalApp(orgId, 0, option.Name, option.Icon, option.LinkUrl)
		if resp.Failure() {
			return 0, resp.Error()
		}
		newFolderId = resp.Data.Id
	}

	return newFolderId, nil
}

//func CreateExternalApp(orgId, userId int64, name, icon, linkUrl string) *permissionvo.CreateLessCodeAppResp {
//	appType := consts.LcAppTypeForFolder
//	return appfacade.CreateLessCodeApp(&permissionvo.CreateLessCodeAppReq{
//		AppType:      &appType,
//		OrgId:        &orgId,
//		UserId:       &userId,
//		Name:         &name,
//		AuthType:     2,
//		Icon:         icon,
//		ExternalApp:  1,
//		LinkUrl:      linkUrl,
//		AddAllMember: true,
//	})
//}

// GetLcAppList 查询 app 列表
func GetLcAppList(orgIds []int64, types []int, names []string) ([]po.LcApp, errs.SystemErrorInfo) {
	apps := make([]po.LcApp, 0)
	cond := db.Cond{
		consts.TcOrgId:   db.In(orgIds),
		consts.TcDelFlag: consts.AppIsNoDelete,
	}
	if len(names) > 0 {
		cond[consts.TcName] = db.In(names)
	}
	if len(types) > 0 {
		cond[consts.TcType] = db.In(types)
	}
	if err := mysql.SelectAllByCond(consts.TableLcApp, cond, &apps); err != nil {
		log.Errorf("[GetLcAppList] err: %v", err)
		return apps, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return apps, nil
}

// UpdateLcAppSortBatch 批量更新 app 的排序值
func UpdateLcAppSortBatch(appMapSort map[int64]int64) errs.SystemErrorInfo {
	if len(appMapSort) < 1 {
		return nil
	}
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Errorf("[UpdateLcAppSortBatch] mysql.GetConnect err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	conn.SetLogging(true)
	appIds := make([]int64, 0)
	// 不同 appId 更新的 sort 不同
	sqlPartStr := strings.Builder{}
	for appId, sort := range appMapSort {
		sqlPartStr.WriteString(fmt.Sprintf(" WHEN %d THEN %d ", appId, sort))
		appIds = append(appIds, appId)
	}
	// update 语句不能带有 limit，所以换一种方式。
	_, err = conn.Exec("UPDATE `lc_app` " +
		" SET " + "sort=  " +
		" CASE id " +
		sqlPartStr.String() +
		"  END " +
		fmt.Sprintf(" WHERE id IN (%s)", str.Int64Implode(appIds, ",")) +
		fmt.Sprintf(" LIMIT %d", len(appIds)),
	)
	if err != nil {
		log.Errorf("[UpdateLcAppSortBatch] err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

// GetAppSortMapOfWorkHour 获取工时目录的 appId 对应到 sort 的 map
// orgMapSort：{appId: sort}
func GetAppSortMapOfWorkHour(orgIds []int64) (map[int64]int64, errs.SystemErrorInfo) {
	mapData := make(map[int64]int64, len(orgIds))
	// 查询组织的 app，查询它们的工时（目录）（type=3 & name='工时'）
	apps, err := GetLcAppList(orgIds, []int{consts.LcAppTypeForFolder}, []string{"工时"})
	if err != nil {
		log.Errorf("[AddTrashMenu] orgIds: %s, err：%v", json.ToJsonIgnoreError(orgIds), err)
		return mapData, err
	}
	orgMapSort := make(map[int64]int64, len(apps))
	for _, app := range apps {
		orgMapSort[app.OrgID] = app.Sort
	}

	return orgMapSort, nil
}

// FilterExistTrashOrgArr 过滤掉已经存在“回收站目录”的组织
func FilterExistTrashOrgArr(orgArr []bo.OrganizationBo) ([]bo.OrganizationBo, errs.SystemErrorInfo) {
	filteredData := make([]bo.OrganizationBo, 0)
	orgIds := make([]int64, len(orgArr))
	for i, org := range orgArr {
		orgIds[i] = org.Id
	}
	if len(orgIds) < 1 {
		return filteredData, nil
	}
	apps, err := GetLcAppList(orgIds, []int{consts.LcAppTypeForFolder}, []string{"回收站"})
	if err != nil {
		log.Errorf("[AddTrashMenu] orgIds: %s, err：%v", json.ToJsonIgnoreError(orgIds), err)
		return filteredData, err
	}
	existOrgIds := make([]int64, 0)
	for _, app := range apps {
		existOrgIds = append(existOrgIds, app.OrgID)
	}
	for _, org := range orgArr {
		if exist, _ := slice.Contain(existOrgIds, org.Id); !exist {
			filteredData = append(filteredData, org)
		}
	}

	return filteredData, nil
}

// GetProjectIdColForIssues 获取任务数据中的 projectId 列
func GetProjectIdColForIssues(issues []map[string]interface{}, key string) ([]int64, errs.SystemErrorInfo) {
	var oriErr error
	proIdMap := make(map[int64]struct{}, 0)
	for _, issue := range issues {
		// parse projectId
		projectId := int64(0)
		if f1, ok := issue[key].(float64); ok {
			proIdStr := strconv.FormatFloat(f1, 'f', -1, 64)
			projectId, oriErr = strconv.ParseInt(proIdStr, 10, 64)
			if oriErr != nil {
				log.Errorf("[GetProjectIdColForIssues] err: %v", oriErr)
				return nil, errs.BuildSystemErrorInfo(errs.TypeConvertError, oriErr)
			}
		}
		proIdMap[projectId] = struct{}{}
	}
	projectIds := make([]int64, 0)
	for proId, _ := range proIdMap {
		projectIds = append(projectIds, proId)
	}

	return projectIds, nil
}

// GetProAppIdMap 通过 projectId 获取对应的 appId
func GetProAppIdMap(orgId int64, projectIds []int64) (map[int64]int64, errs.SystemErrorInfo) {
	resMap := make(map[int64]int64, 0)
	if len(projectIds) < 1 {
		return resMap, nil
	}
	proList, err := GetProjectBoList(orgId, projectIds)
	if err != nil {
		log.Errorf("[GetIssueAppIdMap] GetProjectBoList err: %v", err)
		return resMap, err
	}
	for _, pro := range proList {
		resMap[pro.Id] = pro.AppId
	}

	return resMap, nil
}

// GetWrongAppIdsIssue 查询有错误的任务数据
func GetWrongAppIdsIssue(orgId, summaryAppId int64, page, size int) ([]map[string]interface{}, int64, errs.SystemErrorInfo) {
	return nil, 0, nil
	// appId: []int64{0}  && projectId != 0
	//condArr := []*vo.LessCondsData{
	//	{
	//		// 类型(between,equal,un_equal,gt,gte,in,like,lt,lte,not_in,not_like,not_null,is_null,all_in,values_in)
	//		Type:      "un_equal",
	//		Column:    "projectId",
	//		FieldType: nil,
	//		Value:     0,
	//	}, {
	//		Type:      "values_in",
	//		Values:    []interface{}{"0"},
	//		Column:    "appIds",
	//		FieldType: nil,
	//	},
	//}
	//lessReq := vo.LessCondsData{
	//	Type:  "and",
	//	Left:  nil,
	//	Right: nil,
	//	Conds: condArr,
	//}
	//lessReqParam := formvo.LessIssueListReq{
	//	Condition: lessReq,
	//	AppId:     summaryAppId,
	//	OrgId:     orgId,
	//	Page:      int64(page),
	//	Size:      int64(size),
	//	NeedTotal: true,
	//}
	//lessResp := formfacade.LessIssueList(lessReqParam)
	//if lessResp.Failure() {
	//	log.Errorf("[GetWrongAppIdsIssue] err: %v", lessResp.Error())
	//	return nil, 0, lessResp.Error()
	//}
	//if lessResp.Data.Total > 1000 {
	//	log.Infof("[GetWrongAppIdsIssue] 超过 1000 条错误数据的组织 orgId: %d, err issue num: %d", orgId, lessResp.Data.Total)
	//}
	//
	//return lessResp.Data.List, lessResp.Data.Total, nil
}

// UpdateIterStatusToFixed 跑批处理迭代状态，迭代状态更新为固定的状态值
// 取出对应的 status
// 查询 status 对应的对象，获取对应的 type
// 通过 type 映射出目标状态值，
// 更新 iteration 表的 status
// 更新 ppm_pri_iteration_status_relation 表的 status
func UpdateIterStatusToFixed(orgId int64) errs.SystemErrorInfo {
	iterOrder := 3
	orgAllStatusMap2Type, err := GetOneOrgAllStatusMap2Type(orgId)
	if err != nil {
		log.Errorf("[UpdateIterStatusToFixed] err: %v", err)
		return err
	}
	newStatusIdMapByType := make(map[int]int64, 0)
	newStatusList := consts.IterationStatusList
	for _, item := range newStatusList {
		newStatusIdMapByType[item.Type] = item.ID
	}

	// 查询组织下所有的项目，更新项目对应的迭代
	for page := 1; ; page += 1 {
		proBoList, _, err := GetProjectList(0, db.Cond{
			consts.TcOrgId:         orgId,
			consts.TcProjectTypeId: consts.ProjectTypeAgileId,
			consts.TcIsDelete:      consts.AppIsNoDelete,
		}, nil, []*string{str.ToPtr("id asc")}, 500, page)
		if err != nil {
			log.Errorf("[UpdateIterStatusToFixed] err: %v", err)
			return err
		}
		if len(proBoList) < 1 {
			break
		}
		proIds := make([]int64, 0)
		for _, pro := range proBoList {
			proIds = append(proIds, pro.Id)
		}

		iterList, _, err := GetIterationBoList(1, 1000, db.Cond{
			consts.TcProjectId: db.In(proIds),
			consts.TcIsDelete:  consts.AppIsNoDelete,
		}, &iterOrder)
		if err != nil {
			log.Errorf("[UpdateIterStatusToFixed] GetIterationBoList err: %v", err)
			return err
		}
		if len(*iterList) < 1 {
			continue
		}
		iterIds := make([]int64, 0, len(*iterList))
		statusIds := make([]int64, 0)
		for _, iter := range *iterList {
			statusIds = append(statusIds, iter.Status)
			iterIds = append(iterIds, iter.Id)
		}
		// 查询 iteration_status_relation todo
		iterRelations, err := GetIterationRelationStatusInfoBatch(orgId, iterIds)
		if err != nil {
			log.Errorf("[UpdateIterStatusToFixed] GetIterationRelationStatusInfoBatch err: %v", err)
			return err
		}
		for _, oneBo := range iterRelations {
			statusIds = append(statusIds, oneBo.StatusId)
		}
		statusIds = slice.SliceUniqueInt64(statusIds)
		// 旧状态对应的新状态值 map
		oldStatusIdMap2NewStatusId := GetStatusIdMap2NewStatusId(statusIds, orgAllStatusMap2Type, newStatusIdMapByType)
		// update
		if err := ExecUpdateIterStatus(orgId, statusIds, oldStatusIdMap2NewStatusId); err != nil {
			log.Errorf("[UpdateIterStatusToFixed] orgId: %d, statusIds: %s, err: %v", orgId, json.ToJsonIgnoreError(statusIds), err)
		}
	}

	return nil
}

// ExecUpdateIterStatus 请求 MySQL 执行更新
func ExecUpdateIterStatus(orgId int64, oldStatusId []int64, oldStatusIdMap2NewStatusId map[int64]int64) errs.SystemErrorInfo {
	if len(oldStatusIdMap2NewStatusId) < 1 {
		return nil
	}
	conn, dbErr := mysql.GetConnect()
	if dbErr != nil {
		log.Errorf("[ExecUpdateIterStatus] mysql.GetConnect err: %v", dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	sqlPartStr := strings.Builder{}
	for old, newStatusId := range oldStatusIdMap2NewStatusId {
		sqlPartStr.WriteString(fmt.Sprintf(" WHEN %d THEN %d ", old, newStatusId))
	}
	// update 语句不能带有 limit，所以换一种方式。
	sql1 := "UPDATE `" + consts.TableIteration + "` " +
		" SET " + "`status`=  " +
		" CASE `status` " +
		sqlPartStr.String() +
		" else `status` END " +
		fmt.Sprintf(" WHERE `status` IN (%s)", str.Int64Implode(oldStatusId, ",")) +
		fmt.Sprintf(" AND `org_id` = %d", orgId)
	log.Infof("[ExecUpdateIterStatus] sql1: %s", sql1)
	_, dbErr = conn.Exec(sql1)
	if dbErr != nil {
		log.Errorf("[ExecUpdateIterStatus] update ppm_pri_iteration err: %v", dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	// 更新 ppm_pri_iteration_status_relation 表
	sql2 := "UPDATE `" + consts.TableIterationStatusRelation + "` " +
		" SET " + "status_id=  " +
		" CASE status_id " +
		sqlPartStr.String() +
		" else status_id END " +
		fmt.Sprintf(" WHERE status_id IN (%s)", str.Int64Implode(oldStatusId, ",")) +
		fmt.Sprintf(" AND org_id = %d", orgId)
	log.Infof("[ExecUpdateIterStatus] sql2: %s", sql2)
	_, dbErr = conn.Exec(sql2)
	if dbErr != nil {
		log.Errorf("[ExecUpdateIterStatus] update ppm_pri_iteration_status_relation err: %v", dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return nil
}

func GetStatusIdMap2NewStatusId(statusIds []int64, orgAllStatusMap2Type map[int64]int, newStatusIdMapByType map[int]int64) map[int64]int64 {
	map2NewStatus := make(map[int64]int64, 0)
	for _, statusId := range statusIds {
		if statusType, ok := orgAllStatusMap2Type[statusId]; ok {
			map2NewStatus[statusId] = newStatusIdMapByType[statusType]
		}
	}

	return map2NewStatus
}

func GetOneOrgAllStatusMap2Type(orgId int64) (map[int64]int, errs.SystemErrorInfo) {
	resMap := make(map[int64]int, 0)
	env := os.Getenv(consts.RunEnvKey)
	statusList := make([]bo.CacheProcessStatusBo, 0)
	if env == "test" {
		statusList, _ = GetProcessStatusListFromCache()
	} else {
		//statusResp := processfacade.GetProcessStatusBatch(processvo.GetProcessStatusBatchReqVo{
		//	OrgId: orgId,
		//})
		//if statusResp.Failure() {
		//	log.Errorf("[GetOneOrgAllStatusMap2Type] err: %v", statusResp.Error())
		//	return nil, statusResp.Error()
		//}
		//log.Infof("[GetOneOrgAllStatusMap2Type] orgId: %d, status list: %s", orgId, json.ToJsonIgnoreError(statusResp))
		//statusList = *(statusResp.Data)
	}
	for _, item := range statusList {
		resMap[item.StatusId] = item.StatusType
	}

	return resMap, nil
}

// GetProcessStatusListFromCache test for 14396
func GetProcessStatusListFromCache() ([]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	key := "polaris:processsvc:org_14396:process_status_list"
	processStatusListJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	cacheProcessStatusList := &[]bo.CacheProcessStatusBo{}
	if processStatusListJson != "" {
		err = json.FromJson(processStatusListJson, cacheProcessStatusList)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	}

	return *cacheProcessStatusList, nil
}
