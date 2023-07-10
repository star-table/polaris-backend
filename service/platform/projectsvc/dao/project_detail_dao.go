package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertProjectDetail(po po.PpmProProjectDetail, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("ProjectDetail dao Insert err %v", err)
	}
	return nil
}

//func InsertProjectDetailBatch(pos []po.PpmProProjectDetail, tx ...sqlbuilder.Tx) error {
//	var err error = nil
//	if tx != nil && len(tx) > 0 {
//		batch := tx[0].InsertInto(consts.TableProjectDetail).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
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
//		batch := conn.InsertInto(consts.TableProjectDetail).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
//		go mysql.BatchDone(slice.ToSlice(pos), batch)
//
//		err = batch.Wait()
//		if err != nil {
//			return err
//		}
//	}
//	if err != nil {
//		log.Errorf("ProjectDetail dao InsertBatch err %v", err)
//	}
//	return nil
//}

func InsertProjectDetailBatch(pos []po.PpmProProjectDetail, tx ...sqlbuilder.Tx) error {
	var err error = nil

	isTx := tx != nil && len(tx) > 0

	var batch *sqlbuilder.BatchInserter

	if !isTx {
		//没有事务
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}

		batch = conn.InsertInto(consts.TableProjectDetail).Batch(len(pos))
	}

	if batch == nil {
		batch = tx[0].InsertInto(consts.TableProjectDetail).Batch(len(pos))
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

func UpdateProjectDetail(po po.PpmProProjectDetail, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("ProjectDetail dao Update err %v", err)
	}
	return err
}

func UpdateProjectDetailById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateProjectDetailByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateProjectDetailByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateProjectDetailByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateProjectDetailByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableProjectDetail, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableProjectDetail, cond, upd)
	}
	if err != nil {
		log.Errorf("ProjectDetail dao Update err %v", err)
	}
	return mod, err
}

func DeleteProjectDetailById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateProjectDetailByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteProjectDetailByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateProjectDetailByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectProjectDetailById(id int64) (*po.PpmProProjectDetail, error) {
	po := &po.PpmProProjectDetail{}
	err := mysql.SelectById(consts.TableProjectDetail, id, po)
	if err != nil {
		log.Errorf("ProjectDetail dao SelectById err %v", err)
	}
	return po, err
}

func SelectProjectDetailByIdAndOrg(id int64, orgId int64) (*po.PpmProProjectDetail, error) {
	po := &po.PpmProProjectDetail{}
	err := mysql.SelectOneByCond(consts.TableProjectDetail, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("ProjectDetail dao Select err %v", err)
	}
	return po, err
}

func SelectProjectDetail(cond db.Cond) (*[]po.PpmProProjectDetail, error) {
	pos := &[]po.PpmProProjectDetail{}
	err := mysql.SelectAllByCond(consts.TableProjectDetail, cond, pos)
	if err != nil {
		log.Errorf("ProjectDetail dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneProjectDetail(cond db.Cond) (*po.PpmProProjectDetail, error) {
	po := &po.PpmProProjectDetail{}
	err := mysql.SelectOneByCond(consts.TableProjectDetail, cond, po)
	if err != nil {
		log.Errorf("ProjectDetail dao Select cond:%v, err %v", cond, err)
	}
	return po, err
}

func SelectProjectDetailByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmProProjectDetail, uint64, error) {
	pos := &[]po.PpmProProjectDetail{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableProjectDetail, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("ProjectDetail dao SelectPage err %v", err)
	}
	return pos, total, err
}
