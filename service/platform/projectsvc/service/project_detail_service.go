package service

import (
	"time"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

func ProjectDetail(orgId, projectId int64) (*vo.ProjectDetail, errs.SystemErrorInfo) {

	bos, err := domain.GetProjectDetailByProjectIdBo(projectId, orgId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}
	result := &vo.ProjectDetail{}
	copyErr := copyer.Copy(bos, result)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return result, nil
}

func CreateProjectDetail(currentUserId int64, input vo.CreateProjectDetailReq) (*vo.Void, errs.SystemErrorInfo) {

	//TODO 权限
	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	//}

	entity := &bo.ProjectDetailBo{}
	copyErr := copyer.Copy(input, entity)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectDetail)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	entity.Id = id
	entity.Creator = currentUserId
	entity.Updator = currentUserId

	err1 := domain.CreateProjectDetail(entity)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	return &vo.Void{
		ID: id,
	}, nil
}

func UpdateProjectDetail(orgId, currentUserId int64, input vo.UpdateProjectDetailReq) (*vo.Void, errs.SystemErrorInfo) {
	if input.Notice != nil {
		isRight := format.VerifyProjectNoticeFormat(*input.Notice)
		if !isRight {
			return nil, errs.ProjectNoticeLenError
		}
	}
	entity := &bo.ProjectDetailBo{}
	copyErr := copyer.Copy(input, entity)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	entity.Updator = currentUserId
	entity.UpdateTime = time.Now()

	//是否存在
	projectDetail, err2 := domain.GetProjectDetailBo(entity.Id)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}

	authErr := domain.AuthProject(orgId, currentUserId, projectDetail.ProjectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigModify)
	if authErr != nil {
		log.Error(authErr)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, authErr)
	}

	err1 := domain.UpdateProjectDetail(entity)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}
	if input.IsSyncOutCalendar != nil && *input.IsSyncOutCalendar != projectDetail.IsSyncOutCalendar {
		//删除缓存
		delErr := domain.DeleteProjectCalendarInfo(orgId, projectDetail.ProjectId)
		if *input.IsSyncOutCalendar == consts.IsSyncOutCalendar {
		}
		if delErr != nil {
			return nil, delErr
		}
	}

	if input.Notice != nil {
		asyn.Execute(func() {
			projectInfo, projectInfoErr := domain.GetProjectSimple(orgId, projectDetail.ProjectId)
			if projectInfoErr != nil {
				log.Error(projectInfoErr)
				return
			}
			changeList := []bo.TrendChangeListBo{}
			changeList = append(changeList, bo.TrendChangeListBo{
				Field:     "notice",
				FieldName: consts.ProjectNotice,
				OldValue:  projectDetail.Notice,
				NewValue:  *input.Notice,
			})
			ext := bo.TrendExtensionBo{}
			ext.ObjName = projectInfo.Name
			ext.ChangeList = changeList
			domain.PushProjectTrends(bo.ProjectTrendsBo{
				PushType:   consts.PushTypeUpdateProject,
				OrgId:      orgId,
				ProjectId:  projectDetail.ProjectId,
				OperatorId: currentUserId,
				Ext:        ext,
			})
		})
	}

	return &vo.Void{
		ID: input.ID,
	}, nil
}

func DeleteProjectDetail(currentUserId int64, input vo.DeleteProjectDetailReq) (*vo.Void, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//currentUserId := cacheUserInfo.UserId
	targetId := input.ID

	//TODO 权限
	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	//}

	bo, err1 := domain.GetProjectDetailBo(targetId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	err2 := domain.DeleteProjectDetail(bo, currentUserId)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}

	return &vo.Void{
		ID: targetId,
	}, nil
}
