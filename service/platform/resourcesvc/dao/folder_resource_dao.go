package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertMidTable(po po.PpmResFolderResource, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("MidTable dao Insert err %v", err)
	}
	return nil
}

func SelectMidTablePo(resourceId int64, folderId, orgId int64, tx ...sqlbuilder.Tx) (*po.PpmResFolderResource, errs.SystemErrorInfo) {
	midtablePo := &po.PpmResFolderResource{}
	if tx != nil && len(tx) > 0 {
		err := mysql.TransSelectOneByCond(tx[0], consts.TableFolderResource, db.Cond{
			consts.TcResourceId: resourceId,
			consts.TcOrgId:      orgId,
			consts.TcFolderId:   folderId,
			consts.TcIsDelete:   consts.AppIsNoDelete,
		}, midtablePo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	} else {
		err := mysql.SelectOneByCond(consts.TableFolderResource, db.Cond{
			consts.TcResourceId: resourceId,
			consts.TcOrgId:      orgId,
			consts.TcFolderId:   folderId,
			consts.TcIsDelete:   consts.AppIsNoDelete,
		}, midtablePo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}
	return midtablePo, nil
}

func UpdateMidTable(po po.PpmResFolderResource, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("MidTable dao Update err %v", err)
	}
	return err
}

func UpdateMidTableByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	if tx != nil && len(tx) > 0 {
		_, err := mysql.TransUpdateSmartWithCond(tx[0], consts.TableFolderResource, cond, upd)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		return nil
	} else {
		_, err := mysql.UpdateSmartWithCond(consts.TableFolderResource, cond, upd)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		return nil
	}
}

func RelationIsExist(resourceIds []int64, folderId, orgId int64) (bool, errs.SystemErrorInfo) {
	count, err1 := mysql.SelectCountByCond(consts.TableFolderResource, db.Cond{
		consts.TcFolderId:   folderId,
		consts.TcResourceId: db.In(resourceIds),
		consts.TcOrgId:      orgId,
		consts.TcIsDelete:   consts.AppIsNoDelete,
	})
	if err1 != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}
	if int(count) != len(resourceIds) {
		return false, nil
	}
	return true, nil
}

func ResourceFolderHasRelation(resourceIds []int64, folderId, orgId int64) (bool, errs.SystemErrorInfo) {
	count, err1 := mysql.SelectCountByCond(consts.TableFolderResource, db.Cond{
		consts.TcFolderId:   folderId,
		consts.TcResourceId: db.In(resourceIds),
		consts.TcOrgId:      orgId,
		consts.TcIsDelete:   consts.AppIsNoDelete,
	})
	if err1 != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}
	return count > 0, nil
}

func SelectMidTablePoByFolderId(folderId int64, orgId int64, tx ...sqlbuilder.Tx) (*[]po.PpmResFolderResource, errs.SystemErrorInfo) {
	folderResourcePo := &[]po.PpmResFolderResource{}
	if tx != nil && len(tx) > 0 {
		err1 := mysql.TransSelectAllByCond(tx[0], consts.TableFolderResource, db.Cond{
			consts.TcFolderId: folderId,
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, folderResourcePo)
		if err1 != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
		}
	} else {
		err1 := mysql.SelectAllByCond(consts.TableFolderResource, db.Cond{
			consts.TcFolderId: folderId,
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, folderResourcePo)
		if err1 != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
		}
	}
	return folderResourcePo, nil
}
