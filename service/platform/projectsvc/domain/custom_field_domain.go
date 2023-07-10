package domain

import (
	"strconv"

	"github.com/modern-go/reflect2"

	tableV1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/appvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/appfacade"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/tablefacade"
	consts2 "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func ChangeProjectCustomField(orgId, userId, projectId int64, fieldId int64, isAdd bool, projectObjectTypeId int64) errs.SystemErrorInfo {
	//防止项目成员重复插入
	uid := uuid.NewUuid()
	projectIdStr := strconv.FormatInt(projectId, 10)
	lockKey := consts.AddProjectCustomFieldLock + projectIdStr
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, err)
		return errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	} else {
		return errs.BuildSystemErrorInfo(errs.GetDistributedLockError)
	}

	isExist, err := mysql.IsExistByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:               orgId,
		consts.TcProjectId:           projectId,
		consts.TcIsDelete:            consts.AppIsNoDelete,
		consts.TcRelationType:        consts.IssueRelationTypeCustomField,
		consts.TcProjectObjectTypeId: projectObjectTypeId,
		consts.TcRelationId:          fieldId,
	})
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}

	if isExist {
		if isAdd {
			return nil
		} else {
			_, err := mysql.UpdateSmartWithCond(consts.TableProjectRelation, db.Cond{
				consts.TcOrgId:               orgId,
				consts.TcProjectId:           projectId,
				consts.TcIsDelete:            consts.AppIsNoDelete,
				consts.TcRelationType:        consts.IssueRelationTypeCustomField,
				consts.TcProjectObjectTypeId: projectObjectTypeId,
				consts.TcRelationId:          fieldId,
			}, mysql.Upd{
				consts.TcIsDelete: consts.AppIsDeleted,
				consts.TcUpdator:  userId,
			})
			if err != nil {
				log.Error(err)
				return errs.MysqlOperateError
			}
		}
	} else {
		if isAdd {
			id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
			if err != nil {
				log.Error(err)
				return err
			}

			insertErr := mysql.Insert(&po.PpmProProjectRelation{
				Id:                  id,
				OrgId:               orgId,
				ProjectId:           projectId,
				RelationId:          fieldId,
				RelationType:        consts.IssueRelationTypeCustomField,
				ProjectObjectTypeId: projectObjectTypeId,
				Creator:             userId,
			})
			if insertErr != nil {
				log.Error(insertErr)
				return errs.MysqlOperateError
			}
		} else {
			return nil
		}
	}

	return nil
}

func AddProjectCustomFieldBatch(orgId, userId, projectId int64, fieldIds []int64, projectObjectTypeId int64) errs.SystemErrorInfo {
	//防止项目成员重复插入
	uid := uuid.NewUuid()
	projectIdStr := strconv.FormatInt(projectId, 10)
	lockKey := consts.AddProjectCustomFieldLock + projectIdStr
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, err)
		return errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	} else {
		return errs.BuildSystemErrorInfo(errs.GetDistributedLockError)
	}

	infoPos := &[]po.PpmProProjectRelation{}

	selectErr := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:               orgId,
		consts.TcProjectId:           projectId,
		consts.TcIsDelete:            consts.AppIsNoDelete,
		consts.TcRelationType:        consts.IssueRelationTypeCustomField,
		consts.TcProjectObjectTypeId: projectObjectTypeId,
		consts.TcRelationId:          db.In(fieldIds),
	}, infoPos)
	if selectErr != nil {
		log.Error(selectErr)
		return errs.MysqlOperateError
	}

	existFieldIds := []int64{}
	for _, relation := range *infoPos {
		existFieldIds = append(existFieldIds, relation.RelationId)
	}
	_, addIds := util.GetDifMemberIds(existFieldIds, fieldIds)
	if len(addIds) == 0 {
		return nil
	}

	idResp, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, len(addIds))
	if idErr != nil {
		log.Error(idErr)
		return idErr
	}
	insertPos := []interface{}{}
	allFieldIds := []int64{}
	for i, id := range addIds {
		allFieldIds = append(allFieldIds, id)
		insertPos = append(insertPos, po.PpmProProjectRelation{
			Id:                  idResp.Ids[i].Id,
			OrgId:               orgId,
			ProjectId:           projectId,
			RelationId:          id,
			RelationType:        consts.IssueRelationTypeCustomField,
			ProjectObjectTypeId: projectObjectTypeId,
			Creator:             userId,
		})
	}
	insertErr := mysql.BatchInsert(&po.PpmProProjectRelation{}, insertPos)
	if insertErr != nil {
		log.Error(insertErr)
		return errs.MysqlOperateError
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{
			FieldIds:            allFieldIds,
			ProjectObjectTypeId: projectObjectTypeId,
		}
		projectTrendsBo := bo.ProjectTrendsBo{
			PushType:   consts.PushTypeUseOrgCustomField,
			OrgId:      orgId,
			ProjectId:  projectId,
			OperatorId: userId,
			Ext:        ext,
		}
		PushProjectTrends(projectTrendsBo)
	})

	return nil
}

func GetCustomFieldInfo(orgId, fieldId int64) (*bo.CustomFieldBo, errs.SystemErrorInfo) {
	info := &po.PpmProCustomField{}
	err := mysql.SelectOneByCond(consts.TableCustomField, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    []int64{0, orgId},
		consts.TcId:       fieldId,
	}, info)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.CustomFieldNotExist
		} else {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}
	}

	res := &bo.CustomFieldBo{}
	_ = copyer.Copy(info, res)

	return res, nil
}

func DeleteCustomField(orgId, userId, fieldId int64) errs.SystemErrorInfo {
	_, err := mysql.UpdateSmartWithCond(consts.TableCustomField, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcId:       fieldId,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  userId,
	})

	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}

	return nil
}

func GetProjectCustomField(orgId int64, projectIds []int64, isIgnoreStatus bool) ([]bo.CustomFieldBo, errs.SystemErrorInfo) {
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    db.In([]int64{0, orgId}),
		consts.TcId:       db.In(db.Raw("select relation_id from ppm_pro_project_relation where org_id = ? and relation_type = ? and status = 1 and project_id in ? and is_delete = 2", orgId, consts.IssueRelationTypeCustomField, projectIds)),
	}
	if isIgnoreStatus {
		cond[consts.TcId] = db.In(db.Raw("select relation_id from ppm_pro_project_relation where org_id = ? and relation_type = ? and project_id in ? and is_delete = 2", orgId, consts.IssueRelationTypeCustomField, projectIds))
	}

	pos := &[]po.PpmProCustomField{}
	err := mysql.SelectAllByCond(consts.TableCustomField, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.CustomFieldBo{}
	_ = copyer.Copy(pos, bos)

	return *bos, nil
}

func GetProjectCustomFieldIds(orgId int64, projectIds []int64) ([]int64, errs.SystemErrorInfo) {
	customFieldBos, err := GetProjectCustomField(orgId, projectIds, false)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	ids := make([]int64, 0)
	for _, customField := range customFieldBos {
		ids = append(ids, customField.Id)
	}
	return ids, nil
}

func GetIssuesWorkHourList(orgId int64, projectIds, issueIds []int64) ([]*bo.SimpleWorkHourBo, errs.SystemErrorInfo) {
	pos := []*po.SimpleWorkHourPo{}
	cond := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: db.In(projectIds),
		consts.TcIssueId:   db.In(issueIds),
		consts.TcType:      db.In([]int8{consts2.WorkHourTypeTotalPredict, consts2.WorkHourTypeSubPredict, consts2.WorkHourTypeActual}),
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}
	err := mysql.SelectAllByCond(consts.TableIssueWorkHours, cond, &pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	bos := []*bo.SimpleWorkHourBo{}
	_ = copyer.Copy(pos, &bos)

	return bos, nil
}

// GetCustomFieldInfoFromViewList 从表格视图中取出视图的列信息
func GetCustomFieldInfoFromViewList(orgId, appId, tableId int64, usingType string,
) ([]appvo.GetAppViewListRespDataListItem, *appvo.GetAppViewListRespDataListItemConfig, errs.SystemErrorInfo) {
	// func body start
	var oriErr error
	tableViewList := make([]appvo.GetAppViewListRespDataListItem, 0)
	viewConfigObj := &appvo.GetAppViewListRespDataListItemConfig{}
	// 获取“表格视图”
	// 根据前端展示字段列的规则，优先从“表格视图”中取出字段，如果找不到合适的视图
	// 敏捷项目：使用 config 中，objectType 对应的字段列表，如果没有 objectType 对应的数据，则使用默认的任务内置字段。
	// 通用项目：使用默认的任务内置字段。
	// 表头的顺序存放在视图中，因此要根据视图按顺序组装字段列表
	viewResp := appfacade.GetAppViewList(&appvo.GetAppViewListReq{
		OrgId: orgId,
		AppId: appId,
	})
	if viewResp.Failure() {
		log.Errorf("[GetCustomFieldInfoFromViewList] appfacade.GetAppViewList err: %v", viewResp.Error())
		return tableViewList, viewConfigObj, viewResp.Error()
	}
	viewList := viewResp.Data
	curTableIdStr := strconv.FormatInt(tableId, 10)
	for _, item := range viewList {
		// type 为 1 的表示表格类型的视图，用户自己定义的视图也是表格类型的视图。
		// 为了区分出“表格视图”，暂且只能用 name 进行判断。
		if item.Type != 1 && item.ViewName != "表格视图" {
			continue
		}
		tmpConfigObj := appvo.GetAppViewListRespDataListItemConfig{}
		//tmpViewId, _ := strconv.ParseInt(item.ID, 10, 64)
		// 有的 item.Config 是 json 字符串，有的又是对象，这里做个兼容。
		if configStr, ok := item.Config.(string); ok {
			oriErr = json.FromJson(configStr, &tmpConfigObj)
		} else if configMap, ok := item.Config.(map[string]interface{}); ok {
			configStr := json.ToJsonIgnoreError(configMap)
			oriErr = json.FromJson(configStr, &tmpConfigObj)
		} else {
			configStr := json.ToJsonIgnoreError(item.Config)
			oriErr = json.FromJson(configStr, &tmpConfigObj)
		}
		if oriErr != nil {
			log.Errorf("[GetCustomFieldInfoFromViewList] json.FromJson err: %v", oriErr)
			continue
		}
		// 1 表示表格视图
		if tmpConfigObj.TableId == curTableIdStr && item.Type == 1 {
			tableViewList = append(tableViewList, item)
			viewConfigObj = &tmpConfigObj
			// 一旦找到“表格视图”，则退出。
			break
		}
	}

	// 没有可用的“表格视图”，无法进行导入。但勉强可以进行任务导出。
	if len(tableViewList) < 1 && usingType == consts.FilterColumnUsingTypeImport {
		err := errs.TableViewNotExist
		log.Errorf("[GetCustomFieldInfoFromViewList] err: %v", err)
		return tableViewList, viewConfigObj, nil
	}

	return tableViewList, viewConfigObj, nil
}

// FilterCustomFieldForAgilePro 敏捷项目过滤隐藏字段
func FilterCustomFieldForAgilePro(columnKeyArr []string, columns []*projectvo.TableColumnData,
	allColumnsMap map[string]*projectvo.TableColumnData, needHidden bool, hiddenColumnKeys []interface{},
) []*projectvo.TableColumnData {
	// func body start
	sortedColumns := make([]*projectvo.TableColumnData, 0)
	if len(columnKeyArr) < 1 {
		columnKeyArr = make([]string, 0)
		// 把 columns 中的字段追加在后面。
		for _, column := range columns {
			columnKey := column.Name
			// 如果存在于隐藏字段列表里，则表明需要隐藏该字段
			if !needHidden && hiddenColumnKeys != nil && len(hiddenColumnKeys) > 0 {
				if isExist, _ := slice.Contain(hiddenColumnKeys, columnKey); isExist {
					continue
				}
			}
			if column.IsSys {
				continue
			}
			columnKeyArr = append(columnKeyArr, columnKey)
		}
	}
	for _, columnKey := range columnKeyArr {
		// 如果存在于隐藏字段列表里，则表明需要隐藏该字段
		if !needHidden && hiddenColumnKeys != nil && len(hiddenColumnKeys) > 0 {
			if isExist, _ := slice.Contain(hiddenColumnKeys, columnKey); isExist {
				continue
			}
		}
		if val, ok := allColumnsMap[columnKey]; ok {
			sortedColumns = append(sortedColumns, val)
		}
	}

	return sortedColumns
}

// FilterCustomFieldForCommonPro 通用项目过滤隐藏字段
func FilterCustomFieldForCommonPro(columnKeyArr []string, columns, originalColumns []*projectvo.TableColumnData,
	allColumnsMap map[string]*projectvo.TableColumnData, needHidden bool, hiddenColumnKeys []interface{},
) []*projectvo.TableColumnData {
	// func body start
	tableSortArr := make([]string, 0)
	sortedColumns := make([]*projectvo.TableColumnData, 0)
	// 如果表格视图中没提供列名，则直接使用 columns
	// if len(tableViewList) > 0 {
	if len(columnKeyArr) < 1 {
		for _, column := range originalColumns {
			tableSortArr = append(tableSortArr, column.Name)
		}
	}
	// 通用项目中，先取视图中的 tableOrder 中的字段，作为顺序，其余字段进行**追加**
	// tableSortArr 中可能出现多个一样的，因此对其去重
	tableSortArr = slice.SliceUniqueString(tableSortArr)
	for _, columnKey := range tableSortArr {
		// 如果存在于隐藏字段列表里，则表明需要隐藏该字段
		if !needHidden && hiddenColumnKeys != nil && len(hiddenColumnKeys) > 0 {
			if isExist, _ := slice.Contain(hiddenColumnKeys, columnKey); isExist {
				continue
			}
		}
		if val, ok := allColumnsMap[columnKey]; ok {
			sortedColumns = append(sortedColumns, val)
		}
	}
	// 把 columns 中的字段追加在后面。
	for _, column := range columns {
		columnKey := column.Name
		// 已经存在的字段，则略过
		if isExist, _ := slice.Contain(tableSortArr, columnKey); isExist {
			continue
		}
		// 如果存在于隐藏字段列表里，则表明需要隐藏该字段
		if !needHidden && hiddenColumnKeys != nil && len(hiddenColumnKeys) > 0 {
			if isExist, _ := slice.Contain(hiddenColumnKeys, columnKey); isExist {
				continue
			}
		}
		if column.IsSys {
			continue
		}
		sortedColumns = append(sortedColumns, column)
	}

	return sortedColumns
}

// FilterForBuiltInAndCustomColumns 对任务的字段进行一些过滤
func FilterForBuiltInAndCustomColumns(columns []*projectvo.TableColumnData, proTypeId int64, usingType string, isNeedDocument bool,
) ([]*projectvo.TableColumnData, []*projectvo.TableColumnData, []*projectvo.TableColumnData, errs.SystemErrorInfo) {
	// func body start
	// 对列进行一些特殊 filter
	builtInColumns := make([]*projectvo.TableColumnData, 0)
	customColumns := make([]*projectvo.TableColumnData, 0)
	notSupportCustomFieldTypes := make([]string, 0, len(consts2.NotSupportImportCustomFieldTypes))
	notSupportFieldKeys := make([]string, 0)
	extraRemoveColumn := make([]string, 0)
	if usingType == consts.FilterColumnUsingTypeImport {
		notSupportCustomFieldTypes = append(notSupportCustomFieldTypes, consts2.NotSupportImportCustomFieldTypes...)
		notSupportFieldKeys = consts2.HiddenColumnsInImportTpl
		extraRemoveColumn = []string{consts.BasicFieldCode}
	} else {
		notSupportCustomFieldTypes = append(notSupportCustomFieldTypes, consts2.NotSupportExportCustomFieldTypes...)
		if !isNeedDocument {
			notSupportCustomFieldTypes = append(notSupportCustomFieldTypes, consts.LcColumnFieldTypeImage, consts.LcColumnFieldTypeDocument)
		}
	}

	columns = RemoveColumnByProjectType(columns, proTypeId, extraRemoveColumn...)
	for _, column := range columns {
		columnKey := column.Name
		fieldType := column.Field.Type
		// 暂时不支持导入的字段类型
		if isExist, _ := slice.Contain(notSupportCustomFieldTypes, fieldType); isExist {
			continue
		}
		if isExist, _ := slice.Contain(notSupportFieldKeys, columnKey); isExist {
			continue
		}

		if exist, _ := slice.Contain(consts.DefaultField, columnKey); exist {
			builtInColumns = append(builtInColumns, column)
		} else {
			customColumns = append(customColumns, column)
		}
	}

	return columns, builtInColumns, customColumns, nil
}

// GetColumnMapByColumns 通过字段信息列表组装成对应的 map
func GetColumnMapByColumns(columns []map[string]interface{}) map[string]map[string]interface{} {
	columnsMap := make(map[string]map[string]interface{}, 0)
	for _, column := range columns {
		columnKey := column["key"].(string)
		columnsMap[columnKey] = column
	}

	return columnsMap
}

// GetOneOptionOfPriority 从优先级的 option 中，选出一个作为默认选中的 option
func GetOneOptionOfPriority(optionMap map[string]projectvo.FormConfigColumnFieldPropsOneVal) projectvo.FormConfigColumnFieldPropsOneVal {
	if opValue, ok := optionMap["普通"]; ok {
		return opValue
	} else {
		for _, value := range optionMap {
			return value
		}
	}

	return projectvo.FormConfigColumnFieldPropsOneVal{}
}

func FilterEmptyColumnKey(columnKeyArr []string) []string {
	columnKeyArr = slice.SliceUniqueString(columnKeyArr)
	newSlice := make([]string, 0, len(columnKeyArr))
	for _, val := range columnKeyArr {
		if val == "" {
			continue
		}
		newSlice = append(newSlice, val)
	}

	return newSlice
}

// GetColumnMapByColumnsForTableColumn 通过字段信息列表组装成对应的 map，只不过列信息是改造后查询到的列信息
func GetColumnMapByColumnsForTableColumn(columns []*projectvo.TableColumnData) map[string]*projectvo.TableColumnData {
	columnsMap := make(map[string]*projectvo.TableColumnData, len(columns))
	for _, column := range columns {
		columnKey := column.Name
		columnsMap[columnKey] = column
	}

	return columnsMap
}

// GetTableColumnsForImport 导入导出-获取表头信息
func GetTableColumnsForImport(orgId, tableId int64) ([]*projectvo.TableColumnData, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	columns := make([]*projectvo.TableColumnData, 0)
	columnPropsMap := make(map[string]map[string]interface{}, 0)
	resp := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
		OrgId:  orgId,
		UserId: 0,
		Input: &tableV1.ReadTableSchemasRequest{
			TableIds:           []int64{tableId},
			ColumnIds:          nil, // 不传表示获取表内的所有列信息
			IsNeedDescription:  true,
			IsNeedCommonColumn: false,
		},
	})
	if resp.Failure() {
		log.Errorf("[GetTableColumnsForImport] GetTableColumns err: %v", resp.Error())
		return columns, columnPropsMap, resp.Error()
	}
	for _, item := range resp.Data.Tables {
		if item.TableId == tableId {
			columns = item.Columns
			break
		}
	}
	if len(columns) < 1 {
		return columns, columnPropsMap, nil
	}
	columnPropsMap = GetTableColumnsPropsMap(columns)

	return columns, columnPropsMap, nil
}

// GetTableColumnsPropsMap 将字段列表中的 props 统一放到一个 map 中，方便取用
func GetTableColumnsPropsMap(columns []*projectvo.TableColumnData) map[string]map[string]interface{} {
	// propsMap: {columnKey: obj}
	propsMap := make(map[string]map[string]interface{}, len(columns))
	for _, column := range columns {
		if !reflect2.IsNil(column.Field.Props) {
			propsMap[column.Name] = column.Field.Props
		}
	}

	return propsMap
}
