package service

import (
	"context"
	"testing"

	json2 "github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestDingTalkInfo(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		msgMap := map[string]interface{}{
			"app": "projectsvc",
			"env": "dev",
			"msg": "【警告日志】this is a test warning...",
		}
		msgJson := json2.ToJsonIgnoreError(msgMap)
		param := vo.DingTalkInfoReq{
			Content: msgJson,
			Other:   "",
		}
		err := DingTalkInfo(param)
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("22222222222222-end...")
	}))
}

func TestPostWithTimeout(t *testing.T) {
	//convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
	//	t1 := time.Now()
	//	_, respStatusCode, err := http.PostWithTimeout("https://www.facebook.com", map[string]interface{}{}, `{"testBody"}`, 5)
	//	t2 := time.Since(t1).Milliseconds()
	//	log.Info(respStatusCode)
	//	log.Info(err)
	//	convey.So(t2, convey.ShouldBeGreaterThanOrEqualTo, 5000)
	//}))
}
