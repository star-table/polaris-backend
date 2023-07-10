package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertCountries(po po.PpmCmmCountries, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Countries dao Insert err %v", err)
	}
	return nil
}

func InsertCountriesBatch(pos []po.PpmCmmCountries, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableCountries).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableCountries).Batch(len(pos))
	}
	go func() {
		defer batch.Done()
		for i := range pos {
			batch.Values(pos[i])
		}
	}()
	err = batch.Wait()
	if err != nil {
		log.Errorf("Iteration dao InsertBatch err %v", err)
		return err
	}
	return nil
}

func UpdateCountries(po po.PpmCmmCountries, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Countries dao Update err %v", err)
	}
	return err
}

func UpdateCountriesById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateCountriesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateCountriesByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateCountriesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateCountriesByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableCountries, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableCountries, cond, upd)
	}
	if err != nil {
		log.Errorf("Countries dao Update err %v", err)
	}
	return mod, err
}

func DeleteCountriesById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateCountriesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteCountriesByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateCountriesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectCountriesById(id int64) (*po.PpmCmmCountries, error) {
	po := &po.PpmCmmCountries{}
	err := mysql.SelectById(consts.TableCountries, id, po)
	if err != nil {
		log.Errorf("Countries dao SelectById err %v", err)
	}
	return po, err
}

func SelectCountriesByIdAndOrg(id int64, orgId int64) (*po.PpmCmmCountries, error) {
	po := &po.PpmCmmCountries{}
	err := mysql.SelectOneByCond(consts.TableCountries, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Countries dao Select err %v", err)
	}
	return po, err
}

func SelectCountries(cond db.Cond) (*[]po.PpmCmmCountries, error) {
	pos := &[]po.PpmCmmCountries{}
	err := mysql.SelectAllByCond(consts.TableCountries, cond, pos)
	if err != nil {
		log.Errorf("Countries dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneCountries(cond db.Cond) (*po.PpmCmmCountries, error) {
	po := &po.PpmCmmCountries{}
	err := mysql.SelectOneByCond(consts.TableCountries, cond, po)
	if err != nil {
		log.Errorf("Countries dao Select err %v", err)
	}
	return po, err
}

func SelectCountriesByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmCmmCountries, uint64, error) {
	pos := &[]po.PpmCmmCountries{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableCountries, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Countries dao SelectPage err %v", err)
	}
	return pos, total, err
}
