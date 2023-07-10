package lc_table

import (
	"github.com/spf13/cast"
	"google.golang.org/protobuf/types/known/structpb"
)

//无码中的字段组件封装，定义的结构体
//参考文档 https://wiki.bjx.cloud/pages/viewpage.action?pageId=32145423
//组件都以 `Ct` 作为前缀，表示 Column Type。

type LcCtInput struct {
	Name     string         `json:"name"`
	Label    string         `json:"label"`
	Unique   bool           `json:"unique"`
	Editable bool           `json:"editable"`
	Writable bool           `json:"writable"`
	EnLabel  string         `json:"enLabel"`
	Field    LcCtInputField `json:"field"`
}
type LcCtInputField struct {
	Type  string              `json:"type"`
	Props LcCtInputFieldProps `json:"props"`
}
type LcCtInputFieldProps struct {
	Required    bool                           `json:"required"`
	IsSearch    bool                           `json:"isSearch"`
	Disabled    bool                           `json:"disabled"`
	Hide        bool                           `json:"hide"`
	FieldSearch LcCtInputFieldPropsFieldSearch `json:"fieldSearch"`
	PushMsg     bool                           `json:"pushMsg"`
}
type LcCtInputFieldPropsFieldSearch struct {
	Type string `json:"type"`
	Sort int    `json:"sort"`
}

type LcCtTextarea struct {
	Name     string            `json:"name"`
	Label    string            `json:"label"`
	EnLabel  string            `json:"enLabel"`
	Editable bool              `json:"editable"`
	Writable bool              `json:"writable"`
	Field    LcCtTextareaField `json:"field"`
}

type LcCtTextareaField struct {
	Type  string                 `json:"type"`
	Props LcCtTextareaFieldProps `json:"props"`
}

type LcCtTextareaFieldProps struct {
	Hide    bool `json:"hide"`
	PushMsg bool `json:"pushMsg"`
}

// select 适用于 select groupSelect multiSelect 等，如果缺少 field 可以补充。
type LcCtSelect struct {
	Name     string          `json:"name"`
	Label    string          `json:"label"`
	EnLabel  string          `json:"enLabel"`
	Editable bool            `json:"editable"`
	Writable bool            `json:"writable"`
	Field    LcCtSelectField `json:"field"`
}
type LcCtSelectField struct {
	Type       string               `json:"type"`
	DataType   string               `json:"dataType"`
	CustomType string               `json:"customType"`
	Props      LcCtSelectFieldProps `json:"props"`
}
type LcCtSelectFieldProps struct {
	IsText        bool                       `json:"isText"`
	Disabled      bool                       `json:"disabled"`
	Required      bool                       `json:"required"`
	TitleDisabled bool                       `json:"titleDisabled"`
	TypeDisabled  bool                       `json:"typeDisabled"`
	Default       interface{}                `json:"default"`
	Select        LcCtSelectFieldPropsSelect `json:"select"`
}
type LcCtSelectFieldPropsSelect struct {
	Options []LcCtSelectFieldPropsSelectOptions `json:"options"`
}
type LcCtSelectFieldPropsSelectOptions struct {
	Id    interface{} `json:"id"`
	Color string      `json:"color"`
	Value string      `json:"value"`
}

type LcCtDatepicker struct {
	Name     string              `json:"name"`
	Label    string              `json:"label"`
	EnLabel  string              `json:"enLabel"`
	Editable bool                `json:"editable"`
	Writable bool                `json:"writable"`
	Field    LcCtDatepickerField `json:"field"`
}
type LcCtDatepickerField struct {
	Type  string                   `json:"type"`
	Props LcCtDatepickerFieldProps `json:"props"`
}

type LcCtDatepickerFieldProps struct {
	PushMsg bool `json:"pushMsg"`
}

type LcCtMember struct {
	Name     string          `json:"name"`
	Label    string          `json:"label"`
	EnLabel  string          `json:"enLabel"`
	Editable bool            `json:"editable"`
	Writable bool            `json:"writable"`
	Field    LcCtMemberField `json:"field"`
}
type LcCtMemberField struct {
	Type  string               `json:"type"`
	Props LcCtMemberFieldProps `json:"props"`
}
type LcCtMemberFieldProps struct {
	Multiple          bool         `json:"multiple"`
	Limit             *int         `json:"limit,omitempty"`
	CollaboratorRoles *[]string    `json:"collaboratorRoles"`
	Member            LcPropMember `json:"member"`
	PushMsg           bool         `json:"pushMsg"`
}

type LcCtRelateTable struct {
	Name  string               `json:"name"`
	Label string               `json:"label"`
	Field LcCtRelateTableField `json:"field"`
}
type LcCtRelateTableField struct {
	Type  string                    `json:"type"`
	Props LcCtRelateTableFieldProps `json:"props"`
}
type LcCtRelateTableFieldProps struct {
	RelateTable LcCtRelateTableFieldPropsRelateTable `json:"relateTable"`
	ShowDetails bool                                 `json:"showDetails"`
	FormOrder   int                                  `json:"formOrder"`
	Required    bool                                 `json:"required"`
	TabParam    string                               `json:"tabParam"`
	GroupParam  string                               `json:"groupParam"`
	IsSearch    bool                                 `json:"isSearch"`
	FieldSearch LcCtRelateTableFieldPropsFieldSearch `json:"fieldSearch"`
}
type LcCtRelateTableFieldPropsRelateTable struct {
	AppID          string   `json:"appId"`
	IDField        string   `json:"id_field"`
	DisplayField   string   `json:"display_field"`
	DisplayColumns []string `json:"display_columns"`
}
type LcCtRelateTableFieldPropsFieldSearch struct {
	Type string `json:"type"`
	Sort int    `json:"sort"`
}

type LcFormConfig struct {
	Fields      []LcCommonField `json:"fields"`
	FieldOrders []string        `json:"fieldOrders"`
	ViewOrders  []string        `json:"viewOrders"`
	BaseFields  []string        `json:"baseFields"`
}

type LcCommonField struct {
	Key        string      `json:"key"`
	Name       string      `json:"name"`       // 列的 key
	AliasTitle string      `json:"aliasTitle"` // 列的别名
	Title      string      `json:"title"`
	Label      string      `json:"label"` // 列名称
	EnLabel    string      `json:"enLabel"`
	EnTitle    string      `json:"enTitle"`
	Editable   *bool       `json:"editable,omitempty"`
	Writable   bool        `json:"writable"`
	Field      LcFieldData `json:"field"`
	IsOrg      bool        `json:"isOrg"`
}

type LcOneColumn struct {
	Name       string      `json:"name"`
	Label      string      `json:"label"`
	AliasTitle string      `json:"aliasTitle"`
	EnTitle    string      `json:"enTitle"`
	Field      LcFieldData `json:"field"`
	IsOrg      bool        `json:"isOrg"`
	IsSys      bool        `json:"isSys"`
	Key        string      `json:"key"`   // 字段标识，唯一标识
	Title      string      `json:"title"` // 字段名称
	Unique     bool        `json:"unique"`
	Writable   bool        `json:"writable"`
	Editable   bool        `json:"editable"`
}

type LcFieldData struct {
	Type                string  `json:"type"`
	CustomType          string  `json:"customType"`
	ProjectObjectTypeId int64   `json:"projectObjectTypeId"`
	Props               LcProps `json:"props"`
}

type LcProps struct {
	IsText            bool               `json:"isText"`
	Required          bool               `json:"required"`
	Multiple          bool               `json:"multiple"`
	Member            LcPropMember       `json:"member"`
	Select            LcPropSelect       `json:"select"`
	MultiSelect       LcPropSelect       `json:"multiselect"`
	GroupSelect       *LcPropGroupSelect `json:"groupSelect"`
	InputNumber       LcPropInputNumber  `json:"inputnumber"`
	CollaboratorRoles []string           `json:"collaboratorRoles"`
	PushMsg           bool               `json:"pushMsg"`
	Default           interface{}        `json:"default"`
}

type LcPropGroupSelect struct {
	GroupOptions []LcGroupOptions       `json:"groupOptions"`
	Options      []LcGroupOptionsDetail `json:"options"`
}

type LcGroupOptions struct {
	Id        int                    `json:"id"`
	Value     string                 `json:"value"`
	Children  []LcGroupOptionsDetail `json:"children"`
	Color     string                 `json:"color"`
	FontColor string                 `json:"fontColor"`
}

type LcGroupOptionsDetail struct {
	Color     string `json:"color"`
	FontColor string `json:"fontColor"`
	Id        int    `json:"id"`
	Value     string `json:"value"`
	Sort      int    `json:"sort"`
	ParentId  int    `json:"parentId"`
	TableId   string `json:"tableId,omitempty"`
}

type LcPropMember struct {
	Multiple        bool `json:"multiple"`
	Required        bool `json:"required"`
	IsCollaborators bool `json:"isCollaborators"`
}

type LcPropSelect struct {
	Options []LcOptions `json:"options"`
}

type LcOptions struct {
	Color     string      `json:"color"`
	FontColor string      `json:"fontColor"`
	Id        interface{} `json:"id"`
	Value     string      `json:"value"`
	Sort      *int        `json:"sort,omitempty"`
	Key       string      `json:"key,omitempty"`
}

type LcProjectOptions struct {
	LcOptions
	Children []*LcOptions `json:"children"`
}

type LcPropInputNumber struct {
	Accuracy   string `json:"accuracy"`
	Thousandth bool   `json:"thousandth"`
	Percentage bool   `json:"percentage"`
	Unique     bool   `json:"unique"`
	Required   bool   `json:"required"`
}

type LcDocumentValue struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Suffix string `json:"suffix"`
}

type PolarisPersonFieldValue struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

func ExchangeToLcGroupSelectOptions(props map[string]interface{}, tableId string) []*LcGroupOptionsDetail {
	options := make([]*LcGroupOptionsDetail, 0, 3)
	if m, ok := props["groupSelect"].(map[string]interface{}); ok {
		selectStruct, err := structpb.NewStruct(m)
		if err != nil {
			return options
		}
		if selectStruct.Fields["options"] != nil && selectStruct.Fields["options"].GetListValue() != nil {
			for _, option := range selectStruct.Fields["options"].GetListValue().Values {
				if option.GetStructValue() != nil && option.GetStructValue().Fields != nil {
					fields := option.GetStructValue().Fields
					options = append(options, &LcGroupOptionsDetail{
						Color:     cast.ToString(fields["color"].GetStringValue()),
						FontColor: cast.ToString(fields["fontColor"].GetStringValue()),
						Id:        cast.ToInt(fields["id"].GetNumberValue()),
						Value:     cast.ToString(fields["value"].GetStringValue()),
						ParentId:  cast.ToInt(fields["parentId"].GetNumberValue()),
						TableId:   tableId,
					})
				}
			}
		}
	}

	return options
}
