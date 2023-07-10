package service

import (
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

var Log = logger.GetDefaultLogger()

func AssemblyUserIdInfo(baseUserInfo *bo.BaseUserInfoBo) *vo.UserIDInfo {
	return &vo.UserIDInfo{
		ID:         baseUserInfo.UserId,
		UserID:     baseUserInfo.UserId,
		Name:       baseUserInfo.Name,
		Avatar:     baseUserInfo.Avatar,
		EmplID:     baseUserInfo.OutUserId,
		IsDeleted:  baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted,
		IsDisabled: baseUserInfo.OrgUserStatus == consts.AppStatusDisabled,
	}
}

func checkNeedUpdate(updateFields []string, field string) bool {
	if updateFields == nil || len(updateFields) == 0 {
		return true
	}
	bol, err := slice.Contain(updateFields, field)
	if err != nil {
		return false
	}
	return bol
}

func checkNeedUpdateLc(lcUpdates map[string]interface{}, field string) bool {
	if lcUpdates == nil || len(lcUpdates) == 0 {
		return true
	}
	_, ok := lcUpdates[field]
	return ok
}
