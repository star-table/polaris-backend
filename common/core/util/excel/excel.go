package excel

//
//import (
//	"bytes"
//	"errors"
//	"fmt"
//	"net/http"
//	"strings"
//
//	"github.com/star-table/common/core/logger"
//	"github.com/star-table/common/core/util/json"
//	"github.com/star-table/common/core/util/slice"
//	"github.com/star-table/polaris-backend/common/core/consts"
//	"github.com/star-table/polaris-backend/common/core/errs"
//	"github.com/star-table/polaris-backend/common/core/util/excel/sortedmap"
//	str2 "github.com/star-table/polaris-backend/common/core/util/str"
//	"github.com/star-table/polaris-backend/common/model/bo"
//	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
//	"github.com/tealeg/xlsx/v3"
//)
//
//var log = *logger.GetDefaultLogger()
//
///**
//* GenerateCSVFromXLSXFileWithMap 通过 excel 组装数据
//* excelFileName 本地excel路径名
//* orginFileUrl  远程excel url
//* sheetIndex    第几份excel
//* ignoreRows    忽略行数
//* timeFormatIndex 需要转换时间格式的列
//* parseType 解析类型 array: 以数组的方式解析，map: 以 map 的方式解析
//* headerMapGroup 支持的头部组
//* extraParamMap 额外参数
//	* isCommonPro 是否是通用项目
//*/
//func GenerateCSVFromXLSXFileWithMap(url string, urlType int, sheetIndex int, ignoreRows int, timeFormatIndex []string, parseType string,
//	headerMapGroup map[string]*sortedmap.SortMap, extraParamMap map[string]interface{},
//) (interface{}, int, *sortedmap.SortMap, error) {
//	var error error
//	var xlFile *xlsx.File
//	if urlType == consts.UrlTypeHttpPath {
//		resp, err := http.Get(url)
//		if err != nil {
//			log.Error(err)
//			return nil, ignoreRows, nil, err
//		}
//		resBody := resp.Body
//		buf := new(bytes.Buffer)
//		buf.ReadFrom(resBody)
//
//		xlFile, error = xlsx.OpenBinary(buf.Bytes())
//	} else {
//		// 多2行是注意事项和表头
//		// url = "../../../app/" + url
//		xlFile, error = xlsx.OpenFile(url)
//	}
//	if error != nil {
//		log.Errorf("[GenerateCSVFromXLSXFileWithMap] err: %v", error)
//		return nil, ignoreRows, nil, error
//	}
//	sheetLen := len(xlFile.Sheets)
//	switch {
//	case sheetLen == 0:
//		return nil, ignoreRows, nil, errors.New("This XLSX file contains no sheets.")
//	case sheetIndex >= sheetLen:
//		return nil, ignoreRows, nil, fmt.Errorf("No sheet %d available, please select a sheet between 0 and %d\n", sheetIndex, sheetLen-1)
//	}
//	sheet := xlFile.Sheets[sheetIndex]
//	var result [][]string
//	var headerMap *sortedmap.SortMap
//	switch parseType {
//	case "map":
//		headerMap, oriErr := GetAndCheckHeaderMap(headerMapGroup, sheet, extraParamMap)
//		if oriErr != nil {
//			log.Errorf("[GenerateCSVFromXLSXFileWithMap] err: %v", oriErr)
//			return result, ignoreRows, nil, oriErr
//		}
//		mapArrResult, ignoreRows := assemblyResultWithKey(sheet, headerMap, timeFormatIndex)
//		return mapArrResult, ignoreRows, headerMap, nil
//	default:
//		// result, ignoreRows = assemblyResult(sheet, ignoreRows, timeFormatIndex)
//		return nil, ignoreRows, headerMap, errors.New("不支持的 parseType 类型。")
//	}
//
//	return result, ignoreRows, headerMap, nil
//}
//
//// GetAndCheckHeaderMap 通过导入的 excel 的行获取与之匹配的头部，并检查 excel 中的表头是否和 table 的表头一致。
//// headerMapGroup 支持的头部组
//// isCommonPro 是否是通用项目
//func GetAndCheckHeaderMap(headerMapGroup map[string]*sortedmap.SortMap, sheet *xlsx.Sheet, extraParamMap map[string]interface{},
//) (*sortedmap.SortMap, error) {
//	var headerMap *sortedmap.SortMap
//	projectTypeId := extraParamMap["projectTypeId"].(int64)
//	switch int(projectTypeId) {
//	case consts.ProjectTypeNormalId:
//		// 如果是通用项目，则直接返回确定的头部。
//		if val, ok := headerMapGroup[consts.ProjectObjectTypeLangCodeCommonTask]; ok {
//			headerMap = val
//		}
//	case consts.ProjectTypeAgileId, consts.ProjectTypeCommon2022V47:
//		tableInfo := extraParamMap["tableInfo"].(*bo.ProjectObjectTypeBo)
//		tableName := tableInfo.Name
//		if val, ok := headerMapGroup[tableName]; ok {
//			headerMap = val
//		}
//	}
//	// 任务的基础列/内置列
//	builtinHeader := sortedmap.NewSortMap()
//	for _, columnKey := range headerMap.KeyVec {
//		switch int(projectTypeId) {
//		case consts.ProjectTypeCommon2022V47:
//			// 特殊处理：新项目类型，任务栏不会再出现
//			// Common2022 指代新项目类型
//			if exist, _ := slice.Contain(consts.HiddenColumnNamesForCommon2022Type, columnKey); exist {
//				continue
//			}
//		}
//		builtinHeader.Insert(columnKey, headerMap.GetByKey(columnKey, ""))
//	}
//
//	// 追加自定义的列
//	columns := make([]*projectvo.TableColumnData, 0)
//	if tmpVal, ok := extraParamMap["customColumns"]; ok {
//		if tmpColumns, ok2 := tmpVal.([]*projectvo.TableColumnData); ok2 {
//			columns = tmpColumns
//		} else {
//			oriErr := errors.New("columns assert error.")
//			return headerMap, oriErr
//		}
//	}
//	for _, column := range columns {
//		headerMap.Insert(column.Name, column.Name)
//	}
//
//	// 比较 excel header 是否匹配当前表的列
//	if err := checkExcelHeader(builtinHeader, sheet, extraParamMap); err != nil {
//		log.Errorf("[GetAndCheckHeaderMap] err: %v", err)
//		return headerMap, err
//	}
//
//	return headerMap, nil
//}
//
//// checkExcelHeader 检查 excel 中数据的表头是否和当前表的表头一致
//func checkExcelHeader(builtinHeader *sortedmap.SortMap, sheet *xlsx.Sheet, extraParamMap map[string]interface{}) error {
//	ignoreRows := 1
//	ignoreRows = GetIgnoreRowNum(sheet)
//	ignoreRowsForHeader := ignoreRows - 1
//	projectTypeId := extraParamMap["projectTypeId"].(int64)
//	lang := "cn"
//	if tmpVal, ok := extraParamMap["lang"]; ok {
//		lang = tmpVal.(string)
//	}
//	isOtherLang := false
//	if tmpVal, ok := extraParamMap["isOtherLang"]; ok {
//		isOtherLang = tmpVal.(bool)
//	}
//	otherLanguageMap := make(map[string]string, 0)
//	if isOtherLang {
//		if isOtherLang {
//			if tmpMap, ok1 := consts.LANG_ISSUE_IMPORT_TITLE_MAP[lang]; ok1 {
//				otherLanguageMap = tmpMap
//			}
//		}
//	}
//	excelHeaderNameArr := make([]string, 0)
//	headerRow, err := sheet.Row(ignoreRowsForHeader)
//	if err != nil {
//		log.Errorf("[checkExcelHeader] err: %v", err)
//		return err
//	}
//	headerRow.ForEachCell(func(cell *xlsx.Cell) error {
//		str, _ := cell.FormattedValue()
//		if len(str) > 0 {
//			excelHeaderNameArr = append(excelHeaderNameArr, strings.TrimSpace(str))
//		}
//		return nil
//	})
//
//	customColumns := extraParamMap["customColumns"].([]*projectvo.TableColumnData)
//	tableHeaderNameArr := make([]string, 0)
//	// 内置字段是特定的顺序，必须从 builtinHeader 中取
//	for _, columnKey := range builtinHeader.KeyVec {
//		// tableHeaderNameArr 中存储的是列的名字
//		tableHeaderNameArr = append(tableHeaderNameArr, builtinHeader.GetByKey(columnKey, ""))
//	}
//	for _, column := range customColumns {
//		switch int(projectTypeId) {
//		case consts.ProjectTypeCommon2022V47:
//			// 特殊处理：新项目类型，任务栏不会再出现
//			// Common2022 指代新项目类型
//			// column.Name 即是 column key
//			if exist, _ := slice.Contain(consts.HiddenColumnNamesForCommon2022Type, column.Name); exist {
//				continue
//			}
//		}
//		tableHeaderNameArr = append(tableHeaderNameArr, column.Label)
//	}
//	if len(tableHeaderNameArr) != len(excelHeaderNameArr) {
//		err := errs.ImportExcelNotMatchedWithTable
//		viewHeaderJson := json.ToJsonIgnoreError(tableHeaderNameArr)
//		excelHeaderJson := json.ToJsonIgnoreError(excelHeaderNameArr)
//		log.Errorf("[checkExcelHeader] 表头列数量不一致 err: %v, tableHeaderNameArr:%s, excelHeaderNameArr:%s", err,
//			viewHeaderJson, excelHeaderJson)
//		return errs.BuildSystemErrorInfoWithMessage(err,
//			fmt.Sprintf("导入表格与极星模板的表头（列名）不一致。可能是您下载模板之后，调整过字段（列名）的内容或顺序。请重新下载导入模板，进行导入。"))
//	}
//	for i, oneName := range excelHeaderNameArr {
//		tmpTableName := strings.TrimSpace(tableHeaderNameArr[i])
//		if isOtherLang {
//			if tmpVal, ok2 := otherLanguageMap[tmpTableName]; ok2 {
//				tmpTableName = tmpVal
//			}
//		}
//		if !strings.Contains(oneName, tmpTableName) {
//			err := errs.ImportExcelNotMatchedWithTable
//			log.Errorf("[checkExcelHeader] err: %v, oneName:%v, tmpTableName:%v", err, oneName, tmpTableName)
//			return errs.BuildSystemErrorInfoWithMessage(err, fmt.Sprintf("表头列不同！请检查并修改或者重新下载导入模板。极星表列：[%s],"+
//				"excel 列：[%s]",
//				tmpTableName,
//				oneName))
//		}
//	}
//
//	return nil
//}
//
//// GetIgnoreRowNum 获取忽略的行数，主要存放的是注意事项，或表头
//func GetIgnoreRowNum(sheet *xlsx.Sheet) int {
//	ignoreRows := 1
//	firstRow, err := sheet.Row(0)
//	if err != nil {
//		log.Errorf("[getIgnoreRowNum] err: %v", err)
//	}
//	firstCellValArr := make([]string, 0)
//	firstRow.ForEachCell(func(cell *xlsx.Cell) error {
//		firstCellValArr = append(firstCellValArr, cell.Value)
//		return nil
//	})
//	if sheet.MaxRow > 0 && len(firstCellValArr) > 0 {
//		if len(firstCellValArr) == 1 ||
//			(len(firstCellValArr) > 3 && firstCellValArr[1] == "") {
//			ignoreRows = 2
//		} else {
//			ignoreRows = 1
//		}
//	}
//
//	return ignoreRows
//}
//
//func GetRowsIgnoreErr(sheet *xlsx.Sheet) [][]string {
//	rows := make([][]string, 0)
//	sheet.ForEachRow(func(row *xlsx.Row) error {
//		rowValArr := make([]string, 0)
//		row.ForEachCell(func(cell *xlsx.Cell) error {
//			rowValArr = append(rowValArr, cell.Value)
//			return nil
//		})
//		rows = append(rows, rowValArr)
//		return nil
//	})
//
//	return rows
//}
//
//// 以表头值作为 key 的方式处理
//func assemblyResultWithKey(sheet *xlsx.Sheet, headerMap *sortedmap.SortMap, timeFormatIndex []string) ([]*sortedmap.SortMap, int) {
//	ignoreRows := 1
//	mapArrResult := make([]*sortedmap.SortMap, 0)
//	ignoreRows = GetIgnoreRowNum(sheet)
//
//	err := sheet.ForEachRow(func(row *xlsx.Row) error {
//		key := row.GetCoordinate()
//		if key <= ignoreRows-1 {
//			return nil // continue
//		}
//		vals := make([]string, 0)
//		rowMap := sortedmap.NewSortMap()
//		err := row.ForEachCell(func(cell *xlsx.Cell) error {
//			// 通过 index 获取对应的 key
//			cellIndex, _ := cell.GetCoordinates()
//			tmpKey, _ := headerMap.GetKeyByIndex(cellIndex)
//
//			str, err := cell.FormattedValue()
//			// 去除 cell 中的空白字符
//			str = strings.Trim(str, "\n\r\t ")
//			str = strings.TrimSpace(str)
//			if err != nil {
//				vals = append(vals, err.Error())
//				rowMap.Insert(tmpKey, str)
//				return nil // continue
//			}
//			if str2.InArray(tmpKey, timeFormatIndex) || CheckKeyIsDateTimeV2(tmpKey) {
//				time, err := cell.GetTime(false)
//				if err != nil {
//					log.Infof("[assemblyResultWithKey] GetTime err: %v,cell: [%s]", err, cell.Value)
//				} else {
//					str = time.Format(consts.AppTimeFormat)
//				}
//			}
//			rowMap.Insert(tmpKey, str)
//			return nil
//		})
//		if err != nil {
//			log.Errorf("[assemblyResultWithKey] ForEachCell err: %v", err)
//		}
//		mapArrResult = append(mapArrResult, rowMap)
//		return nil
//	})
//	if err != nil {
//		log.Errorf("[assemblyResultWithKey] ForEachRow err: %v", err)
//	}
//
//	return mapArrResult, ignoreRows
//}
//
//// 检查**自定义字段**的键是否是时间日期类型。
//func CheckKeyIsDateTimeV2(key string) bool {
//	dateTimeFlag := fmt.Sprintf(":Type=%v", "datepicker")
//	return strings.Contains(key, dateTimeFlag)
//}
