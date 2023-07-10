package domain

import (
	"context"
	"testing"

	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestDeleteAttachmentsForColumn(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		orgId := 2373
		userId := 29612
		projectId := 14054
		tableId := 1548943365410656256
		columnId := "_field_1658201035042"
		err := DeleteAttachmentsForColumn(int64(orgId), int64(userId), int64(projectId), int64(tableId), columnId)
		if err != nil {
			t.Log(err)
		}
	}))
}

func TestDeleteAttachmentsForOneTable(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		orgId := 2373
		userId := 29612
		projectId := 14054
		//appId := 1544950152769363970
		tableId := 1548943365410656256
		err := DeleteAttachmentsForOneTable(int64(orgId), int64(userId), int64(projectId), int64(tableId))
		if err != nil {
			t.Log(err)
		}
	}))
}
