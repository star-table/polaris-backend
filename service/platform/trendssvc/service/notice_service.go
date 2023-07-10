package service

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/trendsvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/service/platform/trendssvc/domain"
)

func UnreadNoticeCount(orgId, userId int64) (uint64, errs.SystemErrorInfo) {
	return domain.UnreadNoticeCount(orgId, userId)
}

func GetCallMeCount(req trendsvo.CallMeCountReqVo) (uint64, errs.SystemErrorInfo) {
	return domain.CallMeCount(req)
}

func NoticeList(orgId, userId int64, page, size int, input *vo.NoticeListReq) (*vo.NoticeList, errs.SystemErrorInfo) {
	count, list, err := domain.GetNoticeList(orgId, userId, page, size, input)
	if err != nil {
		return nil, err
	}

	if len(*list) == 0 {
		return &vo.NoticeList{
			Total: int64(count),
			List:  nil,
		}, nil
	}

	creators := []int64{}
	issueIds := []int64{}
	projectIds := []int64{}
	for _, v := range *list {
		creators = append(creators, v.Creator)
		issueIds = append(issueIds, v.IssueId)
		projectIds = append(projectIds, v.ProjectId)
	}

	//用户信息
	creatorResp, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, slice.SliceUniqueInt64(creators))
	if err != nil {
		return nil, err
	}
	creatorInfo := map[int64]vo.UserIDInfo{}
	for _, v := range creatorResp {
		creatorInfo[v.UserId] = vo.UserIDInfo{
			UserID: v.UserId,
			Avatar: v.Avatar,
			Name:   v.Name,
			EmplID: v.OutUserId,
		}
	}
	//任务信息
	issueResp := projectfacade.GetLcIssueInfoBatch(projectvo.GetLcIssueInfoBatchReqVo{OrgId: orgId, IssueIds: slice.SliceUniqueInt64(issueIds)})
	if issueResp.Failure() {
		log.Error(issueResp.Error())
		return nil, issueResp.Error()
	}
	issueInfo := map[int64]string{}
	issueParent := map[int64]int64{}
	for _, v := range issueResp.Data {
		issueInfo[v.Id] = v.Title
		issueParent[v.Id] = v.ParentId
	}
	//项目信息
	projectResp := projectfacade.GetSimpleProjectInfo(projectvo.GetSimpleProjectInfoReqVo{OrgId: orgId, Ids: slice.SliceUniqueInt64(projectIds)})
	if projectResp.Failure() {
		log.Error(projectResp.Error())
		return nil, projectResp.Error()
	}
	projectInfo := map[int64]string{}
	for _, v := range *projectResp.Data {
		projectInfo[v.ID] = v.Name
	}

	noticeVo := &[]*vo.Notice{}
	copyErr := copyer.Copy(list, noticeVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	for k, v := range *noticeVo {
		if user, ok := creatorInfo[v.Creator]; ok {
			(*noticeVo)[k].CreatorInfo = &user
		}
		if issueName, ok := issueInfo[v.IssueID]; ok {
			(*noticeVo)[k].IssueName = issueName
		}
		if parentId, ok := issueParent[v.IssueID]; ok {
			(*noticeVo)[k].ParentIssueID = parentId
		} else {
			(*noticeVo)[k].ParentIssueID = 0
		}
		if projectName, ok := projectInfo[v.ProjectID]; ok {
			(*noticeVo)[k].ProjectName = projectName
		}
	}

	return &vo.NoticeList{
		Total: int64(count),
		List:  *noticeVo,
	}, nil
}
