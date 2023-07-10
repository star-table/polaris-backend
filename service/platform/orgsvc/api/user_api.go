package api

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/ordervo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (GetGreeter) PersonalInfo(req orgvo.PersonalInfoReqVo) orgvo.PersonalInfoRespVo {
	res, err := service.PersonalInfo(req.OrgId, req.UserId, req.SourceChannel)
	respData := vo.PersonalInfo{}
	if err == nil {
		functionArr := res.Functions
		res.Functions = nil
		copyer.Copy(res, &respData)

		functionArrOfOrderPkg := make([]ordervo.FunctionLimitObj, 0, len(functionArr))
		copyer.Copy(functionArr, &functionArrOfOrderPkg)
		respData.Functions = domain.GetFunctionKeyListByFunctions(functionArrOfOrderPkg)
	}

	return orgvo.PersonalInfoRespVo{Err: vo.NewErr(err), PersonalInfo: &respData}
}

func (GetGreeter) PersonalInfoRest(req orgvo.PersonalInfoReqVo) orgvo.PersonalInfoRestRespVo {
	res, err := service.PersonalInfo(req.OrgId, req.UserId, req.SourceChannel)
	if err == nil {
		// functions 信息从 getOrgConfig 接口中获取
		res.Functions = make([]orgvo.FunctionLimitObj, 0)
	}

	return orgvo.PersonalInfoRestRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) GetUserIds(reqVo orgvo.GetUserIdsReqVo) orgvo.GetUserIdsRespVo {
	res, err := service.GetUserIds(reqVo.OrgId, reqVo.CorpId, reqVo.SourceChannel, reqVo.EmpIdsBody.EmpIds)
	return orgvo.GetUserIdsRespVo{Err: vo.NewErr(err), GetUserIds: res}
}

func (GetGreeter) GetUserId(reqVo orgvo.GetUserIdReqVo) orgvo.GetUserIdRespVo {
	res, err := service.GetUserId(reqVo.OrgId, reqVo.CorpId, reqVo.SourceChannel, reqVo.EmpId)
	return orgvo.GetUserIdRespVo{Err: vo.NewErr(err), GetUserId: res}
}

func (GetGreeter) UserConfigInfo(reqVo orgvo.UserConfigInfoReqVo) orgvo.UserConfigInfoRespVo {
	res, err := service.UserConfigInfo(reqVo.OrgId, reqVo.UserId)
	return orgvo.UserConfigInfoRespVo{Err: vo.NewErr(err), UserConfigInfo: res}
}

func (PostGreeter) UpdateUserConfig(req orgvo.UpdateUserConfigReqVo) orgvo.UpdateUserConfigRespVo {
	res, err := service.UpdateUserConfig(req.OrgId, req.UserId, req.UpdateUserConfigReq)
	return orgvo.UpdateUserConfigRespVo{Err: vo.NewErr(err), UpdateUserConfig: res}
}

func (PostGreeter) UpdateUserPcConfig(req orgvo.UpdateUserPcConfigReqVo) orgvo.UpdateUserConfigRespVo {
	res, err := service.UpdateUserPcConfig(req.OrgId, req.UserId, req.UpdateUserPcConfigReq)
	return orgvo.UpdateUserConfigRespVo{Err: vo.NewErr(err), UpdateUserConfig: res}
}

func (PostGreeter) UpdateUserDefaultProjectIdConfig(req orgvo.UpdateUserDefaultProjectIdConfigReqVo) orgvo.UpdateUserConfigRespVo {
	res, err := service.UpdateUserDefaultProjectIdConfig(req.OrgId, req.UserId, req.UpdateUserDefaultProjectIdConfigReq)
	return orgvo.UpdateUserConfigRespVo{Err: vo.NewErr(err), UpdateUserConfig: res}
}

func (GetGreeter) VerifyOrg(reqVo orgvo.VerifyOrgReqVo) vo.BoolRespVo {
	return vo.BoolRespVo{Err: vo.NewErr(nil), IsTrue: service.VerifyOrg(reqVo.OrgId, reqVo.UserId)}
}

func (PostGreeter) VerifyOrgUsers(reqVo orgvo.VerifyOrgUsersReqVo) vo.BoolRespVo {
	return vo.BoolRespVo{Err: vo.NewErr(nil), IsTrue: service.VerifyOrgUsers(reqVo.OrgId, reqVo.Input.UserIds)}
}

func (PostGreeter) VerifyOrgUsersReturnValid(reqVo orgvo.VerifyOrgUsersReqVo) vo.DataRespVo {
	return vo.DataRespVo{Err: vo.NewErr(nil), Data: service.VerifyOrgUsersReturnValid(reqVo.OrgId, reqVo.Input.UserIds)}
}

func (GetGreeter) GetUserInfo(reqVo orgvo.GetUserInfoReqVo) orgvo.GetUserInfoRespVo {
	res, err := service.GetUserInfo(reqVo.OrgId, reqVo.UserId, reqVo.SourceChannel)
	return orgvo.GetUserInfoRespVo{Err: vo.NewErr(err), UserInfo: res}
}

func (GetGreeter) GetOutUserInfoListBySourceChannel(req orgvo.GetOutUserInfoListBySourceChannelReqVo) orgvo.GetOutUserInfoListBySourceChannelRespVo {
	res, err := service.GetOutUserInfoListBySourceChannel(req.SourceChannel, req.Page, req.Size)
	return orgvo.GetOutUserInfoListBySourceChannelRespVo{Err: vo.NewErr(err), UserOutInfoBoList: res}
}

func (PostGreeter) GetOutUserInfoListByUserIds(req orgvo.GetOutUserInfoListByUserIdsReqVo) orgvo.GetOutUserInfoListByUserIdsRespVo {
	res, err := service.GetOutUserInfoListByUserIds(req.UserIds)
	return orgvo.GetOutUserInfoListByUserIdsRespVo{Err: vo.NewErr(err), UserOutInfoBoList: res}
}

func (GetGreeter) GetUserInfoListByOrg(reqVo orgvo.GetUserInfoListReqVo) orgvo.GetUserInfoListRespVo {
	res, err := service.GetUserInfoListByOrg(reqVo.OrgId)
	return orgvo.GetUserInfoListRespVo{Err: vo.NewErr(err), SimpleUserInfo: res}
}

func (PostGreeter) UpdateUserInfo(req orgvo.UpdateUserInfoReqVo) vo.CommonRespVo {
	res, err := service.UpdateUserInfo(req.OrgId, req.UserId, req.UpdateUserInfoReq)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) GetUserInfoByUserIds(input orgvo.GetUserInfoByUserIdsReqVo) orgvo.GetUserInfoByUserIdsListRespVo {
	res, err := service.GetUserInfoByUserIds(input)
	return orgvo.GetUserInfoByUserIdsListRespVo{Err: vo.NewErr(err), GetUserInfoByUserIdsRespVo: res}
}

func (PostGreeter) BatchGetUserDetailInfo(reqVo orgvo.BatchGetUserInfoReq) orgvo.BatchGetUserInfoResp {
	res, err := service.BatchGetUserDetailInfo(reqVo.UserIds)
	return orgvo.BatchGetUserInfoResp{Data: res, Err: vo.NewErr(err)}
}

func (PostGreeter) JudgeUserIsAdmin(reqVo orgvo.JudgeUserIsAdminReq) vo.BoolRespVo {
	return vo.BoolRespVo{Err: vo.NewErr(nil), IsTrue: service.JudgeUserIsAdmin(reqVo.OutOrgId, reqVo.OutUserId, reqVo.SourceChannel)}
}

func (PostGreeter) ExportAddressList(req orgvo.ExportAddressListReqVo) orgvo.ExportAddressListRespVo {
	res, err := service.ExportAddressList(req.OrgId, req.CurrentUserId, req.Data)
	return orgvo.ExportAddressListRespVo{Err: vo.NewErr(err), Data: orgvo.ExportAddressListResp{Url: res}}
}

func (PostGreeter) EmptyUser(req orgvo.EmptyUserReqVo) vo.CommonRespVo {
	res, err := service.EmptyUser(req.OrgId, req.CurrentUserId, req.Data)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

//func (PostGreeter) CreateUser(req orgvo.CreateUserReqVo) vo.CommonRespVo {
//	res, err := service.CreateUser(req.OrgId, req.CurrentUserId, req.Data)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}

//func (PostGreeter) UpdateUser(req orgvo.UpdateUserReqVo) vo.CommonRespVo {
//	res, err := service.UpdateUser(req.OrgId, req.CurrentUserId, req.Data)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}

func (PostGreeter) UserList(req orgvo.UserListReqVo) orgvo.UserListRespVo {
	res, err := service.UserList(req.OrgId, req.Data)
	return orgvo.UserListRespVo{Err: vo.NewErr(err), Data: res}
}

// 通过用户的创建时间查询用户
func (PostGreeter) GetUserListWithCreateTimeRange(req orgvo.GetUserListWithCreateTimeRangeReqVo) orgvo.GetUserListWithCreateTimeRangeRespVo {
	res, err := service.GetUserListWithCreateTimeRange(req.Data)
	return orgvo.GetUserListWithCreateTimeRangeRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) SearchUser(req orgvo.SearchUserReqVo) orgvo.SearchUserRespVo {
	res, err := service.SearchUser(req.OrgId, req.Data)
	return orgvo.SearchUserRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) InviteUser(req orgvo.InviteUserReqVo) orgvo.InviteUserRespVo {
	res, err := service.InviteUser(req.OrgId, req.CurrentUserId, req.Data)
	return orgvo.InviteUserRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) InviteUserList(req orgvo.InviteUserListReqVo) orgvo.InviteUserListRespVo {
	res, err := service.InviteUserList(req.OrgId, req.Data)
	return orgvo.InviteUserListRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) RemoveInviteUser(req orgvo.RemoveInviteUserReqVo) vo.CommonRespVo {
	err := service.RemoveInviteUser(req.OrgId, req.CurrentUserId, req.Data)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}

func (PostGreeter) UserStat(req orgvo.UserStatReqVo) orgvo.UserStatRespVo {
	res, err := service.UserStat(req.OrgId)
	return orgvo.UserStatRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) SyncUserInfoFromFeiShu(req orgvo.SyncUserInfoFromFeiShuReqVo) vo.CommonRespVo {
	err := service.PushSyncMemberDept(req.OrgId, req.CurrentUserId, req.SyncUserInfoFromFeiShuReq)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: &vo.Void{
		ID: 1,
	}}
}

func (PostGreeter) GetUserInfoByFeishuTenantKey(req orgvo.GetUserInfoByFeishuTenantKeyReq) orgvo.GetUserInfoByFeishuTenantKeyResp {
	res, err := service.GetUserInfoByFeishuTenantKey(req.TenantKey, req.OpenId)
	return orgvo.GetUserInfoByFeishuTenantKeyResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetOrgUserIdsByEmIds(req orgvo.GetOrgUserIdsByEmIdsReq) orgvo.GetOrgUserIdsByEmIdsResp {
	res, err := service.GetOrgUserIdsByEmIds(req.OrgId, req.SourceChannel, req.EmpIds)
	return orgvo.GetOrgUserIdsByEmIdsResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) FilterResignedUserIds(req orgvo.FilterResignedUserIdsReqVo) orgvo.FilterResignedUserIdsResp {
	res, err := service.FilterResignedUserIds(req.OrgId, req.UserIds)
	return orgvo.FilterResignedUserIdsResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetOrgUserIds(req orgvo.GetOrgUserIdsReq) orgvo.GetOrgUserIdsResp {
	res, err := service.GetOrgUserIds(req.OrgId)
	return orgvo.GetOrgUserIdsResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetPayRemind(req orgvo.GetPayRemindReq) orgvo.GetPayRemindResp {
	res, err := service.GetPayRemind(req.OrgId, req.UserId)
	return orgvo.GetPayRemindResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetOrgUsersInfoByEmIds(req orgvo.GetOrgUsersInfoByEmIdsReq) orgvo.GetBaseUserInfoBatchRespVo {
	res, err := service.GetOrgUsersInfoByEmIds(req.OrgId, req.SourceChannel, req.EmpIds)
	return orgvo.GetBaseUserInfoBatchRespVo{
		Err:           vo.NewErr(err),
		BaseUserInfos: res,
	}
}

func (GetGreeter) GetOrgUserAndDeptCount(req orgvo.GetOrgUserAndDeptCountReq) orgvo.GetOrgUserAndDeptCountResp {
	res, err := service.GetOrgUserAndDeptCount(req.OrgId)
	return orgvo.GetOrgUserAndDeptCountResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (GetGreeter) SetVisitUserGuideStatus(req orgvo.SetVisitUserGuideStatusReq) vo.BoolRespVo {
	isOk, err := service.SetVisitUserGuideStatus(req)
	return vo.BoolRespVo{
		Err:    vo.NewErr(err),
		IsTrue: isOk,
	}
}

func (PostGreeter) SetVersionVisible(req orgvo.SetVersionReq) vo.BoolRespVo {
	res, err := service.SetVersionVisible(req)
	return vo.BoolRespVo{
		Err:    vo.NewErr(err),
		IsTrue: res,
	}
}

func (PostGreeter) SetUserActivity(req orgvo.SetUserActivityReq) vo.BoolRespVo {
	res, err := service.SetUserActivity(req)
	return vo.BoolRespVo{
		Err:    vo.NewErr(err),
		IsTrue: res,
	}
}

func (GetGreeter) GetVersion(req orgvo.GetVersionReq) orgvo.GetVersionRespVo {
	res, err := service.GetVersion(req)
	return orgvo.GetVersionRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) SetUserViewLocation(req orgvo.SaveViewLocationReqVo) orgvo.SaveViewLocationRespVo {
	resp, err := service.SetUserViewLocation(req)
	return orgvo.SaveViewLocationRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (GetGreeter) GetUserViewLastLocation(req orgvo.GetViewLocationReq) orgvo.GetViewLocationRespVo {
	resp, err := service.GetViewLocationList(req)
	return orgvo.GetViewLocationRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (PostGreeter) DeleteAppUserViewLocation(req orgvo.UpdateUserViewLocationReq) vo.CommonRespVo {
	err := service.DeleteAppUserLocation(req.OrgId, req.UserId, req.AppId)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}

func (PostGreeter) CreateOrgUser(req orgvo.CreateOrgMemberReqVo) vo.CommonRespVo {
	err := service.CreateOrgUser(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
}

func (PostGreeter) UpdateOrgUser(req orgvo.UpdateOrgMemberReqVo) vo.CommonRespVo {
	err := service.UpdateOrgUser(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
}

func (PostGreeter) GetOrgSuperAdminInfo(req orgvo.GetOrgSuperAdminInfoReq) orgvo.GetOrgSuperAdminInfoResp {
	res, err := service.GetOrgSuperAdminInfo(req.OrgId)
	return orgvo.GetOrgSuperAdminInfoResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) UpdateUserToSysManageGroup(req orgvo.UpdateUserToSysManageGroupReq) vo.CommonRespVo {
	err := service.UpdateUserToSysManageGroup(req.OrgId, req.Input.UserIds, req.Input.UpdateType)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}
