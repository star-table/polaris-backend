package orgfacade

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
)

func GetBaseUserInfoBatchRelaxed(orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	respVo := GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: userIds,
	})
	return respVo.BaseUserInfos, respVo.Error()
}
