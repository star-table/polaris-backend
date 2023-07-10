package middlewares

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/star-table/polaris-backend/common/core/errs"

	"github.com/star-table/polaris-backend/app/server/handler"
	"github.com/star-table/polaris-backend/facade/orgfacade"

	"github.com/star-table/common/core/logger"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/gin-gonic/gin"
)

var (
	log = logger.GetDefaultLogger()

	noNeedLoginPath = map[string]struct{}{
		"/newApi/rest/orgsvc/exchangeShareToken":        {},
		"/newApi/rest/orgsvc/getWeiXinRegisterUrl":      {},
		"/newApi/rest/orgsvc/getWeiXinInstallUrl":       {},
		"/newApi/rest/projectsvc/getShareViewInfoByKey": {},
	}

	// 允许调用的路径
	shareTokenApiPath = map[string]struct{}{
		"/api/rest/read/org/columns":               {},
		"/api/rest/project/menu/getConfig":         {},
		"/api/task":                                {},
		"/api/task?name=getIssueWorkHoursInfo":     {},
		"/api/rest/project/*/issue/filter":         {},
		"/api/rest/project/*/issue/filterForIssue": {},
		"/api/rest/project/*/columns":              {},
	}
)

// 极星用户登录验证中间件
func PolarisRestAuth(c *gin.Context) {
	// pmToken := c.GetHeader(consts.AppHeaderTokenName)
	// 在后续的 http 请求中回根据 GinContextKey 的值判断是否是 gin context。
	c.Set("GinContextKey", c)
	if checkIsNeedAuth(c) {
		cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(c)
		//cacheUserInfo := &bo.CacheUserInfoBo{UserId: 2168, OrgId: 1424}
		if err != nil {
			// errStack := TraceId(err.Error())
			log.Errorf("[PolarisRestAuth] err: %v", err)
			c.AbortWithStatusJSON(200, handler.Res{
				Code:    int32(err.Code()),
				Message: err.Message(),
			})
			//return
		} else {
			// 如果是分享类型，只有部分接口能使用
			if cacheUserInfo.TokenType == consts.TokenTypeShare {
				path := strings.Replace(c.Request.URL.Path, "/-1/", "/*/", -1)
				reg := regexp.MustCompile(`/\d*/`)
				path = reg.ReplaceAllString(path, "/*/")
				if _, ok := shareTokenApiPath[path]; !ok {
					c.AbortWithStatusJSON(200, handler.Res{
						Code:    int32(errs.TokenAuthError.Code()),
						Message: errs.TokenAuthError.Message(),
					})
				}
			}
		}
		c.Set(consts.AppUserInfoCacheKey, cacheUserInfo)
	}
	c.Next()
}

func checkIsNeedAuth(c *gin.Context) bool {
	if _, ok := noNeedLoginPath[c.Request.URL.Path]; ok {
		return false
	}

	return true
}

// print stack trace for debug
func Trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(1, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
