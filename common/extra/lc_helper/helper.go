package lc_helper

import (
	"fmt"
	"strings"

	"github.com/star-table/polaris-backend/common/model/vo/projectvo"

	"github.com/star-table/polaris-backend/common/core/consts"
)

var (
	NotJsonColumnMap = map[string]struct{}{
		consts.BasicFieldId:            {},
		consts.BasicFieldOrgId:         {},
		consts.BasicFieldRecycleFlag:   {},
		consts.BasicFieldCreator:       {},
		consts.BasicFieldUpdator:       {},
		consts.BasicFieldCreateTime:    {},
		consts.BasicFieldUpdateTime:    {},
		consts.ProBasicFieldAppId:      {},
		consts.BasicFieldProjectId:     {},
		consts.BasicFieldTableId:       {},
		consts.BasicFieldPath:          {},
		consts.BasicFieldParentId:      {},
		consts.BasicFieldIssueId:       {},
		consts.BasicFieldCollaborators: {},
		consts.BasicFieldOrder:         {},
		consts.BasicFieldCode:          {},
		consts.LcJsonColumn:            {},
	}
)

// 字段名转换为查询无码的SELECT列名
func ConvertToFilterColumnV2(c string, columns map[string]*projectvo.TableColumnData) string {
	switch c {
	case consts.BasicFieldId:
		return c
	case consts.BasicFieldTitle:
		return fmt.Sprintf("\"data\"::jsonb->>'%s' \"%s\"", c, c)
	case consts.BasicFieldRemark:
		return fmt.Sprintf("\"data\"::jsonb->>'%s' \"%s\"", c, c)
	default:
		//if column, ok := columns[c]; ok &&
		//	(column.Field.Type == consts.LcColumnFieldTypeMember || column.Field.Type == consts.LcColumnFieldTypeDept) {
		//	return fmt.Sprintf("case when \"data\"::jsonb->>'%s' is null then null else \"data\"::jsonb->'%s' end \"%s\"", c, c, c)
		//} else {
		//	return fmt.Sprintf("\"data\"::jsonb->'%s' \"%s\"", c, c)
		//}
		if column, ok := columns[c]; ok {
			if column.Field.Type == consts.LcColumnFieldTypeTextarea ||
				column.Field.Type == consts.LcColumnFieldTypeDatepicker {
				return fmt.Sprintf("\"data\"::jsonb->>'%s' \"%s\"", c, c)
			}
		}
		return fmt.Sprintf("case when \"data\"::jsonb->>'%s' is null then null else \"data\"::jsonb->'%s' end \"%s\"", c, c, c)
	}
}

// 字段名转换为查询无码的SELECT列名
func ConvertToFilterColumn(c string) string {
	if _, ok := NotJsonColumnMap[c]; ok {
		return "\"" + c + "\""
	}

	switch c {
	case consts.BasicFieldId, consts.BasicFieldOrgId, consts.BasicFieldRecycleFlag:
		return c
	default:
		return fmt.Sprintf("\"data\" :: jsonb -> '%s' \"%s\"", c, c)
	}
}

//func ConvertToJsonColumn(column *projectvo.TableColumnData) string {
//	if column.Field.Type == consts.LcColumnFieldTypeInput || column.Field.Type == consts.LcColumnFieldTypeRichText ||
//		column.Field.Type == consts.LcColumnFieldTypeEmail || column.Field.Type == consts.LcColumnFieldTypeLink ||
//		column.Field.Type == consts.LcColumnFieldTypeMobile || column.Field.Type == consts.LcColumnFieldTypeTextarea {
//
//		return fmt.Sprintf("\"data\" :: jsonb ->> '%s' \"%s\"", column.Name, column.Name)
//	}
//
//	return fmt.Sprintf("\"data\" :: jsonb -> '%s' \"%s\"", column.Name, column.Name)
//}

// 字段名转换为查询无码的WHERE条件列名
func ConvertToCondColumn(c string) string {
	if _, ok := NotJsonColumnMap[c]; ok {
		return "\"" + c + "\""
	}

	if strings.Contains(c, consts.LcJsonColumn) {
		return c
	}

	switch c {
	case consts.BasicFieldId:
		return c
	default:
		return fmt.Sprintf("\"data\" :: jsonb -> '%s'", c)
	}
}
