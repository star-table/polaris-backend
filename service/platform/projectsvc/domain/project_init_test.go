package domain

import (
	"testing"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
)

func ConfigInit() {
	config.LoadConfig("F:\\polaris-backend\\polaris-server\\configs", "application")

	cache.Set(consts.CacheDingTalkSuiteTicket, "abc")
}

//func TestProjectProjectTypeInit(t *testing.T) {
//	ConfigInit()
//
//	conn, _ := mysql.GetConnect()
//
//	tx, _ := conn.NewTx(nil)
//	defer func() {
//		if tx != nil {
//			if err := tx.Close(); err != nil {
//				log.Info(strs.ObjectToString(err))
//			}
//		}
//	}()
//
//	contextMap := make(map[string]interface{})
//
//	err := ProjectObjectTypeInit(1000, &contextMap, tx)
//	if err != nil {
//		log.Info("初始化project_object_type失败:" + strs.ObjectToString(err))
//		err1 := tx.Rollback()
//		if err1 != nil {
//			log.Info("Rollback error:" + strs.ObjectToString(err1))
//		}
//	} else {
//		err2 := tx.Commit()
//		if err2 != nil {
//			log.Info("Commit error:" + strs.ObjectToString(err2))
//		}
//	}
//}

func TestPriorityInit(t *testing.T) {
	ConfigInit()

	conn, _ := mysql.GetConnect()

	tx, _ := conn.NewTx(nil)
	defer func() {
		if tx != nil {
			if err := tx.Close(); err != nil {
				log.Info(strs.ObjectToString(err))
			}
		}
	}()

	err := PriorityInit(1001, tx)
	if err != nil {
		log.Info("初始化priority失败" + strs.ObjectToString(err))
		err1 := tx.Rollback()
		if err1 != nil {
			log.Info("Rollback error" + strs.ObjectToString(err1))
		}
	} else {
		err2 := tx.Commit()
		if err2 != nil {
			log.Info("Commit error" + strs.ObjectToString(err2))
		}
	}
}

//func TestProjectInit(t *testing.T) {
//	ConfigInit()
//
//	conn, _ := mysql.GetConnect()
//
//	tx, _ := conn.NewTx(nil)
//
//	contextMap := make(map[string]interface{})
//
//	err := ProjectInit(6, &contextMap, tx)
//	if err != nil {
//		log.Info("初始化project失败", err)
//		err1 := tx.Rollback()
//		if err1 != nil {
//			log.Info("Rollback error", err1)
//		}
//	} else {
//		err2 := tx.Commit()
//		if err2 != nil {
//			log.Info("Commit error", err2)
//		}
//	}
//}
