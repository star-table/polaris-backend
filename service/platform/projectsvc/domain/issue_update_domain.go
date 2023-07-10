package domain

import (
	"strings"

	"github.com/star-table/polaris-backend/common/model/vo/lc_table"

	pushV1 "github.com/star-table/interface/golang/push/v1"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/model/bo"
)

// CheckHasUpdateIssueForGroupChat 更新任务时，检查更新了哪些内容 todo
// 参考：service/platform/trendssvc/domain/trends_issue_domain.go@assemblyTrendsBos
//func CheckHasUpdateIssueForGroupChat(issueTrendsBo *bo.IssueTrendsBo) ([]string, errs.SystemErrorInfo) {
//	updateColumns := make([]string, 0)
//	if columns := GetUpdateIssueCustomColumns(issueTrendsBo); len(columns) > 0 {
//		updateColumns = append(updateColumns, columns...)
//	} else if columns := GetUpdateIssueMemberColumns(issueTrendsBo); len(columns) > 0 {
//		updateColumns = append(updateColumns, columns...)
//	} else if columns := GetUpdateIssueFileColumns(issueTrendsBo); len(columns) > 0 {
//		updateColumns = append(updateColumns, columns...)
//	} else if columns := GetUpdateIssueRelatingColumns(issueTrendsBo); len(columns) > 0 {
//		updateColumns = append(updateColumns, columns...)
//	} else if true {
//		// 工时字段判断 todo
//	}
//
//	return updateColumns, nil
//}

// GetUpdateIssueCustomColumns 群聊推送，根据配置检查是否需要推送，自定义列的检查
//func GetUpdateIssueCustomColumns(issueTrendsBo *bo.IssueTrendsBo) []string {
//	changedColumns := make([]string, 0)
//	if len(issueTrendsBo.Ext.ChangeList) == 0 {
//		return changedColumns
//	}
//	for _, changedItem := range issueTrendsBo.Ext.ChangeList {
//		changedColumns = append(changedColumns, changedItem.FieldName)
//	}
//
//	return changedColumns
//}
//
//// GetUpdateIssueMemberColumns 群聊推送，根据配置检查是否需要推送，成员是否更新
//func GetUpdateIssueMemberColumns(issueTrendsBo *bo.IssueTrendsBo) []string {
//	updateColumns := make([]string, 0)
//	if !int64Slice.CompareSliceInt64(issueTrendsBo.BeforeOwner, issueTrendsBo.AfterOwner) {
//		updateColumns = append(updateColumns, consts.BasicFieldOwnerId)
//	}
//	if !int64Slice.CompareSliceInt64(issueTrendsBo.BeforeChangeAuditors, issueTrendsBo.AfterChangeAuditors) {
//		updateColumns = append(updateColumns, consts.BasicFieldAuditorIds)
//	}
//	if !int64Slice.CompareSliceInt64(issueTrendsBo.BeforeChangeFollowers, issueTrendsBo.AfterChangeFollowers) {
//		updateColumns = append(updateColumns, consts.BasicFieldFollowerIds)
//	}
//
//	return updateColumns
//}
//
//// GetUpdateIssueFileColumns 群聊推送，根据配置检查是否需要推送，文件是否更新（附件、图片）
//func GetUpdateIssueFileColumns(issueTrendsBo *bo.IssueTrendsBo) []string {
//	updateColumns := make([]string, 0)
//	for column, beforeVal := range issueTrendsBo.BeforeChangeDocuments {
//		if afterVal, ok := issueTrendsBo.AfterChangeDocuments[column]; ok {
//			if len(afterVal.(map[string][]interface{})) != len(beforeVal.(map[string][]interface{})) {
//				updateColumns = append(updateColumns, column)
//			}
//		} else {
//			updateColumns = append(updateColumns, column)
//		}
//	}
//
//	for column, beforeVal := range issueTrendsBo.BeforeChangeImages {
//		if afterVal, ok := issueTrendsBo.AfterChangeImages[column]; ok {
//			if len(afterVal.([]interface{})) != len(beforeVal.([]interface{})) {
//				updateColumns = append(updateColumns, column)
//			}
//		} else {
//			updateColumns = append(updateColumns, column)
//		}
//	}
//
//	return updateColumns
//}

// GetUpdateIssueRelatingColumns 群聊推送，根据配置检查是否需要推送，关联字段是否更新（前后置、关联）
//func GetUpdateIssueRelatingColumns(issueTrendsBo *bo.IssueTrendsBo) []string {
//	updateColumns := make([]string, 0)
//	for key, beforeVal := range issueTrendsBo.BeforeChangeRelating {
//		if afterVal, ok := issueTrendsBo.AfterChangeRelating[key]; ok {
//			if len(afterVal) != len(beforeVal) {
//				updateColumns = append(updateColumns, consts.BasicFieldRelating)
//			}
//		} else {
//			updateColumns = append(updateColumns, consts.BasicFieldRelating)
//		}
//	}
//
//	for key, beforeVal := range issueTrendsBo.BeforeChangeBaRelating {
//		if afterVal, ok := issueTrendsBo.AfterChangeBaRelating[key]; ok {
//			if len(afterVal) != len(beforeVal) {
//				updateColumns = append(updateColumns, consts.BasicFieldBaRelating)
//			}
//		} else {
//			updateColumns = append(updateColumns, consts.BasicFieldBaRelating)
//		}
//	}
//
//	return updateColumns
//}

// AssemblyCardForUpdateIssueOtherColumn 拼装卡片内容，包括标题、起止时间和一些自定义字段 AssemblyFsCardForUpdateIssueWorkHour
func AssemblyCardForUpdateIssueOtherColumn(headers map[string]lc_table.LcCommonField, changeItem *bo.TrendChangeListBo, cardMeta *pushV1.TemplateCard) {
	curColumn, ok := headers[changeItem.Field]
	if !ok {
		//log.Infof("[AssemblyCardForUpdateIssueOtherColumn] could not match column. columnKey: %s", changeItem.Field)
		return
	}
	displayName := curColumn.AliasTitle
	if displayName == "" {
		displayName = curColumn.Label
	}
	if str.RuneLen(displayName) > consts.ColumnAliasLimit {
		displayName = str.Substr(displayName, 0, consts.ColumnAliasLimit) + consts.CardIssueColumnNameOverflow
	}

	oldValStr := consts.CardDefaultOwnerNameForUpdateIssue
	newValStr := consts.CardDefaultOwnerNameForUpdateIssue
	if strings.Trim(changeItem.OldValue, " \n\t") != "" {
		oldValStr = changeItem.OldValue
	}
	if strings.Trim(changeItem.NewValue, " \n\t") != "" {
		newValStr = changeItem.NewValue
	}

	// 时间去除默认的空时间
	if changeItem.FieldType == consts.LcColumnFieldTypeDatepicker {
		if oldValStr == consts.BlankTime || oldValStr == consts.BlankEmptyTime {
			oldValStr = consts.CardDefaultOwnerNameForUpdateIssue
		}
		if newValStr == consts.BlankTime || newValStr == consts.BlankEmptyTime {
			newValStr = consts.CardDefaultOwnerNameForUpdateIssue
		}
	}
	// 富文本去除标签
	if changeItem.FieldType == consts.LcColumnFieldTypeRichText {
		oldValStr = str.TrimHtml(oldValStr)
		newValStr = str.TrimHtml(newValStr)
		// 去重空格和空行
		oldValStr = str.ReplaceWhiteSpaceCharToSpace(oldValStr)
		newValStr = str.ReplaceWhiteSpaceCharToSpace(newValStr)
		if strings.Trim(oldValStr, " \n\t") == "" {
			oldValStr = consts.CardIssueChangeDescForNotSupportShow
		}
		if strings.Trim(newValStr, " \n\t") == "" {
			newValStr = consts.CardIssueChangeDescForNotSupportShow
		}
	}

	markdownTile := RenderChangeDescForUpdateIssueInCard(oldValStr, newValStr)
	cardMeta.Divs = append(cardMeta.Divs, &pushV1.CardTextDiv{
		Fields: []*pushV1.CardTextField{
			&pushV1.CardTextField{
				Key:   displayName,
				Value: markdownTile,
			},
		},
	})
	//cardElements.MarkdownElements = append(cardElements.MarkdownElements, markdownTile)
	//cardElements.FsElements = append(cardElements.FsElements, vo.CardElementContentModule{
	//	Tag: "div",
	//	Fields: []vo.CardElementField{
	//		{
	//			Text: vo.CardElementText{
	//				Tag:     "lark_md",
	//				Content: markdownTile,
	//			},
	//		},
	//	},
	//})
}
