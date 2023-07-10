package service

import (
	"strings"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/core/yidun"

	"github.com/star-table/common/core/util/temp"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/rand"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/msgfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
)

const defaultAuthCode = "000000"

func SendSMSLoginCode(phoneNumber string) errs.SystemErrorInfo {
	return SendAuthCode(orgvo.SendAuthCodeReqVo{
		Input: vo.SendAuthCodeReq{
			AuthType:    consts.AuthCodeTypeLogin,
			AddressType: consts.ContactAddressTypeMobile,
			Address:     phoneNumber,
		},
	})
}

func SendAuthCode(req orgvo.SendAuthCodeReqVo) errs.SystemErrorInfo {
	input := req.Input

	addressType := input.AddressType
	authType := input.AuthType
	contactAddress := input.Address

	if addressType != consts.ContactAddressTypeMobile && addressType != consts.ContactAddressTypeEmail {
		return errs.NotSupportedContactAddressType
	}

	limitErr := domain.CheckSMSLoginCodeFreezeTime(authType, addressType, contactAddress)
	if limitErr != nil {
		log.Error(limitErr)
		return limitErr
	}

	err := domain.CheckLoginNameIsExist(addressType, contactAddress)
	//如果不是注册，登录，绑定，该账户必须存在
	if authType != consts.AuthCodeTypeRegister && authType != consts.AuthCodeTypeLogin && authType != consts.AuthCodeTypeBind {
		if err != nil {
			log.Error(err)
			if err.Code() == errs.UserNotExist.Code() {
				return errs.NotBindAccountError
			} else {
				return err
			}
		}
	}

	//如果是注册，该账户必须不存在，绑定已经改为随意绑定
	if authType == consts.AuthCodeTypeRegister {
		if err == nil {
			return errs.AccountAlreadyBindError
		}
	}

	authCode := defaultAuthCode
	loginName := contactAddress
	infos := strings.Split(loginName, "-")
	if len(infos) > 1 {
		loginName = infos[1]
	}
	if !IsInWhiteList(loginName) {
		authCode = rand.RandomVerifyCode(6)

		//异步发送
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
				}
			}()
			switch addressType {
			case consts.ContactAddressTypeMobile:
				sendErr := sendSmsAuthCode(authType, loginName, authCode)
				if sendErr != nil {
					log.Error(sendErr)
				}
			case consts.ContactAddressTypeEmail:
				sendErr := sendMailAuthCode(authType, loginName, authCode)
				if sendErr != nil {
					log.Error(sendErr)
				}
			}
		}()
	}
	setFreezeErr := domain.SetSMSLoginCodeFreezeTime(authType, addressType, contactAddress, consts.SMSLoginCodeFreezeTime)
	if setFreezeErr != nil {
		//这里不要影响主流程
		log.Error(setFreezeErr)
	}
	setLoginCode := domain.SetSMSLoginCode(authType, addressType, contactAddress, authCode)
	if setLoginCode != nil {
		//这里不要影响主流程
		log.Error(setLoginCode)
	}
	return nil
}

func sendSmsAuthCode(authType int, mobile string, authCode string) errs.SystemErrorInfo {
	infos := strings.Split(mobile, "-")
	if len(infos) > 1 {
		mobile = infos[1]
	}
	switch authType {
	case consts.AuthCodeTypeLogin:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeLoginAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeRegister:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeRegisterAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeResetPwd:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeResetPwdAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeRetrievePwd:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeRetrievePwdAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeBind:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeBindAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeUnBind:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeUnBindAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	}
	return errs.NotSupportedAuthCodeType
}

func sendSmsByTplWithParam(authType int, mobile string, tplParam map[string]string) errs.SystemErrorInfo {
	switch authType {
	case consts.AuthCodeTypeLogin:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeLoginAuthCode, tplParam)
	case consts.AuthCodeTypeRegister:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeRegisterAuthCode, tplParam)
	case consts.AuthCodeTypeResetPwd:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeResetPwdAuthCode, tplParam)
	case consts.AuthCodeTypeRetrievePwd:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeRetrievePwdAuthCode, tplParam)
	case consts.AuthCodeTypeBind:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeBindAuthCode, tplParam)
	case consts.AuthCodeTypeUnBind:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeUnBindAuthCode, tplParam)
	case consts.AuthCodeTypeInviteUserByPhones: // 待申请模板，批量邀请成员加入团队
		// 参数 authCode 在这里其实是邀请成员的链接
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeInviteUserByPhones, tplParam)
	}

	return errs.NotSupportedAuthCodeType
}

func sendMailAuthCode(authType int, email string, authCode string) errs.SystemErrorInfo {
	emails := []string{email}

	switch authType {
	case consts.AuthCodeTypeLogin:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeLogin, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionLogin,
		}))
	case consts.AuthCodeTypeRegister:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeRegister, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionRegister,
		}))
	case consts.AuthCodeTypeResetPwd:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeResetPwd, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionResetPwd,
		}))
	case consts.AuthCodeTypeRetrievePwd:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeRetrievePwd, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionRetrievePwd,
		}))
	case consts.AuthCodeTypeBind:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeBind, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionBind,
		}))
	case consts.AuthCodeTypeUnBind:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeUnBind, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionUnBind,
		}))
	}
	return errs.NotSupportedAuthCodeType
}

func IsInWhiteList(phoneNumber string) bool {
	whiteList, err := domain.GetPhoneNumberWhiteList()
	if err != nil {
		log.Error(err)
		return false
	}
	for _, v := range whiteList {
		if v == phoneNumber {
			return true
		}
	}

	return false
}

func VerifyCaptcha(captchaID, captchaPassword *string, phoneNumber string, yiDunValidate *string) errs.SystemErrorInfo {
	if IsInWhiteList(phoneNumber) {
		return nil
	}
	if yiDunValidate != nil {
		if strings.TrimSpace(*yiDunValidate) == "" {
			return errs.CaptchaError
		}
		verifyResult, err := yidun.Verify(strings.TrimSpace(*yiDunValidate), "u")
		if err != nil {
			log.Error(err)
			return errs.CaptchaError
		}
		if !verifyResult.Result {
			log.Errorf("yidun verify err: %s", json.ToJsonIgnoreError(verifyResult))
			return errs.CaptchaError
		}
	} else {
		if captchaID == nil || captchaPassword == nil {
			return errs.CaptchaError
		}

		res, err := domain.GetPwdLoginCode(*captchaID)
		if err != nil {
			log.Error(err)
			return err
		}

		clearErr := domain.ClearPwdLoginCode(*captchaID)
		if clearErr != nil {
			log.Error(clearErr)
			return clearErr
		}

		if res != *captchaPassword {
			return errs.CaptchaError
		}
	}

	return nil
}
