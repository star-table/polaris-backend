package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertIterationStat(po po.PpmStaIterationStat, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("IterationStat dao Insert err %v", err)
	}
	return nil
}

//func InsertIterationStatBatch(pos []po.PpmStaIterationStat, tx ...sqlbuilder.Tx) error {
//	var err error = nil
//	if tx != nil && len(tx) > 0 {
//		batch := tx[0].InsertInto(consts.TableIterationStat).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
//		//go iterationStatBatchDone(pos,batch)
//		go mysql.BatchDone(slice.ToSlice(pos), batch)
//
//		err = batch.Wait()
//		if err != nil {
//			return err
//		}
//	} else {
//		conn, err := mysql.GetConnect()
//		if err != nil {
//			return err
//		}
//		batch := conn.InsertInto(consts.TableIterationStat).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
//		//go iterationStatBatchDone(pos,batch)
//		go mysql.BatchDone(slice.ToSlice(pos), batch)
//
//		err = batch.Wait()
//		if err != nil {
//			return err
//		}
//	}
//	if err != nil {
//		log.Errorf("IterationStat dao InsertBatch err %v", err)
//	}
//	return nil
//}

func InsertIterationStatBatch(pos []po.PpmStaIterationStat, tx ...sqlbuilder.Tx) error {
	var err error = nil

	isTx := tx != nil && len(tx) > 0

	var batch *sqlbuilder.BatchInserter

	if !isTx {
		//没有事务
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}

		batch = conn.InsertInto(consts.TableIterationStat).Batch(len(pos))
	}

	if batch == nil {
		batch = tx[0].InsertInto(consts.TableIterationStat).Batch(len(pos))
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

func UpdateIterationStat(po po.PpmStaIterationStat, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("IterationStat dao Update err %v", err)
	}
	return err
}

func UpdateIterationStatById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIterationStatByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIterationStatByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIterationStatByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIterationStatByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableIterationStat, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableIterationStat, cond, upd)
	}
	if err != nil {
		log.Errorf("IterationStat dao Update err %v", err)
	}
	return mod, err
}

func DeleteIterationStatById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIterationStatByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteIterationStatByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIterationStatByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectIterationStatById(id int64) (*po.PpmStaIterationStat, error) {
	po := &po.PpmStaIterationStat{}
	err := mysql.SelectById(consts.TableIterationStat, id, po)
	if err != nil {
		log.Errorf("IterationStat dao SelectById err %v", err)
	}
	return po, err
}

func SelectIterationStatByIdAndOrg(id int64, orgId int64) (*po.PpmStaIterationStat, error) {
	po := &po.PpmStaIterationStat{}
	err := mysql.SelectOneByCond(consts.TableIterationStat, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("IterationStat dao Select err %v", err)
	}
	return po, err
}

func SelectIterationStat(cond db.Cond) (*[]po.PpmStaIterationStat, error) {
	pos := &[]po.PpmStaIterationStat{}
	err := mysql.SelectAllByCond(consts.TableIterationStat, cond, pos)
	if err != nil {
		log.Errorf("IterationStat dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneIterationStat(cond db.Cond) (*po.PpmStaIterationStat, error) {
	po := &po.PpmStaIterationStat{}
	err := mysql.SelectOneByCond(consts.TableIterationStat, cond, po)
	if err != nil {
		log.Errorf("IterationStat dao Select err %v", err)
	}
	return po, err
}

func SelectIterationStatByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmStaIterationStat, uint64, error) {
	pos := &[]po.PpmStaIterationStat{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIterationStat, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("IterationStat dao SelectPage err %v", err)
	}
	return pos, total, err
}
