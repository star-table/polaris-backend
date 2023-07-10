package domain

//func GetPriorityBoList(page uint, size uint, cond db.Cond) (*[]bo.PriorityBo, int64, errs.SystemErrorInfo) {
//	pos, total, err := dao.SelectPriorityByPage(cond, bo.PageBo{
//		Page:  int(page),
//		Size:  int(size),
//		Order: "sort asc",
//	})
//	if err != nil {
//		log.Error(err)
//		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	bos := &[]bo.PriorityBo{}
//
//	copyErr := copyer.Copy(pos, bos)
//	if copyErr != nil {
//		log.Error(copyErr)
//		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
//	}
//	return bos, int64(total), nil
//}
//
//func GetPriorityBo(id int64) (*bo.PriorityBo, errs.SystemErrorInfo) {
//	po, err := dao.SelectPriorityById(id)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
//	}
//	bo := &bo.PriorityBo{}
//	err1 := copyer.Copy(po, bo)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
//	}
//	return bo, nil
//}
//
//func CreatePriority(bo *bo.PriorityBo) errs.SystemErrorInfo {
//	po := &po.PpmPrsPriority{}
//	copyErr := copyer.Copy(bo, po)
//	if copyErr != nil {
//		log.Error(copyErr)
//		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
//	}
//
//	err := dao.InsertPriority(*po)
//	if err != nil {
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//
//	return nil
//}
//
//func UpdatePriority(bo *bo.PriorityBo) errs.SystemErrorInfo {
//	po := &po.PpmPrsPriority{}
//	copyErr := copyer.Copy(bo, po)
//	if copyErr != nil {
//		log.Error(copyErr)
//		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
//	}
//
//	err := dao.UpdatePriority(*po)
//	if err != nil {
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	return nil
//}
//
//func DeletePriority(bo *bo.PriorityBo, operatorId int64) errs.SystemErrorInfo {
//	_, err := dao.DeletePriorityById(bo.Id, operatorId)
//	if err != nil {
//		log.Error(err)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	return nil
//}
//
//func VerifyPriority(orgId int64, typ int, beVerifyId int64) (bool, errs.SystemErrorInfo) {
//	list, err := GetPriorityListByType(orgId, typ)
//	if err != nil {
//		log.Error(err)
//		return false, err
//	}
//	for _, priority := range *list {
//		if priority.Id == beVerifyId {
//			return true, nil
//		}
//	}
//	return false, nil
//}
//
//func InitPriority(orgId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
//	count, err := mysql.SelectCountByCond((&po.PpmPrsPriority{}).TableName(), db.Cond{
//		consts.TcIsDelete: consts.AppIsNoDelete,
//		consts.TcOrgId:    orgId,
//	})
//	if err != nil {
//		log.Error(err)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	if count != 0 {
//		log.Infof("优先级已初始化")
//		return nil
//	}
//
//	priorityInsert := []interface{}{}
//	for _, v := range po.Priority {
//		id, err := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsPriority{}).TableName())
//		if err != nil {
//			log.Error(err)
//			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//		}
//		temp := v
//		temp.OrgId = orgId
//		temp.Id = id
//		priorityInsert = append(priorityInsert, temp)
//	}
//	err = mysql.TransBatchInsert(tx, &po.PpmPrsPriority{}, priorityInsert)
//	if err != nil {
//		log.Error(err)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//
//	return nil
//}
