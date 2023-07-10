package domain

import (
	"strconv"

	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

const RoleInitSql = consts.TemplateDirPrefix + "role_init.template"

func RoleInit(orgId int64, tx sqlbuilder.Tx) (*bo.RoleInitResp, errs.SystemErrorInfo) {
	maps := map[string]interface{}{}
	maps["OrgId"] = orgId

	//ppm_rol_role_group主键封装
	roleGroup := &po.PpmRolRoleGroup{}
	roleGroupIds, err := idfacade.ApplyMultipleIdRelaxed(0, roleGroup.TableName(), "", 3)
	if err != nil {
		return nil, err
	}
	for k, v := range roleGroupIds.Ids {
		maps["RoleGroupId"+strconv.Itoa(k+1)] = v.Id
	}
	log.Infof("ppm_rol_role_group主键分配完成")

	//ppm_rol_role主键封装
	role := &po.PpmRolRole{}
	roleIds, err := idfacade.ApplyMultipleIdRelaxed(0, role.TableName(), "", 14)
	if err != nil {
		return nil, err
	}
	for k, v := range roleIds.Ids {
		maps["RoleId"+strconv.Itoa(k+1)] = v.Id
	}
	//获取组织超级管理员角色id和普通管理员角色id
	orgSuperAdminRoleId := maps["RoleId7"].(int64)
	orgNormalAdminRoleId := maps["RoleId8"].(int64)
	log.Infof("ppm_rol_role主键分配完成")

	//ppm_rol_role_permission_operation主键封装
	rolePermissionOperation := &po.PpmRolRolePermissionOperation{}
	rolePermissionOperationIds, err := idfacade.ApplyMultipleIdRelaxed(0, rolePermissionOperation.TableName(), "", 164)
	if err != nil {
		return nil, err
	}
	for k, v := range rolePermissionOperationIds.Ids {
		maps["PermissionOperation"+strconv.Itoa(k+1)] = v.Id
	}
	log.Infof("ppm_rol_role_permission_operation主键分配完成")

	// TODO 切换lesscode-permission
	//ppm_rol_permission主键封装
	permission := &po.PpmRolPermission{}
	permissionIds, err := idfacade.ApplyMultipleIdRelaxed(0, permission.TableName(), "", 29)
	if err != nil {
		return nil, err
	}
	for k, v := range permissionIds.Ids {
		maps["PermissionId"+strconv.Itoa(k+1)] = v.Id
	}
	log.Infof("ppm_rol_permission主键分配完成")

	// TODO 切换lesscode-permission
	//ppm_rol_permission_operation主键封装
	rolPermissionOperation := &po.PpmRolPermission{}
	rolPermissionOperationIds, err := idfacade.ApplyMultipleIdRelaxed(0, rolPermissionOperation.TableName(), "", 119)
	if err != nil {
		return nil, err
	}
	for k, v := range rolPermissionOperationIds.Ids {
		maps["RolPermissionOperationId"+strconv.Itoa(k+1)] = v.Id
	}
	log.Infof("ppm_rol_permission_operation主键分配完成")

	err = util.ReadAndWrite(RoleInitSql, maps, tx)
	if err != nil {
		return nil, err
	}

	return &bo.RoleInitResp{
		OrgSuperAdminRoleId:  orgSuperAdminRoleId,
		OrgNormalAdminRoleId: orgNormalAdminRoleId,
	}, nil
}

func ChangeDefaultRole() errs.SystemErrorInfo {
	//给组织管理员添加权限（查看用户列表，变更用户部门，部门新增/删除/修改）
	pos := &[]po.PpmRolRole{}
	err := mysql.SelectAllByCond(consts.TableRole, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcLangCode: consts.RoleGroupOrgManager,
	}, pos)
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}

	roleIds := []int64{}
	roleMap := map[int64]int64{}
	for _, role := range *pos {
		roleIds = append(roleIds, role.Id)
		roleMap[role.Id] = role.OrgId
	}

	allCount := len(roleIds)
	if allCount == 0 {
		return nil
	}

	page := 0
	size := 100
	next := true
	for {
		if !next {
			break
		}
		start := page * size
		end := (page + 1) * size
		if end >= allCount {
			end = allCount
			next = false
		}
		curr := roleIds[start:end]

		_, updErr := mysql.UpdateSmartWithCond(consts.TableRolePermissionOperation, db.Cond{
			consts.TcRoleId:       db.In(curr),
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcPermissionId: 12,
		}, mysql.Upd{
			consts.TcOperationCodes: "(View)|(Create)|(Attention)|(UnAttention)|(ModifyField)",
		})
		if updErr != nil {
			log.Error(updErr)
			return errs.MysqlOperateError
		}

		page++
	}

	return nil
}
