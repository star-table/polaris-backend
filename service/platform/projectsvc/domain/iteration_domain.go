package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetIterationBoList(page uint, size uint, cond db.Cond, orderBy *int) (*[]bo.IterationBo, int64, errs.SystemErrorInfo) {
	order := "sort desc"
	if orderBy != nil {
		switch *orderBy {
		case 1:
			order = "create_time asc"
		case 2:
			order = "create_time desc"
		case 3:
			order = "sort asc"
		case 4:
			order = "sort desc"
		default:
			order = "sort desc"
		}
	}
	pos, total, err := dao.SelectIterationByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: order,
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.IterationBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetIterationBo(id int64) (*bo.IterationBo, errs.SystemErrorInfo) {
	po, err := dao.SelectIterationById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.IterationBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func GetIterationBoByOrgId(id int64, orgId int64) (*bo.IterationBo, errs.SystemErrorInfo) {
	po, err := dao.SelectIterationByIdAndOrg(id, orgId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.IterationBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateIteration(bo *bo.IterationBo, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	po := &po.PpmPriIteration{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.InsertIteration(*po, tx)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateIteration(updateBo *bo.IterationUpdateBo) errs.SystemErrorInfo {
	_, err := dao.UpdateIterationById(updateBo.Id, updateBo.Upd)
	//TODO 处理动态

	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func DeleteIteration(bo *bo.IterationBo, operatorId int64) errs.SystemErrorInfo {
	err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除迭代
		_, err := dao.DeleteIterationById(bo.Id, operatorId, tx)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}

		//删除任务和迭代关联
		//_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcOrgId:       bo.OrgId,
		//	consts.TcIterationId: bo.Id,
		//	consts.TcIsDelete:    consts.AppIsNoDelete,
		//}, mysql.Upd{
		//	consts.TcIterationId: 0,
		//	consts.TcUpdator:     operatorId,
		//})
		//if err2 != nil {
		//	log.Error(err2)
		//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		//}
		return nil
	})

	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}
	return nil
}

func GetIterationsInfo(orgId int64, ids []int64) ([]bo.IterationBo, errs.SystemErrorInfo) {
	iterationPos := &[]po.PpmPriIteration{}
	err := mysql.SelectAllByCond(consts.TableIteration, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(ids),
	}, iterationPos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.IterationBo{}
	copyer.Copy(iterationPos, bos)

	return *bos, nil
}

//获取指定迭代前一个迭代sort值
func GetBeforeIterationSort(orgId, projectId int64, afterSort int64) (int64, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return 0, errs.MysqlOperateError
	}
	info := &po.PpmPriIteration{}
	err1 := conn.Select("id", "sort").From(consts.TableIteration).Where(db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
		consts.TcSort:      db.Lt(afterSort),
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}).OrderBy("sort desc").One(info)
	if err1 != nil {
		if err1 == db.ErrNoMoreRows {
			//如果没有比这个更大的
			res := (afterSort/65536 - 1) * 65536
			return res, nil
		}
		log.Error(err1)
		return 0, errs.MysqlOperateError
	}

	return info.Sort, nil
}

//获取指定迭代后一个迭代sort值
func GetAfterIterationSort(orgId, projectId int64, beforeSort int64) (int64, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return 0, errs.MysqlOperateError
	}
	info := &po.PpmPriIteration{}
	err1 := conn.Select("id", "sort").From(consts.TableIteration).Where(db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
		consts.TcSort:      db.Gt(beforeSort),
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}).OrderBy("sort asc").One(info)
	if err1 != nil {
		if err1 == db.ErrNoMoreRows {
			//如果没有比这个更大的
			res := (beforeSort/65536 + 1) * 65536
			return res, nil
		}
		log.Error(err1)
		return 0, errs.MysqlOperateError
	}

	return info.Sort, nil
}

func GetMinIterationSort(orgId, projectId int64) (int64, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return 0, errs.MysqlOperateError
	}
	info := &po.PpmPriIteration{}
	err1 := conn.Select("id", "sort").From(consts.TableIteration).Where(db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}).OrderBy("sort asc").One(info)
	if err1 != nil {
		if err1 == db.ErrNoMoreRows {
			return 0, nil
		}
		log.Error(err1)
		return 0, errs.MysqlOperateError
	}

	return info.Sort, nil
}
