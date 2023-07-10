package service

import (
	"sort"

	"github.com/spf13/cast"

	tablev1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/tablefacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"upper.io/db.v3"
)

func CreateColumn(req projectvo.CreateColumnReqVo) (*projectvo.CreateColumnReply, errs.SystemErrorInfo) {
	//判断项目权限
	authErr := domain.AuthProjectWithAppId(req.OrgId, req.UserId, req.Input.ProjectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigModifyField, req.Input.AppId)
	if authErr != nil {
		log.Errorf("[CreateColumn.AuthProjectWithAppId], err: %v", authErr)
		return nil, authErr
	}

	columnInfo := tablefacade.CreateColumn(req)
	if columnInfo.Code != 0 {
		log.Errorf("[tablefacade.CreateColumn] err: %v", columnInfo.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, columnInfo.Message)
	}

	asyn.Execute(func() {
		recordColumnChangeTrends(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.TableId, req.Input.Column.Name, "", "")
	})

	return columnInfo.Data, nil
}

func CopyColumn(req projectvo.CopyColumnReqVo) (*tablev1.CopyColumnReply, errs.SystemErrorInfo) {
	//判断项目权限
	authErr := domain.AuthProjectWithAppId(req.OrgId, req.UserId, req.Input.ProjectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigModifyField, req.Input.AppId)
	if authErr != nil {
		log.Errorf("[CopyColumn.AuthProjectWithAppId], err: %v", authErr)
		return nil, authErr
	}

	columnInfo := tablefacade.CopyColumn(req)
	if columnInfo.Code != 0 {
		log.Errorf("[tablefacade.CopyColumn] err: %v", columnInfo.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, columnInfo.Message)
	}

	asyn.Execute(func() {
		recordColumnChangeTrends(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.TableId, columnInfo.Data.CreateColumnId, "", "")
	})

	return columnInfo.Data, nil
}

func UpdateColumn(req projectvo.UpdateColumnReqVo) (*projectvo.UpdateColumnReply, errs.SystemErrorInfo) {
	params := req.Input
	if params.Column == nil {
		return nil, errs.ParamError
	}
	//判断项目权限
	authErr := domain.AuthProjectWithAppId(req.OrgId, req.UserId, params.ProjectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigModifyField, params.AppId)
	if authErr != nil {
		log.Errorf("[UpdateColumn.AuthProjectWithAppId], err: %v", authErr)
		return nil, authErr
	}
	//如果是更改状态需要判断有没有删除选项，删除了则需要判断该状态下是否有任务，有的话不允许更新
	if params.Column.Name == consts.BasicFieldIssueStatus {
		oldStatusList, oldStatusListErr := domain.GetTableStatus(req.OrgId, params.TableId)
		if oldStatusListErr != nil {
			log.Errorf("[UpdateColumn]GetTableStatus failed:%v, orgId:%d, tableId:%d", oldStatusListErr, req.OrgId, params.TableId)
			return nil, oldStatusListErr
		}

		newStatusList := domain.GetStatusListFromStatusColumn(params.Column)
		var newStatusIds []int64
		for _, infoBo := range newStatusList {
			newStatusIds = append(newStatusIds, infoBo.ID)
		}
		var deleteIds []int64
		for _, infoBo := range oldStatusList {
			if ok, _ := slice.Contain(newStatusIds, infoBo.ID); !ok {
				deleteIds = append(deleteIds, infoBo.ID)
			}
		}
		if len(deleteIds) > 0 {
			listReq := &tablev1.ListRawRequest{
				FilterColumns: []string{
					"count(*) as count",
				},
				Condition: &tablev1.Condition{
					Type: tablev1.ConditionType_and,
					Conditions: domain.GetNoRecycleCondition(
						domain.GetRowsCondition(consts.BasicFieldTableId, tablev1.ConditionType_equal, cast.ToString(params.TableId), nil),
						domain.GetRowsCondition(consts.BasicFieldIssueStatus, tablev1.ConditionType_in, nil, deleteIds),
					),
				},
				TableId: params.TableId,
			}

			lessResp, err := domain.GetRawRows(req.OrgId, req.UserId, listReq)
			if err != nil {
				log.Errorf("[UpdateColumn] orgId:%d, projectId:%d, LessIssueList failure:%v", req.OrgId, params.ProjectId, err.Error())
				return nil, err
			}
			issueAssignCountBos := []bo.CommonCountBo{}
			err2 := copyer.Copy(lessResp.Data, &issueAssignCountBos)
			if err2 != nil {
				log.Errorf("[UpdateColumn] orgId:%d, projectId:%d, Copy failure:%v", req.OrgId, params.ProjectId, err2)
				return nil, errs.ObjectCopyError
			}
			if len(issueAssignCountBos) > 0 && issueAssignCountBos[0].Count > 0 {
				return nil, errs.RemainIssuesInStatus
			}
		}
	}

	// 校验 && 任务栏
	updateColumnResp := tablefacade.UpdateColumn(req)
	if updateColumnResp.Code != 0 {
		log.Errorf("[tablefacade.UpdateColumn] err: %v", updateColumnResp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, updateColumnResp.Message)
	}

	// 推送
	asyn.Execute(func() {
		recordColumnChangeTrends(req.OrgId, req.UserId, params.ProjectId, req.Input.TableId, "", req.Input.Column.Name, "")

		// 打开/关闭协作人字段，都需要将该字段的人添加/踢出 群
		if req.Input.Column.Field.Type == consts.LcColumnFieldTypeMember {
		}
	})

	return updateColumnResp.Data, nil
}

func DeleteColumn(req projectvo.DeleteColumnReqVo) (*tablev1.DeleteColumnReply, errs.SystemErrorInfo) {
	params := req.Input
	//默认字段不允许删除
	if params.ColumnId == consts.BasicFieldIssueStatus {
		return nil, errs.DefaultFieldError
	}
	//判断项目权限
	authErr := domain.AuthProjectWithAppId(req.OrgId, req.UserId, params.ProjectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigModifyField, params.AppId)
	if authErr != nil {
		log.Errorf("[UpdateColumn.AuthProjectWithAppId], err: %v", authErr)
		return nil, authErr
	}

	// 附件字段删除，把回收站的附件数据删除，相关的关联关系删除
	errSys := domain.DeleteAttachmentsForColumn(req.OrgId, req.UserId, params.ProjectId, params.TableId, params.ColumnId)
	if errSys != nil {
		// 这里不return，附件字段删除 相关的数据处理失败，不应影响删除列的主操作
		log.Errorf("[DeleteColumn] DeleteAttachmentsForColumn err:%v, orgId:%d, userId:%d, projectId:%d, tableId:%d, columnId:%s",
			errSys, req.OrgId, req.UserId, params.ProjectId, params.TableId, params.ColumnId)
	}

	deleteColumnResp := tablefacade.DeleteColumn(req)
	if deleteColumnResp.Code != 0 {
		log.Errorf("[tablefacade.DeleteColumn] err: %v", deleteColumnResp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, deleteColumnResp.Message)
	}

	asyn.Execute(func() {
		recordColumnChangeTrends(req.OrgId, req.UserId, params.ProjectId, req.Input.TableId, "", "", params.ColumnId)
	})

	return deleteColumnResp.Data, nil
}

func UpdateColumnDesc(req projectvo.UpdateColumnDescriptionReqVo) (*tablev1.UpdateColumnDescriptionReply, errs.SystemErrorInfo) {
	params := req.Input
	//判断项目权限
	authErr := domain.AuthProjectWithAppId(req.OrgId, req.UserId, params.ProjectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigModifyField, params.AppId)
	if authErr != nil {
		log.Errorf("[UpdateColumn.AuthProjectWithAppId], err: %v", authErr)
		return nil, authErr
	}

	updateColumnDescResp := tablefacade.UpdateColumnDescription(req)
	if updateColumnDescResp.Code != 0 {
		log.Errorf("[tablefacade.UpdateColumnDesc] err: %v", updateColumnDescResp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, updateColumnDescResp.Message)
	}

	return updateColumnDescResp.Data, nil

}

func recordColumnChangeTrends(orgId, userId, projectId, tableId int64, createColumn, updateColumn, deleteColumn string) {
	ext := bo.TrendExtensionBo{ProjectObjectTypeId: tableId}
	added := make([]string, 0)
	deleted := make([]string, 0)
	updated := make([]string, 0)

	added = append(added, createColumn)
	updated = append(updated, updateColumn)
	deleted = append(deleted, deleteColumn)

	ext.AddedFormFields = added
	ext.DeletedFormFields = deleted
	ext.UpdatedFormFields = updated

	projectTrendsBo := bo.ProjectTrendsBo{
		PushType:   consts.PushTypeUpdateFormField,
		OrgId:      orgId,
		ProjectId:  projectId,
		OperatorId: userId,
		Ext:        ext,
	}
	domain.PushProjectTrends(projectTrendsBo)
}

func GetOneTableColumns(req projectvo.GetTableColumnReq) (*projectvo.TableColumnsTable, errs.SystemErrorInfo) {
	// 获取汇总表
	if req.TableId <= 0 {
		return getSummaryTableColumns(req)
	}

	// 获取普通表
	return getTableColumns(req)
}

// addIterationOptions 添加迭代选项
func addIterationOptions(orgId, projectId int64, column *projectvo.TableColumnData) errs.SystemErrorInfo {
	iterationList, _, err := domain.GetIterationBoList(0, 0, db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
		consts.TcOrgId:     orgId,
	}, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	options := make([]*lc_table.LcOptions, 0, len(*iterationList)+1)
	options = append(options, &lc_table.LcOptions{
		Id:    0,
		Value: "未规划",
	})
	for _, iterationBo := range *iterationList {
		options = append(options, &lc_table.LcOptions{
			Id:    iterationBo.Id,
			Value: iterationBo.Name,
		})
	}

	if column.Field.Props != nil && column.Field.Props["select"] != nil {
		if m, ok := column.Field.Props["select"].(map[string]interface{}); ok {
			m["options"] = options
		}
	}

	return nil
}

func getSummaryTableColumns(req projectvo.GetTableColumnReq) (*projectvo.TableColumnsTable, errs.SystemErrorInfo) {
	summaryAppId, err := domain.GetOrgSummaryAppId(req.OrgId)
	if err != nil {
		log.Errorf("[GetOneTableColumns] 获取汇总表appID异常: %v", err)
		return nil, err
	}
	tablesColumns, err := domain.GetTableColumnByAppId(req.OrgId, req.UserId, summaryAppId, nil)
	if err != nil {
		log.Errorf("[GetTablesColumnsList].GetTableColumnByAppId failed, orgId:%v, appId:%v, err:%v", req.OrgId, summaryAppId, err)
		return nil, err
	}
	if len(tablesColumns.Tables) == 0 {
		return nil, errs.TableNotExist
	}
	if req.NotAllIssue == 1 {
		return tablesColumns.Tables[0], nil
	}
	tablesMap, err := addIssueStatusOptions(req.OrgId, req.UserId, tablesColumns.Tables[0].Columns)
	if err != nil {
		return nil, err
	}
	err = addProjectOptions(req.OrgId, req.UserId, tablesColumns.Tables[0].Columns, tablesMap)

	return tablesColumns.Tables[0], err
}

// 根据tableId获取表
func getTableColumns(req projectvo.GetTableColumnReq) (*projectvo.TableColumnsTable, errs.SystemErrorInfo) {
	columnsResp, err := domain.GetTablesColumnsByTableIds(req.OrgId, req.UserId, []int64{req.TableId}, nil, true)
	if err != nil {
		log.Errorf("[getTableColumns] failed, err:%v", err)
		return nil, err
	}

	if len(columnsResp.Tables) == 0 {
		return nil, errs.TableNotExist
	}

	columns := columnsResp.Tables[0]
	if req.ProjectId > 0 {
		projectInfo, err := domain.GetProjectSimple(req.OrgId, req.ProjectId)
		if err != nil {
			log.Errorf("[getTableColumns] GetProjectSimple failed, orgId:%v, projectId:%v, err:%v", req.OrgId, req.ProjectId, err)
			return nil, errs.ProjectNotExist
		}

		if projectInfo.ProjectTypeId == consts.ProjectTypeAgileId {
			for _, column := range columns.Columns {
				if column.Name == consts.BasicFieldIterationId {
					err = addIterationOptions(req.OrgId, req.ProjectId, column)
					break
				}
			}
		}
		columns.Columns = domain.RemoveColumnByProjectType(columns.Columns, projectInfo.ProjectTypeId)
	}

	return columns, err
}

// addProjectOptions 添加projectId的选项，加入所有project
func addProjectOptions(orgId, userId int64, columns []*projectvo.TableColumnData, tablesMap map[int64][]*lc_table.LcOptions) errs.SystemErrorInfo {
	isFilling := 3
	projectList, err := Projects(projectvo.ProjectsRepVo{
		Page: 0,
		Size: 0,
		ProjectExtraBody: projectvo.ProjectExtraBody{
			Input: &vo.ProjectsReq{
				IsFiling: &isFilling,
			},
			NoNeedRedundancyInfo: true,
			ProjectTypeIds:       []int64{consts.ProjectTypeNormalId, consts.ProjectTypeAgileId, consts.ProjectTypeCommon2022V47},
		},
		OrgId:         orgId,
		UserId:        userId,
		SourceChannel: "",
	})
	if err != nil {
		log.Error(err)
		return err
	}

	var projectOptions []lc_table.LcProjectOptions
	for _, project := range projectList.List {
		child := tablesMap[cast.ToInt64(project.AppID)]
		if child == nil {
			child = make([]*lc_table.LcOptions, 0)
		}
		projectOptions = append(projectOptions, lc_table.LcProjectOptions{
			LcOptions: lc_table.LcOptions{
				Id:    project.ID,
				Value: project.Name,
				Key:   consts.BasicFieldProjectId,
			},
			Children: child,
		})
	}

	for _, column := range columns {
		if column.Name == consts.BasicFieldProjectId {
			column.Field.Type = consts.DataTypeCascader
			column.Field.CustomType = consts.DataTypeCascader
			column.Label = "所属项目与表"
			column.Field.Props = map[string]interface{}{
				consts.DataTypeCascader: map[string]interface{}{"options": projectOptions},
			}
		}
	}

	return nil
}

// 添加所有表头的issueStatus选项
func addIssueStatusOptions(orgId, userId int64, columns []*projectvo.TableColumnData) (map[int64][]*lc_table.LcOptions, errs.SystemErrorInfo) {
	resp := tablefacade.ReadTableSchemasByOrgId(projectvo.GetTableSchemasByOrgIdReq{
		OrgId:     orgId,
		UserId:    userId,
		ColumnIds: []string{consts.BasicFieldIssueStatus},
	})
	if resp.Failure() {
		log.Errorf("[addIssueStatusOptions] tablefacade.ReadTableSchemasByOrgId failed, orgId:%v, err:%v", orgId, resp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, resp.Message)
	}

	tablesMap := make(map[int64][]*lc_table.LcOptions, len(resp.Data.Tables))
	issueStatusOptions := make([]interface{}, 0, len(resp.Data.Tables))
	for _, table := range resp.Data.Tables {
		tablesMap[table.AppId] = append(tablesMap[table.AppId],
			&lc_table.LcOptions{Id: cast.ToString(table.TableId), Value: table.Name, Key: consts.BasicFieldTableId})

		if len(table.Columns) >= 1 && table.Columns[0] != nil && table.Columns[0].Field.Props["groupSelect"] != nil {
			options := lc_table.ExchangeToLcGroupSelectOptions(table.Columns[0].Field.Props, cast.ToString(table.TableId))
			if len(options) > 0 {
				for i := range options {
					issueStatusOptions = append(issueStatusOptions, options[i])
				}
			}
		}
	}

	for _, tables := range tablesMap {
		sort.Slice(tables, func(i, j int) bool {
			return cast.ToString(tables[i].Id) < cast.ToString(tables[j].Id)
		})
	}

	for _, column := range columns {
		if column.Name == consts.BasicFieldIssueStatus {
			if column.Field.Props != nil && column.Field.Props["groupSelect"] != nil {
				if m, ok := column.Field.Props["groupSelect"].(map[string]interface{}); ok {
					m["options"] = issueStatusOptions
					break
				}
			}
		}
	}

	return tablesMap, nil
}

func GetTablesColumns(req projectvo.GetTablesColumnsReq) (*projectvo.TablesColumnsRespData, errs.SystemErrorInfo) {
	columnsResp, err := domain.GetTablesColumnsByTableIds(req.OrgId, req.UserId, req.Input.TableIds, req.Input.ColumnIds, false)
	if err != nil {
		log.Errorf("[GetOneTableColumns] failed, err:%v", err)
		return nil, err
	}

	if len(columnsResp.Tables) == 0 {
		return nil, errs.TableNotExist
	}

	return &projectvo.TablesColumnsRespData{Tables: columnsResp.Tables}, nil
}
