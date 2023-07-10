package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertDepartment(po po.PpmOrgDepartment, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Department dao Insert err %v", err)
	}
	return nil
}

func InsertDepartmentBatch(pos []po.PpmOrgDepartment, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableDepartment).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableDepartment).Batch(len(pos))
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

func UpdateDepartment(po po.PpmOrgDepartment, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Department dao Update err %v", err)
	}
	return err
}

func UpdateDepartmentById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateDepartmentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateDepartmentByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateDepartmentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateDepartmentByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableDepartment, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableDepartment, cond, upd)
	}
	if err != nil {
		log.Errorf("Department dao Update err %v", err)
	}
	return mod, err
}

func DeleteDepartmentById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateDepartmentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteDepartmentByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateDepartmentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectDepartmentById(id int64) (*po.PpmOrgDepartment, error) {
	po := &po.PpmOrgDepartment{}
	err := mysql.SelectById(consts.TableDepartment, id, po)
	if err != nil {
		log.Errorf("Department dao SelectById err %v", err)
	}
	return po, err
}

func SelectDepartmentByIdAndOrg(id int64, orgId int64) (*po.PpmOrgDepartment, error) {
	po := &po.PpmOrgDepartment{}
	err := mysql.SelectOneByCond(consts.TableDepartment, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Department dao Select err %v", err)
	}
	return po, err
}

func SelectDepartment(cond db.Cond) (*[]po.PpmOrgDepartment, error) {
	pos := &[]po.PpmOrgDepartment{}
	err := mysql.SelectAllByCond(consts.TableDepartment, cond, pos)
	if err != nil {
		log.Errorf("Department dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneDepartment(cond db.Cond) (*po.PpmOrgDepartment, error) {
	po := &po.PpmOrgDepartment{}
	err := mysql.SelectOneByCond(consts.TableDepartment, cond, po)
	if err != nil {
		log.Errorf("Department dao Select err %v", err)
	}
	return po, err
}

func SelectDepartmentByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmOrgDepartment, uint64, error) {
	pos := &[]po.PpmOrgDepartment{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableDepartment, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Department dao SelectPage err %v", err)
	}
	return pos, total, err
}
