package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertIteration(po po.PpmPriIteration, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Iteration dao Insert err %v", err)
	}
	return nil
}

//func InsertIterationBatch(pos []po.PpmPriIteration, tx ...sqlbuilder.Tx) error {
//	var err error = nil
//	if tx != nil && len(tx) > 0 {
//		batch := tx[0].InsertInto(consts.TableIteration).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
//		//go iterationBatchDone(pos,batch)
//
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
//		batch := conn.InsertInto(consts.TableIteration).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
//		//go iterationBatchDone(pos,batch)
//
//		go mysql.BatchDone(slice.ToSlice(pos), batch)
//
//		err = batch.Wait()
//		if err != nil {
//			return err
//		}
//	}
//	if err != nil {
//		log.Errorf("Iteration dao InsertBatch err %v", err)
//	}
//	return nil
//}

func InsertIterationBatch(pos []po.PpmPriIteration, tx ...sqlbuilder.Tx) error {
	var err error = nil

	isTx := tx != nil && len(tx) > 0

	var batch *sqlbuilder.BatchInserter

	if !isTx {
		//没有事务
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}

		batch = conn.InsertInto(consts.TableIteration).Batch(len(pos))
	}

	if batch == nil {
		batch = tx[0].InsertInto(consts.TableIteration).Batch(len(pos))
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

func UpdateIteration(po po.PpmPriIteration, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Iteration dao Update err %v", err)
	}
	return err
}

func UpdateIterationById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIterationByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIterationByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIterationByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIterationByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableIteration, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableIteration, cond, upd)
	}
	if err != nil {
		log.Errorf("Iteration dao Update err %v", err)
	}
	return mod, err
}

func DeleteIterationById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIterationByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteIterationByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIterationByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectIterationById(id int64) (*po.PpmPriIteration, error) {
	po := &po.PpmPriIteration{}
	err := mysql.SelectById(consts.TableIteration, id, po)
	if err != nil {
		log.Errorf("Iteration dao SelectById err %v", err)
	}
	return po, err
}

func SelectIterationByIdAndOrg(id int64, orgId int64) (*po.PpmPriIteration, error) {
	po := &po.PpmPriIteration{}
	err := mysql.SelectOneByCond(consts.TableIteration, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Iteration dao Select err %v", err)
	}
	return po, err
}

func SelectIteration(cond db.Cond) (*[]po.PpmPriIteration, error) {
	pos := &[]po.PpmPriIteration{}
	err := mysql.SelectAllByCond(consts.TableIteration, cond, pos)
	if err != nil {
		log.Errorf("Iteration dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneIteration(cond db.Cond) (*po.PpmPriIteration, error) {
	po := &po.PpmPriIteration{}
	err := mysql.SelectOneByCond(consts.TableIteration, cond, po)
	if err != nil {
		log.Errorf("Iteration dao Select err %v", err)
	}
	return po, err
}

func SelectIterationByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmPriIteration, uint64, error) {
	pos := &[]po.PpmPriIteration{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIteration, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Iteration dao SelectPage err %v", err)
	}
	return pos, total, err
}

func JudgeIterationIsExist(orgId, id int64) bool {
	iteration := &po.PpmPriIteration{}
	err := mysql.SelectOneByCond(consts.TableIteration, db.Cond{
		"id":        id,
		"org_id":    orgId,
		"is_delete": consts.AppIsNoDelete,
	}, iteration)
	if err != nil {
		return false
	}

	return true
}
