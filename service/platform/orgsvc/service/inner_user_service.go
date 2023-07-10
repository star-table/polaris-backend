package service

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/businees"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
)

func InnerGetUserInfos(orgId int64, idValues interface{}) ([]*orgvo.InnerUserInfo, errs.SystemErrorInfo) {
	prefix, ids := businees.LcParseIds(idValues)
	log.Infof("[InnerGetUserInfos] LcParseIds, orgId: %v, input: %v, prefix: %v, ids: %v",
		orgId, json.ToJsonIgnoreError(idValues), prefix, ids)

	if prefix == consts.LcCustomFieldDeptType {
		resp := userfacade.GetUserIdsByDeptIds(&uservo.GetUserIdsByDeptIdsReq{
			OrgId:   orgId,
			DeptIds: ids,
		})
		if resp.Failure() {
			log.Errorf("[InnerGetUserInfos] GetUserIdsByDeptIds failed, orgId: %v, deptIds: %v, err: %v", orgId, ids, resp.Error())
			return nil, resp.Error()
		}
		ids = resp.Data.UserIds
	}

	userInfoBos, err := domain.BatchGetUserDetailInfoWithMobile(ids)
	if err != nil {
		log.Errorf("[InnerGetUserInfos] BatchGetUserDetailInfoWithMobile failed, ids: %v, err: %v", ids, err)
		return nil, err
	}

	var userInfos []*orgvo.InnerUserInfo
	for _, u := range userInfoBos {
		userInfos = append(userInfos, &orgvo.InnerUserInfo{
			Id:     u.ID,
			OrgID:  u.OrgID,
			Name:   u.Name,
			Mail:   u.Email,
			Phone:  u.Mobile,
			Avatar: u.Avatar,
		})
	}

	return userInfos, nil
}
