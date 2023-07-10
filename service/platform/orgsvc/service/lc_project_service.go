package service

import (
	"fmt"
	"sync"

	"github.com/star-table/polaris-backend/facade/tablefacade"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"
	"github.com/star-table/polaris-backend/facade/appfacade"
)

var wg sync.WaitGroup

// 项目数据同步到无码- form 字段组装
func GetProjectFormFields(orgId int64) ([]interface{}, []string, errs.SystemErrorInfo) {
	fields := make([]interface{}, 0)
	otherFields := make([]string, 0)
	addField := ""

	editableFlag := true
	addField = consts.ProBasicFieldName
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "input",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldAppId
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "inputnumber",
			CustomType: "inputnumber",
			Props: lc_table.LcProps{
				InputNumber: lc_table.LcPropInputNumber{
					Accuracy:   "1",
					Thousandth: false,
					Percentage: false,
					Unique:     false,
					Required:   false,
				},
			},
		},
	})

	addField = consts.ProBasicFieldProId
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "inputnumber",
			CustomType: "inputnumber",
			Props: lc_table.LcProps{
				InputNumber: lc_table.LcPropInputNumber{
					Accuracy:   "1",
					Thousandth: false,
					Percentage: false,
					Unique:     false,
					Required:   false,
				},
			},
		},
	})

	addField = consts.ProBasicFieldCode
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "input",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldPreCode
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "input",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldOwnerIds
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "member",
			CustomType: "",
			Props: lc_table.LcProps{
				Member: lc_table.LcPropMember{
					Multiple: true,
					Required: false,
				},
			},
		},
	})

	addField = consts.ProBasicFieldParticipantIds
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "member",
			CustomType: "",
			Props: lc_table.LcProps{
				Member: lc_table.LcPropMember{
					Multiple: true,
					Required: false,
				},
			},
		},
	})

	// 外部群聊标识 []string
	otherFields = append(otherFields, consts.ProBasicFieldOutChat)

	// 单选类型
	// 优先级。原来的项目优先级并没有使用上，暂时不配置
	//addField = consts.ProBasicFieldPriorityId
	//fields = append(fields, lc_table.LcCommonField{
	//	Label: consts.ProFormFieldsMap[addField],
	//	Name:  addField,
	//	Field: lc_table.LcFieldData{
	//		Type:       "input",
	//		CustomType: "",
	//		Props:      lc_table.LcProps{},
	//	},
	//})

	addField = consts.ProBasicFieldPlanStartTime
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "datepicker",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldPlanEndTime
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "datepicker",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldPublicStatus
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "select",
			CustomType: "",
			Props: lc_table.LcProps{
				Member: lc_table.LcPropMember{},
				Select: lc_table.LcPropSelect{
					Options: ProFormFieldOptionsForPublicStatus(),
				},
			},
		},
	})

	addField = consts.ProBasicFieldStatus
	//statusOptionList, err := GetProjectStatusOptionList(orgId)
	//if err != nil {
	//	log.Error(err)
	//	return nil, nil, err
	//}
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "select",
			CustomType: "",
			Props: lc_table.LcProps{
				Member: lc_table.LcPropMember{},
				Select: lc_table.LcPropSelect{
					Options: ProFormFieldOptionsForProStatus(),
				},
			},
		},
	})

	addField = consts.ProBasicFieldTemplateFlag
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "input",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldResource
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "input",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldIsFiling
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "select",
			CustomType: "",
			Props: lc_table.LcProps{
				Member: lc_table.LcPropMember{},
				Select: lc_table.LcPropSelect{
					Options: ProFormFieldOptionsForIsFiling(),
				},
			},
		},
	})

	addField = consts.ProBasicFieldIsEnableWorkHours
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "select",
			CustomType: "",
			Props: lc_table.LcProps{
				Member: lc_table.LcPropMember{},
				Select: lc_table.LcPropSelect{
					Options: ProFormFieldOptionsForIsEnableWorkHours(),
				},
			},
		},
	})

	addField = consts.ProBasicFieldRemark
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "input",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldOutCalendarSettings
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "inputnumber",
			CustomType: "inputnumber",
			Props: lc_table.LcProps{
				InputNumber: lc_table.LcPropInputNumber{
					Accuracy:   "1",
					Thousandth: false,
					Percentage: false,
					Unique:     false,
					Required:   false,
				},
			},
		},
	})

	addField = consts.ProBasicFieldOutCalendar
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "input",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldOutChat
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "input",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldChatSettings
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "input",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldCreator
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "member",
			CustomType: "",
			Props: lc_table.LcProps{
				Member: lc_table.LcPropMember{
					Multiple: false,
					Required: false,
				},
			},
		},
	})

	addField = consts.ProBasicFieldUpdator
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "member",
			CustomType: "",
			Props: lc_table.LcProps{
				Member: lc_table.LcPropMember{
					Multiple: false,
					Required: false,
				},
			},
		},
	})

	addField = consts.ProBasicFieldCreateTime
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "datepicker",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	addField = consts.ProBasicFieldUpdateTime
	fields = append(fields, lc_table.LcCommonField{
		Label:    consts.ProFormFieldsMap[addField],
		Name:     addField,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "datepicker",
			CustomType: "",
			Props:      lc_table.LcProps{},
		},
	})

	return fields, otherFields, nil
}

func ProFormFieldOptionsForProStatus() []lc_table.LcOptions {
	projectStatusList := consts.ProjectStatusList
	options := []lc_table.LcOptions{}
	for _, status := range projectStatusList {
		options = append(options, lc_table.LcOptions{
			Color:     "",
			FontColor: "",
			Id:        status.ID,
			Value:     status.Name,
			Sort:      nil,
		})
	}
	return options
}

func ProFormFieldOptionsForIsEnableWorkHours() []lc_table.LcOptions {
	return []lc_table.LcOptions{
		{
			Color:     "",
			FontColor: "",
			Id:        "1",
			Value:     "已开启",
			Sort:      nil,
		}, {
			Color:     "",
			FontColor: "",
			Id:        "2",
			Value:     "未开启",
			Sort:      nil,
		},
	}
}

func ProFormFieldOptionsForIsFiling() []lc_table.LcOptions {
	return []lc_table.LcOptions{
		{
			Color:     "",
			FontColor: "",
			Id:        "1",
			Value:     "已归档",
			Sort:      nil,
		}, {
			Color:     "",
			FontColor: "",
			Id:        "2",
			Value:     "未归档",
			Sort:      nil,
		},
	}
}

func ProFormFieldOptionsForPublicStatus() []lc_table.LcOptions {
	return []lc_table.LcOptions{
		{
			Color:     "",
			FontColor: "",
			Id:        "1",
			Value:     "公开",
			Sort:      nil,
		}, {
			Color:     "",
			FontColor: "",
			Id:        "2",
			Value:     "私密",
			Sort:      nil,
		},
	}
}

// CreateOrgProjectForm 给组织创建**项目表单**
func CreateOrgProjectForm(orgBo bo.OrganizationBo, oldOrgRemarkObj *orgvo.OrgRemarkConfigType) (int64, errs.SystemErrorInfo) {
	// 新建项目表（form），用于存放该组织的所有项目
	appType := consts.LcAppTypeForForm
	formName := fmt.Sprintf("project form for %s", orgBo.Name)
	fields, commonFields, err := GetProjectFormFields(orgBo.Id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	resp := appfacade.CreateLessCodeApp(&permissionvo.CreateLessCodeAppReq{
		OrgId:     &orgBo.Id,
		AppType:   &appType,
		Name:      &formName,
		ExtendsId: 0,
		PkgId:     0,
		UserId:    &orgBo.Owner,
		Config:    json.ToJsonIgnoreError(formvo.LessFormConfigData{Fields: fields, BaseFields: commonFields}),
		ProjectId: 0,
		ParentId:  0, //历史数据没有目录
		Hidden:    0, //是否隐藏？

	})
	if resp.Failure() {
		log.Errorf("[CreateOrgProjectForm][项目同步]失败，组织id:%d, 原因:%s", orgBo.Id, resp.Error())
		return 0, resp.Error()
	}
	// 将创建好的 form 的 id 存入组织的 remark 中
	orgProFormId := resp.Data.Id
	log.Infof("[CreateOrgProjectForm] new project form appId: %d", orgProFormId)
	return orgProFormId, nil
}

// saveOrgFields 保存组织字段
func saveOrgFields(orgId, userId int64) errs.SystemErrorInfo {
	orgFields, errFields := orgCustomFields()
	if errFields != nil {
		log.Error(errFields)
		return errFields
	}

	//将优先级也作为组织字段保存起来
	editableFlag := true
	orgFields = append(orgFields, lc_helper.GetOrgPriorityField())

	//将需求类型，需求来源，缺陷类型，严重程度转为组织字段
	//需求类型
	demandTypeOptions := []lc_table.LcOptions{}
	for k, i := range consts.DemandTypeMap {
		demandTypeOptions = append(demandTypeOptions, lc_table.LcOptions{
			Color:     "",
			FontColor: consts.DefaultSelectFontColor,
			Id:        i,
			Value:     k,
		})
	}
	orgFields = append(orgFields, lc_table.LcCommonField{
		Label:    "需求类型",
		Name:     consts.BasicFieldDemandType,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "select",
			CustomType: "",
			Props: lc_table.LcProps{
				IsText: true,
				Select: lc_table.LcPropSelect{
					Options: demandTypeOptions,
				},
			},
		},
	})

	//需求来源
	demandSourceOptions := []lc_table.LcOptions{}
	for k, i := range consts.DemandSourceMap {
		demandSourceOptions = append(demandSourceOptions, lc_table.LcOptions{
			Color:     "",
			FontColor: consts.DefaultSelectFontColor,
			Id:        i,
			Value:     k,
		})
	}
	orgFields = append(orgFields, lc_table.LcCommonField{
		Label:    "需求来源",
		Name:     consts.BasicFieldDemandSource,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "select",
			CustomType: "",
			Props: lc_table.LcProps{
				IsText: true,
				Select: lc_table.LcPropSelect{
					Options: demandSourceOptions,
				},
			},
		},
	})

	//缺陷类型
	bugTypeOptions := []lc_table.LcOptions{}
	for k, i := range consts.BugTypeMap {
		bugTypeOptions = append(bugTypeOptions, lc_table.LcOptions{
			Color:     "",
			FontColor: consts.DefaultSelectFontColor,
			Id:        i,
			Value:     k,
		})
	}
	orgFields = append(orgFields, lc_table.LcCommonField{
		Label:    "缺陷类型",
		Name:     consts.BasicFieldBugType,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "select",
			CustomType: "",
			Props: lc_table.LcProps{
				IsText: true,
				Select: lc_table.LcPropSelect{
					Options: bugTypeOptions,
				},
			},
		},
	})

	//严重程度
	bugPropertyOptions := []lc_table.LcOptions{}
	for k, i := range consts.BugPropertyMap {
		bugPropertyOptions = append(bugPropertyOptions, lc_table.LcOptions{
			Color:     "",
			FontColor: consts.DefaultSelectFontColor,
			Id:        i,
			Value:     k,
		})
	}
	orgFields = append(orgFields, lc_table.LcCommonField{
		Label:    "严重程度",
		Name:     consts.BasicFieldBugProperty,
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "select",
			CustomType: "",
			Props: lc_table.LcProps{
				IsText: true,
				Select: lc_table.LcPropSelect{
					Options: bugPropertyOptions,
				},
			},
		},
	})

	var columns []interface{}
	for _, field := range orgFields {
		columns = append(columns, field)
	}

	if len(orgFields) > 0 {
		//resp := formfacade.LessBaseFormSave(formvo.LessBaseFormSaveReq{
		//	OrgId:   orgId,
		//	UserId:  userId,
		//	Added:   addedField,
		//	Deleted: nil,
		//})
		resp := tablefacade.InitOrgColumns(orgvo.InitOrgColumnsReq{
			OrgId:  orgId,
			UserId: userId,
			Input:  &orgvo.InitOrgColumnsRequest{Columns: columns},
		})
		if resp.Failure() {
			log.Error(resp.Error())
			return resp.Error()
		}
	}
	return nil
}

// SaveOrgRemarkAndOrgFields 保存组织无码配置和组织字段
func SaveOrgRemarkAndOrgFields(orgId, userId int64, orgRemarkObj *orgvo.OrgRemarkConfigType) (*orgvo.SaveOrgSummaryTableAppIdReqVoData, errs.SystemErrorInfo) {
	// 将汇总表的 appId 存入组织属性（remark）中。如果汇总表等应用不存在，则创建之。
	orgConfig, _, err := SaveOrgSomeTableAppId(orgId, userId, orgRemarkObj)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = saveOrgFields(orgId, userId)
	if err != nil {
		return nil, err
	}

	return orgConfig, nil
}

func orgCustomFields() ([]lc_table.LcCommonField, errs.SystemErrorInfo) {
	fields := []lc_table.LcCommonField{}
	processOptions := []lc_table.PolarisPersonFieldValue{}
	err := json.FromJson(lc_helper.ProcessFieldValue, &processOptions)
	if err != nil {
		log.Error(err)
		return nil, errs.JSONConvertError
	}

	accuracy := "1"
	if len(processOptions) > 1 {
		switch processOptions[len(processOptions)-1].Value {
		case "1":
			accuracy = "1.0"
		case "2":
			accuracy = "1.00"
		case "3":
			accuracy = "1.000"
		default:
			break
		}
	}
	editableFlag := true
	fields = append(fields, lc_table.LcCommonField{
		Name:     "_field_900",
		Label:    "进度",
		Editable: &editableFlag,
		Writable: true,
		Field: lc_table.LcFieldData{
			Type:       "inputnumber",
			CustomType: "inputnumber",
			Props: lc_table.LcProps{
				InputNumber: lc_table.LcPropInputNumber{
					Accuracy:   accuracy,
					Thousandth: false,
					Percentage: true,
					Unique:     false,
					Required:   false,
				},
			},
		},
	})
	storyPointsOptions := []lc_table.LcOptions{}
	err2 := json.FromJson(lc_helper.StoryPointFieldValue, &storyPointsOptions)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.JSONConvertError
	}
	fields = append(fields, lc_table.LcCommonField{
		Name:     "_field_901",
		Label:    "Story Points",
		Editable: &editableFlag,
		Writable: true,
		Field: lc_table.LcFieldData{
			Type:       "select",
			CustomType: "",
			Props: lc_table.LcProps{
				IsText:   false,
				Required: false,
				Member:   lc_table.LcPropMember{},
				Select:   lc_table.LcPropSelect{Options: storyPointsOptions},
			},
		},
	})
	scoreOptions := []lc_table.LcOptions{}
	err3 := json.FromJson(lc_helper.ScoreFieldValue, &scoreOptions)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.JSONConvertError
	}
	fields = append(fields, lc_table.LcCommonField{
		Name:     "_field_902",
		Label:    "评分",
		Writable: true,
		Editable: &editableFlag,
		Field: lc_table.LcFieldData{
			Type:       "select",
			CustomType: "",
			Props: lc_table.LcProps{
				IsText:   false,
				Required: false,
				Member:   lc_table.LcPropMember{},
				Select:   lc_table.LcPropSelect{Options: scoreOptions},
			},
		},
	})
	return fields, nil
}
