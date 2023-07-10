package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertRegions(po po.PpmCmmRegions, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Regions dao Insert err %v", err)
	}
	return nil
}

func InsertRegionsBatch(pos []po.PpmCmmRegions, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableRegions).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableRegions).Batch(len(pos))
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

func UpdateRegions(po po.PpmCmmRegions, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Regions dao Update err %v", err)
	}
	return err
}

func UpdateRegionsById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateRegionsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateRegionsByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateRegionsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateRegionsByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableRegions, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableRegions, cond, upd)
	}
	if err != nil {
		log.Errorf("Regions dao Update err %v", err)
	}
	return mod, err
}

func DeleteRegionsById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateRegionsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteRegionsByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateRegionsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectRegionsById(id int64) (*po.PpmCmmRegions, error) {
	po := &po.PpmCmmRegions{}
	err := mysql.SelectById(consts.TableRegions, id, po)
	if err != nil {
		log.Errorf("Regions dao SelectById err %v", err)
	}
	return po, err
}

func SelectRegionsByIdAndOrg(id int64, orgId int64) (*po.PpmCmmRegions, error) {
	po := &po.PpmCmmRegions{}
	err := mysql.SelectOneByCond(consts.TableRegions, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Regions dao Select err %v", err)
	}
	return po, err
}

func SelectRegions(cond db.Cond) (*[]po.PpmCmmRegions, error) {
	pos := &[]po.PpmCmmRegions{}
	err := mysql.SelectAllByCond(consts.TableRegions, cond, pos)
	if err != nil {
		log.Errorf("Regions dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneRegions(cond db.Cond) (*po.PpmCmmRegions, error) {
	po := &po.PpmCmmRegions{}
	err := mysql.SelectOneByCond(consts.TableRegions, cond, po)
	if err != nil {
		log.Errorf("Regions dao Select err %v", err)
	}
	return po, err
}

func SelectRegionsByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmCmmRegions, uint64, error) {
	pos := &[]po.PpmCmmRegions{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableRegions, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Regions dao SelectPage err %v", err)
	}
	return pos, total, err
}
