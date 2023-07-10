package domain

import (
	"fmt"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//获取组织用户特殊角色
func GetOrgRoleUser(orgId int64, projectId int64, langCodes []string) (*[]bo.RoleUser, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	rolePo := &[]po.RoleUser{}
	cond := db.Cond{
		"r." + consts.TcOrgId:    orgId,
		"g." + consts.TcOrgId:    orgId,
		"u." + consts.TcOrgId:    orgId,
		"r." + consts.TcIsDelete: consts.AppIsNoDelete,
		"g." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"r." + consts.TcStatus:   consts.AppStatusEnable,
		"r.role_group_id":        db.Raw("g.id"),
		"u.role_id":              db.Raw("r.id"),
	}
	if projectId != 0 {
		cond["g."+consts.TcLangCode] = consts.RoleGroupPro
		cond["r."+consts.TcProjectId] = projectId
	} else {
		cond["g."+consts.TcLangCode] = consts.RoleGroupOrg
	}

	if langCodes != nil && len(langCodes) > 0 {
		cond["r."+consts.TcLangCode] = db.In(langCodes)
	}

	err = conn.Select(db.Raw("u.user_id, r.id as role_id, r.name as role_name, r.lang_code")).From("ppm_rol_role as r", "ppm_rol_role_group as g", "ppm_rol_role_user as u").Where(cond).OrderBy("r.id asc").All(rolePo)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	roleBo := &[]bo.RoleUser{}
	copyErr := copyer.Copy(rolePo, roleBo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return roleBo, nil
}

func GetOrgRoleUserBatch(orgIds []int64, projectId int64, langCodes []string) (*[]bo.RoleUserWithOrgId, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	rolePo := &[]po.RoleUserWithOrgId{}
	cond := db.Cond{
		"r." + consts.TcOrgId:    db.In(orgIds),
		"g." + consts.TcOrgId:    db.In(orgIds),
		"u." + consts.TcOrgId:    db.In(orgIds),
		"r." + consts.TcIsDelete: consts.AppIsNoDelete,
		"g." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"r." + consts.TcStatus:   consts.AppStatusEnable,
		"r.role_group_id":        db.Raw("g.id"),
		"u.role_id":              db.Raw("r.id"),
	}
	if projectId != 0 {
		cond["g."+consts.TcLangCode] = consts.RoleGroupPro
		cond["r."+consts.TcProjectId] = projectId
	} else {
		cond["g."+consts.TcLangCode] = consts.RoleGroupOrg
	}

	if langCodes != nil && len(langCodes) > 0 {
		cond["r."+consts.TcLangCode] = db.In(langCodes)
	}

	err = conn.Select(db.Raw("u.org_id, u.user_id, r.id as role_id, r.name as role_name, r.lang_code")).From("ppm_rol_role as r", "ppm_rol_role_group as g", "ppm_rol_role_user as u").Where(cond).OrderBy("r.id asc").All(rolePo)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	roleBos := &[]bo.RoleUserWithOrgId{}
	copyErr := copyer.Copy(rolePo, roleBos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return roleBos, nil
}

func UpdateUserOrgRole(role *bo.RoleBo, orgId int64, currentUserId int64, userIds []int64, projectId *int64) (int64, errs.SystemErrorInfo) {
	lockKey := consts.UpdateUserOrgRoleLock + fmt.Sprintf("%d:%d:%d", role.OrgId, role.Id, role.ProjectId)

	uuid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uuid)
	if lockErr != nil {
		log.Error("获取锁异常")
		return 0, errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}

	if suc {
		//释放锁
		defer func() {
			if _, e := cache.ReleaseDistributedLock(lockKey, uuid); e != nil {
				log.Error(e)
			}
		}()

	} else {
		return 0, errs.BuildSystemErrorInfo(errs.GetDistributedLockError)
	}

	mysqlErr := mysql.TransX(func(tx sqlbuilder.Tx) error {

		//删除旧有组织角色 如果有projectId删除对应project下的所有角色关系 如果没有 删除组织下的所有角色关系
		cond := db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcUserId:   db.In(userIds),
			consts.TcOrgId:    orgId,
		}
		if projectId != nil {
			cond[consts.TcProjectId] = projectId
		}
		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableRoleUser, cond, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  currentUserId,
		})
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		//如果不是通用角色并且不是项目成员(归入特殊角色)，就执行插入操作
		if role.OrgId != 0 && role.LangCode != consts.RoleGroupProMember {
			//增加新角色
			idResp, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableRoleUser, len(userIds))
			if idErr != nil {
				log.Error(idErr)
				return idErr
			}
			pos := []interface{}{}
			for i, id := range userIds {
				pos = append(pos, po.PpmRolRoleUser{
					Id:        idResp.Ids[i].Id,
					OrgId:     role.OrgId,
					ProjectId: role.ProjectId,
					RoleId:    role.Id,
					UserId:    id,
					Creator:   currentUserId,
					Updator:   currentUserId,
				})
			}
			insertErr := mysql.TransBatchInsert(tx, &po.PpmRolRoleUser{}, pos)
			if insertErr != nil {
				log.Error(insertErr)
				log.Error(insertErr)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertErr)
			}
		}
		return err
	})

	if mysqlErr != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	//删除缓存
	trueProjectId := int64(0)
	if projectId != nil {
		trueProjectId = *projectId
	}
	for _, id := range userIds {
		clearErr := ClearUserRoleList(orgId, id, trueProjectId)
		if clearErr != nil {
			log.Error(clearErr)
			return 0, clearErr
		}
	}

	return int64(len(userIds)), nil
}

//添加角色和用户关联
func AddRoleUserRelation(orgId, userId, roleId int64) errs.SystemErrorInfo {
	roleUser := &po.PpmRolRoleUser{}
	err := mysql.SelectOneByCond(consts.TableRoleUser, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   userId,
		consts.TcRoleId:   roleId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, roleUser)
	if err != nil && err == db.ErrNoMoreRows {
		id, err1 := idfacade.ApplyPrimaryIdRelaxed(roleUser.TableName())
		if err1 != nil {
			log.Error(err1)
			return errs.BuildSystemErrorInfo(errs.ApplyIdError, err1)
		}
		roleUser.Id = id
		roleUser.OrgId = orgId
		roleUser.UserId = userId
		roleUser.RoleId = roleId
		roleUser.IsDelete = consts.AppIsNoDelete
		err2 := mysql.Insert(roleUser)
		if err2 != nil {
			log.Error(err2)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
		}
	}
	return nil
}

//移除角色和用户关联
func RemoveRoleUserRelation(orgId int64, userIds []int64, operatorId int64) errs.SystemErrorInfo {
	_, err := mysql.UpdateSmartWithCond(consts.TableRoleUser, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(userIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operatorId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

//移除角色和部门关联
func RemoveRoleDepartmentRelation(orgId int64, deptIds []int64, operatorId int64) errs.SystemErrorInfo {
	_, err := mysql.UpdateSmartWithCond(consts.TableRoleDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcDepartmentId:   db.In(deptIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operatorId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}
