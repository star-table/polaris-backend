package api

import (
	"context"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/extra/gin/util"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

func (r *mutationResolver) ResetPassword(ctx context.Context, input vo.ResetPasswordReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := orgfacade.ResetPassword(orgvo.ResetPasswordReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return &vo.Void{ID: cacheUserInfo.UserId}, nil
}

func (r *mutationResolver) SetPassword(ctx context.Context, input vo.SetPasswordReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := orgfacade.SetPassword(orgvo.SetPasswordReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return &vo.Void{ID: cacheUserInfo.UserId}, nil
}

func (r *mutationResolver) SendSmsLoginCode(ctx context.Context, input vo.SendSmsLoginCodeReq) (*vo.Void, error) {
	respVo := orgfacade.SendSMSLoginCode(orgvo.SendSMSLoginCodeReqVo{
		Input: input,
	})
	return &vo.Void{ID: 0}, respVo.Error()
}

func (r *mutationResolver) SendAuthCode(ctx context.Context, input vo.SendAuthCodeReq) (*vo.Void, error) {
	respVo := orgfacade.SendAuthCode(orgvo.SendAuthCodeReqVo{
		Input: input,
	})
	return &vo.Void{ID: 0}, respVo.Error()
}

func (r *mutationResolver) UserQuit(ctx context.Context) (*vo.Void, error) {
	token, err := util.GetCtxToken(ctx)
	if err != nil {
		//token不存在，直接退出
		return &vo.Void{ID: 1}, nil
	}
	//异步退出，保证接口响应低延迟
	go func() {
		_ = orgfacade.UserQuit(orgvo.UserQuitReqVo{
			Token: token,
		})
	}()
	return &vo.Void{ID: 1}, nil
}

func (r *mutationResolver) UserLogin(ctx context.Context, input vo.UserLoginReq) (*vo.UserLoginResp, error) {
	respVo := orgfacade.UserLogin(orgvo.UserLoginReqVo{
		UserLoginReq: input,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.Data, nil
}

func (r *mutationResolver) RetrievePassword(ctx context.Context, input vo.RetrievePasswordReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := orgfacade.RetrievePassword(orgvo.RetrievePasswordReqVo{
		OrgId: cacheUserInfo.OrgId,
		Input: input,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return &vo.Void{ID: 0}, nil
}
