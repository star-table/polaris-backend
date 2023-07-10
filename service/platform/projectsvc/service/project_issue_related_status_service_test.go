package service

import (
	"context"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProjectIssueRelatedStatus(t *testing.T) {
	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
		res, err := ProjectIssueRelatedStatus(1013, vo.ProjectIssueRelatedStatusReq{ProjectID: 1319, TableID: "1479"})
		t.Log(json.ToJsonIgnoreError(res))
		t.Log(err)
	}))
}
