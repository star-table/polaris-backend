package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func CreateIssueView(po po.PpmPriIssueView, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("PpmPriIssueView dao Insert err %v", err)
	}
	return nil
}

func SelectIssueViewById(id int64) (*po.PpmPriIssueView, error) {
	po := &po.PpmPriIssueView{}
	err := mysql.SelectById(consts.TableIssueView, id, po)
	if err != nil {
		log.Errorf("PpmPriIssueView dao SelectById err %v", err)
	}
	return po, err
}

func SelectIssueViewsByCond(page, size int, cond db.Cond, order interface{}) (int64, []po.PpmPriIssueView, error) {
	poList := make([]po.PpmPriIssueView, 0)
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIssueView, cond, nil, page, size, order, &poList)
	if err == db.ErrNoMoreRows {
		return int64(0), poList, nil
	}
	return int64(total), poList, err
}

func SelectIssueViewsByCondCount(cond db.Cond) (uint64, error) {
	total, err := mysql.SelectCountByCond(consts.TableIssueView, cond)
	if err != nil {
		log.Errorf("SelectIssueViewsByCondCount dao Select err %v", err)
	}

	return total, err
}

func UpdateIssueView(po po.PpmPriIssueView, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("PpmPriIssueView dao Update err %v", err)
	}
	return err
}

func UpdateIssueViewById(orgId, id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIssueViewByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcDelFlag:  consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIssueViewByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableIssueView, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableIssueView, cond, upd)
	}
	if err != nil {
		log.Errorf("UpdateIssueViewByCond dao Update err %v", err)
	}
	return mod, err
}
