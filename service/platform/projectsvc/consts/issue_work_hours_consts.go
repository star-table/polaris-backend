package consts

// 剩余工时计算方式 1动态计算；2手动填写 默认1
const RemainTimeCalTypeAuto = 1
const RemainTimeCalTypeManual = 2
const RemainTimeCalTypeDefault = RemainTimeCalTypeAuto

const (
	// 项目详情中，是否关闭工时功能
	ProjectDetailEnableWorkHoursOn  = 1
	ProjectDetailEnableWorkHoursOff = 2

	// 工时类型:1简单版预估工时（总预估工时），2实际工时，3详细版预估工时（子预估工时）
	WorkHourTypeTotalPredict = 1
	WorkHourTypeActual       = 2
	WorkHourTypeSubPredict   = 3

	// 允许用户填写的最大工时 100w 小时
	MaxWorkHour = 3600000000
)
