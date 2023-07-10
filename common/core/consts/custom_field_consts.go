package consts

const (
	//类型(1文本类型2单选框3多选框4日期选框5人员选择6是非选择7数字框)
	CustomTypeText        = 1
	CustomTypeRadio       = 2
	CustomTypeCheckBox    = 3
	CustomTypeDate        = 4
	CustomTypePersonnel   = 5
	CustomTypeTrueOrFalse = 6
	CustomTypeNumber      = 7

	// 自定义字段是非选择值枚举
	CustomTypeBoolTrueVal = 1
	CustomTypeBoolFalseVal = 2
)

const (
	ProjectViewTypeBoard = 1 //看板
	ProjectViewTypeList  = 2 //列表
	ProjectViewTypeTable = 3 //表格
)

//项目默认定义的展示字段
const (
	ProjectFieldTitle           = "title"           //任务名称
	ProjectFieldStatus          = "status"          //任务状态
	ProjectFieldOwner           = "owner"           //负责人
	ProjectFieldCreateTime      = "createTime"      //创建时间
	ProjectFieldPriority        = "priority"        //优先级
	ProjectFieldPlanStartTime   = "planStartTime"   //开始时间
	ProjectFieldPlanEndTime     = "planEndTime"     //截止时间
	ProjectFieldTag             = "tag"             //标签
	ProjectFieldPlanHour        = "planHour"        //计划工时
	ProjectFieldActualHour      = "actualHour"      //实际工时
	ProjectFieldRelatedIssueNum = "relatedIssueNum" //关联任务数
	ProjectFieldCommentNum      = "commentNum"      //评论数
)

const (
	CustomFieldScopeOrg = 1
	CustomFieldScopePro = 2
	CustomFieldScopeSys = 3
)

var (
	CustomFieldTrueOrFalseValMap = map[string]int64{
		"是": 1,
		"否": 2,
	}
)


const DefaultSelectFontColor = "#3c4152"

//需求类型
var DemandTypeMap = map[string]int64{
	"新功能":1,
	"需求优化":2,
	"交互优化":3,
	"视觉优化":4,
	"其它":5,
}

//需求来源
var DemandSourceMap = map[string]int64{
	"产品经理":1,
	"客户":2,
	"老板":3,
	"市场需求":4,
	"活动需求":5,
	"其它":6,
}

//严重程度
var BugPropertyMap = map[string]int64{
	"建议":1,
	"一般":2,
	"严重":3,
	"致命":4,
}

//缺陷类型
var BugTypeMap = map[string]int64{
	"功能问题":1,
	"性能问题":2,
	"接口问题":3,
	"界面问题":4,
	"兼容问题":5,
	"易用性问题":6,
	"逻辑问题":7,
	"需求问题":8,
	"不确定":9,
}