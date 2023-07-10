package service

import (
	"github.com/star-table/polaris-backend/common/core/errs"
)

func AuthOrgRole(orgId, userId int64, path string, operation string) errs.SystemErrorInfo {
	return Authenticate(orgId, userId, nil, nil, path, operation, nil)
}
