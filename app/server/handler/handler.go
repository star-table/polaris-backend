package handler

import (
	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/polaris-backend/app/api"
	"github.com/star-table/polaris-backend/app/generated"
	sconsts "github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/extra/trace/gin2micro"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
	"github.com/jtolds/gls"
)

// Defining the Graphql handler
func GraphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.GraphQL(generated.NewExecutableSchema(generated.Config{
		Resolvers: &api.Resolver{},
	}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func RestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
		openTraceId := ""
		traceContext, ok := gin2micro.ContextWithSpan(c)
		if ok && traceContext != nil {
			if v := traceContext.Value(gin2micro.OpenTraceKey); v != nil {
				openTraceId = v.(string)
			}
		}
		threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId, sconsts.JaegerContextTraceKey: openTraceId}, func() {
			c.Next()
		})
	}
}

// Defining the Playground handler
func PlaygroundHandler() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/query")
	return func(c *gin.Context) {
		httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
		threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
			h.ServeHTTP(c.Writer, c.Request)
		})
	}
}

func DingTalkCallBackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
		threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
			//dingtalk.DingTalkCallbackHandler(c.Writer, c.Request)
		})
	}
}

func HeartbeatHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "ok")
	}
}
