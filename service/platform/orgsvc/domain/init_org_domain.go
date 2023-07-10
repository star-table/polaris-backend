package domain

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3/lib/sqlbuilder"
)

//留着做对比方便的注释
func InitOrg(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (*po.PpmOrgOrganization, errs.SystemErrorInfo) {
	return nil, nil
}
