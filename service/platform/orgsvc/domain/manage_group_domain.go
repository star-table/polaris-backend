package domain

import (
	"fmt"
	"strconv"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	int64Slice "github.com/star-table/polaris-backend/common/core/util/slice/int64"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// GetOrgSysAdmin 获取组织超管，返回结果 orgId -> adminIds
func GetOrgSysAdmin(orgIds []int64) (map[int64][]int64, errs.SystemErrorInfo) {
	groups := make([]po.LcPerManageGroup, 0)
	err := mysql.SelectAllByCond(consts.TableLcPerManageGroup, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    db.In(orgIds),
		consts.TcLangCode: consts.ManageGroupLangCodeSysAdmin,
	}, &groups)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	results := map[int64][]int64{}
	for _, group := range groups {
		if group.UserIds == nil || *group.UserIds == "[]" {
			continue
		}
		tmpUserIds := make([]int64, 0)
		_ = json.FromJson(*group.UserIds, &tmpUserIds)
		results[group.OrgId] = tmpUserIds
	}
	return results, nil
}

func GetManageGroupList(orgId int64, groupIds []int64) ([]*po.LcPerManageGroup, errs.SystemErrorInfo) {
	groups := make([]*po.LcPerManageGroup, 0)
	err := mysql.SelectAllByCond(consts.TableLcPerManageGroup, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(groupIds),
	}, &groups)
	if err != nil {
		log.Errorf("[GetManageGroupList] err:%v, orgId:%v, groupIds:%v", err, orgId, groupIds)
		return nil, errs.MysqlOperateError
	}
	return groups, nil
}

// 将用户绑定到角色管理组
func BindUserRoleGroups(orgId, operator int64, groupIds []int64, userIds []int64) errs.SystemErrorInfo {
	userIdsStr := json.ToJsonIgnoreError(userIds)
	upd := mysql.Upd{
		consts.TcUserIds: userIdsStr,
		consts.TcUpdator: operator,
		consts.TcVersion: db.Raw("version + 1"),
	}

	_, err := mysql.UpdateSmartWithCond(consts.TableLcPerManageGroup, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       db.In(groupIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)

	if err != nil {
		log.Errorf("[BindUserRoleGroups] err:%v, orgId:%v, groupIds:%v", err, orgId, groupIds)
		return errs.MysqlOperateError
	}
	return nil
}

func GetManageGroupListByUsers(orgId int64, userIds []int64) ([]*po.LcPerManageGroup, errs.SystemErrorInfo) {
	conn, dbErr := mysql.GetConnect()
	if dbErr != nil {
		log.Errorf("[GetManageGroupListByUsers] err:%v", dbErr)
		return nil, errs.MysqlOperateError
	}
	//拼装条件
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}
	userIdsStr := make([]string, len(userIds))
	for i, v := range userIds {
		userIdsStr[i] = strconv.FormatInt(v, 10)
	}
	var groups []*po.LcPerManageGroup
	dbErr = conn.Select(db.Raw("*")).
		From(consts.TableLcPerManageGroup).
		Where(cond).
		And(
			db.Or(
				db.Raw(fmt.Sprintf("JSON_OVERLAPS(`user_ids` -> '$', CAST('%s' AS JSON))", json.ToJsonIgnoreError(userIds))),
				db.Raw(fmt.Sprintf("JSON_OVERLAPS(`user_ids` -> '$', CAST('%s' AS JSON))", json.ToJsonIgnoreError(userIdsStr))),
			),
		).
		All(&groups)
	if dbErr != nil {
		log.Errorf("[GetManageGroupListByUsers] err:%v, orgId:%v, userIds:%v", dbErr, orgId, userIds)
		return nil, errs.MysqlOperateError
	}
	return groups, nil
}

func GetManageGroupListByUser(orgId int64, userId int64) (*po.LcPerManageGroup, errs.SystemErrorInfo) {
	list, err := GetManageGroupListByUsers(orgId, []int64{userId})
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errs.MysqlOperateError
	}
	return list[0], nil
}

// 切换用户管理组
func ChangeUserManageGroup(orgId, operatorId int64, req orgvo.ChangeUserManageGroupReq) errs.SystemErrorInfo {
	manageGroupList, errSys := GetManageGroupList(orgId, req.ManageGroupIds)
	if errSys != nil {
		log.Errorf("[ChangeUserManageGroup] GetManageGroupList err:%v, orgId:%v", errSys, orgId)
		return errSys
	}
	if len(manageGroupList) != len(req.ManageGroupIds) {
		return errs.RoleNotExist
	}
	oldAdminGroups, errSys := GetManageGroupListByUsers(orgId, []int64{req.UserId})
	if errSys != nil {
		log.Errorf("[ChangeUserManageGroup] GetManageGroupListByUsers err:%v, orgId:%v", errSys, orgId)
		return errSys
	}
	oldGroupIds := []int64{}
	for _, group := range oldAdminGroups {
		oldGroupIds = append(oldGroupIds, group.Id)
	}
	_, addIds, delIds := int64Slice.CompareSliceAddDelInt64(req.ManageGroupIds, oldGroupIds)
	if len(addIds) == 0 && len(delIds) == 0 {
		// 说明没有变动
		return nil
	}

	sysGroup, err := GetSysManageGroup(orgId)
	if err != nil {
		log.Errorf("[ChangeUserManageGroup] GetSysManageGroup err:%v, orgId:%v", err, orgId)
		return errs.MysqlOperateError
	}
	if ok, _ := slice.Contain(addIds, sysGroup.Id); ok {
		return errs.DenyChangeSysAdminGroupOfUser
	}
	if ok, _ := slice.Contain(delIds, sysGroup.Id); ok {
		return errs.DenyChangeSysAdminGroupOfUser
	}

	err = mysql.TransX(func(tx sqlbuilder.Tx) error {
		for _, id := range addIds {
			_, err := AppendUserIntoAdminGroup(id, operatorId, req.UserId, tx)
			if err != nil {
				log.Errorf("[ChangeUserManageGroup] AppendUserIntoAdminGroup err:%v, orgId:%v", err, orgId)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
		}

		for _, oneGroup := range oldAdminGroups {
			if ok, _ := slice.Contain(delIds, oneGroup.Id); !ok {
				continue
			}
			tmpUserIds := make([]int64, 0)
			err = json.FromJson(*oneGroup.UserIds, &tmpUserIds)
			if err != nil {
				log.Error(err)
				return errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
			}
			index := int64Slice.GetSearchedIndexArr(tmpUserIds, req.UserId)
			if index != -1 {
				_, err = RemoveUserFromAdminGroup(oneGroup.Id, operatorId, index, tx)
				if err != nil {
					log.Error(err)
					return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Errorf("[ChangeUserManageGroup] update err:%v, orgId:%v", err, orgId)
		return errs.MysqlOperateError
	}

	return nil
}

func GetSysManageGroup(orgId int64) (*po.LcPerManageGroup, error) {
	var group po.LcPerManageGroup
	dbErr := mysql.SelectOneByCond(consts.TableLcPerManageGroup, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcLangCode: consts.ManageGroupSys,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &group)

	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return nil, errs.ManageGroupNotExist
		}
		return nil, dbErr
	}
	return &group, nil
}

// AppendUserIntoAdminGroup 向管理组中增加一个人
func AppendUserIntoAdminGroup(id int64, operateUid int64, userId int64, tx sqlbuilder.Tx) (int64, error) {
	var (
		count int64
		dbErr error
	)
	cond1 := db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	upd1 := mysql.Upd{
		consts.TcUserIds: db.Raw(fmt.Sprintf("JSON_ARRAY_APPEND(`%s`, '$', '%d')", consts.TcUserIds, userId)),
		consts.TcUpdator: operateUid,
		consts.TcVersion: db.Raw("version + 1"),
	}
	if tx != nil {
		count, dbErr = mysql.TransUpdateSmartWithCond(
			tx,
			consts.TableLcPerManageGroup,
			cond1,
			upd1,
		)
	} else {
		count, dbErr = mysql.UpdateSmartWithCond(
			consts.TableLcPerManageGroup,
			cond1,
			upd1,
		)
	}

	if dbErr != nil {
		log.Error(dbErr)
		return 0, dbErr
	}
	return count, nil
}

// RemoveUserFromAdminGroup 删除管理组中的某个人
func RemoveUserFromAdminGroup(id int64, operateUid int64, index int, tx sqlbuilder.Tx) (int64, error) {
	var (
		count int64
		dbErr error
	)
	cond1 := db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	upd1 := mysql.Upd{
		consts.TcUserIds: db.Raw(fmt.Sprintf("json_remove(`%s`, '$[%d]')", consts.TcUserIds, index)),
		consts.TcUpdator: operateUid,
		consts.TcVersion: db.Raw("version + 1"),
	}
	if tx != nil {
		count, dbErr = mysql.TransUpdateSmartWithCond(
			tx,
			consts.TableLcPerManageGroup,
			cond1,
			upd1,
		)
	} else {
		count, dbErr = mysql.UpdateSmartWithCond(
			consts.TableLcPerManageGroup,
			cond1,
			upd1,
		)
	}
	if dbErr != nil {
		log.Error(dbErr)
		return 0, dbErr
	}
	return count, nil
}
