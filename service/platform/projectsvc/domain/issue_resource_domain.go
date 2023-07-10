package domain

import (
	"fmt"
	"time"

	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/datacenter"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/formfacade"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
	"upper.io/db.v3/lib/sqlbuilder"
)

func CreateIssueResource(orgId, userId int64, input *vo.CreateIssueResourceReq) (int64, errs.SystemErrorInfo) {
	bucketName := ""
	md5 := ""
	if input.BucketName != nil {
		bucketName = *input.BucketName
	}
	if input.Md5 != nil {
		md5 = *input.Md5
	}

	resourceType := input.ResourceType
	sourceType := input.PolicyType
	skipCreateRelation := true
	//if input.IssueId == 0 {
	//	skipCreateRelation = true // 没有issueId时跳过创建relation
	//}
	createResourceBo := bo.CreateResourceBo{
		ProjectId:          input.ProjectId,
		IssueId:            input.IssueId,
		OrgId:              orgId,
		Path:               input.ResourcePath,
		Name:               input.FileName,
		Suffix:             input.FileSuffix,
		Bucket:             bucketName,
		Type:               resourceType,
		Size:               input.ResourceSize,
		Md5:                md5,
		OperatorId:         userId,
		SourceType:         &sourceType,
		SkipCreateRelation: skipCreateRelation,
	}
	respVo := resourcefacade.CreateResource(resourcevo.CreateResourceReqVo{CreateResourceBo: createResourceBo})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return 0, respVo.Error()
	}
	resourceId := respVo.ResourceId

	// 不在 任务关系表ppm_pri_issue_relation中 维护资源和任务的关系
	//_, err2 := UpdateIssueRelationSingle(operatorId, issueBo, consts.IssueRelationTypeResource, []int64{resId})
	//if err2 != nil {
	//	log.Error(err2)
	//	return 0, err2
	//}

	if input.IssueId > 0 {
		issueBos, err := GetIssueInfosLc(orgId, userId, []int64{input.IssueId})
		if err != nil || len(issueBos) == 0 {
			return resourceId, nil
		}
		issueBo := issueBos[0]

		asyn.Execute(func() {
			resourceTrend := []bo.ResourceInfoBo{}
			resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
				Name:       input.FileName,
				Url:        input.ResourcePath,
				Size:       input.ResourceSize,
				UploadTime: time.Now(),
				Suffix:     input.FileSuffix,
			})
			issueTrendsBo := &bo.IssueTrendsBo{
				PushType:      consts.PushTypeUploadResource,
				OrgId:         issueBo.OrgId,
				OperatorId:    userId,
				DataId:        issueBo.DataId,
				IssueId:       issueBo.Id,
				ParentIssueId: issueBo.ParentId,
				ProjectId:     issueBo.ProjectId,
				TableId:       issueBo.TableId,
				PriorityId:    issueBo.PriorityId,
				ParentId:      issueBo.ParentId,
				IssueTitle:    issueBo.Title,
				IssueStatusId: issueBo.Status,
				Ext: bo.TrendExtensionBo{
					ObjName:      issueBo.Title,
					ResourceInfo: resourceTrend,
				},
			}
			asyn.Execute(func() {
				PushIssueTrends(issueTrendsBo)
			})
			asyn.Execute(func() {
				PushIssueThirdPlatformNotice(issueTrendsBo)
			})
		})
	}
	return resourceId, nil
}

func DeleteIssueResource(issueBo bo.IssueBo, deletedResourceIds []int64, operatorId int64) errs.SystemErrorInfo {
	return DeleteIssueRelationByIds(operatorId, issueBo, consts.IssueRelationTypeResource, deletedResourceIds)
}

// 将附件字段的 附件资源 关联删除
func RecycleIssueResources(orgId, userId, projectId, issueId int64, deleteResourceIds []int64, recycleVersionId int64, columnId string) errs.SystemErrorInfo {
	// 删除表关联
	if deleteResourceIds == nil || len(deleteResourceIds) == 0 {
		return nil
	}
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 删除任务和附件的关联
		//upd := mysql.Upd{
		//	consts.TcIsDelete: consts.AppIsDeleted,
		//	consts.TcUpdator:  userId,
		//}
		//conds := db.Cond{
		//	consts.TcOrgId:        orgId,
		//	consts.TcIssueId:      issueId,
		//	consts.TcRelationId:   db.In(deleteResourceIds),
		//	consts.TcRelationType: consts.IssueRelationTypeResource,
		//	consts.TcIsDelete:     consts.AppIsNoDelete,
		//}
		//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, conds, upd)
		//if err != nil {
		//	log.Errorf("[RecycleIssueResources] update TableIssueRelation err:%v, orgId:%d, issueId:%d, delResourceIds:%v",
		//		err, orgId, issueId, deleteResourceIds)
		//	return err
		//}
		// 修改资源关联表version为回收站的id，并存入columnId
		resp := resourcefacade.DeleteAttachmentRelation(resourcevo.DeleteAttachmentRelationReq{
			OrgId:  orgId,
			UserId: userId,
			Input: &resourcevo.DeleteAttachmentRelationData{
				ProjectId:        projectId,
				IssueIds:         []int64{issueId},
				RecycleVersionId: recycleVersionId,
				ColumnId:         columnId,
				ResourceIds:      deleteResourceIds,
			},
		})
		if resp.Failure() {
			log.Errorf("[RecycleIssueResources] DeleteAttachmentRelation err:%v, orgId:%d, issueId:%d, delResourceIds:%v",
				resp.Error(), orgId, issueId, deleteResourceIds)
			return resp.Error()
		}

		return nil
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func DeleteLessIssueAttachmentsData(orgId, userId int64, projectId, issueId int64, resourceIds []int64, columnId string) errs.SystemErrorInfo {
	appId, err := GetAppIdFromProjectId(orgId, projectId)
	if err != nil {
		log.Errorf("[DeleteLessIssueAttachmentsData]GetAppIdFromProjectId err:%v", err)
		return err
	}
	condsData := []*vo.LessCondsData{
		{
			Type:   consts.ConditionEqual,
			Value:  orgId,
			Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
		},
		{
			Type:   consts.ConditionEqual,
			Value:  issueId,
			Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
		},
	}

	condition := vo.LessCondsData{
		Type:  consts.ConditionAnd,
		Conds: condsData,
	}
	defaultValue := consts.LcJsonColumn
	for _, resourceId := range resourceIds {
		defaultValue += fmt.Sprintf("#-'{%s,%d}'", columnId, resourceId)
	}
	sets := []datacenter.Set{
		{
			Column:          consts.LcJsonColumn,
			Value:           defaultValue,
			Type:            consts.SetTypeJson,
			Action:          consts.SetActionSet,
			WithoutPretreat: true,
		},
	}
	batchRaw := formfacade.LessUpdateIssueBatchRaw(&formvo.LessUpdateIssueBatchReq{
		OrgId:     orgId,
		AppId:     appId,
		UserId:    userId,
		Condition: condition,
		Sets:      sets,
	})

	if batchRaw.Failure() {
		log.Errorf("[UpdateLessIssueAttachmentsRecycle]LessUpdateIssueBatchRaw err:%v", batchRaw.Error())
		return batchRaw.Error()
	}

	return nil
}
