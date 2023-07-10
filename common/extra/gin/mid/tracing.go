package mid

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/util/network"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/times"
	"github.com/star-table/common/core/util/uuid"
	consts2 "github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/extra/trace/gin2micro"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var log = logger.GetRequestResponseLogger()

func newTraceId() string {
	return fmt.Sprintf("%s_%s", network.GetIntranetIp(), uuid.NewUuid())
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func isBinaryContent(contentType string) bool {
	return strings.Contains(contentType, "image") || strings.Contains(contentType, "video") ||
		strings.Contains(contentType, "audio")
}

func isMultipart(contentType string) bool {
	return strings.Contains(contentType, "multipart/form-data")
}

func StartTrace() gin.HandlerFunc {
	return func(c *gin.Context) {

		traceId := c.Request.Header.Get(consts.TraceIdKey)
		if "" == traceId {
			traceId = newTraceId()
			traceContext, ok := gin2micro.ContextWithSpan(c)
			if ok && traceContext != nil {
				if v := traceContext.Value(gin2micro.OpenTraceKey); v != nil {
					traceId = v.(string)
				}
			}
		}

		contentType := c.ContentType()
		//postForm := ""
		body := ""
		httpContext := model.HttpContext{
			Request:   c.Request,
			Url:       c.Request.RequestURI,
			Method:    c.Request.Method,
			StartTime: times.GetNowMillisecond(),
			EndTime:   0,
			Ip:        c.ClientIP(),
			TraceId:   traceId,
		}

		userToken := c.GetHeader(consts2.AppHeaderTokenName)
		orgId := c.GetHeader(consts2.AppHeaderOrgName)
		env := c.GetHeader(consts2.AppHeaderEnvName)
		plat := c.GetHeader(consts2.AppHeaderPlatName)
		proId := c.GetHeader(consts2.AppHeaderProName)
		ver := c.GetHeader(consts2.AppHeaderVerName)
		lang := c.GetHeader(consts2.AppHeaderLanguage)

		if !isBinaryContent(contentType) && !isMultipart(contentType) {
			// 判断不是上传文件等大消息体，记录消息体日志
			//c.Request.ParseForm()
			//postForm = c.Request.PostForm.Encode()
			data, err := c.GetRawData()
			if err != nil {
				log.Info(err.Error())
			}
			body = string(data)
			// 重新写入body
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}

		ctx := context.WithValue(c.Request.Context(), consts.TraceIdKey, traceId)
		ctx = context.WithValue(ctx, consts.HttpContextKey, httpContext)
		c.Request = c.Request.WithContext(ctx)

		c.Header(consts.TraceIdKey, traceId)

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		urlCount := strs.Len(httpContext.Url)
		if urlCount <= 6 || httpContext.Url[urlCount-6:] != "health" {
			// 仅记录非心跳日志

			httpContext = c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
			httpContext.Status = c.Writer.Status()
			httpContext.EndTime = times.GetNowMillisecond()
			runtime := httpContext.EndTime - httpContext.StartTime
			runtimes := strconv.FormatInt(runtime, 10)
			httpStatus := strconv.Itoa(httpContext.Status)
			respBody := blw.body.String()
			if len(respBody) > 1000 {
				respBody = respBody[:1000]
			}
			msg := fmt.Sprintf("request|traceId=%s|start=%s|ip=%s|contentType=%s|method=%s|url=%s|header={\"ver\":\"%s\","+
				"\"token\":\"%s\",\"plat\":\"%s\",\"env\":\"%s\",\"org\":\"%s\",\"pro\":\"%s\",\"lang\":\"%s\"}|body=%s|------response|end=%s|time=%s|status=%s|body=%s|",
				httpContext.TraceId, times.GetDateTimeStrByMillisecond(httpContext.StartTime), httpContext.Ip, contentType,
				httpContext.Method, httpContext.Url, ver, userToken, plat, env, orgId, proId, lang, strings.ReplaceAll(body, "\n", "\\n"), times.GetDateTimeStrByMillisecond(httpContext.EndTime),
				runtimes, httpStatus, respBody)

			log.InfoH(&httpContext, msg, zap.String("duration", runtimes),
				zap.String("responseCode", httpStatus), zap.String("path", httpContext.Url))
		}
	}
}
