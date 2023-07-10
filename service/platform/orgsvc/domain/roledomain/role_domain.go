package domain

import (
	"fmt"
	"strings"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	po2 "github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//特殊角色名称
var DefaultRoleName = []string{
	"组织超级管理员",
	"组织管理员",
	"组织成员",
	"负责人",
	"项目成员",
}

func CreateRole(orgId, userId int64, input vo.CreateRoleReq, groupId int64) (int64, errs.SystemErrorInfo) {
	//判断是否重名
	newUUID := uuid.NewUuid()
	var projectId int64
	if input.ProjectID != nil {
		projectId = *input.ProjectID
	}
	lockKey := fmt.Sprintf("%s%d:%d", consts.ModifyRoleLock, orgId, projectId)
	suc, err := cache.TryGetDistributedLock(lockKey, newUUID)
	if err != nil {
		log.Error(err)
		return 0, errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, newUUID); err != nil {
				log.Error(err)
			}
		}()
	} else {
		return 0, errs.SystemError
	}

	if ok, _ := slice.Contain(DefaultRoleName, input.Name); ok {
		return 0, errs.DefaultRoleNameErr
	}
	//成功获取锁，先判断是否重名
	_, roleErr := GetRoleByName(orgId, projectId, input.Name)
	if roleErr != nil {
		if roleErr.Code() != errs.RoleNotExist.Code() {
			log.Error(roleErr)
			return 0, roleErr
		}
	} else {
		return 0, errs.RoleNameRepeatErr
	}
	po := &po2.PpmRolRole{}
	copyErr := copyer.Copy(input, po)
	if copyErr != nil {
		log.Error(copyErr)
		return 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	id, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableRole)
	if idErr != nil {
		log.Error(idErr)
		return 0, idErr
	}
	po.Id = id
	po.RoleGroupId = groupId
	po.Creator = userId
	po.OrgId = orgId
	po.Updator = userId
	if input.ProjectID != nil {
		po.ProjectId = *input.ProjectID
	}
	insertErr := mysql.Insert(po)
	if insertErr != nil {
		log.Error(insertErr)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertErr)
	}

	return id, nil
}

func JudgeRoleName(name string) errs.SystemErrorInfo {
	name = strings.TrimSpace(name)
	if len(name) == 0 || strs.Len(name) > 10 {
		log.Error("角色名字不符合规范")
		return errs.RoleNameLenErr
	}

	return nil
}

//没有projectId就传0
func UpdateRole(orgId int64, projectId int64, roleId int64, updateBo bo.UpdateRoleBo, operatorId int64) errs.SystemErrorInfo {
	upd := mysql.Upd{}
	if updateBo.Name != nil {
		name := *updateBo.Name
		//nameErr := JudgeRoleName(name)
		//if nameErr != nil {
		//	return nameErr
		//}
		isNameRight := format.VerifyRoleNameFormat(name)
		if !isNameRight {
			return errs.RoleNameLenErr
		}

		//判断是否重名
		newUUID := uuid.NewUuid()
		lockKey := fmt.Sprintf("%s%d:%d", consts.ModifyRoleLock, orgId, projectId)
		suc, err := cache.TryGetDistributedLock(lockKey, newUUID)
		if err != nil {
			log.Error(err)
			return errs.TryDistributedLockError
		}
		if suc {
			defer func() {
				if _, err := cache.ReleaseDistributedLock(lockKey, newUUID); err != nil {
					log.Error(err)
				}
			}()
		} else {
			return errs.RoleModifyBusy
		}

		if ok, _ := slice.Contain(DefaultRoleName, name); ok {
			return errs.DefaultRoleNameErr
		}
		//成功获取锁，先判断是否重名
		roleBo, roleErr := GetRoleByName(orgId, projectId, name)
		if roleErr != nil {
			if roleErr.Code() == errs.RoleNotExist.Code() {
				upd[consts.TcName] = name
			} else {
				log.Error(roleErr)
			}
		} else if roleBo.Id != roleId { //如果不是同一角色
			return errs.RoleNameRepeatErr
		}
	}
	if len(upd) == 0 {
		return nil
	}

	//加入更新人
	upd[consts.TcUpdator] = operatorId

	_, updateErr := mysql.UpdateSmartWithCond(consts.TableRole, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       roleId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)
	if updateErr != nil {
		log.Error(updateErr)
		return errs.MysqlOperateError
	}
	return nil
}

func DeleteRoles(orgId int64, roleIds []int64, operatorId int64) errs.SystemErrorInfo {
	if len(roleIds) == 0 {
		return nil
	}
	//查询用户角色绑定关系
	roleUserPos := &[]po2.PpmRolRoleUser{}
	selectErr := mysql.SelectAllByCond(consts.TableRoleUser, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcRoleId:   db.In(roleIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, roleUserPos)
	if selectErr != nil {
		log.Error(selectErr)
		return errs.MysqlOperateError
	}
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除角色
		_, updateErr := mysql.TransUpdateSmartWithCond(tx, consts.TableRole, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcId:       db.In(roleIds),
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  operatorId,
		})
		if updateErr != nil {
			log.Error(updateErr)
			return errs.MysqlOperateError
		}

		//删除用户角色绑定关系
		_, updateErr1 := mysql.TransUpdateSmartWithCond(tx, consts.TableRoleUser, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcRoleId:   db.In(roleIds),
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  operatorId,
		})
		if updateErr1 != nil {
			log.Error(updateErr1)
			return errs.MysqlOperateError
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}

	//删除用户角色缓存
	for _, user := range *roleUserPos {
		clearErr := ClearUserRoleList(orgId, user.UserId, user.ProjectId)
		if clearErr != nil {
			log.Error(clearErr)
			return clearErr
		}
	}

	return nil
}

func GetRole(orgId, projectId, roleId int64) (*bo.RoleBo, errs.SystemErrorInfo) {
	roles, err := GetRoleList(orgId, projectId, []int64{roleId}, true)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if len(roles) == 0 {
		log.Error("角色不存在")
		return nil, errs.RoleNotExist
	}
	return &roles[0], nil
}

func GetRoleByName(orgId, projectId int64, roleName string) (*bo.RoleBo, errs.SystemErrorInfo) {
	role := &po2.PpmRolRole{}
	//只关心自定义的角色名称
	err := mysql.SelectOneByCond(consts.TableRole, db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
		consts.TcName:      roleName,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcLangCode:  consts.BlankString,
	}, role)
	if err == db.ErrNoMoreRows {
		log.Error("记录不存在")
		return nil, errs.RoleNotExist
	}
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	roleBo := &bo.RoleBo{}
	_ = copyer.Copy(role, roleBo)
	return roleBo, nil
}

func GetSysDefaultRoles(orgId int64) ([]bo.RoleBo, errs.SystemErrorInfo) {
	rolePos := &[]po2.PpmRolRole{}
	dbErr := mysql.SelectAllByCond(consts.TableRole, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcLangCode: db.NotEq(consts.BlankString),
	}, rolePos)
	if dbErr != nil {
		log.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	result := &[]bo.RoleBo{}
	_ = copyer.Copy(rolePos, result)
	return *result, nil
}

//projectId：传0表示忽略此条件
//org_id = 0:传空忽略此条件
func GetRoleList(orgId int64, projectId int64, roleIds []int64, includeAll bool) ([]bo.RoleBo, errs.SystemErrorInfo) {
	result := &[]bo.RoleBo{}
	//if len(roleIds) == 0 {
	//	return *result, nil
	//}

	//去重
	roleIds = slice.SliceUniqueInt64(roleIds)

	//拼装条件
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if len(roleIds) > 0 {
		cond[consts.TcId] = db.In(roleIds)
	}
	if projectId > 0 {
		projectIds := []int64{projectId}
		//是否查询全部（包括组织下特殊角色）
		if includeAll {
			projectIds = append(projectIds, 0)
		}
		cond[consts.TcProjectId] = db.In(projectIds)
	}
	if orgId > 0 {
		cond[consts.TcOrgId] = db.In([]int64{0, orgId})
	}
	rolePos := &[]po2.PpmRolRole{}
	dbErr := mysql.SelectAllByCond(consts.TableRole, cond, rolePos)
	if dbErr != nil {
		log.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	_ = copyer.Copy(rolePos, result)
	return *result, nil
}

////实时获取角色列表
//func GetRoleListActually(orgId int64) ([]bo.RoleBo, errs.SystemErrorInfo) {
//	roleListPo := &[]po2.PpmRolRole{}
//	err := mysql.SelectAllByCond(consts.TableRole, db.Cond{
//		consts.TcOrgId:    orgId,
//		consts.TcStatus:   consts.AppStatusEnable,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}, roleListPo)
//	if err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//
//	roleListBo := &[]bo.RoleBo{}
//	copyErr := copyer.Copy(roleListPo, roleListBo)
//	if copyErr != nil {
//		log.Error(copyErr)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
//	}
//
//	return *roleListBo, nil
//}
