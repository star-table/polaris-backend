package dao

import (
	"context"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPostGreeter_Departments(t *testing.T) {
	convey.Convey("Test sql", t, test.StartUp(func(ctx context.Context) {

		_, _, _ = OrcConfigPageList(1, 100)

	}))
}
