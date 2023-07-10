package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestPrepareTransferOrg(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(PrepareTransferOrg(10000001, 10000002))
	}))
}
