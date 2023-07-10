package msgvo

import (
	"github.com/star-table/common/core/model"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type PushMsgToMqReqVo struct {
	Msg     model.MqMessage `json:"msg"`
	MsgType int             `json:"msgType"`
	OrgId   int64           `json:"orgId"`
}

type InsertMqConsumeFailMsgReqVo struct {
	Msg     model.MqMessage `json:"msg"`
	MsgType int             `json:"msgType"`
	OrgId   int64           `json:"orgId"`
}

type GetFailMsgListReqVo struct {
	OrgId   *int64 `json:"orgId"`
	MsgType *int   `json:"msgType"`
	Page    *int   `json:"page"`
	Size    *int   `json:"size"`
}

type GetFailMsgListRespVo struct {
	vo.Err
	Data GetFailMsgListRespData `json:"data"`
}

type GetFailMsgListRespData struct {
	MsgList *[]bo.MessageBo `json:"msgList"`
}

type UpdateMsgStatusReqVo struct {
	MsgId     int64 `json:"msgId"`
	NewStatus int   `json:"newStatus"`
}

type WriteSomeFailedMsgReqVo struct {
	OrgId int64                        `json:"orgId"`
	Input *WriteSomeFailedMsgReqVoData `json:"input"`
}

type WriteSomeFailedMsgReqVoData struct {
	MsgType int    `json:"msgType"` // 参考 `consts.PushTypeFsCallbackHandleFailed`
	MsgBody string `json:"msgBody"`
}
