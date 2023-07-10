package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type ProjectInitReqVo struct {
	OrgId int64 `json:"orgId"`
}

type ProjectInitRespVo struct {
	ContextMap map[string]interface{} `json:"data"`

	vo.Err
}