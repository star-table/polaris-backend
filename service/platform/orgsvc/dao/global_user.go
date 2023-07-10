package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetGlobalUser() GlobalUserInterface {
	return globalUserDao
}

type GlobalUserInterface interface {
	// Create 创建globalUser
	Create(m *po.PpmOrgGlobalUser, tx sqlbuilder.Tx) error
	// Delete 删除globalUser
	Delete(id int64, tx sqlbuilder.Tx) error
	// GetGlobalUserByMobile 根据手机号获取globalUser
	GetGlobalUserByMobile(mobile string) (*po.PpmOrgGlobalUser, error)
	// GetGlobalUsersByMobiles 根据手机号列表获取globalUser列表
	GetGlobalUsersByMobiles(mobiles []string) ([]*po.PpmOrgGlobalUser, error)
	// GetGlobalUserById 通过id获取globalUser
	GetGlobalUserById(globalId int64) (*po.PpmOrgGlobalUser, error)
	// GetGlobalUserByWeiXinId 通过id获取globalUser
	GetGlobalUserByWeiXinId(openId string) (*po.PpmOrgGlobalUser, error)
	// GetMobilesMapByIds 获取globalId对应mobile的map
	GetMobilesMapByIds(globalIds []int64) (map[int64]string, error)
	// UpdateGlobalLastLoginInfo 更新最后一次登陆的组织
	UpdateGlobalLastLoginInfo(globalUserId, userId, orgId int64) error
	// UpdateMobile 替换手机号
	UpdateMobile(globalUserId int64, mobile string, tx sqlbuilder.Tx) error

	UpdateWeiXinOpenId(tx sqlbuilder.Tx, mobile string, openId string) error

	UnbindWeiXinOpenId(tx sqlbuilder.Tx, globalUserId int64) error
}

type globalUser struct{}

var globalUserDao = &globalUser{}

func (g globalUser) Create(m *po.PpmOrgGlobalUser, tx sqlbuilder.Tx) error {
	return mysql.TransInsert(tx, m)
}

func (g globalUser) Delete(id int64, tx sqlbuilder.Tx) error {
	_, err := tx.DeleteFrom(po.TableNamePpmOrgGlobalUser).Where(consts.TcId+" = ?", id).Exec()
	return err
}

func (g globalUser) GetGlobalUserByMobile(mobile string) (*po.PpmOrgGlobalUser, error) {
	m := &po.PpmOrgGlobalUser{}
	err := mysql.SelectOneByCond(po.TableNamePpmOrgGlobalUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcMobile:   mobile,
	}, m)

	return m, err
}

func (g globalUser) GetGlobalUserByWeiXinId(openId string) (*po.PpmOrgGlobalUser, error) {
	m := &po.PpmOrgGlobalUser{}
	err := mysql.SelectOneByCond(po.TableNamePpmOrgGlobalUser, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcWeiXinOpenId: openId,
	}, m)

	return m, err
}

func (g globalUser) GetGlobalUsersByMobiles(mobiles []string) ([]*po.PpmOrgGlobalUser, error) {
	ms := make([]*po.PpmOrgGlobalUser, 0, len(mobiles))
	err := mysql.SelectAllByCond(po.TableNamePpmOrgGlobalUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcMobile:   db.In(mobiles),
	}, &ms)

	return ms, err
}

func (g globalUser) GetMobilesMapByIds(globalIds []int64) (map[int64]string, error) {
	ms := make([]*po.PpmOrgGlobalUser, 0, len(globalIds))
	err := mysql.SelectAllByCond(po.TableNamePpmOrgGlobalUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(globalIds),
	}, &ms)
	if err != nil {
		return nil, err
	}

	result := make(map[int64]string, len(ms))
	for _, m := range ms {
		result[m.Id] = m.Mobile
	}

	return result, nil
}

func (g globalUser) GetGlobalUserById(globalId int64) (*po.PpmOrgGlobalUser, error) {
	m := &po.PpmOrgGlobalUser{}
	err := mysql.SelectOneByCond(po.TableNamePpmOrgGlobalUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       globalId,
	}, m)

	return m, err
}

// UpdateGlobalLastLoginInfo 更新下最后一次登陆记录的useId和orgId
func (g globalUser) UpdateGlobalLastLoginInfo(globalUserId, userId, orgId int64) error {
	_, err := mysql.UpdateSmartWithCond(po.TableNamePpmOrgGlobalUser, db.Cond{
		consts.TcId: globalUserId,
	}, mysql.Upd{
		consts.TcLastLoginOrgId:  orgId,
		consts.TcLastLoginUserId: userId,
	},
	)

	return err
}

func (g globalUser) UpdateMobile(globalUserId int64, mobile string, tx sqlbuilder.Tx) error {
	_, err := mysql.TransUpdateSmartWithCond(tx, po.TableNamePpmOrgGlobalUser, db.Cond{
		consts.TcId: globalUserId,
	}, mysql.Upd{
		consts.TcMobile: mobile,
	},
	)

	return err
}

func (g globalUser) UpdateWeiXinOpenId(tx sqlbuilder.Tx, mobile string, openId string) error {
	_, err := mysql.TransUpdateSmartWithCond(tx, po.TableNamePpmOrgGlobalUser, db.Cond{
		consts.TcMobile: mobile,
	}, mysql.Upd{
		consts.TcWeiXinOpenId: openId,
	},
	)

	return err
}

func (g globalUser) UnbindWeiXinOpenId(tx sqlbuilder.Tx, globalUserId int64) error {
	_, err := mysql.TransUpdateSmartWithCond(tx, po.TableNamePpmOrgGlobalUser, db.Cond{
		consts.TcId: globalUserId,
	}, mysql.Upd{
		consts.TcWeiXinOpenId: "",
	},
	)

	return err
}
