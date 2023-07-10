package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertIndustry(po po.PpmCmmIndustry, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Industry dao Insert err %v", err)
	}
	return nil
}

func InsertIndustryBatch(pos []po.PpmCmmIndustry, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableIndustry).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableIndustry).Batch(len(pos))
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

func UpdateIndustry(po po.PpmCmmIndustry, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Industry dao Update err %v", err)
	}
	return err
}

func UpdateIndustryById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIndustryByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIndustryByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIndustryByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIndustryByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableIndustry, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableIndustry, cond, upd)
	}
	if err != nil {
		log.Errorf("Industry dao Update err %v", err)
	}
	return mod, err
}

func DeleteIndustryById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIndustryByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteIndustryByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIndustryByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectIndustryById(id int64) (*po.PpmCmmIndustry, error) {
	po := &po.PpmCmmIndustry{}
	err := mysql.SelectById(consts.TableIndustry, id, po)
	if err != nil {
		log.Errorf("Industry dao SelectById err %v", err)
	}
	return po, err
}

func SelectIndustryByIdAndOrg(id int64, orgId int64) (*po.PpmCmmIndustry, error) {
	po := &po.PpmCmmIndustry{}
	err := mysql.SelectOneByCond(consts.TableIndustry, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Industry dao Select err %v", err)
	}
	return po, err
}

func SelectIndustry(cond db.Cond) (*[]po.PpmCmmIndustry, error) {
	pos := &[]po.PpmCmmIndustry{}
	err := mysql.SelectAllByCond(consts.TableIndustry, cond, pos)
	if err != nil {
		log.Errorf("Industry dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneIndustry(cond db.Cond) (*po.PpmCmmIndustry, error) {
	po := &po.PpmCmmIndustry{}
	err := mysql.SelectOneByCond(consts.TableIndustry, cond, po)
	if err != nil {
		log.Errorf("Industry dao Select err %v", err)
	}
	return po, err
}

func SelectIndustryByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmCmmIndustry, uint64, error) {
	pos := &[]po.PpmCmmIndustry{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIndustry, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Industry dao SelectPage err %v", err)
	}
	return pos, total, err
}

func SelectIndustryBoAllList(cond db.Cond) (*[]po.PpmCmmIndustry, error) {
	pos := &[]po.PpmCmmIndustry{}

	err := mysql.SelectAllByCond(consts.TableProjectObjectType, db.Cond{
		consts.TcIsShow: consts.AppShowEnable,
	}, pos)

	return pos, err
}
