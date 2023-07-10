package api

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (GetGreeter) GetOrgBoList() orgvo.GetOrgBoListRespVo {
	res, err := service.GetOrgBoList()
	return orgvo.GetOrgBoListRespVo{Err: vo.NewErr(err), OrganizationBoList: res}
}

func (PostGreeter) SaveOrgSummaryTableAppId(req orgvo.SaveOrgSummaryTableAppIdReqVo) orgvo.VoidRespVo {
	res, err := service.SaveOrgSummaryTableAppId(req.OrgId, req.UserId, req.Input)
	return orgvo.VoidRespVo{
		Err: vo.NewErr(err),
		Data: orgvo.VoidRespVoData{
			IsOk: res,
		},
	}
}

func (PostGreeter) CreateOrg(req orgvo.CreateOrgReqVo) orgvo.CreateOrgRespVo {
	sourceChannel := ""
	sourcePlatform := ""
	if req.Data.CreateOrgReq.SourceChannel != nil {
		sourceChannel = *req.Data.CreateOrgReq.SourceChannel
	}
	if req.Data.CreateOrgReq.SourcePlatform != nil {
		sourcePlatform = *req.Data.CreateOrgReq.SourcePlatform
	}
	orgId, err := service.CreateOrgBase(req, sourceChannel, sourcePlatform)
	if err != nil {
		log.Error(err)
		return orgvo.CreateOrgRespVo{Err: vo.NewErr(err), Data: orgvo.CreateOrgRespVoData{OrgId: orgId}}
	}

	if req.Data.ImportSampleData == 1 {
		//初始化示例数据
		err1 := service.NewbieGuideInit(orgId, req.UserId, sourceChannel)
		if err1 != nil {
			log.Error(err1)
		}
	}
	return orgvo.CreateOrgRespVo{Err: vo.NewErr(err), Data: orgvo.CreateOrgRespVoData{OrgId: orgId}}
}

func (PostGreeter) PrepareTransferOrg(req orgvo.PrepareTransferOrgReq) orgvo.PrepareTransferOrgRespVo {
	res, err := service.PrepareTransferOrg(req.Input.OriginOrgId, req.Input.ToOrgId)
	return orgvo.PrepareTransferOrgRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) TransferOrg(req orgvo.PrepareTransferOrgReq) vo.CommonRespVo {
	err := service.TransferOrg(req.Input.OriginOrgId, req.Input.ToOrgId)
	return vo.CommonRespVo{
		Err: vo.NewErr(err),
		Void: &vo.Void{
			ID: req.OrgId,
		},
	}
}

func (PostGreeter) UserOrganizationList(req orgvo.UserOrganizationListReqVo) orgvo.UserOrganizationListRespVo {
	res, err := service.UserOrganizationList(req.UserId)
	return orgvo.UserOrganizationListRespVo{Err: vo.NewErr(err), UserOrganizationListResp: res}
}

func (PostGreeter) SwitchUserOrganization(req orgvo.SwitchUserOrganizationReqVo) orgvo.SwitchUserOrganizationRespVo {
	err := service.SwitchUserOrganization(req.OrgId, req.UserId, req.Token)
	return orgvo.SwitchUserOrganizationRespVo{Err: vo.NewErr(err), OrgId: req.OrgId}
}

func (PostGreeter) UpdateOrganizationSetting(req orgvo.UpdateOrganizationSettingReqVo) orgvo.UpdateOrganizationSettingRespVo {
	res, err := service.UpdateOrganizationSetting(req)
	return orgvo.UpdateOrganizationSettingRespVo{Err: vo.NewErr(err), OrgId: res}
}

// UpdateOrgRemarkSetting org 表的 remark 存放着一些配置，该接口用户更新配置中的部分项
func (PostGreeter) UpdateOrgRemarkSetting(req orgvo.UpdateOrgRemarkSettingReqVo) vo.BoolRespVo {
	isOk, err := service.UpdateOrgRemarkSetting(req)
	return vo.BoolRespVo{
		Err:    vo.NewErr(err),
		IsTrue: isOk,
	}
}

func (PostGreeter) OrganizationInfo(req orgvo.OrganizationInfoReqVo) orgvo.OrganizationInfoRespVo {
	res, err := service.OrganizationInfo(req)
	return orgvo.OrganizationInfoRespVo{Err: vo.NewErr(err), OrganizationInfo: res}
}

func (PostGreeter) GetOrgOutInfoByOutOrgId(req orgvo.GetOutOrgInfoByOutOrgIdReqVo) orgvo.GetOutOrgInfoByOutOrgIdRespVo {
	res, err := service.GetOrgOutInfoByOutOrgId(req.OrgId, req.OutOrgId)
	return orgvo.GetOutOrgInfoByOutOrgIdRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) ScheduleOrganizationPageList(req orgvo.ScheduleOrganizationPageListReqVo) orgvo.ScheduleOrganizationPageListRespVo {
	res, err := service.ScheduleOrganizationPageList(req)
	return orgvo.ScheduleOrganizationPageListRespVo{Err: vo.NewErr(err), ScheduleOrganizationPageListResp: res}
}

//通过来源获取组织id列表
func (PostGreeter) GetOrgIdListBySourceChannel(req orgvo.GetOrgIdListBySourceChannelReqVo) orgvo.GetOrgIdListBySourceChannelRespVo {
	res, err := service.GetOrgIdListBySourceChannel(req.Input.SourceChannels, req.Page, req.Size, req.Input.IsPaid)
	return orgvo.GetOrgIdListBySourceChannelRespVo{Err: vo.NewErr(err), Data: orgvo.GetOrgIdListBySourceChannelRespData{OrgIds: res}}
}

// GetOrgIdListByPage 分页获取组织 id 列表
func (PostGreeter) GetOrgIdListByPage(req orgvo.GetOrgIdListByPageReqVo) orgvo.GetOrgIdListBySourceChannelRespVo {
	res, err := service.GetOrgIdListByPage(req.Input, req.Page, req.Size)
	return orgvo.GetOrgIdListBySourceChannelRespVo{Err: vo.NewErr(err), Data: orgvo.GetOrgIdListBySourceChannelRespData{OrgIds: res}}
}

// GetOrgBoListByPage 分页获取组织列表
func (PostGreeter) GetOrgBoListByPage(req orgvo.GetOrgIdListByPageReqVo) orgvo.GetOrgBoListByPageRespVo {
	resData, err := service.GetOrgBoListByPage(req.Input, req.Page, req.Size)
	return orgvo.GetOrgBoListByPageRespVo{Err: vo.NewErr(err), Data: resData}
}

// GetOrgConfig 获取组织配置
func (GetGreeter) GetOrgConfig(req orgvo.GetOrgConfigReq) orgvo.GetOrgConfigResp {
	res, err := service.GetOrgConfig(req.OrgId)
	respData := vo.OrgConfig{}
	copyer.Copy(res, &respData)

	return orgvo.GetOrgConfigResp{Err: vo.NewErr(err), Data: &respData}
}

// GetOrgConfigRest 获取组织配置，用于 restful 路由
func (GetGreeter) GetOrgConfigRest(req orgvo.GetOrgConfigReq) orgvo.GetOrgConfigRestResp {
	res, err := service.GetOrgConfig(req.OrgId)

	return orgvo.GetOrgConfigRestResp{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) GetOutOrgInfoByOrgIdBatch(req orgvo.GetOutOrgInfoByOrgIdBatchReqVo) orgvo.GetOutOrgInfoByOrgIdBatchRespVo {
	res, err := service.GetOutOrgInfoByOrgIdBatch(req.OrgIds)
	return orgvo.GetOutOrgInfoByOrgIdBatchRespVo{Err: vo.NewErr(err), Data: res}
}

// 申请授权
func (PostGreeter) ApplyScopes(req orgvo.ApplyScopesReqVo) orgvo.ApplyScopesRespVo {
	return orgvo.ApplyScopesRespVo{Err: vo.NewErr(nil), Data: nil}
}

func (PostGreeter) CheckSpecificScope(req orgvo.CheckSpecificScopeReqVo) orgvo.CheckSpecificScopeRespVo {
	return orgvo.CheckSpecificScopeRespVo{Err: vo.NewErr(nil), Data: nil}
}

//func (PostGreeter) UpdateOrgBasicShowSetting(req orgvo.UpdateOrgBasicShowSettingReq) vo.CommonRespVo {
//	res, err := service.UpdateOrgBasicShowSetting(req.OrgId, req.UserId, req.Params)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}

// 检查组织拥有者是否是超管，如果不是，则设置为超管。
func (GetGreeter) CheckAndSetSuperAdmin(req orgvo.CheckAndSetSuperAdminReq) vo.BoolRespVo {
	res, err := service.CheckAndSetSuperAdmin(req.OrgId)
	return vo.BoolRespVo{
		Err:    vo.NewErr(err),
		IsTrue: res,
	}
}

func (PostGreeter) StopThirdIntegration(req vo.CommonReqVo) vo.CommonRespVo {
	err := service.StopThirdIntegration(req.OrgId, req.UserId, req.SourceChannel)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: req.OrgId},
	}
}

// SendCardToAdminForUpgrade 成员点击升级极星时，需要发送卡片给管理员
func (PostGreeter) SendCardToAdminForUpgrade(req orgvo.SendCardToAdminForUpgradeReqVo) vo.BoolRespVo {
	return vo.BoolRespVo{
		Err:    vo.NewErr(nil),
		IsTrue: true,
	}
}

func (PostGreeter) GetOrgInfo(req orgvo.GetOrgInfoReq) orgvo.GetOrgInfoResp {
	res, err := service.GetOrgInfoBo(req.OrgId)
	return orgvo.GetOrgInfoResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
