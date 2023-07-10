package service

import (
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func ProjectDayStats(orgId int64, page uint, size uint, params *vo.ProjectDayStatReq) (*vo.ProjectDayStatList, errs.SystemErrorInfo) {
	projectId := params.ProjectID
	projectBo, err1 := domain.GetProjectSimple(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityProject)
	}

	//默认查询十五天
	startTime, endTime, err := domain.CalDateRangeCond(params.StartDate, params.EndDate, 15)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//如果and后的日期是到天的，需要加一天
	condEndTime := endTime.AddDate(0, 0, 1)

	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcAppId:    projectBo.AppId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcDate:     db.Between(startTime.Format(consts.AppDateFormat), condEndTime.Format(consts.AppDateFormat)),
	}

	statDailyAppList, total, errSys := domain.GetStatDailyAppList(page, size, cond)
	if errSys != nil {
		log.Errorf("[ProjectDayStats] GetStatDailyAppList err:%v, orgId:%v, projectId:%v", errSys, orgId, projectId)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, errSys)
	}
	statDailyMap := make(map[int64]*po.LcStatDailyApp, total)
	for _, statPo := range statDailyAppList {
		statDailyMap[statPo.Id] = statPo
	}

	resultList := []*vo.ProjectDayStat{}
	err3 := copyer.Copy(&statDailyAppList, &resultList)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	for _, stat := range resultList {
		if v, ok := statDailyMap[stat.ID]; ok {
			stat.StatDate = types.Time(v.Date)
		}
	}

	return &vo.ProjectDayStatList{
		Total: total,
		List:  resultList,
	}, nil
}

//date: yyyy-MM-dd
func AppendProjectDayStat(projectBo bo.ProjectBo, date string) errs.SystemErrorInfo {
	return domain.AppendProjectDayStat(projectBo, date)
}

// PayLimitNum 组织下已经使用的资源数量统计
func PayLimitNum(orgId int64) (*projectvo.PayLimitNumRespData, errs.SystemErrorInfo) {
	return domain.PayLimitNum(orgId)
}

func GetProjectStatistics(orgId, userId int64, req projectvo.GetProjectStatisticsReq) ([]*projectvo.GetProjectStatisticsData, errs.SystemErrorInfo) {
	params := req.ExtraInfoReq
	projectIds := make([]int64, 0, len(params))
	resourceIds := make([]int64, 0, len(params))
	for _, info := range params {
		projectIds = append(projectIds, info.ProjectId)
		resourceIds = append(resourceIds, info.ResourceId)
	}
	summaryAppId, err := domain.GetOrgSummaryAppId(orgId)
	if err != nil {
		log.Errorf("[GetProjectStatistics]GetOrgSummaryAppId err:%v, orgId:%v", err, orgId)
		return nil, err
	}

	statForProjects, err := domain.GetIssueFinishedStatForProjects(orgId, userId, summaryAppId, projectIds)
	if err != nil {
		log.Errorf("[GetProjectStatistics]GetIssueFinishedStatForProjects err:%v, orgId:%v, userId:%v, projectId:%v",
			err, orgId, userId, projectIds)
		return nil, err
	}
	resourceBosMap := make(map[int64]bo.ResourceBo)
	if len(resourceIds) > 0 {
		resourceResp := resourcefacade.GetResourceById(
			resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: resourcevo.GetResourceByIdReqBody{ResourceIds: resourceIds}},
		)
		if resourceResp.Failure() {
			log.Errorf("[GetProjectStatistics] GetResourceById err:%v, orgId:%v", resourceResp.Error(), orgId)
			return nil, resourceResp.Error()
		}
		for _, res := range resourceResp.ResourceBos {
			resourceBosMap[res.Id] = res
		}
	}

	res := []*projectvo.GetProjectStatisticsData{}

	for _, info := range params {
		data := &projectvo.GetProjectStatisticsData{
			ProjectId:       info.ProjectId,
			CoverResourceId: info.ResourceId,
		}
		if stat, ok := statForProjects[info.ProjectId]; ok {
			data.AllIssues = stat.All
			data.FinishIssues = stat.Finish
			data.OverdueIssues = stat.Overdue
			data.UnfinishedIssues = stat.RelateUnfinish
		}
		if resourceBo, ok := resourceBosMap[info.ResourceId]; ok {
			coverUrl := util.JointUrl(resourceBo.Host, resourceBo.Path)
			resourcePath := util.GetCompressedPath(coverUrl, resourceBo.Type)
			data.CoverResourcePath = resourcePath
		}
		res = append(res, data)
	}
	return res, nil
}
