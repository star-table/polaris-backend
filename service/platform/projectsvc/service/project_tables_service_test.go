package service

import (
	"context"
	"testing"

	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestProjectSupportObjectTypes(t *testing.T) {
	//convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
	//	res, err := ProjectSupportObjectTypes(1872, vo.ProjectSupportObjectTypeListReq{
	//		ProjectID: 12898,
	//	})
	//	if err != nil {
	//		t.Error(err)
	//		return
	//	}
	//	t.Log(json.ToJsonIgnoreError(res))
	//}))
}

func TestCheckAsyncTaskIsExecuting1(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		isExecuting := domain.CheckAsyncTaskIsRunning(2373, 1560151305056518146, 1560151306109194240)
		t.Log(isExecuting)
	}))
}
