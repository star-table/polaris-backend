package domain

import (
	"fmt"
	"testing"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"github.com/smartystreets/goconvey/convey"
	"upper.io/db.v3"
)

func TestTeamInit(t *testing.T) {

	config.LoadEnvConfig("F:\\polaris-backend-clone\\config", "application.common", "local")

	conn, _ := mysql.GetConnect()
	tx, _ := conn.NewTx(nil)

	TeamInit(2, tx)
	tx.Commit()

}

func TestGetCorpAuthInfo(t *testing.T) {
	convey.Convey("测试加载env2", t, func() {
		config.LoadEnvConfig("config", "application", "dev")
		info := &[]po.PpmOrgDepartmentOutInfo{}
		mysql.SelectAllByCond(consts.TableDepartmentOutInfo, db.Cond{}, info)
		fmt.Println(info)
	})
}
