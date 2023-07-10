package mid

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	consts1 "github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/jtolds/gls"
)

var defaultLog = logger.GetDefaultLogger()

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		//ua := c.GetHeader("User-Agent")
		//defaultLog.Infof("user-agent:%s", ua)
		//获取语言
		lang := c.GetHeader(consts.AppHeaderLanguage)
		threadlocal.Mgr.SetValues(gls.Values{consts.AppHeaderLanguage: lang}, func() {
			c.Next()
		})
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,PM-TOKEN,PM-ORG,PM-PRO,PM-ENV,PM-PLAT,PM-VER,PM-TRACE-ID,PM-LANG,ADMIN-TOKEN")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length,  Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,PM-TOKEN,PM-ORG,PM-PRO,PM-ENV,PM-PLAT,PM-VER,PM-TRACE-ID,PM-LANG")

		c.Header("Access-Control-Allow-Credentials", "true")

		//表示隔24个小时才发起预检请求。也就是说，发送两次请求
		c.Header("Access-Control-Max-Age", "86400")

		//fmt.Println(c.Request.Context().Value(consts.TraceIdKey))
		//fmt.Println(c.Request.Context().Value(consts.HttpContextKey).(domains.HttpContext).StartTime)

		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(consts.AppHeaderTokenName)
		//兼容
		if token == "" {
			token = c.GetHeader("Token")
		}
		if token != "" {
			c.Set(consts.AppHeaderTokenName, token)
		}

		version := c.GetHeader(consts.AppHeaderVerName)
		if version != "" && strings.Index(version, ".") != -1 {
			c.Set(consts.AppHeaderVerName, version)
		}

		// 处理请求
		c.Next()
	}
}

func SentryMiddleware(applicationName string, env string, dsn string) gin.HandlerFunc {
	var connectSentryErr error
	// 如果已经有 sentry 实例了，则无需再初始化
	if sentry.CurrentHub().Client() == nil {
		connectSentryErr = sentry.Init(sentry.ClientOptions{
			Dsn:         dsn,
			ServerName:  applicationName,
			Environment: env,
		})
	}
	return func(c *gin.Context) {
		//捕获异常
		defer func() {
			if r := recover(); r != nil {
				errMsg := ""
				if err, ok := r.(error); ok {
					errMsg = err.Error()
				} else {
					errMsg = fmt.Sprint(r)
				}
				event := sentry.NewEvent()
				event.EventID = sentry.EventID(errMsg)
				event.Message = fmt.Sprintf("捕获到的错误：%s, 堆栈:%s\n", errMsg, stack.GetStack())
				event.Environment = env
				event.ServerName = applicationName
				event.Tags[consts1.TraceIdKey] = threadlocal.GetTraceId()
				if connectSentryErr == nil {
					sentry.CaptureEvent(event)
					sentry.Flush(time.Second * 5)
				} else {
					fmt.Println(connectSentryErr)
				}
				fmt.Println(event.Message)

				respObj := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfoWithMessage(errs.ServerError, errMsg))}
				c.String(200, json.ToJsonIgnoreError(respObj))
			}
		}()

		// 处理请求
		c.Next()
	}
}

//func OpenTracingMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		spanCtx, _ := tracing.SpanFromContext(c)
//		if spanCtx != nil {
//			span := opentracing.StartSpan(c.Request.RequestURI, opentracing.ChildOf(spanCtx))
//			ctx = opentracing.ContextWithSpan(ctx, span)
//			defer span.Finish()
//		}
//		// 处理请求
//		c.Next()
//	}
//}
