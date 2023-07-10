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

func GetIssuePropertyById(orgId, sourceId int64, projectId int64) (*bo.IssuePropertyBo, errs.SystemErrorInfo) {
	list, err1 := GetIssuePropertyList(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	for _, source := range list {
		if source.Id == sourceId {
			return &source, nil
		}
	}
	return nil, errs.PropertyIdNotExist
}

func GetIssuePropertyListByProjectObjectTypeId(orgId, projectTypeId int64, projectId int64) ([]bo.IssuePropertyBo, errs.SystemErrorInfo) {
	list, err1 := GetIssuePropertyList(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	result := make([]bo.IssuePropertyBo, 0)
	for _, source := range list {
		if source.ProjectObjectTypeId == projectTypeId {
			result = append(result, source)
		}
	}
	return result, nil
}

func GetIssuePropertyList(orgId int64, projectId int64) ([]bo.IssuePropertyBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheIssuePropertyList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	issuePropertyListBo := &[]bo.IssuePropertyBo{}
	issuePropertyListPo := &[]po.PpmPrsIssueProperty{}
	issuePropertyListJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if issuePropertyListJson != "" {

		err = json.FromJson(issuePropertyListJson, issuePropertyListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return *issuePropertyListBo, nil
	} else {
		err := mysql.SelectAllByCond(consts.TableIssueProperty, db.Cond{
			consts.TcOrgId:     orgId,
			consts.TcProjectId: projectId,
			consts.TcIsDelete:  consts.AppIsNoDelete,
			consts.TcStatus:    consts.AppStatusEnable,
		}, issuePropertyListPo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(issuePropertyListPo, issuePropertyListBo)
		issueSourceListJson, err := json.ToJson(issuePropertyListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, issueSourceListJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return *issuePropertyListBo, nil
	}
}

func GetNewIssuePropertyId(orgId int64, oldProjectId, oldId int64, newProjectId, projectObjectTypeId int64) int64 {
	info, err := GetIssuePropertyById(orgId, oldId, oldProjectId)
	if err != nil {
		log.Error(err)
		return 0
	}

	newInfos, err := GetIssuePropertyListByProjectObjectTypeId(orgId, projectObjectTypeId, newProjectId)
	if err != nil {
		log.Error(err)
		return 0
	}

	for _, newInfo := range newInfos {
		if info.Name == newInfo.Name {
			return newInfo.Id
		}
	}

	return 0
}

func GetIssuePropertyBos(orgId, projectId int64, projectObjectTypeId int64) ([]bo.IssuePropertyBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmPrsIssueProperty{}
	cond := db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
	}
	if projectObjectTypeId > 0 {
		cond[consts.TcProjectObjectTypeId] = projectObjectTypeId
	}
	err := mysql.SelectAllByCond(consts.TableIssueProperty, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.IssuePropertyBo{}
	_ = copyer.Copy(pos, bos)

	return *bos, nil
}
