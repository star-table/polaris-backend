package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertCities(po po.PpmCmmCities, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Cities dao Insert err %v", err)
	}
	return nil
}

func InsertCitiesBatch(pos []po.PpmCmmCities, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableCities).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableCities).Batch(len(pos))
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

func UpdateCities(po po.PpmCmmCities, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Cities dao Update err %v", err)
	}
	return err
}

func UpdateCitiesById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateCitiesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateCitiesByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateCitiesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateCitiesByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableCities, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableCities, cond, upd)
	}
	if err != nil {
		log.Errorf("Cities dao Update err %v", err)
	}
	return mod, err
}

func DeleteCitiesById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateCitiesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteCitiesByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateCitiesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectCitiesById(id int64) (*po.PpmCmmCities, error) {
	po := &po.PpmCmmCities{}
	err := mysql.SelectById(consts.TableCities, id, po)
	if err != nil {
		log.Errorf("Cities dao SelectById err %v", err)
	}
	return po, err
}

func SelectCitiesByIdAndOrg(id int64, orgId int64) (*po.PpmCmmCities, error) {
	po := &po.PpmCmmCities{}
	err := mysql.SelectOneByCond(consts.TableCities, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Cities dao Select err %v", err)
	}
	return po, err
}

func SelectCities(cond db.Cond) (*[]po.PpmCmmCities, error) {
	pos := &[]po.PpmCmmCities{}
	err := mysql.SelectAllByCond(consts.TableCities, cond, pos)
	if err != nil {
		log.Errorf("Cities dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneCities(cond db.Cond) (*po.PpmCmmCities, error) {
	po := &po.PpmCmmCities{}
	err := mysql.SelectOneByCond(consts.TableCities, cond, po)
	if err != nil {
		log.Errorf("Cities dao Select err %v", err)
	}
	return po, err
}

func SelectCitiesByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmCmmCities, uint64, error) {
	pos := &[]po.PpmCmmCities{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableCities, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Cities dao SelectPage err %v", err)
	}
	return pos, total, err
}
