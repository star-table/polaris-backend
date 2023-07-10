package domain

import (
	"time"

	"github.com/star-table/common/core/types"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
)

func SelectIssueAssignCount(orgId int64, projectId int64, rankTop int) ([]bo.IssueAssignCountBo, errs.SystemErrorInfo) {
	return nil, nil
	//finishedIds, err2 := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeComplete)
	//if err2 != nil || len(*finishedIds) == 0 {
	//	log.Error(err2)
	//	return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
	//}
	//conn, err := mysql.GetConnect()
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	//}
	//issueAssignCountBos := &[]bo.IssueAssignCountBo{}
	//err1 := conn.Select(db.Raw("count(*) as count"), consts.TcOwner).From(consts.TableIssue).Where(db.Cond{
	//	consts.TcProjectId: projectId,
	//	consts.TcStatus:    db.NotIn(finishedIds),
	//	consts.TcIsDelete:  consts.AppIsNoDelete,
	//}).GroupBy(consts.TcOwner).OrderBy("count desc").Limit(rankTop).All(issueAssignCountBos)
	//
	//if err1 != nil {
	//	log.Error(err1)
	//	return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	//}

	//appId, appIdErr := GetOrgSummaryAppId(orgId)
	//if appIdErr != nil {
	//	log.Errorf("[SelectIssueAssignTotalCount]GetOrgSummaryAppId failed:%v, orgId:%d", appIdErr, orgId)
	//	return nil, appIdErr
	//}
	//
	//lessReqParam := formvo.LessIssueListReq{
	//	Condition: vo.LessCondsData{
	//		Type: "and",
	//		Conds: []*vo.LessCondsData{
	//			{
	//				Type:   "equal",
	//				Column: consts.BasicFieldProjectId,
	//				Value:  projectId,
	//			},
	//			{
	//				Type:   "not_in",
	//				Column: consts.BasicFieldIssueStatusType,
	//				Values: []interface{}{consts.StatusTypeComplete},
	//			},
	//			{
	//				Type:   "equal",
	//				Column: consts.BasicFieldDelFlag,
	//				Value:  consts.AppIsNoDelete,
	//			},
	//			{
	//				Type:   "equal",
	//				Column: consts.BasicFieldRecycleFlag,
	//				Value:  consts.AppIsNoDelete,
	//			},
	//		},
	//	},
	//	Page: 1,
	//	Size: int64(rankTop),
	//	Groups: []string{
	//		"owner",
	//	},
	//	AppId: appId,
	//	OrgId: orgId,
	//	FilterColumns: []string{
	//		"count(*) as count",
	//		fmt.Sprintf("left(right(jsonb_array_elements(data->'%s')::text,-3),-1) \"%s\"", consts.BasicFieldOwnerId, "owner"),
	//		//fmt.Sprintf("regexp_replace(jsonb_array_elements(data->'%s')::text,'[^0-9]+','','g') \"%s\"", consts.BasicFieldOwnerId, "owner"),
	//	},
	//	Orders: []*vo.LessOrder{
	//		{
	//			Asc:    false,
	//			Column: "count",
	//		},
	//	},
	//}
	//lessResp := formfacade.LessFilterCustomStat(lessReqParam)
	//if lessResp.Failure() {
	//	log.Errorf("[SelectIssueAssignCount] orgId:%d, projectId:%d, LessIssueList failure:%v", orgId, projectId, lessResp.Error())
	//	return nil, lessResp.Error()
	//}
	//issueAssignCountBos := &[]bo.IssueAssignCountBo{}
	//err := copyer.Copy(lessResp.Data, issueAssignCountBos)
	//if err != nil {
	//	log.Errorf("[SelectIssueAssignCount] orgId:%d, projectId:%d, Copy failure:%v", orgId, projectId, err)
	//	return nil, errs.ObjectCopyError
	//}
	//return *issueAssignCountBos, nil
}

func SelectIssueAssignTotalCount(orgId int64, projectId int64) (int64, errs.SystemErrorInfo) {
	return 0, nil
	//appId, appIdErr := GetOrgSummaryAppId(orgId)
	//if appIdErr != nil {
	//	log.Errorf("[SelectIssueAssignTotalCount]GetOrgSummaryAppId failed:%v, orgId:%d", appIdErr, orgId)
	//	return 0, appIdErr
	//}
	//
	//lessReqParam := formvo.LessIssueListReq{
	//	Condition: vo.LessCondsData{
	//		Type: "and",
	//		Conds: []*vo.LessCondsData{
	//			{
	//				Type:   "equal",
	//				Column: consts.BasicFieldProjectId,
	//				Value:  projectId,
	//			},
	//			{
	//				Type:   "not_in",
	//				Column: consts.BasicFieldIssueStatusType,
	//				Values: []interface{}{consts.StatusTypeComplete},
	//			},
	//			{
	//				Type:   "equal",
	//				Column: consts.BasicFieldDelFlag,
	//				Value:  consts.AppIsNoDelete,
	//			},
	//			{
	//				Type:   "equal",
	//				Column: consts.BasicFieldRecycleFlag,
	//				Value:  consts.AppIsNoDelete,
	//			},
	//		},
	//	},
	//	AppId: appId,
	//	OrgId: orgId,
	//}
	//lessResp := formfacade.LessFilterStat(lessReqParam)
	//if lessResp.Failure() {
	//	log.Errorf("[SelectIssueAssignTotalCount] orgId:%d, projectId:%d, LessIssueList failure:%v", orgId, projectId, lessResp.Error())
	//	return 0, lessResp.Error()
	//}

	//return lessResp.Data, nil
}

func GetIssueCountByOwnerId(orgId, ownerId int64) (int64, errs.SystemErrorInfo) {
	//condition := &tablePb.Condition{
	//	Type: tablePb.ConditionType_and,
	//	Conditions: GetNoRecycleCondition(
	//		GetRowsCondition(consts.BasicFieldOwnerId, tablePb.ConditionType_in, nil, []interface{}{fmt.Sprintf("U_%d", ownerId)}),
	//		GetRowsCondition(consts.BasicFieldIssueStatusType, tablePb.ConditionType_not_in, nil, []interface{}{consts.StatusTypeComplete}),
	//	),
	//}
	//
	//lessResp, err := GetRawRows(orgId, ownerId, &tablePb.ListRawRequest{
	//	FilterColumns: []string{`count(*) "all"`},
	//	Condition:     condition,
	//	Page:          1,
	//	Size:          1,
	//})
	//if err != nil {
	//	log.Errorf("[GetIssueCountByOwnerId] orgId:%d, userId:%d, LessIssueList failure:%v", orgId, ownerId, err.Error())
	//	return 0, err
	//}
	//
	//if len(lessResp.Data) > 0 {
	//	return cast.ToInt64(lessResp.Data[0]["all"]), nil
	//}

	return 0, nil
}

// IssueDailyPersonalWorkCompletionStatForLc 任务完成统计图，根据endTime做groupBy，统计当前负责人在指定时间段之内的每天任务完成数量
// 通过无码查询统计
func IssueDailyPersonalWorkCompletionStatForLc(orgId, ownerId int64, startDate *types.Time, endDate *types.Time) ([]bo.IssueDailyPersonalWorkCompletionStatBo, errs.SystemErrorInfo) {
	return nil, nil
	//appId, appIdErr := GetOrgSummaryAppId(orgId)
	//if appIdErr != nil {
	//	log.Errorf("[GetIssueCountByOwnerId]GetOrgSummaryAppId failed:%v, orgId:%d", appIdErr, orgId)
	//	return nil, appIdErr
	//}
	////默认查询七天
	//startTime, endTime, dateErr := CalDateRangeCond(startDate, endDate, 7)
	//if dateErr != nil {
	//	log.Error(dateErr)
	//	return nil, dateErr
	//}
	//// 如果and后的日期是到天的，需要加一天
	//condEndTime := endTime.AddDate(0, 0, 1)
	//
	//lessReqParam := formvo.LessIssueListReq{
	//	Condition: vo.LessCondsData{
	//		Type: "and",
	//		Conds: []*vo.LessCondsData{
	//			{
	//				Column: consts.BasicFieldOwnerId,
	//				Type:   "values_in",
	//				Values: []interface{}{fmt.Sprintf("U_%d", ownerId)},
	//			},
	//			{
	//				Column: consts.BasicFieldIssueStatusType, // 已完成的任务
	//				Type:   "in",
	//				Values: []interface{}{consts.StatusTypeComplete},
	//			},
	//			{
	//				Column: consts.ProBasicFieldIsDelete,
	//				Type:   "equal",
	//				Value:  consts.AppIsNoDelete,
	//			},
	//			{
	//				Column: consts.BasicFieldPlanEndTime, //理论上更新的时候会把这个值更新进去的
	//				Type:   "between",
	//				Values: []interface{}{startTime.Format(consts.AppDateFormat), condEndTime.Format(consts.AppDateFormat)},
	//			},
	//		},
	//	},
	//	Groups: []string{
	//		"\"date\"",
	//	},
	//	FilterColumns: []string{
	//		" count(*) as \"all\" ",
	//		" date((\"data\"::jsonb ->> 'planEndTime')::timestamp) as \"date\" ",
	//	},
	//	AppId:  appId,
	//	OrgId:  orgId,
	//	UserId: ownerId,
	//}
	//// lessResp := formfacade.LessIssueList(lessReqParam)
	//lessResp := formfacade.LessFilterCustomStat(lessReqParam)
	//if lessResp.Failure() {
	//	log.Errorf("[GetIssueCountByOwnerId] orgId:%d, userId:%d, LessIssueList failure:%v", orgId, ownerId, lessResp.Error())
	//	return nil, lessResp.Error()
	//}
	//resultBos := []bo.IssueDailyPersonalWorkCompletionStatBo{}
	//dataJsonStr := json.ToJsonIgnoreError(lessResp.Data)
	//json.FromJson(dataJsonStr, &resultBos)
	//
	//resultBoMap := map[string]bo.IssueDailyPersonalWorkCompletionStatBo{}
	//for _, resultBo := range resultBos {
	//	date := resultBo.Date.Format(consts.AppDateFormat)
	//	resultBoMap[date] = resultBo
	//}
	////处理结果集，补填空缺
	//afterDealResultBos := make([]bo.IssueDailyPersonalWorkCompletionStatBo, 0)
	//cursorTime := *startTime
	//cursorTimeDateFormat := cursorTime.Format(consts.AppDateFormat)
	//cursorTime, _ = time.Parse(consts.AppDateFormat, cursorTimeDateFormat)
	//endTimeDateFormat := endTime.Format(consts.AppDateFormat)
	//refEndTime, _ := time.Parse(consts.AppDateFormat, endTimeDateFormat)
	//for {
	//	if cursorTime.After(refEndTime) {
	//		break
	//	}
	//	var targetBo *bo.IssueDailyPersonalWorkCompletionStatBo = nil
	//	cursorTimeDateFormat := cursorTime.Format(consts.AppDateFormat)
	//	if resultBo, ok := resultBoMap[cursorTimeDateFormat]; ok {
	//		targetBo = &resultBo
	//	} else {
	//		targetBo = &bo.IssueDailyPersonalWorkCompletionStatBo{
	//			Count: 0,
	//			Date:  cursorTime,
	//		}
	//	}
	//	afterDealResultBos = append(afterDealResultBos, *targetBo)
	//	cursorTime = cursorTime.AddDate(0, 0, 1)
	//}
	//
	//return afterDealResultBos, nil
}

//计算日期范围条件
//startDate: 开始时间
//endDate: 结束时间
//cond: 条件
//defaultDay: 默认查询多少天前
func CalDateRangeCond(startDate *types.Time, endDate *types.Time, defaultDay int) (*time.Time, *time.Time, errs.SystemErrorInfo) {
	if startDate == nil && endDate == nil {
		currentTime := time.Now()
		sd := types.Time(currentTime.AddDate(0, 0, -(defaultDay - 1)))
		ed := types.Time(currentTime)
		startDate = &sd
		endDate = &ed
	}
	if startDate != nil && endDate != nil {
	} else {
		log.Error("没有提供明确的时间范围")
		return nil, nil, errs.BuildSystemErrorInfo(errs.DateRangeError)
	}
	startTime := time.Time(*startDate)
	endTime := time.Time(*endDate)

	return &startTime, &endTime, nil
}

//func GetIssueCountByStatus(orgId int64, projectId int64, statusId int64) (uint64, errs.SystemErrorInfo) {
//	count, err := mysql.SelectCountByCond(consts.TableIssue, db.Cond{
//		consts.TcIsDelete:  consts.AppIsNoDelete,
//		consts.TcOrgId:     orgId,
//		consts.TcProjectId: projectId,
//		consts.TcStatus:    statusId,
//	})
//	if err != nil {
//		log.Error(err)
//		return 0, errs.MysqlOperateError
//	}
//
//	return count, nil
//}
