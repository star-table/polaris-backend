package service

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	vo1 "github.com/star-table/polaris-backend/common/model/vo"
)

// 日志发送给到钉钉群
func DingTalkInfo(param vo1.DingTalkInfoReq) errs.SystemErrorInfo {
	//dingTalkLogBotConfig := config.GetLogConfig("dingtalklogbot")
	//if dingTalkLogBotConfig == nil {
	//	log.Errorf("DingTalkLogBotUrl config is nil")
	//	return errs.DingTalkLogBotConfigError
	//}
	//reqUrl := dingTalkLogBotConfig.LogPath
	//if reqUrl == "" || len(reqUrl) < 10 {
	//	log.Errorf("DingTalkLogBotUrl config is error 1")
	//	return errs.DingTalkLogBotConfigError
	//}
	//queryParams := map[string]interface{}{}
	//logContent := "【警告日志】| " + param.Content + param.Other
	//requestBodyMap := map[string]interface{}{
	//	"msgtype": "text",
	//	"text": map[string]string{
	//		"content": logContent,
	//	},
	//}
	//requestBodyJson := json2.ToJsonIgnoreError(requestBodyMap)
	//respBody, respStatusCode, err := http.PostWithTimeout(reqUrl, queryParams, requestBodyJson, 2)
	//if err != nil {
	//	log.Errorf("request [%s] failed, response status code [%d], respBody：[%s], err [%v]", reqUrl, respStatusCode, respBody, err)
	//	return errs.BuildSystemErrorInfo(errs.RequestError, err)
	//}

	return nil
}
