package domain

import (
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/dao"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
)

var log = logger.GetDefaultLogger()

func GetAreaBoList(page uint, size uint, cond db.Cond) (*[]bo.AreaBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectAreaByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.AreaBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetAreaBo(id int64) (*bo.AreaBo, errs.SystemErrorInfo) {
	po, err := dao.SelectAreaById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.AreaBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateArea(bo *bo.AreaBo) errs.SystemErrorInfo {
	po := &po.PpmCmmArea{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.InsertArea(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateArea(bo *bo.AreaBo) errs.SystemErrorInfo {
	po := &po.PpmCmmArea{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.UpdateArea(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func DeleteArea(bo *bo.AreaBo, operatorId int64) errs.SystemErrorInfo {
	_, err := dao.DeleteAreaById(bo.Id, operatorId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}
