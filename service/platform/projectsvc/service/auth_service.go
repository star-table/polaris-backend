package service

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

func AuthProjectPermission(orgId, userId, projectId int64, path string, operation string, authFiling bool) errs.SystemErrorInfo {
	return domain.AuthProjectWithCond(orgId, userId, projectId, path, operation, authFiling, false, 0)
}
