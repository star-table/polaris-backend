package dao

//func InsertPriority(po po2.PpmPrsPriority, tx ...sqlbuilder.Tx) error {
//	var err error = nil
//	if tx != nil && len(tx) > 0 {
//		err = mysql.TransInsert(tx[0], &po)
//	} else {
//		err = mysql.Insert(&po)
//	}
//	if err != nil {
//		log.Errorf("Priority dao Insert err %v", err)
//	}
//	return nil
//}

//func InsertPriorityBatch(pos []po.PpmPrsPriority, tx ...sqlbuilder.Tx) error {
//	var err error = nil
//	if tx != nil && len(tx) > 0 {
//		batch := tx[0].InsertInto(consts.TablePriority).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
//		//go priorityBatchDone(pos,batch)
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
//		batch := conn.InsertInto(consts.TablePriority).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
//		//go priorityBatchDone(pos,batch)
//		go mysql.BatchDone(slice.ToSlice(pos), batch)
//
//		err = batch.Wait()
//		if err != nil {
//			return err
//		}
//	}
//	if err != nil {
//		log.Errorf("Priority dao InsertBatch err %v", err)
//	}
//	return nil
//}

//func InsertPriorityBatch(pos []po2.PpmPrsPriority, tx ...sqlbuilder.Tx) error {
//	var err error = nil
//
//	isTx := tx != nil && len(tx) > 0
//
//	var batch *sqlbuilder.BatchInserter
//
//	if !isTx {
//		//没有事务
//		conn, err := mysql.GetConnect()
//		if err != nil {
//			return err
//		}
//
//		batch = conn.InsertInto(consts.TablePriority).Batch(len(pos))
//	}
//
//	if batch == nil {
//		batch = tx[0].InsertInto(consts.TablePriority).Batch(len(pos))
//	}
//
//	go func() {
//		defer batch.Done()
//		for i := range pos {
//			batch.Values(pos[i])
//		}
//	}()
//
//	err = batch.Wait()
//	if err != nil {
//		log.Errorf("Iteration dao InsertBatch err %v", err)
//		return err
//	}
//	return nil
//}
//
//func UpdatePriority(po po2.PpmPrsPriority, tx ...sqlbuilder.Tx) error {
//	var err error = nil
//	if tx != nil && len(tx) > 0 {
//		err = mysql.TransUpdate(tx[0], &po)
//	} else {
//		err = mysql.Update(&po)
//	}
//	if err != nil {
//		log.Errorf("Priority dao Update err %v", err)
//	}
//	return err
//}
//
//func UpdatePriorityById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
//	return UpdatePriorityByCond(db.Cond{
//		consts.TcId:       id,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}, upd, tx...)
//}
//
//func UpdatePriorityByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
//	return UpdatePriorityByCond(db.Cond{
//		consts.TcId:       id,
//		consts.TcOrgId:    orgId,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}, upd, tx...)
//}
//
//func UpdatePriorityByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
//	var mod int64 = 0
//	var err error = nil
//	if tx != nil && len(tx) > 0 {
//		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TablePriority, cond, upd)
//	} else {
//		mod, err = mysql.UpdateSmartWithCond(consts.TablePriority, cond, upd)
//	}
//	if err != nil {
//		log.Errorf("Priority dao Update err %v", err)
//	}
//	return mod, err
//}
//
//func DeletePriorityById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
//	upd := mysql.Upd{
//		consts.TcIsDelete: consts.AppIsDeleted,
//	}
//	if operatorId > 0 {
//		upd[consts.TcUpdator] = operatorId
//	}
//	return UpdatePriorityByCond(db.Cond{
//		consts.TcId:       id,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}, upd, tx...)
//}
//
//func DeletePriorityByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
//	upd := mysql.Upd{
//		consts.TcIsDelete: consts.AppIsDeleted,
//	}
//	if operatorId > 0 {
//		upd[consts.TcUpdator] = operatorId
//	}
//	return UpdatePriorityByCond(db.Cond{
//		consts.TcId:       id,
//		consts.TcOrgId:    orgId,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}, upd, tx...)
//}
//
//func SelectPriorityById(id int64) (*po2.PpmPrsPriority, error) {
//	po := &po2.PpmPrsPriority{}
//	err := mysql.SelectById(consts.TablePriority, id, po)
//	if err != nil {
//		log.Errorf("Priority dao SelectById err %v", err)
//	}
//	return po, err
//}
//
//func SelectPriorityByIdAndOrg(id int64, orgId int64) (*po2.PpmPrsPriority, error) {
//	po := &po2.PpmPrsPriority{}
//	err := mysql.SelectOneByCond(consts.TablePriority, db.Cond{
//		consts.TcId:       id,
//		consts.TcOrgId:    orgId,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}, po)
//	if err != nil {
//		log.Errorf("Priority dao Select err %v", err)
//	}
//	return po, err
//}
//
//func SelectPriority(cond db.Cond) (*[]po2.PpmPrsPriority, error) {
//	pos := &[]po2.PpmPrsPriority{}
//	err := mysql.SelectAllByCond(consts.TablePriority, cond, pos)
//	if err != nil {
//		log.Errorf("Priority dao SelectList err %v", err)
//	}
//	return pos, err
//}
//
//func SelectOnePriority(cond db.Cond) (*po2.PpmPrsPriority, error) {
//	po := &po2.PpmPrsPriority{}
//	err := mysql.SelectOneByCond(consts.TablePriority, cond, po)
//	if err != nil {
//		log.Errorf("Priority dao Select err %v", err)
//	}
//	return po, err
//}
//
//func SelectPriorityByPage(cond db.Cond, pageBo bo.PageBo) (*[]po2.PpmPrsPriority, uint64, error) {
//	pos := &[]po2.PpmPrsPriority{}
//	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TablePriority, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
//	if err != nil {
//		log.Errorf("Priority dao SelectPage err %v", err)
//	}
//	return pos, total, err
//}
