package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetShareView() shareViewInterface {
	return shareViewDao
}

type shareViewInterface interface {
	GetByViewId(orgId, userId, viewId int64) (*po.PpmShareView, error)

	GetByShareKey(shareKey string) (*po.PpmShareView, error)

	// Create 创建shareView
	Create(m *po.PpmShareView) error

	// Delete 删除shareView
	Delete(userId, viewId int64) error

	UpdatePassword(orgId, userId, viewId int64, password string) error

	ResetShareKeyAndPassword(orgId, userId, viewId int64, newShareKey string) error

	UpdateShareConfig(orgId, userId, viewId int64, config string) error
}

type shareView struct{}

var shareViewDao = &shareView{}

func (s *shareView) GetByViewId(orgId, userId, viewId int64) (*po.PpmShareView, error) {
	m := &po.PpmShareView{}
	err := mysql.SelectOneByCond(po.TableNameShareView, db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcViewId: viewId,
		consts.TcUserId: userId,
	}, m)

	return m, err
}

func (s *shareView) GetByShareKey(shareKey string) (*po.PpmShareView, error) {
	m := &po.PpmShareView{}
	err := mysql.SelectOneByCond(po.TableNameShareView, db.Cond{
		consts.TcShareKey: shareKey,
	}, m)

	return m, err
}

// Create 创建shareView
func (s *shareView) Create(m *po.PpmShareView) error {
	return mysql.Insert(m)
}

// Delete 删除shareView
func (s *shareView) Delete(userId, viewId int64) error {
	return mysql.TransX(func(tx sqlbuilder.Tx) error {
		_, err := tx.DeleteFrom(po.TableNameShareView).Where(consts.TcViewId+" = ? and "+consts.TcUserId+" = ? ", viewId, userId).Exec()
		return err
	})
}

func (s *shareView) UpdatePassword(orgId, userId, viewId int64, password string) error {
	return s.update(orgId, userId, viewId, mysql.Upd{
		consts.TcSharePassword: password,
	})
}

func (s *shareView) ResetShareKeyAndPassword(orgId, userId, viewId int64, newShareKey string) error {
	return s.update(orgId, userId, viewId, mysql.Upd{
		consts.TcShareKey:      newShareKey,
		consts.TcSharePassword: "",
	})
}

func (s *shareView) UpdateShareConfig(orgId, userId, viewId int64, config string) error {
	return s.update(orgId, userId, viewId, mysql.Upd{
		consts.TcConfig: config,
	})
}

func (s *shareView) update(orgId, userId, viewId int64, upd mysql.Upd) error {
	_, err := mysql.UpdateSmartWithCond(po.TableNameShareView, db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcViewId: viewId,
		consts.TcUserId: userId,
	}, upd,
	)

	return err
}
