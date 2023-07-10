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
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

//组织角色组缓存
func GetGroupRoleList(orgId int64) (*[]bo.RoleGroupBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheRoleGroupList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}
	listJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	bo := &[]bo.RoleGroupBo{}
	if listJson != "" {
		err := json.FromJson(listJson, bo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		po := &[]po.PpmRolRoleGroup{}
		selectErr := mysql.SelectAllByCond(consts.TableRoleGroup, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
		}, po)
		if selectErr != nil {
			log.Error(selectErr)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
		}
		_ = copyer.Copy(po, bo)
		listJson, err = json.ToJson(bo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, listJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}

	return bo, nil
}

//删除用户角色列表缓存
func ClearUserRoleList(orgId, userId, projectId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheUserRoleListHash, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}

	_, err := cache.HDel(key, projectId)

	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}


//删除用户角色列表缓存
func ClearDepartmentRoleList(orgId, departmentId, projectId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheDepartmentRoleListHash, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyDepartmentIdConstName: departmentId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}

	_, err := cache.HDel(key, projectId)

	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}
