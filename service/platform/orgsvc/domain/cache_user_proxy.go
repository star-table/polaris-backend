package domain

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/language/english"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/ordervo"
	"github.com/star-table/polaris-backend/facade/orderfacade"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

func GetBaseUserOutInfoBatch(orgId int64, userIds []int64) ([]bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	keys := make([]interface{}, len(userIds))
	for i, userId := range userIds {
		key, _ := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: userId,
		})
		keys[i] = key
	}
	resultList := make([]string, 0)
	if len(keys) > 0 {
		list, err := cache.MGet(keys...)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		resultList = list
	}
	baseUserOutInfoList := make([]bo.BaseUserOutInfoBo, 0)
	validUserIds := map[int64]bool{}
	for _, userInfoJson := range resultList {
		userOutInfoBo := &bo.BaseUserOutInfoBo{}
		err := json.FromJson(userInfoJson, userOutInfoBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		baseUserOutInfoList = append(baseUserOutInfoList, *userOutInfoBo)
		validUserIds[userOutInfoBo.UserId] = true
	}

	missUserIds := make([]int64, 0)
	//找不存在的
	if len(userIds) != len(validUserIds) {
		for _, userId := range userIds {
			if _, ok := validUserIds[userId]; !ok {
				missUserIds = append(missUserIds, userId)
			}
		}
	}

	//批量查外部信息
	outInfos, userErr := GetBaseUserOutInfoByUserIds(orgId, missUserIds)
	if userErr != nil {
		log.Error(userErr)
		return nil, userErr
	}

	if len(outInfos) > 0 {
		baseUserOutInfoList = append(baseUserOutInfoList, outInfos...)
	}

	return baseUserOutInfoList, nil
}

func GetBaseUserInfoBatch(orgId int64, originUserIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	//去重
	userIds := slice.SliceUniqueInt64(originUserIds)

	keys := make([]interface{}, len(userIds))
	for i, userId := range userIds {
		key, _ := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: userId,
		})
		keys[i] = key
	}
	resultList := make([]string, 0)
	if len(keys) > 0 {
		list, err := cache.MGet(keys...)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		resultList = list
	}
	baseUserInfoList := make([]bo.BaseUserInfoBo, 0)
	validUserIds := map[int64]bool{}
	for _, userInfoJson := range resultList {
		userInfoBo := &bo.BaseUserInfoBo{}
		err := json.FromJson(userInfoJson, userInfoBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		baseUserInfoList = append(baseUserInfoList, *userInfoBo)
		validUserIds[userInfoBo.UserId] = true
	}

	log.Infof("from cache %s", json.ToJsonIgnoreError(baseUserInfoList))
	missUserIds := make([]int64, 0)
	//找不存在的
	if len(userIds) != len(validUserIds) {
		for _, userId := range userIds {
			if _, ok := validUserIds[userId]; !ok {
				missUserIds = append(missUserIds, userId)
			}
		}
	}

	missUserInfos, userErr := getLocalBaseUserInfoBatch(orgId, missUserIds)
	if userErr != nil {
		log.Error(userErr)
		return nil, userErr
	}
	if len(missUserInfos) > 0 {
		baseUserInfoList = append(baseUserInfoList, missUserInfos...)
	}
	if ok, _ := slice.Contain(originUserIds, int64(0)); ok {
		baseUserInfoList = append(baseUserInfoList, bo.BaseUserInfoBo{
			UserId:             0,
			OutUserId:          "",
			OrgId:              orgId,
			OutOrgId:           "",
			Name:               english.WordTransLate("未分配"),
			NamePy:             "",
			Avatar:             consts.AvatarForUnallocated,
			HasOutInfo:         false,
			HasOrgOutInfo:      false,
			OutOrgUserId:       "",
			OrgUserIsDelete:    0,
			OrgUserStatus:      0,
			OrgUserCheckStatus: 0,
		})
	}

	//获取用户外部信息
	baseUserOutInfos, err := GetBaseUserOutInfoBatch(orgId, userIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	outInfoMap := maps.NewMap("UserId", baseUserOutInfos)
	for i, _ := range baseUserInfoList {
		userInfo := baseUserInfoList[i]
		if outInfoInterface, ok := outInfoMap[userInfo.UserId]; ok {
			outInfo := outInfoInterface.(bo.BaseUserOutInfoBo)

			userInfo.OutUserId = outInfo.OutUserId
			userInfo.OutOrgUserId = outInfo.OutOrgUserId
			userInfo.OutOrgId = outInfo.OutOrgId
			userInfo.HasOutInfo = outInfo.OutUserId != ""
			userInfo.HasOrgOutInfo = outInfo.OutOrgId != ""
		}
		baseUserInfoList[i] = userInfo
	}

	//按照原始顺序排序
	baseUserInfoMap := map[int64]bo.BaseUserInfoBo{}
	for _, infoBo := range baseUserInfoList {
		baseUserInfoMap[infoBo.UserId] = infoBo
	}

	resBo := []bo.BaseUserInfoBo{}
	for _, id := range originUserIds {
		if info, ok := baseUserInfoMap[id]; ok {
			resBo = append(resBo, info)
		}
	}
	log.Infof("user info: %s", json.ToJsonIgnoreError(resBo))

	return resBo, nil
}

func ClearBaseUserInfo(orgId, userId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

//批量清楚用户缓存信息
func ClearBaseUserInfoBatch(orgId int64, userIds []int64) errs.SystemErrorInfo {
	keys := make([]interface{}, 0)
	for _, userId := range userIds {
		key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: userId,
		})
		if err5 != nil {
			log.Error(err5)
			return err5
		}
		keys = append(keys, key)
	}
	_, err := cache.Del(keys...)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

//sourceChannel可以为空
func GetBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	if userId == 0 {
		//系统创建
		return &bo.BaseUserInfoBo{
			OrgId:  orgId,
			Name:   english.WordTransLate("未分配"),
			Avatar: consts.AvatarForUnallocated,
		}, nil
	}

	key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	baseUserInfoJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	baseUserInfo := &bo.BaseUserInfoBo{}
	if baseUserInfoJson != "" {
		err := json.FromJson(baseUserInfoJson, baseUserInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
	} else {
		userInfo, errorInfo := getLocalBaseUserInfo(orgId, userId, key)
		if errorInfo != nil {
			log.Error(errorInfo)
			return nil, errorInfo
		}
		baseUserInfo = userInfo
	}

	//这里不存缓存，动态获取
	baseUserOutInfo, sysErr := GetBaseUserOutInfo(orgId, userId)
	if sysErr != nil {
		log.Error(sysErr)
		return nil, sysErr
	}
	baseUserInfo.OutUserId = baseUserOutInfo.OutUserId
	baseUserInfo.OutOrgId = baseUserOutInfo.OutOrgId
	baseUserInfo.HasOutInfo = baseUserInfo.OutUserId != ""
	baseUserInfo.HasOrgOutInfo = baseUserInfo.OutOrgId != ""
	baseUserInfo.OutOrgUserId = baseUserOutInfo.OutOrgUserId

	return baseUserInfo, nil
}

func GetBaseUserOutInfo(orgId int64, userId int64) (*bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	if userId == 0 {
		//系统创建
		return &bo.BaseUserOutInfoBo{
			OrgId: orgId,
		}, nil
	}

	key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}
	baseUserOutInfoJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if baseUserOutInfoJson != "" {
		baseUserOutInfo := &bo.BaseUserOutInfoBo{}
		err := json.FromJson(baseUserOutInfoJson, baseUserOutInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		return baseUserOutInfo, nil
	} else {
		//用户外部信息
		userOutInfo := &po.PpmOrgUserOutInfo{}
		_ = mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
			consts.TcUserId:   userId,
		}, userOutInfo)

		outInfo := bo.BaseUserOutInfoBo{
			UserId:       userId,
			OrgId:        orgId,
			OutUserId:    userOutInfo.OutUserId,
			OutOrgUserId: userOutInfo.OutOrgUserId,
		}
		//组织外部信息
		orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
		err = mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
		}, orgOutInfo)
		if err != nil {
			if err == db.ErrNoMoreRows {

			} else {
				log.Error(err)
				return nil, errs.MysqlOperateError
			}
		} else {
			outInfo.OutOrgId = orgOutInfo.OutOrgId
		}

		baseUserOutInfoJson, err := json.ToJson(outInfo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		err = cache.SetEx(key, baseUserOutInfoJson, consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return &outInfo, nil
	}
}

func GetBaseUserOutInfoByUserIds(orgId int64, userIds []int64) ([]bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	log.Infof("批量获取用户外部信息 %d, %s", orgId, json.ToJsonIgnoreError(userIds))

	resultList := make([]bo.BaseUserOutInfoBo, 0)

	if userIds == nil || len(userIds) == 0 {
		return resultList, nil
	}

	//用户外部信息
	userOutInfos := &[]po.PpmOrgUserOutInfo{}
	err := mysql.SelectAllByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(userIds),
	}, userOutInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	//组织外部信息
	orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
	err = mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}, orgOutInfo)

	msetArgs := map[string]string{}
	keys := make([]string, 0)
	for _, userOutInfo := range *userOutInfos {
		key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: userOutInfo.UserId,
		})
		if err5 != nil {
			log.Error(err5)
			return nil, err5
		}
		keys = append(keys, key)

		outInfo := bo.BaseUserOutInfoBo{
			UserId:       userOutInfo.UserId,
			OrgId:        orgId,
			OutUserId:    userOutInfo.OutUserId,
			OutOrgId:     orgOutInfo.OutOrgId,
			OutOrgUserId: userOutInfo.OutOrgUserId,
		}

		resultList = append(resultList, outInfo)
		msetArgs[key] = json.ToJsonIgnoreError(outInfo)
	}

	if len(msetArgs) > 0 {
		err = cache.MSet(msetArgs)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		for _, key := range keys {
			_, _ = cache.Expire(key, consts.GetCacheBaseUserInfoExpire())
		}
	}()

	return resultList, nil
}

//sourceChannel可以为空
func getLocalBaseUserInfo(orgId, userId int64, key string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	user, err := GetUserInfoByUserId(userId)
	if err != nil {
		return nil, err
	}

	baseUserInfo := &bo.BaseUserInfoBo{
		UserId:             user.Id,
		Name:               user.Name,
		NamePy:             user.NamePinyin,
		Avatar:             user.Avatar,
		OrgId:              orgId,
		OrgUserIsDelete:    2,
		OrgUserStatus:      1,
		OrgUserCheckStatus: 1,
	}

	if orgId > 0 {
		newestUserOrganization, err1 := GetUserOrganizationNewestRelation(orgId, userId)
		if err1 != nil {
			log.Error(err1)
			return nil, err1
		}
		baseUserInfo.OrgUserIsDelete = newestUserOrganization.IsDelete
		baseUserInfo.OrgUserStatus = newestUserOrganization.Status
		baseUserInfo.OrgUserCheckStatus = newestUserOrganization.CheckStatus
	}

	baseUserInfoJson, err2 := json.ToJson(baseUserInfo)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err2)
	}
	err2 = cache.SetEx(key, baseUserInfoJson, consts.GetCacheBaseUserInfoExpire())
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err2)
	}
	return baseUserInfo, nil
}

func GetUserInfoByUserId(userId int64) (*po.PpmOrgUser, errs.SystemErrorInfo) {
	user := &po.PpmOrgUser{}
	err := mysql.SelectById(user.TableName(), userId, user)
	if err != nil {
		log.Errorf("[GetUserInfoByUserId] userId: %d, err: %v", userId, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return user, nil
}

func getLocalBaseUserInfoBatch(orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	log.Infof("批量获取用户信息 %d, %s", orgId, json.ToJsonIgnoreError(userIds))

	baseUserInfos := make([]bo.BaseUserInfoBo, 0)

	if userIds == nil || len(userIds) == 0 {
		return baseUserInfos, nil
	}

	users := &[]po.PpmOrgUser{}
	err := mysql.SelectAllByCond(consts.TableUser, db.Cond{
		consts.TcId: db.In(userIds),
	}, users)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	//获取关联列表，要做去重
	userOrganizationPos := &[]po.PpmOrgUserOrganization{}
	_, selectErr := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcUserId: db.In(userIds),
	}, nil, 0, -1, "id asc", userOrganizationPos)
	if selectErr != nil {
		log.Error(selectErr)
		return nil, errs.MysqlOperateError
	}

	//id升序，保留最新: key: userId, value: po
	userOrgMap := map[int64]po.PpmOrgUserOrganization{}
	for _, userOrg := range *userOrganizationPos {
		userOrgMap[userOrg.UserId] = userOrg
	}

	for _, user := range *users {
		baseUserInfo := bo.BaseUserInfoBo{
			UserId: user.Id,
			Name:   user.Name,
			NamePy: user.NamePinyin,
			Avatar: user.Avatar,
			OrgId:  orgId,
		}

		if userOrg, ok := userOrgMap[user.Id]; ok {
			baseUserInfo.OrgUserIsDelete = userOrg.IsDelete
			baseUserInfo.OrgUserStatus = userOrg.Status
			baseUserInfo.OrgUserCheckStatus = userOrg.CheckStatus
		}

		baseUserInfos = append(baseUserInfos, baseUserInfo)
	}

	msetArgs := map[string]string{}
	keys := make([]string, 0)
	for _, baseUserInfo := range baseUserInfos {
		key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: baseUserInfo.UserId,
		})
		if err5 != nil {
			log.Error(err5)
			return nil, err5
		}
		msetArgs[key] = json.ToJsonIgnoreError(baseUserInfo)
		keys = append(keys, key)
	}

	if len(msetArgs) > 0 {
		err = cache.MSet(msetArgs)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		for _, key := range keys {
			_, _ = cache.Expire(key, consts.GetCacheBaseUserInfoExpire())
		}
	}()
	return baseUserInfos, nil
}

func GetUserConfigInfo(orgId int64, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	userConfig, err := getUserConfigInfo(orgId, userId)
	if err != nil {
		userConfig = &bo.UserConfigBo{}
		userConfigBo, err2 := InsertUserConfig(orgId, userId)
		if err2 != nil {
			log.Error(err2)
			return nil, errs.BuildSystemErrorInfo(errs.UserConfigUpdateError, err2)
		}
		err3 := copyer.Copy(userConfigBo, userConfig)
		if err3 != nil {
			log.Error(err3)
			return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err3)
		}
	}
	return userConfig, nil
}

func getUserConfigInfo(orgId int64, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	userConfigJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	userConfigBo := &bo.UserConfigBo{}
	if userConfigJson != "" {
		err := json.FromJson(userConfigJson, userConfigBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		return userConfigBo, nil
	} else {
		userConfig := &po.PpmOrgUserConfig{}
		err = mysql.SelectOneByCond(userConfig.TableName(), db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcUserId:   userId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, userConfig)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(userConfig, userConfigBo)
		userConfigJson, err = json.ToJson(userConfigBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		err = cache.Set(key, userConfigJson)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return userConfigBo, nil
	}
}

func DeleteUserConfigInfo(orgId int64, userId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	return nil
}

func ClearUserCacheInfo(token string) errs.SystemErrorInfo {
	userCacheKey := sconsts.CacheUserToken + token
	_, err := cache.Del(userCacheKey)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func ResetOrgPayNum(orgId int64) error {
	// 查询订单的付费信息
	orderPayInfo := orderfacade.GetOrderPayInfo(ordervo.GetOrderPayInfoReq{OrgId: orgId})
	if orderPayInfo.Failure() {
		log.Errorf("[SetPayUserInfo] GetOrderPayInfo err:%v, orgId:%v", orderPayInfo.Error(), orgId)
		return orderPayInfo.Error()
	}
	payNum := orderPayInfo.Data.PayNum

	payRangeInfoJson, errCache := cache.HGet(consts.CachePayRangeInfo, fmt.Sprintf("%d", orgId))
	if errCache != nil {
		log.Error(errCache)
		return errCache
	}
	if payRangeInfoJson != "" {
		rangeData := bo.PayRangeData{}
		_ = json.FromJson(payRangeInfoJson, &rangeData)
		rangeData.PayNum = payNum
		errCache = cache.HSet(consts.CachePayRangeInfo, fmt.Sprintf("%d", orgId), json.ToJsonIgnoreError(rangeData))
		if errCache != nil {
			log.Error(errCache)
			return errCache
		}
	}

	return nil
}

// ClearAllOrgUserPayCache 清除org的用户的付费缓存，解决付费后延迟问题
func ClearAllOrgUserPayCache(orgId int64) errs.SystemErrorInfo {
	users, err := GetUserOrgInfos(orgId)
	if err != nil {
		log.Errorf("[ClearAllOrgUserPayCache] GetUserOrgInfos orgId:%v, err:%v", orgId, err)
		return err
	}
	keys := make([]interface{}, 0, len(users))
	for _, user := range users {
		key, _ := util.ParseCacheKey(sconsts.CacheUserCheckPay, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: user.UserId,
		})
		keys = append(keys, key)
	}
	if len(keys) > 0 {
		_, err := cache.Del(keys)
		if err != nil {
			log.Errorf("[ClearAllOrgUserPayCache] Del keys:%v, err:%v", keys, err)
		}
	}

	return nil
}

// GetDuration2Noon 计算当前时间到中午或第二天中午的秒数
func GetDuration2Noon() int64 {
	now := time.Now()
	hour := now.Hour()
	durationSec := int64(0)
	if hour <= 12 {
		// 当天中午
		durationSec = int64((12 - hour) * 3600)
	} else {
		// 当第二天中午
		durationSec = int64((hour - 12) * 3600)
	}
	randSec := rand.Intn(3600)

	return int64(randSec) + durationSec
}

func SetShareUrl(key, url string) errs.SystemErrorInfo {
	if key == "" {
		return nil
	}
	cacheErr := cache.SetEx(sconsts.CacheShareUrl+key, url, 60*60*24*3)
	if cacheErr != nil {
		log.Error(cacheErr)
		return errs.CacheProxyError
	}

	return nil
}

func GetShareUrl(key string) (string, errs.SystemErrorInfo) {
	url, cacheErr := cache.Get(sconsts.CacheShareUrl + key)
	if cacheErr != nil {
		log.Error(cacheErr)
		return "", errs.CacheProxyError
	}

	return url, nil
}

func NeedRemindPayExpire(orgId, userId int64, payEndTime time.Time) (bool, errs.SystemErrorInfo) {
	expireTime := payEndTime.Unix() - time.Now().Unix()
	if expireTime <= 0 {
		//超出时间不需要提醒
		return false, nil
	}

	key, err5 := util.ParseCacheKey(sconsts.CachePayExpireRemind, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return false, err5
	}
	isExist, isExistErr := cache.Exist(key)
	if isExistErr != nil {
		log.Error(isExistErr)
		return false, errs.CacheProxyError
	}
	if isExist {
		isRemind, err := cache.HGet(key, strconv.FormatInt(userId, 10))
		if err != nil {
			log.Error(err)
			return false, errs.CacheProxyError
		}
		if isRemind == "" {
			err := cache.HSet(key, strconv.FormatInt(userId, 10), "true")
			if err != nil {
				log.Error(err)
				return false, errs.CacheProxyError
			}
		} else {
			return false, nil
		}
	} else {
		err := cache.HSet(key, strconv.FormatInt(userId, 10), "true")
		if err != nil {
			log.Error(err)
			return false, errs.CacheProxyError
		}
		_, expireErr := cache.Expire(key, expireTime)
		if expireErr != nil {
			log.Error(expireErr)
			return false, errs.CacheProxyError
		}
		return true, nil
	}

	return false, nil
}

func NeedRemindPayOverdue(orgId, userId int64, payEndTime time.Time) (bool, errs.SystemErrorInfo) {
	expireTime := payEndTime.AddDate(0, 0, 1).Unix() - time.Now().Unix()
	if expireTime <= 0 {
		//超出时间不需要提醒
		return false, nil
	}

	key, err5 := util.ParseCacheKey(sconsts.CachePayOverdueRemind, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return false, err5
	}
	isExist, isExistErr := cache.Exist(key)
	if isExistErr != nil {
		log.Error(isExistErr)
		return false, errs.CacheProxyError
	}
	if isExist {
		isRemind, err := cache.HGet(key, strconv.FormatInt(userId, 10))
		if err != nil {
			log.Error(err)
			return false, errs.CacheProxyError
		}
		if isRemind == "" {
			err := cache.HSet(key, strconv.FormatInt(userId, 10), "true")
			if err != nil {
				log.Error(err)
				return false, errs.CacheProxyError
			}
			return true, nil
		} else {
			return false, nil
		}
	} else {
		err := cache.HSet(key, strconv.FormatInt(userId, 10), "true")
		if err != nil {
			log.Error(err)
			return false, errs.CacheProxyError
		}
		_, expireErr := cache.Expire(key, expireTime)
		if expireErr != nil {
			log.Error(expireErr)
			return false, errs.CacheProxyError
		}
		return true, nil
	}

	return false, nil
}
