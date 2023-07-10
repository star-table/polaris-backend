package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertRecycleBin(po po.PpmPrsRecycleBin, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("RecycleBin dao Insert err %v", err)
	}
	return nil
}

func InsertRecycleBinBatch(pos []po.PpmPrsRecycleBin, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableRecycleBin).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableRecycleBin).Batch(len(pos))
	}
	go func() {
		defer batch.Done()
		for i := range pos {
			batch.Values(pos[i])
		}
	}()
	err = batch.Wait()
	if err != nil {
		log.Errorf("RecycleBin dao InsertBatch err %v", err)
		return err
	}
	return nil
}
