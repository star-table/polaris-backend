package service

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

func PAReport(reportMsg orgvo.PAReportMsg) errs.SystemErrorInfo {
	pa := po.PpmOrgPrivatizationAuthority{}
	mysqlErr := mysql.SelectOneByCond(consts.TablePrivatizationAuthority, db.Cond{
		consts.TcSecret: reportMsg.Secret,
	}, &pa)
	if mysqlErr != nil {
		if mysqlErr == db.ErrNoMoreRows {
			// todo 修改pa状态
			return errs.PANotExist
		}
		log.Error(mysqlErr)
		return errs.MysqlOperateError
	}

	prId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TablePrivatizationRecords)
	if err != nil {
		log.Error(err)
		return err
	}
	pr := po.PpmOrgPrivatizationRecords{
		Id:              prId,
		AuthorityId:     pa.Id,
		OrgUserCount:    reportMsg.UserCount,
		OrgDeptCount:    reportMsg.DeptCount,
		OrgProjectCount: reportMsg.ProjectCount,
		OrgIssueCount:   reportMsg.IssueCount,
		OutSideIp:       reportMsg.OutSideIp,
		IsDelete:        consts.AppIsNoDelete,
	}
	_ = mysql.Insert(&pr)
	return nil
}
