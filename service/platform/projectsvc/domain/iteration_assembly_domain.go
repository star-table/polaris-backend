package domain

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"upper.io/db.v3"
)

//statusType定义（1：未开始，2：进行中，3：已完成）
func IterationCondStatusAssembly(cond *db.Cond, orgId int64, statusType int) errs.SystemErrorInfo {
	var statusIds []int64 = nil

	if statusType == 1 {
		statusIds = append(statusIds, consts.StatusNotStart.ID)
	} else if statusType == 2 {
		statusIds = append(statusIds, consts.StatusRunning.ID)
	} else {
		statusIds = append(statusIds, consts.StatusComplete.ID)
	}
	(*cond)[consts.TcStatus] = db.In(statusIds)
	return nil
}
