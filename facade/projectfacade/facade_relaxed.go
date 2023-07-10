package projectfacade

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
)

func GetProjectBoListByProjectTypeLangCodeRelaxed(orgId int64, projectTypeLangCode *string) ([]bo.ProjectBo, errs.SystemErrorInfo) {
	respVo := GetProjectBoListByProjectTypeLangCode(projectvo.GetProjectBoListByProjectTypeLangCodeReqVo{
		OrgId:               orgId,
		ProjectTypeLangCode: projectTypeLangCode,
	})

	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.ProjectBoList, nil
}
