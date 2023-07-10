package handler

//func UploadApkInfoHandler() gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//
//		reqUrl := fmt.Sprintf("%s/api/mct/apk/uploadApkInfo", config.GetPreUrl("compatsvc"))
//
//		reqVo := &compatsvc.UploadApkInfoReqVo{}
//
//		if c.BindJSON(reqVo) != nil {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.JSONConvertError))
//			return
//		}
//
//		//权限校验
//		if reqVo.OrgId != cacheUserInfo.OrgId || reqVo.Uid != cacheUserInfo.UserId {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//		//请求响应
//		respVo := &compatsvc.UploadApkInfoRespVo{}
//
//		respBody, respStatusCode, err := http.Post(reqUrl, nil, json.ToJsonIgnoreError(reqVo))
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//}
//
//func GetAllApkHandler() gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//		uid := c.Query("uid")
//		orgId := c.Query("orgId")
//
//		//权限校验
//		if uid != strconv.FormatInt(cacheUserInfo.UserId, 10) || orgId != strconv.FormatInt(cacheUserInfo.OrgId, 10) {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//
//		respVo := &compatsvc.GetAllApkRespListVo{}
//		reqUrl := fmt.Sprintf("%s/api/mct/apk/getAllApk", config.GetPreUrl("compatsvc"))
//
//		queryParams := map[string]interface{}{}
//
//		queryParams["uid"] = uid
//		queryParams["orgId"] = orgId
//		queryParams["page"] = c.Query("page")
//		queryParams["size"] = c.Query("size")
//
//		fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
//		respBody, respStatusCode, err := http.Get(reqUrl, queryParams)
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//
//	}
//}
//
//func DeleteApkHandler() gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//
//		reqUrl := fmt.Sprintf("%s/api/mct/apk/deleteApk", config.GetPreUrl("compatsvc"))
//		//获取解析参数
//		reqVo := &compatsvc.DeleteApkReqVo{}
//
//		if c.BindJSON(reqVo) != nil {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.JSONConvertError))
//			return
//		}
//
//		//权限校验
//		if reqVo.OrgId != cacheUserInfo.OrgId || reqVo.Uid != cacheUserInfo.UserId {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//		//请求响应
//		respVo := &compatsvc.DeleteApkRespVo{}
//
//		respBody, respStatusCode, err := http.Post(reqUrl, nil, json.ToJsonIgnoreError(reqVo))
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//}
//
//func GetAllDevicesHandler() gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//		uid := c.Query("uid")
//		orgId := c.Query("orgId")
//
//		//权限校验
//		if uid != strconv.FormatInt(cacheUserInfo.UserId, 10) || orgId != strconv.FormatInt(cacheUserInfo.OrgId, 10) {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//
//		respVo := &compatsvc.GetAllDevicesRespListVo{}
//		reqUrl := fmt.Sprintf("%s/api/mct/devices/getAllDevices", config.GetPreUrl("compatsvc"))
//
//		queryParams := map[string]interface{}{}
//
//		queryParams["uid"] = uid
//		queryParams["orgId"] = orgId
//
//		fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
//		respBody, respStatusCode, err := http.Get(reqUrl, queryParams)
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//
//	}
//
//}
//
//func GetDeviceFiltrateHandler() gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//		uid := c.Query("uid")
//		orgId := c.Query("orgId")
//
//		//权限校验
//		if uid != strconv.FormatInt(cacheUserInfo.UserId, 10) || orgId != strconv.FormatInt(cacheUserInfo.OrgId, 10) {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//
//		queryParams := map[string]interface{}{}
//
//		queryParams["uid"] = uid
//		queryParams["orgId"] = orgId
//
//		reqUrl := fmt.Sprintf("%s/api/mct/devices/getDeviceFiltrate", config.GetPreUrl("compatsvc"))
//
//		fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
//
//		respVo := &compatsvc.GetDeviceFiltrateRespVo{}
//
//		respBody, respStatusCode, err := http.Get(reqUrl, queryParams)
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//}
//
//func GetAllDevicesStatusHandler() gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//		uid := c.Query("uid")
//		orgId := c.Query("orgId")
//
//		//权限校验
//		if uid != strconv.FormatInt(cacheUserInfo.UserId, 10) || orgId != strconv.FormatInt(cacheUserInfo.OrgId, 10) {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//
//		queryParams := map[string]interface{}{}
//
//		queryParams["uid"] = uid
//		queryParams["orgId"] = orgId
//
//		reqUrl := fmt.Sprintf("%s/api/mct/devices/getAllDevicesStatus", config.GetPreUrl("compatsvc"))
//
//		fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
//
//		respVo := &compatsvc.GetAllDevicesStatusListRespVo{}
//
//		respBody, respStatusCode, err := http.Get(reqUrl, queryParams)
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//
//}
//
//func StartCompatHandler() gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//
//		reqUrl := fmt.Sprintf("%s/api/mct/devices/startCompat", config.GetPreUrl("compatsvc"))
//		//组装参数
//		reqVo := &compatsvc.StartCompatReqVo{}
//
//		if c.BindJSON(reqVo) != nil {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.JSONConvertError))
//			return
//		}
//
//		//权限校验
//		if reqVo.OrgId != cacheUserInfo.OrgId || reqVo.Uid != cacheUserInfo.UserId {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//		//请求响应
//		respVo := &compatsvc.StartCompatRespVo{}
//
//		respBody, respStatusCode, err := http.Post(reqUrl, nil, json.ToJsonIgnoreError(reqVo))
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//}
//
////报告
//func GetReportHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//
//		reqUrl := fmt.Sprintf("%s/api/mct/report/getReport", config.GetPreUrl("compatsvc"))
//		//组装参数
//		reqVo := &compatsvc.GetReportReqVo{}
//
//		if c.BindJSON(reqVo) != nil {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.JSONConvertError))
//			return
//		}
//
//		//权限校验
//		if reqVo.OrgId != cacheUserInfo.OrgId || reqVo.Uid != cacheUserInfo.UserId {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//		//请求响应
//		respVo := &compatsvc.GetReportListRespVo{}
//
//		respBody, respStatusCode, err := http.Post(reqUrl, nil, json.ToJsonIgnoreError(reqVo))
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//}
//
//func DeleteReportHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//
//		reqUrl := fmt.Sprintf("%s/api/mct/report/deleteReport", config.GetPreUrl("compatsvc"))
//		//组装参数
//		reqVo := &compatsvc.DeleteReportReqVo{}
//
//		if c.BindJSON(reqVo) != nil {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.JSONConvertError))
//			return
//		}
//
//		//权限校验
//		if reqVo.OrgId != cacheUserInfo.OrgId || reqVo.Uid != cacheUserInfo.UserId {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//		//请求响应
//		respVo := &compatsvc.DeleteReportRespVo{}
//
//		respBody, respStatusCode, err := http.Post(reqUrl, nil, json.ToJsonIgnoreError(reqVo))
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//}
//
//func GetReportApkInfoHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//
//		reqUrl := fmt.Sprintf("%s/api/mct/report/getReportApkInfo", config.GetPreUrl("compatsvc"))
//		//组装参数
//		reqVo := &compatsvc.GetReportApkInfoReqVo{}
//
//		if c.BindJSON(reqVo) != nil {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.JSONConvertError))
//			return
//		}
//
//		//权限校验
//		if reqVo.OrgId != cacheUserInfo.OrgId || reqVo.Uid != cacheUserInfo.UserId {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//		//请求响应
//		respVo := &compatsvc.GetReportApkInfoRespVo{}
//
//		respBody, respStatusCode, err := http.Post(reqUrl, nil, json.ToJsonIgnoreError(reqVo))
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//}
//
//func GetReportDetailOverViewHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//
//		reqUrl := fmt.Sprintf("%s/api/mct/report/getReportDetailOverView", config.GetPreUrl("compatsvc"))
//		//组装参数
//		reqVo := &compatsvc.ReportDetailOverViewReqVo{}
//
//		if c.BindJSON(reqVo) != nil {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.JSONConvertError))
//			return
//		}
//
//		//权限校验
//		if reqVo.OrgId != cacheUserInfo.OrgId || reqVo.Uid != cacheUserInfo.UserId {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//		//请求响应
//		respVo := &compatsvc.ReportDetailOverViewRespListVo{}
//
//		respBody, respStatusCode, err := http.Post(reqUrl, nil, json.ToJsonIgnoreError(reqVo))
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//}
//
//func GetReportDetailSingleHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//
//		reqUrl := fmt.Sprintf("%s/api/mct/report/getReportDetailSingle", config.GetPreUrl("compatsvc"))
//		//组装参数
//		reqVo := &compatsvc.GetReportDetailSingleReqVo{}
//
//		if c.BindJSON(reqVo) != nil {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.JSONConvertError))
//			return
//		}
//
//		//权限校验
//		if reqVo.OrgId != cacheUserInfo.OrgId || reqVo.Uid != cacheUserInfo.UserId {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//		//请求响应
//		respVo := &compatsvc.GetReportDetailSingleRespVo{}
//
//		respBody, respStatusCode, err := http.Post(reqUrl, nil, json.ToJsonIgnoreError(reqVo))
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//}
//
//func GetReportDetailErrorHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//
//		reqUrl := fmt.Sprintf("%s/api/mct/report/getReportDetailError", config.GetPreUrl("compatsvc"))
//		//组装参数
//		reqVo := &compatsvc.GetReportDetailErrorReqVo{}
//
//		if c.BindJSON(reqVo) != nil {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.JSONConvertError))
//			return
//		}
//
//		//权限校验
//		if reqVo.OrgId != cacheUserInfo.OrgId || reqVo.Uid != cacheUserInfo.UserId {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//		//请求响应
//		respVo := &compatsvc.GetReportDetailErrorListRespVo{}
//
//		respBody, respStatusCode, err := http.Post(reqUrl, nil, json.ToJsonIgnoreError(reqVo))
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//}
//
//func GetReportDetailPerformanceHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		r := c.Request
//
//		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
//
//		if userErr != nil {
//			log.Error(userErr)
//			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
//			return
//		}
//
//		reqUrl := fmt.Sprintf("%s/api/mct/report/getReportDetailPerformance", config.GetPreUrl("compatsvc"))
//		//组装参数
//		reqVo := &compatsvc.GetReportDetailPerformanceReqVo{}
//
//		if c.BindJSON(reqVo) != nil {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.JSONConvertError))
//			return
//		}
//
//		//权限校验
//		if reqVo.OrgId != cacheUserInfo.OrgId || reqVo.Uid != cacheUserInfo.UserId {
//			c.JSON(200, errs.BuildSystemErrorInfo(errs.TokenAuthError))
//			return
//		}
//		//请求响应
//		respVo := &compatsvc.GetReportDetailPerformanceListRespVo{}
//
//		respBody, respStatusCode, err := http.Post(reqUrl, nil, json.ToJsonIgnoreError(reqVo))
//
//		//Process the response
//		if err != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
//			log.Errorf("request [%s] failed, response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		//接口响应错误
//		if respStatusCode < 200 || respStatusCode > 299 {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
//			log.Errorf("request [%s] failed , response status code [%d], err [%v]", reqUrl, respStatusCode, err)
//			c.JSON(500, respVo)
//			return
//		}
//		jsonConvertErr := json.FromJson(respBody, respVo)
//		if jsonConvertErr != nil {
//			respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
//		}
//		c.JSON(200, respVo)
//	}
//}
