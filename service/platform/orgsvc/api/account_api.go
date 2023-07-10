package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) SetPassword(req orgvo.SetPasswordReqVo) vo.VoidErr {
	err := service.SetPassword(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) ResetPassword(req orgvo.ResetPasswordReqVo) vo.VoidErr {
	err := service.ResetPassword(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) RetrievePassword(req orgvo.RetrievePasswordReqVo) vo.VoidErr {
	err := service.RetrievePassword(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

// UnbindLoginName 目前情况先不允许解绑，要不会导致本地账户丢失，再也无法使用
func (PostGreeter) UnbindLoginName(req orgvo.UnbindLoginNameReqVo) vo.VoidErr {
	err := service.UnbindLoginName(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) BindLoginName(req orgvo.BindLoginNameReqVo) vo.VoidErr {
	err := service.BindLoginName(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) CheckLoginName(req orgvo.CheckLoginNameReqVo) vo.VoidErr {
	err := service.CheckLoginName(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) VerifyOldName(req orgvo.UnbindLoginNameReqVo) vo.VoidErr {
	err := service.VerifyOldName(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) ChangeLoginName(req orgvo.BindLoginNameReqVo) vo.VoidErr {
	err := service.ChangeLoginName(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) DisbandThirdAccount(req vo.CommonReqVo) vo.VoidErr {
	err := service.DisbandThirdAccount(req.OrgId, req.UserId, req.SourceChannel)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) ThirdAccountBindList(req vo.CommonReqVo) orgvo.ThirdAccountListResp {
	res, err := service.ThirdAccountBindList(req.OrgId, req.UserId)
	return orgvo.ThirdAccountListResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
