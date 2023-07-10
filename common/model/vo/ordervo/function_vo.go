package ordervo

import "github.com/star-table/polaris-backend/common/model/vo"

type FunctionReq struct {
	Level int64 `json:"level"`
}

type FunctionResp struct {
	vo.Err
	Data []FunctionLimitObj `json:"data"`
}

type FunctionLimitObj struct {
	Name     string              `json:"name"`
	Key      string              `json:"key"`
	HasLimit bool                `json:"hasLimit"`
	Limit    []FunctionLimitItem `json:"limit"`
}

type FunctionLimitItem struct {
	Typ  string `json:"typ"`
	Num  int    `json:"num"`
	Unit string `json:"unit"`
}
