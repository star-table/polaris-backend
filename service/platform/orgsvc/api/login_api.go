package api

import (
	"strings"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) UserLogin(req orgvo.UserLoginReqVo) orgvo.UserSMSLoginRespVo {
	if req.UserLoginReq.Name != nil && *req.UserLoginReq.Name != "" {
		creatorName := strings.TrimSpace(*req.UserLoginReq.Name)
		//creatorNameLen := str.CountStrByGBK(creatorName)
		//if creatorNameLen == 0 || creatorNameLen > 20{
		//	log.Error("姓名长度错误")
		//	return orgvo.UserSMSLoginRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.UserNameLenError)), NewData: nil}
		//}
		isNameRight := format.VerifyUserNameFormat(creatorName)
		if !isNameRight {
			return orgvo.UserSMSLoginRespVo{Err: vo.NewErr(errs.UserNameLenError), Data: nil}
		}
		req.UserLoginReq.Name = &creatorName
	}

	res, err := service.UserLogin(req.UserLoginReq)
	return orgvo.UserSMSLoginRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) UserQuit(req orgvo.UserQuitReqVo) vo.VoidErr {
	err := service.UserQuit(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) JoinOrgByInviteCode(req orgvo.JoinOrgByInviteCodeReq) vo.VoidErr {
	err := service.JoinOrgByInviteCode(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) ExchangeShareToken(req orgvo.ExchangeShareTokenReq) orgvo.UserSMSLoginRespVo {
	data, err := service.ExchangeShareToken(req)
	return orgvo.UserSMSLoginRespVo{
		Err:  vo.NewErr(err),
		Data: data,
	}
}
