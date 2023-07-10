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
)

func GetIssueSourceBoList(page uint, size uint, cond db.Cond) (*[]bo.IssueSourceBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectIssueSourceByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.IssueSourceBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetIssueSourceBo(id int64) (*bo.IssueSourceBo, errs.SystemErrorInfo) {
	po, err := dao.SelectIssueSourceById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.IssueSourceBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateIssueSource(bo *bo.IssueSourceBo) errs.SystemErrorInfo {
	po := &po.PpmPrsIssueSource{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.InsertIssueSource(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateIssueSource(bo *bo.IssueSourceBo) errs.SystemErrorInfo {
	po := &po.PpmPrsIssueSource{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.UpdateIssueSource(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func DeleteIssueSource(bo *bo.IssueSourceBo, operatorId int64) errs.SystemErrorInfo {
	_, err := dao.DeleteIssueSourceById(bo.Id, operatorId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func GetNewIssueSourceId(oldId int64, projectObjectTypeId int64, orgId int64) int64 {
	info, err := GetIssueSourceBo(oldId)
	if err != nil {
		log.Error(err)
		return 0
	}

	new := &po.PpmPrsIssueSource{}
	err1 := mysql.SelectOneByCond(consts.TableIssueSource, db.Cond{
		consts.TcIsDelete:            consts.AppIsNoDelete,
		consts.TcOrgId:               orgId,
		consts.TcProjectObjectTypeId: projectObjectTypeId,
		consts.TcName:                info.Name,
	}, new)
	if err1 != nil {
		log.Error(err1)
		return 0
	}

	return new.Id
}

func GetIssueSourceBos(orgId, projectId int64, projectObjectTypeId int64) ([]bo.IssueSourceBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmPrsIssueSource{}
	cond := db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
	}
	if projectObjectTypeId > 0 {
		cond[consts.TcProjectObjectTypeId] = projectObjectTypeId
	}
	err := mysql.SelectAllByCond(consts.TableIssueSource, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.IssueSourceBo{}
	_ = copyer.Copy(pos, bos)

	return *bos, nil
}
