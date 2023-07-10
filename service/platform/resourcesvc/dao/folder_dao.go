package dao

import (
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertFolder(po po.PpmResFolder, tx ...sqlbuilder.Tx) error {
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

func FolderIdIsExist(ids []int64, projectId int64, orgId int64, tx ...sqlbuilder.Tx) (bool, errs.SystemErrorInfo) {
	//0为首页
	if len(ids) == 1 && ids[0] == 0 {
		return true, nil
	}
	total := uint64(len(ids))
	if ok, _ := slice.Contain(ids, int64(0)); ok {
		total = total - 1
	}
	if tx != nil && len(tx) > 0 {
		count, err := mysql.TransSelectCountByCond(tx[0], consts.TableFolder, db.Cond{
			consts.TcId:        db.In(ids),
			consts.TcProjectId: projectId,
			consts.TcOrgId:     orgId,
			consts.TcIsDelete:  consts.AppIsNoDelete,
		})
		if err != nil {
			return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		return count == total, nil
	} else {
		count, err1 := mysql.SelectCountByCond(consts.TableFolder, db.Cond{
			consts.TcId:        db.In(ids),
			consts.TcProjectId: projectId,
			consts.TcOrgId:     orgId,
			consts.TcIsDelete:  consts.AppIsNoDelete,
		})
		if err1 != nil {
			return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
		}
		return count == total, nil
	}
}

func SelectFolderByIds(ids []int64, tx ...sqlbuilder.Tx) (*[]po.PpmResFolder, error) {
	pos := &[]po.PpmResFolder{}
	if tx != nil && len(tx) > 0 {
		err := mysql.TransSelectAllByCond(tx[0], consts.TableFolder, db.Cond{
			consts.TcId:       db.In(ids),
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, pos)
		if err != nil {
			log.Error(err)
		}
		return pos, err
	} else {
		err := mysql.SelectAllByCond(consts.TableFolder, db.Cond{
			consts.TcId:       db.In(ids),
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, pos)
		if err != nil {
			log.Error(err)
		}
		return pos, err
	}
}

func UpdateFolder(po po.PpmResFolder, tx ...sqlbuilder.Tx) error {
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

func UpdateFolderByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	if tx != nil && len(tx) > 0 {
		_, err := mysql.TransUpdateSmartWithCond(tx[0], consts.TableFolder, cond, upd)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		return nil
	} else {
		_, err := mysql.UpdateSmartWithCond(consts.TableFolder, cond, upd)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		return nil
	}
}

func SelectFolderByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmResFolder, uint64, error) {
	pos := &[]po.PpmResFolder{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableFolder, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Folder dao SelectPage err %v", err)
	}
	return pos, total, err
}
