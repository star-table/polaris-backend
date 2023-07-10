package pushfacade

import (
	"errors"
	"fmt"

	"github.com/star-table/go-common/pkg/encoding"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/http"
	"github.com/star-table/polaris-backend/common/model/vo"
)

var log = logger.GetDefaultLogger()

var (
	ApiV1Prefix = "/inner/v1"
)

func request(reqUrl string, req, respVo interface{}) vo.Err {
	requestBody := marshalToString(req)
	respBody, respStatusCode, err := http.Post(reqUrl, nil, requestBody)
	//Process the response
	if err != nil {
		log.Errorf("request [%s] failed, body [%q], response status code [%d], err [%v]", reqUrl, requestBody, respStatusCode, err)
		return vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299 {
		voErr := vo.Err{}
		if len(respBody) != 0 {
			_ = unmarshalFromString(respBody, &voErr)
		}
		log.Errorf("request [%s] failed, body [%q], response status code [%d], err [%v]", reqUrl, requestBody, respStatusCode, voErr)
		message := fmt.Sprintf("tablesvc response code %d", respStatusCode)
		if voErr.Message != "" {
			message = voErr.Message
		}
		if voErr.Code != 0 {
			return voErr
		}

		return vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(message)))

	}
	jsonConvertErr := unmarshalFromString(respBody, respVo)
	if jsonConvertErr != nil {
		return vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}

	return vo.NewErr(nil)
}

func marshalToString(v interface{}) []byte {
	bts, _ := encoding.GetJsonCodec().Marshal(v)
	return bts
}

func unmarshalFromString(str []byte, v interface{}) error {
	return encoding.GetJsonCodec().Unmarshal(str, v)
}
