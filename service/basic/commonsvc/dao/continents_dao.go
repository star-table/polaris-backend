package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertContinents(po po.PpmCmmContinents, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Continents dao Insert err %v", err)
	}
	return nil
}

func InsertContinentsBatch(pos []po.PpmCmmContinents, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableContinents).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableContinents).Batch(len(pos))
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

func UpdateContinents(po po.PpmCmmContinents, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Continents dao Update err %v", err)
	}
	return err
}

func UpdateContinentsById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateContinentsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateContinentsByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateContinentsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateContinentsByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableContinents, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableContinents, cond, upd)
	}
	if err != nil {
		log.Errorf("Continents dao Update err %v", err)
	}
	return mod, err
}

func DeleteContinentsById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateContinentsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteContinentsByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateContinentsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectContinentsById(id int64) (*po.PpmCmmContinents, error) {
	po := &po.PpmCmmContinents{}
	err := mysql.SelectById(consts.TableContinents, id, po)
	if err != nil {
		log.Errorf("Continents dao SelectById err %v", err)
	}
	return po, err
}

func SelectContinentsByIdAndOrg(id int64, orgId int64) (*po.PpmCmmContinents, error) {
	po := &po.PpmCmmContinents{}
	err := mysql.SelectOneByCond(consts.TableContinents, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Continents dao Select err %v", err)
	}
	return po, err
}

func SelectContinents(cond db.Cond) (*[]po.PpmCmmContinents, error) {
	pos := &[]po.PpmCmmContinents{}
	err := mysql.SelectAllByCond(consts.TableContinents, cond, pos)
	if err != nil {
		log.Errorf("Continents dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneContinents(cond db.Cond) (*po.PpmCmmContinents, error) {
	po := &po.PpmCmmContinents{}
	err := mysql.SelectOneByCond(consts.TableContinents, cond, po)
	if err != nil {
		log.Errorf("Continents dao Select err %v", err)
	}
	return po, err
}

func SelectContinentsByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmCmmContinents, uint64, error) {
	pos := &[]po.PpmCmmContinents{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableContinents, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Continents dao SelectPage err %v", err)
	}
	return pos, total, err
}
