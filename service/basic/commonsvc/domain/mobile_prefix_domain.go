package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/dao"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
)

func GetMobilePrefixBoList(page uint, size uint, cond db.Cond) (*[]bo.MobilePrefixBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectMobilePrefixByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.MobilePrefixBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetMobilePrefixBo(id int64) (*bo.MobilePrefixBo, errs.SystemErrorInfo) {
	po, err := dao.SelectMobilePrefixById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.MobilePrefixBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateMobilePrefix(bo *bo.MobilePrefixBo) errs.SystemErrorInfo {
	po := &po.PpmCmmMobilePrefix{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.InsertMobilePrefix(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateMobilePrefix(bo *bo.MobilePrefixBo) errs.SystemErrorInfo {
	po := &po.PpmCmmMobilePrefix{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.UpdateMobilePrefix(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func DeleteMobilePrefix(bo *bo.MobilePrefixBo, operatorId int64) errs.SystemErrorInfo {
	_, err := dao.DeleteMobilePrefixById(bo.Id, operatorId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}
