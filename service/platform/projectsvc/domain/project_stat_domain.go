package domain

import (
	"time"

	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//date: yyyy-MM-dd
func AppendProjectDayStat(projectBo bo.ProjectBo, date string, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	statDate, timeParseError := time.Parse(consts.AppDateFormat, date)
	if timeParseError != nil {
		log.Error(timeParseError)
		return errs.BuildSystemErrorInfo(errs.SystemError)
	}

	orgId := projectBo.OrgId
	projectId := projectBo.Id

	id, err3 := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectDayStat)
	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.ApplyIdError, err3)
	}

	projectCacheBo, err := LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	issueStatusStatBos, err1 := GetIssueStatusStatWithLc(bo.IssueStatusStatCondBo{
		OrgId:     orgId,
		ProjectId: &projectId,
	})
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	projectDayStatPo := &po.PpmStaProjectDayStat{}
	projectDayStatPo.Id = id
	projectDayStatPo.OrgId = orgId
	projectDayStatPo.ProjectId = projectId
	projectDayStatPo.StatDate = statDate
	projectDayStatPo.Status = projectCacheBo.Status

	issueStatusStatMap := maps.NewMap("ProjectTypeId", issueStatusStatBos)
	ext := bo.StatExtBo{}
	ext.Issue = bo.StatIssueExtBo{
		Data: issueStatusStatMap,
	}
	projectDayStatPo.Ext = json.ToJsonIgnoreError(ext)
	//封装状态
	assemblyProjectStat(issueStatusStatBos, projectDayStatPo)

	err5 := dao.InsertProjectDayStat(*projectDayStatPo, tx...)
	if err5 != nil {
		log.Error(err5)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func assemblyProjectStat(issueStatusStatBos []bo.IssueStatusStatBo, projectDayStatPo *po.PpmStaProjectDayStat) {
	for _, statBo := range issueStatusStatBos {
		switch statBo.ProjectTypeLangCode {
		case consts.ProjectObjectTypeLangCodeDemand:
			projectDayStatPo.DemandCount += statBo.IssueCount
			projectDayStatPo.DemandWaitCount += statBo.IssueWaitCount
			projectDayStatPo.DemandRunningCount += statBo.IssueRunningCount
			projectDayStatPo.DemandEndCount += statBo.IssueEndCount
			projectDayStatPo.DemandOverdueCount += statBo.IssueOverdueCount
		case consts.ProjectObjectTypeLangCodeTask:
			projectDayStatPo.TaskCount += statBo.IssueCount
			projectDayStatPo.TaskWaitCount += statBo.IssueWaitCount
			projectDayStatPo.TaskRunningCount += statBo.IssueRunningCount
			projectDayStatPo.TaskEndCount += statBo.IssueEndCount
			projectDayStatPo.TaskOverdueCount += statBo.IssueOverdueCount
		case consts.ProjectObjectTypeLangCodeBug:
			projectDayStatPo.BugCount += statBo.IssueCount
			projectDayStatPo.BugWaitCount += statBo.IssueWaitCount
			projectDayStatPo.BugRunningCount += statBo.IssueRunningCount
			projectDayStatPo.BugEndCount += statBo.IssueEndCount
			projectDayStatPo.BugOverdueCount += statBo.IssueOverdueCount
		case consts.ProjectObjectTypeLangCodeTestTask:
			projectDayStatPo.TesttaskCount += statBo.IssueCount
			projectDayStatPo.TesttaskWaitCount += statBo.IssueWaitCount
			projectDayStatPo.TesttaskRunningCount += statBo.IssueRunningCount
			projectDayStatPo.TesttaskEndCount += statBo.IssueEndCount
			projectDayStatPo.TesttaskOverdueCount += statBo.IssueOverdueCount
		}
		projectDayStatPo.IssueCount += statBo.IssueCount
		projectDayStatPo.IssueWaitCount += statBo.IssueWaitCount
		projectDayStatPo.IssueRunningCount += statBo.IssueRunningCount
		projectDayStatPo.IssueEndCount += statBo.IssueEndCount
		projectDayStatPo.IssueOverdueCount += statBo.IssueOverdueCount
		projectDayStatPo.StoryPointCount += statBo.StoryPointCount
		projectDayStatPo.StoryPointWaitCount += statBo.StoryPointWaitCount
		projectDayStatPo.StoryPointRunningCount += statBo.StoryPointRunningCount
		projectDayStatPo.StoryPointEndCount += statBo.StoryPointEndCount
	}
}

func GetProjectCountByOwnerId(orgId, ownerId int64) (int64, errs.SystemErrorInfo) {
	finishedIds, _ := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryProject, consts.StatusTypeComplete)
	if len(finishedIds) == 0 {
		log.Errorf("[GetProjectCountByOwnerId] GetStatusIdsByCategory err: not finished status id")
		return 0, errs.BuildSystemErrorInfo(errs.CacheProxyError, errs.ProcessStatusNotExist)
	}

	total, err1 := mysql.SelectCountByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcOwner:    ownerId,
		consts.TcStatus:   db.NotIn(finishedIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcIsFiling: consts.AppIsNotFilling,
	})
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return int64(total), nil
}

func PayLimitNum(orgId int64) (*projectvo.PayLimitNumRespData, errs.SystemErrorInfo) {
	projectNum, err := mysql.SelectCountByCond(consts.TableProject, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcTemplateFlag: consts.TemplateFalse,
	})
	if err != nil {
		return nil, errs.MysqlOperateError
	}

	projectPos := []*po.PpmProProject{}
	err = mysql.SelectAllByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcTemplateFlag: consts.TemplateTrue,
	}, &projectPos)
	if err != nil {
		log.Errorf("[PayLimitNum] err:%v, orgId:%v", err, orgId)
		return nil, errs.MysqlOperateError
	}
	projectIds := make([]int64, 0, len(projectPos))
	for _, pro := range projectPos {
		projectIds = append(projectIds, pro.Id)
	}

	//issueNum, err := mysql.SelectCountByCond(consts.TableIssue, db.Cond{
	//	consts.TcOrgId:     orgId,
	//	consts.TcProjectId: db.NotIn(db.Raw("select id from ppm_pro_project where org_id = ? and is_delete = 2 and template_flag = 1", orgId)),
	//})
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.MysqlOperateError
	//}

	var condition *tablePb.Condition
	if len(projectIds) > 0 {
		condition = &tablePb.Condition{}
		condition.Type = tablePb.ConditionType_and
		condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_not_in, nil, projectIds))
	}
	issueNum, errSys := getRowsCountByCondition(orgId, 0, condition)
	if errSys != nil {
		log.Errorf("[PayLimitNum] getRowsCountByCondition err:%v, orgId:%v", err, orgId)
		return nil, errSys
	}

	// 加上被删除的任务
	deletedIssueNum, errCache := GetDeletedIssueNum(orgId)
	if errCache != nil {
		log.Errorf("[PayLimitNum] GetDeletedIssueNum err:%v, orgId:%v", errCache, orgId)
	}

	log.Infof("[PayLimitNum] issueNum:%v, deletedIssueNum:%v, orgId:%v", issueNum, deletedIssueNum, orgId)

	resourceResp := resourcefacade.CacheResourceSize(resourcevo.CacheResourceSizeReq{OrgId: orgId})
	if resourceResp.Failure() {
		log.Error(resourceResp.Error())
		return nil, resourceResp.Error()
	}

	return &projectvo.PayLimitNumRespData{
		ProjectNum: int64(projectNum),
		IssueNum:   issueNum + deletedIssueNum,
		FileSize:   resourceResp.Size,
		DashNum:    int64(projectNum), // 现有仪表盘数量等于项目数。后续如果规则发生变化，这里的逻辑也会发生变化
	}, nil
}
