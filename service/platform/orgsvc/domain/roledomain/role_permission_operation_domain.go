package domain

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	po2 "github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func UpdateRolePermissionOperation(orgId int64, userId int64, roleId int64, permissionOperation map[int64][]string, roleLangCode string, projectId int64) errs.SystemErrorInfo {
	insertPo := []interface{}{}

	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableRolePermissionOperation, len(permissionOperation))
	if err != nil {
		log.Error(err)
		return err
	}
	k := 0
	var oldPermission []int64
	for i, i2 := range permissionOperation {
		oldPermission = append(oldPermission, i)
		//获取相关的权限详情
		info, err := GetPermissionById(i)
		if err != nil {
			return err
		}
		length := len(i2)
		var operationCodes string
		if length == 0 {
			continue
		} else if length == 1 {
			operationCodes = i2[0]
		} else {
			for _, s := range i2 {
				operationCodes += fmt.Sprintf("(%s)|", s)
			}
			operationCodes = operationCodes[0 : len(operationCodes)-1]
		}
		path := strings.ReplaceAll(info.Path, "{org_id}", strconv.FormatInt(orgId, 10))
		path = strings.ReplaceAll(path, "{pro_id}", strconv.FormatInt(0, 10))
		temp := po2.PpmRolRolePermissionOperation{
			Id:             ids.Ids[k].Id,
			OrgId:          orgId,
			RoleId:         roleId,
			PermissionId:   i,
			PermissionPath: path,
			OperationCodes: operationCodes,
			Creator:        userId,
		}
		if roleLangCode == consts.RoleGroupProMember {
			temp.ProjectId = projectId
		}
		insertPo = append(insertPo, temp)
		k++
	}

	//上锁
	lockKey := fmt.Sprintf("%s%d:%d", consts.ModifyRolePermissionLock, orgId, roleId)
	lockUuid := uuid.NewUuid()

	suc, lockErr := cache.TryGetDistributedLock(lockKey, lockUuid)
	if lockErr != nil {
		log.Error(lockErr)
		return errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, lockUuid); err != nil {
				log.Error(err)
			}
		}()
	}
	//删除旧权限,增加新权限
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(oldPermission) > 0 {
			cond := db.Cond{
				consts.TcIsDelete:     consts.AppIsNoDelete,
				consts.TcRoleId:       roleId,
				consts.TcPermissionId: db.In(oldPermission),
				consts.TcOrgId:        orgId,
			}
			if roleLangCode == consts.RoleGroupProMember {
				cond[consts.TcProjectId] = projectId
			}
			_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableRolePermissionOperation, cond, mysql.Upd{
				consts.TcIsDelete: consts.AppIsDeleted,
				consts.TcUpdator:  userId,
			})
			if err != nil {
				log.Error(err)
				return err
			}
		}

		err := mysql.TransBatchInsert(tx, &po2.PpmRolRolePermissionOperation{}, insertPo)
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})
	if transErr != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}

	return nil
}

func CopyPermission(permissionList *[]bo.RolePermissionOperationBo, roleId int64, orgId int64) errs.SystemErrorInfo {
	if permissionList == nil || len(*permissionList) == 0 {
		return nil
	}

	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableRolePermissionOperation, len(*permissionList))
	if err != nil {
		log.Error(err)
		return err
	}
	insertPo := []interface{}{}
	for i, operationBo := range *permissionList {
		operationPo := &po2.PpmRolRolePermissionOperation{}
		copyErr := copyer.Copy(operationBo, operationPo)
		if copyErr != nil {
			log.Error(copyErr)
			return errs.ObjectCopyError
		}
		operationPo.Id = ids.Ids[i].Id
		operationPo.RoleId = roleId
		operationPo.OrgId = orgId
		insertPo = append(insertPo, *operationPo)
	}

	insertErr := mysql.BatchInsert(&po2.PpmRolRolePermissionOperation{}, insertPo)
	if insertErr != nil {
		log.Error(insertErr)
		return errs.MysqlOperateError
	}

	return nil
}
