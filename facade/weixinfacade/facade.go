package weixinfacade

import (
	"errors"
	"fmt"

	"github.com/star-table/common/core/util/json"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/polaris-backend/common/core/errs"

	"github.com/star-table/polaris-backend/common/core/util/http"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
)

var log = logger.GetDefaultLogger()
var url = "https://api.weixin.qq.com"

type errInterface interface {
	GetCode() int
}

func GetAccessTokenByCode(code string) (*orgvo.AccessTokenResp, error) {
	respVo := &orgvo.AccessTokenResp{}

	reqUrl := fmt.Sprintf("%s/sns/oauth2/access_token", url)
	queryParams := map[string]interface{}{
		"code":       code,
		"appid":      config.GetConfig().PersonWeiXin.AppId,
		"secret":     config.GetConfig().PersonWeiXin.AppSecret,
		"grant_type": "authorization_code",
	}
	respBody, respStatusCode, err := http.Get(reqUrl, queryParams)

	return respVo, getResult(respBody, respStatusCode, err, respVo)
}

func GetQrCode(accessToken string, scene string) (*orgvo.WeiXinQrCodeReply, error) {
	req := &orgvo.WeiXinQrCodeReq{
		ExpireSeconds: 604800,
		ActionName:    "QR_STR_SCENE",
		ActionInfo:    &orgvo.ActionInfo{Scene: &orgvo.Scene{SceneStr: scene}},
	}
	respVo := &orgvo.WeiXinQrCodeReply{}
	body, _ := json.Marshal(req)

	reqUrl := fmt.Sprintf("%s/cgi-bin/qrcode/create?access_token=%s", url, accessToken)
	respBody, respStatusCode, err := http.Post(reqUrl, map[string]interface{}{}, body)

	return respVo, getResult(respBody, respStatusCode, err, respVo)
}

func getResult(respBody []byte, respStatusCode int, err error, resp errInterface) error {
	//Process the response
	if err != nil {
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", url, respStatusCode, err)
		return err
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299 {
		return errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("response code %d", respStatusCode)))
	}
	jsonConvertErr := json.Unmarshal(respBody, resp)
	if jsonConvertErr != nil {
		return errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr)
	}
	if resp.GetCode() != 0 {
		log.Errorf("[GetAccessTokenByCode] resp:%v", resp)
		return errs.ServerError
	}

	return nil
}
