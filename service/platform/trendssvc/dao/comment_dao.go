package dao

import (
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/trendssvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertComment(po po.PpmTreComment, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Comment dao Insert err " + strs.ObjectToString(err))
	}
	return nil
}

//func InsertCommentBatch(pos []po.PpmTreComment, tx ...sqlbuilder.Tx) error{
//	var err error = nil
//	if tx != nil && len(tx) > 0{
//		batch := tx[0].InsertInto(consts.TableComment).Batch(len(pos))
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
//		batch := conn.InsertInto(consts.TableComment).Batch(len(pos))
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
//		log.Errorf("Comment dao InsertBatch err %v", err)
//	}
//	return nil
//}

func InsertCommentBatch(pos []po.PpmTreComment, tx ...sqlbuilder.Tx) error {
	var err error = nil

	isTx := tx != nil && len(tx) > 0

	var batch *sqlbuilder.BatchInserter

	if !isTx {
		//没有事务
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}

		batch = conn.InsertInto(consts.TableComment).Batch(len(pos))
	}

	if batch == nil {
		batch = tx[0].InsertInto(consts.TableComment).Batch(len(pos))
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

func UpdateComment(po po.PpmTreComment, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Comment dao Update err %v", err)
	}
	return err
}

func UpdateCommentById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateCommentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateCommentByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateCommentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateCommentByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableComment, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableComment, cond, upd)
	}
	if err != nil {
		log.Errorf("Comment dao Update err %v", err)
	}
	return mod, err
}

func DeleteCommentById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateCommentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteCommentByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateCommentByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectCommentById(id int64) (*po.PpmTreComment, error) {
	po := &po.PpmTreComment{}
	err := mysql.SelectById(consts.TableComment, id, po)
	if err != nil {
		log.Errorf("Comment dao SelectById err %v", err)
	}
	return po, err
}

func SelectCommentByIdAndOrg(id int64, orgId int64) (*po.PpmTreComment, error) {
	po := &po.PpmTreComment{}
	err := mysql.SelectOneByCond(consts.TableComment, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Comment dao Select err %v", err)
	}
	return po, err
}

func SelectComment(cond db.Cond) (*[]po.PpmTreComment, error) {
	pos := &[]po.PpmTreComment{}
	err := mysql.SelectAllByCond(consts.TableComment, cond, pos)
	if err != nil {
		log.Errorf("Comment dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneComment(cond db.Cond) (*po.PpmTreComment, error) {
	po := &po.PpmTreComment{}
	err := mysql.SelectOneByCond(consts.TableComment, cond, po)
	if err != nil {
		log.Errorf("Comment dao Select err %v", err)
	}
	return po, err
}

func SelectCommentByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmTreComment, uint64, error) {
	pos := &[]po.PpmTreComment{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableComment, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Comment dao SelectPage err %v", err)
	}
	return pos, total, err
}
