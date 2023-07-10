package service

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cast"

	"github.com/star-table/polaris-backend/common/model/vo/formvo"

	"github.com/star-table/polaris-backend/common/core/util/http"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/core/util/excel/sortedmap"
	lang2 "github.com/star-table/polaris-backend/common/core/util/lang"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/appvo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/appfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	consts2 "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

// ImportIssues 通过 excel 数据获取任务内置字段的头部，然后追加“任务类型”下的自定义字段头部。
// 经过任务状态改造后，ProjectObjectTypeID 参数被改为 tableId。
// 通用项目可以支持一个 excel 中导入任务到多个不同的任务栏（同项目）中
func ImportIssues(orgId, currentUserId int64, input vo.ImportIssuesReq) (projectvo.ImportIssuesRespVoData, errs.SystemErrorInfo) {
	result := projectvo.ImportIssuesRespVoData{}

	// 权限检查
	err := domain.AuthProject(orgId, currentUserId, input.ProjectID, consts.RoleOperationPathOrgProProConfig, consts.OperationProIssue4Import)
	if err != nil {
		log.Errorf("[ImportIssues] domain.AuthProject err: %v", err)
		return result, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	tableIdInt64 := cast.ToInt64(input.TableID)
	if tableIdInt64 < 1 {
		return result, errs.InvalidTableId
	}
	importIssue, err := domain.NewImportIssue(orgId, input.ProjectID, tableIdInt64, input.URL)
	if err != nil {
		return result, err
	}
	projectInfo := importIssue.ProjectInfo

	if isExecuting := domain.CheckAsyncTaskIsRunning(orgId, projectInfo.AppId, tableIdInt64); isExecuting {
		err := errs.ImportAsyncTaskIsExecuting
		log.Errorf("[ImportIssues] domain.CheckAsyncTaskImportIsRunning err: %v", err)
		return result, err
	}

	// 导入时，一般是导入到表格视图中，但是无法通过 projectId 得到其默认的表格视图。默认使用“表格视图”
	param := map[string]string{
		"usingType": consts.FilterColumnUsingTypeImport,
	}
	_, _, allColumns, err := GetBuiltInAndCustomColumns(orgId, input.ProjectID, tableIdInt64,
		projectInfo.ProjectTypeId, true, param, false)
	if err != nil {
		log.Error(err)
		return result, err
	}

	err = importIssue.SetColumns(allColumns)
	if err != nil {
		log.Error(err)
		return result, err
	}

	data, ignoreRows, headerMap, importErr := importIssue.ParseExcel()
	if importErr != nil {
		log.Errorf("[ImportIssues] excel.GenerateCSVFromXLSXFileWithMap err: %v", importErr)
		return result, errs.BuildSystemErrorInfo(errs.InvalidImportFile, importErr)
	}

	groupData := data.([]*sortedmap.SortMap)
	importData, err := domain.FilterEmptyRowForImportIssue(groupData)
	if err != nil {
		log.Errorf("domain.FilterEmptyRowForImportIssue err: %v", err)
		return result, err
	}

	// 检查任务总数限制
	authFunctionErr := domain.AuthPayTask(orgId, consts.FunctionTaskLimit, len(importData))
	if authFunctionErr != nil {
		log.Errorf("[ImportIssues] domain.AuthPayTask err: %v", authFunctionErr)
		return result, authFunctionErr
	}

	// 根据导入的 cell 值，和已知的用户数据匹配处理
	count, err := handleImportData(orgId, currentUserId, importData, &input, importIssue, ignoreRows, headerMap)
	if err != nil {
		log.Errorf("[ImportIssues] handleImportData err: %v", err)
		return result, err
	}
	result.Count = count
	result.TaskId = importIssue.AsyncTaskId
	result.UnknownImportColumns = importIssue.UnknownImportColumns

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = projectInfo.Name
		domain.PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeCreateIssueBatch,
			OrgId:      orgId,
			ProjectId:  input.ProjectID,
			OperatorId: currentUserId,
			Ext:        ext,
		})
	})

	return result, nil
}

// handleImportData 处理导入的数据。其中，导入数据的每一行是以 map 的方式存储。
func handleImportData(orgId, userId int64, importData []*sortedmap.SortMap, input *vo.ImportIssuesReq, importIssue *domain.ImportIssue, ignoreRows int, headerMap *sortedmap.SortMap) (int64, errs.SystemErrorInfo) {
	colIndexMap := headerMap.GetMapKeyByKey()
	errTipsPrevFormat := domain.WordTranslate(" 第%d行,第%d列错误：%s")

	var errStrList []string
	var datas []map[string]interface{}
	var lastParentData map[string]interface{}           // 最近的一条父任务
	var lastParentDataChildren []map[string]interface{} // 最近的一条父任务的子任务集
	var totalCount int64

	for k, v := range importData {
		// 任务类型：是否是子任务。
		isChild := v.GetByKey(consts.HelperFieldIsParentDesc, "") == domain.WordTranslate("子任务")

		// 第一行不能是子任务。
		if lastParentData == nil && isChild {
			return 0, errs.ChildIssueForFirst
		}

		data, err := importIssue.HandleImportData(v, colIndexMap)
		if err != nil {
			log.Errorf("[handleImportData] buildImportData err: %v", err)
			return 0, err
		}
		totalCount += 1

		if isChild {
			// 子任务
			lastParentDataChildren = append(lastParentDataChildren, data)
			lastParentData[consts.TempFieldChildren] = lastParentDataChildren
			lastParentData[consts.TempFieldChildrenCount] = cast.ToInt64(lastParentData[consts.TempFieldChildrenCount]) + 1
		} else {
			// 父任务
			lastParentData = data
			lastParentDataChildren = data[consts.TempFieldChildren].([]map[string]interface{})
			datas = append(datas, data)
		}

		if len(importIssue.ErrMsgList) > 0 {
			var tmpErrColNumList []int
			for colNum, _ := range importIssue.ErrMsgList {
				tmpErrColNumList = append(tmpErrColNumList, colNum)
			}
			sort.Ints(tmpErrColNumList)
			for _, colNum := range tmpErrColNumList {
				if msgArr, ok := importIssue.ErrMsgList[colNum]; ok {
					for _, oneMsg := range msgArr {
						errStrList = append(errStrList, fmt.Sprintf(errTipsPrevFormat, k+ignoreRows+1, colNum, domain.WordTranslate(oneMsg)))
					}
				}
			}
		}
	}
	if len(errStrList) > 0 {
		err := errs.BuildSystemErrorInfoWithMessage(errs.FileParseFail, json.ToJsonIgnoreError(errStrList))
		log.Errorf("[handleImportData] err: %v", err)
		return 0, err
	}

	// 异步任务
	tableIds := []string{cast.ToString(importIssue.TableInfo.TableId)}
	err := domain.CreateAsyncTask(orgId, totalCount, importIssue.AsyncTaskId, map[string]string{
		consts.AsyncTaskHashPartKeyOfTableIds: json.ToJsonIgnoreError(tableIds),
	})
	if err != nil {
		log.Errorf("[handleImportData] CreateAsyncTask err: %v", err)
		return 0, err
	}

	// 异步批量创建
	req := &projectvo.BatchCreateIssueReqVo{
		OrgId:     orgId,
		UserId:    userId,
		AppId:     importIssue.ProjectInfo.AppId,
		ProjectId: input.ProjectID,
		TableId:   importIssue.TableInfo.TableId,
		Data:      datas,
	}
	AsyncBatchCreateIssue(req, true, &projectvo.TriggerBy{
		TriggerBy: consts.TriggerByImport,
	}, importIssue.AsyncTaskId)
	//	// 创建任务不成功，则释放缓存
	//	domain.ClearAsyncTaskCacheForCreateIssue(orgId, infoMap.AsyncTaskId)
	//	log.Errorf("[ImportIssues] DealCreateIssue err: %v", err)
	//	return 0, err
	//}

	return int64(len(importData)), nil
}

// ExportIssueTemplate 导出任务模板
func ExportIssueTemplate(orgId int64, projectId int64, tableId int64) (string, errs.SystemErrorInfo) {
	createExcel, err := domain.NewImportIssueTemplate(orgId, projectId, tableId)
	if err != nil {
		log.Errorf("[ExportIssueTemplate] NewFile orgId:%v,projectId:%v,tableId:%v, err:%v", orgId, projectId, tableId, err)
		return "", err
	}

	// 获取表头
	// 下载一个模板时，也需要指定特定的视图。默认使用“表格视图”
	param := map[string]string{"usingType": "import"}
	builtInColumns, customColumns, allColumns, busiErr := GetBuiltInAndCustomColumns(orgId, projectId, tableId, createExcel.ProjectInfo.ProjectTypeId, true, param, false)
	if busiErr != nil {
		log.Error(busiErr)
		return "", busiErr
	}
	busiErr = createExcel.SetColumns(mergeColumns(builtInColumns, customColumns))
	if busiErr != nil {
		log.Error(busiErr)
		return "", busiErr
	}

	skippedColumnNames := make([]string, 0, 10)
	for _, column := range allColumns {
		isIn := false
		for _, data := range createExcel.Columns {
			if column.Name == data.Name {
				isIn = true
			}
		}
		if !isIn {
			skippedColumnNames = append(skippedColumnNames, "【"+column.Label+"】")
		}
	}

	_ = createExcel.SetNoticeText(createExcel.ProjectInfo.ProjectTypeId, skippedColumnNames)
	_ = createExcel.SetExcelHeadKeyAndLockSomeRow()
	createExcel.SetExcelHeadForImportTpl()
	enumListForParentFlag := []string{domain.WordTranslate("父任务"), domain.WordTranslate("子任务")}
	createExcel.SetValidateDropList(fmt.Sprintf("%c", 'B'), enumListForParentFlag)

	err = createExcel.Save()
	if err != nil {
		businessErr := errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
		log.Errorf("[ExportIssueTemplate] err: %v", businessErr)
		return "", businessErr
	}

	return createExcel.FileUrl, nil
}

// ExportUserOrDeptSameNameList 导出同名的用户、部门列表
func ExportUserOrDeptSameNameList(orgId int64, projectId int64) (string, errs.SystemErrorInfo) {
	excelUrl := ""
	excelFileName := "重名的成员和部门列表.xlsx"
	if lang2.IsEnglish() {
		excelFileName = "duplicate-member-or-dept-list.xlsx"
	}
	// 根据数据列表生成 excel 文件路径信息
	relatePath, fileDir, fileName, err := GetExportFileInfo(orgId, projectId, excelFileName)
	if err != nil {
		log.Error(err)
		return excelUrl, err
	}
	fullFileName := fileDir + "/" + fileName

	deptTransData, userTransData, err := domain.GetSameNameUserAndDeptSimpleInfo(orgId, "all")
	if err != nil {
		log.Errorf("[ExportUserOrDeptSameNameList] err: %v, orgId: %d", err, orgId)
		return excelUrl, err
	}

	//if len(userInfoResp.SimpleUserInfo) < 1 && len(sameNameDeptList) < 1 {
	//	return excelUrl, errs.DeptAndUserNoSameName
	//}

	if err = domain.WriteUserOrDeptSameNameListToExcel(orgId, fullFileName, userTransData, deptTransData); err != nil {
		log.Errorf("[ExportUserOrDeptSameNameList] domain.WriteUserOrDeptSameNameListToExcel err: %v", err)
		return excelUrl, err
	}
	excelUrl = config.GetOSSConfig().LocalDomain + relatePath + "/" + fileName

	return excelUrl, nil
}

// GetUserOrDeptSameNameList 获取同名数据
// dataType 数据类型, `user` 同名用户; `dept` 同名部门; all 所有
func GetUserOrDeptSameNameList(orgId int64, dataType string) (projectvo.GetUserOrDeptSameNameListRespData, errs.SystemErrorInfo) {
	result := projectvo.GetUserOrDeptSameNameListRespData{}
	userTransData, deptTransData, err := domain.GetSameNameUserAndDeptSimpleInfo(orgId, dataType)
	if err != nil {
		log.Errorf("[GetUserOrDeptSameNameList] err: %v, orgId: %d", err, orgId)
		return result, err
	}
	list := make([]projectvo.GetUserOrDeptSameNameListItem, 0, len(userTransData))
	switch dataType {
	case "user":
		for _, item := range userTransData {
			list = append(list, projectvo.GetUserOrDeptSameNameListItem{
				Id:   strconv.FormatInt(item.Id, 10),
				Name: item.Name,
			})
		}
	case "dept":
		for _, item := range deptTransData {
			list = append(list, projectvo.GetUserOrDeptSameNameListItem{
				Id:   strconv.FormatInt(item.Id, 10),
				Name: item.Name,
			})
		}
	}
	result.List = list

	return result, nil
}

// GetExportFileInfo 获取导出的文件相关信息
// fileName string 导出的文件名，如：“成员导入模板.xlsx”
func GetExportFileInfo(orgId, projectId int64, fileName string) (relatePath, fileDir, outFileName string, err errs.SystemErrorInfo) {
	relatePath = "/org_" + strconv.FormatInt(orgId, 10) + "/project_" + strconv.FormatInt(projectId, 10)
	// log.Info(config.GetOSSConfig())
	fileDir = strings.TrimRight(config.GetOSSConfig().RootPath, "/") + relatePath
	mkdirErr := os.MkdirAll(fileDir, 0777)
	if mkdirErr != nil {
		log.Errorf("[GetExportFileInfo] err: %v", mkdirErr)
		return "", "", "", errs.BuildSystemErrorInfo(errs.IssueDomainError, mkdirErr)
	}
	outFileName = fileName
	return relatePath, fileDir, outFileName, nil
}

// GetBuiltInAndCustomColumns 获取内置字段列表和自定义字段列表，并返回项目的所有列信息
// needHidden 为 true 表示返回隐藏的字段。对于导入，最好拿到所有的字段信息（包含隐藏），用于匹配。
// proTypeId 项目类型：1 通用；2 敏捷
// extraParam 额外的参数，其中有个**默认**值：{usingType: import}
// 	* usingType 表示用途，`import` 导入；`export` 导出；
func GetBuiltInAndCustomColumns(orgId, projectId, tableId int64, proTypeId int64, needHidden bool, extraParam map[string]string, isNeedDocument bool,
) ([]*projectvo.TableColumnData, []*projectvo.TableColumnData, []*projectvo.TableColumnData, errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	usingType := consts.FilterColumnUsingTypeImport
	if extraParam["usingType"] != "" {
		usingType = extraParam["usingType"]
	}

	// 如果不是有效的 tableId，则直接返回
	if tableId < 1 {
		return nil, nil, nil, errs.ParamError
	}
	// 向模板中插入自定义字段列，这些列从无码接口中获取
	builtInColumns := make([]*projectvo.TableColumnData, 0)
	customColumns := make([]*projectvo.TableColumnData, 0)
	// 过滤完排好序的字段名列表
	sortedColumns := make([]*projectvo.TableColumnData, 0)
	appId, err := domain.GetAppIdFromProjectId(orgId, projectId)
	if err != nil {
		log.Errorf("[GetBuiltInAndCustomColumns] domain.GetAppIdFromProjectId err: %v", err)
		return nil, nil, nil, err
	}

	columns, _, err := domain.GetTableColumnsForImport(orgId, tableId)
	if err != nil {
		log.Errorf("[GetBuiltInAndCustomColumns] GetTableColumnsForImport err: %v", err)
		return builtInColumns, customColumns, nil, err
	}
	originalColumns := columns

	// 根据前端展示字段列的规则，优先从“表格视图”中取出字段，如果找不到合适的视图
	// 	* 敏捷项目：使用 config 中，objectTypeId 对应的字段列表，如果没有 objectTypeId 对应的数据，则使用 columns 中的列（去除 isSys 为 true 的列）
	// 	* 通用项目：如果表格视图不存在，则直接取出 columns 中的列作为要任务数据列（去除 isSys 为 true 的列）
	// 详细规则详见：https://startable.feishu.cn/wiki/wikcn6uGzbNJrtarGvWDi6MTuFe#HiN3w7
	_, viewConfigObj, err := domain.GetCustomFieldInfoFromViewList(orgId, appId, tableId, usingType)
	if err != nil {
		log.Errorf("[GetBuiltInAndCustomColumns] domain.GetCustomFieldInfoFromViewList err: %v", err)
		return nil, nil, nil, err
	}
	hiddenColumnKeys := viewConfigObj.HiddenColumnIds

	// 将 column 列表转成 key 为键的 map
	columnsMap := domain.GetColumnMapByColumnsForTableColumn(columns)

	switch int(proTypeId) {
	case consts.ProjectTypeNormalId:
		// 通用项目
		// 优先使用视图中的有序列名
		columnKeyArr := make([]string, 0)
		if viewConfigObj != nil && len(viewConfigObj.TableOrder) > 0 {
			columnKeyArr = viewConfigObj.TableOrder
		}
		columnKeyArr = domain.FilterEmptyColumnKey(columnKeyArr)
		sortedColumns = domain.FilterCustomFieldForCommonPro(columnKeyArr, columns, originalColumns, columnsMap, needHidden, hiddenColumnKeys)
	case consts.ProjectTypeAgileId, consts.ProjectTypeCommon2022V47:
		// 敏捷项目，从 columns 中取
		sortedColumns = domain.FilterCustomFieldForAgilePro(nil, columns, columnsMap, needHidden, hiddenColumnKeys)
	default:
		sortedColumns = columns
	}
	columns = sortedColumns
	if len(columns) == 0 {
		err = errs.IssueColumnIsEmpty
		log.Errorf("[GetBuiltInAndCustomColumns] err: %v", err)
		return builtInColumns, customColumns, columns, err
	}

	// 对列进行一些特殊 filter
	columns, builtInColumns, customColumns, err = domain.FilterForBuiltInAndCustomColumns(columns, proTypeId, usingType, isNeedDocument)
	if err != nil {
		log.Errorf("[GetBuiltInAndCustomColumns] err: %v", err)
		return builtInColumns, customColumns, originalColumns, err
	}

	return builtInColumns, customColumns, columns, nil
}

// ExportData 导出任务数据
func ExportData(orgId int64, req projectvo.ExportIssueReqVo) (string, errs.SystemErrorInfo) {
	input := req.Input
	curUserId := req.UserId
	projectId := input.ProjectId
	iterationId := input.IterationId
	tableIdStr := input.TableId
	tableId, oriErr := strconv.ParseInt(tableIdStr, 10, 64)
	if oriErr != nil {
		log.Errorf("[ExportData] err: %v", oriErr)
		return "", errs.BuildSystemErrorInfo(errs.TypeConvertError, oriErr)
	}
	exportObj, businessErr := domain.NewExportIssue(orgId, projectId, tableId, iterationId)
	if businessErr != nil {
		log.Errorf("[ExportData] err: %v, projectId: %d", businessErr, projectId)
		return "", businessErr
	}

	log.Infof("[ExportData] projectId: %d, isCommon: %v", projectId, exportObj.ProjectInfo.ProjectTypeId == consts.ProjectTypeNormalId)
	// 权限检查
	//busiErr := domain.AuthProject(orgId, curUserId, projectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProIssue4Export)
	//if busiErr != nil {
	//	log.Error(busiErr)
	//	return "", errs.BuildSystemErrorInfo(errs.Unauthorized, busiErr)
	//}
	// 获取自定义字段
	exportColumns, err1 := getExportColumns(orgId, projectId, tableId, exportObj.ProjectInfo.ProjectTypeId, req.Input.IsNeedDocument)
	if err1 != nil {
		log.Errorf("[ExportData] getExportColumns:%v", err1)
		return "", err1
	}
	err1 = exportObj.SetColumns(exportColumns)
	if err1 != nil {
		log.Errorf("[ExportData] SetColumns:%v", err1)
		return "", err1
	}

	// 设定表头
	setHeaderForExportIssue(exportObj)

	// 查询所有任务
	allIssueList, businessErr := GetExportIssues(orgId, curUserId, tableId,
		exportObj, &input)
	if businessErr != nil {
		log.Errorf("[ExportData] GetExportIssues err: %v", businessErr)
		return "", businessErr
	}

	// 下载所有图片缩略图
	downLoadExportImage(exportObj, allIssueList)

	// 获取父任务
	issuesParent := make([]*bo.ExportIssueBo, 0)
	issueIds := make([]int64, 0, len(allIssueList))
	childrenMap := make(map[int64][]*bo.ExportIssueBo, len(issuesParent))
	for _, issue := range allIssueList {
		issueIds = append(issueIds, issue.Id)
		if issue.ParentId > 0 {
			childrenMap[issue.ParentId] = append(childrenMap[issue.ParentId], issue)
		} else {
			issuesParent = append(issuesParent, issue)
		}
	}

	// 拼装sheet数据
	for _, issue := range issuesParent {
		insertCell(exportObj, issue, exportColumns)
		insertChildIssue(exportObj, issue, childrenMap, exportColumns)
	}

	err := exportObj.Fd.SaveAs(exportObj.FileNameWithPath)
	if err != nil {
		log.Errorf("[ExportData] err: %v", err)
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	return exportObj.FileUrl, nil
}

func getExportColumns(orgId, projectId, tableId, projectTypeId int64, isNeedDocument bool) ([]*projectvo.TableColumnData, errs.SystemErrorInfo) {
	// 获取自定义字段
	param := map[string]string{"usingType": "export"}
	builtInColumns, customColumns, _, err1 := GetBuiltInAndCustomColumns(orgId, projectId, tableId, projectTypeId, true, param, isNeedDocument)
	if err1 != nil {
		log.Errorf("[getExportColumns] GetBuiltInAndCustomColumns err:%v", err1)
		return nil, err1
	}

	return mergeColumns(builtInColumns, customColumns), nil
}

func mergeColumns(builtInColumns, customColumns []*projectvo.TableColumnData) []*projectvo.TableColumnData {
	// 父子描述放在第二个位置
	exportColumns := make([]*projectvo.TableColumnData, 0, len(builtInColumns)+len(customColumns))
	if len(builtInColumns) > 0 {
		exportColumns = append(exportColumns, builtInColumns[0])
		exportColumns = append(exportColumns, consts2.ParentDescColumn)
		exportColumns = append(exportColumns, builtInColumns[1:]...)
	} else {
		exportColumns = append(exportColumns, consts2.ParentDescColumn)
	}

	exportColumns = append(exportColumns, customColumns...)

	return exportColumns
}

// 每次插入后，回再次检查是否还有子任务。
func insertChildIssue(sheet *domain.ExportIssue, issueInfo *bo.ExportIssueBo,
	issuesChildrenMap map[int64][]*bo.ExportIssueBo, customColumns []*projectvo.TableColumnData,
) {
	currentBo := issueInfo
	curChildren := issuesChildrenMap[currentBo.Id]

	for _, infoBo := range curChildren {
		insertCell(sheet, infoBo, customColumns)
		// 针对当前子任务，再次检查是否有子任务。
		insertChildIssue(sheet, infoBo, issuesChildrenMap, customColumns)
	}
}

func GetExportIssues(orgId, curUserId, tableId int64,
	exportObj *domain.ExportIssue, input *projectvo.ExportIssueReqVoData,
) ([]*bo.ExportIssueBo, errs.SystemErrorInfo) {

	allColumnKeys := make([]string, 0, len(exportObj.Columns))
	for _, item := range exportObj.Columns {
		if item.Name != "" {
			allColumnKeys = append(allColumnKeys, item.Name)
		}
	}

	// 目前是导出全量的。最多能导出 8k 条。
	page, size := int64(1), int64(8000)
	queryParam := &vo.HomeIssueInfoReq{
		LessConds: &vo.LessCondsData{
			Type:  "and",
			Conds: make([]*vo.LessCondsData, 0, 8),
		},
		FilterColumns: allColumnKeys, //infoMapObj.AllColumnKeys
	}
	queryParam.LessConds.Conds = append(queryParam.LessConds.Conds, &vo.LessCondsData{
		Column: consts.BasicFieldRecycleFlag, // lc_helper.ConvertToCondColumn(consts.BasicFieldRecycleFlag),
		Type:   consts.ConditionEqual,
		Value:  consts.DeleteFlagNotDel,
	})

	if tableId > 0 {
		queryParam.LessConds.Conds = append(queryParam.LessConds.Conds, &vo.LessCondsData{
			Column: consts.BasicFieldTableId, // fmt.Sprintf("\"data\"::jsonb -> '%s'", consts.BasicFieldTableId),
			Type:   consts.ConditionEqual,
			Value:  fmt.Sprintf("%d", tableId), // fmt.Sprintf("%d", tableId),
		})
	}
	// 我协作的
	if domain.CheckNeedCollaborateFilter(orgId, curUserId, tableId, exportObj.ProjectInfo) {
		isMemberResp := appfacade.IsAppMember(appvo.IsAppMemberReq{
			AppId:  exportObj.ProjectInfo.AppId,
			OrgId:  orgId,
			UserId: curUserId,
		})
		if isMemberResp.Failure() {
			log.Error(isMemberResp.Error())
			return nil, isMemberResp.Error()
		}
		// 不是项目成员，则只能看到他协作的任务；否则可以看到项目表下的所有任务
		if !isMemberResp.Data {
			resp := orgfacade.GetUserDeptIds(&orgvo.GetUserDeptIdsReq{
				OrgId:  orgId,
				UserId: curUserId,
			})
			if resp.Failure() {
				log.Errorf("[LcHomeIssuesForProject] orgId:%v, curUserId:%v, GetUserDeptIds failure:%v", orgId, curUserId, resp.Error())
				return nil, resp.Error()
			}
			deptIdsResp := resp

			// 筛选条件：我协作的任务
			queryParam.LessConds.Conds = append(queryParam.LessConds.Conds,
				domain.GetCollaborateDataQueryConds(orgId, curUserId, deptIdsResp.Data.DeptIds))
		}
	}
	if input != nil {
		if input.Condition != nil {
			queryParam.LessConds = input.Condition
		}
		if len(input.Orders) > 0 {
			queryParam.LessOrder = input.Orders
		}
		if len(input.FilterColumns) > 0 {
			queryParam.FilterColumns = input.FilterColumns
		}
	}
	if queryParam.FilterColumns != nil && len(queryParam.FilterColumns) > 0 {
		filterColumns := []string{
			lc_helper.ConvertToFilterColumn(consts.BasicFieldId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldAuditStatus),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldParentId),
		}
		for _, column := range queryParam.FilterColumns {
			if exportObj.ColumnsMap[column] != nil {
				filterColumns = append(filterColumns, lc_helper.ConvertToFilterColumn(column))
			}
		}
		queryParam.FilterColumns = slice.SliceUniqueString(filterColumns)
	}

	lessResp, _ := domain.GetIssueListWithRefColumn(formvo.LessIssueListReq{
		Condition:     *queryParam.LessConds,
		AppId:         exportObj.ProjectInfo.AppId,
		OrgId:         orgId,
		UserId:        curUserId,
		Page:          page,
		Size:          size,
		FilterColumns: queryParam.FilterColumns,
		TableId:       tableId,
		Export:        true,
		NeedRefColumn: true,
	}, exportObj.Columns, "")
	if lessResp.Failure() {
		log.Errorf("[GetIssueList] GetIssueListWithRefColumn failed, orgId: %d, condition: %s, err: %v", orgId,
			json.ToJsonIgnoreError(queryParam.LessConds), lessResp.Error())
		return nil, lessResp.Error()
	}

	// 任务数据转换
	list := make([]*bo.ExportIssueBo, 0, int(lessResp.Data.Count))
	dataList := make([]map[string]interface{}, 0, lessResp.Data.Count)
	err := json.Unmarshal(lessResp.Data.Data, &dataList)
	if err != nil {
		log.Errorf("[GetRowsExpand] Unmarshal err: %v", err)
		return nil, errs.JSONConvertError
	}
	for _, issueM := range dataList {
		// 处理自定义字段
		list = append(list, exportObj.TransferLcDataIntoExportData(exportObj, issueM, lessResp.UserDepts))
	}

	return list, nil
}

func downLoadExportImage(sheet *domain.ExportIssue, allIssueList []*bo.ExportIssueBo) {
	needDownLoadsMap := make(map[string]string)
	for _, issueBo := range allIssueList {
		for _, document := range issueBo.Documents {
			if _, ok := consts2.ExportDownloadTypesMap[document.Suffix]; ok {
				document.Name = sheet.ImageDownloadDir + document.Name
				needDownLoadsMap[document.Path+"?x-oss-process=image/resize,w_48,h_48"] = document.Name
			}
		}
	}
	if len(needDownLoadsMap) > 0 {
		http.NewDownloader(needDownLoadsMap, 20).DownloadAll()
	}
}

// 设置表头
// 包含任务内置字段以及任务的自定义字段。
func setHeaderForExportIssue(sheet *domain.ExportIssue) {
	// 增加一行，作为表头
	row := sheet.AddRow()
	for _, column := range sheet.Columns {
		cell := row.AddCell()
		cellValue := domain.GetColumnNameWithAlias(column)
		cell.SetCell(cellValue)
	}
}

func insertCell(sheet *domain.ExportIssue, issueInfo *bo.ExportIssueBo, exportColumns []*projectvo.TableColumnData,
) {
	row := sheet.AddRow()
	// 追加自定义字段值
	for _, column := range exportColumns {
		columnKey := column.Name
		cell := row.AddCell()
		if val, ok := issueInfo.ExportValues[columnKey]; ok {
			if column.Field.Type == consts.LcColumnFieldTypeDocument || column.Field.Type == consts.LcColumnFieldTypeImage {
				if val != nil {
					cell.AddImage(val.([]*bo.DocumentBo))
				}
			} else if column.Field.Type == consts.LcColumnFieldTypeConditionRef {
				if documents, ok := val.([]*bo.DocumentBo); ok {
					cell.AddImage(documents)
				} else {
					cell.SetCell(cast.ToString(val))
				}
			} else {
				cell.SetCell(cast.ToString(val))
			}
		}
	}
}

//// ImportIssues2022 导入任务，支持新增或更新模式 todo
//func ImportIssues2022(orgId, curUserId int64, input *projectvo.ImportIssues2022ReqVoData) (projectvo.ImportIssuesRespVoData, errs.SystemErrorInfo) {
//	var err errs.SystemErrorInfo
//	result := projectvo.ImportIssuesRespVoData{}
//	proAppId, _ := strconv.ParseInt(input.AppIdStr, 10, 64)
//	tableIdInt64, _ := strconv.ParseInt(input.TableIdStr, 10, 64)
//
//	projectBo, err := domain.GetProjectByAppId(orgId, proAppId)
//	if err != nil {
//		log.Errorf("[AuthProjectWithCondByAppId] GetProjectByAppId err: %v", err)
//		return result, err
//	}
//
//	// 权限检查
//	err = domain.AuthProject(orgId, curUserId, projectBo.Id, "", consts.OperationProIssue4Import)
//	if err != nil {
//		log.Errorf("[ImportIssues] domain.AuthProject err: %v", err)
//		return result, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
//	}
//
//	tableInfo, err := domain.GetTableInfo(orgId, proAppId, tableIdInt64)
//	if err != nil {
//		log.Errorf("[ImportIssues] tableId: %d, GetTableInfo err: %v", tableIdInt64, err)
//		return result, err
//	}
//	// 导入时，一般是导入到表格视图中，但是无法通过 projectId 得到其默认的表格视图。默认使用“表格视图”
//	param := map[string]string{
//		"usingType": consts.FilterColumnUsingTypeImport,
//	}
//	_, customFields, allColumns, err := GetBuiltInAndCustomColumns(orgId, projectBo.Id, tableIdInt64,
//		projectBo.ProjectTypeId, true, param, false)
//	if err != nil {
//		log.Error(err)
//		return result, err
//	}
//	// 收集拼接任务的信息
//	infoMap, err := domain.GetIssueInfoMapForImportIssue(orgId, proAppId, tableIdInt64, projectBo, nil)
//	if err != nil {
//		log.Errorf("[ImportIssues] domain.PrepareImportInfo err: %v", err)
//		return result, err
//	}
//	infoMap.OperateUserId = curUserId
//	infoMap.AllColumns = allColumns
//	infoMap.CustomColumns = customFields
//	allColumnMap := make(map[string]*projectvo.TableColumnData, len(allColumns))
//	for _, column := range allColumns {
//		columnKey := column.Name
//		allColumnMap[columnKey] = column
//	}
//	infoMap.AllColumnMap = allColumnMap
//	infoMap.ProjectInfo = projectBo
//	infoMap.TableInfo = tableInfo
//
//	dataHeader := sortedmap.NewSortMap()
//	for _, column := range input.ColumnIds {
//		dataHeader.Insert(column, column)
//	}
//	infoMap.DataHeader = dataHeader
//
//	switch input.Action {
//	case "create":
//		result, err = ImportIssues2022ForCreate(orgId, curUserId, input, infoMap)
//	case "update":
//		result, err = ImportIssues2022ForUpdate(orgId, curUserId, input, infoMap)
//	default:
//		return result, errs.ParamError
//	}
//	if err != nil {
//		log.Errorf("[ImportIssues2022] err: %v, orgId: %d", err, orgId)
//		return result, err
//	}
//
//	return result, nil
//}
//
//func ImportIssues2022ForCreate(orgId, curUserId int64, input *projectvo.ImportIssues2022ReqVoData,
//	infoBox *projectvo.IssueInfoMapForImportIssue,
//) (projectvo.ImportIssuesRespVoData, errs.SystemErrorInfo) {
//	result := projectvo.ImportIssuesRespVoData{}
//
//	dataList := make([]*sortedmap.SortMap, 0, len(input.List))
//	for _, item := range input.List {
//		oneRecord := sortedmap.NewSortMap()
//		for _, column := range input.ColumnIds {
//			oneRecord.Insert(column, item[column])
//		}
//		dataList = append(dataList, oneRecord)
//	}
//	count, err := handleImportData(orgId, curUserId, dataList, &vo.ImportIssuesReq{
//		ProjectID: infoBox.ProjectInfo.Id,
//		TableID:   strconv.FormatInt(infoBox.TableInfo.Id, 10),
//	}, infoBox, 0, infoBox.DataHeader)
//	if err != nil {
//		log.Errorf("[ImportIssues] handleImportData err: %v", err)
//		return result, err
//	}
//	result.Count = count
//	result.TaskId = infoBox.AsyncTaskId
//
//	return result, nil
//}
//
//// ImportIssues2022ForUpdate 任务导入-更新
//func ImportIssues2022ForUpdate(orgId, curUserId int64, input *projectvo.ImportIssues2022ReqVoData, infoBox *projectvo.IssueInfoMapForImportIssue,
//) (projectvo.ImportIssuesRespVoData, errs.SystemErrorInfo) {
//	result := projectvo.ImportIssuesRespVoData{}
//	// 解析 dataId，并将其转换为 issueId，再调用方法批量更新
//	dataList, err := domain.DecodeDataIdForImportIssue(infoBox, input.List)
//	if err != nil {
//		log.Errorf("[ImportIssues2022ForUpdate] DecodeDataIdForImportIssue err: %v", err)
//		return result, err
//	}
//	// 转换数据
//	dataList, err = domain.TransferImportDataForUpdate(dataList, infoBox)
//	if err != nil {
//		log.Errorf("[ImportIssues2022ForUpdate] TransferImportDataForUpdate err: %v", err)
//		return result, err
//	}
//	// 批量更新任务。批量更新需支持部分成功、部分失败 todo
//	err = BatchUpdateIssue(&projectvo.BatchUpdateIssueReqInnerVo{
//		OrgId:     orgId,
//		UserId:    curUserId,
//		AppId:     infoBox.ProjectInfo.AppId,
//		ProjectId: infoBox.ProjectInfo.Id,
//		TableId:   infoBox.TableInfo.Id,
//		Data:      dataList,
//	}, true, consts.TriggerByNormal)
//	if err != nil {
//		log.Errorf("[ImportIssues2022ForUpdate] BatchUpdateIssue err: %v", err)
//		return result, err
//	}
//
//	return result, nil
//}
