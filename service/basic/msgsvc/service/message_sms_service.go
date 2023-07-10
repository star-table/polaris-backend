package service

import (
	"github.com/star-table/common/core/mail"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/sdk/sms"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/vo/msgvo"
	"strings"
)

func SendLoginSMS(req msgvo.SendLoginSMSReqVo) errs.SystemErrorInfo {
	code := req.Code
	phoneNumber := req.PhoneNumber

	resp, err := sms.SendSMS(phoneNumber, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeLoginAuthCode, map[string]string{
		consts.SMSParamsNameCode: code,
	})
	if err != nil {
		log.Infof("登录验证码发送失败， 手机号 %s, 错误信息: %v", phoneNumber, err)
		return errs.BuildSystemErrorInfo(errs.SMSLoginCodeSendError, err)
	}

	if !strings.EqualFold(resp.Code, "OK") {
		log.Infof("登录验证码发送失败， 手机号 %s, 错误信息: %v", phoneNumber, resp.Message)

		if strings.EqualFold(resp.Code, "isv.MOBILE_NUMBER_ILLEGAL") {
			return errs.BuildSystemErrorInfo(errs.SMSPhoneNumberFormatError)
		}
		if strings.EqualFold(resp.Code, "isv.BUSINESS_LIMIT_CONTROL") {
			return errs.BuildSystemErrorInfo(errs.SMSSendLimitError)
		}
		return errs.BuildSystemErrorInfo(errs.SMSLoginCodeSendError)
	}

	log.Infof("登录验证码发送成功， 手机号 %s, code %s", phoneNumber, code)
	return nil
}

func SendSMS(req msgvo.SendSMSReqVo) errs.SystemErrorInfo{
	input := req.Input
	phoneNumber := input.Mobile

	resp, err := sms.SendSMS(phoneNumber, input.SignName, input.TemplateCode, input.Params)
	if err != nil {
		log.Infof("登录验证码发送失败， 手机号 %s, 错误信息: %v", phoneNumber, err)
		return errs.BuildSystemErrorInfo(errs.SMSLoginCodeSendError, err)
	}

	if !strings.EqualFold(resp.Code, "OK") {
		log.Infof("登录验证码发送失败， 手机号 %s, 错误信息: %v", phoneNumber, resp.Message)

		if strings.EqualFold(resp.Code, "isv.MOBILE_NUMBER_ILLEGAL") {
			return errs.BuildSystemErrorInfo(errs.SMSPhoneNumberFormatError)
		}
		if strings.EqualFold(resp.Code, "isv.BUSINESS_LIMIT_CONTROL") {
			return errs.BuildSystemErrorInfo(errs.SMSSendLimitError)
		}
		return errs.BuildSystemErrorInfo(errs.SMSLoginCodeSendError)
	}

	log.Infof("登录验证码发送成功， 手机号 %s, params %s", phoneNumber, json.ToJsonIgnoreError(input.Params))
	return nil
}

func SendMail(req msgvo.SendMailReqVo) errs.SystemErrorInfo{
	input := req.Input
	emails := input.Emails
	subject := strings.TrimSpace(input.Subject)
	content := input.Content

	verifySuc := format.VerifyEmailFormat(emails...)
	if ! verifySuc{
		log.Errorf("邮箱格式验证未通过 %s", json.ToJsonIgnoreError(emails))
		return errs.EmailFormatErr
	}

	if strs.Len(subject) == 0{
		return errs.EmailSubjectEmptyErr
	}

	err := mail.SendMail(emails, subject, content)
	if err != nil{
		log.Error(err)
		return errs.EmailSendErr
	}
	return nil
}
