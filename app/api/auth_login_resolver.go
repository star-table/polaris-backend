package api

import (
	"context"

	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

func (r *mutationResolver) UnbindLoginName(ctx context.Context, input vo.UnbindLoginNameReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := orgfacade.UnbindLoginName(orgvo.UnbindLoginNameReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return &vo.Void{ID: cacheUserInfo.UserId}, nil
}

func (r *mutationResolver) BindLoginName(ctx context.Context, input vo.BindLoginNameReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := orgfacade.BindLoginName(orgvo.BindLoginNameReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return &vo.Void{ID: cacheUserInfo.UserId}, nil
}

func (r *mutationResolver) CheckLoginName(ctx context.Context, input vo.CheckLoginNameReq) (*vo.Void, error) {
	respVo := orgfacade.CheckLoginName(orgvo.CheckLoginNameReqVo{
		Input: input,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return &vo.Void{ID: 0}, nil
}

func (r *mutationResolver) VerifyOldName(ctx context.Context, input vo.UnbindLoginNameReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := orgfacade.VerifyOldName(orgvo.UnbindLoginNameReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return &vo.Void{ID: cacheUserInfo.UserId}, nil
}

func (r *mutationResolver) ChangeLoginName(ctx context.Context, input vo.BindLoginNameReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := orgfacade.ChangeLoginName(orgvo.BindLoginNameReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return &vo.Void{ID: cacheUserInfo.UserId}, nil
}

func (r *mutationResolver) AuthFsCode(ctx context.Context, input vo.FeiShuAuthReq) (*vo.FeiShuAuthCodeResp, error) {
	respVo := orgfacade.FeiShuAuthCode(input)
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.Auth, nil
}

func (r *mutationResolver) AuthFs(ctx context.Context, input vo.FeiShuAuthReq) (*vo.FeiShuAuthResp, error) {
	respVo := orgfacade.FeiShuAuth(input)
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.Auth, nil
}

func (r *queryResolver) CheckTokenValidity(ctx context.Context) (*vo.CheckTokenValidityResp, error) {
	userInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return &vo.CheckTokenValidityResp{
			ID: 401,
		}, nil
	}
	return &vo.CheckTokenValidityResp{
		OrgID: userInfo.OrgId,
	}, nil
}
