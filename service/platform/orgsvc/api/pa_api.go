package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) PAReport(req orgvo.PAReportMsgReqVo) vo.VoidErr {
	err := service.PAReport(req.Body)
	return vo.VoidErr{Err: vo.NewErr(err)}
}
