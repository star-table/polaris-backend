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

func GetProjectDetailBoList(page uint, size uint, cond db.Cond) (*[]bo.ProjectDetailBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectProjectDetailByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.ProjectDetailBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetProjectDetailBo(id int64) (*bo.ProjectDetailBo, errs.SystemErrorInfo) {
	po, err := dao.SelectProjectDetailById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.ProjectDetailBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func GetProjectDetails(orgId int64, projectIds []int64) (*[]bo.ProjectDetailBo, errs.SystemErrorInfo) {
	po := &[]po.PpmProProjectDetail{}
	err := mysql.SelectAllByCond(consts.TableProjectDetail, db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcOrgId:     orgId,
		consts.TcProjectId: db.In(projectIds),
	}, po)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bo := &[]bo.ProjectDetailBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func GetProjectDetailByProjectIdBo(projectId, orgId int64) (*bo.ProjectDetailBo, errs.SystemErrorInfo) {
	po, err := dao.SelectOneProjectDetail(db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.ProjectDetailBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateProjectDetail(bo *bo.ProjectDetailBo) errs.SystemErrorInfo {
	po := &po.PpmProProjectDetail{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.InsertProjectDetail(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateProjectDetail(bo *bo.ProjectDetailBo) errs.SystemErrorInfo {
	po := &po.PpmProProjectDetail{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.UpdateProjectDetail(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func DeleteProjectDetail(bo *bo.ProjectDetailBo, operatorId int64) errs.SystemErrorInfo {
	_, err := dao.DeleteProjectDetailById(bo.Id, operatorId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

// CheckProIsEnablePrivacyMode 检查项目是否开启了隐私模式
func CheckProIsEnablePrivacyMode(fsChatSettingPtr *string) bool {
	return false
}
