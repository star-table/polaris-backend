package msgvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type GetMQTTChannelKeyReqVo struct {
	UserId int64                   `json:"userId"`
	OrgId  int64                   `json:"orgId"`
	Input  vo.GetMQTTChannelKeyReq `json:"input"`
}

type GetMQTTChannelKeyRespVo struct {
	vo.Err
	Data *vo.GetMQTTChannelKeyResp `json:"data"`
}

type PushMsgToMqttReqVo struct {
	Msg bo.MQTTDataRefreshNotice `json:"msg"`
}

type PushRemindMsgToMqttReqVo struct {
	UserId int64               `json:"userId"`
	OrgId  int64               `json:"orgId"`
	Msg    bo.MQTTRemindNotice `json:"msg"`
}

type MqttPublishData struct {
	Channel string `json:"channel"`
	Body    []byte `json:"body"`
}

type MqttPublishReqVo struct {
	Data *MqttPublishData `json:data`
}
