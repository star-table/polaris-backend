package consts

var (
	// 任务状态
	LcIssueStatusMap = map[string]string{
		"todo":     "待处理",
		"doing":    "进行中",
		"finished": "已完成",
	}
	// 任务优先级
	LcIssuePriorityMap = map[string]string{
		"p1": "最低",
		"p2": "较低",
		"p3": "普通",
		"p4": "较高",
		"p5": "最高",
	}
)

const (
	// 标签表的名字字段名
	//LcTagTableColumnNameOfName = "_field_polaris_tag_name"
	//LcTagTableColumnNameOfId   = "id"

	//LcBarTableColumnNameOfName      = "_field_polaris_bar_name"
	//LcPriorityTableColumnNameOfName = "_field_polaris_priority_name"
	//LcSceneTableColumnNameOfName    = "_field_polaris_scene_name"

	// 字段类型
	// 无码表格自定义字段类型标识：部门
	// 无码表格自定义字段类型标识：成员
	LcColumnFieldTypeSelect         = "select"         // 单选
	LcColumnFieldTypeMultiSelect    = "multiselect"    // 多选
	LcColumnFieldTypeGroupSelect    = "groupSelect"    // 分组单选
	LcColumnFieldTypeTextarea       = "textarea"       // 多行文本
	LcColumnFieldTypeInput          = "input"          // 单行文本
	LcColumnFieldTypeRichText       = "richtext"       // 富文本
	LcColumnFieldTypeMember         = "member"         // 成员
	LcColumnFieldTypeDept           = "dept"           // 部门
	LcColumnFieldTypeDatepicker     = "datepicker"     // 时间
	LcColumnFieldTypeAmount         = "amount"         // 金额
	LcColumnFieldTypeInputNumber    = "inputnumber"    // 数字
	LcColumnFieldTypeMobile         = "mobile"         // 手机号
	LcColumnFieldTypeLink           = "link"           // 链接
	LcColumnFieldTypeEmail          = "email"          // 邮箱
	LcColumnFieldTypeDocument       = "document"       // 附件
	LcColumnFieldTypeImage          = "image"          // 图片
	LcColumnFieldTypeWorkHour       = "workHour"       // 工时
	LcColumnFieldTypeRelating       = "relating"       // 关联
	LcColumnFieldTypeBaRelating     = "baRelating"     // 前后置
	LcColumnFieldTypeConditionRef   = "conditionRef"   // 条件引用
	LcColumnFieldTypeFormula        = "formula"        // 公式
	LcColumnFieldTypeSingleRelating = "singleRelating" // 单向关联

	// 机器人个人推送 - 描述任务字段值变化类型
	TrendChangeListBoChangeTypeDescAdd    = "新增"
	TrendChangeListBoChangeTypeDescUpload = "上传"
	TrendChangeListBoChangeTypeDescDelete = "删除"
	TrendChangeListBoChangeTypeDescUpdate = "更新"
)

var (
	// 机器人个人推送-这些类型的列发生变化，需要组装详情展示
	FsRobotPersonIssueUpdateContentColumnsForText = []string{
		LcColumnFieldTypeInput, LcColumnFieldTypeTextarea, LcColumnFieldTypeTextarea,
		LcColumnFieldTypeDatepicker, LcColumnFieldTypeSelect, LcColumnFieldTypeMultiSelect, LcColumnFieldTypeGroupSelect,
		LcColumnFieldTypeAmount, LcColumnFieldTypeInputNumber,
		LcColumnFieldTypeMobile, LcColumnFieldTypeEmail, LcColumnFieldTypeLink,
		LcColumnFieldTypeMember, LcColumnFieldTypeDept,
	}
	// 附件、图片等
	FsRobotPersonIssueUpdateContentColumnsForFile = []string{
		LcColumnFieldTypeDocument, LcColumnFieldTypeImage,
	}
	// 前后置、关联
	FsRobotPersonIssueUpdateContentColumnsForRelate = []string{
		LcColumnFieldTypeBaRelating, LcColumnFieldTypeRelating,
	}
)

const (
	// MaxTowTableColumns 两个表相乘的数量要小于5万，要不就不拿数据了
	MaxTowTableColumns = 50000
)
