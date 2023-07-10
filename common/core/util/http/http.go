package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/util/times"
	sconsts "github.com/star-table/polaris-backend/common/core/consts"
	pkgErr "github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/uber/jaeger-client-go"
)

const defaultContentType = "application/json"

var httpClient = &http.Client{}
var log = logger.GetDefaultLogger()

type HeaderOption struct {
	Name  string
	Value string
}

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(30) * time.Second,
	}
}

func responseHandle(resp *http.Response, err error) ([]byte, int, error) {
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, resp.StatusCode, err
	}
	return b, resp.StatusCode, nil
}

// 不可存在结构体值，否则会报错。
func ConvertToQueryParams(params map[string]interface{}) string {
	if params == nil {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("?")
	for k, v := range params {
		vStr := serialToString(v)
		buffer.WriteString(fmt.Sprintf("%s=%v&", k, vStr))
	}
	buffer.Truncate(buffer.Len() - 1)
	return buffer.String()
}

func isNil(dest interface{}) bool {
	if dest == nil {
		return true
	}
	v := reflect.ValueOf(dest)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return false
}

// 将数据转化为字符串
func serialToString(dest interface{}) string {
	var key string
	if isNil(dest) {
		return key
	}
	switch dest.(type) {
	case float64:
		key = decimal.NewFromFloat(dest.(float64)).String()
	case *float64:
		if dest.(*float64) != nil {
			key = decimal.NewFromFloat(*dest.(*float64)).String()
		}
	case float32:
		key = decimal.NewFromFloat32(dest.(float32)).String()
	case *float32:
		if dest.(*float32) != nil {
			key = decimal.NewFromFloat32(*dest.(*float32)).String()
		}
	case int:
		key = strconv.Itoa(dest.(int))
	case *int:
		if dest.(*int) != nil {
			key = strconv.Itoa(*dest.(*int))
		}
	case uint:
		key = strconv.Itoa(int(dest.(uint)))
	case *uint:
		key = strconv.Itoa(int(*dest.(*uint)))
	case int8:
		key = strconv.Itoa(int(dest.(int8)))
	case *int8:
		if dest.(*int8) != nil {
			key = strconv.Itoa(int(*dest.(*int8)))
		}
	case uint8:
		key = strconv.Itoa(int(dest.(uint8)))
	case *uint8:
		if dest.(*uint8) != nil {
			key = strconv.Itoa(int(*dest.(*uint8)))
		}
	case int16:
		key = strconv.Itoa(int(dest.(int16)))
	case *int16:
		if dest.(*int16) != nil {
			key = strconv.Itoa(int(*dest.(*int16)))
		}
	case uint16:
		key = strconv.Itoa(int(dest.(uint16)))
	case *uint16:
		if dest.(*uint16) != nil {
			key = strconv.Itoa(int(*dest.(*uint16)))
		}
	case int32:
		key = strconv.Itoa(int(dest.(int32)))
	case *int32:
		if dest.(*int32) != nil {
			key = strconv.Itoa(int(*dest.(*int32)))
		}
	case uint32:
		key = strconv.Itoa(int(dest.(uint32)))
	case *uint32:
		if dest.(*uint32) != nil {
			key = strconv.Itoa(int(*dest.(*uint32)))
		}
	case int64:
		key = strconv.FormatInt(dest.(int64), 10)
	case *int64:
		if dest.(*int64) != nil {
			key = strconv.FormatInt(*dest.(*int64), 10)
		}
	case uint64:
		key = strconv.FormatUint(dest.(uint64), 10)
	case *uint64:
		if dest.(*uint64) != nil {
			key = strconv.FormatUint(*dest.(*uint64), 10)
		}
	case string:
		key = dest.(string)
	case *string:
		if dest.(*string) != nil {
			key = *dest.(*string)
		}
	case []byte:
		key = string(dest.([]byte))
	case *[]byte:
		if dest.(*[]byte) != nil {
			key = string(*dest.(*[]byte))
		}
	case bool:
		if dest.(bool) {
			key = "true"
		} else {
			key = "false"
		}
	case *bool:
		if dest.(*bool) != nil {
			if *dest.(*bool) {
				key = "true"
			} else {
				key = "false"
			}
		}
	default:
	}
	return key
}

func Request(method string, url string, params map[string]interface{}, body []byte, headerOptions ...HeaderOption) ([]byte, int, error) {
	fullUrl := url + ConvertToQueryParams(params)
	req, err := http.NewRequest(method, fullUrl, bytes.NewReader(body))
	if err != nil {
		return nil, 0, pkgErr.Wrap(err, "http request error")
	}

	if len(body) > 0 {
		req.Header.Set("Content-Type", defaultContentType)
	}
	req.Header.Set(consts.TraceIdKey, threadlocal.GetTraceId())
	req.Header.Set(sconsts.HttpHeaderKratosTraceId, threadlocal.GetTraceId())
	req.Header.Set(sconsts.AppHeaderLanguage, threadlocal.GetValue(sconsts.AppHeaderLanguage))
	req.Header.Set(jaeger.TraceContextHeaderName, threadlocal.GetValue(sconsts.JaegerContextTraceKey))

	for _, headerOption := range headerOptions {
		req.Header.Set(headerOption.Name, headerOption.Value)
	}

	log.Infof("[HTTP] %s | request [%s] starting | body [%q] | headers [%v]", method, fullUrl, body, req.Header)

	start := times.GetNowMillisecond()
	resp, err := httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	end := times.GetNowMillisecond()
	timeConsuming := strconv.FormatInt(end-start, 10)

	respBody, httpCode, err := responseHandle(resp, err)
	cutRespBody := respBody
	if len(cutRespBody) > 5000 {
		cutRespBody = cutRespBody[:5000]
	}
	log.Infof("[HTTP] %s | request [%s] successful | body [%q] | headers [%v] | response status code [%d] | response body [%q] | time-consuming [%s]", method, fullUrl, body, req.Header, httpCode, cutRespBody, timeConsuming)
	return respBody, httpCode, err
}

func Get(url string, params map[string]interface{}, headerOptions ...HeaderOption) ([]byte, int, error) {
	return Request(sconsts.HttpMethodGet, url, params, nil, headerOptions...)
}

func Post(url string, params map[string]interface{}, body []byte, headerOptions ...HeaderOption) ([]byte, int, error) {
	return Request(sconsts.HttpMethodPost, url, params, body, headerOptions...)
}

func Put(url string, params map[string]interface{}, body []byte, headerOptions ...HeaderOption) ([]byte, int, error) {
	return Request(sconsts.HttpMethodPut, url, params, body, headerOptions...)
}

func Delete(url string, params map[string]interface{}, body []byte, headerOptions ...HeaderOption) ([]byte, int, error) {
	return Request(sconsts.HttpMethodDelete, url, params, body, headerOptions...)
}

//func PostWithTimeout(url string, params map[string]interface{}, body []byte, timeout uint32, headerOptions ...HeaderOption) ([]byte, int, error) {
//	method := sconsts.HttpMethodPost
//	fullUrl := url + ConvertToQueryParams(params)
//	req, err := http.NewRequest(method, fullUrl, bytes.NewReader(body))
//	if err != nil {
//		return nil, 0, pkgErr.Wrap(err, "http post request error")
//	}
//
//	req.Header.Set("Content-Type", defaultContentType)
//	req.Header.Set(consts.TraceIdKey, threadlocal.GetTraceId())
//	req.Header.Set(sconsts.AppHeaderLanguage, threadlocal.GetValue(sconsts.AppHeaderLanguage))
//	req.Header.Set(jaeger.TraceContextHeaderName, threadlocal.GetValue(sconsts.JaegerContextTraceKey))
//	for _, headerOption := range headerOptions {
//		req.Header.Set(headerOption.Name, headerOption.Value)
//	}
//
//	headers, _ := json.Marshal(req.Header)
//	log.Infof("[HTTP] %s | request [%s] starting | body [%q] | headers [%q]", method, fullUrl, body, headers)
//
//	start := times.GetNowMillisecond()
//	resp, err := getHttpClientWithTimeout(timeout).Do(req)
//	if resp != nil {
//		defer resp.Body.Close()
//	}
//	end := times.GetNowMillisecond()
//	timeConsuming := strconv.FormatInt(end-start, 10)
//
//	respBody, httpCode, err := responseHandle(resp, err)
//	cutRespBody := respBody
//	if len(cutRespBody) > 5000 {
//		cutRespBody = cutRespBody[:5000]
//	}
//	log.Infof("[HTTP] %s | request [%s] successful | body [%q] | headers [%q] | response status code [%d] | response body [%q] | time-consuming [%s]", fullUrl, body, headers, httpCode, cutRespBody, timeConsuming)
//	return respBody, httpCode, err
//}
