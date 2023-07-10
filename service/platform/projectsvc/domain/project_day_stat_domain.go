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

func GetProjectDayStatBoList(page uint, size uint, cond db.Cond) (*[]bo.ProjectDayStatBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectProjectDayStatByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.ProjectDayStatBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetProjectDayStatBo(id int64) (*bo.ProjectDayStatBo, errs.SystemErrorInfo) {
	po, err := dao.SelectProjectDayStatById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.ProjectDayStatBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateProjectDayStat(bo *bo.ProjectDayStatBo) errs.SystemErrorInfo {
	po := &po.PpmStaProjectDayStat{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.InsertProjectDayStat(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateProjectDayStat(bo *bo.ProjectDayStatBo) errs.SystemErrorInfo {
	po := &po.PpmStaProjectDayStat{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.UpdateProjectDayStat(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func DeleteProjectDayStat(bo *bo.ProjectDayStatBo, operatorId int64) errs.SystemErrorInfo {
	_, err := dao.DeleteProjectDayStatById(bo.Id, operatorId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func GetStatDailyAppList(page uint, size uint, cond db.Cond) ([]*po.LcStatDailyApp, int64, errs.SystemErrorInfo) {
	statAppPos := []*po.LcStatDailyApp{}
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	err = conn.Select(db.Raw("id, org_id, project_id, `date`, sum(issue_count) as issue_count, " +
		"sum(issue_not_start_count) as issue_not_start_count, sum(issue_running_count) as issue_running_count," +
		"sum(issue_overdue_count) as issue_overdue_count, sum(issue_complete_count) as issue_complete_count, creator, " +
		"create_time , updator, update_time")).
		From(consts.TableStatDailyApp).
		Where(cond).GroupBy(consts.TcDate).All(&statAppPos)
	//total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableStatDailyApp, cond, nil, int(page), int(size), "", &statAppPos)
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return statAppPos, int64(len(statAppPos)), nil
}
