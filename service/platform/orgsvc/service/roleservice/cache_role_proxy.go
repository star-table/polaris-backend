package service

import (
	"strconv"
	"strings"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/bo"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	domain "github.com/star-table/polaris-backend/service/platform/orgsvc/domain/roledomain"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

//查询系统定义角色
func GetRoleByLangCode(orgId int64, langCode string) (*bo.RoleBo, errs.SystemErrorInfo) {
	if langCode == consts.BlankString {
		log.Error("查询langCode不能为空")
		return nil, errs.BaseDomainError
	}
	key, err5 := util.ParseCacheKey(sconsts.CacheRoleListHash, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}
	roleListJson, err := cache.HGet(key, langCode)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	roleBo := &bo.RoleBo{}
	if roleListJson != "" {
		err := json.FromJson(roleListJson, roleBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		rolePo := &po.PpmRolRole{}
		err := mysql.SelectOneByCond(consts.TableRole, db.Cond{
			consts.TcOrgId:    db.In([]int64{0, orgId}),
			consts.TcStatus:   consts.AppStatusEnable,
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcLangCode: langCode,
		}, rolePo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(rolePo, roleBo)
		roleListJson, err = json.ToJson(roleBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.HSet(key, langCode, roleListJson)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}

	return roleBo, nil
}

func GetRoleOperationByCode(code string) (*bo.OperationBo, errs.SystemErrorInfo) {
	roleOperationList, err := GetRoleOperationList()
	if err != nil {
		return nil, err
	}
	for _, roleOperation := range *roleOperationList {
		if strings.EqualFold(roleOperation.Code, code) {
			return &roleOperation, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.RoleOperationNotExist)
}

func GetRoleOperationList() (*[]bo.OperationBo, errs.SystemErrorInfo) {
	key := sconsts.CacheRoleOperationList

	roleOperationListJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if roleOperationListJson != "" {
		roleOperationList := &[]bo.OperationBo{}
		err := json.FromJson(roleOperationListJson, roleOperationList)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return roleOperationList, nil
	} else {
		roleOperationListPo := &[]po.PpmRolOperation{}
		roleOperationListBo := &[]bo.OperationBo{}
		err := mysql.SelectAllByCond(consts.TableRoleOperation, db.Cond{
			consts.TcStatus:   consts.AppStatusEnable,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, roleOperationListPo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(roleOperationListPo, roleOperationListBo)
		roleOperationListJson, err = json.ToJson(roleOperationListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, roleOperationListJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return roleOperationListBo, nil
	}
}

func GetRolePermissionOperationListByPath(orgId int64, roleIds []int64, path string, projectId int64) (*[]bo.RolePermissionOperationBo, errs.SystemErrorInfo) {
	rolePermissionOperationList := &[]bo.RolePermissionOperationBo{}

	roleIds = slice.SliceUniqueInt64(roleIds)
	for _, roleId := range roleIds {
		list, err := GetRolePermissionOperationList(orgId, roleId, projectId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for _, operation := range *list {
			if strings.EqualFold(operation.PermissionPath, path) || strings.HasPrefix(path, operation.PermissionPath+"/") {
				*rolePermissionOperationList = append(*rolePermissionOperationList, operation)
			}
		}
	}
	return rolePermissionOperationList, nil
}

func GetRolePermissionOperationList(orgId, roleId int64, projectId int64) (*[]bo.RolePermissionOperationBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheRolePermissionOperationList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
		consts.CacheKeyRoleIdConstName:    roleId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	rolePermissionOperationListJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if rolePermissionOperationListJson != "" {
		rolePermissionOperationList := &[]bo.RolePermissionOperationBo{}
		err := json.FromJson(rolePermissionOperationListJson, rolePermissionOperationList)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return rolePermissionOperationList, nil
	} else {
		rolePermissionOperationListPo := &[]po.PpmRolRolePermissionOperation{}
		rolePermissionOperationListBo := &[]bo.RolePermissionOperationBo{}
		cond := db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcRoleId:   roleId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}
		if projectId != 0 {
			cond[consts.TcProjectId] = db.In([]int64{0, projectId})
		} else {
			cond[consts.TcProjectId] = 0
		}

		err := mysql.SelectAllByCond(consts.TableRolePermissionOperation, cond, rolePermissionOperationListPo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}

		dealedList := getNeedRolePermissionOperationList(*rolePermissionOperationListPo, projectId)
		_ = copyer.Copy(dealedList, rolePermissionOperationListBo)

		rolePermissionOperationListJson, err = json.ToJson(rolePermissionOperationListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, rolePermissionOperationListJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return rolePermissionOperationListBo, nil
	}
}

//获取需要的权限列表（如果有项目id的就用项目里的，没有就用通用的）
func getNeedRolePermissionOperationList(list []po.PpmRolRolePermissionOperation, projectId int64) []po.PpmRolRolePermissionOperation {
	newArr := []po.PpmRolRolePermissionOperation{}
	if projectId != 0 {
		for _, operation := range list {
			if operation.ProjectId == projectId {
				newArr = append(newArr, operation)
			}
		}
	}

	if len(newArr) == 0 {
		newArr = list
	}

	return newArr
}

func ClearRolePermissionOperationList(orgId, roleId, projectId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheRolePermissionOperationList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
		consts.CacheKeyRoleIdConstName:    roleId,
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

func GetUserRoleList(orgId, userId int64, projectId int64) (*[]bo.RoleUserBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheUserRoleListHash, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	roleUserListJson, err := cache.HGet(key, strconv.FormatInt(projectId, 10))
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if roleUserListJson != "" {
		roleUserList := &[]bo.RoleUserBo{}
		err := json.FromJson(roleUserListJson, roleUserList)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return roleUserList, nil
	} else {
		roleUserListPo := &[]po.PpmRolRoleUser{}
		roleUserListBo := &[]bo.RoleUserBo{}
		err := mysql.SelectAllByCond(consts.TableRoleUser, db.Cond{
			consts.TcOrgId:     orgId,
			consts.TcUserId:    userId,
			consts.TcProjectId: projectId,
			consts.TcIsDelete:  consts.AppIsNoDelete,
		}, roleUserListPo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(roleUserListPo, roleUserListBo)
		roleUserListJson, err = json.ToJson(roleUserListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.HSet(key, strconv.FormatInt(projectId, 10), roleUserListJson)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return roleUserListBo, nil
	}
}

func GetUserRoleListByProjectId(orgId, userId, projectId int64) ([]int64, errs.SystemErrorInfo) {
	userRoleList, err := GetUserRoleList(orgId, userId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	roleIds := []int64{}
	for _, userBo := range *userRoleList {
		roleIds = append(roleIds, userBo.RoleId)
	}
	return slice.SliceUniqueInt64(roleIds), nil
}

func GetCompensatoryRolePermissionPaths(orgId int64) ([]string, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheCompensatoryRolePermissionPathList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}
	cacheJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if cacheJson != "" {
		resultPaths := &[]string{}
		err := json.FromJson(cacheJson, resultPaths)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return *resultPaths, nil
	} else {
		paths := make([]string, 0)
		pathMap := map[string]string{}
		if len(sconsts.RolePermissionOperationDefineMap) == 0 {
			return paths, nil
		}
		for k, _ := range sconsts.RolePermissionOperationDefineMap {
			v := strings.ReplaceAll(k, "{org_id}", strconv.FormatInt(orgId, 10))
			paths = append(paths, v)
			pathMap[v] = k
		}
		rolePermissionOperationListPo := &[]po.PpmRolRolePermissionOperation{}

		//这边将已删除的也查出来，防止已经补偿但是用户主动删除
		err := mysql.SelectAllByCond(consts.TableRolePermissionOperation, db.Cond{
			consts.TcOrgId:          orgId,
			consts.TcPermissionPath: db.In(paths),
		}, rolePermissionOperationListPo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		resultPaths := make([]string, 0)
		for _, po := range *rolePermissionOperationListPo {
			if path, ok := pathMap[po.PermissionPath]; ok {
				if exist, _ := slice.Contain(resultPaths, path); !exist {
					resultPaths = append(resultPaths, path)
				}
			}
		}
		cacheErr := cache.Set(key, json.ToJsonIgnoreError(resultPaths))
		if cacheErr != nil {
			log.Error(cacheErr)
			return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
		}
		return resultPaths, nil
	}
}

func SetCompensatoryRolePermissionPaths(orgId int64, paths []string) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheCompensatoryRolePermissionPathList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	cacheErr := cache.Set(key, json.ToJsonIgnoreError(paths))
	if cacheErr != nil {
		log.Error(cacheErr)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError)
	}
	return nil
}

func GetRoleListByGroup(orgId int64, langCode string, projectId int64) ([]bo.RoleBo, errs.SystemErrorInfo) {
	list, err := domain.GetGroupRoleList(orgId)
	if err != nil {
		return nil, err
	}
	var groupId int64
	for _, v := range *list {
		if v.LangCode == langCode {
			groupId = v.Id
			break
		}
	}
	if groupId == 0 {
		return nil, errs.OrgRoleGroupNotExist
	}

	roleList, err := domain.GetRoleList(orgId, projectId, []int64{}, true)
	if err != nil {
		return nil, err
	}

	roleBo := []bo.RoleBo{}
	for _, v := range roleList {
		if v.RoleGroupId == groupId {
			roleBo = append(roleBo, v)
		}
	}

	return roleBo, nil
}

func ClearUserRoleList(orgId int64, userIds []int64, projectId int64) errs.SystemErrorInfo {
	for _, userId := range userIds {
		err := domain.ClearUserRoleList(orgId, userId, projectId)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

// 去除了 rolesvc 后，导致引入的东西不存在，方法先注释
func GetDepartmentRoleList(orgId, departmentId int64, projectId int64) {
	/*
	key, err5 := util.ParseCacheKey(sconsts.CacheDepartmentRoleListHash, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyDepartmentIdConstName: departmentId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	roleDepartmentListJson, err := cache.HGet(key, strconv.FormatInt(projectId, 10))
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if roleDepartmentListJson != "" {
		roleDepartmentList := &[]bo.RoleDepartmentBo{}
		err := json.FromJson(roleDepartmentListJson, roleDepartmentList)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return roleDepartmentList, nil
	} else {
		roleDepartmentListPo := &[]po.PpmRolRoleDepartment{}
		roleDepartmentListBo := &[]bo.RoleDepartmentBo{}
		err := mysql.SelectAllByCond(consts.TableRoleDepartment, db.Cond{
			consts.TcOrgId:     orgId,
			consts.TcDepartmentId:    departmentId,
			consts.TcProjectId: projectId,
			consts.TcIsDelete:  consts.AppIsNoDelete,
		}, roleDepartmentListPo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(roleDepartmentListPo, roleDepartmentListBo)
		roleDepartmentListJson, err = json.ToJson(roleDepartmentListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.HSet(key, strconv.FormatInt(projectId, 10), roleDepartmentListJson)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return roleDepartmentListBo, nil
	}
	*/
}
