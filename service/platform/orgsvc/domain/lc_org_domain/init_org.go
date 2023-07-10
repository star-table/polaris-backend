package lc_org_domain

import (
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
)

var (
	log = logger.GetDefaultLogger()
)

//// 将组织设为付费组织
//func SetOrgPaid(orgId int64, needSetToPaid int) errs.SystemErrorInfo {
//	if needSetToPaid != 1 {
//		return nil
//	}
//	orcInfo := po.PpmOrcConfig{}
//	queryCond := db.Cond{
//		consts.TcIsDelete: consts.AppIsNoDelete,
//		consts.TcOrgId:    orgId,
//	}
//	err := mysql.SelectOneByCond(consts.TableOrgConfig, queryCond, &orcInfo)
//	if err == db.ErrNoMoreRows {
//		// 插入新的 org config 记录
//		_, err := LcOrgConfigInfoInit(orgId)
//		if err != nil {
//			return err
//		}
//		return nil
//	} else if err != nil {
//		log.Error(err)
//		return errs.BuildSystemErrorInfoWithMessage(errs.OrgConfigNotExist, "org config 不存在。")
//	}
//	endTime := time.Now().AddDate(100, 0, 0)
//	_, err = mysql.UpdateSmartWithCond(consts.TableOrgConfig, queryCond, mysql.Upd{
//		consts.TcPayLevel:   consts.PayLevelEnterprise,
//		consts.TcPayEndTime: endTime,
//	})
//	if err != nil {
//		log.Error(err)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	// 删除组织信息缓存
//	busiErr := ClearCacheBaseOrgInfo(orgId)
//	if busiErr != nil {
//		log.Errorf("ClearCache BaseOrgInfo err: %q\n", busiErr)
//		return errs.BuildSystemErrorInfo(errs.RedisOperateError, busiErr)
//	}
//	clearErr := ClearOrgConfig(orgId)
//	if clearErr != nil {
//		log.Error(clearErr)
//	}
//	return nil
//}

func ClearCacheBaseOrgInfo(orgId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheBaseOrgInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}

	_, err := cache.Del(key)

	if err != nil {
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}

func ClearOrgConfig(orgId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheOrgPayFunction, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})

	if err5 != nil {
		log.Error(err5)
		return err5
	}

	_, err := cache.Del(key)

	if err != nil {
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}

// 增加 orc 配置记录。
//func LcOrgConfigInfoInit(orgId int64) (int64, errs.SystemErrorInfo) {
//	sysConfig := &po.PpmOrcConfig{}
//	payLevel := &po.PpmBasPayLevel{}
//	// id 为 4 的表示 VIP 版。
//	err := mysql.SelectById(payLevel.TableName(), 4, payLevel)
//	if err != nil {
//		log.Error(err)
//		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	orgConfigId, err1 := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrgConfig)
//	if err1 != nil {
//		log.Error(err1)
//		return 0, err1
//	}
//
//	sysConfig.Id = orgConfigId
//	sysConfig.OrgId = orgId
//	sysConfig.TimeZone = "Asia/Shanghai"
//	sysConfig.TimeDifference = "+08:00"
//	sysConfig.PayLevel = consts.PayLevelEnterprise
//	sysConfig.PayStartTime = time.Now()
//	// 大数据用户购买了融合版本极星，因此直接给 100 年的时间限制。
//	sysConfig.PayEndTime = time.Now().Add(time.Duration(payLevel.Duration) * time.Second * 100)
//	sysConfig.Language = "zh-CN"
//	sysConfig.RemindSendTime = "09:00:00"
//	sysConfig.DatetimeFormat = "yyyy-MM-dd HH:mm:ss"
//	sysConfig.PasswordLength = 6
//	sysConfig.PasswordRule = 1
//	sysConfig.MaxLoginFailCount = 0
//	sysConfig.Status = consts.AppStatusEnable
//	err2 := mysql.Insert(sysConfig)
//	if err2 != nil {
//		log.Error("LcOrgConfigInfoInit 组织初始化，添加组织配置信息时异常: " + strs.ObjectToString(err2))
//		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
//	}
//
//	return orgConfigId, nil
//}
