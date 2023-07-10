package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"upper.io/db.v3"
)

//获取未完成的的迭代列表
func GetNotCompletedIterationBoList(orgId int64, projectId int64) ([]bo.IterationBo, errs.SystemErrorInfo) {
	statusIds, _ := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryProject, consts.StatusTypeComplete)
	bos := &[]bo.IterationBo{}

	pos, err1 := dao.SelectIteration(db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
		consts.TcStatus:    db.NotIn(statusIds),
		consts.TcIsDelete:  consts.AppIsNoDelete,
	})
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}

	err2 := copyer.Copy(pos, bos)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *bos, nil
}
