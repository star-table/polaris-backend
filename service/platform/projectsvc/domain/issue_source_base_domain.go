package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func SourceExist(orgId, sourceId int64) bool {
	isExist, err := mysql.IsExistByCond(consts.TableIssueSource, db.Cond{
		consts.TcId:       sourceId,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
	})
	if err != nil {
		return false
	}

	return isExist
}

func GetIssueSourceInfo(orgId int64, sourceIds []int64) ([]bo.IssueSourceBo, errs.SystemErrorInfo) {
	info := &[]po.PpmPrsIssueSource{}
	err := mysql.SelectAllByCond(consts.TableIssueSource, db.Cond{
		consts.TcId:       db.In(sourceIds),
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
	}, info)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	infoBo := &[]bo.IssueSourceBo{}
	copyErr := copyer.Copy(info, infoBo)
	if copyErr != nil {
		return nil, errs.ObjectCopyError
	}

	return *infoBo, nil
}
