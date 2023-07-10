package domain

import (
	"time"

	"github.com/star-table/polaris-backend/common/model/bo/status"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func StatisticIterationCountGroupByStatus(orgId, projectId int64) (*bo.IterationStatusTypeCountBo, errs.SystemErrorInfo) {
	var notStartTotal, processingTotal, finishedTotal int64 = 0, 0, 0
	iterationQueryCond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if projectId > 0 {
		iterationQueryCond[consts.TcProjectId] = projectId
	}
	iterationPos, err := dao.SelectIteration(iterationQueryCond)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	cacheStatusBos := consts.IterationStatusList
	statusCacheMap := maps.NewMap("Id", cacheStatusBos)

	for _, iterationPo := range *iterationPos {
		// 通过迭代的状态值，匹配处状态对象
		value, ok := statusCacheMap[iterationPo.Status]
		if !ok {
			continue
		}

		statusBo := value.(status.StatusInfoBo)
		switch statusBo.Type {
		case consts.StatusTypeNotStart:
			notStartTotal++
		case consts.StatusTypeRunning:
			processingTotal++
		case consts.StatusTypeComplete:
			finishedTotal++
		}
	}

	return &bo.IterationStatusTypeCountBo{
		NotStartTotal:   notStartTotal,
		ProcessingTotal: processingTotal,
		FinishedTotal:   finishedTotal,
	}, nil
}

//date: yyyy-MM-dd
func AppendIterationStat(iterationBo bo.IterationBo, date string, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	statDate, timeParseError := time.Parse(consts.AppDateFormat, date)
	if timeParseError != nil {
		log.Error(timeParseError)
		return errs.BuildSystemErrorInfo(errs.SystemError)
	}

	orgId := iterationBo.OrgId
	projectId := iterationBo.ProjectId
	iterationId := iterationBo.Id

	id, err3 := idfacade.ApplyPrimaryIdRelaxed(consts.TableIterationStat)
	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.ApplyIdError, err3)
	}

	projectCacheBo, err := LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	// 需要一个无码版的任务统计
	issueStatusStatBos, err1 := GetIssueStatusStatWithLc(bo.IssueStatusStatCondBo{
		OrgId:       orgId,
		ProjectId:   &projectId,
		IterationId: &iterationId,
	})
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	iterationStatPo := &po.PpmStaIterationStat{}
	iterationStatPo.Id = id
	iterationStatPo.OrgId = orgId
	iterationStatPo.ProjectId = projectId
	iterationStatPo.IterationId = iterationId
	iterationStatPo.StatDate = statDate
	iterationStatPo.Status = projectCacheBo.Status

	issueStatusStatMap := maps.NewMap("ProjectTypeId", issueStatusStatBos)
	ext := bo.StatExtBo{}
	ext.Issue = bo.StatIssueExtBo{
		Data: issueStatusStatMap,
	}
	iterationStatPo.Ext = json.ToJsonIgnoreError(ext)
	//封装状态
	assemblyIterationStat(issueStatusStatBos, iterationStatPo)

	err5 := dao.InsertIterationStat(*iterationStatPo, tx...)
	if err5 != nil {
		log.Error(err5)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

// AppendIterationStatWithLc 查询无码的数据，进而计算新的迭代状态并保存
func AppendIterationStatWithLc(iterationBo bo.IterationBo, date string, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	statDate, timeParseError := time.Parse(consts.AppDateFormat, date)
	if timeParseError != nil {
		log.Error(timeParseError)
		return errs.BuildSystemErrorInfo(errs.SystemError)
	}

	orgId := iterationBo.OrgId
	projectId := iterationBo.ProjectId
	iterationId := iterationBo.Id

	id, err3 := idfacade.ApplyPrimaryIdRelaxed(consts.TableIterationStat)
	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.ApplyIdError, err3)
	}

	projectCacheBo, err := LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	// 需要一个无码版的任务统计
	issueStatusStatBos, err1 := GetIssueStatusStatWithLc(bo.IssueStatusStatCondBo{
		OrgId:       orgId,
		ProjectId:   &projectId,
		IterationId: &iterationId,
	})
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	iterationStatPo := &po.PpmStaIterationStat{}
	iterationStatPo.Id = id
	iterationStatPo.OrgId = orgId
	iterationStatPo.ProjectId = projectId
	iterationStatPo.IterationId = iterationId
	iterationStatPo.StatDate = statDate
	iterationStatPo.Status = projectCacheBo.Status

	issueStatusStatMap := maps.NewMap("ProjectTypeId", issueStatusStatBos)
	ext := bo.StatExtBo{}
	ext.Issue = bo.StatIssueExtBo{
		Data: issueStatusStatMap,
	}
	iterationStatPo.Ext = json.ToJsonIgnoreError(ext)
	//封装状态
	assemblyIterationStat(issueStatusStatBos, iterationStatPo)

	err5 := dao.InsertIterationStat(*iterationStatPo, tx...)
	if err5 != nil {
		log.Error(err5)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func assemblyIterationStat(issueStatusStatBos []bo.IssueStatusStatBo, iterationStatPo *po.PpmStaIterationStat) {
	for _, statBo := range issueStatusStatBos {
		switch statBo.ProjectTypeLangCode {
		case consts.ProjectObjectTypeLangCodeDemand:
			iterationStatPo.DemandCount += statBo.IssueCount
			iterationStatPo.DemandWaitCount += statBo.IssueWaitCount
			iterationStatPo.DemandRunningCount += statBo.IssueRunningCount
			iterationStatPo.DemandEndCount += statBo.IssueEndCount
			iterationStatPo.DemandOverdueCount += statBo.IssueOverdueCount
		case consts.ProjectObjectTypeLangCodeTask:
			iterationStatPo.TaskCount += statBo.IssueCount
			iterationStatPo.TaskWaitCount += statBo.IssueWaitCount
			iterationStatPo.TaskRunningCount += statBo.IssueRunningCount
			iterationStatPo.TaskEndCount += statBo.IssueEndCount
			iterationStatPo.TaskOverdueCount += statBo.IssueOverdueCount
		case consts.ProjectObjectTypeLangCodeBug:
			iterationStatPo.BugCount += statBo.IssueCount
			iterationStatPo.BugWaitCount += statBo.IssueWaitCount
			iterationStatPo.BugRunningCount += statBo.IssueRunningCount
			iterationStatPo.BugEndCount += statBo.IssueEndCount
			iterationStatPo.BugOverdueCount += statBo.IssueOverdueCount
		case consts.ProjectObjectTypeLangCodeTestTask:
			iterationStatPo.TesttaskCount += statBo.IssueCount
			iterationStatPo.TesttaskWaitCount += statBo.IssueWaitCount
			iterationStatPo.TesttaskRunningCount += statBo.IssueRunningCount
			iterationStatPo.TesttaskEndCount += statBo.IssueEndCount
			iterationStatPo.TesttaskOverdueCount += statBo.IssueOverdueCount
		}
		iterationStatPo.IssueCount += statBo.IssueCount
		iterationStatPo.IssueWaitCount += statBo.IssueWaitCount
		iterationStatPo.IssueRunningCount += statBo.IssueRunningCount
		iterationStatPo.IssueEndCount += statBo.IssueEndCount
		iterationStatPo.IssueOverdueCount += statBo.IssueOverdueCount
		iterationStatPo.StoryPointCount += statBo.StoryPointCount
		iterationStatPo.StoryPointWaitCount += statBo.StoryPointWaitCount
		iterationStatPo.StoryPointRunningCount += statBo.StoryPointRunningCount
		iterationStatPo.StoryPointEndCount += statBo.StoryPointEndCount
	}
}
