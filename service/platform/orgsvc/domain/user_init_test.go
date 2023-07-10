package domain

import (
	"fmt"
	"testing"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/sdk/dingtalk"
)

//func TestUserInit(t *testing.T) {
//
//	config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\polaris-server\\configs", "application")
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
//	cache.Set(consts.CacheDingTalkSuiteTicket, "abc")
//	_, err := UserInitByOrg("manager5225", "ding8ac2bab2b708b3cc35c2f4657eb6378f", 0, tx)
//
//	if err != nil {
//		err2 := tx.Rollback()
//		if err2 != nil {
//			log.Info(strs.ObjectToString(err))
//		}
//	} else {
//		err2 := tx.Commit()
//		if err2 != nil {
//			log.Info(strs.ObjectToString(err))
//		}
//	}
//}

func TestGetUserList(t *testing.T) {
	config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\polaris-server\\configs", "application")

	c, _ := dingtalk.GetDingTalkClient("ding8ac2bab2b708b3cc35c2f4657eb6378f", "_")
	fmt.Println(c.GetDepMemberIds("1"))
}
