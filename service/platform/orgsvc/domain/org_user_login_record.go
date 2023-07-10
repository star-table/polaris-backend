package domain

import (
	"time"

	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
)

//添加登录记录
func AddLoginRecord(orgId, userId int64, sourceChannel string) errs.SystemErrorInfo {
	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrgUserLoginRecord)
	if err != nil {
		log.Error(err)
		return err
	}
	userLoginRecord := po.PpmOrgUserLoginRecord{
		Id:            id,
		OrgId:         orgId,
		UserId:        userId,
		LoginTime:     time.Now(),
		SourceChannel: sourceChannel,
		Creator:       userId,
		Updator:       userId,
	}

	insertErr := mysql.Insert(&userLoginRecord)
	if insertErr != nil {
		log.Error(insertErr)
		return errs.MysqlOperateError
	}

	return nil
}
