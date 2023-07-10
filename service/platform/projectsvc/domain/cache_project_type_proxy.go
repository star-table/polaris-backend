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

func GetProjectTypeList(orgId int64) (*[]bo.ProjectTypeBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectTypeList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	projectTypeListPo := &[]po.PpmPrsProjectType{}
	projectTypeListBo := &[]bo.ProjectTypeBo{}
	projectTypeListJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if projectTypeListJson != "" {

		err = json.FromJson(projectTypeListJson, projectTypeListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return projectTypeListBo, nil
	} else {
		err := mysql.SelectAllByCond(consts.TableProjectType, db.Cond{
			consts.TcOrgId:    db.In([]int64{orgId, 0}),
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   consts.AppStatusEnable,
		}, projectTypeListPo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(projectTypeListPo, projectTypeListBo)
		projectTypeListJson, err := json.ToJson(projectTypeListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, projectTypeListJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return projectTypeListBo, nil
	}
}

func GetProjectTypeByLangCode(orgId int64, langCode string) (*bo.ProjectTypeBo, errs.SystemErrorInfo) {
	list, err := GetProjectTypeList(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, projectType := range *list {
		if projectType.LangCode == langCode {
			return &projectType, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.ProjectTypeNotExist)
}
