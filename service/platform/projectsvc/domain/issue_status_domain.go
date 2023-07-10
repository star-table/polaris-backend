package domain

import (
	tableV1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo/status"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/tablefacade"
)

func GetTableStatus(orgId int64, tableId int64) ([]status.StatusInfoBo, errs.SystemErrorInfo) {
	if tableId == 0 {
		// 找汇总表的状态
		summeryTableResp := tablefacade.GetSummeryTableId(projectvo.GetSummeryTableIdReqVo{
			OrgId:  orgId,
			UserId: 0,
			Input:  &tableV1.ReadSummeryTableIdRequest{},
		})
		if summeryTableResp.Failure() {
			log.Errorf("[GetTableStatus] failed, orgId: %d, tableId: %d, err: %v", orgId, tableId, summeryTableResp.Error())
			return nil, summeryTableResp.Error()
		}
		tableId = summeryTableResp.Data.TableId
	}
	tableColumnsResp := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
		OrgId:  orgId,
		UserId: 0,
		Input: &tableV1.ReadTableSchemasRequest{
			TableIds:  []int64{tableId},
			ColumnIds: []string{consts.BasicFieldIssueStatus},
		},
	})
	if tableColumnsResp.Failure() {
		log.Errorf("[GetTableStatus] failed, orgId: %d, tableId: %d, err: %v", orgId, tableId, tableColumnsResp.Error())
		return nil, tableColumnsResp.Error()
	}
	if len(tableColumnsResp.Data.Tables) < 1 {
		err := errs.TableNotExist
		log.Errorf("[GetTableStatus] query table column err: not found table or table column, orgId: %d, tableId: %d", orgId, tableId)
		return nil, err
	}
	for _, tableSchema := range tableColumnsResp.Data.Tables {
		// curTableId := tableSchema.TableId
		columns := tableSchema.Columns
		columnMap := GetColumnMapByColumnsForTableMode(columns)
		if statusColumn, ok := columnMap[consts.BasicFieldIssueStatus]; ok {
			return GetStatusListFromStatusColumn(statusColumn), nil
		}
	}
	return []status.StatusInfoBo{}, nil
}
