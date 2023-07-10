package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertArea(po po.PpmCmmArea, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Area dao Insert err %v", err)
	}
	return nil
}

func InsertAreaBatch(pos []po.PpmCmmArea, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableArea).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableArea).Batch(len(pos))
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

func UpdateArea(po po.PpmCmmArea, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Area dao Update err %v", err)
	}
	return err
}

func UpdateAreaById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateAreaByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateAreaByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateAreaByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateAreaByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableArea, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableArea, cond, upd)
	}
	if err != nil {
		log.Errorf("Area dao Update err %v", err)
	}
	return mod, err
}

func DeleteAreaById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateAreaByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteAreaByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateAreaByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectAreaById(id int64) (*po.PpmCmmArea, error) {
	po := &po.PpmCmmArea{}
	err := mysql.SelectById(consts.TableArea, id, po)
	if err != nil {
		log.Errorf("Area dao SelectById err %v", err)
	}
	return po, err
}

func SelectAreaByIdAndOrg(id int64, orgId int64) (*po.PpmCmmArea, error) {
	po := &po.PpmCmmArea{}
	err := mysql.SelectOneByCond(consts.TableArea, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Area dao Select err %v", err)
	}
	return po, err
}

func SelectArea(cond db.Cond) (*[]po.PpmCmmArea, error) {
	pos := &[]po.PpmCmmArea{}
	err := mysql.SelectAllByCond(consts.TableArea, cond, pos)
	if err != nil {
		log.Errorf("Area dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneArea(cond db.Cond) (*po.PpmCmmArea, error) {
	po := &po.PpmCmmArea{}
	err := mysql.SelectOneByCond(consts.TableArea, cond, po)
	if err != nil {
		log.Errorf("Area dao Select err %v", err)
	}
	return po, err
}

func SelectAreaByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmCmmArea, uint64, error) {
	pos := &[]po.PpmCmmArea{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableArea, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Area dao SelectPage err %v", err)
	}
	return pos, total, err
}
