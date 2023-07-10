package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertProjectDayStat(po po.PpmStaProjectDayStat, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("ProjectDayStat dao Insert err %v", err)
	}
	return nil
}

func InsertProjectDayStatBatch(pos []po.PpmStaProjectDayStat, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableProjectDayStat).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableProjectDayStat).Batch(len(pos))
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

func UpdateProjectDayStat(po po.PpmStaProjectDayStat, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("ProjectDayStat dao Update err %v", err)
	}
	return err
}

func UpdateProjectDayStatById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateProjectDayStatByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateProjectDayStatByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateProjectDayStatByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateProjectDayStatByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableProjectDayStat, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableProjectDayStat, cond, upd)
	}
	if err != nil {
		log.Errorf("ProjectDayStat dao Update err %v", err)
	}
	return mod, err
}

func DeleteProjectDayStatById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateProjectDayStatByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteProjectDayStatByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateProjectDayStatByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectProjectDayStatById(id int64) (*po.PpmStaProjectDayStat, error) {
	po := &po.PpmStaProjectDayStat{}
	err := mysql.SelectById(consts.TableProjectDayStat, id, po)
	if err != nil {
		log.Errorf("ProjectDayStat dao SelectById err %v", err)
	}
	return po, err
}

func SelectProjectDayStatByIdAndOrg(id int64, orgId int64) (*po.PpmStaProjectDayStat, error) {
	po := &po.PpmStaProjectDayStat{}
	err := mysql.SelectOneByCond(consts.TableProjectDayStat, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("ProjectDayStat dao Select err %v", err)
	}
	return po, err
}

func SelectProjectDayStat(cond db.Cond) (*[]po.PpmStaProjectDayStat, error) {
	pos := &[]po.PpmStaProjectDayStat{}
	err := mysql.SelectAllByCond(consts.TableProjectDayStat, cond, pos)
	if err != nil {
		log.Errorf("ProjectDayStat dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneProjectDayStat(cond db.Cond) (*po.PpmStaProjectDayStat, error) {
	po := &po.PpmStaProjectDayStat{}
	err := mysql.SelectOneByCond(consts.TableProjectDayStat, cond, po)
	if err != nil {
		log.Errorf("ProjectDayStat dao Select err %v", err)
	}
	return po, err
}

func SelectProjectDayStatByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmStaProjectDayStat, uint64, error) {
	pos := &[]po.PpmStaProjectDayStat{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableProjectDayStat, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("ProjectDayStat dao SelectPage err %v", err)
	}
	return pos, total, err
}
