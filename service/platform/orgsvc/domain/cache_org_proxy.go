package domain

import (
	"fmt"
	"time"

	"github.com/spf13/cast"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/types"
	"github.com/star-table/polaris-backend/common/model/vo/ordervo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orderfacade"

	msgPb "github.com/star-table/interface/golang/msg/v1"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/facade/common/report"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/businees"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetOrgIdByOutOrgId(outOrgId string) (int64, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheOutOrgIdRelationId, map[string]interface{}{
		consts.CacheKeyOutOrgIdConstName: outOrgId,
	})
	if err5 != nil {
		log.Error(err5)
		return 0, err5
	}
	orgIdInfoJson, err := cache.Get(key)
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if orgIdInfoJson != "" {
		orgIdInfo := &bo.OrgIdInfo{}
		err := json.FromJson(orgIdInfoJson, orgIdInfo)
		if err != nil {
			return 0, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return orgIdInfo.OrgId, nil
	} else {
		orgBo, err := GetOrgByOutOrgId(outOrgId)
		if err != nil {
			log.Error(err)
			return 0, err
		}
		orgIdInfo := bo.OrgIdInfo{
			OutOrgId: outOrgId,
			OrgId:    orgBo.Id,
		}
		orgIdInfoJson = json.ToJsonIgnoreError(orgIdInfo)
		err1 := cache.SetEx(key, orgIdInfoJson, consts.GetCacheBaseExpire())
		if err1 != nil {
			return 0, errs.BuildSystemErrorInfo(errs.RedisOperateError, err1)
		}
		return orgIdInfo.OrgId, nil
	}
}

func GetBaseOrgOutInfo(orgId int64) (*bo.BaseOrgOutInfoBo, errs.SystemErrorInfo) {
	key, sysErr := util.ParseCacheKey(sconsts.CacheBaseOrgOutInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if sysErr != nil {
		log.Error(sysErr)
		return nil, sysErr
	}
	outOrgInfoJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if outOrgInfoJson != "" {
		orgOutInfoBo := &bo.BaseOrgOutInfoBo{}
		err := json.FromJson(outOrgInfoJson, orgOutInfoBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return orgOutInfoBo, nil
	} else {
		var orgOutInfos []*po.PpmOrgOrganizationOutInfo
		err = mysql.SelectAllByCond("ppm_org_organization_out_info", db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: db.In([]int{0, consts.AppIsNoDelete}),
		}, &orgOutInfos)
		if err != nil || len(orgOutInfos) == 0 {
			log.Infof("get org out info orgId: %d err: %v", orgId, err)
			return nil, errs.OrgOutInfoNotExist
		}
		var orgOutInfo *po.PpmOrgOrganizationOutInfo
		for _, outInfo := range orgOutInfos {
			if outInfo.OutOrgId != "" { // 优先拿有out_org_id的
				orgOutInfo = outInfo
				break
			}
		}
		if orgOutInfo == nil { // 没有就拿第一个
			orgOutInfo = orgOutInfos[0]
		}
		orgOutInfoBo := &bo.BaseOrgOutInfoBo{
			OrgId:         orgId,
			OutOrgId:      orgOutInfo.OutOrgId,
			SourceChannel: orgOutInfo.SourceChannel,
		}
		outOrgInfoJson = json.ToJsonIgnoreError(orgOutInfoBo)
		err = cache.SetEx(key, outOrgInfoJson, consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return orgOutInfoBo, nil
	}
}

func GetBaseOrgInfoByOutOrgId(outOrgId string) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	orgId, err := GetOrgIdByOutOrgId(outOrgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GetBaseOrgInfo(orgId)
}

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

// GetBaseOrgInfo usercenter有一份一样的代码，导致缓存有问题，所以这个方法目前先要兼容这些问题
func GetBaseOrgInfo(orgId int64) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	baseOrgInfo := &bo.BaseOrgInfoBo{}

	key, sysErr := util.ParseCacheKey(sconsts.CacheBaseOrgInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if sysErr != nil {
		log.Error(sysErr)
		return nil, sysErr
	}

	// 拿缓存
	if value, err := cache.Get(key); err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	} else if value != "" {
		if err = json.FromJson(value, baseOrgInfo); err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	}

	// 没拿到缓存或者缓存有问题，从DB拿数据
	if baseOrgInfo.OrgId == 0 || baseOrgInfo.SourceChannel == "" {
		orgOutInfo, sysErr := GetBaseOrgOutInfo(orgId)
		if sysErr != nil {
			log.Errorf("[GetBaseOrgInfo] GetBaseOrgOutInfo orgId:%v, err:%v", orgId, sysErr.Error())
			return nil, sysErr
		}

		var orgInfo po.PpmOrgOrganization
		if err := mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
			consts.TcId:       orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, &orgInfo); err != nil {
			log.Errorf("[GetBaseOrgInfo] get ppm_org_organization orgId:%v, err:%v", orgId, strs.ObjectToString(err))
			return nil, errs.BuildSystemErrorInfo(errs.OrgNotExist)
		}

		baseOrgInfo.OrgId = orgId
		baseOrgInfo.OrgName = orgInfo.Name
		baseOrgInfo.OrgOwnerId = orgInfo.Owner
		baseOrgInfo.Creator = orgInfo.Creator
		baseOrgInfo.OutOrgId = orgOutInfo.OutOrgId
		baseOrgInfo.SourceChannel = orgOutInfo.SourceChannel // 这个才是目前正确的source_channel

		log.Infof("[GetBaseOrgInfo] getOrgInfo: %v", json.ToJsonIgnoreError(baseOrgInfo))

		// 写入缓存
		if value, err := json.ToJson(baseOrgInfo); err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		} else {
			if err = cache.SetEx(key, value, consts.GetCacheBaseExpire()); err != nil {
				return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
			}
		}
	}

	return baseOrgInfo, nil
}

func GetOutDeptAndInnerDept(orgId int64, tx *sqlbuilder.Tx) (map[string]int64, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheDeptRelation, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	deptRelationListJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	deptRelationList := &map[string]int64{}
	if deptRelationListJson != "" {
		err = json.FromJson(deptRelationListJson, deptRelationList)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return *deptRelationList, nil
	} else {
		deptOurInfoList := &[]po.PpmOrgDepartmentOutInfo{}

		selectErr := queryDepartmentOutInfWithTx(tx, deptOurInfoList, orgId)

		log.Info("部门关联关系: " + strs.ObjectToString(deptOurInfoList))
		if selectErr != nil {
			return *deptRelationList, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
		}
		for _, v := range *deptOurInfoList {
			(*deptRelationList)[v.OutOrgDepartmentId] = v.DepartmentId
		}
		deptRelationListJson, err := json.ToJson(deptRelationList)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, deptRelationListJson, consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}

		return *deptRelationList, nil
	}
}

func queryDepartmentOutInfWithTx(tx *sqlbuilder.Tx, deptOurInfoList *[]po.PpmOrgDepartmentOutInfo, orgId int64) error {
	var selectErr error
	if tx != nil {
		//TODO 未定义TransSelectAllByCond，先使用SelectAllByCond(不使用事务不会有问题)
		selectErr = mysql.TransSelectAllByCond(*tx, consts.TableDepartmentOutInfo, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, deptOurInfoList)
	} else {
		selectErr = mysql.SelectAllByCond(consts.TableDepartmentOutInfo, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, deptOurInfoList)
	}
	return selectErr
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

// GetOrgConfigRich 获取组织配置
func GetOrgConfigRich(orgId int64) (*orgvo.OrgConfig, errs.SystemErrorInfo) {
	info, infoErr := GetOrgConfig(orgId)
	if infoErr != nil {
		log.Errorf("[GetOrgConfig] GetOrgConfig err: %v, orgId: %d", infoErr, orgId)
		return nil, infoErr
	}

	payLevel := info.PayLevel
	if info.PayLevel == consts.PayLevelDouble11Activity {
		payLevel = consts.PayLevelFlagship
	}

	org, err1 := GetOrgBoById(orgId)
	if err1 != nil {
		log.Errorf("[GetOrgConfig] GetOrgBoById err: %v, orgId: %d", err1, orgId)
		return nil, err1
	}
	appDeployType := GetAppDeployType(config.GetConfig().Application.RunMode)
	isPrivateDeploy := appDeployType == "private"
	res := &orgvo.OrgConfig{
		ID:              info.Id,
		OrgID:           info.OrgId,
		PayLevel:        payLevel,
		PayStartTime:    types.Time(info.PayStartTime),
		PayEndTime:      types.Time(info.PayEndTime),
		PayLevelTrue:    0,
		CreateTime:      types.Time(info.CreateTime),
		OrgMemberNumber: 0,
		IsGrayLevel:     false,
		SummaryAppID:    "",
		BasicShowSetting: &orgvo.BasicShowSetting{
			WorkBenchShow: true,  //默认工作栏展示
			SideBarShow:   false, //默认侧边栏收起
			MirrorStat:    false, //默认不统计镜像
		},
		Logo:          org.LogoUrl,
		AppDeployType: appDeployType,
		RemainDays:    GetRemainDays(info.PayEndTime),
	}

	// 付费过双11 活动 6.6折的 团队标识  2022-11-11
	isActivity11, errSys := CheckIsActivity11()
	if errSys != nil {
		log.Errorf("[GetOrgConfig] GetOrgBoById err: %v, orgId: %d", errSys, orgId)
		return nil, errSys
	}
	if isActivity11 && !CheckIsPrivateDeploy() {
		if info.PayLevel == consts.PayLevelDouble11Activity {
			res.IsPayActivity11 = consts.ActivityNotFinished
		} else {
			res.IsPayActivity11 = consts.ActivityFinished
		}
	} else {
		res.IsPayActivity11 = 0
	}

	if info.BasicShowSetting != "" {
		basicShowSettingVo := &orgvo.BasicShowSetting{}
		err := json.FromJson(info.BasicShowSetting, &basicShowSettingVo)
		if err != nil {
			log.Errorf("[GetOrgConfig] err: %v", err)
		} else {
			res.BasicShowSetting = basicShowSettingVo
		}
	}
	// 如果是私有化部署，则重置 payLevel
	// 私有化部署**只**由部署配置中的 Application.RunMode 决定
	if isPrivateDeploy {
		RedirectPayLevelInfoForPrivate(info)
		res.PayLevel = info.PayLevel
		// 私有化部署暂时没有到期时间
		res.PayEndTime = types.Time(info.PayEndTime)
	} else {
		// 如果不是私有化部署，则修正 payLevel 值
		//if res.PayLevel >= consts.PayLevelFlagship {
		//	res.PayLevel = consts.PayLevelFlagship
		//}
	}

	//org, orgErr := GetOrgBoById(orgId)
	//if orgErr != nil {
	//	return nil, orgErr
	//}
	remarkStr := org.Remark
	remarkObj := orgvo.OrgRemarkConfigType{}
	_ = json.FromJson(remarkStr, &remarkObj)
	res.SummaryAppID = fmt.Sprintf("%d", remarkObj.OrgSummaryTableAppId)

	count, err := GetOrgMemberCount(orgId)
	if err != nil {
		log.Errorf("[GetOrgConfig] err: %v, orgId: %d", err, orgId)
		return nil, err
	}

	// 加上outOrgId
	baseOrgInfo, err := GetBaseOrgInfo(orgId)
	if err != nil {
		log.Errorf("[GetOrgConfig] GetBaseOrgInfo err: %v, orgId: %d", err, orgId)
		return nil, err
	}
	res.OutOrgId = baseOrgInfo.OutOrgId

	res.OrgMemberNumber = int64(count)

	// 处理一下付费版和试用版
	res.PayLevelTrue = res.PayLevel
	if businees.CheckIsPaidVer(res.PayLevel) {
		// 检查是否是试用订单，仅适用于飞书
		// checkResp := orderfacade.CheckFsOrgIsInTrial(ordervo.CheckFsOrgIsInTrialReq{
		// 	OrgId: orgId,
		// })
		// if checkResp.Failure() {
		// 	log.Errorf("[GetOrgConfig] CheckFsOrgIsInTrial err: %v, orgId: %d", checkResp.Error(), orgId)
		// 	return res, checkResp.Error()
		// }
		// if checkResp.IsTrue {
		// 	res.PayLevelTrue = 3
		// }

		// 时间小于16天，判断为试用版；只有试用订单，会导致 org config 中 payStartTime 和 payEndTime 间隔 15 天。
		if info.PayStartTime.AddDate(0, 0, 16).After(info.PayEndTime) {
			res.PayLevelTrue = 3
			res.IsEvaluate = true
		}
	}

	// 判断是否是灰度企业
	res.IsGrayLevel = JudgeGrayLevelOrg(orgId)

	// 查询该组织支持的功能信息
	functionsResp := orderfacade.GetFunctionByLevel(ordervo.FunctionReq{Level: int64(res.PayLevel)})
	if functionsResp.Failure() {
		log.Errorf("[PersonalInfo] GetFunctionByLevel err: %v, level: %d, orgId: %d", functionsResp.Error(), res.PayLevel, orgId)
		return nil, functionsResp.Error()
	}
	functions := make([]orgvo.FunctionLimitObj, 0, len(functionsResp.Data))
	copyer.Copy(functionsResp.Data, &functions)
	res.Functions = functions

	return res, nil
}

func GetOrgConfig(orgId int64) (*bo.OrgConfigBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheOrgPayFunction, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	infoJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	res := &bo.OrgConfigBo{}
	if infoJson != "" {
		err = json.FromJson(infoJson, res)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		orgConfig := &po.PpmOrcConfig{}
		err := mysql.SelectOneByCond(consts.TableOrgConfig, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
		}, orgConfig)
		if err != nil {
			if err == db.ErrNoMoreRows {
				return nil, errs.OrgNotExist
			} else {
				log.Error(err)
				return nil, errs.MysqlOperateError
			}
		}
		_ = copyer.Copy(orgConfig, res)
		if CheckIsPrivateDeploy() {
			RedirectPayLevelInfoForPrivate(res)
		}
		err = cache.SetEx(key, json.ToJsonIgnoreError(res), consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
	}

	// 有效期判断
	if businees.CheckIsPaidVer(res.PayLevel) &&
		res.PayEndTime.Before(time.Now()) {
		// 过了有效期需要降级
		_, err := mysql.UpdateSmartWithCond(consts.TableOrgConfig, db.Cond{
			consts.TcId: res.Id,
		}, mysql.Upd{
			consts.TcPayLevel: consts.PayLevelStandard,
		})
		if err != nil {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}

		res.PayLevel = consts.PayLevelStandard
		err = cache.SetEx(key, json.ToJsonIgnoreError(res), consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}

		clearErr := ClearOrgConfig(orgId)
		if clearErr != nil {
			log.Error(clearErr)
		}

		// 上报事件
		asyn.Execute(func() {
			orgConfig, err := GetOrgConfigRich(orgId)
			if err != nil {
				return
			}
			orgEvent := &commonvo.OrgEvent{}
			orgEvent.OrgId = orgId
			orgEvent.New = orgConfig

			openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
			openTraceIdStr := cast.ToString(openTraceId)

			report.ReportOrgEvent(msgPb.EventType_OrgConfigUpdated, openTraceIdStr, orgEvent, true)
		})
	}
	return res, nil
}

func GetWhiteListVipOrg(orgId int64) bool {
	key := sconsts.CacheVipOrg
	value, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return false
	}
	if value == "" {
		return false
	}

	result := &[]int64{}
	err = json.FromJson(value, result)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return false
	}

	ok, _ := slice.Contain(*result, orgId)

	return ok
}

func JudgeGrayLevelOrg(orgId int64) bool {
	key := sconsts.CacheGrayLevelOrg
	value, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return false
	}
	if value == "" {
		return false
	}

	result := &[]int64{}
	err = json.FromJson(value, result)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return false
	}

	ok, _ := slice.Contain(*result, orgId)

	return ok
}

func GetOrgAppTicketFromCache(key string) (*bo.AppTicketBo, errs.SystemErrorInfo) {
	baseOrgInfoJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	baseOrgInfo := &bo.AppTicketBo{}
	if baseOrgInfoJson != "" {
		err := json.FromJson(baseOrgInfoJson, baseOrgInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		return nil, errs.AppTicketNotAllocated
	}

	return baseOrgInfo, nil
}

func GetOrgAppTicket(orgId int64) (*bo.AppTicketBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheOrgAppTicket, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	baseOrgInfo, err := GetOrgAppTicketFromCache(key)
	if err != nil {
		//防止缓存击穿
		lockId := uuid.NewUuid()
		suc, err := cache.TryGetDistributedLock(consts.GetOrgAppTicketLock, lockId)
		if err != nil {
			log.Error("获取锁异常")
			return nil, errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
		}
		if suc {
			//释放锁
			defer func() {
				if _, e := cache.ReleaseDistributedLock(consts.GetOrgAppTicketLock, lockId); e != nil {
					log.Error(e)
				}
			}()

			baseOrgInfo, err := GetOrgAppTicketFromCache(key)
			if err == nil {
				return baseOrgInfo, nil
			}

			info := &po.PpmOrgSecret{}
			selectErr := mysql.SelectOneByCond(consts.TableOrgSecret, db.Cond{
				consts.TcOrgId:    orgId,
				consts.TcIsDelete: consts.AppIsNoDelete,
			}, info)
			if selectErr != nil {
				if selectErr == db.ErrNoMoreRows {
					//如果没有就添加
					id, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrgSecret)
					if idErr != nil {
						log.Error(idErr)
						return nil, idErr
					}

					key := uuid.NewUuid()
					secret := uuid.NewUuid()
					err1 := mysql.Insert(&po.PpmOrgSecret{
						Id:     id,
						OrgId:  orgId,
						Key:    key,
						Secret: secret,
					})
					if err1 != nil {
						log.Error(err1)
						return nil, errs.MysqlOperateError
					}

					baseOrgInfo = &bo.AppTicketBo{
						AppId:     key,
						AppSecret: secret,
					}
				} else {
					log.Error(selectErr)
					return nil, errs.MysqlOperateError
				}
			} else {
				baseOrgInfo = &bo.AppTicketBo{
					AppId:     info.Key,
					AppSecret: info.Secret,
				}
			}

			cacheErr := cache.SetEx(key, json.ToJsonIgnoreError(baseOrgInfo), consts.GetCacheBaseExpire())
			if cacheErr != nil {
				log.Error(cacheErr)
			}

			return baseOrgInfo, nil
		} else {
			return GetOrgAppTicketFromCache(key)
		}
	}

	return baseOrgInfo, nil
}
