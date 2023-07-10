package domain

import (
	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	sconsts "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
)

func DeleteIssueSourceListCache(orgId int64) errs.SystemErrorInfo {

	key, err := util.ParseCacheKey(sconsts.CacheIssueSourceList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})

	if err != nil {
		log.Error(err)
		return err
	}

	_, err1 := cache.Del(key)

	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	return nil
}
