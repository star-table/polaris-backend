package pushfacade

import (
	"fmt"

	"github.com/star-table/polaris-backend/facade"

	pushPb "github.com/star-table/interface/golang/push/v1"
	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
)

func GenerateCard(req *pushPb.GenerateCardReq) *commonvo.GenerateCardResp {
	respVo := &commonvo.GenerateCardResp{Data: &pushPb.GenerateCardReply{}}
	reqUrl := fmt.Sprintf("%s%s/generate/card", config.GetPreUrl(consts.ServicePush), ApiV1Prefix)
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func PushMqtt(req *pushPb.PushMqttReq) *commonvo.PushMqttResp {
	respVo := &commonvo.PushMqttResp{Data: &pushPb.PushMqttReply{}}
	reqUrl := fmt.Sprintf("%s%s/push/mqtt", config.GetPreUrl(consts.ServicePush), ApiV1Prefix)
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GenerateMqttKey(req *pushPb.GenerateMqttKeyReq) *commonvo.GenerateMqttKeyResp {
	respVo := &commonvo.GenerateMqttKeyResp{Data: &pushPb.GenerateMqttKeyReply{}}
	reqUrl := fmt.Sprintf("%s%s/generate/mqtt/key", config.GetPreUrl(consts.ServicePush), ApiV1Prefix)
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}
