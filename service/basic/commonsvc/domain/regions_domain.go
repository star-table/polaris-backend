package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/dao"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/po"
	"upper.io/db.v3"
)

func GetRegionsBoList(page uint, size uint, cond db.Cond) (*[]bo.RegionsBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectRegionsByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.RegionsBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetRegionBoListByCond(cond db.Cond) (*[]bo.RegionsBo, errs.SystemErrorInfo) {
	pos, err := dao.SelectRegions(cond)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.RegionsBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, nil
}

func GetRegionsBo(id int64) (*bo.RegionsBo, errs.SystemErrorInfo) {
	po, err := dao.SelectRegionsById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.RegionsBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateRegions(bo *bo.RegionsBo) errs.SystemErrorInfo {
	po := &po.PpmCmmRegions{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.InsertRegions(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateRegions(bo *bo.RegionsBo) errs.SystemErrorInfo {
	po := &po.PpmCmmRegions{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.UpdateRegions(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func DeleteRegions(bo *bo.RegionsBo, operatorId int64) errs.SystemErrorInfo {
	_, err := dao.DeleteRegionsById(bo.Id, operatorId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}
