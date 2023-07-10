package service

import (
	"context"
	"testing"

	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetRecycleList(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {
		resp, err := GetRecycleList(1242, 23131, 8526, 1, 1, 30)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestDeleteWorkHourForIssues1(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {
		err := domain.DeleteWorkHourForIssues(1000023, []int64{5846}, 2022042901)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("--end--")
	}))
}

func TestRecoverWorkHourForIssues1(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {
		err := domain.RecoverWorkHours(1000023, []int64{5846}, 2022042901)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("--end--")
	}))
}
