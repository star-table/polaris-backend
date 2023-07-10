package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/library/db/mysql"
	domain "github.com/star-table/polaris-backend/service/platform/orgsvc/domain/roledomain"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestRoleInit(t *testing.T) {
	config.LoadConfig("F:\\polaris-backend\\polaris-server\\configs", "application")

	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("权限init", func() {
			conn, _ := mysql.GetConnect()
			tx, _ := conn.NewTx(nil)
			_, err := domain.RoleInit(1, tx)
			if err != nil {
				tx.Rollback()
				fmt.Println(err)
			}
			tx.Commit()
		})
	}))
}

func TestChangeDefaultRole(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		t.Log(getDataIdByIssueId(2373, 0, 10236478))
	}))
}
