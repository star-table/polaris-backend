package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/dao"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
)

func GetCitiesBoList(page uint, size uint, cond db.Cond) (*[]bo.CitiesBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectCitiesByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.CitiesBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetCityBoListByCond(cond db.Cond) (*[]bo.CitiesBo, errs.SystemErrorInfo) {
	pos, err := dao.SelectCities(cond)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.CitiesBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, nil
}

func GetCitiesBo(id int64) (*bo.CitiesBo, errs.SystemErrorInfo) {
	po, err := dao.SelectCitiesById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.CitiesBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateCities(bo *bo.CitiesBo) errs.SystemErrorInfo {
	po := &po.PpmCmmCities{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.InsertCities(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateCities(bo *bo.CitiesBo) errs.SystemErrorInfo {
	po := &po.PpmCmmCities{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.UpdateCities(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func DeleteCities(bo *bo.CitiesBo, operatorId int64) errs.SystemErrorInfo {
	_, err := dao.DeleteCitiesById(bo.Id, operatorId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}
