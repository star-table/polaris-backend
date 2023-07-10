package util

import (
	"github.com/star-table/common/core/util/temp"
	"github.com/star-table/polaris-backend/common/core/errs"
)

func ParseCacheKey(key string, params map[string]interface{}) (string, errs.SystemErrorInfo) {

	target, err := temp.Render(key, params)
	if err != nil {
		log.Error(err)
		return "", errs.BuildSystemErrorInfo(errs.TemplateRenderError, err)
	}
	return target, nil
}
