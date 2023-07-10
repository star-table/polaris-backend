package api

import (
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/extra/gin/mvc"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/service"
)

var log = logger.GetDefaultLogger()

type PostGreeter struct {
	mvc.Greeter
}

type GetGreeter struct {
	mvc.Greeter
}

func (GetGreeter) Health() string {
	return "ok"
}

// 发送日志到钉钉群里
func (PostGreeter) DingTalkInfo(req commonvo.DingTalkInfoReqVo) vo.BoolRespVo {
	logData := req.LogData
	reqParam := vo.DingTalkInfoReq{
		Content: json.ToJsonIgnoreError(logData),
		Other:   "",
	}
	isOk := false
	err := service.DingTalkInfo(reqParam)
	if err == nil {
		isOk = true
	}
	return vo.BoolRespVo{
		Err:    vo.NewErr(err),
		IsTrue: isOk,
	}
}
