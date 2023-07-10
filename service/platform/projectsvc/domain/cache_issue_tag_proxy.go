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

func IssueTagStatByCache(orgId, projectId int64) ([]bo.IssueTagStatBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectTagInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}
	resJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.RedisOperateError
	}

	resBo := &[]bo.IssueTagStatBo{}
	if resJson != "" {
		err := json.FromJson(resJson, resBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return *resBo, nil
	} else {
		conn, err := mysql.GetConnect()
		if err != nil {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}
		resPo := &[]po.IssueTagStat{}
		selectErr := conn.Select(db.Raw("count(i.issue_id) AS total, t.id as tag_id")).From("ppm_pri_tag t").
			LeftJoin("ppm_pri_issue_tag i").On("t.id = i.tag_id").And("i.is_delete = 2").
			Where(db.Cond{
				"t." + consts.TcIsDelete:  consts.AppIsNoDelete,
				"t." + consts.TcOrgId:     orgId,
				"t." + consts.TcProjectId: projectId,
			}).GroupBy("t." + consts.TcId).All(resPo)
		if selectErr != nil {
			log.Error(selectErr)
			return nil, errs.MysqlOperateError
		}

		err = cache.SetEx(key, json.ToJsonIgnoreError(resPo), 120)
		if err != nil {
			log.Error(err)
			return nil, errs.CacheProxyError
		}

		_ = copyer.Copy(resPo, resBo)

		return *resBo, nil
	}
}
