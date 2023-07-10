package domain

import (
	tableV1 "github.com/star-table/interface/golang/table/v1"
	tablev1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/tablefacade"
	"github.com/spf13/cast"
)

func GetTableInfo(orgId int64, appId, tableId int64) (*projectvo.TableMetaData, errs.SystemErrorInfo) {
	resp := tablefacade.ReadTables(projectvo.GetTablesReqVo{
		OrgId:  orgId,
		UserId: 0,
		Input: &tableV1.ReadTablesRequest{
			AppId: appId,
		},
	})
	if resp.Failure() {
		log.Errorf("[GetTableInfo]GetTablesOfApp failed:%v, orgId:%d, tableId:%d", resp.Error(), orgId, tableId)
		return nil, resp.Error()
	}
	for _, table := range resp.Data.Tables {
		if table.TableId == tableId {
			return table, nil
		}
	}

	return nil, errs.TableNotExist
}

func GetAppTableList(orgId int64, appId int64) ([]*projectvo.TableMetaData, errs.SystemErrorInfo) {
	resp := tablefacade.ReadTables(projectvo.GetTablesReqVo{
		OrgId:  orgId,
		UserId: 0,
		Input: &tableV1.ReadTablesRequest{
			AppId: appId,
		},
	})
	if resp.Failure() {
		log.Errorf("[GetTableInfo]GetTablesOfApp failed:%v, orgId:%d, appId:%d", resp.Error(), orgId, appId)
		return nil, resp.Error()
	}

	return resp.Data.Tables, nil
}

// GetTableListByProAppIds 通过项目对应的 appIds，查询其下的 table list
func GetTableListByProAppIds(orgId int64, proAppIds []int64) ([]*projectvo.AppTables, errs.SystemErrorInfo) {
	tables := make([]*projectvo.AppTables, 0, len(proAppIds))
	if len(proAppIds) < 1 {
		return tables, nil
	}

	tablesResp := tablefacade.ReadTablesByApps(projectvo.ReadTablesByAppsReqVo{
		OrgId:  orgId,
		UserId: 0,
		Input: &tableV1.ReadTablesByAppsRequest{
			AppIds: proAppIds,
		},
	})
	if tablesResp.Failure() {
		log.Errorf("[GetTableListByIds] ReadTablesByApps err: %v", tablesResp.Error())
		return nil, tablesResp.Error()
	}

	return tablesResp.Data.AppsTables, nil
}

//// GetTableListMapByProjectIds 通过 projectIds 获取对应的 table 列表，并按 appId 分组
//func GetTableListMapByProjectIds(orgId int64, projectIds []int64) (map[int64][]*projectvo.TableMetaData, errs.SystemErrorInfo) {
//	tableListMap := make(map[int64][]*projectvo.TableMetaData, 0)
//	proAppIds, err := GetProjectAppIdsByProjectIds(orgId, projectIds)
//	if err != nil {
//		log.Errorf("[GetTableListMapByProjectIds] err: %v", err)
//		return tableListMap, err
//	}
//	appsTables, err := GetTableListByProAppIds(orgId, proAppIds)
//	if err != nil {
//		log.Errorf("[GetTableListMapByProAppIds] err: %v", err)
//		return tableListMap, err
//	}
//	for _, obj := range appsTables {
//		tableListMap[obj.AppId] = obj.Tables
//	}
//
//	return tableListMap, nil
//}

// GetTableByTableId 通过表格id，查询tablelist
func GetTableByTableId(orgId int64, userId int64, tableId int64) (*projectvo.TableMetaData, errs.SystemErrorInfo) {
	tablesResp := tablefacade.ReadOneTable(projectvo.GetTableInfoReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &tableV1.ReadTableRequest{
			TableId: tableId,
		},
	})
	if tablesResp.Failure() {
		log.Errorf("[GetTableByTableId] ReadOneTable err: %v", tablesResp.Error())
		return nil, tablesResp.Error()
	}

	return tablesResp.Data.Table, nil
}

// GetTableListMapByProAppIds 通过项目对应的 appIds, 查询其下的 table list，返回以 appId 为 key，值为 tables list 的 map
func GetTableListMapByProAppIds(orgId int64, proAppIds []int64) (map[int64][]*projectvo.TableMetaData, errs.SystemErrorInfo) {
	resMap := make(map[int64][]*projectvo.TableMetaData, len(proAppIds))
	if len(proAppIds) < 1 {
		return resMap, nil
	}
	appsTables, err := GetTableListByProAppIds(orgId, proAppIds)
	if err != nil {
		log.Errorf("[GetTableListMapByProAppIds] err: %v", err)
		return resMap, err
	}

	for _, obj := range appsTables {
		resMap[obj.AppId] = obj.Tables
	}

	return resMap, nil
}

func CheckTableName(orgId, userId, appId int64, name string) errs.SystemErrorInfo {
	// 校验是否表是否重名
	tableList := tablefacade.ReadTables(projectvo.GetTablesReqVo{
		OrgId:  orgId,
		UserId: userId,
		Input:  &tablev1.ReadTablesRequest{AppId: appId},
	})
	if tableList.Failure() {
		log.Errorf("[domain.CheckTableName] tablefacade.ReadTables failed, orgId:%v, userId:%v, appId:%v, err:%v", orgId, userId, appId, tableList.Error())
		return errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, "tablefacade.CreateTable error")
	}

	for _, table := range tableList.Data.Tables {
		if name == table.Name {
			return errs.TableSameName
		}
	}
	return nil
}

func GetRowsExpand(req *formvo.LessIssueListReq) *formvo.LessIssueListResp {
	result := &formvo.LessIssueListResp{}
	reply := GetRows(req)

	if reply.Failure() {
		log.Errorf("[GetRowsExpand] list req:%v, err: %v", req, reply.Error())
		result.Err = reply.Err
		return result
	}

	result.UserDepts = reply.UserDepts
	result.Data = reply.Data
	result.List = make([]map[string]interface{}, 0, reply.Data.Count)
	err := json.Unmarshal(reply.Data.Data, &result.List)
	if err != nil {
		log.Errorf("[GetRowsExpand] Unmarshal data:%v, err: %v", reply.Data.Data, err)
		result.Err = vo.NewErr(errs.JSONConvertError)
		return result
	}

	if reply.Data.RelateData != nil {
		relateDatas := make(map[string]map[string]interface{}, reply.Data.Count)
		err = json.Unmarshal(reply.Data.RelateData, &relateDatas)
		if err != nil {
			log.Errorf("[GetRowsExpand] Unmarshal data:%v, err: %v", reply.Data.RelateData, err)
			result.Err = vo.NewErr(errs.JSONConvertError)
			return result
		}
		for _, m := range result.List {
			id := cast.ToString(m[consts.BasicFieldDataId])
			if dm, ok := relateDatas[id]; ok {
				for k, v := range dm {
					if m2, ok2 := m[k].(map[string]interface{}); ok2 && (m2["linkTo"] != nil || m2["linkFrom"] != nil) {
						m2["linkNames"] = v
					} else {
						m[k] = v
					}
				}
			}
		}
	}

	return result
}

func GetRows(req *formvo.LessIssueListReq, appAuthInfo ...string) *projectvo.ListRowsReply {
	req.Condition.Conds = append(req.Condition.Conds, &vo.LessCondsData{
		Type:   "equal",
		Value:  req.OrgId,
		Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
	})

	listReq := &tablev1.ListRequest{
		Query:        json.ToJsonIgnoreError(req),
		TableId:      req.TableId,
		NeedChangeId: req.NeedChangeId,
	}
	if len(appAuthInfo) > 0 {
		listReq.AppAuthData = appAuthInfo[0]
	}
	reply := tablefacade.List(req.OrgId, req.UserId, listReq)

	if reply.Failure() {
		log.Errorf("[GetRows] list req:%v, err: %v", req, reply.Error())
		return reply
	}

	userDepts, err2 := getUserDepts(req.OrgId, reply.Data.UserIds, reply.Data.DeptIds)
	if err2 != nil {
		log.Errorf("[GetRowsExpand] Unmarshal data:%v, err: %v", reply.Data.Data, err2)
		reply.Err = vo.NewErr(err2)
		return reply
	}
	reply.UserDepts = userDepts

	return reply
}

func GetRawRows(orgId, userId int64, req *tablev1.ListRawRequest) (*formvo.LessIssueRawListResp, errs.SystemErrorInfo) {
	condition := &tablev1.Condition{Type: tablev1.ConditionType_and, Conditions: []*tablev1.Condition{
		GetRowsCondition(consts.BasicFieldOrgId, tablev1.ConditionType_equal, orgId, nil),
	}}
	if req.Condition != nil {
		condition.Conditions = append(condition.Conditions, req.Condition)
	}
	req.Condition = condition

	reply := tablefacade.ListRaw(orgId, userId, req)
	if reply.Failure() {
		log.Errorf("[GetRawRows] ListRaw req:%v, err: %v, stack:%v", req, reply.Error(), stack.GetStack())
		return nil, reply.Error()
	}

	result := &formvo.LessIssueRawListResp{Data: []map[string]interface{}{}}
	err := json.Unmarshal(reply.Data.Data, &result.Data)
	if err != nil {
		log.Errorf("[GetRawRows] Unmarshal data:%v, err: %v", reply.Data.Data, err)
		return nil, errs.JSONConvertError
	}

	return result, nil
}

func GetRowsCondition(columnId string, columnType tablev1.ConditionType, value, values interface{}) *tablev1.Condition {
	if value != nil {
		if vs, ok := value.([]interface{}); ok {
			values = vs
		} else {
			values = []interface{}{value}
		}
	}

	if columnType != tablev1.ConditionType_raw_sql {
		columnId = lc_helper.ConvertToCondColumn(columnId)
	}

	return &tablev1.Condition{
		Type:   columnType,
		Value:  json.ToJsonIgnoreError(values),
		Column: columnId,
	}
}

func GetRowsAndCondition(conditions ...*tablev1.Condition) *tablev1.Condition {
	return &tablev1.Condition{Type: tablev1.ConditionType_and, Conditions: conditions}
}

func GetNoRecycleCondition(conditions ...*tablev1.Condition) []*tablev1.Condition {
	conditions = append(conditions,
		GetRowsCondition(consts.BasicFieldRecycleFlag, tablev1.ConditionType_equal, consts.AppIsNoDelete, nil),
	)

	return conditions
}
