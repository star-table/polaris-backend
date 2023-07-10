package datacenter

import "github.com/star-table/polaris-backend/common/model/vo"

type QueryReq struct {
	From      []Table          `json:"from"`
	Condition vo.LessCondsData `json:"condition"`
	Limit     int              `json:"limit"`
	Offset    int              `json:"offset"`
	Columns   []string         `json:"columns"`
	Orders    []Order          `json:"orders"`
	Groups    []string         `json:"groups"`
}

type Table struct {
	Type     string `json:"type"`
	Schema   string `json:"schema"`
	Database string `json:"database"`
}

type Order struct {
	Column string `json:"column"`
	IsAsc  bool   `json:"is_asc"`
}

type QueryResp struct {
	vo.Err
	Data []map[string]interface{} `json:"data"`
}

type Set struct {
	Column          string      `json:"column"`
	Value           interface{} `json:"value"`
	Type            int         `json:"type"`
	Action          int         `json:"action"`
	WithoutPretreat bool        `json:"withoutPretreat"`
}
