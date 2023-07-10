package lc_helper

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
)

func NewLcTableConfig(fields []interface{}) lc_table.TableConfig {
	return lc_table.TableConfig{
		Fields: fields,
	}
}

// GetIssueStatusList 汇总表-任务状态：option list
func GetIssueStatusList() []lc_table.LcCtSelectFieldPropsSelectOptions {
	var StatusMap = consts.LcIssueStatusMap
	list := make([]lc_table.LcCtSelectFieldPropsSelectOptions, 0)
	for k, v := range StatusMap {
		list = append(list, lc_table.LcCtSelectFieldPropsSelectOptions{
			Id:    k,
			Color: "",
			Value: v,
		})
	}
	return list
}

// GetIssuePriorityList 汇总表-优先级：option list
func GetIssuePriorityList() []lc_table.LcCtSelectFieldPropsSelectOptions {
	var StatusMap = consts.LcIssuePriorityMap
	list := make([]lc_table.LcCtSelectFieldPropsSelectOptions, 0)
	for k, v := range StatusMap {
		list = append(list, lc_table.LcCtSelectFieldPropsSelectOptions{
			Id:    k,
			Color: "",
			Value: v,
		})
	}
	return list
}
