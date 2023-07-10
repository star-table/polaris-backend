package service

import (
	"context"
	"testing"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/domain"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetPolicy1(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(2739)
		res, err := domain.GetPayFunctionLimitResourceNum(orgId, consts.FunctionAttachSizeLimit)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}
