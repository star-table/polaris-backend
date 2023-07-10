package domain

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/star-table/polaris-backend/common/model/vo/uservo"

	"upper.io/db.v3"

	lang2 "github.com/star-table/polaris-backend/common/core/util/lang"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/xuri/excelize/v2"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/md5"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	consts2 "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type InAndExportBase struct {
	OrgId            int64
	FileName         string
	FileNameWithPath string
	FileUrl          string
	Fd               *excelize.File
	ActiveSheet      string
	ProjectInfo      *bo.ProjectBo
	TableInfo        *projectvo.TableMetaData
	Columns          []*projectvo.TableColumnData
	ColumnsMap       map[string]*projectvo.TableColumnData
	Col              int // 当前所在的列
	Row              int // 当前所在的行
	RowMaxHeight     float64
}

func (ei *InAndExportBase) SetColumns(columns []*projectvo.TableColumnData) errs.SystemErrorInfo {
	ei.Columns = columns
	ei.ColumnsMap = make(map[string]*projectvo.TableColumnData, len(columns))

	for _, column := range columns {
		// 如果是迭代，则设置一下迭代选项
		if column.Name == consts.BasicFieldIterationId {
			options, err := ei.getIterationOptions()
			if err != nil {
				return err
			}
			column.Field.Props["select"] = map[string]interface{}{"options": options}
		}
		ei.ColumnsMap[column.Name] = column
	}

	return nil
}

func (ei *InAndExportBase) getIterationOptions() ([]*projectvo.ColumnSelectOption, errs.SystemErrorInfo) {
	// 迭代信息
	iterResp, _, err1 := GetIterationBoList(1, 5000, db.Cond{
		consts.TcOrgId:     ei.OrgId,
		consts.TcProjectId: ei.ProjectInfo.Id,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}, nil)
	if err1 != nil {
		log.Errorf("[NewExportIssue] GetIterationBoList :%v", err1)
		return nil, err1
	}

	options := make([]*projectvo.ColumnSelectOption, 0, len(*iterResp))
	for _, iterationBo := range *iterResp {
		options = append(options, &projectvo.ColumnSelectOption{
			Id:    iterationBo.Id,
			Value: iterationBo.Name,
		})
	}
	options = append(options, &projectvo.ColumnSelectOption{Id: 0, Value: "未规划"})

	return options, nil
}

func (ei *InAndExportBase) AddSheet(sheetName string) {
	ei.ActiveSheet = sheetName
}

func (ei *InAndExportBase) AddRow() *InAndExportBase {
	ei.Row += 1
	ei.RowMaxHeight = 0
	ei.Col = 0
	return ei
}

func (ei *InAndExportBase) AddCell() *InAndExportBase {
	ei.Col += 1
	return ei
}

func (ei *InAndExportBase) SetCell(value interface{}) {
	colName, _ := excelize.ColumnNumberToName(ei.Col)
	axis := fmt.Sprintf("%s%d", colName, ei.Row)
	ei.Fd.SetCellValue(ei.ActiveSheet, axis, value)
}

func (ei *InAndExportBase) AddImage(documents []*bo.DocumentBo) {
	padding := 5
	imageWidth := 48
	maxSizeOfRow := 3
	widthMMToPx := 7.6
	heightMMToPx := 1.2
	colName, _ := excelize.ColumnNumberToName(ei.Col)
	axis := fmt.Sprintf("%s%d", colName, ei.Row)
	count := len(documents)
	if count > 0 {
		sort.Slice(documents, func(i, j int) bool {
			return documents[i].Id < documents[j].Id
		})
		lastColName, _ := excelize.ColumnNumberToName(ei.Col - 1)
		lastWidth, _ := ei.Fd.GetColWidth(ei.ActiveSheet, lastColName)
		lastHeight, _ := ei.Fd.GetRowHeight(ei.ActiveSheet, ei.Row-1)
		maxWidth := float64(imageWidth*maxSizeOfRow+((count-1)%maxSizeOfRow)*padding+30) / widthMMToPx
		maxHeight := float64(imageWidth*((count-1)/maxSizeOfRow+1)+((count-1)/maxSizeOfRow)*padding+15) / heightMMToPx
		if maxHeight > ei.RowMaxHeight {
			ei.RowMaxHeight = maxHeight
		} else {
			maxHeight = ei.RowMaxHeight
		}

		ei.Fd.SetColWidth(ei.ActiveSheet, colName, colName, maxWidth)
		ei.Fd.SetRowHeight(ei.ActiveSheet, ei.Row, maxHeight)
		ei.Fd.SetColWidth(ei.ActiveSheet, lastColName, lastColName, maxWidth)
		ei.Fd.SetRowHeight(ei.ActiveSheet, ei.Row-1, maxHeight)

		for index, document := range documents {
			x := int(float64((index%maxSizeOfRow)*padding+(index%maxSizeOfRow)*imageWidth)) + 10
			y := (index/maxSizeOfRow)*(imageWidth+padding) + 5
			fileName := document.Name
			_, err := os.Stat(fileName)
			if err != nil {
				fileName = "./config/static/icon/" + document.Icon
			}
			formatStyle := fmt.Sprintf(`{
        		"x_offset": %v,
        		"y_offset": %v,
        		"hyperlink": "%s",
        		"hyperlink_type": "External",
        		"print_obj": false,
        		"lock_aspect_ratio": true,
        		"locked": false,
        		"positioning": "oneCell",
				"autofit": false
   			 }`, x, y, document.DisplayPath)
			ei.Fd.AddPicture(ei.ActiveSheet, axis, fileName, formatStyle)
		}

		ei.Fd.SetColWidth(ei.ActiveSheet, lastColName, lastColName, lastWidth)
		ei.Fd.SetRowHeight(ei.ActiveSheet, ei.Row-1, lastHeight)
	}
}

func (ei *InAndExportBase) SetProjectAndTableInfo(orgId, projectId, tableId int64) (err errs.SystemErrorInfo) {
	ei.ProjectInfo, err = GetProject(orgId, projectId)
	if err != nil {
		log.Error(err)
		return errs.ProjectNotExist
	}
	ei.TableInfo, err = GetTableInfo(orgId, ei.ProjectInfo.AppId, tableId)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (ei *InAndExportBase) GetColumnOptions(column *projectvo.TableColumnData) []*projectvo.ColumnSelectOption {
	var options []*projectvo.ColumnSelectOption
	props := column.Field.Props[column.Field.Type]
	if m, ok := props.(map[string]interface{}); ok && m["options"] != nil {
		options = make([]*projectvo.ColumnSelectOption, 0, 10)
		_ = json.FromJson(json.ToJsonIgnoreError(m["options"]), &options)
	}

	return options
}

// ExportIssue 任务导出
type ExportIssue struct {
	InAndExportBase
	ImageDownloadDir string
	Iteration        *bo.IterationBo
}

// NewExportIssue 初始化任务导出实例
func NewExportIssue(orgId int64, projectId int64, tableId, iterationId int64) (*ExportIssue, errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	ei := &ExportIssue{
		InAndExportBase: InAndExportBase{
			OrgId:       orgId,
			ActiveSheet: "Sheet1",
			Fd:          excelize.NewFile(),
		},
	}
	ei.Fd.NewSheet("Sheet1")

	err = ei.SetProjectAndTableInfo(orgId, projectId, tableId)
	if err != nil {
		return nil, err
	}

	if iterationId > 0 {
		ei.Iteration, err = GetIterationBoByOrgId(iterationId, orgId)
		if err != nil {
			log.Error(err)
			return ei, err
		}
	}

	return ei, ei.MakeDirForExport()
}

func (ei *ExportIssue) MakeDirForExport() errs.SystemErrorInfo {
	tempDownloadDir := "tmp_img"
	relatePath := "/org_" + strconv.FormatInt(ei.OrgId, 10) + "/project_" + strconv.FormatInt(ei.ProjectInfo.Id, 10)
	excelDir := config.GetOSSConfig().RootPath + relatePath
	mkdirErr := os.MkdirAll(excelDir+"/"+tempDownloadDir, 0777)
	if mkdirErr != nil {
		log.Error(mkdirErr)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, mkdirErr)
	}
	fileName := ei.ProjectInfo.Name + "_" + ei.TableInfo.Name + WordTranslate("统计") + ".xlsx"
	if ei.Iteration != nil {
		fileName = ei.ProjectInfo.Name + "_" + ei.Iteration.Name + "_" + ei.TableInfo.Name + WordTranslate("统计") + ".xlsx"
	}
	fileName = strings.Replace(fileName, "/", "_", -1)
	ei.FileNameWithPath = excelDir + "/" + fileName
	ei.FileUrl = config.GetOSSConfig().LocalDomain + relatePath + "/" + fileName
	ei.ImageDownloadDir = excelDir + "/" + tempDownloadDir + "/"

	return nil
}

func (ei *ExportIssue) TransferLcDataIntoExportData(exportObj *ExportIssue, lcData map[string]interface{}, userDepts map[string]*uservo.MemberDept) *bo.ExportIssueBo {
	exportIssueBo := &bo.ExportIssueBo{
		Id:           cast.ToInt64(lcData[consts.BasicFieldIssueId]),
		ParentId:     cast.ToInt64(lcData[consts.BasicFieldParentId]),
		IterationId:  cast.ToInt64(lcData[consts.BasicFieldIterationId]),
		ExportValues: make(map[string]interface{}, len(exportObj.Columns)),
	}

	for _, column := range exportObj.Columns {
		columnKey := column.Name
		propsMap := column.Field.Props
		val := lcData[columnKey]
		var exportValue interface{}

		switch column.Field.Type {
		case consts.LcColumnFieldTypeGroupSelect: // 分组单选
			exportValue = ei.getGroupSelectExportInfo(val, column, exportObj.ProjectInfo, lcData)
		case consts.LcColumnFieldTypeInputNumber: // 数字
			if val == nil || val == "null" {
				continue
			}
			exportValue = ei.getNumberExportInfo(val, propsMap)
		case consts.LcColumnFieldTypeAmount: // 货币 // https://wiki.bjx.cloud/pages/viewpage.action?pageId=32145423
			if val == nil || val == "null" {
				continue
			}
			exportValue = ei.getAmountExportInfo(val, propsMap)
		case consts.LcColumnFieldTypeSelect: // 单选。单选值可能是 int64 也可能是 string
			exportValue = ei.getSelectExportInfo(val, column)
		case consts.LcColumnFieldTypeMultiSelect: // 多选
			exportValue = ei.getMultiSelectExportInfo(val, column)
		case consts.LcColumnFieldTypeDept:
			exportValue = ei.getDeptExportInfo(val, userDepts)
		case consts.LcColumnFieldTypeMember:
			exportValue = ei.getMemberExportInfo(val, userDepts)
		case consts.LcColumnFieldTypeDatepicker:
			exportValue = ei.getDatepickerExportInfo(val, columnKey)
		case consts.LcColumnFieldTypeDocument, consts.LcColumnFieldTypeImage:
			tempDocuments := ei.getDocumentExportInfo(val)
			exportValue = tempDocuments
			exportIssueBo.Documents = append(exportIssueBo.Documents, tempDocuments...)
		case consts.LcColumnFieldTypeBaRelating:
			exportValue = ei.getBaRelatingExportInfo(val)
		case consts.LcColumnFieldTypeRelating:
			exportValue = ei.getRelatingExportInfo(val)
		case consts.LcColumnFieldTypeWorkHour:
			exportValue = ei.getWorkHourExportInfo(val)
		case consts.LcColumnFieldTypeRichText:
			exportValue = ei.getRichTextExportInfo(val)
		case consts.LcColumnFieldTypeConditionRef:
			tempValue, tempDocuments := ei.getConditionRefExportInfo(exportObj, val, propsMap, userDepts)
			exportIssueBo.Documents = append(exportIssueBo.Documents, tempDocuments...)
			if len(tempDocuments) > 0 {
				exportValue = tempDocuments
			} else {
				exportValue = tempValue
			}
		default:
			if columnKey == consts.HelperFieldIsParentDesc {
				exportValue = ei.getParentDescExportInfo(lcData)
			} else {
				exportValue = cast.ToString(val)
			}
		}
		exportIssueBo.ExportValues[columnKey] = exportValue
		//}
	}

	return exportIssueBo
}

func (ei *ExportIssue) getRichTextExportInfo(val interface{}) string {
	regx, err := regexp.Compile("<[^>]*>")
	if err != nil {
		return ""
	}
	valStr := regx.ReplaceAllString(cast.ToString(val), "")
	valStr = strings.Replace(valStr, "\n", "", -1)
	valStr = strings.Replace(valStr, "\t", "", -1)
	valStr = strings.Replace(valStr, "\r", "", -1)

	return valStr
}

func (ei *ExportIssue) getBaRelatingExportInfo(val interface{}) string {
	relating := &bo.RelatingIssue{}
	err := json.FromJson(json.ToJsonIgnoreError(val), relating)
	if err != nil {
		return ""
	}
	if len(relating.LinkTo) == 0 && len(relating.LinkFrom) == 0 {
		return ""
	}

	return fmt.Sprintf(consts2.BaRelatingExportTemplate, len(relating.LinkTo), len(relating.LinkFrom))
}

func (ei *ExportIssue) getRelatingExportInfo(val interface{}) string {
	relating := &bo.RelatingIssue{}
	err := json.FromJson(json.ToJsonIgnoreError(val), relating)
	if err != nil {
		return ""
	}

	total := len(relating.LinkTo) + len(relating.LinkFrom)
	if total == 0 {
		return ""
	}

	return fmt.Sprintf(consts2.RelatingExportTemplate, total)
}

func (ei *ExportIssue) getWorkHourExportInfo(val interface{}) string {
	workHour := &bo.WorkHour{}
	err := json.FromJson(json.ToJsonIgnoreError(val), workHour)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(consts2.WorkHourExportTemplate, workHour.PlanHour, workHour.ActualHour)
}

func (ei *ExportIssue) getParentDescExportInfo(lcData map[string]interface{}) string {
	valStr := ""
	parentId := cast.ToInt64(lcData[consts.BasicFieldParentId])
	if parentId == 0 {
		valStr = WordTranslate("父任务")
	} else {
		valStr = WordTranslate("子任务")
	}

	return valStr
}

func (ei *ExportIssue) getGroupSelectExportInfo(val interface{}, column *projectvo.TableColumnData,
	projectInfo *bo.ProjectBo, lcData map[string]interface{}) string {
	valId := cast.ToString(val)
	options := ei.GetColumnOptions(column)
	optionsMap := make(map[string]*projectvo.ColumnSelectOption, 0)
	for _, option := range options {
		optionsMap[cast.ToString(option.Id)] = option
	}
	valStr := ""
	// 任务状态特殊逻辑，判断确认状态
	if projectInfo.ProjectTypeId == consts.ProjectTypeNormalId {
		// 对于通用项目，状态值有两种：未开始、已完成。如果是已完成，则还需判断它的 auditStatus 吗，如果不为 3（审批通过），则表示待确认。
		if statusBo, ok := optionsMap[valId]; ok {
			valStr = statusBo.Value
			if cast.ToInt(statusBo.ParentId) == consts.StatusTypeComplete {
				if cast.ToInt(lcData[consts.BasicFieldAuditStatus]) != consts.AuditStatusPass {
					valStr = consts.WaitConfirmStatusName
				}
			}
		}
	} else {
		if targetOption, ok := optionsMap[valId]; ok {
			valStr = targetOption.Value
		}
	}

	// 翻译
	if tmpMap, ok1 := consts.LANG_STATUS_MAP[lang2.GetLang()]; ok1 {
		if newStr, ok := tmpMap[valStr]; ok {
			valStr = newStr
		}
	}

	return valStr
}

func (ei *ExportIssue) getNumberExportInfo(val interface{}, propsMap map[string]interface{}) string {
	valFloat := cast.ToFloat64(val)
	renderedValStr := ""
	// 根据字段规则进行展示
	propsStr := json.ToJsonIgnoreError(propsMap)
	propsObj := projectvo.FormConfigColumnFieldInputNumberProps{}
	json.FromJson(propsStr, &propsObj)
	// 精度
	accuracyNum := 0
	coefficientNum := 1 // 系数，为了不四舍五入。比如精度为1.0，则乘以10，向下取整，再除以该系数 10.
	switch propsObj.Inputnumber.Accuracy {
	case "1.0":
		accuracyNum = 1
		coefficientNum = 10
	case "1.00":
		accuracyNum = 2
		coefficientNum = 100
	case "1.000":
		accuracyNum = 3
		coefficientNum = 1000
	case "1.0000":
		accuracyNum = 4
		coefficientNum = 10000
	default:
		accuracyNum = 0
	}
	// 先乘以一个系数，再除以它。去掉多余的后缀，用于去除四舍五入。
	valFloat, _ = decimal.NewFromFloat(float64(coefficientNum)).Mul(decimal.NewFromFloat(valFloat)).Float64()
	valFloat = math.Floor(valFloat)
	// 千分位
	if propsObj.Inputnumber.Thousandth {
		// f1, _ := decimal.NewFromFloat(valFloat).Round(int32(accuracyNum)).Float64()
		f1, _ := decimal.NewFromFloat(valFloat).Div(decimal.NewFromFloat(float64(coefficientNum))).RoundBank(int32(accuracyNum)).Float64()
		// 打印出千分位的值 https://www.cnblogs.com/DillGao/p/8986602.html
		p := message.NewPrinter(language.English)
		renderedValStr = p.Sprintf("%."+strconv.Itoa(accuracyNum)+"f", f1)
	} else {
		renderedValStr = decimal.NewFromFloat(valFloat).Div(decimal.NewFromFloat(float64(coefficientNum))).StringFixedBank(int32(accuracyNum))
	}
	// 百分比符号
	if propsObj.Inputnumber.Percentage {
		renderedValStr += "%"
	}

	return renderedValStr
}

func (ei *ExportIssue) getAmountExportInfo(val interface{}, propsMap map[string]interface{}) string {
	valFloat := cast.ToFloat64(val)
	renderedValStr := ""
	// 根据字段规则进行展示
	propsStr := json.ToJsonIgnoreError(propsMap)
	propsObj := projectvo.FormConfigColumnFieldAmountProps{}
	json.FromJson(propsStr, &propsObj)
	// 货币的渲染根据三个规则：符号、前后缀标识、精度
	// 精度
	switch propsObj.Amount.Accuracy {
	case "1.0":
	case "1.00":
	case "1.000":
	case "1.0000":
	default:
	}
	unitFlag := propsObj.Amount.Sign
	switch propsObj.Amount.Location {
	case "suffix": // 后缀
		renderedValStr = fmt.Sprintf("%[1]v%[2]s", valFloat, unitFlag)
	default: // 默认：前缀
		renderedValStr = fmt.Sprintf("%[2]s%[1]v", valFloat, unitFlag)
	}

	return renderedValStr
}

func (ei *ExportIssue) getSelectExportInfo(val interface{}, column *projectvo.TableColumnData) string {
	options := ei.GetColumnOptions(column)
	optionsMap := make(map[string]*projectvo.ColumnSelectOption, 0)
	for _, option := range options {
		optionsMap[cast.ToString(option.Id)] = option
	}

	valNameStr := ""
	if tmpOption, ok := optionsMap[cast.ToString(val)]; ok {
		valNameStr = tmpOption.Value
	}

	return valNameStr
}

func (ei *ExportIssue) getMultiSelectExportInfo(val interface{}, column *projectvo.TableColumnData) string {
	valJson := json.ToJsonIgnoreError(val)
	if valJson == "" {
		return ""
	}

	options := ei.GetColumnOptions(column)
	optionsMap := make(map[interface{}]*projectvo.ColumnSelectOption, 0)
	for _, option := range options {
		optionsMap[option.Id] = option
	}
	valArr := make([]interface{}, 0)
	valNameArr := make([]string, 0)
	json.FromJson(valJson, &valArr)
	for _, oneVal := range valArr {
		if tmpOption, ok := optionsMap[oneVal]; ok {
			valNameArr = append(valNameArr, tmpOption.Value)
		}
	}

	return strings.Join(valNameArr, "||")
}

func (ei *ExportIssue) getDeptExportInfo(val interface{}, userDepts map[string]*uservo.MemberDept) string {
	stringVals := cast.ToStringSlice(val)
	deptNameArr := make([]string, 0, len(stringVals))
	for _, stringVal := range stringVals {
		if !strings.Contains(stringVal, consts.LcCustomFieldDeptType) {
			stringVal = consts.LcCustomFieldDeptType + stringVal
		}
		if userDepts[stringVal] != nil {
			deptNameArr = append(deptNameArr, userDepts[stringVal].Name)
		}
	}

	return strings.Join(deptNameArr, consts2.ExcelJoinFlag)
}

func (ei *ExportIssue) getConditionRefExportInfo(exportObj *ExportIssue, val interface{}, propsMap map[string]interface{}, userDepts map[string]*uservo.MemberDept) (string, []*bo.DocumentBo) {
	title := ""
	allDocuments := make([]*bo.DocumentBo, 0)
	if m, ok := propsMap[consts.LcColumnFieldTypeConditionRef].(map[string]interface{}); ok && m["column"] != nil {
		column := &projectvo.TableColumnData{}
		_ = json.FromJson(json.ToJsonIgnoreError(m["column"]), column)
		switch values := val.(type) {
		case []interface{}:
			titles := make([]string, 0, len(values))
			for _, dataValue := range values {
				tempExportObjet := &ExportIssue{
					InAndExportBase: InAndExportBase{
						Columns:     []*projectvo.TableColumnData{column},
						ProjectInfo: exportObj.ProjectInfo,
					},
				}
				tempBo := ei.TransferLcDataIntoExportData(tempExportObjet, map[string]interface{}{column.Name: dataValue}, userDepts)

				allDocuments = append(allDocuments, tempBo.Documents...)
				if str := cast.ToString(tempBo.ExportValues[column.Name]); str != "" {
					titles = append(titles, str)
				}
			}
			title = strings.Join(titles, consts2.ExcelJoinFlag)
		default:
			title = cast.ToString(val)
		}
	}

	return title, allDocuments
}

func (ei *ExportIssue) getDocumentExportInfo(val interface{}) []*bo.DocumentBo {
	documents := make([]*bo.DocumentBo, 0)
	if val == nil {
		return documents
	}

	documents = ei.getDocumentsFromValue(val)
	newDocuments := make([]*bo.DocumentBo, 0, len(documents))
	for _, document := range documents {
		if document.RecycleFlag == consts.AppIsNoDelete {
			if document.ResourceType == consts.FsResource {
				ei.setDocumentSuffixAndIcon(document, "/", "docs")
			} else {
				ei.setDocumentSuffixAndIcon(document, ".", "txt")
			}
			document.DisplayPath = document.Path
			if _, ok := consts2.ExportOfficeFileType[document.Suffix]; ok {
				document.DisplayPath = consts2.OfficeViewUrl + document.Path
			} else if document.Suffix == "pdf" {
				document.DisplayPath = consts2.PdfVideUrl + document.Path
			}
			newDocuments = append(newDocuments, document)
		}
	}

	return newDocuments
}

func (ei *ExportIssue) getDocumentsFromValue(val interface{}) []*bo.DocumentBo {
	documents := make([]*bo.DocumentBo, 0)
	if val == nil {
		return documents
	}
	valJson := json.ToJsonIgnoreError(val)
	if valJson == "" {
		return documents
	}

	switch val.(type) {
	case []interface{}:
		_ = json.FromJson(valJson, &documents)
	case map[string]interface{}:
		documentsMap := make(map[string]*bo.DocumentBo)
		_ = json.FromJson(valJson, &documentsMap)
		for _, documentBo := range documentsMap {
			documents = append(documents, documentBo)
		}
	}

	return documents
}

func (ei *ExportIssue) setDocumentSuffixAndIcon(document *bo.DocumentBo, sep, defaultSuffix string) {
	paths := strings.Split(document.Path, sep)
	suffix := strings.ToLower(paths[len(paths)-1])
	if png, ok := consts2.ExportFileTypeIconMap[suffix]; ok {
		document.Suffix = suffix
		document.Icon = png
	} else {
		document.Suffix = defaultSuffix
		document.Icon = consts2.ExportFileTypeIconMap[defaultSuffix]
	}
	document.Name = md5.Md5V(document.Path) + "." + document.Suffix
}

// getDatepickerExportInfo 从 pg 中查询出来后，转换成对象，返回给调用者
func (ei *ExportIssue) getDatepickerExportInfo(val interface{}, columnKey string) string {
	resTimeStr := cast.ToString(val)
	// 开始时间和结束时间特殊逻辑，不设置的填了1970年，需要改成不设置不要给默认值了，蛋疼
	if columnKey == consts.BasicFieldPlanEndTime || columnKey == consts.BasicFieldPlanStartTime {
		if resTimeStr < consts.BlankElasticityTime {
			resTimeStr = ""
		}
	}
	return resTimeStr
}

// getMemberExportInfo 处理成员数据，转换为对象
func (ei *ExportIssue) getMemberExportInfo(val interface{}, userDepts map[string]*uservo.MemberDept) string {
	stringVals := cast.ToStringSlice(val)
	if len(stringVals) == 0 {
		stringVal := cast.ToString(val)
		if stringVal != "" {
			stringVals = append(stringVals, cast.ToString(val))
		}
	}
	names := make([]string, 0, len(stringVals))
	for _, stringVal := range stringVals {
		if !strings.Contains(stringVal, consts.LcCustomFieldUserType) {
			stringVal = consts.LcCustomFieldUserType + stringVal
		}
		if userDepts[stringVal] != nil {
			names = append(names, userDepts[stringVal].Name)
		}
	}

	return strings.Join(names, "||")
}
