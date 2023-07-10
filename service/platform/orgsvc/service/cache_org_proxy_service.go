package service

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
)

func GetBaseOrgInfo(orgId int64) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseOrgInfo(orgId)
}

func GetBaseUserInfoByEmpId(orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return &bo.BaseUserInfoBo{}, nil
}

func GetBaseUserInfoByEmpIdBatch(orgId int64, input orgvo.GetBaseUserInfoByEmpIdBatchReqVoInput) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseUserInfoByOpenIdBatch(orgId, input.OpenIds)
}

func GetUserConfigInfo(orgId int64, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	return domain.GetUserConfigInfo(orgId, userId)
}

func GetUserConfigInfoBatch(orgId int64, input *orgvo.GetUserConfigInfoBatchReqVoInput) ([]bo.UserConfigBo, errs.SystemErrorInfo) {
	return domain.GetUserConfigInfoBatch(orgId, input.UserIds)
}

func GetBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseUserInfo(orgId, userId)
}

func GetDingTalkBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return &bo.BaseUserInfoBo{}, nil
}

func GetBaseUserInfoBatch(orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseUserInfoBatch(orgId, userIds)
}

func SetShareUrl(key, url string) errs.SystemErrorInfo {
	return domain.SetShareUrl(key, url)
}

func GetShareUrl(key string) (string, errs.SystemErrorInfo) {
	return domain.GetShareUrl(key)
}

func ClearOrgUsersPayCache(orgId int64) errs.SystemErrorInfo {
	return domain.ClearAllOrgUserPayCache(orgId)
}
