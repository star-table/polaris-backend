package service

import (
	"context"
	"github.com/star-table/polaris-backend/common/model/vo/msgvo"
	"testing"

	"github.com/star-table/polaris-backend/service/basic/msgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestFixAddOrderForFeiShu(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		if err := FixAddOrderForFeiShu(msgvo.FixAddOrderForFeiShuReqData{
			StartTime: nil,
			EndTime:   nil,
		}); err != nil {
			t.Error(err)
			return
		}
		t.Log("...end...")
	}))
}
