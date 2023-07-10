package facade

import (
	"errors"
	"fmt"
	"strconv"

	"google.golang.org/protobuf/encoding/protojson"

	"google.golang.org/protobuf/proto"

	"github.com/star-table/go-common/utils/unsafe"

	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/http"
	"github.com/star-table/polaris-backend/common/model/vo"
)

var log = logger.GetDefaultLogger()

func Request(method string, reqUrl string, params map[string]interface{}, headers []http.HeaderOption, req, resp interface{}) vo.Err {
	var (
		requestBody    []byte
		respBody       []byte
		respStatusCode int
		err            error
	)

	if bts, ok := req.([]byte); ok {
		requestBody = bts
	} else {
		requestBody, _ = json.Marshal(req)
	}

	switch method {
	case consts.HttpMethodGet:
		respBody, respStatusCode, err = http.Get(reqUrl, params, headers...)
	case consts.HttpMethodPost:
		respBody, respStatusCode, err = http.Post(reqUrl, params, requestBody, headers...)
	case consts.HttpMethodPut:
		respBody, respStatusCode, err = http.Put(reqUrl, params, requestBody, headers...)
	case consts.HttpMethodDelete:
		respBody, respStatusCode, err = http.Delete(reqUrl, params, requestBody, headers...)
	default:
		log.Errorf("request [%s] failed, unknown http method: %s", reqUrl, method)
		return vo.NewErr(errs.BuildSystemErrorInfoWithMessage(errs.ServerError, fmt.Sprintf("unknown http method: %s", method)))
	}
	if err != nil {
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
		return vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
	}

	if respStatusCode < 200 || respStatusCode > 299 {
		voErr := vo.Err{}
		if len(respBody) != 0 {
			_ = json.Unmarshal(respBody, &voErr)
		}
		if voErr.Failure() {
			log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, voErr.Error())
			return voErr
		}

		message := fmt.Sprintf("response code %d", respStatusCode)
		if voErr.Message != "" {
			message = voErr.Message
		}
		return vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(message)))
	}

	if s, ok := resp.(*string); ok {
		*s = unsafe.BytesString(respBody)
	} else {
		err = json.Unmarshal(respBody, resp)
		if err != nil {
			log.Errorf("request [%s] failed, resp json decode error. resp: %q, err [%v]", reqUrl, respBody, err)
			return vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, err))
		}
	}

	return vo.NewErr(nil)
}

func RequestWithCommonHeader(orgId int64, userId int64, method string, reqUrl string, params map[string]interface{}, headers []http.HeaderOption, req, resp interface{}) vo.Err {
	var (
		respBody       []byte
		respStatusCode int
		err            error
	)

	requestBody, _ := json.Marshal(req)

	headers = append(headers, http.HeaderOption{Name: consts.AppHeaderXMdOrgId, Value: strconv.FormatInt(orgId, 10)})
	headers = append(headers, http.HeaderOption{Name: consts.AppHeaderXMdUserId, Value: strconv.FormatInt(userId, 10)})

	switch method {
	case consts.HttpMethodGet:
		respBody, respStatusCode, err = http.Get(reqUrl, params, headers...)
	case consts.HttpMethodPost:
		respBody, respStatusCode, err = http.Post(reqUrl, params, requestBody, headers...)
	case consts.HttpMethodPut:
		respBody, respStatusCode, err = http.Put(reqUrl, params, requestBody, headers...)
	case consts.HttpMethodDelete:
		respBody, respStatusCode, err = http.Delete(reqUrl, params, requestBody, headers...)
	default:
		log.Errorf("request [%s] failed, unknown http method: %s", reqUrl, method)
		return vo.NewErr(errs.BuildSystemErrorInfoWithMessage(errs.ServerError, fmt.Sprintf("unknown http method: %s", method)))
	}
	if err != nil {
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
		return vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
	}

	if respStatusCode < 200 || respStatusCode > 299 {
		voErr := vo.Err{}
		if len(respBody) != 0 {
			_ = json.Unmarshal(respBody, &voErr)
		}
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, voErr.Error().Error())
		if voErr.Failure() {
			return voErr
		}

		message := fmt.Sprintf("response code %d", respStatusCode)
		if voErr.Message != "" {
			message = voErr.Message
		}
		return vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(message)))
	}

	if bts, ok := resp.(*[]byte); ok {
		*bts = respBody
	} else {
		err = json.Unmarshal(respBody, resp)
		if err != nil {
			log.Errorf("request [%s] failed, resp json decode error. resp: %q, err [%v]", reqUrl, respBody, err)
			return vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, err))
		}
	}

	return vo.NewErr(nil)
}

func RequestWithCommonHeaderGrpc(orgId int64, userId int64, method string, reqUrl string, params map[string]interface{}, headers []http.HeaderOption, req interface{}, resp proto.Message) vo.Err {
	var (
		respBody       []byte
		respStatusCode int
		err            error
	)

	requestBody, _ := json.Marshal(req)

	headers = append(headers, http.HeaderOption{Name: consts.AppHeaderXMdOrgId, Value: strconv.FormatInt(orgId, 10)})
	headers = append(headers, http.HeaderOption{Name: consts.AppHeaderXMdUserId, Value: strconv.FormatInt(userId, 10)})

	switch method {
	case consts.HttpMethodGet:
		respBody, respStatusCode, err = http.Get(reqUrl, params, headers...)
	case consts.HttpMethodPost:
		respBody, respStatusCode, err = http.Post(reqUrl, params, requestBody, headers...)
	case consts.HttpMethodPut:
		respBody, respStatusCode, err = http.Put(reqUrl, params, requestBody, headers...)
	case consts.HttpMethodDelete:
		respBody, respStatusCode, err = http.Delete(reqUrl, params, requestBody, headers...)
	default:
		log.Errorf("request [%s] failed, unknown http method: %s", reqUrl, method)
		return vo.NewErr(errs.BuildSystemErrorInfoWithMessage(errs.ServerError, fmt.Sprintf("unknown http method: %s", method)))
	}
	if err != nil {
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
		return vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
	}

	if respStatusCode < 200 || respStatusCode > 299 {
		log.Errorf("request [%s] failed, response status code [%d], resp [%q]", reqUrl, respStatusCode, respBody)
		message := fmt.Sprintf("response code %d", respStatusCode)
		return vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(message)))
	}

	err = protojson.Unmarshal(respBody, resp)
	if err != nil {
		log.Errorf("request [%s] failed, resp json decode error. resp [%q], err [%v]", reqUrl, respBody, err)
		return vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, err))
	}

	return vo.NewErr(nil)
}

func ToJsonString(obj interface{}) string {
	return json.ToJsonIgnoreError(obj)
}
