package msgfacade

import (
	"github.com/star-table/common/core/model"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/msgvo"
)

func PushMsgToMqRelaxed(msg model.MqMessage, msgType int, orgId int64) errs.SystemErrorInfo{
	respVo := PushMsgToMq(msgvo.PushMsgToMqReqVo{Msg: msg, MsgType: msgType, OrgId: orgId})
	if respVo.Failure() {
		return respVo.Error()
	}
	return nil
}

func InsertMqConsumeFailMsgRelaxed(msg model.MqMessage, msgType int, orgId int64) errs.SystemErrorInfo{
	respVo := InsertMqConsumeFailMsg(msgvo.InsertMqConsumeFailMsgReqVo{Msg: msg, MsgType: msgType, OrgId: orgId})
	if respVo.Failure() {
		return respVo.Error()
	}
	return nil
}


func SendSMSRelaxed(mobile string, signName string, templateCode string, params map[string]string) errs.SystemErrorInfo {
	respVo := SendSMS(msgvo.SendSMSReqVo{
		Input: msgvo.SendSMSReqVoReqData{
			Mobile: mobile,
			SignName: signName,
			TemplateCode: templateCode,
			Params: params,
		},
	})
	if respVo.Failure() {
		return respVo.Error()
	}
	return nil
}

func SendMailRelaxed(emails []string, subject string, content string) errs.SystemErrorInfo{
	respVo := SendMail(msgvo.SendMailReqVo{
		Input: msgvo.SendMailReqData{
			Emails: emails,
			Subject: subject,
			Content: content,
		},
	})
	if respVo.Failure() {
		return respVo.Error()
	}
	return nil
}