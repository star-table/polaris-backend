package service

import (
	"context"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUploadOssByFsImageKey(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		t.Log(  UploadOssByFsImageKey(1004, "img_ead2e693-606e-40a0-a979-80243a19ee8g", false))
	}))
}
