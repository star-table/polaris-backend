package domain

import (
	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	int64Slice "github.com/star-table/polaris-backend/common/core/util/slice/int64"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/facade/resourcefacade"
	"github.com/spf13/cast"
)

func AddIssueTrends(issueTrendsBo bo.IssueTrendsBo) {
	var trendsBos []bo.TrendsBo = nil

	assemblyError := assemblyTrendsBos(issueTrendsBo, &trendsBos)
	if assemblyError != nil || trendsBos == nil {
		return
	}

	for i, _ := range trendsBos {
		trendsBos[i].CreateTime = types.Time(issueTrendsBo.OperateTime)
	}

	err := CreateTrendsBatch(trendsBos)
	if err != nil {
		log.Error(err)
	}
}

func assemblyTrendsBos(issueTrendsBo bo.IssueTrendsBo, trendsBos *[]bo.TrendsBo) errs.SystemErrorInfo {
	trendsType := issueTrendsBo.PushType
	var tBos []bo.TrendsBo
	var err1 errs.SystemErrorInfo = nil

	log.Infof("[assemblyTrendsBos] %v", json.ToJsonIgnoreError(issueTrendsBo))

	if trendsType == consts.PushTypeCreateIssue {
		*trendsBos, err1 = assemblyCreateIssueTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeUpdateIssue {
		// 更新任务
		// 因为工时变更是单独产生的。防止变更工时时，产生两个动态。
		if !issueTrendsBo.UpdateWorkHour {
			tBos, err1 = assemblyUpdateIssueTrends(issueTrendsBo)
			if err1 != nil {
				log.Errorf(consts.Assembly_UpdateIssueTrends_eorror_printf, err1)
				return err1
			}
			*trendsBos = append(*trendsBos, tBos...)
		}

		// 关联、前后置
		tBos, err1 = assemblyUpdateRelationIssueTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_UpdateIssueTrends_eorror_printf, err1)
			return err1
		}
		*trendsBos = append(*trendsBos, tBos...)

		// 附件
		tBos, err1 = assemblyDocumentsTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_UpdateIssueTrends_eorror_printf, err1)
			return err1
		}
		*trendsBos = append(*trendsBos, tBos...)

		// 图片
		tBos, err1 = assemblyPicturesTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_UpdateIssueTrends_eorror_printf, err1)
			return err1
		}
		*trendsBos = append(*trendsBos, tBos...)

		// 工时
		tBos, err1 = assemblyTrendForWorkHour(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_UpdateIssueTrends_eorror_printf, err1)
			return err1
		}
		*trendsBos = append(*trendsBos, tBos...)

	} else if trendsType == consts.PushTypeDeleteIssue {
		*trendsBos, err1 = assemblyDeleteIssueTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeIssueComment {
		*trendsBos, err1 = assemblyCreateIssueComment(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeUpdateIssueProjectTable {
		*trendsBos, err1 = assemblyUpdateIssueProjectTable(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeUpdateIssueParent {
		*trendsBos, err1 = assemblyUpdateIssueParent(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeAddIssueTag {
		*trendsBos, err1 = assemblyAddTag(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeDeleteIssueTag {
		*trendsBos, err1 = assemblyDeleteTag(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeRecoverIssue {
		*trendsBos, err1 = assemblyRecoverIssueTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeReferResource {
		*trendsBos, err1 = assemblyReferResource(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeDeleteReferResource {
		*trendsBos, err1 = assemblyDeleteReferResource(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeCreateWorkHour {
		*trendsBos, err1 = assemblyTrendForCreateWorkHour(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeUpdateWorkHour {
		*trendsBos, err1 = assemblyTrendForUpdateWorkHour(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeDeleteWorkHour {
		*trendsBos, err1 = assemblyTrendForDeleteWorkHour(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeAuditIssue {
		*trendsBos, err1 = assemblyAuditIssueTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeWithdrawIssue {
		*trendsBos, err1 = assemblyWithdrawIssueTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeIssueStartChat {
		*trendsBos, err1 = assemblyStartIssueChat(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	}

	return nil
}

// 检查 IssueTrendsBo 中的信息，判断出是新增/更新/删除
func assemblyTrendForWorkHour(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := make([]bo.TrendsBo, 0)
	_, add, del := int64Slice.CompareSliceAddDelInt64(issueTrendsBo.AfterWorkHourIds, issueTrendsBo.BeforeWorkHourIds)
	// 取交集，有交集意味着这些工时记录中有字段（工时起至时间、执行人、描述等）发生变化
	editRecordIds := int64Slice.Int64Intersect(issueTrendsBo.AfterWorkHourIds, issueTrendsBo.BeforeWorkHourIds)
	if len(add) > 0 {
		trendsTiles, err := assemblyTrendForCreateWorkHour(issueTrendsBo)
		if err != nil {
			log.Errorf("[assemblyTrendForWorkHour] assemblyTrendForCreateWorkHourerr: %v, issueId: %v", err, issueTrendsBo.IssueId)
			return trendsBos, err
		}
		trendsBos = append(trendsBos, trendsTiles...)
	}
	if len(del) > 0 {
		trendsTiles, err := assemblyTrendForDeleteWorkHour(issueTrendsBo)
		if err != nil {
			log.Errorf("[assemblyTrendForWorkHour] assemblyTrendForDeleteWorkHour err: %v, issueId: %v", err, issueTrendsBo.IssueId)
			return trendsBos, err
		}
		trendsBos = append(trendsBos, trendsTiles...)
	}
	if len(editRecordIds) > 0 {
		trendsTiles, err := assemblyTrendForUpdateWorkHour(issueTrendsBo)
		if err != nil {
			log.Errorf("[assemblyTrendForWorkHour] assemblyTrendForUpdateWorkHour err: %v, issueId: %v", err, issueTrendsBo.IssueId)
			return trendsBos, err
		}
		trendsBos = append(trendsBos, trendsTiles...)
	}

	return trendsBos, nil
}

// 新增工时，增加动态，新增的动态数据（工时）存放在 NewValue 字段中
func assemblyTrendForCreateWorkHour(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := make([]bo.TrendsBo, 0)
	ext := issueTrendsBo.Ext
	ext.IssueType = "T"
	ext.ObjName = issueTrendsBo.IssueTitle
	operatorId := issueTrendsBo.OperatorId
	newObject := vo.OneWorkHourRecord{}
	json.FromJson(issueTrendsBo.NewValue, &newObject)
	//拼装动态(这里要求实时)
	trendsBo := bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationCreate,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeWorkHour,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeCreateWorkHour,
		RelationObjType: consts.TrendsOperObjectTypeWorkHour,
		NewValue:        &issueTrendsBo.NewValue,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(ext),
		Creator:         operatorId,
	}
	trendsBos = append(trendsBos, trendsBo)
	return trendsBos, nil
}

// 编辑工时，新增动态。变更的数据位于 ext.ChangeList 字段中
func assemblyTrendForUpdateWorkHour(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := make([]bo.TrendsBo, 0)
	ext := issueTrendsBo.Ext
	ext.IssueType = "T"
	ext.ObjName = issueTrendsBo.IssueTitle
	operatorId := issueTrendsBo.OperatorId
	newObject := vo.OneWorkHourRecord{}
	json.FromJson(issueTrendsBo.NewValue, &newObject)
	//拼装动态(这里要求实时)
	trendsBo := bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeWorkHour,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeUpdateIssueWorkHour,
		RelationObjType: consts.TrendsOperObjectTypeWorkHour,
		NewValue:        &issueTrendsBo.NewValue,
		OldValue:        &issueTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(ext),
		Creator:         operatorId,
	}
	trendsBos = append(trendsBos, trendsBo)

	return trendsBos, nil
}

// 删除工时，新增动态。
func assemblyTrendForDeleteWorkHour(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := make([]bo.TrendsBo, 0)
	ext := issueTrendsBo.Ext
	ext.IssueType = "T"
	ext.ObjName = issueTrendsBo.IssueTitle
	operatorId := issueTrendsBo.OperatorId
	newObject := vo.OneWorkHourRecord{}
	json.FromJson(issueTrendsBo.NewValue, &newObject)
	//拼装动态(这里要求实时)
	trendsBo := bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationDelete,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeWorkHour,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeDeleteIssueWorkHour,
		RelationObjType: consts.TrendsOperObjectTypeWorkHour,
		NewValue:        &issueTrendsBo.NewValue,
		OldValue:        &issueTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(ext),
		Creator:         operatorId,
	}
	trendsBos = append(trendsBos, trendsBo)

	return trendsBos, nil
}

// 创建评论
func assemblyCreateIssueComment(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := make([]bo.TrendsBo, 0)
	ext := issueTrendsBo.Ext
	ext.IssueType = "T"
	ext.ObjName = issueTrendsBo.IssueTitle

	operatorId := issueTrendsBo.OperatorId

	commentBo := ext.CommentBo
	commentId := commentBo.Id

	//拼装动态(这里要求实时)
	trendsBo := bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationCreate,
		OperObjId:       commentId,
		OperObjType:     consts.TrendsOperObjectTypeComment,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeCreateIssueComment,
		RelationObjType: consts.TrendsOperObjectTypeIssue,
		NewValue:        &commentBo.Content,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(ext),
		Creator:         operatorId,
	}
	trendsBos = append(trendsBos, trendsBo)
	return trendsBos, nil
}

// 创建任务
func assemblyCreateIssueTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	var trendsBos []bo.TrendsBo
	trendsBo := bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationCreate,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.ParentIssueId,
		RelationType:    consts.TrendsRelationTypeCreateIssue,
		RelationObjType: consts.TrendsOperObjectTypeIssue,
		NewValue:        &issueTrendsBo.NewValue,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	}

	trendsBos = append(trendsBos, trendsBo)
	return trendsBos, nil
}

// 更新任务
func assemblyUpdateIssueTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	var trendsBos []bo.TrendsBo
	if len(issueTrendsBo.Ext.ChangeList) == 0 {
		return trendsBos, nil
	}
	trendsBo := bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: issueTrendsBo.OperateObjProperty,
		RelationObjId:   0,
		RelationType:    consts.TrendsRelationTypeUpdateIssue,
		RelationObjType: consts.BlankString,
		NewValue:        &issueTrendsBo.NewValue,
		OldValue:        &issueTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	}
	trendsBos = append(trendsBos, trendsBo)
	return trendsBos, nil
}

func assemblyRelatingOpRelationType(columnType string, isLinkTo, isAdd bool) (string, string, string) {
	switch columnType {
	case tablePb.ColumnType_relating.String():
		if isAdd {
			return consts.RoleOperationBind, consts.TrendsRelationTypeAddRelationIssue, consts.TrendsRelationTypeAddRelationIssueByOther
		} else {
			return consts.RoleOperationUnbind, consts.TrendsRelationTypeDeleteRelationIssue, consts.TrendsRelationTypeDeleteRelationIssueByOther
		}
	case tablePb.ColumnType_baRelating.String():
		if isLinkTo {
			if isAdd {
				return consts.RoleOperationBind, consts.TrendsRelationTypeAddBeforeIssue, consts.TrendsRelationTypeAddBeforeIssueByOther
			} else {
				return consts.RoleOperationUnbind, consts.TrendsRelationTypeDeleteBeforeIssue, consts.TrendsRelationTypeDeleteBeforeIssueByOther
			}
		} else {
			if isAdd {
				return consts.RoleOperationBind, consts.TrendsRelationTypeAddAfterIssue, consts.TrendsRelationTypeAddAfterIssueByOther
			} else {
				return consts.RoleOperationUnbind, consts.TrendsRelationTypeDeleteAfterIssue, consts.TrendsRelationTypeDeleteAfterIssueByOther
			}
		}
	case tablePb.ColumnType_singleRelating.String():
		if isAdd {
			return consts.RoleOperationBind, consts.TrendsRelationTypeAddSingleRelationIssue, ""
		} else {
			return consts.RoleOperationUnbind, consts.TrendsRelationTypeDeleteSingleRelationIssue, ""
		}
	default:
		return "", "", ""
	}
}

//func assemblyRelatingInfo(i int, operCode, relationType, relationTypePeer *string) {
//	if i == 0 || i == 1 || i == 2 || i == 3 {
//		if i == 0 || i == 2 {
//			*operCode = consts.RoleOperationUnbind
//			*relationType = consts.TrendsRelationTypeDeleteRelationIssue
//			*relationTypePeer = consts.TrendsRelationTypeDeleteRelationIssueByOther
//		} else if i == 1 || i == 3 {
//			*operCode = consts.RoleOperationBind
//			*relationType = consts.TrendsRelationTypeAddRelationIssue
//			*relationTypePeer = consts.TrendsRelationTypeAddRelationIssueByOther
//		}
//	} else if i == 4 || i == 5 {
//		if i == 4 {
//			*operCode = consts.RoleOperationUnbind
//			*relationType = consts.TrendsRelationTypeDeleteBeforeIssue
//			*relationTypePeer = consts.TrendsRelationTypeDeleteBeforeIssueByOther
//		} else if i == 5 {
//			*operCode = consts.RoleOperationBind
//			*relationType = consts.TrendsRelationTypeAddBeforeIssue
//			*relationTypePeer = consts.TrendsRelationTypeAddBeforeIssueByOther
//		}
//	} else if i == 6 || i == 7 {
//		if i == 6 {
//			*operCode = consts.RoleOperationUnbind
//			*relationType = consts.TrendsRelationTypeDeleteAfterIssue
//			*relationTypePeer = consts.TrendsRelationTypeDeleteAfterIssueByOther
//		} else if i == 7 {
//			*operCode = consts.RoleOperationBind
//			*relationType = consts.TrendsRelationTypeAddAfterIssue
//			*relationTypePeer = consts.TrendsRelationTypeAddAfterIssueByOther
//		}
//	}
//}

// 删除任务
func assemblyDeleteIssueTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := []bo.TrendsBo{}
	ext := bo.TrendExtensionBo{}
	ext.IssueType = "T"
	ext.ObjName = issueTrendsBo.IssueTitle
	trendsBos = append(trendsBos, bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationDelete,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.ParentIssueId,
		RelationType:    consts.TrendsRelationTypeDeleteIssue,
		RelationObjType: consts.TrendsOperObjectTypeIssue,
		NewValue:        nil,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(ext),
		Creator:         issueTrendsBo.OperatorId,
	})

	return trendsBos, nil
}

func assemblyUpdateRelationIssueTrendsInner(issueTrendsBo *bo.IssueTrendsBo,
	issueInfoById map[int64]*vo.Issue, columnId string, issueId int64,
	op, relationType, relationTypePeer string,
	trendsBos []bo.TrendsBo) []bo.TrendsBo {
	if _, ok := issueInfoById[issueId]; !ok {
		return trendsBos
	}

	//主动
	issueTrendsBo.Ext.RelationIssue = &bo.RelationIssue{
		ID:    issueInfoById[issueId].ID,
		Title: issueInfoById[issueId].Title,
	}
	trendsBos = append(trendsBos, bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        op,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		RelationObjId:   issueInfoById[issueId].ID,
		RelationType:    relationType,
		RelationObjType: columnId,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	})

	//被动
	if relationTypePeer != "" {
		issueTrendsBo.Ext.RelationIssue = &bo.RelationIssue{
			ID:    issueTrendsBo.IssueId,
			Title: issueInfoById[issueTrendsBo.IssueId].Title,
		}
		trendsBos = append(trendsBos, bo.TrendsBo{
			OrgId:           issueTrendsBo.OrgId,
			Module1:         consts.TrendsModuleOrg,
			Module2Id:       issueTrendsBo.ProjectId,
			Module2:         consts.TrendsModuleProject,
			Module3Id:       issueInfoById[issueId].ID,
			Module3:         consts.TrendsModuleIssue,
			OperCode:        op,
			OperObjId:       issueInfoById[issueId].ID,
			OperObjType:     consts.TrendsOperObjectTypeIssue,
			RelationObjId:   issueTrendsBo.IssueId,
			RelationType:    relationTypePeer,
			RelationObjType: consts.BlankString,
			Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
			Creator:         issueTrendsBo.OperatorId,
		})
	}

	return trendsBos
}

// 更新关联前后置
func assemblyUpdateRelationIssueTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := []bo.TrendsBo{}
	var allRelateIssueIds []int64
	allRelateIssueIds = append(allRelateIssueIds, issueTrendsBo.IssueId)
	for _, relatingChange := range issueTrendsBo.RelatingChange {
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkToAdd...)
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkToDel...)
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkFromAdd...)
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkFromDel...)
	}
	for _, relatingChange := range issueTrendsBo.BaRelatingChange {
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkToAdd...)
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkToDel...)
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkFromAdd...)
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkFromDel...)
	}
	for _, relatingChange := range issueTrendsBo.SingleRelatingChange {
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkToAdd...)
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkToDel...)
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkFromAdd...)
		allRelateIssueIds = append(allRelateIssueIds, relatingChange.LinkFromDel...)
	}
	allRelateIssueIds = slice.SliceUniqueInt64(allRelateIssueIds)
	if len(allRelateIssueIds) == 0 {
		return trendsBos, nil
	}

	issueInfoById := map[int64]*vo.Issue{}
	issueInfo := projectfacade.GetIssueInfoList(projectvo.IssueInfoListReqVo{
		OrgId:    issueTrendsBo.OrgId,
		UserId:   issueTrendsBo.OperatorId,
		IssueIds: allRelateIssueIds,
	})
	if issueInfo.Failure() {
		log.Error(issueInfo.Error())
		return nil, issueInfo.Error()
	}
	for i, v := range issueInfo.IssueInfos {
		issueInfoById[cast.ToInt64(v.ID)] = &issueInfo.IssueInfos[i]
	}

	// 关联
	for columnId, relatingChange := range issueTrendsBo.RelatingChange {
		columnType := tablePb.ColumnType_relating.String()

		op, relationType, relationTypePeer := assemblyRelatingOpRelationType(columnType, true, true)
		for _, id := range relatingChange.LinkToAdd {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}

		op, relationType, relationTypePeer = assemblyRelatingOpRelationType(columnType, true, false)
		for _, id := range relatingChange.LinkToDel {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}

		op, relationType, relationTypePeer = assemblyRelatingOpRelationType(columnType, false, true)
		for _, id := range relatingChange.LinkFromAdd {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}

		op, relationType, relationTypePeer = assemblyRelatingOpRelationType(columnType, false, false)
		for _, id := range relatingChange.LinkFromDel {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}
	}

	// 前后置
	for columnId, relatingChange := range issueTrendsBo.BaRelatingChange {
		columnType := tablePb.ColumnType_baRelating.String()

		op, relationType, relationTypePeer := assemblyRelatingOpRelationType(columnType, true, true)
		for _, id := range relatingChange.LinkToAdd {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}

		op, relationType, relationTypePeer = assemblyRelatingOpRelationType(columnType, true, false)
		for _, id := range relatingChange.LinkToDel {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}

		op, relationType, relationTypePeer = assemblyRelatingOpRelationType(columnType, false, true)
		for _, id := range relatingChange.LinkFromAdd {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}

		op, relationType, relationTypePeer = assemblyRelatingOpRelationType(columnType, false, false)
		for _, id := range relatingChange.LinkFromDel {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}
	}

	// 单向关联
	for columnId, relatingChange := range issueTrendsBo.SingleRelatingChange {
		columnType := tablePb.ColumnType_singleRelating.String()

		op, relationType, relationTypePeer := assemblyRelatingOpRelationType(columnType, true, true)
		for _, id := range relatingChange.LinkToAdd {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}

		op, relationType, relationTypePeer = assemblyRelatingOpRelationType(columnType, true, false)
		for _, id := range relatingChange.LinkToDel {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}

		op, relationType, relationTypePeer = assemblyRelatingOpRelationType(columnType, false, true)
		for _, id := range relatingChange.LinkFromAdd {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}

		op, relationType, relationTypePeer = assemblyRelatingOpRelationType(columnType, false, false)
		for _, id := range relatingChange.LinkFromDel {
			trendsBos = assemblyUpdateRelationIssueTrendsInner(&issueTrendsBo, issueInfoById, columnId, id, op, relationType, relationTypePeer, trendsBos)
		}
	}

	log.Infof("[assemblyUpdateRelationIssueTrends] type:%v issue:%v trends:%v", issueTrendsBo.PushType, issueTrendsBo.IssueId, json.ToJsonIgnoreError(trendsBos))
	return trendsBos, nil
}

// 引用附件
func assemblyReferResource(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBo := &bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeReferResource,
		RelationObjType: consts.BlankString,
		//NewValue:        &commentBo.Content,
		//OldValue:        nil,
		Ext:     json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator: issueTrendsBo.OperatorId,
	}

	return []bo.TrendsBo{*trendsBo}, nil
}

// 删除引用附件
func assemblyDeleteReferResource(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBo := &bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeDelReferResource,
		RelationObjType: consts.BlankString,
		//NewValue:        &commentBo.Content,
		//OldValue:        nil,
		Ext:     json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator: issueTrendsBo.OperatorId,
	}

	return []bo.TrendsBo{*trendsBo}, nil
}

// 变更任务栏
func assemblyUpdateIssueProjectTable(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBo := &bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeUpdateIssueProjectTable,
		RelationObjType: consts.BlankString,
		NewValue:        &issueTrendsBo.NewValue,
		OldValue:        &issueTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	}

	return []bo.TrendsBo{*trendsBo}, nil
}

// 变更父任务
func assemblyUpdateIssueParent(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	var trendsBos []bo.TrendsBo

	fromParentId := cast.ToInt64(issueTrendsBo.OldValue)
	targetParentId := cast.ToInt64(issueTrendsBo.NewValue)

	// 本任务动态
	trendsBos = append(trendsBos, bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationType:    consts.TrendsRelationTypeUpdateIssueParent,
		NewValue:        &issueTrendsBo.NewValue,
		OldValue:        &issueTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	})

	ext := bo.TrendExtensionBo{}
	ext.IssueType = "T"
	ext.ObjName = issueTrendsBo.IssueTitle
	issueIdStr := cast.ToString(issueTrendsBo.IssueId)

	// 原父任务
	if fromParentId > 0 {
		trendsBos = append(trendsBos, bo.TrendsBo{
			OrgId:           issueTrendsBo.OrgId,
			Module1:         consts.TrendsModuleOrg,
			Module2Id:       issueTrendsBo.ProjectId,
			Module2:         consts.TrendsModuleProject,
			Module3Id:       0,
			Module3:         consts.TrendsModuleIssue,
			OperCode:        consts.RoleOperationModify,
			OperObjId:       fromParentId,
			OperObjType:     consts.TrendsOperObjectTypeIssue,
			OperObjProperty: consts.BlankString,
			RelationType:    consts.TrendsRelationTypeDeleteChildIssue,
			RelationObjId:   fromParentId,
			RelationObjType: consts.TrendsOperObjectTypeIssue,
			NewValue:        nil,
			OldValue:        &issueIdStr,
			Ext:             json.ToJsonIgnoreError(ext),
			Creator:         issueTrendsBo.OperatorId,
		})
	}

	// 新父任务
	if targetParentId > 0 {
		trendsBos = append(trendsBos, bo.TrendsBo{
			OrgId:           issueTrendsBo.OrgId,
			Module1:         consts.TrendsModuleOrg,
			Module2Id:       issueTrendsBo.ProjectId,
			Module2:         consts.TrendsModuleProject,
			Module3Id:       0,
			Module3:         consts.TrendsModuleIssue,
			OperCode:        consts.RoleOperationModify,
			OperObjId:       targetParentId,
			OperObjType:     consts.TrendsOperObjectTypeIssue,
			OperObjProperty: consts.BlankString,
			RelationType:    consts.TrendsRelationTypeCreateChildIssue,
			RelationObjId:   targetParentId,
			RelationObjType: consts.TrendsOperObjectTypeIssue,
			NewValue:        &issueIdStr,
			OldValue:        nil,
			Ext:             json.ToJsonIgnoreError(ext),
			Creator:         issueTrendsBo.OperatorId,
		})
	}

	return trendsBos, nil
}

// 增加任务标签
func assemblyAddTag(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBo := &bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeAddIssueTag,
		RelationObjType: consts.BlankString,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	}

	return []bo.TrendsBo{*trendsBo}, nil
}

// 删除任务标签
func assemblyDeleteTag(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBo := &bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeDeleteIssueTag,
		RelationObjType: consts.BlankString,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	}

	return []bo.TrendsBo{*trendsBo}, nil
}

// 恢复任务
func assemblyRecoverIssueTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := []bo.TrendsBo{}
	ext := bo.TrendExtensionBo{}
	ext.IssueType = "T"
	ext.ObjName = issueTrendsBo.IssueTitle
	trendsBos = append(trendsBos, bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationCreate,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.ParentIssueId,
		RelationType:    consts.TrendsRelationTypeRecoverIssue,
		RelationObjType: consts.TrendsOperObjectTypeIssue,
		NewValue:        nil,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(ext),
		Creator:         issueTrendsBo.OperatorId,
	})

	return trendsBos, nil
}

// 审核任务
func assemblyAuditIssueTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	var trendsBos []bo.TrendsBo
	trendsBos = append(trendsBos, bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationCheck,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: issueTrendsBo.OperateObjProperty,
		RelationObjId:   issueTrendsBo.ParentIssueId,
		RelationType:    consts.TrendsRelationTypeAuditIssue,
		RelationObjType: consts.TrendsOperObjectTypeIssue,
		NewValue:        &issueTrendsBo.NewValue,
		OldValue:        &issueTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	})

	return trendsBos, nil
}

// 撤回任务
func assemblyWithdrawIssueTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	var trendsBos []bo.TrendsBo
	trendsBos = append(trendsBos, bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationWithdraw,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: issueTrendsBo.OperateObjProperty,
		RelationObjId:   issueTrendsBo.ParentIssueId,
		RelationType:    consts.TrendsRelationTypeWithdrawIssue,
		RelationObjType: consts.TrendsOperObjectTypeIssue,
		NewValue:        &issueTrendsBo.NewValue,
		OldValue:        &issueTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	})

	return trendsBos, nil
}

// 发起群聊讨论动态
func assemblyStartIssueChat(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := make([]bo.TrendsBo, 0)
	ext := issueTrendsBo.Ext
	ext.ObjName = issueTrendsBo.IssueTitle
	trendsBos = append(trendsBos, bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationStartIssueChat,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.ParentIssueId,
		RelationType:    consts.TrendsRelationTypeStartIssueChat,
		RelationObjType: consts.TrendsOperObjectTypeIssue,
		NewValue:        &issueTrendsBo.NewValue,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(ext),
		Creator:         issueTrendsBo.OperatorId,
	})

	return trendsBos, nil
}

func assemblyDocumentsTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	orgId := issueTrendsBo.OrgId
	issueId := issueTrendsBo.IssueId

	trendsBos := []bo.TrendsBo{}
	beforeChangeDocuments := issueTrendsBo.BeforeChangeDocuments
	afterChangeDocuments := issueTrendsBo.AfterChangeDocuments

	log.Infof("[assemblyDocumentsTrends], orgId:%v, issueId:%v, beforeChangeDocuments:%v, afterChangeDocuments:%v",
		orgId, issueId, beforeChangeDocuments, afterChangeDocuments)

	newDocumentsMap := make(map[string]lc_table.LcDocumentValue)
	oldDocumentsMap := make(map[string]lc_table.LcDocumentValue)

	newResourceIds := []int64{}
	oldResourceIds := []int64{}

	for _, v := range beforeChangeDocuments {
		oldDocuments := make(map[string]lc_table.LcDocumentValue)
		err := copyer.Copy(v, &oldDocuments)
		if err != nil {
			return trendsBos, errs.ObjectCopyError
		}
		for k, m := range oldDocuments {
			oldResourceIds = append(oldResourceIds, m.Id)
			oldDocumentsMap[k] = m
		}
	}
	for _, v := range afterChangeDocuments {
		newDocuments := make(map[string]lc_table.LcDocumentValue)
		err := copyer.Copy(v, &newDocuments)
		if err != nil {
			return trendsBos, errs.ObjectCopyError
		}
		for k, m := range newDocuments {
			newResourceIds = append(newResourceIds, m.Id)
			newDocumentsMap[k] = m
		}
	}

	// new和old全为nil，或者两个len相同，就跳过
	if (beforeChangeDocuments == nil && afterChangeDocuments == nil) || (len(newDocumentsMap) == len(oldDocumentsMap)) {
		return trendsBos, nil
	}
	resourceInfos := []bo.ResourceInfoBo{}

	ext := issueTrendsBo.Ext
	ext.ObjName = issueTrendsBo.IssueTitle

	relationType := ""
	// 对比新旧的ids
	_, addIds, delIds := int64Slice.CompareSliceAddDelInt64(newResourceIds, oldResourceIds)
	if len(addIds) > 0 {
		// 上传附件
		relationType = consts.TrendsRelationTypeUploadResource
		resourceList, err := getResourceBoList(orgId, addIds)
		if err != nil {
			log.Errorf("[assemblyDocumentsTrends]getResourceBoList err:%v", err)
			return nil, err
		}
		resourceInfos = resourceList
		log.Infof("[assemblyDocumentsTrends] orgId:%v, issueId:%v, 要上传的附件信息: %v", orgId, issueId, resourceInfos)
	}

	if len(delIds) > 0 {
		// 删除附件
		relationType = consts.TrendsRelationTypeDeleteResource
		resourceList, err := getResourceBoList(orgId, delIds)
		if err != nil {
			log.Errorf("[assemblyDocumentsTrends]getResourceBoList err:%v", err)
			return nil, err
		}
		resourceInfos = resourceList

		log.Infof("[assemblyDocumentsTrends] orgId:%v, issueId:%v, 要删除的附件信息: %v", orgId, issueId, resourceInfos)
	}

	ext.ResourceInfo = resourceInfos
	trendsBos = append(trendsBos, bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    relationType,
		RelationObjType: consts.BlankString,
		NewValue:        nil,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(ext),
		Creator:         issueTrendsBo.OperatorId,
	})
	return trendsBos, nil
}

// 暂时做成原来的图片动态展示的格式: 修改xx为xxx
func assemblyPicturesTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := []bo.TrendsBo{}
	beforeChangeImages := issueTrendsBo.BeforeChangeImages
	afterChangeImages := issueTrendsBo.AfterChangeImages

	log.Infof("[assemblyPicturesTrends], orgId:%v, issueId:%v, beforeChangeImages:%v, afterChangeImages:%v",
		issueTrendsBo.OrgId, issueTrendsBo.IssueId, beforeChangeImages, afterChangeImages)

	newImages := make([]lc_table.LcDocumentValue, 0)
	oldImages := make([]lc_table.LcDocumentValue, 0)

	columnId := ""

	for col, v := range beforeChangeImages {
		err := copyer.Copy(v, &oldImages)
		if err != nil {
			return trendsBos, errs.ObjectCopyError
		}
		columnId = col
	}
	for _, v := range afterChangeImages {
		err := copyer.Copy(v, &newImages)
		if err != nil {
			return trendsBos, errs.ObjectCopyError
		}
	}

	if (beforeChangeImages == nil && afterChangeImages == nil) || len(newImages) == len(oldImages) {
		return trendsBos, nil
	}

	// 获取新老数据的name
	newValue := ""
	oldValue := ""
	for _, v := range oldImages {
		oldValue += v.Name + " "
	}
	for _, v := range newImages {
		newValue += v.Name + " "
	}
	changeList := []bo.TrendChangeListBo{}
	changeList = append(changeList, bo.TrendChangeListBo{
		Field:     columnId,
		FieldType: consts.LcColumnFieldTypeImage,
		FieldName: "图片",
		OldValue:  oldValue,
		NewValue:  newValue,
	})

	issueTrendsBo.Ext.ChangeList = changeList

	trendsBos = append(trendsBos, bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: issueTrendsBo.OperateObjProperty,
		RelationObjId:   0,
		RelationType:    consts.TrendsRelationTypeUpdateIssue,
		RelationObjType: consts.BlankString,
		//NewValue:        &newValue,
		//OldValue:        &oldValue,
		Ext:     json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator: issueTrendsBo.OperatorId,
	})

	return trendsBos, nil
}

func getResourceBoList(orgId int64, resourceIds []int64) ([]bo.ResourceInfoBo, errs.SystemErrorInfo) {
	resourceInfos := []bo.ResourceInfoBo{}
	isDelete := consts.AppIsNoDelete
	resourceBoList := resourcefacade.GetResourceBoList(resourcevo.GetResourceBoListReqVo{
		Page: 0,
		Size: 0,
		Input: resourcevo.GetResourceBoListCond{
			OrgId:       orgId,
			ResourceIds: &resourceIds,
			IsDelete:    &isDelete,
		},
	})
	if resourceBoList.Failure() {
		log.Errorf("[getResourceBoList] resourcefacade.GetResourceBoList failed, err:%v, orgId:%d, resourceIds:%v",
			resourceBoList.Failure(), orgId, resourceIds)
		return nil, resourceBoList.Error()
	}
	for _, resource := range *resourceBoList.ResourceBos {
		resourceInfos = append(resourceInfos, bo.ResourceInfoBo{
			Id:         resource.Id,
			Url:        resource.Host + resource.Path,
			Name:       resource.Name,
			Size:       resource.Size,
			UploadTime: resource.CreateTime,
			Suffix:     resource.Suffix,
		})
	}

	return resourceInfos, nil
}
