package openapi

import (
	"strconv"
	"strings"

	"github.com/star-table/go-common/utils/unsafe"

	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/gin-gonic/gin"
)

var log = logger.GetDefaultLogger()

const (
	// OK ok
	OK int32 = 0

	// RequestErr request error
	RequestErr int32 = -400

	// ServerErr server error
	ServerErr int32 = -500

	contextErrCode = "context/err/code"
)

type res struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Fail(c *gin.Context, err errs.SystemErrorInfo) {
	errCode := RequestErr
	errMsg := "request err"

	if err != nil {
		errCode = int32(err.Code())
		errMsg = err.Message()
	}
	c.Set(contextErrCode, errCode)
	c.Data(200, "application/json;charset=utf8",
		json.ToJsonBytesIgnoreError(res{
			Code:    errCode,
			Message: errMsg,
		}))
}

func Suc(c *gin.Context, data interface{}) {
	code := OK
	c.Set(contextErrCode, code)
	c.Data(200, "application/json;charset=utf8",
		json.ToJsonBytesIgnoreError(res{
			Code: code,
			Data: data,
		}))
}

func SuccessJson(c *gin.Context, data string) {
	code := OK
	c.Set(contextErrCode, code)
	c.Data(200, "application/json;charset=utf8", unsafe.StringBytes(data))
}

// parseOpenAuthInfo 获取Open认证信息
func ParseOpenAuthInfo(c *gin.Context) (*orgvo.OpenAPIAuthData, errs.SystemErrorInfo) {
	accessToken := c.GetHeader("Authorization")
	orgID := c.GetHeader("X-Tenant-Id")
	if accessToken == "" || orgID == "" {
		return nil, errs.OpenAccessTokenInvalid
	}
	index := strings.Index(accessToken, "Bearer ")
	if index == -1 {
		return nil, errs.OpenAccessTokenInvalid
	}
	accessToken = accessToken[index+7:]

	orgId, err := strconv.ParseInt(orgID, 10, 64)
	if err != nil {
		return nil, errs.OpenAccessTokenInvalid
	}

	openApiAuthResp := orgfacade.OpenAPIAuth(orgvo.OpenAPIAuthReq{
		AccessToken: accessToken,
		OrgID:       orgId,
	})
	if openApiAuthResp.Failure() {
		log.Error(openApiAuthResp.Message)
		return nil, openApiAuthResp.Error()
	}
	return openApiAuthResp.Data, nil
}

func ParseInt(str string) int {
	return int(ParseInt64(str))
}

func ParseInt64(str string) int64 {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return v
}
