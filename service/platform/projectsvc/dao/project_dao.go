package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func SelectOneProject(cond db.Cond) (*po.PpmProProject, error) {
	po := &po.PpmProProject{}
	err := mysql.SelectOneByCond(consts.TableProject, cond, po)
	if err != nil {
		log.Errorf("Project dao Select err %v", err)
	}
	return po, err
}

func InsertProject(po po.PpmProProject, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Project dao Insert err %v", err)
	}
	return nil
}

func InsertProjectBatch(pos []po.PpmProProject, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableProject).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableProject).Batch(len(pos))
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

func UpdateProject(po po.PpmProProject, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Project dao Update err %v", err)
	}
	return err
}

func UpdateProjectById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateProjectByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateProjectByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateProjectByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateProjectByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableProject, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableProject, cond, upd)
	}
	if err != nil {
		log.Errorf("Project dao Update err %v", err)
	}
	return mod, err
}

func DeleteProjectById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateProjectByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteProjectByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateProjectByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectProjectById(id int64) (*po.PpmProProject, error) {
	po := &po.PpmProProject{}
	err := mysql.SelectById(consts.TableProject, id, po)
	if err != nil {
		log.Errorf("Project dao SelectById err %v", err)
	}
	return po, err
}

func SelectProjectByIdAndOrg(id int64, orgId int64) (*po.PpmProProject, error) {
	po := &po.PpmProProject{}
	err := mysql.SelectOneByCond(consts.TableProject, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Project dao Select err %v", err)
	}
	return po, err
}

func SelectProject(cond db.Cond) (*[]po.PpmProProject, error) {
	pos := &[]po.PpmProProject{}
	err := mysql.SelectAllByCond(consts.TableProject, cond, pos)
	if err != nil {
		log.Errorf("Project dao SelectList err %v", err)
	}
	return pos, err
}

func SelectProjectByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmProProject, uint64, error) {
	pos := &[]po.PpmProProject{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableProject, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Project dao SelectPage err %v", err)
	}
	return pos, total, err
}
