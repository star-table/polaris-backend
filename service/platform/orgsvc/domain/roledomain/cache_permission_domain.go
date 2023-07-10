package domain

import (
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	bo2 "github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

var log = logger.GetDefaultLogger()

// TODO 切换lesscode-permission
func GetPermissionList() (*[]bo2.PermissionBo, errs.SystemErrorInfo) {
	key := sconsts.CachePermissionList
	listJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	bo := &[]bo2.PermissionBo{}
	if listJson != "" {
		err := json.FromJson(listJson, bo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		po := &[]po.PpmRolPermission{}
		selectErr := mysql.SelectAllByCond(consts.TablePermission, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   consts.AppStatusEnable,
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

// TODO 切换lesscode-permission
func GetPermissionByType(permissionType int) ([]bo2.PermissionBo, errs.SystemErrorInfo) {
	res := []bo2.PermissionBo{}
	list, err := GetPermissionList()
	if err != nil {
		log.Error(err)
		return res, err
	}
	for _, v := range *list {
		if v.IsShow == consts.AppShowEnable && v.Type == permissionType {
			res = append(res, v)
		}
	}

	return res, nil
}

// TODO 切换lesscode-permission
//获取项目权限项（只取部分）
func GetProjectPermission(orgId, projectId int64) ([]bo2.PermissionBo, errs.SystemErrorInfo) {
	res := []bo2.PermissionBo{}
	list, err := GetPermissionList()
	if err != nil {
		log.Error(err)
		return res, err
	}
	selectedPermission := []string{
		//任务
		consts.PermissionProIssue4,
		//任务栏
		consts.PermissionOrgProjectObjectType,
		//文件管理
		consts.PermissionProFile,
		//标签管理
		consts.PermissionProTag,
		//附件管理
		consts.PermissionProAttachment,
		//项目成员管理
		consts.PermissionProMember,
		//项目相关
		consts.PermissionProConfig,
		//项目角色管理
		consts.PermissionProRole,
	}
	if projectId != 0 {
		//获取项目信息
		projectInfo := projectfacade.GetCacheProjectInfo(projectvo.GetCacheProjectInfoReqVo{
			ProjectId:projectId,
			OrgId:         orgId,
		})
		if projectInfo.Failure() {
			log.Error(projectInfo.Error())
			return nil, projectInfo.Error()
		}
		//判断项目是否是敏捷项目
		if projectInfo.ProjectCacheBo.ProjectType == consts.ProjectTypeAgileId {
			selectedPermission = append(selectedPermission, consts.PermissionProIteration)
		}
	}

	for _, v := range *list {
		if v.IsShow == consts.AppShowEnable && v.Type == consts.PermissionTypePro {
			if ok, _ := slice.Contain(selectedPermission, v.LangCode); ok {
				res = append(res, v)
			}
		}
	}

	return res, nil
}

// TODO 切换lesscode-permission
func GetPermissionById(permissionId int64) (bo2.PermissionBo, errs.SystemErrorInfo) {
	res := bo2.PermissionBo{}
	list, err := GetPermissionList()
	if err != nil {
		log.Error(err)
		return res, err
	}

	for _, bo := range *list {
		if bo.Id == permissionId {
			return bo, err
		}
	}

	return res, errs.PermissionNotExist
}
