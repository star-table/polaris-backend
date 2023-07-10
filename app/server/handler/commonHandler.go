package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/star-table/polaris-backend/common/model/vo"

	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/facade"

	"github.com/star-table/polaris-backend/common/core/errs"

	"github.com/star-table/common/core/config"

	"github.com/gin-gonic/gin"
)

var CommonHandler = &commonHandler{}

type commonHandler struct {
}

func (ch *commonHandler) Handler(ctx *gin.Context) {
	url := ch.getUrl(ctx)
	if url == "" {
		Fail(ctx, errs.BuildSystemErrorInfoWithMessage(errs.PathNotFound, ""))
		return
	}

	var (
		respVo  string
		respErr vo.Err
	)
	switch ctx.Request.Method {
	case http.MethodPost:
		if ctx.Request.Body == nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		defer ctx.Request.Body.Close()
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			Fail(ctx, errs.BuildSystemErrorInfoWithMessage(errs.ParamError, err.Error()))
			return
		}
		respErr = facade.Request(consts.HttpMethodPost, url, map[string]interface{}{}, nil, body, &respVo)
	default:
		respErr = facade.Request(consts.HttpMethodGet, url, map[string]interface{}{}, nil, nil, &respVo)
	}

	if respErr.Failure() {
		Fail(ctx, respErr.Error())
	}

	SuccessJson(ctx, respVo)
}

func (ch *commonHandler) getUrl(c *gin.Context) string {
	var (
		orgId, userId int64
	)
	cacheUserInfo, err := GetCacheUserInfo(c)
	if err == nil {
		orgId = cacheUserInfo.OrgId
		userId = cacheUserInfo.UserId
	}

	allSvc := map[string]string{
		"orgsvc":      config.GetPreUrl("orgsvc"),
		"projectsvc":  config.GetPreUrl("projectsvc"),
		"trendssvc":   config.GetPreUrl("trendssvc"),
		"resourcesvc": config.GetPreUrl("resourcesvc"),
	}
	for key, url := range allSvc {
		if strings.HasPrefix(c.Request.RequestURI, "/newApi/rest/"+key) {
			url = url + strings.Replace(c.Request.RequestURI, "/newApi/rest/", "/api/", 1)
			if strings.Contains(c.Request.RequestURI, "?") {
				url = url + fmt.Sprintf("&orgId=%d&userId=%d", orgId, userId)
			} else {
				url = url + fmt.Sprintf("?orgId=%d&userId=%d", orgId, userId)
			}
			return url
		}
	}

	return ""
}
