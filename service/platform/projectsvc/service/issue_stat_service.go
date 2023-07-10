package service

import (
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/date"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"upper.io/db.v3"
)

func IssueAssignRank(orgId int64, projectId int64, rankTop int) ([]*vo.IssueAssignRankInfo, errs.SystemErrorInfo) {
	projectExist := domain.JudgeProjectIsExist(orgId, projectId)
	if !projectExist {
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityProject)
	}

	issueAssignCountBos, err := domain.SelectIssueAssignCount(orgId, projectId, rankTop)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	issueAssignRankInfos := make([]*vo.IssueAssignRankInfo, len(issueAssignCountBos))
	rankTotalCount := int64(0)

	ownerIds := make([]int64, 0)
	for _, issueAssignCountBo := range issueAssignCountBos {
		ownerIds = append(ownerIds, issueAssignCountBo.Owner)
	}

	ownerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, ownerIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	ownerMap := maps.NewMap("UserId", ownerInfos)

	for i, issueAssignCountBo := range issueAssignCountBos {
		if ownerInfoInterface, ok := ownerMap[issueAssignCountBo.Owner]; ok {
			ownerInfo := ownerInfoInterface.(bo.BaseUserInfoBo)
			issueAssignRankInfos[i] = &vo.IssueAssignRankInfo{
				Name:                 ownerInfo.Name,
				Avatar:               ownerInfo.Avatar,
				EmplID:               ownerInfo.OutUserId,
				IncompleteissueCount: issueAssignCountBo.Count,
			}
			rankTotalCount += issueAssignCountBo.Count
		} else {
			log.Errorf("用户 %d 信息不存在，组织id %d", issueAssignCountBo.Owner, orgId)
		}
	}
	//加入:其它统计,公式：未完成的任务总数 - rank榜总数
	issueAssignTotalCount, err := domain.SelectIssueAssignTotalCount(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	otherCount := issueAssignTotalCount - rankTotalCount
	issueAssignRankInfos = append(issueAssignRankInfos, &vo.IssueAssignRankInfo{
		Name:                 "其它",
		IncompleteissueCount: otherCount,
	})
	return issueAssignRankInfos, nil
}

func IssueAndProjectCountStat(reqVo projectvo.IssueAndProjectCountStatReqVo) (*vo.IssueAndProjectCountStatResp, errs.SystemErrorInfo) {
	userId := reqVo.UserId
	orgId := reqVo.OrgId

	projectNotCompletedCount, err := domain.GetProjectCountByOwnerId(orgId, userId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	issueNotCompletedCount, err := domain.GetIssueCountByOwnerId(orgId, userId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//
	//参与的归档项目数量,参与项目数量
	filingCount, notfilingCount, err := assemblyProjectCount(orgId, reqVo.UserId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &vo.IssueAndProjectCountStatResp{
		IssueNotCompletedCount:         issueNotCompletedCount,
		ProjectNotCompletedCount:       projectNotCompletedCount,
		ParticipantsProjectCount:       notfilingCount,
		FilingParticipantsProjectCount: filingCount,
	}, nil
}

func assemblyProjectCount(orgId, userId int64) (int64, int64, errs.SystemErrorInfo) {
	relationBo, err := domain.GetProjectRelationByTypeAndUserId(orgId, userId, []int64{consts.IssueRelationTypeOwner, consts.IssueRelationTypeParticipant})
	if err != nil {
		return 0, 0, nil
	}

	//获取处理去重后的projectId
	if relationBo == nil || len(*relationBo) < 1 {
		return 0, 0, nil
	}

	allProjectIds := make([]int64, 0, len(*relationBo))

	for _, value := range *relationBo {
		allProjectIds = append(allProjectIds, value.ProjectId)
	}
	//去重
	distinctProjectIds := slice.SliceUniqueInt64(allProjectIds)

	//查询出所有的任务 然后内存处理判断是否是归档的
	conds := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       db.In(distinctProjectIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}

	bos, _, err := domain.GetProjectList(userId, conds, nil, nil, 0, 0)

	if err != nil {
		return 0, 0, err
	}

	var filingCount int64
	var notfilingCount int64

	for _, value := range bos {

		if value.IsFiling == consts.ProjectIsFiling {
			filingCount++
			continue
		}
		notfilingCount++
	}

	return filingCount, notfilingCount, nil

}

func IssueDailyPersonalWorkCompletionStat(reqVo projectvo.IssueDailyPersonalWorkCompletionStatReqVo) (*vo.IssueDailyPersonalWorkCompletionStatResp, errs.SystemErrorInfo) {
	params := reqVo.Input
	orgId := reqVo.OrgId
	userId := reqVo.UserId

	startDate := (*types.Time)(nil)
	endDate := (*types.Time)(nil)
	if params != nil {
		startDate = params.StartDate
		endDate = params.EndDate
	}

	statBos, err := domain.IssueDailyPersonalWorkCompletionStatForLc(orgId, userId, startDate, endDate)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	statVos := make([]*vo.IssueDailyPersonalWorkCompletionStatData, 0)
	if statBos != nil && len(statBos) > 0 {
		for _, statBo := range statBos {
			statVos = append(statVos, &vo.IssueDailyPersonalWorkCompletionStatData{
				StatDate:       date.Format(statBo.Date),
				CompletedCount: statBo.Count,
			})
		}
	}
	return &vo.IssueDailyPersonalWorkCompletionStatResp{
		List: statVos,
	}, nil
}

//func GetIssueCountByStatus(orgId int64, projectId int64, statusId int64) (uint64, errs.SystemErrorInfo) {
//	return domain.GetIssueCountByStatus(orgId, projectId, statusId)
//}

func AuthProject(orgId, userId, projectId int64, path string, operation string) errs.SystemErrorInfo {
	return domain.AuthProject(orgId, userId, projectId, path, operation)
}
