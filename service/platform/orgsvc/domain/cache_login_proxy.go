package domain

import (
	"strings"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
)

const LoginCodeFreezeTimeBlankObj = "BLANK"

func CheckSMSLoginCodeFreezeTime(authType, addressType int, address string) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.CacheSmsSendLoginCodeFreezeTime, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       address,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	exist, err := cache.Exist(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	if exist {
		return errs.BuildSystemErrorInfo(errs.SMSSendTimeLimitError)
	}
	return nil
}

func ClearSMSLoginCodeFreezeTime(authType, addressType int, phoneNumber string) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.CacheSmsSendLoginCodeFreezeTime, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func SetSMSLoginCodeFreezeTime(authType, addressType int, phoneNumber string, minute int) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.CacheSmsSendLoginCodeFreezeTime, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	err := cache.SetEx(key, LoginCodeFreezeTimeBlankObj, int64(60*minute))
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func ClearSMSLoginCode(authType, addressType int, phoneNumber string) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.CacheSmsLoginCode, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func ClearPwdLoginCode(loginName string) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.CacheLoginGraphCode, map[string]interface{}{
		consts.CacheKeyLoginNameConstName: loginName,
	})
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func SetSMSLoginCode(authType, addressType int, phoneNumber string, authCode string) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.CacheSmsLoginCode, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	err := cache.SetEx(key, authCode, int64(60*5))
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func GetSMSLoginCode(authType, addressType int, phoneNumber string) (string, errs.SystemErrorInfo) {
	key, _ := util.ParseCacheKey(sconsts.CacheSmsLoginCode, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	value, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return "", errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if value == "" {
		log.Error("登录验证码为空: " + phoneNumber)
		return "", errs.BuildSystemErrorInfo(errs.SMSLoginCodeInvalid)
	}
	return value, nil
}

func GetPwdLoginCode(loginName string) (string, errs.SystemErrorInfo) {
	key, _ := util.ParseCacheKey(sconsts.CacheLoginGraphCode, map[string]interface{}{
		consts.CacheKeyLoginNameConstName: loginName,
	})
	value, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return "", errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if value == "" {
		log.Error("登录验证码为空" + loginName)
		return "", errs.BuildSystemErrorInfo(errs.SMSLoginCodeInvalid)
	}
	return value, nil
}

func SetPwdLoginCode(loginName, code string) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.CacheLoginGraphCode, map[string]interface{}{
		consts.CacheKeyLoginNameConstName: loginName,
	})
	err := cache.SetEx(key, code, int64(60*5))
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	return nil
}

func SetPwdLoginCodeVerifyFreeze(authType, addressType int, loginName string) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.CachePwdLoginCodeVerifyFreeze, map[string]interface{}{
		consts.CacheKeyLoginNameConstName:   loginName,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	err := cache.SetEx(key, "v", 3600)
	if err != nil {
		log.Error(err)
		return errs.RedisOperateError
	}
	return nil
}

func GetPwdLoginCodeVerifyFreeze(authType, addressType int, loginName string) (string, errs.SystemErrorInfo) {
	key, _ := util.ParseCacheKey(sconsts.CachePwdLoginCodeVerifyFreeze, map[string]interface{}{
		consts.CacheKeyLoginNameConstName:   loginName,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	value, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return "", errs.RedisOperateError
	}
	return value, nil
}

func SetPwdLoginCodeVerifyFailTimesIncrement(authType, addressType int, loginName string) (int64, errs.SystemErrorInfo) {
	key, _ := util.ParseCacheKey(sconsts.CachePwdLoginCodeVerifyFailTimes, map[string]interface{}{
		consts.CacheKeyLoginNameConstName:   loginName,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	c, err := cache.Incrby(key, 1)
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return c, nil
}

func ClearPwdLoginCodeVerifyFailTimesIncrement(authType, addressType int, loginName string) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.CachePwdLoginCodeVerifyFailTimes, map[string]interface{}{
		consts.CacheKeyLoginNameConstName:   loginName,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func SetSMSLoginCodeVerifyFailTimesIncrement(authType, addressType int, phoneNumber string) (int64, errs.SystemErrorInfo) {
	key, _ := util.ParseCacheKey(sconsts.CacheSmsLoginCodeVerifyFailTimes, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	c, err := cache.Incrby(key, 1)
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return c, nil
}

func ClearSMSLoginCodeVerifyFailTimesIncrement(authType, addressType int, phoneNumber string) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.CacheSmsLoginCodeVerifyFailTimes, map[string]interface{}{
		consts.CacheKeyPhoneConstName:       phoneNumber,
		consts.CacheKeyAuthTypeConstName:    authType,
		consts.CacheKeyAddressTypeConstName: addressType,
	})
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func ClearSMSLoginCodeAllCache(authType, addressType int, phoneNumber string) errs.SystemErrorInfo {
	err := ClearSMSLoginCode(authType, addressType, phoneNumber)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	err = ClearSMSLoginCodeVerifyFailTimesIncrement(authType, addressType, phoneNumber)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func GetPhoneNumberWhiteList() ([]string, errs.SystemErrorInfo) {
	key := sconsts.CachePhoneNumberWhiteList
	value, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if value == "" {
		//如果没有就设置
		defaultWhiteList, setErr := SetPhoneNumberWhiteList()
		if setErr != nil {
			return nil, setErr
		}

		return defaultWhiteList, nil
	}

	result := &[]string{}
	err = json.FromJson(value, result)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
	}

	return *result, nil
}

func SetPhoneNumberWhiteList() ([]string, errs.SystemErrorInfo) {
	key := sconsts.CachePhoneNumberWhiteList
	whiteList := []string{
		"18891111111",
		"18891111112",
		"18891111113",
		"18891111114",
		"18891111115",
		"18616215755",
		"18117493299",
		"15618960078",
		"17621142248",
		"18516232262",
		"18917633966",
		"13037568259",
		"18917637631",
		"15618931187",
		"18221304331",
		"13621822254",
	}
	whiteListJson, err := json.ToJson(whiteList)
	if err != nil {
		log.Error(err)
		return whiteList, errs.BuildSystemErrorInfo(errs.JSONConvertError)
	}
	err = cache.SetEx(key, whiteListJson, 60*60*24*30)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return whiteList, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	return whiteList, nil
}

func AuthCodeVerify(authType, addressType int, contactAddress, authCode string) errs.SystemErrorInfo {
	localCode, err := GetSMSLoginCode(authType, addressType, contactAddress)
	if err != nil {
		log.Error(err)
		return err
	}

	if !strings.EqualFold(localCode, authCode) {
		//计数，防止刷接口
		times, err := SetSMSLoginCodeVerifyFailTimesIncrement(authType, addressType, contactAddress)
		if err != nil {
			return err
		}
		if times > 60 {
			//如果失败五次，则清空其他缓存，让用户重新发送
			err := ClearSMSLoginCodeAllCache(authType, addressType, contactAddress)
			if err != nil {
				log.Error(err)
			}
			return errs.BuildSystemErrorInfo(errs.SMSLoginCodeVerifyFailTimesOverLimit)
		}
		return errs.BuildSystemErrorInfo(errs.SMSLoginCodeNotMatch)
	}

	//登录成功
	//清理验证码相关缓存
	err = ClearSMSLoginCodeAllCache(authType, addressType, contactAddress)
	if err != nil {
		log.Error(err)
	}
	return nil
}

func ClearChangeLoginNameSign(orgId, userId int64, addressType int) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.ChangeLoginNameSign, map[string]interface{}{
		consts.CacheKeyOfOrg:         orgId,
		consts.CacheKeyOfUser:        userId,
		consts.CacheKeyOfAddressType: addressType,
	})
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func GetChangeLoginNameSign(orgId, userId int64, addressType int) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.ChangeLoginNameSign, map[string]interface{}{
		consts.CacheKeyOfOrg:         orgId,
		consts.CacheKeyOfUser:        userId,
		consts.CacheKeyOfAddressType: addressType,
	})
	value, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return errs.RedisOperateError
	}
	if value == "" {
		return errs.ChangeLoginNameInvalid
	}
	return nil
}

func SetChangeLoginNameSign(orgId, userId int64, addressType int) errs.SystemErrorInfo {
	key, _ := util.ParseCacheKey(sconsts.ChangeLoginNameSign, map[string]interface{}{
		consts.CacheKeyOfOrg:         orgId,
		consts.CacheKeyOfUser:        userId,
		consts.CacheKeyOfAddressType: addressType,
	})
	err := cache.SetEx(key, "true", int64(60*5))
	if err != nil {
		log.Error(err)
		return errs.RedisOperateError
	}

	return nil
}
