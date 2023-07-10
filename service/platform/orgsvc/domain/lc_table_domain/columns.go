package lc_table_domain

//无码中的字段组件封装
//参考文档 https://wiki.bjx.cloud/pages/viewpage.action?pageId=32145423
//组件都以 `Ct` 作为前缀，表示 Column Type。

/*
// 已迁移到 common/extra/lc_helper
// “项目类型”的字段类型
// 和@千源 约定的特殊定制，“项目类型”字段，其值是项目 id。
func GetLcCtProjectType(columnName, label string) lc_table.LcCtInput {
	return lc_table.LcCtInput{
		Name:   columnName,
		Label:  label,
		Unique: true,
		Field: lc_table.LcCtInputField{
			Type:  "project",
			Props: lc_table.LcCtInputFieldProps{
				Required:    true,
				IsSearch:    true,
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
			Type:  "input",
			Props: lc_table.LcCtInputFieldProps{
				Required:    true,
				IsSearch:    true,
				FieldSearch: lc_table.LcCtInputFieldPropsFieldSearch{
					Type: "formSelect",
					Sort: 1,
				},
			},
		},
	}
}

func GetLcCtInputFull(columnName, label string, enLabel string, unique bool, required bool, editable bool, disabled bool, hide bool) lc_table.LcCtInput {
	return lc_table.LcCtInput{
		Name:   columnName,
		Label:  label,
		EnLabel: enLabel,
		Unique: unique,
		Editable: editable,
		Field: lc_table.LcCtInputField{
			Type:  "input",
			Props: lc_table.LcCtInputFieldProps{
				Required:    required,
				IsSearch:    true,
				Hide: hide,
				Disabled: disabled,
				FieldSearch: lc_table.LcCtInputFieldPropsFieldSearch{
					Type: "formSelect",
					Sort: 1,
				},
			},
		},
	}
}

func GetLcCtSelect(columnName, label string, enLabel string, customType string, options []lc_table.LcCtSelectFieldPropsSelectOptions, editable bool, isText bool, disabled bool, required bool) lc_table.LcCtSelect {
	return lc_table.LcCtSelect{
		Name:  columnName,
		Label: label,
		EnLabel: enLabel,
		Editable: editable,
		Field: lc_table.LcCtSelectField{
			Type:  "select",
			CustomType:customType,
			Props: lc_table.LcCtSelectFieldProps{
				IsText: isText,
				Disabled: disabled,
				Required: required,
				Select: lc_table.LcCtSelectFieldPropsSelect{
					Options: options,
				},
			},
		},
	}
}

func GetLcCtGroupSelect(columnName, label string, customType string, options []lc_table.LcCtSelectFieldPropsSelectOptions, required, editable bool, isText bool, disabled bool, titleDisabled, typeDisabled bool) lc_table.LcCtSelect {
	return lc_table.LcCtSelect{
		Name:  columnName,
		Label: label,
		Editable: editable,
		Field: lc_table.LcCtSelectField{
			Type:  "groupSelect",
			CustomType:customType,
			Props: lc_table.LcCtSelectFieldProps{
				IsText: isText,
				Disabled: disabled,
				TitleDisabled: titleDisabled,
				TypeDisabled: typeDisabled,
				Required: required,
				Select: lc_table.LcCtSelectFieldPropsSelect{
					Options: options,
				},
			},
		},
	}
}

func GetLcCtDatepicker(columnName, label string, enLabel string, editable bool) lc_table.LcCtDatepicker {
	return lc_table.LcCtDatepicker{
		Name:  columnName,
		Label: label,
		Editable: editable,
		Field: lc_table.LcCtDatepickerField{
			Type:  "datepicker",
			Props: lc_table.LcCtDatepickerFieldProps{},
		},
	}
}

func GetLcCtMember(columnName, label string, enLabel string, editable bool, multiple bool, limit int, hasDefaultCollaboratorRoles bool) lc_table.LcCtMember {
	data := lc_table.LcCtMember{
		Name:  columnName,
		Label: label,
		EnLabel: enLabel,
		Editable: editable,
		Field: lc_table.LcCtMemberField{
			Type:  "member",
			Props: lc_table.LcCtMemberFieldProps{
				Multiple:multiple,
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
			Type:  "relateTable",
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
*/
