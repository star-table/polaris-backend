package domain

import (
	tableV1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/tablefacade"
)

// GetTableColumnByAppId 通appId获取表的列信息
func GetTableColumnByAppId(orgId int64, userId int64, appId int64, columnIds []string) (*projectvo.TablesColumnsRespData, errs.SystemErrorInfo) {
	tablesColumns := tablefacade.ReadTableSchemasByAppId(projectvo.GetTableSchemasByAppIdReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &tableV1.ReadTableSchemasByAppIdRequest{
			AppId:     appId,
			ColumnIds: columnIds,
		},
	})
	if tablesColumns.Failure() {
		log.Errorf("[GetTableColumnByAppId] tablefacade.ReadTableSchemasByAppId failed, orgId:%v, appId:%v, err:%v", orgId, appId, tablesColumns.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, tablesColumns.Message)
	}

	return tablesColumns.Data, nil
}

// GetTablesColumnsByTableIds 通过tableIds获取表头信息
func GetTablesColumnsByTableIds(orgId, userId int64, tableIds []int64, columnIds []string, isNeedDescription bool) (*projectvo.TablesColumnsRespData, errs.SystemErrorInfo) {
	tablesColumns := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &tableV1.ReadTableSchemasRequest{
			TableIds:          tableIds,
			ColumnIds:         columnIds,
			IsNeedDescription: isNeedDescription,
		},
	})

	if tablesColumns.Failure() {
		log.Errorf("[GetTableColumnConfig] tablefacade.GetTableColumns failed, orgId:%v, tableIds:%v, err:%v", orgId, tableIds, tablesColumns.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, tablesColumns.Message)
	}

	return tablesColumns.Data, nil
}

// RemoveColumnByProjectType 有些项目类型不需要某些表头
func RemoveColumnByProjectType(columns []*projectvo.TableColumnData, projectType int64, extraColumns ...string) []*projectvo.TableColumnData {
	removeColumnIds := GetNoNeedColumnByProjectType(projectType)
	removeColumnIds = append(removeColumnIds, extraColumns...)
	newColumns := make([]*projectvo.TableColumnData, 0, len(columns))
	for _, column := range columns {
		isIn := false
		for _, id := range removeColumnIds {
			if id == column.Name {
				isIn = true
				break
			}
		}
		if !isIn {
			newColumns = append(newColumns, column)
		}
	}

	return newColumns
}

func GetNoNeedColumnByProjectType(projectType int64) []string {
	switch projectType {
	case consts.ProjectTypeNormalId:
		return consts.HiddenColumnNameForCommon
	case consts.ProjectTypeCommon2022V47:
		return consts.HiddenColumnNamesForCommon2022Type
	case consts.ProjectTypeAgileId:
		return consts.HiddenColumnNameForAgileId
	case consts.ProjectTypeEmpty:
		return consts.HiddenColumnNameForEmpty
	}

	return []string{}
}

// GetTableColumnConfig 获取一个表的指定列信息
// @columnIds 获取的列的 key/name（列的标识符）
func GetTableColumnConfig(orgId int64, tableId int64, columnIds []string, isNeedCommonColumn bool) (*projectvo.TableColumnConfigInfo, errs.SystemErrorInfo) {
	tableColumns := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
		OrgId:  orgId,
		UserId: 0,
		Input: &tableV1.ReadTableSchemasRequest{
			TableIds:           []int64{tableId},
			ColumnIds:          columnIds,
			IsNeedDescription:  true,
			IsNeedCommonColumn: isNeedCommonColumn,
		},
	})

	if tableColumns.Failure() {
		log.Errorf("[GetTableColumnConfig] tablefacade.GetTableColumns failed, orgId:%v, tableId:%v, err:%v", orgId, tableId, tableColumns.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, tableColumns.Message)
	}
	if len(tableColumns.Data.Tables) < 1 {
		err := errs.TableNotExist
		log.Errorf("[GetTableColumnConfig] orgId:%v, tableId:%v, err:%v", orgId, tableId, err)
		return nil, err
	}
	columnList := tableColumns.Data.Tables[0].Columns

	columnSlice := []*projectvo.TableColumnData{}
	for _, column := range columnList {
		columnSlice = append(columnSlice, &projectvo.TableColumnData{
			Name:              column.Name,
			Label:             column.Label,
			AliasTitle:        column.AliasTitle,
			Description:       column.Description,
			IsSys:             column.IsSys,
			IsOrg:             column.IsOrg,
			Writable:          column.Writable,
			Editable:          column.Editable,
			Unique:            column.Unique,
			UniquePreHandler:  column.UniquePreHandler,
			SensitiveStrategy: column.SensitiveStrategy,
			SensitiveFlag:     column.SensitiveFlag,
			Field: projectvo.TableColumnField{
				Type:       column.Field.Type,
				CustomType: column.Field.CustomType,
				DataType:   column.Field.DataType,
				Props:      column.Field.Props,
			},
		})
	}

	return &projectvo.TableColumnConfigInfo{
		TableId: tableId,
		Columns: columnSlice,
	}, nil
}

// GetTableColumnsMap 获取一个表的列信息，以列的 name 为 key，返回 map
func GetTableColumnsMap(orgId int64, tableId int64, columnIds []string, flags ...bool) (map[string]*projectvo.TableColumnData, errs.SystemErrorInfo) {
	if tableId == 0 {
		summeryTableResp := tablefacade.GetSummeryTableId(projectvo.GetSummeryTableIdReqVo{
			OrgId:  orgId,
			UserId: 0,
			Input:  &tableV1.ReadSummeryTableIdRequest{},
		})
		if summeryTableResp.Failure() {
			log.Errorf("[GetTableColumnsMap] failed, orgId: %d, tableId: %d, err: %v", orgId, tableId, summeryTableResp.Error())
		} else {
			tableId = summeryTableResp.Data.TableId
		}
	}

	var isNeedCommonColumn bool
	if len(flags) > 0 {
		isNeedCommonColumn = flags[0]
	}

	resMap := make(map[string]*projectvo.TableColumnData, 0)
	columnInfo, err := GetTableColumnConfig(orgId, tableId, columnIds, isNeedCommonColumn)
	if err != nil {
		log.Errorf("[GetTableColumnsMap] err: %v", err)
		return resMap, err
	}
	columns := columnInfo.Columns
	for _, column := range columns {
		resMap[column.Name] = column
	}

	return resMap, nil
}
