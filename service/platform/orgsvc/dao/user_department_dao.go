package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertUserDepartment(po po.PpmOrgUserDepartment, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("UserDepartment dao Insert err %v", err)
	}
	return nil
}

func InsertUserDepartmentBatch(pos []po.PpmOrgUserDepartment, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableUserDepartment).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableUserDepartment).Batch(len(pos))
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

func UpdateUserDepartment(po po.PpmOrgUserDepartment, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("UserDepartment dao Update err %v", err)
	}
	return err
}

func UpdateUserDepartmentById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateUserDepartmentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateUserDepartmentByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateUserDepartmentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateUserDepartmentByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableUserDepartment, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableUserDepartment, cond, upd)
	}
	if err != nil {
		log.Errorf("UserDepartment dao Update err %v", err)
	}
	return mod, err
}

func DeleteUserDepartmentById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateUserDepartmentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteUserDepartmentByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateUserDepartmentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectUserDepartmentById(id int64) (*po.PpmOrgUserDepartment, error) {
	po := &po.PpmOrgUserDepartment{}
	err := mysql.SelectById(consts.TableUserDepartment, id, po)
	if err != nil {
		log.Errorf("UserDepartment dao SelectById err %v", err)
	}
	return po, err
}

func SelectUserDepartmentByIdAndOrg(id int64, orgId int64) (*po.PpmOrgUserDepartment, error) {
	po := &po.PpmOrgUserDepartment{}
	err := mysql.SelectOneByCond(consts.TableUserDepartment, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("UserDepartment dao Select err %v", err)
	}
	return po, err
}

func SelectUserDepartment(cond db.Cond) (*[]po.PpmOrgUserDepartment, error) {
	pos := &[]po.PpmOrgUserDepartment{}
	err := mysql.SelectAllByCond(consts.TableUserDepartment, cond, pos)
	if err != nil {
		log.Errorf("UserDepartment dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneUserDepartment(cond db.Cond) (*po.PpmOrgUserDepartment, error) {
	po := &po.PpmOrgUserDepartment{}
	err := mysql.SelectOneByCond(consts.TableUserDepartment, cond, po)
	if err != nil {
		log.Errorf("UserDepartment dao Select err %v", err)
	}
	return po, err
}

func SelectUserDepartmentByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmOrgUserDepartment, uint64, error) {
	pos := &[]po.PpmOrgUserDepartment{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserDepartment, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("UserDepartment dao SelectPage err %v", err)
	}
	return pos, total, err
}
