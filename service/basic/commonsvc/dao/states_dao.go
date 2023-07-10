package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertStates(po po.PpmCmmStates, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("States dao Insert err %v", err)
	}
	return nil
}

func InsertStatesBatch(pos []po.PpmCmmStates, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableStates).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableStates).Batch(len(pos))
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

func UpdateStates(po po.PpmCmmStates, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("States dao Update err %v", err)
	}
	return err
}

func UpdateStatesById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateStatesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateStatesByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateStatesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateStatesByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableStates, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableStates, cond, upd)
	}
	if err != nil {
		log.Errorf("States dao Update err %v", err)
	}
	return mod, err
}

func DeleteStatesById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateStatesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteStatesByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateStatesByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectStatesById(id int64) (*po.PpmCmmStates, error) {
	po := &po.PpmCmmStates{}
	err := mysql.SelectById(consts.TableStates, id, po)
	if err != nil {
		log.Errorf("States dao SelectById err %v", err)
	}
	return po, err
}

func SelectStatesByIdAndOrg(id int64, orgId int64) (*po.PpmCmmStates, error) {
	po := &po.PpmCmmStates{}
	err := mysql.SelectOneByCond(consts.TableStates, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("States dao Select err %v", err)
	}
	return po, err
}

func SelectStates(cond db.Cond) (*[]po.PpmCmmStates, error) {
	pos := &[]po.PpmCmmStates{}
	err := mysql.SelectAllByCond(consts.TableStates, cond, pos)
	if err != nil {
		log.Errorf("States dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneStates(cond db.Cond) (*po.PpmCmmStates, error) {
	po := &po.PpmCmmStates{}
	err := mysql.SelectOneByCond(consts.TableStates, cond, po)
	if err != nil {
		log.Errorf("States dao Select err %v", err)
	}
	return po, err
}

func SelectStatesByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmCmmStates, uint64, error) {
	pos := &[]po.PpmCmmStates{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableStates, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("States dao SelectPage err %v", err)
	}
	return pos, total, err
}
