package domain

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	consts2 "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"

	"github.com/araddon/dateparse"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/excel/sortedmap"
	"github.com/star-table/polaris-backend/common/core/util/format"
	lang2 "github.com/star-table/polaris-backend/common/core/util/lang"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/permissionfacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/tealeg/xlsx/v2"
	"github.com/xuri/excelize/v2"
)

// ImportIssue 任务导入
type ImportIssue struct {
	InAndExportBase
	FilePath             string
	ColumnKeyRowTh       int // 列 key 所在的行；索引值是从 1 开始
	DataStartRowTh       int // excel 中正式数据开始的行；索引值是从 1 开始
	ColumnMaps           map[string]*projectvo.TableColumnData
	FileColumnsArr       []*projectvo.TableColumnData // 当前导入的 excel 的列信息
	UnknownImportColumns []string
	Rows                 [][]string // 当前 excel 文件解析出的数据，包含表头等行
	OptionsMap           map[string]map[string]*projectvo.ColumnSelectOption
	// 用户信息
	UserInfoMap        map[string]bo.SimpleUserWithUserIdBo `json:"userInfoMap"`
	UserIdMapKeyByName map[string]int64                     `json:"userIdMapKeyByName"`
	// 重名的用户
	DuplicateUserNames []string `json:"duplicateUserNames"`
	// 重名的部门
	DuplicateDeptNames []string `json:"duplicateDeptNames"`
	DeptMap            map[string]orgvo.SimpleDeptForCustomField
	AsyncTaskId        string `json:"asyncTaskId"`

	ErrMsgList map[int][]string
}

func NewImportIssue(orgId, projectId, tableId int64, filePath string) (*ImportIssue, errs.SystemErrorInfo) {
	importIssue := &ImportIssue{
		InAndExportBase: InAndExportBase{
			OrgId: orgId,
		},
		FilePath:       filePath,
		ColumnKeyRowTh: 2,
		DataStartRowTh: 4,
		ColumnMaps:     make(map[string]*projectvo.TableColumnData),
		OptionsMap:     make(map[string]map[string]*projectvo.ColumnSelectOption),
		ErrMsgList:     make(map[int][]string),
	}

	err := importIssue.SetProjectAndTableInfo(orgId, projectId, tableId)
	if err != nil {
		return nil, err
	}

	err = importIssue.PrepareImportInfo()
	if err != nil {
		return nil, err
	}

	return importIssue, nil
}

func (ii *ImportIssue) SetColumns(columns []*projectvo.TableColumnData) errs.SystemErrorInfo {
	err := ii.InAndExportBase.SetColumns(columns)
	if err != nil {
		return err
	}

	for _, column := range columns {
		ii.ColumnMaps[column.Name] = column
		options := ii.GetColumnOptions(column)
		for _, option := range options {
			if ii.OptionsMap[column.Name] == nil {
				ii.OptionsMap[column.Name] = make(map[string]*projectvo.ColumnSelectOption)
			}
			ii.OptionsMap[column.Name][option.Value] = option
			if lanMap, ok := consts.AllTranslate[lang2.GetLang()]; ok && lanMap[option.Value] != "" {
				ii.OptionsMap[column.Name][lanMap[option.Value]] = option
			}
		}
	}

	return nil
}

// ParseExcel 解析导入任务场景的 excel
// 基于 excelize 库：https://github.com/qax-os/excelize  https://xuri.me/excelize
// @param dataList 返回解析后的数据集合
// @param ignoreRowNum 忽略的行数。两种情况：1.第一行的表头需要忽略；2.存在注意事项行也需要忽略；
func (ii *ImportIssue) ParseExcel() (dataList interface{}, ignoreRowNum int, headerMap *sortedmap.SortMap, businessErr errs.SystemErrorInfo) {
	defaultSheet := "Sheet1"
	ignoreRowNum = ii.DataStartRowTh - 1
	fd, err := excelize.OpenFile(ii.FilePath, excelize.Options{
		RawCellValue: true,
	})
	if err != nil {
		businessErr = errs.ReadExcelFailed
		log.Errorf("[ParseExcelForImportIssue] err: %v", errs.BuildSystemErrorInfo(errs.ReadExcelFailed, err))
		return
	}
	defer func() {
		if err := fd.Close(); err != nil {
			log.Errorf("[ParseExcelForImportIssue] fd.Close err: %v", err)
		}
	}()
	ii.Fd = fd
	// 获取当前表的列信息
	if _, businessErr = ii.GetColumnInfoInTable(); businessErr != nil {
		log.Errorf("[ParseExcelForImportIssue] err: %v", businessErr)
		return
	}
	// 获取 Sheet1 上所有行
	rows, err := fd.GetRows(defaultSheet, excelize.Options{
		RawCellValue: true,
	})
	if err != nil {
		businessErr = errs.ReadExcelFailed
		log.Errorf("[ParseExcelForImportIssue] err: %v", errs.BuildSystemErrorInfo(errs.ReadExcelFailed, err))
		return
	}
	ii.Rows = rows
	// 解析 header key
	headerMap = ii.ParseColumnKey()
	// 如果 ”key 列“没有任务状态（issueStatus），则视为不是新模板，此时用户需要下载新模板
	if len(ii.FileColumnsArr) < 1 || headerMap.GetByKey(consts.BasicFieldTitle, "") == "" {
		// 如果一行合法字段都匹配不上，则不能导入。说明不是下载模板进行导入
		businessErr = errs.HasNoMatchedColumn
		log.Errorf("[ParseExcelForImportIssue] err: %v", businessErr)
		return
	}
	// 解析任务数据
	dataList = ii.ParseDataList()

	return
}

func (ii *ImportIssue) GetColumnInfoInTable() (map[string]*projectvo.TableColumnData, errs.SystemErrorInfo) {
	columnMap := make(map[string]*projectvo.TableColumnData, 0)
	for _, column := range ii.Columns {
		columnMap[column.Name] = column
	}
	ii.ColumnMaps = columnMap
	// 追加一个特殊的列：父子任务。这个列不在表中
	ii.ColumnMaps[consts.HelperFieldIsParentDesc] = consts2.ParentDescColumn

	return columnMap, nil
}

// ParseColumnKey 解析列 key 行，返回对应的列信息数组
func (ii *ImportIssue) ParseColumnKey() *sortedmap.SortMap {
	var keyRow []string
	var nameRow []string

	for ri, row := range ii.Rows {
		// 只有第 2 行才是列的 key
		if ri == ii.ColumnKeyRowTh-1 {
			keyRow = row

			// 保持和 excel 文件的列一致
			ii.FileColumnsArr = make([]*projectvo.TableColumnData, len(row))
		}
		if ri == ii.ColumnKeyRowTh {
			nameRow = row
		}

		if keyRow != nil && nameRow != nil {
			break
		}
	}

	headMap := sortedmap.NewSortMap()
	for i, cellVal := range keyRow {
		cellVal = strings.TrimSpace(cellVal)
		if oneColumn, ok := ii.ColumnMaps[cellVal]; ok {
			headMap.Insert(oneColumn.Name, oneColumn.Label)
			ii.FileColumnsArr[i] = oneColumn
		} else {
			if i < len(nameRow) {
				ii.UnknownImportColumns = append(ii.UnknownImportColumns, nameRow[i])
				ii.FileColumnsArr[i] = nil
			}
		}
	}

	return headMap
}

func (ii *ImportIssue) ParseDataList() []*sortedmap.SortMap {
	entryArr := make([]*sortedmap.SortMap, 0, len(ii.Rows))
	for ri, row := range ii.Rows {
		// ri 从 0 开始
		if ri < ii.DataStartRowTh-1 {
			continue
		}
		entry := sortedmap.NewSortMap()
		for ci, cell := range row {
			if ci >= len(ii.FileColumnsArr) {
				// 多出的列值忽略，不做解析
				break
			}
			// 获取 cell 对应的列 key
			curColumnInfo := ii.FileColumnsArr[ci]
			if curColumnInfo == nil {
				continue
			}
			columnKey := curColumnInfo.Name
			entry.Insert(columnKey, strings.TrimSpace(cell))
		}
		entryArr = append(entryArr, entry)
	}

	return entryArr
}

// SetImportColumnValueForTitle 导入任务时，处理标题
func (ii *ImportIssue) SetImportColumnValueForTitle(data map[string]interface{}, cellValue string, colIndexMap map[string]int) {
	issueTitle := cellValue
	resTitle := strings.TrimSpace(issueTitle)
	isTitleRight := format.VerifyIssueNameFormat(resTitle) && len(resTitle) > 0
	if !isTitleRight {
		colNumTh := colIndexMap[consts.BasicFieldTitle]
		ii.ErrMsgList[colNumTh+1] = append(ii.ErrMsgList[colNumTh+1], errs.IssueTitleCannotEmpty.Message())
	}
	data[consts.BasicFieldTitle] = resTitle
}

// SetImportColumnValueForStatus 导入任务时，设置任务状态。导入时，将 excel 的 cell 中的值设置到任务对象中
func (ii *ImportIssue) SetImportColumnValueForStatus(data map[string]interface{}, cellValue string) {
	resStatus := int64(0)
	resAuditStatus := 0
	auditorIds := cast.ToStringSlice(data[consts.BasicFieldAuditorIds])
	optionsMap := ii.OptionsMap[consts.BasicFieldIssueStatus]
	issueStatus := cellValue

	if optionsMap[issueStatus] != nil {
		resStatus = cast.ToInt64(optionsMap[issueStatus].Id)
	} else {
		// 待确认的状态处理
		isWaitConfirmStr := issueStatus == consts.WaitConfirmStatusName
		if lang2.IsEnglish() {
			isWaitConfirmStr = issueStatus == consts.WaitConfirmStatusEnName
		}
		if ii.ProjectInfo.ProjectTypeId == consts.ProjectTypeNormalId && isWaitConfirmStr {
			if len(auditorIds) > 0 {
				resAuditStatus = consts.AuditStatusNotView
			} else {
				resAuditStatus = consts.AuditStatusPass
			}
			// 通过 type 获取已完成状态
			for _, option := range optionsMap {
				if cast.ToInt64(option.ParentId) == consts.StatusTypeComplete {
					resStatus = cast.ToInt64(option.Id)
					break
				}
			}
		}
	}

	data[consts.BasicFieldIssueStatus] = resStatus
	if resAuditStatus > 0 {
		data[consts.BasicFieldAuditStatus] = resAuditStatus
	}
}

// SetPlanTimesForImport 导入任务时，处理起止时间的匹配
func (ii *ImportIssue) SetPlanTimesForImport(data map[string]interface{}, row *sortedmap.SortMap, colIndexMap map[string]int) {
	var err error
	planStartTime := consts.BlankTimeObject
	planEndTime := consts.BlankTimeObject

	// 计划开始时间
	planStartTimeStr := strings.TrimSpace(row.GetByKey(consts.BasicFieldPlanStartTime, ""))
	if planStartTimeStr != "" {
		if planStartTime, err = ParseExcelDatetime(planStartTimeStr); err != nil {
			msg := "开始时间不正确。" + err.Error()
			log.Error(msg)
			colNumTh := colIndexMap[consts.BasicFieldPlanStartTime]
			ii.ErrMsgList[colNumTh+1] = append(ii.ErrMsgList[colNumTh+1], consts.ImportUserErrOfStartTimeStyleErr)
		}
	}

	// 计划结束时间
	planEndTimeStr := strings.TrimSpace(row.GetByKey(consts.BasicFieldPlanEndTime, ""))
	if planEndTimeStr != "" {
		if planEndTime, err = ParseExcelDatetime(planEndTimeStr); err != nil {
			msg := "截止时间不正确。" + err.Error()
			log.Error(msg)
			colNumTh := colIndexMap[consts.BasicFieldPlanEndTime]
			ii.ErrMsgList[colNumTh+1] = append(ii.ErrMsgList[colNumTh+1], consts.ImportUserErrOfEndTimeLessErr)
		}
	}

	// 对于起止时间都填写了，才会比较截止时间是否位于开始时间之后
	if planStartTime.After(consts.BlankTimeObject) &&
		planEndTime.After(consts.BlankTimeObject) &&
		planEndTime.Before(planStartTime) {
		colNumTh := colIndexMap[consts.BasicFieldPlanEndTime]
		ii.ErrMsgList[colNumTh+1] = append(ii.ErrMsgList[colNumTh+1], consts.ImportUserErrOfEndTimeLessErr)
	}

	if planStartTime.After(consts.BlankTimeObject) {
		if column, ok := ii.ColumnMaps[consts.BasicFieldPlanStartTime]; ok {
			log.Infof("[SetPlanTimesForImport] column: %v, time: %v", json.ToJsonIgnoreError(column), planStartTime)
			if !cast.ToBool(column.Field.Props["showTime"]) { // 未打开时间显示
				if planStartTime.Hour() == 0 && planStartTime.Minute() == 0 && planStartTime.Second() == 0 { // 未设置时分秒
					planStartTime = time.Date(planStartTime.Year(), planStartTime.Month(), planStartTime.Day(), 9, 0, 0, 0, planStartTime.Location()) // 设置成09:00:00
				}
			}
		}
		data[consts.BasicFieldPlanStartTime] = planStartTime.Format(consts.AppTimeFormat)
	}
	if planEndTime.After(consts.BlankTimeObject) {
		if column, ok := ii.ColumnMaps[consts.BasicFieldPlanEndTime]; ok {
			log.Infof("[SetPlanTimesForImport] column: %v, time: %v", json.ToJsonIgnoreError(column), planEndTime)
			if !cast.ToBool(column.Field.Props["showTime"]) { // 未打开时间显示
				if planEndTime.Hour() == 0 && planEndTime.Minute() == 0 && planEndTime.Second() == 0 { // 未设置时分秒
					planEndTime = time.Date(planEndTime.Year(), planEndTime.Month(), planEndTime.Day(), 18, 0, 0, 0, planEndTime.Location()) // 设置成18:00:00
				}
			}
		}
		data[consts.BasicFieldPlanEndTime] = planEndTime.Format(consts.AppTimeFormat)
	}
}

// FilterEmptyRowForImportIssue 导入任务时，检查每一行是否有数据，如果是空行则忽略
func FilterEmptyRowForImportIssue(groupData []*sortedmap.SortMap) ([]*sortedmap.SortMap, errs.SystemErrorInfo) {
	//手动处理空行
	importData := make([]*sortedmap.SortMap, 0)
	for _, v := range groupData {
		tempLen := len(v.DataMap)
		// 检查每一个 cell 是否是空
		for _, val := range v.DataMap {
			if val == "" {
				tempLen--
			}
		}
		// 如果一行数据中，全部是空 cell，则跳过这一行
		if tempLen > 0 {
			importData = append(importData, v)
		}
	}
	if len(importData) == 0 {
		return importData, errs.ImportDataEmpty
	}
	if len(importData) > consts.MaxRawSupportForImportIssue {
		return importData, errs.TooLargeImportData
	}

	return importData, nil
}

// PrepareImportInfo 导入任务时，组装用于匹配任务信息的 map 数据
func (ii *ImportIssue) PrepareImportInfo() errs.SystemErrorInfo {
	var err errs.SystemErrorInfo

	// 人员信息。需针对同名部分做处理。
	ii.UserInfoMap, ii.UserIdMapKeyByName, ii.DuplicateUserNames, err = GetUserMapForImportIssue(ii.OrgId)
	if err != nil {
		log.Errorf("[PrepareImportInfo] GetUserMapForImportIssue orgId:%v, err: %v", ii.OrgId, err)
		return err
	}
	ii.DeptMap, err = GetDeptListForImportIssue(ii.OrgId)
	if err != nil {
		log.Errorf("[PrepareImportInfo] GetDeptListForImportIssue orgId:%v, err: %v", ii.OrgId, err)
		return err
	}

	ii.DuplicateDeptNames, err = GetSameNameDeptNames(ii.OrgId)
	if err != nil {
		log.Errorf("[PrepareImportInfo] GetSameNameDeptNames err: %v", err)
		return err
	}

	// 导入任务时，生成一个异步任务id
	ii.AsyncTaskId = GenAsyncTaskIdForImport(ii.ProjectInfo.AppId, ii.TableInfo.TableId)

	return nil
}

// SetDeptDataForImport 导入时，将单元格中的“部门名”转换为部门 id 值
func (ii *ImportIssue) SetDeptDataForImport(data map[string]interface{}, cellValue string, column *projectvo.TableColumnData) errs.SystemErrorInfo {
	deptIdArr := make([]int64, 0)
	resultDeptArr := make([]orgvo.SimpleDeptForCustomField, 0)
	// 获取字段的一些属性
	propsStr := json.ToJsonIgnoreError(column.Field.Props)
	propsObj := projectvo.FormConfigColumnFieldDeptProps{}
	json.FromJson(propsStr, &propsObj)

	deptNames := strings.Split(cellValue, "||")
	deptNames = slice.SliceUniqueString(deptNames)
	for _, tmpName := range deptNames {
		// 如果值同名，则提示用户使用“推荐样式”填值
		if exist, _ := slice.Contain(ii.DuplicateDeptNames, tmpName); exist {
			err := errs.DeptDuplicateWhenImport
			log.Errorf("[TransferCustomColumnOfDept] err: %v", err)
			return err
		}
		if val, ok := ii.DeptMap[tmpName]; ok {
			deptIdArr = append(deptIdArr, val.ID)
			resultDeptArr = append(resultDeptArr, val)
			// 如果不是多选，则 break，只取一条
			if propsObj.Dept.Multiple != true {
				break
			}
		}
	}
	data[column.Name] = deptIdArr

	return nil
}

// SetMemberDataForImport 导入时，成员类型值的转换。包含创建人、关注人 等等
func (ii *ImportIssue) SetMemberDataForImport(data map[string]interface{}, cellValue string, column *projectvo.TableColumnData) errs.SystemErrorInfo {
	var resVal interface{}
	field := column.Field
	// 获取字段的一些属性
	propsStr := json.ToJsonIgnoreError(field.Props)
	propsObj := projectvo.FormConfigColumnFieldMemberProps{}
	json.FromJson(propsStr, &propsObj)
	valUserIdStrArr := make([]string, 0)
	valUserArr := make([]bo.SimpleUserWithUserIdBo, 0)
	idArr := make([]int64, 0)
	// 多个成员用 `||` 分隔
	userNames := strings.Split(cellValue, "||")
	userNames = slice.SliceUniqueString(userNames)
	for _, tmpName := range userNames {
		// 如果值同名，则提示用户使用“推荐样式”填值
		if exist, _ := slice.Contain(ii.DuplicateUserNames, tmpName); exist {
			err := errs.MemberDuplicateWhenImport
			log.Errorf("[SetMemberDataForImport] err: %v", err)
			return err
		}
		if simpleUser, ok := ii.UserInfoMap[tmpName]; !ok {
			continue
		} else {
			valUserIdStrArr = append(valUserIdStrArr, fmt.Sprintf("U_%v", simpleUser.Id))
			valUserArr = append(valUserArr, simpleUser)
			idArr = append(idArr, simpleUser.Id)
			// 如果不是多选，则 break，只取一条
			// if propsObj.Member.Multiple != true {
			if !(propsObj.Multiple || propsObj.Member.Multiple) {
				break
			}
		}
	}
	// 创建者、更新者直接存储 id
	if propsObj.Member.IsUseCreator || propsObj.Member.IsUseUpdator {
		if valUserIdStrArr != nil && len(valUserIdStrArr) > 0 {
			resVal = idArr[0]
		}
	} else {
		resVal = valUserIdStrArr
	}
	data[column.Name] = resVal

	return nil
}

// ParseCustomColumnMobileForImport 导入时，转换单元格中的“手机号”值
func (ii *ImportIssue) ParseCustomColumnMobileForImport(cellValue string) string {
	// 如果是以 `+86` 开头的，则直接显示，如果没有，则补上 `+86`
	defaultPrefix := "+86"
	if !strings.HasPrefix(cellValue, defaultPrefix) {
		cellValue = defaultPrefix + " " + cellValue
	}

	return cellValue
}

// SetGroupSelectDataForImport 导入时，转换单元格中的“分组单选”值
func (ii *ImportIssue) SetGroupSelectDataForImport(data map[string]interface{}, column *projectvo.TableColumnData, cellValue string,
	colIndexMap map[string]int,
) errs.SystemErrorInfo {
	var resVal interface{}
	optionsMap := ii.OptionsMap[column.Name]
	if targetOption, ok := optionsMap[cellValue]; ok {
		resVal = targetOption.Id
	} else {
		// 分组单选：保持 form 的规则一致，如果无法匹配，则直接报错提示
		errMsg := fmt.Sprintf("分组单选值[%s]不匹配。", cellValue)
		colNumTh := colIndexMap[column.Name]
		ii.ErrMsgList[colNumTh+1] = append(ii.ErrMsgList[colNumTh+1], errMsg)
	}
	data[column.Name] = resVal

	return nil
}

// SetMultiSelectDataForImport 导入时，转换单元格中的“多选”值
func (ii *ImportIssue) SetMultiSelectDataForImport(data map[string]interface{}, cellValue string, column *projectvo.TableColumnData) {
	// 多选，value 是 “值 id 的数组json字符串”，如："[\"4cfc87be-7881-6e77-2888-0589beb3260f\",\"840f7e30-75e2-edfa-b6fb-1cf6c8fb41fb\"]"
	optionsMap := ii.OptionsMap[column.Name]
	valueList := make([]interface{}, 0)
	titles := strings.Split(cellValue, "||")
	for _, oneVal := range titles {
		valueObj, isExists := optionsMap[oneVal]
		if !isExists {
			continue
		}
		valueList = append(valueList, valueObj.Id)
	}
	data[column.Name] = valueList
}

func (ii *ImportIssue) SetSelectDataForImport(data map[string]interface{}, cellValue string, column *projectvo.TableColumnData) {
	optionsMap := ii.OptionsMap[column.Name]
	if o, ok := optionsMap[cellValue]; ok {
		data[column.Name] = o.Id
	}
}

// SetPickerDataForImport 导入时，转换单元格中的“日期时间”值
func (ii *ImportIssue) SetPickerDataForImport(data map[string]interface{}, cellValue string, column *projectvo.TableColumnData, colIndexMap map[string]int) {
	datetime1, err := ParseExcelDatetime(cellValue)
	if err != nil {
		log.Error(err)
		colNumTh := colIndexMap[column.Name]
		ii.ErrMsgList[colNumTh+1] = append(ii.ErrMsgList[colNumTh+1], consts.ImportIssueErrOfDateTime)
	} else {
		log.Infof("[SetPickerDataForImport] column: %v, time: %v", json.ToJsonIgnoreError(column), datetime1)
		if datetime1.After(consts.BlankTimeObject) {
			if !cast.ToBool(column.Field.Props["showTime"]) { // 未打开时间显示
				if datetime1.Hour() == 0 && datetime1.Minute() == 0 && datetime1.Second() == 0 { // 未设置时分秒
					datetime1 = time.Date(datetime1.Year(), datetime1.Month(), datetime1.Day(), 18, 0, 0, 0, datetime1.Location()) // 设置成18:00:00
				}
			}
		}
		data[column.Name] = datetime1.Format(consts.AppTimeFormat)
	}
}

func (ii *ImportIssue) HandleImportData(row *sortedmap.SortMap, colIndexMap map[string]int) (map[string]interface{}, errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	data := map[string]interface{}{}
	// 特殊处理计划开始时间和结束时间
	ii.SetPlanTimesForImport(data, row, colIndexMap)
	if auditorIds := cast.ToString(row.DataMap[consts.BasicFieldAuditorIds]); auditorIds != "" && ii.ColumnMaps[consts.BasicFieldAuditorIds] != nil {
		_ = ii.SetMemberDataForImport(data, auditorIds, ii.ColumnMaps[consts.BasicFieldAuditorIds])
	}

	for columnId, value := range row.DataMap {
		if columnId == consts.BasicFieldPlanStartTime || columnId == consts.BasicFieldPlanEndTime ||
			columnId == consts.HelperFieldIsParentDesc || columnId == consts.BasicFieldAuditorIds {
			continue
		}
		if column, ok := ii.ColumnMaps[columnId]; ok {
			valueStr := cast.ToString(value)
			if valueStr == "" {
				continue
			}
			switch column.Field.Type {
			case consts.LcColumnFieldTypeMember:
				err = ii.SetMemberDataForImport(data, valueStr, column)
			case consts.LcColumnFieldTypeDept:
				err = ii.SetDeptDataForImport(data, valueStr, column)
			case consts.LcColumnFieldTypeAmount:
				ii.SetAmountDataForImport(data, valueStr, column)
			case consts.LcColumnFieldTypeInputNumber:
				ii.SetNumberDataForImport(data, valueStr, column)
			case consts.LcColumnFieldTypeSelect:
				ii.SetSelectDataForImport(data, valueStr, column)
			case consts.LcColumnFieldTypeMultiSelect:
				ii.SetMultiSelectDataForImport(data, valueStr, column)
			case consts.LcColumnFieldTypeDatepicker:
				ii.SetPickerDataForImport(data, valueStr, column, colIndexMap)
			case consts.LcColumnFieldTypeGroupSelect:
				if columnId == consts.BasicFieldIssueStatus {
					ii.SetImportColumnValueForStatus(data, valueStr)
				} else {
					err = ii.SetGroupSelectDataForImport(data, column, valueStr, colIndexMap)
				}
			default:
				if columnId == consts.BasicFieldTitle {
					ii.SetImportColumnValueForTitle(data, valueStr, colIndexMap)
				} else {
					data[columnId] = valueStr
				}
			}
		}
	}

	data[consts.TempFieldChildren] = []map[string]interface{}{}
	data[consts.TempFieldChildrenCount] = 0
	//if data[consts.BasicFieldIterationId] == nil && ii.ProjectInfo.ProjectTypeId == consts.ProjectTypeAgileId {
	//	data[consts.BasicFieldIterationId] = 0
	//}

	return data, err
}

// SetAmountDataForImport 自定义字段（货币）值的转换：将渲染后的值转换为无码字段的值。如：excel 中的 `$1.45`，根据规则将其转化为：`1.45`
func (ii *ImportIssue) SetAmountDataForImport(data map[string]interface{}, cellVal string, column *projectvo.TableColumnData) {
	// 如果是货币类型的自定义字段，则需要将单位去除即可
	propsStr := json.ToJsonIgnoreError(column.Field.Props)
	propsObj := projectvo.FormConfigColumnFieldAmountProps{}
	json.FromJson(propsStr, &propsObj)
	unitFlag := propsObj.Amount.Sign
	resultVal1 := strings.Trim(cellVal, unitFlag)
	// 根据规则，转换字段的精度
	resultVal2, err := strconv.ParseFloat(resultVal1, 64)
	if err != nil {
		// 如果异常，先见到那处理，赋值为 0
		log.Errorf("[SetAmountDataForImport] ParseFloat err: %v", err)
		resultVal2 = 0
	}
	accuracyNum := 0
	switch propsObj.Amount.Accuracy {
	case "1.0":
		accuracyNum = 1
	case "1.00":
		accuracyNum = 2
	case "1.000":
		accuracyNum = 3
	case "1.0000":
		accuracyNum = 4
	default:
	}
	data[column.Name] = format.Round(resultVal2, accuracyNum)
}

// SetNumberDataForImport 导入时，将数值类型的值转换为无码存储的值
func (ii *ImportIssue) SetNumberDataForImport(data map[string]interface{}, cellVal string, column *projectvo.TableColumnData) {
	propsStr := json.ToJsonIgnoreError(column.Field.Props)
	propsObj := projectvo.FormConfigColumnFieldInputNumberProps{}
	json.FromJson(propsStr, &propsObj)
	// 去除百分号
	if propsObj.Inputnumber.Percentage {
		cellVal = strings.Replace(cellVal, "%", "", 1)
	}
	// 千分位的去除
	if propsObj.Inputnumber.Thousandth {
		cellVal = strings.ReplaceAll(cellVal, ",", "")
	}
	// 精度
	accuracyNum := 0
	switch propsObj.Inputnumber.Accuracy {
	case "1.0":
		accuracyNum = 1
	case "1.00":
		accuracyNum = 2
	case "1.000":
		accuracyNum = 3
	case "1.0000":
		accuracyNum = 4
	default:
		accuracyNum = 0
	}
	dec, _ := decimal.NewFromString(cellVal)
	valFloat, _ := dec.Round(int32(accuracyNum)).Float64()
	data[column.Name] = valFloat
}

type ImportIssueTemplate struct {
	InAndExportBase
	Lang           string // 语言。用于多语言适配。en:英文；zh：中文
	IsEnglish      bool
	DataStartRowTh int // 正式数据开始的行。头部、注意事项行不计算在内。一般，导入模板文件中，注意事项行是第1行。
	DataEndRowTh   int
}

func NewImportIssueTemplate(orgId, projectId, tableId int64) (*ImportIssueTemplate, errs.SystemErrorInfo) {
	defaultSheet := "Sheet1"
	importIssueTemplate := &ImportIssueTemplate{
		InAndExportBase: InAndExportBase{
			OrgId:       orgId,
			ActiveSheet: defaultSheet,
		},
		Lang:           lang2.GetLang(),
		IsEnglish:      lang2.IsEnglish(),
		DataStartRowTh: 4,
		DataEndRowTh:   4 + consts.MaxRawSupportForImportIssue,
	}
	newFd := excelize.NewFile()
	newFd.NewSheet(defaultSheet)
	newFd.SetActiveSheet(0)
	importIssueTemplate.Fd = newFd

	err := importIssueTemplate.SetProjectAndTableInfo(orgId, projectId, tableId)
	if err != nil {
		return nil, err
	}

	otherLanguageMap := make(map[string]string, 0)
	relatePath := "/org_" + strconv.FormatInt(orgId, 10) + "/project_" + strconv.FormatInt(projectId, 10)
	excelDir := config.GetOSSConfig().RootPath + relatePath
	mkdirErr := os.MkdirAll(excelDir, 0777)
	if mkdirErr != nil {
		log.Error(mkdirErr)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, mkdirErr)
	}
	keyWord := "导入模板"
	if importIssueTemplate.IsEnglish {
		keyWord = "Import_template"
	}
	fileName := keyWord + ".xlsx"
	if importIssueTemplate.ProjectInfo.ProjectTypeId != consts.ProjectTypeNormalId && tableId != 0 {
		tableInfo, err := GetTableInfo(orgId, importIssueTemplate.ProjectInfo.AppId, tableId)
		if err != nil {
			log.Errorf("[ExportIssueTemplate] err: %v", err)
			return nil, err
		}
		if importIssueTemplate.IsEnglish {
			if tmpMap, ok1 := consts.LANG_ISSUE_OBJECT_TYPE_MAP[importIssueTemplate.Lang]; ok1 {
				otherLanguageMap = tmpMap
			}
			if tmpVal, ok2 := otherLanguageMap[tableInfo.Name]; ok2 {
				tableInfo.Name = tmpVal
			}
		}
		fileName = keyWord + "_" + tableInfo.Name + ".xlsx"
	}
	importIssueTemplate.FileNameWithPath = excelDir + "/" + fileName
	importIssueTemplate.FileUrl = config.GetOSSConfig().LocalDomain + relatePath + "/" + fileName

	return importIssueTemplate, nil
}

// SetNoticeText 在 A1 设置”注意事项“
func (nef *ImportIssueTemplate) SetNoticeText(projectTypeId int64, ignoreColumnNames []string) errs.SystemErrorInfo {
	nef.AddRow()

	isEnglish := nef.IsEnglish
	axis := nef.GetAxis("A", 1)
	cellA1Text := GetImportNoticeTextA1(projectTypeId)
	if isEnglish {
		cellA1Text = GetImportNoticeTextA1LangEn(projectTypeId)
	}
	if err := nef.Fd.SetCellRichText(nef.ActiveSheet, axis, cellA1Text); err != nil {
		businessErr := errs.BuildSystemErrorInfo(errs.FileSaveFail, err)
		log.Errorf("[SetNoticeText] err: %v", businessErr)
		return businessErr
	}
	// 设置样式
	style, err := nef.Fd.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "top",
			WrapText:   true,
		},
	})
	if err != nil {
		businessErr := errs.BuildSystemErrorInfo(errs.FileSaveFail, err)
		log.Errorf("[SetNoticeText] err: %v", businessErr)
		return businessErr
	}
	nef.Fd.SetCellStyle(nef.ActiveSheet, axis, axis, style)
	// 设置 A1 的宽高
	if err = nef.Fd.SetColWidth(nef.ActiveSheet, "A", "A", 83); err != nil {
		businessErr := errs.BuildSystemErrorInfo(errs.FileSaveFail, err)
		log.Errorf("[SetNoticeText] err: %v", businessErr)
		return businessErr
	}
	if err = nef.Fd.SetRowHeight(nef.ActiveSheet, 1, 190); err != nil {
		businessErr := errs.BuildSystemErrorInfo(errs.FileSaveFail, err)
		log.Errorf("[SetNoticeText] err: %v", businessErr)
		return businessErr
	}

	return nef.SetNoticeTextB1(ignoreColumnNames)
}

// SetNoticeTextB1 在 B1 设置”注意事项“
func (nef *ImportIssueTemplate) SetNoticeTextB1(ignoreColumnNames []string) errs.SystemErrorInfo {
	if len(ignoreColumnNames) == 0 {
		return nil
	}

	sort.Strings(ignoreColumnNames)
	// B1
	axis := nef.GetAxis("B", 1)
	cellB1Text := GetImportNoticeTextB1(ignoreColumnNames)
	if err := nef.Fd.SetCellRichText(nef.ActiveSheet, axis, cellB1Text); err != nil {
		businessErr := errs.BuildSystemErrorInfo(errs.FileSaveFail, err)
		log.Errorf("[SetNoticeTextB1] err: %v", businessErr)
		return businessErr
	}
	// 设置样式
	style, err := nef.Fd.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "top",
			WrapText:   true,
		},
	})
	if err != nil {
		businessErr := errs.BuildSystemErrorInfo(errs.FileSaveFail, err)
		log.Errorf("[SetNoticeText] err: %v", businessErr)
		return businessErr
	}
	nef.Fd.SetCellStyle(nef.ActiveSheet, axis, axis, style)
	// 设置 B 的宽
	if err = nef.Fd.SetColWidth(nef.ActiveSheet, "B", "B", 22); err != nil {
		businessErr := errs.BuildSystemErrorInfo(errs.FileSaveFail, err)
		log.Errorf("[SetNoticeText] err: %v", businessErr)
		return businessErr
	}
	return nil
}

// SetExcelHeadForImportTpl 任务导入模板，设置头部数据，列名称
func (nef *ImportIssueTemplate) SetExcelHeadForImportTpl() {
	nef.AddRow()
	for _, column := range nef.Columns {
		cell := nef.AddCell()
		cellValue := GetColumnNameWithAlias(column)
		cell.SetCell(cellValue)
		options := nef.GetColumnOptions(column)
		if len(options) > 0 {
			enumValList := make([]string, 0, len(options))
			for _, option := range options {
				enumValList = append(enumValList, option.Value)
			}
			colName, _ := excelize.ColumnNumberToName(nef.Col)
			nef.SetValidateDropList(colName, enumValList)
		} else if column.Name == consts.HelperFieldIsParentDesc {
			colName, _ := excelize.ColumnNumberToName(nef.Col)
			nef.SetValidateDropList(colName, []string{WordTranslate("父任务"), WordTranslate("子任务")})
		}
	}
}

// SetExcelHeadKeyAndLockSomeRow 向导入模板中设置列的 key，并隐藏和锁定特殊的行
func (nef *ImportIssueTemplate) SetExcelHeadKeyAndLockSomeRow() errs.SystemErrorInfo {
	nef.AddRow()
	for _, column := range nef.Columns {
		nef.AddCell().SetCell(column.Name)
	}

	// 数据行不锁定
	unLockStyle, _ := nef.Fd.NewStyle(&excelize.Style{
		Protection: &excelize.Protection{
			Locked: false,
		},
	})
	nef.Fd.SetRowStyle(nef.ActiveSheet, nef.DataStartRowTh, nef.DataEndRowTh, unLockStyle)

	// 影藏key行
	nef.Fd.SetRowVisible(nef.ActiveSheet, nef.Row, false)

	// 保护工作表
	nef.Fd.ProtectSheet(nef.ActiveSheet, &excelize.FormatSheetProtection{
		AlgorithmName: "SHA-512",
		Password:      "$^&FD&^F^",
		EditScenarios: false,
		DeleteColumns: true,
		InsertColumns: true,
	})

	return nil
}

func (nef *ImportIssueTemplate) SetValidateDropList(colFlag string, enumValList []string) {
	dvRange := excelize.NewDataValidation(true)
	startAxis, _ := excelize.JoinCellName(colFlag, nef.DataStartRowTh)
	endAxis, _ := excelize.JoinCellName(colFlag, nef.DataEndRowTh)
	dvRange.Sqref = fmt.Sprintf("%s:%s", startAxis, endAxis)
	dvRange.SetDropList(enumValList)
	nef.Fd.AddDataValidation(nef.ActiveSheet, dvRange)
}

func (nef *ImportIssueTemplate) GetAxis(col string, row int) string {
	axis, _ := excelize.JoinCellName(col, row)

	return axis
}

func (nef *ImportIssueTemplate) Save() errs.SystemErrorInfo {
	if err := nef.Fd.SaveAs(nef.FileNameWithPath); err != nil {
		businessErr := errs.BuildSystemErrorInfo(errs.FileSaveFail, err)
		log.Errorf("[SetExcelHeadForImportTpl] err: %v", businessErr)
		return businessErr
	}

	return nil
}

// GetColumnNameWithAlias 获取带有别名的字段名。如：标题（标题别名）
func GetColumnNameWithAlias(column *projectvo.TableColumnData) string {
	originName := column.Label
	if otherMap, ok := consts.LANG_ISSUE_IMPORT_TITLE_MAP[lang2.GetLang()]; ok {
		if tempValue, ok1 := otherMap[originName]; ok1 {
			originName = tempValue
		}
	}

	if column.AliasTitle != "" {
		return fmt.Sprintf("%s（%s）", originName, column.AliasTitle)
	} else {
		return originName
	}
}

func WordTranslate(word string) string {
	dic := consts.LANG_IMPORT_ERR_MAP["en"]
	if lang2.IsEnglish() {
		if name, ok := dic[word]; ok {
			return name
		}
	}

	return word
}

func ParseExcelDatetime(cellVal string) (time.Time, error) {
	t, err := dateparse.ParseAny(cast.ToString(cellVal))
	if err != nil {
		defaultRes := time.Time{}
		f, err := strconv.ParseFloat(cellVal, 64)
		if err != nil {
			log.Errorf("[ParseExcelDatetime] err: %v", err)
			return defaultRes, err
		}
		dt, err := excelize.ExcelDateToTime(f, false)
		if err != nil {
			log.Errorf("[ParseExcelDatetime] err: %v", err)
			return defaultRes, nil
		}

		return dt, nil
	}

	return t, nil
}

// GetUserDeptIdsMap 获取用户的部门id列表，一个用户可能在多个部门中
func GetUserDeptIdsMap(orgId int64, userIds []int64) (map[int64][]int64, errs.SystemErrorInfo) {
	resp := userfacade.GetAllUserByIds(orgId, userIds)
	if resp.Failure() {
		log.Errorf("[GetUserDeptIdsMap] err: %v", resp.Error())
		return nil, resp.Error()
	}
	resultMap := make(map[int64][]int64, 0)
	for _, user := range resp.Data {
		deptIdArr := make([]int64, 0)
		for _, dept := range user.DeptList {
			deptIdArr = append(deptIdArr, dept.DepartmentId)
		}
		if _, ok := resultMap[user.Id]; ok {
			resultMap[user.Id] = append(resultMap[user.Id], deptIdArr...)
		} else {
			resultMap[user.Id] = deptIdArr
		}
	}

	return resultMap, nil
}

// GetUserDeptFullNameArr 获取用户的完整部门名称列表，一个用户可能在多个部门中。
func GetUserDeptFullNameArr(orgId int64, userIds []int64) (map[int64][]string, errs.SystemErrorInfo) {
	// 用户的部门列表，如：{211: ["org1/dept1/dept2"]}
	resultMap := make(map[int64][]string, 0)
	userDeptIdMap, err := GetUserDeptIdsMap(orgId, userIds)
	if err != nil {
		log.Error(err)
		return resultMap, err
	}
	deptIds := make([]int64, 0)
	for _, item := range userDeptIdMap {
		deptIds = append(deptIds, item...)
	}
	deptIds = slice.SliceUniqueInt64(deptIds)
	if len(deptIds) < 1 {
		return resultMap, nil
	}
	// 获取部门的全路径部门名称信息，如：职能部/技术部/运维组
	fullDeptNameMap, err := GetDeptFullNameByIds(orgId, deptIds, "/")
	if err != nil {
		log.Errorf("[GetUserDeptFullNameArr] err: %v", err)
		return resultMap, nil
	}
	for userId, deptIdArr := range userDeptIdMap {
		if _, ok := resultMap[userId]; !ok {
			resultMap[userId] = make([]string, 0)
		}
		for _, deptId := range deptIdArr {
			if fullDeptNameMap[deptId] != "" {
				resultMap[userId] = append(resultMap[userId], fullDeptNameMap[deptId])
			}
		}
	}

	return resultMap, nil
}

// GetUserDeptFullNameStr 获取用户的部门名称，字符串拼接，如：`"职能部/技术部；职能部/财务部"`
// sep 拼接的分隔符，如：`；`
func GetUserDeptFullNameStr(orgId int64, userIds []int64, sep string) (map[int64]string, errs.SystemErrorInfo) {
	resultMap := make(map[int64]string, 0)
	deptNameMap, err := GetUserDeptFullNameArr(orgId, userIds)
	if err != nil {
		log.Errorf("[GetUserDeptFullNameStr] err: %v", err)
		return resultMap, err
	}
	for userId, deptNameArr := range deptNameMap {
		resultMap[userId] = strings.Join(deptNameArr, sep)
	}
	return resultMap, nil
}

// GetDeptFullNameByIds 获取部门的全路径部门名称信息，如：职能部/技术部/运维组
// sep 拼接多部门名称时的分隔符，如："/"
func GetDeptFullNameByIds(orgId int64, deptIds []int64, sep string) (map[int64]string, errs.SystemErrorInfo) {
	resultMap := make(map[int64]string, 0)
	deptNameMapResp := userfacade.GetFullDeptNameArrByIds(orgId, deptIds)
	if deptNameMapResp.Failure() {
		log.Errorf("[GetDeptFullNameByIds] err: %v", deptNameMapResp.Error())
		return resultMap, deptNameMapResp.Error()
	}
	for deptId, deptNameArr := range deptNameMapResp.Data {
		resultMap[deptId] = strings.Join(deptNameArr, sep)
	}

	return resultMap, nil
}

// GetUserOrDeptNameMap 获取用户所在部门或部门数据的全路径
// ids 可能是用户id，也可能是部门id
func GetUserOrDeptNameMap(orgId int64, dataList []bo.UserOrDeptInfoItem, idType string) (map[int64]string, errs.SystemErrorInfo) {
	var businessErr errs.SystemErrorInfo

	nameMap := make(map[int64]string, 0)
	switch idType {
	case "dept":
		deptIds := make([]int64, len(dataList))
		for i, item := range dataList {
			deptIds[i] = item.Id
		}
		nameMap, businessErr = GetDeptFullNameByIds(orgId, deptIds, "/")
		if businessErr != nil {
			log.Errorf("[GetUserOrDeptNameMap] err: %v", businessErr)
			return nameMap, businessErr
		}
	case "user":
		// 获取用户的部门
		userIds := make([]int64, len(dataList))
		for i, item := range dataList {
			userIds[i] = item.Id
		}
		// 获取完整的部门名称列表
		nameMap, businessErr = GetUserDeptFullNameStr(orgId, userIds, "；")
		if businessErr != nil {
			log.Errorf("[GetUserOrDeptNameMap] err: %v", businessErr)
			return nameMap, businessErr
		}
	}

	return nameMap, nil
}

// GetSameNameDeptList 获取有同名情况的部门列表
func GetSameNameDeptList(orgId int64) ([]uservo.GetDeptListRespDataItem, errs.SystemErrorInfo) {
	multiDeptIds := make([]int64, 0)
	multiDeptList := make([]uservo.GetDeptListRespDataItem, 0)
	deptResp := userfacade.GetDeptList(orgId)
	if deptResp.Failure() {
		log.Errorf("[GetSameNameDeptList] err: %v", deptResp.Error())
		return multiDeptList, deptResp.Error()
	}
	// 找出同名的组织
	multiMap := make(map[string][]uservo.GetDeptListRespDataItem, 0)
	multiDeptMapById := make(map[int64]uservo.GetDeptListRespDataItem, 0)
	for _, dept := range deptResp.Data {
		multiDeptMapById[dept.Id] = dept

		if _, ok := multiMap[dept.Name]; ok {
			multiMap[dept.Name] = append(multiMap[dept.Name], dept)
		} else {
			multiMap[dept.Name] = []uservo.GetDeptListRespDataItem{
				dept,
			}
		}
	}
	for _, deptArr := range multiMap {
		if len(deptArr) > 1 {
			for _, dept := range deptArr {
				multiDeptIds = append(multiDeptIds, dept.Id)
			}
		}
	}
	for _, deptId := range multiDeptIds {
		if dept, ok := multiDeptMapById[deptId]; ok {
			multiDeptList = append(multiDeptList, dept)
		}
	}

	return multiDeptList, nil
}

// GetSameNameDeptNames 获取同名部门的名称列表
func GetSameNameDeptNames(orgId int64) ([]string, errs.SystemErrorInfo) {
	deptNameArr := make([]string, 0)
	deptArr, err := GetSameNameDeptList(orgId)
	if err != nil {
		log.Errorf("[GetSameNameDeptNames] GetSameNameDeptList err: %v", err)
		return deptNameArr, err
	}
	for _, dept := range deptArr {
		deptNameArr = append(deptNameArr, dept.Name)
	}

	return deptNameArr, nil
}

// GetSameNameUserList 通过用户列表去除同名的用户
func GetSameNameUserList(userList []bo.SimpleUserInfoBo) ([]bo.UserOrDeptInfoItem, errs.SystemErrorInfo) {
	multiIds := make([]int64, 0)
	multiUserList := make([]bo.UserOrDeptInfoItem, 0)
	multiMap := make(map[string][]bo.SimpleUserInfoBo, 0)
	multiUserMapById := make(map[int64]bo.SimpleUserInfoBo, 0)

	for _, item := range userList {
		multiUserMapById[item.Id] = item

		if _, ok := multiMap[item.Name]; ok {
			multiMap[item.Name] = append(multiMap[item.Name], item)
		} else {
			multiMap[item.Name] = []bo.SimpleUserInfoBo{
				item,
			}
		}
	}
	for _, userArr := range multiMap {
		if len(userArr) > 1 {
			for _, user := range userArr {
				multiIds = append(multiIds, user.Id)
			}
		}
	}
	for _, id := range multiIds {
		if user, ok := multiUserMapById[id]; ok {
			multiUserList = append(multiUserList, bo.UserOrDeptInfoItem{
				Id:   user.Id,
				Name: user.Name,
			})
		}
	}

	return multiUserList, nil
}

func GetSameNameUserIdList(userList []bo.SimpleUserInfoBo) ([]int64, errs.SystemErrorInfo) {
	sameNameUsers, err := GetSameNameUserList(userList)
	if err != nil {
		log.Errorf("[GetSameNameUserIdList] GetSameNameUserList err: %v", err)
		return nil, err
	}
	ids := make([]int64, 0)
	for _, user := range sameNameUsers {
		ids = append(ids, user.Id)
	}

	return ids, nil
}

// GetSameNameUserAndDeptSimpleInfo 获取同名的用户和部门
// dataType 获取的数据类型, `user` 同名用户; `dept` 同名部门; all 所有
func GetSameNameUserAndDeptSimpleInfo(orgId int64, dataType string) ([]bo.UserOrDeptInfoItem, []bo.UserOrDeptInfoItem,
	errs.SystemErrorInfo,
) {
	var err errs.SystemErrorInfo
	deptTransData, userTransData := make([]bo.UserOrDeptInfoItem, 0), make([]bo.UserOrDeptInfoItem, 0)
	// 处理导出部门数据，只导出名字重复的部分
	if dataType == "all" || dataType == "dept" {
		sameNameDeptList, err := GetSameNameDeptList(orgId)
		if err != nil {
			log.Errorf("[GetSameNameUserAndDeptSimpleInfo] err: %v", err)
			return userTransData, deptTransData, err
		}
		copyer.Copy(sameNameDeptList, &deptTransData)
	}

	// 用户数据，只导出名字重复的部分
	if dataType == "all" || dataType == "user" {
		userInfoResp := orgfacade.GetUserInfoListByOrg(orgvo.GetUserInfoListReqVo{
			OrgId: orgId,
		})
		if userInfoResp.Failure() {
			log.Errorf("[GetSameNameUserAndDeptSimpleInfo] err: %v", userInfoResp.Error())
			return userTransData, deptTransData, userInfoResp.Error()
		}
		if len(userInfoResp.SimpleUserInfo) > 0 {
			dataList := userInfoResp.SimpleUserInfo
			userTransData, err = GetSameNameUserList(dataList)
			if err != nil {
				log.Errorf("[GetSameNameUserAndDeptSimpleInfo] err: %v", err)
				return userTransData, deptTransData, err
			}
		}
	}

	return userTransData, deptTransData, nil
}

// RenderSameNameKey 获取用户/部门名称及其id组成的值。导入时，对于同名用户，推荐用户以 `张三#123#` 的样式作为成员名导入，部门同理。
func RenderSameNameKey(userId int64, userName string) string {
	return fmt.Sprintf("%s#%d#", userName, userId)
}

// WriteUserOrDeptSameNameListToExcel 写入 excel 数据
// typeStr 数据类型：部门（dept），成员（user）
func WriteUserOrDeptSameNameListToExcel(orgId int64, fullFileName string, userList, deptList []bo.UserOrDeptInfoItem) errs.SystemErrorInfo {
	// 列：重名部门/成员 | 类型 | 所处部门 | 部门/成员ID | 建议导入样式
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	var businessErr errs.SystemErrorInfo

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		log.Errorf("[WriteUserOrDeptSameNameListToExcel] err: %v", err)
		return errs.BuildSystemErrorInfo(errs.InviteImportTplGenExcelErr, err)
	}
	row = sheet.AddRow()

	cell = row.AddCell()
	cell.Value = "重名部门/成员"
	cell = row.AddCell()
	cell.Value = "类型"
	cell = row.AddCell()
	cell.Value = "所处部门"
	cell = row.AddCell()
	cell.Value = "部门/成员ID"
	cell = row.AddCell()
	cell.Value = "建议导入样式"

	isEnglish := lang2.IsEnglish()
	langTextMap := make(map[string]string, 0)
	if isEnglish {
		if otherMap, ok := consts.LANG_ISSUE_IMPORT_TITLE_MAP[lang2.GetLang()]; ok {
			langTextMap = otherMap
			for i, c := range row.Cells {
				if tempValue, ok1 := otherMap[c.Value]; ok1 {
					row.Cells[i].Value = tempValue
				}
			}
		}
	}

	deptNameMap := make(map[int64]string, 0)
	if len(deptList) > 0 {
		deptNameMap, businessErr = GetUserOrDeptNameMap(orgId, deptList, "dept")
		if businessErr != nil {
			log.Errorf("[WriteUserOrDeptSameNameListToExcel] err: %v", businessErr)
			return businessErr
		}
		for _, item := range deptList {
			row = sheet.AddRow()

			cell = row.AddCell() // 重名部门/成员
			cell.Value = item.Name
			cell = row.AddCell() // 类型
			cell.Value = GetTypeNameForExportSameUserDept("dept", isEnglish, langTextMap)
			cell = row.AddCell() // 所处部门
			cell.Value = deptNameMap[item.Id]
			cell = row.AddCell() // 部门/成员ID
			cell.Value = strconv.FormatInt(item.Id, 10)
			cell = row.AddCell() // 建议导入样式
			cell.Value = RenderSameNameKey(item.Id, item.Name)
		}
	}
	if len(userList) > 0 {
		deptNameMap, businessErr = GetUserOrDeptNameMap(orgId, userList, "user")
		if businessErr != nil {
			log.Errorf("[WriteUserOrDeptSameNameListToExcel] err: %v", businessErr)
			return businessErr
		}
		// 数据列表写入
		for _, item := range userList {
			row = sheet.AddRow()

			cell = row.AddCell() // 重名部门/成员
			cell.Value = item.Name
			cell = row.AddCell() // 类型
			cell.Value = GetTypeNameForExportSameUserDept("user", isEnglish, langTextMap)
			cell = row.AddCell() // 所处部门
			cell.Value = deptNameMap[item.Id]
			cell = row.AddCell() // 部门/成员ID
			cell.Value = strconv.FormatInt(item.Id, 10)
			cell = row.AddCell() // 建议导入样式
			cell.Value = RenderSameNameKey(item.Id, item.Name)
		}
	}

	err = file.Save(fullFileName)
	if err != nil {
		log.Errorf("[WriteUserOrDeptSameNameListToExcel] err: %v", err)
		return errs.BuildSystemErrorInfo(errs.InviteImportTplGenExcelErr, err)
	}

	return nil
}

// GetTypeNameForExportSameUserDept 导出同名用户或部门时，渲染“类型” cell 中的值
func GetTypeNameForExportSameUserDept(typeFlag string, isEnglish bool, langTextMap map[string]string) string {
	name := ""
	switch typeFlag {
	case "user":
		name = "用户"
	case "dept":
		name = "部门"
	}
	if isEnglish {
		if tempValue, ok1 := langTextMap[name]; ok1 {
			name = tempValue
		}
	}

	return name
}

// GetImportNoticeTextB1 任务导入时的注意事项文案
func GetImportNoticeTextB1(columnNames []string) []excelize.RichTextRun {
	richText := []excelize.RichTextRun{}
	richText = append(richText, excelize.RichTextRun{
		Text: "以下字段类型暂不支持导入",
		Font: &excelize.Font{
			Color:  "3368FF",
			Family: "Times New Roman",
			Size:   11,
			Bold:   true,
		},
	})
	richText = append(richText, excelize.RichTextRun{
		Text: "\n已被隐藏：",
		Font: &excelize.Font{
			Color:  "3368FF",
			Family: "Times New Roman",
			Size:   11,
			Bold:   true,
		},
	})

	for i, d := range columnNames {
		richText = append(richText, excelize.RichTextRun{
			Text: fmt.Sprintf("\n%d. %s", i+1, d),
			Font: &excelize.Font{
				Family: "Times New Roman",
				Size:   11,
			},
		})
	}
	return richText
}

// GetImportNoticeText 任务导入时的注意事项文案
func GetImportNoticeTextA1(proTypeId int64) []excelize.RichTextRun {
	richText := []excelize.RichTextRun{}
	richText = append(richText, excelize.RichTextRun{
		Text: "导入必读：",
		Font: &excelize.Font{
			Color:  "3368FF",
			Family: "Times New Roman",
			Bold:   true,
			Size:   11,
		},
	})
	desc := []string{
		"\n1. 下载模板后，仅可整理行数据，极星协作和Excel内均不可再进行字段（列名）的调整，包括内容、顺序的调整，此模板表头已被锁定不可编辑",
		"\n2.【任务标题】为必填项，且不可超过500字",
		"\n3. 父子类型若为子任务，当前子任务会跟随上一个父级任务",
		"\n4. 任务总数不超过300条",
		"\n5. 选择项目前只支持添加项目内目标字段配置的选项，多个选项请用“||”符号隔开，如“PC端||移动端||小程序”",
	}
	if int(proTypeId) != consts.ProjectTypeEmpty {
		desc = append(desc,
			"\n6. 日期字段如“开始时间”字段：填写格式为YYYY-MM-DD XX:XX，如“2020/04/10 12:00”",
			"\n7. 截止时间不能小于开始时间",
			"\n8. 负责人若不填写或组织内不存在该人员姓名，那么本条任务的负责人将默认为未分配")
		if int(proTypeId) == consts.ProjectTypeNormalId {
			desc = append(desc, "\n9.需要确认人执行确认的任务，请添加确认人，并在任务状态栏输入“待确认”")
		}
	}

	for _, d := range desc {
		richText = append(richText, excelize.RichTextRun{
			Text: d,
			Font: &excelize.Font{
				Family: "Times New Roman",
				Size:   11,
			},
		})
	}
	return richText
}

// GetImportNoticeTextA1LangEn 任务导入时的注意事项文案英文版
func GetImportNoticeTextA1LangEn(proTypeId int64) []excelize.RichTextRun {
	richText := []excelize.RichTextRun{}
	richText = append(richText, excelize.RichTextRun{
		Text: "Important Tips：",
		Font: &excelize.Font{
			Color:  "339FFF",
			Family: "Times New Roman",
		},
	})
	desc := []string{
		"\n1. After downloading the template, only row data can be organized. No field (column name) adjustment can be made in Polestar Collaboration and Excel, including content and order adjustment",
		"\n2. title is required",
		"\n3. If the parent-child type is a subtask, the current subtask will follow the previous parent task",
		"\n4. The total number of tasks should not exceed 300",
		"\n5. The task title should not exceed 50 words",
		"\n6. Only the option to add the target field configuration within the project is supported before selecting the project. Multiple options: please use \"||\" symbol to separate them, such as \"PC terminal||mobile terminal||applet\"",
		"\n7. Date field, such as \"start time\": fill in the format of yyyy-mm-dd XX: XX, such as \"2020/04/10 12:00\"",
		"\n8. The deadline cannot be less than the start time",
		"\n9. If the assignee does not fill in or does not exist in the organization, the assignee of this task will be the \"unallocated\" by default",
	}
	if int(proTypeId) == consts.ProjectTypeNormalId {
		desc = append(desc, "\n10. If you need a task to be confirmed by someone,  add the person as confirmer, and enter \"To confirm\" in the task status bar")
	}

	for _, d := range desc {
		richText = append(richText, excelize.RichTextRun{
			Text: d,
			Font: &excelize.Font{
				Family: "Times New Roman",
			},
		})
	}
	return richText
}

// GetUserMapForImportIssue 导入任务时，获取用于匹配成员名的用户 map 信息。map 中还设定了针对重名用户特殊的 key 信息映射。
func GetUserMapForImportIssue(orgId int64) (map[string]bo.SimpleUserWithUserIdBo, map[string]int64, []string, errs.SystemErrorInfo) {
	userNameMapUid := map[string]int64{}
	sameNames := make([]string, 0)
	userInfoMap := make(map[string]bo.SimpleUserWithUserIdBo, 0)
	userInfo := orgfacade.GetUserInfoListByOrg(orgvo.GetUserInfoListReqVo{
		OrgId: orgId,
	})
	if userInfo.Failure() {
		log.Errorf("[GetUserMapForImportIssue] orgfacade.GetUserInfoListByOrg err: %v", userInfo.Error())
		return userInfoMap, userNameMapUid, sameNames, userInfo.Error()
	}
	// 针对同名部分做处理。先获取同名的用户列表，将其按特定格式追加到 userInfoMap 中
	sameNameUserIds, err := GetSameNameUserIdList(userInfo.SimpleUserInfo)
	if err != nil {
		log.Errorf("[GetUserMapForImportIssue] domain.GetSameNameUserIdList err: %v", err)
		return userInfoMap, userNameMapUid, sameNames, err
	}
	for _, v := range userInfo.SimpleUserInfo {
		userNameMapUid[v.Name] = v.Id
		tmpUserBo := bo.SimpleUserWithUserIdBo{
			Id:     v.Id,
			UserId: v.Id,
			Name:   v.Name,
			Avatar: v.Avatar,
			Status: v.Status,
			Type:   consts.LcCustomFieldUserType,
		}
		userInfoMap[v.Name] = tmpUserBo
		if exist, _ := slice.Contain(sameNameUserIds, v.Id); exist {
			uniqueNameKey := RenderSameNameKey(v.Id, v.Name)
			userNameMapUid[uniqueNameKey] = v.Id
			userInfoMap[uniqueNameKey] = tmpUserBo
			sameNames = append(sameNames, v.Name)
		}
	}

	return userInfoMap, userNameMapUid, sameNames, nil
}

// CheckNeedCollaborateFilter 检查是否需要协作人条件进行查询。不是管理员，则需要协作人条件
func CheckNeedCollaborateFilter(orgId, currentUserId, tableId int64, projectBo *bo.ProjectBo) bool {
	optAuthResp := permissionfacade.GetAppAuth(orgId, projectBo.AppId, tableId, currentUserId)
	if optAuthResp.Failure() {
		log.Errorf("[LcHomeIssuesForProject] orgId:%d, userId:%d, appId:%v, GetAppAuthWithoutCollaborator failure:%v", orgId, currentUserId, projectBo.AppId, optAuthResp.Error())
		return false
	}
	appAuthInfo := &optAuthResp.Data
	isAdmin := appAuthInfo.HasAppRootPermission || appAuthInfo.SysAdmin || appAuthInfo.OrgOwner || appAuthInfo.
		AppOwner || currentUserId == projectBo.Owner

	return !isAdmin
}
