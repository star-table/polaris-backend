package issue_export

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
)

/// 任务导出

func ExportIssueWithCustomField() {

}

// 获取任务导出的字段列表
//func GetExportFields(orgId, projectId int64) ([]*vo.GetExportFieldsRespFieldsItem, errs.SystemErrorInfo) {
//	inBuildList, err := GetExportFieldsFixed()
//	if err != nil {
//		return nil, err
//	}
//	customList, err := GetExportFieldsCustom(orgId, projectId)
//	if err != nil {
//		return nil, err
//	}
//	inBuildList = append(inBuildList, customList...)
//	return inBuildList, nil
//}

//func GetExportFieldsCustom(orgId, projectId int64) ([]*vo.GetExportFieldsRespFieldsItem, errs.SystemErrorInfo) {
//	list := make([]*vo.GetExportFieldsRespFieldsItem, 0)
//	page := 1
//	size := 1000
//	isCurrentProject := 1
//	customFieldsObj, err := service.CustomFieldList(orgId, 1, 1000, vo.CustomFieldListReq{
//		Page:                 &page,
//		Size:                 &size,
//		ProjectID:            &projectId,
//		ProjectObjectTypeID:  nil,
//		IsUsedCurrentProject: &isCurrentProject,
//		IsOrgField:           nil,
//		Name:                 nil,
//		OrderType:            nil,
//	})
//	if err != nil {
//		return nil, err
//	}
//	if customFieldsObj.Total < 1 {
//		return list, nil
//	}
//	for _, fieldObj := range customFieldsObj.List {
//		list = append(list, &vo.GetExportFieldsRespFieldsItem{
//			FieldID:    fieldObj.ID,
//			Name:       fieldObj.Name,
//			IsMust:     false,// 自定义字段可以都不要，用户自己选择。
//			DefineType: 11,
//		})
//	}
//	return list, nil
//}

// 获取内置的固定字段列表
func GetExportFieldsFixed() ([]*vo.GetExportFieldsRespFieldsItem, errs.SystemErrorInfo) {
	list := make([]*vo.GetExportFieldsRespFieldsItem, 0)
	// 字段名：是否必选。例如 id: false，表示 id 字段不是必选。导出时，用户可以选择不导出该字段值。
	fixedColumns := map[string]bool{
		"id":            false,
		"code":          false,
		"title":         true,
		"issueType":     false,
		"isFiling":      false,
		"owner":         false,
		"priority":      false,
		"describe":      false,
		"sourceId":      false,
		"planStartTime": false,
		"planEndTime":   false,
		"follower":      false,
	}
	for columnName, isMust := range fixedColumns {
		list = append(list, &vo.GetExportFieldsRespFieldsItem{
			FieldID:    0,
			Name:       columnName,
			IsMust:     isMust,
			DefineType: 10,
		})
	}
	return list, nil
}
