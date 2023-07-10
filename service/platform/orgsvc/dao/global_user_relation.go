package dao

import (
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetGlobalUserRelation() GlobalUserRelationInterface {
	return globalUserRelationDao
}

type GlobalUserRelationInterface interface {
	// CreateRelations 创建关联关系
	CreateRelations(relations []*po.PpmOrgGlobalUserRelation, tx sqlbuilder.Tx) error
	// DeleteRelationsByGlobalUserId 删除关联关系
	DeleteRelationsByGlobalUserId(globalUserId int64, tx sqlbuilder.Tx) error
	// DeleteRelationsByUserId 根据userId删除关联关系
	DeleteRelationsByUserId(userId int64, tx sqlbuilder.Tx) error
	// GetGlobalUserIdByUserId 根据userId获取globalId
	GetGlobalUserIdByUserId(userId int64) (int64, error)
	// GetGlobalUserIdsMapByUserIds 根据userIds获取globalIds
	GetGlobalUserIdsMapByUserIds(userIds []int64) (map[int64][]int64, error)
	// GetUserIdsByGlobalUserId 获取一个globalUser下所有绑定的用户
	GetUserIdsByGlobalUserId(globalUserId int64) ([]int64, error)
	// GetUserIdsByUserId 根据userId，获取其他一起绑定的userId
	GetUserIdsByUserId(userId int64) ([]int64, error)
	// GetRelationsByGlobalUserIds 获取关系列表
	GetRelationsByGlobalUserIds(globalUserIds []int64) ([]*po.PpmOrgGlobalUserRelation, error)
	// UpdateUsersToNewGlobalId 替换绑定关系
	UpdateUsersToNewGlobalId(oldGlobalId, newGlobalId int64, tx sqlbuilder.Tx) error
}

type globalUserRelation struct{}

var globalUserRelationDao = &globalUserRelation{}

func (g globalUserRelation) CreateRelations(relations []*po.PpmOrgGlobalUserRelation, tx sqlbuilder.Tx) error {
	return mysql.TransBatchInsert(tx, &po.PpmOrgGlobalUserRelation{}, slice.ToSlice(relations))
}

func (g globalUserRelation) UpdateUsersToNewGlobalId(oldGlobalId, newGlobalId int64, tx sqlbuilder.Tx) error {
	_, err := mysql.TransUpdateSmartWithCond(tx, po.TableNamePpmOrgGlobalUserRelation, db.Cond{
		consts.TcGlobalUserId: oldGlobalId,
	}, mysql.Upd{
		consts.TcGlobalUserId: newGlobalId,
	},
	)

	return err
}

// DeleteRelationsByGlobalUserId 解除一个g_user下所有绑定的用户
func (g globalUserRelation) DeleteRelationsByGlobalUserId(globalUserId int64, tx sqlbuilder.Tx) error {
	_, err := tx.DeleteFrom(po.TableNamePpmOrgGlobalUserRelation).Where(consts.TcGlobalUserId+" = ?", globalUserId).Exec()
	return err
}

// DeleteRelationsByUserId 解除单一渠道用户的绑定
func (g globalUserRelation) DeleteRelationsByUserId(userId int64, tx sqlbuilder.Tx) error {
	_, err := tx.DeleteFrom(po.TableNamePpmOrgGlobalUserRelation).Where(consts.TcUserId+" = ?", userId).Exec()
	return err
}

func (g globalUserRelation) GetGlobalUserIdByUserId(userId int64) (int64, error) {
	m := &po.PpmOrgGlobalUserRelation{}
	err := mysql.SelectOneByCond(po.TableNamePpmOrgGlobalUserRelation, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUserId:   userId,
	}, m)

	return m.GlobalUserId, err
}

func (g globalUserRelation) GetGlobalUserIdsMapByUserIds(userIds []int64) (map[int64][]int64, error) {
	ms := make([]*po.PpmOrgGlobalUserRelation, 0, len(userIds))
	err := mysql.SelectAllByCond(po.TableNamePpmOrgGlobalUserRelation, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUserId:   db.In(userIds),
	}, &ms)

	result := make(map[int64][]int64, len(ms))
	for _, m := range ms {
		result[m.GlobalUserId] = append(result[m.GlobalUserId], m.UserId)
	}

	return result, err
}

func (g globalUserRelation) GetUserIdsByUserId(userId int64) ([]int64, error) {
	globalUserId, err := g.GetGlobalUserIdByUserId(userId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return []int64{userId}, nil
		}
		return nil, err
	}

	userIds, err := g.GetUserIdsByGlobalUserId(globalUserId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return []int64{userId}, nil
		}
		return nil, err
	}

	return userIds, nil
}

func (g globalUserRelation) GetUserIdsByGlobalUserId(globalUserId int64) ([]int64, error) {
	ms := make([]*po.PpmOrgGlobalUserRelation, 0, 1)
	err := mysql.SelectAllByCond(po.TableNamePpmOrgGlobalUserRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcGlobalUserId: globalUserId,
	}, &ms)

	return g.getUserIds(ms), err
}

func (g globalUserRelation) GetRelationsByGlobalUserIds(globalUserIds []int64) ([]*po.PpmOrgGlobalUserRelation, error) {
	ms := make([]*po.PpmOrgGlobalUserRelation, 0, 1)
	err := mysql.SelectAllByCond(po.TableNamePpmOrgGlobalUserRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcGlobalUserId: db.In(globalUserIds),
	}, &ms)

	return ms, err
}

func (g globalUserRelation) getUserIds(ms []*po.PpmOrgGlobalUserRelation) []int64 {
	userIds := make([]int64, 0, len(ms))
	for _, m := range ms {
		userIds = append(userIds, m.UserId)
	}

	return userIds
}
