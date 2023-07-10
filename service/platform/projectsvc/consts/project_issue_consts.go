package consts

import (
	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
)

const (
	// 一些内置字段的 id。
	FilterTitleId         = -1
	FilterNumberId        = -2
	FilterOwnerId         = -3  // 负责人
	FilterStatusId        = -4  //状态，
	FilterPlanEndTimeId   = -6  //截止时间
	FilterPriorityId      = -7  // 优先级
	FilterFollowerId      = -8  // 关注人
	FilterPlanStartTimeId = -9  // 开始时间
	FilterIterationId     = -10 // 迭代
	FilterTaskColumnId    = -11 // 任务栏
	FilterDemandTypeId    = -12 // 需求类型
	FilterDemandSourceId  = -13 // 需求来源
	FilterBugTypeId       = -14 // 缺陷类型
	FilterSeverityId      = -15 // 严重程度
	FilterTagId           = -16 // 标签
	FilterAuditorId       = -26 //确认人
	FilterIsOverdue       = -29 // 逾期任务
	FilterStatusType      = -31 // 状态类型的筛选
)

var (
	// 是否是逾期任务的筛选值映射
	IssueIsOverdueMap = map[string]int8{
		"已逾期": 1,
		"未逾期": 2,
	}

	// 通用项目下，任务如果设定了**确认人**，则此时有“待确认”状态。因此筛选中追加一项“待确认”的筛选。
	IssueCommonProExtraStatus = map[int]string{
		-1: "待确认",
	}

	// 导入暂不支持的自定义字段类型列表
	NotSupportImportCustomFieldTypes = []string{
		consts.LcColumnFieldTypeImage,    // 图片
		consts.LcColumnFieldTypeRichText, // 富文本
		consts.LcColumnFieldTypeDocument, // 附件

		consts.LcColumnFieldTypeWorkHour,
		consts.LcColumnFieldTypeRelating,
		consts.LcColumnFieldTypeBaRelating,
		consts.LcColumnFieldTypeConditionRef,
		consts.LcColumnFieldTypeFormula, // 公式
		consts.LcColumnFieldTypeSingleRelating,
		tablePb.ColumnType_reference.String(),
	}

	NotSupportExportCustomFieldTypes = []string{
		//consts.LcColumnFieldTypeRichText, // 富文本
		consts.LcColumnFieldTypeFormula, // 公式
	}

	//NotSupportImportCustomFieldTypesOfInt32 = []tableV1.ColumnType{
	//	tableV1.ColumnType_image,    // 图片
	//	tableV1.ColumnType_richtext, // 富文本
	//	tableV1.ColumnType_document, // 附件
	//
	//	tableV1.ColumnType_workHour,   // 工时
	//	tableV1.ColumnType_relating,   // 关联
	//	tableV1.ColumnType_baRelating, // 前后置
	//}

	// 下载模板时，无需显示的字段（系统生成值的字段）
	HiddenColumnsInImportTpl = []string{
		consts.BasicFieldCreator,
		consts.BasicFieldUpdator,
		"createTime",
		"updateTime",
		consts.BasicFieldCode,
	}

	// 敏捷项目任务的状态类型筛选支持以下值
	AgilityProIssueStatusType = map[int]string{
		-1: "待确认",
		1:  "未开始",
		2:  "已完成",
		3:  "未开始",
		4:  "进行中",
		5:  "已逾期",
	}

	ExportDownloadTypesMap = map[string]struct{}{
		"gif":  {},
		"png":  {},
		"jpg":  {},
		"jpeg": {},
		"webp": {},
	}

	ExportFileTypeIconMap = map[string]string{
		"gif":       "gif.png",
		"tiff":      "gif.png",
		"html":      "html.png",
		"jpg":       "jpg.png",
		"webp":      "jpg.png",
		"png":       "png.png",
		"mp3":       "mp3.png",
		"pdf":       "pdf.png",
		"ppt":       "ppt.png",
		"pptx":      "ppt.png",
		"psd":       "psd.png",
		"txt":       "txt.png",
		"mp4":       "video.png",
		"avi":       "video.png",
		"mpg":       "video.png",
		"rmv":       "video.png",
		"mov":       "video.png",
		"zip":       "zip.png",
		"rar":       "zip.png",
		"xls":       "excel.png",
		"xlsx":      "excel.png",
		"doc":       "word.png",
		"docx":      "word.png",
		"csv":       "word.png",
		"docs":      "file_doc_nor.png",
		"sheets":    "file_sheet_nor.png",
		"mindnotes": "file_mindnote_nor.png",
		"mindnote":  "file_mindnote_nor.png",
		"bitable":   "file_bitable_nor.png",
		"base":      "file_bitable_nor.png",
	}

	ExportOfficeFileType = map[string]struct{}{
		"xls":  {},
		"xlsx": {},
		"doc":  {},
		"docx": {},
		"csv":  {},
		"ppt":  {},
		"pptx": {},
	}

	OfficeViewUrl = "https://view.officeapps.live.com/op/view.aspx?src="
	PdfVideUrl    = "https://pdfviewer.startable.cn/?file="

	ParentDescColumn = &projectvo.TableColumnData{
		Name:  consts.HelperFieldIsParentDesc,
		Label: "父子类型",
	}

	BaRelatingExportTemplate = "前置%d条/后置%d条"
	RelatingExportTemplate   = "已关联%d个任务"
	WorkHourExportTemplate   = "计划工时%vh/实际工时%vh"

	ExcelJoinFlag = "||"
)
