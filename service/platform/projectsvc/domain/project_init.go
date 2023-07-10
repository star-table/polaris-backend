package domain

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"upper.io/db.v3/lib/sqlbuilder"
)

//优先级初始化
func PriorityInit(orgId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	//log.Infof("start:优先级初始化")
	//err := InitPriority(orgId, tx)
	//if err != nil {
	//	return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	//}
	//log.Infof("success:优先级初始化")
	return nil
}
