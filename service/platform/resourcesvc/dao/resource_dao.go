package dao

import (
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

func InsertResource(po po.PpmResResource, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Resource dao Insert err %v", err)
	}
	return nil
}
func ResourceIdIsExist(resourceIds []int64, orgId int64, projectId int64) (bool, errs.SystemErrorInfo) {
	isExist, err1 := mysql.IsExistByCond(consts.TableResource, db.Cond{
		consts.TcId: db.In(resourceIds),
		//consts.TcProjectId: projectId,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})

	if err1 != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}
	return isExist, nil
}

//func InsertResourceBatch(pos []po.PpmResResource, tx ...sqlbuilder.Tx) error{
//	var err error = nil
//	if tx != nil && len(tx) > 0{
//		batch := tx[0].InsertInto(consts.TableResource).Batch(len(pos))
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
//		batch := conn.InsertInto(consts.TableResource).Batch(len(pos))
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
//		log.Errorf("Resource dao InsertBatch err %v", err)
//	}
//	return nil
//}

func InsertResourceBatch(pos []po.PpmResResource, tx ...sqlbuilder.Tx) error {
	var err error = nil

	isTx := tx != nil && len(tx) > 0

	var batch *sqlbuilder.BatchInserter

	if !isTx {
		//没有事务
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}

		batch = conn.InsertInto(consts.TableResource).Batch(len(pos))
	}

	if batch == nil {
		batch = tx[0].InsertInto(consts.TableResource).Batch(len(pos))
	}

	go func() {
		defer batch.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()
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

func UpdateResource(po po.PpmResResource, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Resource dao Update err %v", err)
	}
	return err
}

func UpdateResourceById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateResourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateResourceByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateResourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateResourceByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableResource, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableResource, cond, upd)
	}
	if err != nil {
		log.Errorf("Resource dao Update err %v", err)
	}
	return mod, err
}

func DeleteResourceById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateResourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteResourceByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateResourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectResourceById(id int64, tx ...sqlbuilder.Tx) (*po.PpmResResource, error) {
	po := &po.PpmResResource{}
	if tx != nil && len(tx) > 0 {
		err := mysql.TransSelectById(tx[0], consts.TableResource, id, po)
		if err != nil {
			log.Errorf("Resource dao SelectById err %v", err)
		}
		return po, err
	} else {
		err := mysql.SelectById(consts.TableResource, id, po)
		if err != nil {
			log.Errorf("Resource dao SelectById err %v", err)
		}
		return po, err
	}
}

func SelectResourceByIdAndOrg(id int64, orgId int64) (*po.PpmResResource, error) {
	po := &po.PpmResResource{}
	err := mysql.SelectOneByCond(consts.TableResource, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Resource dao Select err %v", err)
	}
	return po, err
}

func SelectResource(cond db.Cond) ([]*po.PpmResResource, error) {
	pos := []*po.PpmResResource{}
	err := mysql.SelectAllByCond(consts.TableResource, cond, &pos)
	if err != nil {
		log.Errorf("Resource dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneResource(cond db.Cond) (*po.PpmResResource, error) {
	po := &po.PpmResResource{}
	err := mysql.SelectOneByCond(consts.TableResource, cond, po)
	if err != nil {
		log.Errorf("Resource dao Select err %v", err)
	}
	return po, err
}

func SelectResourceByPage(cond db.Cond, union *db.Union, pageBo bo.PageBo) (*[]po.PpmResResource, uint64, error) {
	pos := &[]po.PpmResResource{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableResource, cond, union, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Resource dao SelectPage err %v", err)
	}
	return pos, total, err
}
