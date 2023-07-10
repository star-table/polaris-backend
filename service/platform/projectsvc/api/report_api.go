package api

import (
	msgPb "github.com/star-table/interface/golang/msg/v1"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/facade/common/report"
)

func (PostGreeter) ReportAppEvent(req *commonvo.ReportAppEventReq) *vo.CommonRespVo {
	report.ReportAppEvent(msgPb.EventType(req.EventType), req.TraceId, req.AppEvent)
	return &vo.CommonRespVo{}
}

func (PostGreeter) ReportTableEvent(req *commonvo.ReportTableEventReq) *vo.CommonRespVo {
	report.ReportTableEvent(msgPb.EventType(req.EventType), req.TraceId, req.TableEvent)
	return &vo.CommonRespVo{}
}
