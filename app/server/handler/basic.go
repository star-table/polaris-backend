package handler

import (
	"github.com/star-table/go-common/utils/unsafe"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/gin-gonic/gin"
)

const (
	// OK ok
	OK int32 = 0

	// RequestErr request error
	RequestErr int32 = -400

	// ServerErr server error
	ServerErr int32 = -500

	contextErrCode = "context/err/code"
)

type Res struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

// 从上下文中获取缓存数据
func GetCacheUserInfo(c *gin.Context) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	res := &bo.CacheUserInfoBo{}
	val, exist := c.Get(consts.AppUserInfoCacheKey)
	if !exist {
		return nil, errs.TokenAuthError
	}
	res = val.(*bo.CacheUserInfoBo)
	return res, nil
}

func Fail(c *gin.Context, err errs.SystemErrorInfo) {
	errCode := RequestErr
	errMsg := "request err"

	if err != nil {
		errCode = int32(err.Code())
		errMsg = err.Message()
	}
	c.Set(contextErrCode, errCode)
	c.JSON(200, Res{
		Code:    errCode,
		Success: false,
		Message: errMsg,
	})
}

func Success(c *gin.Context, data interface{}) {
	code := OK
	c.Set(contextErrCode, code)
	c.Data(200, "application/json;charset=utf8",
		json.ToJsonBytesIgnoreError(Res{
			Code:    code,
			Success: true,
			Data:    data,
		}))
}

func SuccessJson(c *gin.Context, data string) {
	code := OK
	c.Set(contextErrCode, code)
	c.Data(200, "application/json;charset=utf8", unsafe.StringBytes(data))
}
