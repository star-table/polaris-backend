package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

func SetLabConfig(orgId, userId int64, input orgvo.SetLabReq) errs.SystemErrorInfo {
	// 更新需要查询数据库，有就更新，没有就插入新数据

	configJson, err3 := json.ToJson(input)
	if err3 != nil {
		return errs.BuildSystemErrorInfoWithMessage(errs.JSONConvertError, err3.Error())
	}

	pos := po.PpmOrgLabConfig{}
	err1 := mysql.SelectOneByCond(pos.TableName(), db.Cond{
		consts.TcOrgId: orgId,
	}, &pos)
	if err1 != nil && err1 != db.ErrNoMoreRows {
		return errs.MysqlOperateError
	}

	if pos.Config == "" {
		// 插入数据
		labConfig := bo.LabConfigBo{
			OrgId:  orgId,
			Config: configJson,
		}
		err := insertLabConfig(&labConfig)
		if err != nil {
			log.Errorf("[SetLabConfig] 保存labConfig失败, orgId: %d, 错误: %v", orgId, err)
			return nil
		}
	} else {
		// 需要更新

		oldSetting := &orgvo.SetLabReq{}
		err1 = json.FromJson(pos.Config, oldSetting)
		if err1 != nil {
			log.Errorf("[updateLabConfig] 解析老配置失败, orgId: %d, err: %v", orgId, err1)
			return errs.JSONConvertError
		}
		// 自动化开关 开->关
		//if oldSetting.AutomationSwitch && !input.AutomationSwitch {
		//	resp := automationfacade.GlobalSwitchOff(orgId, userId)
		//	if resp.Failure() {
		//		log.Errorf("[updateLabConfig] automation GlobalSwitchOff failed, orgId: %d, err: %v", orgId, resp.Error())
		//		return resp.Error()
		//	}
		//}

		update, err := updateLabConfig(orgId, configJson)
		if err != nil && update != true {
			log.Errorf("[updateLabConfig] 更新失败, orgId: %d, err: %v", orgId, err)
			return nil
		}
	}
	// 清除缓存
	clearErr := clearLabConfig(orgId)
	if clearErr != nil {
		log.Errorf("[SetLabConfig] 清除缓存失败, orgId: %d, 错误: %v", orgId, clearErr)
	}
	return nil
}

func GetLabConfig(orgId int64) (*bo.GetLabConfigBo, errs.SystemErrorInfo) {
	key, err := util.ParseCacheKey(sconsts.CacheLabConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err != nil {
		return nil, err
	}
	labConfigJson, err2 := cache.Get(key)
	if err2 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err2)
	}
	labConfigBo := bo.GetLabConfigBo{}
	if labConfigJson != "" {
		err3 := json.FromJson(labConfigJson, &labConfigBo)
		if err3 != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		labConfigPo := po.PpmOrgLabConfig{}
		err4 := mysql.SelectOneByCond(labConfigPo.TableName(), db.Cond{
			consts.TcOrgId: orgId,
		}, &labConfigPo)
		if err4 != nil && err4 != db.ErrNoMoreRows {
			log.Errorf("[GetLabConfig]查询实验室开关错误, orgId: %d, err: %v", orgId, err4)
			return nil, errs.MysqlOperateError
		}
		if labConfigPo.Config == "" {
			// 这里可以设置 开关关闭的缓存
			labConfigBo.WorkBenchShow = false
			labConfigBo.ProOverview = false
			labConfigBo.SideBarShow = false
			labConfigBo.EmptyApp = false
			//labConfigBo.AutomationSwitch = false
			//labConfigBo.SubmitButton = false
			labConfigPo.Config = json.ToJsonIgnoreError(labConfigBo)
			cacheSetErr := cache.SetEx(key, labConfigPo.Config, consts.CacheBaseExpire)
			if cacheSetErr != nil {
				return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, cacheSetErr)
			}
			return &bo.GetLabConfigBo{
				WorkBenchShow: false,
				ProOverview:   false,
				SideBarShow:   false,
				EmptyApp:      false,
				//AutomationSwitch: false,
				//SubmitButton:  false,
			}, nil
		}
		err6 := json.FromJson(labConfigPo.Config, &labConfigBo)
		if err6 != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}

		err5 := cache.SetEx(key, labConfigPo.Config, consts.CacheBaseExpire)
		if err5 != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err5)
		}
	}

	return &labConfigBo, nil
}

func insertLabConfig(req *bo.LabConfigBo) errs.SystemErrorInfo {
	insert := po.PpmOrgLabConfig{}
	err1 := copyer.Copy(req, &insert)
	if err1 != nil {
		if err1 != nil {
			log.Errorf("[insertLabConfig] 拷贝错误: %v", err1)
			return errs.BuildSystemErrorInfo(errs.SystemError, err1)
		}
	}
	err2 := mysql.Insert(&insert)
	if err2 != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}
	return nil
}

func updateLabConfig(orgId int64, config string) (bool, errs.SystemErrorInfo) {
	upds := mysql.Upd{
		consts.TcLabConfig: config,
	}
	_, err := mysql.UpdateSmartWithCond(consts.TableLabConfig, db.Cond{
		consts.TcOrgId: orgId,
	}, upds)
	if err != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return true, nil
}

func clearLabConfig(orgId int64) errs.SystemErrorInfo {
	key, err := util.ParseCacheKey(sconsts.CacheLabConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err != nil {
		return err
	}
	_, err2 := cache.Del(key)
	if err2 != nil {
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err2)
	}
	return nil
}
