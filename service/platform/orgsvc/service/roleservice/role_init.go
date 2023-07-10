package service

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	domain "github.com/star-table/polaris-backend/service/platform/orgsvc/domain/roledomain"
	"upper.io/db.v3/lib/sqlbuilder"
)

func RoleInit(orgId int64) (*bo.RoleInitResp, errs.SystemErrorInfo) {
	var resp *bo.RoleInitResp = nil
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		roleInitResp, err := domain.RoleInit(orgId, tx)
		resp = roleInitResp
		return err
	})
	if err != nil {
		log.Error(err)
		return nil, errs.SystemError
	}
	return resp, nil
}

func ChangeDefaultRole() errs.SystemErrorInfo {
	return domain.ChangeDefaultRole()
}
