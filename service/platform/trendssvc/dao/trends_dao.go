package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/trendssvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertTrends(po po.PpmTreTrends, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Trends dao Insert err %v", err)
	}
	return nil
}

//func InsertTrendsBatch(pos []po.PpmTreTrends, tx ...sqlbuilder.Tx) error{
//	var err error = nil
//	if tx != nil && len(tx) > 0{
//		batch := tx[0].InsertInto(consts.TableTrends).Batch(len(pos))
//		go func() {
//			defer batch.Done()
//			for i := range pos {
//				batch.Values(pos[i])
//			}
//		}()
//		err = batch.Wait()
//		if err != nil {
//			return err
//		}
//	}else{
//		conn, err := mysql.GetConnect()
//		if err != nil {
//			return err
//		}
//		batch := conn.InsertInto(consts.TableTrends).Batch(len(pos))
//		go func() {
//			defer batch.Done()
//			for i := range pos {
//				batch.Values(pos[i])
//			}
//		}()
//		err = batch.Wait()
//		if err != nil {
//			return err
//		}
//	}
//	if err != nil{
//		log.Errorf("Trends dao InsertBatch err %v", err)
//	}
//	return nil
//}

func InsertTrendsBatch(pos []po.PpmTreTrends, tx ...sqlbuilder.Tx) error {
	var err error = nil

	isTx := tx != nil && len(tx) > 0

	var batch *sqlbuilder.BatchInserter

	if !isTx {
		//没有事务
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}

		batch = conn.InsertInto(consts.TableTrends).Batch(len(pos))
	}

	if batch == nil {
		batch = tx[0].InsertInto(consts.TableTrends).Batch(len(pos))
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

func UpdateTrends(po po.PpmTreTrends, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Trends dao Update err %v", err)
	}
	return err
}

func UpdateTrendsById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateTrendsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateTrendsByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateTrendsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateTrendsByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableTrends, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableTrends, cond, upd)
	}
	if err != nil {
		log.Errorf("Trends dao Update err %v", err)
	}
	return mod, err
}

func DeleteTrendsById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateTrendsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteTrendsByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateTrendsByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectTrendsById(id int64) (*po.PpmTreTrends, error) {
	po := &po.PpmTreTrends{}
	err := mysql.SelectById(consts.TableTrends, id, po)
	if err != nil {
		log.Errorf("Trends dao SelectById err %v", err)
	}
	return po, err
}

func SelectTrendsByIdAndOrg(id int64, orgId int64) (*po.PpmTreTrends, error) {
	po := &po.PpmTreTrends{}
	err := mysql.SelectOneByCond(consts.TableTrends, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Trends dao Select err %v", err)
	}
	return po, err
}

func SelectTrends(cond db.Cond) (*[]po.PpmTreTrends, error) {
	pos := &[]po.PpmTreTrends{}
	err := mysql.SelectAllByCond(consts.TableTrends, cond, pos)
	if err != nil {
		log.Errorf("Trends dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneTrends(cond db.Cond) (*po.PpmTreTrends, error) {
	po := &po.PpmTreTrends{}
	err := mysql.SelectOneByCond(consts.TableTrends, cond, po)
	if err != nil {
		log.Errorf("Trends dao Select err %v", err)
	}
	return po, err
}

func SelectTrendsByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmTreTrends, uint64, error) {
	pos := &[]po.PpmTreTrends{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableTrends, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Trends dao SelectPage err %v", err)
	}
	return pos, total, err
}
