package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetThirdLogin() ThirdLoginInterface {
	return thirdLoginDao
}

type ThirdLoginInterface interface {
	GetThirdLoginInfo(thirdUserId string) (*po.PpmOrgThirdLoginUser, error)
	GetThirdLoginInfoByUserIds(userIds []int64) ([]*po.PpmOrgThirdLoginUser, error)
	AddThirdLoginInfo(tx sqlbuilder.Tx, info *po.PpmOrgThirdLoginUser) error
	UnBindThirdLogin(tx sqlbuilder.Tx, thirdId string) error
}

type thirdLogin struct{}

var thirdLoginDao = &thirdLogin{}

func (t *thirdLogin) GetThirdLoginInfo(thirdUserId string) (*po.PpmOrgThirdLoginUser, error) {
	m := &po.PpmOrgThirdLoginUser{}
	err := mysql.SelectOneByCond(po.TableNamePpmOrgThirdLoginUser, db.Cond{
		consts.TcOutUserId:      thirdUserId,
		consts.TcSourcePlatForm: consts.AppSourcePlatformPersonWeixin,
	}, m)

	return m, err
}

func (t *thirdLogin) AddThirdLoginInfo(tx sqlbuilder.Tx, info *po.PpmOrgThirdLoginUser) error {
	return mysql.TransInsert(tx, info)
}

func (t *thirdLogin) UnBindThirdLogin(tx sqlbuilder.Tx, TcThirdUserId string) error {
	_, err := tx.DeleteFrom(po.TableNamePpmOrgThirdLoginUser).Where(consts.TcThirdUserId+" = ? ", TcThirdUserId).Exec()
	if err != nil {
		log.Errorf("[UnBindThirdLogin] delete err:%v", err)
		return errs.MysqlOperateError
	}
	return nil
}

func (t *thirdLogin) GetThirdLoginInfoByUserIds(userIds []int64) ([]*po.PpmOrgThirdLoginUser, error) {
	m := make([]*po.PpmOrgThirdLoginUser, 0, 1)
	err := mysql.SelectAllByCond(po.TableNamePpmOrgThirdLoginUser, db.Cond{
		consts.TcUserId: db.In(userIds),
	}, &m)

	return m, err
}
