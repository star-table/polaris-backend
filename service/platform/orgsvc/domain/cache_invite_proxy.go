package domain

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
)

func SetUserInviteCodeInfo(inviteCode string, inviteInfo bo.InviteInfoBo) errs.SystemErrorInfo {
	inviteInfoJson, err := json.ToJson(inviteInfo)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.JSONConvertError)
	}
	err = cache.SetEx(sconsts.CacheUserInviteCode+inviteCode, inviteInfoJson, int64(sconsts.CacheUserInviteCodeExpire))
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

// 多个手机号邀请时的邀请信息缓存
func SetUserInviteInfoForPhones(inviteCode string, inviteInfo interface{}, expiredDuration int64) errs.SystemErrorInfo {
	inviteInfoJson, err := json.ToJson(inviteInfo)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.JSONConvertError)
	}
	err = cache.SetEx(sconsts.CacheUserInviteCode+inviteCode, inviteInfoJson, expiredDuration)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func GetUserInviteCodeInfo(inviteCode string) (*bo.InviteInfoBo, errs.SystemErrorInfo) {
	inviteInfoJson, err := cache.Get(sconsts.CacheUserInviteCode + inviteCode)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	log.Infof("邀请code %s 对应的邀请信息 %s", inviteCode, inviteInfoJson)
	if inviteInfoJson == "" {
		return nil, errs.BuildSystemErrorInfo(errs.InviteCodeInvalid)
	}
	inviteInfo := &bo.InviteInfoBo{}
	err = json.FromJson(inviteInfoJson, inviteInfo)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
	}
	return inviteInfo, nil
}
