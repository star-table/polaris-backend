package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertMobilePrefix(po po.PpmCmmMobilePrefix, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("MobilePrefix dao Insert err %v", err)
	}
	return nil
}

func InsertMobilePrefixBatch(pos []po.PpmCmmMobilePrefix, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableMobilePrefix).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableMobilePrefix).Batch(len(pos))
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

func UpdateMobilePrefix(po po.PpmCmmMobilePrefix, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("MobilePrefix dao Update err %v", err)
	}
	return err
}

func UpdateMobilePrefixById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateMobilePrefixByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateMobilePrefixByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateMobilePrefixByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateMobilePrefixByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableMobilePrefix, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableMobilePrefix, cond, upd)
	}
	if err != nil {
		log.Errorf("MobilePrefix dao Update err %v", err)
	}
	return mod, err
}

func DeleteMobilePrefixById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateMobilePrefixByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteMobilePrefixByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateMobilePrefixByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectMobilePrefixById(id int64) (*po.PpmCmmMobilePrefix, error) {
	po := &po.PpmCmmMobilePrefix{}
	err := mysql.SelectById(consts.TableMobilePrefix, id, po)
	if err != nil {
		log.Errorf("MobilePrefix dao SelectById err %v", err)
	}
	return po, err
}

func SelectMobilePrefixByIdAndOrg(id int64, orgId int64) (*po.PpmCmmMobilePrefix, error) {
	po := &po.PpmCmmMobilePrefix{}
	err := mysql.SelectOneByCond(consts.TableMobilePrefix, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("MobilePrefix dao Select err %v", err)
	}
	return po, err
}

func SelectMobilePrefix(cond db.Cond) (*[]po.PpmCmmMobilePrefix, error) {
	pos := &[]po.PpmCmmMobilePrefix{}
	err := mysql.SelectAllByCond(consts.TableMobilePrefix, cond, pos)
	if err != nil {
		log.Errorf("MobilePrefix dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneMobilePrefix(cond db.Cond) (*po.PpmCmmMobilePrefix, error) {
	po := &po.PpmCmmMobilePrefix{}
	err := mysql.SelectOneByCond(consts.TableMobilePrefix, cond, po)
	if err != nil {
		log.Errorf("MobilePrefix dao Select err %v", err)
	}
	return po, err
}

func SelectMobilePrefixByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmCmmMobilePrefix, uint64, error) {
	pos := &[]po.PpmCmmMobilePrefix{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableMobilePrefix, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("MobilePrefix dao SelectPage err %v", err)
	}
	return pos, total, err
}
