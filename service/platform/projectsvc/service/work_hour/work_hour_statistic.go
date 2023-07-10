package work_hour

import (
	"fmt"
	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/tealeg/xlsx/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

// 工时统计数据的导出。
func ExportWorkHourStatistic(orgId int64, param vo.GetWorkHourStatisticReq) (*vo.ExportWorkHourStatisticResp, errs.SystemErrorInfo) {
	// 前端分页一般会进行分页，但导出任务是对筛选条件下的所有数据进行导出，因此要全量，并且是从第一页开始。这里直接用一个比较大的数，模糊这个边界。
	resizeCount := int64(100_0000)
	param.Size = &resizeCount
	firstPage := int64(1)
	param.Page = &firstPage
	// 参数校验和校准。
	if err := CheckAndFixGetStatisticDataParam(&param); err != nil {
		return nil, err
	}
	dataWrapped, err := GetWorkHourStatistic(orgId, param)
	if err != nil {
		return nil, err
	}
	relatePath, fileDir, fileName, err := GetExportFileInfo(orgId)
	url := config.GetOSSConfig().LocalDomain + relatePath + "/" + fileName
	// 日期字符串数组
	dateStrArr := GetDateListByDateRange(*param.StartTime, *param.EndTime)
	GenExcelFileForData(fileDir+"/"+fileName, dataWrapped, dateStrArr)
	return &vo.ExportWorkHourStatisticResp{
		URL: url,
	}, nil
}

// 生成 excel 文件
func GenExcelFileForData(fileWithPath string, workHourInfo *vo.GetWorkHourStatisticResp, dateStrArr []string) errs.SystemErrorInfo {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var err error

	dataList := workHourInfo.GroupStatisticList
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}
	// 渲染 excel 的头部
	row = sheet.AddRow()
	RenderExcelHeader(row, dateStrArr)
	// 渲染 excel 的主体数据
	RenderExcelBody(sheet, dataList)
	// 渲染总计数据
	RenderSummaryData(sheet, assemblySummaryData(dateStrArr, workHourInfo))
	// 文件保存
	err = file.Save(fileWithPath)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.FileSaveFail, err)
	}
	return nil
}

func assemblySummaryData(dateStrArr []string, workHourInfo *vo.GetWorkHourStatisticResp) *vo.OnePersonWorkHourStatisticInfo {
	dateWorkHourList := make([]*vo.OneDateWorkHour, len(dateStrArr))
	for i := 0; i < len(dateStrArr); i++ {
		tmpActualSum := float64(0)
		for _, item := range workHourInfo.GroupStatisticList {
			tmpNum, err := strconv.ParseFloat(item.DateWorkHourList[i].Time, 64)
			if err != nil {
				tmpNum = 0
			}
			tmpActualSum += tmpNum
		}
		// 换算成秒
		timeSeconds := tmpActualSum * 3600
		dateWorkHourList[i] = &vo.OneDateWorkHour{
			Date:    dateStrArr[i],
			WeekDay: 0, // 这个字段在后面的计算中没有用到，就先不处理。
			Time:    format.FormatNeedTimeIntoString(int64(timeSeconds)),
		}
	}
	return &vo.OnePersonWorkHourStatisticInfo{
		WorkerID:         0,
		Name:             "总计",
		PredictHourTotal: workHourInfo.Summary.PredictTotal,
		ActualHourTotal:  workHourInfo.Summary.ActualTotal,
		DateWorkHourList: dateWorkHourList,
	}
}

// 渲染总计数据
func RenderSummaryData(sheet *xlsx.Sheet, summaryData *vo.OnePersonWorkHourStatisticInfo) {
	var row *xlsx.Row
	var cell *xlsx.Cell
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = summaryData.Name
	cell = row.AddCell()
	cell.Value = summaryData.PredictHourTotal
	cell = row.AddCell()
	cell.Value = summaryData.ActualHourTotal
	// 装载日期内的工时数据
	for _, itemInOneDay := range summaryData.DateWorkHourList {
		cell = row.AddCell()
		cell.Value = itemInOneDay.Time
	}
}

func RenderExcelHeader(row *xlsx.Row, dateStrList []string) {
	var cell *xlsx.Cell
	cell = row.AddCell()
	cell.Value = "成员"
	cell = row.AddCell()
	cell.Value = "计划总工时"
	cell = row.AddCell()
	cell.Value = "实际总工时"
	for _, dateStr := range dateStrList {
		weekStr, monthlyDate, _ := GetWeekAndDate(dateStr)
		headerStr := fmt.Sprintf("%[1]s %[2]s", weekStr, monthlyDate)
		cell = row.AddCell()
		cell.Value = headerStr
	}
}

// 通过标准的日期字符串，获取对应的星期和日期。如将 `2020-12-04 00:00:00` 抓换为 `周一`、`11/29`。
func GetWeekAndDate(datetimeStr string) (weekStr, monthlyDate string, err errs.SystemErrorInfo) {
	t, _ := time.Parse(consts.AppTimeFormat, datetimeStr)
	weekNum := int(t.Weekday())
	switch weekNum {
	case 1:
		weekStr = "周一"
	case 2:
		weekStr = "周二"
	case 3:
		weekStr = "周三"
	case 4:
		weekStr = "周四"
	case 5:
		weekStr = "周五"
	case 6:
		weekStr = "周六"
	case 0:
		weekStr = "周日"
	}
	monthlyDate = t.Format("01/02")
	return
}

func RenderExcelBody(sheet *xlsx.Sheet, dataList []*vo.OnePersonWorkHourStatisticInfo) {
	var row *xlsx.Row
	var cell *xlsx.Cell
	for _, rowData := range dataList {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = rowData.Name
		// 装载 2 列汇总
		cell = row.AddCell()
		cell.Value = rowData.PredictHourTotal
		cell = row.AddCell()
		cell.Value = rowData.ActualHourTotal
		// 装载日期内的工时数据
		for _, itemInOneDay := range rowData.DateWorkHourList {
			cell = row.AddCell()
			cell.Value = itemInOneDay.Time
		}
	}
}

// 获取导出的文件相关信息
func GetExportFileInfo(orgId int64) (relatePath, fileDir, fileName string, err errs.SystemErrorInfo) {
	curMonthlyDate := time.Now().Format("200601")
	relatePath = "/work_hour/org_" + strconv.FormatInt(orgId, 10) + "/date_" + curMonthlyDate
	fileDir = strings.TrimRight(config.GetOSSConfig().RootPath, "/") + relatePath
	mkdirErr := os.MkdirAll(fileDir, 0777)
	if mkdirErr != nil {
		log.Error(mkdirErr)
		return "", "", "", errs.BuildSystemErrorInfo(errs.IssueDomainError, mkdirErr)
	}
	fileName = "工时统计.xlsx"
	return relatePath, fileDir, fileName, nil
}
