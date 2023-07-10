package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertResourceRelation(po po.PpmResResourceRelation, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("ResourceRelation dao Insert err %v", err)
	}
	return nil
}

func UpdateResourceRelationByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	if tx != nil && len(tx) > 0 {
		_, err := mysql.TransUpdateSmartWithCond(tx[0], consts.TableResourceRelation, cond, upd)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		return nil
	} else {
		_, err := mysql.UpdateSmartWithCond(consts.TableResourceRelation, cond, upd)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		return nil
	}
}

func SelectResourceRelationByCond(cond db.Cond, tx sqlbuilder.Tx) ([]*po.PpmResResourceRelation, error) {
	pos := make([]*po.PpmResResourceRelation, 0, 5)
	err := mysql.TransSelectAllByCond(tx, consts.TableResourceRelation, cond, &pos)
	if err != nil {
		log.Errorf("SelectResourceRelationByCond dao SelectList err %v", err)
	}
	return pos, err
}
