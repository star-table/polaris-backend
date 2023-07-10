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
	sconsts "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetMenu(orgId int64, appId int64) (*bo.GetMenuInfoBo, errs.SystemErrorInfo) {
	key, err := util.ParseCacheKey(sconsts.CacheProjectMenuConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
		consts.CacheKeyAppIdConstName: appId,
	})
	if err != nil {
		return nil, err
	}
	menuJson, err2 := cache.Get(key)
	if err2 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err2)
	}
	menuInfoBo := bo.GetMenuInfoBo{}
	if menuJson != "" {
		err3 := json.FromJson(menuJson, &menuInfoBo)
		if err3 != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		menuPo := po.PpmProProjectMenuConfig{}
		err4 := mysql.SelectOneByCond(menuPo.TableName(), db.Cond{
			consts.TcOrgId: orgId,
			consts.TcAppId: appId,
		}, &menuPo)
		if err4 != nil && err4 != db.ErrNoMoreRows {
			log.Errorf("查询菜单错误, orgId: %v, appId: %v, 错误: %v", orgId, appId, err4)
			return nil, errs.MysqlOperateError
		}
		menuConfigMap := map[string]interface{}{}
		if menuPo.Config == "" {
			menuInfoBo.Config = menuConfigMap
		}else {
			err44 := json.FromJson(menuPo.Config, &menuConfigMap)
			if err44 != nil {
				return nil, errs.JSONConvertError
			}
		}
		menuInfoBo.OrgId = orgId
		menuInfoBo.AppId = appId
		menuInfoBo.Config = menuConfigMap
		menuJson = json.ToJsonIgnoreError(&menuInfoBo)
		err5 := cache.SetEx(key, menuJson, consts.CacheBaseExpire)
		if err5 != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err5)
		}
	}

	return &menuInfoBo, nil
}

func SaveMenu(orgId int64, appId int64, config string) (int64, errs.SystemErrorInfo) {
	// 更新需要查询数据库，有就更新，没有就插入新数据
	pos := po.PpmProProjectMenuConfig{}
	err2 := mysql.SelectOneByCond(consts.TableProjectMenu, db.Cond{
		consts.TcOrgId: orgId,
		consts.TcAppId: appId,
	}, &pos)
	if err2 != nil && err2 != db.ErrNoMoreRows {
		return 0, errs.MysqlOperateError
	}
	if pos.Config == "" {
		menuConfig := bo.ProjectMenuConfigBo{
			OrgId:  orgId,
			AppId:  appId,
			Config: config,
		}
		err := insertProjectMenuConfig(&menuConfig)
		if err != nil {
			log.Errorf("[SaveMenu] 保存菜单失败, orgId: %v, appId: %v, 错误: %v", orgId, appId, err)
			return 0, err
		}
	} else {
		// 需要更新
		update, err := updateProjectMenuConfig(orgId, appId, config)
		if err != nil && update != true {
			log.Errorf("[updateProjectMenuConfig] 更新失败, orgId: %v, appId: %v, 错误: %v", orgId, appId, err)
			return 0, err
		}
	}
	clearErr := clearMenuCache(orgId, appId)
	if clearErr != nil {
		log.Errorf("[clearMenuCache] 清除缓存失败, orgId: %v, appId: %v, 错误: %v", orgId, appId, clearErr)
	}
	return appId, nil
}

func clearMenuCache(orgId int64, appId int64) errs.SystemErrorInfo {
	key, err := util.ParseCacheKey(sconsts.CacheProjectMenuConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
		consts.CacheKeyAppIdConstName: appId,
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

func insertProjectMenuConfig(req *bo.ProjectMenuConfigBo) errs.SystemErrorInfo {
	insert := po.PpmProProjectMenuConfig{}
	err1 := copyer.Copy(req, &insert)
	if err1 != nil {
		log.Errorf("[insertProjectMenuConfig] 拷贝错误: %v", err1)
		return errs.BuildSystemErrorInfo(errs.SystemError, err1)
	}
	err2 := mysql.Insert(&insert)
	if err2 != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}
	return nil
}

func updateProjectMenuConfig(orgId int64, appId int64, config string) (bool, errs.SystemErrorInfo) {
	upds := mysql.Upd{
		"config": config,
	}
	_, err := mysql.UpdateSmartWithCond(consts.TableProjectMenu, db.Cond{
		consts.TcOrgId: orgId,
		consts.TcAppId: appId,
	}, upds)
	if err != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return true, nil
}
