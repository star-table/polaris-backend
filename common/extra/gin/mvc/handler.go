package mvc

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/util/json"
	sconsts "github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/extra/trace/gin2micro"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/jtolds/gls"
	"github.com/opentracing/opentracing-go"
)

var log = logger.GetDefaultLogger()

type GinHandler struct {
	engine   *gin.Engine
	Greeters []interface{}
}

func NewGinHandler(engine *gin.Engine) GinHandler {
	return GinHandler{
		engine:   engine,
		Greeters: make([]interface{}, 0),
	}
}

func (h *GinHandler) RegisterGreeter(greeter interface{}) {
	h.Greeters = append(h.Greeters, greeter)

	greeterValue := reflect.ValueOf(greeter)
	greeterType := reflect.TypeOf(greeter)
	greeterInfo := GetGreeterInfo(greeter)

	methodNum := greeterType.NumMethod()

	for i := 0; i < methodNum; i++ {
		methodName := greeterType.Method(i).Name
		method := greeterValue.MethodByName(methodName)
		h.handleGreeterFunc(greeterInfo, methodName, method.Interface())
	}
}

func (h *GinHandler) handleGreeterFunc(greeter Greeter, funcName string, greeterFunc interface{}) {
	applicationName := greeter.ApplicationName
	version := greeter.Version
	httpType := greeter.HttpType
	switch httpType {
	case http.MethodGet:
		h.engine.GET(BuildApi(applicationName, version, funcName), h.tracingHandler(greeterFunc))
	case http.MethodPost:
		h.engine.POST(BuildApi(applicationName, version, funcName), h.tracingHandler(greeterFunc))
	default:
		log.Infof("not support http type %s", httpType)
	}
}

func (h *GinHandler) tracingHandler(targetFunc interface{}) gin.HandlerFunc {
	handler := h.restHandler(targetFunc)
	return func(c *gin.Context) {
		//token := c.GetHeader(consts.APP_HEADER_TOKEN_NAME)
		httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
		openTraceId := ""
		var span opentracing.Span
		traceContext, ok := gin2micro.ContextWithSpan(c)
		if ok && traceContext != nil {
			if v := traceContext.Value(gin2micro.OpenTraceKey); v != nil {
				openTraceId = v.(string)
			}
			span = opentracing.SpanFromContext(traceContext)
		}
		threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId, sconsts.JaegerContextTraceKey: openTraceId, sconsts.JaegerContextSpanKey: span}, func() {
			handler(c)
		})
	}
}

func (h *GinHandler) restHandler(targetFunc interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var respObj interface{} = nil

		reqInfo, err := NewRequestInfo(ctx)
		if err != nil {
			log.Error(err)
			respObj = vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfoWithMessage(errs.ServerError, err.Error()))}
		}
		startTime := time.Now()
		respData, err := InjectFunc(targetFunc, *reqInfo)
		elapsed := time.Since(startTime)

		if err != nil {
			log.Error(err)
			respObj = vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfoWithMessage(errs.ServerError, err.Error()))}
		} else {
			if len(respData) > 0 {
				respObj = respData[0].Interface()
			}
		}
		if respObj != nil {
			var respBody = ""
			if s, ok := respObj.(string); ok {
				respBody = s
			} else {
				respBody = json.ToJsonIgnoreError(respObj)
			}

			respContent(ctx, 200, respBody)

			//截取下日志长度
			cutRespBody := respBody
			if len(cutRespBody) > 5000 {
				cutRespBody = cutRespBody[:5000]
			}
			// 无效日志太多，健康检查不进行日志的打印
			if !strings.Contains(ctx.Request.URL.String(), "/health") {
				log.Infof("请求处理完成，总耗时-> [%dms], url-> [%s], respBody-> [%s]", elapsed/1e6, ctx.Request.URL, cutRespBody)
			}
		}
	}
}

func respContent(ctx *gin.Context, httpStatus int, resp string) {
	ctx.Status(httpStatus)

	if !bodyAllowedForStatus(httpStatus) {
		render.String{Format: resp}.WriteContentType(ctx.Writer)
		ctx.Writer.WriteHeaderNow()
		return
	}
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(ctx.Writer, resp)
}

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

func NewRequestInfo(ctx *gin.Context) (*RequestInfo, error) {
	reqInfo := RequestInfo{
		Parameters: map[string]string{},
		Ctx:        ctx.Request.Context(),
	}
	bytes, err := ctx.GetRawData()
	if err != nil {
		return nil, err
	}
	reqInfo.Body = string(bytes)
	for key, params := range ctx.Request.URL.Query() {
		reqInfo.Parameters[key] = params[0]
	}
	return &reqInfo, nil
}

//func (h *GinHandler) Post(targetFunc interface{}){
//	h.engine.POST(BuildApiByFunc(h.applicationName, h.version, targetFunc), RestHandler(targetFunc))
//}
//
//func (h *GinHandler) Get(targetFunc interface{}){
//	h.engine.GET(BuildApiByFunc(h.applicationName, h.version, targetFunc), RestHandler(targetFunc))
//}
//
//func (h *GinHandler) PostVersion(version string, targetFunc interface{}){
//	h.engine.POST(BuildApiByFunc(h.applicationName, version, targetFunc), RestHandler(targetFunc))
//}
//
//func (h *GinHandler) GetVersion(version string, targetFunc interface{}){
//	h.engine.GET(BuildApiByFunc(h.applicationName, version, targetFunc), RestHandler(targetFunc))
//}

//func BuildApiByFunc(applicationName, version string, targetFunc interface{}) string{
//	funcPath := runtime.FuncForPC(reflect.ValueOf(targetFunc).Pointer()).Name()
//	strs := strings.Split(funcPath, ".")
//	funcName := strs[len(strs) - 1]
//	if strings.Index(funcName, "-") > -1{
//		strs = strings.Split(funcName, "-")
//		funcName = strs[0]
//	}
//	return BuildApi(applicationName, version, funcName)
//}

func BuildApi(applicationName, version string, funcName string) string {
	apiName := str.LcFirst(funcName)
	if version == "" {
		return fmt.Sprintf("api/%s/%s", applicationName, apiName)
	} else {
		return fmt.Sprintf("api/%s/%s/%s", applicationName, version, apiName)
	}
}
