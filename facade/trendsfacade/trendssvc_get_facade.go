package trendsfacade

//func TrendList(req trendsvo.TrendListReqVo) trendsvo.TrendListRespVo {
//	respVo := &trendsvo.TrendListRespVo{}
//
//	reqUrl := fmt.Sprintf("%s/api/trendssvc/trendList", config.GetPreUrl("trendssvc"))
//	queryParams := map[string]interface{}{}
//	queryParams["orgId"] = req.OrgId
//	queryParams["userId"] = req.UserId
//	requestBody := json.ToJsonIgnoreError(req.Input)
//	fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
//	fullUrl += "|" + requestBody
//
//	respBody, respStatusCode, err := http.Get(reqUrl, queryParams)
//
//	//Process the response
//	if err != nil {
//		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
//		return *respVo
//	}
//	//接口响应错误
//	if respStatusCode < 200 || respStatusCode > 299 {
//		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("trendssvc response code %d", respStatusCode))))
//		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
//		return *respVo
//	}
//	jsonConvertErr := json.FromJson(respBody, respVo)
//	if jsonConvertErr != nil {
//		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//	}
//	return *respVo
//}
