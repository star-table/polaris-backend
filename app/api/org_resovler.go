package api

import (
	"context"

	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/extra/gin/util"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

func (r *mutationResolver) CreateOrg(ctx context.Context, input vo.CreateOrgReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserWithoutOrgVerifyRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	token, err1 := util.GetCtxToken(ctx)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}

	importSampleData := 0
	if input.ImportSampleData != nil {
		importSampleData = *input.ImportSampleData
	}
	respVo := orgfacade.CreateOrg(orgvo.CreateOrgReqVo{Data: orgvo.CreateOrgReqVoData{
		CreatorId:        cacheUserInfo.UserId,
		CreateOrgReq:     input,
		UserToken:        token,
		ImportSampleData: importSampleData,
	},
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return &vo.Void{ID: respVo.Data.OrgId}, nil

}

//切换用户组织
func (r *mutationResolver) SwitchUserOrganization(ctx context.Context, input vo.SwitchUserOrganizationReq) (*vo.Void, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserWithoutOrgVerifyRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	token, err1 := util.GetCtxToken(ctx)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}

	if cacheUserInfo.OrgId == input.OrgID {
		log.Infof("不需要切换组织")
		return &vo.Void{ID: input.OrgID}, nil
	}

	respVo := orgfacade.SwitchUserOrganization(orgvo.SwitchUserOrganizationReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  input.OrgID,
		Token:  token,
	})

	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return &vo.Void{ID: respVo.OrgId}, nil

}

//组织设置
func (r *mutationResolver) UpdateOrganizationSetting(ctx context.Context, input vo.UpdateOrganizationSettingsReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := orgfacade.UpdateOrganizationSetting(orgvo.UpdateOrganizationSettingReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
	})

	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}

	return &vo.Void{ID: respVo.OrgId}, nil
}

func (r *queryResolver) OrganizationInfo(ctx context.Context, input vo.OrganizationInfoReq) (*vo.OrganizationInfoResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
		OrgId:  input.OrgID,
		UserId: cacheUserInfo.UserId,
	})

	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}

	return respVo.OrganizationInfo, nil

}

//用户组织列表
func (r *queryResolver) UserOrganizationList(ctx context.Context) (*vo.UserOrganizationListResp, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserWithoutOrgVerifyRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := orgfacade.UserOrganizationList(orgvo.UserOrganizationListReqVo{
		UserId: cacheUserInfo.UserId,
	})

	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.UserOrganizationListResp, nil
}

func (r *queryResolver) GetOrgConfig(ctx context.Context) (*vo.OrgConfig, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserWithoutOrgVerifyRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := orgfacade.GetOrgConfig(orgvo.GetOrgConfigReq{OrgId: cacheUserInfo.OrgId})
	if resp.Failure() {
		log.Error(resp.Message)
		return nil, resp.Error()
	}

	return resp.Data, nil
}

func (r *queryResolver) GetAppTicket(ctx context.Context) (*vo.GetAppTicketResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := orgfacade.GetAppTicket(orgvo.GetAppTicketReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}

	return respVo.Data, nil
}

func (r *mutationResolver) InitFsAccount(ctx context.Context, input vo.InitFeiShuAccountReq) (*vo.FeiShuAuthCodeResp, error) {
	respVo := orgfacade.InitFeiShuAccount(orgvo.InitFsAccountReqVo{
		Input: input,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.Auth, nil
}

func (r *queryResolver) GetJsAPITicket(ctx context.Context) (*vo.GetJsAPITicketResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := orgfacade.GetJsAPITicket(vo.CommonReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Data, nil
}

func (r *mutationResolver) BoundFs(ctx context.Context, input vo.BoundFeiShuReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	token, err1 := util.GetCtxToken(ctx)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	respVo := orgfacade.BoundFeiShu(orgvo.BoundFsReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Token:  token,
		Input:  input,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.Void, nil
}

func (r *mutationResolver) BoundFsAccount(ctx context.Context, input vo.BoundFeiShuAccountReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	token, err1 := util.GetCtxToken(ctx)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	respVo := orgfacade.BoundFeiShuAccount(orgvo.BoundFsAccountReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Token:  token,
		Input:  input,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.Void, nil
}
