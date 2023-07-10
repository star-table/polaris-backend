package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetIterationStatBoList(page uint, size uint, cond db.Cond) (*[]bo.IterationStatBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectIterationStatByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: consts.TcStatDate + " asc",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.IterationStatBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetIterationStatBo(id int64) (*bo.IterationStatBo, errs.SystemErrorInfo) {
	po, err := dao.SelectIterationStatById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.IterationStatBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateIterationStat(bo *bo.IterationStatBo) errs.SystemErrorInfo {
	po := &po.PpmStaIterationStat{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.InsertIterationStat(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateIterationStat(bo *bo.IterationStatBo) errs.SystemErrorInfo {
	po := &po.PpmStaIterationStat{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.UpdateIterationStat(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func DeleteIterationStat(bo *bo.IterationStatBo, operatorId int64) errs.SystemErrorInfo {
	_, err := dao.DeleteIterationStatById(bo.Id, operatorId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}
