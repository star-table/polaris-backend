package projectvo

import (
	"fmt"

	tablev1 "github.com/star-table/interface/golang/table/v1"

	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
)

type SaveFormHeaderReq struct {
	OrgId         int64                 `json:"orgId"`
	UserId        int64                 `json:"userId"`
	SourceChannel string                `json:"sourceChannel"`
	InputAppId    int64                 `json:"inputAppId"`
	Params        vo.SaveFormHeaderData `json:"params"`
}

type SaveFormHeaderResp struct {
	vo.Err
	Data *vo.SaveFormHeaderRespData `json:"data"`
}

type GetFormConfigReq struct {
	OrgId               int64 `json:"orgId"`
	UserId              int64 `json:"userId"`
	ProjectId           int64 `json:"projectId"`
	ProjectObjectTypeId int64 `json:"projectObjectTypeId"`
}

type GetFormConfigResp struct {
	vo.Err
	Timestamp interface{}     `json:"timestamp"`
	Data      *FormConfigData `json:"data"`
}

type BizFormResp struct {
	vo.Err
	Timestamp interface{} `json:"timestamp"`
	Data      *BizForm    `json:"data"`
}

type BizForm struct {
	OrgId              int64                             `json:"orgId"`
	AppId              int64                             `json:"appId"`
	AppName            string                            `json:"appName"`
	FormId             int64                             `json:"formId"`
	ExtendsId          int64                             `json:"extendsId"`
	ExtendsFormId      int64                             `json:"extendsFormId"`
	ExtendsFieldParams map[string]lc_table.LcCommonField `json:"extendsFieldParams"`
	FieldParams        map[string]lc_table.LcCommonField `json:"fieldParams"`
}

func (form BizForm) GetTableName() string {
	if form.ExtendsFormId > 0 {
		return fmt.Sprintf("_form_%d_%d", form.OrgId, form.ExtendsFormId)
	}
	return fmt.Sprintf("_form_%d_%d", form.OrgId, form.FormId)
}

// 改成调用go-table的ReadTableSchemas，只返回 columns
type FormConfigData struct {
	AppId int64 `json:"appId"`
	//Fields       []map[string]interface{} `json:"fields"`
	Columns []map[string]interface{} `json:"columns"`
	//FieldOrders  interface{}              `json:"fieldOrders"`
	//CustomConfig map[string][]interface{} `json:"customConfig"`
	//BaseFields   []string                 `json:"baseFields"`
}

type FormConfigColumn struct {
	Children interface{}           `json:"children"`
	Editable bool                  `json:"editable"`
	Field    FormConfigColumnField `json:"field"`
	IsOrg    bool                  `json:"isOrg"`
	IsSys    bool                  `json:"isSys"`
	Key      string                `json:"key"`
	Render   interface{}           `json:"render"`
	Rules    interface{}           `json:"rules"`
	Title    string                `json:"title"`
	Unique   bool                  `json:"unique"`
	Writable bool                  `json:"writable"`
}

type FormConfigColumnField struct {
	AsyncData           interface{} `json:"asyncData"`
	CustomType          interface{} `json:"customType"`
	DataRely            interface{} `json:"dataRely"`
	DataType            interface{} `json:"dataType"`
	ProjectObjectTypeID interface{} `json:"projectObjectTypeId"`
	Props               interface{} `json:"props"` // 不同类型的自定义字段，它们的 props 属性值也不同
	Type                string      `json:"type"`
}

// 多选类型的 props 的类型
// 由于前端缘故，Select 中的 option 的 id 的类型是不固定的，可能是 int，也可能是 string。为了做适配，需要用 interface{} 表示。
type FormConfigColumnFieldMultiselectProps struct {
	Check       bool                                                           `json:"disabled"`
	IsCustom    bool                                                           `json:"isCustom"`
	Multiple    bool                                                           `json:"multiple"`
	Multiselect FormConfigColumnFieldMultiselectPropsMultiselectForInterfaceId `json:"multiselect"`
	Select      FormConfigColumnFieldMultiselectPropsMultiselectForInterfaceId `json:"select"` // 单选、多选大部分结构类似，因此复用一下。有 Select 时，没有 Multiselect
	IsSearch    bool                                                           `json:"isSearch"`
	Required    bool                                                           `json:"required"`
}

type FormConfigColumnFieldMultiselectPropsMultiselect struct {
	Options []FormConfigColumnFieldPropsMultiselectOption `json:"options"`
}

type FormConfigColumnFieldMultiselectPropsMultiselectForInterfaceId struct {
	Options []ColumnSelectOption `json:"options"`
}

type TableColumnFieldSelect struct {
	Select TableColumnFieldSelectOption `json:"select"`
}

type TableColumnFieldSelectOption struct {
	Options []FormConfigColumnFieldPropsMultiselectOption `json:"options"`
}

type FormConfigColumnFieldPropsMultiselectOption struct {
	Color string `json:"color"`
	Id    int64  `json:"id"`
	Value string `json:"value"`
}

type ColumnSelectOption struct {
	Id       interface{} `json:"id"`
	Value    string      `json:"value"`
	Color    string      `json:"color"`
	ParentId interface{} `json:"parentId"`
}

type FormConfigColumnFieldPropsOneVal struct {
	// 字段id
	FieldID string `json:"fieldId"`
	// 自定义字段中 select 类型的 option 的 id 值是 interface{} 类型
	FieldIdOfInterface interface{} `json:"fieldIdOfInterface"`
	// 字段值
	Value interface{} `json:"value"`
	// 名称
	Title string `json:"title"`
	Color string `json:"color"`
}

// inputnumber 数字的 props 类型
type FormConfigColumnFieldInputNumberProps struct {
	Checked     bool                                             `json:"checked"`
	Hide        bool                                             `json:"hide"`
	Inputnumber FormConfigColumnFieldInputNumberPropsInputNumber `json:"inputnumber"`
	IsCustom    bool                                             `json:"isCustom"`
	Multiple    bool                                             `json:"multiple"`
	Required    bool                                             `json:"required"`
}

type FormConfigColumnFieldInputNumberPropsInputNumber struct {
	Accuracy   string `json:"accuracy"`
	Percentage bool   `json:"percentage"`
	Thousandth bool   `json:"thousandth"`
	Unique     bool   `json:"unique"`
}

// 货币的 props 类型
type FormConfigColumnFieldAmountProps struct {
	Amount   FormConfigColumnFieldAmountPropsAmount `json:"amount"`
	Checked  bool                                   `json:"checked"`
	IsCustom bool                                   `json:"isCustom"`
	Multiple bool                                   `json:"multiple"`
	Required bool                                   `json:"required"`
}

type FormConfigColumnFieldAmountPropsAmount struct {
	Accuracy string `json:"accuracy"`
	Location string `json:"location"`
	Sign     string `json:"sign"`
	Unique   bool   `json:"unique"`
}

type FormConfigColumnFieldMemberProps struct {
	Member       FormConfigColumnFieldMemberPropsMember `json:"member"`
	Checked      bool                                   `json:"checked"`
	IsCustom     bool                                   `json:"isCustom"`
	Multiple     bool                                   `json:"multiple"`
	Required     bool                                   `json:"required"`
	IsUseUpdator bool                                   `json:"isUseUpdator"`
	IsUseCreator bool                                   `json:"isUseCreator"`
}

type FormConfigColumnFieldMemberPropsMember struct {
	Multiple     bool `json:"multiple"`
	IsUseUpdator bool `json:"isUseUpdator"`
	IsUseCreator bool `json:"isUseCreator"`
}

type FormConfigColumnFieldDeptProps struct {
	Dept     FormConfigColumnFieldDeptPropsDept `json:"dept"`
	Disabled bool                               `json:"disabled"`
	Multiple bool                               `json:"multiple"`
	Required bool                               `json:"required"`
}

type FormConfigColumnFieldDeptPropsDept struct {
	Multiple bool `json:"multiple"`
}

// 分组单选的 props 结构
type FormConfigColumnFieldGroupSelectProps struct {
	GroupSelect FormConfigColumnFieldGroupSelectPropsGroupSelect `json:"groupSelect"`
	IsCustom    bool                                             `json:"isCustom"`
	Multiple    bool                                             `json:"multiple"`
	Required    bool                                             `json:"required"`
}

type FormConfigColumnFieldGroupSelectPropsGroupSelect struct {
	GroupOptions []FormConfigColumnFieldGroupSelectPropsGroupSelectGroupOption `json:"groupOptions"`
	Options      []FormConfigColumnFieldGroupSelectPropsGroupSelectOption      `json:"options"`
}

type FormConfigColumnFieldGroupSelectPropsGroupSelectGroupOption struct {
	Id       string                                                             `json:"id"`
	Value    string                                                             `json:"value"`
	Children []FormConfigColumnFieldGroupSelectPropsGroupSelectGroupOptionChild `json:"children"`
}

type FormConfigColumnFieldGroupSelectPropsGroupSelectOption struct {
	Color    string `json:"color"`
	ID       int64  `json:"id"` // 改造后，option 的 id 为数值
	ParentID int    `json:"parentId"`
	Value    string `json:"value"`
}

type FormConfigColumnFieldGroupSelectPropsGroupSelectGroupOptionChild struct {
	Color    string `json:"color"`
	ID       string `json:"id"`
	ParentID string `json:"parentId"`
	Value    string `json:"value"`
}

type FormConfigColumnFieldDatepickerProps struct {
	Datepicker      FormConfigColumnFieldDatepickerPropsDatepicker `json:"datepicker"`
	Checked         bool                                           `json:"checked"`
	IsCustom        bool                                           `json:"isCustom"`
	Multiple        bool                                           `json:"multiple"`
	Required        bool                                           `json:"required"`
	Hide            bool                                           `json:"hide"`
	IsCreateSys     bool                                           `json:"isCreateSys"`
	IsUseUpdateTime bool                                           `json:"isUseUpdateTime"`
	IsUseCreateTime bool                                           `json:"isUseCreateTime"`
	IsUseUpdate     bool                                           `json:"isUseUpdate"`
	IsUseCreate     bool                                           `json:"isUseCreate"`
}

type FormConfigColumnFieldDatepickerPropsDatepicker struct {
	IsUseUpdateTime bool `json:"isUseUpdateTime"`
	IsUseCreateTime bool `json:"isUseCreateTime"`
}

type GetFormConfigBatchReq struct {
	OrgId  int64                  `json:"orgId"`
	UserId int64                  `json:"userId"`
	Data   GetFormConfigBatchData `json:"data"`
}

type GetFormConfigBatchData struct {
	AppIds []string `json:"appIds"`
}

type GetTableColumnBatchResp struct {
	vo.Err
	Data []tablev1.TableSchema `json:"data"`
}

type GetFormConfigBatchResp struct {
	vo.Err
	Data map[string]map[int64]FormConfigData `json:"data"`
}

type GetFormConfigBatchLcResp struct {
	vo.Err
	Timestamp interface{}      `json:"timestamp"`
	Data      []FormConfigData `json:"data"`
}

type TableColumnData struct {
	Name              string           `json:"name"`
	Label             string           `json:"label"`
	AliasTitle        string           `json:"aliasTitle"`
	Description       string           `json:"description"`
	IsSys             bool             `json:"isSys"`
	IsOrg             bool             `json:"isOrg"`
	Writable          bool             `json:"writable"`
	Editable          bool             `json:"editable"`
	Unique            bool             `json:"unique"`
	UniquePreHandler  string           `json:"uniquePreHandler"`
	SensitiveStrategy string           `json:"sensitiveStrategy"`
	SensitiveFlag     int32            `json:"sensitiveFlag"`
	Field             TableColumnField `json:"field"`
}

type TableColumnField struct {
	Type       string                 `json:"type,omitempty"`
	CustomType string                 `json:"customType"`
	DataType   string                 `json:"dataType,omitempty"`
	Props      map[string]interface{} `json:"props"`
	RefSetting map[string]interface{} `json:"refSetting"`
}

type TableColumnConfigInfo struct {
	TableId int64              `json:"tableId"`
	Columns []*TableColumnData `json:"columns"`
	//CustomConfig map[string][]string `json:"customConfig"`
}

//type TablesColumnsReq struct {
//	OrgId  int64 `json:"orgId"`
//	UserId int64 `json:"userId"`
//	Input  *tablev1.ReadTableSchemasRequest
//}

type TableColumnsResp struct {
	vo.Err
	Data *TableColumnsTable `json:"data"`
}

type TablesColumnsResp struct {
	vo.Err
	Data *TablesColumnsRespData `json:"data"`
}

type TablesColumnsRespData struct {
	Tables []*TableColumnsTable `json:"tables"`
}

type TableColumnsTable struct {
	AppId   int64              `json:"appId,string"`
	TableId int64              `json:"tableId,string"`
	Name    string             `json:"name"`
	Columns []*TableColumnData `json:"columns"`
}

type GetTableColumnReq struct {
	OrgId       int64 `json:"orgId"`
	UserId      int64 `json:"userId"`
	ProjectId   int64 `json:"projectId"`
	TableId     int64 `json:"tableId"`
	NotAllIssue int64 `json:"NotAllIssue"` // 如果为1，是指拿汇总表普通字段，不拿projects和迭代啥的
}

type GetTableColumnRespVo struct {
	vo.Err
	Data *TableColumnConfigInfo `json:"data"`
}

type GetTablesColumnsReq struct {
	OrgId  int64               `json:"orgId"`
	UserId int64               `json:"userId"`
	Input  *TablesColumnsInput `json:"input"`
}

type TablesColumnsInput struct {
	TableIds  []int64  `json:"tableIds"`
	ColumnIds []string `json:"columnIds"`
}

// 更新任务时需要用到的一些信息
type UpdateIssueHelperInfo struct {
	TableColumns TablesColumnsRespData             `json:"tableColumns"`
	ColumnsMap   map[string]lc_table.LcCommonField `json:"columnsMap"`
}
