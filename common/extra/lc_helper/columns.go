package lc_helper

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
)

//无码中的字段组件封装
//参考文档 https://wiki.bjx.cloud/pages/viewpage.action?pageId=32145423
//组件都以 `Ct` 作为前缀，表示 Column Type。

// GetLcCtProjectType “项目类型”的字段类型
// 和@千源 约定的特殊定制，“项目类型”字段，其值是项目 id。
func GetLcCtProjectType(columnName, label string) lc_table.LcCtInput {
	return lc_table.LcCtInput{
		Name:   columnName,
		Label:  label,
		Unique: true,
		Field: lc_table.LcCtInputField{
			Type: "project",
			Props: lc_table.LcCtInputFieldProps{
				Required: true,
				IsSearch: true,
				FieldSearch: lc_table.LcCtInputFieldPropsFieldSearch{
					Type: "formSelect",
					Sort: 1,
				},
			},
		},
	}
}

// 单行文本
func GetLcCtInput(columnName, label string) lc_table.LcCtInput {
	return lc_table.LcCtInput{
		Name:   columnName,
		Label:  label,
		Unique: true,
		Field: lc_table.LcCtInputField{
			Type: "input",
			Props: lc_table.LcCtInputFieldProps{
				Required: true,
				IsSearch: true,
				FieldSearch: lc_table.LcCtInputFieldPropsFieldSearch{
					Type: "formSelect",
					Sort: 1,
				},
			},
		},
	}
}

func GetLcCtTextArea(columnName, label, enLabel string, editable bool, writable, hide, pushMsg bool) lc_table.LcCtTextarea {
	return lc_table.LcCtTextarea{
		Name:     columnName,
		Label:    label,
		EnLabel:  enLabel,
		Editable: editable,
		Writable: writable,
		Field: lc_table.LcCtTextareaField{
			Type: "textarea",
			Props: lc_table.LcCtTextareaFieldProps{
				Hide:    hide,
				PushMsg: pushMsg,
			},
		},
	}
}

func GetLcCtInputFull(columnName, label string, enLabel string, unique bool, required bool, editable bool, writable bool, disabled bool, hide bool, isPushMsg bool) lc_table.LcCtInput {
	return lc_table.LcCtInput{
		Name:     columnName,
		Label:    label,
		EnLabel:  enLabel,
		Unique:   unique,
		Editable: editable,
		Writable: writable,
		Field: lc_table.LcCtInputField{
			Type: "input",
			Props: lc_table.LcCtInputFieldProps{
				Required: required,
				IsSearch: true,
				Hide:     hide,
				Disabled: disabled,
				PushMsg:  isPushMsg,
				FieldSearch: lc_table.LcCtInputFieldPropsFieldSearch{
					Type: "formSelect",
					Sort: 1,
				},
			},
		},
	}
}

func GetLcCtSelect(columnName, label string, enLabel string, customType string, options []lc_table.LcCtSelectFieldPropsSelectOptions,
	editable bool, writable bool, isText bool, disabled bool, required bool,
) lc_table.LcCtSelect {
	var defaultValue interface{}
	if len(options) > 0 {
		defaultValue = options[0].Id
	}
	return lc_table.LcCtSelect{
		Name:     columnName,
		Label:    label,
		EnLabel:  enLabel,
		Editable: editable,
		Writable: writable,
		Field: lc_table.LcCtSelectField{
			Type:       "select",
			CustomType: customType,
			Props: lc_table.LcCtSelectFieldProps{
				IsText:   isText,
				Disabled: disabled,
				Required: required,
				Select: lc_table.LcCtSelectFieldPropsSelect{
					Options: options,
				},
				Default: defaultValue,
			},
		},
	}
}

// GetLcCtGroupSelect 分组单选 groupSelect
func GetLcCtGroupSelect(columnKey, label string, customType string, groupSelect lc_table.LcPropGroupSelect,
	require, isPushMsg bool) lc_table.LcOneColumn {
	var defaultValue interface{}
	if len(groupSelect.GroupOptions) > 0 {
		if len(groupSelect.GroupOptions[0].Children) > 0 {
			defaultValue = groupSelect.GroupOptions[0].Children[0].Id
		}
	}
	return lc_table.LcOneColumn{
		Name:       columnKey,
		Label:      label,
		Key:        columnKey,
		Title:      label,
		AliasTitle: "",
		EnTitle:    "",
		Writable:   true,
		Editable:   true,
		Field: lc_table.LcFieldData{
			Type:       "groupSelect",
			CustomType: customType,
			Props: lc_table.LcProps{
				GroupSelect: &groupSelect,
				Required:    require,
				PushMsg:     isPushMsg,
				Default:     defaultValue,
			},
		},
	}
}

func GetLcCtDatepicker(columnName, label string, enLabel string, editable, isPushMsg bool) lc_table.LcCtDatepicker {
	return lc_table.LcCtDatepicker{
		Name:     columnName,
		Label:    label,
		Editable: editable,
		Writable: true,
		Field: lc_table.LcCtDatepickerField{
			Type:  "datepicker",
			Props: lc_table.LcCtDatepickerFieldProps{PushMsg: isPushMsg},
		},
	}
}

func GetLcCtMember(columnName, label string, enLabel string, editable bool, writable bool, multiple bool, limit int, hasDefaultCollaboratorRoles, isPushMsg bool) lc_table.LcCtMember {
	data := lc_table.LcCtMember{
		Name:     columnName,
		Label:    label,
		EnLabel:  enLabel,
		Editable: editable,
		Writable: writable,
		Field: lc_table.LcCtMemberField{
			Type: "member",
			Props: lc_table.LcCtMemberFieldProps{
				Multiple: multiple,
				PushMsg:  isPushMsg,
			},
		},
	}
	if limit > 0 {
		data.Field.Props.Limit = &limit
	}

	if hasDefaultCollaboratorRoles {
		data.Field.Props.CollaboratorRoles = &[]string{"-1"}
	}

	return data
}

func GetLcCtRelateTable(columnName, label string, relateAppId string, displayField string, displayColumns []string) lc_table.LcCtRelateTable {
	return lc_table.LcCtRelateTable{
		Name:  columnName,
		Label: label,
		Field: lc_table.LcCtRelateTableField{
			Type: "relateTable",
			Props: lc_table.LcCtRelateTableFieldProps{
				RelateTable: lc_table.LcCtRelateTableFieldPropsRelateTable{
					AppID:          relateAppId,
					IDField:        "id",
					DisplayField:   displayField,
					DisplayColumns: displayColumns,
				},
				ShowDetails: true,
				FormOrder:   0,
				Required:    true,
				TabParam:    "",
				GroupParam:  "",
				IsSearch:    true,
				FieldSearch: lc_table.LcCtRelateTableFieldPropsFieldSearch{
					Type: "formSelect",
					Sort: 1,
				},
			},
		},
	}
}

func GetDocumentColumn() interface{} {
	c := map[string]interface{}{}
	json.Unmarshal([]byte(`{
	"name": "document", 
	"field": {
		"type": "document", 
		"customType":"",
		"dataType": "CUSTOM",
		"props": {
			"checked": true,
			"disabled": false, 
			"hide": false,
			"multiple": false, 
			"required": false
		}
	}, 
	"label": "附件", 
	"editable": true, 
	"writable": true
}`), &c)

	return c
}

func GetSelectColumn() interface{} {
	c := map[string]interface{}{}
	json.Unmarshal([]byte(`{"name":"_field_select","label":"单选","aliasTitle":"","description":"","isSys":false,"isOrg":false,"writable":true,"editable":true,"unique":false,"uniquePreHandler":"","sensitiveStrategy":"","sensitiveFlag":0,"field":{"type":"select","customType":"","dataType":"STRING","props":{"required":false,"select":{"options":[]},"collaboratorRoles":null,"disabled":false,"multiple":false,"options":[],"pushMsg":false},"refSetting":null}}`), &c)

	return c
}

func GetMultiSelectColumn() interface{} {
	c := map[string]interface{}{}
	json.Unmarshal([]byte(`{"name":"_field_multi_select","label":"多选","aliasTitle":"","description":"","isSys":false,"isOrg":false,"writable":true,"editable":true,"unique":false,"uniquePreHandler":"","sensitiveStrategy":"","sensitiveFlag":0,"field":{"type":"multiselect","customType":"","dataType":"STRING","props":{"collaboratorRoles":null,"disabled":false,"multiple":false,"multiselect":{"options":[]},"options":[],"pushMsg":false,"required":false},"refSetting":null}}`), &c)

	return c
}

// GetOrgPriorityField 组织字段
func GetOrgPriorityField() lc_table.LcCommonField {
	editableFlag := true
	priorityOptions := []lc_table.LcOptions{
		{
			FontColor: "#FFFFFF",
			Id:        1,
			Value:     "最高",
			Color:     "#FF5037",
		},
		{
			FontColor: "#FFFFFF",
			Id:        2,
			Value:     "较高",
			Color:     "#FFC700",
		},
		{
			FontColor: "#FFFFFF",
			Id:        3,
			Value:     "普通",
			Color:     "#67D287",
		},
		{
			FontColor: "#FFFFFF",
			Id:        4,
			Value:     "较低",
			Color:     "#5991FF",
		},
		{
			FontColor: "#FFFFFF",
			Id:        5,
			Value:     "最低",
			Color:     "#CACACA",
		},
	}

	return lc_table.LcCommonField{
		Label:    "优先级",
		Name:     consts.BasicFieldPriority,
		Writable: true,
		IsOrg:    true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "select",
			CustomType: "select",
			Props: lc_table.LcProps{
				Select: lc_table.LcPropSelect{
					Options: priorityOptions,
				},
			},
		},
	}
}
