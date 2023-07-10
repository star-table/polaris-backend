package api

import (
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) SendSMSLoginCode(req orgvo.SendSMSLoginCodeReqVo) vo.VoidErr {
	phoneNumber := req.Input.PhoneNumber
	//phone format check
	verifyErr := service.VerifyCaptcha(req.Input.CaptchaID, req.Input.CaptchaPassword, req.Input.PhoneNumber, req.Input.YidunValidate)
	if verifyErr != nil {
		return vo.VoidErr{Err: vo.NewErr(verifyErr)}
	}

	err := service.SendSMSLoginCode(phoneNumber)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) SendAuthCode(req orgvo.SendAuthCodeReqVo) vo.VoidErr {
	if ok, _ := slice.Contain([]int{consts.AuthCodeTypeBind, consts.AuthCodeTypeUnBind}, req.Input.AuthType); !ok {
		verifyErr := service.VerifyCaptcha(req.Input.CaptchaID, req.Input.CaptchaPassword, req.Input.Address, req.Input.YidunValidate)
		if verifyErr != nil {
			return vo.VoidErr{Err: vo.NewErr(verifyErr)}
		}
	}

	err := service.SendAuthCode(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) GetPwdLoginCode(req orgvo.GetPwdLoginCodeReqVo) orgvo.GetPwdLoginCodeRespVo {
	res, err := domain.GetPwdLoginCode(req.CaptchaId)
	return orgvo.GetPwdLoginCodeRespVo{Err: vo.NewErr(err), CaptchaPassword: res}
}

func (PostGreeter) SetPwdLoginCode(req orgvo.SetPwdLoginCodeReqVo) vo.VoidErr {
	return vo.VoidErr{Err: vo.NewErr(domain.SetPwdLoginCode(req.CaptchaId, req.CaptchaPassword))}
}
