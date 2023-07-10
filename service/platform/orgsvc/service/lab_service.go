package service

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
)

func SetLabConfig(orgId, userId int64, input orgvo.SetLabReq) (*orgvo.SetLabResp, errs.SystemErrorInfo) {
	userAuthorityResp := userfacade.GetUserAuthority(orgId, userId)
	if userAuthorityResp.Failure() {
		log.Errorf("[SetLabConfig.GetUserAuthority] failed, orgId:%v, userId:%v, err:%v", orgId, userId, userAuthorityResp.Error())
		return nil, errs.BuildSystemErrorInfo(errs.UserDomainError, userAuthorityResp.Error())
	}
	isOrgOwner := userAuthorityResp.Data.IsOrgOwner
	isSysAdmin := userAuthorityResp.Data.IsSysAdmin
	isSubAdmin := userAuthorityResp.Data.IsSubAdmin
	if !isOrgOwner && !isSysAdmin && !isSubAdmin {
		return nil, errs.SetLabNoPermissionErr
	}

	err := domain.SetLabConfig(orgId, userId, input)
	if err != nil {
		log.Errorf("[SetLabConfig]错误, err: %v", err)
		return nil, err
	}
	return &orgvo.SetLabResp{OrgId: orgId}, nil
}

func GetLabConfig(input orgvo.GetLabReqVo) (*orgvo.GetLabResp, errs.SystemErrorInfo) {
	resp, err := domain.GetLabConfig(input.OrgId)
	if err != nil {
		log.Errorf("[GetLabConfig] err:%v", err)
		return nil, err
	}
	return &orgvo.GetLabResp{
		//WorkBenchShow:    resp.WorkBenchShow,
		ProOverview: resp.ProOverview,
		SideBarShow: resp.SideBarShow,
		EmptyApp:    resp.EmptyApp,
		//AutomationSwitch: resp.AutomationSwitch,
		//SubmitButton:  resp.SubmitButton,
	}, nil
}
