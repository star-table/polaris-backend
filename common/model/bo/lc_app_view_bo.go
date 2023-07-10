package bo

type AppViewConfig struct {
	Orders              []AppViewOrder `json:"orders"`
	Condition           interface{}    `json:"condition"`
	ProjectObjectTypeId interface{}    `json:"projectObjectTypeId"`
	RealCondition       interface{}    `json:"realCondition"`
	Group               interface{}    `json:"group"`
	OrderParams         interface{}    `json:"orderParams"`
	HiddenColumnIds     []string       `json:"hiddenColumnIds"`
	TableId             string         `json:"tableId"`
}

type LcAppView struct {
	ViewName string        `json:"viewName"`
	Type     int           `json:"type"`
	Config   AppViewConfig `json:"config"`
}

type LcAppViewJson struct {
	ViewName string `json:"viewName"`
	Type     int    `json:"type"`
	Config   string `json:"config"`
}

type AppViewOrder struct {
	Asc    bool   `json:"asc"`
	Column string `json:"column"`
}

type AppViewCondition struct {
	Column    int64           `json:"column"`
	FieldType string          `json:"fieldType"`
	Option    []AppViewOption `json:"option"`
	Type      string          `json:"type"`
	Values    []string        `json:"values"`
}

type AppViewOption struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

type AppViewParam struct {
	Conds []AppViewParamConds `json:"conds"`
}

type AppViewParamConds struct {
	Type   string `json:"type"`
	Value  string `json:"value"`
	Column int64  `json:"column"`
}

type AppViewLessShowCondition struct {
	Column    string                     `json:"column"`
	FieldType string                     `json:"fieldType"`
	Option    []map[string]interface{}   `json:"option"`
	Type      string                     `json:"type"`
	Values    []interface{}              `json:"values"`
	Value     interface{}                `json:"value"`
	Props     AppViewLessProp            `json:"props"`
	Conds     []AppViewLessShowCondition `json:"conds"`
}

type AppViewLessProp struct {
	DatePicker AppViewLessPropSelectType `json:"datePicker"`
}

type AppViewLessPropSelectType struct {
	SelectType string `json:"selectType"`
}

type AppViewLessCondition struct {
	Column    string                 `json:"column"`
	FieldType string                 `json:"fieldType"`
	Type      string                 `json:"type"`
	Value     interface{}            `json:"value"`
	Values    []interface{}          `json:"values"`
	Conds     []AppViewLessCondition `json:"conds"`
}

type MemberConditionValues struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type AppViewMultiConditionForProject struct {
	Type  string                 `json:"type"`
	Conds []AppViewLessCondition `json:"conds"`
}
