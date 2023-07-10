package projectvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type IterationStatsReqVo struct {
	Page  uint                `json:"page"`
	Size  uint                `json:"size"`
	Input vo.IterationStatReq `json:"input"`
	OrgId int64               `json:"orgId"`
}

type IterationStatsRespVo struct {
	vo.Err
	IterationStats *vo.IterationStatList `json:"data"`
}

type AppendIterationStatReqVo struct {
	IterationBo bo.IterationBo `json:"iterationBo"`
	Date        string         `json:"date"`
}
