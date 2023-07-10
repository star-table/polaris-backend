package service

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	service "github.com/star-table/polaris-backend/service/platform/orgsvc/service/roleservice"
)

func AuthOrgRole(orgId, userId int64, path string, operation string) errs.SystemErrorInfo {
	return service.Authenticate(orgId, userId, nil, nil, path, operation, nil)
}
