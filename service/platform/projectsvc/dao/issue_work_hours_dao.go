package dao

import (
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3/lib/sqlbuilder"
)

// 插入一条数据
func InsertIssueWorkHours(po po.PpmPriIssueWorkHours, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("issue_work_hours dao Insert err %v", err)
		return err
	}
	return nil
}

// 插入多条
func InsertMultiIssueWorkHours(list []*po.PpmPriIssueWorkHours, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransBatchInsert(tx[0], &po.PpmPriIssueWorkHours{}, slice.ToSlice(list))
	} else {
		err = mysql.BatchInsert(&po.PpmPriIssueWorkHours{}, slice.ToSlice(list))
	}
	if err != nil {
		log.Errorf("issue_work_hours dao Insert err %v", err)
		return err
	}
	return nil
}
