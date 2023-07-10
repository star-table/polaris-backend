package domain

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/star-table/common/core/util/id/snowflake"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cast"
	"upper.io/db.v3/lib/sqlbuilder"
)

type _makeInfo struct {
	originOrgId   int64
	newOrgId      int64
	originUserIds []int64
	newUserIds    []int64
}

func TestTransferOrgToOtherPlatform(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		makeInfo := &_makeInfo{
			originOrgId: 10000001,
			newOrgId:    10000002,
		}
		err := _makeFakeInfo(makeInfo)
		if err != nil {
			log.Error(err)
			return
		}

		//fmt.Println(TransferOrgToOtherPlatform(makeInfo.originOrgId, makeInfo.newOrgId))
	}))
}

// 制造假数据
func _makeFakeInfo(makeInfo *_makeInfo) error {
	dbErr := _delFakeInfo(makeInfo)
	if dbErr != nil {
		return dbErr
	}

	dbErr = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := _makeUserInfo(makeInfo, tx)
		if err != nil {
			return err
		}

		err = _makeOrgInfo(makeInfo, tx)
		if err != nil {
			return err
		}

		return _makeDepartmentInfo(makeInfo, tx)
	})

	return dbErr
}

func _makeOrgInfo(makeInfo *_makeInfo, tx sqlbuilder.Tx) error {
	err := _createOrgInfo(makeInfo.originOrgId, makeInfo.originUserIds, "fs", tx)
	if err != nil {
		return err
	}

	return _createOrgInfo(makeInfo.newOrgId, makeInfo.newUserIds, "ding", tx)
}

func _createOrgInfo(orgId int64, userIds []int64, sourceChannel string, tx sqlbuilder.Tx) error {
	org := &po.PpmOrgOrganization{}
	org.Id = orgId
	org.Status = consts.AppStatusEnable
	org.IsDelete = consts.AppIsNoDelete
	org.Creator = 1
	org.Owner = 1
	org.Updator = 1
	org.Name = cast.ToString(orgId)
	org.SourceChannel = sourceChannel
	org.SourcePlatform = sourceChannel

	orgOutInfo := &po.PpmOrgOrganizationOutInfo{
		Id:             snowflake.Id(),
		OrgId:          orgId,
		OutOrgId:       cast.ToString(orgId),
		SourceChannel:  sourceChannel,
		SourcePlatform: sourceChannel,
		TenantCode:     cast.ToString(orgId),
		Name:           cast.ToString(orgId),
		Status:         consts.AppStatusEnable,
		IsDelete:       consts.AppIsNoDelete,
	}
	for _, userId := range userIds {
		userOrgInfo := &po.PpmOrgUserOrganization{
			Id:          snowflake.Id(),
			OrgId:       orgId,
			UserId:      userId,
			CheckStatus: consts.AppCheckStatusSuccess,
			UseStatus:   consts.AppStatusDisabled,
			Status:      consts.AppStatusEnable,
		}
		err := mysql.TransInsert(tx, userOrgInfo)
		if err != nil {
			return err
		}
	}

	err := mysql.TransInsert(tx, org)
	if err != nil {
		return err
	}

	return mysql.TransInsert(tx, orgOutInfo)
}

func _makeUserInfo(makeInfo *_makeInfo, tx sqlbuilder.Tx) error {
	for i := 0; i < 5; i++ {
		userId, err := _createUser(makeInfo.originOrgId, "fs", tx)
		if err != nil {
			return err
		}
		makeInfo.originUserIds = append(makeInfo.originUserIds, userId)

		newUserId, err := _createUser(makeInfo.newOrgId, "ding", tx)
		if err != nil {
			return err
		}
		makeInfo.newUserIds = append(makeInfo.newUserIds, newUserId)

		if i < 3 {
			err = _createGlobalUser(userId, newUserId, tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func _createUser(orgId int64, sourceChannel string, tx sqlbuilder.Tx) (int64, error) {
	user := &po.PpmOrgUser{
		Id:             snowflake.Id(),
		OrgId:          orgId,
		Name:           "1111",
		NamePinyin:     "111",
		LoginName:      "1111",
		SourceChannel:  sourceChannel,
		SourcePlatform: sourceChannel,
	}
	userOut := &po.PpmOrgUserOutInfo{
		Id:             snowflake.Id(),
		OrgId:          orgId,
		UserId:         user.Id,
		SourcePlatform: sourceChannel,
		SourceChannel:  sourceChannel,
		OutOrgUserId:   cast.ToString(orgId),
		OutUserId:      cast.ToString(orgId),
	}

	err := mysql.TransInsert(tx, user)
	if err != nil {
		return user.Id, err
	}

	return user.Id, mysql.TransInsert(tx, userOut)
}

func _createGlobalUser(userId1, userId2 int64, tx sqlbuilder.Tx) error {
	globalId := snowflake.Id()
	err := dao.GetGlobalUser().Create(&po.PpmOrgGlobalUser{
		Id:              globalId,
		Mobile:          cast.ToString(globalId)[3:],
		LastLoginUserId: userId1,
		LastLoginOrgId:  userId1,
	}, tx)
	if err != nil {
		return err
	}

	return dao.GetGlobalUserRelation().CreateRelations([]*po.PpmOrgGlobalUserRelation{
		{
			Id:           snowflake.Id(),
			GlobalUserId: globalId,
			UserId:       userId1,
		},
		{
			Id:           snowflake.Id(),
			GlobalUserId: globalId,
			UserId:       userId2,
		},
	}, tx)
}

func _makeDepartmentInfo(makeInfo *_makeInfo, tx sqlbuilder.Tx) error {
	for i := 0; i < 5; i++ {
		_, err := _createDepartmentInfo(makeInfo.originOrgId, makeInfo.originUserIds[i], "fs", tx)
		if err != nil {
			return err
		}

		_, err = _createDepartmentInfo(makeInfo.newOrgId, makeInfo.newUserIds[i], "ding", tx)
		if err != nil {
			return err
		}

	}

	return nil
}

func _createDepartmentInfo(orgId, userId int64, sourceChannel string, tx sqlbuilder.Tx) (int64, error) {
	departmentInfo := &po.PpmOrgDepartment{
		Id:             snowflake.Id(),
		OrgId:          orgId,
		Name:           cast.ToString(orgId),
		Code:           cast.ToString(orgId),
		SourcePlatform: sourceChannel,
		SourceChannel:  sourceChannel,
		Creator:        userId,
		Updator:        userId,
	}
	departmentOutInfo := &po.PpmOrgDepartmentOutInfo{
		Id:                 snowflake.Id(),
		OrgId:              orgId,
		DepartmentId:       departmentInfo.Id,
		SourcePlatform:     sourceChannel,
		SourceChannel:      sourceChannel,
		OutOrgDepartmentId: cast.ToString(orgId),
		Name:               "xxxx",
		Creator:            userId,
		Updator:            userId,
	}
	departmentUserInfo := &po.PpmOrgUserDepartment{
		Id:           snowflake.Id(),
		OrgId:        orgId,
		UserId:       userId,
		DepartmentId: departmentInfo.Id,
		Creator:      userId,
		Updator:      userId,
	}

	err := mysql.TransInsert(tx, departmentInfo)
	if err != nil {
		return 0, err
	}

	err = mysql.TransInsert(tx, departmentOutInfo)
	if err != nil {
		return 0, err
	}

	err = mysql.TransInsert(tx, departmentUserInfo)
	if err != nil {
		return 0, err
	}

	return departmentInfo.Id, nil
}

// 删除假数据
func _delFakeInfo(makeInfo *_makeInfo) error {
	tables := []string{consts.TableOrganizationOutInfo, consts.TableUserOrganization, consts.TableUser, consts.TableUserOutInfo, consts.TableDepartment, consts.TableDepartmentOutInfo, consts.TableUserDepartment}
	dbErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		ids := strings.Join(cast.ToStringSlice([]int64{makeInfo.originOrgId, makeInfo.newOrgId}), ",")
		_, err := tx.DeleteFrom(consts.TableOrganization).Where(fmt.Sprintf("id in(%v)", ids)).Exec()
		if err != nil {
			return err
		}

		for _, table := range tables {
			err = _delByOrgIds(table, []int64{makeInfo.originOrgId, makeInfo.newOrgId}, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return dbErr
}

func _delByOrgIds(tableName string, ids []int64, tx sqlbuilder.Tx) error {
	idsStr := strings.Join(cast.ToStringSlice(ids), ",")
	_, err := tx.DeleteFrom(tableName).Where(fmt.Sprintf("org_id in(%v)", idsStr)).Exec()
	return err
}
