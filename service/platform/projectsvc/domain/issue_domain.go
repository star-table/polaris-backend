package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	automationPb "github.com/star-table/interface/golang/automation/v1"

	"github.com/star-table/polaris-backend/common/core/util/jsonx"

	msgPb "github.com/star-table/interface/golang/msg/v1"
	pushPb "github.com/star-table/interface/golang/push/v1"
	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/businees"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/core/util/convert"
	"github.com/star-table/polaris-backend/common/core/util/date"
	"github.com/star-table/polaris-backend/common/core/util/format"
	slice2 "github.com/star-table/polaris-backend/common/core/util/slice"
	int642 "github.com/star-table/polaris-backend/common/core/util/slice/int64"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/bo/status"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/common/model/vo/datacenter"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
	"github.com/star-table/polaris-backend/facade/common/report"
	"github.com/star-table/polaris-backend/facade/formfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
	"github.com/star-table/polaris-backend/facade/tablefacade"
	consts2 "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"github.com/spf13/cast"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetCalendarInfo(orgId, projectId int64) (isOk bool, info bo.CalendarInfo) {
	isOk = false
	projectCalendarInfo, err := GetProjectCalendarInfo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return
	}
	// 判断略显麻烦，先去掉。
	//if !CheckSyncOutCalendarValid(projectCalendarInfo.IsSyncOutCalendar) {
	//	log.Error("项目设置导出日历的配置不合法。")
	//	return
	//}
	if projectCalendarInfo.CalendarId == consts.BlankString {
		log.Info("无对应项目日历或未设置导出日历")
		return
	}
	orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("组织外部信息不存在 %v", err)
		return
	}
	isOk = true
	info.OrgId = orgId
	info.OutOrgId = orgBaseInfo.OutOrgId
	info.CalendarId = projectCalendarInfo.CalendarId
	info.Creator = projectCalendarInfo.Creator
	info.SourceChannel = orgBaseInfo.SourceChannel
	info.SyncCalendarFlag = projectCalendarInfo.IsSyncOutCalendar
	return
}

func DeleteRelationByDeleteMember(tx sqlbuilder.Tx, delMembers []interface{}, projectOwner int64, projectId int64, orgId int64, currentUserId int64) errs.SystemErrorInfo {
	//issueRelation := &[]po.PpmPriIssueRelation{}
	//err := mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
	//	consts.TcOrgId:        orgId,
	//	consts.TcProjectId:    projectId,
	//	consts.TcIsDelete:     consts.AppIsNoDelete,
	//	consts.TcRelationType: consts.IssueRelationTypeOwner,
	//	consts.TcRelationId:   db.In(delMembers),
	//}, issueRelation)
	//
	//if err != nil {
	//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}
	////原有负责人、参与人直接从任务中移除
	//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
	//	consts.TcOrgId:        orgId,
	//	consts.TcProjectId:    projectId,
	//	consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeOwner}),
	//	consts.TcIsDelete:     consts.AppIsNoDelete,
	//	consts.TcRelationId:   db.In(delMembers),
	//}, mysql.Upd{
	//	consts.TcIsDelete:   consts.AppIsDeleted,
	//	consts.TcUpdator:    currentUserId,
	//	consts.TcUpdateTime: time.Now(),
	//})
	//
	//if err != nil {
	//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}
	////任务负责人则移除并调整为项目负责人
	//length := len(*issueRelation)
	//if length > 0 {
	//	issueIds := make([]int64, length)
	//	issuePoInfo := make([]interface{}, length)
	//	for k, v := range *issueRelation {
	//		issueIds = append(issueIds, v.IssueId)
	//		id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueRelation)
	//		if err != nil {
	//			mysql.Rollback(tx)
	//			return errs.BuildSystemErrorInfo(errs.ApplyIdError)
	//		}
	//		issuePoInfo[k] = po.PpmPriIssueRelation{
	//			Id:           id,
	//			OrgId:        orgId,
	//			IssueId:      v.IssueId,
	//			RelationType: consts.IssueRelationTypeOwner,
	//			RelationId:   projectOwner,
	//			Creator:      currentUserId,
	//			CreateTime:   time.Now(),
	//			Updator:      currentUserId,
	//			UpdateTime:   time.Now(),
	//			IsDelete:     consts.AppIsNoDelete,
	//		}
	//	}
	//
	//	err = mysql.TransBatchInsert(tx, &po.PpmPriIssueRelation{}, issuePoInfo)
	//	if err != nil {
	//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//	}
	//	_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
	//		consts.TcId: db.In(issueIds),
	//	}, mysql.Upd{
	//		consts.TcOwner: projectOwner,
	//	})
	//	if err != nil {
	//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//	}
	//}

	//原有参与人直接从任务中移除
	_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeFollower}),
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationId:   db.In(delMembers),
	}, mysql.Upd{
		consts.TcIsDelete:   consts.AppIsDeleted,
		consts.TcUpdator:    currentUserId,
		consts.TcUpdateTime: time.Now(),
	})

	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

// GetIssueFinishedStatForProjects 获取多个项目下任务完成进度。逾期任务数、已完成未完成任务数等
// 参考 `SelectIssueAssignCount` 函数
func GetIssueFinishedStatForProjects(orgId int64, curUserId int64, summaryAppId int64, projectIds []int64) (map[int64]bo.IssueStatistic, errs.SystemErrorInfo) {
	issueStat := map[int64]bo.IssueStatistic{}
	// 从无码接口查询两次：
	// 查询逾期未完成任务数、已完成任务数，并按项目分组
	issueFinishedOrNotRes, err := StatisticProIssueFinishedOrNot(orgId, curUserId, summaryAppId, projectIds)
	if err != nil {
		log.Errorf("[GetIssueFinishedStatForProjects] orgId:%d, StatisticProIssueFinishedOrNot err: %v", orgId, err)
		return issueStat, err
	}
	dataJsonStr := json.ToJsonIgnoreError(issueFinishedOrNotRes)
	statistic := []bo.IssueStatistic{}
	json.FromJson(dataJsonStr, &statistic)
	for _, v := range statistic {
		issueStat[v.ProjectId] = v
	}

	// 查询我关注或者我负责的未完成的任务数，按项目分组
	//relateRes, err := StatisticProIssuesRelateSomeone(orgId, curUserId, summaryAppId, projectIds)
	//if err != nil {
	//	log.Errorf("[GetIssueFinishedStatForProjects] orgId:%d, GetIssueFinishedStatForProjects err: %v", orgId, err)
	//	return issueStat, err
	//}
	//for _, total := range relateRes {
	//	unfinish := total.Total
	//	if temp, ok := issueStat[total.ProjectId]; ok {
	//		// 我相关的未完成的任务统计
	//		temp.RelateUnfinish = unfinish
	//		issueStat[total.ProjectId] = temp
	//	}
	//}
	//log.Infof("[GetIssueFinishedStatForProjects] issueStat: %s", json.ToJsonIgnoreError(issueStat))

	return issueStat, nil
}

// StatisticProIssueFinishedOrNot 查询逾期未完成任务数、已完成任务数，并按项目分组
func StatisticProIssueFinishedOrNot(orgId int64, curUserId int64, summaryAppId int64, projectIds []int64) ([]map[string]interface{}, errs.SystemErrorInfo) {
	nowTimeDateStr := time.Now().Format(consts.AppTimeFormat)
	finishedTypeIds := []int64{consts.StatusTypeComplete}
	unfinishedTypeIds := []int64{consts.StatusTypeNotStart, consts.StatusTypeRunning}
	// data::jsonb->statusType 方式的查询，会将数值转为字符串，因此，in 查询后的列举值必须是带引号的字符串，如：`in ('7', '8')`
	finishedTypeIdsStr := "'" + str.Int64Implode(finishedTypeIds, "','") + "'"
	unfinishedTypeIdsStr := "'" + str.Int64Implode(unfinishedTypeIds, "','") + "'"

	condition := &tablePb.Condition{
		Type: tablePb.ConditionType_and,
		Conditions: GetNoRecycleCondition(
			GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_in, nil, projectIds),
			//GetRowsCondition(consts.BasicFieldPlanEndTime, tablePb.ConditionType_between, nil, []interface{}{consts.BlankTime, nowTimeDateStr}),
		),
	}

	req := &tablePb.ListRawRequest{
		DbType: tablePb.DbType_slave1,
		FilterColumns: []string{
			" count(*) as \"all\" ",
			fmt.Sprintf("\"%s\"", consts.BasicFieldProjectId),
			fmt.Sprintf(" sum(case when \"data\" :: jsonb -> '%s' in (%s) and \"data\" :: jsonb ->> '%s' >'%s' and \"data\" :: jsonb ->> '%s' <'%s' then 1 else 0 end) \"overdue\" ",
				consts.BasicFieldIssueStatusType, unfinishedTypeIdsStr, consts.BasicFieldPlanEndTime, "1970-01-01 00:00:00", consts.BasicFieldPlanEndTime, nowTimeDateStr),
			fmt.Sprintf(" sum(case when \"data\" :: jsonb -> '%s' in (%s) then 1 else 0 end) \"finish\" ", consts.BasicFieldIssueStatusType, finishedTypeIdsStr),
		},
		Condition: condition,
		Groups: []string{
			// group by 用别名也可
			fmt.Sprintf("\"%s\"", consts.BasicFieldProjectId),
		},
		Page: 1,
		Size: int32(len(projectIds) + 10),
	}

	lessResp, err := GetRawRows(orgId, curUserId, req)
	if err != nil {
		log.Errorf("[GetIssueFinishedStatForProjects] orgId:%d, err: %v", orgId, err.Error())
		return nil, err
	}

	return lessResp.Data, nil
}

// StatisticProIssuesRelateSomeone 查询一些项目下，和我相关的任务数统计。查询我关注或者我负责的未完成的任务数，按项目分组
func StatisticProIssuesRelateSomeone(orgId int64, curUserId int64, summaryAppId int64, projectIds []int64) ([]bo.RelateIssueTotal, errs.SystemErrorInfo) {
	//nowTimeDateStr := time.Now().Format(consts.AppTimeFormat)

	values := make([]interface{}, 0, 1)
	uidStr := fmt.Sprintf("U_%d", curUserId)
	values = append(values, uidStr)
	condition := &tablePb.Condition{
		Type: tablePb.ConditionType_and,
		Conditions: GetNoRecycleCondition(
			GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_in, nil, projectIds),
			GetRowsCondition(consts.BasicFieldIssueStatusType, tablePb.ConditionType_in, nil, []interface{}{consts.StatusTypeNotStart, consts.StatusTypeRunning}),
			//GetRowsCondition(consts.BasicFieldPlanEndTime, tablePb.ConditionType_between, nil, []interface{}{consts.BlankTime, nowTimeDateStr}),
			&tablePb.Condition{
				Type: tablePb.ConditionType_or,
				Conditions: []*tablePb.Condition{
					GetRowsCondition(consts.BasicFieldOwnerId, tablePb.ConditionType_values_in, nil, values),
					GetRowsCondition(consts.BasicFieldFollowerIds, tablePb.ConditionType_values_in, nil, values),
				},
			},
		),
	}

	req := &tablePb.ListRawRequest{
		DbType: tablePb.DbType_slave1,
		FilterColumns: []string{
			" count(*) as \"total\" ",
			fmt.Sprintf("\"%s\"", consts.BasicFieldProjectId),
		},
		Condition: condition,
		Groups: []string{
			// group by 用别名也可
			fmt.Sprintf("\"%s\"", consts.BasicFieldProjectId),
		},
		Page:    1,
		Size:    int32(len(projectIds) + 10), // int64(len(projectIds)) 防止为 0，所以 `+ 10`
		TableId: 0,
	}

	lessResp, err := GetRawRows(orgId, curUserId, req)
	if err != nil {
		log.Errorf("[StatisticProIssuesRelateSomeone] orgId:%d, err: %v", orgId, err.Error())
		return nil, err
	}
	dataJsonStr := json.ToJsonIgnoreError(lessResp.Data)
	statistic := []bo.RelateIssueTotal{}
	json.FromJson(dataJsonStr, &statistic)
	log.Infof("[StatisticProIssuesRelateSomeone] query res: %s", dataJsonStr)
	log.Infof("[StatisticProIssuesRelateSomeone] after from json res: %s", json.ToJsonIgnoreError(statistic))

	return statistic, nil
}

//func GetRelateIssueNum(orgId int64, projectIds []int64, currentUserId int64) (map[int64]int64, errs.SystemErrorInfo) {
//	conn, err := mysql.GetConnect()
//	if err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	res := &[]bo.RelateIssueTotal{}
//	_ = conn.Select(db.Raw("project_id, count(distinct issue_id) as total")).From(consts.TableIssueRelation).Where(db.Cond{
//		consts.TcIsDelete:     consts.AppIsNoDelete,
//		consts.TcOrgId:        orgId,
//		consts.TcProjectId:    db.In(projectIds),
//		consts.TcRelationId:   currentUserId,
//		consts.TcRelationType: db.In([]int{consts.IssueRelationTypeOwner, consts.IssueRelationTypeFollower}),
//	}).GroupBy(consts.TcProjectId).All(res)
//	stat := map[int64]int64{}
//	for _, total := range *res {
//		stat[total.ProjectId] = total.Total
//	}
//
//	return stat, nil
//}

//// 通过条件获取任务 id 列表
//func GetIssueIdListByCond(orgId int64, issueCond db.Cond) ([]*bo.IssueIdBo, errs.SystemErrorInfo) {
//	conn, err := mysql.GetConnect()
//	if err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	defaultCond := db.Cond{
//		consts.TcOrgId:    orgId,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}
//	if len(issueCond) > 0 {
//		for key, val := range issueCond {
//			defaultCond[key] = val
//		}
//	}
//	res := []*bo.IssueIdBo{}
//	_ = conn.Select(db.Raw("id")).From(consts.TableIssue).Where(defaultCond).All(&res)
//	return res, nil
//}

// GetIssueIdList 通过无码查询任务 id 数据
func GetIssueIdList(orgId int64, conds []*tablePb.Condition, page, size int64) ([]*bo.IssueIdData, errs.SystemErrorInfo) {
	list, err := GetIssueInfosMapLc(orgId, 0, &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: conds,
	}, []string{lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId)}, page, size)
	if err != nil {
		return nil, err
	}

	issueIds := make([]*bo.IssueIdData, 0, len(list))
	if len(list) == 0 {
		return issueIds, nil
	}

	err2 := copyer.Copy(list, &issueIds)
	if err2 != nil {
		log.Errorf("[GetIssueIdList] json转换错误, err:%v", err2)
		return nil, errs.JSONConvertError
	}
	return issueIds, nil
}

func GetRunningIterationStat(projectIds []int64, orgId int64) (map[int64]bo.IssueStatistic, errs.SystemErrorInfo) {
	//获取迭代进行中状态的id集合
	//iterStatusIds := make([]int64, 0, 1)
	//iterStatusIds = append(iterStatusIds, consts.StatusRunning.ID)
	iterStatusId := consts.StatusRunning.ID

	iterations := &[]po.PpmPriIteration{}
	err1 := mysql.SelectAllByCond(consts.TableIteration, db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcStatus:    iterStatusId,
		consts.TcId:        db.In(db.Raw("select max(id) from ppm_pri_iteration where org_id=? and status = ? and is_delete = 2 GROUP BY project_id", orgId, iterStatusId)),
		consts.TcProjectId: db.In(projectIds),
	}, iterations)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.MysqlOperateError
	}

	result := map[int64]bo.IssueStatistic{}
	if len(*iterations) == 0 {
		return result, nil
	}
	allIterationIds := []int64{}
	midInfo := map[int64]bo.IssueStatistic{}
	for _, iteration := range *iterations {
		allIterationIds = append(allIterationIds, iteration.Id)
		midInfo[iteration.Id] = bo.IssueStatistic{
			ProjectId:     iteration.ProjectId,
			IterationId:   iteration.Id,
			IterationName: iteration.Name,
		}
	}
	// 直接根据迭代筛选，而非 projectId
	statInfo, statErr := GetIssueInfoByIterationForLc(allIterationIds, orgId, 0)
	if statErr != nil {
		log.Error(statErr)
		return nil, statErr
	}

	for _, statistic := range statInfo {
		if _, ok := midInfo[statistic.IterationId]; ok {
			temp := bo.IssueStatistic{
				ProjectId:     midInfo[statistic.IterationId].ProjectId,
				IterationId:   midInfo[statistic.IterationId].IterationId,
				IterationName: midInfo[statistic.IterationId].IterationName,
				Overdue:       statistic.Overdue,
				All:           statistic.All,
				Finish:        statistic.Finish,
			}
			result[midInfo[statistic.IterationId].ProjectId] = temp
		}
	}

	return result, nil
}

// GetIssueInfoByIterationForLc 统计迭代下完成的、逾期的任务数，按迭代分组，通过无码接口查询。
func GetIssueInfoByIterationForLc(iterationIds []int64, orgId int64, projectId int64) (map[int64]bo.IssueStatistic, errs.SystemErrorInfo) {
	issueStat := map[int64]bo.IssueStatistic{}
	statistic := make([]*bo.IssueStatistic, 0)
	finishedTypeIds := []int64{consts.StatusTypeComplete}
	unfinishedTypeIds := []int64{consts.StatusTypeNotStart, consts.StatusTypeRunning}
	// data::jsonb->statusType 方式的查询，会将数值转为字符串，因此，in 查询后的列举值必须是带引号的字符串，如：`in ('7', '8')`
	finishedTypeIdsStr := "'" + str.Int64Implode(finishedTypeIds, "','") + "'"
	unfinishedTypeIdsStr := "'" + str.Int64Implode(unfinishedTypeIds, "','") + "'"
	nowTimeDateStr := time.Now().Format(consts.AppTimeFormat)

	condition := &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: GetNoRecycleCondition(),
	}

	if projectId != 0 {
		condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_equal, projectId, nil))
		condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldIterationId, tablePb.ConditionType_in, nil, iterationIds))
	} else {
		condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldIterationId, tablePb.ConditionType_in, nil, iterationIds))
	}

	req := &tablePb.ListRawRequest{
		DbType: tablePb.DbType_slave1,
		FilterColumns: []string{
			"count(*) as all",
			fmt.Sprintf("\"data\" :: jsonb -> '%s' \"%s\"", consts.BasicFieldIterationId, consts.BasicFieldIterationId),
			fmt.Sprintf(" sum(case when \"data\" :: jsonb -> '%s' in (%s) and \"data\" :: jsonb ->> '%s' >'%s' and \"data\" :: jsonb ->> '%s' <'%s' then 1 else 0 end) \"overdue\" ",
				consts.BasicFieldIssueStatusType, unfinishedTypeIdsStr, consts.BasicFieldPlanEndTime, consts.BlankTime, consts.BasicFieldPlanEndTime, nowTimeDateStr),
			fmt.Sprintf(" sum(case when \"data\" :: jsonb -> '%s' in (%s) then 1 else 0 end) \"finish\" ", consts.BasicFieldIssueStatusType, finishedTypeIdsStr),
		},
		Condition: condition,
		Groups: []string{
			"\"" + consts.BasicFieldIterationId + "\"",
		},
	}

	lessResp, err := GetRawRows(orgId, 0, req)
	if err != nil {
		log.Errorf("[GetIssueInfoByIterationForLc] orgId:%d, projectId:%d, LessIssueList failure:%v", orgId, projectId, err)
		return nil, err
	}

	err2 := copyer.Copy(lessResp.Data, &statistic)
	if err2 != nil {
		log.Errorf("[GetIssueInfoByIterationForLc] orgId:%d, projectId:%d, Copy failure:%v", orgId, projectId, err2)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err2)
	}
	for _, v := range statistic {
		issueStat[v.IterationId] = *v
	}

	return issueStat, nil
}

// GetIssueAndChildrenIds 查询任务以及所有子任务的任务ID
func GetIssueAndChildrenIds(orgId int64, issueIds []int64) ([]int64, errs.SystemErrorInfo) {
	if len(issueIds) == 0 {
		return nil, nil
	}

	allIds := issueIds

	// 首先查找有子任务的任务，即父任务
	list, errSys := GetIssueInfosMapLc(orgId, 0,
		GetRowsCondition(consts.BasicFieldParentId, tablePb.ConditionType_in, nil, issueIds),
		[]string{"DISTINCT " + lc_helper.ConvertToFilterColumn(consts.BasicFieldParentId)}, -1, -1, true)
	if errSys != nil {
		return nil, errSys
	}
	if len(list) > 0 {
		parentIds := make([]*bo.ParentIdData, 0, len(list))
		if err := copyer.Copy(list, &parentIds); err != nil {
			log.Errorf("[GetIssueAndChildrenIds] json转换错误, err:%v", err)
			return nil, errs.JSONConvertError
		}

		// 根据path模糊查找所有子任务
		conditions := make([]*tablePb.Condition, 0, len(parentIds))
		for _, parentId := range parentIds {
			conditions = append(conditions, GetRowsCondition(consts.BasicFieldPath, tablePb.ConditionType_like, fmt.Sprintf(",%d,", parentId.Id), nil))
		}
		issueIdList, errSys := GetIssueIdList(orgId, []*tablePb.Condition{{Type: tablePb.ConditionType_or, Conditions: conditions}}, -1, -1)
		if errSys != nil {
			log.Errorf("[GetIssueAndChildrenIds] GetIssueIdList err:%v, orgId:%v, issueIds:%v", errSys, orgId, issueIds)
			return nil, errSys
		}

		for _, child := range issueIdList {
			allIds = append(allIds, child.Id)
		}
		allIds = slice.SliceUniqueInt64(allIds)
	}

	return allIds, nil
}

func DeleteIssue(issueBo *bo.IssueBo, operatorId int64, sourceChannel string, takeChildren *bool) ([]int64, errs.SystemErrorInfo) {
	orgId := issueBo.OrgId
	issueId := issueBo.Id

	beforeParticipantIds, err1 := GetIssueRelationIdsByRelateType(orgId, issueId, consts.IssueRelationTypeParticipant)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}
	beforeFollowerIds, err1 := GetIssueRelationIdsByRelateType(orgId, issueId, consts.IssueRelationTypeFollower)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}
	appId, appIdErr := GetAppIdFromProjectId(orgId, issueBo.ProjectId)
	if appIdErr != nil {
		log.Error(appIdErr)
		return nil, appIdErr
	}

	var updateBatchReqs []*formvo.LessUpdateIssueBatchReq

	//issueChild := &[]po.PpmPriIssue{}
	//err := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
	//	consts.TcOrgId:    orgId,
	//	consts.TcPath:     db.Like(fmt.Sprintf("%s%d,%s", issueBo.Path, issueId, "%")),
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//}, issueChild)

	condition := &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: GetNoRecycleCondition(GetRowsCondition(consts.BasicFieldPath, tablePb.ConditionType_like, fmt.Sprintf(",%d,", issueId), nil)),
	}
	filterColumns := []string{
		lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
	}
	issueChildInfos, errSys := GetIssueInfosMapLc(orgId, operatorId, condition, filterColumns, -1, -1)
	if errSys != nil {
		log.Errorf("[DeleteIssue] GetIssueInfosMapLc err:%v, orgId:%v, issueId:%v", errSys, orgId, issueId)
		return nil, errSys
	}
	issueChildIds := make([]int64, 0, len(issueChildInfos))
	for _, data := range issueChildInfos {
		issueChildIds = append(issueChildIds, cast.ToInt64(data[consts.BasicFieldIssueId]))
	}

	summaryAppId, errSys := GetOrgSummaryAppId(orgId)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}

	// 拉取原始数据
	condition = &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: GetNoRecycleCondition(GetRowsCondition(consts.BasicFieldIssueId, tablePb.ConditionType_equal, issueId, nil)),
	}
	issueDatas, errSys := GetIssueInfosMapLc(orgId, operatorId, condition, nil, -1, -1)
	if errSys != nil || len(issueDatas) == 0 {
		log.Errorf("[DeleteIssue] GetIssueInfosMapLc err:%v, orgId:%v, issueId:%v", errSys, orgId, issueId)
		return nil, errSys
	}
	issueData := issueDatas[0]

	// 获取表头
	fromTableSchema, errSys := GetTableColumnsMap(orgId, issueBo.TableId, nil, true)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}
	var allRelatingIssueIds []int64
	for columnId, column := range fromTableSchema {
		if column.Field.Type == tablePb.ColumnType_relating.String() ||
			column.Field.Type == tablePb.ColumnType_baRelating.String() ||
			column.Field.Type == tablePb.ColumnType_singleRelating.String() {
			if v, ok := issueData[columnId]; ok {
				oldRelating := &bo.RelatingIssue{}
				jsonx.Copy(v, oldRelating)
				allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkTo)...)
				allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkFrom)...)
			}
		}
	}
	allRelatingIssueIds = slice.SliceUniqueInt64(allRelatingIssueIds)
	if len(allRelatingIssueIds) > 0 {
		updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
			OrgId:  orgId,
			AppId:  summaryAppId, // 目标项目appId
			UserId: operatorId,
			Condition: vo.LessCondsData{
				Type: consts.ConditionAnd,
				Conds: []*vo.LessCondsData{
					{
						Type:   consts.ConditionEqual,
						Value:  orgId,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
					},
					{
						Type:   consts.ConditionIn,
						Values: allRelatingIssueIds,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
					},
				},
			},
			Sets: []datacenter.Set{
				{
					Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldUpdateTime),
					Value:           "NOW()",
					Type:            consts.SetTypeNormal,
					Action:          consts.SetActionSet,
					WithoutPretreat: true,
				},
			},
		})
	}

	conn, err := mysql.GetConnect()
	if err != nil {
		log.Errorf(consts.DBOpenErrorSentence, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	tx, err := conn.NewTx(nil)
	if err != nil {
		log.Errorf(consts.TxOpenErrorSentence, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	defer mysql.Close(conn, tx)

	allIssueId := []int64{issueId}
	childIssueIds := []int64{}
	if takeChildren == nil || *takeChildren {
		// 默认删除子任务
		//for _, issue := range *issueChild {
		//	allIssueId = append(allIssueId, issue.Id)
		//	childIssueIds = append(childIssueIds, issue.Id)
		//}
		allIssueId = append(allIssueId, issueChildIds...)
		childIssueIds = append(childIssueIds, issueChildIds...)
	} else {
		newPath := "0,"

		// 如果不删除子任务，就把顶级子任务变为父任务
		likeCond := fmt.Sprintf("%s%d,%s", issueBo.Path, issueId, "%")
		//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcPath:     db.Like(likeCond),
		//	consts.TcOrgId:    orgId,
		//	consts.TcParentId: issueId,
		//}, mysql.Upd{
		//	consts.TcParentId: 0,
		//	consts.TcPath:     db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, issueBo.Path, newPath)),
		//	consts.TcUpdator:  operatorId,
		//})
		//if err != nil {
		//	log.Error(err)
		//	return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		//}

		// 子任务的子任务，更新path
		//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcPath:     db.Like(likeCond),
		//	consts.TcOrgId:    orgId,
		//	consts.TcParentId: db.NotEq(issueId),
		//}, mysql.Upd{
		//	consts.TcPath:    db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, issueBo.Path, newPath)),
		//	consts.TcUpdator: operatorId,
		//})
		//if err != nil {
		//	log.Error(err)
		//	return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		//}

		// 无码的 刷顶级子任务的parentId
		updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
			OrgId:  issueBo.OrgId,
			AppId:  appId,
			UserId: operatorId,
			Condition: vo.LessCondsData{
				Type: consts.ConditionAnd,
				Conds: []*vo.LessCondsData{
					{
						Type:   consts.ConditionEqual,
						Value:  issueBo.OrgId,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
					},
					{
						Type:   consts.ConditionEqual,
						Value:  issueId,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldParentId),
					},
				},
			},
			Sets: []datacenter.Set{
				{
					Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldParentId),
					Value:           0,
					Type:            consts.SetTypeNormal,
					Action:          consts.SetActionSet,
					WithoutPretreat: false,
				},
			},
		})
		// 无码的 刷子任务的path
		updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
			OrgId:  issueBo.OrgId,
			AppId:  appId,
			UserId: operatorId,
			Condition: vo.LessCondsData{
				Type: consts.ConditionAnd,
				Conds: []*vo.LessCondsData{
					{
						Type:   consts.ConditionEqual,
						Value:  issueBo.OrgId,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
					},
					{
						Type:   consts.ConditionLike,
						Value:  likeCond,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
					},
				},
			},
			Sets: []datacenter.Set{
				{
					Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
					Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, issueBo.Path, newPath),
					Type:            consts.SetTypeNormal,
					Action:          consts.SetActionSet,
					WithoutPretreat: true,
				},
			},
		})

	}

	recycleVersionId, versionErr := AddRecycleRecord(orgId, operatorId, issueBo.ProjectId, []int64{issueId}, consts.RecycleTypeIssue, tx)
	if versionErr != nil {
		log.Error(versionErr)
		return nil, versionErr
	}
	//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
	//	consts.TcId:       db.In(allIssueId),
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//}, mysql.Upd{
	//	consts.TcUpdator:  operatorId,
	//	consts.TcVersion:  recycleVersionId,
	//	consts.TcIsDelete: consts.AppIsDeleted,
	//})
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssueDetail, db.Cond{
	//	consts.TcOrgId:    orgId,
	//	consts.TcIssueId:  db.In(allIssueId),
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//}, mysql.Upd{
	//	consts.TcUpdator:  operatorId,
	//	consts.TcVersion:  recycleVersionId,
	//	consts.TcIsDelete: consts.AppIsDeleted,
	//})
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	err3 := DeleteAllIssueRelation(tx, operatorId, orgId, allIssueId, recycleVersionId)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}

	// 更新无码
	for _, req := range updateBatchReqs {
		resp := formfacade.LessUpdateIssueBatchRaw(req)
		if resp.Failure() {
			log.Error(resp.Error())
			return nil, resp.Error()
		}
	}

	//无码删除
	resp := formfacade.LessRecycleIssue(formvo.LessRecycleIssueReq{
		AppId:    appId,
		OrgId:    orgId,
		UserId:   operatorId,
		IssueIds: allIssueId,
		TableId:  issueBo.TableId,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}

	err = tx.Commit()
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	//文件关联
	resourceRelationResp := resourcefacade.DeleteResourceRelation(resourcevo.DeleteResourceRelationReqVo{
		OrgId:  orgId,
		UserId: operatorId,
		Input: resourcevo.DeleteResourceRelationData{
			ProjectId:        issueBo.ProjectId,
			IssueIds:         allIssueId,
			RecycleVersionId: recycleVersionId,
			SourceTypes:      []int{consts.OssPolicyTypeIssueResource, consts.OssPolicyTypeLesscodeResource},
		},
	})
	if resourceRelationResp.Failure() {
		log.Error(resourceRelationResp.Error())
		return nil, resourceRelationResp.Error()
	}
	// 工时的删除
	if err := DeleteWorkHourForIssues(orgId, []int64{issueId}, recycleVersionId); err != nil {
		log.Errorf("[DeleteIssue] issueId: %d, DeleteWorkHourForIssues err: %v", issueId, err)
		return nil, err
	}

	asyn.Execute(func() {
		ownerIds, errSys := businees.LcMemberToUserIdsWithError(issueBo.OwnerId)
		if errSys != nil {
			log.Errorf("原始的数据:%v, err:%v", issueBo.OwnerId, err)
			return
		}
		blank := []int64{}
		issueTrendsBo := &bo.IssueTrendsBo{
			PushType:                 consts.PushTypeDeleteIssue,
			OrgId:                    orgId,
			OperatorId:               operatorId,
			DataId:                   issueBo.DataId,
			IssueId:                  issueBo.Id,
			ParentIssueId:            issueBo.ParentId,
			ProjectId:                issueBo.ProjectId,
			TableId:                  issueBo.TableId,
			PriorityId:               issueBo.PriorityId,
			ParentId:                 issueBo.ParentId,
			IssueTitle:               issueBo.Title,
			IssueStatusId:            issueBo.Status,
			BeforeOwner:              ownerIds,
			AfterOwner:               blank,
			BeforeChangeFollowers:    *beforeFollowerIds,
			AfterChangeFollowers:     blank,
			BeforeChangeParticipants: *beforeParticipantIds,
			AfterChangeParticipants:  blank,
			IssueChildren:            childIssueIds,
			SourceChannel:            sourceChannel,
		}
		asyn.Execute(func() {
			PushIssueTrends(issueTrendsBo)
		})
		asyn.Execute(func() {
			PushIssueThirdPlatformNotice(issueTrendsBo)
		})
	})

	return childIssueIds, nil
}

func DeleteIssueBatch(orgId int64, issueAndChildrenBos []*bo.IssueBo, allIssueId []int64, operatorId int64, sourceChannel string, projectId int64, tableId int64) ([]int64, errs.SystemErrorInfo) {
	needChangeParentIds := make([]int64, 0) //需要转化为独立任务（父任务id转为0）
	needChangePath := make([]string, 0)
	pathIssueIdMap := make(map[string]int64, 0)

	allInputIssues := make([]*bo.IssueBo, 0)
	allRemainIssues := make([]*bo.IssueBo, 0)
	issueBos := make([]*bo.IssueBo, 0) //需要移动的所有父任务（包含变更后成为父任务的）
	allParentIds := make([]int64, 0)
	couldDeleteIssueIds := make([]int64, 0)

	var relatingColumnIds []string
	tableSchema, errSys := GetTableColumnsMap(orgId, tableId, nil, true)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}
	for columnId, column := range tableSchema {
		if column.Field.Type == tablePb.ColumnType_relating.String() ||
			column.Field.Type == tablePb.ColumnType_baRelating.String() ||
			column.Field.Type == tablePb.ColumnType_singleRelating.String() {
			relatingColumnIds = append(relatingColumnIds, columnId)
		}
	}
	appId, appIdErr := GetAppIdFromProjectId(orgId, projectId)
	if appIdErr != nil {
		log.Error(appIdErr)
		return nil, appIdErr
	}

	// issueAndChildrenBos 中可能包含只有父任务的id，也可能只有单独的子任务id。因为可以只删除父任务，而不删除子任务。即存在以下几种情况：
	// 	* 父任务、子任务一起删除
	// 	* 只有子任务删除
	// 	* 只有父任务删除
	// 		* 这种情况需要**注意**，删除父任务之后，需要将其子任务的 parentId 改为 0；

	// issueAndChildrenBos 是**有权限**删除的任务列表，如果它为空，则表示没有能被删除的任务。此时可以直接返回。
	if len(issueAndChildrenBos) < 1 {
		log.Infof("[DeleteIssueBatch] 没有可以删除的任务。orgId: %d", orgId)
		return nil, nil
	}
	for _, issueBo := range issueAndChildrenBos {
		if ok, _ := slice.Contain(allIssueId, issueBo.Id); ok {
			allInputIssues = append(allInputIssues, issueBo)
			couldDeleteIssueIds = append(couldDeleteIssueIds, issueBo.Id)
			if issueBo.ParentId != 0 {
				if ok1, _ := slice.Contain(allIssueId, issueBo.ParentId); !ok1 {
					issueBos = append(issueBos, issueBo)
					allParentIds = append(allParentIds, issueBo.Id)
				}
			} else {
				issueBos = append(issueBos, issueBo)
				allParentIds = append(allParentIds, issueBo.Id)
			}
		} else {
			allRemainIssues = append(allRemainIssues, issueBo)
			//需要转化为父任务的逻辑（当前任务不需要被删除，但是父任务需要删除）
			if ok, _ := slice.Contain(allIssueId, issueBo.ParentId); ok {
				needChangeParentIds = append(needChangeParentIds, issueBo.Id)
				//curPath := fmt.Sprintf("%s%d,", issueBo.Path, issueBo.Id)
				curPath := issueBo.Path
				needChangePath = append(needChangePath, curPath)
				pathIssueIdMap[curPath] = issueBo.Id
			}
		}
	}
	if len(couldDeleteIssueIds) < 1 {
		log.Info("没有可以删除的任务。")
		return nil, nil
	}

	var updateBatchReqs []*formvo.LessUpdateIssueBatchReq

	// 需要移动的任务涉及的关联引用的任务ID
	if len(relatingColumnIds) > 0 {
		var allRelatingIssueIds []int64
		for _, issueBo := range allInputIssues {
			for _, columnId := range relatingColumnIds {
				if v, ok := issueBo.LessData[columnId]; ok {
					oldRelating := &bo.RelatingIssue{}
					jsonx.Copy(v, oldRelating)
					allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkTo)...)
					allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkFrom)...)
				}
			}
		}
		allRelatingIssueIds = slice.SliceUniqueInt64(allRelatingIssueIds)
		var relatingIssueIds []int64
		for _, id := range allRelatingIssueIds {
			if ok, _ := slice.Contain(couldDeleteIssueIds, id); !ok {
				relatingIssueIds = append(relatingIssueIds, id)
			}
		}
		if len(relatingIssueIds) > 0 {
			updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
				OrgId:  orgId,
				AppId:  appId,
				UserId: operatorId,
				Condition: vo.LessCondsData{
					Type: consts.ConditionAnd,
					Conds: []*vo.LessCondsData{
						{
							Type:   consts.ConditionEqual,
							Value:  orgId,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
						},
						{
							Type:   consts.ConditionIn,
							Values: relatingIssueIds,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
						},
					},
				},
				Sets: []datacenter.Set{
					{
						Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldUpdateTime),
						Value:           "NOW()",
						Type:            consts.SetTypeNormal,
						Action:          consts.SetActionSet,
						WithoutPretreat: true,
					},
				},
			})
		}
	}

	//查找只需要修改path的子任务
	changePathIssueMap := make(map[string][]int64, 0)
	childIssueIds := map[int64][]int64{}
	changePathIssueBoMap := make(map[string][]*bo.IssueBo, 0)
	needChangePathMap := make(map[string]int, 0)
	for _, s := range needChangePath {
		needChangePathMap[s] = strings.Count(s, ",")
	}
	//从最长的path去匹配，用最长的去替换path
	for _, issueBo := range allInputIssues {
		if ok1, _ := slice.Contain(allIssueId, issueBo.ParentId); ok1 {
			max := 0
			curPath := ""
			for s, i := range needChangePathMap {
				if strings.Contains(issueBo.Path, s) {
					if i > max {
						curPath = s
						max = i
					}
				}
			}
			if curPath != "" {
				changePathIssueMap[curPath] = append(changePathIssueMap[curPath], issueBo.Id)
				changePathIssueBoMap[curPath] = append(changePathIssueBoMap[curPath], issueBo)
				if parentId, ok := pathIssueIdMap[curPath]; ok {
					childIssueIds[parentId] = append(childIssueIds[parentId], issueBo.Id)
				}
			}
		}
	}

	for _, issueBo := range allRemainIssues {
		if ok1, _ := slice.Contain(allIssueId, issueBo.ParentId); !ok1 {
			max := 0
			curPath := ""
			for s, i := range needChangePathMap {
				if strings.Contains(issueBo.Path, s) {
					if i > max {
						curPath = s
						max = i
					}
				}
			}
			if curPath != "" {
				changePathIssueMap[curPath] = append(changePathIssueMap[curPath], issueBo.Id)
			}
		}
	}
	log.Infof("[DeleteIssueBatch] needChangePathMap: %v, changePathIssueMap: %v", json.ToJsonIgnoreError(needChangePathMap), json.ToJsonIgnoreError(changePathIssueMap))

	childrenMap := map[int64][]*bo.IssueBo{}
	for s, int64s := range changePathIssueBoMap {
		idArr := strings.Split(s[0:len(s)-1], ",")
		if len(idArr) >= 2 { // fix panic
			parentId, err := strconv.ParseInt(idArr[len(idArr)-2], 10, 64)
			if err != nil {
				log.Error(err)
				return nil, errs.ParamError
			}
			childrenMap[parentId] = int64s
		}
	}

	conn, err := mysql.GetConnect()
	if err != nil {
		log.Errorf(consts.DBOpenErrorSentence, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	tx, err := conn.NewTx(nil)
	if err != nil {
		log.Errorf(consts.TxOpenErrorSentence, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	defer mysql.Close(conn, tx)

	newPath := "0,"
	updForm := make([]map[string]interface{}, 0)
	if len(needChangeParentIds) > 0 {
		//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcOrgId: orgId,
		//	consts.TcId:    db.In(needChangeParentIds),
		//}, mysql.Upd{
		//	consts.TcParentId: 0,
		//	consts.TcPath:     newPath,
		//})
		//if err != nil {
		//	log.Error(err)
		//	return nil, errs.MysqlOperateError
		//}
		// 刷无码
		for _, id := range needChangeParentIds {
			updForm = append(updForm, map[string]interface{}{
				consts.BasicFieldIssueId:  id,
				consts.BasicFieldParentId: 0,
				consts.BasicFieldPath:     newPath,
			})
		}

		for i, s := range changePathIssueMap {
			//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
			//	consts.TcOrgId: orgId,
			//	consts.TcId:    db.In(s),
			//}, mysql.Upd{
			//	consts.TcPath: db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, i, newPath)),
			//})
			//if err != nil {
			//	log.Error(err)
			//	return nil, errs.MysqlOperateError
			//}

			// 刷无码
			updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
				OrgId:  orgId,
				AppId:  appId,
				UserId: operatorId,
				Condition: vo.LessCondsData{
					Type: consts.ConditionAnd,
					Conds: []*vo.LessCondsData{
						{
							Type:   consts.ConditionEqual,
							Value:  orgId,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
						},
						{
							Type:   consts.ConditionIn,
							Values: s,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
						},
					},
				},
				Sets: []datacenter.Set{
					{
						Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
						Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, i, newPath),
						Type:            consts.SetTypeNormal,
						Action:          consts.SetActionSet,
						WithoutPretreat: true,
					},
				},
			})
		}
	}

	// 无码更新
	resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
		AppId:   appId,
		OrgId:   orgId,
		UserId:  operatorId,
		TableId: tableId,
		Form:    updForm,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}

	// 无码更新
	for _, req := range updateBatchReqs {
		batchResp := formfacade.LessUpdateIssueBatchRaw(req)
		if batchResp.Failure() {
			log.Error(batchResp.Error())
			return nil, batchResp.Error()
		}
	}

	recycleVersionId, versionErr := AddRecycleRecord(orgId, operatorId, projectId, allParentIds, consts.RecycleTypeIssue, tx)
	if versionErr != nil {
		log.Error(versionErr)
		return nil, versionErr
	}
	err2 := deleteIssuesAndRelation(orgId, projectId, appId, tableId, operatorId, recycleVersionId, couldDeleteIssueIds, tx)
	if err2 != nil {
		return nil, err2
	}

	err = tx.Commit()
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	//文件关联
	resourceRelationResp := resourcefacade.DeleteResourceRelation(resourcevo.DeleteResourceRelationReqVo{
		OrgId:  orgId,
		UserId: operatorId,
		Input: resourcevo.DeleteResourceRelationData{
			ProjectId:        projectId,
			IssueIds:         couldDeleteIssueIds,
			RecycleVersionId: recycleVersionId,
			SourceTypes:      []int{consts.OssPolicyTypeIssueResource, consts.OssPolicyTypeLesscodeResource},
		},
	})
	if resourceRelationResp.Failure() {
		log.Error(resourceRelationResp.Error())
		return nil, resourceRelationResp.Error()
	}
	// 删除工时
	if err := DeleteWorkHourForIssues(orgId, couldDeleteIssueIds, recycleVersionId); err != nil {
		log.Errorf("[DeleteIssueBatch] orgId: %d, err: %v", orgId, err)
		return nil, err
	}

	asyn.Execute(func() {
		beforeFollowerIds, err1 := GetRelationInfoByIssueIds(allParentIds, []int{consts.IssueRelationTypeFollower})
		if err1 != nil {
			log.Error(err1)
			return
		}
		beforeOwnerIds, err1 := GetRelationInfoByIssueIds(allParentIds, []int{consts.IssueRelationTypeOwner})
		if err1 != nil {
			log.Error(err1)
			return
		}
		ownerMap := map[int64][]int64{}
		for _, ownerId := range beforeOwnerIds {
			ownerMap[ownerId.IssueId] = append(ownerMap[ownerId.IssueId], ownerId.RelationId)
		}
		followerMap := map[int64][]int64{}
		for _, id := range beforeFollowerIds {
			followerMap[id.IssueId] = append(followerMap[id.IssueId], id.RelationId)
		}
		blank := []int64{}
		for _, issueBo := range issueBos {
			issueTrendsBo := &bo.IssueTrendsBo{
				PushType:             consts.PushTypeDeleteIssue,
				OrgId:                orgId,
				OperatorId:           operatorId,
				DataId:               issueBo.DataId,
				IssueId:              issueBo.Id,
				ParentIssueId:        issueBo.ParentId,
				ProjectId:            issueBo.ProjectId,
				TableId:              tableId,
				PriorityId:           issueBo.PriorityId,
				ParentId:             issueBo.ParentId,
				IssueTitle:           issueBo.Title,
				IssueStatusId:        issueBo.Status,
				AfterOwner:           blank,
				AfterChangeFollowers: blank,
				SourceChannel:        sourceChannel,
			}
			if _, ok := ownerMap[issueBo.Id]; ok {
				issueTrendsBo.BeforeOwner = ownerMap[issueBo.Id]
			}
			if _, ok := childIssueIds[issueBo.Id]; ok {
				issueTrendsBo.IssueChildren = childIssueIds[issueBo.Id]
			}
			if _, ok := followerMap[issueBo.Id]; ok {
				issueTrendsBo.BeforeChangeFollowers = followerMap[issueBo.Id]
			}

			asyn.Execute(func() {
				PushIssueTrends(issueTrendsBo)
			})
			asyn.Execute(func() {
				PushIssueThirdPlatformNotice(issueTrendsBo)
			})
		}
	})

	return couldDeleteIssueIds, nil
}

func deleteIssuesAndRelation(orgId, projectId, appId, tableId, operatorId, recycleVersionId int64, issueIds []int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
	//	consts.TcOrgId:    orgId,
	//	consts.TcId:       db.In(issueIds),
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//}, mysql.Upd{
	//	consts.TcUpdator:  operatorId,
	//	consts.TcVersion:  recycleVersionId,
	//	consts.TcIsDelete: consts.AppIsDeleted,
	//})
	//if err != nil {
	//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	//issueDetail := &po.PpmPriIssueDetail{}
	//_, err = mysql.TransUpdateSmartWithCond(tx, issueDetail.TableName(), db.Cond{
	//	consts.TcOrgId:    orgId,
	//	consts.TcIssueId:  db.In(issueIds),
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//}, mysql.Upd{
	//	consts.TcUpdator:  operatorId,
	//	consts.TcVersion:  recycleVersionId,
	//	consts.TcIsDelete: consts.AppIsDeleted,
	//})
	//if err != nil {
	//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	err3 := DeleteAllIssueRelation(tx, operatorId, orgId, issueIds, recycleVersionId)
	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}

	//无码回收
	recycleResp := formfacade.LessRecycleIssue(formvo.LessRecycleIssueReq{
		AppId:    appId,
		OrgId:    orgId,
		UserId:   operatorId,
		IssueIds: issueIds,
		TableId:  tableId,
	})
	if recycleResp.Failure() {
		log.Error(recycleResp.Error())
	}

	return nil
}

// DeleteTableIssues 删除一个table，把所有任务都删除，不进入回收站
func DeleteTableIssues(orgId, operatorId, projectId, appId, tableId int64, templateFlag int) ([]int64, errs.SystemErrorInfo) {
	conditions := []*tablePb.Condition{
		GetRowsCondition(consts.BasicFieldTableId, tablePb.ConditionType_equal, cast.ToString(tableId), nil),
	}
	list, err := GetIssueIdList(orgId, conditions, 0, 0)
	if err != nil {
		return nil, err
	}

	issueIds := make([]int64, 0, len(list))
	for _, data := range list {
		issueIds = append(issueIds, data.Id)
	}
	// 删除附件和任务关系、删除附件关联关系、删除回收站附件
	errSys := DeleteAttachmentsForOneTable(orgId, operatorId, projectId, tableId)
	if errSys != nil {
		log.Errorf("[DeleteTableIssues] DeleteAttachmentsForOneTable err:%v, orgId:%v, userId:%v, tableId:%v",
			errSys, orgId, operatorId, tableId)
		return nil, errSys
	}

	err2 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(issueIds) > 0 {
			err := deleteIssuesAndRelation(orgId, projectId, appId, tableId, operatorId, 0, issueIds, tx)
			if err != nil {
				return err
			}

			err = deleteRecycleRecord(orgId, operatorId, projectId, tableId, issueIds, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err2 != nil {
		log.Errorf("[DeleteTableIssues] err:%v", err2)
		return nil, errs.MysqlOperateError
	}

	if templateFlag == consts.TemplateFalse {
		errCache := SetDeletedIssuesNum(orgId, int64(len(issueIds)))
		if errCache != nil {
			log.Errorf("[DeleteTableIssues]SetDeletedIssuesNum err:%v", errCache)
		}
	}

	return issueIds, nil
}

func deleteRecycleRecord(orgId, userId, projectId, tableId int64, issueIds []int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	_, recordsErr := mysql.TransUpdateSmartWithCond(tx, consts.TableRecycleBin, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: consts.RecycleTypeIssue,
		consts.TcRelationId:   db.In(issueIds),
		consts.TcOrgId:        orgId,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  userId,
	})
	if recordsErr != nil {
		log.Errorf("[deleteRecycleRecord] UpdateSmartWithCond failed:%v", recordsErr)
		return errs.MysqlOperateError
	}

	return nil
}

func GetHomeIssuePriorityInfoBo(orgId, tableId, priorityId int64) (*bo.HomeIssuePriorityInfoBo, errs.SystemErrorInfo) {
	if tableId == 0 || priorityId == 0 {
		return &bo.HomeIssuePriorityInfoBo{}, nil
	}

	priority, err := GetPriorityById(orgId, tableId, priorityId)
	if err != nil {
		log.Error(err)
		//return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	} else {
		priorityInfo := ConvertPriorityCacheInfoToHomeIssuePriorityInfo(*priority)
		return priorityInfo, nil
	}
	return &bo.HomeIssuePriorityInfoBo{}, nil
}

func GetPriorityById(orgId, tableId, priorityId int64) (*bo.PriorityBo, errs.SystemErrorInfo) {
	prioritys, err := GetPriorityListFromTableSchema(orgId, tableId)
	if err != nil {
		return nil, err
	}
	for _, p := range prioritys {
		if p.Id == priorityId {
			return p, nil
		}
	}

	return nil, errs.BuildSystemErrorInfo(errs.ObjectRecordNotFoundError, fmt.Errorf("tableId:%v, priorityId:%v, list:%v", tableId, priorityId, prioritys))
}

// GetPriorityListFromTableSchema 从 table 的 schema 中查询出优先级列表
func GetPriorityListFromTableSchema(orgId int64, tableId int64) ([]*bo.PriorityBo, errs.SystemErrorInfo) {
	columnsInfo, err := GetTableColumnConfig(orgId, tableId, []string{consts.BasicFieldPriority}, false)
	if err != nil {
		log.Errorf("[GetPriorityListFromTableSchema] err: %v", err)
		return nil, err
	}
	columnMap := make(map[string]*projectvo.TableColumnData, 0)
	for _, column := range columnsInfo.Columns {
		columnMap[column.Name] = column
	}
	returnRes := make([]*bo.PriorityBo, 0, 6)
	if priorityColumn, ok := columnMap[consts.BasicFieldPriority]; ok {
		propsMap := priorityColumn.Field.Props
		propsJson := json.ToJsonIgnoreError(propsMap)
		propsObj := projectvo.FormConfigColumnFieldMultiselectProps{}
		if err := json.FromJson(propsJson, &propsObj); err != nil {
			log.Errorf("[GetPriorityListFromTableSchema] err: %v", err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}

		for _, option := range propsObj.Select.Options {
			tmpOpId := cast.ToInt64(option.Id)
			returnRes = append(returnRes, &bo.PriorityBo{
				Id:       tmpOpId,
				OrgId:    orgId,
				LangCode: "",
				Name:     option.Value,
				BgStyle:  option.Color,
			})
		}
	}

	return returnRes, nil
}

func ConvertPriorityCacheInfoToHomeIssuePriorityInfo(priorityCacheInfo bo.PriorityBo) *bo.HomeIssuePriorityInfoBo {
	priorityInfo := &bo.HomeIssuePriorityInfoBo{}
	priorityInfo.ID = priorityCacheInfo.Id
	priorityInfo.Name = priorityCacheInfo.Name
	priorityInfo.FontStyle = priorityCacheInfo.FontStyle
	priorityInfo.BgStyle = priorityCacheInfo.BgStyle
	return priorityInfo
}

func ConvertStatusInfoToHomeIssueStatusInfo(statusInfo bo.CacheProcessStatusBo) *status.StatusInfoBo {
	homeStatusInfo := &status.StatusInfoBo{}
	homeStatusInfo.ID = statusInfo.StatusId
	homeStatusInfo.Name = statusInfo.Name
	homeStatusInfo.Type = statusInfo.StatusType
	homeStatusInfo.BgStyle = statusInfo.BgStyle
	homeStatusInfo.FontStyle = statusInfo.FontStyle
	homeStatusInfo.Sort = statusInfo.Sort
	return homeStatusInfo
}

func GetHomeProjectInfoBo(orgId, projectId int64) (*bo.HomeIssueProjectInfoBo, errs.SystemErrorInfo) {
	projectCacheInfo, err := LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	projectInfo := ConvertProjectCacheInfoToHomeIssueProjectInfo(*projectCacheInfo)
	return projectInfo, nil
}

func ConvertProjectCacheInfoToHomeIssueProjectInfo(projectCacheInfo bo.ProjectAuthBo) *bo.HomeIssueProjectInfoBo {
	projectInfo := &bo.HomeIssueProjectInfoBo{}
	projectInfo.Name = projectCacheInfo.Name
	projectInfo.ID = projectCacheInfo.Id
	projectInfo.IsFilling = projectCacheInfo.IsFilling
	projectInfo.ProjectTypeId = projectCacheInfo.ProjectType
	return projectInfo
}

func ConvertBaseUserInfoToHomeIssueOwnerInfo(baseUserInfo bo.BaseUserInfoBo) *bo.HomeIssueOwnerInfoBo {
	ownerInfo := &bo.HomeIssueOwnerInfoBo{}
	ownerInfo.ID = baseUserInfo.UserId
	ownerInfo.UserId = baseUserInfo.UserId
	ownerInfo.Name = baseUserInfo.Name
	ownerInfo.Avatar = &baseUserInfo.Avatar
	ownerInfo.IsDeleted = baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted
	ownerInfo.IsDisabled = baseUserInfo.OrgUserStatus == consts.AppStatusDisabled
	return ownerInfo
}

func ConvertIssueTagBosToMapGroupByIssueId(issueTagBos []bo.IssueTagBo) maps.LocalMap {
	issueTagBoMap := maps.LocalMap{}
	if issueTagBos != nil && len(issueTagBos) > 0 {
		for _, issueTagBo := range issueTagBos {
			issueId := issueTagBo.IssueId
			if issueTagsInterface, ok := issueTagBoMap[issueId]; ok {
				if issueTags, ok := issueTagsInterface.(*[]bo.IssueTagBo); ok {
					*issueTags = append(*issueTags, issueTagBo)
				}
			} else {
				issueTagBoMap[issueId] = &[]bo.IssueTagBo{issueTagBo}
			}
		}
	}
	return issueTagBoMap
}

func ConvertIssueTagBosToHomeIssueTagBos(issueTagBos []bo.IssueTagBo) []bo.HomeIssueTagInfoBo {
	homeIssueTagBos := make([]bo.HomeIssueTagInfoBo, 0)
	if issueTagBos != nil && len(issueTagBos) > 0 {
		for _, issueTagBo := range issueTagBos {
			homeIssueTagBos = append(homeIssueTagBos, bo.HomeIssueTagInfoBo{
				ID:        issueTagBo.TagId,
				Name:      issueTagBo.TagName,
				FontStyle: issueTagBo.FontStyle,
				BgStyle:   issueTagBo.BgStyle,
			})
		}
	}
	return homeIssueTagBos
}

func NoticeIssueAudit(orgId int64, userId int64, orgBaseInfo *bo.BaseOrgInfoBo, issueBo *bo.IssueBo, notNeedNoticeIds []int64) {
	if orgBaseInfo.OutOrgId == "" {
		return
	}

	var allUserIds []int64
	for _, auditorId := range issueBo.AuditorIdsI64 {
		if ok, _ := slice.Contain(notNeedNoticeIds, auditorId); ok {
			continue
		}
		allUserIds = append(allUserIds, auditorId)
	}

	allUserIds = slice.SliceUniqueInt64(allUserIds)
	if len(allUserIds) == 0 {
		return
	}

	// 查询用户信息
	resp := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: allUserIds,
	})
	if resp.Failure() {
		log.Errorf("[NoticeIssueAudit] GetBaseUserInfoBatch err: %v", resp.Error())
		return
	}
	if len(resp.BaseUserInfos) == 0 {
		return
	}
	userInfoMap := map[int64]*bo.BaseUserInfoBo{}
	for i, info := range resp.BaseUserInfos {
		userInfoMap[info.UserId] = &resp.BaseUserInfos[i]
	}

	// 机器人将信息回复给用户
	var allOutUserIds []string
	for _, id := range allUserIds {
		if userId == id {
			//自己操作的就不需要推送给自己了
			continue
		}
		if info, ok := userInfoMap[id]; ok {
			if info.OutUserId != "" {
				allOutUserIds = append(allOutUserIds, info.OutUserId)
			}
		}
	}
	if len(allOutUserIds) == 0 {
		return
	}
	//card.SendCardMeta(orgBaseInfo.SourceChannel, orgBaseInfo.OutOrgId, cd, allOutUserIds)
}

// GetTodoCard
func GetTodoCard(orgId, userId int64, sourceChannel string, todo *automationPb.Todo) (*pushPb.TemplateCard, errs.SystemErrorInfo) {
	return nil, nil
}

// GetCardInfoForStartAudit 获取并向 cardInfo 设置组装卡片时的信息
// starter 审批发起人的 id
func GetCardInfoForStartAudit(starter int64, issue *bo.IssueBo) (*projectvo.BaseInfoBoxForIssueCard, errs.SystemErrorInfo) {
	infoObj := &projectvo.BaseInfoBoxForIssueCard{}

	infoObj.IssueInfo = *issue
	infoObj.OperateUserId = starter
	// 查询任务负责人
	userIds := make([]int64, 0, len(infoObj.IssueInfo.OwnerIdI64)+1)
	userIds = append(userIds, starter)
	userIds = append(userIds, infoObj.IssueInfo.OwnerIdI64...)
	userInfoArr, err := orgfacade.GetBaseUserInfoBatchRelaxed(issue.OrgId, userIds)
	if err != nil {
		log.Errorf("[GetAndSetCardInfoForUrgeAudit] 查询组织 %d 用户信息出现异常 %v", issue.OrgId, err)
		return infoObj, err
	}
	userMap := make(map[int64]bo.BaseUserInfoBo, len(userInfoArr))
	for _, user := range userInfoArr {
		userMap[user.UserId] = user
	}
	if opUser, ok := userMap[starter]; ok {
		infoObj.OperateUser = opUser
	}
	//ownerUserNameArr := make([]string, 0)
	ownerInfos := make([]*bo.BaseUserInfoBo, 0)
	for _, uid := range infoObj.IssueInfo.OwnerIdI64 {
		if user, ok := userMap[uid]; ok {
			//ownerUserNameArr = append(ownerUserNameArr, user.Name)
			ownerInfos = append(ownerInfos, &user)
		}
	}
	infoObj.OwnerInfos = ownerInfos

	projectInfo, err1 := LoadProjectAuthBo(issue.OrgId, issue.ProjectId)
	if err1 != nil {
		log.Error(err1)
		return infoObj, err1
	}
	infoObj.ProjectAuthBo = *projectInfo

	tableColumns, err := GetTableColumnsMap(issue.OrgId, infoObj.IssueInfo.TableId, nil)
	if err != nil {
		log.Errorf("[GetAndSetCardInfoForUrgeAudit] 获取表头失败 org:%d proj:%d table:%d starter: %d, err: %v",
			issue.OrgId, issue.ProjectId, infoObj.IssueInfo.TableId, starter, err)
		return infoObj, err
	}
	infoObj.ProjectTableColumn = tableColumns
	headers := make(map[string]lc_table.LcCommonField, 0)
	copyer.Copy(tableColumns, &headers)
	infoObj.TableColumnMap = headers

	// 查询父任务信息
	if issue.ParentId > 0 {
		parents, err := GetIssueInfosLc(issue.OrgId, 0, []int64{issue.ParentId})
		if err != nil {
			log.Errorf("[GetCardInfoForStartAudit] GetIssueInfosLc err: %v, parentId: %d", err, issue.ParentId)
			return nil, err
		}
		if len(parents) > 0 {
			parent := parents[0]
			infoObj.ParentIssue = *parent
		}
	}

	// 查询任务所在表名
	if infoObj.IssueInfo.TableId != 0 {
		tableInfo, err := GetTableByTableId(issue.OrgId, starter, infoObj.IssueInfo.TableId)
		if err != nil {
			log.Errorf("[GetAndSetCardInfoForUrgeAudit] GetTableByTableId err:%v, userId:%d, tableId:%d", err, starter,
				infoObj.IssueInfo.TableId)
			return infoObj, err
		}
		infoObj.IssueTableInfo = *tableInfo
	} else {
		err := errs.InvalidTableId
		log.Errorf("[GetAndSetCardInfoForUrgeAudit] err: %v, issueId: %d", err, issue.Id)
		return infoObj, err
	}

	orgBaseResp, errSys := orgfacade.GetBaseOrgInfoRelaxed(issue.OrgId)
	if errSys != nil {
		log.Errorf("[GetCardInfoForStartAudit] err:%v, orgId:%v", errSys, issue.OrgId)
		return infoObj, errSys
	}
	infoObj.SourceChannel = orgBaseResp.SourceChannel

	// 查看详情、PC端查看等按钮的 url 链接
	links := GetIssueLinks(orgBaseResp.SourceChannel, issue.OrgId, issue.Id)
	infoObj.IssueInfoUrl = links.SideBarLink
	infoObj.IssuePcUrl = links.Link

	return infoObj, nil
}

////返回要切换的项目id， 如果是0，表示不需要切换
////第二个参数是项目对象类型的项目id
//func CheckProjectObjectTypeSwitchProject(projectId, targetProjectId, orgId int64) (int64, int64, errs.SystemErrorInfo) {
//	switchProjectId := int64(0)
//	if projectId != targetProjectId {
//		switchProjectId = targetProjectId
//	}
//	return switchProjectId, targetProjectId, nil
//}

func ChangeIssueChildren(operatorId int64, issueBo *bo.IssueBo, parentId int64, parentPath string) ([]int64, []*bo.IssueBo, errs.SystemErrorInfo) {
	if issueBo.ParentId == parentId {
		return nil, nil, nil
	}
	fromParentId := issueBo.ParentId
	targetParentId := parentId

	//判断父任务和当前任务是否有交叉（即：父任务是不是当前任务的下级）
	parentArr := strings.Split(parentPath, ",")
	if ok, _ := slice.Contain(parentArr, strconv.FormatInt(issueBo.Id, 10)); ok {
		return nil, nil, errs.ParentIsChildIssue
	}

	//childPos := &[]po.PpmPriIssue{}
	////查询是否存在子任务
	//err := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
	//	consts.TcPath:     db.Like(fmt.Sprintf("%s%d,%s", issueBo.Path, issueBo.Id, "%")),
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//	consts.TcOrgId:    issueBo.OrgId,
	//}, childPos)
	//if err != nil {
	//	log.Error(err)
	//	return nil, nil, errs.MysqlOperateError
	//}
	//
	////子任务集合
	//childrenIds := []int64{}
	//if len(*childPos) > 0 {
	//	for _, data := range *childPos {
	//		childrenIds = append(childrenIds, data.Id)
	//	}
	//}

	filterColumns := []string{
		lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
		lc_helper.ConvertToFilterColumn(consts.TcPath),
	}
	condition := &tablePb.Condition{Type: tablePb.ConditionType_and}
	condition.Conditions = GetNoRecycleCondition(
		GetRowsCondition(consts.TcPath, tablePb.ConditionType_like, fmt.Sprintf(",%d,", issueBo.Id), nil))
	childInfos, errSys := GetIssueInfosMapLc(issueBo.OrgId, operatorId, condition, filterColumns, -1, -1)
	if errSys != nil {
		log.Errorf("[ChangeIssueChildren] GetIssueInfosMapLc err:%v, orgId:%v, issueId:%v", errSys, issueBo.OrgId, issueBo.Id)
		return nil, nil, errSys
	}
	var issueBos []*bo.IssueBo
	childrenIds := make([]int64, 0, len(childInfos))
	for _, data := range childInfos {
		issue, errSys := ConvertIssueDataToIssueBo(data)
		if errSys != nil {
			log.Errorf("[ChangeIssueChildren]ConvertIssueDataToIssueBo err:%v, orgId:%v, issueId:%v", errSys, issueBo.OrgId, issueBo.Id)
			return nil, nil, errSys
		}
		childrenIds = append(childrenIds, issue.Id)
		issueBos = append(issueBos, issue)
	}

	newPath := "0,"
	if parentId != 0 {
		newPath = parentPath + strconv.FormatInt(parentId, 10) + ","
	}

	//判断层级是否超出9级
	if parentId != 0 {
		if strings.Count(newPath, ",") > consts.IssueLevel {
			return nil, nil, errs.IssueLevelOutLimit
		}

		//转化后相差的等级
		changeLevel := strings.Count(newPath, ",") - strings.Count(issueBo.Path, ",")
		for _, issue := range issueBos {
			if strings.Count(issue.Path, ",")+changeLevel > consts.IssueLevel {
				return nil, nil, errs.IssueLevelOutLimit
			}
		}
	}

	appId, appIdErr := GetAppIdFromProjectId(issueBo.OrgId, issueBo.ProjectId)
	if appIdErr != nil {
		log.Error(appIdErr)
		return nil, nil, appIdErr
	}

	// 变更之前先查询老数据信息
	allIssueIds := append(childrenIds, issueBo.Id)
	oldIssueBos, errSys := GetIssueInfosLc(issueBo.OrgId, operatorId, allIssueIds, GetMoveRelateColumnIds()...)
	if errSys != nil {
		log.Errorf("[ChangeIssueChildren] GetIssueInfosLc err: %v", errSys)
		return nil, nil, errSys
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 更新父任务
		upd := mysql.Upd{
			consts.TcParentId: parentId,
			consts.TcPath:     newPath,
			consts.TcUpdator:  operatorId,
		}
		//_, updateErr := mysql.UpdateSmartWithCond(consts.TableIssue, db.Cond{
		//	consts.TcId: issueBo.Id,
		//}, upd)
		//if updateErr != nil {
		//	log.Error(updateErr)
		//	return errs.MysqlOperateError
		//}

		// 更新子任务
		//_, updateErr = mysql.UpdateSmartWithCond(consts.TableIssue, db.Cond{
		//	consts.TcId: db.In(childrenIds),
		//}, mysql.Upd{
		//	consts.TcPath:    db.Raw("replace(`path`, \"" + issueBo.Path + "\", '" + newPath + "')"),
		//	consts.TcUpdator: operatorId,
		//})
		//if updateErr != nil {
		//	log.Error(updateErr)
		//	return errs.MysqlOperateError
		//}

		// 无码更新: 父任务的Path和ParentId
		updData := slice2.CaseCamelCopy(upd)
		updData[consts.BasicFieldIssueId] = issueBo.Id
		resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
			AppId:   appId,
			OrgId:   issueBo.OrgId,
			UserId:  operatorId,
			TableId: issueBo.TableId,
			Form: []map[string]interface{}{
				updData,
			},
		})
		if resp.Failure() {
			log.Error(resp.Error())
			return resp.Error()
		}

		// 无码更新：子任务的Path
		batchResp := formfacade.LessUpdateIssueBatchRaw(&formvo.LessUpdateIssueBatchReq{
			OrgId:  issueBo.OrgId,
			AppId:  appId,
			UserId: operatorId,
			Condition: vo.LessCondsData{
				Type: consts.ConditionAnd,
				Conds: []*vo.LessCondsData{
					{
						Type:   consts.ConditionEqual,
						Value:  issueBo.OrgId,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
					},
					{
						Type:   consts.ConditionIn,
						Values: childrenIds,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
					},
				},
			},
			Sets: []datacenter.Set{
				{
					Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
					Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, issueBo.Path, newPath),
					Type:            consts.SetTypeNormal,
					Action:          consts.SetActionSet,
					WithoutPretreat: true,
				},
			},
		})
		if batchResp.Failure() {
			log.Error(batchResp.Error())
			return batchResp.Error()
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return nil, nil, errs.MysqlOperateError
	}

	issueBo.ParentId = parentId
	issueBo.Path = newPath

	// 保存动态
	asyn.Execute(func() {
		var oldName, newName string
		if fromParentId > 0 {
			issues, errSys := GetIssueInfosLc(issueBo.OrgId, operatorId, []int64{fromParentId}, lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle))
			if errSys != nil {
				log.Error(errSys)
				return
			}
			if len(issues) > 0 {
				oldName = issues[0].Title
				if strings.Trim(oldName, " ") == "" {
					oldName = "未命名记录"
				}
			}
		}
		if targetParentId > 0 {
			issues, errSys := GetIssueInfosLc(issueBo.OrgId, operatorId, []int64{targetParentId}, lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle))
			if errSys != nil {
				log.Error(errSys)
				return
			}
			if len(issues) > 0 {
				newName = issues[0].Title
				if strings.Trim(newName, " ") == "" {
					newName = "未命名记录"
				}
			}
		}
		pushType := consts.PushTypeUpdateIssueParent
		var changeList []bo.TrendChangeListBo
		changeList = append(changeList, bo.TrendChangeListBo{
			Field:     consts.BasicFieldParentId,
			FieldName: consts.Parent,
			OldValue:  oldName,
			NewValue:  newName,
		})

		issueTrendsBo := &bo.IssueTrendsBo{
			PushType:      pushType,
			OrgId:         issueBo.OrgId,
			OperatorId:    operatorId,
			DataId:        issueBo.DataId,
			IssueId:       issueBo.Id,
			ParentIssueId: issueBo.ParentId,
			ProjectId:     issueBo.ProjectId,
			TableId:       issueBo.TableId,
			PriorityId:    issueBo.PriorityId,
			IssueTitle:    issueBo.Title,
			ParentId:      issueBo.ParentId,
			OldValue:      cast.ToString(fromParentId),
			NewValue:      cast.ToString(targetParentId),
			Ext: bo.TrendExtensionBo{
				ObjName:    issueBo.Title,
				ChangeList: changeList,
			},
		}
		PushIssueTrends(issueTrendsBo)
	})

	return allIssueIds, oldIssueBos, nil
}

// UpdateIssueProjectTable 更新任务所属项目/表
func UpdateIssueProjectTable(orgId, operatorId int64, issueBo *bo.IssueBo, targetTableId int64, targetProjectId int64, takeChildren bool, updatePath bool) errs.SystemErrorInfo {
	newUUID := uuid.NewUuid()
	lockKey := fmt.Sprintf("%s%d", consts.IssueRelateOperationLock, issueBo.Id)
	suc, lockErr := cache.TryGetDistributedLock(lockKey, newUUID)
	if lockErr != nil {
		log.Error(lockErr)
		return errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, newUUID); err != nil {
				log.Error(err)
			}
		}()
	} else {
		//未获取到锁，直接响应错误信息
		return errs.MoveIssueFail
	}

	fromProjectId := issueBo.ProjectId
	fromTableId := issueBo.TableId
	issueId := issueBo.Id
	var fromProjectBo *bo.ProjectBo
	var targetProjectBo *bo.ProjectBo
	var fromTable *projectvo.TableMetaData
	var targetTable *projectvo.TableMetaData
	var errSys errs.SystemErrorInfo
	newPath := "0,"

	isChangeProject := fromProjectId != targetProjectId
	isChangeTable := fromTableId != targetTableId
	if !isChangeProject && !isChangeTable {
		return nil
	}

	// 验证权限
	if isChangeProject && targetProjectId > 0 {
		errSys = AuthProject(orgId, operatorId, targetProjectId, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Create)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}
	}

	// from/target project
	if fromProjectId > 0 {
		fromProjectBo, errSys = GetProject(orgId, fromProjectId)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}
	}
	if targetProjectId == fromProjectId {
		targetProjectBo = fromProjectBo
	} else if targetProjectId != fromProjectId && targetProjectId > 0 {
		targetProjectBo, errSys = GetProject(orgId, targetProjectId)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}
	}

	// from/target table
	if fromTableId > 0 {
		fromTable, errSys = GetTableByTableId(orgId, operatorId, fromTableId)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}
	}
	if targetTableId == fromTableId {
		targetTable = fromTable
	} else if targetTableId != fromTableId && targetTableId > 0 {
		targetTable, errSys = GetTableByTableId(orgId, operatorId, targetTableId)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}
	}

	// target app
	var fromAppId int64
	var targetAppId int64
	var summaryAppId int64
	if fromProjectBo != nil {
		fromAppId = fromProjectBo.AppId
	}
	if targetProjectBo != nil {
		targetAppId = targetProjectBo.AppId
	}
	if targetAppId == 0 {
		summaryAppId, errSys = GetOrgSummaryAppId(orgId)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}
		targetAppId = summaryAppId
	}
	if fromAppId == 0 {
		summaryAppId, errSys = GetOrgSummaryAppId(orgId)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}
		fromAppId = summaryAppId
	}
	//if targetAppId == 0 {
	//	targetAppId = summaryAppId
	//}

	upd := mysql.Upd{
		consts.TcUpdator:   operatorId,
		consts.TcProjectId: targetProjectId,
		consts.TcTableId:   targetTableId,
	}

	// 任务状态
	allNewStatus, errSys := GetTableStatus(orgId, targetTableId)
	if errSys != nil {
		log.Error(errSys)
		return errSys
	}
	newStatusType := consts.StatusTypeNotStart
	var newStatus int64
	if isChangeTable {
		for _, status := range allNewStatus {
			if status.Type == newStatusType {
				newStatus = status.ID
				upd[consts.TcStatus] = status.ID
				break
			}
		}
	}

	// 确认状态
	if isChangeTable {
		if targetProjectBo != nil && targetProjectBo.ProjectTypeId == consts.ProjectTypeNormalId {
			upd[consts.TcAuditStatus] = consts.AuditStatusNotView
		} else {
			upd[consts.TcAuditStatus] = consts.AuditStatusNoNeed
		}
	}

	// 迭代
	if isChangeProject {
		upd[consts.TcIterationId] = 0
	}

	allIssueIds := []int64{issueId}

	// 查询子任务
	var childIssuePos []*po.PpmPriIssue
	if takeChildren {
		//err := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
		//	consts.TcPath:     db.Like(fmt.Sprintf("%s%d,%s", issueBo.Path, issueBo.Id, "%")),
		//	consts.TcIsDelete: consts.AppIsNoDelete,
		//	consts.TcOrgId:    orgId,
		//}, &childIssuePos)
		//if err != nil {
		//	log.Error(err)
		//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		//}
		condition := &tablePb.Condition{
			Type:       tablePb.ConditionType_and,
			Conditions: GetNoRecycleCondition(GetRowsCondition(consts.BasicFieldPath, tablePb.ConditionType_like, fmt.Sprintf(",%d,", issueId), nil)),
		}
		filterColumns := []string{
			lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
		}
		issueChildInfos, errSys := GetIssueInfosMapLc(orgId, operatorId, condition, filterColumns, -1, -1)
		if errSys != nil {
			log.Errorf("[UpdateIssueProjectTable] GetIssueInfosMapLc err:%v, orgId:%v, issueId:%v", errSys, orgId, issueId)
			return errSys
		}

		for _, data := range issueChildInfos {
			childIssuePos = append(childIssuePos, &po.PpmPriIssue{
				Id:       cast.ToInt64(data[consts.BasicFieldIssueId]),
				ParentId: cast.ToInt64(data[consts.BasicFieldParentId]),
			})
		}

		if len(childIssuePos) > 0 {
			for _, data := range childIssuePos {
				allIssueIds = append(allIssueIds, data.Id)
			}
		}
	}

	// 开始写入DB
	var updForm []map[string]interface{}  // 移动项目表的任务
	var updForm2 []map[string]interface{} // 不移动项目表的任务
	var updateBatchReqs []*formvo.LessUpdateIssueBatchReq
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 更新本任务
		parentUpd := slice2.CopyUpd(upd)
		if updatePath {
			parentUpd[consts.TcParentId] = 0
			parentUpd[consts.TcPath] = newPath
		}
		//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcId: issueBo.Id,
		//}, parentUpd)
		//if err != nil {
		//	log.Error(err)
		//	return err
		//}
		parentUpdForm := slice2.CaseCamelCopy(parentUpd)
		parentUpdForm[consts.BasicFieldIssueId] = issueBo.Id
		if isChangeTable {
			parentUpdForm[consts.BasicFieldIssueStatus] = newStatus
			parentUpdForm[consts.BasicFieldIssueStatusType] = consts.StatusTypeNotStart
		}
		parentUpdForm[consts.BasicFieldAppId] = cast.ToString(targetAppId)
		parentUpdForm[consts.BasicFieldTableId] = cast.ToString(targetTableId) // 转成字符串
		delete(parentUpdForm, consts.TcStatus)
		updForm = append(updForm, parentUpdForm)

		// 处理子任务
		if len(childIssuePos) > 0 {
			if takeChildren {
				// 移动子任务
				for _, issue := range childIssuePos {
					childUpd := slice2.CopyUpd(upd)
					//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
					//	consts.TcId: issue.Id,
					//}, childUpd)
					//if err != nil {
					//	log.Error(err)
					//	return err
					//}
					childUpdForm := slice2.CaseCamelCopy(childUpd)
					childUpdForm[consts.BasicFieldIssueId] = issue.Id
					if isChangeTable {
						childUpdForm[consts.BasicFieldIssueStatus] = newStatus
						childUpdForm[consts.BasicFieldIssueStatusType] = consts.StatusTypeNotStart
					}
					childUpdForm[consts.BasicFieldAppId] = cast.ToString(targetAppId)
					childUpdForm[consts.BasicFieldTableId] = cast.ToString(targetTableId) // 转成字符串
					delete(parentUpdForm, consts.TcStatus)
					delete(childUpdForm, consts.TcPath)
					updForm = append(updForm, childUpdForm)
				}

				// 更新所有子任务path
				if updatePath {
					likeCond := fmt.Sprintf("%s%d,%s", issueBo.Path, issueBo.Id, "%")
					//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
					//	consts.TcPath:  db.Like(likeCond),
					//	consts.TcOrgId: orgId,
					//}, mysql.Upd{
					//	consts.TcPath:    db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, issueBo.Path, newPath)),
					//	consts.TcUpdator: operatorId,
					//})
					//if err != nil {
					//	log.Error(err)
					//	return err
					//}
					// 无码的请求体
					updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
						OrgId:  issueBo.OrgId,
						AppId:  targetAppId,
						UserId: operatorId,
						Condition: vo.LessCondsData{
							Type: consts.ConditionAnd,
							Conds: []*vo.LessCondsData{
								{
									Type:   consts.ConditionEqual,
									Value:  issueBo.OrgId,
									Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
								},
								{
									Type:   consts.ConditionLike,
									Value:  likeCond,
									Column: lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
								},
							},
						},
						Sets: []datacenter.Set{
							{
								Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
								Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, issueBo.Path, newPath),
								Type:            consts.SetTypeNormal,
								Action:          consts.SetActionSet,
								WithoutPretreat: true,
							},
						},
					})
				}
			} else {
				// 不带子任务一起移动的话，就把顶级子任务变为父任务
				//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
				//	consts.TcPath:     db.Like(fmt.Sprintf("%s%d,%s", issueBo.Path, issueBo.Id, "%")),
				//	consts.TcOrgId:    orgId,
				//	consts.TcParentId: issueBo.Id,
				//}, mysql.Upd{
				//	consts.TcParentId: 0,
				//	consts.TcPath:     db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, issueBo.Path, newPath)),
				//	consts.TcUpdator:  operatorId,
				//})
				//if err != nil {
				//	log.Error(err)
				//	return err
				//}
				for _, issue := range childIssuePos {
					if issue.ParentId == issueId {
						updForm2 = append(updForm2, map[string]interface{}{
							consts.BasicFieldIssueId:  issue.Id,
							consts.BasicFieldParentId: 0,
						})
					}
				}

				// 子任务的子任务path
				likeCond := fmt.Sprintf("%s%d,%s", issueBo.Path, issueBo.Id, "%")
				//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
				//	consts.TcPath:     likeCond,
				//	consts.TcOrgId:    orgId,
				//	consts.TcParentId: db.NotEq(issueBo.Id),
				//}, mysql.Upd{
				//	consts.TcPath:    db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, issueBo.Path, newPath)),
				//	consts.TcUpdator: operatorId,
				//})
				//if err != nil {
				//	log.Error(err)
				//	return err
				//}
				// 无码的请求体
				updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
					OrgId:  issueBo.OrgId,
					AppId:  fromAppId,
					UserId: operatorId,
					Condition: vo.LessCondsData{
						Type: consts.ConditionAnd,
						Conds: []*vo.LessCondsData{
							{
								Type:   consts.ConditionEqual,
								Value:  issueBo.OrgId,
								Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
							},
							{
								Type:   consts.ConditionUnEqual,
								Value:  issueBo.Id,
								Column: lc_helper.ConvertToCondColumn(consts.BasicFieldParentId),
							},
							{
								Type:   consts.ConditionLike,
								Value:  likeCond,
								Column: lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
							},
						},
					},
					Sets: []datacenter.Set{
						{
							Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
							Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, issueBo.Path, newPath),
							Type:            consts.SetTypeNormal,
							Action:          consts.SetActionSet,
							WithoutPretreat: true,
						},
					},
				})
			}
		}

		// 无码更新
		if len(updForm) > 0 {
			resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
				AppId:   targetAppId,
				OrgId:   orgId,
				UserId:  operatorId,
				TableId: targetTableId,
				Form:    updForm,
			})
			if resp.Failure() {
				log.Error(resp.Error())
				return resp.Error()
			}
		}
		if len(updForm2) > 0 {
			resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
				AppId:   fromAppId,
				OrgId:   orgId,
				UserId:  operatorId,
				TableId: fromTableId,
				Form:    updForm2,
			})
			if resp.Failure() {
				log.Error(resp.Error())
				return resp.Error()
			}
		}

		// 无码更新
		for _, req := range updateBatchReqs {
			batchResp := formfacade.LessUpdateIssueBatchRaw(req)
			if batchResp.Failure() {
				log.Error(batchResp.Error())
				return batchResp.Error()
			}
		}

		return nil
	})
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}

	// 更新任务工时关联的 projectId
	if isChangeProject {
		errSys = UpdateWorkHourProjectId(orgId, fromProjectId, targetProjectId, allIssueIds, false)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}
	}

	// 移动任务附件
	if isChangeProject {
		errSys = switchIssueResource(orgId, allIssueIds, operatorId, targetProjectId)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}
	}

	if isChangeProject && targetProjectId > 0 {
		// 处理任务成员和关注人
		errSys = switchIssueMember(orgId, allIssueIds, operatorId, targetProjectId)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}

	}

	//asyn.Execute(func() {
	//	if isChangeTable {
	//		if err := IssueChatInviteUsersWhenMoveIssue(orgId, []int64{issueId}, targetTableId); err != nil {
	//			log.Errorf("[UpdateIssueProjectTable] IssueChatInviteUsersWhenMoveIssue err: %v, issueId: %d", err, issueId)
	//			return
	//		}
	//	}
	//})

	asyn.Execute(func() {
		if isChangeProject || isChangeTable {
			pushType := consts.PushTypeUpdateIssueProjectTable
			var changeList []bo.TrendChangeListBo
			if isChangeProject {
				var oldName, newName string
				if fromProjectBo != nil {
					oldName = fromProjectBo.Name
				}
				if targetProjectBo != nil {
					newName = targetProjectBo.Name
				}
				changeList = append(changeList, bo.TrendChangeListBo{
					Field:     consts.BasicFieldProjectId,
					FieldName: consts.Project,
					OldValue:  oldName,
					NewValue:  newName,
				})
			}
			if isChangeTable {
				var oldName, newName string
				if fromTable != nil {
					oldName = fromTable.Name
				}
				if targetTable != nil {
					newName = targetTable.Name
				}
				changeList = append(changeList, bo.TrendChangeListBo{
					Field:     consts.BasicFieldTableId,
					FieldName: consts.Table,
					OldValue:  oldName,
					NewValue:  newName,
				})
			}

			oldValue := bo.ProjectTableBo{}
			oldValue.ProjectId = fromProjectId
			oldValue.TableId = fromTableId

			newValue := bo.ProjectTableBo{}
			newValue.ProjectId = targetProjectId
			newValue.TableId = targetTableId

			issueTrendsBo := &bo.IssueTrendsBo{
				PushType:      pushType,
				OrgId:         issueBo.OrgId,
				OperatorId:    operatorId,
				DataId:        issueBo.DataId,
				IssueId:       issueBo.Id,
				ParentIssueId: issueBo.ParentId,
				ProjectId:     targetProjectId,
				TableId:       targetTableId,
				PriorityId:    issueBo.PriorityId,
				IssueTitle:    issueBo.Title,
				ParentId:      issueBo.ParentId,
				OldValue:      json.ToJsonIgnoreError(oldValue),
				NewValue:      json.ToJsonIgnoreError(newValue),
				Ext: bo.TrendExtensionBo{
					ObjName:    issueBo.Title,
					ChangeList: changeList,
				},
			}
			PushIssueTrends(issueTrendsBo)

			//asyn.Execute(func() {
			//	PushIssueThirdPlatformNotice(issueTrendsBo)
			//})
			//asyn.Execute(func() {
			//	//推送群聊卡片
			//	PushInfoToChat(issueBo.OrgId, targetProjectId, issueTrendsBo)
			//})
		}
	})
	return nil
}

//// MoveIssueProTable 移动记录 选择移动过去的内容
//func MoveIssueProTable(orgId, operatorId int64, issueBo *bo.IssueBo, targetTableId int64, targetProjectId int64,
//	issueTitle string, chooseFields []string, childrenIds []int64) errs.SystemErrorInfo {
//	newUUID := uuid.NewUuid()
//	lockKey := fmt.Sprintf("%s%d", consts.IssueRelateOperationLock, issueBo.Id)
//	suc, lockErr := cache.TryGetDistributedLock(lockKey, newUUID)
//	if lockErr != nil {
//		log.Error(lockErr)
//		return errs.TryDistributedLockError
//	}
//	if suc {
//		defer func() {
//			if _, err := cache.ReleaseDistributedLock(lockKey, newUUID); err != nil {
//				log.Error(err)
//			}
//		}()
//	} else {
//		//未获取到锁，直接响应错误信息
//		return errs.MoveIssueFail
//	}
//
//	fromProjectId := issueBo.ProjectId
//	fromTableId := issueBo.TableId
//	issueId := issueBo.Id
//	var fromProjectBo *bo.ProjectBo
//	var targetProjectBo *bo.ProjectBo
//	var fromTable *projectvo.TableMetaData
//	var targetTable *projectvo.TableMetaData
//	var errSys errs.SystemErrorInfo
//	newPath := "0,"
//
//	isChangeProject := fromProjectId != targetProjectId
//	isChangeTable := fromTableId != targetTableId
//	if !isChangeProject && !isChangeTable {
//		return nil
//	}
//
//	// from/target project
//	if fromProjectId > 0 {
//		fromProjectBo, errSys = GetProjectSimple(orgId, fromProjectId)
//		if errSys != nil {
//			log.Error(errSys)
//			return errSys
//		}
//	}
//	if targetProjectId == fromProjectId {
//		targetProjectBo = fromProjectBo
//	} else if targetProjectId != fromProjectId && targetProjectId > 0 {
//		targetProjectBo, errSys = GetProjectSimple(orgId, targetProjectId)
//		if errSys != nil {
//			log.Error(errSys)
//			return errSys
//		}
//	}
//
//	// from/target table
//	if fromTableId > 0 {
//		fromTable, errSys = GetTableByTableId(orgId, operatorId, fromTableId)
//		if errSys != nil {
//			log.Error(errSys)
//			return errSys
//		}
//	}
//	if targetTableId == fromTableId {
//		targetTable = fromTable
//	} else if targetTableId != fromTableId && targetTableId > 0 {
//		targetTable, errSys = GetTableByTableId(orgId, operatorId, targetTableId)
//		if errSys != nil {
//			log.Error(errSys)
//			return errSys
//		}
//	}
//
//	// target app
//	var fromAppId int64
//	var targetAppId int64
//	var summaryAppId int64
//	if fromProjectBo != nil {
//		fromAppId = fromProjectBo.AppId
//	}
//	if targetProjectBo != nil {
//		targetAppId = targetProjectBo.AppId
//	}
//	if targetAppId == 0 {
//		summaryAppId, errSys = GetOrgSummaryAppId(orgId)
//		if errSys != nil {
//			log.Error(errSys)
//			return errSys
//		}
//		targetAppId = summaryAppId
//	}
//	if fromAppId == 0 {
//		summaryAppId, errSys = GetOrgSummaryAppId(orgId)
//		if errSys != nil {
//			log.Error(errSys)
//			return errSys
//		}
//		fromAppId = summaryAppId
//	}
//	//if targetAppId == 0 {
//	//	targetAppId = summaryAppId
//	//}
//
//	// 本条被移动的记录的更新
//	upd := mysql.Upd{
//		consts.TcUpdator:   operatorId,
//		consts.TcProjectId: targetProjectId,
//		consts.TcTableId:   targetTableId,
//		consts.TcTitle:     issueTitle,
//		consts.TcPath:      newPath,
//		consts.TcParentId:  0,
//	}
//
//	// 任务状态
//	allNewStatus, errSys := GetTableStatus(orgId, targetTableId)
//	if errSys != nil {
//		log.Error(errSys)
//		return errSys
//	}
//	newStatusType := consts.StatusTypeNotStart
//	var newStatus int64
//	if isChangeTable {
//		for _, statusInfoBo := range allNewStatus {
//			if statusInfoBo.Type == newStatusType {
//				newStatus = statusInfoBo.ID
//				upd[consts.TcStatus] = statusInfoBo.ID
//				break
//			}
//		}
//	}
//
//	// 确认状态
//	if isChangeTable {
//		if targetProjectBo != nil && targetProjectBo.ProjectTypeId == consts.ProjectTypeNormalId {
//			upd[consts.TcAuditStatus] = consts.AuditStatusNotView
//		} else {
//			upd[consts.TcAuditStatus] = consts.AuditStatusNoNeed
//		}
//	}
//
//	// 迭代
//	if issueBo.IterationId != 0 {
//		// 如果关联了迭代
//		if !str.CheckStrInArray(chooseFields, consts.BasicFieldIterationId) {
//			upd[consts.TcIterationId] = 0
//		}
//	}
//
//	// 开始时间
//	if !str.CheckStrInArray(chooseFields, consts.BasicFieldPlanStartTime) {
//		upd[consts.TcPlanStartTime] = consts.BlankTime
//	}
//
//	// 截止时间
//	if !str.CheckStrInArray(chooseFields, consts.BasicFieldPlanStartTime) {
//		upd[consts.TcPlanEndTime] = consts.BlankTime
//	}
//
//	// 所有子记录id
//	allChildIds := []int64{}
//	// 真正需要携带过去的子任务id
//	allRealNeedMoveChildIds := []int64{}
//	allRealNeedMoveChildIds = append(allRealNeedMoveChildIds, childrenIds...)
//
//	childIssuePos := []po.PpmPriIssue{}
//	childIssueMap := map[int64]po.PpmPriIssue{}
//
//	// 查询该条记录下 所有子任务
//	err := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
//		consts.TcIsDelete: consts.AppIsNoDelete,
//		consts.TcOrgId:    orgId,
//		//consts.TcParentId: issueId,
//		consts.TcPath: db.Like(fmt.Sprintf("%s%d,%s", issueBo.Path, issueId, "%")),
//	}, &childIssuePos)
//	if err != nil {
//		log.Errorf("[MoveIssueProTable] 查询子任务err:%v", err)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	for _, childIssue := range childIssuePos {
//		allChildIds = append(allChildIds, childIssue.Id)
//		childIssueMap[childIssue.Id] = childIssue
//	}
//
//	// 传过来的子记录
//	chooseChildIssueBos := []po.PpmPriIssue{}
//	for _, child := range childrenIds {
//		if v, ok := childIssueMap[child]; ok {
//			chooseChildIssueBos = append(chooseChildIssueBos, v)
//		}
//	}
//	// 传过来的子记录的子记录
//	leftChildIds := businees.DifferenceInt64Set(allChildIds, childrenIds)
//	leftChildIssues := []po.PpmPriIssue{}
//	for _, id := range leftChildIds {
//		if v, ok := childIssueMap[id]; ok {
//			leftChildIssues = append(leftChildIssues, v)
//		}
//	}
//
//	needTakeLeftChildIds := []int64{}
//	// 剩下的子任务中的path如果包含childrenIds的path，说明是需要携带过去的id
//	for _, id := range childrenIds {
//		//path := fmt.Sprintf("0,%d,", id)
//		for _, leftChildIssue := range leftChildIssues {
//			if strings.Contains(leftChildIssue.Path, fmt.Sprintf("%d", id)) {
//				needTakeLeftChildIds = append(needTakeLeftChildIds, leftChildIssue.Id)
//			}
//		}
//	}
//
//	needTakeLeftChildIds = slice.SliceUniqueInt64(needTakeLeftChildIds)
//	allRealNeedMoveChildIds = append(allRealNeedMoveChildIds, needTakeLeftChildIds...)
//	allRealNeedMoveChildIds = slice.SliceUniqueInt64(allRealNeedMoveChildIds)
//	// 不需要移动的子记录
//	allNoNeedMoveChildIds := businees.DifferenceInt64Set(allChildIds, allRealNeedMoveChildIds)
//
//	var updateMoveForm []map[string]interface{}   // 移动的任务
//	var updateNoMoveForm []map[string]interface{} // 不需要移动的任务
//	// 批量更新无码数据，更新path
//	var updateBatchForm []formvo.LessUpdateIssueBatchReq
//
//	// 所有需要移动的记录 包括父记录和子记录
//	allMoveIssueIds := []int64{}
//	allMoveIssueIds = append(allMoveIssueIds, issueId)
//	allMoveIssueIds = append(allMoveIssueIds, allRealNeedMoveChildIds...)
//
//	err = mysql.TransX(func(tx sqlbuilder.Tx) error {
//		// 更新本条被移动的记录
//		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{consts.TcId: issueId}, upd)
//		if err != nil {
//			log.Error(err)
//			return err
//		}
//		// 构建本条记录的无码更新条件
//		parentForm := slice2.CaseCamelCopy(upd)
//		parentForm[consts.BasicFieldIssueId] = issueBo.Id
//		if isChangeTable {
//			parentForm[consts.BasicFieldIssueStatus] = newStatus
//			parentForm[consts.BasicFieldIssueStatusType] = consts.StatusTypeNotStart
//		}
//		parentForm[consts.BasicFieldAppIds] = []string{cast.ToString(targetAppId)}
//		parentForm[consts.BasicFieldTableId] = cast.ToString(targetTableId) // 转成字符串
//		delete(parentForm, consts.TcStatus)
//		updateMoveForm = append(updateMoveForm, parentForm)
//
//		// 更新子任务的数据
//
//		// 更新该条记录的子记录部分数据
//		// 区分 需要更新移动的子记录和不需要移动的子记录
//
//		// 构建子记录的更新数据
//		childUpd := slice2.CopyUpd(upd)
//		// 子记录的title不需要使用传过来的标题
//		delete(childUpd, consts.TcTitle)
//		delete(childUpd, consts.TcPath)
//		delete(childUpd, consts.TcParentId)
//
//		if len(allRealNeedMoveChildIds) > 0 {
//			_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
//				consts.TcOrgId:    orgId,
//				consts.TcIsDelete: consts.AppIsNoDelete,
//				consts.TcId:       db.In(childrenIds),
//			}, childUpd)
//			if err != nil {
//				log.Error(err)
//				return err
//			}
//			likeCond := fmt.Sprintf("%s%d,%s", issueBo.Path, issueBo.Id, "%")
//			_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
//				consts.TcPath:  db.Like(likeCond),
//				consts.TcOrgId: orgId,
//			}, mysql.Upd{
//				consts.TcPath:    db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, issueBo.Path, newPath)),
//				consts.TcUpdator: operatorId,
//			})
//			if err != nil {
//				log.Error(err)
//				return err
//			}
//
//			updateBatchForm = append(updateBatchForm, formvo.LessUpdateIssueBatchReq{
//				OrgId:  orgId,
//				AppId:  targetAppId,
//				UserId: operatorId,
//				Condition: vo.LessCondsData{
//					Type:   "like",
//					Value:  likeCond,
//					Column: lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
//				},
//				Sets: []datacenter.Set{
//					datacenter.Set{
//						Column:          consts.LcJsonColumn,
//						Value:           fmt.Sprintf("%s || jsonb_build_object('%s', replace(data->>'%s','%s','%s'))", consts.LcJsonColumn, consts.BasicFieldPath, consts.BasicFieldPath, issueBo.Path, newPath),
//						Type:            consts.SetTypeJson,
//						Action:          consts.SetActionSet,
//						WithoutPretreat: true,
//					},
//				},
//			})
//		}
//
//		if len(allNoNeedMoveChildIds) > 0 {
//			// 没有移动的子任务，将顶级子任务变为父任务
//			_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
//				consts.TcPath:     db.Like(fmt.Sprintf("%s%d,%s", issueBo.Path, issueBo.Id, "%")),
//				consts.TcOrgId:    orgId,
//				consts.TcId:       db.In(allNoNeedMoveChildIds),
//				consts.TcParentId: issueBo.Id,
//			}, mysql.Upd{
//				consts.TcParentId: 0,
//				consts.TcPath:     db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, issueBo.Path, newPath)),
//				consts.TcUpdator:  operatorId,
//			})
//			if err != nil {
//				log.Error(err)
//				return err
//			}
//
//			// 处理子任务的子任务path
//			//likeCond := fmt.Sprintf("%s%d,%s", issueBo.Path, issueBo.Id, "%")
//			_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
//				//consts.TcPath:     likeCond,
//				consts.TcOrgId:    orgId,
//				consts.TcId:       db.In(allNoNeedMoveChildIds),
//				consts.TcParentId: db.NotEq(issueBo.Id),
//			}, mysql.Upd{
//				consts.TcPath:    db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, issueBo.Path, newPath)),
//				consts.TcUpdator: operatorId,
//			})
//			if err != nil {
//				log.Error(err)
//				return err
//			}
//
//			// 更新path的无码条件组装
//			values := []interface{}{}
//			for _, id := range allNoNeedMoveChildIds {
//				values = append(values, id)
//			}
//			updateBatchForm = append(updateBatchForm, formvo.LessUpdateIssueBatchReq{
//				OrgId:  orgId,
//				AppId:  fromAppId,
//				UserId: operatorId,
//				Condition: vo.LessCondsData{
//					Type: "and",
//					Conds: []*vo.LessCondsData{
//						&vo.LessCondsData{
//							Type:   "un_equal",
//							Value:  issueBo.Id,
//							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldParentId),
//						},
//						&vo.LessCondsData{
//							Type:   "in",
//							Values: values,
//							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
//						},
//					},
//				},
//				Sets: []datacenter.Set{
//					datacenter.Set{
//						Column:          consts.LcJsonColumn,
//						Value:           fmt.Sprintf("%s || jsonb_build_object('%s', replace(data->>'%s','%s','%s'))", consts.LcJsonColumn, consts.BasicFieldPath, consts.BasicFieldPath, issueBo.Path, newPath),
//						Type:            consts.SetTypeJson,
//						Action:          consts.SetActionSet,
//						WithoutPretreat: true,
//					},
//				},
//			})
//		}
//
//		// 更新remark
//		if !str.CheckStrInArray(chooseFields, consts.BasicFieldRemark) {
//			parentForm[consts.BasicFieldRemark] = consts.BlankString
//			_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssueDetail, db.Cond{
//				consts.TcOrgId:    orgId,
//				consts.TcIsDelete: consts.AppIsNoDelete,
//				consts.TcIssueId:  db.In(allMoveIssueIds),
//			}, mysql.Upd{
//				consts.TcUpdator:   operatorId,
//				consts.TcProjectId: targetProjectId,
//				consts.TcRemark:    consts.BlankString,
//			})
//			if err != nil {
//				log.Error(err)
//				return err
//			}
//		}
//
//		// 更新issueRelation表
//		if !str.CheckStrInArray(chooseFields, consts.BasicFieldOwnerId) {
//			parentForm[consts.BasicFieldOwnerId] = []string{}
//			_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
//				consts.TcOrgId:        orgId,
//				consts.TcIssueId:      db.In(allMoveIssueIds),
//				consts.TcRelationType: consts.IssueRelationTypeOwner,
//				consts.TcIsDelete:     consts.AppIsNoDelete,
//			}, mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
//			if err != nil {
//				log.Error(err)
//				return err
//			}
//		}
//		if !str.CheckStrInArray(chooseFields, consts.BasicFieldFollowerIds) {
//			parentForm[consts.BasicFieldFollowerIds] = []string{}
//			_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
//				consts.TcOrgId:        orgId,
//				consts.TcIssueId:      db.In(allMoveIssueIds),
//				consts.TcRelationType: consts.IssueRelationTypeFollower,
//				consts.TcIsDelete:     consts.AppIsNoDelete,
//			}, mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
//			if err != nil {
//				log.Error(err)
//				return err
//			}
//		}
//		if !str.CheckStrInArray(chooseFields, consts.BasicFieldAuditorIds) {
//			parentForm[consts.BasicFieldAuditorIds] = []string{}
//			parentForm[consts.BasicFieldAuditStatusDetail] = map[string]int{}
//			//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
//			//	consts.TcOrgId:        orgId,
//			//	consts.TcIssueId:      db.In(allMoveIssueIds),
//			//	consts.TcRelationType: consts.IssueRelationTypeAuditor,
//			//	consts.TcIsDelete:     consts.AppIsNoDelete,
//			//}, mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
//			//if err != nil {
//			//	log.Error(err)
//			//	return err
//			//}
//		}
//
//		// 优先级
//		if !str.CheckStrInArray(chooseFields, consts.BasicFieldPriority) {
//			if _, ok := issueBo.LessData[consts.BasicFieldPriority]; ok {
//				parentForm[consts.BasicFieldPriority] = nil
//			}
//		}
//
//		// 工时
//		if !str.CheckStrInArray(chooseFields, consts.BasicFieldWorkHour) {
//			if _, ok := issueBo.LessData[consts.BasicFieldWorkHour]; ok {
//				parentForm[consts.BasicFieldWorkHour] = consts.DefaultWorkHour
//			}
//		}
//
//		appIdUpdateForm := map[int64][]map[string]interface{}{}
//
//		// 关联
//		if !str.CheckStrInArray(chooseFields, consts.BasicFieldRelating) {
//			relatingForm, err := getNoMoveRelatingForm(orgId, operatorId, issueBo, &issueBo.RelatingIssue, consts.BasicFieldRelating)
//			if err != nil {
//				log.Errorf("[MoveIssueProTable] err:%v", err)
//				return err
//			}
//			for appId, form := range relatingForm {
//				if _, ok := appIdUpdateForm[appId]; ok {
//					appIdUpdateForm[appId] = append(appIdUpdateForm[appId], form...)
//				} else {
//					appIdUpdateForm[appId] = form
//				}
//			}
//		}
//
//		// 前后置
//		if !str.CheckStrInArray(chooseFields, consts.BasicFieldBaRelating) {
//			baRelatingForm, err := getNoMoveRelatingForm(orgId, operatorId, issueBo, &issueBo.BaRelatingIssue, consts.BasicFieldBaRelating)
//			if err != nil {
//				log.Errorf("[MoveIssueProTable] err:%v", err)
//				return err
//			}
//			for appId, form := range baRelatingForm {
//				if _, ok := appIdUpdateForm[appId]; ok {
//					appIdUpdateForm[appId] = append(appIdUpdateForm[appId], form...)
//				} else {
//					appIdUpdateForm[appId] = form
//				}
//			}
//		}
//
//		for _, id := range allNoNeedMoveChildIds {
//			if childIssue, ok := childIssueMap[id]; ok {
//				if childIssue.ParentId == issueId {
//					updateNoMoveForm = append(updateNoMoveForm, map[string]interface{}{
//						consts.BasicFieldIssueId:  childIssue.Id,
//						consts.BasicFieldParentId: 0,
//					})
//				}
//			}
//		}
//
//		for _, id := range allRealNeedMoveChildIds {
//			childForm := map[string]interface{}{}
//			copyer.Copy(parentForm, &childForm)
//			childForm[consts.BasicFieldIssueId] = id
//			delete(childForm, consts.BasicFieldTitle)
//			delete(childForm, consts.BasicFieldParentId)
//			delete(childForm, consts.BasicFieldPath)
//			updateMoveForm = append(updateMoveForm, childForm)
//		}
//
//		if len(updateMoveForm) > 0 {
//			resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
//				AppId:   targetAppId,
//				OrgId:   orgId,
//				UserId:  operatorId,
//				TableId: targetTableId,
//				Form:    updateMoveForm,
//			})
//			if resp.Failure() {
//				log.Errorf("[MoveIssueProTable] err:%v, updateMoveForm:%v", resp.Error(), json.ToJsonIgnoreError(updateMoveForm))
//				return resp.Error()
//			}
//		}
//
//		if len(updateNoMoveForm) > 0 {
//			resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
//				AppId:   fromAppId,
//				OrgId:   orgId,
//				UserId:  operatorId,
//				TableId: fromTableId,
//				Form:    updateNoMoveForm,
//			})
//			if resp.Failure() {
//				log.Errorf("[MoveIssueProTable] err:%v, updateNoMoveForm:%v", resp.Error(), json.ToJsonIgnoreError(updateNoMoveForm))
//				return resp.Error()
//			}
//		}
//
//		for appId, form := range appIdUpdateForm {
//			resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
//				AppId:  appId,
//				OrgId:  orgId,
//				UserId: operatorId,
//				Form:   form,
//			})
//			if resp.Failure() {
//				log.Errorf("[MoveIssueProTable] err:%v, updateNoMoveForm:%v", resp.Error(), json.ToJsonIgnoreError(updateNoMoveForm))
//				return resp.Error()
//			}
//		}
//
//		for _, req := range updateBatchForm {
//			batchResp := formfacade.LessUpdateIssueBatchRaw(req)
//			if batchResp.Failure() {
//				log.Error(batchResp.Error())
//				return batchResp.Error()
//			}
//		}
//
//		return nil
//	})
//	if err != nil {
//		log.Error(err)
//		return errs.MysqlOperateError
//	}
//
//	// 更新任务工时关联的 projectId
//	if isChangeProject {
//		isDeleteWorkHour := false
//		if !str.CheckStrInArray(chooseFields, consts.BasicFieldWorkHour) {
//			// 不移动工时，就需要把工时关联删除
//			isDeleteWorkHour = true
//		}
//		errSys = UpdateWorkHourProjectId(orgId, fromProjectId, targetProjectId, allMoveIssueIds, isDeleteWorkHour)
//		if errSys != nil {
//			log.Error(errSys)
//			return errSys
//		}
//	}
//
//	// 移动任务附件
//	if isChangeProject {
//		errSys = switchIssueResource(orgId, allMoveIssueIds, operatorId, targetProjectId)
//		if errSys != nil {
//			log.Error(errSys)
//			return errSys
//		}
//	}
//
//	if isChangeProject && targetProjectId > 0 {
//		// 处理任务成员和关注人
//		errSys = switchIssueMember(orgId, allMoveIssueIds, operatorId, targetProjectId)
//		if errSys != nil {
//			log.Error(errSys)
//			return errSys
//		}
//
//		// 更新飞书日历
//		go func() {
//			defer func() {
//				if r := recover(); r != nil {
//					log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
//				}
//			}()
//			calendarErr := SwitchCalendar(orgId, fromProjectId, allMoveIssueIds, operatorId, targetProjectId)
//			if calendarErr != nil {
//				log.Error(calendarErr)
//				return
//			}
//		}()
//	}
//
//	asyn.Execute(func() {
//		if isChangeTable {
//			if err := IssueChatInviteUsersWhenMoveIssue(orgId, []int64{issueId}, targetTableId); err != nil {
//				log.Errorf("[MoveIssueProTable] IssueChatInviteUsersWhenMoveIssue err: %v, issueId: %d", err, issueId)
//				return
//			}
//		}
//	})
//
//	asyn.Execute(func() {
//		if isChangeProject || isChangeTable {
//			pushType := consts.PushTypeUpdateIssueProjectTable
//			var changeList []bo.TrendChangeListBo
//			if isChangeProject {
//				var oldName, newName string
//				if fromProjectBo != nil {
//					oldName = fromProjectBo.Name
//				}
//				if targetProjectBo != nil {
//					newName = targetProjectBo.Name
//				}
//				changeList = append(changeList, bo.TrendChangeListBo{
//					Field:     consts.BasicFieldProjectId,
//					FieldName: consts.Project,
//					OldValue:  oldName,
//					NewValue:  newName,
//				})
//			}
//			if isChangeTable {
//				var oldName, newName string
//				if fromTable != nil {
//					oldName = fromTable.Name
//				}
//				if targetTable != nil {
//					newName = targetTable.Name
//				}
//				changeList = append(changeList, bo.TrendChangeListBo{
//					Field:     consts.BasicFieldTableId,
//					FieldName: consts.Table,
//					OldValue:  oldName,
//					NewValue:  newName,
//				})
//			}
//
//			oldValue := bo.ProjectTableBo{}
//			oldValue.ProjectId = fromProjectId
//			oldValue.TableId = fromTableId
//
//			newValue := bo.ProjectTableBo{}
//			newValue.ProjectId = targetProjectId
//			newValue.TableId = targetTableId
//
//			issueTrendsBo := &bo.IssueTrendsBo{
//				PushType:      pushType,
//				OrgId:         issueBo.OrgId,
//				OperatorId:    operatorId,
//				IssueId:       issueBo.Id,
//				ParentIssueId: issueBo.ParentId,
//				ProjectId:     targetProjectId,
//				TableId:       targetTableId,
//				PriorityId:    issueBo.PriorityId,
//				IssueTitle:    issueBo.Title,
//				ParentId:      issueBo.ParentId,
//				OldValue:      json.ToJsonIgnoreError(oldValue),
//				NewValue:      json.ToJsonIgnoreError(newValue),
//				Ext: bo.TrendExtensionBo{
//					ObjName:    issueBo.Title,
//					ChangeList: changeList,
//				},
//			}
//			PushIssueTrends(issueTrendsBo)
//
//			//asyn.Execute(func() {
//			//	PushIssueThirdPlatformNotice(issueTrendsBo)
//			//})
//			//asyn.Execute(func() {
//			//	//推送群聊卡片
//			//	PushInfoToChat(issueBo.OrgId, targetProjectId, issueTrendsBo)
//			//})
//		}
//	})
//	return nil
//}

// 批量更新任务所属项目/表
func MoveIssueProTableBatch(orgId, operatorId int64, issueAndChildrenBos []*bo.IssueBo, inputIssueIds []int64,
	fromTableId int64, fromProjectId int64, targetTableId int64, targetProjectId int64, chooseFields []string, issueTitleMap map[int64]string) ([]int64, errs.SystemErrorInfo) {
	allMoveIssues := make([]*bo.IssueBo, 0) //除去鉴权失败的剩余所有符合条件的有效移动任务
	allMoveIds := make([]int64, 0)          //除去鉴权失败的剩余所有符合条件的有效移动任务id
	allRemainIds := make([]int64, 0)        //除去鉴权失败的剩余所有符合条件的有效移动任务id
	for _, childrenBo := range issueAndChildrenBos {
		if ok, _ := slice.Contain(inputIssueIds, childrenBo.Id); ok {
			allMoveIssues = append(allMoveIssues, childrenBo)
			allMoveIds = append(allMoveIds, childrenBo.Id)
		} else {
			allRemainIds = append(allRemainIds, childrenBo.Id)
		}
	}
	if len(allMoveIds) == 0 {
		return allMoveIds, nil
	}

	var fromProjectBo *bo.ProjectBo
	var targetProjectBo *bo.ProjectBo
	var fromTable *projectvo.TableMetaData
	var targetTable *projectvo.TableMetaData
	var relatingColumnIds []string
	var errSys errs.SystemErrorInfo
	rootPath := "0,"

	isChangeProject := fromProjectId != targetProjectId
	isChangeTable := fromTableId != targetTableId
	if !isChangeProject && !isChangeTable {
		return []int64{}, nil
	}

	// 验证权限
	if isChangeProject && targetProjectId > 0 {
		//校验当前用户有没有该项目的创建权限
		errSys = AuthProject(orgId, operatorId, targetProjectId, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Create)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
	}

	// from/target project
	if fromProjectId > 0 {
		fromProjectBo, errSys = GetProjectSimple(orgId, fromProjectId)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
	}
	if targetProjectId == fromProjectId {
		targetProjectBo = fromProjectBo
	} else if targetProjectId != fromProjectId && targetProjectId > 0 {
		targetProjectBo, errSys = GetProjectSimple(orgId, targetProjectId)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
	}

	// from/target table
	if fromTableId > 0 {
		fromTable, errSys = GetTableByTableId(orgId, operatorId, fromTableId)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}

		fromTableSchema, errSys := GetTableColumnsMap(orgId, fromTableId, nil, true)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
		for columnId, column := range fromTableSchema {
			if column.Field.Type == tablePb.ColumnType_relating.String() ||
				column.Field.Type == tablePb.ColumnType_baRelating.String() ||
				column.Field.Type == tablePb.ColumnType_singleRelating.String() {
				relatingColumnIds = append(relatingColumnIds, columnId)
			}
		}
	}
	if targetTableId == fromTableId {
		targetTable = fromTable
	} else if targetTableId != fromTableId && targetTableId > 0 {
		targetTable, errSys = GetTableByTableId(orgId, operatorId, targetTableId)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
	}

	// target app
	var fromAppId int64
	var targetAppId int64
	var summaryAppId int64
	if fromProjectBo != nil {
		fromAppId = fromProjectBo.AppId
	}
	if targetProjectBo != nil {
		targetAppId = targetProjectBo.AppId
	}
	summaryAppId, errSys = GetOrgSummaryAppId(orgId)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}
	if fromAppId == 0 {
		fromAppId = summaryAppId
	}
	if targetAppId == 0 {
		targetAppId = summaryAppId
	}

	nowTimeStr := time.Now().Format(consts.AppTimeFormat)

	changeToRootIds := make([]int64, 0) // 需要转化为顶级父任务的
	targetChangePath := make(map[string]string)
	fromChangePath := make(map[string]string)

	for _, issueBo := range issueAndChildrenBos {
		if ok, _ := slice.Contain(allMoveIds, issueBo.Id); ok {
			// 需要移动的任务
			// 是子任务
			if issueBo.ParentId != 0 {
				// 如果父任务不在需要移动的里面，该任务需要转为顶级任务并移动到新项目表
				if ok, _ = slice.Contain(allMoveIds, issueBo.ParentId); !ok {
					changeToRootIds = append(changeToRootIds, issueBo.Id)
					targetChangePath[fmt.Sprintf("%s%d,", issueBo.Path, issueBo.Id)] = fmt.Sprintf("0,%d,", issueBo.Id)
				}
			}

		} else {
			// 不需要移动的任务
			// 是子任务
			if issueBo.ParentId != 0 {
				// 如果父任务在需要移动的里面，该任务需要转为顶级任务并保留在原项目表
				if ok, _ = slice.Contain(allMoveIds, issueBo.ParentId); ok {
					changeToRootIds = append(changeToRootIds, issueBo.Id)
					fromChangePath[fmt.Sprintf("%s%d,", issueBo.Path, issueBo.Id)] = fmt.Sprintf("0,%d,", issueBo.Id)
				}
			}
		}
	}

	var updateBatchReqs []*formvo.LessUpdateIssueBatchReq

	// 需要移动的任务涉及的关联引用的任务ID
	if len(relatingColumnIds) > 0 {
		var allRelatingIssueIds []int64
		for _, issueBo := range issueAndChildrenBos {
			if ok, _ := slice.Contain(allMoveIds, issueBo.Id); ok {
				// 需要移动的任务
				for _, columnId := range relatingColumnIds {
					if v, ok := issueBo.LessData[columnId]; ok {
						oldRelating := &bo.RelatingIssue{}
						jsonx.Copy(v, oldRelating)
						allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkTo)...)
						allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkFrom)...)
					}
				}
			}
		}
		allRelatingIssueIds = slice.SliceUniqueInt64(allRelatingIssueIds)
		var relatingIssueIds []int64
		for _, id := range allRelatingIssueIds {
			if ok, _ := slice.Contain(allMoveIds, id); !ok {
				relatingIssueIds = append(relatingIssueIds, id)
			}
		}
		if len(relatingIssueIds) > 0 {
			updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
				OrgId:  orgId,
				AppId:  summaryAppId, // 目标项目appId
				UserId: operatorId,
				Condition: vo.LessCondsData{
					Type: consts.ConditionAnd,
					Conds: []*vo.LessCondsData{
						{
							Type:   consts.ConditionEqual,
							Value:  orgId,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
						},
						{
							Type:   consts.ConditionIn,
							Values: relatingIssueIds,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
						},
					},
				},
				Sets: []datacenter.Set{
					{
						Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldUpdateTime),
						Value:           "NOW()",
						Type:            consts.SetTypeNormal,
						Action:          consts.SetActionSet,
						WithoutPretreat: true,
					},
				},
			})
		}
	}

	// 更新带上项目id，防止并发导致脏数据
	upd := mysql.Upd{
		consts.TcUpdator:   operatorId,
		consts.TcProjectId: targetProjectId,
		consts.TcTableId:   targetTableId,
	}

	// 任务状态，如果目标项目是空项目，没有issueStatus，不去请求了
	projectTypeId := int64(0)
	if targetProjectBo != nil {
		projectTypeId = targetProjectBo.ProjectTypeId
	}
	var allNewStatus []status.StatusInfoBo
	if projectTypeId != consts.ProjectTypeEmpty {
		allNewStatus, errSys = GetTableStatus(orgId, targetTableId)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
	}

	newStatusType := consts.StatusTypeNotStart
	var newStatus int64
	if isChangeTable {
		for _, status := range allNewStatus {
			if status.Type == newStatusType {
				newStatus = status.ID
				upd[consts.TcStatus] = status.ID
				break
			}
		}
	}

	// 确认状态
	if isChangeTable {
		if targetProjectBo != nil && targetProjectBo.ProjectTypeId == consts.ProjectTypeNormalId {
			upd[consts.TcAuditStatus] = consts.AuditStatusNotView
		} else {
			upd[consts.TcAuditStatus] = consts.AuditStatusNoNeed
		}
	}

	// 迭代信息
	//if isChangeProject {
	//	upd[consts.TcIterationId] = 0
	//}
	// 迭代
	if issueAndChildrenBos[0].IterationId != 0 {
		if !str.CheckStrInArray(chooseFields, consts.BasicFieldIterationId) {
			upd[consts.TcIterationId] = 0
		}
	}

	// 开始时间
	if !str.CheckStrInArray(chooseFields, consts.BasicFieldPlanStartTime) {
		upd[consts.TcPlanStartTime] = consts.BlankTime
	}

	// 截止时间
	if !str.CheckStrInArray(chooseFields, consts.BasicFieldPlanStartTime) {
		upd[consts.TcPlanEndTime] = consts.BlankTime
	}

	var updForm []map[string]interface{}     // 移动项目表的任务
	var updFormFrom []map[string]interface{} // 不移动项目表的任务

	var (
		err     errs.SystemErrorInfo
		columns []*projectvo.TableColumnData
	)
	if fromTableId != 0 {
		columns, err = GetTableColumns(orgId, operatorId, fromTableId, false)
		if err != nil {
			log.Errorf("[MoveIssueProTableBatch] GetTableColumns tableId:%v, err:%v", fromTableId, err)
			return nil, err
		}
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 刷parentId path
		if len(changeToRootIds) > 0 {
			//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
			//	consts.TcOrgId: orgId,
			//	consts.TcId:    db.In(changeToRootIds),
			//}, mysql.Upd{
			//	consts.TcParentId: 0,
			//	consts.TcPath:     rootPath,
			//})
			//if err != nil {
			//	log.Error(err)
			//	return err
			//}

			// 移动到目标项目表里的任务，刷path
			for oldPath, newPath := range targetChangePath {
				//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
				//	consts.TcOrgId: orgId,
				//	consts.TcId:    db.In(allMoveIds),
				//}, mysql.Upd{
				//	consts.TcPath: db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, oldPath, newPath)),
				//})
				//if err != nil {
				//	log.Error(err)
				//	return err
				//}

				// 无码的请求体
				updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
					OrgId:  orgId,
					AppId:  targetAppId, // 目标项目appId
					UserId: operatorId,
					Condition: vo.LessCondsData{
						Type: consts.ConditionAnd,
						Conds: []*vo.LessCondsData{
							{
								Type:   consts.ConditionEqual,
								Value:  orgId,
								Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
							},
							{
								Type:   consts.ConditionIn,
								Values: allMoveIds,
								Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
							},
						},
					},
					Sets: []datacenter.Set{
						{
							Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
							Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, oldPath, newPath),
							Type:            consts.SetTypeNormal,
							Action:          consts.SetActionSet,
							WithoutPretreat: true,
						},
					},
				})
			}
			// 留在原项目表里的任务，刷path
			for oldPath, newPath := range fromChangePath {
				//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
				//	consts.TcOrgId: orgId,
				//	consts.TcId:    db.In(allRemainIds),
				//}, mysql.Upd{
				//	consts.TcPath: db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, oldPath, newPath)),
				//})
				//if err != nil {
				//	log.Error(err)
				//	return err
				//}

				// 无码的请求体
				updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
					OrgId:  orgId,
					AppId:  fromAppId, // 原项目appId
					UserId: operatorId,
					Condition: vo.LessCondsData{
						Type: consts.ConditionAnd,
						Conds: []*vo.LessCondsData{
							{
								Type:   consts.ConditionEqual,
								Value:  orgId,
								Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
							},
							{
								Type:   consts.ConditionIn,
								Values: allRemainIds,
								Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
							},
						},
					},
					Sets: []datacenter.Set{
						{
							Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
							Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, oldPath, newPath),
							Type:            consts.SetTypeNormal,
							Action:          consts.SetActionSet,
							WithoutPretreat: true,
						},
					},
				})
			}
		}

		parentForm := slice2.CaseCamelCopy(upd)
		if isChangeTable {
			parentForm[consts.BasicFieldMoveTime] = nowTimeStr
		}

		for _, id := range changeToRootIds {
			noNeedMoveForm := map[string]interface{}{}
			copyer.Copy(parentForm, &noNeedMoveForm)
			delete(noNeedMoveForm, consts.BasicFieldTableId)
			delete(noNeedMoveForm, consts.BasicFieldProjectId)
			delete(noNeedMoveForm, consts.TcUpdator)
			delete(noNeedMoveForm, consts.TcStatus)
			noNeedMoveForm[consts.BasicFieldIssueId] = id
			noNeedMoveForm[consts.BasicFieldParentId] = 0
			noNeedMoveForm[consts.BasicFieldPath] = rootPath
			updFormFrom = append(updFormFrom, noNeedMoveForm)
			//updFormFrom = append(updFormFrom, map[string]interface{}{
			//	consts.BasicFieldIssueId:  id,
			//	consts.BasicFieldParentId: 0,
			//	consts.BasicFieldPath:     rootPath,
			//})
		}

		// 更新remark
		if !str.CheckStrInArray(chooseFields, consts.BasicFieldRemark) {
			parentForm[consts.BasicFieldRemark] = consts.BlankString
			//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueDetail, db.Cond{
			//	consts.TcOrgId:    orgId,
			//	consts.TcIsDelete: consts.AppIsNoDelete,
			//	consts.TcIssueId:  db.In(allMoveIds),
			//}, mysql.Upd{
			//	consts.TcUpdator:   operatorId,
			//	consts.TcProjectId: targetProjectId,
			//	consts.TcRemark:    consts.BlankString,
			//})
			//if err != nil {
			//	log.Error(err)
			//	return err
			//}
			//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssueDetail, db.Cond{
			//	consts.TcOrgId:    orgId,
			//	consts.TcIsDelete: consts.AppIsNoDelete,
			//	consts.TcIssueId:  db.In(allRemainIds),
			//}, mysql.Upd{
			//	consts.TcUpdator: operatorId,
			//	consts.TcRemark:  consts.BlankString,
			//})
			//if err != nil {
			//	log.Error(err)
			//	return err
			//}
		}

		// 更新issueRelation表
		if !str.CheckStrInArray(chooseFields, consts.BasicFieldOwnerId) {
			parentForm[consts.BasicFieldOwnerId] = []string{}
			_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
				consts.TcOrgId:        orgId,
				consts.TcIssueId:      db.In(inputIssueIds),
				consts.TcRelationType: consts.IssueRelationTypeOwner,
				consts.TcIsDelete:     consts.AppIsNoDelete,
			}, mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
			if err != nil {
				log.Error(err)
				return err
			}
		}
		if !str.CheckStrInArray(chooseFields, consts.BasicFieldFollowerIds) {
			parentForm[consts.BasicFieldFollowerIds] = []string{}
			_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
				consts.TcOrgId:        orgId,
				consts.TcIssueId:      db.In(inputIssueIds),
				consts.TcRelationType: consts.IssueRelationTypeFollower,
				consts.TcIsDelete:     consts.AppIsNoDelete,
			}, mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
			if err != nil {
				log.Error(err)
				return err
			}
		}
		if !str.CheckStrInArray(chooseFields, consts.BasicFieldAuditorIds) {
			parentForm[consts.BasicFieldAuditorIds] = []string{}
			parentForm[consts.BasicFieldAuditStatusDetail] = map[string]int{}
			//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
			//	consts.TcOrgId:        orgId,
			//	consts.TcIssueId:      db.In(inputIssueIds),
			//	consts.TcRelationType: consts.IssueRelationTypeAuditor,
			//	consts.TcIsDelete:     consts.AppIsNoDelete,
			//}, mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
			//if err != nil {
			//	log.Error(err)
			//	return err
			//}
		}

		if !str.CheckStrInArray(chooseFields, consts.BasicFieldPriority) {
			parentForm[consts.BasicFieldPriority] = nil
		}

		// 工时
		if !str.CheckStrInArray(chooseFields, consts.BasicFieldWorkHour) {
			parentForm[consts.BasicFieldWorkHour] = consts.DefaultWorkHour
		}

		//appIdUpdateForm := map[int64][]map[string]interface{}{}
		// 关联
		if !str.CheckStrInArray(chooseFields, consts.BasicFieldRelating) {
			parentForm[consts.BasicFieldRelating] = &bo.RelatingIssue{LinkTo: []string{}, LinkFrom: []string{}}
			//for _, issueBo := range issueAndChildrenBos {
			//	relatingForm, err := getNoMoveRelatingForm(orgId, operatorId, issueBo, &issueBo.RelatingIssue, consts.BasicFieldRelating)
			//	if err != nil {
			//		log.Errorf("[MoveIssueProTable] err:%v", err)
			//		return err
			//	}
			//	for appId, form := range relatingForm {
			//		if _, ok := appIdUpdateForm[appId]; ok {
			//			appIdUpdateForm[appId] = append(appIdUpdateForm[appId], form...)
			//		} else {
			//			appIdUpdateForm[appId] = form
			//		}
			//	}
			//}
		}
		// 单向关联和双向关联移动后，都需要删除，目前还没有办法保留
		for _, column := range columns {
			if column.Field.Type == tablePb.ColumnType_singleRelating.String() {
				parentForm[column.Name] = &bo.RelatingIssue{LinkTo: []string{}, LinkFrom: []string{}}
			}
			if column.Field.Type == tablePb.ColumnType_relating.String() {
				props := column.Field.Props[column.Field.Type]
				if m, ok := props.(map[string]interface{}); ok {
					tableId := cast.ToInt64(m[consts.RelateTableId])
					if tableId != 0 {
						parentForm[column.Name] = &bo.RelatingIssue{LinkTo: []string{}, LinkFrom: []string{}}
					}
				}
			}
		}

		// 前后置
		if !str.CheckStrInArray(chooseFields, consts.BasicFieldBaRelating) {
			parentForm[consts.BasicFieldBaRelating] = &bo.RelatingIssue{LinkTo: []string{}, LinkFrom: []string{}}
			//for _, issueBo := range issueAndChildrenBos {
			//	baRelatingForm, err := getNoMoveRelatingForm(orgId, operatorId, issueBo, &issueBo.BaRelatingIssue, consts.BasicFieldBaRelating)
			//	if err != nil {
			//		log.Errorf("[MoveIssueProTable] err:%v", err)
			//		return err
			//	}
			//	for appId, form := range baRelatingForm {
			//		if _, ok := appIdUpdateForm[appId]; ok {
			//			appIdUpdateForm[appId] = append(appIdUpdateForm[appId], form...)
			//		} else {
			//			appIdUpdateForm[appId] = form
			//		}
			//	}
			//}
		}

		// 刷upd
		//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcId:       db.In(allMoveIds),
		//	consts.TcIsDelete: consts.AppIsNoDelete,
		//}, upd)
		//if err != nil {
		//	log.Error(err)
		//	return err
		//}
		for _, id := range allMoveIds {
			tempUpdForm := map[string]interface{}{}
			copyer.Copy(parentForm, &tempUpdForm)
			tempUpdForm[consts.BasicFieldIssueId] = id
			if isChangeTable {
				tempUpdForm[consts.BasicFieldIssueStatus] = newStatus
				tempUpdForm[consts.BasicFieldIssueStatusType] = consts.StatusTypeNotStart
			}
			tempUpdForm[consts.BasicFieldAppId] = cast.ToString(targetAppId)
			tempUpdForm[consts.BasicFieldTableId] = cast.ToString(targetTableId) // 转成字符串
			delete(tempUpdForm, consts.TcStatus)
			if v, ok := issueTitleMap[id]; ok {
				tempUpdForm[consts.BasicFieldTitle] = v
			}
			updForm = append(updForm, tempUpdForm)
		}

		// 无码更新: 先在原项目表中刷新parentId/path等信息，再去目标项目表刷projectId/tableId，这样update最终在正确的table表头中执行
		if len(updFormFrom) > 0 {
			resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
				AppId:   fromAppId,
				OrgId:   orgId,
				UserId:  operatorId,
				TableId: fromTableId,
				Form:    updFormFrom,
			})
			if resp.Failure() {
				log.Error(resp.Error())
				return resp.Error()
			}
		}
		if len(updForm) > 0 {
			resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
				AppId:   targetAppId,
				OrgId:   orgId,
				UserId:  operatorId,
				TableId: targetTableId,
				Form:    updForm,
			})
			if resp.Failure() {
				log.Error(resp.Error())
				return resp.Error()
			}
		}

		//for appId, form := range appIdUpdateForm {
		//	resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
		//		AppId:  appId,
		//		OrgId:  orgId,
		//		UserId: operatorId,
		//		Form:   form,
		//	})
		//	if resp.Failure() {
		//		log.Errorf("[MoveIssueProTableBatch] err:%v", resp.Error())
		//		return resp.Error()
		//	}
		//}
		// 无码更新
		for _, req := range updateBatchReqs {
			batchResp := formfacade.LessUpdateIssueBatchRaw(req)
			if batchResp.Failure() {
				log.Error(batchResp.Error())
				return batchResp.Error()
			}
		}
		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}

	// 更新任务工时关联的 projectId
	if isChangeProject {
		isDeleteWorkHour := false
		if !str.CheckStrInArray(chooseFields, consts.BasicFieldWorkHour) {
			// 不移动工时，就需要把工时关联删除
			isDeleteWorkHour = true
		}
		errSys = UpdateWorkHourProjectId(orgId, fromProjectId, targetProjectId, allMoveIds, isDeleteWorkHour)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
	}

	// 移动任务附件
	//if isChangeProject {
	//	errSys = switchIssueResource(orgId, allMoveIds, operatorId, targetProjectId)
	//	if errSys != nil {
	//		log.Error(errSys)
	//		return nil, errSys
	//	}
	//}

	if isChangeProject && targetProjectId > 0 {
		// 处理任务成员和关注人
		errSys = switchIssueMember(orgId, allMoveIds, operatorId, targetProjectId)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}

	}

	//asyn.Execute(func() {
	//	if isChangeTable {
	//		if err := IssueChatInviteUsersWhenMoveIssue(orgId, allMoveIds, targetTableId); err != nil {
	//			log.Errorf("[MoveIssueProTableBatch] IssueChatInviteUserWhenMoveIssue err: %v, issueIds: %s", err, json.ToJsonIgnoreError(allMoveIds))
	//			return
	//		}
	//	}
	//})

	asyn.Execute(func() {
		if isChangeProject || isChangeTable {
			pushType := consts.PushTypeUpdateIssueProjectTable
			for _, issueBo := range allMoveIssues {
				changeList := []bo.TrendChangeListBo{}
				//if isChangeProject {
				var oldProjectName, newProjectName string
				if fromProjectBo != nil {
					oldProjectName = fromProjectBo.Name
				}
				if targetProjectBo != nil {
					newProjectName = targetProjectBo.Name
				}
				changeList = append(changeList, bo.TrendChangeListBo{
					Field:     consts.BasicFieldProjectId,
					FieldName: consts.Project,
					OldValue:  oldProjectName,
					NewValue:  newProjectName,
				})
				//}
				//if isChangeTable {
				var oldTableName, newTableName string
				if fromTable != nil {
					oldTableName = fromTable.Name
				}
				if targetTable != nil {
					newTableName = targetTable.Name
				}
				changeList = append(changeList, bo.TrendChangeListBo{
					Field:     consts.BasicFieldTableId,
					FieldName: consts.Table,
					OldValue:  oldTableName,
					NewValue:  newTableName,
				})
				//}

				oldValue := bo.ProjectTableBo{}
				oldValue.ProjectId = fromProjectId
				oldValue.TableId = fromTableId

				newValue := bo.ProjectTableBo{}
				newValue.ProjectId = targetProjectId
				newValue.TableId = targetTableId

				issueTrendsBo := &bo.IssueTrendsBo{
					PushType:      pushType,
					OrgId:         issueBo.OrgId,
					OperatorId:    operatorId,
					DataId:        issueBo.DataId,
					IssueId:       issueBo.Id,
					ParentIssueId: issueBo.ParentId,
					ProjectId:     targetProjectId,
					TableId:       targetTableId,
					PriorityId:    issueBo.PriorityId,
					IssueTitle:    issueBo.Title,
					ParentId:      issueBo.ParentId,
					OldValue:      json.ToJsonIgnoreError(oldValue),
					NewValue:      json.ToJsonIgnoreError(newValue),
					Ext: bo.TrendExtensionBo{
						ObjName:    issueBo.Title,
						ChangeList: changeList,
					},
				}
				PushIssueTrends(issueTrendsBo)

				//asyn.Execute(func() {
				//	PushIssueThirdPlatformNotice(issueTrendsBo)
				//})
				//asyn.Execute(func() {
				//	//推送群聊卡片
				//	PushInfoToChat(issueBo.OrgId, targetProjectId, issueTrendsBo)
				//})
			}
		}
	})
	return allMoveIds, nil
}

// 切换任务资源所属项目
func switchIssueResource(orgId int64, issueIds []int64, operatorId int64, projectId int64) errs.SystemErrorInfo {
	resp := resourcefacade.UpdateResourceRelationProjectId(resourcevo.UpdateResourceRelationProjectIdReqVo{
		OrgId:  orgId,
		UserId: operatorId,
		Input: resourcevo.UpdateResourceRelationProjectIdData{
			ProjectId: projectId,
			IssueIds:  issueIds,
		},
	})
	return resp.Error()
}

// 切换任务相关成员：负责人直接带过去，其余的如果是目标项目的人则带过去
func switchIssueMember(orgId int64, issueIds []int64, operatorId int64, projectId int64) errs.SystemErrorInfo {
	// 获取目标项目的所有成员
	var newMemberList []po.PpmProProjectRelation
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: db.In(consts.MemberRelationTypeList),
	}, &newMemberList)
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}
	var newMemberIds []int64
	for _, relation := range newMemberList {
		newMemberIds = append(newMemberIds, relation.RelationId)
	}

	// 获取目标任务的所有负责人
	var ownerList []po.PpmPriIssueRelation
	err = mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      db.In(issueIds),
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeOwner}),
	}, &ownerList)
	// 需要新增到目标项目的人
	var ownerIds []int64
	for _, relation := range ownerList {
		if ok, _ := slice.Contain(newMemberIds, relation.RelationId); !ok {
			ownerIds = append(ownerIds, relation.RelationId)
		}
	}

	// 不属于新项目成员的关注人和参与人直接去除(和负责人相关的不移出，因为会加入新的项目)
	allRelateIds := append(newMemberIds, ownerIds...)
	_, updateErr := mysql.UpdateSmartWithCond(consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      db.In(issueIds),
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower}),
		consts.TcRelationId:   db.NotIn(allRelateIds),
	}, mysql.Upd{
		consts.TcUpdator:  operatorId,
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	if updateErr != nil {
		log.Error(updateErr)
		return errs.MysqlOperateError
	}
	// 负责人移动到新项目
	moveErr := UpdateProjectRelation(operatorId, orgId, projectId, consts.IssueRelationTypeOwner, ownerIds)
	if moveErr != nil {
		log.Error(moveErr)
		return moveErr
	}

	// 属于新项目的关注人和参与人带过去，负责人直接带过去
	upd := mysql.Upd{
		consts.TcUpdator:   operatorId,
		consts.TcProjectId: projectId,
	}
	conn, err := mysql.GetConnect()
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	_, err = conn.Update(consts.TableIssueRelation).Set(upd).Where(db.And(
		db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIssueId:  db.In(issueIds),
			consts.TcIsDelete: consts.AppIsNoDelete,
		},
		db.Or(
			db.Cond{
				consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower}),
				consts.TcRelationId:   db.In(allRelateIds),
			},
			db.Cond{
				consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeOwner}),
			},
		),
	)).Exec()
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return errs.IssueRelationUpdateError
	}
	return nil
}

//func GetDefaultPriority(orgId int64) (bo.PriorityBo, errs.SystemErrorInfo) {
//	resBo := bo.PriorityBo{}
//	priorityList, err := GetPriorityListByType(orgId, consts.PriorityTypeIssue)
//	if err != nil {
//		log.Error(err)
//		return resBo, err
//	}
//	for _, priorityBo := range *priorityList {
//		if priorityBo.LangCode == consts.PriorityIssueCommon {
//			resBo = priorityBo
//			break
//		}
//	}
//	if resBo.Id == 0 && len(*priorityList) > 0 {
//		resBo = (*priorityList)[0]
//	}
//	return resBo, nil
//}

func SetWorkHourMap(orgId int64, projectIds, issueIds []int64, workHourMap map[int64]bo.IssueWorkHour) {
	// 查询项目下这些任务的所有工时信息。
	workHourList, err := GetIssuesWorkHourList(orgId, projectIds, issueIds)
	if err != nil {
		log.Error(err)
		return
	}
	workerIds := []int64{}
	// 以 issueId 为 key
	var workHourListKeyByIssueId = make(map[int64][]*bo.SimpleWorkHourBo, 0)
	for _, item := range workHourList {
		workerIds = append(workerIds, item.WorkerId)
		if _, ok := workHourListKeyByIssueId[item.IssueId]; ok {
			workHourListKeyByIssueId[item.IssueId] = append(workHourListKeyByIssueId[item.IssueId], item)
		} else {
			workHourListKeyByIssueId[item.IssueId] = []*bo.SimpleWorkHourBo{item}
		}
	}
	workerIds = int642.ArrayUnique(workerIds)
	// 查询用户信息
	// 查询工时执行人的信息
	userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, workerIds)
	if err != nil {
		log.Error(err)
		return
	}
	userInfoMap := map[int64]*bo.BaseUserInfoBo{}
	for _, user := range userInfos {
		copiedUser := user
		userInfoMap[user.UserId] = &copiedUser
	}
	for issueId, tmpWorkHourList := range workHourListKeyByIssueId {
		// 标识某个任务计算预估工时，是以总预估方式计算，还是以子预估方式计算。
		predictIsSub := false
		// workerId 为 key 的记录
		oneGroupList := map[int64]*bo.WorkHourForSomeoneBo{}
		// 从 tmpWorkHourList 中，按 worker 进行分组
		for _, workHour := range tmpWorkHourList {
			if _, ok := oneGroupList[workHour.WorkerId]; !ok {
				oneGroupList[workHour.WorkerId] = &bo.WorkHourForSomeoneBo{
					TotalPredict: []*bo.SimpleWorkHourBo{},
					SubPredict:   []*bo.SimpleWorkHourBo{},
					Actual:       []*bo.SimpleWorkHourBo{},
				}
			}
			switch workHour.Type {
			case consts2.WorkHourTypeTotalPredict:
				oneGroupList[workHour.WorkerId].TotalPredict = append(oneGroupList[workHour.WorkerId].TotalPredict, workHour)
			case consts2.WorkHourTypeSubPredict:
				oneGroupList[workHour.WorkerId].SubPredict = append(oneGroupList[workHour.WorkerId].SubPredict, workHour)
				predictIsSub = true
			case consts2.WorkHourTypeActual:
				oneGroupList[workHour.WorkerId].Actual = append(oneGroupList[workHour.WorkerId].Actual, workHour)
			}
		}
		// 如果有子预估工时，则总预估工时用 SubPredict 统计，否则使用 TotalPredict 统计。todo
		tmpIssueData := &bo.IssueWorkHour{
			PredictWorkHour: "0",
			ActualWorkHour:  "0",
			PredictList:     []*bo.PredictListItem{},
			ActualList:      []*bo.ActualListItem{},
		}
		oneIssuePredictInt := uint32(0)
		oneIssueActualInt := uint32(0)
		for workerId, item := range oneGroupList {
			var userInfo *bo.BaseUserInfoBo
			if val, ok := userInfoMap[workerId]; ok {
				userInfo = val
			} else {
				log.Errorf("用户找不到 id 为 %[1]v 的用户信息。", workerId)
				continue
			}
			onePersonPredictInt := uint32(0)
			onePersonActualInt := uint32(0)
			// 默认用总预估记录作为这个人的预估工时。
			tmpPredictList := item.TotalPredict
			// 如果有子预估工时，则使用子预估工时计算这个人的预估工时。
			if len(item.SubPredict) > 0 {
				tmpPredictList = item.SubPredict
			}
			for _, oneWorkHour := range tmpPredictList {
				onePersonPredictInt += oneWorkHour.NeedTime
			}
			if predictIsSub {
				tmpIssueData.PredictList = append(tmpIssueData.PredictList, &bo.PredictListItem{
					Name:     userInfo.Name,
					WorkHour: format.FormatNeedTimeIntoString(int64(onePersonPredictInt)),
				})
			}
			// 实际工时需要按日期维度分组。
			var tmpDateList []*bo.ActualWorkHourDateItem
			tmpDateMap := map[string]*bo.ActualWorkHourDateItemForInt{}
			for _, oneWorkHour := range item.Actual {
				onePersonActualInt += oneWorkHour.NeedTime
				dateStr := time.Unix(int64(oneWorkHour.StartTime), 0).Format(consts.AppDateFormat)
				if _, ok := tmpDateMap[dateStr]; ok {
					tmpDateMap[dateStr].WorkHourInt += oneWorkHour.NeedTime
				} else {
					tmpDateMap[dateStr] = &bo.ActualWorkHourDateItemForInt{
						Date:        dateStr,
						WorkHourInt: oneWorkHour.NeedTime,
					}
				}
			}
			for dateKey, val := range tmpDateMap {
				tmpDateList = append(tmpDateList, &bo.ActualWorkHourDateItem{
					Date:     dateKey,
					WorkHour: format.FormatNeedTimeIntoString(int64(val.WorkHourInt)),
				})
			}
			tmpIssueData.ActualList = append(tmpIssueData.ActualList, &bo.ActualListItem{
				Name:                   userInfo.Name,
				ActualWorkHourDateList: tmpDateList,
			})
			oneIssuePredictInt += onePersonPredictInt
			oneIssueActualInt += onePersonActualInt
		}
		tmpIssueData.PredictWorkHour = format.FormatNeedTimeIntoString(int64(oneIssuePredictInt))
		tmpIssueData.ActualWorkHour = format.FormatNeedTimeIntoString(int64(oneIssueActualInt))
		workHourMap[issueId] = *tmpIssueData
	}
}

//func GetTrulyIssueInfo(orgId, issueId int64) (*bo.IssueBo, errs.SystemErrorInfo) {
//	issue := &po.PpmPriIssue{}
//	err := mysql.SelectOneByCond(issue.TableName(), db.Cond{
//		consts.TcId:    issueId,
//		consts.TcOrgId: orgId,
//	}, issue)
//	if err != nil {
//		return nil, errs.IllegalityIssue
//	}
//	issueBo := &bo.IssueBo{}
//	err1 := copyer.Copy(issue, issueBo)
//	if err1 != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
//	}
//	return issueBo, nil
//}

//func GetIssueBo(orgId, issueId int64, isDelete ...int) (*bo.IssueBo, errs.SystemErrorInfo) {
//	var isDeleteI interface{}
//	isDeleteI = consts.AppIsNoDelete
//	if isDelete != nil && len(isDelete) > 0 {
//		isDeleteI = db.In(isDelete)
//	}
//	issue := &po.PpmPriIssue{}
//	err := mysql.SelectOneByCond(issue.TableName(), db.Cond{
//		consts.TcId:       issueId,
//		consts.TcOrgId:    orgId,
//		consts.TcIsDelete: isDeleteI,
//	}, issue)
//	if err != nil {
//		return nil, errs.IllegalityIssue
//	}
//	issueBo := &bo.IssueBo{}
//	err1 := copyer.Copy(issue, issueBo)
//	if err1 != nil {
//		log.Errorf("[GetIssueBo] err: %v", err1)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
//	}
//	return issueBo, nil
//}

// GetIssueBoIncludeDeleted 查询任务信息，包含删除和未删除的。
//func GetIssueBoIncludeDeleted(orgId, issueId int64) (*bo.IssueBo, errs.SystemErrorInfo) {
//	issue := &po.PpmPriIssue{}
//	err := mysql.SelectOneByCond(issue.TableName(), db.Cond{
//		consts.TcId:    issueId,
//		consts.TcOrgId: orgId,
//	}, issue)
//	if err != nil {
//		return nil, errs.IllegalityIssue
//	}
//	issueBo := &bo.IssueBo{}
//	err1 := copyer.Copy(issue, issueBo)
//	if err1 != nil {
//		log.Errorf("[GetIssueBoIncludeDeleted] err: %v", err1)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
//	}
//	return issueBo, nil
//}

//func GetCountIssueByProjectObjectTypeId(projectObjecTypeId, projectId, orgId int64) (bool, errs.SystemErrorInfo) {
//	issue := &po.PpmPriIssue{}
//	count, err := mysql.SelectCountByCond(issue.TableName(), db.Cond{
//		consts.TcProjectObjectTypeId: projectObjecTypeId,
//		consts.TcProjectId:           projectId,
//		consts.TcOrgId:               orgId,
//		consts.TcIsDelete:            consts.AppIsNoDelete,
//	})
//
//	if err != nil {
//		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	result := int(count) > 0
//	return result, nil
//}

func GetIssueAuthBo(issueBo *bo.IssueBo, currentUserId int64) (*bo.IssueAuthBo, errs.SystemErrorInfo) {
	issueId := issueBo.Id

	ownerIds, err2 := GetIssueRelationIdsByRelateType(issueBo.OrgId, issueId, consts.IssueRelationTypeOwner)
	if err2 != nil {
		log.Error(err2)
		return nil, err2
	}
	issueAuthBo := &bo.IssueAuthBo{
		Id:        issueId,
		Owner:     *ownerIds,
		Creator:   issueBo.Creator,
		ProjectId: issueBo.ProjectId,
		Status:    issueBo.Status,
		TableId:   issueBo.TableId,
	}

	// 2019-12-24-nico: 产品想要支持父任务负责人操作所有子任务
	if issueBo.ParentId != 0 {
		issueBos, err := GetIssueInfosLc(issueBo.OrgId, currentUserId, []int64{issueBo.ParentId})
		if err != nil {
			log.Errorf("[GetIssueAuthBo] err: %v, parentId: %d", err, issueBo.ParentId)
			return issueAuthBo, err
		}
		if len(issueBos) <= 0 {
			return issueAuthBo, err
		}
		parentIssueBo := issueBos[0]
		parentOwnerIds, errSys := businees.LcMemberToUserIdsWithError(parentIssueBo.OwnerId)
		if errSys != nil {
			log.Errorf("[GetIssueAuthBo] parent issue ownerId: %v, err: %v", parentIssueBo.OwnerId, errSys)
			return nil, errSys
		}

		if ok, _ := slice.Contain(parentOwnerIds, currentUserId); ok {
			issueAuthBo.Owner = append(issueAuthBo.Owner, currentUserId)
		}
	}

	return issueAuthBo, nil
}

//func GetIssueBoList(issueListCond bo.IssueBoListCond) ([]bo.IssueBo, errs.SystemErrorInfo) {
//	issueList := &[]*po.PpmPriIssue{}
//
//	cond := db.Cond{
//		consts.TcIsDelete: consts.AppIsNoDelete,
//		consts.TcOrgId:    issueListCond.OrgId,
//	}
//	if issueListCond.ProjectId != nil {
//		cond[consts.TcProjectId] = *issueListCond.ProjectId
//	}
//	if issueListCond.ProjectIds != nil && len(issueListCond.ProjectIds) > 0 {
//		cond[consts.TcProjectId] = db.In(issueListCond.ProjectIds)
//	}
//	if issueListCond.IterationId != nil {
//		cond[consts.TcIterationId] = *issueListCond.IterationId
//	}
//	if issueListCond.Ids != nil {
//		cond[consts.TcId] = db.In(issueListCond.Ids)
//	}
//	err1 := mysql.SelectAllByCond(consts.TableIssue, cond, issueList)
//
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
//	}
//
//	bos := &[]bo.IssueBo{}
//	err2 := copyer.Copy(issueList, bos)
//	if err2 != nil {
//		log.Error(err2)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err2)
//	}
//	return *bos, nil
//}

// ConvertMapIntoIssueDetailAndUnionBo 将无码查询出来的任务信息转换为任务对象
func ConvertMapIntoIssueDetailAndUnionBo(issueInfoMap map[string]interface{}) bo.IssueAndDetailUnionBo {
	loc, _ := time.LoadLocation("Local")

	issueObj := bo.IssueAndDetailUnionBo{}
	i, ok := issueInfoMap["issueId"]
	if !ok {
		return issueObj
	}
	issueId := int64(0)
	if id, ok := i.(int64); ok {
		issueId = id
	} else if id, ok1 := i.(int); ok1 {
		issueId = int64(id)
	} else if id, ok1 := i.(float64); ok1 {
		//理论上 map 解析json 会将int转为float64
		issueIdStr := strconv.FormatFloat(id, 'f', -1, 64)
		parseId, err := strconv.ParseInt(issueIdStr, 10, 64)
		if err != nil {
			log.Errorf("[ConvertMapIntoIssueDetailAndUnionBo] convert issueId err: %v", err)
		} else {
			issueId = parseId
		}
	}
	issueObj.IssueId = issueId
	auditStatus := 0
	if val, ok := issueInfoMap[consts.BasicFieldAuditStatus]; ok && val != nil {
		if valInt, ok1 := val.(int); ok1 {
			auditStatus = valInt
		} else if valF, ok1 := val.(float64); ok1 {
			valStr := strconv.FormatFloat(valF, 'f', -1, 64)
			valInt, err := strconv.ParseInt(valStr, 10, 64)
			if err != nil {
				log.Errorf("[ConvertMapIntoIssueDetailAndUnionBo] convert auditStatus err: %v", err)
			} else {
				auditStatus = int(valInt)
			}
		}
	}
	issueObj.IssueAuditStatus = auditStatus

	issueStatus := int64(0)
	if val, ok := issueInfoMap[consts.BasicFieldIssueStatus]; ok && val != nil {
		if valInt, ok1 := val.(int); ok1 {
			issueStatus = int64(valInt)
		} else if valF, ok1 := val.(float64); ok1 {
			valStr := strconv.FormatFloat(valF, 'f', -1, 64)
			valInt, err := strconv.ParseInt(valStr, 10, 64)
			if err != nil {
				log.Errorf("[ConvertMapIntoIssueDetailAndUnionBo] convert issueStatus err: %v", err)
			} else {
				issueStatus = valInt
			}
		}
	}
	issueObj.IssueStatusId = issueStatus

	if val, ok := issueInfoMap[consts.BasicFieldIssueStatusType]; ok && val != nil {
		if valF, ok1 := val.(float64); ok1 {
			valStr := strconv.FormatFloat(valF, 'f', -1, 64)
			issueStatusType, err := strconv.ParseInt(valStr, 10, 64)
			if err != nil {
				log.Errorf("[ConvertMapIntoIssueDetailAndUnionBo] convert issueStatusType err: %v", err)
			} else {
				issueObj.IssueStatusType = int(issueStatusType)
			}
		}
	}
	if val, ok := issueInfoMap["tableId"]; ok && val != nil {
		if valStr, ok1 := val.(string); ok1 {
			tableId, err := strconv.ParseInt(valStr, 10, 64)
			if err != nil {
				log.Errorf("[ConvertMapIntoIssueDetailAndUnionBo] convert tableId err: %v", err)
			} else {
				issueObj.TableId = tableId
			}
		}
	}
	if val, ok := issueInfoMap[consts.ProBasicFieldCreateTime]; ok && val != nil {
		timeStr := val.(string)
		createTime, err := time.ParseInLocation(consts.AppTimeFormat, timeStr, loc)
		if err != nil {
			log.Errorf("[ConvertMapIntoIssueDetailAndUnionBo] convert createTime err: %v", err)
		}
		issueObj.CreateTime = createTime
	}
	if val := cast.ToString(issueInfoMap[consts.BasicFieldPlanEndTime]); val != "" {
		planEndTime, err := date.StrToTimeWithLoc(val, loc)
		if err != nil {
			log.Errorf("[ConvertMapIntoIssueDetailAndUnionBo] issueId: %d, timeStr: %v, convert planEndTime err1: %v",
				issueId, val, err)
		}
		issueObj.PlanEndTime = planEndTime
	}
	if val, ok := issueInfoMap["endTime"]; ok && val != nil {
		timeStr := val.(string)
		endTime, err := date.StrToTimeWithLoc(timeStr, loc)
		if err != nil {
			log.Errorf("[ConvertMapIntoIssueDetailAndUnionBo] convert endTime err: %v", err)
		}
		issueObj.EndTime = endTime
	}
	if val, ok := issueInfoMap["ownerChangeTime"]; ok && val != nil {
		timeStr := val.(string)
		ownerChangeTime, err := date.StrToTimeWithLoc(timeStr, loc)
		if err != nil {
			log.Errorf("[ConvertMapIntoIssueDetailAndUnionBo] convert ownerChangeTime err: %v", err)
		}
		issueObj.OwnerChangeTime = ownerChangeTime
	}

	return issueObj
}

//// UpdateIssueForTrueProjectValue 更新 `bo.IssueAndDetailUnionBo` 中的 projectObjectTypeId 值
//// 注意：这个方法会更新 list 的值。
//func UpdateIssueForTrueProjectValue(projectId int64, list *[]bo.IssueAndDetailUnionBo) errs.SystemErrorInfo {
//	relations, err := GetIssueProRelationsByProject(projectId)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	relationsMap := make(map[int64]bo.IssueRelationBo, len(relations))
//	for _, item := range relations {
//		relationsMap[item.IssueId] = item
//	}
//	var curProjectObjectTypeId int64
//	for index, item := range *list {
//		curProjectObjectTypeId = 0
//		if tmpRelationBo, ok := relationsMap[item.IssueId]; ok {
//			curProjectObjectTypeId = tmpRelationBo.RelationId
//		}
//		(*list)[index].IssueProjectObjectTypeId = curProjectObjectTypeId
//	}
//	return nil
//}

//func GetToNowCondIssueCount(projectId int64, timePoint time.Time, processStatus []int64) (int, errs.SystemErrorInfo) {
//	issue := &po.PpmPriIssue{}
//
//	//获取当天00:00
//	now := time.Now()
//	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
//
//	count, err := mysql.SelectCountByCond(issue.TableName(), db.Cond{
//		consts.TcProjectId: projectId,
//		consts.TcEndTime:   db.Between(startTime, timePoint),
//		consts.TcStatus:    db.In(processStatus),
//		consts.TcIsDelete:  consts.AppIsNoDelete,
//	})
//
//	if err != nil {
//		log.Error(err)
//		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//
//	return int(count), nil
//}

//func GetUnFinishIssueCount(projectId int64, processStatus []int64) (int, errs.SystemErrorInfo) {
//	issue := &po.PpmPriIssue{}
//	count, err := mysql.SelectCountByCond(issue.TableName(), db.Cond{
//		consts.TcStatus:    db.In(processStatus),
//		consts.TcIsDelete:  consts.AppIsNoDelete,
//		consts.TcProjectId: projectId,
//	})
//
//	if err != nil {
//		log.Error(err)
//		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//	return int(count), nil
//}

//func GetOverdueCondIssueCount(projectId int64, timePoint time.Time, processStatus []int64) (int, errs.SystemErrorInfo) {
//	issue := &po.PpmPriIssue{}
//	count, err := mysql.SelectCountByCond(issue.TableName(), db.Cond{
//		consts.TcProjectId:   projectId,
//		consts.TcStatus:      db.In(processStatus),
//		consts.TcIsDelete:    consts.AppIsNoDelete,
//		consts.TcPlanEndTime: db.Before(timePoint),
//	})
//
//	if err != nil {
//		log.Error(err)
//		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//	return int(count), nil
//}

//func GetAllChildrenIssues(orgId int64, parentId int64, childrenIds []int64) ([]bo.IssueBo, errs.SystemErrorInfo) {
//	issueChild := &[]po.PpmPriIssue{}
//	err := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
//		consts.TcOrgId:    orgId,
//		consts.TcParentId: parentId,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//		consts.TcId:       db.In(childrenIds),
//	}, issueChild)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.MysqlOperateError
//	}
//
//	result := &[]bo.IssueBo{}
//	copyErr := copyer.Copy(issueChild, result)
//	if copyErr != nil {
//		log.Error(copyErr)
//		return nil, errs.ObjectCopyError
//	}
//
//	if len(*result) == 0 {
//		return *result, nil
//	}
//	allIssueIds := []int64{}
//	for _, issueBo := range *result {
//		allIssueIds = append(allIssueIds, issueBo.Id)
//	}
//
//	details := &[]po.PpmPriIssueDetail{}
//	err1 := mysql.SelectAllByCond(consts.TableIssueDetail, db.Cond{
//		consts.TcOrgId:    orgId,
//		consts.TcIssueId:  db.In(allIssueIds),
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}, details)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.MysqlOperateError
//	}
//
//	detailMap := maps.NewMap("IssueId", details)
//	for i, issueBo := range *result {
//		if detail, ok := detailMap[issueBo.Id]; ok {
//			detailInfo := detail.(po.PpmPriIssueDetail)
//			(*result)[i].Remark = *detailInfo.Remark
//		}
//	}
//
//	return *result, nil
//}

// GetIssueInfoMapLcByIssueId 从无码获取单条数据
//func GetIssueInfoMapLcByIssueId(orgId, userId int64, issueId int64, filterColumns ...string) (map[string]interface{}, errs.SystemErrorInfo) {
//	list, err := GetIssueInfosMapLcByIssueIds(orgId, userId, []int64{issueId}, filterColumns...)
//	if err != nil {
//		return nil, err
//	}
//	if len(list) > 0 {
//		return list[0], nil
//	}
//
//	return nil, errs.IssueNotExist
//}

// GetIssueInfosMapLcByIssueIds 从无码获取多条数据
func GetIssueInfosMapLcByIssueIds(orgId, userId int64, issueIds []int64, filterColumns ...string) ([]map[string]interface{}, errs.SystemErrorInfo) {
	condition := GetRowsCondition(consts.BasicFieldIssueId, tablePb.ConditionType_in, nil, issueIds)
	return GetIssueInfosMapLc(orgId, userId, condition, filterColumns, 1, 2000)
}

// GetIssueInfosMapLcByDataIds 从无码获取多条数据
func GetIssueInfosMapLcByDataIds(orgId, userId int64, dataIds []int64, filterColumns ...string) ([]map[string]interface{}, errs.SystemErrorInfo) {
	condition := GetRowsCondition(consts.BasicFieldId, tablePb.ConditionType_in, nil, dataIds)
	return GetIssueInfosMapLc(orgId, userId, condition, filterColumns, 1, 2000)
}

// GetIssueInfosMapLcByIssueIds 从无码获取多条数据
//func GetIssueInfosMapLcByIssueIdsFilter(orgId, userId, appId, tableId int64, issueIds []int64, filterColumns ...string) ([]map[string]interface{}, errs.SystemErrorInfo) {
//	var issueIdIs []interface{}
//	for _, issueId := range issueIds {
//		issueIdIs = append(issueIdIs, fmt.Sprintf("'%d'", issueId))
//	}
//	noPretreat := true
//	lessReq := vo.LessCondsData{
//		Type:       "in",
//		Values:     issueIdIs,
//		Column:     consts.BasicFieldIssueId,
//		NoPretreat: &noPretreat,
//	}
//	return GetIssueInfosMapLcFilter(orgId, userId, appId, tableId, lessReq, filterColumns, 1, 2000)
//}

// GetIssueInfoMapLcByDataId 从无码获取单条数据
//func GetIssueInfoMapLcByDataId(orgId, userId int64, dataId int64, filterColumns ...string) (map[string]interface{}, errs.SystemErrorInfo) {
//	list, err := GetIssueInfosMapLcByDataIds(orgId, userId, []int64{dataId}, filterColumns...)
//	if err != nil {
//		return nil, err
//	}
//	if len(list) > 0 {
//		return list[0], nil
//	}
//
//	return nil, errs.IssueNotExist
//}

func filterColumnUnique(s []string, noExtraIdFlag ...bool) []string {
	if s == nil {
		return nil
	}
	res := make([]string, 0)
	exist := make(map[string]bool)
	for _, s2 := range s {
		if _, ok := exist[s2]; ok {
			continue
		}
		res = append(res, s2)
		exist[s2] = true
	}
	if len(noExtraIdFlag) == 0 {
		if _, ok := exist[consts.BasicFieldId]; !ok {
			res = append(res, lc_helper.ConvertToFilterColumn(consts.BasicFieldId))
		}
		if _, ok := exist[consts.BasicFieldIssueId]; !ok {
			res = append(res, lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId))
		}
	}
	return res
}

func GetTableIssueMaxOrder(orgId, userId int64) (float64, errs.SystemErrorInfo) {
	lessResp, err := GetRawRows(orgId, userId, &tablePb.ListRawRequest{
		FilterColumns: []string{lc_helper.ConvertToFilterColumn(consts.BasicFieldOrder)},
		Orders: []*tablePb.Order{
			{
				Asc:    false,
				Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrder),
			},
		},
		Page: 1,
		Size: 1,
	})
	if err != nil {
		log.Errorf("[GetTableIssueLastOrder] LessIssueRawList failed, orgId: %d,  err: %v", orgId, err.Error())
		return 0, err
	}
	if len(lessResp.Data) == 0 {
		return 0, nil
	}

	return cast.ToFloat64(lessResp.Data[0][consts.BasicFieldOrder]), nil
}

func GetIssueInfosMapLc(orgId, userId int64, condition *tablePb.Condition, filterColumns []string, page, size int64, noExtraIdFlag ...bool) ([]map[string]interface{}, errs.SystemErrorInfo) {
	filterColumns = filterColumnUnique(filterColumns, noExtraIdFlag...)
	req := &tablePb.ListRawRequest{
		FilterColumns: filterColumns,
		Condition:     condition,
		Page:          int32(page),
		Size:          int32(size),
	}

	reply, err := GetRawRows(orgId, userId, req)
	if err != nil {
		return nil, err
	}

	return reply.Data, nil
}

//func GetIssueInfosMapLcFilter(orgId, userId, appId, tableId int64, condition vo.LessCondsData, filterColumns []string, page, size int64) ([]map[string]interface{}, errs.SystemErrorInfo) {
//	filterColumns = filterColumnUnique(filterColumns)
//	lessResp := formfacade.LessIssueList(formvo.LessIssueListReq{
//		Condition:     condition,
//		OrgId:         orgId,
//		AppId:         appId,
//		TableId:       tableId,
//		UserId:        userId,
//		FilterColumns: filterColumns,
//		Page:          page,
//		Size:          size,
//	})
//	if lessResp.Failure() {
//		log.Errorf("[GetIssueInfosMapLc] LessIssueList failed, orgId: %d, condition: %s, err: %v", orgId,
//			json.ToJsonIgnoreError(condition), lessResp.Error())
//		return nil, lessResp.Error()
//	}
//
//	return lessResp.Data.List, nil
//}

func ConvertIssueDataToIssueBo(data map[string]interface{}) (*bo.IssueBo, errs.SystemErrorInfo) {
	issueBo := &bo.IssueBo{}
	err := copyer.Copy(data, issueBo)
	if err != nil {
		log.Errorf("[ConvertIssueDataToIssueBo] json转换错误, err:%v", err)
		return nil, errs.JSONConvertError
	}
	issueBo.Id = cast.ToInt64(data[consts.BasicFieldIssueId])
	//issueBo.IssueId = cast.ToInt64(data[consts.BasicFieldIssueId])
	issueBo.DataId = cast.ToInt64(data[consts.BasicFieldId])
	issueBo.TableId = cast.ToInt64(data[consts.BasicFieldTableId])
	issueBo.Status = cast.ToInt64(data[consts.BasicFieldIssueStatus])
	issueBo.AppId = cast.ToInt64(data[consts.BasicFieldAppId])
	if len(issueBo.Path) == 0 {
		issueBo.Path = "0,"
	}
	issueBo.OwnerIdI64 = businees.LcMemberToUserIds(issueBo.OwnerId)
	issueBo.AuditorIdsI64 = businees.LcMemberToUserIds(issueBo.AuditorIds)
	issueBo.FollowerIdsI64 = businees.LcMemberToUserIds(issueBo.FollowerIds)
	if issueBo.AuditStatusDetail == nil {
		issueBo.AuditStatusDetail = make(map[string]int)
	}
	issueBo.LessData = data
	return issueBo, nil
}

func GetIssueInfoLc(orgId, userId int64, issueId int64, filterColumns ...string) (*bo.IssueBo, errs.SystemErrorInfo) {
	issueLcDatas, errSys := GetIssueInfosMapLcByIssueIds(orgId, userId, []int64{issueId}, filterColumns...)
	if errSys != nil {
		log.Errorf("[GetIssueInfosLc] GetIssueInfosMapLcByIssueIds err: %v, issueId: %v", errSys, issueId)
		return nil, errSys
	}

	if len(issueLcDatas) < 1 {
		log.Errorf("[GetIssueInfosLc] GetIssueInfosMapLcByIssueIds err: %v, issueId: %v", errSys, issueId)
		return nil, errs.IssueNotExist
	}

	var issueBo *bo.IssueBo
	for _, data := range issueLcDatas {
		issueBo, errSys = ConvertIssueDataToIssueBo(data)
		if errSys != nil {
			log.Errorf("[GetIssueInfosLc] ConvertIssueDataToIssueBo err: %v, issueId: %v", errSys, issueId)
			return nil, errSys
		}
		break
	}
	return issueBo, nil
}

func GetIssueInfosLc(orgId, userId int64, issueIds []int64, filterColumns ...string) ([]*bo.IssueBo, errs.SystemErrorInfo) {
	issueLcDatas, errSys := GetIssueInfosMapLcByIssueIds(orgId, userId, issueIds, filterColumns...)
	if errSys != nil {
		log.Errorf("[GetIssueInfosLc] GetIssueInfosMapLcByIssueIds err: %v, issueIds: %v", errSys,
			json.ToJsonIgnoreError(issueIds))
		return nil, errSys
	}

	var issueBos []*bo.IssueBo
	for _, data := range issueLcDatas {
		issueBo, errSys := ConvertIssueDataToIssueBo(data)
		if errSys != nil {
			log.Errorf("[GetIssueInfosLc] ConvertIssueDataToIssueBo err: %v, issueIds: %s", errSys,
				json.ToJsonIgnoreError(issueIds))
			return nil, errSys
		}
		issueBos = append(issueBos, issueBo)
	}
	return issueBos, nil
}

func GetIssueInfosLcByDataIds(orgId, userId int64, dataIds []int64) ([]*bo.IssueBo, errs.SystemErrorInfo) {
	issueLcDatas, errSys := GetIssueInfosMapLcByDataIds(orgId, userId, dataIds)
	if errSys != nil {
		return nil, errSys
	}

	var issueBos []*bo.IssueBo
	for _, data := range issueLcDatas {
		issueBo, errSys := ConvertIssueDataToIssueBo(data)
		if errSys != nil {
			return nil, errSys
		}
		issueBos = append(issueBos, issueBo)
	}
	return issueBos, nil
}

func GetMoveRelateColumnIds() []string {
	return []string{
		lc_helper.ConvertToFilterColumn(consts.BasicFieldOrgId),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldAppId),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectId),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldParentId),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldPath),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldOrder),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle),
	}
}

func ReportMoveEvent(userId int64, oldIssue *bo.IssueBo, newData map[string]interface{}, userDepts map[string]*uservo.MemberDept) {
	var columnIds []string
	dataId := cast.ToInt64(newData[consts.BasicFieldId])
	issueId := cast.ToInt64(newData[consts.BasicFieldIssueId])
	orgId := cast.ToInt64(newData[consts.BasicFieldOrgId])
	appId := cast.ToInt64(newData[consts.BasicFieldAppId])
	projectId := cast.ToInt64(newData[consts.BasicFieldProjectId])
	tableId := cast.ToInt64(newData[consts.BasicFieldTableId])
	parentId := cast.ToInt64(newData[consts.BasicFieldParentId])
	path := cast.ToString(newData[consts.BasicFieldPath])
	order := cast.ToInt64(newData[consts.BasicFieldOrder])
	if oldIssue.ProjectId != projectId || oldIssue.AppId != appId {
		columnIds = append(columnIds, consts.BasicFieldProjectId)
		columnIds = append(columnIds, consts.BasicFieldAppId)
	}
	if oldIssue.TableId != tableId {
		columnIds = append(columnIds, consts.BasicFieldTableId)
	}
	if oldIssue.ParentId != parentId {
		columnIds = append(columnIds, consts.BasicFieldParentId)
	}
	if oldIssue.Path != path {
		columnIds = append(columnIds, consts.BasicFieldPath)
	}
	if oldIssue.Order != order {
		columnIds = append(columnIds, consts.BasicFieldOrder)
	}
	if len(columnIds) == 0 {
		return
	}

	oldData := oldIssue.LessData

	// 转成前端识别的模式。。
	AssembleDataIds(oldData)
	AssembleDataIds(newData)

	e := &commonvo.DataEvent{
		OrgId:          orgId,
		AppId:          appId,
		ProjectId:      projectId,
		TableId:        tableId,
		DataId:         dataId,
		IssueId:        issueId,
		UserId:         userId,
		Old:            oldData,
		New:            newData,
		UpdatedColumns: columnIds,
		UserDepts:      userDepts,
	}
	openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
	openTraceIdStr := cast.ToString(openTraceId)
	report.ReportDataEvent(msgPb.EventType_DataMoved, openTraceIdStr, e)
}

//func ArchiveIssue(orgId int64, issueBos []bo.IssueBo, operatorId int64, sourceChannel string) (map[int64][]int64, errs.SystemErrorInfo) {
//	allIssueId := []int64{}
//	originIds := []int64{}
//	for _, issueBo := range issueBos {
//		allIssueId = append(allIssueId, issueBo.Id)
//		originIds = append(originIds, issueBo.Id)
//	}
//
//	issueAndChildren, issueErr := GetIssueAndChildren(orgId, allIssueId, nil)
//	if issueErr != nil {
//		log.Error(issueErr)
//		return nil, issueErr
//	}
//	childIssueIds := map[int64][]int64{}
//	for _, issue := range issueAndChildren {
//		allIssueId = append(allIssueId, issue.Id)
//		for _, id := range originIds {
//			if strings.Contains(issue.Path, fmt.Sprintf(",%d,", id)) {
//				childIssueIds[id] = append(childIssueIds[id], issue.Id)
//			}
//		}
//	}
//	_, err := mysql.UpdateSmartWithCond(consts.TableIssue, db.Cond{
//		consts.TcId:       db.In(allIssueId),
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}, mysql.Upd{
//		consts.TcIsFiling: consts.ProjectIsFiling,
//		consts.TcUpdator:  operatorId,
//	})
//	if err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//
//	asyn.Execute(func() {
//		for _, issueBo := range issueBos {
//			ownerIds, err3 := GetIssueRelationIdsByRelateType(orgId, issueBo.Id, consts.IssueRelationTypeOwner)
//			if err3 != nil {
//				return
//			}
//			issueTrendsBo := &bo.IssueTrendsBo{
//				PushType:      consts.PushTypeArchiveIssue,
//				OrgId:         orgId,
//				OperatorId:    operatorId,
//				IssueId:       issueBo.Id,
//				ParentIssueId: issueBo.ParentId,
//				ProjectId:     issueBo.ProjectId,
//				PriorityId:    issueBo.PriorityId,
//				ParentId:      issueBo.ParentId,
//
//				IssueTitle:    issueBo.Title,
//				IssueStatusId: issueBo.Status,
//				BeforeOwner:   *ownerIds,
//				AfterOwner:    []int64{0},
//				SourceChannel: sourceChannel,
//				TableId:       issueBo.TableId,
//			}
//			if _, ok := childIssueIds[issueBo.Id]; ok {
//				issueTrendsBo.IssueChildren = childIssueIds[issueBo.Id]
//			}
//
//			asyn.Execute(func() {
//				PushIssueTrends(issueTrendsBo)
//			})
//			asyn.Execute(func() {
//				PushIssueThirdPlatformNotice(issueTrendsBo)
//			})
//		}
//	})
//
//	return childIssueIds, nil
//}

//func GetIssueIdsByOrgId(orgId int64, page, size int) ([]int64, int64, errs.SystemErrorInfo) {
//	list := make([]int64, 0)
//	cond1 := db.Cond{
//		consts.TcOrgId: orgId,
//	}
//	bos, total, err1 := SelectList(cond1, nil, page, size, "id desc", true)
//	for _, bo := range *bos {
//		list = append(list, bo.Id)
//	}
//	if err1 != nil {
//		log.Error(err1)
//		return list, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
//	}
//	return list, total, nil
//}

//// CheckCreateSourceIsCopy 检查任务创建来源是否是复制
//func CheckCreateSourceIsCopy(lessCodeData map[string]interface{}) bool {
//	if len(lessCodeData) < 1 {
//		return false
//	}
//	if val, ok := lessCodeData[consts.HelperFieldCreateFrom]; ok {
//		return val == consts.CreateIssueSourceCopy
//	}
//
//	return false
//}

//func GetSomeInfoForCreateIssue(orgId, projectId, tableId int64) (*projectvo.GetSomeInfoForCreateIssueBo, errs.SystemErrorInfo) {
//	info := &projectvo.GetSomeInfoForCreateIssueBo{
//		ProjectId: projectId,
//		TableId:   tableId,
//	}
//
//	// 获取项目
//	if projectId > 0 {
//		project, errSys := GetProjectSimple(orgId, projectId)
//		if errSys != nil {
//			log.Error(errSys)
//			return nil, errSys
//		}
//		info.Project = project
//	}
//
//	// 获取表头
//	if tableId == 0 {
//		summeryTableResp := tablefacade.GetSummeryTableId(projectvo.GetSummeryTableIdReqVo{
//			OrgId:  orgId,
//			UserId: 0,
//			Input:  &tablePb.ReadSummeryTableIdRequest{},
//		})
//		if summeryTableResp.Failure() {
//			log.Errorf("[GetSomeInfoForCreateIssue] failed, orgId: %d, tableId: %d, err: %v", orgId, tableId, summeryTableResp.Error())
//			return nil, summeryTableResp.Error()
//		}
//		tableId = summeryTableResp.NewData.TableId
//	}
//	columnsData, errSys := GetTablesColumnsByTableIds(orgId, 0, []int64{tableId}, nil, true)
//	if errSys != nil {
//		log.Errorf("[GetSomeInfoForCreateIssue] GetTableColumns failed:%v, orgId:%d, tableId:%d", errSys, orgId, tableId)
//		return nil, errSys
//	}
//	if len(columnsData.Tables) < 1 {
//		log.Errorf("[GetSomeInfoForCreateIssue] GetTableColumns not found err: %v, orgId:%d, tableId:%d", errSys, orgId, tableId)
//		return nil, errs.TableNotExist
//	}
//	columns := columnsData.Tables[0].Columns
//	lcCommonFields := make([]lc_table.LcCommonField, 0, len(columns))
//	err := copyer.Copy(columns, &lcCommonFields)
//	if err != nil {
//		log.Errorf("[GetSomeInfoForCreateIssue] json转换错误: %v", err)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
//	}
//	headers := map[string]lc_table.LcCommonField{}
//	for _, col := range lcCommonFields {
//		headers[col.Name] = col
//	}
//	info.ColumnMap = headers
//
//	return info, nil
//}

// DeleteIssueInLcByProject 根据项目 id 删除存于无码的任务数据（软删除）
func DeleteIssueInLcByProject(orgId, opUserId int64, projectId int64, templateFlag int) errs.SystemErrorInfo {
	reply := tablefacade.DeleteRows(orgId, opUserId, &tablePb.DeleteRequest{Condition: &tablePb.Condition{Type: tablePb.ConditionType_and, Conditions: []*tablePb.Condition{
		GetRowsCondition(consts.BasicFieldOrgId, tablePb.ConditionType_equal, orgId, nil),
		GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_equal, projectId, nil),
	}}})
	if reply.Failure() {
		return reply.Error()
	}

	if templateFlag == consts.TemplateFalse {
		errCache := SetDeletedIssuesNum(orgId, reply.Data.Count)
		if errCache != nil {
			log.Errorf("[DeleteIssueInLcByProject] SetDeletedIssuesNum err:%v, orgId:%v", errCache, orgId)
		}
	}

	return nil
}

//// 获取复制
//func getCopyRelatingForm(orgId, userId int64, addDataId int64, relatingValue interface{}, relatingType string) (map[int64][]map[string]interface{}, errs.SystemErrorInfo) {
//	allDataIds := []int64{}
//	linkToDataIds := []int64{}
//	linkFromDataIds := []int64{}
//	oldRelating := bo.RelatingIssue{}
//	copyer.Copy(relatingValue, &oldRelating)
//	if len(oldRelating.LinkTo) > 0 {
//		linkToDataIds = convert.ToInt64Slice(oldRelating.LinkTo)
//	}
//	if len(oldRelating.LinkFrom) > 0 {
//		linkFromDataIds = convert.ToInt64Slice(oldRelating.LinkFrom)
//	}
//	allDataIds = append(allDataIds, linkToDataIds...)
//	allDataIds = append(allDataIds, linkFromDataIds...)
//	allDataIds = slice.SliceUniqueInt64(allDataIds)
//	lcIssues, errSys := GetIssueInfosLcByDataIds(orgId, userId, allDataIds)
//	if errSys != nil {
//		log.Errorf("[getCopyRelatingForm]GetIssueInfosLcByDataIds err:%v", errSys)
//		return nil, errSys
//	}
//	updateForm := []map[string]interface{}{}
//	appIdUpdateMap := map[int64][]map[string]interface{}{}
//	// 根据appId不同分组，构建批量更新条件
//	for _, iss := range lcIssues {
//		if int642.InArray(iss.DataId, linkToDataIds) {
//			// 原来记录的linkTo中dataIds对应记录的linkFrom需要 添加该复制过来记录的dataId
//			relatingIssue := &bo.RelatingIssue{}
//			if relatingType == consts.BasicFieldRelating {
//				relatingIssue = &iss.RelatingIssue
//			} else {
//				relatingIssue = &iss.BaRelatingIssue
//			}
//			if relatingIssue != nil {
//				if relatingIssue.LinkTo == nil {
//					relatingIssue.LinkTo = make([]string, 0)
//				}
//				if relatingIssue.LinkFrom == nil {
//					relatingIssue.LinkFrom = make([]string, 0)
//				}
//				updateLinkFrom := relatingIssue.LinkFrom
//				updateLinkFrom = append(updateLinkFrom, fmt.Sprintf("%d", addDataId))
//				updateLinkTo := slice.SliceUniqueString(relatingIssue.LinkTo)
//				updateLinkFrom = slice.SliceUniqueString(updateLinkFrom)
//				updateForm = append(updateForm, map[string]interface{}{
//					consts.BasicFieldIssueId: iss.Id,
//					relatingType: bo.RelatingIssue{
//						LinkTo:   updateLinkTo,
//						LinkFrom: updateLinkFrom,
//					},
//				})
//			}
//		}
//		if int642.InArray(iss.DataId, linkFromDataIds) {
//			relatingIssue := &bo.RelatingIssue{}
//			if relatingType == consts.BasicFieldRelating {
//				relatingIssue = &iss.RelatingIssue
//			} else {
//				relatingIssue = &iss.BaRelatingIssue
//			}
//			if relatingIssue != nil {
//				if relatingIssue.LinkTo == nil {
//					relatingIssue.LinkTo = make([]string, 0)
//				}
//				if relatingIssue.LinkFrom == nil {
//					relatingIssue.LinkFrom = make([]string, 0)
//				}
//				updatedLinkTo := relatingIssue.LinkTo
//				updatedLinkTo = append(updatedLinkTo, fmt.Sprintf("%d", addDataId))
//				updatedLinkTo = slice.SliceUniqueString(updatedLinkTo)
//				updatedLinkFrom := slice.SliceUniqueString(relatingIssue.LinkFrom)
//				updateForm = append(updateForm, map[string]interface{}{
//					consts.BasicFieldIssueId: iss.Id,
//					relatingType: bo.RelatingIssue{
//						LinkTo:   updatedLinkTo,
//						LinkFrom: updatedLinkFrom,
//					},
//				})
//			}
//		}
//		if _, ok := appIdUpdateMap[iss.AppId]; ok {
//			appIdUpdateMap[iss.AppId] = append(appIdUpdateMap[iss.AppId], updateForm...)
//		} else {
//			appIdUpdateMap[iss.AppId] = updateForm
//		}
//	}
//
//	return appIdUpdateMap, nil
//}

// 没有移动过去，需要将linkTo中对应记录的linkFrom删除该记录的dataId，linkFrom中dataIds对应的所有记录的linkTo也需要删除对应的dataId
func getNoMoveRelatingForm(orgId, userId int64, issueBo *bo.IssueBo, relating *bo.RelatingIssue, relatingType string) (map[int64][]map[string]interface{}, errs.SystemErrorInfo) {
	if relating == nil {
		return nil, nil
	}
	allIssueIds := []int64{}
	linkToDataIds := []int64{}
	linkFromDataIds := []int64{}
	if relating.LinkTo != nil && len(relating.LinkTo) > 0 {
		linkToDataIds = convert.ToInt64Slice(relating.LinkTo)
	}
	if relating.LinkFrom != nil && len(relating.LinkFrom) > 0 {
		linkFromDataIds = convert.ToInt64Slice(relating.LinkFrom)
	}
	allIssueIds = append(allIssueIds, linkToDataIds...)
	allIssueIds = append(allIssueIds, linkFromDataIds...)
	allIssueIds = slice.SliceUniqueInt64(allIssueIds)
	if len(allIssueIds) == 0 {
		return nil, nil
	}
	lcIssues, errSys := GetIssueInfosLc(orgId, userId, allIssueIds)
	if errSys != nil {
		log.Errorf("[getUpdateRelating] err:%v", errSys)
		return nil, errSys
	}
	log.Infof("[getNoMoveRelatingForm] ids: %v, issues: %v", json.ToJsonIgnoreError(allIssueIds), json.ToJsonIgnoreError(lcIssues))
	updateForm := []map[string]interface{}{}
	appIdUpdateForm := map[int64][]map[string]interface{}{}
	for _, iss := range lcIssues {
		if iss.Id == 0 {
			continue
		}
		if int642.InArray(iss.Id, linkToDataIds) {
			// 该记录的linkTo中dataIds对应记录的linkFrom需要删除 该记录的dataId
			relatingIssue := &bo.RelatingIssue{}
			if relatingType == consts.BasicFieldRelating {
				relatingIssue = &iss.RelatingIssue
			} else {
				relatingIssue = &iss.BaRelatingIssue
			}
			if relatingIssue != nil {
				if relatingIssue.LinkTo == nil {
					relatingIssue.LinkTo = make([]string, 0)
				}
				if relatingIssue.LinkFrom == nil {
					relatingIssue.LinkFrom = make([]string, 0)
				}
				updatedLinkFrom := str.DeleteSliceElement(relatingIssue.LinkFrom, fmt.Sprintf("%d", issueBo.Id))
				updatedLinkFrom = slice.SliceUniqueString(updatedLinkFrom)
				updatedLinkTo := slice.SliceUniqueString(relatingIssue.LinkTo)
				updateForm = append(updateForm, map[string]interface{}{
					consts.BasicFieldIssueId: iss.Id,
					relatingType: bo.RelatingIssue{
						LinkTo:   updatedLinkTo,
						LinkFrom: updatedLinkFrom,
					},
				})
			}
		} else if int642.InArray(iss.Id, linkFromDataIds) {
			relatingIssue := &bo.RelatingIssue{}
			if relatingType == consts.BasicFieldRelating {
				relatingIssue = &iss.RelatingIssue
			} else {
				relatingIssue = &iss.BaRelatingIssue
			}
			if relatingIssue != nil {
				if relatingIssue.LinkTo == nil {
					relatingIssue.LinkTo = make([]string, 0)
				}
				if relatingIssue.LinkFrom == nil {
					relatingIssue.LinkFrom = make([]string, 0)
				}
				updatedLinkTo := str.DeleteSliceElement(relatingIssue.LinkTo, fmt.Sprintf("%d", issueBo.Id))
				updatedLinkTo = slice.SliceUniqueString(updatedLinkTo)
				updatedLinkFrom := slice.SliceUniqueString(relatingIssue.LinkFrom)
				updateForm = append(updateForm, map[string]interface{}{
					consts.BasicFieldIssueId: iss.Id,
					relatingType: bo.RelatingIssue{
						LinkTo:   updatedLinkTo,
						LinkFrom: updatedLinkFrom,
					},
				})
			}
		}
		if _, ok := appIdUpdateForm[iss.AppId]; ok {
			appIdUpdateForm[iss.AppId] = append(appIdUpdateForm[iss.AppId], updateForm...)
		} else {
			appIdUpdateForm[iss.AppId] = updateForm
		}
	}

	return appIdUpdateForm, nil
}
