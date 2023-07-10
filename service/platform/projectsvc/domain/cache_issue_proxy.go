package domain

import (
	"github.com/spf13/cast"
	"strconv"
	"time"

	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	sconsts "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
)

// 检查**任务审批**是否需要催办
func NeedUrgeIssue(orgId, projectId, issueId int64) (bool, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheUrgeIssueAudit, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
		consts.CacheKeyIssueIdConstName:   issueId,
	})
	if err5 != nil {
		log.Error(err5)
		return false, err5
	}
	infoJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return false, errs.CacheProxyError
	}

	if infoJson == "" {
		return true, nil
	} else {
		return false, nil
	}
}

// 检查**任务**是否需要催办
func NeedUrgeForOnlyIssue(orgId, projectId, issueId int64) (bool, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheUrgeIssue, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
		consts.CacheKeyIssueIdConstName:   issueId,
	})
	if err5 != nil {
		log.Error(err5)
		return false, err5
	}
	infoJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return false, errs.CacheProxyError
	}

	if infoJson == "" {
		return true, nil
	} else {
		return false, nil
	}
}

func UpdateUrgeIssueTime(orgId, projectId, issueId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheUrgeIssueAudit, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
		consts.CacheKeyIssueIdConstName:   issueId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	err := cache.SetEx(key, strconv.FormatInt(time.Now().Unix(), 10), 60*10)
	if err != nil {
		log.Error(err)
		return errs.CacheProxyError
	}

	return nil
}

// **任务**催办时间记录/更新
func UpdateUrgeOnlyIssueTime(orgId, projectId, issueId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheUrgeIssue, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
		consts.CacheKeyIssueIdConstName:   issueId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	err := cache.SetEx(key, strconv.FormatInt(time.Now().Unix(), 10), 60*10)
	if err != nil {
		log.Error(err)
		return errs.CacheProxyError
	}

	return nil
}

func GetLastUrgeIssue(orgId, projectId, issueId int64, keyTpl string) (int64, errs.SystemErrorInfo) {
	if len(keyTpl) < 1 {
		keyTpl = sconsts.CacheUrgeIssueAudit
	}
	key, err5 := util.ParseCacheKey(keyTpl, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
		consts.CacheKeyIssueIdConstName:   issueId,
	})
	if err5 != nil {
		log.Error(err5)
		return 0, err5
	}
	infoJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return 0, errs.CacheProxyError
	}

	if infoJson == "" {
		return int64(0), nil
	} else {
		now, err := strconv.ParseInt(infoJson, 10, 64)
		if err != nil {
			log.Error(err)
			return 0, errs.JSONConvertError
		}
		return now, nil
	}
}

func CheckIssuesAllowUrge(orgId, projectId int64, issueIds []int64, keyTpl string) ([]bool, errs.SystemErrorInfo) {
	if len(keyTpl) < 1 {
		keyTpl = sconsts.CacheUrgeIssueAudit
	}
	var keys []interface{}
	for _, issueId := range issueIds {
		key, err := util.ParseCacheKey(keyTpl, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyIssueIdConstName:   issueId,
		})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		keys = append(keys, key)
	}

	infoJsons, err := cache.MGetFull(keys...)
	if err != nil {
		log.Error(err)
		return nil, errs.CacheProxyError
	}

	if len(infoJsons) == 0 {
		return nil, nil
	}

	var rs []bool
	for _, i := range infoJsons {
		if i == nil {
			rs = append(rs, true)
		} else {
			rs = append(rs, false)
		}
	}
	return rs, nil
}

func UpdateIssuesUrgeTime(orgId, projectId int64, issueIds []int64, keyTpl string) errs.SystemErrorInfo {
	if len(keyTpl) < 1 {
		keyTpl = sconsts.CacheUrgeIssueAudit
	}
	var keys []string
	for _, issueId := range issueIds {
		key, err := util.ParseCacheKey(keyTpl, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyIssueIdConstName:   issueId,
		})
		if err != nil {
			log.Error(err)
			return err
		}
		keys = append(keys, key)
	}

	for _, key := range keys {
		err := cache.SetEx(key, strconv.FormatInt(time.Now().Unix(), 10), 60*10)
		if err != nil {
			log.Error(err)
			return errs.CacheProxyError
		}
	}
	return nil
}

func SetDeletedIssuesNum(orgId, issueNums int64) errs.SystemErrorInfo {
	_, err := cache.HINCRBY(sconsts.CacheDeletedIssueNums, strconv.FormatInt(orgId, 10), issueNums)
	if err != nil {
		log.Errorf("[SetDeletedIssuesNum] err:%v", err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}
	return nil
}

func GetDeletedIssueNum(orgId int64) (int64, errs.SystemErrorInfo) {
	res, err := cache.HMGet(sconsts.CacheDeletedIssueNums, strconv.FormatInt(orgId, 10))
	if err != nil {
		log.Errorf("[GetDeletedIssueNum] err:%v", err)
		return 0, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}
	if v, ok := res[strconv.FormatInt(orgId, 10)]; ok {
		return cast.ToInt64(&v), nil
	}
	return 0, nil
}
