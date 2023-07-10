package domain

import (
	"context"
	"testing"

	"github.com/star-table/polaris-backend/service/platform/resourcesvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetResourceRelations(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		orgId := 2373
		isDelete := 1
		versionId := 7600
		projectId := 14213
		resourceIds := []int64{1004464}

		relations, err := GetResourceRelations(int64(orgId), int64(projectId), versionId, isDelete, nil, resourceIds, true)
		if err != nil {
			t.Error(err)
		}
		t.Log(relations)
	}))

}
