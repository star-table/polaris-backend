package lc_helper

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
)

// GetDefaultGroupSelectForIssueStatus 任务状态-默认选项配置
func GetDefaultGroupSelectForIssueStatus() lc_table.LcPropGroupSelect {
	options := make([]lc_table.LcGroupOptionsDetail, 0, 3)
	groups := make([]lc_table.LcGroupOptions, 0, 3)
	json.FromJson(DefaultGroupSelectJsonForIssueStatus, &groups)
	for _, group := range groups {
		for _, item := range group.Children {
			options = append(options, lc_table.LcGroupOptionsDetail{
				Id:        item.Id,
				Value:     item.Value,
				ParentId:  item.ParentId,
				Color:     item.Color,
				FontColor: item.FontColor,
			})
		}
	}

	return lc_table.LcPropGroupSelect{
		GroupOptions: groups,
		Options:      options,
	}
}

// GetDefaultSelectOptionsForTaskBar 任务栏-单选-默认选项值
func GetDefaultSelectOptionsForTaskBar() []lc_table.LcCtSelectFieldPropsSelectOptions {
	options := make([]lc_table.LcCtSelectFieldPropsSelectOptions, 0)
	json.FromJson(DefaultSelectJsonForTaskBar, &options)

	return options
}

// GetDefaultSelectOptionsForIterationId 迭代-单选-默认选项值
func GetDefaultSelectOptionsForIterationId() []lc_table.LcCtSelectFieldPropsSelectOptions {
	options := make([]lc_table.LcCtSelectFieldPropsSelectOptions, 0)
	json.FromJson(DefaultSelectJsonForIterationId, &options)

	return options
}
