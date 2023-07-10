package service

import (
	"context"
	"testing"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/getsentry/sentry-go"
	"github.com/smartystreets/goconvey/convey"
)

//func TestProjectTypeCategory(t *testing.T) {
//	convey.Convey("Test ProjectTypeCategory", t, test.StartUp(func(ctx context.Context) {
//		res, err := ProjectTypeCategory(1001)
//		t.Log(json.ToJsonIgnoreError(res))
//		assert.Equal(t, err, nil)
//	}))
//}
//
//func TestProjectTypeList(t *testing.T) {
//	convey.Convey("Test ProjectTypeList", t, test.StartUp(func(ctx context.Context) {
//		res, err := ProjectTypeList(1, 1872)
//		if err != nil {
//			log.Error(err)
//			return
//		}
//		t.Log(json.ToJsonIgnoreError(res))
//	}))
//}

func TestSendSentry(t *testing.T) {
	convey.Convey("TestSendSentry", t, test.StartUp(func(ctx context.Context) {
		sentryInit()
		log.SetExtraLoggerOption("sentryClient", sentry.CurrentHub().Client())
		log.Info("init sentryClient ok")
		log.Error(errs.FileReadFail)
		//log.Error(errs.FileNotExist)
	}))
}

func sentryInit() {
	sentryConfig := config.GetSentryConfig()
	sentryDsn := ""
	if sentryConfig != nil {
		sentryDsn = sentryConfig.Dsn
	}
	dsn := sentryDsn
	applicationName := config.GetApplication().Name
	env := "dev"
	_ = sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		ServerName:  applicationName,
		Environment: env,
	})
}
