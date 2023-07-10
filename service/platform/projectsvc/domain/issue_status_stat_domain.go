package domain

import (
	"fmt"
	"time"

	tablePb "github.com/star-table/interface/golang/table/v1"

	"github.com/star-table/common/core/util/json"

	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/date"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
)

// GetIssueStatusStatWithLc 查询和计算迭代下的任务统计信息
// 主要统计各个 table 的各种任务数
func GetIssueStatusStatWithLc(issueStatusStatCondBo bo.IssueStatusStatCondBo) ([]bo.IssueStatusStatBo, errs.SystemErrorInfo) {
	orgId := issueStatusStatCondBo.OrgId
	tmpParam := bo.IssueBoListCond{
		OrgId:        orgId,
		ProjectId:    issueStatusStatCondBo.ProjectId,
		IterationId:  issueStatusStatCondBo.IterationId,
		RelationType: issueStatusStatCondBo.RelationType,
		UserId:       &issueStatusStatCondBo.UserId,
		ProAppId:     nil,
	}
	if issueStatusStatCondBo.ProjectId == nil && orgId < 1 {
		log.Error("[GetIssueStatusStatWithLc] projectId/orgId 为必传参数。")
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.InputParamEmpty, "orgId 和 projectId 必传其一。")
	}
	if issueStatusStatCondBo.ProjectId != nil && *issueStatusStatCondBo.ProjectId > 0 {
		proAppId, err := GetAppIdFromProjectId(issueStatusStatCondBo.OrgId, *issueStatusStatCondBo.ProjectId)
		if err != nil {
			log.Errorf("[GetIssueStatusStatWithLc] GetAppIdFromProjectId err: %v", err)
			return nil, err
		}
		tmpParam.ProAppId = &proAppId
	} else {
		summaryAppId, err := GetOrgSummaryAppId(orgId)
		if err != nil {
			log.Errorf("[GetIssueStatusStatWithLc] GetAppIdFromProjectId err: %v", err)
			return nil, err
		}
		tmpParam.ProAppId = &summaryAppId
	}

	if issueStatusStatCondBo.UserId < 1 {
		tmpParam.UserId = nil
	}

	return getIssueStatusStat(tmpParam)
}

// GetIssueAndDetailUnionBoListFromLc 从无码接口查询任务数据
func getIssueStatusStat(issueListCond bo.IssueBoListCond) ([]bo.IssueStatusStatBo, errs.SystemErrorInfo) {
	condition := &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: GetNoRecycleCondition(),
	}

	if issueListCond.ProjectId != nil {
		condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_equal, *issueListCond.ProjectId, nil))
	}
	if issueListCond.RelationType != nil && issueListCond.UserId != nil {
		tmpCondArr := AssemblyIssueRelationCondForLcQueryIssue(map[int]interface{}{
			*issueListCond.RelationType: *issueListCond.UserId,
		})
		condition.Conditions = append(condition.Conditions, tmpCondArr...)
	}
	if issueListCond.IterationId != nil {
		condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldIterationId, tablePb.ConditionType_equal, *issueListCond.IterationId, nil))
	}

	todayZero := date.GetZeroTime(time.Now())
	todayZeroStr := todayZero.Format(consts.AppTimeFormat)
	nowStr := time.Now().Format(consts.AppTimeFormat)
	tomorrowZero := todayZero.AddDate(0, 0, 1)
	tomorrowZeroStr := tomorrowZero.Format(consts.AppTimeFormat)
	theDayAfterTomorrowZeroStr := tomorrowZero.AddDate(0, 0, 1).Format(consts.AppTimeFormat)
	filterColumns := []string{
		`count(1) as "issueCount"`, // 总数
		`sum(case when data->'issueStatusType'='1' then 1 else 0 end) as "issueWaitCount"`,
		`sum(case when data->'issueStatusType'='2' then 1 else 0 end) as "issueRunningCount"`,
		`sum(case when data->'issueStatusType'='3' and data->'auditStatus' in ('-1','3') then 1 else 0 end) as "issueEndCount"`,
		`sum(case when data->'issueStatusType'='3' and data->'auditStatus' in ('1','2') then 1 else 0 end) as "waitConfirmedCount"`,
		fmt.Sprintf(`sum(case when data->'issueStatusType'='3' and data->'auditStatus' in ('-1','3') and data->'endTime' is not null and data->'endTime'>='"%s"' and data->'endTime'<'"%s"' then 1 else 0 end) as "issueEndTodayCount"`, todayZeroStr, tomorrowZeroStr),
		fmt.Sprintf(`sum(case when data->'issueStatusType'<>'3' and data->'planEndTime' is not null and data->'planEndTime'>'"1970-01-01 00:00:00"' and data->'planEndTime'<'"%s"' then 1 else 0 end) as "issueOverdueCount"`, nowStr),
		fmt.Sprintf(`sum(case when "createTime" is not null and "createTime">='%s' and "createTime"<'%s' then 1 else 0 end) as "todayCreateCount"`, todayZeroStr, tomorrowZeroStr),
		fmt.Sprintf(`sum(case when data->'ownerChangeTime' is not null and data->'ownerChangeTime'>'"%s"' and data->'ownerChangeTime'<'"%s"' then 1 else 0 end) as "todayCount"`, todayZeroStr, tomorrowZeroStr),
		`sum(case when data->'issueStatusType'='3' and data->'planEndTime' is not null and data->'planEndTime'>'"1970-01-01 00:00:00"' and data->'endTime' is not null and data->'endTime'>'"1970-01-01 00:00:00"' and data->'endTime'>data->'planEndTime' then 1 else 0 end) as "issueOverdueEndCount"`,
		fmt.Sprintf(`sum(case when data->'issueStatusType'='3' and data->'auditStatus' in ('-1','3') and data->'planEndTime' is not null and data->'planEndTime'>='"%s"' and data->'planEndTime'<'"%s"' then 1 else 0 end) as "issueEndTodayCount"`, todayZeroStr, tomorrowZeroStr),
		fmt.Sprintf(`sum(case when data->'issueStatusType'<>'3' and data->'planEndTime' is not null and data->'planEndTime'>='"%s"' and data->'planEndTime'<'"%s"' then 1 else 0 end) as "issueOverdueTomorrowCount"`, todayZeroStr, theDayAfterTomorrowZeroStr),
		fmt.Sprintf(`sum(case when data->'issueStatusType'<>'3' and data->'planEndTime' is not null and data->'planEndTime'>='"%s"' and data->'planEndTime'<'"%s"' then 1 else 0 end) as "issueOverdueTodayCount"`, todayZeroStr, tomorrowZeroStr),
	}

	req := &tablePb.ListRawRequest{
		DbType:        tablePb.DbType_slave1,
		FilterColumns: filterColumns,
		Condition:     condition,
		Groups: []string{
			// group by 用别名也可
			fmt.Sprintf("\"%s\"", consts.BasicFieldProjectId),
		},
		Page: 1,
		Size: 1,
	}
	lessResp, err := GetRawRows(issueListCond.OrgId, 0, req)
	if err != nil {
		log.Errorf("[GetIssueFinishedStatForProjects] orgId:%d, err: %v", issueListCond.OrgId, err.Error())
		return nil, err
	}

	result := make([]bo.IssueStatusStatBo, 0, 1)
	_ = json.FromJson(json.ToJsonIgnoreError(lessResp.Data), &result)
	return result, nil
}

func dealStatForLc(issueAndDetailUnionBos []bo.IssueAndDetailUnionBo, tablesMap maps.LocalMap) ([]bo.IssueStatusStatBo, errs.SystemErrorInfo) {
	extData := map[int64]bo.IssueStatusStatBo{}
	for _, targetBo := range issueAndDetailUnionBos {
		tableObjIf, ok := tablesMap[targetBo.TableId]
		if !ok {
			continue
		}
		tableObj := tableObjIf.(projectvo.TableMetaData)
		data, ok := extData[targetBo.TableId]
		if !ok {
			data = bo.IssueStatusStatBo{
				ProjectTypeId:       targetBo.TableId,
				ProjectTypeName:     tableObj.Name,
				ProjectTypeLangCode: "",
			}
		}

		data.IssueCount++
		data.StoryPointCount++

		todayZero := date.GetZeroTime(time.Now())
		now := time.Now()
		tomorrowZero := todayZero.AddDate(0, 0, 1)
		theDayAfterTomorrowZero := tomorrowZero.AddDate(0, 0, 1)

		if targetBo.IssueStatusType == consts.StatusTypeNotStart {
			data.IssueWaitCount++
			data.StoryPointWaitCount++
		} else if targetBo.IssueStatusType == consts.StatusTypeRunning {
			data.IssueRunningCount++
			data.StoryPointRunningCount++
		} else if targetBo.IssueStatusType == consts.StatusTypeComplete {
			//待确认的(通用项目才有待确认)
			if targetBo.IssueStatusId == int64(26) {
				data.IssueEndCount++
				data.StoryPointEndCount++
				if targetBo.IssueAuditStatus == consts.AuditStatusNotView {
					data.WaitConfirmedCount++
				} else {
					//今日完成
					if targetBo.EndTime.After(todayZero) && targetBo.EndTime.Before(tomorrowZero) {
						data.IssueEndTodayCount++
					}
				}
			} else {
				data.IssueEndCount++
				data.StoryPointEndCount++

				//今日完成
				if targetBo.EndTime.After(todayZero) && targetBo.EndTime.Before(tomorrowZero) {
					data.IssueEndTodayCount++
				}
			}
		}

		if targetBo.CreateTime.After(todayZero) && targetBo.CreateTime.Before(tomorrowZero) {
			//今日创建
			data.TodayCreateCount++
		}

		if targetBo.IssueStatusType != consts.StatusTypeComplete {
			if targetBo.PlanEndTime.After(consts.BlankTimeObject) && targetBo.PlanEndTime.Before(now) {
				//已逾期
				data.IssueOverdueCount++
			}
			if targetBo.PlanEndTime.Equal(now) || (targetBo.PlanEndTime.After(now) && targetBo.PlanEndTime.Before(tomorrowZero)) {
				//今日到期
				data.IssueOverdueTodayCount++
			}
			//明日逾期
			if targetBo.PlanEndTime.Equal(tomorrowZero) || (targetBo.PlanEndTime.After(tomorrowZero) && targetBo.PlanEndTime.Before(theDayAfterTomorrowZero)) {
				//即将逾期
				data.IssueOverdueTomorrowCount++
			}

			if targetBo.OwnerChangeTime.After(todayZero) && targetBo.OwnerChangeTime.Before(tomorrowZero) {
				//今日指派给我
				data.TodayCount++
			}

		}
		//逾期完成
		if targetBo.IssueStatusType == consts.StatusTypeComplete && targetBo.PlanEndTime.Before(targetBo.EndTime) && targetBo.PlanEndTime.After(consts.BlankTimeObject) {
			data.IssueOverdueEndCount++
		}

		extData[targetBo.TableId] = data
	}

	statBos := make([]bo.IssueStatusStatBo, 0)
	for _, v := range extData {
		statBos = append(statBos, v)
	}

	return statBos, nil
}
