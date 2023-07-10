package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/msgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertMessage(po po.PpmTakMessage, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Message dao Insert err %v", err)
	}
	return nil
}

func InsertMessageBatch(pos []po.PpmTakMessage, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableMessage).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableMessage).Batch(len(pos))
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

func UpdateMessage(po po.PpmTakMessage, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Message dao Update err %v", err)
	}
	return err
}

func UpdateMessageById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateMessageByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateMessageByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateMessageByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateMessageByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableMessage, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableMessage, cond, upd)
	}
	if err != nil {
		log.Errorf("Message dao Update err %v", err)
	}
	return mod, err
}

func DeleteMessageById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateMessageByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteMessageByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateMessageByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectMessageById(id int64) (*po.PpmTakMessage, error) {
	po := &po.PpmTakMessage{}
	err := mysql.SelectById(consts.TableMessage, id, po)
	if err != nil {
		log.Errorf("Message dao SelectById err %v", err)
	}
	return po, err
}

func SelectMessageByIdAndOrg(id int64, orgId int64) (*po.PpmTakMessage, error) {
	po := &po.PpmTakMessage{}
	err := mysql.SelectOneByCond(consts.TableMessage, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Message dao Select err %v", err)
	}
	return po, err
}

func SelectMessage(cond db.Cond) (*[]po.PpmTakMessage, error) {
	pos := &[]po.PpmTakMessage{}
	err := mysql.SelectAllByCond(consts.TableMessage, cond, pos)
	if err != nil {
		log.Errorf("Message dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneMessage(cond db.Cond) (*po.PpmTakMessage, error) {
	po := &po.PpmTakMessage{}
	err := mysql.SelectOneByCond(consts.TableMessage, cond, po)
	if err != nil {
		log.Errorf("Message dao Select err %v", err)
	}
	return po, err
}

func SelectMessageByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmTakMessage, uint64, error) {
	pos := &[]po.PpmTakMessage{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableMessage, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Message dao SelectPage err %v", err)
	}
	return pos, total, err
}
